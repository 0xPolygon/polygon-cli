package loadtest

import (
	"time"

	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/ethereum/go-ethereum/common"
)

// Sample represents a single load test request/response.
type Sample struct {
	GoRoutineID int64
	RequestID   int64
	RequestTime time.Time
	WaitTime    time.Duration
	Receipt     string
	IsError     bool
	Nonce       uint64
}

// BlockSummary holds data about a single block's transactions.
type BlockSummary struct {
	Block     *rpctypes.RawBlockResponse
	Receipts  map[common.Hash]rpctypes.RawTxReceipt
	Latencies map[uint64]time.Duration
}

// Latency holds min, median, and max latency values.
type Latency struct {
	Min    float64
	Median float64
	Max    float64
}

// Summary holds summary data for a single block.
type Summary struct {
	BlockNumber uint64
	Time        time.Time
	GasLimit    uint64
	GasUsed     uint64
	NumTx       int
	Utilization float64
	Latencies   Latency
}

// SummaryOutput holds the complete summary output data.
type SummaryOutput struct {
	Summaries          []Summary
	SuccessfulTx       int64
	TotalTx            int64
	TotalMiningTime    time.Duration
	TotalGasUsed       uint64
	TransactionsPerSec float64
	GasPerSecond       float64
	Latencies          Latency
}

// BlobCommitment holds blob transaction commitment data.
type BlobCommitment struct {
	Blob          [131072]byte // kzg4844.Blob size
	Commitment    [48]byte     // kzg4844.Commitment size
	Proof         [48]byte     // kzg4844.Proof size
	VersionedHash common.Hash
}
