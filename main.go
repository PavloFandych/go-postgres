package main

import (
	"go-postgres/router"
	"go-postgres/storage/config"
	"log"
	"net/http"
)

func main() {
	defer config.ShutdownDb()
	log.Println("Starting server on the port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router.Router()))
}
