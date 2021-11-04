package main

import (
	"fmt"
	"testing"

	"github.com/andkolbe/chirper-go/internal/config"
	"github.com/gorilla/mux"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	router := routes(&app)

	switch v := router.(type) {
	case *mux.Router:
		// do nothing, test passed
	default:
		t.Error(fmt.Sprintf("type is not *mux.Router, type is %T", v))
	}
}