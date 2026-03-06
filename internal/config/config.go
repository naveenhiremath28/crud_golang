package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Config holds all application configuration
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWKSURL    string
	VaultURL   string
	VaultToken string
	AppHost    string
	AppPort    string
}

// Load reads environment variables and returns a Config instance
func Load(log *zap.SugaredLogger) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found, using system environment variables")
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "postgres"),
		JWKSURL:    getEnv("JWKS_URL", "http://localhost:8083/realms/employee-realm/protocol/openid-connect/certs"),
		VaultURL:   getEnv("VAULT_URL", "http://localhost:8200"),
		VaultToken: getEnv("VAULT_TOKEN", "root"),
		AppHost:    getEnv("APP_HOST", "localhost"),
		AppPort:    getEnv("APP_PORT", "8080"),
	}

	log.Info("Configuration loaded successfully")
	return cfg, nil
}

// getEnv retrieves an environment variable or returns a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
