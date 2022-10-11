package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/valyala/fasthttp"
)

func main() {
	port := envOr("PORT", "8080")
	log.Printf("Port: %v", port)

	server := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			path := string(ctx.Path())
			if path == "/favicon.ico" {
				ctx.Response.Header.Set(fasthttp.HeaderCacheControl, "Cache-Control: public, max-age=31536000")
				ctx.SetStatusCode(fasthttp.StatusNoContent)
				return
			}
			if strings.HasPrefix(path, "/s/") {
				ctx.SetBody([]byte(strings.TrimPrefix(path, "/s/")))
				return
			}
			if strings.HasPrefix(path, "/url") {
				fastProxy(ctx)
				return
			}
			ctx.SetBody([]byte("go away"))
		},
	}
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		_ = server.Shutdown()
	}()

	if err := server.ListenAndServe(":" + port); err != nil {
		log.Fatalf("server start0 fail, %v", err)
	}
}

func envOr(env, def string) string {
	e := os.Getenv(env)
	if e == "" {
		return def
	}
	return e
}
