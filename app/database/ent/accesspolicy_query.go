// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/willie-lin/cloud-terminal/app/database/ent/accesspolicy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/predicate"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
)

// AccessPolicyQuery is the builder for querying AccessPolicy entities.
type AccessPolicyQuery struct {
	config
	ctx        *QueryContext
	order      []accesspolicy.OrderOption
	inters     []Interceptor
	predicates []predicate.AccessPolicy
	withTenant *TenantQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the AccessPolicyQuery builder.
func (apq *AccessPolicyQuery) Where(ps ...predicate.AccessPolicy) *AccessPolicyQuery {
	apq.predicates = append(apq.predicates, ps...)
	return apq
}

// Limit the number of records to be returned by this query.
func (apq *AccessPolicyQuery) Limit(limit int) *AccessPolicyQuery {
	apq.ctx.Limit = &limit
	return apq
}

// Offset to start from.
func (apq *AccessPolicyQuery) Offset(offset int) *AccessPolicyQuery {
	apq.ctx.Offset = &offset
	return apq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (apq *AccessPolicyQuery) Unique(unique bool) *AccessPolicyQuery {
	apq.ctx.Unique = &unique
	return apq
}

// Order specifies how the records should be ordered.
func (apq *AccessPolicyQuery) Order(o ...accesspolicy.OrderOption) *AccessPolicyQuery {
	apq.order = append(apq.order, o...)
	return apq
}

// QueryTenant chains the current query on the "tenant" edge.
func (apq *AccessPolicyQuery) QueryTenant() *TenantQuery {
	query := (&TenantClient{config: apq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := apq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := apq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(accesspolicy.Table, accesspolicy.FieldID, selector),
			sqlgraph.To(tenant.Table, tenant.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, accesspolicy.TenantTable, accesspolicy.TenantPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(apq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first AccessPolicy entity from the query.
// Returns a *NotFoundError when no AccessPolicy was found.
func (apq *AccessPolicyQuery) First(ctx context.Context) (*AccessPolicy, error) {
	nodes, err := apq.Limit(1).All(setContextOp(ctx, apq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{accesspolicy.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (apq *AccessPolicyQuery) FirstX(ctx context.Context) *AccessPolicy {
	node, err := apq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first AccessPolicy ID from the query.
// Returns a *NotFoundError when no AccessPolicy ID was found.
func (apq *AccessPolicyQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = apq.Limit(1).IDs(setContextOp(ctx, apq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{accesspolicy.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (apq *AccessPolicyQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := apq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single AccessPolicy entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one AccessPolicy entity is found.
// Returns a *NotFoundError when no AccessPolicy entities are found.
func (apq *AccessPolicyQuery) Only(ctx context.Context) (*AccessPolicy, error) {
	nodes, err := apq.Limit(2).All(setContextOp(ctx, apq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{accesspolicy.Label}
	default:
		return nil, &NotSingularError{accesspolicy.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (apq *AccessPolicyQuery) OnlyX(ctx context.Context) *AccessPolicy {
	node, err := apq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only AccessPolicy ID in the query.
// Returns a *NotSingularError when more than one AccessPolicy ID is found.
// Returns a *NotFoundError when no entities are found.
func (apq *AccessPolicyQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = apq.Limit(2).IDs(setContextOp(ctx, apq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{accesspolicy.Label}
	default:
		err = &NotSingularError{accesspolicy.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (apq *AccessPolicyQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := apq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of AccessPolicies.
func (apq *AccessPolicyQuery) All(ctx context.Context) ([]*AccessPolicy, error) {
	ctx = setContextOp(ctx, apq.ctx, ent.OpQueryAll)
	if err := apq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*AccessPolicy, *AccessPolicyQuery]()
	return withInterceptors[[]*AccessPolicy](ctx, apq, qr, apq.inters)
}

// AllX is like All, but panics if an error occurs.
func (apq *AccessPolicyQuery) AllX(ctx context.Context) []*AccessPolicy {
	nodes, err := apq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of AccessPolicy IDs.
func (apq *AccessPolicyQuery) IDs(ctx context.Context) (ids []uuid.UUID, err error) {
	if apq.ctx.Unique == nil && apq.path != nil {
		apq.Unique(true)
	}
	ctx = setContextOp(ctx, apq.ctx, ent.OpQueryIDs)
	if err = apq.Select(accesspolicy.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (apq *AccessPolicyQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := apq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (apq *AccessPolicyQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, apq.ctx, ent.OpQueryCount)
	if err := apq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, apq, querierCount[*AccessPolicyQuery](), apq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (apq *AccessPolicyQuery) CountX(ctx context.Context) int {
	count, err := apq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (apq *AccessPolicyQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, apq.ctx, ent.OpQueryExist)
	switch _, err := apq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (apq *AccessPolicyQuery) ExistX(ctx context.Context) bool {
	exist, err := apq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the AccessPolicyQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (apq *AccessPolicyQuery) Clone() *AccessPolicyQuery {
	if apq == nil {
		return nil
	}
	return &AccessPolicyQuery{
		config:     apq.config,
		ctx:        apq.ctx.Clone(),
		order:      append([]accesspolicy.OrderOption{}, apq.order...),
		inters:     append([]Interceptor{}, apq.inters...),
		predicates: append([]predicate.AccessPolicy{}, apq.predicates...),
		withTenant: apq.withTenant.Clone(),
		// clone intermediate query.
		sql:  apq.sql.Clone(),
		path: apq.path,
	}
}

// WithTenant tells the query-builder to eager-load the nodes that are connected to
// the "tenant" edge. The optional arguments are used to configure the query builder of the edge.
func (apq *AccessPolicyQuery) WithTenant(opts ...func(*TenantQuery)) *AccessPolicyQuery {
	query := (&TenantClient{config: apq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	apq.withTenant = query
	return apq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.AccessPolicy.Query().
//		GroupBy(accesspolicy.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (apq *AccessPolicyQuery) GroupBy(field string, fields ...string) *AccessPolicyGroupBy {
	apq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &AccessPolicyGroupBy{build: apq}
	grbuild.flds = &apq.ctx.Fields
	grbuild.label = accesspolicy.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//	}
//
//	client.AccessPolicy.Query().
//		Select(accesspolicy.FieldCreatedAt).
//		Scan(ctx, &v)
func (apq *AccessPolicyQuery) Select(fields ...string) *AccessPolicySelect {
	apq.ctx.Fields = append(apq.ctx.Fields, fields...)
	sbuild := &AccessPolicySelect{AccessPolicyQuery: apq}
	sbuild.label = accesspolicy.Label
	sbuild.flds, sbuild.scan = &apq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a AccessPolicySelect configured with the given aggregations.
func (apq *AccessPolicyQuery) Aggregate(fns ...AggregateFunc) *AccessPolicySelect {
	return apq.Select().Aggregate(fns...)
}

func (apq *AccessPolicyQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range apq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, apq); err != nil {
				return err
			}
		}
	}
	for _, f := range apq.ctx.Fields {
		if !accesspolicy.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if apq.path != nil {
		prev, err := apq.path(ctx)
		if err != nil {
			return err
		}
		apq.sql = prev
	}
	return nil
}

func (apq *AccessPolicyQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*AccessPolicy, error) {
	var (
		nodes       = []*AccessPolicy{}
		_spec       = apq.querySpec()
		loadedTypes = [1]bool{
			apq.withTenant != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*AccessPolicy).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &AccessPolicy{config: apq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, apq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := apq.withTenant; query != nil {
		if err := apq.loadTenant(ctx, query, nodes,
			func(n *AccessPolicy) { n.Edges.Tenant = []*Tenant{} },
			func(n *AccessPolicy, e *Tenant) { n.Edges.Tenant = append(n.Edges.Tenant, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (apq *AccessPolicyQuery) loadTenant(ctx context.Context, query *TenantQuery, nodes []*AccessPolicy, init func(*AccessPolicy), assign func(*AccessPolicy, *Tenant)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[uuid.UUID]*AccessPolicy)
	nids := make(map[uuid.UUID]map[*AccessPolicy]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(accesspolicy.TenantTable)
		s.Join(joinT).On(s.C(tenant.FieldID), joinT.C(accesspolicy.TenantPrimaryKey[0]))
		s.Where(sql.InValues(joinT.C(accesspolicy.TenantPrimaryKey[1]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(accesspolicy.TenantPrimaryKey[1]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(uuid.UUID)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := *values[0].(*uuid.UUID)
				inValue := *values[1].(*uuid.UUID)
				if nids[inValue] == nil {
					nids[inValue] = map[*AccessPolicy]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*Tenant](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "tenant" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}

func (apq *AccessPolicyQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := apq.querySpec()
	_spec.Node.Columns = apq.ctx.Fields
	if len(apq.ctx.Fields) > 0 {
		_spec.Unique = apq.ctx.Unique != nil && *apq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, apq.driver, _spec)
}

func (apq *AccessPolicyQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(accesspolicy.Table, accesspolicy.Columns, sqlgraph.NewFieldSpec(accesspolicy.FieldID, field.TypeUUID))
	_spec.From = apq.sql
	if unique := apq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if apq.path != nil {
		_spec.Unique = true
	}
	if fields := apq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, accesspolicy.FieldID)
		for i := range fields {
			if fields[i] != accesspolicy.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := apq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := apq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := apq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := apq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (apq *AccessPolicyQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(apq.driver.Dialect())
	t1 := builder.Table(accesspolicy.Table)
	columns := apq.ctx.Fields
	if len(columns) == 0 {
		columns = accesspolicy.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if apq.sql != nil {
		selector = apq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if apq.ctx.Unique != nil && *apq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range apq.predicates {
		p(selector)
	}
	for _, p := range apq.order {
		p(selector)
	}
	if offset := apq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := apq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// AccessPolicyGroupBy is the group-by builder for AccessPolicy entities.
type AccessPolicyGroupBy struct {
	selector
	build *AccessPolicyQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (apgb *AccessPolicyGroupBy) Aggregate(fns ...AggregateFunc) *AccessPolicyGroupBy {
	apgb.fns = append(apgb.fns, fns...)
	return apgb
}

// Scan applies the selector query and scans the result into the given value.
func (apgb *AccessPolicyGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, apgb.build.ctx, ent.OpQueryGroupBy)
	if err := apgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AccessPolicyQuery, *AccessPolicyGroupBy](ctx, apgb.build, apgb, apgb.build.inters, v)
}

func (apgb *AccessPolicyGroupBy) sqlScan(ctx context.Context, root *AccessPolicyQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(apgb.fns))
	for _, fn := range apgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*apgb.flds)+len(apgb.fns))
		for _, f := range *apgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*apgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := apgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// AccessPolicySelect is the builder for selecting fields of AccessPolicy entities.
type AccessPolicySelect struct {
	*AccessPolicyQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (aps *AccessPolicySelect) Aggregate(fns ...AggregateFunc) *AccessPolicySelect {
	aps.fns = append(aps.fns, fns...)
	return aps
}

// Scan applies the selector query and scans the result into the given value.
func (aps *AccessPolicySelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, aps.ctx, ent.OpQuerySelect)
	if err := aps.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AccessPolicyQuery, *AccessPolicySelect](ctx, aps.AccessPolicyQuery, aps, aps.inters, v)
}

func (aps *AccessPolicySelect) sqlScan(ctx context.Context, root *AccessPolicyQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(aps.fns))
	for _, fn := range aps.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*aps.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := aps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}