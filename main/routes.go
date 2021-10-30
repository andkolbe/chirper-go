package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func routes(repo *Repository) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/users", repo.usersIndex)
	// mux.Get("/users/:id", usersShow)
	// http.HandleFunc("/users/create", usersCreate)

	return mux
}