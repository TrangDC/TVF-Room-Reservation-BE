package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.48

import (
	"context"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/ent"
	graphql1 "gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/graphql"
)

// ID is the resolver for the id field.
func (r *roomResolver) ID(ctx context.Context, obj *ent.Room) (string, error) {
	return obj.ID.String(), nil
}

// OfficeID is the resolver for the office_id field.
func (r *roomResolver) OfficeID(ctx context.Context, obj *ent.Room) (string, error) {
	return obj.OfficeID.String(), nil
}

// Room returns graphql1.RoomResolver implementation.
func (r *Resolver) Room() graphql1.RoomResolver { return &roomResolver{r} }

type roomResolver struct{ *Resolver }