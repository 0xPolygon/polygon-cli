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
