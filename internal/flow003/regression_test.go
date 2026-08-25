package flow003

import (
	"meter-sync/internal/domain"
	"meter-sync/internal/service"
	"meter-sync/internal/store"
	"path/filepath"
	"testing"
)

func Test674BusinessRegression(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc, err := service.New(st, service.FixedClock{Value: "t"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Register("meter-1", "Maker", "M1", "SYNC-1", "u", 100); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Register("meter-2", "Maker", "M2", "SYNC-2", "u", 200); err != nil {
		t.Fatal(err)
	}
	records, err := svc.SyncAmounts([]domain.ImportRow{{ID: "meter-1", Reference: "SYNC-1", Amount: 111}, {ID: "meter-2", Reference: "SYNC-2", Amount: 222}}, "u")
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("records=%v", records)
	}
	if records[0].Amount != 111 {
		t.Fatalf("first amount=%d", records[0].Amount)
	}
	if records[1].Amount != 222 {
		t.Fatalf("second amount=%d", records[1].Amount)
	}
}
