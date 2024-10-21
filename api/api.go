package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin0b101/book-rental-api/service/book"
	"github.com/martin0b101/book-rental-api/service/user"
	"github.com/redis/go-redis/v9"
)

type ApiServer struct {
	address string
	database *sql.DB
	cache *redis.Client
	ctx context.Context
}


func NewApiServer(address string, database *sql.DB, redisCache *redis.Client, ctx context.Context) (*ApiServer) {
	return &ApiServer{
		address: address,
		database:  database,
		cache: redisCache,
		ctx: ctx,
	}
}

func (s *ApiServer) Run() error {
	router := gin.Default()

	// register user store 
	userStore := user.NewStore(s.database)
	userHandler := user.NewHandler(userStore, s.cache, s.ctx)
	userHandler.RegisterRoutes(router)

	// register book store
	bookStore := book.NewStore(s.database)
	bookHandler := book.NewHandler(bookStore, s.cache, s.ctx)
	bookHandler.RegisterRoutes(router)
	
	log.Println("Listening on ", s.address)

	return http.ListenAndServe(s.address, router)
}