package connector

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/ssh"
)

// DirectConnector 实现直接连通目标主机（如 SSH 转发或直接代理）的连接器
type DirectConnector struct{}

// NewDirectConnector 创建 DirectConnector 实例
func NewDirectConnector() *DirectConnector {
	return &DirectConnector{}
}

func (c *DirectConnector) Name() string {
	return "direct"
}

type directConnection struct {
	client *ssh.Client
}

func (d *directConnection) SSHClient() *ssh.Client {
	return d.client
}

func (d *directConnection) Close() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}

// Connect 使用资源认证信息直接拨号建立 SSH 连接
func (c *DirectConnector) Connect(ctx context.Context, req *ConnectRequest) (Connection, error) {
	if req == nil {
		return nil, fmt.Errorf("connect request cannot be nil")
	}

	addr := fmt.Sprintf("%s:%d", req.Target.Host, req.Target.Port)
	username := "root"
	password := ""
	privateKey := ""

	if req.AuthData != nil {
		if u, ok := req.AuthData["username"].(string); ok && u != "" {
			username = u
		}
		if p, ok := req.AuthData["password"].(string); ok && p != "" {
			password = p
		}
		if k, ok := req.AuthData["ssh_key"].(string); ok && k != "" {
			privateKey = k
		}
	}

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

	hostKeyCallback := ssh.InsecureIgnoreHostKey()
	if req.Target.HostKey != "" {
		parsed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(req.Target.HostKey))
		if err != nil {
			log.Printf("Warning: failed to parse hostKey for %s: %v, falling back to insecure", addr, err)
		} else {
			hostKeyCallback = ssh.FixedHostKey(parsed)
		}
	}

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            auths,
		HostKeyCallback: hostKeyCallback,
		Timeout:         10 * time.Second,
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return &directConnection{client: client}, nil
}

// ContainerSSHConfig 构建指向目标直接 SSH 端口的 proxy backend 配置
func (c *DirectConnector) ContainerSSHConfig(ctx context.Context, req *ConnectRequest) (*ContainerSSHConfig, error) {
	if req == nil {
		return nil, fmt.Errorf("connect request cannot be nil")
	}

	addr := fmt.Sprintf("%s:%d", req.Target.Host, req.Target.Port)
	return &ContainerSSHConfig{
		Backend: "proxy",
		Proxy: &ProxyConfig{
			Address: addr,
		},
	}, nil
}
