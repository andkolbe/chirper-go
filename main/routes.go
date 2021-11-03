package main

import (
	"net/http"

	"github.com/andkolbe/chirper-go/internal/handlers"
	"github.com/gorilla/mux"
)

func routes(repo *handlers.Repository) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", repo.Home).Methods("GET")
	router.HandleFunc("/about", repo.About).Methods("GET")

	u := router.PathPrefix("/users").Subrouter()
	u.HandleFunc("/", repo.GetAllUsersHandler).Methods("GET")
	u.HandleFunc("/{id}", repo.GetUserByIDHandler).Methods("GET")
	u.HandleFunc("/", repo.RegisterNewUserHandler).Methods("POST")
	u.HandleFunc("/{id}", repo.UpdateUserHandler).Methods("PUT")
	u.HandleFunc("/{id}", repo.DeleteUserHandler).Methods("DELETE")

	c := router.PathPrefix("/chirps").Subrouter()
	c.HandleFunc("/", repo.GetAllChirpsHandler).Methods("GET")
	c.HandleFunc("/{id}", repo.GetChirpByIDHandler).Methods("GET")
	c.HandleFunc("/", repo.CreateNewChirpHandler).Methods("POST")
	c.HandleFunc("/{id}", repo.UpdateChirpHandler).Methods("PUT")
	c.HandleFunc("/{id}", repo.DeleteChirpHandler).Methods("DELETE")

	router.HandleFunc("/login", repo.Login).Methods("POST")

	return router
}