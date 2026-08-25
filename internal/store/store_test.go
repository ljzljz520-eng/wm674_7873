package store

import (
	"meter-sync/internal/domain"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "meters.db")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record, _ := domain.NewRecord("r1", "Maker", "M1", "S1", "u", "t", 42)
	if err := st.SaveRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	st, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	loaded, err := st.GetRecord("r1")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Amount != 42 || loaded.Manufacturer != "Maker" {
		t.Fatalf("loaded=%+v", loaded)
	}
}

func TestStoreListsAndUpdates(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	r, _ := domain.NewRecord("r1", "Maker", "M1", "S1", "u", "t", 4)
	if err = st.SaveRecord(r); err != nil {
		t.Fatal(err)
	}
	updated, err := st.UpdateAmount("r1", 1, 8, "u", "t2")
	if err != nil || updated.Amount != 8 {
		t.Fatalf("updated=%+v err=%v", updated, err)
	}
	records, err := st.ListRecords(domain.Query{Limit: 1})
	if err != nil || len(records) != 1 {
		t.Fatalf("records=%v err=%v", records, err)
	}
}
