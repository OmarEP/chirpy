package main

import (
	"net/http"
)

func main() {
	// Create a new http.ServeMux
	mux := http.NewServeMux()

	// Wrap that mux in a custom middleware function that adds 
	// CORS headers to the response (see the tip below on how to do that).
	corsMux := middlewareCors(mux)

	// Create a new http.Server and use the corsMux as the handler
	server := &http.Server{
		Addr: "localhost:8080",
		Handler: corsMux,
	}
	// Use the server's ListenAndServe method to start the server
	server.ListenAndServe()
	// Build and run your server (e.g. go build -o out && ./out)
	
	// Open http://localhost:8080 in your browser. You should see a 404 error because we haven't connected any handler logic yet. 
	// Don't worry, that's what is expected for the tests to pass for now.

	// Paste the URL of your web server (e.g. http://localhost:8080)
	// into the text box and run the HTTP tests.
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
