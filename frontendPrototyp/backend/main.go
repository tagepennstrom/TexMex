package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Allow specific methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow specific headers
}

func compileLatex(w http.ResponseWriter, r *http.Request) {
	enableCORS(w) // Add CORS headers

	if r.Method == http.MethodOptions {
		return // Handle preflight requests
	}

	// Read LaTeX code from request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Write LaTeX code to a .tex file
	texFile := "document.tex"
	err = ioutil.WriteFile(texFile, body, 0644)
	if err != nil {
		http.Error(w, "Failed to write .tex file", http.StatusInternalServerError)
		return
	}

	// Run LaTeX compiler (use tectonic or pdflatex)
	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", texFile)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		http.Error(w, "Compilation failed", http.StatusInternalServerError)
		return
	}

	// Check if the PDF file was created
	pdfFile := "document.pdf"
	if _, err := os.Stat(pdfFile); os.IsNotExist(err) {
		http.Error(w, "PDF file not found after compilation", http.StatusInternalServerError)
		return
	}

	// Send back the correct PDF URL
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"pdfUrl": "/pdf"}`) // Ensure the URL matches the /pdf endpoint
}

func servePdf(w http.ResponseWriter, r *http.Request) {
	enableCORS(w) // Add CORS headers

	// Log the request to verify the endpoint is being hit
	fmt.Println("Serving PDF file...")

	// Serve the PDF file with the correct content type
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, "document.pdf")
}

func main() {
	http.HandleFunc("/compile", compileLatex)
	http.HandleFunc("/pdf", servePdf)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
