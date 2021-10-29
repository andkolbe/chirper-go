package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// we can only use string and int safely because we set NOT NULL constraints on all of the columns on the table
type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Created_At time.Time `json:"created_at"`
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("unable to load .env file")
	}
}

var db *sql.DB

func init() {
	LoadEnv()
	
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASS")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("DB_NAME")

	URL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, DBNAME)

	var err error
	// Get a database handle
	db, err = sql.Open("mysql", URL)
	if err != nil {
		log.Fatal(err)
	}

	// because sql.Open doesn't actually check a connection, we also call DB.Ping() to make sure that everything works OK on startup
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to DB!")
}

func main() {
	http.HandleFunc("/users", usersIndex)
	http.HandleFunc("/users/show", usersShow)

	http.ListenAndServe(":3000", nil)
}

func usersIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Created_At)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		fmt.Printf("%q, %s, %s, %s", user.ID, user.Name, user.Email, user.Password)
	}
}

func usersShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return 
	}

	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	user := new(User)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Created_At)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return 
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return 
	}

	fmt.Fprintf(w, "%q, %s, %s", user.ID, user.Name, user.Email)
}