package handlers

import "github.com/andkolbe/chirper-go/internal/models"

// A handler responds to an HTTP request
// It is responsible for writing response headers and bodies

// All the dependencies for our handlers are explicitly defined in one place
// models.UserModel is a dependency of the Repository struct
type Repository struct {
	dbmodel models.DBModel
}

// creates a new repository
// the Repository type is populated with all of the info received as parameters and it handed back as a pointer to Repository
func NewRepo(dbmod models.DBModel) *Repository {
	return &Repository{
		dbmodel: dbmod,
	}
}

// response format
type response struct {
    ID      int64  `json:"id,omitempty"`
    Message string `json:"message,omitempty"`
}