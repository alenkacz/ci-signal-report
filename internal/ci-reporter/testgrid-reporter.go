package cireporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// RequestTestgridOverview this function is used to accumulate a summary of testgrid
func RequestTestgridOverview(meta Meta) ([]TestGridStatistics, error) {
	// The report checks master-blocking and master-informing
	requiredJobs := []testgridJob{
		{OutputName: "Master-Blocking", URLName: "sig-release-master-blocking", Emoji: masterBlockingEmoji},
		{OutputName: "Master-Informing", URLName: "sig-release-master-informing", Emoji: masterInformingEmoji},
	}

	// If a release version got specified add additional jobs to report
	if meta.Env.ReleaseVersion != "" {
		requiredJobsVersion := []testgridJob{
			{OutputName: fmt.Sprintf("%s-blocking", meta.Env.ReleaseVersion), URLName: fmt.Sprintf("sig-release-%s-blocking", meta.Env.ReleaseVersion), Emoji: masterBlockingEmoji},
			{OutputName: fmt.Sprintf("%s-informing", meta.Env.ReleaseVersion), URLName: fmt.Sprintf("sig-release-%s-informing", meta.Env.ReleaseVersion), Emoji: masterInformingEmoji},
		}
		for i := range requiredJobsVersion {
			requiredJobs = append(requiredJobs, requiredJobsVersion[i])
		}
	}

	testgridStats := make([]TestGridStatistics, 0)
	for _, job := range requiredJobs {
		// Request Testgrid subpage summary data
		jobs, err := requestTestgridSiteSummary(job)
		if err != nil {
			return nil, err
		}
		statistics := getStatistics(jobs)
		statistics.Name = job.OutputName
		statistics.Emoji = job.Emoji
		testgridStats = append(testgridStats, statistics)
	}
	return testgridStats, nil
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
func getStatistics(jobs map[string]testgridOverview) TestGridStatistics {
	result := TestGridStatistics{}
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
