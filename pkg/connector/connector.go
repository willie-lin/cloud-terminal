package connector

import (
	"context"

	"golang.org/x/crypto/ssh"
)

// TargetInfo 目标网络地址与主机公钥信息
type TargetInfo struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	HostKey string `json:"hostKey,omitempty"`
}

// ConnectRequest 建立连接或生成配置所需的统一请求参数
type ConnectRequest struct {
	SessionID       string                 `json:"sessionId"`
	ResourceURN     string                 `json:"resourceUrn"`
	Protocol        string                 `json:"protocol"` // "ssh", "mysql", "redis", "k8s-service" 等
	Target          TargetInfo             `json:"target"`
	AuthData        map[string]interface{} `json:"authData,omitempty"`
	ResourceDetails map[string]interface{} `json:"resourceDetails,omitempty"`
}

// ProxyConfig 代理模式参数
type ProxyConfig struct {
	Address string `json:"address,omitempty"`
}

// ResourceLimit 容器资源限制
type ResourceLimit struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// DockerConfig Docker 后端配置
type DockerConfig struct {
	Image       string            `json:"image"`
	Cmd         []string          `json:"cmd,omitempty"`
	Entrypoint  []string          `json:"entrypoint,omitempty"`
	NetworkMode string            `json:"networkMode,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Resources   *ResourceLimit    `json:"resources,omitempty"`
}

// KubernetesConfig Kubernetes 后端配置
type KubernetesConfig struct {
	Namespace      string            `json:"namespace,omitempty"`
	PodName        string            `json:"podName,omitempty"`
	Image          string            `json:"image,omitempty"`
	Cmd            []string          `json:"cmd,omitempty"`
	Env            map[string]string `json:"env,omitempty"`
	ServiceAccount string            `json:"serviceAccount,omitempty"`
}

// ContainerSSHConfig 供 ContainerSSH Webhook 使用的后端配置
type ContainerSSHConfig struct {
	Backend    string            `json:"backend,omitempty"` // "docker" | "proxy" | "kubernetes"
	Proxy      *ProxyConfig      `json:"proxy,omitempty"`
	Docker     *DockerConfig     `json:"docker,omitempty"`
	Kubernetes *KubernetesConfig `json:"kubernetes,omitempty"`
}

// Connection 连接句柄接口
type Connection interface {
	SSHClient() *ssh.Client
	Close() error
}

// Connector 核心连接器抽象接口
type Connector interface {
	Name() string
	Connect(ctx context.Context, req *ConnectRequest) (Connection, error)
	ContainerSSHConfig(ctx context.Context, req *ConnectRequest) (*ContainerSSHConfig, error)
}
