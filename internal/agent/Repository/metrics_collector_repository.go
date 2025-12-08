package repository

import (
	"sync"

	agentiface "github.com/bryzgaalov/metrics-collector/internal/agent/interface"
)

type MetricsCollector struct {
	m        sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
}

func NewMetricsCollector() agentiface.MetricsCollectorRepository {
	return &MetricsCollector{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (s *MetricsCollector) SetGauge(name string, value float64) {
	s.m.Lock()
	s.gauges[name] = value
	s.m.Unlock()
}

func (s *MetricsCollector) AddCounter(name string, delta int64) {
	s.m.Lock()
	s.counters[name] += delta
	s.m.Unlock()
}

func (s *MetricsCollector) Snapshot() (map[string]float64, map[string]int64) {
	s.m.RLock()
	defer s.m.RUnlock()

	g := make(map[string]float64, len(s.gauges))
	for k, v := range s.gauges {
		g[k] = v
	}

	c := make(map[string]int64, len(s.counters))
	for k, v := range s.counters {
		c[k] = v
	}

	return g, c
}
