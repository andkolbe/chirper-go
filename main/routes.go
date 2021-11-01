package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func routes(repo *Repository) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/users", repo.ShowAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", repo.ShowOneUserByID).Methods("GET")
	// http.HandleFunc("/users/create", usersCreate)

	return router
}