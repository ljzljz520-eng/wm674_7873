package service

import (
	"fmt"
	"sort"

	"meter-sync/internal/domain"
)

type Summary struct {
	Total, Draft, Reviewed, Published, Archived int
	Amount                                      int64
}

func Summarize(records []domain.Record) Summary {
	summary := Summary{}
	for _, record := range records {
		summary.Total++
		summary.Amount += record.Amount
		switch record.Status {
		case domain.StatusDraft:
			summary.Draft++
		case domain.StatusReviewed:
			summary.Reviewed++
		case domain.StatusPublished:
			summary.Published++
		case domain.StatusArchived:
			summary.Archived++
		}
	}
	return summary
}
func (s Summary) String() string {
	return fmt.Sprintf("total=%d draft=%d reviewed=%d published=%d archived=%d amount=%d", s.Total, s.Draft, s.Reviewed, s.Published, s.Archived, s.Amount)
}
func SortByAmount(records []domain.Record) []domain.Record {
	copyRecords := append([]domain.Record(nil), records...)
	sort.SliceStable(copyRecords, func(i, j int) bool {
		if copyRecords[i].Amount == copyRecords[j].Amount {
			return copyRecords[i].ID < copyRecords[j].ID
		}
		return copyRecords[i].Amount > copyRecords[j].Amount
	})
	return copyRecords
}
