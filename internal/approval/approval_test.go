package approval

import (
	"meter-sync/internal/domain"
	"testing"
)

func TestApprovalQueue(t *testing.T) {
	queue := New()
	request, err := queue.Submit(domain.Record{ID: "r", Version: 1}, "reviewer", "amount", "t")
	if err != nil {
		t.Fatal(err)
	}
	if err = queue.Decide(request.ID, Approved, "t2"); err != nil {
		t.Fatal(err)
	}
	if !queue.HasApproved("r") {
		t.Fatal("not approved")
	}
}
