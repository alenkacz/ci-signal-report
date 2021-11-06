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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// TestgridReport used to implement RequestData & Print for testgrid report data
type TestgridReport struct {
	ReportData ReportData
}

// RequestData this function is used to accumulate a summary of testgrid
func (r *TestgridReport) RequestData(meta Meta, wg *sync.WaitGroup) ReportData {
	// The report checks master-blocking and master-informing
	requiredJobs := []testgridJob{
		{OutputName: "Master-Blocking", URLName: string(sigReleaseMasterBlocking), Emoji: masterBlockingEmoji},
		{OutputName: "Master-Informing", URLName: string(sigReleaseMasterInforming), Emoji: masterInformingEmoji},
	}

	// If a release version got specified add additional jobs to report
	if len(meta.Flags.ReleaseVersion) > 0 {
		for _, r := range meta.Flags.ReleaseVersion {
			requiredJobs = append(requiredJobs, testgridJob{OutputName: fmt.Sprintf("%s-blocking", r), URLName: fmt.Sprintf("sig-release-%s-blocking", r), Emoji: masterBlockingEmoji})
			requiredJobs = append(requiredJobs, testgridJob{OutputName: fmt.Sprintf("%s-informing", r), URLName: fmt.Sprintf("sig-release-%s-informing", r), Emoji: masterInformingEmoji})
		}
	}

	return meta.DataPostProcessing(r, testgridReport, assembleTestgridRequests(meta, requiredJobs), wg)
}

// Print extends TestgridReport and prints report data to the console
func (r *TestgridReport) Print(meta Meta, reportData ReportData) {
	for _, reportField := range reportData.Data {
		headerLine := fmt.Sprintf("\n\n%s Tests in %s", reportField.Emoji, reportField.Title)
		if meta.Flags.EmojisOff {
			headerLine = fmt.Sprintf("\n\nTests in %s", reportField.Title)
		}
		for _, stat := range reportField.Records {
			if stat.ID == testgridReportSummary {
				fmt.Println(headerLine)
				for _, note := range stat.Notes {
					fmt.Println("- " + note)
				}
				fmt.Print("\n")
				if !meta.Flags.ShortOn {
					fmt.Print("\nFAILING & FLAKY JOBS:\n")
				}
			} else if stat.ID == testgridReportDetails {
				if meta.Flags.EmojisOff {
					fmt.Printf("%s severity:%d, %s\n", stat.Status, stat.Severity, stat.Title)
				} else {
					fmt.Printf("%s %s %s\n", stat.Status, stat.Highlight, stat.Title)
				}
				fmt.Printf("- %s\n", stat.URL)
				for _, note := range stat.Notes {
					fmt.Printf("- %s\n", note)
				}
			}
		}
	}
}

// PutData extends TestgridReport and stores the data at runtime to the struct val ReportData
func (r *TestgridReport) PutData(reportData ReportData) {
	r.ReportData = reportData
}

// GetData extends TestgridReport and returns the data that has been stored at runtime int the struct val ReportData (counter to SaveData/1)
func (r TestgridReport) GetData() ReportData {
	return r.ReportData
}

func assembleTestgridRequests(meta Meta, requiredJobs []testgridJob) chan ReportDataField {
	c := make(chan ReportDataField)
	go func() {
		defer close(c)
		wg := sync.WaitGroup{}
		for _, j := range requiredJobs {
			wg.Add(1)
			go func(job testgridJob) {
				jobBaseURL := fmt.Sprintf("https://testgrid.k8s.io/%s", job.URLName)
				jobsData, err := reqTestgridSiteData(job, jobBaseURL)
				if err != nil {
					log.Fatalf("error %v", err)
				}
				records := []ReportDataRecord{getSummary(jobsData)}

				if !meta.Flags.ShortOn {
					for jobName, jobData := range jobsData {
						if jobData.OverallStatus != passing {
							records = append(records, getDetails(jobName, jobData, jobBaseURL, meta.Flags.EmojisOff))
						}
					}
				}

				reportData := ReportDataField{
					Emoji:   job.Emoji,
					Title:   job.OutputName,
					Records: records,
				}
				c <- reportData
				wg.Done()
			}(j)
		}
		wg.Wait()
	}()
	return c
}

// This function is used to request job summary data from a testgrid subpage
func reqTestgridSiteData(job testgridJob, jobBaseURL string) (TestgridData, error) {
	// This url points to testgrid/summary which returns a JSON document
	url := fmt.Sprintf("%s/summary", jobBaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// Parse body form http request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Unmarshal JSON from body into TestgridJobsOverview struct
	jobs, err := UnmarshalTestgrid(body)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// This function is used to count up the status from testgrid tests
func getSummary(jobs map[string]testgridValue) ReportDataRecord {
	result := ReportDataRecord{ID: testgridReportSummary}
	statuses := map[overallStatus]int{total: len(jobs), passing: 0, failing: 0, flaky: 0, stale: 0}
	for _, v := range jobs {
		if v.OverallStatus == passing {
			statuses[passing]++
		} else if v.OverallStatus == failing {
			statuses[failing]++
		} else if v.OverallStatus == flaky {
			statuses[flaky]++
		} else {
			statuses[stale]++
		}
	}
	result.Notes = append(result.Notes, fmt.Sprintf("%d jobs %s", statuses[total], strings.ToLower(string(total))))
	result.Notes = append(result.Notes, fmt.Sprintf("%d jobs %s", statuses[passing], strings.ToLower(string(passing))))
	result.Notes = append(result.Notes, fmt.Sprintf("%d jobs %s", statuses[flaky], strings.ToLower(string(flaky))))
	result.Notes = append(result.Notes, fmt.Sprintf("%d jobs %s", statuses[failing], strings.ToLower(string(failing))))
	if statuses[stale] != 0 {
		result.Notes = append(result.Notes, fmt.Sprintf("%d jobs %s", statuses[stale], strings.ToLower(string(stale))))
	}
	return result
}

// This function is used get additional information about testgrid jobs
func getDetails(jobName string, jobData testgridValue, jobBaseURL string, emojisOff bool) ReportDataRecord {
	result := ReportDataRecord{ID: testgridReportDetails}
	result.Status = string(jobData.OverallStatus)
	result.Title = jobName
	result.URL = fmt.Sprintf("%s#%s", jobBaseURL, jobName)

	// If the status is failing give information about failing tests
	if jobData.OverallStatus == failing {
		// Filter sigs
		sigRegex := regexp.MustCompile(`sig-[a-zA-Z]+`)
		sigsInvolved := map[string]int{}
		for _, test := range jobData.Tests {
			sigs := sigRegex.FindAllString(test.TestName, -1)
			for _, sig := range sigs {
				sigsInvolved[sig] = sigsInvolved[sig] + 1
			}
		}
		sigs := reflect.ValueOf(sigsInvolved).MapKeys()

		result.Notes = append(result.Notes, fmt.Sprintf("Sig's involved %v", sigs))
		result.Notes = append(result.Notes, fmt.Sprintf("Currently %d test are failing", len(jobData.Tests)))
	}

	const (
		testgridRegexRecentRuns   = "runs"
		testgridRegexRecentPasses = "passes"

		thresholdWarning  float64 = 0.5 // 0.0 ... 0.5 -> warning
		thresholdInfo     float64 = 0.8 // 0.5 ... 0.8 -> info
		newTestThreshhold float64 = 5.0 // if 0.0 ... 5.0 -> new test
	)
	// This regex filters the latest executions
	// e.g. "8 of 9 (88.9%) recent columns passed (19455 of 19458 or 100.0% cells)" -> 8 passes of 9 runs recently
	latestExec := getRegexParams(fmt.Sprintf(`(?P<%s>\d{1,2})\sof\s(?P<%s>\d{1,2})`, testgridRegexRecentPasses, testgridRegexRecentRuns), jobData.Status)
	testgridRegexRecentPassesFloat, err := strconv.ParseFloat(latestExec[testgridRegexRecentPasses], 64)
	if err != nil {
		fmt.Println(err)
	}
	testgridRegexRecentRunsFloat, err := strconv.ParseFloat(latestExec[testgridRegexRecentRuns], 64)
	if err != nil {
		fmt.Println(err)
	}

	highlightEmoji := ""
	if jobData.OverallStatus == failing {
		highlightEmoji = statusFailingEmoji
	} else {
		highlightEmoji = statusFlakyEmoji
	}
	recentSuccessRate := testgridRegexRecentPassesFloat / testgridRegexRecentRunsFloat
	severity := Severity(0)
	if testgridRegexRecentRunsFloat <= newTestThreshhold {
		severity = LightSeverity
		highlightEmoji = statusNewEmoji
	} else {
		if recentSuccessRate <= thresholdWarning {
			severity = HighSeverity
		} else if recentSuccessRate <= thresholdInfo {
			severity = MediumSeverity
		} else {
			severity = LightSeverity
		}
	}

	result.Severity = severity
	result.Highlight = strings.Repeat(highlightEmoji, int(severity))

	result.Notes = append(result.Notes, fmt.Sprintf("%s of %s passed recently", latestExec[testgridRegexRecentPasses], latestExec[testgridRegexRecentRuns]))
	return result
}

// Parses string with the given regular expression and returns the group values defined in the expression.
// e.g. `(?P<Year>\d{4})-(?P<Month>\d{2})-(?P<Day>\d{2})` + `2015-05-27` -> map[Year:2015 Month:05 Day:27]
func getRegexParams(regEx, s string) (paramsMap map[string]string) {
	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(s)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

type testgridJob struct {
	OutputName string
	URLName    string
	Emoji      string
}

// The types below reflect testgrid summary json (e.g. https://testgrid.k8s.io/sig-release-master-informing/summary)

// TestgridData contains all jobs under one specific field like 'sig-release-master-informing'
type TestgridData map[string]testgridValue

// UnmarshalTestgrid []byte into TestgridData
func UnmarshalTestgrid(data []byte) (TestgridData, error) {
	var r TestgridData
	err := json.Unmarshal(data, &r)
	return r, err
}

// Marshal TestgridData struct into []bytes
func (r *TestgridData) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// testgridValue information about a specifc job
type testgridValue struct {
	Alert               string        `json:"alert"`
	LastRunTimestamp    int64         `json:"last_run_timestamp"`
	LastUpdateTimestamp int64         `json:"last_update_timestamp"`
	LatestGreen         string        `json:"latest_green"`
	OverallStatus       overallStatus `json:"overall_status"`
	Status              string        `json:"status"`
	Tests               []test        `json:"tests"`
	DashboardName       dashboardName `json:"dashboard_name"`
	Healthiness         healthiness   `json:"healthiness"`
	BugURL              string        `json:"bug_url"`
}

// healthiness
type healthiness struct {
	Tests             []interface{} `json:"tests"`
	PreviousFlakiness int64         `json:"previousFlakiness"`
}

type test struct {
	DisplayName    string        `json:"display_name"`
	TestName       string        `json:"test_name"`
	FailCount      int64         `json:"fail_count"`
	FailTimestamp  int64         `json:"fail_timestamp"`
	PassTimestamp  int64         `json:"pass_timestamp"`
	BuildLink      string        `json:"build_link"`
	BuildURLText   string        `json:"build_url_text"`
	BuildLinkText  string        `json:"build_link_text"`
	FailureMessage string        `json:"failure_message"`
	LinkedBugs     []interface{} `json:"linked_bugs"`
	FailTestLink   string        `json:"fail_test_link"`
}

type dashboardName string

const (
	sigReleaseMasterInforming dashboardName = "sig-release-master-informing"
	sigReleaseMasterBlocking  dashboardName = "sig-release-master-blocking"
)

type overallStatus string

const (
	total   overallStatus = "TOTAL"
	failing overallStatus = "FAILING"
	flaky   overallStatus = "FLAKY"
	passing overallStatus = "PASSING"
	stale   overallStatus = "STALE"
)

// This information is used internally to differentiate between summary and detail ReportDataRecords
const (
	testgridReportSummary = 0
	testgridReportDetails = 1
)
