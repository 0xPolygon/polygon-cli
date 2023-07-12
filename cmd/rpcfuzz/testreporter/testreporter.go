// Package testreporter provides the utilities to capture, report, and log test results.
package testreporter

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rs/zerolog/log"
)

type (
	TestResult struct {
		Name                string          `json:"name"`
		Method              string          `json:"method"`
		Args                [][]interface{} `json:"args"`
		Result              []interface{}   `json:"result"`
		Errors              []error         `json:"errors"`
		NumberOfTestsRan    int             `json:"numberOfTestsRan"`
		NumberOfTestsPassed int             `json:"numberOfTestsPassed"`
		NumberOfTestsFailed int             `json:"numberOfTestsFailed"`
	}
	TestResults struct {
		Tests       []TestResult
		TableWriter table.Writer
	}
)

func New(testName string, testMethod string, numOfTestRuns int) TestResult {
	return TestResult{
		Name:             testName,
		Method:           testMethod,
		NumberOfTestsRan: numOfTestRuns,
	}
}

func (tr *TestResult) Pass(args []interface{}, result interface{}, err error) {
	tr.NumberOfTestsPassed++
	tr.Args = append(tr.Args, args)
	tr.Result = append(tr.Result, result)
	tr.Errors = append(tr.Errors, err)
}

func (tr *TestResult) Fail(args []interface{}, result interface{}, err error) {
	tr.NumberOfTestsFailed++
	tr.Args = append(tr.Args, args)
	tr.Result = append(tr.Result, result)
	tr.Errors = append(tr.Errors, err)
}

func (trs *TestResults) AddTestResult(testResult TestResult) {
	trs.Tests = append(trs.Tests, testResult)
}

func createAndWriteToFile(filePath string, content []byte) {
	directoryPath := filepath.Dir(filePath)
	if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
		log.Error().Err(err).Msg("Error while trying to create file directory")
		return
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Error().Err(err).Msg("Error while trying to create file")
		return
	}
	defer f.Close()

	bufWriter := bufio.NewWriter(f)
	_, err = bufWriter.Write(content)
	if err != nil {
		log.Error().Err(err).Msg("Error while trying to write to file")
		return
	}

	bufWriter.Flush()
}

func (trs *TestResults) ExportResultToJSON(filePath string) {
	jsonContent, err := json.MarshalIndent(trs.Tests, "", "\t")
	if err != nil {
		log.Error().Err(err).Msg("Error while trying to marshal test results to json")
		return
	}
	createAndWriteToFile(filePath, jsonContent)
}

func (trs *TestResults) ExportResultToCSV(filePath string) {
	createAndWriteToFile(filePath, []byte(trs.TableWriter.RenderCSV()))
}

func (trs *TestResults) ExportResultToMarkdown(filePath string) {
	createAndWriteToFile(filePath, []byte(trs.TableWriter.RenderMarkdown()))
}

func (trs *TestResults) ExportResultToHTML(filePath string) {
	createAndWriteToFile(filePath, []byte(trs.TableWriter.RenderHTML()))
}

func (trs *TestResults) GenerateTabularResult() {
	trs.TableWriter = table.NewWriter()
	trs.TableWriter.AppendHeader(table.Row{"Name", "Method", "Test(s) Passed", "Test(s) Ran"})
	trs.TableWriter.SortBy([]table.SortBy{
		{Name: "Method", Mode: table.Asc},
		{Name: "Name", Mode: table.Asc},
	})

	for _, currTestResult := range trs.Tests {
		trs.TableWriter.AppendRow([]interface{}{currTestResult.Name, currTestResult.Method, currTestResult.NumberOfTestsPassed, currTestResult.NumberOfTestsRan})
	}
}

func (trs *TestResults) PrintTabularResult() {
	trs.TableWriter.SetOutputMirror(os.Stdout)
	trs.TableWriter.SetStyle(table.StyleColoredBright)
	trs.TableWriter.Render()
}

func (trs *TestResults) PrintSummary() {
	for _, tr := range trs.Tests {
		log.Info().Str("Method", tr.Method).Msgf("%d/%d Test(s) Passed", tr.NumberOfTestsPassed, tr.NumberOfTestsRan)
	}
}
