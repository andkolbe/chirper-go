package main

import (
	"log"
	"net/http"
	"os"

	"github.com/andkolbe/chirper-go/internal/config"
	"github.com/andkolbe/chirper-go/internal/driver"
	"github.com/andkolbe/chirper-go/internal/env"
	"github.com/andkolbe/chirper-go/internal/handlers"
	"github.com/andkolbe/chirper-go/internal/models"
	"github.com/andkolbe/chirper-go/internal/render"
)

func main() {
	env.LoadEnv()
	PORT := os.Getenv("PORT")
	URL := os.Getenv("URL")
	if PORT == "" || URL == "" {
		log.Fatal("env variables are not set")
	}

	var app config.AppConfig

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot crete template cache")
	}
	// store the template cache in the app
	app.TemplateCache = templateCache

	// connect to db
	db, err := driver.DBConnect(URL)
	if err != nil {
		log.Fatal(err)
	}
	
	// Initalise an instance of Repository with a models.UserModel instance (which in turn wraps the connection pool)
	repo := handlers.NewRepo(models.DBModel{DB: db})

	mux := routes(repo)

	http.ListenAndServe("127.0.0.1:"+PORT, mux)
}