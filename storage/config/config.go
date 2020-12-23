package config

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres golang driver
	"log"
	"os"
	"time"
)

var (
	db *sql.DB
)

func ConnectDb() {
	if db == nil {
		if environmentError := godotenv.Load(".properties"); environmentError != nil {
			log.Fatal("Error loading .properties file")
		}
		database, openDbError := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
		if openDbError != nil {
			panic(openDbError)
		}
		if pingError := database.Ping(); pingError != nil {
			panic(pingError)
		}
		log.Println("Connection has been established.")
		database.SetMaxOpenConns(5)
		database.SetMaxIdleConns(5)
		database.SetConnMaxLifetime(10 * time.Second)
		db = database
	}
}

func GetDB() *sql.DB {
	return db
}

func ShutdownDb() {
	if closeError := db.Close(); closeError != nil {
		log.Fatal("Unable to close db")
	}
}
