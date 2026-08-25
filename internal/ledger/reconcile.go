package ledger

import (
	"fmt"
	"meter-sync/internal/domain"
	"sort"
)

type Difference struct {
	RecordID string
	Expected int64
	Actual   int64
	Delta    int64
	Balanced bool
}

func Compare(record domain.Record, entries []Entry) Difference {
	var actual int64
	for _, entry := range entries {
		if entry.RecordID == record.ID {
			actual += entry.Amount
		}
	}
	return Difference{RecordID: record.ID, Expected: record.Amount, Actual: actual, Delta: actual - record.Amount, Balanced: actual == record.Amount}
}
func CompareMany(records []domain.Record, entries []Entry) []Difference {
	result := make([]Difference, 0, len(records))
	for _, record := range records {
		result = append(result, Compare(record, entries))
	}
	sort.Slice(result, func(i, j int) bool { return result[i].RecordID < result[j].RecordID })
	return result
}
func Unbalanced(differences []Difference) []Difference {
	result := make([]Difference, 0)
	for _, difference := range differences {
		if !difference.Balanced {
			result = append(result, difference)
		}
	}
	return result
}
func ReconciliationMessage(difference Difference) string {
	if difference.Balanced {
		return fmt.Sprintf("%s balanced", difference.RecordID)
	}
	return fmt.Sprintf("%s delta=%d", difference.RecordID, difference.Delta)
}
func TotalDelta(differences []Difference) int64 {
	var total int64
	for _, difference := range differences {
		total += difference.Delta
	}
	return total
}
