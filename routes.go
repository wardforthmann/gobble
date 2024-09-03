package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func addRoutes(
	r *chi.Mux,
	creds map[string]string,
) {

	if len(creds) > 0 {
		r.With(middleware.BasicAuth("", creds)).Get("/*", showFiles)
	} else {
		r.Get("/*", showFiles)
	}

	r.With(statusCodeHandler()).Post("/*", handlePost)
}
