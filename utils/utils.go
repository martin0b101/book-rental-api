package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseJSON(request *http.Request, payload any) error {
	if request.Body == nil{
		fmt.Println("missing payload")
	}
	return json.NewDecoder(request.Body).Decode(payload)
}


// func WriteJSON(writer http.ResponseWriter, )