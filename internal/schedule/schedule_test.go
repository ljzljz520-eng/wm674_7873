package schedule

import "testing"

func TestPlanner(t *testing.T) {
	planner := New()
	if err := planner.Add(Job{ID: "j", Frequency: Daily, Enabled: true, NextRun: "2026-01-01"}); err != nil {
		t.Fatal(err)
	}
	if len(planner.Due("2026-01-02")) != 1 {
		t.Fatal("due")
	}
	if !planner.Disable("j") || planner.ActiveCount() != 0 {
		t.Fatal("disable")
	}
}
