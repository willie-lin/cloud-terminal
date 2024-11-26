// Code generated by ent, DO NOT EDIT.

package usergroup

import (
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the usergroup type in the database.
	Label = "user_group"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldGroupName holds the string denoting the group_name field in the database.
	FieldGroupName = "group_name"
	// Table holds the table name of the usergroup in the database.
	Table = "user_groups"
)

// Columns holds all SQL columns for usergroup fields.
var Columns = []string{
	FieldID,
	FieldGroupName,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the UserGroup queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByGroupName orders the results by the group_name field.
func ByGroupName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGroupName, opts...).ToFunc()
}
