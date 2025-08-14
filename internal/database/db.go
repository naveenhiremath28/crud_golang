package database

import (
	"log"
	"fmt"
	// "database/sql"
    // _ "github.com/lib/pq"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Connect() error {
	var err error
	if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_NAME", "mydb"),
	)

	// DB, err = sql.Open("postgres", connStr)
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	}
	// if err = DB.Ping(); err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }
	fmt.Println("Database Connected Successfully..!")
	return nil
}