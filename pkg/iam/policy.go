package iam

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ─── Action 匹配（含 NotAction） ──────────────────────────────

// matchAction 检查 action 是否匹配 statement 的 Actions/NotActions
func matchAction(stmt *Statement, action string) (matched bool) {
	// Actions 匹配
	if len(stmt.Actions) > 0 {
		for _, p := range stmt.Actions {
			if globMatch(action, p) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// NotActions 排除
	if len(stmt.NotActions) > 0 {
		for _, p := range stmt.NotActions {
			if globMatch(action, p) {
				return false // 被排除
			}
		}
	}

	return len(stmt.Actions) > 0 || len(stmt.NotActions) > 0
}

// ─── Resource 匹配（含 NotResource） ─────────────────────────

// matchResource 检查 resource URN 是否匹配 statement 的 Resources/NotResources
func matchResource(stmt *Statement, resourceURN string) bool {
	// Resources 匹配
	matched := false
	if len(stmt.Resources) > 0 {
		for _, p := range stmt.Resources {
			if p == "*" || urnMatch(resourceURN, p) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// NotResources 排除
	if len(stmt.NotResources) > 0 {
		for _, p := range stmt.NotResources {
			if p == "*" || urnMatch(resourceURN, p) {
				return false
			}
		}
	}

	return len(stmt.Resources) > 0 || len(stmt.NotResources) > 0
}

// ─── Policy 变量替换 ─────────────────────────────────────────

var varPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// ExpandVariables 展开 Policy 中的变量引用
// 支持：${aws:username}, ${aws:userid}, ${aws:CurrentTime}, ${resource:Urn}, ${resource:Tag/*}
func ExpandVariables(s string, req *Request) string {
	return varPattern.ReplaceAllStringFunc(s, func(match string) string {
		inner := match[2 : len(match)-1] // 去掉 ${}
		parts := strings.SplitN(inner, ":", 2)
		if len(parts) < 2 {
			return match
		}
		switch parts[0] {
		case "aws":
			return expandAWSVar(parts[1], req)
		case "resource":
			return expandResourceVar(parts[1], req)
		default:
			if req.Context != nil {
				if v, ok := req.Context[inner]; ok {
					return fmt.Sprint(v)
				}
			}
			return match
		}
	})
}

func expandAWSVar(key string, req *Request) string {
	switch key {
	case "username":
		return ctxString(req, CtxUsername)
	case "userid":
		return ctxString(req, CtxUserID)
	case "rolename":
		return ctxString(req, CtxRoleName)
	case "CurrentTime":
		return time.Now().Format(time.RFC3339)
	case "EpochTime":
		return fmt.Sprintf("%d", time.Now().Unix())
	case "SourceIp":
		return ctxString(req, CtxSourceIP)
	default:
		return "${aws:" + key + "}"
	}
}

func expandResourceVar(key string, req *Request) string {
	switch key {
	case "Urn":
		return req.ResourceURN
	default:
		if strings.HasPrefix(key, "Tag/") {
			tagKey := strings.TrimPrefix(key, "Tag/")
			if req.Context != nil {
				fullKey := CtxResourceTag + "/" + tagKey
				if v, ok := req.Context[fullKey]; ok {
					return fmt.Sprint(v)
				}
			}
		}
		return "${resource:" + key + "}"
	}
}

// ─── Policy 校验 ─────────────────────────────────────────────

// ValidatePolicy 校验策略文档的合法性
func ValidatePolicy(p Policy) []error {
	var errs []error
	if p.ID == "" {
		errs = append(errs, fmt.Errorf("policy ID is required"))
	}
	if len(p.Statements) == 0 {
		errs = append(errs, fmt.Errorf("policy %s: at least one statement required", p.ID))
	}
	for i, stmt := range p.Statements {
		if stmt.Effect != EffectAllow && stmt.Effect != EffectDeny {
			errs = append(errs, fmt.Errorf("policy %s statement[%d]: invalid effect %q", p.ID, i, stmt.Effect))
		}
		if len(stmt.Actions) == 0 && len(stmt.NotActions) == 0 {
			errs = append(errs, fmt.Errorf("policy %s statement[%d]: Actions or NotActions required", p.ID, i))
		}
	}
	return errs
}

// ValidateRequest 校验请求
func ValidateRequest(req *Request) []error {
	var errs []error
	if req.PrincipalID == "" {
		errs = append(errs, fmt.Errorf("PrincipalID is required"))
	}
	if req.Action == "" {
		errs = append(errs, fmt.Errorf("Action is required"))
	}
	if req.ResourceURN == "" {
		errs = append(errs, fmt.Errorf("ResourceURN is required"))
	}
	return errs
}

// ─── PassRole 辅助 ────────────────────────────────────────────

// PassRoleRequest 构造一个 iam:PassRole 鉴权请求
// service 是目标服务名（如 "ec2"、"lambda"），传入 "" 表示不限制
func PassRoleRequest(principalID, roleUrn, service string) *Request {
	req := NewRequest(principalID, "iam:PassRole", roleUrn)
	if service != "" {
		req.WithContext("iam:PassedToService", service)
	}
	return req
}

// BuildRoleURN 构造 Role 的 URN
func BuildRoleURN(env, region, roleName string) string {
	return "urn:ct:" + env + ":" + region + ":role:" + roleName
}
