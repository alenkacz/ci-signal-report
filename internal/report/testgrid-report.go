package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/leonardpahlke/ci-signal-report/internal/config"
	"github.com/leonardpahlke/ci-signal-report/internal/models"
)

// This function is used to accumulate a summary of testgrid
func RequestTestgridOverview(meta config.Meta) ([]models.TestGridStatistics, error) {
	// The report checks master-blocking and master-informing
	requiredJobs := []models.TestgridJob{
		{OutputName: "Master-Blocking", UrlName: "sig-release-master-blocking"},
		{OutputName: "Master-Informing", UrlName: "sig-release-master-informing"},
	}

	// If a release version got specified add additional jobs to report
	if meta.Env.ReleaseVersion != "" {
		requiredJobsVersion := []models.TestgridJob{
			{OutputName: meta.Env.ReleaseVersion + "-blocking", UrlName: "sig-release-" + meta.Env.ReleaseVersion + "-blocking"},
			{OutputName: meta.Env.ReleaseVersion + "-informing", UrlName: "sig-release-" + meta.Env.ReleaseVersion + "-informing"},
		}
		for i := range requiredJobsVersion {
			requiredJobs = append(requiredJobs, requiredJobsVersion[i])
		}
	}

	testgridStats := make([]models.TestGridStatistics, 0)
	for _, kubeJob := range requiredJobs {
		// Request Testgrid subpage summary data
		jobs, err := requestTestgridSiteSummary(kubeJob)
		if err != nil {
			return nil, err
		}
		statistics := getStatistics(jobs)
		statistics.Name = kubeJob.OutputName
		testgridStats = append(testgridStats, statistics)
	}
	return testgridStats, nil
}

// This function is used to request job summary data from a testgrid subpage
func requestTestgridSiteSummary(kubeJob models.TestgridJob) (models.TestgrudJobsOverview, error) {
	// This url points to testgrid/summary which returns a JSON document
	url := fmt.Sprintf("https://testgrid.k8s.io/%s/summary", kubeJob.UrlName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// Parse body form http request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Unmarshal JSON from body into models.TestgrudJobsOverview struct
	jobs := make(models.TestgrudJobsOverview)
	err = json.Unmarshal(body, &jobs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// This function is used to count up the status from testgrid tests
func getStatistics(jobs map[string]models.TestgridOverview) models.TestGridStatistics {
	result := models.TestGridStatistics{}
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
