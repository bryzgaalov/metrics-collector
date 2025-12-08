package interfaces

import (
	"context"

	model "github.com/bryzgaalov/metrics-collector/internal/model"
)

type MetricsService interface {
	UpdateMetric(ctx context.Context, metric model.Metrics) error
	GetMetric(ctx context.Context, metricName string) (model.Metrics, bool)
	ListMetrics(ctx context.Context) (map[string]model.Metrics, error)
}
