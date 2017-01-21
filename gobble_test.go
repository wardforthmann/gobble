package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
)

func TestHandlePost(t *testing.T) {
	req, err := http.NewRequest("GET", "/", strings.NewReader("Test POST please ignore"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlePost)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

}
