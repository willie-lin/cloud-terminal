// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
)

// Tenant is the model entity for the Tenant schema.
type Tenant struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the TenantQuery when eager-loading is set.
	Edges        TenantEdges `json:"edges"`
	selectValues sql.SelectValues
}

// TenantEdges holds the relations/edges for other nodes in the graph.
type TenantEdges struct {
	// Users holds the value of the users edge.
	Users []*User `json:"users,omitempty"`
	// Roles holds the value of the roles edge.
	Roles []*Role `json:"roles,omitempty"`
	// Resources holds the value of the resources edge.
	Resources []*Resource `json:"resources,omitempty"`
	// Permissions holds the value of the permissions edge.
	Permissions []*Permission `json:"permissions,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// UsersOrErr returns the Users value or an error if the edge
// was not loaded in eager-loading.
func (e TenantEdges) UsersOrErr() ([]*User, error) {
	if e.loadedTypes[0] {
		return e.Users, nil
	}
	return nil, &NotLoadedError{edge: "users"}
}

// RolesOrErr returns the Roles value or an error if the edge
// was not loaded in eager-loading.
func (e TenantEdges) RolesOrErr() ([]*Role, error) {
	if e.loadedTypes[1] {
		return e.Roles, nil
	}
	return nil, &NotLoadedError{edge: "roles"}
}

// ResourcesOrErr returns the Resources value or an error if the edge
// was not loaded in eager-loading.
func (e TenantEdges) ResourcesOrErr() ([]*Resource, error) {
	if e.loadedTypes[2] {
		return e.Resources, nil
	}
	return nil, &NotLoadedError{edge: "resources"}
}

// PermissionsOrErr returns the Permissions value or an error if the edge
// was not loaded in eager-loading.
func (e TenantEdges) PermissionsOrErr() ([]*Permission, error) {
	if e.loadedTypes[3] {
		return e.Permissions, nil
	}
	return nil, &NotLoadedError{edge: "permissions"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Tenant) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case tenant.FieldName, tenant.FieldDescription:
			values[i] = new(sql.NullString)
		case tenant.FieldCreatedAt, tenant.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case tenant.FieldID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Tenant fields.
func (t *Tenant) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case tenant.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				t.ID = *value
			}
		case tenant.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				t.CreatedAt = value.Time
			}
		case tenant.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				t.UpdatedAt = value.Time
			}
		case tenant.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				t.Name = value.String
			}
		case tenant.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				t.Description = value.String
			}
		default:
			t.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Tenant.
// This includes values selected through modifiers, order, etc.
func (t *Tenant) Value(name string) (ent.Value, error) {
	return t.selectValues.Get(name)
}

// QueryUsers queries the "users" edge of the Tenant entity.
func (t *Tenant) QueryUsers() *UserQuery {
	return NewTenantClient(t.config).QueryUsers(t)
}

// QueryRoles queries the "roles" edge of the Tenant entity.
func (t *Tenant) QueryRoles() *RoleQuery {
	return NewTenantClient(t.config).QueryRoles(t)
}

// QueryResources queries the "resources" edge of the Tenant entity.
func (t *Tenant) QueryResources() *ResourceQuery {
	return NewTenantClient(t.config).QueryResources(t)
}

// QueryPermissions queries the "permissions" edge of the Tenant entity.
func (t *Tenant) QueryPermissions() *PermissionQuery {
	return NewTenantClient(t.config).QueryPermissions(t)
}

// Update returns a builder for updating this Tenant.
// Note that you need to call Tenant.Unwrap() before calling this method if this Tenant
// was returned from a transaction, and the transaction was committed or rolled back.
func (t *Tenant) Update() *TenantUpdateOne {
	return NewTenantClient(t.config).UpdateOne(t)
}

// Unwrap unwraps the Tenant entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (t *Tenant) Unwrap() *Tenant {
	_tx, ok := t.config.driver.(*txDriver)
	if !ok {
		panic("ent: Tenant is not a transactional entity")
	}
	t.config.driver = _tx.drv
	return t
}

// String implements the fmt.Stringer.
func (t *Tenant) String() string {
	var builder strings.Builder
	builder.WriteString("Tenant(")
	builder.WriteString(fmt.Sprintf("id=%v, ", t.ID))
	builder.WriteString("created_at=")
	builder.WriteString(t.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(t.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(t.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(t.Description)
	builder.WriteByte(')')
	return builder.String()
}

// Tenants is a parsable slice of Tenant.
type Tenants []*Tenant
