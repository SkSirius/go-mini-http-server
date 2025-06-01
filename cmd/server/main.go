package main

import (
	"net/http"

	"github.com/sksirius/mini-http-server/internal/handlers"
	"github.com/sksirius/mini-http-server/internal/middleware"
	"github.com/sksirius/mini-http-server/internal/router"
	"github.com/sksirius/mini-http-server/internal/server"
)

func main() {
	chain := middleware.Chain(
		middleware.RecoveryMiddleware,
		middleware.LoggingMiddleware,
	)

	r := router.New()
	r.Handle("GET", "/hello", chain(http.HandlerFunc(handlers.HelloHandler)))
	r.Handle("GET", "/time", chain(http.HandlerFunc(handlers.TimeHandler)))
	r.Handle("POST", "/echo", chain(http.HandlerFunc(handlers.EchoHandler)))

	s := server.New(":8080", r)
	s.Start()
	s.GracefulShutdown()
}
