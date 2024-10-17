package book

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/martin0b101/book-rental-api/types"
)



type Handler struct {
	store types.BookStore
}

func NewHandler(store types.BookStore) *Handler{
	return &Handler{store: store}
}

func (handler *Handler) RegisterRoutes(router *mux.Router){
	router.HandleFunc("/books", handler.getAvailableBooks).Methods("GET")
	router.HandleFunc("/book/borrow", handler.borrowBook).Methods("POST")
	router.HandleFunc("/book/return", handler.returnBook).Methods("POST")


}

func (handler *Handler) getAvailableBooks(writer http.ResponseWriter, request *http.Request){

	books, err := handler.store.GetAvailableBooks()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	jsonResponse, jsonError := json.Marshal(types.Response{
		Status: http.StatusOK,
		Error: false,
		Data: books,
	})

	if jsonError != nil{
		http.Error(writer, jsonError.Error(), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonResponse)
}



func (handler *Handler) borrowBook(writer http.ResponseWriter, request *http.Request){
	var borrowRequest types.BookActionRequest
	err := json.NewDecoder(request.Body).Decode(&borrowRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	book, errBorrow :=handler.store.BorrowBook(borrowRequest.BookId, borrowRequest.UserId)

	if errBorrow != nil{
		http.Error(writer, errBorrow.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, jsonError := json.Marshal(types.Response{
		Status: http.StatusOK,
		Error: false,
		Data: book,
	})

	if jsonError != nil{
		http.Error(writer, jsonError.Error(), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonResponse)
}

func (handler *Handler) returnBook(writer http.ResponseWriter, request *http.Request){

	var returnBookRequest types.BookActionRequest
	err := json.NewDecoder(request.Body).Decode(&returnBookRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	
	book, errReturn := handler.store.ReturnBook(returnBookRequest.BookId, returnBookRequest.UserId)

	if errReturn != nil{
		http.Error(writer, errReturn.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, jsonError := json.Marshal(types.Response{
		Status: http.StatusOK,
		Error: false,
		Data: book,
	})

	if jsonError != nil{
		http.Error(writer, jsonError.Error(), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonResponse)
}