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
	reportInterval := flag.Duration("r", constants.ReportInterval, "Report interval in seconds")
	pollInterval := flag.Duration("p", constants.PollInterval, "Poll interval in seconds")
	flag.Parse()

	storage := repository.NewMetricsCollector()

	tick := 0
	reportsEvery := int(*reportInterval / *pollInterval)

	for {
		service.CollectRuntimeMetrics(storage)

		tick++
		if tick%reportsEvery == 0 {
			_ = service.SendAllMetricsToServer(storage, *address)
		}

		time.Sleep(*pollInterval)
	}

}
