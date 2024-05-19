package middleware

import "net/http"

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request content type is 'application/json'
		if r.Header.Get("Content-Type") != "" && r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type header is not application/json", http.StatusUnsupportedMediaType)
			return
		}

		// Set the response writer content type to 'application/json'
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
