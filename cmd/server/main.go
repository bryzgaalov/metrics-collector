package main

import (
	"flag"
	"net"
	"net/http"

	constants "github.com/bryzgaalov/metrics-collector/internal/agent/Constants"
	"github.com/bryzgaalov/metrics-collector/internal/handler"
	"github.com/bryzgaalov/metrics-collector/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	address := flag.String("a", constants.ServerBaseURL, "HTTP server address")
	flag.Parse()

	mem := repository.NewMemStorage()
	handler.Storage = mem

	r := chi.NewRouter()
	r.Post("/update/*", handler.MetricsHandler)
	r.Get("/values", handler.MetricsListHandler)
	r.Get("/value/{type}/{name}", handler.MetricValueHandler)
	r.Get("/", handler.BaseHTMLHandler)

	listenAddr := *address
	if host, port, err := net.SplitHostPort(*address); err == nil && port != "" {
		_ = host
		listenAddr = ":" + port
	}

	err := http.ListenAndServe(listenAddr, r)
	if err != nil {
		panic(err)
	}

}
