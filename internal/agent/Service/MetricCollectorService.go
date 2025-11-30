package service

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"

	"github.com/bryzgaalov/metrics-collector/internal/agent/interface"
)

func CollectRuntimeMetrics(repo agentiface.MetricsCollectorRepository) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	repo.SetGauge("Alloc", float64(m.Alloc))
	repo.SetGauge("BuckHashSys", float64(m.BuckHashSys))
	repo.SetGauge("Frees", float64(m.Frees))
	repo.SetGauge("GCCPUFraction", m.GCCPUFraction)
	repo.SetGauge("GCSys", float64(m.GCSys))
	repo.SetGauge("HeapAlloc", float64(m.HeapAlloc))
	repo.SetGauge("HeapIdle", float64(m.HeapIdle))
	repo.SetGauge("HeapInuse", float64(m.HeapInuse))
	repo.SetGauge("HeapObjects", float64(m.HeapObjects))
	repo.SetGauge("HeapReleased", float64(m.HeapReleased))
	repo.SetGauge("HeapSys", float64(m.HeapSys))
	repo.SetGauge("LastGC", float64(m.LastGC))
	repo.SetGauge("Lookups", float64(m.Lookups))
	repo.SetGauge("MCacheInuse", float64(m.MCacheInuse))
	repo.SetGauge("MCacheSys", float64(m.MCacheSys))
	repo.SetGauge("MSpanInuse", float64(m.MSpanInuse))
	repo.SetGauge("MSpanSys", float64(m.MSpanSys))
	repo.SetGauge("Mallocs", float64(m.Mallocs))
	repo.SetGauge("NextGC", float64(m.NextGC))
	repo.SetGauge("NumForcedGC", float64(m.NumForcedGC))
	repo.SetGauge("NumGC", float64(m.NumGC))
	repo.SetGauge("OtherSys", float64(m.OtherSys))
	repo.SetGauge("PauseTotalNs", float64(m.PauseTotalNs))
	repo.SetGauge("StackInuse", float64(m.StackInuse))
	repo.SetGauge("StackSys", float64(m.StackSys))
	repo.SetGauge("Sys", float64(m.Sys))
	repo.SetGauge("TotalAlloc", float64(m.TotalAlloc))

	repo.SetGauge("RandomValue", rand.Float64())

	repo.AddCounter("PollCount", 1)
}

func SendAllMetricsToServer(repo agentiface.MetricsCollectorRepository, baseURL string) error {
	gauges, counters := repo.Snapshot()
	var lastErr error

	for name, value := range gauges {
		url := fmt.Sprintf("%s/update/gauge/%s/%s",
			baseURL,
			name,
			strconv.FormatFloat(value, 'f', -1, 64),
		)

		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			lastErr = err
			continue
		}
		_ = resp.Body.Close()
	}

	for name, value := range counters {
		url := fmt.Sprintf("%s/update/counter/%s/%d",
			baseURL,
			name,
			value,
		)

		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			lastErr = err
			continue
		}
		_ = resp.Body.Close()
	}

	return lastErr
}
