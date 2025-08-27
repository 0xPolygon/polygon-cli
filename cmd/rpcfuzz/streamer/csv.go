package streamer

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"strconv"
)

type CSVStreamer struct {
	writer        io.Writer
	csvWriter     *csv.Writer
	quiet         bool
	headerWritten bool
}

func NewCSVStreamer(writer io.Writer, quiet bool) *CSVStreamer {
	return &CSVStreamer{
		writer:    writer,
		csvWriter: csv.NewWriter(writer),
		quiet:     quiet,
	}
}

func (c *CSVStreamer) StreamTestExecution(exec TestExecution) error {
	if c.quiet {
		return nil
	}

	if !c.headerWritten {
		c.csvWriter.Write([]string{"test_name", "method", "args", "result", "status", "duration_ms", "timestamp", "error"})
		c.headerWritten = true
	}

	durationMs := strconv.FormatInt(exec.Duration.Milliseconds(), 10)
	argsJSON, _ := json.Marshal(exec.Args)
	resultJSON, _ := json.Marshal(exec.Result)
	record := []string{
		exec.TestName,
		exec.Method,
		string(argsJSON),
		string(resultJSON),
		exec.Status,
		durationMs,
		exec.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		exec.Error,
	}

	err := c.csvWriter.Write(record)
	c.csvWriter.Flush()
	return err
}

func (c *CSVStreamer) StreamSummary(summary TestSummary) error {
	// CSV doesn't need interim summaries, just final
	return nil
}

func (c *CSVStreamer) StreamFinalSummary(summaries []TestSummary) error {
	c.csvWriter.Flush()
	return c.csvWriter.Error()
}
