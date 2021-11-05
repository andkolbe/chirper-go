package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andkolbe/chirper-go/internal/config/helpers"
	"github.com/andkolbe/chirper-go/internal/forms"
	"github.com/andkolbe/chirper-go/internal/models"
	"github.com/andkolbe/chirper-go/internal/render"
	"github.com/gorilla/mux"
)

// GET /chirps
func (repo *Repository) GetAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := repo.dbmodel.GetAllChirps()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chirps)
}

// GET /chirps/{id}
func (repo *Repository) GetChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	chirp, err := repo.dbmodel.GetChirpByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chirp)
}

// POST /chirps
func (repo *Repository) PostNewChirpHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	chirp := models.Chirp{
		UserID:   r.Form.Get("userid"),
		Content:  r.Form.Get("content"),
		Location: r.Form.Get("location"),
	}

	// pass the post form data into the empty form that is initialized on the handler
	form := forms.New(r.PostForm)

	// checks on form fields we created the forms package
	form.Required("userid", "content", "location")
	form.MinLength("content", 5)

	if !form.Valid() {
		// data comes from models.templatedata
		data := make(map[string]interface{})
		data["chirp"] = chirp

		// reload the page with the form and data passed to the form
		render.Template(w, r, "new_chirp.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// add chirp into session
	repo.App.Session.Put(r.Context(), "chirp", chirp)

	// any time your site receives a POST request, you should direct users to another page with a HTTP redirect so they can't accidently click on the submit twice
	http.Redirect(w, r, "/chirps/summary", http.StatusSeeOther)

	// var chirp models.Chirp

	// json.NewDecoder(r.Body).Decode(&chirp)

	// insertID := repo.dbmodel.CreateNewChirp(chirp)

	// res := response {
	// 	ID: insertID,
	// 	Message: "Chirp created successfully!",
	// }

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(res)

}

// PUT /chirps/{id}
func (repo *Repository) UpdateChirpHandler(w http.ResponseWriter, r *http.Request) {
	var chirp models.Chirp
	json.NewDecoder(r.Body).Decode(&chirp)

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	updatedRows := repo.dbmodel.UpdateChirp(chirp, id)

	msg := fmt.Sprintf("Chirp updated successfully. Total rows affected %v", updatedRows)
	intID, _ := strconv.Atoi(id)
	res := response{
		ID:      int64(intID),
		Message: msg,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// DELETE /users/{id}
func (repo *Repository) DeleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	deletedRows := repo.dbmodel.DeleteChirp(id)

	msg := fmt.Sprintf("Chirp deleted successfully. Total rows affected %v", deletedRows)
	intID, _ := strconv.Atoi(id)

	res := response{
		ID:      int64(intID),
		Message: msg,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
