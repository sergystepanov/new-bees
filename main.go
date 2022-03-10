package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	defaultPort            = "8080"
	defaultEncoding        = "windows-1251"
	defaultCharsetResponse = "text/html; charset=" + defaultEncoding

	messageGoAway = "Go away!"
)

var client *http.Client

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
		log.Println("$PORT will be set as default:", defaultPort)
	}

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client = &http.Client{
		Transport: customTransport,
	}

	http.HandleFunc("/", root)
	http.HandleFunc("/s/", status)
	http.HandleFunc("/url", url)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Couldn't start the server: %s", err)
	}
}

// Handles root requests
func root(w http.ResponseWriter, _ *http.Request) { _, _ = fmt.Fprintf(w, messageGoAway) }

// Handles status requests
func status(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "%s", strings.TrimPrefix(r.URL.Path, "/s/"))
}

// Handles URL requests
func url(w http.ResponseWriter, r *http.Request) {
	// parse POST params
	if err := r.ParseForm(); err != nil {
		log.Printf("couldn't parse params, %v", err)
		return
	}

	url := r.Form.Get("_url")
	if url == "" {
		log.Printf("bad params!")
		return
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("couldn't get, %v", err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("has close error, %v", err)
			return
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("couldn't read, %v", err)
		return
	}

	w.Header().Set("Content-Type", defaultCharsetResponse)
	_, _ = fmt.Fprintf(w, "%s", body)
}
