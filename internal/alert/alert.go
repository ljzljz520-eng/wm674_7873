package alert

import (
	"fmt"
	"meter-sync/internal/domain"
	"sort"
	"strings"
)

type Severity string

const (
	Info     Severity = "info"
	Warning  Severity = "warning"
	Critical Severity = "critical"
)

type Alert struct {
	ID        string
	RecordID  string
	Severity  Severity
	Code      string
	Message   string
	Resolved  bool
	CreatedAt string
}
type Center struct{ alerts map[string]Alert }

func New() *Center { return &Center{alerts: make(map[string]Alert)} }
func (c *Center) Raise(record domain.Record, severity Severity, code, message, stamp string) Alert {
	id := fmt.Sprintf("%s:%s:%d", record.ID, code, record.Version)
	alert := Alert{ID: id, RecordID: record.ID, Severity: severity, Code: code, Message: message, CreatedAt: stamp}
	c.alerts[id] = alert
	return alert
}
func (c *Center) Resolve(id string) bool {
	alert, ok := c.alerts[id]
	if !ok {
		return false
	}
	alert.Resolved = true
	c.alerts[id] = alert
	return true
}
func (c *Center) Get(id string) (Alert, bool) { alert, ok := c.alerts[id]; return alert, ok }
func (c *Center) List(recordID string, includeResolved bool) []Alert {
	result := make([]Alert, 0)
	for _, alert := range c.alerts {
		if recordID != "" && alert.RecordID != recordID {
			continue
		}
		if !includeResolved && alert.Resolved {
			continue
		}
		result = append(result, alert)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
func (c *Center) CountOpen(recordID string) int { return len(c.List(recordID, false)) }
func (c *Center) Summary(recordID string) string {
	alerts := c.List(recordID, true)
	if len(alerts) == 0 {
		return "no alerts"
	}
	parts := make([]string, len(alerts))
	for i, alert := range alerts {
		state := "open"
		if alert.Resolved {
			state = "resolved"
		}
		parts[i] = fmt.Sprintf("%s:%s:%s", alert.Code, alert.Severity, state)
	}
	return strings.Join(parts, " | ")
}
