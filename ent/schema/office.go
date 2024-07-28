package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Office holds the schema definition for the Office entity.
type Office struct {
	ent.Schema
}

// Fields of the Office.
func (Office) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Annotations(entgql.OrderField("ID")),

		field.String("name").
			NotEmpty().
			MaxLen(255).
			Annotations(entgql.OrderField("NAME")),

		field.String("description").
			Optional(),
	}
}

// Edges of the Office.
func (Office) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("rooms", Room.Type),
		edge.To("bookings", Booking.Type),
	}
}
