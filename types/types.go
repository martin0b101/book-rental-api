package types

import (
	"errors"
	"time"

)

type UserStore interface {
	GetUsers() ([]User, error)
	CreateUser(RegisterUserRequest) (*User, error)
}

type BookStore interface {
	GetAvailableBooks() ([]Book, error)
	BorrowBook(bookId int, userId int) (*Book, error)
	ReturnBook(bookId int, userId int) (*Book, error)
}

type User struct {
	Id int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}

type Book struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Quantity int `json:"quantity"`
}

type Borrow struct {
	Id *int `json:"id"`
	UserId int `json:"user_id"`
	BookId int `json:"book_id"`
	BorrowedQuantity int `json:"borrowed_quantity"`
	BorrowedAt time.Time `json:"borrowed_at"`
	ReturnedAt *time.Time `json:"returned_at"`
}


type BookActionRequest struct {
	BookId int `json:"book_id"`
	UserId int `json:"user_id"`
}

type RegisterUserRequest struct{
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}


type Response struct {
	Status int `json:"status"`
	Error bool `json:"error"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}


// Errors
var NotFoundError = errors.New("Not found")