package main

import (
	"net/http"

	"github.com/bryzgaalov/metrics-collector/internal/handler"
	"github.com/bryzgaalov/metrics-collector/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	mem := repository.NewMemStorage()
	handler.Storage = mem

	r := chi.NewRouter()
	r.Post("/update/*", handler.MetricsHandler)
	r.Get("/values", handler.MetricsListHandler)
	r.Get("/value/{type}/{name}", handler.MetricValueHandler)
	r.Get("/", handler.BaseHTMLHandler)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}

}
