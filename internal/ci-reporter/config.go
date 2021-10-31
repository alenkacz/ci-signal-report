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
	// shortens the report output (less details)
	ShortOn bool
	// prints emojis
	EmojisOff bool
	// specifies a specific release version that should be included in the report like "1.22"
	ReleaseVersion string
}

// Meta meta struct to use ci-reporter functions
type Meta struct {
	Env          metaEnv
	Flags        metaFlags
	GitHubClient *github.Client
}

// SetMeta this function is used to set meta information that is being needed to generate ci-signal-report
func SetMeta() Meta {
	// Flags
	// -short default: off
	isFlagShortSet := flag.Bool("short", false, "a short report for mails and slack")

	// -emoji-off - default : on
	isFlagEmojiOff := flag.Bool("emoji-off", false, "toggel if emojis should not be printed out")

	// TODO: check input with regex
	// -v default: ""
	releaseVersion := flag.String("v", "", "specify a release version to get additional report data")
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
			ReleaseVersion: *releaseVersion,
		},
		GitHubClient: ghClient,
	}
}
