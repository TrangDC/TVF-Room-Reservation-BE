// Code generated by ent, DO NOT EDIT.

package user

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.User {
	return predicate.User(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldName, v))
}

// WorkEmail applies equality check predicate on the "work_email" field. It's identical to WorkEmailEQ.
func WorkEmail(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldWorkEmail, v))
}

// Oid applies equality check predicate on the "oid" field. It's identical to OidEQ.
func Oid(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldOid, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldName, v))
}

// WorkEmailEQ applies the EQ predicate on the "work_email" field.
func WorkEmailEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldWorkEmail, v))
}

// WorkEmailNEQ applies the NEQ predicate on the "work_email" field.
func WorkEmailNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldWorkEmail, v))
}

// WorkEmailIn applies the In predicate on the "work_email" field.
func WorkEmailIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldWorkEmail, vs...))
}

// WorkEmailNotIn applies the NotIn predicate on the "work_email" field.
func WorkEmailNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldWorkEmail, vs...))
}

// WorkEmailGT applies the GT predicate on the "work_email" field.
func WorkEmailGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldWorkEmail, v))
}

// WorkEmailGTE applies the GTE predicate on the "work_email" field.
func WorkEmailGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldWorkEmail, v))
}

// WorkEmailLT applies the LT predicate on the "work_email" field.
func WorkEmailLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldWorkEmail, v))
}

// WorkEmailLTE applies the LTE predicate on the "work_email" field.
func WorkEmailLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldWorkEmail, v))
}

// WorkEmailContains applies the Contains predicate on the "work_email" field.
func WorkEmailContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldWorkEmail, v))
}

// WorkEmailHasPrefix applies the HasPrefix predicate on the "work_email" field.
func WorkEmailHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldWorkEmail, v))
}

// WorkEmailHasSuffix applies the HasSuffix predicate on the "work_email" field.
func WorkEmailHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldWorkEmail, v))
}

// WorkEmailEqualFold applies the EqualFold predicate on the "work_email" field.
func WorkEmailEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldWorkEmail, v))
}

// WorkEmailContainsFold applies the ContainsFold predicate on the "work_email" field.
func WorkEmailContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldWorkEmail, v))
}

// OidEQ applies the EQ predicate on the "oid" field.
func OidEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldOid, v))
}

// OidNEQ applies the NEQ predicate on the "oid" field.
func OidNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldOid, v))
}

// OidIn applies the In predicate on the "oid" field.
func OidIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldOid, vs...))
}

// OidNotIn applies the NotIn predicate on the "oid" field.
func OidNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldOid, vs...))
}

// OidGT applies the GT predicate on the "oid" field.
func OidGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldOid, v))
}

// OidGTE applies the GTE predicate on the "oid" field.
func OidGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldOid, v))
}

// OidLT applies the LT predicate on the "oid" field.
func OidLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldOid, v))
}

// OidLTE applies the LTE predicate on the "oid" field.
func OidLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldOid, v))
}

// OidContains applies the Contains predicate on the "oid" field.
func OidContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldOid, v))
}

// OidHasPrefix applies the HasPrefix predicate on the "oid" field.
func OidHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldOid, v))
}

// OidHasSuffix applies the HasSuffix predicate on the "oid" field.
func OidHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldOid, v))
}

// OidEqualFold applies the EqualFold predicate on the "oid" field.
func OidEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldOid, v))
}

// OidContainsFold applies the ContainsFold predicate on the "oid" field.
func OidContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldOid, v))
}

// HasBookings applies the HasEdge predicate on the "bookings" edge.
func HasBookings() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, BookingsTable, BookingsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasBookingsWith applies the HasEdge predicate on the "bookings" edge with a given conditions (other predicates).
func HasBookingsWith(preds ...predicate.Booking) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := newBookingsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasUserRoles applies the HasEdge predicate on the "user_roles" edge.
func HasUserRoles() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, UserRolesTable, UserRolesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserRolesWith applies the HasEdge predicate on the "user_roles" edge with a given conditions (other predicates).
func HasUserRolesWith(preds ...predicate.UserRole) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := newUserRolesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.User) predicate.User {
	return predicate.User(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.User) predicate.User {
	return predicate.User(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.User) predicate.User {
	return predicate.User(sql.NotPredicates(p))
}
