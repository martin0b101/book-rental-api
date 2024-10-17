package main

import (
	"database/sql"
	"log"

	"github.com/martin0b101/book-rental-api/api"
	"github.com/martin0b101/book-rental-api/config"
	"github.com/martin0b101/book-rental-api/database"
)

func main() {

	db, err := database.NewPostgresSQLStorage(config.Config{
		Host: config.Envs.Host,
		Port: config.Envs.Port,
		DBName: config.Envs.DBName,
		User: config.Envs.User,
		Password: config.Envs.Password,
	})

	if err != nil{
		log.Fatal(err)
	}

	initDatabase(db)

	server := api.NewApiServer(":8080", db)
	if err := server.Run(); err != nil{
		log.Fatal(err)
	}
}


func initDatabase(db *sql.DB){
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}