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
	// JsonOut specifies if the output should be in json format
	JsonOut bool
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
	isFlagShortSet := flag.Bool("short", false, "a short report for mails and slack")

	// -emoji-off - default : on
	isFlagEmojiOff := flag.Bool("emoji-off", false, "toggel if emojis should not be printed out")

	// -v default: ""
	releaseVersion := flag.String("v", "", "specify a release version to get additional report data")

	// -emoji-off - default : off
	isJsonOut := flag.Bool("json", false, "toggel if output should be printed in json format")
	flag.Parse()

	var env metaEnv
	err := envconfig.Process("", &env)
	if err != nil {
		// "Make sure to provide a GITHUB_AUTH_TOKEN, received an error during env decoding"
		panic(err)
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
			JsonOut:        *isJsonOut,
		},
		GitHubClient:       ghClient,
		DataPostProcessing: dataPostProcessing,
	}
}

// This function is used to split release version input ("1.22, 1.21" => ["1.22", "1.21"])
func splitReleaseVersionInput(input string) []string {
	re, err := regexp.Compile(`\d.\d\d`)
	if err != nil {
		log.Fatal(err)
	}

	releaseVersion := []string{}

	for _, e := range strings.Split(input, ",") {
		if e != "" {
			trimStr := strings.TrimSpace(e)
			fmt.Println(trimStr)
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
