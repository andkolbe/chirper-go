package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	{"post new chirp", "/chirps/new", "POST", []postData{
		{key: "userid", value: "1"},
		{key: "content", value: "Hello World"},
		{key: "location", value: "Birmingham"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	// httptest has a test server built into it
	testServer := httptest.NewTLSServer(routes)
	// close the testServer when this function is done running
	defer testServer.Close()

	for _, test := range theTests {
		if test.method == "GET" {
			// make a request as though we were a client
			response, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", test.name, test.expectedStatusCode, response.StatusCode)
			}
		} else { // POST
			// construct a variable that is in the format that the test server expects to receive
			// url.Values is a built in type, part of the standard library, that holds information as a POST request
			values := url.Values{}
			// populate values with our entries
			for _, param := range test.params {
				values.Add(param.key, param.value)
			}
			// call the testServer again but with a post form
			response, err := testServer.Client().PostForm(testServer.URL + test.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", test.name, test.expectedStatusCode, response.StatusCode)
			}
		}
	}
}
