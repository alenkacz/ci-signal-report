package main

import (
	"fmt"

	ci_reporter "github.com/leonardpahlke/ci-signal-report/internal/ci-reporter"
)

func main() {
	meta := ci_reporter.SetMeta()

	/* GitHub Report */
	// 1. request github data
	listGithubIssueOverview, err := ci_reporter.RequestGitHubCardsData(meta)
	if err != nil {
		fmt.Printf("Error RequestGitHubCardsData %v", err)
	}
	// 2. print github data
	ci_reporter.PrintGitHubCards(meta.Flags.Short, listGithubIssueOverview)

	/* TestGrid Report */
	// 1. request testgrid data
	testgridOverview, err := ci_reporter.RequestTestgridOverview(meta)
	if err != nil {
		fmt.Printf("Error RequestTestgridOverview %v", err)
	}

	// 2. print testgrid data
	ci_reporter.PrintTestGridOverview(testgridOverview)
}
