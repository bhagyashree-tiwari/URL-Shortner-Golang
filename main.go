package main

import (
	// for generating random keys
	// "crypto/rand"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type URLShortner struct {
	urls map[string]string
}

// Implement URL shortening
func (us *URLShortner) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	// Generate a unique shortened key for the original URL
	shortKey := generateShortKey()
	us.urls[shortKey] = originalURL

	// construct the full shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:8080/short/%s", shortKey)

	// Render the HTML response with the shortened URL
	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
	<h2>URL Shortner</h2>
	<p>Original URL: %s</p>
	<p>Shortened URL: <a href ="%s"> %s </a></p>
	<form method="post" action="/shorten">
		<input type="text" name="url" placeholder="Enter a URL">
		<input type="submit" value="Shorten">
	</form>
	`, originalURL, shortenedURL, shortenedURL)
	fmt.Fprint(w, responseHTML)
}

// Implement URl redirection
func (us *URLShortner) handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len("/short/"):]
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	// Retrieve the original URL from the `urls` map using the shortened key
	originalURL, found := us.urls[shortKey]
	if !found {
		http.Error(w, "Shortened Key not found", http.StatusNotFound)
		return
	}

	// Redirect the user to the original URL
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)

}

// Generate Short Keys
func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]

	}
	return string(shortKey)
}

// Main function and server setup
func main() {
	shortner := &URLShortner{
		urls: make(map[string]string),
	}

	http.HandleFunc("/shorten", shortner.handleShorten)
	http.HandleFunc("/short/", shortner.handleRedirect)

	fmt.Println("URL Shortner is running on :8080")
	http.ListenAndServe(":8080", nil)
}
