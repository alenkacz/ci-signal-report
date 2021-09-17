package transformer

import (
	"fmt"

	"github.com/leonardpahlke/ci-signal-report/internal/models"
)

// This function is used to print github cards to the console
func PrintGitHubCards(shortReport bool, cardsAfterType []models.GithubIssueCardSummary) {
	for _, e := range cardsAfterType {
		fmt.Println("\n----------\n" + e.CardsTitle)
		for k, v := range e.ListGithubIssueOverview {
			fmt.Printf("SIG %s\n", k)
			for _, i := range v {
				fmt.Printf("#%d %s %s\n", i.Id, i.Url, i.Title)
			}
			fmt.Println()
		}
	}
}

// This function is used to print testgrid data to the console
func PrintTestGridOverview(testgridStats []models.TestGridStatistics) {
	for _, stat := range testgridStats {
		fmt.Printf("Failures in %s\n", stat.Name)
		fmt.Printf("\t%d jobs total\n", stat.Total)
		fmt.Printf("\t%d are passing\n", stat.Passing)
		fmt.Printf("\t%d are flaking\n", stat.Flaking)
		fmt.Printf("\t%d are failing\n", stat.Failing)
		fmt.Printf("\t%d are stale\n", stat.Stale)
		fmt.Print("\n\n")
	}
}
