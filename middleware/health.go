package middleware

import (
	"context"
	"database/sql"

	"net/http"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/internal/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// dbCheck checks if the database is alive
func dbCheck(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, "SELECT 1")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		return nil
	}

	return rows.Err()
}

// LivelinessCheckMiddleware is a middleware to check if the server is alive
func ReadinessCheckMiddleware(db *sql.DB, logger *zap.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if err := dbCheck(ctx, db); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.WrapGQLInternalError(ctx))
			return
		}

		c.Status(http.StatusOK)
	}
}

// LivelinessCheckMiddleware is a middleware to check if the server is alive
func LivelinessCheckMiddleware(db *sql.DB, logger *zap.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if err := dbCheck(ctx, db); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, util.WrapGQLInternalError(ctx))
			return
		}
		c.Status(http.StatusOK)
	}
}
