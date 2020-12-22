package router

import (
	"go-postgres/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/user/{id}", controller.Get).Methods("GET")
	router.HandleFunc("/api/user", controller.GetAll).Methods("GET")
	router.HandleFunc("/api/user", controller.Create).Methods("POST")
	router.HandleFunc("/api/user/{id}", controller.Update).Methods("PUT")
	router.HandleFunc("/api/user/{id}", controller.Delete).Methods("DELETE")
	return router
}
