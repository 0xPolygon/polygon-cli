// Package testreporter provides the utilities to capture, report, and log test results.
package testreporter

import (
	"github.com/rs/zerolog/log"
)

type (
	TestResult struct {
		Name                string
		Method              string
		Args                [][]interface{}
		Result              []interface{}
		Errors              []error
		NumberOfTestsRan    int
		NumberOfTestsPassed int
		NumberOfTestsFailed int
	}
	TestResults []TestResult
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

// TODO:
// export to json
// export to csv
func (trs *TestResults) PrintTabularResult() {
}

func (trs *TestResults) PrintSummary() {
	for _, tr := range *trs {
		log.Info().Str("Method", tr.Method).Msgf("%d/%d Test(s) Passed", tr.NumberOfTestsPassed, tr.NumberOfTestsRan)
	}
}
