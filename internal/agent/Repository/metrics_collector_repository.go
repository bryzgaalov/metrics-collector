package repository

import agentiface "github.com/bryzgaalov/metrics-collector/internal/agent/interface"

type MetricsCollector struct {
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
	s.gauges[name] = value
}

func (s *MetricsCollector) AddCounter(name string, delta int64) {
	s.counters[name] += delta
}

func (s *MetricsCollector) Snapshot() (map[string]float64, map[string]int64) {
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
