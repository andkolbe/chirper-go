package main

import (
	"encoding/json"
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
func (repo *Repository) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := repo.users.GetAllUsers()
	if err != nil {
        log.Println(err)
        http.Error(w, http.StatusText(500), 500)
        return
    }

	// set the response header to send back json
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(users)
}

func (repo *Repository) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	// get and store any params on the request in a variable 
	vars := mux.Vars(r)
	// pull the id out of the Vars map
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

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)
}

func (repo *Repository) CreateNewUserHandler(w http.ResponseWriter, r *http.Request) {

	// user := models.User{
	// 	Name: r.FormValue("name"), // r.Form.Get("first_name") matches the name="first_name" field on the html page
	// 	Email:     r.FormValue("email"),
	// 	Password:  r.FormValue("password"),
	// }

	// if user.Name == "" || user.Email == "" || user.Password == "" {
	// 	http.Error(w, http.StatusText(400), 400)
	// 	return
	// }

	// user, err := repo.users.PostNewUser(user)
	// if err != nil {
    //     log.Println(err)
    //     http.Error(w, http.StatusText(500), 500)
    //     return
    // }

	// fmt.Fprintf(w, "%q, %s, %s", user.ID, user.Name, user.Email)

	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	// 	http.Error(w, http.StatusText(500), 500)
	// 	return
	// }
	// fmt.Fprintf(w, "User created successfully! (%d row affected)\n", rowsAffected)
}

func (repo *Repository) UpdateUserHandler (w http.ResponseWriter, r *http.Request) {

}

func (repo *Repository) DeleteUserHandler (w http.ResponseWriter, r *http.Request) {
	
}