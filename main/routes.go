package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func routes(repo *Repository) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/users", repo.GetAllUsersHandler).Methods("GET")
	router.HandleFunc("/users/{id}", repo.GetUserByIDHandler).Methods("GET")
	router.HandleFunc("/users", repo.CreateNewUserHandler).Methods("POST")
	// router.HandleFunc("/users/{id}", repo.UpdateUserHandler).Methods("PUT")
	router.HandleFunc("/users/{id}", repo.DeleteUserHandler).Methods("DELETE")

	return router
}