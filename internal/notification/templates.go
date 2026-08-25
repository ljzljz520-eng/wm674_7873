package notification

import (
	"fmt"
	"meter-sync/internal/domain"
	"strings"
)

type Template struct {
	Name    string
	Subject string
	Body    string
	Enabled bool
}

func DefaultTemplates() []Template {
	return []Template{{Name: "review-request", Subject: "Review meter sync %s", Body: "Please review synchronization record %s for %s.", Enabled: true}, {Name: "publish-notice", Subject: "Published meter sync %s", Body: "Synchronization record %s is published.", Enabled: true}, {Name: "archive-notice", Subject: "Archived meter sync %s", Body: "Synchronization record %s is archived.", Enabled: true}, {Name: "amount-alert", Subject: "Amount check for %s", Body: "Amount %d requires attention.", Enabled: true}}
}
func Render(template Template, record domain.Record) Message {
	subject := fmt.Sprintf(template.Subject, record.ID)
	body := fmt.Sprintf(template.Body, record.ID, record.Manufacturer)
	if strings.Contains(template.Name, "amount") {
		body = fmt.Sprintf(template.Body, record.Amount)
	}
	return Message{RecordID: record.ID, Subject: subject, Body: body, Channel: Console}
}
func FindTemplate(name string) Template {
	for _, template := range DefaultTemplates() {
		if template.Name == name {
			return template
		}
	}
	return Template{}
}
func EnabledTemplates() []Template {
	result := make([]Template, 0)
	for _, template := range DefaultTemplates() {
		if template.Enabled {
			result = append(result, template)
		}
	}
	return result
}
