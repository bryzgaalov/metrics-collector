package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	interfaces "github.com/bryzgaalov/metrics-collector/internal/interface"
	models "github.com/bryzgaalov/metrics-collector/internal/model"
)

var Storage interfaces.MetricsStorageInterface

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Only text/plain is supported!", http.StatusUnsupportedMediaType)
		return
	}

	respString := strings.TrimPrefix(r.URL.Path, `/update/`)
	partsStrings := strings.Split(respString, `/`)
	if len(partsStrings) != 3 {
		http.Error(w, "Only /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ> is supported!", http.StatusUnsupportedMediaType)
	}

	metricType := partsStrings[0]
	metricName := partsStrings[1]
	rawValue := partsStrings[2]
	if metricType != models.Gauge && metricType != models.Counter {
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		return
	}

	if metricName == "" {
		http.Error(w, "unknown metric type", http.StatusNotFound)
		return
	}

	var metric models.Metrics

	switch metricType {
	case models.Gauge:
		val, err := strconv.ParseFloat(rawValue, 64)
		if err != nil {
			http.Error(w, "gauge value must be float", http.StatusBadRequest)
			return
		}
		metric = models.Metrics{
			ID:    metricName,
			MType: models.Gauge,
		}
		metric.Value = &val

	case models.Counter:
		delta, err := strconv.ParseInt(rawValue, 10, 64)
		if err != nil {
			http.Error(w, "counter value must be int", http.StatusBadRequest)
			return
		}

		if existing, ok := Storage.Get(metricName); ok && existing.Delta != nil {
			delta = *existing.Delta + delta
		}

		metric = models.Metrics{
			ID:    metricName,
			MType: models.Counter,
		}
		metric.Delta = &delta
	}

	if err := Storage.Save(metric); err != nil {
		http.Error(w, "failed to save metric", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Length", r.Header.Get("Content-Length"))
	w.Header().Set("Date", time.Now().UTC().Format(time.RFC1123))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func MetricsListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	all := Storage.GetAll()

	w.Header().Set("Content-Type", "application/json")

	data, err := json.MarshalIndent(all, "", "  ") // <- красиво форматируем
	if err != nil {
		http.Error(w, "failed to marshal metrics", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
