package models

type TestgridJob struct {
	OutputName string
	UrlName    string
}

type TestGridStatistics struct {
	Name    string
	Total   int
	Passing int
	Flaking int
	Failing int
	Stale   int
}

type TestgrudJobsOverview = map[string]TestgridOverview
type TestgridOverview struct {
	OverallStatus string `json:"overall_status"`
}
