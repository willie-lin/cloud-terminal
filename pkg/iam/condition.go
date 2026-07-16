package iam

import (
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MatchConditions 检查请求是否满足所有 Conditions
// 同一个 Condition 内的所有条件 AND 运算
func MatchConditions(cond *Condition, req *Request) bool {
	if cond == nil {
		return true
	}
	return matchIP(cond, req) &&
		matchString(cond, req) &&
		matchArn(cond, req) &&
		matchDate(cond, req) &&
		matchBool(cond, req) &&
		matchNumeric(cond, req) &&
		matchNull(cond, req) &&
		matchMFA(cond, req) &&
		matchSet(cond, req) &&
		matchTags(cond, req)
}

// ─── IP ────────────────────────────────────────────────────────

func matchIP(cond *Condition, req *Request) bool {
	for key, cidrs := range cond.IpAddress {
		val := ctxString(req, key)
		if val == "" || !matchCIDR(val, cidrs) {
			return false
		}
	}
	for key, cidrs := range cond.NotIpAddress {
		val := ctxString(req, key)
		if val != "" && matchCIDR(val, cidrs) {
			return false
		}
	}
	return true
}

func matchCIDR(ip string, cidrs []string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	for _, c := range cidrs {
		if strings.Contains(c, "/") {
			_, n, err := net.ParseCIDR(c)
			if err == nil && n.Contains(parsed) {
				return true
			}
		} else if c == ip {
			return true
		}
	}
	return false
}

// ─── String ────────────────────────────────────────────────────

func matchString(cond *Condition, req *Request) bool {
	for k, v := range cond.StringEquals {
		if ctxString(req, k) != v {
			return false
		}
	}
	for k, v := range cond.StringNotEquals {
		if ctxString(req, k) == v {
			return false
		}
	}
	for k, v := range cond.StringEqualsIgnoreCase {
		if !strings.EqualFold(ctxString(req, k), v) {
			return false
		}
	}
	for k, p := range cond.StringLike {
		if !globMatch(ctxString(req, k), p) {
			return false
		}
	}
	for k, p := range cond.StringNotLike {
		if globMatch(ctxString(req, k), p) {
			return false
		}
	}
	return true
}

// globMatch * 和 ? 通配匹配
func globMatch(value, pattern string) bool {
	if pattern == "*" {
		return true
	}
	re := "^"
	for _, ch := range pattern {
		switch ch {
		case '*':
			re += ".*"
		case '?':
			re += "."
		case '.', '+', '(', ')', '|', '[', ']', '{', '}', '^', '$', '\\':
			re += "\\" + string(ch)
		default:
			re += string(ch)
		}
	}
	re += "$"
	matched, _ := regexp.MatchString(re, value)
	return matched
}

// ─── ARN/URN ───────────────────────────────────────────────────

func matchArn(cond *Condition, req *Request) bool {
	for k, p := range cond.ArnEquals {
		if !urnMatch(ctxString(req, k), p) {
			return false
		}
	}
	for k, p := range cond.ArnLike {
		if !urnMatch(ctxString(req, k), p) {
			return false
		}
	}
	for k, p := range cond.ArnNotEquals {
		if urnMatch(ctxString(req, k), p) {
			return false
		}
	}
	for k, p := range cond.ArnNotLike {
		if urnMatch(ctxString(req, k), p) {
			return false
		}
	}
	return true
}

// urnMatch URN 通配匹配
func urnMatch(value, pattern string) bool {
	if pattern == "*" || pattern == value {
		return true
	}
	p := strings.SplitN(pattern, ":", 6)
	v := strings.SplitN(value, ":", 6)
	if len(p) != 6 || len(v) != 6 {
		// 不足 6 段时，如果 pattern 是 "*" 则匹配
		if pattern == "*" {
			return true
		}
		return false
	}
	// 所有 6 段都支持 glob 通配
	// 这样 "k8s-*" 能匹配 "k8s-service"
	// 最后一段（name）还支持含 "/" 的路径匹配
	for i := 0; i < 6; i++ {
		if p[i] == "*" {
			continue
		}
		if !globMatch(v[i], p[i]) {
			return false
		}
	}
	return true
}

// ─── Date ──────────────────────────────────────────────────────

func matchDate(cond *Condition, req *Request) bool {
	now := time.Now()

	check := func(key, dateStr string, after, equal bool) bool {
		if dateStr == "" {
			return true
		}
		val := ctxTime(req, key, now)
		t, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return false
		}
		if after && equal {
			return !val.Before(t) // >=
		}
		if after {
			return val.After(t)
		}
		if equal {
			return val.Equal(t)
		}
		return val.Before(t)
	}

	for k, v := range cond.DateGreaterThan {
		if !check(k, v, true, false) {
			return false
		}
	}
	for k, v := range cond.DateLessThan {
		if !check(k, v, false, false) {
			return false
		}
	}
	for k, v := range cond.DateEquals {
		if !check(k, v, false, true) {
			return false
		}
	}
	for k, v := range cond.DateNotEquals {
		if check(k, v, false, true) {
			return false
		}
	}
	if cond.DateBetween != nil {
		if cond.DateBetween.After != "" {
			if !check(cond.DateBetween.Key, cond.DateBetween.After, true, false) {
				return false
			}
		}
		if cond.DateBetween.Before != "" {
			if !check(cond.DateBetween.Key, cond.DateBetween.Before, false, false) {
				return false
			}
		}
	}
	return true
}

// ─── Bool ──────────────────────────────────────────────────────

func matchBool(cond *Condition, req *Request) bool {
	for k, expect := range cond.Bool {
		if ctxBool(req, k) != expect {
			return false
		}
	}
	return true
}

// ─── Numeric ───────────────────────────────────────────────────

func matchNumeric(cond *Condition, req *Request) bool {
	for k, v := range cond.NumericEquals {
		if ctxFloat(req, k) != v {
			return false
		}
	}
	for k, v := range cond.NumericNotEquals {
		if ctxFloat(req, k) == v {
			return false
		}
	}
	for k, v := range cond.NumericLessThan {
		if ctxFloat(req, k) >= v {
			return false
		}
	}
	for k, v := range cond.NumericGreaterThan {
		if ctxFloat(req, k) <= v {
			return false
		}
	}
	return true
}

// ─── Null ──────────────────────────────────────────────────────

func matchNull(cond *Condition, req *Request) bool {
	// AWS IAM 语义：
	//   Null: {"key": true}  → key 必须不存在（is null）
	//   Null: {"key": false} → key 必须存在（is not null）
	for k, shouldBeNull := range cond.Null {
		_, exists := req.Context[k]
		if shouldBeNull && exists {
			return false // key 存在但要求为空
		}
		if !shouldBeNull && !exists {
			return false // key 不存在但要求非空
		}
	}
	return true
}

// ─── MFA ───────────────────────────────────────────────────────

func matchMFA(cond *Condition, req *Request) bool {
	if cond.RequireMFA == nil {
		return true
	}
	if !*cond.RequireMFA {
		return true
	}
	return ctxBool(req, CtxMFAAuth)
}

// ─── Set 运算符 ────────────────────────────────────────────────

func matchSet(cond *Condition, req *Request) bool {
	if cond.ForAllValues != nil {
		reqVals := ctxStringSlice(req, cond.ForAllValues.Key)
		for _, rv := range reqVals {
			found := false
			for _, pv := range cond.ForAllValues.Values {
				if rv == pv {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}
	if cond.ForAnyValue != nil {
		reqVals := ctxStringSlice(req, cond.ForAnyValue.Key)
		for _, rv := range reqVals {
			for _, pv := range cond.ForAnyValue.Values {
				if rv == pv {
					return true
				}
			}
		}
		return false
	}
	return true
}

// ─── ABAC Tags ─────────────────────────────────────────────────

func matchTags(cond *Condition, req *Request) bool {
	for k, v := range cond.PrincipalTag {
		if ctxString(req, CtxPrincipalTag+"/"+k) != v {
			return false
		}
	}
	for k, v := range cond.ResourceTag {
		if ctxString(req, CtxResourceTag+"/"+k) != v {
			return false
		}
	}
	for k, v := range cond.RequestTag {
		if ctxString(req, CtxRequestTag+"/"+k) != v {
			return false
		}
	}
	return true
}

// ─── Context 辅助 ──────────────────────────────────────────────

func ctxString(req *Request, key string) string {
	if req.Context == nil {
		return ""
	}
	v, _ := req.Context[key].(string)
	return v
}

func ctxBool(req *Request, key string) bool {
	if req.Context == nil {
		return false
	}
	v, _ := req.Context[key].(bool)
	return v
}

func ctxFloat(req *Request, key string) float64 {
	if req.Context == nil {
		return 0
	}
	switch v := req.Context[key].(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	}
	return 0
}

func ctxTime(req *Request, key string, fallback time.Time) time.Time {
	if req.Context == nil {
		return fallback
	}
	switch v := req.Context[key].(type) {
	case time.Time:
		return v
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			return t
		}
	}
	return fallback
}

func ctxStringSlice(req *Request, key string) []string {
	if req.Context == nil {
		return nil
	}
	switch v := req.Context[key].(type) {
	case []string:
		return v
	case []interface{}:
		out := make([]string, len(v))
		for i, x := range v {
			out[i], _ = x.(string)
		}
		return out
	}
	return nil
}
