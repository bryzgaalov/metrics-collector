package repository

import (
	"context"
	"testing"

	models "github.com/bryzgaalov/metrics-collector/internal/model"
)

func TestNewMemStorage_NotNil(t *testing.T) {
	s := NewMemStorage()
	if s == nil {
		t.Fatalf("expected NewMemStorage to return non-nil storage")
	}
	if s.data == nil {
		t.Fatalf("expected internal data map to be initialized")
	}
}

func TestMemStorage_SaveAndGet(t *testing.T) {
	s := NewMemStorage()
	metric := models.Metrics{
		ID:    "Alloc",
		MType: models.Gauge,
		Value: float64Ptr(123.456),
	}

	err := s.Save(context.Background(), metric)
	if err != nil {
		t.Fatalf("unexpected error from Save: %v", err)
	}

	got, ok := s.Get(context.Background(), "Alloc")
	if !ok {
		t.Fatalf("expected metric 'Alloc' to exist, but it was not found")
	}
	if *got.Value != 123.456 {
		t.Fatalf("expected value 123.456, got %v", *got.Value)
	}
}

func float64Ptr(v float64) *float64 { return &v }

func TestMemStorage_SaveOverridesExisting(t *testing.T) {
	s := NewMemStorage()

	m1 := models.Metrics{ID: "Counter", MType: models.Counter, Delta: int64Ptr(2)}
	m2 := models.Metrics{ID: "Counter", MType: models.Counter, Delta: int64Ptr(10)}

	_ = s.Save(context.Background(), m1)
	_ = s.Save(context.Background(), m2)

	got, ok := s.Get(context.Background(), "Counter")
	if !ok {
		t.Fatalf("expected metric 'Counter' to exist")
	}
	if *got.Delta != 10 {
		t.Fatalf("expected counter to be overwritten with 10, got %v", *got.Delta)
	}
}

func int64Ptr(v int64) *int64 { return &v }

func TestMemStorage_GetNotExisting(t *testing.T) {
	s := NewMemStorage()

	_, ok := s.Get(context.Background(), "UnknownMetric")
	if ok {
		t.Fatalf("expected ok==false for missing metric")
	}
}
