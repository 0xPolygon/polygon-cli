package streamer

import (
	"time"
)

// Individual test execution result - streamed immediately
type TestExecution struct {
	TestName  string        `json:"test_name"`
	Method    string        `json:"method"`
	Args      []any         `json:"args"`
	Result    any           `json:"result,omitempty"`
	Error     string        `json:"error,omitempty"`
	Status    string        `json:"status"` // "pass" or "fail"
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

// Summary stats - no individual data stored
type TestSummary struct {
	TestName      string        `json:"test_name"`
	Method        string        `json:"method"`
	TestsRan      int           `json:"tests_ran"`
	TestsPassed   int           `json:"tests_passed"`
	TestsFailed   int           `json:"tests_failed"`
	SuccessRate   float64       `json:"success_rate"`
	TotalDuration time.Duration `json:"total_duration"`
}
