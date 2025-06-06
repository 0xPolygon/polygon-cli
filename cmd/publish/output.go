package publish

import (
	"fmt"
	"sync/atomic"
	"time"
)

const (
	InputDataSourceFile  = "file"
	InputDataSourceArgs  = "args"
	InputDataSourceStdin = "stdin"
)

type inputDataSource string

type output struct {
	InputDataSource       inputDataSource
	InputDataCount        atomic.Uint64
	ValidInputs           atomic.Uint64
	InvalidInputs         atomic.Uint64
	StartTime             time.Time
	EndTime               time.Time
	TxsSentSuccessfully   atomic.Uint64
	TxsSentUnsuccessfully atomic.Uint64
}

func (s *output) Start() {
	s.StartTime = time.Now()
}

func (s *output) Stop() {
	s.EndTime = time.Now()
}

func (s *output) Print() {
	elapsed := s.EndTime.Sub(s.StartTime)
	elapsedSeconds := elapsed.Seconds()
	if elapsedSeconds == 0 {
		elapsedSeconds = 1
	}
	txSent := s.TxsSentSuccessfully.Load() + s.TxsSentUnsuccessfully.Load()
	txsSendPerSecond := 0.0
	if elapsedSeconds > 0.0001 {
		txsSendPerSecond = float64(txSent) / elapsedSeconds
	}
	successRatio := float64(0)
	if txSent > 0 {
		successRatio = float64(s.TxsSentSuccessfully.Load()) / float64(txSent) * 100
	}

	summaryString := fmt.Sprintf(`-----------------------------------
              Summary              
-----------------------------------
Concurrency: %d
JobQueueSize: %d
RateLimit: %d
-----------------------------------
Input Data Source: %s
Input Data Count: %d
Valid Inputs: %d
Invalid Inputs: %d
-----------------------------------
Elapsed Time: %s
Txs Sent: %d
Txs Sent Per Second: %.2f
Txs Sent Successfully: %d
Txs Sent Unsuccessfully: %d
Success Ratio: %.2f%%
-----------------------------------`,
		*publishInputArgs.concurrency,
		*publishInputArgs.jobQueueSize,
		*publishInputArgs.rateLimit,

		s.InputDataSource,
		s.InputDataCount.Load(),
		s.ValidInputs.Load(),
		s.InvalidInputs.Load(),

		elapsed,
		txSent,
		txsSendPerSecond,
		s.TxsSentSuccessfully.Load(),
		s.TxsSentUnsuccessfully.Load(),
		successRatio)

	fmt.Println(summaryString)
}
