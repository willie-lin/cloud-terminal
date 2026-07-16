package iam

import "context"

// ─── Provider 接口 ─────────────────────────────────────────

// Provider 策略提供者接口
type Provider interface {
	// GetEvalChain 获取指定主体的身份策略链（SCP + Boundary + IdentityPolicies）
	GetEvalChain(ctx context.Context, principalID string) (*EvalChain, error)

	// GetResourcePolicy 获取目标资源的资源策略（ResourcePolicy）
	GetResourcePolicy(ctx context.Context, resourceURN string) (*Policy, error)

	// GetRoleTrustPolicy 获取角色的信任策略（用于 PassRole / AssumeRole 检查）
	GetRoleTrustPolicy(ctx context.Context, roleID string) (interface{}, error)
}

// ProviderFunc 函数式 Provider
type ProviderFunc func(ctx context.Context, principalID string) (*EvalChain, error)

func (f ProviderFunc) GetEvalChain(ctx context.Context, principalID string) (*EvalChain, error) {
	return f(ctx, principalID)
}

func (ProviderFunc) GetResourcePolicy(_ context.Context, _ string) (*Policy, error) {
	return nil, nil
}

func (ProviderFunc) GetRoleTrustPolicy(_ context.Context, _ string) (interface{}, error) {
	return nil, nil
}

// ─── MemoryProvider ────────────────────────────────────────

// MemoryProvider 内存策略提供者（测试/默认环境用）
type MemoryProvider struct {
	chain            *EvalChain
	resourcePolicies map[string]Policy
	trustPolicies    map[string]interface{}
}

func NewMemoryProvider() *MemoryProvider {
	return &MemoryProvider{
		chain:            &EvalChain{},
		resourcePolicies: make(map[string]Policy),
		trustPolicies:    make(map[string]interface{}),
	}
}

func (p *MemoryProvider) GetEvalChain(_ context.Context, _ string) (*EvalChain, error) {
	return p.chain, nil
}

func (p *MemoryProvider) GetResourcePolicy(_ context.Context, resourceURN string) (*Policy, error) {
	if pol, ok := p.resourcePolicies[resourceURN]; ok {
		return &pol, nil
	}
	return nil, nil
}

func (p *MemoryProvider) GetRoleTrustPolicy(_ context.Context, roleID string) (interface{}, error) {
	if tp, ok := p.trustPolicies[roleID]; ok {
		return tp, nil
	}
	return nil, nil
}

func (p *MemoryProvider) AddIdentityPolicy(pol Policy) {
	p.chain.AddIdentityPolicy(pol)
}

func (p *MemoryProvider) SetResourcePolicy(urn string, pol Policy) {
	p.resourcePolicies[urn] = pol
}

func (p *MemoryProvider) SetPermissionBoundary(pol Policy) {
	p.chain.PermissionBoundary = &pol
}

func (p *MemoryProvider) SetTrustPolicy(roleID string, tp interface{}) {
	p.trustPolicies[roleID] = tp
}

// ─── Evaluator ─────────────────────────────────────────────

// Evaluator 带 Provider 的一步评估器
type Evaluator struct {
	engine   *Engine
	provider Provider
}

// NewEvaluator 创建带数据源的评估器
func NewEvaluator(provider Provider) *Evaluator {
	return &Evaluator{
		engine:   New(),
		provider: provider,
	}
}

// Evaluate 从 Provider 获取策略链并执行评估
func (ev *Evaluator) Evaluate(ctx context.Context, req *Request) (*Result, error) {
	chain, err := ev.provider.GetEvalChain(ctx, req.PrincipalID)
	if err != nil {
		return nil, err
	}

	// 如果指定了资源 URN，加载资源策略
	if req.ResourceURN != "" {
		rp, err := ev.provider.GetResourcePolicy(ctx, req.ResourceURN)
		if err != nil {
			return nil, err
		}
		chain.ResourcePolicy = rp
	}

	result := ev.engine.Evaluate(req, chain)

	// 审计日志
	if audit := AuditFromContext(ctx); audit != nil {
		audit.Log(AuditEntry{
			PrincipalID: req.PrincipalID,
			Action:      req.Action,
			ResourceURN: req.ResourceURN,
			Decision:    result.Decision,
			MatchedBy:   result.MatchedBy,
			PolicyID:    result.PolicyID,
			Reason:      result.Reason,
		})
	}

	return result, nil
}

// Simulate 模拟评估
func (ev *Evaluator) Simulate(ctx context.Context, req *Request, actions, resources []string) ([]SimulateResult, error) {
	chain, err := ev.provider.GetEvalChain(ctx, req.PrincipalID)
	if err != nil {
		return nil, err
	}
	return ev.engine.Simulate(req, chain, actions, resources), nil
}
