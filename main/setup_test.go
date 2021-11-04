package main

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// everything inside this function will run before any of the tests are run
func TestMain(m *testing.M) {

	LoadEnv()

	// before this function exists, run all of the tests
	os.Exit(m.Run())
}

type myHandler struct{}

// all we need to do to create a handler interface is to add a ServeHTTP function that takes in a ResponseWriter and a *Request
func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("unable to load .env file")
	}
}