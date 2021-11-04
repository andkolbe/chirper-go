package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// create a table test so we don't have to write tests for each handler individually

// holds whatever we are posting to a page
type postData struct {
	key   string
	value string
}

// initialize a struct and populate it with data at the same time
var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"new chirp page", "/chirps/new", "GET", []postData{}, http.StatusOK},
	{"edit chirp page", "/chirps/edit", "GET", []postData{}, http.StatusOK},
	{"show chirp page", "/chirps/show", "GET", []postData{}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	// httptest has a test server built into it
	testServer := httptest.NewTLSServer(routes)
	// close the testServer when this function is done running
	defer testServer.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			// make a request as though we were a client
			response, err := testServer.Client().Get(testServer.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, response.StatusCode)
			}
		} else {

		}
	}
}
