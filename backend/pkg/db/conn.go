package db

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/e-notary-bprs/backend/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := url.URL{
		Scheme: "postgres",
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		User:   url.UserPassword(cfg.User, cfg.Password),
		Path:   cfg.Name,
	}
	query := dsn.Query()
	query.Set("sslmode", cfg.SSLMode)
	dsn.RawQuery = query.Encode()
	database, err := sql.Open("pgx", dsn.String())
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
