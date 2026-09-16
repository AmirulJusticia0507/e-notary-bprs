package db

import (
	"database/sql"
	"fmt"

	"github.com/e-notary-bprs/backend/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := cfg.DSN()
	database, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(5)
	if err := database.Ping(); err != nil {
		database.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return database, nil
}
