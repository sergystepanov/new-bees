package main

import (
	"bytes"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/valyala/fasthttp"
)

func main() {
	port := envOr("PORT", "8080")
	log.Printf("Port: %v", port)

	routeFavIco := []byte("/favicon.ico")
	routeS := []byte("/s/")
	routeURL := []byte("/url")

	respGoAway := []byte("go away")

	server := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			path := ctx.Path()
			if bytes.Equal(path, routeFavIco) {
				ctx.Response.Header.Set(fasthttp.HeaderCacheControl, "Cache-Control: public, max-age=31536000")
				ctx.SetStatusCode(fasthttp.StatusNoContent)
				return
			}
			if bytes.HasPrefix(path, routeS) {
				ctx.SetBody(bytes.TrimPrefix(path, routeS))
				return
			}
			if bytes.HasPrefix(path, routeURL) {
				fastProxy(ctx)
				return
			}
			ctx.SetBody(respGoAway)
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
