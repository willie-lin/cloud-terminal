// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/willie-lin/cloud-terminal/app/database/ent/auditlog"
)

// AuditLog is the model entity for the AuditLog schema.
type AuditLog struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Timestamp holds the value of the "timestamp" field.
	Timestamp time.Time `json:"timestamp,omitempty"`
	// ActorID holds the value of the "actor_id" field.
	ActorID int `json:"actor_id,omitempty"`
	// ActorUsername holds the value of the "actor_username" field.
	ActorUsername string `json:"actor_username,omitempty"`
	// Action holds the value of the "action" field.
	Action string `json:"action,omitempty"`
	// ResourceID holds the value of the "resource_id" field.
	ResourceID int `json:"resource_id,omitempty"`
	// ResourceType holds the value of the "resource_type" field.
	ResourceType string `json:"resource_type,omitempty"`
	// Details holds the value of the "details" field.
	Details map[string]interface{} `json:"details,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the AuditLogQuery when eager-loading is set.
	Edges        AuditLogEdges `json:"edges"`
	selectValues sql.SelectValues
}

// AuditLogEdges holds the relations/edges for other nodes in the graph.
type AuditLogEdges struct {
	// User holds the value of the user edge.
	User []*User `json:"user,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading.
func (e AuditLogEdges) UserOrErr() ([]*User, error) {
	if e.loadedTypes[0] {
		return e.User, nil
	}
	return nil, &NotLoadedError{edge: "user"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*AuditLog) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case auditlog.FieldDetails:
			values[i] = new([]byte)
		case auditlog.FieldActorID, auditlog.FieldResourceID:
			values[i] = new(sql.NullInt64)
		case auditlog.FieldActorUsername, auditlog.FieldAction, auditlog.FieldResourceType:
			values[i] = new(sql.NullString)
		case auditlog.FieldCreatedAt, auditlog.FieldUpdatedAt, auditlog.FieldTimestamp:
			values[i] = new(sql.NullTime)
		case auditlog.FieldID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the AuditLog fields.
func (al *AuditLog) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case auditlog.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				al.ID = *value
			}
		case auditlog.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				al.CreatedAt = value.Time
			}
		case auditlog.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				al.UpdatedAt = value.Time
			}
		case auditlog.FieldTimestamp:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field timestamp", values[i])
			} else if value.Valid {
				al.Timestamp = value.Time
			}
		case auditlog.FieldActorID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field actor_id", values[i])
			} else if value.Valid {
				al.ActorID = int(value.Int64)
			}
		case auditlog.FieldActorUsername:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field actor_username", values[i])
			} else if value.Valid {
				al.ActorUsername = value.String
			}
		case auditlog.FieldAction:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field action", values[i])
			} else if value.Valid {
				al.Action = value.String
			}
		case auditlog.FieldResourceID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field resource_id", values[i])
			} else if value.Valid {
				al.ResourceID = int(value.Int64)
			}
		case auditlog.FieldResourceType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field resource_type", values[i])
			} else if value.Valid {
				al.ResourceType = value.String
			}
		case auditlog.FieldDetails:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field details", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &al.Details); err != nil {
					return fmt.Errorf("unmarshal field details: %w", err)
				}
			}
		default:
			al.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the AuditLog.
// This includes values selected through modifiers, order, etc.
func (al *AuditLog) Value(name string) (ent.Value, error) {
	return al.selectValues.Get(name)
}

// QueryUser queries the "user" edge of the AuditLog entity.
func (al *AuditLog) QueryUser() *UserQuery {
	return NewAuditLogClient(al.config).QueryUser(al)
}

// Update returns a builder for updating this AuditLog.
// Note that you need to call AuditLog.Unwrap() before calling this method if this AuditLog
// was returned from a transaction, and the transaction was committed or rolled back.
func (al *AuditLog) Update() *AuditLogUpdateOne {
	return NewAuditLogClient(al.config).UpdateOne(al)
}

// Unwrap unwraps the AuditLog entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (al *AuditLog) Unwrap() *AuditLog {
	_tx, ok := al.config.driver.(*txDriver)
	if !ok {
		panic("ent: AuditLog is not a transactional entity")
	}
	al.config.driver = _tx.drv
	return al
}

// String implements the fmt.Stringer.
func (al *AuditLog) String() string {
	var builder strings.Builder
	builder.WriteString("AuditLog(")
	builder.WriteString(fmt.Sprintf("id=%v, ", al.ID))
	builder.WriteString("created_at=")
	builder.WriteString(al.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(al.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("timestamp=")
	builder.WriteString(al.Timestamp.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("actor_id=")
	builder.WriteString(fmt.Sprintf("%v", al.ActorID))
	builder.WriteString(", ")
	builder.WriteString("actor_username=")
	builder.WriteString(al.ActorUsername)
	builder.WriteString(", ")
	builder.WriteString("action=")
	builder.WriteString(al.Action)
	builder.WriteString(", ")
	builder.WriteString("resource_id=")
	builder.WriteString(fmt.Sprintf("%v", al.ResourceID))
	builder.WriteString(", ")
	builder.WriteString("resource_type=")
	builder.WriteString(al.ResourceType)
	builder.WriteString(", ")
	builder.WriteString("details=")
	builder.WriteString(fmt.Sprintf("%v", al.Details))
	builder.WriteByte(')')
	return builder.String()
}

// AuditLogs is a parsable slice of AuditLog.
type AuditLogs []*AuditLog