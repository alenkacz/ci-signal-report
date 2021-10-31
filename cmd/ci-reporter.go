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

package main

import (
	"fmt"

	ci_reporter "github.com/leonardpahlke/ci-signal-report/internal/ci-reporter"
)

func main() {
	meta := ci_reporter.SetMeta()

	// GitHub Report
	listGithubIssueOverview, err := ci_reporter.ReqGitHubData(meta)
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
