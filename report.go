package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
		fmt.Printf("error when querying cards overview, exiting")
		os.Exit(1)
	}
	printJobsStatistics()
}

const (
	newCards                = 4212817
	underInvestigationCards = 4212819
	observingCards          = 4212821
	resolvedCards           = 5891313
)

type issueOverview struct {
	url   string
	id    int
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
	resolvedCardsOverview, err := getCardsFromColumn(resolvedCards, client)
	if err != nil {
		return err
	}

	printCards(newCardsOverview, investigationCardsOverview, observingCardsOverview, resolvedCardsOverview)
	return nil
}

func printCards(new []*issueOverview, investigation []*issueOverview, observing []*issueOverview, resolved []*issueOverview) {
	fmt.Println("Resolved")
	for _, v := range resolved {
		fmt.Printf("SIG %s\n", v.sig)
		fmt.Printf("#%d %s %s\n", v.id, v.url, v.title)
		fmt.Println()
	}

	fmt.Println("In flight")
	for _, v := range observing {
		fmt.Printf("SIG %s\n", v.sig)
		fmt.Printf("#%d %s %s\n", v.id, v.url, v.title)
		fmt.Println()
	}
	for _, v := range investigation {
		fmt.Printf("SIG %s\n", v.sig)
		fmt.Printf("#%d %s %s\n", v.id, v.url, v.title)
		fmt.Println()
	}

	fmt.Println("New/Not Yet Started")
	for _, v := range new {
		fmt.Printf("SIG %s\n", v.sig)
		fmt.Printf("#%d %s %s\n", v.id, v.url, v.title)
		fmt.Println()
	}
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
		issueUrlParts := strings.Split(issueUrl, "/")
		issueId, err := strconv.Atoi(issueUrlParts[len(issueUrlParts)-1])
		if err != nil {
			fmt.Printf("Cannot parse issue ID from url %s", issueUrl)
			return nil, err
		}
		issue, _, err := client.Issues.Get(context.Background(), "kubernetes", "kubernetes", issueId)
		if err != nil {
			fmt.Printf("Cannot query issue id %d", issueId)
			return nil, err
		}
		overview := issueOverview{
			url:   issueUrl,
			id:    issueId,
			title: cleanTitle(*issue.Title),
		}
		for _, v := range issue.Labels {
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