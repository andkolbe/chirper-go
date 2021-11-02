package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/andkolbe/chirper-go/internal/models"
	"github.com/gorilla/mux"
)

// GET /users
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
	// send all the users as json in the response
	// the http.ResponseWriter is an io.Writer
	json.NewEncoder(w).Encode(users)
}

// GET /users/{id}
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

// POST /users
func (repo *Repository) CreateNewUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// take the data from the request body and convert it into the models.User instance variable 
	// the http.Request body is an io.Reader
	json.NewDecoder(r.Body).Decode(&user)

	insertID := repo.users.CreateNewUser(user)

	res := response {
		ID: insertID,
		Message: "User created successfully!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}

// PUT /users/{id}
func (repo *Repository) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)


	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	updatedRows := repo.users.UpdateUser(user, id)

	msg := fmt.Sprintf("User updated successfully. Total rows affected %v", updatedRows)
	intID, _ := strconv.Atoi(id)
	res := response{
        ID:      int64(intID),
        Message: msg,
    }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// DELETE /users/{id}
func (repo *Repository) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	deletedRows := repo.users.DeleteUser(id)

	msg := fmt.Sprintf("User deleted successfully. Total rows affected %v", deletedRows)
	intID, _ := strconv.Atoi(id)

	res := response{
        ID:      int64(intID),
        Message: msg,
    }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
