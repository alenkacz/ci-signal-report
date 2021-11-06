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
	"sync"
)

// TestgridReport used to implement RequestData & Print for testgrid report data
type TestgridReport struct {
	ReportData ReportData
}

// RequestTestgridOverview this function is used to accumulate a summary of testgrid
func (r *TestgridReport) RequestData(meta Meta, wg *sync.WaitGroup) ReportData {
	// The report checks master-blocking and master-informing
	requiredJobs := []testgridJob{
		{OutputName: "Master-Blocking", URLName: "sig-release-master-blocking", Emoji: masterBlockingEmoji},
		{OutputName: "Master-Informing", URLName: "sig-release-master-informing", Emoji: masterInformingEmoji},
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
		headerLine := fmt.Sprintf("%s Tests in %s", reportField.Emoji, reportField.Title)
		if meta.Flags.EmojisOff {
			headerLine = fmt.Sprintf("Tests in %s", reportField.Title)
		}
		for _, stat := range reportField.Records {
			fmt.Println(headerLine)
			fmt.Printf("\t%d jobs total\n", stat.Total)
			fmt.Printf("\t%d are passing\n", stat.Passing)
			fmt.Printf("\t%d are flaking\n", stat.Flaking)
			fmt.Printf("\t%d are failing\n", stat.Failing)
			fmt.Printf("\t%d are stale\n", stat.Stale)
			fmt.Print("\n\n")
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
				jobs, err := requestTestgridSiteSummary(job)
				if err != nil {
					log.Fatalf("error %v", err)
				}
				reportData := ReportDataField{
					Emoji:   job.Emoji,
					Title:   job.OutputName,
					Records: []ReportDataRecord{getStatistics(jobs)},
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
func requestTestgridSiteSummary(job testgridJob) (testgridJobsOverview, error) {
	// This url points to testgrid/summary which returns a JSON document
	url := fmt.Sprintf("https://testgrid.k8s.io/%s/summary", job.URLName)
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
	jobs := make(testgridJobsOverview)
	err = json.Unmarshal(body, &jobs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// This function is used to count up the status from testgrid tests
func getStatistics(jobs map[string]testgridOverview) ReportDataRecord {
	result := ReportDataRecord{}
	for _, v := range jobs {
		if v.OverallStatus == "PASSING" {
			result.Passing++
		} else if v.OverallStatus == "FAILING" {
			result.Failing++
		} else if v.OverallStatus == "FLAKY" {
			result.Flaking++
		} else {
			result.Stale++
		}
		result.Total++
	}
	return result
}

type testgridJob struct {
	OutputName string
	URLName    string
	Emoji      string
}

// TestGridStatistics information as summary about a testgrid area (like master-blocking)
type TestGridStatistics struct {
	Emoji   string
	Name    string
	Total   int
	Passing int
	Flaking int
	Failing int
	Stale   int
}

type testgridJobsOverview = map[string]testgridOverview
type testgridOverview struct {
	OverallStatus string `json:"overall_status"`
}
