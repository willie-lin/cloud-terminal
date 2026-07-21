package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ContainerConnector 实现基于 ContainerSSH 沙箱容器（Docker/K8s）运行环境的连接器
type ContainerConnector struct {
	defaultImage string
}

// NewContainerConnector 创建 ContainerConnector 实例
func NewContainerConnector(defaultImage string) *ContainerConnector {
	if defaultImage == "" {
		defaultImage = "cloud-terminal/connector:latest"
	}
	return &ContainerConnector{defaultImage: defaultImage}
}

func (c *ContainerConnector) Name() string {
	return "container"
}

func (c *ContainerConnector) Connect(ctx context.Context, req *ConnectRequest) (Connection, error) {
	return nil, fmt.Errorf("container connector builds execution sandboxes via ContainerSSH gateway; direct SSH dial is not applicable")
}

// ContainerSSHConfig 根据协议和资源参数动态构建 ContainerSSH 沙箱容器参数
func (c *ContainerConnector) ContainerSSHConfig(ctx context.Context, req *ConnectRequest) (*ContainerSSHConfig, error) {
	if req == nil {
		return nil, fmt.Errorf("connect request cannot be nil")
	}

	cfg := &ContainerSSHConfig{
		Backend: "docker",
		Docker: &DockerConfig{
			Image:       c.defaultImage,
			NetworkMode: "host",
			Resources: &ResourceLimit{
				CPU:    "1",
				Memory: "512M",
			},
			Env: make(map[string]string),
		},
	}

	switch strings.ToLower(req.Protocol) {
	case "mysql":
		cfg.Docker.Image = "cloud-terminal/mysql-client:latest"
		if req.AuthData != nil {
			if u, ok := req.AuthData["username"].(string); ok {
				cfg.Docker.Env["DB_USER"] = u
			}
			if p, ok := req.AuthData["password"].(string); ok {
				cfg.Docker.Env["DB_PASS"] = p
			}
		}
		cfg.Docker.Cmd = []string{"mysql", "-h", req.Target.Host, "-P", strconv.Itoa(req.Target.Port)}
	case "redis":
		cfg.Docker.Image = "cloud-terminal/redis-client:latest"
		cfg.Docker.Cmd = []string{"redis-cli", "-h", req.Target.Host, "-p", strconv.Itoa(req.Target.Port)}
	case "k8s-service":
		cfg.Docker.Image = "cloud-terminal/kubectl:latest"
		if req.ResourceDetails != nil {
			if ns, ok := req.ResourceDetails["namespace"].(string); ok {
				cfg.Docker.Env["KUBE_NAMESPACE"] = ns
			}
			if sa, ok := req.ResourceDetails["service_account"].(string); ok {
				cfg.Docker.Env["KUBE_SERVICE_ACCOUNT"] = sa
			}
		}
	case "ssh":
		cfg.Docker.Image = c.defaultImage
		if req.AuthData != nil {
			if key, ok := req.AuthData["ssh_key"].(string); ok {
				cfg.Docker.Env["SSH_KEY"] = key
			}
			if u, ok := req.AuthData["username"].(string); ok {
				cfg.Docker.Env["SSH_USER"] = u
			}
		}
		cfg.Docker.Cmd = []string{"ssh", "-p", strconv.Itoa(req.Target.Port), req.Target.Host}
	default:
		cfg.Docker.Image = c.defaultImage
	}

	// 注入基础目标环境变量
	cfg.Docker.Env["TARGET_HOST"] = req.Target.Host
	cfg.Docker.Env["TARGET_PORT"] = strconv.Itoa(req.Target.Port)
	if req.ResourceURN != "" {
		cfg.Docker.Env["TARGET_URN"] = req.ResourceURN
	}
	if req.SessionID != "" {
		cfg.Docker.Env["SESSION_ID"] = req.SessionID
	}

	// 注入额外细节变量
	if req.ResourceDetails != nil {
		for k, v := range req.ResourceDetails {
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

	return cfg, nil
}
