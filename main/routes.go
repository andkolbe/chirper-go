package main

import (
	"net/http"

	"github.com/andkolbe/chirper-go/internal/config"
	"github.com/andkolbe/chirper-go/internal/handlers"
	"github.com/gorilla/mux"
)

func routes(app *config.AppConfig) http.Handler {
	router := mux.NewRouter()

	// middleware
	router.Use(NoSurf)
	router.Use(SessionLoad)

	router.HandleFunc("/", handlers.Repo.HomePage).Methods("GET")
	router.HandleFunc("/about", handlers.Repo.AboutPage).Methods("GET")

	router.HandleFunc("/login", handlers.Repo.LoginPage).Methods("GET")
	router.HandleFunc("/login", handlers.Repo.Login).Methods("POST")
	router.HandleFunc("/logout", handlers.Repo.Logout).Methods("GET")

	u := router.PathPrefix("/users").Subrouter()
	u.HandleFunc("/", handlers.Repo.GetAllUsersHandler).Methods("GET")
	u.HandleFunc("/{id}", handlers.Repo.GetUserByIDHandler).Methods("GET")
	u.HandleFunc("/", handlers.Repo.RegisterNewUserHandler).Methods("POST")
	u.HandleFunc("/{id}", handlers.Repo.UpdateUserHandler).Methods("PUT")
	u.HandleFunc("/{id}", handlers.Repo.DeleteUserHandler).Methods("DELETE")

	ac := router.PathPrefix("/api/chirps").Subrouter()
	ac.HandleFunc("/", handlers.Repo.GetAllChirpsHandler).Methods("GET")
	ac.HandleFunc("/{id}", handlers.Repo.GetChirpByIDHandler).Methods("GET")
	ac.HandleFunc("/", handlers.Repo.PostNewChirpHandler).Methods("POST")
	ac.HandleFunc("/{id}", handlers.Repo.UpdateChirpHandler).Methods("PUT")
	ac.HandleFunc("/{id}", handlers.Repo.DeleteChirpHandler).Methods("DELETE")

	c := router.PathPrefix("/chirps").Subrouter()
	c.HandleFunc("/new", handlers.Repo.NewChirpPage).Methods("GET")
	c.HandleFunc("/new", handlers.Repo.PostNewChirpHandler).Methods("POST")
	c.HandleFunc("/edit", handlers.Repo.EditChirpPage).Methods("GET")
	c.HandleFunc("/show", handlers.Repo.ShowOneChirpPage).Methods("GET")

	c.HandleFunc("/summary", handlers.Repo.ChirpSummary).Methods("GET")

	admin := router.PathPrefix("/admin").Subrouter()
	admin.Use(Auth)
	admin.HandleFunc("/profile", handlers.Repo.ProfilePage).Methods("GET")


	// create a file server - a place to get static files from
	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return router
}