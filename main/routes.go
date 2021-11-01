package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func routes(repo *Repository) http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/users", repo.usersIndex).Methods("GET")
	mux.HandleFunc("/users/{id}", repo.usersShow).Methods("GET")
	// http.HandleFunc("/users/create", usersCreate)

	return mux
}