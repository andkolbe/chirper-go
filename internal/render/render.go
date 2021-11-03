package render

import (
	"fmt"
	"html/template"
	"net/http"
)

func Template(w http.ResponseWriter, tmpl string) {
	// Parse the template that is passed in, execute it, add the given data (if any), and add it to the response
	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template", err)
		return
	}
}