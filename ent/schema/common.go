package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// CommonMixin holds the schema definition for the jobTitle entity.
type CommonMixin struct {
	mixin.Schema
}

// Fields of the CommonMixin. NOTE : Part of the public API ( ultimately exposed to end users
func (CommonMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Annotations(entgql.OrderField("ID")),
		field.Time("created_at").Default(time.Now).Immutable().Annotations(entgql.OrderField("CREATED_AT")),
		field.Time("updated_at").Optional().Annotations(entgql.OrderField("UPDATED_AT")),
		field.Time("deleted_at").Optional().Annotations(entgql.OrderField("DELETED_AT")),
	}
}

type SlugMixin struct {
	mixin.Schema
}

func (SlugMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("slug").MaxLen(255).Unique().Annotations(entgql.OrderField("SLUG")),
	}
}
