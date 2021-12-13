package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-postgres/model"
	"go-postgres/storage"
	"log"
	"net/http"
	"strconv"
)

func Create(responseWriter http.ResponseWriter, request *http.Request) {
	var user model.User
	if decodeError := json.NewDecoder(request.Body).Decode(&user); decodeError != nil {
		log.Fatalf("Unable to decode the request body.  %v", decodeError)
	}
	userId := storage.InsertUser(user)
	response := model.Response{
		Id:      userId,
		Message: "User has been created successfully",
	}
	if encodeError := json.NewEncoder(responseWriter).Encode(response); encodeError != nil {
		log.Fatalf("Unable to encode the response. %v", encodeError)
	}
}

func Get(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, toIntError := strconv.Atoi(params["id"])
	if toIntError != nil {
		log.Fatalf("Unable to convert the string into int.  %v", toIntError)
	}
	user, getUserError := storage.GetUser(int64(id))
	if getUserError != nil {
		log.Fatalf("Unable to get user. %v", getUserError)
	}
	if encodeError := json.NewEncoder(responseWriter).Encode(user); encodeError != nil {
		log.Fatalf("Unable to encode the response. %v", encodeError)
	}
}

func GetAll(responseWriter http.ResponseWriter, _ *http.Request) {
	users := storage.GetAllUsers()
	if len(users) > 0 {
		encodeError := json.NewEncoder(responseWriter).Encode(users)
		if encodeError != nil {
			log.Fatalf("Unable to encode the response. %v", encodeError)
		}
	} else {
		_ = json.NewEncoder(responseWriter).Encode([]string{})
	}

}

func Update(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, toIntError := strconv.Atoi(params["id"])
	if toIntError != nil {
		log.Fatalf("Unable to convert the string into int.  %v", toIntError)
	}
	var user model.User
	if decodeError := json.NewDecoder(request.Body).Decode(&user); decodeError != nil {
		log.Fatalf("Unable to decode the request body.  %v", decodeError)
	}
	updatedRows := storage.UpdateUser(int64(id), user)
	response := model.Response{
		Id:      int64(id),
		Message: fmt.Sprintf("User updated successfully. Total rows/records affected %v", updatedRows),
	}
	if encodeError := json.NewEncoder(responseWriter).Encode(response); encodeError != nil {
		log.Fatalf("Unable to encode the response. %v", encodeError)
	}
}

func Delete(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, toIntError := strconv.Atoi(params["id"])
	if toIntError != nil {
		log.Fatalf("Unable to convert the string into int.  %v", toIntError)
	}
	deletedRows := storage.DeleteUser(int64(id))
	response := model.Response{
		Id:      int64(id),
		Message: fmt.Sprintf("User updated successfully. Total rows/records affected %v", deletedRows),
	}
	if encodeError := json.NewEncoder(responseWriter).Encode(response); encodeError != nil {
		log.Fatalf("Unable to encode the response. %v", encodeError)
	}
}
