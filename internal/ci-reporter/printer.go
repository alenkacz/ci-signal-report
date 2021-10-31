/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cireporter

import (
	"fmt"
	"strings"
)

// PrintGitHubCards this function is used to print github cards to the console
func PrintGitHubCards(flags metaFlags, cardsAfterType []GithubSummary) {
	for _, e := range cardsAfterType {
		headerLine := fmt.Sprintf("\n\n%s %s", e.Emoji, strings.ToUpper(e.CardsTitle))
		if flags.EmojisOff {
			headerLine = fmt.Sprintf("\n\n%s", strings.ToUpper(e.CardsTitle))
		}
		fmt.Println(headerLine)
		for k, v := range e.ListGithubIssueOverview {
			fmt.Printf("SIG %s\n", k)
			for _, i := range v {
				fmt.Printf("- #%d %s %s\n", i.ID, i.URL, i.Title)
			}
			fmt.Println()
		}
	}
}

// PrintTestGridOverview this function is used to print testgrid data to the console
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
