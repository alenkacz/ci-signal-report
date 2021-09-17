package report

import (
	"fmt"

	"github.com/alenkacz/ci-signal-report/internal/models"
)

func PrintGhCards(shortReport bool, new map[string][]*models.GhIssueOverview, investigation map[string][]*models.GhIssueOverview, observing map[string][]*models.GhIssueOverview, resolved map[string][]*models.GhIssueOverview) {
	fmt.Println("New/Not Yet Started")
	for k, v := range new {
		fmt.Printf("SIG %s\n", k)
		for _, i := range v {
			fmt.Printf("#%d %s %s\n", i.Id, i.Url, i.Title)
		}
		fmt.Println()
	}

	fmt.Println("In flight")
	for k, v := range investigation {
		fmt.Printf("SIG %s\n", k)
		for _, i := range v {
			fmt.Printf("#%d %s %s\n", i.Id, i.Url, i.Title)
		}
		fmt.Println()
	}

	if shortReport == false {
		fmt.Println("Observing")
		for k, v := range observing {
			fmt.Printf("SIG %s\n", k)
			for _, i := range v {
				fmt.Printf("#%d %s %s\n", i.Id, i.Url, i.Title)
			}
			fmt.Println()
		}

		fmt.Println("Resolved")
		for k, v := range resolved {
			fmt.Printf("SIG %s\n", k)
			for _, i := range v {
				fmt.Printf("#%d %s %s\n", i.Id, i.Url, i.Title)
			}
			fmt.Println()
		}
	}

}
