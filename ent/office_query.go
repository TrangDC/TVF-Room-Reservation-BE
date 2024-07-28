// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/booking"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/office"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/predicate"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/room"
)

// OfficeQuery is the builder for querying Office entities.
type OfficeQuery struct {
	config
	ctx          *QueryContext
	order        []office.OrderOption
	inters       []Interceptor
	predicates   []predicate.Office
	withRooms    *RoomQuery
	withBookings *BookingQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the OfficeQuery builder.
func (oq *OfficeQuery) Where(ps ...predicate.Office) *OfficeQuery {
	oq.predicates = append(oq.predicates, ps...)
	return oq
}

// Limit the number of records to be returned by this query.
func (oq *OfficeQuery) Limit(limit int) *OfficeQuery {
	oq.ctx.Limit = &limit
	return oq
}

// Offset to start from.
func (oq *OfficeQuery) Offset(offset int) *OfficeQuery {
	oq.ctx.Offset = &offset
	return oq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (oq *OfficeQuery) Unique(unique bool) *OfficeQuery {
	oq.ctx.Unique = &unique
	return oq
}

// Order specifies how the records should be ordered.
func (oq *OfficeQuery) Order(o ...office.OrderOption) *OfficeQuery {
	oq.order = append(oq.order, o...)
	return oq
}

// QueryRooms chains the current query on the "rooms" edge.
func (oq *OfficeQuery) QueryRooms() *RoomQuery {
	query := (&RoomClient{config: oq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := oq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := oq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(office.Table, office.FieldID, selector),
			sqlgraph.To(room.Table, room.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, office.RoomsTable, office.RoomsColumn),
		)
		fromU = sqlgraph.SetNeighbors(oq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryBookings chains the current query on the "bookings" edge.
func (oq *OfficeQuery) QueryBookings() *BookingQuery {
	query := (&BookingClient{config: oq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := oq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := oq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(office.Table, office.FieldID, selector),
			sqlgraph.To(booking.Table, booking.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, office.BookingsTable, office.BookingsColumn),
		)
		fromU = sqlgraph.SetNeighbors(oq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Office entity from the query.
// Returns a *NotFoundError when no Office was found.
func (oq *OfficeQuery) First(ctx context.Context) (*Office, error) {
	nodes, err := oq.Limit(1).All(setContextOp(ctx, oq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{office.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (oq *OfficeQuery) FirstX(ctx context.Context) *Office {
	node, err := oq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Office ID from the query.
// Returns a *NotFoundError when no Office ID was found.
func (oq *OfficeQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = oq.Limit(1).IDs(setContextOp(ctx, oq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{office.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (oq *OfficeQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := oq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Office entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Office entity is found.
// Returns a *NotFoundError when no Office entities are found.
func (oq *OfficeQuery) Only(ctx context.Context) (*Office, error) {
	nodes, err := oq.Limit(2).All(setContextOp(ctx, oq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{office.Label}
	default:
		return nil, &NotSingularError{office.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (oq *OfficeQuery) OnlyX(ctx context.Context) *Office {
	node, err := oq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Office ID in the query.
// Returns a *NotSingularError when more than one Office ID is found.
// Returns a *NotFoundError when no entities are found.
func (oq *OfficeQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = oq.Limit(2).IDs(setContextOp(ctx, oq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{office.Label}
	default:
		err = &NotSingularError{office.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (oq *OfficeQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := oq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Offices.
func (oq *OfficeQuery) All(ctx context.Context) ([]*Office, error) {
	ctx = setContextOp(ctx, oq.ctx, "All")
	if err := oq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Office, *OfficeQuery]()
	return withInterceptors[[]*Office](ctx, oq, qr, oq.inters)
}

// AllX is like All, but panics if an error occurs.
func (oq *OfficeQuery) AllX(ctx context.Context) []*Office {
	nodes, err := oq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Office IDs.
func (oq *OfficeQuery) IDs(ctx context.Context) (ids []uuid.UUID, err error) {
	if oq.ctx.Unique == nil && oq.path != nil {
		oq.Unique(true)
	}
	ctx = setContextOp(ctx, oq.ctx, "IDs")
	if err = oq.Select(office.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (oq *OfficeQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := oq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (oq *OfficeQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, oq.ctx, "Count")
	if err := oq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, oq, querierCount[*OfficeQuery](), oq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (oq *OfficeQuery) CountX(ctx context.Context) int {
	count, err := oq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (oq *OfficeQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, oq.ctx, "Exist")
	switch _, err := oq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (oq *OfficeQuery) ExistX(ctx context.Context) bool {
	exist, err := oq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the OfficeQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (oq *OfficeQuery) Clone() *OfficeQuery {
	if oq == nil {
		return nil
	}
	return &OfficeQuery{
		config:       oq.config,
		ctx:          oq.ctx.Clone(),
		order:        append([]office.OrderOption{}, oq.order...),
		inters:       append([]Interceptor{}, oq.inters...),
		predicates:   append([]predicate.Office{}, oq.predicates...),
		withRooms:    oq.withRooms.Clone(),
		withBookings: oq.withBookings.Clone(),
		// clone intermediate query.
		sql:  oq.sql.Clone(),
		path: oq.path,
	}
}

// WithRooms tells the query-builder to eager-load the nodes that are connected to
// the "rooms" edge. The optional arguments are used to configure the query builder of the edge.
func (oq *OfficeQuery) WithRooms(opts ...func(*RoomQuery)) *OfficeQuery {
	query := (&RoomClient{config: oq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	oq.withRooms = query
	return oq
}

// WithBookings tells the query-builder to eager-load the nodes that are connected to
// the "bookings" edge. The optional arguments are used to configure the query builder of the edge.
func (oq *OfficeQuery) WithBookings(opts ...func(*BookingQuery)) *OfficeQuery {
	query := (&BookingClient{config: oq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	oq.withBookings = query
	return oq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Office.Query().
//		GroupBy(office.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (oq *OfficeQuery) GroupBy(field string, fields ...string) *OfficeGroupBy {
	oq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &OfficeGroupBy{build: oq}
	grbuild.flds = &oq.ctx.Fields
	grbuild.label = office.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.Office.Query().
//		Select(office.FieldName).
//		Scan(ctx, &v)
func (oq *OfficeQuery) Select(fields ...string) *OfficeSelect {
	oq.ctx.Fields = append(oq.ctx.Fields, fields...)
	sbuild := &OfficeSelect{OfficeQuery: oq}
	sbuild.label = office.Label
	sbuild.flds, sbuild.scan = &oq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a OfficeSelect configured with the given aggregations.
func (oq *OfficeQuery) Aggregate(fns ...AggregateFunc) *OfficeSelect {
	return oq.Select().Aggregate(fns...)
}

func (oq *OfficeQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range oq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, oq); err != nil {
				return err
			}
		}
	}
	for _, f := range oq.ctx.Fields {
		if !office.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if oq.path != nil {
		prev, err := oq.path(ctx)
		if err != nil {
			return err
		}
		oq.sql = prev
	}
	return nil
}

func (oq *OfficeQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Office, error) {
	var (
		nodes       = []*Office{}
		_spec       = oq.querySpec()
		loadedTypes = [2]bool{
			oq.withRooms != nil,
			oq.withBookings != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Office).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Office{config: oq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, oq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := oq.withRooms; query != nil {
		if err := oq.loadRooms(ctx, query, nodes,
			func(n *Office) { n.Edges.Rooms = []*Room{} },
			func(n *Office, e *Room) { n.Edges.Rooms = append(n.Edges.Rooms, e) }); err != nil {
			return nil, err
		}
	}
	if query := oq.withBookings; query != nil {
		if err := oq.loadBookings(ctx, query, nodes,
			func(n *Office) { n.Edges.Bookings = []*Booking{} },
			func(n *Office, e *Booking) { n.Edges.Bookings = append(n.Edges.Bookings, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (oq *OfficeQuery) loadRooms(ctx context.Context, query *RoomQuery, nodes []*Office, init func(*Office), assign func(*Office, *Room)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*Office)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(room.FieldOfficeID)
	}
	query.Where(predicate.Room(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(office.RoomsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.OfficeID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "office_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (oq *OfficeQuery) loadBookings(ctx context.Context, query *BookingQuery, nodes []*Office, init func(*Office), assign func(*Office, *Booking)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*Office)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(booking.FieldOfficeID)
	}
	query.Where(predicate.Booking(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(office.BookingsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.OfficeID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "office_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (oq *OfficeQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := oq.querySpec()
	_spec.Node.Columns = oq.ctx.Fields
	if len(oq.ctx.Fields) > 0 {
		_spec.Unique = oq.ctx.Unique != nil && *oq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, oq.driver, _spec)
}

func (oq *OfficeQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(office.Table, office.Columns, sqlgraph.NewFieldSpec(office.FieldID, field.TypeUUID))
	_spec.From = oq.sql
	if unique := oq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if oq.path != nil {
		_spec.Unique = true
	}
	if fields := oq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, office.FieldID)
		for i := range fields {
			if fields[i] != office.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := oq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := oq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := oq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := oq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (oq *OfficeQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(oq.driver.Dialect())
	t1 := builder.Table(office.Table)
	columns := oq.ctx.Fields
	if len(columns) == 0 {
		columns = office.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if oq.sql != nil {
		selector = oq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if oq.ctx.Unique != nil && *oq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range oq.predicates {
		p(selector)
	}
	for _, p := range oq.order {
		p(selector)
	}
	if offset := oq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := oq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// OfficeGroupBy is the group-by builder for Office entities.
type OfficeGroupBy struct {
	selector
	build *OfficeQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ogb *OfficeGroupBy) Aggregate(fns ...AggregateFunc) *OfficeGroupBy {
	ogb.fns = append(ogb.fns, fns...)
	return ogb
}

// Scan applies the selector query and scans the result into the given value.
func (ogb *OfficeGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ogb.build.ctx, "GroupBy")
	if err := ogb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*OfficeQuery, *OfficeGroupBy](ctx, ogb.build, ogb, ogb.build.inters, v)
}

func (ogb *OfficeGroupBy) sqlScan(ctx context.Context, root *OfficeQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(ogb.fns))
	for _, fn := range ogb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*ogb.flds)+len(ogb.fns))
		for _, f := range *ogb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*ogb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ogb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// OfficeSelect is the builder for selecting fields of Office entities.
type OfficeSelect struct {
	*OfficeQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (os *OfficeSelect) Aggregate(fns ...AggregateFunc) *OfficeSelect {
	os.fns = append(os.fns, fns...)
	return os
}

// Scan applies the selector query and scans the result into the given value.
func (os *OfficeSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, os.ctx, "Select")
	if err := os.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*OfficeQuery, *OfficeSelect](ctx, os.OfficeQuery, os, os.inters, v)
}

func (os *OfficeSelect) sqlScan(ctx context.Context, root *OfficeQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(os.fns))
	for _, fn := range os.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*os.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := os.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
