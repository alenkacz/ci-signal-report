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
	"log"
	"sync"
)

// Reports
const (
	githubReport   = "github"
	testgridReport = "testgrid"
)

// Emojis
const (
	inFlightEmoji        = "\U0001F6EB"
	notYetStartedEmoji   = "\U0001F914"
	observingEmoji       = "\U0001F440"
	resolvedEmoji        = "\U0001F389"
	masterBlockingEmoji  = "\U0001F525"
	masterInformingEmoji = "\U0001F4A1"
	statusFailingEmoji   = "\U0001F534"
	statusFlakyEmoji     = "\U0001F535"
	statusNewEmoji       = "\U00002728"
	statusOldEmoji       = "\U0001F319"
)

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	// colorYellow = "\033[33m"
	colorBlue = "\033[34m"
	// colorPurple = "\033[35m"
	// colorCyan   = "\033[36m"
	// colorWhite  = "\033[37m"
)

// Severity used to rank report records
type Severity int

// HighSeverity, MediumSeverity, LightSeverity used to rank report records from 0...3
const (
	HighSeverity   Severity = 3
	MediumSeverity Severity = 2
	LightSeverity  Severity = 1
)

// CIReport this interface to implement Reporters
type CIReport interface {
	RequestData(meta Meta, wg *sync.WaitGroup) ReportData
	Print(meta Meta, reportData ReportData)
	PutData(reportData ReportData)
	GetData() ReportData
}

// UnmarshalReport transforms a json obj into a Report struct
func UnmarshalReport(data []byte) (Report, error) {
	var r Report
	err := json.Unmarshal(data, &r)
	return r, err
}

// Marshal method to transform a Report into JSON format
func (r *Report) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// PrintJSON pretty print json to console
func (r *Report) PrintJSON() {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal Report %v", err)
	}
	fmt.Print(string(b))
}

// Report wraps multiple report data objects
type Report []ReportData

// ReportData that contains multiple data fields
type ReportData struct {
	Data []ReportDataField `json:"data"`
	// Name like 'github' or 'testgrid'
	Name string `json:"name"`
}

// ReportDataField one field of a report that contains multiple records
type ReportDataField struct {
	Emoji   string             `json:"emoji"`
	Title   string             `json:"title"`
	Records []ReportDataRecord `json:"records"`
}

// ReportDataRecord that contain specifc information about a testgrid job or about a github issue (flexible)
type ReportDataRecord struct {
	// record url
	URL string `json:"url"`
	// record identifier
	ID int64 `json:"id"`
	// record title
	Title string `json:"title"`
	// k8s sig reference
	Sig string `json:"sig"`
	// collection of additional information
	Notes []string `json:"notes"`
	// record status
	Status string `json:"status"`
	// can be set to show importance
	Severity Severity `json:"severity"`
	// can be set to highlight the record (with an emoji for example)
	Highlight string `json:"highlight"`
}
