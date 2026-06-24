package postgres

import (
	"data/internal/config"

	"github.com/jmoiron/sqlx"
	// Register PostgreSQL driver for database/sql.
	_ "github.com/lib/pq"
)

func NewClient(cfg *config.PostgresConfig) (*sqlx.DB, error) {
	client, err := sqlx.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	err = client.Ping()
	if err != nil {
		return nil, err
	}
	return client, nil
}
