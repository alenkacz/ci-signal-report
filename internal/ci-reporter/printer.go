package ci_reporter

import (
	"fmt"
	"strings"
)

// This function is used to print github cards to the console
func PrintGitHubCards(flags metaFlags, cardsAfterType []GithubIssueCardSummary) {
	for _, e := range cardsAfterType {
		headerLine := fmt.Sprintf("\n\n%s %s", e.Emoji, strings.ToUpper(e.CardsTitle))
		if flags.EmojisOff {
			headerLine = fmt.Sprintf("\n\n%s", strings.ToUpper(e.CardsTitle))
		}
		fmt.Println(headerLine)
		for k, v := range e.ListGithubIssueOverview {
			fmt.Printf("SIG %s\n", k)
			for _, i := range v {
				fmt.Printf("- #%d %s %s\n", i.Id, i.Url, i.Title)
			}
			fmt.Println()
		}
	}
}

// This function is used to print testgrid data to the console
func PrintTestGridOverview(flags metaFlags, testgridStats []TestGridStatistics) {
	for _, stat := range testgridStats {
		headerLine := fmt.Sprintf("%s Tests in %s\n", stat.Emoji, stat.Name)
		if flags.EmojisOff {
			headerLine = fmt.Sprintf("Tests in %s\n", stat.Name)
		}
		fmt.Println(headerLine)
		fmt.Printf("\t%d jobs total\n", stat.Total)
		fmt.Printf("\t%d are passing\n", stat.Passing)
		fmt.Printf("\t%d are flaking\n", stat.Flaking)
		fmt.Printf("\t%d are failing\n", stat.Failing)
		fmt.Printf("\t%d are stale\n", stat.Stale)
		fmt.Print("\n\n")
	}
}
