package domain

type AuditEvent struct {
	ID         string `json:"id"`
	RecordID   string `json:"record_id"`
	Action     string `json:"action"`
	Actor      string `json:"actor"`
	Details    string `json:"details"`
	OccurredAt string `json:"occurred_at"`
}

func NewAudit(id, recordID, action, actor, details, stamp string) AuditEvent {
	return AuditEvent{ID: id, RecordID: recordID, Action: action, Actor: actor, Details: details, OccurredAt: stamp}
}

type Workflow struct {
	ID        string   `json:"id"`
	RecordID  string   `json:"record_id"`
	Name      string   `json:"name"`
	State     string   `json:"state"`
	Steps     []string `json:"steps"`
	UpdatedAt string   `json:"updated_at"`
}

func NewWorkflow(id, recordID, name, stamp string) Workflow {
	return Workflow{ID: id, RecordID: recordID, Name: name, State: "created", Steps: []string{"create"}, UpdatedAt: stamp}
}

func (w Workflow) Advance(step, stamp string) Workflow {
	w.Steps = append(append([]string(nil), w.Steps...), step)
	w.State = step
	w.UpdatedAt = stamp
	return w
}

type Attachment struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Name      string `json:"name"`
	MediaType string `json:"media_type"`
	Digest    string `json:"digest"`
	AddedAt   string `json:"added_at"`
}

func NewAttachment(id, recordID, name, mediaType, digest, stamp string) Attachment {
	return Attachment{ID: id, RecordID: recordID, Name: name, MediaType: mediaType, Digest: digest, AddedAt: stamp}
}
