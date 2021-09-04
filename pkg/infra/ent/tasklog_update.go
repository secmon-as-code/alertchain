// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/alertchain/pkg/infra/ent/annotation"
	"github.com/m-mizutani/alertchain/pkg/infra/ent/execlog"
	"github.com/m-mizutani/alertchain/pkg/infra/ent/predicate"
	"github.com/m-mizutani/alertchain/pkg/infra/ent/tasklog"
)

// TaskLogUpdate is the builder for updating TaskLog entities.
type TaskLogUpdate struct {
	config
	hooks    []Hook
	mutation *TaskLogMutation
}

// Where appends a list predicates to the TaskLogUpdate builder.
func (tlu *TaskLogUpdate) Where(ps ...predicate.TaskLog) *TaskLogUpdate {
	tlu.mutation.Where(ps...)
	return tlu
}

// AddAnnotatedIDs adds the "annotated" edge to the Annotation entity by IDs.
func (tlu *TaskLogUpdate) AddAnnotatedIDs(ids ...int) *TaskLogUpdate {
	tlu.mutation.AddAnnotatedIDs(ids...)
	return tlu
}

// AddAnnotated adds the "annotated" edges to the Annotation entity.
func (tlu *TaskLogUpdate) AddAnnotated(a ...*Annotation) *TaskLogUpdate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return tlu.AddAnnotatedIDs(ids...)
}

// AddExecLogIDs adds the "exec_logs" edge to the ExecLog entity by IDs.
func (tlu *TaskLogUpdate) AddExecLogIDs(ids ...int) *TaskLogUpdate {
	tlu.mutation.AddExecLogIDs(ids...)
	return tlu
}

// AddExecLogs adds the "exec_logs" edges to the ExecLog entity.
func (tlu *TaskLogUpdate) AddExecLogs(e ...*ExecLog) *TaskLogUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return tlu.AddExecLogIDs(ids...)
}

// Mutation returns the TaskLogMutation object of the builder.
func (tlu *TaskLogUpdate) Mutation() *TaskLogMutation {
	return tlu.mutation
}

// ClearAnnotated clears all "annotated" edges to the Annotation entity.
func (tlu *TaskLogUpdate) ClearAnnotated() *TaskLogUpdate {
	tlu.mutation.ClearAnnotated()
	return tlu
}

// RemoveAnnotatedIDs removes the "annotated" edge to Annotation entities by IDs.
func (tlu *TaskLogUpdate) RemoveAnnotatedIDs(ids ...int) *TaskLogUpdate {
	tlu.mutation.RemoveAnnotatedIDs(ids...)
	return tlu
}

// RemoveAnnotated removes "annotated" edges to Annotation entities.
func (tlu *TaskLogUpdate) RemoveAnnotated(a ...*Annotation) *TaskLogUpdate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return tlu.RemoveAnnotatedIDs(ids...)
}

// ClearExecLogs clears all "exec_logs" edges to the ExecLog entity.
func (tlu *TaskLogUpdate) ClearExecLogs() *TaskLogUpdate {
	tlu.mutation.ClearExecLogs()
	return tlu
}

// RemoveExecLogIDs removes the "exec_logs" edge to ExecLog entities by IDs.
func (tlu *TaskLogUpdate) RemoveExecLogIDs(ids ...int) *TaskLogUpdate {
	tlu.mutation.RemoveExecLogIDs(ids...)
	return tlu
}

// RemoveExecLogs removes "exec_logs" edges to ExecLog entities.
func (tlu *TaskLogUpdate) RemoveExecLogs(e ...*ExecLog) *TaskLogUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return tlu.RemoveExecLogIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tlu *TaskLogUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(tlu.hooks) == 0 {
		affected, err = tlu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TaskLogMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			tlu.mutation = mutation
			affected, err = tlu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(tlu.hooks) - 1; i >= 0; i-- {
			if tlu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = tlu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tlu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (tlu *TaskLogUpdate) SaveX(ctx context.Context) int {
	affected, err := tlu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tlu *TaskLogUpdate) Exec(ctx context.Context) error {
	_, err := tlu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tlu *TaskLogUpdate) ExecX(ctx context.Context) {
	if err := tlu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tlu *TaskLogUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   tasklog.Table,
			Columns: tasklog.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: tasklog.FieldID,
			},
		},
	}
	if ps := tlu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if tlu.mutation.AnnotatedCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tasklog.AnnotatedTable,
			Columns: []string{tasklog.AnnotatedColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: annotation.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tlu.mutation.RemovedAnnotatedIDs(); len(nodes) > 0 && !tlu.mutation.AnnotatedCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tasklog.AnnotatedTable,
			Columns: []string{tasklog.AnnotatedColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: annotation.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tlu.mutation.AnnotatedIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tasklog.AnnotatedTable,
			Columns: []string{tasklog.AnnotatedColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: annotation.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if tlu.mutation.ExecLogsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tasklog.ExecLogsTable,
			Columns: []string{tasklog.ExecLogsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: execlog.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tlu.mutation.RemovedExecLogsIDs(); len(nodes) > 0 && !tlu.mutation.ExecLogsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tasklog.ExecLogsTable,
			Columns: []string{tasklog.ExecLogsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: execlog.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tlu.mutation.ExecLogsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tasklog.ExecLogsTable,
			Columns: []string{tasklog.ExecLogsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: execlog.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tlu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{tasklog.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// TaskLogUpdateOne is the builder for updating a single TaskLog entity.
type TaskLogUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TaskLogMutation
}

// AddAnnotatedIDs adds the "annotated" edge to the Annotation entity by IDs.
func (tluo *TaskLogUpdateOne) AddAnnotatedIDs(ids ...int) *TaskLogUpdateOne {
	tluo.mutation.AddAnnotatedIDs(ids...)
	return tluo
}

// AddAnnotated adds the "annotated" edges to the Annotation entity.
func (tluo *TaskLogUpdateOne) AddAnnotated(a ...*Annotation) *TaskLogUpdateOne {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return tluo.AddAnnotatedIDs(ids...)
}

// AddExecLogIDs adds the "exec_logs" edge to the ExecLog entity by IDs.
func (tluo *TaskLogUpdateOne) AddExecLogIDs(ids ...int) *TaskLogUpdateOne {
	tluo.mutation.AddExecLogIDs(ids...)
	return tluo
}

// AddExecLogs adds the "exec_logs" edges to the ExecLog entity.
func (tluo *TaskLogUpdateOne) AddExecLogs(e ...*ExecLog) *TaskLogUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return tluo.AddExecLogIDs(ids...)
}

// Mutation returns the TaskLogMutation object of the builder.
func (tluo *TaskLogUpdateOne) Mutation() *TaskLogMutation {
	return tluo.mutation
}

// ClearAnnotated clears all "annotated" edges to the Annotation entity.
func (tluo *TaskLogUpdateOne) ClearAnnotated() *TaskLogUpdateOne {
	tluo.mutation.ClearAnnotated()
	return tluo
}

// RemoveAnnotatedIDs removes the "annotated" edge to Annotation entities by IDs.
func (tluo *TaskLogUpdateOne) RemoveAnnotatedIDs(ids ...int) *TaskLogUpdateOne {
	tluo.mutation.RemoveAnnotatedIDs(ids...)
	return tluo
}

// RemoveAnnotated removes "annotated" edges to Annotation entities.
func (tluo *TaskLogUpdateOne) RemoveAnnotated(a ...*Annotation) *TaskLogUpdateOne {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return tluo.RemoveAnnotatedIDs(ids...)
}

// ClearExecLogs clears all "exec_logs" edges to the ExecLog entity.
func (tluo *TaskLogUpdateOne) ClearExecLogs() *TaskLogUpdateOne {
	tluo.mutation.ClearExecLogs()
	return tluo
}

// RemoveExecLogIDs removes the "exec_logs" edge to ExecLog entities by IDs.
func (tluo *TaskLogUpdateOne) RemoveExecLogIDs(ids ...int) *TaskLogUpdateOne {
	tluo.mutation.RemoveExecLogIDs(ids...)
	return tluo
}

// RemoveExecLogs removes "exec_logs" edges to ExecLog entities.
func (tluo *TaskLogUpdateOne) RemoveExecLogs(e ...*ExecLog) *TaskLogUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return tluo.RemoveExecLogIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tluo *TaskLogUpdateOne) Select(field string, fields ...string) *TaskLogUpdateOne {
	tluo.fields = append([]string{field}, fields...)
	return tluo
}

// Save executes the query and returns the updated TaskLog entity.
func (tluo *TaskLogUpdateOne) Save(ctx context.Context) (*TaskLog, error) {
	var (
		err  error
		node *TaskLog
	)
	if len(tluo.hooks) == 0 {
		node, err = tluo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TaskLogMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			tluo.mutation = mutation
			node, err = tluo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(tluo.hooks) - 1; i >= 0; i-- {
			if tluo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = tluo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, tluo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (tluo *TaskLogUpdateOne) SaveX(ctx context.Context) *TaskLog {
	node, err := tluo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tluo *TaskLogUpdateOne) Exec(ctx context.Context) error {
	_, err := tluo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tluo *TaskLogUpdateOne) ExecX(ctx context.Context) {
	if err := tluo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tluo *TaskLogUpdateOne) sqlSave(ctx context.Context) (_node *TaskLog, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   tasklog.Table,
			Columns: tasklog.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: tasklog.FieldID,
			},
		},
	}
	id, ok := tluo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing TaskLog.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := tluo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, tasklog.FieldID)
		for _, f := range fields {
			if !tasklog.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != tasklog.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tluo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if tluo.mutation.AnnotatedCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tasklog.AnnotatedTable,
			Columns: []string{tasklog.AnnotatedColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: annotation.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tluo.mutation.RemovedAnnotatedIDs(); len(nodes) > 0 && !tluo.mutation.AnnotatedCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tasklog.AnnotatedTable,
			Columns: []string{tasklog.AnnotatedColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: annotation.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tluo.mutation.AnnotatedIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tasklog.AnnotatedTable,
			Columns: []string{tasklog.AnnotatedColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: annotation.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if tluo.mutation.ExecLogsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tasklog.ExecLogsTable,
			Columns: []string{tasklog.ExecLogsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: execlog.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tluo.mutation.RemovedExecLogsIDs(); len(nodes) > 0 && !tluo.mutation.ExecLogsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tasklog.ExecLogsTable,
			Columns: []string{tasklog.ExecLogsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: execlog.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tluo.mutation.ExecLogsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tasklog.ExecLogsTable,
			Columns: []string{tasklog.ExecLogsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: execlog.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &TaskLog{config: tluo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tluo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{tasklog.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}