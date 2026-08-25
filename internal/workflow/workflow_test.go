package workflow

import (
	"meter-sync/internal/domain"
	"meter-sync/internal/service"
	"meter-sync/internal/store"
	"path/filepath"
	"strings"
	"testing"
)

func workflowService(t *testing.T) *service.Service {
	t.Helper()
	st, err := store.Open(filepath.Join(t.TempDir(), "db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	svc, err := service.New(st, service.FixedClock{Value: "t"})
	if err != nil {
		t.Fatal(err)
	}
	return svc
}
func TestWorkflowCreateReviewArchive(t *testing.T) {
	svc := workflowService(t)
	life := Lifecycle{Service: svc}
	record, err := life.CreateReviewPublishArchive("r", "Maker", "M", "S", "u", 9)
	if err != nil || record.Status != "archived" {
		t.Fatalf("record=%+v err=%v", record, err)
	}
}
func TestWorkflowSearchUpdatePublish(t *testing.T) {
	svc := workflowService(t)
	if _, err := svc.Register("r", "Maker", "M", "S", "u", 9); err != nil {
		t.Fatal(err)
	}
	flow := SearchUpdate{Service: svc}
	if _, err := flow.FindAndUpdate("Maker", 10, "u"); err != nil {
		t.Fatal(err)
	}
	if _, err := flow.PublishUpdated("r", "u"); err != nil {
		t.Fatal(err)
	}
}

func TestWorkflowImportReport(t *testing.T) {
	svc := workflowService(t)
	flow := ImportReport{Service: svc}
	result, err := flow.ImportCSV(strings.NewReader("r1,Maker,M,S1,11\nr2,Maker,N,S2,13\n"), "u")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Accepted) != 2 || len(result.Rejected) != 0 {
		t.Fatalf("result=%+v", result)
	}
	summary, err := flow.Report(domain.Query{Manufacturer: "Maker"})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Total != 2 || summary.Amount != 24 {
		t.Fatalf("summary=%+v", summary)
	}
}
