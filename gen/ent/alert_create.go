// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/alertchain/gen/ent/alert"
	"github.com/m-mizutani/alertchain/gen/ent/attribute"
	"github.com/m-mizutani/alertchain/gen/ent/job"
	"github.com/m-mizutani/alertchain/gen/ent/reference"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

// AlertCreate is the builder for creating a Alert entity.
type AlertCreate struct {
	config
	mutation *AlertMutation
	hooks    []Hook
}

// SetTitle sets the "title" field.
func (ac *AlertCreate) SetTitle(s string) *AlertCreate {
	ac.mutation.SetTitle(s)
	return ac
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (ac *AlertCreate) SetNillableTitle(s *string) *AlertCreate {
	if s != nil {
		ac.SetTitle(*s)
	}
	return ac
}

// SetDescription sets the "description" field.
func (ac *AlertCreate) SetDescription(s string) *AlertCreate {
	ac.mutation.SetDescription(s)
	return ac
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (ac *AlertCreate) SetNillableDescription(s *string) *AlertCreate {
	if s != nil {
		ac.SetDescription(*s)
	}
	return ac
}

// SetDetector sets the "detector" field.
func (ac *AlertCreate) SetDetector(s string) *AlertCreate {
	ac.mutation.SetDetector(s)
	return ac
}

// SetNillableDetector sets the "detector" field if the given value is not nil.
func (ac *AlertCreate) SetNillableDetector(s *string) *AlertCreate {
	if s != nil {
		ac.SetDetector(*s)
	}
	return ac
}

// SetStatus sets the "status" field.
func (ac *AlertCreate) SetStatus(ts types.AlertStatus) *AlertCreate {
	ac.mutation.SetStatus(ts)
	return ac
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (ac *AlertCreate) SetNillableStatus(ts *types.AlertStatus) *AlertCreate {
	if ts != nil {
		ac.SetStatus(*ts)
	}
	return ac
}

// SetSeverity sets the "severity" field.
func (ac *AlertCreate) SetSeverity(t types.Severity) *AlertCreate {
	ac.mutation.SetSeverity(t)
	return ac
}

// SetNillableSeverity sets the "severity" field if the given value is not nil.
func (ac *AlertCreate) SetNillableSeverity(t *types.Severity) *AlertCreate {
	if t != nil {
		ac.SetSeverity(*t)
	}
	return ac
}

// SetDetectedAt sets the "detected_at" field.
func (ac *AlertCreate) SetDetectedAt(i int64) *AlertCreate {
	ac.mutation.SetDetectedAt(i)
	return ac
}

// SetNillableDetectedAt sets the "detected_at" field if the given value is not nil.
func (ac *AlertCreate) SetNillableDetectedAt(i *int64) *AlertCreate {
	if i != nil {
		ac.SetDetectedAt(*i)
	}
	return ac
}

// SetCreatedAt sets the "created_at" field.
func (ac *AlertCreate) SetCreatedAt(i int64) *AlertCreate {
	ac.mutation.SetCreatedAt(i)
	return ac
}

// SetClosedAt sets the "closed_at" field.
func (ac *AlertCreate) SetClosedAt(i int64) *AlertCreate {
	ac.mutation.SetClosedAt(i)
	return ac
}

// SetNillableClosedAt sets the "closed_at" field if the given value is not nil.
func (ac *AlertCreate) SetNillableClosedAt(i *int64) *AlertCreate {
	if i != nil {
		ac.SetClosedAt(*i)
	}
	return ac
}

// SetID sets the "id" field.
func (ac *AlertCreate) SetID(ti types.AlertID) *AlertCreate {
	ac.mutation.SetID(ti)
	return ac
}

// AddAttributeIDs adds the "attributes" edge to the Attribute entity by IDs.
func (ac *AlertCreate) AddAttributeIDs(ids ...int) *AlertCreate {
	ac.mutation.AddAttributeIDs(ids...)
	return ac
}

// AddAttributes adds the "attributes" edges to the Attribute entity.
func (ac *AlertCreate) AddAttributes(a ...*Attribute) *AlertCreate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return ac.AddAttributeIDs(ids...)
}

// AddReferenceIDs adds the "references" edge to the Reference entity by IDs.
func (ac *AlertCreate) AddReferenceIDs(ids ...int) *AlertCreate {
	ac.mutation.AddReferenceIDs(ids...)
	return ac
}

// AddReferences adds the "references" edges to the Reference entity.
func (ac *AlertCreate) AddReferences(r ...*Reference) *AlertCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return ac.AddReferenceIDs(ids...)
}

// AddJobIDs adds the "jobs" edge to the Job entity by IDs.
func (ac *AlertCreate) AddJobIDs(ids ...int) *AlertCreate {
	ac.mutation.AddJobIDs(ids...)
	return ac
}

// AddJobs adds the "jobs" edges to the Job entity.
func (ac *AlertCreate) AddJobs(j ...*Job) *AlertCreate {
	ids := make([]int, len(j))
	for i := range j {
		ids[i] = j[i].ID
	}
	return ac.AddJobIDs(ids...)
}

// Mutation returns the AlertMutation object of the builder.
func (ac *AlertCreate) Mutation() *AlertMutation {
	return ac.mutation
}

// Save creates the Alert in the database.
func (ac *AlertCreate) Save(ctx context.Context) (*Alert, error) {
	var (
		err  error
		node *Alert
	)
	ac.defaults()
	if len(ac.hooks) == 0 {
		if err = ac.check(); err != nil {
			return nil, err
		}
		node, err = ac.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AlertMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = ac.check(); err != nil {
				return nil, err
			}
			ac.mutation = mutation
			if node, err = ac.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(ac.hooks) - 1; i >= 0; i-- {
			if ac.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ac.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ac.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ac *AlertCreate) SaveX(ctx context.Context) *Alert {
	v, err := ac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ac *AlertCreate) Exec(ctx context.Context) error {
	_, err := ac.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ac *AlertCreate) ExecX(ctx context.Context) {
	if err := ac.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ac *AlertCreate) defaults() {
	if _, ok := ac.mutation.Status(); !ok {
		v := alert.DefaultStatus
		ac.mutation.SetStatus(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ac *AlertCreate) check() error {
	if _, ok := ac.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "status"`)}
	}
	if _, ok := ac.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "created_at"`)}
	}
	return nil
}

func (ac *AlertCreate) sqlSave(ctx context.Context) (*Alert, error) {
	_node, _spec := ac.createSpec()
	if err := sqlgraph.CreateNode(ctx, ac.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		_node.ID = _spec.ID.Value.(types.AlertID)
	}
	return _node, nil
}

func (ac *AlertCreate) createSpec() (*Alert, *sqlgraph.CreateSpec) {
	var (
		_node = &Alert{config: ac.config}
		_spec = &sqlgraph.CreateSpec{
			Table: alert.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: alert.FieldID,
			},
		}
	)
	if id, ok := ac.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := ac.mutation.Title(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: alert.FieldTitle,
		})
		_node.Title = value
	}
	if value, ok := ac.mutation.Description(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: alert.FieldDescription,
		})
		_node.Description = value
	}
	if value, ok := ac.mutation.Detector(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: alert.FieldDetector,
		})
		_node.Detector = value
	}
	if value, ok := ac.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: alert.FieldStatus,
		})
		_node.Status = value
	}
	if value, ok := ac.mutation.Severity(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: alert.FieldSeverity,
		})
		_node.Severity = value
	}
	if value, ok := ac.mutation.DetectedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: alert.FieldDetectedAt,
		})
		_node.DetectedAt = value
	}
	if value, ok := ac.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: alert.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := ac.mutation.ClosedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: alert.FieldClosedAt,
		})
		_node.ClosedAt = value
	}
	if nodes := ac.mutation.AttributesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   alert.AttributesTable,
			Columns: []string{alert.AttributesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: attribute.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ac.mutation.ReferencesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   alert.ReferencesTable,
			Columns: []string{alert.ReferencesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: reference.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ac.mutation.JobsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   alert.JobsTable,
			Columns: []string{alert.JobsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: job.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// AlertCreateBulk is the builder for creating many Alert entities in bulk.
type AlertCreateBulk struct {
	config
	builders []*AlertCreate
}

// Save creates the Alert entities in the database.
func (acb *AlertCreateBulk) Save(ctx context.Context) ([]*Alert, error) {
	specs := make([]*sqlgraph.CreateSpec, len(acb.builders))
	nodes := make([]*Alert, len(acb.builders))
	mutators := make([]Mutator, len(acb.builders))
	for i := range acb.builders {
		func(i int, root context.Context) {
			builder := acb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*AlertMutation)
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
					_, err = mutators[i+1].Mutate(root, acb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, acb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, acb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (acb *AlertCreateBulk) SaveX(ctx context.Context) []*Alert {
	v, err := acb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (acb *AlertCreateBulk) Exec(ctx context.Context) error {
	_, err := acb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (acb *AlertCreateBulk) ExecX(ctx context.Context) {
	if err := acb.Exec(ctx); err != nil {
		panic(err)
	}
}