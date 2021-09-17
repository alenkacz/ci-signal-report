package config

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-github/v34/github"
	"golang.org/x/oauth2"
)

type metaEnv struct {
	GithubToken    string
	ReleaseVersion string
}

type metaFlags struct {
	Short bool
}

type Meta struct {
	Env          metaEnv
	Flags        metaFlags
	GitHubClient *github.Client
}

func SetMeta() Meta {
	// Flags
	isFlagShortSet := flag.Bool("short", false, "a short report for mails and slack")
	flag.Parse()

	// Environment vairables
	githubApiToken := os.Getenv("GITHUB_AUTH_TOKEN")
	releaseVersion := os.Getenv("RELEASE_VERSION")
	if githubApiToken == "" {
		fmt.Printf("Please provide GITHUB_AUTH_TOKEN env variable to be able to pull cards from the github board")
		os.Exit(1)
	}

	// Setup github client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubApiToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)

	// Set meta data
	return Meta{
		Env: metaEnv{
			GithubToken:    githubApiToken,
			ReleaseVersion: releaseVersion,
		},
		Flags: metaFlags{
			Short: *isFlagShortSet,
		},
		GitHubClient: ghClient,
	}
}
