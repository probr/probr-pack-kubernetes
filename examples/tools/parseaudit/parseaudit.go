// This tool parses audit logs and flatten data in csv format
// Sample usage:
// 	go run ./examples/tools/parseaudit/parseaudit.go 							>> (will use default output directory from config to find json files)
// 	go run ./examples/tools/parseaudit/parseaudit.go "/path/to/auditlogs/" 		>> (will use cli arg as path to json files)

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/citihub/probr/audit"
	"github.com/citihub/probr/config"
)

type flatObj struct {
	SourceFile      string
	ProbeName       string
	ScenarioID      int
	ScenarioName    string
	StepID          int
	StepDescription string
	StepFunction    string
	StepPayload     string
}

func main() {

	config.Init("")
	auditLogsDir := config.Vars.AuditDir()

	if len(os.Args[1:]) > 0 {
		auditLogsDir = os.Args[1] // Override config default
	}

	fmt.Println(fmt.Printf("Parsing audit logs from: '%s' ...", auditLogsDir))

	files, err := getAllFiles(auditLogsDir)
	if err != nil {
		log.Panicf("failed reading directory: %s", err)
	}
	fmt.Printf("Number of files in current directory: %d \n", len(files))

	// Files
	for _, fileInfo := range files {

		if !strings.HasSuffix(fileInfo.Name(), ".json") { //Only parse json files
			continue
		}

		var flatRows []flatObj

		filePath := path.Join(auditLogsDir, fileInfo.Name())
		probeAudit, err := deserializeJSON(filePath)
		check(err)

		// Scenarios
		for scenarioID, scenario := range probeAudit.Scenarios {

			// Steps
			for stepID, step := range scenario.Steps {
				var flatRow flatObj

				flatRow.SourceFile = filePath // FlatRow

				flatRow.ScenarioID = scenarioID
				flatRow.ScenarioName = scenario.Name

				flatRow.StepID = stepID
				flatRow.StepDescription = step.Description
				flatRow.StepFunction = step.Name

				flatRows = append(flatRows, flatRow)
			}
		}

		// Sort collection
		sort.Sort(byScenarioAndStepID(flatRows))

		// Write file
		csvFilePath := filePath + ".csv"
		writeCSVFile(csvFilePath, flatRows, true)
		fmt.Println("Created: ", csvFilePath)
	}
}

func getAllFiles(dir string) ([]os.FileInfo, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Panicf("failed reading directory: %s", err)
	}
	return entries, err
}

func check(e error) {
	if e != nil {
		log.Panicf("failed to perform action: %s", e)
	}
}

func deserializeJSON(filePath string) (audit.ProbeAudit, error) {
	data, err := ioutil.ReadFile(filePath)
	check(err)

	// Deserialize Json
	var probeAudit audit.ProbeAudit
	err = json.Unmarshal(data, &probeAudit)
	check(err)

	return probeAudit, err
}

func writeCSVFile(filePath string, rows []flatObj, addHeader bool) error {
	csvFile, err := os.Create(filePath)
	check(err)
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)

	// Header
	if addHeader {
		var row []string
		row = append(row, "SourceFile")
		row = append(row, "ScenarioID")
		row = append(row, "ScenarioName")
		row = append(row, "StepID")
		row = append(row, "StepDescription")
		row = append(row, "StepFunction")
		writer.Write(row)
	}

	for _, r := range rows {
		var row []string
		row = append(row, r.SourceFile)
		row = append(row, strconv.Itoa(r.ScenarioID))
		row = append(row, r.ScenarioName)
		row = append(row, strconv.Itoa(r.StepID))
		row = append(row, r.StepDescription)
		row = append(row, r.StepFunction)
		writer.Write(row)
	}

	writer.Flush()

	return err
}

// Implementing Sort interface

type byScenarioAndStepID []flatObj

func (s byScenarioAndStepID) Len() int {
	return len(s)
}
func (s byScenarioAndStepID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byScenarioAndStepID) Less(i, j int) bool {
	if s[i].ScenarioID == s[j].ScenarioID {
		return s[i].StepID < s[j].StepID
	}
	return s[i].ScenarioID < s[j].ScenarioID
}

// Implementing Sort interface
