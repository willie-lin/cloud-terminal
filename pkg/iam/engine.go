package iam

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Engine 策略鉴权引擎
type Engine struct{}

// New 创建鉴权引擎
func New() *Engine { return &Engine{} }

// ─── 单次评估 ──────────────────────────────────────────────────

// Evaluate 执行一次完整鉴权
// 使用所有策略层（SCP → Boundary → Identity → Resource → Session）
func (e *Engine) Evaluate(req *Request, chain *EvalChain) *Result {
	// 预处理：展开变量
	chain = expandChainVars(chain, req)

	// 0️⃣  请求校验
	if errs := ValidateRequest(req); len(errs) > 0 {
		return &Result{Decision: DecisionDeny, MatchedBy: "validation", Reason: fmt.Sprintf("invalid request: %v", errs)}
	}

	// 1️⃣  SCP 层：组织级控制
	if r := evalLayer(req, chain.OrganizationSCP, LayerSCP); r != nil {
		return r
	}

	// 2️⃣  Permission Boundary
	if chain.PermissionBoundary != nil {
		r := evalBoundary(req, chain.PermissionBoundary)
		if r != nil {
			return r
		}
	}

	// 3️⃣  Identity Policies
	identityResult := evalLayer(req, chain.IdentityPolicies, LayerIdentity)

	// 4️⃣  Resource Policy
	resourceResult := evalSinglePolicy(req, chain.ResourcePolicy, LayerResource)

	// 5️⃣  Session Policies
	sessionResult := evalLayer(req, chain.SessionPolicies, LayerSession)

	// ── 最终决策 ──
	return finalDecision(identityResult, resourceResult, sessionResult)
}

// ─── 各层级评估 ────────────────────────────────────────────────

// evalLayer 评估一个策略层（多个 Policy）
// 如果找到显式 Deny，立即返回 Deny
// 如果找到显式 Allow，保留但不立即返回（后续层可能 Deny）
func evalLayer(req *Request, policies []Policy, layer PolicyLayer) *Result {
	sorted := sortByPriority(policies)

	// 先扫所有 Deny
	for _, p := range sorted {
		for i, stmt := range p.Statements {
			if stmt.Effect != EffectDeny {
				continue
			}
			if !matchAction(&stmt, req.Action) {
				continue
			}
			if !matchResource(&stmt, req.ResourceURN) {
				continue
			}
			if !MatchConditions(stmt.Condition, req) {
				continue
			}
			return &Result{
				Decision: DecisionDeny, MatchedBy: layer.String(),
				PolicyID: p.ID, Statement: i, Reason: fmt.Sprintf("%s explicit deny", layer),
			}
		}
	}

	// 再扫所有 Allow（只找第一个匹配）
	for _, p := range sorted {
		for i, stmt := range p.Statements {
			if stmt.Effect != EffectAllow {
				continue
			}
			if !matchAction(&stmt, req.Action) {
				continue
			}
			if !matchResource(&stmt, req.ResourceURN) {
				continue
			}
			if !MatchConditions(stmt.Condition, req) {
				continue
			}
			return &Result{
				Decision: DecisionAllow, MatchedBy: layer.String(),
				PolicyID: p.ID, Statement: i, Reason: fmt.Sprintf("%s explicit allow", layer),
			}
		}
	}

	return nil // 本层未决定
}

// evalSinglePolicy 评估单个策略
func evalSinglePolicy(req *Request, p *Policy, layer PolicyLayer) *Result {
	if p == nil {
		return nil
	}
	return evalLayer(req, []Policy{*p}, layer)
}

// evalBoundary 评估权限边界
// 边界内的 Allow 放行，边界外的 Deny
func evalBoundary(req *Request, boundary *Policy) *Result {
	if boundary == nil {
		return nil
	}
	result := evalSinglePolicy(req, boundary, LayerPermissionBoundary)
	if result == nil {
		// 边界里没有显式语句匹配 — 超出边界
		return &Result{
			Decision: DecisionDeny, MatchedBy: LayerPermissionBoundary.String(),
			Reason: "not allowed by permission boundary",
		}
	}
	if result.Decision == DecisionDeny {
		return result
	}
	return nil // 边界内允许，继续
}

// ─── 最终决策 ──────────────────────────────────────────────────

// finalDecision 综合各层结果做出最终决策
//
//	逻辑（AWS IAM 等价）：
//	  - Session Deny → Deny（Session 只能缩小权限）
//	  - Session Allow + Identity/Resource Allow → Allow
//	  - Identity Allow OR Resource Allow → Allow
//	  - 其他 → Deny
func finalDecision(identity, resource, session *Result) *Result {
	// 任何层显式 Deny → 立即拒绝
	for _, r := range []*Result{identity, resource, session} {
		if r != nil && r.Decision == DecisionDeny {
			return r
		}
	}

	// 身份层或资源层任意 Allow → 允许
	identityAllow := identity != nil && identity.Decision == DecisionAllow
	resourceAllow := resource != nil && resource.Decision == DecisionAllow

	if identityAllow || resourceAllow {
		// Session 存在时必须也 Allow（Session 只能缩小权限）
		if session != nil && session.Decision != DecisionAllow {
			return &Result{
				Decision: DecisionDeny, MatchedBy: LayerSession.String(),
				Reason: "session policy does not allow this action",
			}
		}
		// 返回优先的身份层结果，身份层不存在时返回资源层结果
		if identityAllow {
			return identity
		}
		return resource
	}

	return &Result{Decision: DecisionDeny, Reason: "default deny (no matching policy)"}
}

// ─── 工具 ──────────────────────────────────────────────────────

func sortByPriority(policies []Policy) []Policy {
	out := make([]Policy, len(policies))
	copy(out, policies)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Priority < out[j].Priority
	})
	return out
}

func expandChainVars(chain *EvalChain, req *Request) *EvalChain {
	if chain == nil {
		return &EvalChain{}
	}
	c := *chain
	for i := range c.OrganizationSCP {
		c.OrganizationSCP[i] = expandPolicyVars(c.OrganizationSCP[i], req)
	}
	if c.PermissionBoundary != nil {
		p := expandPolicyVars(*c.PermissionBoundary, req)
		c.PermissionBoundary = &p
	}
	for i := range c.IdentityPolicies {
		c.IdentityPolicies[i] = expandPolicyVars(c.IdentityPolicies[i], req)
	}
	if c.ResourcePolicy != nil {
		p := expandPolicyVars(*c.ResourcePolicy, req)
		c.ResourcePolicy = &p
	}
	for i := range c.SessionPolicies {
		c.SessionPolicies[i] = expandPolicyVars(c.SessionPolicies[i], req)
	}
	return &c
}

func expandPolicyVars(p Policy, req *Request) Policy {
	// 深拷贝 Statements
	stmts := make([]Statement, len(p.Statements))
	for i, s := range p.Statements {
		actions := make([]string, len(s.Actions))
		copy(actions, s.Actions)
		notActions := make([]string, len(s.NotActions))
		copy(notActions, s.NotActions)
		resources := make([]string, len(s.Resources))
		copy(resources, s.Resources)
		notResources := make([]string, len(s.NotResources))
		copy(notResources, s.NotResources)

		for j, a := range actions {
			actions[j] = ExpandVariables(a, req)
		}
		for j, r := range resources {
			resources[j] = ExpandVariables(r, req)
		}

		stmts[i] = Statement{
			Effect:       s.Effect,
			Actions:      actions,
			NotActions:   notActions,
			Resources:    resources,
			NotResources: notResources,
			Condition:    s.Condition,
		}
	}
	return Policy{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Statements:  stmts,
		Priority:    p.Priority,
		Version:     p.Version,
	}
}

// ─── Simulate API ──────────────────────────────────────────────

// Simulate 批量模拟评估多个 Action/Resource 组合
// 等价于 AWS IAM SimulatePrincipalPolicy
func (e *Engine) Simulate(req *Request, chain *EvalChain, actions, resources []string) []SimulateResult {
	var results []SimulateResult
	for _, action := range actions {
		for _, resource := range resources {
			r := &Request{
				PrincipalID: req.PrincipalID,
				Action:      action,
				ResourceURN: resource,
				Context:     req.Context,
			}
			result := e.Evaluate(r, chain)
			results = append(results, SimulateResult{
				Action:      action,
				ResourceURN: resource,
				Decision:    result.Decision,
				MatchedBy:   result.MatchedBy,
				PolicyID:    result.PolicyID,
				Reason:      result.Reason,
			})
		}
	}
	return results
}

// ─── TrustPolicy 评估 ─────────────────────────────────────────

// CanAssumeRole 检查一个 Principal 是否被信任可以 Assume 指定 Role
// trustPolicyJSON 是 Role 实体中 trust_policy 字段的 JSON 值
func (e *Engine) CanAssumeRole(principalID string, trustPolicyJSON interface{}, req *Request) *Result {
	if trustPolicyJSON == nil {
		return &Result{Decision: DecisionDeny, MatchedBy: "TrustPolicy", Reason: "no trust policy"}
	}

	// 解析 trust_policy JSON
	data, err := json.Marshal(trustPolicyJSON)
	if err != nil {
		return &Result{Decision: DecisionDeny, MatchedBy: "TrustPolicy", Reason: "invalid trust policy json"}
	}

	var stmts []Statement
	if err := json.Unmarshal(data, &stmts); err != nil {
		// 尝试按 Policy 结构解析
		var tp struct {
			Version    string      `json:"version"`
			Statements []Statement `json:"statements"`
		}
		if err := json.Unmarshal(data, &tp); err != nil {
			return &Result{Decision: DecisionDeny, MatchedBy: "TrustPolicy", Reason: "cannot parse trust policy"}
		}
		stmts = tp.Statements
	}

	// TrustPolicy 只检查 Effect=Allow + Principal 匹配
	// 不检查 Action/Resource（信任策略不是权限策略）
	for i, stmt := range stmts {
		if stmt.Effect != EffectAllow {
			continue
		}
		// 检查 Principal 是否包含当前用户
		if stmt.Principal != nil {
			if !matchPrincipal(stmt.Principal, principalID) {
				continue
			}
		}
		// 信任策略中的 Condition 限制
		if !MatchConditions(stmt.Condition, req) {
			continue
		}
		return &Result{
			Decision: DecisionAllow, MatchedBy: "TrustPolicy",
			Statement: i, Reason: "trust policy allows assume",
		}
	}

	return &Result{Decision: DecisionDeny, MatchedBy: "TrustPolicy", Reason: "not trusted by any statement"}
}

// matchPrincipal 检查 Principal 是否匹配
func matchPrincipal(p *Principal, principalID string) bool {
	if p == nil {
		return false
	}
	for _, id := range p.AWS {
		if id == "*" || id == principalID {
			return true
		}
	}
	for _, id := range p.Canonical {
		if id == principalID {
			return true
		}
	}
	return false
}

// CheckPassRole 检查用户是否有权限将指定 Role 传递给服务
// AWS IAM 的 iam:PassRole 权限检查
// 用户必须拥有对目标 Role 的 iam:PassRole 操作权限
// 可选支持 iam:PassedToService 条件来限制目标服务
func (e *Engine) CheckPassRole(req *Request, chain *EvalChain, roleUrn string) *Result {
	if chain == nil {
		return &Result{Decision: DecisionDeny, MatchedBy: "PassRole", Reason: "no eval chain provided"}
	}

	// 构造 PassRole 请求
	passReq := &Request{
		PrincipalID: req.PrincipalID,
		Action:      "iam:PassRole",
		ResourceURN: roleUrn,
		Context:     req.Context,
	}

	// 走完整评估链检查
	result := e.Evaluate(passReq, chain)
	if result.Decision == DecisionAllow {
		return &Result{
			Decision: DecisionAllow, MatchedBy: "PassRole",
			PolicyID: result.PolicyID, Reason: "iam:PassRole allowed",
		}
	}

	return &Result{
		Decision: DecisionDeny, MatchedBy: "PassRole",
		Reason: "no iam:PassRole permission for this role",
	}
}
