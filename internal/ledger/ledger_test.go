package ledger

import (
	"meter-sync/internal/domain"
	"testing"
)

func TestLedgerEntries(t *testing.T) {
	ledger := New()
	record := domain.Record{ID: "r", Amount: 10, Currency: "CNY", SyncReference: "S", Version: 1}
	entry, err := ledger.Add(record, "sync", "t")
	if err != nil {
		t.Fatal(err)
	}
	if ledger.Total("r") != 10 {
		t.Fatal("total")
	}
	if err := ledger.MarkReconciled(entry.ID); err != nil {
		t.Fatal(err)
	}
	if len(ledger.Unreconciled()) != 0 {
		t.Fatal("unreconciled")
	}
}
func TestRuleSet(t *testing.T) {
	rules := DefaultRules()
	if !rules.Valid(10, "CNY") {
		t.Fatal(rules.Explain(10, "CNY"))
	}
	if rules.Valid(0, "CNY") {
		t.Fatal("zero should fail")
	}
}
