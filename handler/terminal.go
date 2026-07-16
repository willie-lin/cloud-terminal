package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/ssh"

	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/resource"
	"github.com/willie-lin/cloud-terminal/pkg/sts"
)

// ─── WebSocket 升级器 ─────────────────────────────────────────

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true // 开发环境允许所有来源；生产环境应限制
	},
}

// ─── WebSocket 终端 Handler ─────────────────────────────────────

// TerminalWebSocket 处理 WebSocket 终端连接
// 流程：验证 JWT → 查询资源 → 建立 SSH 隧道 → 双向转发
func TerminalWebSocket(client *ent.Client, stsService *sts.Service) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// 1. 获取 URN 参数
		urn := c.QueryParam("urn")
		if urn == "" {
			return c.String(http.StatusBadRequest, "missing urn parameter")
		}

		// 2. 验证 JWT
		tokenStr := c.QueryParam("token")
		if tokenStr == "" {
			// 也尝试从 Authorization header 获取
			tokenStr = c.Request().Header.Get("Authorization")
			if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
				tokenStr = tokenStr[7:]
			} else {
				tokenStr = ""
			}
		}
		if tokenStr == "" {
			return c.String(http.StatusUnauthorized, "missing token")
		}

		claims, err := stsService.ValidateToken(tokenStr)
		if err != nil {
			log.Printf("Terminal: invalid token: %v", err)
			return c.String(http.StatusUnauthorized, "invalid token")
		}

		// 验证 URN 匹配
		if claims.ResourceURN != "" && claims.ResourceURN != urn {
			log.Printf("Terminal: token URN %s != requested URN %s", claims.ResourceURN, urn)
			return c.String(http.StatusForbidden, "token URN mismatch")
		}

		// 3. 查询资源
		r, err := client.Resource.Query().
			Where(resource.Urn(urn)).
			Only(c.Request().Context())
		if err != nil {
			if ent.IsNotFound(err) {
				return c.String(http.StatusNotFound, "resource not found")
			}
			log.Printf("Terminal: query resource error: %v", err)
			return c.String(http.StatusInternalServerError, "query resource error")
		}

		// 4. 升级为 WebSocket
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			log.Printf("Terminal: upgrade error: %v", err)
			return err
		}

		// 5. 建立 SSH 连接
		sshClient, err := dialSSH(r.IP, r.Port, r.AuthData)
		if err != nil {
			log.Printf("Terminal: ssh dial error: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("SSH connection failed: %v\n", err)))
			conn.Close()
			return nil
		}
		defer sshClient.Close()

		// 6. 打开 SSH session
		session, err := sshClient.NewSession()
		if err != nil {
			log.Printf("Terminal: ssh session error: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("SSH session failed: %v\n", err)))
			conn.Close()
			return nil
		}
		defer session.Close()

		// 7. 设置终端模式
		session.RequestPty("xterm-256color", 40, 80, ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		})

		// 8. 获取 SSH 管道
		stdin, err := session.StdinPipe()
		if err != nil {
			log.Printf("Terminal: stdin pipe error: %v", err)
			return nil
		}
		stdout, err := session.StdoutPipe()
		if err != nil {
			log.Printf("Terminal: stdout pipe error: %v", err)
			return nil
		}
		stderr, err := session.StderrPipe()
		if err != nil {
			log.Printf("Terminal: stderr pipe error: %v", err)
			return nil
		}

		// 9. 启动 shell
		if err := session.Shell(); err != nil {
			log.Printf("Terminal: shell error: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Shell failed: %v\n", err)))
			conn.Close()
			return nil
		}

		// 10. 双向数据转发
		errChan := make(chan error, 3)

		// SSH → WebSocket
		go func() {
			merged := io.MultiReader(stdout, stderr)
			buf := make([]byte, 4096)
			for {
				n, err := merged.Read(buf)
				if err != nil {
					errChan <- err
					return
				}
				if n > 0 {
					if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
						errChan <- err
						return
					}
				}
			}
		}()

		// WebSocket → SSH
		go func() {
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					errChan <- err
					return
				}
				if _, err := stdin.Write(msg); err != nil {
					errChan <- err
					return
				}
			}
		}()

		// 等待任一方向出错或关闭
		<-errChan

		return nil
	}
}

// ─── SSH 拨号 ─────────────────────────────────────────────────

// dialSSH 使用资源认证信息建立 SSH 连接
func dialSSH(ip string, port int, authData map[string]interface{}) (*ssh.Client, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)

	// 从 authData 提取认证信息
	username := "root"
	password := ""
	privateKey := ""

	if authData != nil {
		if u, ok := authData["username"].(string); ok && u != "" {
			username = u
		}
		if p, ok := authData["password"].(string); ok && p != "" {
			password = p
		}
		if k, ok := authData["ssh_key"].(string); ok && k != "" {
			privateKey = k
		}
	}

	// 构建认证方法
	var auths []ssh.AuthMethod

	if privateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(privateKey))
		if err == nil {
			auths = append(auths, ssh.PublicKeys(signer))
		}
	}

	if password != "" {
		auths = append(auths, ssh.Password(password))
	}

	if len(auths) == 0 {
		return nil, fmt.Errorf("no authentication method available for %s", addr)
	}

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 开发环境跳过主机密钥检查
		Timeout:         10 * time.Second,
	}

	return ssh.Dial("tcp", addr, config)
}

// ─── 消息类型 ─────────────────────────────────────────────────

// TerminalMessage WebSocket 终端消息
type TerminalMessage struct {
	Type string          `json:"type"` // "input", "resize", "ping"
	Data json.RawMessage `json:"data"`
}

// ResizeMessage 终端尺寸调整消息
type ResizeMessage struct {
	Cols int `json:"cols"`
	Rows int `json:"rows"`
}
