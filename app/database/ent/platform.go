// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/willie-lin/cloud-terminal/app/database/ent/platform"
)

// Platform is the model entity for the Platform schema.
type Platform struct {
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
	// Region holds the value of the "region" field.
	Region string `json:"region,omitempty"`
	// Version holds the value of the "version" field.
	Version string `json:"version,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PlatformQuery when eager-loading is set.
	Edges        PlatformEdges `json:"edges"`
	selectValues sql.SelectValues
}

// PlatformEdges holds the relations/edges for other nodes in the graph.
type PlatformEdges struct {
	// Tenants holds the value of the tenants edge.
	Tenants []*Tenant `json:"tenants,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// TenantsOrErr returns the Tenants value or an error if the edge
// was not loaded in eager-loading.
func (e PlatformEdges) TenantsOrErr() ([]*Tenant, error) {
	if e.loadedTypes[0] {
		return e.Tenants, nil
	}
	return nil, &NotLoadedError{edge: "tenants"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Platform) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case platform.FieldName, platform.FieldDescription, platform.FieldRegion, platform.FieldVersion:
			values[i] = new(sql.NullString)
		case platform.FieldCreatedAt, platform.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case platform.FieldID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Platform fields.
func (pl *Platform) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case platform.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				pl.ID = *value
			}
		case platform.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				pl.CreatedAt = value.Time
			}
		case platform.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				pl.UpdatedAt = value.Time
			}
		case platform.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				pl.Name = value.String
			}
		case platform.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				pl.Description = value.String
			}
		case platform.FieldRegion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field region", values[i])
			} else if value.Valid {
				pl.Region = value.String
			}
		case platform.FieldVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field version", values[i])
			} else if value.Valid {
				pl.Version = value.String
			}
		default:
			pl.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Platform.
// This includes values selected through modifiers, order, etc.
func (pl *Platform) Value(name string) (ent.Value, error) {
	return pl.selectValues.Get(name)
}

// QueryTenants queries the "tenants" edge of the Platform entity.
func (pl *Platform) QueryTenants() *TenantQuery {
	return NewPlatformClient(pl.config).QueryTenants(pl)
}

// Update returns a builder for updating this Platform.
// Note that you need to call Platform.Unwrap() before calling this method if this Platform
// was returned from a transaction, and the transaction was committed or rolled back.
func (pl *Platform) Update() *PlatformUpdateOne {
	return NewPlatformClient(pl.config).UpdateOne(pl)
}

// Unwrap unwraps the Platform entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pl *Platform) Unwrap() *Platform {
	_tx, ok := pl.config.driver.(*txDriver)
	if !ok {
		panic("ent: Platform is not a transactional entity")
	}
	pl.config.driver = _tx.drv
	return pl
}

// String implements the fmt.Stringer.
func (pl *Platform) String() string {
	var builder strings.Builder
	builder.WriteString("Platform(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pl.ID))
	builder.WriteString("created_at=")
	builder.WriteString(pl.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(pl.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(pl.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(pl.Description)
	builder.WriteString(", ")
	builder.WriteString("region=")
	builder.WriteString(pl.Region)
	builder.WriteString(", ")
	builder.WriteString("version=")
	builder.WriteString(pl.Version)
	builder.WriteByte(')')
	return builder.String()
}

// Platforms is a parsable slice of Platform.
type Platforms []*Platform