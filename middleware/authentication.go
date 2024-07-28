package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/azuread"
	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func AuthenticateMiddleware(oauthClient azuread.AzureADOAuth, db *sql.DB, logger *zap.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// Skip pre-flight request
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		token := ParseBearerTokenFromRequest(c.Request)
		if len(token) == 0 || oauthClient.VerifyAccessToken(ctx, token) != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.WrapGQLUnauthorizedError(ctx))
			return
		}
		tokenData, err := oauthClient.DecodeToken(ctx, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.WrapGQLUnauthorizedError(ctx))
			return
		}

		var id uuid.UUID
		_ = db.QueryRow("SELECT id FROM users WHERE oid = $1", tokenData.ObjectID).Scan(&id)
		if id == uuid.Nil {
			_, err = db.Query("WITH upsert AS ( UPDATE users SET name = $2, work_email = $3 WHERE oid = $1 RETURNING * ) "+
				"INSERT INTO users (oid, name, work_email) SELECT $1, $2, $3 WHERE NOT EXISTS ( SELECT 1 FROM upsert );", tokenData.ObjectID,
				tokenData.Name,
				tokenData.PreferredUsername)
			if err != nil {
				logger.Error("Failed to upsert user", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, util.WrapGQLUnauthorizedError(ctx))
				return
			}
			_ = db.QueryRow("SELECT id FROM users WHERE oid = $1", tokenData.ObjectID).Scan(&id)

			// Assign default "user" role to the new user
			_, err = db.Exec(`
				INSERT INTO user_roles (id, user_id, role_id)
				SELECT uuid_generate_v4(), $1, roles.id
				FROM roles
				WHERE roles.machine_name = 'user'
				AND NOT EXISTS (
					SELECT 1 FROM user_roles
					WHERE user_id = $1 AND role_id = roles.id
				)
			`, id)
			if err != nil {
				logger.Error("", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, util.WrapGQLUnauthorizedError(ctx))
				return
			}

		}

		// Get the roles of the user
		rows, err := db.Query("SELECT r.machine_name FROM roles AS r JOIN user_roles AS ur ON r.id = ur.role_id WHERE ur.user_id = $1", id)
		if err != nil {
			logger.Error("Failed to get user roles", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.WrapGQLInternalError(ctx))
			return
		}
		defer rows.Close()

		var roles []interface{}
		for rows.Next() {
			var role interface{}
			if err := rows.Scan(&role); err != nil {
				logger.Error("Failed to scan role", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, util.WrapGQLInternalError(ctx))
				return
			}
			roles = append(roles, role)
		}
		if err := rows.Err(); err != nil {
			logger.Error("Failed to iterate over roles", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.WrapGQLInternalError(ctx))
			return
		}

		// Add user id and roles to the context
		ctx = context.WithValue(ctx, "user_id", id)
		ctx = context.WithValue(ctx, "roles", roles)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// ParseBearerTokenFromRequest parses the bearer token from request
func ParseBearerTokenFromRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:6]) == "BEARER" {
		// Default jwt token
		return authHeader[7:]
	}

	return ""
}
