package interfaces

import (
	models "github.com/bryzgaalov/metrics-collector/internal/model"
)

type MetricsStorageInterface interface {
	Save(metric models.Metrics) error
	Get(id string) (models.Metrics, bool)
	GetAll() map[string]models.Metrics
}
