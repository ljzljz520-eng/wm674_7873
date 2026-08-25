package policy

import (
	"meter-sync/internal/domain"
	"testing"
)

func TestDefaultPolicy(t *testing.T) {
	record := domain.Record{ID: "r", SyncReference: "S", Currency: "CNY", Version: 1}
	if err := Default().Check(record); err != nil {
		t.Fatal(err)
	}
	record.Currency = "USD"
	if err := Default().Check(record); err == nil {
		t.Fatal("expected policy failure")
	}
}
