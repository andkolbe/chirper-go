package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func routes(env *Env) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/users", env.usersIndex)
	// mux.Get("/users/:id", usersShow)
	// http.HandleFunc("/users/create", usersCreate)

	return mux
}