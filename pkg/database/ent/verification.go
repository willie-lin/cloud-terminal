// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"github.com/facebook/ent/dialect/sql"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent/verification"
)

// Verification is the model entity for the Verification schema.
type Verification struct {
	config
	// ID of the ent.
	ID int `json:"id,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Verification) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case verification.FieldID:
			values[i] = &sql.NullInt64{}
		default:
			return nil, fmt.Errorf("unexpected column %q for type Verification", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Verification fields.
func (v *Verification) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case verification.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			v.ID = int(value.Int64)
		}
	}
	return nil
}

// Update returns a builder for updating this Verification.
// Note that you need to call Verification.Unwrap() before calling this method if this Verification
// was returned from a transaction, and the transaction was committed or rolled back.
func (v *Verification) Update() *VerificationUpdateOne {
	return (&VerificationClient{config: v.config}).UpdateOne(v)
}

// Unwrap unwraps the Verification entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (v *Verification) Unwrap() *Verification {
	tx, ok := v.config.driver.(*txDriver)
	if !ok {
		panic("ent: Verification is not a transactional entity")
	}
	v.config.driver = tx.drv
	return v
}

// String implements the fmt.Stringer.
func (v *Verification) String() string {
	var builder strings.Builder
	builder.WriteString("Verification(")
	builder.WriteString(fmt.Sprintf("id=%v", v.ID))
	builder.WriteByte(')')
	return builder.String()
}

// Verifications is a parsable slice of Verification.
type Verifications []*Verification

func (v Verifications) config(cfg config) {
	for _i := range v {
		v[_i].config = cfg
	}
}