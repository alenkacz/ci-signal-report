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
	GithubToken    string `envconfig:"GITHUB_AUTH_TOKEN" required:"true"`
	ReleaseVersion string `envconfig:"RELEASE_VERSION"`
}

// Flags that can be set using the ci-reporter
type metaFlags struct {
	// shortens the report output (less details)
	ShortOn bool
	// prints emojis
	EmojisOff bool
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
	// default: off
	isFlagShortSet := flag.Bool("short", false, "a short report for mails and slack")

	// default : on
	isFlagEmojiOff := flag.Bool("emoji-off", false, "toggel if emojis should not be printed out")
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
			ShortOn:   *isFlagShortSet,
			EmojisOff: *isFlagEmojiOff,
		},
		GitHubClient: ghClient,
	}
}
