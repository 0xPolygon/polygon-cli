package streamer

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type MarkdownStreamer struct {
	writer        io.Writer
	quiet         bool
	headerWritten bool
	footerNeeded  bool
}

func NewMarkdownStreamer(writer io.Writer, quiet bool) *MarkdownStreamer {
	return &MarkdownStreamer{
		writer: writer,
		quiet:  quiet,
	}
}

func (m *MarkdownStreamer) StreamTestExecution(exec TestExecution) error {
	if m.quiet {
		return nil
	}

	if !m.headerWritten {
		m.writeHeader()
		m.headerWritten = true
		m.footerNeeded = true
	}

	statusIcon := "✅"
	if exec.Status == "fail" {
		statusIcon = "❌"
	}

	errorText := exec.Error
	if errorText == "" {
		errorText = "-"
	}
	// Escape markdown special characters in error text
	errorText = strings.ReplaceAll(errorText, "|", "\\|")

	argsJSON, _ := json.Marshal(exec.Args)
	resultJSON, _ := json.Marshal(exec.Result)

	_, err := fmt.Fprintf(m.writer, "| %s | %s %s | `%s` | `%s` | %s | %v | %s | %s |\n",
		exec.TestName,
		statusIcon,
		exec.Method,
		string(argsJSON), string(resultJSON),
		exec.Status,
		exec.Duration.String(),
		exec.Timestamp.Format("15:04:05"),
		errorText,
	)

	return err
}

func (m *MarkdownStreamer) StreamSummary(summary TestSummary) error {
	// Markdown doesn't need interim summaries in table format
	return nil
}

func (m *MarkdownStreamer) StreamFinalSummary(summaries []TestSummary) error {
	if m.footerNeeded {
		m.writeFooter(summaries)
	}
	return nil
}

func (m *MarkdownStreamer) writeHeader() {
	fmt.Fprint(m.writer, `# RPC Fuzz Test Results

  ## Test Executions

  | Test Name | Method | Args | Result | Status | Duration | Time | Error |
  |-----------|--------|------|--------|--------|----------|------|-------|
  `)
}

func (m *MarkdownStreamer) writeFooter(summaries []TestSummary) {
	fmt.Fprint(m.writer, `

  ## Summary

  `)

	totalRan, totalPassed, totalFailed := 0, 0, 0
	for _, s := range summaries {
		totalRan += s.TestsRan
		totalPassed += s.TestsPassed
		totalFailed += s.TestsFailed

		successRate := s.SuccessRate * 100
		fmt.Fprintf(m.writer, "- **%s**: %d/%d passed (%.1f%%)\n",
			s.TestName,
			s.TestsPassed,
			s.TestsRan,
			successRate,
		)
	}

	overallRate := float64(totalPassed) / float64(totalRan) * 100
	fmt.Fprintf(m.writer, `
  ---
  **Overall**: %d tests completed, %d passed, %d failed (%.1f%% success rate)

  *Generated at: %s*
  `,
		totalRan, totalPassed, totalFailed, overallRate,
		time.Now().Format("2006-01-02 15:04:05"),
	)
}
