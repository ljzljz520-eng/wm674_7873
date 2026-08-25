package service

import (
	"strings"
	"testing"
)

func TestParseCSV(t *testing.T) {
	rows, err := ParseCSV(strings.NewReader("r1,Maker,M,S1,10\nr2,Maker2,M2,S2,20\n"))
	if err != nil || len(rows) != 2 || rows[1].Amount != 20 {
		t.Fatalf("rows=%v err=%v", rows, err)
	}
}
func TestFormatCSV(t *testing.T) {
	text := FormatCSV(nil)
	if text != "id,manufacturer,model,reference,amount,status\n" {
		t.Fatal(text)
	}
}
