package storage

import (
	"database/sql"
	"go-postgres/model"
	"go-postgres/storage/config"
	"log"
)

const (
	getAllSQL         = "SELECT * FROM users ORDER BY userid"
	getUsedByIdSQL    = "SELECT * FROM users WHERE userid=$1"
	insertSQL         = "INSERT INTO users (name, location, age) VALUES ($1, $2, $3) RETURNING userid"
	updateUserByIdSQL = "UPDATE users SET name=$2, location=$3, age=$4 WHERE userid=$1"
	deleteUserByIdSQL = "DELETE FROM users WHERE userid=$1"
)

var db *sql.DB

func init() {
	config.ConnectDb()
	db = config.GetDB()
}

func InsertUser(user model.User) int64 {
	var id int64
	row := db.QueryRow(insertSQL, user.Name, user.Location, user.Age)
	if scanError := row.Scan(&id); scanError != nil {
		log.Fatalf("Unable to execute the insertSQL. %v", scanError)
	}
	log.Printf("Single record has been inserted with id: %v", id)
	return id
}

func GetUser(id int64) (model.User, error) {
	user := model.User{}
	row := db.QueryRow(getUsedByIdSQL, id)
	scanError := row.Scan(&user.Id, &user.Name, &user.Age, &user.Location)
	switch scanError {
	case sql.ErrNoRows:
		log.Println("No rows have been returned")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan the row %v", scanError)
	}
	return user, scanError
}

func GetAllUsers() []model.User {
	var users []model.User
	rows, queryError := db.Query(getAllSQL)
	if queryError != nil {
		log.Fatalf("Unable to execute the insertSQL. %v", queryError)
	}
	for rows.Next() {
		var user model.User
		if scanError := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Location); scanError != nil {
			log.Fatalf("Unable to scan the row. %v", scanError)
		}
		users = append(users, user)
	}
	if closeError := rows.Close(); closeError != nil {
		log.Fatal("Unable to close rows.")
	}
	return users
}

func UpdateUser(id int64, user model.User) int64 {
	result, executionError := db.Exec(updateUserByIdSQL, id, user.Name, user.Location, user.Age)
	if executionError != nil {
		log.Fatalf("Unable to execute the insertSQL. %v", executionError)
	}
	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffectedError != nil {
		log.Fatalf("Error while checking the affected rows. %v", rowsAffectedError)
	}
	log.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected
}

func DeleteUser(id int64) int64 {
	result, executionError := db.Exec(deleteUserByIdSQL, id)
	if executionError != nil {
		log.Fatalf("Unable to execute the insertSQL. %v", executionError)
	}
	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffectedError != nil {
		log.Fatalf("Error while checking the affected rows. %v", rowsAffectedError)
	}
	log.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected
}
