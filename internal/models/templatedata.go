package models

import "github.com/andkolbe/chirper-go/internal/forms"

// holds any type of data we send from handlers to templates
type TemplateData struct {
	// use maps because we will could send more than one string, int, whatever to the templates
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
