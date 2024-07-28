package directives

import (
	"context"
	"net/http"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/util"

	"github.com/99designs/gqlgen/graphql"
)

func HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, roles []string) (interface{}, error) {
	userRoles, ok := ctx.Value("roles").([]interface{})
	if !ok {
		return nil, util.WrapGQLError(ctx, "missing roles in context", http.StatusConflict, util.ErrorFlagValidateFail)
	}

	for _, role := range roles {
		for _, userRole := range userRoles {
			if userRole == role {
				return next(ctx)
			}
		}
	}

	return nil, util.WrapGQLError(ctx, "you don't have permission to access on this server", http.StatusForbidden, util.ErrorFlagValidateFail)
}
