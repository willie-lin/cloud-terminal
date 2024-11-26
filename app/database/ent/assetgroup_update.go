// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/willie-lin/cloud-terminal/app/database/ent/assetgroup"
	"github.com/willie-lin/cloud-terminal/app/database/ent/predicate"
)

// AssetGroupUpdate is the builder for updating AssetGroup entities.
type AssetGroupUpdate struct {
	config
	hooks    []Hook
	mutation *AssetGroupMutation
}

// Where appends a list predicates to the AssetGroupUpdate builder.
func (agu *AssetGroupUpdate) Where(ps ...predicate.AssetGroup) *AssetGroupUpdate {
	agu.mutation.Where(ps...)
	return agu
}

// Mutation returns the AssetGroupMutation object of the builder.
func (agu *AssetGroupUpdate) Mutation() *AssetGroupMutation {
	return agu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (agu *AssetGroupUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, agu.sqlSave, agu.mutation, agu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (agu *AssetGroupUpdate) SaveX(ctx context.Context) int {
	affected, err := agu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (agu *AssetGroupUpdate) Exec(ctx context.Context) error {
	_, err := agu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (agu *AssetGroupUpdate) ExecX(ctx context.Context) {
	if err := agu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (agu *AssetGroupUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(assetgroup.Table, assetgroup.Columns, sqlgraph.NewFieldSpec(assetgroup.FieldID, field.TypeInt))
	if ps := agu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if n, err = sqlgraph.UpdateNodes(ctx, agu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{assetgroup.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	agu.mutation.done = true
	return n, nil
}

// AssetGroupUpdateOne is the builder for updating a single AssetGroup entity.
type AssetGroupUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *AssetGroupMutation
}

// Mutation returns the AssetGroupMutation object of the builder.
func (aguo *AssetGroupUpdateOne) Mutation() *AssetGroupMutation {
	return aguo.mutation
}

// Where appends a list predicates to the AssetGroupUpdate builder.
func (aguo *AssetGroupUpdateOne) Where(ps ...predicate.AssetGroup) *AssetGroupUpdateOne {
	aguo.mutation.Where(ps...)
	return aguo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (aguo *AssetGroupUpdateOne) Select(field string, fields ...string) *AssetGroupUpdateOne {
	aguo.fields = append([]string{field}, fields...)
	return aguo
}

// Save executes the query and returns the updated AssetGroup entity.
func (aguo *AssetGroupUpdateOne) Save(ctx context.Context) (*AssetGroup, error) {
	return withHooks(ctx, aguo.sqlSave, aguo.mutation, aguo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (aguo *AssetGroupUpdateOne) SaveX(ctx context.Context) *AssetGroup {
	node, err := aguo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (aguo *AssetGroupUpdateOne) Exec(ctx context.Context) error {
	_, err := aguo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aguo *AssetGroupUpdateOne) ExecX(ctx context.Context) {
	if err := aguo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (aguo *AssetGroupUpdateOne) sqlSave(ctx context.Context) (_node *AssetGroup, err error) {
	_spec := sqlgraph.NewUpdateSpec(assetgroup.Table, assetgroup.Columns, sqlgraph.NewFieldSpec(assetgroup.FieldID, field.TypeInt))
	id, ok := aguo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "AssetGroup.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := aguo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, assetgroup.FieldID)
		for _, f := range fields {
			if !assetgroup.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != assetgroup.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := aguo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	_node = &AssetGroup{config: aguo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, aguo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{assetgroup.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	aguo.mutation.done = true
	return _node, nil
}
