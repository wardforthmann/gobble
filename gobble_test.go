package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"os"
)

func init() {
	err := os.MkdirAll("testPublic", 0644)
	if err != nil {
		panic("unable to create dir")
	}
	os.Chdir("testPublic")
}

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

func TestHandlePostWithDir(t *testing.T) {
	req, err := http.NewRequest("GET", "/?dir=testDir", strings.NewReader("Test POST please ignore"))
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

func TestHandlePostWithStatusCode(t *testing.T) {
	req, err := http.NewRequest("GET", "/?status_code=418", strings.NewReader("Test POST please ignore"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	//Chi would normally handle the statusCodeHandler function with its "with" function but we arent going through chi here
	handler := statusCodeHandler(http.HandlerFunc(handlePost))

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTeapot {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusTeapot)
	}
}