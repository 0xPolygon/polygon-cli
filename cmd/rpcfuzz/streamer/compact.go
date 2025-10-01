package streamer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type CompactStreamer struct {
	writer io.Writer
	quiet  bool
}

func NewCompactStreamer(writer io.Writer, quiet bool) *CompactStreamer {
	return &CompactStreamer{writer: writer, quiet: quiet}
}

func (c *CompactStreamer) StreamTestExecution(exec TestExecution) error {
	if c.quiet {
		return nil // Don't output individual results in quiet mode
	}

	status := "✅"
	if exec.Status == "fail" {
		status = "❌"
	}
	argsJSON, _ := json.Marshal(exec.Args)
	resultJSON, _ := json.Marshal(exec.Result)
	_, err := fmt.Fprintf(c.writer, "%s %s %s (%v) | Args: %s | Result: %s\n",
		status, exec.Method, exec.Status, exec.Duration, string(argsJSON), string(resultJSON))
	return err
}

func (c *CompactStreamer) StreamSummary(summary TestSummary) error {
	_, err := fmt.Fprintf(os.Stderr, "\rProgress: %s - %d/%d passed (%.1f%%)",
		summary.TestName, summary.TestsPassed, summary.TestsRan, summary.SuccessRate*100)
	return err
}

func (c *CompactStreamer) StreamFinalSummary(summaries []TestSummary) error {
	fmt.Fprintf(os.Stderr, "\n") // Clear progress line

	totalRan, totalPassed, totalFailed := 0, 0, 0
	for _, s := range summaries {
		totalRan += s.TestsRan
		totalPassed += s.TestsPassed
		totalFailed += s.TestsFailed
	}

	successRate := float64(totalPassed) / float64(totalRan) * 100
	_, err := fmt.Fprintf(c.writer, "\nFinal Summary: %d tests completed, %d passed, %d failed (%.1f%% success rate)\n",
		totalRan, totalPassed, totalFailed, successRate)
	return err
}
