package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin0b101/book-rental-api/service/book"
	"github.com/martin0b101/book-rental-api/service/user"
)

type ApiServer struct {
	address string
	database *sql.DB
}


func NewApiServer(address string, database *sql.DB) (*ApiServer) {
	return &ApiServer{
		address: address,
		database:  database,
	}
}

func (s *ApiServer) Run() error {
	router := gin.Default()

	// register user store 
	userStore := user.NewStore(s.database)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(router)

	// register book store
	bookStore := book.NewStore(s.database)
	bookHandler := book.NewHandler(bookStore)
	bookHandler.RegisterRoutes(router)
	

	log.Println("Listening on ", s.address)

	return http.ListenAndServe(s.address, router)
}