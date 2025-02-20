// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/willie-lin/cloud-terminal/app/database/ent/account"
	"github.com/willie-lin/cloud-terminal/app/database/ent/predicate"
	"github.com/willie-lin/cloud-terminal/app/database/ent/resource"
)

// ResourceUpdate is the builder for updating Resource entities.
type ResourceUpdate struct {
	config
	hooks    []Hook
	mutation *ResourceMutation
}

// Where appends a list predicates to the ResourceUpdate builder.
func (ru *ResourceUpdate) Where(ps ...predicate.Resource) *ResourceUpdate {
	ru.mutation.Where(ps...)
	return ru
}

// SetUpdatedAt sets the "updated_at" field.
func (ru *ResourceUpdate) SetUpdatedAt(t time.Time) *ResourceUpdate {
	ru.mutation.SetUpdatedAt(t)
	return ru
}

// SetName sets the "name" field.
func (ru *ResourceUpdate) SetName(s string) *ResourceUpdate {
	ru.mutation.SetName(s)
	return ru
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ru *ResourceUpdate) SetNillableName(s *string) *ResourceUpdate {
	if s != nil {
		ru.SetName(*s)
	}
	return ru
}

// SetType sets the "type" field.
func (ru *ResourceUpdate) SetType(s string) *ResourceUpdate {
	ru.mutation.SetType(s)
	return ru
}

// SetNillableType sets the "type" field if the given value is not nil.
func (ru *ResourceUpdate) SetNillableType(s *string) *ResourceUpdate {
	if s != nil {
		ru.SetType(*s)
	}
	return ru
}

// SetRrn sets the "rrn" field.
func (ru *ResourceUpdate) SetRrn(s string) *ResourceUpdate {
	ru.mutation.SetRrn(s)
	return ru
}

// SetNillableRrn sets the "rrn" field if the given value is not nil.
func (ru *ResourceUpdate) SetNillableRrn(s *string) *ResourceUpdate {
	if s != nil {
		ru.SetRrn(*s)
	}
	return ru
}

// SetProperties sets the "properties" field.
func (ru *ResourceUpdate) SetProperties(m map[string]interface{}) *ResourceUpdate {
	ru.mutation.SetProperties(m)
	return ru
}

// ClearProperties clears the value of the "properties" field.
func (ru *ResourceUpdate) ClearProperties() *ResourceUpdate {
	ru.mutation.ClearProperties()
	return ru
}

// SetTags sets the "tags" field.
func (ru *ResourceUpdate) SetTags(m map[string]string) *ResourceUpdate {
	ru.mutation.SetTags(m)
	return ru
}

// ClearTags clears the value of the "tags" field.
func (ru *ResourceUpdate) ClearTags() *ResourceUpdate {
	ru.mutation.ClearTags()
	return ru
}

// SetDescription sets the "description" field.
func (ru *ResourceUpdate) SetDescription(s string) *ResourceUpdate {
	ru.mutation.SetDescription(s)
	return ru
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (ru *ResourceUpdate) SetNillableDescription(s *string) *ResourceUpdate {
	if s != nil {
		ru.SetDescription(*s)
	}
	return ru
}

// ClearDescription clears the value of the "description" field.
func (ru *ResourceUpdate) ClearDescription() *ResourceUpdate {
	ru.mutation.ClearDescription()
	return ru
}

// AddAccountIDs adds the "account" edge to the Account entity by IDs.
func (ru *ResourceUpdate) AddAccountIDs(ids ...uuid.UUID) *ResourceUpdate {
	ru.mutation.AddAccountIDs(ids...)
	return ru
}

// AddAccount adds the "account" edges to the Account entity.
func (ru *ResourceUpdate) AddAccount(a ...*Account) *ResourceUpdate {
	ids := make([]uuid.UUID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return ru.AddAccountIDs(ids...)
}

// AddChildIDs adds the "children" edge to the Resource entity by IDs.
func (ru *ResourceUpdate) AddChildIDs(ids ...uuid.UUID) *ResourceUpdate {
	ru.mutation.AddChildIDs(ids...)
	return ru
}

// AddChildren adds the "children" edges to the Resource entity.
func (ru *ResourceUpdate) AddChildren(r ...*Resource) *ResourceUpdate {
	ids := make([]uuid.UUID, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ru.AddChildIDs(ids...)
}

// AddParentIDs adds the "parent" edge to the Resource entity by IDs.
func (ru *ResourceUpdate) AddParentIDs(ids ...uuid.UUID) *ResourceUpdate {
	ru.mutation.AddParentIDs(ids...)
	return ru
}

// AddParent adds the "parent" edges to the Resource entity.
func (ru *ResourceUpdate) AddParent(r ...*Resource) *ResourceUpdate {
	ids := make([]uuid.UUID, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ru.AddParentIDs(ids...)
}

// Mutation returns the ResourceMutation object of the builder.
func (ru *ResourceUpdate) Mutation() *ResourceMutation {
	return ru.mutation
}

// ClearAccount clears all "account" edges to the Account entity.
func (ru *ResourceUpdate) ClearAccount() *ResourceUpdate {
	ru.mutation.ClearAccount()
	return ru
}

// RemoveAccountIDs removes the "account" edge to Account entities by IDs.
func (ru *ResourceUpdate) RemoveAccountIDs(ids ...uuid.UUID) *ResourceUpdate {
	ru.mutation.RemoveAccountIDs(ids...)
	return ru
}

// RemoveAccount removes "account" edges to Account entities.
func (ru *ResourceUpdate) RemoveAccount(a ...*Account) *ResourceUpdate {
	ids := make([]uuid.UUID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return ru.RemoveAccountIDs(ids...)
}

// ClearChildren clears all "children" edges to the Resource entity.
func (ru *ResourceUpdate) ClearChildren() *ResourceUpdate {
	ru.mutation.ClearChildren()
	return ru
}

// RemoveChildIDs removes the "children" edge to Resource entities by IDs.
func (ru *ResourceUpdate) RemoveChildIDs(ids ...uuid.UUID) *ResourceUpdate {
	ru.mutation.RemoveChildIDs(ids...)
	return ru
}

// RemoveChildren removes "children" edges to Resource entities.
func (ru *ResourceUpdate) RemoveChildren(r ...*Resource) *ResourceUpdate {
	ids := make([]uuid.UUID, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ru.RemoveChildIDs(ids...)
}

// ClearParent clears all "parent" edges to the Resource entity.
func (ru *ResourceUpdate) ClearParent() *ResourceUpdate {
	ru.mutation.ClearParent()
	return ru
}

// RemoveParentIDs removes the "parent" edge to Resource entities by IDs.
func (ru *ResourceUpdate) RemoveParentIDs(ids ...uuid.UUID) *ResourceUpdate {
	ru.mutation.RemoveParentIDs(ids...)
	return ru
}

// RemoveParent removes "parent" edges to Resource entities.
func (ru *ResourceUpdate) RemoveParent(r ...*Resource) *ResourceUpdate {
	ids := make([]uuid.UUID, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ru.RemoveParentIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ru *ResourceUpdate) Save(ctx context.Context) (int, error) {
	ru.defaults()
	return withHooks(ctx, ru.sqlSave, ru.mutation, ru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ru *ResourceUpdate) SaveX(ctx context.Context) int {
	affected, err := ru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ru *ResourceUpdate) Exec(ctx context.Context) error {
	_, err := ru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ru *ResourceUpdate) ExecX(ctx context.Context) {
	if err := ru.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ru *ResourceUpdate) defaults() {
	if _, ok := ru.mutation.UpdatedAt(); !ok {
		v := resource.UpdateDefaultUpdatedAt()
		ru.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ru *ResourceUpdate) check() error {
	if v, ok := ru.mutation.Name(); ok {
		if err := resource.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Resource.name": %w`, err)}
		}
	}
	if v, ok := ru.mutation.GetType(); ok {
		if err := resource.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Resource.type": %w`, err)}
		}
	}
	if v, ok := ru.mutation.Rrn(); ok {
		if err := resource.RrnValidator(v); err != nil {
			return &ValidationError{Name: "rrn", err: fmt.Errorf(`ent: validator failed for field "Resource.rrn": %w`, err)}
		}
	}
	return nil
}

func (ru *ResourceUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := ru.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(resource.Table, resource.Columns, sqlgraph.NewFieldSpec(resource.FieldID, field.TypeUUID))
	if ps := ru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ru.mutation.UpdatedAt(); ok {
		_spec.SetField(resource.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := ru.mutation.Name(); ok {
		_spec.SetField(resource.FieldName, field.TypeString, value)
	}
	if value, ok := ru.mutation.GetType(); ok {
		_spec.SetField(resource.FieldType, field.TypeString, value)
	}
	if value, ok := ru.mutation.Rrn(); ok {
		_spec.SetField(resource.FieldRrn, field.TypeString, value)
	}
	if value, ok := ru.mutation.Properties(); ok {
		_spec.SetField(resource.FieldProperties, field.TypeJSON, value)
	}
	if ru.mutation.PropertiesCleared() {
		_spec.ClearField(resource.FieldProperties, field.TypeJSON)
	}
	if value, ok := ru.mutation.Tags(); ok {
		_spec.SetField(resource.FieldTags, field.TypeJSON, value)
	}
	if ru.mutation.TagsCleared() {
		_spec.ClearField(resource.FieldTags, field.TypeJSON)
	}
	if value, ok := ru.mutation.Description(); ok {
		_spec.SetField(resource.FieldDescription, field.TypeString, value)
	}
	if ru.mutation.DescriptionCleared() {
		_spec.ClearField(resource.FieldDescription, field.TypeString)
	}
	if ru.mutation.AccountCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resource.AccountTable,
			Columns: resource.AccountPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.RemovedAccountIDs(); len(nodes) > 0 && !ru.mutation.AccountCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resource.AccountTable,
			Columns: resource.AccountPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.AccountIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resource.AccountTable,
			Columns: resource.AccountPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ru.mutation.ChildrenCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resource.ChildrenTable,
			Columns: resource.ChildrenPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(resource.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.RemovedChildrenIDs(); len(nodes) > 0 && !ru.mutation.ChildrenCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resource.ChildrenTable,
			Columns: resource.ChildrenPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(resource.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.ChildrenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resource.ChildrenTable,
			Columns: resource.ChildrenPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(resource.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ru.mutation.ParentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   resource.ParentTable,
			Columns: resource.ParentPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(resource.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.RemovedParentIDs(); len(nodes) > 0 && !ru.mutation.ParentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   resource.ParentTable,
			Columns: resource.ParentPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(resource.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   resource.ParentTable,
			Columns: resource.ParentPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(resource.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{resource.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ru.mutation.done = true
	return n, nil
}

// ResourceUpdateOne is the builder for updating a single Resource entity.
type ResourceUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ResourceMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (ruo *ResourceUpdateOne) SetUpdatedAt(t time.Time) *ResourceUpdateOne {
	ruo.mutation.SetUpdatedAt(t)
	return ruo
}

// SetName sets the "name" field.
func (ruo *ResourceUpdateOne) SetName(s string) *ResourceUpdateOne {
	ruo.mutation.SetName(s)
	return ruo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (ruo *ResourceUpdateOne) SetNillableName(s *string) *ResourceUpdateOne {
	if s != nil {
		ruo.SetName(*s)
	}
	return ruo
}

// SetType sets the "type" field.
func (ruo *ResourceUpdateOne) SetType(s string) *ResourceUpdateOne {
	ruo.mutation.SetType(s)
	return ruo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (ruo *ResourceUpdateOne) SetNillableType(s *string) *ResourceUpdateOne {
	if s != nil {
		ruo.SetType(*s)
	}
	return ruo
}

// SetRrn sets the "rrn" field.
func (ruo *ResourceUpdateOne) SetRrn(s string) *ResourceUpdateOne {
	ruo.mutation.SetRrn(s)
	return ruo
}

// SetNillableRrn sets the "rrn" field if the given value is not nil.
func (ruo *ResourceUpdateOne) SetNillableRrn(s *string) *ResourceUpdateOne {
	if s != nil {
		ruo.SetRrn(*s)
	}
	return ruo
}

// SetProperties sets the "properties" field.
func (ruo *ResourceUpdateOne) SetProperties(m map[string]interface{}) *ResourceUpdateOne {
	ruo.mutation.SetProperties(m)
	return ruo
}

// ClearProperties clears the value of the "properties" field.
func (ruo *ResourceUpdateOne) ClearProperties() *ResourceUpdateOne {
	ruo.mutation.ClearProperties()
	return ruo
}

// SetTags sets the "tags" field.
func (ruo *ResourceUpdateOne) SetTags(m map[string]string) *ResourceUpdateOne {
	ruo.mutation.SetTags(m)
	return ruo
}

// ClearTags clears the value of the "tags" field.
func (ruo *ResourceUpdateOne) ClearTags() *ResourceUpdateOne {
	ruo.mutation.ClearTags()
	return ruo
}

// SetDescription sets the "description" field.
func (ruo *ResourceUpdateOne) SetDescription(s string) *ResourceUpdateOne {
	ruo.mutation.SetDescription(s)
	return ruo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (ruo *ResourceUpdateOne) SetNillableDescription(s *string) *ResourceUpdateOne {
	if s != nil {
		ruo.SetDescription(*s)
	}
	return ruo
}

// ClearDescription clears the value of the "description" field.
func (ruo *ResourceUpdateOne) ClearDescription() *ResourceUpdateOne {
	ruo.mutation.ClearDescription()
	return ruo
}

// AddAccountIDs adds the "account" edge to the Account entity by IDs.
func (ruo *ResourceUpdateOne) AddAccountIDs(ids ...uuid.UUID) *ResourceUpdateOne {
	ruo.mutation.AddAccountIDs(ids...)
	return ruo
}

// AddAccount adds the "account" edges to the Account entity.
func (ruo *ResourceUpdateOne) AddAccount(a ...*Account) *ResourceUpdateOne {
	ids := make([]uuid.UUID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return ruo.AddAccountIDs(ids...)
}

// AddChildIDs adds the "children" edge to the Resource entity by IDs.
func (ruo *ResourceUpdateOne) AddChildIDs(ids ...uuid.UUID) *ResourceUpdateOne {
	ruo.mutation.AddChildIDs(ids...)
	return ruo
}

// AddChildren adds the "children" edges to the Resource entity.
func (ruo *ResourceUpdateOne) AddChildren(r ...*Resource) *ResourceUpdateOne {
	ids := make([]uuid.UUID, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ruo.AddChildIDs(ids...)
}

// AddParentIDs adds the "parent" edge to the Resource entity by IDs.
func (ruo *ResourceUpdateOne) AddParentIDs(ids ...uuid.UUID) *ResourceUpdateOne {
	ruo.mutation.AddParentIDs(ids...)
	return ruo
}

// AddParent adds the "parent" edges to the Resource entity.
func (ruo *ResourceUpdateOne) AddParent(r ...*Resource) *ResourceUpdateOne {
	ids := make([]uuid.UUID, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ruo.AddParentIDs(ids...)
}

// Mutation returns the ResourceMutation object of the builder.
func (ruo *ResourceUpdateOne) Mutation() *ResourceMutation {
	return ruo.mutation
}

// ClearAccount clears all "account" edges to the Account entity.
func (ruo *ResourceUpdateOne) ClearAccount() *ResourceUpdateOne {
	ruo.mutation.ClearAccount()
	return ruo
}

// RemoveAccountIDs removes the "account" edge to Account entities by IDs.
func (ruo *ResourceUpdateOne) RemoveAccountIDs(ids ...uuid.UUID) *ResourceUpdateOne {
	ruo.mutation.RemoveAccountIDs(ids...)
	return ruo
}

// RemoveAccount removes "account" edges to Account entities.
func (ruo *ResourceUpdateOne) RemoveAccount(a ...*Account) *ResourceUpdateOne {
	ids := make([]uuid.UUID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return ruo.RemoveAccountIDs(ids...)
}

// ClearChildren clears all "children" edges to the Resource entity.
func (ruo *ResourceUpdateOne) ClearChildren() *ResourceUpdateOne {
	ruo.mutation.ClearChildren()
	return ruo
}

// RemoveChildIDs removes the "children" edge to Resource entities by IDs.
func (ruo *ResourceUpdateOne) RemoveChildIDs(ids ...uuid.UUID) *ResourceUpdateOne {
	ruo.mutation.RemoveChildIDs(ids...)
	return ruo
}

// RemoveChildren removes "children" edges to Resource entities.
func (ruo *ResourceUpdateOne) RemoveChildren(r ...*Resource) *ResourceUpdateOne {
	ids := make([]uuid.UUID, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ruo.RemoveChildIDs(ids...)
}

// ClearParent clears all "parent" edges to the Resource entity.
func (ruo *ResourceUpdateOne) ClearParent() *ResourceUpdateOne {
	ruo.mutation.ClearParent()
	return ruo
}

// RemoveParentIDs removes the "parent" edge to Resource entities by IDs.
func (ruo *ResourceUpdateOne) RemoveParentIDs(ids ...uuid.UUID) *ResourceUpdateOne {
	ruo.mutation.RemoveParentIDs(ids...)
	return ruo
}

// RemoveParent removes "parent" edges to Resource entities.
func (ruo *ResourceUpdateOne) RemoveParent(r ...*Resource) *ResourceUpdateOne {
	ids := make([]uuid.UUID, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ruo.RemoveParentIDs(ids...)
}

// Where appends a list predicates to the ResourceUpdate builder.
func (ruo *ResourceUpdateOne) Where(ps ...predicate.Resource) *ResourceUpdateOne {
	ruo.mutation.Where(ps...)
	return ruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ruo *ResourceUpdateOne) Select(field string, fields ...string) *ResourceUpdateOne {
	ruo.fields = append([]string{field}, fields...)
	return ruo
}

// Save executes the query and returns the updated Resource entity.
func (ruo *ResourceUpdateOne) Save(ctx context.Context) (*Resource, error) {
	ruo.defaults()
	return withHooks(ctx, ruo.sqlSave, ruo.mutation, ruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ruo *ResourceUpdateOne) SaveX(ctx context.Context) *Resource {
	node, err := ruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ruo *ResourceUpdateOne) Exec(ctx context.Context) error {
	_, err := ruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ruo *ResourceUpdateOne) ExecX(ctx context.Context) {
	if err := ruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ruo *ResourceUpdateOne) defaults() {
	if _, ok := ruo.mutation.UpdatedAt(); !ok {
		v := resource.UpdateDefaultUpdatedAt()
		ruo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ruo *ResourceUpdateOne) check() error {
	if v, ok := ruo.mutation.Name(); ok {
		if err := resource.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Resource.name": %w`, err)}
		}
	}
	if v, ok := ruo.mutation.GetType(); ok {
		if err := resource.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Resource.type": %w`, err)}
		}
	}
	if v, ok := ruo.mutation.Rrn(); ok {
		if err := resource.RrnValidator(v); err != nil {
			return &ValidationError{Name: "rrn", err: fmt.Errorf(`ent: validator failed for field "Resource.rrn": %w`, err)}
		}
	}
	return nil
}

func (ruo *ResourceUpdateOne) sqlSave(ctx context.Context) (_node *Resource, err error) {
	if err := ruo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(resource.Table, resource.Columns, sqlgraph.NewFieldSpec(resource.FieldID, field.TypeUUID))
	id, ok := ruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Resource.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, resource.FieldID)
		for _, f := range fields {
			if !resource.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != resource.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ruo.mutation.UpdatedAt(); ok {
		_spec.SetField(resource.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := ruo.mutation.Name(); ok {
		_spec.SetField(resource.FieldName, field.TypeString, value)
	}
	if value, ok := ruo.mutation.GetType(); ok {
		_spec.SetField(resource.FieldType, field.TypeString, value)
	}
	if value, ok := ruo.mutation.Rrn(); ok {
		_spec.SetField(resource.FieldRrn, field.TypeString, value)
	}
	if value, ok := ruo.mutation.Properties(); ok {
		_spec.SetField(resource.FieldProperties, field.TypeJSON, value)
	}
	if ruo.mutation.PropertiesCleared() {
		_spec.ClearField(resource.FieldProperties, field.TypeJSON)
	}
	if value, ok := ruo.mutation.Tags(); ok {
		_spec.SetField(resource.FieldTags, field.TypeJSON, value)
	}
	if ruo.mutation.TagsCleared() {
		_spec.ClearField(resource.FieldTags, field.TypeJSON)
	}
	if value, ok := ruo.mutation.Description(); ok {
		_spec.SetField(resource.FieldDescription, field.TypeString, value)
	}
	if ruo.mutation.DescriptionCleared() {
		_spec.ClearField(resource.FieldDescription, field.TypeString)
	}
	if ruo.mutation.AccountCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resource.AccountTable,
			Columns: resource.AccountPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.RemovedAccountIDs(); len(nodes) > 0 && !ruo.mutation.AccountCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resource.AccountTable,
			Columns: resource.AccountPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.AccountIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resource.AccountTable,
			Columns: resource.AccountPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ruo.mutation.ChildrenCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resource.ChildrenTable,
			Columns: resource.ChildrenPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(resource.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.RemovedChildrenIDs(); len(nodes) > 0 && !ruo.mutation.ChildrenCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resource.ChildrenTable,
			Columns: resource.ChildrenPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(resource.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.ChildrenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resource.ChildrenTable,
			Columns: resource.ChildrenPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(resource.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ruo.mutation.ParentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   resource.ParentTable,
			Columns: resource.ParentPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(resource.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.RemovedParentIDs(); len(nodes) > 0 && !ruo.mutation.ParentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   resource.ParentTable,
			Columns: resource.ParentPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(resource.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   resource.ParentTable,
			Columns: resource.ParentPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(resource.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Resource{config: ruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{resource.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ruo.mutation.done = true
	return _node, nil
}
