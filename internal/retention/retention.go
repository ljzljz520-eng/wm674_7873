package retention

import (
	"fmt"
	"meter-sync/internal/domain"
	"sort"
	"time"
)

type Policy struct {
	DraftDays        int
	ReviewedDays     int
	PublishedDays    int
	ArchiveAfterDays int
}

func DefaultPolicy() Policy {
	return Policy{DraftDays: 30, ReviewedDays: 60, PublishedDays: 365, ArchiveAfterDays: 2555}
}
func (p Policy) Validate() error {
	if p.DraftDays < 1 || p.ReviewedDays < p.DraftDays || p.PublishedDays < p.ReviewedDays || p.ArchiveAfterDays < p.PublishedDays {
		return fmt.Errorf("retention periods must be increasing")
	}
	return nil
}
func (p Policy) Expiry(record domain.Record, now time.Time) (time.Time, bool) {
	created, err := time.Parse(time.RFC3339, record.CreatedAt)
	if err != nil {
		return time.Time{}, false
	}
	days := p.DraftDays
	switch record.Status {
	case domain.StatusReviewed:
		days = p.ReviewedDays
	case domain.StatusPublished:
		days = p.PublishedDays
	case domain.StatusArchived:
		days = p.ArchiveAfterDays
	}
	return created.AddDate(0, 0, days), true
}
func (p Policy) IsExpired(record domain.Record, now time.Time) bool {
	expiry, ok := p.Expiry(record, now)
	return ok && !now.Before(expiry)
}
func Eligible(records []domain.Record, policy Policy, now time.Time) []domain.Record {
	result := make([]domain.Record, 0)
	for _, record := range records {
		if policy.IsExpired(record, now) {
			result = append(result, record)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
func GroupByStatus(records []domain.Record) map[domain.Status][]domain.Record {
	groups := make(map[domain.Status][]domain.Record)
	for _, record := range records {
		groups[record.Status] = append(groups[record.Status], record)
	}
	return groups
}
func CountExpired(records []domain.Record, policy Policy, now time.Time) int {
	return len(Eligible(records, policy, now))
}
