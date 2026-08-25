package service

import (
	"meter-sync/internal/domain"
	"testing"
)

func TestBatchValidation(t *testing.T) {
	svc := testService(t)
	summary := svc.DryRun([]domain.ImportRow{{ID: "", Amount: 1}, {ID: "r", Manufacturer: "M", Model: "X", Reference: "S", Amount: -1}})
	if summary.Accepted != 0 || summary.Rejected != 2 {
		t.Fatalf("summary=%+v", summary)
	}
}
func TestProcessBatch(t *testing.T) {
	svc := testService(t)
	summary := svc.ProcessBatch([]domain.ImportRow{{ID: "r", Manufacturer: "M", Model: "X", Reference: "S", Amount: 5}}, "u")
	if summary.Accepted != 1 || summary.TotalAmount != 5 {
		t.Fatalf("summary=%+v", summary)
	}
}
