package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andkolbe/chirper-go/internal/models"
	"github.com/gorilla/mux"
)

// All the dependencies for our handlers are explicitly defined in one place
type Repository struct {
    users models.UserModel
}

// sends a HTTP response listing all users
func (repo *Repository) ShowAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repo.users.GetAllUsers()
	if err != nil {
        log.Println(err)
        http.Error(w, http.StatusText(500), 500)
        return
    }

	for _, user := range users {
		fmt.Fprintf(w, "%s, %s", user.Name, user.Email)
	}
}

func (repo *Repository) ShowOneUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return 
	}

	user, err := repo.users.GetUserByID(id)
	if err != nil {
        log.Println(err)
        http.Error(w, http.StatusText(500), 500)
        return
    }

	fmt.Fprintf(w, "%s, %s", user.Name, user.Email)
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