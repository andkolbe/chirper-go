package forms

// hold all potential form field errors in a map
// each form field is a different key, and the potential errors are the values
// []string because there could be more than one error on a field
type errors map[string][]string

// adds an error message for a given form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// returns the first error message
func (e errors) Get(field string) string {
	errorString := e[field]

	// if there are no errors, return nothing
	if len(errorString) == 0 {
		return ""
	}

	return errorString[0] // if there are multiple errors, display the first one
}