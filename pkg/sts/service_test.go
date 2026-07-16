package sts

import (
	"context"
	"testing"
	"time"

	"github.com/willie-lin/cloud-terminal/pkg/iam"
)

func TestIssueAndValidate(t *testing.T) {
	s := New([]byte("test-secret-key-32-bytes-long!!!"))

	chain := &iam.EvalChain{
		IdentityPolicies: []iam.Policy{{
			ID: "allow-all",
			Statements: []iam.Statement{
				{Effect: iam.EffectAllow, Actions: []string{"*"}, Resources: []string{"*"}},
			},
		}},
	}

	resp, err := s.IssueToken(context.Background(), &IssueRequest{
		UserID:      "user-123",
		ResourceURN: "urn:ct:prod:sh:mysql:db-001",
		TTL:         30 * time.Minute,
	}, iam.New(), chain)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Token == "" {
		t.Fatal("expected non-empty token")
	}
	if resp.SessionID == "" {
		t.Fatal("expected non-empty session id")
	}
	if resp.ExpiresAt.Before(time.Now()) {
		t.Fatal("token should expire in the future")
	}

	claims, err := s.ValidateToken(resp.Token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != "user-123" {
		t.Errorf("expected user-123, got %s", claims.UserID)
	}
	if claims.ResourceURN != "urn:ct:prod:sh:mysql:db-001" {
		t.Errorf("expected urn:ct:prod:sh:mysql:db-001, got %s", claims.ResourceURN)
	}
}

func TestExpiredToken(t *testing.T) {
	s := New([]byte("test-secret-key-32-bytes-long!!!"))
	s.defaultTTL = 1 // 1 nanosecond → 立即过期

	chain := &iam.EvalChain{
		IdentityPolicies: []iam.Policy{{
			ID: "allow-all",
			Statements: []iam.Statement{
				{Effect: iam.EffectAllow, Actions: []string{"*"}, Resources: []string{"*"}},
			},
		}},
	}

	resp, err := s.IssueToken(context.Background(), &IssueRequest{
		UserID:      "user-123",
		ResourceURN: "urn:ct:prod:sh:mysql:db-001",
	}, iam.New(), chain)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	_, err = s.ValidateToken(resp.Token)
	if err == nil {
		t.Error("expired token should be rejected")
	}
}

func TestRevokeToken(t *testing.T) {
	s := New([]byte("test-secret-key-32-bytes-long!!!"))

	chain := &iam.EvalChain{
		IdentityPolicies: []iam.Policy{{
			ID: "allow-all",
			Statements: []iam.Statement{
				{Effect: iam.EffectAllow, Actions: []string{"*"}, Resources: []string{"*"}},
			},
		}},
	}

	resp, err := s.IssueToken(context.Background(), &IssueRequest{
		UserID:      "user-123",
		ResourceURN: "urn:ct:prod:sh:mysql:db-001",
	}, iam.New(), chain)
	if err != nil {
		t.Fatal(err)
	}

	s.RevokeSession(resp.SessionID)

	_, err = s.ValidateToken(resp.Token)
	if err == nil {
		t.Error("revoked token should be rejected")
	}
}

func TestSessionPolicy(t *testing.T) {
	s := New([]byte("test-secret-key-32-bytes-long!!!"))

	chain := &iam.EvalChain{
		IdentityPolicies: []iam.Policy{{
			ID: "allow-all",
			Statements: []iam.Statement{
				{Effect: iam.EffectAllow, Actions: []string{"*"}, Resources: []string{"*"}},
			},
		}},
	}

	sessionPolicy := []iam.Statement{
		{
			Effect:    iam.EffectAllow,
			Actions:   []string{"resource:connect"},
			Resources: []string{"urn:ct:prod:sh:mysql:db-001"},
		},
	}

	resp, err := s.IssueToken(context.Background(), &IssueRequest{
		UserID:        "user-123",
		ResourceURN:   "urn:ct:prod:sh:mysql:db-001",
		SessionPolicy: sessionPolicy,
	}, iam.New(), chain)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := s.ValidateToken(resp.Token)
	if err != nil {
		t.Fatal(err)
	}

	pol, err := claims.GetSessionPolicy()
	if err != nil {
		t.Fatal(err)
	}
	if len(pol) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(pol))
	}
	if pol[0].Effect != iam.EffectAllow {
		t.Errorf("expected Allow effect")
	}

	newChain := BuildEvalChain(claims, chain)
	if len(newChain.SessionPolicies) != 1 {
		t.Errorf("expected 1 session policy, got %d", len(newChain.SessionPolicies))
	}
}

func TestIAMDeny(t *testing.T) {
	s := New([]byte("test-secret-key-32-bytes-long!!!"))
	chain := &iam.EvalChain{}

	_, err := s.IssueToken(context.Background(), &IssueRequest{
		UserID:      "user-123",
		ResourceURN: "urn:ct:prod:sh:mysql:db-001",
	}, iam.New(), chain)
	if err == nil {
		t.Error("should deny when IAM denies")
	}
}

func TestWrongSecret(t *testing.T) {
	s1 := New([]byte("secret-one-!!!!!!!!!!!!!!!!!!!!!!"))
	s2 := New([]byte("secret-two-!!!!!!!!!!!!!!!!!!!!!!"))

	chain := &iam.EvalChain{
		IdentityPolicies: []iam.Policy{{
			ID: "allow-all",
			Statements: []iam.Statement{
				{Effect: iam.EffectAllow, Actions: []string{"*"}, Resources: []string{"*"}},
			},
		}},
	}

	resp, err := s1.IssueToken(context.Background(), &IssueRequest{
		UserID:      "user-123",
		ResourceURN: "urn:ct:prod:sh:mysql:db-001",
	}, iam.New(), chain)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s2.ValidateToken(resp.Token)
	if err == nil {
		t.Error("should reject token signed with different secret")
	}
}

func TestMaxTTL(t *testing.T) {
	s := New([]byte("test-secret-key-32-bytes-long!!!"))

	chain := &iam.EvalChain{
		IdentityPolicies: []iam.Policy{{
			ID: "allow-all",
			Statements: []iam.Statement{
				{Effect: iam.EffectAllow, Actions: []string{"*"}, Resources: []string{"*"}},
			},
		}},
	}

	resp, err := s.IssueToken(context.Background(), &IssueRequest{
		UserID:      "user-123",
		ResourceURN: "urn:ct:prod:sh:mysql:db-001",
		TTL:         100 * time.Hour,
	}, iam.New(), chain)
	if err != nil {
		t.Fatal(err)
	}

	maxExpiry := time.Now().Add(MaxTTL).Add(1 * time.Minute)
	if resp.ExpiresAt.After(maxExpiry) {
		t.Errorf("TTL should be capped at %v, got %v", MaxTTL, resp.ExpiresAt)
	}
}
