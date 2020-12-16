package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"github.com/gorilla/mux" // used to get the params from the route
	"go-postgres/models"     // models package where User schema is defined
	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"strconv"  // package used to covert string into int type

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	return db
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var user models.User
	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	insertID := insertUser(user)
	res := response{
		ID:      insertID,
		Message: "User created successfully",
	}
	_ = json.NewEncoder(w).Encode(res)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get the userid from the request params, key is "id"
	params := mux.Vars(r)
	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	user, err := getUser(int64(id))
	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}
	// send the response
	_ = json.NewEncoder(w).Encode(user)
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the users in the db
	users, err := getAllUsers()
	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}
	// send all the users as response
	_ = json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// get the userid from the request params, key is "id"
	params := mux.Vars(r)
	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	var user models.User
	// decode the json request to user
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	updatedRows := updateUser(int64(id), user)
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", updatedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	// send the response
	_ = json.NewEncoder(w).Encode(res)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// get the userid from the request params, key is "id"
	params := mux.Vars(r)
	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	deletedRows := deleteUser(int64(id))
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", deletedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	// send the response
	_ = json.NewEncoder(w).Encode(res)
}

func insertUser(user models.User) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO users (name, location, age) VALUES ($1, $2, $3) RETURNING userid`
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, user.Name, user.Location, user.Age).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	fmt.Printf("Inserted a single record %v", id)
	return id
}

func getUser(id int64) (models.User, error) {
	db := createConnection()
	defer db.Close()
	var user models.User
	sqlStatement := `SELECT * FROM users WHERE userid=$1`
	row := db.QueryRow(sqlStatement, id)
	// unmarshal the row object to user
	err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Location)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}
	return user, err
}

func getAllUsers() ([]models.User, error) {
	db := createConnection()
	defer db.Close()
	var users []models.User
	sqlStatement := `SELECT * FROM users`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		// unmarshal the row object to user
		err = rows.Scan(&user.ID, &user.Name, &user.Age, &user.Location)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		// append the user in the users slice
		users = append(users, user)
	}
	return users, err
}

func updateUser(id int64, user models.User) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `UPDATE users SET name=$2, location=$3, age=$4 WHERE userid=$1`
	res, err := db.Exec(sqlStatement, id, user.Name, user.Location, user.Age)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected
}

func deleteUser(id int64) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `DELETE FROM users WHERE userid=$1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected
}
