package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andkolbe/chirper-go/driver"
	"github.com/andkolbe/chirper-go/env"
	"github.com/andkolbe/chirper-go/models"
)

func main() {
	env.LoadEnv()
	
	var err error

	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASS")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("DB_NAME")

	// connect to db
	URL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, DBNAME)
	models.DB, err = driver.DBConnect(URL)
	if err != nil {
		log.Fatal(err)
	}
	
	http.HandleFunc("/users", usersIndex)
	http.HandleFunc("/users/:id", usersShow)
	// http.HandleFunc("/users/create", usersCreate)

	http.ListenAndServe(":3000", nil)
}

func usersIndex(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAllUsers()
	if err != nil {
        log.Println(err)
        http.Error(w, http.StatusText(500), 500)
        return
    }

	for _, user := range users {
		fmt.Printf("%q, %s, %s, %s", user.ID, user.Name, user.Email, user.Password)
	}
}

func usersShow(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return 
	}

	user, err := models.GetUserByID(id)
	if err != nil {
        log.Println(err)
        http.Error(w, http.StatusText(500), 500)
        return
    }

	fmt.Fprintf(w, "%q, %s, %s", user.ID, user.Name, user.Email)
}

// func usersCreate(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
// 		return
// 	}

// 	user := models.User{
// 		Name: r.FormValue("name"), // r.Form.Get("first_name") matches the name="first_name" field on the html page
// 		Email:     r.FormValue("email"),
// 		Password:  r.FormValue("password"),
// 	}

// 	if user.Name == "" || user.Email == "" || user.Password == "" {
// 		http.Error(w, http.StatusText(400), 400)
// 		return
// 	}

// 	user, err := models.PostNewUser(user)
// 	if err != nil {
//         log.Println(err)
//         http.Error(w, http.StatusText(500), 500)
//         return
//     }

// 	fmt.Fprintf(w, "%q, %s, %s", user.ID, user.Name, user.Email)

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		http.Error(w, http.StatusText(500), 500)
// 		return
// 	}
// 	fmt.Fprintf(w, "User created successfully! (%d row affected)\n", rowsAffected)
// }
