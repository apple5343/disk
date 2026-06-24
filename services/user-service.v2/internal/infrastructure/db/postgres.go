package db

import (
	"database/sql"
	"fmt"

	// pg driver.
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(dbURL string) (*sql.DB, error) {
	sqlDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}

	return sqlDB, nil
}
