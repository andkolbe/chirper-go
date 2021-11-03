package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/andkolbe/chirper-go/internal/config"
)

// FuncMap is a map of functions that can be used in a template
// we can add functionality to Golang templates that aren't built into the language
var functions = template.FuncMap{}

var app *config.AppConfig

// sets the config for the render package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// renders templates using html/templates
func Template(w http.ResponseWriter, tmpl string) {
	var templateCache map[string]*template.Template

	// if we are using the cache, use it, otherwise rebuild it
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
		log.Fatal("could not get template from template cache")
	}

	// holds bytes. Put the parsed template from memory into bytes
	// write to the buffer instead of straight to the response writer so we can check for an error, and determine where it came from more easily
	buf := new(bytes.Buffer)

	// take the template, execute it, don't pass it any data, and store the value in the buf variable
	_ = t.Execute(buf, nil)

	// write the buf to the response writer
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}

}

// creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	// create a map to store templates for quick lookups
	// key: name of the template, value: the template itself 
	templateCache := map[string]*template.Template{}

	// store the filepath of every template that ends in .page.html in a variable
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return templateCache, err
	}

	// range through all of the pages and create a template out of each one
	for _, page := range pages {
		// extract the page name (last element) out of the filepath 
		name := filepath.Base(page)

		fmt.Println("Page is currently", page)

		// create templates based on the page name, add in any outside functions we've created, and parse the page
		t, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		// determine if there are any layouts in our app that match with the template we just created
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return templateCache, err
		}

		// if there is a match, then the length of matches will be greater than zero
		if len(matches) > 0 {
			// go to the template and parse the layout
			t, err = t.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return templateCache, err
			}
		}

		// take the template and add it to the cache
		templateCache[name] = t
	}

	return templateCache, nil
}