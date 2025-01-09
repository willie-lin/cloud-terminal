package utils

const (
	PlatformName     = "CloudSecPlatform"
	ManagementTenant = "systemTenant"
	Account          = "systemAccount"
)

const (
	UserResourceType     = "user"
	ConfigResourceType   = "config"
	TenantResourceType   = "tenant"
	AccountResourceType  = "account"
	ProjectResourceType  = "project"
	AuditLogResourceType = "audit-log"
	RoleResourceType     = "role"
	PolicyResourceType   = "policy"
)

const (
	// ActionConfigRead ActionConfig
	ActionConfigRead           = "config:read"
	ActionConfigUpdateDatabase = "config:update:database"
	ActionConfigUpdateNetwork  = "config:update:network"

	// ActionTenantCreate Tenant Actions
	ActionTenantCreate = "tenant:create"
	ActionTenantRead   = "tenant:read"
	ActionTenantUpdate = "tenant:update"
	ActionTenantDelete = "tenant:delete"

	// ActionUserCreate User Actions
	ActionUserCreate = "user:create"
	ActionUserRead   = "user:read"
	ActionUserUpdate = "user:update"
	ActionUserDelete = "user:delete"

	// ActionRoleCreate ActionUserCreate User Actions
	ActionRoleCreate = "role:create"
	ActionRoleRead   = "role:read"
	ActionRoleUpdate = "role:update"
	ActionRoleDelete = "role:delete"

	// ActionProjectCreate Project Actions
	ActionProjectCreate = "project:create"
	ActionProjectRead   = "project:read"
	ActionProjectUpdate = "project:update"
	ActionProjectDelete = "project:delete"

	// ActionAuditLogRead Audit Log Actions
	ActionAuditLogRead   = "audit-log:read"
	ActionAuditLogExport = "audit-log:export"
)

const (
	ResourceUserAll        = "rrn:" + PlatformName + ":*:*:user:*"    // 匹配所有租户和账户的用户
	ResourceUserSpecific   = "rrn:" + PlatformName + ":%s:%s:user:%s" // 匹配特定租户和账户的用户
	ResourceConfigAll      = "rrn:" + PlatformName + ":*:*:config:*"
	ResourceConfigSpecific = "rrn:" + PlatformName + ":%s:%s:config:%s"
	ResourceTenantAll      = "rrn:" + PlatformName + ":*:*:tenant:*"
	ResourceAccountAll     = "rrn:" + PlatformName + ":*:*:account:*"
	ResourceProjectAll     = "rrn:" + PlatformName + ":*:*:project:*"
	ResourceAuditLogAll    = "rrn:" + PlatformName + ":*:*:audit-log:*"
	ResourceRoleAll        = "rrn:" + PlatformName + ":*:*:role:*"
	ResourcePolicyAll      = "rrn:" + PlatformName + ":*:*:policy:*"
)

//const (
//	// ResourceConfigAll Resource RRN Patterns
//	ResourceConfigAll   = "rrn:CloudSecPlatform:management:system:config:*"
//	ResourceUserAll     = "rrn:CloudSecPlatform:management:system:user:*"
//	ResourceTenantAll   = "rrn:CloudSecPlatform:management:system:tenant:*"
//	ResourceProjectAll  = "rrn:CloudSecPlatform:management:system:project:*"
//	ResourceAuditLogAll = "rrn:CloudSecPlatform:management:system:audit-log:*"
//)
