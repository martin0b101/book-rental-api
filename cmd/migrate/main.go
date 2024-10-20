package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/martin0b101/book-rental-api/config"
	"github.com/martin0b101/book-rental-api/database"
)

func main() {
	db, err := database.NewPostgresSQLStorage(config.Config{
		Host:     config.Envs.Host,
		Port:     config.Envs.Port,
		DBName:   config.Envs.DBName,
		User:     config.Envs.User,
		Password: config.Envs.Password,
	})

	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		log.Fatal(err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"PostgreSQL",
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		err := migration.Up()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Migration up successful.")
	}

	if cmd == "down" {
		err := migration.Down()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Migration down successful.")
	}
}
