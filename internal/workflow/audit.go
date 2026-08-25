package workflow

import (
	"fmt"
	"meter-sync/internal/domain"
	"meter-sync/internal/service"
	"strings"
)

type AuditTrail struct{ Service *service.Service }

func (a AuditTrail) Timeline(recordID string) ([]domain.AuditEvent, error) {
	events, err := a.Service.Audits(recordID)
	if err != nil {
		return nil, err
	}
	return events, nil
}
func (a AuditTrail) Describe(recordID string) (string, error) {
	events, err := a.Timeline(recordID)
	if err != nil {
		return "", err
	}
	if len(events) == 0 {
		return "no audit events", nil
	}
	parts := make([]string, len(events))
	for i, event := range events {
		parts[i] = fmt.Sprintf("%s:%s:%s", event.Action, event.Actor, event.OccurredAt)
	}
	return strings.Join(parts, " | "), nil
}
func (a AuditTrail) HasAction(recordID, action string) bool {
	events, err := a.Timeline(recordID)
	if err != nil {
		return false
	}
	for _, event := range events {
		if event.Action == action {
			return true
		}
	}
	return false
}
