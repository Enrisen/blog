package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Define handler for the root path
	http.HandleFunc("/", homeHandler)

	// Start the server
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// homeHandler renders the home page template
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure only the root path is handled
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmplPath := filepath.Join("ui", "html", "home.tmpl")

	// Parse the template file
	ts, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template
	err = ts.Execute(w, nil) // Pass nil for data for now
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
