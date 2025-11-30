package repository

import (
	_interface "github.com/bryzgaalov/metrics-collector/internal/agent/interface"
	"testing"
)

func TestNewMetricsCollector_ImplementsInterfaceAndNotNil(t *testing.T) {
	var repo _interface.MetricsCollectorRepository = NewMetricsCollector()

	if repo == nil {
		t.Fatalf("expected NewMetricsCollector to return non-nil value")
	}
}
func TestNewMetricsCollector_EmptySnapshot(t *testing.T) {
	repo := NewMetricsCollector()

	gauges, counters := repo.Snapshot()

	if len(gauges) != 0 {
		t.Fatalf("expected empty gauges on new collector, got %d", len(gauges))
	}
	if len(counters) != 0 {
		t.Fatalf("expected empty counters on new collector, got %d", len(counters))
	}
}
func TestNewMetricsCollector_InstancesAreIsolated(t *testing.T) {
	first := NewMetricsCollector()
	second := NewMetricsCollector()

	first.SetGauge("Alloc", 100)
	first.AddCounter("PollCount", 10)

	g1, c1 := first.Snapshot()
	g2, c2 := second.Snapshot()

	if g1["Alloc"] != 100 {
		t.Fatalf("expected first collector Alloc=100, got %v", g1["Alloc"])
	}

	if _, ok := g2["Alloc"]; ok {
		t.Fatalf("expected second collector to have no 'Alloc' gauge, but it exists")
	}

	if c1["PollCount"] != 10 {
		t.Fatalf("expected first collector PollCount=10, got %v", c1["PollCount"])
	}

	if _, ok := c2["PollCount"]; ok {
		t.Fatalf("expected second collector to have no 'PollCount' counter, but it exists")
	}
}

func TestSetGauge_AddsNewMetric(t *testing.T) {
	repo := NewMetricsCollector()

	repo.SetGauge("Alloc", 123.456)

	gauges, _ := repo.Snapshot()
	val, ok := gauges["Alloc"]
	if !ok {
		t.Fatalf("expected gauge 'Alloc' to exist after SetGauge, but it was not found")
	}
	if val != 123.456 {
		t.Fatalf("expected gauge 'Alloc' to be 123.456, got %v", val)
	}
}
func TestSetGauge_OverridesExistingMetric(t *testing.T) {
	repo := NewMetricsCollector()

	repo.SetGauge("Alloc", 10.0)
	repo.SetGauge("Alloc", 20.5)

	gauges, _ := repo.Snapshot()
	val, ok := gauges["Alloc"]
	if !ok {
		t.Fatalf("expected gauge 'Alloc' to exist, but it was not found")
	}
	if val != 20.5 {
		t.Fatalf("expected gauge 'Alloc' to be 20.5 after override, got %v", val)
	}
}
func TestSetGauge_EmptyNameAllowed(t *testing.T) {
	repo := NewMetricsCollector()

	repo.SetGauge("", 42.0)

	gauges, _ := repo.Snapshot()
	val, ok := gauges[""]
	if !ok {
		t.Fatalf("expected gauge with empty name to exist, but it was not found")
	}
	if val != 42.0 {
		t.Fatalf("expected gauge with empty name to be 42.0, got %v", val)
	}
}

func TestAddCounter_AddsNewCounter(t *testing.T) {
	repo := NewMetricsCollector()

	repo.AddCounter("PollCount", 5)

	_, counters := repo.Snapshot()
	val, ok := counters["PollCount"]
	if !ok {
		t.Fatalf("expected counter 'PollCount' to exist after AddCounter, but it was not found")
	}
	if val != 5 {
		t.Fatalf("expected counter 'PollCount' to be 5, got %v", val)
	}
}
func TestAddCounter_AccumulatesDelta(t *testing.T) {
	repo := NewMetricsCollector()

	repo.AddCounter("PollCount", 3)
	repo.AddCounter("PollCount", 7)

	_, counters := repo.Snapshot()
	val := counters["PollCount"]
	if val != 10 {
		t.Fatalf("expected counter 'PollCount' to be 10 after two calls, got %v", val)
	}
}
func TestAddCounter_AllowsNegativeDelta(t *testing.T) {
	repo := NewMetricsCollector()

	repo.AddCounter("Errors", 10)
	repo.AddCounter("Errors", -4)

	_, counters := repo.Snapshot()
	val := counters["Errors"]
	if val != 6 {
		t.Fatalf("expected counter 'Errors' to be 6 after adding 10 and -4, got %v", val)
	}
}
func TestAddCounter_EmptyNameAllowed(t *testing.T) {
	repo := NewMetricsCollector()

	repo.AddCounter("", 1)

	_, counters := repo.Snapshot()
	val, ok := counters[""]
	if !ok {
		t.Fatalf("expected counter with empty name to exist, but it was not found")
	}
	if val != 1 {
		t.Fatalf("expected counter with empty name to be 1, got %v", val)
	}
}

func TestSnapshot_Empty(t *testing.T) {
	repo := NewMetricsCollector()

	g, c := repo.Snapshot()

	if len(g) != 0 {
		t.Fatalf("expected empty gauges map, got %d elements", len(g))
	}

	if len(c) != 0 {
		t.Fatalf("expected empty counters map, got %d elements", len(c))
	}
}
func TestSnapshot_ReturnsExistingValues(t *testing.T) {
	repo := NewMetricsCollector()

	repo.SetGauge("Alloc", 100.5)
	repo.AddCounter("PollCount", 3)

	g, c := repo.Snapshot()

	if g["Alloc"] != 100.5 {
		t.Fatalf("expected gauge Alloc = 100.5, got %v", g["Alloc"])
	}

	if c["PollCount"] != 3 {
		t.Fatalf("expected counter PollCount = 3, got %v", c["PollCount"])
	}
}
func TestSnapshot_ReturnsCopies(t *testing.T) {
	repo := NewMetricsCollector()

	repo.SetGauge("Alloc", 10.0)
	repo.AddCounter("PollCount", 1)

	g, c := repo.Snapshot()

	g["Alloc"] = 9999     // меняем snapshot
	c["PollCount"] = 7777 // меняем snapshot

	g2, c2 := repo.Snapshot()

	if g2["Alloc"] == 9999 {
		t.Fatalf("expected original gauge to be unchanged, but Alloc was modified")
	}

	if c2["PollCount"] == 7777 {
		t.Fatalf("expected original counter to remain unchanged, but PollCount was modified")
	}
}
func TestSnapshot_SeparatesGaugeAndCounter(t *testing.T) {
	repo := NewMetricsCollector()

	repo.SetGauge("HeapAlloc", 42.0)
	repo.AddCounter("Errors", 13)

	g, c := repo.Snapshot()

	if _, ok := g["Errors"]; ok {
		t.Fatalf("expected no counters inside gauges map")
	}

	if _, ok := c["HeapAlloc"]; ok {
		t.Fatalf("expected no gauges inside counters map")
	}
}
