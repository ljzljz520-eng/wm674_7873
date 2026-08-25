package service

import (
	"meter-sync/internal/domain"
	"meter-sync/internal/store"
	"path/filepath"
	"testing"
)

func testService(t *testing.T) *Service {
	t.Helper()
	st, err := store.Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	svc, err := New(st, FixedClock{Value: "t"})
	if err != nil {
		t.Fatal(err)
	}
	return svc
}
func TestRegisterReviewPublish(t *testing.T) {
	svc := testService(t)
	if _, err := svc.Register("r", "Maker", "M", "S", "u", 5); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Review("r", "reviewer", "ok"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Publish("r", "publisher"); err != nil {
		t.Fatal(err)
	}
}
func TestImportAndSearch(t *testing.T) {
	svc := testService(t)
	result := svc.Import([]domain.ImportRow{{ID: "r1", Manufacturer: "A", Model: "M", Reference: "S1", Amount: 1}, {ID: "r2", Manufacturer: "B", Model: "M", Reference: "S2", Amount: 2}}, "u")
	if len(result.Accepted) != 2 || len(result.Rejected) != 0 {
		t.Fatalf("result=%+v", result)
	}
	records, err := svc.Search(domain.Query{Manufacturer: "a"})
	if err != nil || len(records) != 1 {
		t.Fatalf("records=%v err=%v", records, err)
	}
}
