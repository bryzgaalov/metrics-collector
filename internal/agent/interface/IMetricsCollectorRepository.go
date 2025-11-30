package _interface

type MetricsCollectorRepository interface {
	SetGauge(name string, value float64)
	AddCounter(name string, delta int64)
	Snapshot() (map[string]float64, map[string]int64)
}
