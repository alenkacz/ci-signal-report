package main

import (
	"fmt"
	"os"

	"github.com/alenkacz/ci-signal-report/internal/config"
	"github.com/alenkacz/ci-signal-report/internal/report"
)

func main() {
	meta := config.SetMeta()

	err := report.PrintCardsOverview(meta)
	if err != nil {
		fmt.Printf("error when querying cards overview, exiting: %v\n", err)
		os.Exit(1)
	}

	report.PrintTestgridJobsStatistics(meta)
}
