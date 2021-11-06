package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/andkolbe/chirper-go/internal/config"
	"github.com/andkolbe/chirper-go/internal/models"
	"github.com/justinas/nosurf"
)

// FuncMap is a map of functions that can be used in a template
// we can add functionality to Golang templates that aren't built into the language
var functions = template.FuncMap{}

var app *config.AppConfig

var pathToTemplates = "./templates"

// sets the config for the render package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// adds data that is available to every page in the app without having to manually add it in to every page ourselves
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// we want every page to have a CSRF token attached to it
	td.CSRFToken = nosurf.Token(r)

	// pop the flash, warning, and/or error message out of the session if there is one
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}

	return td
}

// renders templates using html/templates
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var templateCache map[string]*template.Template

	// if we are in production, use the template cache 
	// Otherwise, in development, rebuild it on every request so we don't have to restart the server for every change
	if app.UseCache {
		// get the template cache (that is initialized in main.go) from the app config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	// pull the template out of the cache
	t, ok := templateCache[tmpl]
	// if the template doesn't exist in the cache
	if !ok {
		return errors.New("can't get template from cache")
	}

	// holds bytes. Put the parsed template from memory into bytes
	// write to the buffer instead of straight to the response writer so we can check for an error, and determine where it came from more easily
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	// take the template, execute it, pass it data, and store the value in the buf variable
	_ = t.Execute(buf, td)

	// write the buf to the response writer
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}

	return nil
}

// creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	// create a map to store templates for quick lookups
	// key: name of the template, value: the template itself 
	templateCache := map[string]*template.Template{}

	// store the filepath of every template that ends in .page.html in a variable
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return templateCache, err
	}

	// range through all of the pages and create a template out of each one
	for _, page := range pages {
		// extract the page name (last element) out of the filepath 
		name := filepath.Base(page)

		// create templates based on the page name, add in any outside functions we've created, and parse the page
		t, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		// determine if there are any layouts in our app that match with the template we just created
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return templateCache, err
		}

		// if there is a match, then the length of matches will be greater than zero
		if len(matches) > 0 {
			// go to the template and parse the layout
			t, err = t.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return templateCache, err
			}
		}

		// take the template and add it to the cache
		templateCache[name] = t
	}

	return templateCache, nil
}