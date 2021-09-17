package report

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/alenkacz/ci-signal-report/internal/config"
	"github.com/alenkacz/ci-signal-report/internal/models"
	"github.com/google/go-github/v34/github"
)

func PrintCardsOverview(meta config.Meta) error {
	var listGithubIssueOverview []map[string][]*models.GhIssueOverview
	for _, e := range config.GithubIssueCardConfigs {
		if !(e.OmitWithFlagShort && meta.Flags.Short) {
			cardsOverview, err := reqGhCardsFromColumn(config.GhNewCards, meta.GitHubClient, meta.Env.GithubToken)
			if err != nil {
				return err
			}
			listGithubIssueOverview = append(listGithubIssueOverview, groupByCards(cardsOverview))
		}
	}

	newCardsOverview, err := reqGhCardsFromColumn(config.GhNewCards, meta.GitHubClient, meta.Env.GithubToken)
	if err != nil {
		return err
	}
	investigationCardsOverview, err := reqGhCardsFromColumn(config.GhUnderInvestigationCards, meta.GitHubClient, meta.Env.GithubToken)
	if err != nil {
		return err
	}
	observingCardsOverview, err := reqGhCardsFromColumn(config.GhObservingCards, meta.GitHubClient, meta.Env.GithubToken)
	if err != nil {
		return err
	}
	resolvedCards, err := findResolvedCardsColumnId(meta.GitHubClient)
	if err != nil {
		return err
	}
	resolvedCardsOverview, err := reqGhCardsFromColumn(resolvedCards, meta.GitHubClient, meta.Env.GithubToken)
	if err != nil {
		return err
	}

	PrintGhCards(meta.Flags.Short, groupByCards(newCardsOverview), groupByCards(investigationCardsOverview), groupByCards(observingCardsOverview), groupByCards(resolvedCardsOverview))

	return nil
}

func RequestGithubData() {

}

func findResolvedCardsColumnId(client *github.Client) (int64, error) {
	cards, _, err := client.Projects.ListProjectColumns(context.Background(), config.GhCiSignalBoardProjectId, &github.ListOptions{})
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

func groupByCards(issues []*models.GhIssueOverview) map[string][]*models.GhIssueOverview {
	result := make(map[string][]*models.GhIssueOverview)
	for _, issue := range issues {
		_, ok := result[issue.Sig]
		if !ok {
			result[issue.Sig] = make([]*models.GhIssueOverview, 0)
		}
		result[issue.Sig] = append(result[issue.Sig], issue)
	}
	return result
}

func reqGhCardsFromColumn(cardsId int64, client *github.Client, token string) ([]*models.GhIssueOverview, error) {
	cards, _, err := client.Projects.ListProjectCards(context.Background(), cardsId, &github.ProjectCardListOptions{})
	if err != nil {
		fmt.Printf("error when querying cards %v", err)
		return nil, err
	}
	issues := make([]*models.GhIssueOverview, 0)
	for _, c := range cards {
		if c.ContentURL != nil {
			issueUrl := *c.ContentURL
			issueDetail, err := requestGhIssueDetail(issueUrl, token)
			if err != nil {
				return nil, err
			}

			overview := models.GhIssueOverview{
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

func requestGhIssueDetail(url string, authToken string) (*models.GhIssueDetail, error) {
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
	var result models.GhIssueDetail
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
