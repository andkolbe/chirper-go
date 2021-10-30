package main

import (
	"log"
	"net/http"
	"os"

	"github.com/andkolbe/chirper-go/internal/driver"
	"github.com/andkolbe/chirper-go/internal/env"
	"github.com/andkolbe/chirper-go/internal/models"
)

// All the dependencies for our handlers are explicitly defined in one place
type Repository struct {
    users models.UserModel
}

func main() {
	env.LoadEnv()
	PORT := os.Getenv("PORT")
	URL := os.Getenv("URL")
	if PORT == "" || URL == "" {
		log.Fatal("env variables are not set")
	}

	var err error

	// connect to db
	db, err := driver.DBConnect(URL)
	if err != nil {
		log.Fatal(err)
	}
	
	// Initalise Repository with a models.UserModel instance (which in turn wraps the connection pool)
    repo := &Repository{
		users: models.UserModel{DB: db},
	}

	mux := routes(repo)

	log.Println("Starting web server")

	http.ListenAndServe(":"+PORT, mux)
}

