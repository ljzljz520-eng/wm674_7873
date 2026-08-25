package notification

import (
	"meter-sync/internal/domain"
	"testing"
)

func TestOutbox(t *testing.T) {
	outbox := New()
	message := outbox.Compose(domain.Record{ID: "r", Manufacturer: "M", MeterModel: "X", SyncReference: "S", Amount: 1, Status: domain.StatusDraft}, "registered")
	if outbox.PendingCount() != 1 {
		t.Fatal("pending")
	}
	if err := outbox.Send(message.ID, "t"); err != nil {
		t.Fatal(err)
	}
	if len(outbox.Sent("r")) != 1 {
		t.Fatal("sent")
	}
}
