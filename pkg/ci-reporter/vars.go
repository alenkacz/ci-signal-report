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

// Github card ids
const (
	newCardsID                   = 4212817
	underInvestigationCardsID    = 4212819
	observingCardsID             = 4212821
	githubCiSignalBoardProjectID = 2093513
)

// Emojis
const (
	inFlightEmoji        = "\U0001F6EB"
	notYetStartedEmoji   = "\U0001F914"
	observingEmoji       = "\U0001F440"
	resolvedEmoji        = "\U0001F389"
	masterBlockingEmoji  = "\U000026D4"
	masterInformingEmoji = "\U0001F4A1"
)

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

// PrintJson pretty print json to console
func (r *Report) PrintJson() {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal Report %v", err)
	}
	fmt.Print(string(b))
}

// Report wraps multiple report data objects
type Report []ReportData

// ReportData
type ReportData struct {
	Data []ReportDataField `json:"data"`
	// Name like 'github' or 'testgrid'
	Name string `json:"name"`
}

// ReportDataField
type ReportDataField struct {
	Emoji   string             `json:"emoji"`
	Title   string             `json:"title"`
	Records []ReportDataRecord `json:"records"`
}

// ReportDataRecord
type ReportDataRecord struct {
	Total   int `json:"total"`
	Passing int `json:"passing"`
	Flaking int `json:"flaking"`
	Failing int `json:"failing"`
	Stale   int `json:"stale"`

	URL   string `json:"url"`
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Sig   string `json:"sig"`
}
