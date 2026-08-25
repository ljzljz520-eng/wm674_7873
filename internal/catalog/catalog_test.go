package catalog

import "testing"

func TestCatalogValidation(t *testing.T) {
	profile, err := Find("HD")
	if err != nil {
		t.Fatal(err)
	}
	if err := Validate(profile, "HX-100", "CNY"); err != nil {
		t.Fatal(err)
	}
	if err := Validate(profile, "unknown", "CNY"); err == nil {
		t.Fatal("expected model failure")
	}
}
func TestRegistry(t *testing.T) {
	registry := NewRegistry()
	if _, ok := registry.Get("hd"); !ok {
		t.Fatal("lookup")
	}
	if len(registry.Codes()) < 8 {
		t.Fatal("codes")
	}
}
