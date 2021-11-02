package main

import (
	"net/http"

	"github.com/andkolbe/chirper-go/internal/handlers"
	"github.com/gorilla/mux"
)

func routes(repo *handlers.Repository) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/users", repo.GetAllUsersHandler).Methods("GET")
	router.HandleFunc("/users/{id}", repo.GetUserByIDHandler).Methods("GET")
	router.HandleFunc("/users", repo.CreateNewUserHandler).Methods("POST")
	router.HandleFunc("/users/{id}", repo.UpdateUserHandler).Methods("PUT")
	router.HandleFunc("/users/{id}", repo.DeleteUserHandler).Methods("DELETE")

	router.HandleFunc("/chirps", repo.GetAllChirpsHandler).Methods("GET")
	router.HandleFunc("/chirps/{id}", repo.GetChirpByIDHandler).Methods("GET")
	router.HandleFunc("/chirps", repo.CreateNewChirpHandler).Methods("POST")
	router.HandleFunc("/chirps/{id}", repo.UpdateChirpHandler).Methods("PUT")


	return router
}