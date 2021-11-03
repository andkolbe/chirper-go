package render

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// FuncMap is a map of functions that can be used in a template
// we can add functionality to Golang templates that aren't built into the language
var functions = template.FuncMap{}

// renders templates using html/templates
func Template(w http.ResponseWriter, tmpl string) {

	_, err := TemplateTest(w)
	if err != nil {
		fmt.Println("Error getting template cache:", err)
	}

	// Parse the template that is passed in, execute it, add the given data (if any), and add it to the response
	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	err = parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template", err)
		return
	}
}

func TemplateTest(w http.ResponseWriter) (map[string]*template.Template, error) {
	// create a map to store templates for quick lookups
	// key: name of the template, value: the template itself 
	myCache := map[string]*template.Template{}

	// store the filepath of every template that ends in .page.html in a variable
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	// range through all of the pages and create a template out of each one
	for _, page := range pages {
		// extract the page name (last element) out of the filepath 
		name := filepath.Base(page)

		fmt.Println("Page is currently", page)

		// create templates based on the page name, add in any outside functions we've created, and parse the page
		t, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// determine if there are any layouts in our app that match with the template we just created
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		// if there is a match, then the length of matches will be greater than zero
		if len(matches) > 0 {
			// go to the template and parse the layout
			t, err = t.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		// take the template and add it to the cache
		myCache[name] = t
	}

	return myCache, nil
}