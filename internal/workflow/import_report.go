package workflow

import (
	"io"
	"meter-sync/internal/domain"
	"meter-sync/internal/service"
)

type ImportReport struct{ Service *service.Service }

func (w ImportReport) ImportCSV(input io.Reader, actor string) (domain.ImportResult, error) {
	rows, err := service.ParseCSV(input)
	if err != nil {
		return domain.ImportResult{}, err
	}
	return w.Service.Import(rows, actor), nil
}
func (w ImportReport) Report(query domain.Query) (service.Summary, error) {
	records, err := w.Service.Search(query)
	if err != nil {
		return service.Summary{}, err
	}
	return service.Summarize(records), nil
}
