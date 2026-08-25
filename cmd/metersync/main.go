package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"meter-sync/internal/config"
	"meter-sync/internal/domain"
	"meter-sync/internal/service"
	"meter-sync/internal/store"
	"meter-sync/internal/workflow"
)

func main() {
	settings := config.Default()
	flag.StringVar(&settings.DatabasePath, "db", settings.DatabasePath, "database path")
	flag.StringVar(&settings.Actor, "actor", settings.Actor, "operator")
	flag.Parse()
	settings.Command = flag.Arg(0)
	settings = settings.Normalize()
	st, err := store.Open(settings.DatabasePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer st.Close()
	svc, err := service.New(st, service.FixedClock{Value: "2026-01-01T00:00:00Z"})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := run(settings.Command, settings.Actor, svc); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(command, actor string, svc *service.Service) error {
	switch command {
	case "demo":
		life := workflow.Lifecycle{Service: svc}
		record, err := life.CreateReviewPublishArchive("meter-001", "Huadian", "HX-100", "SYNC-001", actor, 12800)
		if err != nil {
			return err
		}
		return printJSON(record)
	case "register":
		record, err := svc.Register("meter-001", "Huadian", "HX-100", "SYNC-001", actor, 12800)
		if err != nil {
			return err
		}
		return printJSON(record)
	case "search":
		records, err := svc.Search(domain.Query{})
		if err != nil {
			return err
		}
		return printJSON(records)
	case "report":
		records, err := svc.Search(domain.Query{})
		if err != nil {
			return err
		}
		fmt.Println(service.Summarize(records))
		return nil
	default:
		return fmt.Errorf("unknown command %q", command)
	}
}

func printJSON(value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
