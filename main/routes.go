package main

import (
	"net/http"

	"github.com/andkolbe/chirper-go/internal/handlers"
	"github.com/gorilla/mux"
)

func routes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.Repo.Home).Methods("GET")
	router.HandleFunc("/about", handlers.Repo.About).Methods("GET")

	u := router.PathPrefix("/users").Subrouter()
	u.HandleFunc("/", handlers.Repo.GetAllUsersHandler).Methods("GET")
	u.HandleFunc("/{id}", handlers.Repo.GetUserByIDHandler).Methods("GET")
	u.HandleFunc("/", handlers.Repo.RegisterNewUserHandler).Methods("POST")
	u.HandleFunc("/{id}", handlers.Repo.UpdateUserHandler).Methods("PUT")
	u.HandleFunc("/{id}", handlers.Repo.DeleteUserHandler).Methods("DELETE")

	c := router.PathPrefix("/chirps").Subrouter()
	c.HandleFunc("/", handlers.Repo.GetAllChirpsHandler).Methods("GET")
	c.HandleFunc("/{id}", handlers.Repo.GetChirpByIDHandler).Methods("GET")
	c.HandleFunc("/", handlers.Repo.CreateNewChirpHandler).Methods("POST")
	c.HandleFunc("/{id}", handlers.Repo.UpdateChirpHandler).Methods("PUT")
	c.HandleFunc("/{id}", handlers.Repo.DeleteChirpHandler).Methods("DELETE")

	router.HandleFunc("/login", handlers.Repo.Login).Methods("POST")

	return router
}