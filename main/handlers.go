package main

import (
	"encoding/json"
	"io/ioutil"
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
	// pull the id value out of the Vars map
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
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
        log.Println(err)
        http.Error(w, http.StatusText(500), 500)
        return
    }

	var user models.User
	// unmarshal the JSON data retrieved from the body into the new user instance 
	json.Unmarshal(reqBody, &user)

	repo.users.CreateNewUser(user)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
	w.Write([]byte("user created"))



	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	// 	http.Error(w, http.StatusText(500), 500)
	// 	return
	// }
	// fmt.Fprintf(w, "User created successfully! (%d row affected)\n", rowsAffected)
}

// func (repo *Repository) UpdateUserHandler (w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// id := vars["id"]
// 	if id == "" {
// 		http.Error(w, http.StatusText(400), 400)
// 		return 
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(user)
// }

func (repo *Repository) DeleteUserHandler (w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return 
	}

	repo.users.DeleteUser(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user deleted"))
}