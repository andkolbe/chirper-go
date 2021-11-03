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

	// create an instance of the AppConfig that different parts of the app can use
	var app config.AppConfig

	// initialize the template cache when the application starts
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot crete template cache")
	}
	// store the template cache in an instance of the AppConfig
	app.TemplateCache = templateCache

	// lets us change the templates in development without having to restart the server every time
	app.UseCache = false

	// gives our render package access to everything inside of AppConfig  (it needs the template cache we initialized here)
	render.NewTemplates(&app)

	// connect to db
	db, err := driver.DBConnect(URL)
	if err != nil {
		log.Fatal(err)
	}
	
	// Initalise an instance of Repository with a models.UserModel instance (which in turn wraps the connection pool)
	repo := handlers.NewRepo(models.DBModel{DB: db}, &app)
	// gives our handlers package access to everything inside of AppConfig 
	handlers.NewHandlers(repo)

	mux := routes()

	http.ListenAndServe("127.0.0.1:"+PORT, mux)
}