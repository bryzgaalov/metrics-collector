package repository

import (
	models "github.com/bryzgaalov/metrics-collector/internal/model"
)

type MemStorage struct {
	data map[string]models.Metrics
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]models.Metrics),
	}
}

func (m *MemStorage) Save(metric models.Metrics) error {
	m.data[metric.ID] = metric
	return nil
}

func (m *MemStorage) Get(id string) (models.Metrics, bool) {
	v, ok := m.data[id]
	return v, ok
}

func (m *MemStorage) GetAll() map[string]models.Metrics {
	return m.data
}
