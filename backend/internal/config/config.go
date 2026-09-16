package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config berisi semua konfigurasi aplikasi.
type Config struct {
	Database     DatabaseConfig
	Server       ServerConfig
	JWT          JWTConfig
	CBS          CBSConfig
	BPN          BPNConfig
}

type CBSConfig struct {
	BaseURL    string
	APIKey     string
	ClientCode string
}

type BPNConfig struct {
	BaseURL    string
	APIKey     string
	SecretKey  string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	DBURL    string
}

type ServerConfig struct {
	Host string
	Port int
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

// Load membaca config dari environment file dan environment variables.
func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "bprs_enotary"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			DBURL:    firstEnv("DATABASE_URL", "POSTGRES_URL"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvAsInt("PORT", getEnvAsInt("SERVER_PORT", 8080)),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "e-notary-bprs-secret-key"),
			Expiry: getEnvAsDuration("JWT_EXPIRY", 24*time.Hour),
		},
		CBS: CBSConfig{
			BaseURL:    getEnv("CBS_BASE_URL", ""),
			APIKey:     getEnv("CBS_API_KEY", ""),
			ClientCode: getEnv("CBS_CLIENT_CODE", ""),
		},
		BPN: BPNConfig{
			BaseURL:    getEnv("BPN_BASE_URL", ""),
			APIKey:     getEnv("BPN_API_KEY", ""),
			SecretKey:  getEnv("BPN_SECRET_KEY", ""),
		},
	}
	return cfg, nil
}

// DSN mengembalikan PostgreSQL connection string.
func (c *DatabaseConfig) DSN() string {
	if c.DBURL != "" {
		return c.DBURL
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return ""
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if v, err := strconv.Atoi(value); err == nil {
		return v
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if v, err := time.ParseDuration(value); err == nil {
		return v
	}
	return fallback
}
