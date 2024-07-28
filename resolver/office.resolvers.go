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
func (r *officeResolver) ID(ctx context.Context, obj *ent.Office) (string, error) {
	return obj.ID.String(), nil
}

// Office returns graphql1.OfficeResolver implementation.
func (r *Resolver) Office() graphql1.OfficeResolver { return &officeResolver{r} }

type officeResolver struct{ *Resolver }