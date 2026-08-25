package notification

import (
	"errors"
	"fmt"
	"meter-sync/internal/domain"
	"sort"
	"strings"
)

type Channel string

const (
	Email   Channel = "email"
	Console Channel = "console"
	Webhook Channel = "webhook"
)

type Message struct {
	ID        string
	RecordID  string
	Channel   Channel
	Recipient string
	Subject   string
	Body      string
	Sent      bool
	SentAt    string
}
type Outbox struct{ messages map[string]Message }

func New() *Outbox { return &Outbox{messages: make(map[string]Message)} }
func (o *Outbox) Queue(record domain.Record, channel Channel, recipient, subject, body string) Message {
	id := fmt.Sprintf("message:%s:%s:%d", record.ID, channel, record.Version)
	message := Message{ID: id, RecordID: record.ID, Channel: channel, Recipient: recipient, Subject: subject, Body: body}
	o.messages[id] = message
	return message
}
func (o *Outbox) Get(id string) (Message, bool) { message, ok := o.messages[id]; return message, ok }
func (o *Outbox) Send(id, stamp string) error {
	message, ok := o.messages[id]
	if !ok {
		return errors.New("message not found")
	}
	if strings.TrimSpace(message.Recipient) == "" {
		return errors.New("recipient is required")
	}
	message.Sent = true
	message.SentAt = stamp
	o.messages[id] = message
	return nil
}
func (o *Outbox) Pending(recordID string) []Message {
	result := make([]Message, 0)
	for _, message := range o.messages {
		if !message.Sent && (recordID == "" || message.RecordID == recordID) {
			result = append(result, message)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
func (o *Outbox) Sent(recordID string) []Message {
	result := make([]Message, 0)
	for _, message := range o.messages {
		if message.Sent && (recordID == "" || message.RecordID == recordID) {
			result = append(result, message)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
func (o *Outbox) PendingCount() int { return len(o.Pending("")) }
func (o *Outbox) SendAll(stamp string) int {
	count := 0
	for _, message := range o.Pending("") {
		if o.Send(message.ID, stamp) == nil {
			count++
		}
	}
	return count
}
func (o *Outbox) Compose(record domain.Record, action string) Message {
	subject := fmt.Sprintf("meter sync %s: %s", action, record.ID)
	body := fmt.Sprintf("manufacturer=%s model=%s reference=%s amount=%d status=%s", record.Manufacturer, record.MeterModel, record.SyncReference, record.Amount, record.Status)
	return o.Queue(record, Console, "operator", subject, body)
}
