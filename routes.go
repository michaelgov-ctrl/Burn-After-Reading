package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthcheck", app.healthCheckHandler)
	mux.HandleFunc("GET /v1/messages/*", app.retrieveAndRemoveMessageHandler)
	mux.HandleFunc("POST /v1/messages", app.postMessageHandler)
	mux.Handle("GET /metrics", promhttp.Handler())

	return app.recoverPanic(mux)
}
