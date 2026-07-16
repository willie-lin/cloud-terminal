package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/resource"
	"github.com/willie-lin/cloud-terminal/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
)

// ContainerSSHAuthRequest ContainerSSH 认证请求
type ContainerSSHAuthRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
}

// ContainerSSHAuthResponse ContainerSSH 认证响应
type ContainerSSHAuthResponse struct {
	Authenticated bool   `json:"authenticated"`
	SessionID     string `json:"sessionId,omitempty"`
}

// AuthWebhook ContainerSSH 认证 Webhook
func AuthWebhook(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		req := new(ContainerSSHAuthRequest)
		if err := c.Bind(&req); err != nil {
			log.Printf("Error binding auth request: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}

		u, err := client.User.Query().
			Where(user.UsernameEQ(req.Username)).
			Only(c.Request().Context())
		if err != nil {
			if ent.IsNotFound(err) {
				return c.JSON(http.StatusOK, ContainerSSHAuthResponse{Authenticated: false})
			}
			log.Printf("Error querying user for auth: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}

		// SSH 公钥认证优先
		if req.PublicKey != "" {
			if u.SSHPublicKey == req.PublicKey {
				return c.JSON(http.StatusOK, ContainerSSHAuthResponse{
					Authenticated: true,
					SessionID:     u.ID,
				})
			}
			return c.JSON(http.StatusOK, ContainerSSHAuthResponse{Authenticated: false})
		}

		// 密码认证
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

// ContainerSSHConfigRequest ContainerSSH 配置请求
type ContainerSSHConfigRequest struct {
	Username  string `json:"username"`
	SessionID string `json:"sessionId"`
	IPAddress string `json:"ipAddress,omitempty"`
	URN       string `json:"urn,omitempty"`
}

// ContainerConfig 容器配置
type ContainerConfig struct {
	Docker *DockerConfig `json:"docker,omitempty"`
}

// DockerConfig Docker 容器配置
type DockerConfig struct {
	Image       string            `json:"image"`
	Cmd         []string          `json:"cmd,omitempty"`
	Entrypoint  []string          `json:"entrypoint,omitempty"`
	NetworkMode string            `json:"networkMode,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Resources   *ResourceLimit    `json:"resources,omitempty"`
}

// ResourceLimit 资源限制
type ResourceLimit struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// ConfigWebhook ContainerSSH 配置 Webhook
func ConfigWebhook(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		req := new(ContainerSSHConfigRequest)
		if err := c.Bind(&req); err != nil {
			log.Printf("Error binding config request: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}

		config := &ContainerConfig{
			Docker: &DockerConfig{
				Image:       "cloud-terminal/connector:latest",
				NetworkMode: "host",
				Resources: &ResourceLimit{
					CPU:    "1",
					Memory: "512M",
				},
			},
		}

		if req.URN == "" {
			return c.JSON(http.StatusOK, config)
		}

		r, err := client.Resource.Query().
			Where(resource.Urn(req.URN)).
			Only(c.Request().Context())
		if err != nil {
			log.Printf("Resource not found for URN %s: %v", req.URN, err)
			return c.JSON(http.StatusOK, config)
		}

		// 根据资源类型配置镜像
		switch r.Type {
		case "mysql":
			config.Docker.Image = "cloud-terminal/mysql-client:latest"
		case "redis":
			config.Docker.Image = "cloud-terminal/redis-client:latest"
		case "k8s-service":
			config.Docker.Image = "cloud-terminal/kubectl:latest"
		case "ssh":
			config.Docker.Image = "cloud-terminal/connector:latest"
			if r.AuthData != nil {
				if key, ok := r.AuthData["ssh_key"]; ok {
					config.Docker.Env = map[string]string{
						"SSH_KEY": key.(string),
					}
				}
			}
		}

		if config.Docker.Env == nil {
			config.Docker.Env = make(map[string]string)
		}
		config.Docker.Env["TARGET_HOST"] = r.IP
		config.Docker.Env["TARGET_PORT"] = strconv.Itoa(r.Port)

		return c.JSON(http.StatusOK, config)
	}
}
