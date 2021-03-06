package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// holds all of the information associated with the form both when it is initialized and after the form is submitted
type Form struct {
	url.Values // the data that gets passed to the form
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

// checks for required fields
// variadic function. Can pass as many fields as you want into it
func (f *Form) Required(fields ...string) {
	// range through all the fields 
	for _, field := range fields {
		value := f.Get(field)
		// if a field is blank, add an error to that field
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "this field cannot be blank")
		}
	}
}

// checks that the form field is included in the request
func (f *Form) Has(field string) bool {
	formField := f.Get(field)

	return formField != ""
}

// checks that the form field has a min length
func (f *Form) MinLength(field string, length int) bool {
	formField := f.Get(field)

	if len(formField) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}

	return true
}

// checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "invalid email address")
	}
}