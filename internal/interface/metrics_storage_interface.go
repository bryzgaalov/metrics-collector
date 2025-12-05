package interfaces

import (
	"context"

	models "github.com/bryzgaalov/metrics-collector/internal/model"
)

type MetricsStorage interface {
	Save(ctx context.Context, metric models.Metrics) error
	Get(ctx context.Context, id string) (models.Metrics, bool)
	GetAll(ctx context.Context) (map[string]models.Metrics, error)
}
