package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/andkolbe/chirper-go/internal/config/helpers"
	"github.com/andkolbe/chirper-go/internal/forms"
	"github.com/andkolbe/chirper-go/internal/models"
	"github.com/andkolbe/chirper-go/internal/render"
	"github.com/gorilla/mux"
)

// GET /users
func (repo *Repository) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := repo.dbmodel.GetAllUsers()
	if err != nil {
		helpers.ServerError(w, err)		
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
	// read and store any variable specified on the route attached to the request
	vars := mux.Vars(r)
	// pull the id value out of the Vars map
	id := vars["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	user, err := repo.dbmodel.GetUserByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// POST /users
func (repo *Repository) RegisterNewUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// take the data from the request body and convert it into the models.User instance variable 
	// the http.Request body is an io.Reader
	json.NewDecoder(r.Body).Decode(&user)

	insertID := repo.dbmodel.RegisterNewUser(user)

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

	updatedRows := repo.dbmodel.UpdateUser(user, id)

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

	deletedRows := repo.dbmodel.DeleteUser(id)

	msg := fmt.Sprintf("User deleted successfully. Total rows affected %v", deletedRows)
	intID, _ := strconv.Atoi(id)

	res := response{
        ID:      int64(intID),
        Message: msg,
    }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Login
func (repo *Repository) Login(w http.ResponseWriter, r *http.Request) {
	// renew the session on every login/logout to prevent session fixation attack
	_ = repo.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	// make sure our form has the necessary parameters
	form := forms.New(r.PostForm)
	// use the server side validation we wrote in the forms package
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "login.page.html", &models.TemplateData{
			Form: form,
		})
		return 
	}

	// var user models.User
	// json.NewDecoder(r.Body).Decode(&user)

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	id, err := repo.dbmodel.AuthenticateUser(email, password)
	if err != nil {
		log.Println(err)
		repo.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return 
	}

	// res := response {
	// 	ID: int64(id),
	// 	Message: "Logged in successfully!",
	// }
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(res)

	// put the user in the session
	repo.App.Session.Put(r.Context(), "user_id", id)
	repo.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}