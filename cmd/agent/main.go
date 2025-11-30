package main

import (
	constants "github.com/bryzgaalov/metrics-collector/internal/agent/Constants"
	"github.com/bryzgaalov/metrics-collector/internal/agent/Repository"
	"github.com/bryzgaalov/metrics-collector/internal/agent/Service"

	"time"
)

func main() {
	storage := repository.NewMetricsCollector()

	tick := 0
	reportsEvery := int(constants.ReportInterval / constants.PollInterval)

	for {
		service.CollectRuntimeMetrics(storage)

		tick++
		if tick%reportsEvery == 0 {
			service.SendAllMetricsToServer(storage, constants.ServerBaseURL)
		}

		time.Sleep(constants.PollInterval)
	}

}
