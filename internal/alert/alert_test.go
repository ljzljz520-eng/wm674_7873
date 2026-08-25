package alert

import (
	"meter-sync/internal/domain"
	"testing"
)

func TestAlertCenter(t *testing.T) {
	center := New()
	alert := center.Raise(domain.Record{ID: "r", Version: 1}, Warning, "AMOUNT", "check amount", "t")
	if center.CountOpen("r") != 1 {
		t.Fatal("open")
	}
	if !center.Resolve(alert.ID) {
		t.Fatal("resolve")
	}
	if center.CountOpen("r") != 0 {
		t.Fatal("still open")
	}
}
