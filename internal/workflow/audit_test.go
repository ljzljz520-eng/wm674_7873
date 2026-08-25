package workflow

import (
	"meter-sync/internal/service"
	"meter-sync/internal/store"
	"path/filepath"
	"testing"
)

func TestAuditTrail(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc, _ := service.New(st, service.FixedClock{Value: "t"})
	_, _ = svc.Register("r", "Maker", "M", "S", "u", 1)
	trail := AuditTrail{Service: svc}
	if !trail.HasAction("r", "register") {
		t.Fatal("missing action")
	}
	text, err := trail.Describe("r")
	if err != nil || text == "" {
		t.Fatalf("text=%q err=%v", text, err)
	}
}
