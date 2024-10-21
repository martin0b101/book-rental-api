package book

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
	store types.BookStore
	cache *redis.Client
	ctx context.Context
}

func NewHandler(store types.BookStore, cache *redis.Client, ctx context.Context) *Handler{
	return &Handler{
		store: store,
		cache: cache,
		ctx: ctx,
	}
}

func (handler *Handler) RegisterRoutes(router *gin.Engine){
	router.GET("/books", handler.getAvailableBooks)
	router.POST("/book/borrow", handler.borrowBook)
	router.POST("/book/return", handler.returnBook)
}

func (handler *Handler) getAvailableBooks(c *gin.Context){

	value, err := handler.cache.Get(handler.ctx, "Books").Result()

	if err == nil {
		var books []types.Book
		json.Unmarshal([]byte(value), &books)
		c.JSON(http.StatusOK, types.Response{
			Status:  http.StatusOK,
			Error:   false,
			Data:    books,
		})
		return
	}

	books, err := handler.store.GetAvailableBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Status:  http.StatusInternalServerError,
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	jsonBooks, err := json.Marshal(books)
	if err != nil {
		log.Printf("Error: parsing from redis %s", err.Error())
	}
	err = handler.cache.Set(handler.ctx, "Books", jsonBooks, time.Minute*5).Err()
	if err != nil {
		log.Printf("Error: setting value to redis %s", err.Error())
	}

	c.JSON(http.StatusOK, types.Response{
		Status:  http.StatusOK,
		Error:   false,
		Data:    books,
	})
}



func (handler *Handler) borrowBook(c *gin.Context){

	var borrowRequest types.BookActionRequest
	err := c.BindJSON(&borrowRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status:  http.StatusBadRequest,
			Error:   true,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	errBorrow := handler.store.BorrowBook(borrowRequest.BookId, borrowRequest.UserId)

	if errBorrow == types.NotFoundError{
		c.JSON(http.StatusNotFound, types.Response{
			Status:  http.StatusNotFound,
			Error:   true,
			Message: errBorrow.Error(),
		})
		return
	}

	if errBorrow != nil{
		c.JSON(http.StatusInternalServerError, types.Response{
			Status:  http.StatusInternalServerError,
			Error:   true,
			Message: errBorrow.Error(),
		})
		return
	}

	utils.DeleteFromRedis(handler.ctx, handler.cache, "Books") 

	c.JSON(http.StatusOK, types.Response{
		Status:  http.StatusOK,
		Error:   false,
		Message: "Book successfully rented.",
	})
}


func (handler *Handler) returnBook(c *gin.Context){


	var returnBookRequest types.BookActionRequest
	err := c.BindJSON(&returnBookRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status:  http.StatusBadRequest,
			Error:   true,
			Message: err.Error(),
		})
		return
	}
	
	errReturn := handler.store.ReturnBook(returnBookRequest.BookId, returnBookRequest.UserId)

	if errReturn != nil{
		c.JSON(http.StatusInternalServerError, types.Response{
			Status:  http.StatusInternalServerError,
			Error:   true,
			Message: errReturn.Error(),
		})
		return
	}

	utils.DeleteFromRedis(handler.ctx, handler.cache, "Books") 

	c.JSON(http.StatusOK, types.Response{
		Status:  http.StatusOK,
		Error:   false,
		Message: "Book successfully returned.",
	})
}