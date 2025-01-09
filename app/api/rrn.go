package api

import (
	"fmt"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"regexp"
	"strings"
	"sync"
)

var (
	regexCache   = make(map[string]*regexp.Regexp)
	regexCacheMu sync.RWMutex
)

func generateRRN(platform, tenant, account, resourceType, resourceName string) string {
	return fmt.Sprintf("rrn:%s:%s:%s:%s:%s", platform, tenant, account, resourceType, resourceName)
}

func ParseRRN(rrn string) (platform, tenant, account, resourceType, resourceName string, err error) {
	parts := strings.Split(rrn, ":")
	if len(parts) != 5 || parts[0] != "rrn" { //修改为5部分
		return "", "", "", "", "", fmt.Errorf("invalid RRN format: %s (expected rrn:platform:tenant:account:resourceType:resourceName)", rrn)
	}
	return parts[1], parts[2], parts[3], parts[4], parts[5], nil
}

func MatchRRN(permRRN, resourceRRN string) bool {
	regexCacheMu.RLock()
	re, ok := regexCache[permRRN]
	regexCacheMu.RUnlock()

	if !ok {
		permRegex := "^" + strings.ReplaceAll(permRRN, "*", ".*") + "$"
		var err error
		re, err = regexp.Compile(permRegex)
		if err != nil {
			fmt.Println(fmt.Errorf("invalid regex: %w", err))
			return false // 正则表达式编译失败，匹配失败
		}

		regexCacheMu.Lock()
		regexCache[permRRN] = re
		regexCacheMu.Unlock()
	}

	return re.MatchString(resourceRRN)
}

func CheckPermission(policy *ent.AccessPolicy, rrn string) bool {
	for _, statement := range policy.Statements {
		if statement.Effect == "Deny" { // Deny 优先
			for _, resource := range statement.Resources {
				if MatchRRN(resource, rrn) {
					return false // 显式拒绝
				}
			}
		}
	}
	for _, statement := range policy.Statements { // 再检查 Allow
		if statement.Effect == "Allow" {
			for _, resource := range statement.Resources {
				if MatchRRN(resource, rrn) {
					return true // 允许
				}
			}
		}
	}
	return false // 默认拒绝
}
