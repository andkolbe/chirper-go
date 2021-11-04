package render

import (
	"net/http"
	"testing"

	"github.com/andkolbe/chirper-go/internal/models"
)

func TestAddDefaultData(t *testing.T)  {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	// put something in the session to test it
	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func TestTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	templateCache, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	// put the templateCache into the app variable
	app.TemplateCache = templateCache

	// need the request as a parameter to Template
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	// need a  responseWriter as a parameter to Template
	var ww myWriter

	err = Template(&ww, r, "home.page.html", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = Template(&ww, r, "non-existent.page.html", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()

	// put session data in the context
	// alexedwards/scs uses the X-Session header to find the identifier for the session 
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))

	// now that we have the session in the context, we need to put the context back into the request
	// calling session.Load will look for the header and create it if it does not exist
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func CreateTestTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}