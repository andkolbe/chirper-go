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

// Create a custom Env struct which holds a connection pool.
// All the dependencies for our handlers are explicitly defined in one place
type Env struct {
    users models.UserModel
}

func main() {
	env.LoadEnv()
	PORT := os.Getenv("PORT")
	URL := os.Getenv("URL")
	if PORT == "" || URL == "" {
		log.Fatal("env variables are not set")
	}

	var err error

	// connect to db
	db, err := driver.DBConnect(URL)
	if err != nil {
		log.Fatal(err)
	}
	
	// Initalise Env with a models.UserModel instance (which in turn wraps the connection pool).
	// we can see what values they have at runtime by simply looking at how it is initialised
    env := &Env{
		users: models.UserModel{DB: db},
	}

	mux := routes(env)

	log.Println("Starting web server")

	http.ListenAndServe(":"+PORT, mux)
}

// sends a HTTP response listing all users
func (env *Env) usersIndex(w http.ResponseWriter, r *http.Request) {
	users, err := env.users.GetAllUsers()
	if err != nil {
        log.Println(err)
        http.Error(w, http.StatusText(500), 500)
        return
    }

	for _, user := range users {
		fmt.Printf("%q, %s, %s, %s", user.ID, user.Name, user.Email, user.Password)
	}
}

// func usersShow(w http.ResponseWriter, r *http.Request) {
// 	id := r.FormValue("id")
// 	if id == "" {
// 		http.Error(w, http.StatusText(400), 400)
// 		return 
// 	}

// 	user, err := models.GetUserByID(id)
// 	if err != nil {
//         log.Println(err)
//         http.Error(w, http.StatusText(500), 500)
//         return
//     }

// 	fmt.Fprintf(w, "%q, %s, %s", user.ID, user.Name, user.Email)
// }

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
