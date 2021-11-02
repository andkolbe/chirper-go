package handlers

import (
	"encoding/json"
	"log"
	"net/http"
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