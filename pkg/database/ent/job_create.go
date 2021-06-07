// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/willie-lin/cloud-terminal/pkg/database/ent/job"
)

// JobCreate is the builder for creating a Job entity.
type JobCreate struct {
	config
	mutation *JobMutation
	hooks    []Hook
}

// SetCronjobid sets the "cronjobid" field.
func (jc *JobCreate) SetCronjobid(i int) *JobCreate {
	jc.mutation.SetCronjobid(i)
	return jc
}

// SetName sets the "name" field.
func (jc *JobCreate) SetName(s string) *JobCreate {
	jc.mutation.SetName(s)
	return jc
}

// SetFunc sets the "func" field.
func (jc *JobCreate) SetFunc(s string) *JobCreate {
	jc.mutation.SetFunc(s)
	return jc
}

// SetCron sets the "cron" field.
func (jc *JobCreate) SetCron(s string) *JobCreate {
	jc.mutation.SetCron(s)
	return jc
}

// SetMode sets the "mode" field.
func (jc *JobCreate) SetMode(s string) *JobCreate {
	jc.mutation.SetMode(s)
	return jc
}

// SetResourceIds sets the "resourceIds" field.
func (jc *JobCreate) SetResourceIds(s string) *JobCreate {
	jc.mutation.SetResourceIds(s)
	return jc
}

// SetStatus sets the "status" field.
func (jc *JobCreate) SetStatus(s string) *JobCreate {
	jc.mutation.SetStatus(s)
	return jc
}

// SetMetadata sets the "metadata" field.
func (jc *JobCreate) SetMetadata(s string) *JobCreate {
	jc.mutation.SetMetadata(s)
	return jc
}

// SetCreatedAt sets the "created_at" field.
func (jc *JobCreate) SetCreatedAt(t time.Time) *JobCreate {
	jc.mutation.SetCreatedAt(t)
	return jc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (jc *JobCreate) SetNillableCreatedAt(t *time.Time) *JobCreate {
	if t != nil {
		jc.SetCreatedAt(*t)
	}
	return jc
}

// SetUpdatedAt sets the "updated_at" field.
func (jc *JobCreate) SetUpdatedAt(t time.Time) *JobCreate {
	jc.mutation.SetUpdatedAt(t)
	return jc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (jc *JobCreate) SetNillableUpdatedAt(t *time.Time) *JobCreate {
	if t != nil {
		jc.SetUpdatedAt(*t)
	}
	return jc
}

// SetID sets the "id" field.
func (jc *JobCreate) SetID(s string) *JobCreate {
	jc.mutation.SetID(s)
	return jc
}

// Mutation returns the JobMutation object of the builder.
func (jc *JobCreate) Mutation() *JobMutation {
	return jc.mutation
}

// Save creates the Job in the database.
func (jc *JobCreate) Save(ctx context.Context) (*Job, error) {
	var (
		err  error
		node *Job
	)
	jc.defaults()
	if len(jc.hooks) == 0 {
		if err = jc.check(); err != nil {
			return nil, err
		}
		node, err = jc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*JobMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = jc.check(); err != nil {
				return nil, err
			}
			jc.mutation = mutation
			node, err = jc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(jc.hooks) - 1; i >= 0; i-- {
			mut = jc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, jc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (jc *JobCreate) SaveX(ctx context.Context) *Job {
	v, err := jc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// defaults sets the default values of the builder before save.
func (jc *JobCreate) defaults() {
	if _, ok := jc.mutation.CreatedAt(); !ok {
		v := job.DefaultCreatedAt()
		jc.mutation.SetCreatedAt(v)
	}
	if _, ok := jc.mutation.UpdatedAt(); !ok {
		v := job.DefaultUpdatedAt()
		jc.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (jc *JobCreate) check() error {
	if _, ok := jc.mutation.Cronjobid(); !ok {
		return &ValidationError{Name: "cronjobid", err: errors.New("ent: missing required field \"cronjobid\"")}
	}
	if _, ok := jc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New("ent: missing required field \"name\"")}
	}
	if _, ok := jc.mutation.Func(); !ok {
		return &ValidationError{Name: "func", err: errors.New("ent: missing required field \"func\"")}
	}
	if _, ok := jc.mutation.Cron(); !ok {
		return &ValidationError{Name: "cron", err: errors.New("ent: missing required field \"cron\"")}
	}
	if _, ok := jc.mutation.Mode(); !ok {
		return &ValidationError{Name: "mode", err: errors.New("ent: missing required field \"mode\"")}
	}
	if _, ok := jc.mutation.ResourceIds(); !ok {
		return &ValidationError{Name: "resourceIds", err: errors.New("ent: missing required field \"resourceIds\"")}
	}
	if _, ok := jc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New("ent: missing required field \"status\"")}
	}
	if _, ok := jc.mutation.Metadata(); !ok {
		return &ValidationError{Name: "metadata", err: errors.New("ent: missing required field \"metadata\"")}
	}
	if _, ok := jc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New("ent: missing required field \"created_at\"")}
	}
	if _, ok := jc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New("ent: missing required field \"updated_at\"")}
	}
	return nil
}

func (jc *JobCreate) sqlSave(ctx context.Context) (*Job, error) {
	_node, _spec := jc.createSpec()
	if err := sqlgraph.CreateNode(ctx, jc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return _node, nil
}

func (jc *JobCreate) createSpec() (*Job, *sqlgraph.CreateSpec) {
	var (
		_node = &Job{config: jc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: job.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: job.FieldID,
			},
		}
	)
	if id, ok := jc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := jc.mutation.Cronjobid(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: job.FieldCronjobid,
		})
		_node.Cronjobid = value
	}
	if value, ok := jc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: job.FieldName,
		})
		_node.Name = value
	}
	if value, ok := jc.mutation.Func(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: job.FieldFunc,
		})
		_node.Func = value
	}
	if value, ok := jc.mutation.Cron(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: job.FieldCron,
		})
		_node.Cron = value
	}
	if value, ok := jc.mutation.Mode(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: job.FieldMode,
		})
		_node.Mode = value
	}
	if value, ok := jc.mutation.ResourceIds(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: job.FieldResourceIds,
		})
		_node.ResourceIds = value
	}
	if value, ok := jc.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: job.FieldStatus,
		})
		_node.Status = value
	}
	if value, ok := jc.mutation.Metadata(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: job.FieldMetadata,
		})
		_node.Metadata = value
	}
	if value, ok := jc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: job.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := jc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: job.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	return _node, _spec
}

// JobCreateBulk is the builder for creating many Job entities in bulk.
type JobCreateBulk struct {
	config
	builders []*JobCreate
}

// Save creates the Job entities in the database.
func (jcb *JobCreateBulk) Save(ctx context.Context) ([]*Job, error) {
	specs := make([]*sqlgraph.CreateSpec, len(jcb.builders))
	nodes := make([]*Job, len(jcb.builders))
	mutators := make([]Mutator, len(jcb.builders))
	for i := range jcb.builders {
		func(i int, root context.Context) {
			builder := jcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*JobMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, jcb.builders[i+1].mutation)
				} else {
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, jcb.driver, &sqlgraph.BatchCreateSpec{Nodes: specs}); err != nil {
						if cerr, ok := isSQLConstraintError(err); ok {
							err = cerr
						}
					}
				}
				mutation.done = true
				if err != nil {
					return nil, err
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, jcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (jcb *JobCreateBulk) SaveX(ctx context.Context) []*Job {
	v, err := jcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}