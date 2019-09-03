package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/google/go-github/v27/github"
	"golang.org/x/oauth2"
)

type requiredJob struct {
	OutputName string
	UrlName    string
}

func main() {
	githubApiToken := os.Getenv("GITHUB_AUTH_TOKEN")
	if githubApiToken == "" {
		fmt.Printf("Please provide GITHUB_AUTH_TOKEN env variable to be able to pull cards from the github board")
		os.Exit(1)
	}

	err := printCardsOverview(githubApiToken)
	if err != nil {
		fmt.Printf("error when querying cards overview, exiting: %v\n", err)
		os.Exit(1)
	}
	printJobsStatistics()
}

const (
	newCards                = 4212817
	underInvestigationCards = 4212819
	observingCards          = 4212821
	ciSignalBoardProjectId  = 2093513
)

type issueOverview struct {
	url   string
	id    int64
	title string
	sig   string
}

func printCardsOverview(token string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	newCardsOverview, err := getCardsFromColumn(newCards, client)
	if err != nil {
		return err
	}
	investigationCardsOverview, err := getCardsFromColumn(underInvestigationCards, client)
	if err != nil {
		return err
	}
	observingCardsOverview, err := getCardsFromColumn(observingCards, client)
	if err != nil {
		return err
	}
	resolvedCards, err := findResolvedCardsColumns(client)
	if err != nil {
		return err
	}
	resolvedCardsOverview, err := getCardsFromColumn(resolvedCards, client)
	if err != nil {
		return err
	}

	printCards(groupByCards(newCardsOverview), groupByCards(investigationCardsOverview), groupByCards(observingCardsOverview), groupByCards(resolvedCardsOverview))
	return nil
}

func findResolvedCardsColumns(client *github.Client) (int64, error) {
	opt := &github.ListOptions{}
	columns, _, err := client.Projects.ListProjectColumns(context.Background(), ciSignalBoardProjectId, opt)
	if err != nil {
		return 0, err
	}
	resolvedColumns := make([]*github.ProjectColumn, 0)
	for _, v := range columns {
		if v.Name != nil && strings.HasPrefix(*v.Name, "Resolved") {
			resolvedColumns = append(resolvedColumns, v)
		}
	}
	sort.Slice(resolvedColumns, func(i, j int) bool {
		return resolvedColumns[i].GetID() < resolvedColumns[j].GetID()
	})
	return resolvedColumns[0].GetID(), err
}

func printCards(new map[string][]*issueOverview, investigation map[string][]*issueOverview, observing map[string][]*issueOverview, resolved map[string][]*issueOverview) {
	fmt.Println("Resolved")
	for k, v := range resolved {
		fmt.Printf("SIG %s\n", k)
		for _, i := range v {
			fmt.Printf("#%d %s %s\n", i.id, i.url, i.title)
		}
		fmt.Println()
	}

	fmt.Println("In flight")
	for k, v := range observing {
		fmt.Printf("SIG %s\n", k)
		for _, i := range v {
			fmt.Printf("#%d %s %s\n", i.id, i.url, i.title)
		}
		fmt.Println()
	}
	for k, v := range investigation {
		fmt.Printf("SIG %s\n", k)
		for _, i := range v {
			fmt.Printf("#%d %s %s\n", i.id, i.url, i.title)
		}
		fmt.Println()
	}

	fmt.Println("New/Not Yet Started")
	for k, v := range new {
		fmt.Printf("SIG %s\n", k)
		for _, i := range v {
			fmt.Printf("#%d %s %s\n", i.id, i.url, i.title)
		}
		fmt.Println()
	}
}

func groupByCards(issues []*issueOverview) map[string][]*issueOverview {
	result := make(map[string][]*issueOverview)
	for _, i := range issues {
		_, ok := result[i.sig]
		if !ok {
			result[i.sig] = make([]*issueOverview, 0)
		}
		result[i.sig] = append(result[i.sig], i)
	}
	return result
}

func getCardsFromColumn(cardsId int64, client *github.Client) ([]*issueOverview, error) {
	opt := &github.ProjectCardListOptions{}
	cards, _, err := client.Projects.ListProjectCards(context.Background(), cardsId, opt)
	if err != nil {
		fmt.Printf("error when querying cards %v", err)
		return nil, err
	}
	issues := make([]*issueOverview, 0)
	for _, c := range cards {
		issueUrl := *c.ContentURL
		issueDetail, err := getIssueDetail(issueUrl)
		if err != nil {
			return nil, err
		}
		overview := issueOverview{
			url:   issueDetail.HtmlUrl,
			id:    issueDetail.Number,
			title: cleanTitle(issueDetail.Title),
		}
		for _, v := range issueDetail.Labels {
			if strings.Contains(*v.Name, "sig/") {
				overview.sig = strings.Title(strings.Replace(*v.Name, "sig/", "", -1))
				if strings.EqualFold(overview.sig, "cli") {
					overview.sig = strings.ToUpper(overview.sig)
				}
				if strings.EqualFold(overview.sig, "cluster-lifecycle") {
					overview.sig = strings.ToLower(overview.sig)
				}
				break
			}
		}
		issues = append(issues, &overview)
	}
	return issues, nil
}

type IssueDetail struct {
	Number  int64          `json:"number"`
	HtmlUrl string         `json:"html_url"`
	Title   string         `json:"title"`
	Labels  []github.Label `json:"labels,omitempty"`
}

func getIssueDetail(url string) (*IssueDetail, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}
	var result IssueDetail
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func cleanTitle(title string) string {
	title = strings.Replace(title, "[Failing Test]", "", -1)
	return title
}

func printJobsStatistics() {
	requiredJobs := []requiredJob{
		{OutputName: "Master-Blocking", UrlName: "sig-release-master-blocking"},
		{OutputName: "Master-Informing", UrlName: "sig-release-master-informing"},
		{OutputName: "1.16-blocking", UrlName: "sig-release-1.16-blocking"},
		{OutputName: "1.16-informing", UrlName: "sig-release-1.16-blocking"},
	}

	result := make([]statistics, 0)
	for _, kubeJob := range requiredJobs {
		url := fmt.Sprintf("https://testgrid.k8s.io/%s/summary", kubeJob.UrlName)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		jobs, err := getStatsFromJson(body)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		statistics := getStatistics(jobs)
		statistics.name = kubeJob.OutputName
		result = append(result, statistics)
	}

	prettyPrint(result)
}

func prettyPrint(stats []statistics) {
	for _, stat := range stats {
		fmt.Printf("Failures in %s\n", stat.name)
		fmt.Printf("\t%d jobs total\n", stat.total)
		fmt.Printf("\t%d are passing\n", stat.passing)
		fmt.Printf("\t%d are flaking\n", stat.flaking)
		fmt.Printf("\t%d are failing\n", stat.failing)
		fmt.Printf("\t%d are stale\n", stat.stale)
		fmt.Print("\n\n")
	}
}

func getStatistics(jobs map[string]overview) statistics {
	result := statistics{}
	for _, v := range jobs {
		if v.OverallStatus == "PASSING" {
			result.passing++
		} else if v.OverallStatus == "FAILING" {
			result.failing++
		} else if v.OverallStatus == "FLAKY" {
			result.flaking++
		} else {
			result.stale++
		}
		result.total++
	}
	return result
}

type statistics struct {
	name    string
	total   int
	passing int
	flaking int
	failing int
	stale   int
}
type jobsOverview = map[string]overview
type overview struct {
	OverallStatus string `json:"overall_status"`
}

func getStatsFromJson(body []byte) (jobsOverview, error) {
	result := make(jobsOverview)
	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
