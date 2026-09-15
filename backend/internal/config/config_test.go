package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	for _, k := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE", "SERVER_HOST", "SERVER_PORT", "JWT_SECRET", "JWT_EXPIRY"} {
		_ = os.Unsetenv(k)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("Database.Host = %q, want localhost", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("Database.Port = %d, want 5432", cfg.Database.Port)
	}
	if cfg.Database.Name != "bprs_enotary" {
		t.Errorf("Database.Name = %q, want bprs_enotary", cfg.Database.Name)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %d, want 8080", cfg.Server.Port)
	}
	if cfg.JWT.Expiry != 24*time.Hour {
		t.Errorf("JWT.Expiry = %v, want 24h", cfg.JWT.Expiry)
	}
}

func TestDatabaseDSN(t *testing.T) {
	cfg := &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Name:     "bprs_enotary",
		SSLMode:  "disable",
	}
	want := "postgres://postgres:postgres@localhost:5432/bprs_enotary?sslmode=disable"
	if got := cfg.DSN(); got != want {
		t.Errorf("DSN() = %q, want %q", got, want)
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("DB_HOST", "db.internal")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_NAME", "custom_db")
	t.Setenv("JWT_EXPIRY", "2h")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Database.Host != "db.internal" {
		t.Errorf("Database.Host = %q, want db.internal", cfg.Database.Host)
	}
	if cfg.Database.Port != 5433 {
		t.Errorf("Database.Port = %d, want 5433", cfg.Database.Port)
	}
	if cfg.JWT.Expiry != 2*time.Hour {
		t.Errorf("JWT.Expiry = %v, want 2h", cfg.JWT.Expiry)
	}
}
