package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	// we need to create a handler to pass to NoSurf so it can hand back a handler
	var myH myHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing, this is what we expect. Test passes. Yay
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, type is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, type is %T", v))
	}
}