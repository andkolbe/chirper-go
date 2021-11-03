package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/andkolbe/chirper-go/internal/models"
)

// A handler responds to an HTTP request
// It is responsible for writing response headers and bodies

// All the dependencies for our handlers are explicitly defined in one place
// models.UserModel is a dependency of the Repository struct
type Repository struct {
	dbmodel models.DBModel
}

// creates a new repository
// the Repository type is populated with all of the info received as parameters and it handed back as a pointer to Repository
func NewRepo(dbm models.DBModel) *Repository {
	return &Repository{
		dbmodel: dbm,
	}
}

// response format
type response struct {
    ID      int64  `json:"id,omitempty"`
    Message string `json:"message,omitempty"`
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.page.html")
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.page.html")
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	// first step. Parse the template
	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template", err)
		return 
	}
}