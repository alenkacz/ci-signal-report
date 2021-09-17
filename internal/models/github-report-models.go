package models

import "github.com/google/go-github/v34/github"

type GhIssueOverview struct {
	Url   string
	Id    int64
	Title string
	Sig   string
}

type GhIssueDetail struct {
	Number  int64          `json:"number"`
	HtmlUrl string         `json:"html_url"`
	Title   string         `json:"title"`
	Labels  []github.Label `json:"labels,omitempty"`
}

type GithubIssueCardConfig struct {
	CardsTitle        string
	CardId            int
	OmitWithFlagShort bool
}

type GithubIssueCardSummary struct {
	CardsTitle              string
	ListGithubIssueOverview map[string][]*GhIssueOverview
}
