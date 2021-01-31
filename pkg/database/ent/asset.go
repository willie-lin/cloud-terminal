// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebook/ent/dialect/sql"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent/asset"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent/session"
)

// Asset is the model entity for the Asset schema.
type Asset struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// IP holds the value of the "ip" field.
	IP string `json:"ip,omitempty"`
	// Protocol holds the value of the "protocol" field.
	Protocol string `json:"protocol,omitempty"`
	// Port holds the value of the "port" field.
	Port int `json:"port,omitempty"`
	// AccountType holds the value of the "account_type" field.
	AccountType string `json:"account_type,omitempty"`
	// Username holds the value of the "username" field.
	Username string `json:"username,omitempty"`
	// Password holds the value of the "password" field.
	Password string `json:"password,omitempty"`
	// CredentialID holds the value of the "credential_id" field.
	CredentialID string `json:"credential_id,omitempty"`
	// PrivateKey holds the value of the "private_key" field.
	PrivateKey string `json:"private_key,omitempty"`
	// Passphrase holds the value of the "passphrase" field.
	Passphrase string `json:"passphrase,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Active holds the value of the "active" field.
	Active bool `json:"active,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Tags holds the value of the "tags" field.
	Tags string `json:"tags,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the AssetQuery when eager-loading is set.
	Edges          AssetEdges `json:"edges"`
	session_assets *string
}

// AssetEdges holds the relations/edges for other nodes in the graph.
type AssetEdges struct {
	// Sessions holds the value of the sessions edge.
	Sessions *Session
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// SessionsOrErr returns the Sessions value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AssetEdges) SessionsOrErr() (*Session, error) {
	if e.loadedTypes[0] {
		if e.Sessions == nil {
			// The edge sessions was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: session.Label}
		}
		return e.Sessions, nil
	}
	return nil, &NotLoadedError{edge: "sessions"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Asset) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case asset.FieldActive:
			values[i] = &sql.NullBool{}
		case asset.FieldPort:
			values[i] = &sql.NullInt64{}
		case asset.FieldID, asset.FieldName, asset.FieldIP, asset.FieldProtocol, asset.FieldAccountType, asset.FieldUsername, asset.FieldPassword, asset.FieldCredentialID, asset.FieldPrivateKey, asset.FieldPassphrase, asset.FieldDescription, asset.FieldTags:
			values[i] = &sql.NullString{}
		case asset.FieldCreatedAt, asset.FieldUpdatedAt:
			values[i] = &sql.NullTime{}
		case asset.ForeignKeys[0]: // session_assets
			values[i] = &sql.NullString{}
		default:
			return nil, fmt.Errorf("unexpected column %q for type Asset", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Asset fields.
func (a *Asset) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case asset.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				a.ID = value.String
			}
		case asset.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				a.Name = value.String
			}
		case asset.FieldIP:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field ip", values[i])
			} else if value.Valid {
				a.IP = value.String
			}
		case asset.FieldProtocol:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field protocol", values[i])
			} else if value.Valid {
				a.Protocol = value.String
			}
		case asset.FieldPort:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field port", values[i])
			} else if value.Valid {
				a.Port = int(value.Int64)
			}
		case asset.FieldAccountType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field account_type", values[i])
			} else if value.Valid {
				a.AccountType = value.String
			}
		case asset.FieldUsername:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field username", values[i])
			} else if value.Valid {
				a.Username = value.String
			}
		case asset.FieldPassword:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field password", values[i])
			} else if value.Valid {
				a.Password = value.String
			}
		case asset.FieldCredentialID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field credential_id", values[i])
			} else if value.Valid {
				a.CredentialID = value.String
			}
		case asset.FieldPrivateKey:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field private_key", values[i])
			} else if value.Valid {
				a.PrivateKey = value.String
			}
		case asset.FieldPassphrase:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field passphrase", values[i])
			} else if value.Valid {
				a.Passphrase = value.String
			}
		case asset.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				a.Description = value.String
			}
		case asset.FieldActive:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field active", values[i])
			} else if value.Valid {
				a.Active = value.Bool
			}
		case asset.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				a.CreatedAt = value.Time
			}
		case asset.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				a.UpdatedAt = value.Time
			}
		case asset.FieldTags:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field tags", values[i])
			} else if value.Valid {
				a.Tags = value.String
			}
		case asset.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field session_assets", values[i])
			} else if value.Valid {
				a.session_assets = new(string)
				*a.session_assets = value.String
			}
		}
	}
	return nil
}

// QuerySessions queries the "sessions" edge of the Asset entity.
func (a *Asset) QuerySessions() *SessionQuery {
	return (&AssetClient{config: a.config}).QuerySessions(a)
}

// Update returns a builder for updating this Asset.
// Note that you need to call Asset.Unwrap() before calling this method if this Asset
// was returned from a transaction, and the transaction was committed or rolled back.
func (a *Asset) Update() *AssetUpdateOne {
	return (&AssetClient{config: a.config}).UpdateOne(a)
}

// Unwrap unwraps the Asset entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (a *Asset) Unwrap() *Asset {
	tx, ok := a.config.driver.(*txDriver)
	if !ok {
		panic("ent: Asset is not a transactional entity")
	}
	a.config.driver = tx.drv
	return a
}

// String implements the fmt.Stringer.
func (a *Asset) String() string {
	var builder strings.Builder
	builder.WriteString("Asset(")
	builder.WriteString(fmt.Sprintf("id=%v", a.ID))
	builder.WriteString(", name=")
	builder.WriteString(a.Name)
	builder.WriteString(", ip=")
	builder.WriteString(a.IP)
	builder.WriteString(", protocol=")
	builder.WriteString(a.Protocol)
	builder.WriteString(", port=")
	builder.WriteString(fmt.Sprintf("%v", a.Port))
	builder.WriteString(", account_type=")
	builder.WriteString(a.AccountType)
	builder.WriteString(", username=")
	builder.WriteString(a.Username)
	builder.WriteString(", password=")
	builder.WriteString(a.Password)
	builder.WriteString(", credential_id=")
	builder.WriteString(a.CredentialID)
	builder.WriteString(", private_key=")
	builder.WriteString(a.PrivateKey)
	builder.WriteString(", passphrase=")
	builder.WriteString(a.Passphrase)
	builder.WriteString(", description=")
	builder.WriteString(a.Description)
	builder.WriteString(", active=")
	builder.WriteString(fmt.Sprintf("%v", a.Active))
	builder.WriteString(", created_at=")
	builder.WriteString(a.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(a.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", tags=")
	builder.WriteString(a.Tags)
	builder.WriteByte(')')
	return builder.String()
}

// Assets is a parsable slice of Asset.
type Assets []*Asset

func (a Assets) config(cfg config) {
	for _i := range a {
		a[_i].config = cfg
	}
}