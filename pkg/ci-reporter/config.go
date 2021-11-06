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
	"context"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/google/go-github/v34/github"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
)

// Environment variables that can be set using the ci-reporter
type metaEnv struct {
	GithubToken string `envconfig:"GITHUB_AUTH_TOKEN" required:"true"`
}

// Flags that can be set using the ci-reporter
type metaFlags struct {
	// ShortOn tells if the report should be shortended (with less details)
	ShortOn bool
	// EmojisOff tells if emojis should be printed
	EmojisOff bool
	// specifies a specific release version that should be included in the report like "1.22" or "1.22, 1.21"
	ReleaseVersion []string
	// JSONOut specifies if the output should be in json format
	JSONOut bool
	// Specify a report (if this is specified only one report will be printed e.g. SpecificReport: 'github' -> github report)
	SpecificReport string
}

// Meta meta struct to use ci-reporter functions
type Meta struct {
	Env                metaEnv
	Flags              metaFlags
	GitHubClient       *github.Client
	DataPostProcessing func(CIReport, string, chan ReportDataField, *sync.WaitGroup) ReportData
}

func dataPostProcessing(r CIReport, reportName string, chanReportDataField chan ReportDataField, wg *sync.WaitGroup) ReportData {
	reportData := ReportData{
		Data: []ReportDataField{},
		Name: reportName,
	}
	for reportDataField := range chanReportDataField {
		reportData.Data = append(reportData.Data, reportDataField)
	}
	r.PutData(reportData)
	wg.Done()
	return reportData
}

// SetMeta this function is used to set meta information that is being needed to generate ci-signal-report
func SetMeta() Meta {
	// Flags
	// -short default: off
	isFlagShortSet := flag.Bool("short", false, "Shortens the report")

	// -emoji-off - default : on
	isFlagEmojiOff := flag.Bool("emoji-off", false, "Remove emojis from report print-out")

	// -v default: ""
	releaseVersion := flag.String("v", "", "Adds specific K8s release version to the report (like -v '1.22, 1.21' or -v 1.22)")

	// -emoji-off - default : off
	isJSONOut := flag.Bool("json", false, "Report gets printed out in json format")

	// -emoji-off - default : off
	specificReport := flag.String("report", "", fmt.Sprintf("Specify report, options: '%s', '%s'", githubReport, testgridReport))

	flag.Parse()

	var env metaEnv
	err := envconfig.Process("", &env)
	if err != nil {
		// "Make sure to provide a GITHUB_AUTH_TOKEN, received an error during env decoding"
		log.Fatalf("Error processing flags.\n[ERROR] %v", err)
	}

	// Setup github client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: env.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)

	// Set meta data
	return Meta{
		Env: env,
		Flags: metaFlags{
			ShortOn:        *isFlagShortSet,
			EmojisOff:      *isFlagEmojiOff,
			ReleaseVersion: splitReleaseVersionInput(*releaseVersion),
			JSONOut:        *isJSONOut,
			SpecificReport: *specificReport,
		},
		GitHubClient:       ghClient,
		DataPostProcessing: dataPostProcessing,
	}
}

// GetReporters used to get reporters that implement methods like RequestData and Print
func (m Meta) GetReporters() []CIReport {
	if m.Flags.SpecificReport == "" {
		return []CIReport{&GithubReport{}, &TestgridReport{}}
	} else if m.Flags.SpecificReport == githubReport {
		return []CIReport{&GithubReport{}}
	} else if m.Flags.SpecificReport == testgridReport {
		return []CIReport{&TestgridReport{}}
	} else {
		log.Fatalf("Information given via flag -report does not match options [%s, %s]", githubReport, testgridReport)
	}
	return nil
}

// This function is used to split release version input ("1.22, 1.21" => ["1.22", "1.21"])
func splitReleaseVersionInput(input string) []string {
	re := regexp.MustCompile(`\d.\d\d`)
	releaseVersion := []string{}

	for _, e := range strings.Split(input, ",") {
		if e != "" {
			trimStr := strings.TrimSpace(e)
			found := re.MatchString(trimStr)

			if found {
				releaseVersion = append(releaseVersion, trimStr)
			} else {
				fmt.Printf("%s does not match\n", trimStr)
			}
		}
	}
	return releaseVersion
}
