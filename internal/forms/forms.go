package forms

import (
	"net/http"
	"net/url"
)

// holds all of the information associated with the form both when it is initialized and after the form is submitted
type Form struct {
	url.Values
	Errors errors
}

// initializes a form struct that only contains the data that was passed in. No (potential) errors yet 
// pass it nil because it has no data attached to it when it is initialized. Only when it is submitted does the form contain data
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}), // initialize with empty errors map 
	}
}

// checks that the form field is included in the request
func (f *Form) Has(field string, r *http.Request) bool {
	formField := r.Form.Get(field)

	// return true if the form field isn't blank
	return formField != ""
}