package database

import (
	"log"
	"fmt"
	// "database/sql"
    // _ "github.com/lib/pq"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	var err error
	connStr := "host=localhost port=5431 user=postgres password=postgres dbname=go_test sslmode=disable"
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