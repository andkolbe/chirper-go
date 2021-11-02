package main

import (
	"log"
	"net/http"
	"os"

	"github.com/andkolbe/chirper-go/internal/driver"
	"github.com/andkolbe/chirper-go/internal/env"
	"github.com/andkolbe/chirper-go/internal/models"
)

func main() {
	env.LoadEnv()
	PORT := os.Getenv("PORT")
	URL := os.Getenv("URL")
	if PORT == "" || URL == "" {
		log.Fatal("env variables are not set")
	}

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

	http.ListenAndServe("127.0.0.1:"+PORT, mux)
}