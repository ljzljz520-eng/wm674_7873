package retention

import (
	"meter-sync/internal/domain"
	"testing"
	"time"
)

func TestRetention(t *testing.T) {
	policy := DefaultPolicy()
	record := domain.Record{ID: "r", Status: domain.StatusDraft, CreatedAt: "2020-01-01T00:00:00Z"}
	if !policy.IsExpired(record, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("not expired")
	}
	if policy.Validate() != nil {
		t.Fatal("policy")
	}
}
