package main

import (
	"fmt"

	"github.com/leonardpahlke/ci-signal-report/internal/config"
	"github.com/leonardpahlke/ci-signal-report/internal/report"
	"github.com/leonardpahlke/ci-signal-report/internal/transformer"
)

func main() {
	meta := config.SetMeta()

	/* GitHub Report */
	// 1. request github data
	listGithubIssueOverview, err := report.RequestGitHubCardsData(meta)
	if err != nil {
		fmt.Printf("Error RequestGitHubCardsData %v", err)
	}
	// 2. print github data
	transformer.PrintGitHubCards(meta.Flags.Short, listGithubIssueOverview)

	/* TestGrid Report */
	// 1. request testgrid data
	testgridOverview, err := report.RequestTestgridOverview(meta)
	if err != nil {
		fmt.Printf("Error RequestTestgridOverview %v", err)
	}

	// 2. print testgrid data
	transformer.PrintTestGridOverview(testgridOverview)
}
