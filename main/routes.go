package main

import (
	"net/http"

	"github.com/andkolbe/chirper-go/internal/handlers"
	"github.com/gorilla/mux"
)

func routes() http.Handler {
	router := mux.NewRouter()

	// middleware
	router.Use(NoSurf)
	router.Use(SessionLoad)

	router.HandleFunc("/", handlers.Repo.Home).Methods("GET")
	router.HandleFunc("/about", handlers.Repo.About).Methods("GET")

	u := router.PathPrefix("/users").Subrouter()
	u.HandleFunc("/", handlers.Repo.GetAllUsersHandler).Methods("GET")
	u.HandleFunc("/{id}", handlers.Repo.GetUserByIDHandler).Methods("GET")
	u.HandleFunc("/", handlers.Repo.RegisterNewUserHandler).Methods("POST")
	u.HandleFunc("/{id}", handlers.Repo.UpdateUserHandler).Methods("PUT")
	u.HandleFunc("/{id}", handlers.Repo.DeleteUserHandler).Methods("DELETE")

	ac := router.PathPrefix("/api/chirps").Subrouter()
	ac.HandleFunc("/", handlers.Repo.GetAllChirpsHandler).Methods("GET")
	ac.HandleFunc("/{id}", handlers.Repo.GetChirpByIDHandler).Methods("GET")
	ac.HandleFunc("/", handlers.Repo.CreateNewChirpHandler).Methods("POST")
	ac.HandleFunc("/{id}", handlers.Repo.UpdateChirpHandler).Methods("PUT")
	ac.HandleFunc("/{id}", handlers.Repo.DeleteChirpHandler).Methods("DELETE")

	c := router.PathPrefix("/chirps").Subrouter()
	c.HandleFunc("/new", handlers.Repo.NewChirp).Methods("GET")
	c.HandleFunc("/edit", handlers.Repo.EditChirp).Methods("GET")
	c.HandleFunc("/show", handlers.Repo.ShowOneChirp).Methods("GET")



	router.HandleFunc("/login", handlers.Repo.Login).Methods("POST")

	// create a file server - a place to get static files from
	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return router
}