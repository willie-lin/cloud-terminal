// Code generated by ent, DO NOT EDIT.

package role

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the role type in the database.
	Label = "role"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldIsDisabled holds the string denoting the is_disabled field in the database.
	FieldIsDisabled = "is_disabled"
	// FieldIsDefault holds the string denoting the is_default field in the database.
	FieldIsDefault = "is_default"
	// EdgeAccount holds the string denoting the account edge name in mutations.
	EdgeAccount = "account"
	// EdgeUsers holds the string denoting the users edge name in mutations.
	EdgeUsers = "users"
	// EdgeAccessPolicies holds the string denoting the access_policies edge name in mutations.
	EdgeAccessPolicies = "access_policies"
	// EdgeParentRole holds the string denoting the parent_role edge name in mutations.
	EdgeParentRole = "parent_role"
	// EdgeChildRoles holds the string denoting the child_roles edge name in mutations.
	EdgeChildRoles = "child_roles"
	// Table holds the table name of the role in the database.
	Table = "roles"
	// AccountTable is the table that holds the account relation/edge.
	AccountTable = "roles"
	// AccountInverseTable is the table name for the Account entity.
	// It exists in this package in order to avoid circular dependency with the "account" package.
	AccountInverseTable = "accounts"
	// AccountColumn is the table column denoting the account relation/edge.
	AccountColumn = "account_roles"
	// UsersTable is the table that holds the users relation/edge.
	UsersTable = "users"
	// UsersInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UsersInverseTable = "users"
	// UsersColumn is the table column denoting the users relation/edge.
	UsersColumn = "user_role"
	// AccessPoliciesTable is the table that holds the access_policies relation/edge. The primary key declared below.
	AccessPoliciesTable = "role_access_policies"
	// AccessPoliciesInverseTable is the table name for the AccessPolicy entity.
	// It exists in this package in order to avoid circular dependency with the "accesspolicy" package.
	AccessPoliciesInverseTable = "access_policies"
	// ParentRoleTable is the table that holds the parent_role relation/edge. The primary key declared below.
	ParentRoleTable = "role_child_roles"
	// ChildRolesTable is the table that holds the child_roles relation/edge. The primary key declared below.
	ChildRolesTable = "role_child_roles"
)

// Columns holds all SQL columns for role fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldName,
	FieldDescription,
	FieldIsDisabled,
	FieldIsDefault,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "roles"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"account_roles",
}

var (
	// AccessPoliciesPrimaryKey and AccessPoliciesColumn2 are the table columns denoting the
	// primary key for the access_policies relation (M2M).
	AccessPoliciesPrimaryKey = []string{"role_id", "access_policy_id"}
	// ParentRolePrimaryKey and ParentRoleColumn2 are the table columns denoting the
	// primary key for the parent_role relation (M2M).
	ParentRolePrimaryKey = []string{"role_id", "parent_role_id"}
	// ChildRolesPrimaryKey and ChildRolesColumn2 are the table columns denoting the
	// primary key for the child_roles relation (M2M).
	ChildRolesPrimaryKey = []string{"role_id", "parent_role_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultIsDisabled holds the default value on creation for the "is_disabled" field.
	DefaultIsDisabled bool
	// DefaultIsDefault holds the default value on creation for the "is_default" field.
	DefaultIsDefault bool
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Role queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByIsDisabled orders the results by the is_disabled field.
func ByIsDisabled(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsDisabled, opts...).ToFunc()
}

// ByIsDefault orders the results by the is_default field.
func ByIsDefault(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsDefault, opts...).ToFunc()
}

// ByAccountField orders the results by account field.
func ByAccountField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newAccountStep(), sql.OrderByField(field, opts...))
	}
}

// ByUsersCount orders the results by users count.
func ByUsersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newUsersStep(), opts...)
	}
}

// ByUsers orders the results by users terms.
func ByUsers(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newUsersStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByAccessPoliciesCount orders the results by access_policies count.
func ByAccessPoliciesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newAccessPoliciesStep(), opts...)
	}
}

// ByAccessPolicies orders the results by access_policies terms.
func ByAccessPolicies(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newAccessPoliciesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByParentRoleCount orders the results by parent_role count.
func ByParentRoleCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newParentRoleStep(), opts...)
	}
}

// ByParentRole orders the results by parent_role terms.
func ByParentRole(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newParentRoleStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByChildRolesCount orders the results by child_roles count.
func ByChildRolesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newChildRolesStep(), opts...)
	}
}

// ByChildRoles orders the results by child_roles terms.
func ByChildRoles(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newChildRolesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newAccountStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(AccountInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, AccountTable, AccountColumn),
	)
}
func newUsersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(UsersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, UsersTable, UsersColumn),
	)
}
func newAccessPoliciesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(AccessPoliciesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, AccessPoliciesTable, AccessPoliciesPrimaryKey...),
	)
}
func newParentRoleStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(Table, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, ParentRoleTable, ParentRolePrimaryKey...),
	)
}
func newChildRolesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(Table, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, ChildRolesTable, ChildRolesPrimaryKey...),
	)
}
