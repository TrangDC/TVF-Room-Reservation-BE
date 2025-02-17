// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/booking"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/office"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/room"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent/user"
)

// Booking is the model entity for the Booking schema.
type Booking struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	// Slug holds the value of the "slug" field.
	Slug string `json:"slug,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// StartDate holds the value of the "start_date" field.
	StartDate time.Time `json:"start_date,omitempty"`
	// EndDate holds the value of the "end_date" field.
	EndDate time.Time `json:"end_date,omitempty"`
	// IsRepeat holds the value of the "is_repeat" field.
	IsRepeat bool `json:"is_repeat,omitempty"`
	// UserID holds the value of the "user_id" field.
	UserID uuid.UUID `json:"user_id,omitempty"`
	// OfficeID holds the value of the "office_id" field.
	OfficeID uuid.UUID `json:"office_id,omitempty"`
	// RoomID holds the value of the "room_id" field.
	RoomID uuid.UUID `json:"room_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the BookingQuery when eager-loading is set.
	Edges        BookingEdges `json:"edges"`
	selectValues sql.SelectValues
}

// BookingEdges holds the relations/edges for other nodes in the graph.
type BookingEdges struct {
	// User holds the value of the user edge.
	User *User `json:"user,omitempty"`
	// Office holds the value of the office edge.
	Office *Office `json:"office,omitempty"`
	// Room holds the value of the room edge.
	Room *Room `json:"room,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e BookingEdges) UserOrErr() (*User, error) {
	if e.User != nil {
		return e.User, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: user.Label}
	}
	return nil, &NotLoadedError{edge: "user"}
}

// OfficeOrErr returns the Office value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e BookingEdges) OfficeOrErr() (*Office, error) {
	if e.Office != nil {
		return e.Office, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: office.Label}
	}
	return nil, &NotLoadedError{edge: "office"}
}

// RoomOrErr returns the Room value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e BookingEdges) RoomOrErr() (*Room, error) {
	if e.Room != nil {
		return e.Room, nil
	} else if e.loadedTypes[2] {
		return nil, &NotFoundError{label: room.Label}
	}
	return nil, &NotLoadedError{edge: "room"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Booking) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case booking.FieldIsRepeat:
			values[i] = new(sql.NullBool)
		case booking.FieldSlug, booking.FieldTitle:
			values[i] = new(sql.NullString)
		case booking.FieldCreatedAt, booking.FieldUpdatedAt, booking.FieldDeletedAt, booking.FieldStartDate, booking.FieldEndDate:
			values[i] = new(sql.NullTime)
		case booking.FieldID, booking.FieldUserID, booking.FieldOfficeID, booking.FieldRoomID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Booking fields.
func (b *Booking) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case booking.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				b.ID = *value
			}
		case booking.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				b.CreatedAt = value.Time
			}
		case booking.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				b.UpdatedAt = value.Time
			}
		case booking.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				b.DeletedAt = value.Time
			}
		case booking.FieldSlug:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field slug", values[i])
			} else if value.Valid {
				b.Slug = value.String
			}
		case booking.FieldTitle:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field title", values[i])
			} else if value.Valid {
				b.Title = value.String
			}
		case booking.FieldStartDate:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field start_date", values[i])
			} else if value.Valid {
				b.StartDate = value.Time
			}
		case booking.FieldEndDate:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field end_date", values[i])
			} else if value.Valid {
				b.EndDate = value.Time
			}
		case booking.FieldIsRepeat:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_repeat", values[i])
			} else if value.Valid {
				b.IsRepeat = value.Bool
			}
		case booking.FieldUserID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field user_id", values[i])
			} else if value != nil {
				b.UserID = *value
			}
		case booking.FieldOfficeID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field office_id", values[i])
			} else if value != nil {
				b.OfficeID = *value
			}
		case booking.FieldRoomID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field room_id", values[i])
			} else if value != nil {
				b.RoomID = *value
			}
		default:
			b.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Booking.
// This includes values selected through modifiers, order, etc.
func (b *Booking) Value(name string) (ent.Value, error) {
	return b.selectValues.Get(name)
}

// QueryUser queries the "user" edge of the Booking entity.
func (b *Booking) QueryUser() *UserQuery {
	return NewBookingClient(b.config).QueryUser(b)
}

// QueryOffice queries the "office" edge of the Booking entity.
func (b *Booking) QueryOffice() *OfficeQuery {
	return NewBookingClient(b.config).QueryOffice(b)
}

// QueryRoom queries the "room" edge of the Booking entity.
func (b *Booking) QueryRoom() *RoomQuery {
	return NewBookingClient(b.config).QueryRoom(b)
}

// Update returns a builder for updating this Booking.
// Note that you need to call Booking.Unwrap() before calling this method if this Booking
// was returned from a transaction, and the transaction was committed or rolled back.
func (b *Booking) Update() *BookingUpdateOne {
	return NewBookingClient(b.config).UpdateOne(b)
}

// Unwrap unwraps the Booking entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (b *Booking) Unwrap() *Booking {
	_tx, ok := b.config.driver.(*txDriver)
	if !ok {
		panic("ent: Booking is not a transactional entity")
	}
	b.config.driver = _tx.drv
	return b
}

// String implements the fmt.Stringer.
func (b *Booking) String() string {
	var builder strings.Builder
	builder.WriteString("Booking(")
	builder.WriteString(fmt.Sprintf("id=%v, ", b.ID))
	builder.WriteString("created_at=")
	builder.WriteString(b.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(b.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("deleted_at=")
	builder.WriteString(b.DeletedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("slug=")
	builder.WriteString(b.Slug)
	builder.WriteString(", ")
	builder.WriteString("title=")
	builder.WriteString(b.Title)
	builder.WriteString(", ")
	builder.WriteString("start_date=")
	builder.WriteString(b.StartDate.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("end_date=")
	builder.WriteString(b.EndDate.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("is_repeat=")
	builder.WriteString(fmt.Sprintf("%v", b.IsRepeat))
	builder.WriteString(", ")
	builder.WriteString("user_id=")
	builder.WriteString(fmt.Sprintf("%v", b.UserID))
	builder.WriteString(", ")
	builder.WriteString("office_id=")
	builder.WriteString(fmt.Sprintf("%v", b.OfficeID))
	builder.WriteString(", ")
	builder.WriteString("room_id=")
	builder.WriteString(fmt.Sprintf("%v", b.RoomID))
	builder.WriteByte(')')
	return builder.String()
}

// Bookings is a parsable slice of Booking.
type Bookings []*Booking
