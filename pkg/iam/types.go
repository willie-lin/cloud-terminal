package iam

import "context"

import "fmt"

// ─── Effect ─────────────────────────────────────────────────────

// Effect 策略效果
type Effect string

const (
	EffectAllow Effect = "Allow"
	EffectDeny  Effect = "Deny"
)

// ─── Decision ───────────────────────────────────────────────────

// Decision 鉴权决策
type Decision string

const (
	DecisionAllow Decision = "Allow"
	DecisionDeny  Decision = "Deny"
)

func (d Decision) String() string { return string(d) }

// IsAllowed 是否允许
func (d Decision) IsAllowed() bool { return d == DecisionAllow }

// ─── Request ────────────────────────────────────────────────────

// Request 鉴权请求
type Request struct {
	PrincipalID string
	Action      string
	ResourceURN string
	Context     map[string]interface{} // 运行时上下文（IP, 时间, MFA, 标签...）
}

func NewRequest(principalID, action, resourceURN string) *Request {
	return &Request{
		PrincipalID: principalID,
		Action:      action,
		ResourceURN: resourceURN,
		Context:     make(map[string]interface{}),
	}
}

// WithContext 设置上下文键值
func (r *Request) WithContext(key string, val interface{}) *Request {
	if r.Context == nil {
		r.Context = make(map[string]interface{})
	}
	r.Context[key] = val
	return r
}

// ─── Result ─────────────────────────────────────────────────────

// Result 单次鉴权结果
type Result struct {
	Decision  Decision `json:"decision"`
	MatchedBy string   `json:"matched_by,omitempty"` // 匹配的策略层
	PolicyID  string   `json:"policy_id,omitempty"`
	Statement int      `json:"statement,omitempty"`
	Reason    string   `json:"reason"`
}

func (r *Result) String() string {
	if r == nil {
		return "Deny (nil result)"
	}
	if r.Decision == DecisionAllow {
		return fmt.Sprintf("Allow [%s](policy=%s stmt=%d): %s", r.MatchedBy, r.PolicyID, r.Statement, r.Reason)
	}
	return fmt.Sprintf("Deny [%s]: %s", r.MatchedBy, r.Reason)
}

// SimulateResult 批量模拟结果
type SimulateResult struct {
	Action      string   `json:"action"`
	ResourceURN string   `json:"resource_urn"`
	Decision    Decision `json:"decision"`
	MatchedBy   string   `json:"matched_by,omitempty"`
	PolicyID    string   `json:"policy_id,omitempty"`
	Reason      string   `json:"reason"`
}

// ─── Statement ──────────────────────────────────────────────────

// Statement 策略语句
type Statement struct {
	Effect      Effect     `json:"effect"`
	Actions     []string   `json:"actions"`                // 允许/拒绝的操作列表
	NotActions  []string   `json:"not_actions,omitempty"`  // 排除的操作列表
	Resources   []string   `json:"resources"`              // 允许/拒绝的资源列表（URN）
	NotResources []string  `json:"not_resources,omitempty"`// 排除的资源列表
	Condition   *Condition `json:"condition,omitempty"`    // 条件约束
	Principal   *Principal `json:"principal,omitempty"`    // 资源策略中使用，指定主体
}

// Principal 主体（用于 Resource-based Policy）
type Principal struct {
	AWS       []string `json:"aws,omitempty"`        // 用户/角色 ID
	Canonical []string `json:"canonical,omitempty"`
}

// ─── Policy ─────────────────────────────────────────────────────

// Policy 策略文档
type Policy struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Statements  []Statement `json:"statements"`
	Priority    int         `json:"priority"`
	Version     string      `json:"version,omitempty"` // 策略版本（预留）
}

// ─── 策略层标识 ────────────────────────────────────────────────

// PolicyLayer 策略在评估链中的层级
type PolicyLayer string

const (
	LayerSCP                PolicyLayer = "SCP"
	LayerPermissionBoundary PolicyLayer = "PermissionBoundary"
	LayerIdentity           PolicyLayer = "Identity"
	LayerResource           PolicyLayer = "Resource"
	LayerSession            PolicyLayer = "Session"
)

func (l PolicyLayer) String() string { return string(l) }

// ─── Condition ─────────────────────────────────────────────────

// Condition 策略条件
// 全部字段可选，nil 表示不检查该条件
type Condition struct {
	// IP 地址条件
	IpAddress    map[string][]string `json:"ip_address,omitempty"`
	NotIpAddress map[string][]string `json:"not_ip_address,omitempty"`

	// 字符串条件
	StringEquals           map[string]string `json:"string_equals,omitempty"`
	StringNotEquals        map[string]string `json:"string_not_equals,omitempty"`
	StringEqualsIgnoreCase map[string]string `json:"string_equals_ignore_case,omitempty"`
	StringLike             map[string]string `json:"string_like,omitempty"`
	StringNotLike          map[string]string `json:"string_not_like,omitempty"`

	// URN/ARN 条件（使用 URN 通配匹配）
	ArnEquals    map[string]string `json:"arn_equals,omitempty"`
	ArnLike      map[string]string `json:"arn_like,omitempty"`
	ArnNotEquals map[string]string `json:"arn_not_equals,omitempty"`
	ArnNotLike   map[string]string `json:"arn_not_like,omitempty"`

	// 时间条件
	DateGreaterThan map[string]string `json:"date_greater_than,omitempty"` // RFC3339
	DateLessThan    map[string]string `json:"date_less_than,omitempty"`    // RFC3339
	DateEquals      map[string]string `json:"date_equals,omitempty"`
	DateNotEquals   map[string]string `json:"date_not_equals,omitempty"`
	DateBetween     *DateBetween      `json:"date_between,omitempty"`

	// Bool 条件
	Bool map[string]bool `json:"bool,omitempty"`

	// 数字条件
	NumericEquals       map[string]float64 `json:"numeric_equals,omitempty"`
	NumericNotEquals    map[string]float64 `json:"numeric_not_equals,omitempty"`
	NumericLessThan     map[string]float64 `json:"numeric_less_than,omitempty"`
	NumericGreaterThan  map[string]float64 `json:"numeric_greater_than,omitempty"`

	// Null 条件：检查 context 中的 key 是否存在
	Null map[string]bool `json:"null,omitempty"`

	// MFA 条件
	RequireMFA *bool `json:"require_mfa,omitempty"`

	// Set 运算符
	ForAllValues *ForAllValues `json:"for_all_values,omitempty"`
	ForAnyValue  *ForAnyValue  `json:"for_any_value,omitempty"`

	// ABAC 标签条件
	PrincipalTag map[string]string `json:"principal_tag,omitempty"` // 主体标签匹配
	ResourceTag  map[string]string `json:"resource_tag,omitempty"`  // 资源标签匹配
	RequestTag   map[string]string `json:"request_tag,omitempty"`   // 请求标签匹配
}

// DateBetween 时间范围
type DateBetween struct {
	Key    string `json:"key"`
	After  string `json:"after,omitempty"`  // RFC3339
	Before string `json:"before,omitempty"` // RFC3339
}

// ForAllValues Set 集合运算符：请求中的所有值必须在策略值中
type ForAllValues struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// ForAnyValue Set 集合运算符：请求中的至少一个值在策略值中
type ForAnyValue struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// ─── 评估链（Evaluation Chain） ────────────────────────────────

// EvalChain 完整评估链，对应 AWS IAM 评估层次
type EvalChain struct {
	// OrganizationSCP 组织级控制策略（最高层，仅能 Deny）
	OrganizationSCP []Policy `json:"organization_scp,omitempty"`

	// PermissionBoundary 权限边界（定义最大许可范围）
	PermissionBoundary *Policy `json:"permission_boundary,omitempty"`

	// IdentityPolicies 身份策略（User → Group → Role 继承链）
	IdentityPolicies []Policy `json:"identity_policies,omitempty"`

	// ResourcePolicy 资源策略（附加在被访问资源上）
	ResourcePolicy *Policy `json:"resource_policy,omitempty"`

	// SessionPolicies 会话策略（Role Assume 时传入的限制）
	SessionPolicies []Policy `json:"session_policies,omitempty"`
}

// AddIdentityPolicy 添加身份策略
func (c *EvalChain) AddIdentityPolicy(p Policy) {
	c.IdentityPolicies = append(c.IdentityPolicies, p)
}

// AddSCP 添加组织级策略
func (c *EvalChain) AddSCP(p Policy) {
	c.OrganizationSCP = append(c.OrganizationSCP, p)
}

// ─── 上下文常量 ────────────────────────────────────────────────

// 预定义的上下文键
const (
	CtxSourceIP     = "source_ip"
	CtxCurrentTime  = "current_time"
	CtxEpochTime    = "epoch_time"
	CtxUsername     = "aws:username"
	CtxUserID       = "aws:userid"
	CtxRoleName     = "aws:rolename"
	CtxMFAAuth      = "mfa"
	CtxUserAgent    = "user_agent"
	CtxPrincipalTag = "aws:PrincipalTag" // 前缀
	CtxResourceTag  = "aws:ResourceTag"  // 前缀
	CtxRequestTag   = "aws:RequestTag"   // 前缀
)

// ─── 审计日志 ─────────────────────────────────────────────────

// AuditEntry 一次鉴权的完整审计记录
type AuditEntry struct {
	PrincipalID string   `json:"principal_id"`
	Action      string   `json:"action"`
	ResourceURN string   `json:"resource_urn"`
	Decision    Decision `json:"decision"`
	MatchedBy   string   `json:"matched_by,omitempty"`
	PolicyID    string   `json:"policy_id,omitempty"`
	Reason      string   `json:"reason"`
}

// AuditLogger 审计日志接口
type AuditLogger interface {
	Log(entry AuditEntry)
}

// AuditContext 从 context 中提取审计日志器
type auditContextKey struct{}

func WithAuditLogger(ctx context.Context, logger AuditLogger) context.Context {
	return context.WithValue(ctx, auditContextKey{}, logger)
}

func AuditFromContext(ctx context.Context) AuditLogger {
	if ctx == nil {
		return nil
	}
	v, _ := ctx.Value(auditContextKey{}).(AuditLogger)
	return v
}
