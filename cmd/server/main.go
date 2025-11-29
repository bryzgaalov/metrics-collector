package main

import (
	"net/http"

	"github.com/bryzgaalov/metrics-collector/internal/handler"
	"github.com/bryzgaalov/metrics-collector/internal/repository"
)

func main() {
	mem := repository.NewMemStorage()
	handler.Storage = mem

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, handler.MetricsHandler)
	mux.HandleFunc("/values", handler.MetricsListHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}
