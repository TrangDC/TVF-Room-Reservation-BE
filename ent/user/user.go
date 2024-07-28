// Code generated by ent, DO NOT EDIT.

package user

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldWorkEmail holds the string denoting the work_email field in the database.
	FieldWorkEmail = "work_email"
	// FieldOid holds the string denoting the oid field in the database.
	FieldOid = "oid"
	// EdgeBookings holds the string denoting the bookings edge name in mutations.
	EdgeBookings = "bookings"
	// EdgeUserRoles holds the string denoting the user_roles edge name in mutations.
	EdgeUserRoles = "user_roles"
	// Table holds the table name of the user in the database.
	Table = "users"
	// BookingsTable is the table that holds the bookings relation/edge.
	BookingsTable = "bookings"
	// BookingsInverseTable is the table name for the Booking entity.
	// It exists in this package in order to avoid circular dependency with the "booking" package.
	BookingsInverseTable = "bookings"
	// BookingsColumn is the table column denoting the bookings relation/edge.
	BookingsColumn = "user_id"
	// UserRolesTable is the table that holds the user_roles relation/edge.
	UserRolesTable = "user_roles"
	// UserRolesInverseTable is the table name for the UserRole entity.
	// It exists in this package in order to avoid circular dependency with the "userrole" package.
	UserRolesInverseTable = "user_roles"
	// UserRolesColumn is the table column denoting the user_roles relation/edge.
	UserRolesColumn = "user_id"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldWorkEmail,
	FieldOid,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// WorkEmailValidator is a validator for the "work_email" field. It is called by the builders before save.
	WorkEmailValidator func(string) error
	// OidValidator is a validator for the "oid" field. It is called by the builders before save.
	OidValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the User queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByWorkEmail orders the results by the work_email field.
func ByWorkEmail(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldWorkEmail, opts...).ToFunc()
}

// ByOid orders the results by the oid field.
func ByOid(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOid, opts...).ToFunc()
}

// ByBookingsCount orders the results by bookings count.
func ByBookingsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newBookingsStep(), opts...)
	}
}

// ByBookings orders the results by bookings terms.
func ByBookings(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newBookingsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByUserRolesCount orders the results by user_roles count.
func ByUserRolesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newUserRolesStep(), opts...)
	}
}

// ByUserRoles orders the results by user_roles terms.
func ByUserRoles(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newUserRolesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newBookingsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(BookingsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, BookingsTable, BookingsColumn),
	)
}
func newUserRolesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(UserRolesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, UserRolesTable, UserRolesColumn),
	)
}
