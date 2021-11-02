package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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

// GET /users/{id}
func (repo *Repository) GetChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	// get and store any params on the request
	vars := mux.Vars(r)
	// pull the id value out of the Vars map
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