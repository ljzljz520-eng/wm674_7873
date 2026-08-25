package domain

import "testing"

func TestRecordTransitions(t *testing.T) {
	record, err := NewRecord("r1", "Maker", "M1", "S1", "u", "t", 10)
	if err != nil {
		t.Fatal(err)
	}
	record, err = record.MarkReviewed("reviewer", "ok", "t2")
	if err != nil || record.Status != StatusReviewed {
		t.Fatalf("review=%+v err=%v", record, err)
	}
	record, err = record.MarkPublished("t3")
	if err != nil || record.Status != StatusPublished {
		t.Fatalf("publish=%+v err=%v", record, err)
	}
	record, err = record.MarkArchived("t4")
	if err != nil || record.Status != StatusArchived {
		t.Fatalf("archive=%+v err=%v", record, err)
	}
}

func TestQueryMatches(t *testing.T) {
	record := Record{ID: "r", Manufacturer: "Alpha Meter", MeterModel: "X", SyncReference: "ref", Amount: 1, Currency: "CNY", Status: StatusDraft, Version: 1}
	if !(Query{Manufacturer: "alpha"}).Matches(record) {
		t.Fatal("manufacturer query")
	}
	if (Query{Status: StatusPublished}).Matches(record) {
		t.Fatal("status query")
	}
}
