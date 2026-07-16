package iam

import (
	"context"
	"testing"
	"time"
)

// ─── 辅助 ──────────────────────────────────────────────────

func allow(id string) Policy {
	return Policy{ID: id, Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"*"}, Resources: []string{"*"}},
	}}
}

func deny(id string) Policy {
	return Policy{ID: id, Statements: []Statement{
		{Effect: EffectDeny, Actions: []string{"*"}, Resources: []string{"*"}},
	}}
}

func chain(identity ...Policy) *EvalChain {
	return &EvalChain{IdentityPolicies: identity}
}

// ─── 三定律 ────────────────────────────────────────────────

func TestDefaultDeny(t *testing.T) {
	e := New()
	r := e.Evaluate(NewRequest("u1", "connect", "urn:ct:dev:default:ssh:s1"), &EvalChain{})
	if r.Decision != DecisionDeny || r.Reason != "default deny (no matching policy)" {
		t.Fatalf("expected default deny, got %s", r)
	}
}

func TestExplicitAllow(t *testing.T) {
	e := New()
	r := e.Evaluate(NewRequest("u1", "connect", "urn:ct:dev:default:ssh:s1"), chain(allow("p1")))
	if r.Decision != DecisionAllow {
		t.Fatalf("expected allow: %s", r)
	}
	if r.PolicyID != "p1" {
		t.Fatalf("expected p1, got %s", r.PolicyID)
	}
}

func TestExplicitDenyOverrides(t *testing.T) {
	e := New()
	r := e.Evaluate(NewRequest("u1", "connect", "x"), chain(
		allow("allow-all"),
		deny("deny-all"),
	))
	if r.Decision != DecisionDeny || r.PolicyID != "deny-all" {
		t.Fatalf("explicit deny should win: %s", r)
	}
}

// ─── 通配匹配 ──────────────────────────────────────────────

func TestWildcardAction(t *testing.T) {
	e := New()
	pol := Policy{ID: "p1", Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"resource:*"}, Resources: []string{"*"}},
	}}
	tests := []struct {
		action string
		allow  bool
	}{
		{"resource:connect", true},
		{"resource:delete", true},
		{"admin:delete", false},
	}
	for _, tc := range tests {
		r := e.Evaluate(NewRequest("u", tc.action, "*"), chain(pol))
		if (r.Decision == DecisionAllow) != tc.allow {
			t.Errorf("action=%s expect allow=%v got %s", tc.action, tc.allow, r)
		}
	}
}

func TestWildcardResource(t *testing.T) {
	e := New()
	pol := Policy{ID: "p1", Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"*"}, Resources: []string{"urn:ct:prod:*:*:*"}},
	}}
	if r := e.Evaluate(NewRequest("u", "x", "urn:ct:prod:sh:mysql:db"), chain(pol)); r.Decision != DecisionAllow {
		t.Errorf("prod resource should allow: %s", r)
	}
	if r := e.Evaluate(NewRequest("u", "x", "urn:ct:dev:sh:mysql:db"), chain(pol)); r.Decision != DecisionDeny {
		t.Errorf("dev resource should deny: %s", r)
	}
}

// ─── NotAction / NotResource ───────────────────────────────

func TestNotAction(t *testing.T) {
	e := New()
	pol := Policy{ID: "p1", Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"*"}, Resources: []string{"*"}},
		{Effect: EffectDeny, NotActions: []string{"resource:list"}, Resources: []string{"*"}},
	}}
	if r := e.Evaluate(NewRequest("u", "resource:list", "x"), chain(pol)); r.Decision != DecisionAllow {
		t.Errorf("resource:list should be excluded from deny: %s", r)
	}
	if r := e.Evaluate(NewRequest("u", "resource:delete", "x"), chain(pol)); r.Decision != DecisionDeny {
		t.Errorf("resource:delete should be denied: %s", r)
	}
}

// ─── 优先级 ────────────────────────────────────────────────

func TestPriority(t *testing.T) {
	e := New()
	r := e.Evaluate(NewRequest("u", "x", "x"), chain(
		Policy{ID: "low", Priority: 100, Statements: []Statement{
			{Effect: EffectAllow, Actions: []string{"*"}, Resources: []string{"*"}},
		}},
		Policy{ID: "high", Priority: 1, Statements: []Statement{
			{Effect: EffectDeny, Actions: []string{"*"}, Resources: []string{"*"}},
		}},
	))
	if r.Decision != DecisionDeny || r.PolicyID != "high" {
		t.Errorf("high priority deny should win: %s", r)
	}
}

// ─── SCP 层 ────────────────────────────────────────────────

func TestSCP(t *testing.T) {
	e := New()
	r := e.Evaluate(NewRequest("u", "x", "x"), &EvalChain{
		OrganizationSCP: []Policy{deny("scp-deny")},
		IdentityPolicies: []Policy{allow("allow-all")},
	})
	if r.Decision != DecisionDeny || r.PolicyID != "scp-deny" {
		t.Errorf("SCP deny should win: %s", r)
	}
}

// ─── Permission Boundary ──────────────────────────────────

func TestPermissionBoundary(t *testing.T) {
	e := New()
	boundary := Policy{ID: "boundary", Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"resource:connect"}, Resources: []string{"*"}},
	}}

	// 在边界内
	r := e.Evaluate(NewRequest("u", "resource:connect", "x"), &EvalChain{
		PermissionBoundary: &boundary,
		IdentityPolicies:   []Policy{allow("allow-all")},
	})
	if r.Decision != DecisionAllow {
		t.Errorf("within boundary should allow: %s", r)
	}

	// 超出边界
	r2 := e.Evaluate(NewRequest("u", "admin:delete", "x"), &EvalChain{
		PermissionBoundary: &boundary,
		IdentityPolicies:   []Policy{allow("allow-all")},
	})
	if r2.Decision != DecisionDeny {
		t.Errorf("outside boundary should deny: %s", r2)
	}
}

// ─── Resource Policy ──────────────────────────────────────

func TestResourcePolicy(t *testing.T) {
	e := New()
	rp := Policy{ID: "bucket-policy", Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"s3:GetObject"}, Resources: []string{"urn:ct:prod:*:s3:my-bucket/*"}},
	}}
	r := e.Evaluate(NewRequest("cross-account-user", "s3:GetObject", "urn:ct:prod:sh:s3:my-bucket/file.txt"), &EvalChain{
		ResourcePolicy: &rp,
	})
	if r.Decision != DecisionAllow {
		t.Errorf("resource policy should allow: %s", r)
	}
}

// ─── Session Policy ───────────────────────────────────────

func TestSessionPolicy(t *testing.T) {
	e := New()
	r := e.Evaluate(NewRequest("u", "connect", "x"), &EvalChain{
		IdentityPolicies: []Policy{allow("allow-all")},
		SessionPolicies:  []Policy{allow("session")},
	})
	if r.Decision != DecisionAllow {
		t.Errorf("identity+session allow: %s", r)
	}

	// Session 限制后拒绝
	r2 := e.Evaluate(NewRequest("u", "admin:delete", "x"), &EvalChain{
		IdentityPolicies: []Policy{allow("allow-all")},
		SessionPolicies: []Policy{{ID: "session", Statements: []Statement{
			{Effect: EffectDeny, Actions: []string{"admin:*"}, Resources: []string{"*"}},
		}}},
	})
	if r2.Decision != DecisionDeny {
		t.Errorf("session deny should win: %s", r2)
	}
}

// ─── Condition ─────────────────────────────────────────────

func TestConditionIpAddress(t *testing.T) {
	e := New()
	pol := Policy{ID: "p1", Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"*"}, Resources: []string{"*"},
			Condition: &Condition{IpAddress: map[string][]string{CtxSourceIP: {"10.0.0.0/8"}}},
		},
	}}
	r := e.Evaluate(NewRequest("u", "x", "*").WithContext(CtxSourceIP, "10.0.1.100"), chain(pol))
	if r.Decision != DecisionAllow {
		t.Errorf("internal IP allow: %s", r)
	}
	r2 := e.Evaluate(NewRequest("u", "x", "*").WithContext(CtxSourceIP, "8.8.8.8"), chain(pol))
	if r2.Decision != DecisionDeny {
		t.Errorf("external IP deny: %s", r2)
	}
}

func TestConditionMFA(t *testing.T) {
	mfa := true
	e := New()
	pol := Policy{ID: "p1", Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"admin:*"}, Resources: []string{"*"},
			Condition: &Condition{RequireMFA: &mfa},
		},
	}}
	r := e.Evaluate(NewRequest("u", "admin:delete", "*").WithContext(CtxMFAAuth, true), chain(pol))
	if r.Decision != DecisionAllow {
		t.Errorf("MFA allow: %s", r)
	}
	r2 := e.Evaluate(NewRequest("u", "admin:delete", "*").WithContext(CtxMFAAuth, false), chain(pol))
	if r2.Decision != DecisionDeny {
		t.Errorf("no MFA deny: %s", r2)
	}
}

func TestConditionDateBetween(t *testing.T) {
	pol := Policy{ID: "p1", Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"*"}, Resources: []string{"*"},
			Condition: &Condition{
				DateBetween: &DateBetween{
					Key: CtxCurrentTime,
					After:  time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
					Before: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
				},
			},
		},
	}}
	e := New()
	r := e.Evaluate(NewRequest("u", "x", "*"), chain(pol))
	if r.Decision != DecisionAllow {
		t.Errorf("within window: %s", r)
	}
}

func TestConditionStringLike(t *testing.T) {
	pol := Policy{ID: "p1", Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"*"}, Resources: []string{"*"},
			Condition: &Condition{StringLike: map[string]string{CtxUserID: "team-*"}},
		},
	}}
	e := New()
	r := e.Evaluate(NewRequest("u", "x", "*").WithContext(CtxUserID, "team-alpha"), chain(pol))
	if r.Decision != DecisionAllow {
		t.Errorf("team member allow: %s", r)
	}
	r2 := e.Evaluate(NewRequest("u", "x", "*").WithContext(CtxUserID, "external"), chain(pol))
	if r2.Decision != DecisionDeny {
		t.Errorf("non-team deny: %s", r2)
	}
}

func TestConditionArnLike(t *testing.T) {
	pol := Policy{ID: "p1", Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"*"}, Resources: []string{"*"},
			Condition: &Condition{ArnLike: map[string]string{"source_vpc": "urn:ct:prod:*:*:*"}},
		},
	}}
	e := New()
	r := e.Evaluate(NewRequest("u", "x", "*").WithContext("source_vpc", "urn:ct:prod:sh:vpc:main"), chain(pol))
	if r.Decision != DecisionAllow {
		t.Errorf("prod VPC allow: %s", r)
	}
	r2 := e.Evaluate(NewRequest("u", "x", "*").WithContext("source_vpc", "urn:ct:dev:bj:vpc:test"), chain(pol))
	if r2.Decision != DecisionDeny {
		t.Errorf("dev VPC deny: %s", r2)
	}
}

func TestConditionNull(t *testing.T) {
	// Null: {"mfa": true} → AWS 语义：mfa KEY 必须不存在
	pol := Policy{ID: "p1", Statements: []Statement{
		{Effect: EffectDeny, Actions: []string{"*"}, Resources: []string{"*"},
			Condition: &Condition{Null: map[string]bool{"mfa": true}},
		},
	}}
	e := New()
	// mfa key 存在 → Null 条件不匹配 → deny 未触发 → 无 allow → 默认拒绝
	r := e.Evaluate(NewRequest("u", "x", "*").WithContext("mfa", true), chain(pol))
	if r.Decision != DecisionDeny || r.Reason != "default deny (no matching policy)" {
		t.Errorf("MFA exists, Null unmatched, expected default deny: %s", r)
	}
	// mfa key 不存在 → Null 条件匹配 → deny 触发
	r2 := e.Evaluate(NewRequest("u", "x", "*"), chain(pol))
	if r2.Decision != DecisionDeny || r2.Reason != "Identity explicit deny" {
		t.Errorf("no MFA, Null matched, expected deny: %s", r2)
	}
}

// ─── Simulate API ─────────────────────────────────────────

func TestSimulate(t *testing.T) {
	e := New()
	results := e.Simulate(NewRequest("u", "", ""), chain(allow("p1")),
		[]string{"connect", "delete"},
		[]string{"urn:ct:prod:sh:mysql:db", "urn:ct:dev:default:ssh:s1"},
	)
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Decision != DecisionAllow {
			t.Errorf("all should allow: %s/%s = %s", r.Action, r.ResourceURN, r.Decision)
		}
	}
}

// ─── MemoryProvider ───────────────────────────────────────

func TestMemoryProvider(t *testing.T) {
	prov := NewMemoryProvider()
	prov.AddIdentityPolicy(allow("p1"))
	ev := NewEvaluator(prov)
	r, err := ev.Evaluate(context.Background(), NewRequest("u", "x", "*"))
	if err != nil {
		t.Fatal(err)
	}
	if r.Decision != DecisionAllow {
		t.Errorf("expected allow: %s", r)
	}
}

// ─── Policy 变量 ──────────────────────────────────────────

func TestPolicyVariables(t *testing.T) {
	e := New()
	pol := Policy{ID: "p1", Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"*"},
			Resources: []string{"urn:ct:prod:*:*:${aws:username}"},
		},
	}}
	r := e.Evaluate(NewRequest("u1", "x", "urn:ct:prod:sh:ssh:u1").WithContext(CtxUsername, "u1"), chain(pol))
	if r.Decision != DecisionAllow {
		t.Errorf("variable expansion should allow: %s", r)
	}
	r2 := e.Evaluate(NewRequest("u2", "x", "urn:ct:prod:sh:ssh:u1").WithContext(CtxUsername, "u2"), chain(pol))
	if r2.Decision != DecisionDeny {
		t.Errorf("wrong user should deny: %s", r2)
	}
}

func TestExpandCurrentTime(t *testing.T) {
	result := ExpandVariables("access-${aws:CurrentTime}", NewRequest("u", "x", "*"))
	if result == "access-${aws:CurrentTime}" {
		t.Error("CurrentTime should have been expanded")
	}
}

// ─── ABAC Tags ────────────────────────────────────────────

func TestPrincipalTagCondition(t *testing.T) {
	pol := Policy{ID: "p1", Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"*"}, Resources: []string{"*"},
			Condition: &Condition{PrincipalTag: map[string]string{"department": "engineering"}},
		},
	}}
	e := New()
	r := e.Evaluate(NewRequest("u", "x", "*").WithContext(CtxPrincipalTag+"/department", "engineering"), chain(pol))
	if r.Decision != DecisionAllow {
		t.Errorf("engineering tag allow: %s", r)
	}
	r2 := e.Evaluate(NewRequest("u", "x", "*").WithContext(CtxPrincipalTag+"/department", "sales"), chain(pol))
	if r2.Decision != DecisionDeny {
		t.Errorf("sales tag deny: %s", r2)
	}
}

// ─── Edge Cases ───────────────────────────────────────────

func TestEmptyChain(t *testing.T) {
	e := New()
	r := e.Evaluate(NewRequest("u", "x", "y"), &EvalChain{})
	if r.Decision != DecisionDeny {
		t.Errorf("empty chain should deny: %s", r)
	}
}

func TestMultipleStatements(t *testing.T) {
	e := New()
	pol := Policy{ID: "p1", Statements: []Statement{
		{Effect: EffectDeny, Actions: []string{"admin:*"}, Resources: []string{"*"}},
		{Effect: EffectAllow, Actions: []string{"*"}, Resources: []string{"*"}},
	}}
	r := e.Evaluate(NewRequest("u", "admin:delete", "x"), chain(pol))
	if r.Decision != DecisionDeny {
		t.Errorf("admin actions deny: %s", r)
	}
	r2 := e.Evaluate(NewRequest("u", "resource:list", "x"), chain(pol))
	if r2.Decision != DecisionAllow {
		t.Errorf("non-admin allow: %s", r2)
	}
}

func TestForAnyValue(t *testing.T) {
	pol := Policy{ID: "p1", Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"*"}, Resources: []string{"*"},
			Condition: &Condition{
				ForAnyValue: &ForAnyValue{Key: "services", Values: []string{"mysql", "redis"}},
			},
		},
	}}
	e := New()
	r := e.Evaluate(NewRequest("u", "x", "*").WithContext("services", []string{"mysql"}), chain(pol))
	if r.Decision != DecisionAllow {
		t.Errorf("mysql in services should allow: %s", r)
	}
	r2 := e.Evaluate(NewRequest("u", "x", "*").WithContext("services", []string{"kafka"}), chain(pol))
	if r2.Decision != DecisionDeny {
		t.Errorf("kafka not in services should deny: %s", r2)
	}
}

func TestValidatePolicy(t *testing.T) {
	errs := ValidatePolicy(Policy{ID: "", Statements: nil})
	if len(errs) == 0 {
		t.Error("should validate ID required")
	}
	errs = ValidatePolicy(Policy{ID: "p1", Statements: []Statement{
		{Effect: "Invalid", Actions: []string{"*"}, Resources: []string{"*"}},
	}})
	if len(errs) == 0 {
		t.Error("should validate effect")
	}
	errs = ValidatePolicy(Policy{ID: "p1", Statements: []Statement{
		{Effect: EffectAllow, Actions: []string{"*"}, Resources: []string{"*"}},
	}})
	if len(errs) != 0 {
		t.Errorf("valid policy should pass: %v", errs)
	}
}

// ─── TrustPolicy ──────────────────────────────────────────────

func TestCanAssumeRole(t *testing.T) {
	e := New()

	// 信任策略：允许 user-123 和所有 team-* 用户 Assume
	trustPolicy := []Statement{
		{
			Effect: EffectAllow,
			Principal: &Principal{
				AWS: []string{"user-123", "team-admin"},
			},
		},
		{
			Effect: EffectAllow,
			Principal: &Principal{
				AWS: []string{"*"},
			},
			Condition: &Condition{
				StringLike: map[string]string{CtxUsername: "team-*"},
			},
		},
	}

	// user-123 直接信任
	r := e.CanAssumeRole("user-123", trustPolicy, NewRequest("user-123", "sts:AssumeRole", "urn:ct:*:*:role:admin"))
	if r.Decision != DecisionAllow {
		t.Errorf("user-123 should be trusted: %s", r)
	}

	// team-alpha 通过通配 + 条件信任
	r2 := e.CanAssumeRole("team-alpha", trustPolicy, NewRequest("team-alpha", "sts:AssumeRole", "").WithContext(CtxUsername, "team-alpha"))
	if r2.Decision != DecisionAllow {
		t.Errorf("team-alpha should be trusted via condition: %s", r2)
	}

	// evil-user 不被信任
	r3 := e.CanAssumeRole("evil-user", trustPolicy, NewRequest("evil-user", "sts:AssumeRole", ""))
	if r3.Decision != DecisionDeny {
		t.Errorf("evil-user should be denied: %s", r3)
	}

	// 空信任策略
	r4 := e.CanAssumeRole("user-123", nil, NewRequest("user-123", "sts:AssumeRole", ""))
	if r4.Decision != DecisionDeny {
		t.Errorf("nil trust policy should deny: %s", r4)
	}
}

// ─── ABAC Tags 加载 ──────────────────────────────────────────

func TestLoadUserTags(t *testing.T) {
	attrs := map[string]interface{}{
		"department": "engineering",
		"level":      "senior",
	}
	ctx := LoadUserTags(attrs)
	if ctx[CtxPrincipalTag+"/department"] != "engineering" {
		t.Errorf("expected engineering tag, got %v", ctx[CtxPrincipalTag+"/department"])
	}
	if ctx[CtxPrincipalTag+"/level"] != "senior" {
		t.Errorf("expected senior tag, got %v", ctx[CtxPrincipalTag+"/level"])
	}
}

func TestLoadResourceTags(t *testing.T) {
	details := map[string]interface{}{
		"environment": "production",
		"backup":      true,
	}
	ctx := LoadResourceTags(details)
	if ctx[CtxResourceTag+"/environment"] != "production" {
		t.Errorf("expected production tag, got %v", ctx[CtxResourceTag+"/environment"])
	}
}

// ─── PassRole ────────────────────────────────────────────────

func TestPassRole(t *testing.T) {
	e := New()

	// 用户有 iam:PassRole 在 dev 环境的 role 上
	chain := &EvalChain{
		IdentityPolicies: []Policy{
			{ID: "allow-passrole", Statements: []Statement{
				{Effect: EffectAllow, Actions: []string{"iam:PassRole"},
					Resources: []string{"urn:ct:dev:*:role:*"}},
			}},
		},
	}

	// dev 环境 Role 允许 Pass
	r := e.CheckPassRole(NewRequest("u1", "", ""), chain, "urn:ct:dev:default:role:my-role")
	if r.Decision != DecisionAllow {
		t.Errorf("dev role should be passable: %s", r)
	}

	// prod 环境 Role 拒绝 Pass
	r2 := e.CheckPassRole(NewRequest("u1", "", ""), chain, "urn:ct:prod:default:role:my-role")
	if r2.Decision != DecisionDeny {
		t.Errorf("prod role should NOT be passable: %s", r2)
	}

	// nil chain → Deny
	r3 := e.CheckPassRole(NewRequest("u1", "", ""), nil, "urn:ct:dev:default:role:my-role")
	if r3.Decision != DecisionDeny {
		t.Errorf("nil chain should deny: %s", r3)
	}
}
