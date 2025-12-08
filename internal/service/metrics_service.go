package service

import (
	"context"

	interfaces "github.com/bryzgaalov/metrics-collector/internal/interface"
	model "github.com/bryzgaalov/metrics-collector/internal/model"
)

type MetricsService struct {
	storage interfaces.MetricsStorage
}

func NewMetricsService(storage interfaces.MetricsStorage) *MetricsService {
	return &MetricsService{storage: storage}
}

func (s *MetricsService) GetMetric(ctx context.Context, metricName string) (model.Metrics, bool) {
	return s.storage.Get(ctx, metricName)
}
func (s *MetricsService) ListMetrics(ctx context.Context) (map[string]model.Metrics, error) {
	return s.storage.GetAll(ctx)
}

func (s *MetricsService) UpdateMetric(ctx context.Context, metric model.Metrics) error {
	return s.storage.Save(ctx, metric)
}
