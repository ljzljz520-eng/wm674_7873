package catalog

import "testing"

func TestCompatibility(t *testing.T) {
	result := CheckCompatibility("HD", "HX-100", "CNY")
	if !result.Compatible {
		t.Fatal(result)
	}
	if CheckCompatibility("HD", "bad", "CNY").Compatible {
		t.Fatal("bad model")
	}
}
