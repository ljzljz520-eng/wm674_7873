package schedule

import (
	"fmt"
	"meter-sync/internal/domain"
	"sort"
	"strings"
)

type Frequency string

const (
	Daily  Frequency = "daily"
	Weekly Frequency = "weekly"
	Manual Frequency = "manual"
)

type Job struct {
	ID        string
	Name      string
	Frequency Frequency
	Enabled   bool
	LastRun   string
	NextRun   string
	Owner     string
	RecordIDs []string
}
type Planner struct{ jobs map[string]Job }

func New() *Planner { return &Planner{jobs: make(map[string]Job)} }
func (p *Planner) Add(job Job) error {
	if strings.TrimSpace(job.ID) == "" {
		return fmt.Errorf("job id is required")
	}
	if job.Frequency != Daily && job.Frequency != Weekly && job.Frequency != Manual {
		return fmt.Errorf("unsupported frequency %s", job.Frequency)
	}
	if _, ok := p.jobs[job.ID]; ok {
		return fmt.Errorf("job %s exists", job.ID)
	}
	job.RecordIDs = append([]string(nil), job.RecordIDs...)
	p.jobs[job.ID] = job
	return nil
}
func (p *Planner) Get(id string) (Job, bool) { job, ok := p.jobs[id]; return job, ok }
func (p *Planner) Enable(id string) bool {
	job, ok := p.jobs[id]
	if !ok {
		return false
	}
	job.Enabled = true
	p.jobs[id] = job
	return true
}
func (p *Planner) Disable(id string) bool {
	job, ok := p.jobs[id]
	if !ok {
		return false
	}
	job.Enabled = false
	p.jobs[id] = job
	return true
}
func (p *Planner) Due(stamp string) []Job {
	result := make([]Job, 0)
	for _, job := range p.jobs {
		if job.Enabled && (job.NextRun == "" || job.NextRun <= stamp) {
			result = append(result, job)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
func (p *Planner) Complete(id, last, next string) error {
	job, ok := p.jobs[id]
	if !ok {
		return fmt.Errorf("job %s not found", id)
	}
	job.LastRun = last
	job.NextRun = next
	p.jobs[id] = job
	return nil
}
func (p *Planner) Assign(id string, records []domain.Record) error {
	job, ok := p.jobs[id]
	if !ok {
		return fmt.Errorf("job %s not found", id)
	}
	job.RecordIDs = job.RecordIDs[:0]
	for _, record := range records {
		if record.ID != "" {
			job.RecordIDs = append(job.RecordIDs, record.ID)
		}
	}
	p.jobs[id] = job
	return nil
}
func (p *Planner) ActiveCount() int {
	count := 0
	for _, job := range p.jobs {
		if job.Enabled {
			count++
		}
	}
	return count
}
