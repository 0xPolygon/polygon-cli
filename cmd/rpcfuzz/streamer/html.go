package streamer

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"time"
)

type HTMLStreamer struct {
	writer        io.Writer
	quiet         bool
	headerWritten bool
	footerNeeded  bool
}

func NewHTMLStreamer(writer io.Writer, quiet bool) *HTMLStreamer {
	return &HTMLStreamer{
		writer: writer,
		quiet:  quiet,
	}
}

func (h *HTMLStreamer) StreamTestExecution(exec TestExecution) error {
	if h.quiet {
		return nil
	}

	if !h.headerWritten {
		h.writeHeader()
		h.headerWritten = true
		h.footerNeeded = true
	}

	statusClass := "pass"
	statusIcon := "✅"
	if exec.Status == "fail" {
		statusClass = "fail"
		statusIcon = "❌"
	}

	argsJSON, _ := json.Marshal(exec.Args)
	resultJSON, _ := json.Marshal(exec.Result)

	_, err := fmt.Fprintf(h.writer,
		`    <tr class="%s">
          <td>%s</td>
          <td>%s %s</td>
          <td>%s</td>
          <td>%v</td>
          <td>%s</td>
          <td>%s</td>
          <td><pre>%s</pre></td>
          <td><pre>%s</pre></td>
      </tr>
  `,
		statusClass,
		html.EscapeString(exec.TestName),
		statusIcon,
		html.EscapeString(exec.Method),
		html.EscapeString(exec.Status),
		exec.Duration.String(),
		exec.Timestamp.Format("15:04:05"),
		html.EscapeString(exec.Error),
		html.EscapeString(string(argsJSON)),
		html.EscapeString(string(resultJSON)),
	)

	return err
}

func (h *HTMLStreamer) StreamSummary(summary TestSummary) error {
	// HTML doesn't need interim summaries in table format
	return nil
}

func (h *HTMLStreamer) StreamFinalSummary(summaries []TestSummary) error {
	if h.footerNeeded {
		h.writeFooter(summaries)
	}
	return nil
}

func (h *HTMLStreamer) writeHeader() {
	fmt.Fprint(h.writer, `<!DOCTYPE html>
  <html>
  <head>
      <title>RPC Fuzz Test Results</title>
      <style>
          body { font-family: Arial, sans-serif; margin: 20px; }
          table { border-collapse: collapse; width: 100%; }
          th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
          th { background-color: #f2f2f2; }
          .pass { background-color: #f0fff0; }
          .fail { background-color: #fff0f0; }
          .summary { margin-top: 20px; padding: 10px; background-color: #f9f9f9; }
      </style>
  </head>
  <body>
      <h1>RPC Fuzz Test Results</h1>
      <table>
          <thead>
              <tr>
                <th>Test Name</th>
                <th>Method</th>
                <th>Status</th>
                <th>Duration</th>
                <th>Time</th>
                <th>Error</th>
                <th>Args</th>
                <th>Result</th>
             </tr>
          </thead>
          <tbody>
  `)
}

func (h *HTMLStreamer) writeFooter(summaries []TestSummary) {
	fmt.Fprint(h.writer, `        </tbody>
      </table>
      <div class="summary">
          <h2>Summary</h2>
  `)

	totalRan, totalPassed, totalFailed := 0, 0, 0
	for _, s := range summaries {
		totalRan += s.TestsRan
		totalPassed += s.TestsPassed
		totalFailed += s.TestsFailed

		successRate := s.SuccessRate * 100
		fmt.Fprintf(h.writer,
			`        <p><strong>%s</strong>: %d/%d passed (%.1f%%)</p>
  `,
			html.EscapeString(s.TestName),
			s.TestsPassed,
			s.TestsRan,
			successRate,
		)
	}

	overallRate := float64(totalPassed) / float64(totalRan) * 100
	fmt.Fprintf(h.writer,
		`        <hr>
          <p><strong>Overall</strong>: %d tests completed, %d passed, %d failed (%.1f%% success rate)</p>
          <p>Generated at: %s</p>
      </div>
  </body>
  </html>
  `,
		totalRan, totalPassed, totalFailed, overallRate,
		time.Now().Format("2006-01-02 15:04:05"),
	)
}
