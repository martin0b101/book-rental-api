package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/martin0b101/book-rental-api/types"
	"github.com/martin0b101/book-rental-api/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler{
	return &Handler{store: store}
}

func (handler *Handler) RegisterRoutes(router *mux.Router){
	router.HandleFunc("/register", handler.registerUser).Methods("POST")
	router.HandleFunc("/users", handler.getUsers).Methods("GET")
}

func (handler *Handler) registerUser(writer http.ResponseWriter, request *http.Request){

	var payload types.User
	if err := utils.ParseJSON(request, payload); err != nil {

	}	

}

func (handler *Handler) getUsers(writer http.ResponseWriter, request *http.Request){

	users, err := handler.store.GetUsers()

	if err != nil{
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	jsonResponse, jsonError := json.Marshal(types.Response{
		Status: http.StatusOK,
		Error: false,
		Data: users,
	})

	if jsonError != nil{
		http.Error(writer, jsonError.Error(), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonResponse)
}