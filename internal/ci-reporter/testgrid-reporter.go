package ci_reporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// This function is used to accumulate a summary of testgrid
func RequestTestgridOverview(meta CiReporterMeta) ([]TestGridStatistics, error) {
	// The report checks master-blocking and master-informing
	requiredJobs := []TestgridJob{
		{OutputName: "Master-Blocking", UrlName: "sig-release-master-blocking", Emoji: MasterBlockingEmoji},
		{OutputName: "Master-Informing", UrlName: "sig-release-master-informing", Emoji: MasterInformingEmoji},
	}

	// If a release version got specified add additional jobs to report
	if meta.Env.ReleaseVersion != "" {
		requiredJobsVersion := []TestgridJob{
			{OutputName: meta.Env.ReleaseVersion + "-blocking", UrlName: "sig-release-" + meta.Env.ReleaseVersion + "-blocking"},
			{OutputName: meta.Env.ReleaseVersion + "-informing", UrlName: "sig-release-" + meta.Env.ReleaseVersion + "-informing"},
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
func requestTestgridSiteSummary(job TestgridJob) (TestgridJobsOverview, error) {
	// This url points to testgrid/summary which returns a JSON document
	url := fmt.Sprintf("https://testgrid.k8s.io/%s/summary", job.UrlName)
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
	jobs := make(TestgridJobsOverview)
	err = json.Unmarshal(body, &jobs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// This function is used to count up the status from testgrid tests
func getStatistics(jobs map[string]TestgridOverview) TestGridStatistics {
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

type TestgridJob struct {
	OutputName string
	UrlName    string
	Emoji      string
}

type TestGridStatistics struct {
	Emoji   string
	Name    string
	Total   int
	Passing int
	Flaking int
	Failing int
	Stale   int
}

type TestgridJobsOverview = map[string]TestgridOverview
type TestgridOverview struct {
	OverallStatus string `json:"overall_status"`
}
