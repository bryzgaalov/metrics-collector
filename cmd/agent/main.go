package main

import (
	"flag"
	"log"

	constants "github.com/bryzgaalov/metrics-collector/internal/agent/Constants"
	"github.com/bryzgaalov/metrics-collector/internal/agent/Repository"
	"github.com/bryzgaalov/metrics-collector/internal/agent/Service"
	"github.com/caarlos0/env/v6"

	"time"
)

type Config struct {
	Address       string `env:"ADDRESS"`
	ReportSeconds int    `env:"REPORT_INTERVAL"`
	PollSeconds   int    `env:"POLL_INTERVAL"`
}

func main() {
	address := flag.String("a", constants.ServerBaseURL, "HTTP server address")
	rSec := flag.Int("r", int(constants.ReportInterval/time.Second), "Report interval in seconds")
	pSec := flag.Int("p", int(constants.PollInterval/time.Second), "Poll interval in seconds")
	flag.Parse()

	var cfg Config
	errParsing := env.Parse(&cfg)
	if errParsing != nil {
		log.Fatal(errParsing)
	}
	if cfg.Address != "" {
		address = &cfg.Address
	}
	if cfg.PollSeconds != 0 {
		pSec = &cfg.PollSeconds
	}
	if cfg.ReportSeconds != 0 {
		rSec = &cfg.ReportSeconds
	}

	reportInterval := time.Duration(*rSec) * time.Second
	pollInterval := time.Duration(*pSec) * time.Second
	storage := repository.NewMetricsCollector()

	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)

	for {
		select {
		case <-pollTicker.C:
			service.CollectRuntimeMetrics(storage)

		case <-reportTicker.C:
			service.SendAllMetricsToServer(storage, *address)
		}
	}

}
