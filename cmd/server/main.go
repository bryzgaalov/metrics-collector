package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bryzgaalov/metrics-collector/internal/handler"
	"github.com/bryzgaalov/metrics-collector/internal/repository"
	"github.com/bryzgaalov/metrics-collector/internal/service"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
)

type NetAddress struct {
	Host string
	Port int
}

type Config struct {
	Address string `env:"ADDRESS"`
}

func (a NetAddress) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = port
	return nil
}

var addr = &NetAddress{
	Host: "localhost",
	Port: 8080,
}

func main() {
	flag.Var(addr, "a", "server address in form host:port")
	flag.Parse()

	mem := repository.NewMemStorage()
	handler.Service = service.NewMetricsService(mem)

	r := chi.NewRouter()
	r.Post("/update/*", handler.MetricsHandler)
	r.Get("/values", handler.MetricsListHandler)
	r.Get("/value/{type}/{name}", handler.MetricValueHandler)
	r.Get("/", handler.BaseHTMLHandler)

	listenAddr := ":" + strconv.Itoa(addr.Port)
	var cfg Config
	errParsing := env.Parse(&cfg)
	if errParsing != nil {
		log.Fatal(errParsing)
	}
	if cfg.Address != "" {
		listenAddr = cfg.Address
	}

	err := http.ListenAndServe(listenAddr, r)
	if err != nil {
		panic(err)
	}

}
