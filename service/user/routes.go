package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin0b101/book-rental-api/types"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler{
	return &Handler{store: store}
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

	c.JSON(http.StatusOK, types.Response{
		Status: http.StatusOK,
		Error: false,
		Data: user,
	})
}

func (handler *Handler) getUsers(c *gin.Context){

	users, err := handler.store.GetUsers()

	if err != nil{
		c.JSON(http.StatusInternalServerError, types.Response{
			Status:  http.StatusInternalServerError,
			Error:   true,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Status:  http.StatusOK,
		Error:   false,
		Data:    users,
	})
}