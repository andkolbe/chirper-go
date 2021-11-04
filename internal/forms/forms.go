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

// returns true if there are no errors on the form
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// initializes a form struct that only contains the data that was passed in. No (potential) errors yet 
// pass it nil on the handlers because it has no data attached to it when it is initialized. Only when it is submitted does the form contain data
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}), // initialize with empty errors map 
	}
}

// checks that the form field is included in the request
func (f *Form) Has(field string, r *http.Request) bool {
	formField := r.Form.Get(field)

	if formField == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}

	return true
}