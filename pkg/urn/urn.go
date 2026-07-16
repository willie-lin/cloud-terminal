package urn

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	Prefix    = "urn:ct"
	Separator = ":"
)

var validURN = regexp.MustCompile(`^urn:ct:[a-zA-Z0-9_-]+:[a-zA-Z0-9_-]+:[a-zA-Z0-9_-]+:[a-zA-Z0-9_.-]+$`)

// ResourceType 统一资源类型常量
const (
	TypeMySQL      = "mysql"
	TypeRedis      = "redis"
	TypeK8sService = "k8s-service"
	TypeSSH        = "ssh"
	TypeRDP        = "rdp"
	TypeVNC        = "vnc"
	TypeTelnet     = "telnet"
	TypeHTTP       = "http"
	TypeCustom     = "custom"
)

// Build 构造 URN：urn:ct:<env>:<region>:<type>:<name>
func Build(env, region, resourceType, name string) (string, error) {
	if env == "" {
		return "", fmt.Errorf("urn: env must not be empty")
	}
	if region == "" {
		region = "default"
	}
	if resourceType == "" {
		return "", fmt.Errorf("urn: resourceType must not be empty")
	}
	if name == "" {
		return "", fmt.Errorf("urn: name must not be empty")
	}
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s",
		Prefix, env, region, resourceType, name,
	), nil
}

// MustBuild 构造 URN，失败则 panic
func MustBuild(env, region, resourceType, name string) string {
	u, err := Build(env, region, resourceType, name)
	if err != nil {
		panic(err)
	}
	return u
}

// Info 解析后的 URN 结构
type Info struct {
	Env          string
	Region       string
	ResourceType string
	Name         string
	Raw          string
}

// Parse 解析 URN 字符串
func Parse(urnStr string) (*Info, error) {
	parts := strings.SplitN(urnStr, Separator, 6)
	if len(parts) != 6 || parts[0] != "urn" || parts[1] != "ct" {
		return nil, fmt.Errorf("invalid URN format: %s (expected urn:ct:<env>:<region>:<type>:<name>)", urnStr)
	}
	return &Info{
		Env:          parts[2],
		Region:       parts[3],
		ResourceType: parts[4],
		Name:         parts[5],
		Raw:          urnStr,
	}, nil
}

// MustParse 解析 URN，失败则 panic
func MustParse(urnStr string) *Info {
	info, err := Parse(urnStr)
	if err != nil {
		panic(err)
	}
	return info
}

// Validate 校验 URN 格式是否正确
func Validate(urnStr string) bool {
	return validURN.MatchString(urnStr)
}

// Match 通配符匹配：* 匹配任意段
// "urn:ct:prod:*:mysql:*" → 匹配所有 prod 环境的 mysql 资源
// "urn:ct:*:*:*:*" → 匹配所有资源
func Match(pattern, urnStr string) bool {
	pParts := strings.SplitN(pattern, Separator, 6)
	uParts := strings.SplitN(urnStr, Separator, 6)
	if len(pParts) != 6 || len(uParts) != 6 {
		return false
	}
	for i := 0; i < 6; i++ {
		if pParts[i] == "*" {
			continue
		}
		if pParts[i] != uParts[i] {
			return false
		}
	}
	return true
}

// ResourceTypeFromAction 根据操作名推断资源类型
// 如 "resource:connect" → "resource", "mysql:query" → "mysql"
func ResourceTypeFromAction(action string) string {
	parts := strings.SplitN(action, ":", 2)
	if len(parts) > 1 {
		return parts[0]
	}
	return action
}

// BuildResourcePattern 构造用于策略匹配的资源通配模式
// 如 BuildResourcePattern("prod", "*", "mysql", "*") → "urn:ct:prod:*:mysql:*"
func BuildResourcePattern(env, region, resourceType, name string) string {
	if env == "" {
		env = "*"
	}
	if region == "" {
		region = "*"
	}
	if resourceType == "" {
		resourceType = "*"
	}
	if name == "" {
		name = "*"
	}
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s", Prefix, env, region, resourceType, name)
}
