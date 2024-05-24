package main

import (
	"log"
	"net/http"
)

func main() {
	// Serve the biding.html file at the root URL
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "biding.html")
	})

	// Serve static files from the "static" directory
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting server on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
