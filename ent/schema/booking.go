package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Booking struct {
	ent.Schema
}

func (Booking) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Annotations(entgql.OrderField("ID")),

		field.String("title").
			MaxLen(255).
			NotEmpty().
			Annotations(entgql.OrderField("TITLE")),

		field.Time("start_date"),

		field.Time("end_date"),

		field.Bool("is_repeat").
			Optional().
			Default(false),

		field.UUID("user_id", uuid.UUID{}),

		field.UUID("office_id", uuid.UUID{}),

		field.UUID("room_id", uuid.UUID{}),
	}
}

func (Booking) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("bookings").
			Unique().
			Field("user_id").
			Required(),

		edge.From("office", Office.Type).
			Ref("bookings").
			Unique().
			Field("office_id").
			Required(),

		edge.From("room", Room.Type).
			Ref("bookings").
			Unique().
			Field("room_id").
			Required(),
	}
}

func (Booking) Mixin() []ent.Mixin {
	return []ent.Mixin{
		CommonMixin{},
		SlugMixin{},
	}
}
