package main

import (
	"flag"
	"net/http"

	"github.com/bryzgaalov/metrics-collector/internal/handler"
	"github.com/bryzgaalov/metrics-collector/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	address := flag.String("a", "localhost:8080", "HTTP server address")
	flag.Parse()

	mem := repository.NewMemStorage()
	handler.Storage = mem

	r := chi.NewRouter()
	r.Post("/update/*", handler.MetricsHandler)
	r.Get("/values", handler.MetricsListHandler)
	r.Get("/value/{type}/{name}", handler.MetricValueHandler)
	r.Get("/", handler.BaseHTMLHandler)

	err := http.ListenAndServe(*address, r)
	if err != nil {
		panic(err)
	}

}
