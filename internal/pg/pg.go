package pg

import (
	"database/sql"
	"time"

	"gitlab.techvify.dev/its/internship/q2-2024/project/meeting-room-reservation/meeting-room-reservation-be/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// NewDBConnection creates a new DB connection.
func NewDBConnection(config config.PostgresConfig, logger *zap.Logger) (*sql.DB, error) {
	// Connect to pgx.
	db, err := sql.Open("pgx", config.ConnectionString)
	logger.Debug("Connecting to database", zap.String("connection_string", config.ConnectionString))
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(29 * time.Minute) //Azure's default is 30 mins, so we set it to 29 mins to be safe.

	// Pings the database.
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
