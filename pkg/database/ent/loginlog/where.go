// Code generated by entc, DO NOT EDIT.

package loginlog

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// UserID applies equality check predicate on the "user_id" field. It's identical to UserIDEQ.
func UserID(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUserID), v))
	})
}

// ClientIP applies equality check predicate on the "client_ip" field. It's identical to ClientIPEQ.
func ClientIP(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldClientIP), v))
	})
}

// ClentUsetAgent applies equality check predicate on the "clent_uset_agent" field. It's identical to ClentUsetAgentEQ.
func ClentUsetAgent(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldClentUsetAgent), v))
	})
}

// LoginTime applies equality check predicate on the "login_time" field. It's identical to LoginTimeEQ.
func LoginTime(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLoginTime), v))
	})
}

// LogoutTime applies equality check predicate on the "logout_time" field. It's identical to LogoutTimeEQ.
func LogoutTime(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLogoutTime), v))
	})
}

// Remember applies equality check predicate on the "remember" field. It's identical to RememberEQ.
func Remember(v bool) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRemember), v))
	})
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.LoginLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.LoginLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.LoginLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.LoginLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdatedAt), v))
	})
}

// UserIDEQ applies the EQ predicate on the "user_id" field.
func UserIDEQ(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUserID), v))
	})
}

// UserIDNEQ applies the NEQ predicate on the "user_id" field.
func UserIDNEQ(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUserID), v))
	})
}

// UserIDIn applies the In predicate on the "user_id" field.
func UserIDIn(vs ...string) predicate.LoginLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUserID), v...))
	})
}

// UserIDNotIn applies the NotIn predicate on the "user_id" field.
func UserIDNotIn(vs ...string) predicate.LoginLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUserID), v...))
	})
}

// UserIDGT applies the GT predicate on the "user_id" field.
func UserIDGT(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUserID), v))
	})
}

// UserIDGTE applies the GTE predicate on the "user_id" field.
func UserIDGTE(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUserID), v))
	})
}

// UserIDLT applies the LT predicate on the "user_id" field.
func UserIDLT(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUserID), v))
	})
}

// UserIDLTE applies the LTE predicate on the "user_id" field.
func UserIDLTE(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUserID), v))
	})
}

// UserIDContains applies the Contains predicate on the "user_id" field.
func UserIDContains(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldUserID), v))
	})
}

// UserIDHasPrefix applies the HasPrefix predicate on the "user_id" field.
func UserIDHasPrefix(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldUserID), v))
	})
}

// UserIDHasSuffix applies the HasSuffix predicate on the "user_id" field.
func UserIDHasSuffix(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldUserID), v))
	})
}

// UserIDEqualFold applies the EqualFold predicate on the "user_id" field.
func UserIDEqualFold(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldUserID), v))
	})
}

// UserIDContainsFold applies the ContainsFold predicate on the "user_id" field.
func UserIDContainsFold(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldUserID), v))
	})
}

// ClientIPEQ applies the EQ predicate on the "client_ip" field.
func ClientIPEQ(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldClientIP), v))
	})
}

// ClientIPNEQ applies the NEQ predicate on the "client_ip" field.
func ClientIPNEQ(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldClientIP), v))
	})
}

// ClientIPIn applies the In predicate on the "client_ip" field.
func ClientIPIn(vs ...string) predicate.LoginLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldClientIP), v...))
	})
}

// ClientIPNotIn applies the NotIn predicate on the "client_ip" field.
func ClientIPNotIn(vs ...string) predicate.LoginLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldClientIP), v...))
	})
}

// ClientIPGT applies the GT predicate on the "client_ip" field.
func ClientIPGT(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldClientIP), v))
	})
}

// ClientIPGTE applies the GTE predicate on the "client_ip" field.
func ClientIPGTE(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldClientIP), v))
	})
}

// ClientIPLT applies the LT predicate on the "client_ip" field.
func ClientIPLT(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldClientIP), v))
	})
}

// ClientIPLTE applies the LTE predicate on the "client_ip" field.
func ClientIPLTE(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldClientIP), v))
	})
}

// ClientIPContains applies the Contains predicate on the "client_ip" field.
func ClientIPContains(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldClientIP), v))
	})
}

// ClientIPHasPrefix applies the HasPrefix predicate on the "client_ip" field.
func ClientIPHasPrefix(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldClientIP), v))
	})
}

// ClientIPHasSuffix applies the HasSuffix predicate on the "client_ip" field.
func ClientIPHasSuffix(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldClientIP), v))
	})
}

// ClientIPEqualFold applies the EqualFold predicate on the "client_ip" field.
func ClientIPEqualFold(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldClientIP), v))
	})
}

// ClientIPContainsFold applies the ContainsFold predicate on the "client_ip" field.
func ClientIPContainsFold(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldClientIP), v))
	})
}

// ClentUsetAgentEQ applies the EQ predicate on the "clent_uset_agent" field.
func ClentUsetAgentEQ(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldClentUsetAgent), v))
	})
}

// ClentUsetAgentNEQ applies the NEQ predicate on the "clent_uset_agent" field.
func ClentUsetAgentNEQ(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldClentUsetAgent), v))
	})
}

// ClentUsetAgentIn applies the In predicate on the "clent_uset_agent" field.
func ClentUsetAgentIn(vs ...string) predicate.LoginLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldClentUsetAgent), v...))
	})
}

// ClentUsetAgentNotIn applies the NotIn predicate on the "clent_uset_agent" field.
func ClentUsetAgentNotIn(vs ...string) predicate.LoginLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldClentUsetAgent), v...))
	})
}

// ClentUsetAgentGT applies the GT predicate on the "clent_uset_agent" field.
func ClentUsetAgentGT(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldClentUsetAgent), v))
	})
}

// ClentUsetAgentGTE applies the GTE predicate on the "clent_uset_agent" field.
func ClentUsetAgentGTE(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldClentUsetAgent), v))
	})
}

// ClentUsetAgentLT applies the LT predicate on the "clent_uset_agent" field.
func ClentUsetAgentLT(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldClentUsetAgent), v))
	})
}

// ClentUsetAgentLTE applies the LTE predicate on the "clent_uset_agent" field.
func ClentUsetAgentLTE(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldClentUsetAgent), v))
	})
}

// ClentUsetAgentContains applies the Contains predicate on the "clent_uset_agent" field.
func ClentUsetAgentContains(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldClentUsetAgent), v))
	})
}

// ClentUsetAgentHasPrefix applies the HasPrefix predicate on the "clent_uset_agent" field.
func ClentUsetAgentHasPrefix(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldClentUsetAgent), v))
	})
}

// ClentUsetAgentHasSuffix applies the HasSuffix predicate on the "clent_uset_agent" field.
func ClentUsetAgentHasSuffix(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldClentUsetAgent), v))
	})
}

// ClentUsetAgentEqualFold applies the EqualFold predicate on the "clent_uset_agent" field.
func ClentUsetAgentEqualFold(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldClentUsetAgent), v))
	})
}

// ClentUsetAgentContainsFold applies the ContainsFold predicate on the "clent_uset_agent" field.
func ClentUsetAgentContainsFold(v string) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldClentUsetAgent), v))
	})
}

// LoginTimeEQ applies the EQ predicate on the "login_time" field.
func LoginTimeEQ(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLoginTime), v))
	})
}

// LoginTimeNEQ applies the NEQ predicate on the "login_time" field.
func LoginTimeNEQ(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLoginTime), v))
	})
}

// LoginTimeIn applies the In predicate on the "login_time" field.
func LoginTimeIn(vs ...time.Time) predicate.LoginLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldLoginTime), v...))
	})
}

// LoginTimeNotIn applies the NotIn predicate on the "login_time" field.
func LoginTimeNotIn(vs ...time.Time) predicate.LoginLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldLoginTime), v...))
	})
}

// LoginTimeGT applies the GT predicate on the "login_time" field.
func LoginTimeGT(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLoginTime), v))
	})
}

// LoginTimeGTE applies the GTE predicate on the "login_time" field.
func LoginTimeGTE(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLoginTime), v))
	})
}

// LoginTimeLT applies the LT predicate on the "login_time" field.
func LoginTimeLT(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLoginTime), v))
	})
}

// LoginTimeLTE applies the LTE predicate on the "login_time" field.
func LoginTimeLTE(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLoginTime), v))
	})
}

// LogoutTimeEQ applies the EQ predicate on the "logout_time" field.
func LogoutTimeEQ(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLogoutTime), v))
	})
}

// LogoutTimeNEQ applies the NEQ predicate on the "logout_time" field.
func LogoutTimeNEQ(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLogoutTime), v))
	})
}

// LogoutTimeIn applies the In predicate on the "logout_time" field.
func LogoutTimeIn(vs ...time.Time) predicate.LoginLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldLogoutTime), v...))
	})
}

// LogoutTimeNotIn applies the NotIn predicate on the "logout_time" field.
func LogoutTimeNotIn(vs ...time.Time) predicate.LoginLog {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.LoginLog(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldLogoutTime), v...))
	})
}

// LogoutTimeGT applies the GT predicate on the "logout_time" field.
func LogoutTimeGT(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLogoutTime), v))
	})
}

// LogoutTimeGTE applies the GTE predicate on the "logout_time" field.
func LogoutTimeGTE(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLogoutTime), v))
	})
}

// LogoutTimeLT applies the LT predicate on the "logout_time" field.
func LogoutTimeLT(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLogoutTime), v))
	})
}

// LogoutTimeLTE applies the LTE predicate on the "logout_time" field.
func LogoutTimeLTE(v time.Time) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLogoutTime), v))
	})
}

// RememberEQ applies the EQ predicate on the "remember" field.
func RememberEQ(v bool) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRemember), v))
	})
}

// RememberNEQ applies the NEQ predicate on the "remember" field.
func RememberNEQ(v bool) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldRemember), v))
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.LoginLog) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.LoginLog) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.LoginLog) predicate.LoginLog {
	return predicate.LoginLog(func(s *sql.Selector) {
		p(s.Not())
	})
}