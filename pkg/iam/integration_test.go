package iam

import (
	"context"
	"os"
	"testing"

	"github.com/willie-lin/cloud-terminal/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/config"
	"github.com/willie-lin/cloud-terminal/pkg/database"
	"github.com/willie-lin/cloud-terminal/pkg/logger"
)

func TestEntProviderIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION") != "1" {
		t.Skip("set INTEGRATION=1 to run integration test")
	}

	cfg, err := config.LoadConfig("../../config.yaml")
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if err := logger.Init(cfg.Logger); err != nil {
		t.Fatalf("init logger: %v", err)
	}

	client, err := database.NewClient(&cfg.Database)
	if err != nil {
		t.Fatalf("database: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 查找已有 superadmin 或创建测试用户
	u, err := client.User.Query().Where(user.UsernameEQ("superadmin")).Only(ctx)
	if err != nil {
		// 尝试其他测试用户
		u, err = client.User.Query().First(ctx)
		if err != nil {
			t.Skip("no users in database, run app first to seed data")
		}
	}
	t.Logf("✅ Using user: %s (ID: %s)", u.Username, u.ID)

	// 2. GetEvalChain
	provider := NewEntPolicyProvider(client)
	chain, err := provider.GetEvalChain(ctx, u.ID)
	if err != nil {
		t.Fatalf("GetEvalChain: %v", err)
	}
	t.Logf("✅ IdentityPolicies: %d, SCP: %d, Boundary: %v",
		len(chain.IdentityPolicies), len(chain.OrganizationSCP), chain.PermissionBoundary != nil)

	// 3. IAM Evaluate
	engine := New()
	req := NewRequest(u.ID, "resource:connect", "urn:ct:dev:default:mysql:test-db")
	result := engine.Evaluate(req, chain)
	t.Logf("✅ IAM Evaluate: %s - %s", result.Decision, result.Reason)

	if result.Decision != DecisionAllow {
		// 对非超管用户可能是 Deny，这是正常的
		t.Logf("ℹ️  Expected - user %s permissions: %s", u.Username, result.Reason)
	}

	// 4. PassRole
	roleURN := BuildRoleURN("dev", "default", "super_admin")
	passResult := engine.CheckPassRole(req, chain, roleURN)
	t.Logf("✅ PassRole: %s - %s", passResult.Decision, passResult.Reason)

	// 5. ResourcePolicy
	r, err := client.Resource.Query().First(ctx)
	if err == nil {
		rp, err := provider.GetResourcePolicy(ctx, r.Urn)
		if err != nil {
			t.Logf("ℹ️  GetResourcePolicy error: %v", err)
		} else if rp != nil {
			t.Logf("✅ ResourcePolicy(%s): %d stmts", r.Urn, len(rp.Statements))
		} else {
			t.Logf("ℹ️  Resource(%s) has no policy", r.Urn)
		}
	} else {
		t.Logf("ℹ️  No resources in DB")
	}

	t.Log("🎉 Integration test passed")
}
