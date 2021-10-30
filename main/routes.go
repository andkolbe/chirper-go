package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/users", usersIndex)
	mux.Get("/users/:id", usersShow)
	// http.HandleFunc("/users/create", usersCreate)

	return mux
}