package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/martin0b101/book-rental-api/types"
	"github.com/martin0b101/book-rental-api/utils"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	store types.UserStore
	cache *redis.Client
	ctx context.Context
}

func NewHandler(store types.UserStore, cache *redis.Client, ctx context.Context) *Handler{
	return &Handler{
		store: store,
		cache: cache,
		ctx: ctx,
	}
}

func (handler *Handler) RegisterRoutes(router *gin.Engine){
	router.POST("/register", handler.registerUser)
	router.GET("/users", handler.getUsers)
}

func (handler *Handler) registerUser(c *gin.Context){

	var registerRequest types.RegisterUserRequest

	err := c.BindJSON(&registerRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status:  http.StatusBadRequest,
			Error:   true,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	user, errCreate := handler.store.CreateUser(registerRequest)
	if errCreate != nil{
		c.JSON(http.StatusInternalServerError, types.Response{
			Status:  http.StatusInternalServerError,
			Error:   true,
			Message: errCreate.Error(),
			Data:    nil,
		})
		return
	}

	
	utils.DeleteFromRedis(handler.ctx, handler.cache, "Users") 

	c.JSON(http.StatusOK, types.Response{
		Status: http.StatusOK,
		Error: false,
		Data: user,
	})
}

func (handler *Handler) getUsers(c *gin.Context){

	value, err := handler.cache.Get(handler.ctx, "Users").Result()

	if err == nil {
		var users []types.User
		json.Unmarshal([]byte(value), &users)
		c.JSON(http.StatusOK, types.Response{
			Status:  http.StatusOK,
			Error:   false,
			Data:    users,
		})
		return
	}
	users, err := handler.store.GetUsers()

	if err != nil{
		c.JSON(http.StatusInternalServerError, types.Response{
			Status:  http.StatusInternalServerError,
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	jsonUsers, err := json.Marshal(users)
	if err != nil {
		log.Printf("Error: parsing from redis %s", err.Error())
	}
	err = handler.cache.Set(handler.ctx, "Users", jsonUsers, time.Minute*5).Err()

	if err != nil{
		log.Printf("Error: setting value to redis %s", err.Error())
	}

	c.JSON(http.StatusOK, types.Response{
		Status:  http.StatusOK,
		Error:   false,
		Data:    users,
	})
}