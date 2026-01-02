package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
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
}

// Load reads environment variables and returns a Config instance
func Load() (*Config, error) {
	// Best-effort loading of .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "postgres"),
		JWKSURL:    getEnv("JWKS_URL", "http://localhost:8083/realms/employee-realm/protocol/openid-connect/certs"),
		VaultURL:   getEnv("VAULT_URL", "http://localhost:8200"),
		VaultToken: getEnv("VAULT_TOKEN", "root"),
	}, nil
}

// getEnv retrieves an environment variable or returns a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
