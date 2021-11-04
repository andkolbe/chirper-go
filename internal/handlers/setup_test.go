package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andkolbe/chirper-go/internal/config"
	"github.com/andkolbe/chirper-go/internal/driver"
	"github.com/andkolbe/chirper-go/internal/models"
	"github.com/andkolbe/chirper-go/internal/render"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// we have to import everything our handlers need to work to be able to test them

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	LoadEnv()
	PORT := os.Getenv("PORT")
	URL := os.Getenv("URL")
	if PORT == "" || URL == "" {
		log.Fatal("env variables are not set")
	}

	gob.Register(models.Chirp{})

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true 
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction 

	app.Session = session

	// use CreateTestTemplateCache()
	templateCache, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = templateCache

	// set to true so the cache isn't rebuilt on every request 
	// if it was false, it would call CreateTemplateCache() and use the wrong pathToTemplates
	app.UseCache = true

	render.NewTemplates(&app)

	db, err := driver.DBConnect(URL)
	if err != nil {
		log.Fatal(err)
	}
	
	repo := NewRepo(models.DBModel{DB: db}, &app)
	NewHandlers(repo)

	router := mux.NewRouter()

	router.Use(SessionLoad)

	router.HandleFunc("/", Repo.HomePage).Methods("GET")
	router.HandleFunc("/about", Repo.AboutPage).Methods("GET")

	u := router.PathPrefix("/users").Subrouter()
	u.HandleFunc("/", Repo.GetAllUsersHandler).Methods("GET")
	u.HandleFunc("/{id}", Repo.GetUserByIDHandler).Methods("GET")
	u.HandleFunc("/", Repo.RegisterNewUserHandler).Methods("POST")
	u.HandleFunc("/{id}", Repo.UpdateUserHandler).Methods("PUT")
	u.HandleFunc("/{id}", Repo.DeleteUserHandler).Methods("DELETE")

	ac := router.PathPrefix("/api/chirps").Subrouter()
	ac.HandleFunc("/", Repo.GetAllChirpsHandler).Methods("GET")
	ac.HandleFunc("/{id}", Repo.GetChirpByIDHandler).Methods("GET")
	ac.HandleFunc("/", Repo.PostNewChirpHandler).Methods("POST")
	ac.HandleFunc("/{id}", Repo.UpdateChirpHandler).Methods("PUT")
	ac.HandleFunc("/{id}", Repo.DeleteChirpHandler).Methods("DELETE")

	c := router.PathPrefix("/chirps").Subrouter()
	c.HandleFunc("/new", Repo.NewChirpPage).Methods("GET")
	c.HandleFunc("/new", Repo.PostNewChirpHandler).Methods("POST")
	c.HandleFunc("/edit", Repo.EditChirpPage).Methods("GET")
	c.HandleFunc("/show", Repo.ShowOneChirpPage).Methods("GET")

	c.HandleFunc("/summary", Repo.ChirpSummary).Methods("GET")

	router.HandleFunc("/login", Repo.Login).Methods("POST")

	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return router
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return templateCache, err
	}

	
	for _, page := range pages {
		name := filepath.Base(page)

		t, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return templateCache, err
		}

		if len(matches) > 0 {
			t, err = t.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = t
	}

	return templateCache, nil
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("unable to load .env file")
	}
}