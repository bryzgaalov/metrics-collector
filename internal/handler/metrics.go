package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bryzgaalov/metrics-collector/internal/View"
	interfaces "github.com/bryzgaalov/metrics-collector/internal/interface"
	models "github.com/bryzgaalov/metrics-collector/internal/model"
	"github.com/go-chi/chi/v5"
)

var Storage interfaces.MetricsStorageInterface
var indexTmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Metrics</title>
</head>
<body>
    <h1>Metrics</h1>
    <table border="1">
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Value</th>
        </tr>
        {{range .}}
        <tr>
            <td>{{.Name}}</td>
            <td>{{.Type}}</td>
            <td>{{.Value}}</td>
        </tr>
        {{end}}
    </table>
</body>
</html>
`))

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r) {
		return
	}

	tail, ok := parseURL(w, r)
	if !ok {
		return
	}

	switch len(tail) {
	case 0, 1:
		http.Error(w, "metric name not provided", http.StatusNotFound)
		return
	case 2:
		http.Error(w, "metric value not provided", http.StatusBadRequest)
		return
	case 3:
	default:
		http.Error(w, "too many path segments", http.StatusBadRequest)
		return
	}

	metricType := tail[0]
	metricName := tail[1]
	rawValue := tail[2]

	if metricName == "" {
		http.Error(w, "metric name not provided", http.StatusNotFound)
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
	default:
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		return
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
func MetricValueHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	if metricType == "" {
		http.Error(w, "metric type is required", http.StatusBadRequest)
		return
	}
	if metricType != models.Gauge && metricType != models.Counter {
		http.Error(w, "unknown metric type", http.StatusBadRequest)
		return
	}

	metricName := chi.URLParam(r, "name")
	if metricName == "" {
		http.Error(w, "metric name is required", http.StatusBadRequest)
		return
	}

	metric, ok := Storage.Get(metricName)
	if !ok {
		http.Error(w, "metric not found", http.StatusNotFound)
		return
	}
	if string(metric.MType) != metricType {
		http.Error(w, "metric type mismatch", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	switch metric.MType {
	case models.Gauge:
		if metric.Value == nil {
			http.Error(w, "metric value is nil", http.StatusInternalServerError)
			return
		}
		_, err := fmt.Fprintf(w, "%f", *metric.Value)
		if err != nil {
			return
		}

	case models.Counter:
		if metric.Delta == nil {
			http.Error(w, "metric delta is nil", http.StatusInternalServerError)
			return
		}
		_, err := fmt.Fprintf(w, "%d", *metric.Delta)
		if err != nil {
			return
		}
	}
}
func BaseHTMLHandler(w http.ResponseWriter, r *http.Request) {
	metrics := Storage.GetAll()
	var list []View.MetricView
	_ = r.URL.Path // Сделал заглушку, потому что ругалось на параметр r т.к. он был не используемым

	for name, m := range metrics {
		mv := View.MetricView{
			Name: name,
			Type: m.MType,
		}

		switch m.MType {
		case models.Gauge:
			if m.Value != nil {
				mv.Value = strconv.FormatFloat(*m.Value, 'f', -1, 64)
			}
		case models.Counter:
			if m.Delta != nil {
				mv.Value = strconv.FormatInt(*m.Delta, 10)
			}

		}

		list = append(list, mv)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := indexTmpl.Execute(w, list); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		return
	}
}

func validateRequest(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return false
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Only text/plain is supported!", http.StatusUnsupportedMediaType)
		return false
	}
	return true
}
func parseURL(w http.ResponseWriter, r *http.Request) ([]string, bool) {
	rawPath := r.URL.Path
	path := strings.Trim(rawPath, "/")
	segments := strings.Split(path, "/")

	if len(segments) == 0 || segments[0] != "update" {
		http.NotFound(w, r)
		return nil, false
	}

	tail := segments[1:]
	return tail, true
}
