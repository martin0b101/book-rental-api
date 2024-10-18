package book

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin0b101/book-rental-api/types"
)



type Handler struct {
	store types.BookStore
}

func NewHandler(store types.BookStore) *Handler{
	return &Handler{store: store}
}

func (handler *Handler) RegisterRoutes(router *gin.Engine){
	router.GET("/books", handler.getAvailableBooks)
	router.POST("/book/borrow", handler.borrowBook)
	router.POST("/book/return", handler.returnBook)


}

func (handler *Handler) getAvailableBooks(c *gin.Context){

	books, err := handler.store.GetAvailableBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			Status:  http.StatusInternalServerError,
			Error:   true,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, types.Response{
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

	book, errBorrow := handler.store.BorrowBook(borrowRequest.BookId, borrowRequest.UserId)

	if errBorrow == types.NotFoundError{
		c.JSON(http.StatusNotFound, types.Response{
			Status:  http.StatusNotFound,
			Error:   true,
			Message: errBorrow.Error(),
			Data:    nil,
		})
		return
	}

	if errBorrow != nil{
		c.JSON(http.StatusInternalServerError, types.Response{
			Status:  http.StatusInternalServerError,
			Error:   true,
			Message: errBorrow.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, types.Response{
		Status:  http.StatusOK,
		Error:   false,
		Data:    book,
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
			Data:    nil,
		})
		return
	}
	
	book, errReturn := handler.store.ReturnBook(returnBookRequest.BookId, returnBookRequest.UserId)

	if errReturn != nil{
		c.JSON(http.StatusInternalServerError, types.Response{
			Status:  http.StatusInternalServerError,
			Error:   true,
			Message: errReturn.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, types.Response{
		Status:  http.StatusOK,
		Error:   false,
		Data:    book,
	})
}