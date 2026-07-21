package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"

	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/resource"
	"github.com/willie-lin/cloud-terminal/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/connector"
	"github.com/willie-lin/cloud-terminal/pkg/crypto"
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
			result, err := h.evaluator.Evaluate(ctx, &iam.Request{
				PrincipalID: req.SessionID,
				Action:      "resource:connect",
				ResourceURN: req.URN,
			})
			if err != nil {
				pkglogger.Warn("ConfigWebhookV2: IAM evaluate error", zap.Error(err))
				return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
			}
			if result.Decision != iam.DecisionAllow {
				pkglogger.Warn("ConfigWebhookV2: IAM denied",
					zap.String("user", req.Username),
					zap.String("urn", req.URN),
					zap.String("reason", result.Reason),
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
		config := h.buildContainerConfig(c.Request().Context(), resourceURN, req.SessionID)
		return c.JSON(http.StatusOK, config)
	}
}

// buildContainerConfig 根据资源和认证信息调用 ContainerConnector 构建容器运行沙箱参数
func (h *ContainerSSHHandler) buildContainerConfig(ctx context.Context, urn string, sessionID string) *connector.ContainerSSHConfig {
	connMgr := connector.NewContainerConnector("cloud-terminal/connector:latest")

	if urn == "" {
		cfg, _ := connMgr.ContainerSSHConfig(ctx, &connector.ConnectRequest{SessionID: sessionID})
		return cfg
	}

	r, err := h.client.Resource.Query().
		Where(resource.Urn(urn)).
		Only(ctx)
	if err != nil {
		log.Printf("ConfigWebhookV2: resource not found for URN %s: %v", urn, err)
		cfg, _ := connMgr.ContainerSSHConfig(ctx, &connector.ConnectRequest{SessionID: sessionID, ResourceURN: urn})
		return cfg
	}

	decAuth, err := crypto.DecryptAuthData(r.AuthData)
	if err != nil {
		log.Printf("ConfigWebhookV2: decrypt auth_data error: %v", err)
	}

	req := &connector.ConnectRequest{
		SessionID:   sessionID,
		ResourceURN: urn,
		Protocol:    string(r.Type),
		Target: connector.TargetInfo{
			Host:    r.IP,
			Port:    r.Port,
			HostKey: r.HostKey,
		},
		AuthData:        decAuth,
		ResourceDetails: r.Details,
	}

	cfg, err := connMgr.ContainerSSHConfig(ctx, req)
	if err != nil {
		log.Printf("ConfigWebhookV2: build container config error: %v", err)
		// 遇到错误时回退默认配置
		cfg, _ = connMgr.ContainerSSHConfig(ctx, &connector.ConnectRequest{SessionID: sessionID, ResourceURN: urn})
	}
	return cfg
}
