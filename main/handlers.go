package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/andkolbe/chirper-go/internal/models"
	"github.com/gorilla/mux"
)

// A handler responds to an HTTP request

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
	// NewEncoder provides better performance than json.Unmarshal as it does not have to buffer the output into an in memory slice of bytes 
	// this reduces allocations and the overheads of the service
	// send all the users as response
	json.NewEncoder(w).Encode(users)
}

func (repo *Repository) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	// get and store any params on the request
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
	var user models.User
	
	json.NewDecoder(r.Body).Decode(&user)

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

func (repo *Repository) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user models.User
	json.Unmarshal(reqBody, &user)

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	repo.users.UpdateUser(user, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (repo *Repository) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
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
