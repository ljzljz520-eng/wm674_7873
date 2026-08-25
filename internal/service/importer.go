package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"meter-sync/internal/domain"
)

func ParseCSV(input io.Reader) ([]domain.ImportRow, error) {
	reader := csv.NewReader(input)
	reader.FieldsPerRecord = 5
	rows := make([]domain.ImportRow, 0)
	for line := 1; ; line++ {
		values, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
		amount, err := strconv.ParseInt(strings.TrimSpace(values[4]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d amount: %w", line, err)
		}
		rows = append(rows, domain.ImportRow{ID: strings.TrimSpace(values[0]), Manufacturer: strings.TrimSpace(values[1]), Model: strings.TrimSpace(values[2]), Reference: strings.TrimSpace(values[3]), Amount: amount})
	}
	return rows, nil
}

func FormatCSV(rows []domain.Record) string {
	var builder strings.Builder
	builder.WriteString("id,manufacturer,model,reference,amount,status\n")
	for _, r := range rows {
		fmt.Fprintf(&builder, "%s,%s,%s,%s,%d,%s\n", r.ID, r.Manufacturer, r.MeterModel, r.SyncReference, r.Amount, r.Status)
	}
	return builder.String()
}
