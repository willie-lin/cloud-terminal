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
	entresource "github.com/willie-lin/cloud-terminal/ent/resource"
	entsession "github.com/willie-lin/cloud-terminal/ent/session"
	"github.com/willie-lin/cloud-terminal/pkg/connector"
	"github.com/willie-lin/cloud-terminal/pkg/crypto"
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
// 流程：验证 JWT → 检查资源 → 记录 Session+AuditLog → 建立 SSH 隧道 → 双向转发 → 关闭时更新记录
func TerminalWebSocket(client *ent.Client, stsService *sts.Service) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx := c.Request().Context()

		// 1. 获取 URN 参数
		urn := c.QueryParam("urn")
		if urn == "" {
			return c.String(http.StatusBadRequest, "missing urn parameter")
		}

		// 2. 验证 JWT
		tokenStr := c.QueryParam("token")
		if tokenStr == "" {
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

		// 3. 查询资源并检查状态
		r, err := client.Resource.Query().
			Where(entresource.Urn(urn)).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return c.String(http.StatusNotFound, "resource not found")
			}
			log.Printf("Terminal: query resource error: %v", err)
			return c.String(http.StatusInternalServerError, "query resource error")
		}
		if r.Status == entresource.StatusInactive {
			log.Printf("Terminal: resource %s is inactive", urn)
			return c.String(http.StatusForbidden, "resource is inactive")
		}

		// 4. 记录 Session 开始
		now := time.Now()
		remoteAddr := c.Request().RemoteAddr

		sessionRec, err := client.Session.Create().
			SetSessionID(claims.SessionID).
			SetPrincipalUrn("urn:ct::user:" + claims.UserID).
			SetResourceUrn(urn).
			SetMode(entsession.ModeProxy).
			SetStatus(entsession.StatusActive).
			SetStartedAt(now).
			SetRemoteAddress(remoteAddr).
			Save(ctx)
		if err != nil {
			log.Printf("Terminal: create session record error: %v", err)
			// non-fatal: continue even if session logging fails
		}

		// 5. 记录审计日志开始
		auditRec, err := client.AuditLog.Create().
			SetSessionID(claims.SessionID).
			SetUsername(claims.Subject).
			SetAction("resource:connect").
			SetResult("success").
			SetStartedAt(now).
			SetResourceUrnSnapshot(urn).
			SetDetail(map[string]interface{}{
				"resource_ip":   r.IP,
				"resource_port": r.Port,
				"resource_type": r.Type,
				"remote_addr":   remoteAddr,
			}).
			Save(ctx)
		if err != nil {
			log.Printf("Terminal: create audit log error: %v", err)
		}

		// 6. 升级为 WebSocket
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			log.Printf("Terminal: upgrade error: %v", err)
			// 回滚 session 记录
			if sessionRec != nil {
				client.Session.UpdateOneID(sessionRec.ID).SetStatus(entsession.StatusError).Save(ctx)
			}
			return err
		}

		// 7. 解密敏感认证信息并建立连接
		decAuth, err := crypto.DecryptAuthData(r.AuthData)
		if err != nil {
			log.Printf("Terminal: decrypt auth_data error: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Credential decryption failed: %v\n", err)))
			if sessionRec != nil {
				client.Session.UpdateOneID(sessionRec.ID).SetStatus(entsession.StatusError).SetEndedAt(time.Now()).Save(ctx)
			}
			if auditRec != nil {
				client.AuditLog.UpdateOneID(auditRec.ID).SetResult("failure").SetEndedAt(time.Now()).Save(ctx)
			}
			conn.Close()
			return nil
		}

		directConn := connector.NewDirectConnector()
		connHandle, err := directConn.Connect(ctx, &connector.ConnectRequest{
			SessionID:       claims.SessionID,
			ResourceURN:     urn,
			Protocol:        string(r.Type),
			Target: connector.TargetInfo{
				Host:    r.IP,
				Port:    r.Port,
				HostKey: r.HostKey,
			},
			AuthData:        decAuth,
			ResourceDetails: r.Details,
		})
		if err != nil {
			log.Printf("Terminal: connector dial error: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Connection failed: %v\n", err)))
			if sessionRec != nil {
				client.Session.UpdateOneID(sessionRec.ID).SetStatus(entsession.StatusError).SetEndedAt(time.Now()).Save(ctx)
			}
			if auditRec != nil {
				client.AuditLog.UpdateOneID(auditRec.ID).SetResult("failure").SetEndedAt(time.Now()).Save(ctx)
			}
			conn.Close()
			return nil
		}
		defer connHandle.Close()

		sshClient := connHandle.SSHClient()
		if sshClient == nil {
			conn.WriteMessage(websocket.TextMessage, []byte("SSH client is nil\n"))
			conn.Close()
			return nil
		}

		// 8. 打开 SSH session
		session, err := sshClient.NewSession()
		if err != nil {
			log.Printf("Terminal: ssh session error: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("SSH session failed: %v\n", err)))
			if sessionRec != nil {
				client.Session.UpdateOneID(sessionRec.ID).SetStatus(entsession.StatusError).SetEndedAt(time.Now()).Save(ctx)
			}
			if auditRec != nil {
				client.AuditLog.UpdateOneID(auditRec.ID).SetResult("failure").SetEndedAt(time.Now()).Save(ctx)
			}
			conn.Close()
			return nil
		}
		defer session.Close()

		// 9. 设置终端模式
		session.RequestPty("xterm-256color", 40, 80, ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		})

		// 10. 获取 SSH 管道
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

		// 11. 启动 shell
		if err := session.Shell(); err != nil {
			log.Printf("Terminal: shell error: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Shell failed: %v\n", err)))
			if sessionRec != nil {
				client.Session.UpdateOneID(sessionRec.ID).SetStatus(entsession.StatusError).SetEndedAt(time.Now()).Save(ctx)
			}
			if auditRec != nil {
				client.AuditLog.UpdateOneID(auditRec.ID).SetResult("failure").SetEndedAt(time.Now()).Save(ctx)
			}
			conn.Close()
			return nil
		}

		// 12. 双向数据转发
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

		// 13. 更新 Session 和审计记录为已关闭
		endTime := time.Now()
		if sessionRec != nil {
			client.Session.UpdateOneID(sessionRec.ID).
				SetStatus(entsession.StatusClosed).
				SetEndedAt(endTime).
				Save(ctx)
		}
		if auditRec != nil {
			client.AuditLog.UpdateOneID(auditRec.ID).
				SetEndedAt(endTime).
				Save(ctx)
		}

		return nil
	}
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
