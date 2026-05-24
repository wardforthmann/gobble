package main

import (
	"net/http"
)

func addRoutes(
	r *http.ServeMux,
	creds map[string]string,
) {
	var showFilesHandler http.Handler = http.HandlerFunc(showFiles)
	if len(creds) > 0 {
		showFilesHandler = basicAuthMiddleware(creds)(showFilesHandler)
	}
	r.Handle("GET /{path...}", showFilesHandler)

	postHandler := statusCodeHandler()(http.HandlerFunc(handlePost))
	r.Handle("POST /{path...}", postHandler)
}

func basicAuthMiddleware(creds map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok || creds[user] != pass {
				w.Header().Set("WWW-Authenticate", `Basic realm=""`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
