package ci_reporter

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/google/go-github/v34/github"
)

func RequestGitHubCardsData(meta CiReporterMeta) ([]GithubIssueCardSummary, error) {
	resolvedCardsId, err := findResolvedCardsWithProjectId(meta.GitHubClient, GithubCiSignalBoardProjectId)
	if err != nil {
		return nil, err
	}
	var githubIssueCardConfigs = []GithubIssueCardConfig{
		{
			CardsTitle:        "New/Not Yet Started",
			CardId:            GithubNewCardsId,
			OmitWithFlagShort: false,
			Emoji:             NotYetStartedEmoji,
		},
		{
			CardsTitle:        "In flight",
			CardId:            GithubUnderInvestigationCardsId,
			OmitWithFlagShort: false,
			Emoji:             InFlightEmoji,
		},
		{
			CardsTitle:        "Observing",
			CardId:            GithubObservingCardsId,
			OmitWithFlagShort: true,
			Emoji:             ObservingEmoji,
		},
		{
			CardsTitle:        "Resolved",
			CardId:            int(resolvedCardsId),
			OmitWithFlagShort: true,
			Emoji:             ResolvedEmoji,
		},
	}

	var listGithubIssueOverview []GithubIssueCardSummary
	for _, e := range githubIssueCardConfigs {
		if !(e.OmitWithFlagShort && meta.Flags.ShortOn) {
			cardsOverview, err := reqGhCardsFromColumn(int64(e.CardId), meta.GitHubClient, meta.Env.GithubToken)
			if err != nil {
				return nil, err
			}
			listGithubIssueOverview = append(listGithubIssueOverview, GithubIssueCardSummary{
				CardsTitle:              e.CardsTitle,
				Emoji:                   e.Emoji,
				ListGithubIssueOverview: groupByCards(cardsOverview),
			})
		}
	}
	return listGithubIssueOverview, nil
}

func findResolvedCardsWithProjectId(client *github.Client, projectId int64) (int64, error) {
	cards, _, err := client.Projects.ListProjectColumns(context.Background(), projectId, &github.ListOptions{})
	if err != nil {
		return 0, err
	}
	resolvedColumns := make([]*github.ProjectColumn, 0)
	for _, v := range cards {
		if v.Name != nil && *v.Name == "Resolved" {
			resolvedColumns = append(resolvedColumns, v)
		}
	}
	sort.Slice(resolvedColumns, func(i, j int) bool {
		return resolvedColumns[i].GetID() < resolvedColumns[j].GetID()
	})
	return resolvedColumns[0].GetID(), err
}

func reqGhCardsFromColumn(cardsId int64, client *github.Client, token string) ([]*GhIssueOverview, error) {
	cards, _, err := client.Projects.ListProjectCards(context.Background(), cardsId, &github.ProjectCardListOptions{})
	if err != nil {
		fmt.Printf("error when querying cards %v", err)
		return nil, err
	}
	issues := make([]*GhIssueOverview, 0)
	for _, c := range cards {
		if c.ContentURL != nil {
			issueUrl := *c.ContentURL
			issueDetail, err := requestGhIssueDetail(issueUrl, token)
			if err != nil {
				return nil, err
			}

			overview := GhIssueOverview{
				Url:   issueDetail.HtmlUrl,
				Id:    issueDetail.Number,
				Title: strings.Replace(issueDetail.Title, "[Failing Test]", "", -1),
			}
			for _, v := range issueDetail.Labels {
				if strings.Contains(*v.Name, "sig/") {
					overview.Sig = strings.Title(strings.Replace(*v.Name, "sig/", "", -1))
					if strings.EqualFold(overview.Sig, "cli") {
						overview.Sig = strings.ToUpper(overview.Sig)
					}
					if strings.EqualFold(overview.Sig, "cluster-lifecycle") {
						overview.Sig = strings.ToLower(overview.Sig)
					}
					break
				}
			}
			issues = append(issues, &overview)
		}
	}
	return issues, nil
}

func groupByCards(issues []*GhIssueOverview) map[string][]*GhIssueOverview {
	result := make(map[string][]*GhIssueOverview)
	for _, issue := range issues {
		_, ok := result[issue.Sig]
		if !ok {
			result[issue.Sig] = make([]*GhIssueOverview, 0)
		}
		result[issue.Sig] = append(result[issue.Sig], issue)
	}
	return result
}

func requestGhIssueDetail(url string, authToken string) (*GhIssueDetail, error) {
	// Create a Bearer string by appending string access token
	var bearerHeader = "Bearer " + authToken

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// add authorization header to the req
	req.Header.Add("Authorization", bearerHeader)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}
	var result GhIssueDetail
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

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
	Emoji             string
	OmitWithFlagShort bool
}

type GithubIssueCardSummary struct {
	CardsTitle              string
	Emoji                   string
	ListGithubIssueOverview map[string][]*GhIssueOverview
}
