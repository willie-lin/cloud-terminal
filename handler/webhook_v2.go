package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"

	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/resource"
	"github.com/willie-lin/cloud-terminal/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/iam"
	pkglogger "github.com/willie-lin/cloud-terminal/pkg/logger"
	"github.com/willie-lin/cloud-terminal/pkg/sts"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
)

// ─── ContainerSSH v2 Webhook ─────────────────────────────────

// ContainerSSH 认证/配置处理器，集成 STS + IAM 引擎
type ContainerSSHHandler struct {
	client     *ent.Client
	stsService *sts.Service
	evaluator  *iam.Evaluator
}

func NewContainerSSHHandler(client *ent.Client, stsService *sts.Service, evaluator *iam.Evaluator) *ContainerSSHHandler {
	return &ContainerSSHHandler{
		client:     client,
		stsService: stsService,
		evaluator:  evaluator,
	}
}

// ─── 认证 Webhook ─────────────────────────────────────────────

// AuthWebhookV2 增强版认证 Webhook
// 支持三种认证方式：STS Token / 密码 / SSH 公钥
func (h *ContainerSSHHandler) AuthWebhookV2() echo.HandlerFunc {
	return func(c *echo.Context) error {
		req := new(ContainerSSHAuthRequest)
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}

		log.Printf("AuthWebhookV2: username=%s, hasPassword=%v, hasPublicKey=%v",
			req.Username, req.Password != "", req.PublicKey != "")

		// 尝试用 STS Token 认证（username 为 token）
		if strings.HasPrefix(req.Username, "ct-sts-") {
			claims, err := h.stsService.ValidateToken(req.Password)
			if err != nil {
				return c.JSON(http.StatusOK, ContainerSSHAuthResponse{Authenticated: false})
			}
			// 验证 token 对应的用户
			if claims.UserID == req.Username[7:] {
				return c.JSON(http.StatusOK, ContainerSSHAuthResponse{
					Authenticated: true,
					SessionID:     claims.SessionID,
				})
			}
			return c.JSON(http.StatusOK, ContainerSSHAuthResponse{Authenticated: false})
		}

		// 传统认证：密码/SSH 公钥
		u, err := h.client.User.Query().
			Where(user.UsernameEQ(req.Username)).
			Only(c.Request().Context())
		if err != nil {
			if ent.IsNotFound(err) {
				return c.JSON(http.StatusOK, ContainerSSHAuthResponse{Authenticated: false})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}

		// SSH 公钥
		if req.PublicKey != "" && u.SSHPublicKey == req.PublicKey {
			return c.JSON(http.StatusOK, ContainerSSHAuthResponse{
				Authenticated: true,
				SessionID:     u.ID,
			})
		}

		// 密码
		if req.Password == "" {
			return c.JSON(http.StatusOK, ContainerSSHAuthResponse{Authenticated: false})
		}
		if err := utils.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
			return c.JSON(http.StatusOK, ContainerSSHAuthResponse{Authenticated: false})
		}

		return c.JSON(http.StatusOK, ContainerSSHAuthResponse{
			Authenticated: true,
			SessionID:     u.ID,
		})
	}
}

// ─── 配置 Webhook ─────────────────────────────────────────────

// ConfigWebhookV2 增强版配置 Webhook
// 根据 IAM 评估 + 资源信息动态生成容器配置
func (h *ContainerSSHHandler) ConfigWebhookV2() echo.HandlerFunc {
	return func(c *echo.Context) error {
		req := new(ContainerSSHConfigRequest)
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}

		ctx := c.Request().Context()

		// 1. IAM 鉴权
		if req.URN != "" && req.SessionID != "" {
			chain, err := h.evaluator.Evaluate(ctx, &iam.Request{
				PrincipalID: req.SessionID,
				Action:      "resource:connect",
				ResourceURN: req.URN,
			})
			if err != nil {
				pkglogger.Warn("ConfigWebhookV2: IAM evaluate error", zap.Error(err))
				return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
			}
			if chain.Decision != iam.DecisionAllow {
				pkglogger.Warn("ConfigWebhookV2: IAM denied",
					zap.String("user", req.Username),
					zap.String("urn", req.URN),
					zap.String("reason", chain.Reason),
				)
				return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
			}
		}

		// 2. 查询资源（URL 或用户名中提取）
		resourceURN := req.URN
		if resourceURN == "" && req.Username != "" {
			// 尝试从 username 提取 URN（格式：urn:ct:...）
			if strings.HasPrefix(req.Username, "urn:ct:") {
				resourceURN = req.Username
			}
		}

		// 3. 生成容器配置
		config := h.buildContainerConfig(c.Request().Context(), resourceURN)
		return c.JSON(http.StatusOK, config)
	}
}

// buildContainerConfig 根据资源和认证信息构建容器配置
func (h *ContainerSSHHandler) buildContainerConfig(ctx context.Context, urn string) *ContainerConfig {
	cfg := &ContainerConfig{
		Docker: &DockerConfig{
			Image:       "cloud-terminal/connector:latest",
			NetworkMode: "host",
			Resources: &ResourceLimit{
				CPU:    "1",
				Memory: "512M",
			},
			Env: make(map[string]string),
		},
	}

	if urn == "" {
		return cfg
	}

	r, err := h.client.Resource.Query().
		Where(resource.Urn(urn)).
		Only(ctx)
	if err != nil {
		log.Printf("ConfigWebhookV2: resource not found for URN %s: %v", urn, err)
		return cfg
	}

	// 根据资源类型选择镜像
	switch r.Type {
	case "mysql":
		cfg.Docker.Image = "cloud-terminal/mysql-client:latest"
		if r.AuthData != nil {
			if u, ok := r.AuthData["username"].(string); ok {
				cfg.Docker.Env["DB_USER"] = u
			}
			if p, ok := r.AuthData["password"].(string); ok {
				cfg.Docker.Env["DB_PASS"] = p
			}
		}
		cfg.Docker.Cmd = []string{"mysql", "-h", r.IP, "-P", strconv.Itoa(r.Port)}
	case "redis":
		cfg.Docker.Image = "cloud-terminal/redis-client:latest"
		cfg.Docker.Cmd = []string{"redis-cli", "-h", r.IP, "-p", strconv.Itoa(r.Port)}
	case "k8s-service":
		cfg.Docker.Image = "cloud-terminal/kubectl:latest"
		if r.Details != nil {
			if ns, ok := r.Details["namespace"].(string); ok {
				cfg.Docker.Env["KUBE_NAMESPACE"] = ns
			}
			if sa, ok := r.Details["service_account"].(string); ok {
				cfg.Docker.Env["KUBE_SERVICE_ACCOUNT"] = sa
			}
		}
	case "ssh":
		cfg.Docker.Image = "cloud-terminal/connector:latest"
		if r.AuthData != nil {
			if key, ok := r.AuthData["ssh_key"].(string); ok {
				cfg.Docker.Env["SSH_KEY"] = key
			}
			if u, ok := r.AuthData["username"].(string); ok {
				cfg.Docker.Env["SSH_USER"] = u
			}
		}
		cfg.Docker.Cmd = []string{"ssh", fmt.Sprintf("%s@%s", r.IP, strconv.Itoa(r.Port))}
	default:
		cfg.Docker.Image = "cloud-terminal/connector:latest"
	}

	// 注入目标地址
	cfg.Docker.Env["TARGET_HOST"] = r.IP
	cfg.Docker.Env["TARGET_PORT"] = strconv.Itoa(r.Port)
	cfg.Docker.Env["TARGET_URN"] = r.Urn

	// 从 details 注入额外环境变量
	if r.Details != nil {
		for k, v := range r.Details {
			key := "RESOURCE_" + strings.ToUpper(strings.ReplaceAll(k, ".", "_"))
			switch val := v.(type) {
			case string:
				cfg.Docker.Env[key] = val
			case float64:
				cfg.Docker.Env[key] = strconv.FormatFloat(val, 'f', -1, 64)
			case bool:
				cfg.Docker.Env[key] = strconv.FormatBool(val)
			default:
				if data, err := json.Marshal(v); err == nil {
					cfg.Docker.Env[key] = string(data)
				}
			}
		}
	}

	return cfg
}
