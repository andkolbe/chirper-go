package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/chirper-go/internal/config"
	"github.com/andkolbe/chirper-go/internal/driver"
	"github.com/andkolbe/chirper-go/internal/env"
	"github.com/andkolbe/chirper-go/internal/handlers"
	"github.com/andkolbe/chirper-go/internal/models"
	"github.com/andkolbe/chirper-go/internal/render"
)

// create an instance of the AppConfig that different parts of the app can use
var app config.AppConfig

var session *scs.SessionManager

func main() {
	env.LoadEnv()
	PORT := os.Getenv("PORT")
	URL := os.Getenv("URL")
	if PORT == "" || URL == "" {
		log.Fatal("env variables are not set")
	}

	// change this to true when in production
	app.InProduction = false

	// initialize session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // persist the session even if the browser window is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // use https

	// make our scs.SessionManager instance variable available in the AppConfig so it can be used anywhere else in the application
	app.Session = session

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