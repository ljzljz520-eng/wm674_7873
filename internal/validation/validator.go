package validation

import (
	"fmt"
	"meter-sync/internal/domain"
	"strings"
)

type Issue struct {
	Field   string
	Message string
}

func ValidateRecord(record domain.Record) []Issue {
	issues := make([]Issue, 0)
	if strings.TrimSpace(record.Manufacturer) == "" {
		issues = append(issues, Issue{"manufacturer", "is required"})
	}
	if strings.TrimSpace(record.MeterModel) == "" {
		issues = append(issues, Issue{"model", "is required"})
	}
	if record.Amount < 0 {
		issues = append(issues, Issue{"amount", "must not be negative"})
	}
	if record.Status == domain.StatusArchived && record.ArchivedAt == "" {
		issues = append(issues, Issue{"archived_at", "required for archived record"})
	}
	return issues
}
func Explain(issues []Issue) string {
	if len(issues) == 0 {
		return "valid"
	}
	parts := make([]string, len(issues))
	for i, issue := range issues {
		parts[i] = fmt.Sprintf("%s %s", issue.Field, issue.Message)
	}
	return strings.Join(parts, "; ")
}
func IsReadyForPublication(record domain.Record) bool {
	return record.Status == domain.StatusReviewed && len(ValidateRecord(record)) == 0
}
