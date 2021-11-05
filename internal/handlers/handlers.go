package handlers

import (
	"log"
	"net/http"

	"github.com/andkolbe/chirper-go/internal/config"
	"github.com/andkolbe/chirper-go/internal/forms"
	"github.com/andkolbe/chirper-go/internal/models"
	"github.com/andkolbe/chirper-go/internal/render"
)

// A handler responds to an HTTP request
// It is responsible for writing response headers and bodies

// All the dependencies for our handlers are explicitly defined in one place
// models.UserModel is a dependency of the Repository struct
type Repository struct {
	dbmodel models.DBModel
	App     *config.AppConfig
}

// the repository used by the handlers
// pass by reference to Repository, since all references are pointing to the same place (the address in memory where the repo lives), Repo is never out of sync
var Repo *Repository

// creates a new repository
// the Repository type is populated with all of the info received as parameters and it handed back as a pointer to Repository
func NewRepo(dbm models.DBModel, a *config.AppConfig) *Repository {
	return &Repository{
		dbmodel: dbm,
		App:     a,
	}
}

// sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// Home Page
func (repo *Repository) HomePage(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.html", &models.TemplateData{})
}

// About Page
func (repo *Repository) AboutPage(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.html", &models.TemplateData{})
}

// Show One Chirp Page
func (repo *Repository) ShowOneChirpPage(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "show_chirp.page.html", &models.TemplateData{})
}

// New Chirp Page
func (repo *Repository) NewChirpPage(w http.ResponseWriter, r *http.Request) {
	// initialize an empty chirp model and pass it to the data when the page loads so we can use it later on
	var emptyChirp models.Chirp
	data := make(map[string]interface{})
	data["chirp"] = emptyChirp

	render.Template(w, r, "new_chirp.page.html", &models.TemplateData{
		Form: forms.New(nil), // include an empty form with the template
		Data: data,
	})
}

// Edit Chirp Page
func (repo *Repository) EditChirpPage(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "edit_chirp.page.html", &models.TemplateData{})
}

// SAME LOGIC FOR PROFILE PAGE
func (repo *Repository) ChirpSummary(w http.ResponseWriter, r *http.Request) {
	// pull the data called "chirp" out of the session and type assert it to type models.Chirp
	chirp, ok := repo.App.Session.Get(r.Context(), "chirp").(models.Chirp)
	if !ok {
		repo.App.ErrorLog.Println("Can't get error from session")
		repo.App.Session.Put(r.Context(), "error", "can't get chirp from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// remove the chirp from the session
	repo.App.Session.Remove(r.Context(), "chirp")

	// add chirp that was pulled out of the session into the template data
	data := make(map[string]interface{})
	data["chirp"] = chirp

	// render the page with the data passed in
	render.Template(w, r, "chirp-summary.page.html", &models.TemplateData{
		Data: data,
	})
}