package api

import (
	"fmt"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"regexp"
	"strings"
)

//type PolicyStatement struct {
//	Effect    string      `json:"Effect"`
//	Actions   []string    `json:"Action"`
//	Resources []string    `json:"Resource"`
//	Condition interface{} `json:"Condition,omitempty"`
//}

// GenerateRRN creates a Resource Reference Name (RRN) for a given resource
func GenerateRRN(platform, tenant, account, resourceType, resourceName string) string {
	return fmt.Sprintf("rrn:%s:%s:%s:%s:%s", platform, tenant, account, resourceType, resourceName)
}

//// ParseRRN parses a Resource Reference Name (RRN) into its components
//func ParseRRN(rrn string) (platform, tenant, account, resourceType, resourceName string, err error) {
//	parts := strings.Split(rrn, ":")
//	if len(parts) != 6 || parts[0] != "rrn" {
//		return "", "", "", "", "", fmt.Errorf("invalid RRN format")
//	}
//	return parts[1], parts[2], parts[3], parts[4], parts[5], nil
//}
//
//// CheckPermission checks if a given AccessPolicy has permission for a specific RRN
//func CheckPermission(policy *ent.AccessPolicy, rrn string) bool {
//	for _, statement := range policy.Statements {
//		if statement.Effect == "Allow" { // 只检查 "Allow" 语句
//			for _, resource := range statement.Resources {
//				if MatchRRN(resource, rrn) {
//					return true
//				}
//			}
//		}
//	}
//	return false
//}

// MatchRRN checks if a permission RRN matches a resource RRN
//func MatchRRN(permRRN, resourceRRN string) bool {
//	permParts := strings.Split(permRRN, ":")
//	resourceParts := strings.Split(resourceRRN, ":")
//
//	if len(permParts) != len(resourceParts) {
//		return false
//	}
//
//	for i := range permParts {
//		if permParts[i] != "*" && permParts[i] != resourceParts[i] {
//			return false
//		}
//	}
//
//	return true
//}

// ParseRRN parses a Resource Reference Name (RRN) into its components
func ParseRRN(rrn string) (platform, tenant, account, resourceType, resourceName string, err error) {
	parts := strings.Split(rrn, ":")
	if len(parts) != 6 || parts[0] != "rrn" {
		return "", "", "", "", "", fmt.Errorf("invalid RRN format: %s", rrn) // 更详细的错误信息
	}
	return parts[1], parts[2], parts[3], parts[4], parts[5], nil
}

// CheckPermission checks if a given AccessPolicy has permission for a specific RRN
func CheckPermission(policy *ent.AccessPolicy, rrn string) bool {
	for _, statement := range policy.Statements {
		if statement.Effect == "Allow" {
			for _, resource := range statement.Resources {
				if MatchRRN(resource, rrn) {
					return true
				}
			}
		}
	}
	return false
}

// MatchRRN checks if a permission RRN matches a resource RRN (using regexp)
func MatchRRN(permRRN, resourceRRN string) bool {
	// 将 * 转换为 .* 以支持正则表达式匹配
	permRegex := "^" + strings.ReplaceAll(permRRN, "*", ".*") + "$"
	matched, _ := regexp.MatchString(permRegex, resourceRRN) // 忽略错误，如果正则表达式无效则匹配失败
	return matched
}

// 示例用法
//func main() {
//	policy := &ent.AccessPolicy{
//		Statements: []PolicyStatement{ // 使用切片
//			{
//				Effect: "Allow",
//				Actions: []string{"read"},
//				Resources: []string{
//					"rrn:myplatform:tenant1:account1:database/*",
//					"rrn:myplatform:tenant2:*:api/v1/*",
//				},
//			},
//	}
//	rrn1 := "rrn:myplatform:tenant1:account1:database/users"
//	rrn2 := "rrn:myplatform:tenant2:account3:api/v1/get-user-data"
//	rrn3 := "rrn:myplatform:tenant3:account4:api/v2/get-user-data"
//
//	fmt.Println(CheckPermission(policy, rrn1)) // true
//	fmt.Println(CheckPermission(policy, rrn2)) // true
//	fmt.Println(CheckPermission(policy, rrn3)) // false
//}
