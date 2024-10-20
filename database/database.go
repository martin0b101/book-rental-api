package database

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
	"github.com/martin0b101/book-rental-api/config"
)


func NewPostgresSQLStorage(config config.Config) (*sql.DB, error) {
	// Connection string based on provided configuration
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	// Open a connection to PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}


	
	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL")
	return db, nil
}