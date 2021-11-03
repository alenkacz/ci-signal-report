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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/google/go-github/v34/github"
)

// GithubReport used to implement RequestData & Print for github report data
type GithubReport struct {
	ReportData ReportData
}

// RequestData this function is used to get github report data
func (r *GithubReport) RequestData(meta Meta, wg *sync.WaitGroup) ReportData {
	resolvedCardsID, err := findCardsID(meta.GitHubClient, githubCiSignalBoardProjectID, "Resolved")
	if err != nil {
		log.Fatalf("Could not find github card %v", err)
		return ReportData{}
	}
	var githubIssueCardConfigs = []githubIssueCardConfig{
		{
			CardsTitle:        "New/Not Yet Started",
			CardID:            newCardsID,
			OmitWithFlagShort: false,
			Emoji:             notYetStartedEmoji,
		},
		{
			CardsTitle:        "In flight",
			CardID:            underInvestigationCardsID,
			OmitWithFlagShort: false,
			Emoji:             inFlightEmoji,
		},
		{
			CardsTitle:        "Observing",
			CardID:            observingCardsID,
			OmitWithFlagShort: true,
			Emoji:             observingEmoji,
		},
		{
			CardsTitle:        "Resolved",
			CardID:            int(resolvedCardsID),
			OmitWithFlagShort: true,
			Emoji:             resolvedEmoji,
		},
	}

	// DataPostProcessing collects data requested via assembleGithubRequests/2 and returns ReportData
	return meta.DataPostProcessing(r, "github", assembleGithubRequests(meta, githubIssueCardConfigs), wg)
}

// Print extends GithubReport and prints report data to the console
func (r GithubReport) Print(meta Meta, reportData ReportData) {
	// Print regular out
	for _, data := range reportData.Data {
		// Prepare header
		headerLine := fmt.Sprintf("\n\n%s %s", data.Emoji, strings.ToUpper(data.Title))
		if meta.Flags.EmojisOff {
			headerLine = fmt.Sprintf("\n\n%s", strings.ToUpper(data.Title))
		}
		fmt.Println(headerLine)

		// sort report data after SIG
		dataSortedAfterSigs := make(map[string][]ghIssueOverview)
		for _, v := range data.Records {
			dataSortedAfterSigs[v.Sig] = append(dataSortedAfterSigs[v.Sig], ghIssueOverview{
				URL:   v.URL,
				ID:    v.ID,
				Title: v.Title,
			})
		}

		// print sorted report data
		for k, v := range dataSortedAfterSigs {
			fmt.Printf("SIG %s\n", k)
			for _, i := range v {
				fmt.Printf("- #%d %s %s\n", i.ID, i.URL, i.Title)
			}
			fmt.Println()
		}
	}
}

// PutData extends GithubReport and stores the data at runtime to the struct val ReportData
func (r *GithubReport) PutData(reportData ReportData) {
	r.ReportData = reportData
}

// GetData extends GithubReport and returns the data that has been stored at runtime int the struct val ReportData (counter to SaveData/1)
func (r GithubReport) GetData() ReportData {
	return r.ReportData
}

// run all github requests to assemble data
func assembleGithubRequests(meta Meta, githubIssueCardConfigs []githubIssueCardConfig) chan ReportDataField {
	c := make(chan ReportDataField)
	go func() {
		defer close(c)
		var wg sync.WaitGroup
		for _, cardCfg := range githubIssueCardConfigs {
			wg.Add(1)
			go sortDataIntoDataRecord(meta, c, &wg, cardCfg)
		}
		wg.Wait()
	}()
	return c
}

func sortDataIntoDataRecord(meta Meta, c chan ReportDataField, wg *sync.WaitGroup, cardCfg githubIssueCardConfig) {
	if !(cardCfg.OmitWithFlagShort && meta.Flags.ShortOn) {
		reportDataRecord := []ReportDataRecord{}
		// request github data
		for issue := range assembleCardRequests(meta, int64(cardCfg.CardID)) {
			// transform data structure
			reportDataRecord = append(reportDataRecord, ReportDataRecord{
				URL:   issue.URL,
				ID:    issue.ID,
				Title: issue.Title,
				Sig:   issue.Sig,
			})
		}
		// send data through channel; data infos gathered
		c <- ReportDataField{
			Emoji:   cardCfg.Emoji,
			Title:   cardCfg.CardsTitle,
			Records: reportDataRecord,
		}
	}
	wg.Done()
}

// run a github card requests to assemble cards data
func assembleCardRequests(meta Meta, cardsID int64) chan ghIssueOverview {
	// int64(e.CardID), meta.GitHubClient, meta.Env.GithubToken
	c := make(chan ghIssueOverview)

	cards, _, err := meta.GitHubClient.Projects.ListProjectCards(context.Background(), cardsID, &github.ProjectCardListOptions{})
	if err != nil {
		log.Printf("error when querying cards %v", err)
	}

	go func() {
		defer close(c)
		var wg sync.WaitGroup
		for _, card := range cards {
			wg.Add(1)
			go func(card *github.ProjectCard, token string) {
				if card.ContentURL != nil {
					issueDetail, err := requestGhIssueDetail(*card.ContentURL, token)
					if err != nil {
						log.Printf("Error on requesting github card information.\n[ERROR] %v", err)
					}

					overview := ghIssueOverview{
						URL:   issueDetail.HTMLURL,
						ID:    issueDetail.Number,
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
					c <- overview
				}
				wg.Done()
			}(card, meta.Env.GithubToken)
		}
		wg.Wait()
	}()
	return c
}

func findCardsID(client *github.Client, projectID int64, keyword string) (int64, error) {
	cards, _, err := client.Projects.ListProjectColumns(context.Background(), projectID, &github.ListOptions{})
	if err != nil {
		return 0, err
	}
	resolvedColumns := make([]*github.ProjectColumn, 0)
	for _, v := range cards {
		if v.Name != nil && *v.Name == keyword {
			resolvedColumns = append(resolvedColumns, v)
		}
	}
	sort.Slice(resolvedColumns, func(i, j int) bool {
		return resolvedColumns[i].GetID() < resolvedColumns[j].GetID()
	})
	return resolvedColumns[0].GetID(), err
}

func requestGhIssueDetail(url string, authToken string) (ghIssueDetail, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ghIssueDetail{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))

	// Send http request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%v", err)
		return ghIssueDetail{}, err
	}
	var result ghIssueDetail
	err = json.Unmarshal(body, &result)
	if err != nil {
		return ghIssueDetail{}, err
	}
	return result, nil
}

// Internal types

// ghIssueOverview information about a specific github issue
type ghIssueOverview struct {
	URL   string
	ID    int64
	Title string
	Sig   string
}

type ghIssueDetail struct {
	Number  int64          `json:"number"`
	HTMLURL string         `json:"html_url"`
	Title   string         `json:"title"`
	Labels  []github.Label `json:"labels,omitempty"`
}

type githubIssueCardConfig struct {
	CardsTitle        string
	CardID            int
	Emoji             string
	OmitWithFlagShort bool
}
