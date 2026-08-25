package service

import (
	"meter-sync/internal/domain"
	"testing"
)

func TestMetrics(t *testing.T) {
	metrics := Metrics([]domain.Record{{Manufacturer: "A", Amount: 2, Version: 1}, {Manufacturer: "A", Amount: 3, Version: 2}, {Manufacturer: "B", Amount: 8, Version: 1}})
	largest, ok := LargestMetric(metrics)
	if !ok || largest.Manufacturer != "B" {
		t.Fatalf("largest=%v ok=%v", largest, ok)
	}
}
