// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/booking"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/office"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/room"
)

// OfficeCreate is the builder for creating a Office entity.
type OfficeCreate struct {
	config
	mutation *OfficeMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (oc *OfficeCreate) SetName(s string) *OfficeCreate {
	oc.mutation.SetName(s)
	return oc
}

// SetDescription sets the "description" field.
func (oc *OfficeCreate) SetDescription(s string) *OfficeCreate {
	oc.mutation.SetDescription(s)
	return oc
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (oc *OfficeCreate) SetNillableDescription(s *string) *OfficeCreate {
	if s != nil {
		oc.SetDescription(*s)
	}
	return oc
}

// SetID sets the "id" field.
func (oc *OfficeCreate) SetID(u uuid.UUID) *OfficeCreate {
	oc.mutation.SetID(u)
	return oc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (oc *OfficeCreate) SetNillableID(u *uuid.UUID) *OfficeCreate {
	if u != nil {
		oc.SetID(*u)
	}
	return oc
}

// AddRoomIDs adds the "rooms" edge to the Room entity by IDs.
func (oc *OfficeCreate) AddRoomIDs(ids ...uuid.UUID) *OfficeCreate {
	oc.mutation.AddRoomIDs(ids...)
	return oc
}

// AddRooms adds the "rooms" edges to the Room entity.
func (oc *OfficeCreate) AddRooms(r ...*Room) *OfficeCreate {
	ids := make([]uuid.UUID, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return oc.AddRoomIDs(ids...)
}

// AddBookingIDs adds the "bookings" edge to the Booking entity by IDs.
func (oc *OfficeCreate) AddBookingIDs(ids ...uuid.UUID) *OfficeCreate {
	oc.mutation.AddBookingIDs(ids...)
	return oc
}

// AddBookings adds the "bookings" edges to the Booking entity.
func (oc *OfficeCreate) AddBookings(b ...*Booking) *OfficeCreate {
	ids := make([]uuid.UUID, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return oc.AddBookingIDs(ids...)
}

// Mutation returns the OfficeMutation object of the builder.
func (oc *OfficeCreate) Mutation() *OfficeMutation {
	return oc.mutation
}

// Save creates the Office in the database.
func (oc *OfficeCreate) Save(ctx context.Context) (*Office, error) {
	oc.defaults()
	return withHooks(ctx, oc.sqlSave, oc.mutation, oc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (oc *OfficeCreate) SaveX(ctx context.Context) *Office {
	v, err := oc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (oc *OfficeCreate) Exec(ctx context.Context) error {
	_, err := oc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (oc *OfficeCreate) ExecX(ctx context.Context) {
	if err := oc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (oc *OfficeCreate) defaults() {
	if _, ok := oc.mutation.ID(); !ok {
		v := office.DefaultID()
		oc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (oc *OfficeCreate) check() error {
	if _, ok := oc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Office.name"`)}
	}
	if v, ok := oc.mutation.Name(); ok {
		if err := office.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Office.name": %w`, err)}
		}
	}
	return nil
}

func (oc *OfficeCreate) sqlSave(ctx context.Context) (*Office, error) {
	if err := oc.check(); err != nil {
		return nil, err
	}
	_node, _spec := oc.createSpec()
	if err := sqlgraph.CreateNode(ctx, oc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	oc.mutation.id = &_node.ID
	oc.mutation.done = true
	return _node, nil
}

func (oc *OfficeCreate) createSpec() (*Office, *sqlgraph.CreateSpec) {
	var (
		_node = &Office{config: oc.config}
		_spec = sqlgraph.NewCreateSpec(office.Table, sqlgraph.NewFieldSpec(office.FieldID, field.TypeUUID))
	)
	if id, ok := oc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := oc.mutation.Name(); ok {
		_spec.SetField(office.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := oc.mutation.Description(); ok {
		_spec.SetField(office.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if nodes := oc.mutation.RoomsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   office.RoomsTable,
			Columns: []string{office.RoomsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(room.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := oc.mutation.BookingsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   office.BookingsTable,
			Columns: []string{office.BookingsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(booking.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OfficeCreateBulk is the builder for creating many Office entities in bulk.
type OfficeCreateBulk struct {
	config
	err      error
	builders []*OfficeCreate
}

// Save creates the Office entities in the database.
func (ocb *OfficeCreateBulk) Save(ctx context.Context) ([]*Office, error) {
	if ocb.err != nil {
		return nil, ocb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ocb.builders))
	nodes := make([]*Office, len(ocb.builders))
	mutators := make([]Mutator, len(ocb.builders))
	for i := range ocb.builders {
		func(i int, root context.Context) {
			builder := ocb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*OfficeMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ocb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ocb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
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
		if _, err := mutators[0].Mutate(ctx, ocb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ocb *OfficeCreateBulk) SaveX(ctx context.Context) []*Office {
	v, err := ocb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ocb *OfficeCreateBulk) Exec(ctx context.Context) error {
	_, err := ocb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ocb *OfficeCreateBulk) ExecX(ctx context.Context) {
	if err := ocb.Exec(ctx); err != nil {
		panic(err)
	}
}