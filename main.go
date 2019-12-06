package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	defaultPort     = "8080"
	defaultEncoding = "windows-1251"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
		log.Println("$PORT will be set as default:", defaultPort)
	}

	http.HandleFunc("/", root)
	http.HandleFunc("/s/", status)
	http.HandleFunc("/url", url)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Couldn't start the server: %s", err)
	}
}

// Handles root requests
func root(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Go away!")
}

// Handles status requests
func status(w http.ResponseWriter, r *http.Request) {
	word := strings.TrimPrefix(r.URL.Path, "/s/")
	_, _ = fmt.Fprintf(w, "%s", word)
}

// Handles URL requests
func url(w http.ResponseWriter, r *http.Request) {
	// parse POST params
	err := r.ParseForm()
	if err != nil {
		log.Println("Couldn't parse params:", err)
		return
	}

	url := r.Form.Get("_url")
	log.Println(url)
	if url == "" {
		log.Println("No params.")
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Has error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Has read error:", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset="+defaultEncoding)
	_, _ = fmt.Fprintf(w, "%s", body)
}
