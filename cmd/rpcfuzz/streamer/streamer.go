package streamer

type OutputStreamer interface {
	StreamTestExecution(exec TestExecution) error
	StreamSummary(summary TestSummary) error
	StreamFinalSummary(summaries []TestSummary) error
}
