package main

import (
	"fmt"

	ci_reporter "github.com/leonardpahlke/ci-signal-report/internal/ci-reporter"
)

func main() {
	meta := ci_reporter.SetMeta()

	// GitHub Report
	listGithubIssueOverview, err := ci_reporter.RequestGitHubCardsData(meta)
	if err != nil {
		fmt.Printf("Error RequestGitHubCardsData %v", err)
	}
	ci_reporter.PrintGitHubCards(meta.Flags, listGithubIssueOverview)

	// TestGrid Report
	testgridOverview, err := ci_reporter.RequestTestgridOverview(meta)
	if err != nil {
		fmt.Printf("Error RequestTestgridOverview %v", err)
	}
	ci_reporter.PrintTestGridOverview(meta.Flags, testgridOverview)
}
