package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	port := envOr("PORT", "8080")
	log.Printf("Port: %v", port)

	h := http.NewServeMux()

	h.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, "Go away!")
	})
	h.HandleFunc("/s/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "%s", strings.TrimPrefix(r.URL.Path, "/s/"))
	})
	h.HandleFunc("/url", proxy())

	server := &http.Server{Addr: ":" + port, Handler: h}
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Printf("Terminate")
		_ = server.Shutdown(context.Background())
	}()

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server: %s", err)
	}
}
