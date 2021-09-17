package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alenkacz/ci-signal-report/internal/config"
	"github.com/alenkacz/ci-signal-report/internal/models"
)

func PrintTestgridJobsStatistics(meta config.Meta) {
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
		statistics.Name = kubeJob.OutputName
		testgridStats = append(testgridStats, statistics)
	}

	for _, stat := range testgridStats {
		fmt.Printf("Failures in %s\n", stat.Name)
		fmt.Printf("\t%d jobs total\n", stat.Total)
		fmt.Printf("\t%d are passing\n", stat.Passing)
		fmt.Printf("\t%d are flaking\n", stat.Flaking)
		fmt.Printf("\t%d are failing\n", stat.Failing)
		fmt.Printf("\t%d are stale\n", stat.Stale)
		fmt.Print("\n\n")
	}
}

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

func getStatsFromJson(body []byte) (models.TestgrudJobsOverview, error) {
	result := make(models.TestgrudJobsOverview)
	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
