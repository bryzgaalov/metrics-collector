package main

import (
	"flag"

	constants "github.com/bryzgaalov/metrics-collector/internal/agent/Constants"
	"github.com/bryzgaalov/metrics-collector/internal/agent/Repository"
	"github.com/bryzgaalov/metrics-collector/internal/agent/Service"

	"time"
)

func main() {
	address := flag.String("a", constants.ServerBaseURL, "HTTP server address")
	rSec := flag.Int("r", int(constants.ReportInterval/time.Second), "Report interval in seconds")
	pSec := flag.Int("p", int(constants.PollInterval/time.Second), "Poll interval in seconds")
	flag.Parse()

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
