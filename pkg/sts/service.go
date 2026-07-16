// Package sts — Security Token Service
// 对标 AWS STS：签发临时凭据，用于 ContainerSSH 连接鉴权
package sts

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/willie-lin/cloud-terminal/pkg/iam"
)

// ─── 默认配置 ──────────────────────────────────────────────

const (
	DefaultTTL = 1 * time.Hour
	MaxTTL     = 8 * time.Hour
)

// ─── Service ────────────────────────────────────────────────

// Service STS 令牌服务
type Service struct {
	secret     []byte
	defaultTTL time.Duration
	maxTTL     time.Duration

	mu      sync.RWMutex
	revoked map[string]time.Time // sessionID → revokedAt
}

// New 创建 STS 服务
// secret 用于 JWT 签名（至少 32 字节）
func New(secret []byte) *Service {
	if len(secret) == 0 {
		secret = make([]byte, 32)
		rand.Read(secret)
	}
	return &Service{
		secret:     secret,
		defaultTTL: DefaultTTL,
		maxTTL:     MaxTTL,
		revoked:    make(map[string]time.Time),
	}
}

// ─── 请求与响应 ────────────────────────────────────────────

// IssueRequest 签发 Token 的请求
type IssueRequest struct {
	UserID      string          // 用户 ID
	ResourceURN string          // 目标资源 URN
	TTL         time.Duration  // 过期时间（0 使用默认）
	SessionPolicy []iam.Statement // 会话策略（可选，缩小权限用）
}

// IssueResponse 签发 Token 的响应
type IssueResponse struct {
	Token     string    `json:"token"`
	SessionID string    `json:"session_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ─── JWT Claims ────────────────────────────────────────────

// SessionClaims JWT 声明
type SessionClaims struct {
	jwt.RegisteredClaims
	SessionID     string `json:"sid"`
	UserID        string `json:"uid"`
	ResourceURN   string `json:"urn"`
	SessionPolicy string `json:"pol,omitempty"` // JSON 序列化的会话策略
}

// ─── IssueToken ────────────────────────────────────────────

// IssueToken 签发临时凭据
// 先通过 IAM 引擎鉴权，通过后签发 JWT
func (s *Service) IssueToken(ctx context.Context, req *IssueRequest, engine *iam.Engine, chain *iam.EvalChain) (*IssueResponse, error) {
	// 1. IAM 鉴权
	iamReq := iam.NewRequest(req.UserID, "resource:connect", req.ResourceURN)
	iamReq.Context = extractContext(ctx)
	result := engine.Evaluate(iamReq, chain)
	if result.Decision != iam.DecisionAllow {
		return nil, fmt.Errorf("sts: access denied: %s", result.Reason)
	}

	// 2. 计算过期时间
	ttl := req.TTL
	if ttl <= 0 {
		ttl = s.defaultTTL
	}
	if ttl > s.maxTTL {
		ttl = s.maxTTL
	}
	now := time.Now()
	expiresAt := now.Add(ttl)

	// 3. 生成 SessionID
	sessionID := uuid.New().String()

	// 4. 序列化 SessionPolicy（可选）
	var policyJSON string
	if len(req.SessionPolicy) > 0 {
		data, err := json.Marshal(req.SessionPolicy)
		if err != nil {
			return nil, fmt.Errorf("sts: marshal session policy: %w", err)
		}
		policyJSON = string(data)
	}

	// 5. 签发 JWT
	claims := &SessionClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        sessionID,
			Subject:   req.UserID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    "cloud-terminal",
		},
		SessionID:     sessionID,
		UserID:        req.UserID,
		ResourceURN:   req.ResourceURN,
		SessionPolicy: policyJSON,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.secret)
	if err != nil {
		return nil, fmt.Errorf("sts: sign token: %w", err)
	}

	return &IssueResponse{
		Token:     tokenStr,
		SessionID: sessionID,
		ExpiresAt: expiresAt,
	}, nil
}

// ─── ValidateToken ─────────────────────────────────────────

// ValidateToken 验证并解析临时凭据
// 检查：签名 + 过期时间 + 是否被撤销
func (s *Service) ValidateToken(tokenStr string) (*SessionClaims, error) {
	// 解析 JWT
	token, err := jwt.ParseWithClaims(tokenStr, &SessionClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("sts: unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("sts: invalid token: %w", err)
	}

	claims, ok := token.Claims.(*SessionClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("sts: invalid token claims")
	}

	// 检查撤销
	s.mu.RLock()
	_, revoked := s.revoked[claims.SessionID]
	s.mu.RUnlock()
	if revoked {
		return nil, fmt.Errorf("sts: session revoked: %s", claims.SessionID)
	}

	return claims, nil
}

// ─── RevokeSession ─────────────────────────────────────────

// RevokeSession 撤销一次会话
// 撤销后的 Token 将无法通过 ValidateToken
func (s *Service) RevokeSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.revoked[sessionID] = time.Now()
}

// ─── GetSessionPolicy ──────────────────────────────────────

// GetSessionPolicy 从 claims 中反序列化会话策略
func (c *SessionClaims) GetSessionPolicy() ([]iam.Statement, error) {
	if c.SessionPolicy == "" {
		return nil, nil
	}
	var stmts []iam.Statement
	if err := json.Unmarshal([]byte(c.SessionPolicy), &stmts); err != nil {
		return nil, err
	}
	return stmts, nil
}

// BuildEvalChain 用 STS claims + session policy 构建评估链
// 用于后续的 ContainerSSH webhook 鉴权
func BuildEvalChain(claims *SessionClaims, baseChain *iam.EvalChain) *iam.EvalChain {
	if baseChain == nil {
		baseChain = &iam.EvalChain{}
	}
	chain := *baseChain // shallow copy

	// 如果有 session policy，加到 SessionPolicies 层
	if claims.SessionPolicy != "" {
		var stmts []iam.Statement
		if err := json.Unmarshal([]byte(claims.SessionPolicy), &stmts); err == nil && len(stmts) > 0 {
			chain.SessionPolicies = append(chain.SessionPolicies, iam.Policy{
				ID:         "session-" + claims.SessionID,
				Name:       "Session Policy",
				Statements: stmts,
			})
		}
	}

	return &chain
}

// ─── 辅助 ──────────────────────────────────────────────────

func extractContext(ctx context.Context) map[string]interface{} {
	m := make(map[string]interface{})
	if ctx == nil {
		return m
	}
	// 从 context 中提取 IAM 需要的上下文
	// 调用方可以通过 context.WithValue 注入
	if v := ctx.Value("source_ip"); v != nil {
		m["source_ip"] = v
	}
	if v := ctx.Value("mfa"); v != nil {
		m["mfa"] = v
	}
	return m
}
