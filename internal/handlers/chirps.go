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

// GET /chirps
func (repo *Repository) GetAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := repo.dbmodel.GetAllChirps()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
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
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chirp)
}

// POST /chirps
func (repo *Repository) CreateNewChirpHandler(w http.ResponseWriter, r *http.Request) {
	var chirp models.Chirp
	
	json.NewDecoder(r.Body).Decode(&chirp)

	insertID := repo.dbmodel.CreateNewChirp(chirp)

	res := response {
		ID: insertID,
		Message: "Chirp created successfully!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
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