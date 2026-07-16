package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/accesspolicy"
	"github.com/willie-lin/cloud-terminal/ent/group"
	"github.com/willie-lin/cloud-terminal/ent/resource"
	"github.com/willie-lin/cloud-terminal/ent/role"
	"github.com/willie-lin/cloud-terminal/ent/tenant"
	"github.com/willie-lin/cloud-terminal/ent/user"
)

// ─── EntPolicyProvider ────────────────────────────────────────

// EntPolicyProvider 基于 ent 数据库的策略提供者
// 完整实现五层评估链的数据加载：
//
//	SCP（租户级）→ PermissionBoundary → IdentityPolicies（User→Group→Role）→ ResourcePolicy
type EntPolicyProvider struct {
	client *ent.Client
}

func NewEntPolicyProvider(client *ent.Client) *EntPolicyProvider {
	return &EntPolicyProvider{client: client}
}

// GetEvalChain 加载指定用户的完整策略链（不含 ResourcePolicy）
func (p *EntPolicyProvider) GetEvalChain(ctx context.Context, principalID string) (*EvalChain, error) {
	chain := &EvalChain{}

	// 0️⃣ 加载用户实体（含 Group 边缘）
	u, err := p.client.User.Query().
		Where(user.IDEQ(principalID)).
		WithGroup().
		WithRoles(func(rq *ent.RoleQuery) {
			rq.WithAccessPolicies()
			rq.WithPermissionsBoundary()
		}).
		WithAccessPolicies().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return chain, nil
		}
		return nil, fmt.Errorf("query user %s: %w", principalID, err)
	}

	// 1️⃣ SCP：从用户所属租户加载组织级策略
	scpPolicies, err := p.loadSCP(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("load SCP: %w", err)
	}
	chain.OrganizationSCP = scpPolicies

	// 2️⃣ IdentityPolicies：User 直接 → Group → Role
	// User 直接策略
	for _, ap := range u.Edges.AccessPolicies {
		if p.isPolicyActive(ap) {
			chain.AddIdentityPolicy(p.toPolicy(ap))
		}
	}

	// Group 策略（用户所属组的策略）
	if u.Edges.Group != nil {
		groupPolicies, err := p.client.AccessPolicy.Query().
			Where(accesspolicy.HasGroupsWith(group.IDEQ(u.Edges.Group.ID))).
			All(ctx)
		if err != nil {
			return nil, fmt.Errorf("query group policies: %w", err)
		}
		for _, ap := range groupPolicies {
			if p.isPolicyActive(ap) {
				chain.AddIdentityPolicy(p.toPolicy(ap))
			}
		}
	}

	// Role 策略 + 权限边界
	for _, r := range u.Edges.Roles {
		if r.Edges.PermissionsBoundary != nil {
			bp := p.toPolicy(r.Edges.PermissionsBoundary)
			chain.PermissionBoundary = &bp
		}
		for _, ap := range r.Edges.AccessPolicies {
			if p.isPolicyActive(ap) {
				chain.AddIdentityPolicy(p.toPolicy(ap))
			}
		}
	}

	return chain, nil
}

// GetResourcePolicy 加载目标资源的资源策略（ResourcePolicy）
func (p *EntPolicyProvider) GetResourcePolicy(ctx context.Context, resourceURN string) (*Policy, error) {
	r, err := p.client.Resource.Query().
		Where(resource.Urn(resourceURN)).
		WithPolicies().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("query resource %s: %w", resourceURN, err)
	}

	if len(r.Edges.Policies) == 0 {
		return nil, nil
	}

	// 合并资源上的所有策略为一个 Policy 对象
	var allStmts []Statement
	policyID := "rp:" + r.ID
	for _, ap := range r.Edges.Policies {
		if !p.isPolicyActive(ap) {
			continue
		}
		for _, s := range ap.Statements {
			allStmts = append(allStmts, Statement{
				Effect:    Effect(s.Effect),
				Actions:   s.Actions,
				Resources: s.Resources,
				Condition: convertJSONCondition(s.Condition),
			})
		}
	}

	if len(allStmts) == 0 {
		return nil, nil
	}

	return &Policy{
		ID:         policyID,
		Name:       "resource-policy:" + r.Urn,
		Statements: allStmts,
	}, nil
}

// GetRoleTrustPolicy 获取角色的信任策略（用于 PassRole / AssumeRole 检查）
func (p *EntPolicyProvider) GetRoleTrustPolicy(ctx context.Context, roleID string) (interface{}, error) {
	r, err := p.client.Role.Query().
		Where(role.IDEQ(roleID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("query role %s: %w", roleID, err)
	}
	return r.TrustPolicy, nil
}

// CheckPassRole 检查用户是否有权将指定角色传递给服务
// 同时检查：① User/Role 必须有 iam:PassRole 权限 ② Role 的 trust_policy 必须信任该用户
func (p *EntPolicyProvider) CheckPassRole(ctx context.Context, userID, roleID, service string) (*Result, error) {
	engine := New()

	// 1️⃣ 获取用户策略链
	chain, err := p.GetEvalChain(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2️⃣ 获取角色的 trust_policy
	tp, err := p.GetRoleTrustPolicy(ctx, roleID)
	if err != nil {
		return nil, err
	}

	roleURN := BuildRoleURN("*", "*", roleID)

	// 3️⃣ 检查 iam:PassRole 权限
	passReq := PassRoleRequest(userID, roleURN, service)
	passResult := engine.CheckPassRole(passReq, chain, roleURN)
	if passResult.Decision != DecisionAllow {
		return passResult, nil
	}

	// 4️⃣ 检查信任策略（CanAssumeRole）
	if tp != nil {
		trustResult := engine.CanAssumeRole(userID, tp, passReq)
		if trustResult.Decision != DecisionAllow {
			return trustResult, nil
		}
	}

	return &Result{Decision: DecisionAllow, MatchedBy: "PassRole", Reason: "iam:PassRole allowed + trust policy matched"}, nil
}

// ─── 内部辅助 ────────────────────────────────────────────────

// loadSCP 从用户所属租户加载 SCP
func (p *EntPolicyProvider) loadSCP(ctx context.Context, u *ent.User) ([]Policy, error) {
	// 通过用户所属 Group 查找关联的租户
	// 当前 Schema: User → Group, Group → Tenant 关系未直接建立
	// 使用 Resource→Tenant 的规则：查找用户可访问的资源所在租户的策略
	// 简化方案：从所有活跃租户加载策略作为 SCP
	tenants, err := p.client.Tenant.Query().
		Where(tenant.StatusEQ(tenant.StatusActive)).
		WithAccessPolicies().
		All(ctx)
	if err != nil {
		return nil, err
	}

	var scps []Policy
	for _, t := range tenants {
		for _, ap := range t.Edges.AccessPolicies {
			if !p.isPolicyActive(ap) {
				continue
			}
			scps = append(scps, p.toPolicy(ap))
		}
	}
	return scps, nil
}

// isPolicyActive 检查策略是否在有效期内
func (p *EntPolicyProvider) isPolicyActive(ap *ent.AccessPolicy) bool {
	now := time.Now()
	if ap.EffectiveDate != nil && now.Before(*ap.EffectiveDate) {
		return false
	}
	if ap.ExpiryDate != nil && now.After(*ap.ExpiryDate) {
		return false
	}
	return true
}

// toPolicy 将 ent AccessPolicy 转换为 iam.Policy
func (p *EntPolicyProvider) toPolicy(ap *ent.AccessPolicy) Policy {
	var stmts []Statement
	for _, s := range ap.Statements {
		stmt := Statement{
			Effect:    Effect(s.Effect),
			Actions:   s.Actions,
			Resources: s.Resources,
		}
		if s.Condition != nil {
			stmt.Condition = convertJSONCondition(s.Condition)
		}
		stmts = append(stmts, stmt)
	}
	return Policy{
		ID:         ap.ID,
		Name:       ap.Name,
		Statements: stmts,
		Priority:   ap.Priority,
	}
}

// ─── JSON Condition 转换 ──────────────────────────────────────

func convertJSONCondition(raw interface{}) *Condition {
	if raw == nil {
		return nil
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	var c Condition
	if err := json.Unmarshal(data, &c); err != nil {
		return nil
	}
	return &c
}

// ─── ABAC 标签加载 ────────────────────────────────────────────

// LoadUserTags 将用户 attributes 加载为 ABAC 标签
func LoadUserTags(attrs map[string]interface{}) map[string]interface{} {
	ctx := make(map[string]interface{})
	if attrs == nil {
		return ctx
	}
	for k, v := range attrs {
		ctx[CtxPrincipalTag+"/"+k] = v
	}
	return ctx
}

// LoadResourceTags 将资源 details 加载为 ABAC 标签
func LoadResourceTags(details map[string]interface{}) map[string]interface{} {
	ctx := make(map[string]interface{})
	if details == nil {
		return ctx
	}
	for k, v := range details {
		ctx[CtxResourceTag+"/"+k] = v
	}
	return ctx
}

// ─── Ent 查询辅助 ─────────────────────────────────────────────

// RoleByID 按 ID 查询角色（含 trust_policy）
func RoleByID(ctx context.Context, client *ent.Client, roleID string) (*ent.Role, error) {
	return client.Role.Query().Where(role.IDEQ(roleID)).Only(ctx)
}

// ResourceByURN 按 URN 查询资源（含 policies 边缘）
func ResourceByURN(ctx context.Context, client *ent.Client, urn string) (*ent.Resource, error) {
	return client.Resource.Query().
		Where(resource.Urn(urn)).
		WithPolicies().
		WithTenant().
		Only(ctx)
}

// UserByIDWithChain 按 ID 查询用户（含完整策略链所需的边缘）
func UserByIDWithChain(ctx context.Context, client *ent.Client, userID string) (*ent.User, error) {
	return client.User.Query().
		Where(user.IDEQ(userID)).
		WithGroup(func(gq *ent.GroupQuery) {
			gq.WithAccessPolicies()
		}).
		WithRoles(func(rq *ent.RoleQuery) {
			rq.WithAccessPolicies()
			rq.WithPermissionsBoundary()
		}).
		WithAccessPolicies().
		Only(ctx)
}

// GroupsByUserID 查询用户所属的所有组（含策略）
func GroupsByUserID(ctx context.Context, client *ent.Client, userID string) ([]*ent.Group, error) {
	return client.Group.Query().
		Where(group.HasUsersWith(user.IDEQ(userID))).
		WithAccessPolicies().
		All(ctx)
}

// ─── 批量查询 ────────────────────────────────────────────────

// AllActivePolicies 获取所有活跃策略
func AllActivePolicies(ctx context.Context, client *ent.Client) ([]*ent.AccessPolicy, error) {
	now := time.Now()
	return client.AccessPolicy.Query().
		Where(
			accesspolicy.Or(
				accesspolicy.EffectiveDateIsNil(),
				accesspolicy.EffectiveDateLTE(now),
			),
			accesspolicy.Or(
				accesspolicy.ExpiryDateIsNil(),
				accesspolicy.ExpiryDateGTE(now),
			),
		).
		All(ctx)
}

// ─── SQL 辅助（用于复杂查询） ────────────────────────────────

// UserEffectivePolicies 通过原生 SQL 获取用户的有效策略列表
// 返回所有直接/间接附加到用户的策略 ID 列表
func (p *EntPolicyProvider) UserEffectivePolicies(ctx context.Context, userID string) ([]string, error) {
	var policyIDs []string
	err := p.client.AccessPolicy.Query().
		Where(accesspolicy.Or(
			// 用户直接策略
			accesspolicy.HasUsersWith(user.IDEQ(userID)),
			// 用户组策略
			accesspolicy.HasGroupsWith(group.HasUsersWith(user.IDEQ(userID))),
			// 用户角色策略
			accesspolicy.HasRolesWith(role.HasUsersWith(user.IDEQ(userID))),
		)).
		Select("id").
		Scan(ctx, &policyIDs)
	if err != nil {
		return nil, err
	}
	return policyIDs, nil
}
