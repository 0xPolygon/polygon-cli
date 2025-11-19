package streamer

import (
	"encoding/json"
	"io"
)

type JSONStreamer struct {
	writer io.Writer
}

func NewJSONStreamer(writer io.Writer) *JSONStreamer {
	return &JSONStreamer{writer: writer}
}

func (j *JSONStreamer) StreamTestExecution(exec TestExecution) error {
	data := map[string]any{
		"type": "test_execution",
		"data": exec,
	}
	return json.NewEncoder(j.writer).Encode(data)
}

func (j *JSONStreamer) StreamSummary(summary TestSummary) error {
	data := map[string]any{
		"type": "test_summary",
		"data": summary,
	}
	return json.NewEncoder(j.writer).Encode(data)
}

func (j *JSONStreamer) StreamFinalSummary(summaries []TestSummary) error {
	data := map[string]any{
		"type": "final_summary",
		"data": summaries,
	}
	return json.NewEncoder(j.writer).Encode(data)
}
