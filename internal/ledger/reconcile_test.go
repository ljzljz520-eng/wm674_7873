package ledger

import (
	"meter-sync/internal/domain"
	"testing"
)

func TestReconcile(t *testing.T) {
	record := domain.Record{ID: "r", Amount: 10}
	differences := CompareMany([]domain.Record{record}, []Entry{{RecordID: "r", Amount: 9}})
	if len(Unbalanced(differences)) != 1 || TotalDelta(differences) != -1 {
		t.Fatal(differences)
	}
}
