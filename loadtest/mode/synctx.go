package mode

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/montanaflynn/stats"
	"github.com/rs/zerolog/log"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
)

// EIP-7966 error codes returned by eth_sendRawTransactionSync. These are
// method-specific codes, not the standard JSON-RPC ones, and each carries a
// data field: the transaction hash for timeout and queued, the expected nonce
// for a nonce gap.
const (
	// SyncErrTimeout means the transaction reached the mempool but no receipt
	// was produced before the timeout elapsed. The transaction may still be
	// mined afterwards.
	SyncErrTimeout = 4
	// SyncErrQueued means the transaction was not added to the mempool because
	// it is not ready for execution.
	SyncErrQueued = 5
	// SyncErrNonceGap means the transaction was not added to the mempool
	// because of a nonce gap. The error data carries the expected nonce.
	SyncErrNonceGap = 6
)

// Codes bor returns from this method beyond the EIP-7966 set. Bor implements
// only code 4 from the EIP; a transaction it refuses to accept fails in
// SubmitTransaction first and comes back with one of these instead, so they are
// classified separately rather than lumped in with unknown errors.
const (
	SyncErrBorNonceTooLow       = -38010
	SyncErrBorNonceTooHigh      = -38011
	SyncErrBorIntrinsicGas      = -38013
	SyncErrBorInsufficientFunds = -38014
	SyncErrBorClientLimit       = -38026
)

// SyncReceipt is the receipt returned by eth_sendRawTransactionSync. EIP-7966
// specifies "the transaction receipt object as defined by the
// eth_getTransactionReceipt method", but a node answering before canonical
// inclusion may leave the block fields empty, so every field is optional here
// rather than decoded into types.Receipt, which requires them.
type SyncReceipt struct {
	TransactionHash common.Hash     `json:"transactionHash"`
	Status          *hexutil.Uint64 `json:"status"`
	GasUsed         *hexutil.Uint64 `json:"gasUsed"`
	BlockNumber     *hexutil.Big    `json:"blockNumber"`
	BlockHash       *common.Hash    `json:"blockHash"`

	// Preconfirmation is bor's marker for a receipt produced by its
	// preconfirmation pipeline rather than canonical block import. It is not
	// part of EIP-7966, so nodes that do not preconfirm simply omit it.
	Preconfirmation *bool `json:"preconfirmation"`
}

// Speculative reports whether the receipt is a preconfirmation rather than a
// canonical one.
//
// Bor states this outright with a "preconfirmation" field, and that answer is
// taken when present. EIP-7966 itself defines no such marker, so for every
// other node this falls back to reading the block fields: bor's preconfirmation
// receipts carry no block hash, and a receipt with no block cannot be canonical
// on any implementation. A node that both omits the marker and fills in a
// speculative block number is indistinguishable from one answering canonically;
// pair with --wait-for-receipt to confirm canonical inclusion independently.
func (r *SyncReceipt) Speculative() bool {
	if r == nil {
		return false
	}
	if r.Preconfirmation != nil {
		return *r.Preconfirmation
	}
	if r.BlockNumber == nil || r.BlockNumber.ToInt().Sign() == 0 {
		return true
	}
	return r.BlockHash == nil || *r.BlockHash == (common.Hash{})
}

// Succeeded reports whether the receipt carries a success status. A receipt
// with no status field at all is not treated as a success.
func (r *SyncReceipt) Succeeded() bool {
	return r != nil && r.Status != nil && uint64(*r.Status) == types.ReceiptStatusSuccessful
}

// SendRawTransactionSync submits a signed transaction with
// eth_sendRawTransactionSync (EIP-7966), which returns once the node has a
// receipt for it or the timeout elapses, and records the outcome on
// deps.SyncTracker when one is set.
//
// The returned error is the RPC error verbatim, so the caller records the same
// failure the node reported. Note that a timeout is an error even though the
// transaction is in the mempool and may still be mined.
func SendRawTransactionSync(ctx context.Context, deps *Dependencies, cfg *config.Config, tx *types.Transaction) error {
	rawTx, err := tx.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	// Omit the timeout entirely when unset so the node applies its own default
	// (EIP-7966 recommends 2s; bor defaults to 20s and clamps anything above
	// its --rpc.txsync.maxtimeout instead of rejecting it).
	//
	// Encoding: bor takes this parameter as *hexutil.Uint64, which unmarshals
	// only from a quoted hex quantity and rejects a bare JSON number, so hex is
	// the default. EIP-7966 describes the parameter as an integer and
	// spec-literal servers reject the hex form, hence --sync-tx-timeout-int.
	params := []any{hexutil.Encode(rawTx)}
	if timeout := cfg.SyncTxTimeout; timeout > 0 {
		ms := uint64(timeout.Milliseconds())
		if cfg.SyncTxTimeoutInt {
			params = append(params, ms)
		} else {
			params = append(params, hexutil.Uint64(ms))
		}
	}

	var receipt *SyncReceipt
	start := time.Now()
	err = deps.SendRPCClient.CallContext(ctx, &receipt, "eth_sendRawTransactionSync", params...)
	elapsed := time.Since(start)

	deps.SyncTracker.Record(tx.Hash(), receipt, elapsed, err)

	return err
}

// SyncErrorCode returns the EIP-7966 error code carried by an
// eth_sendRawTransactionSync failure, and whether the error had one at all.
func SyncErrorCode(err error) (int, bool) {
	var rpcErr ethrpc.Error
	if errors.As(err, &rpcErr) {
		return rpcErr.ErrorCode(), true
	}
	return 0, false
}

// SyncTxResult holds the outcome of a single synchronous submission.
type SyncTxResult struct {
	TxHash      string `json:"tx_hash"`
	DurationMs  int64  `json:"duration_ms"`
	Speculative bool   `json:"speculative,omitempty"`
	BlockNumber uint64 `json:"block_number,omitempty"`
	GasUsed     uint64 `json:"gas_used,omitempty"`
	Status      uint64 `json:"status,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Error       string `json:"error,omitempty"`
}

// SyncTracker aggregates eth_sendRawTransactionSync outcomes across routines.
// A nil tracker records nothing, so callers never need to check for one.
type SyncTracker struct {
	submitted   atomic.Uint64
	receipts    atomic.Uint64
	speculative atomic.Uint64
	canonical   atomic.Uint64
	reverted    atomic.Uint64
	noStatus    atomic.Uint64

	timeouts  atomic.Uint64
	queued    atomic.Uint64
	nonceGaps atomic.Uint64
	rejected  atomic.Uint64
	otherErrs atomic.Uint64

	mu      sync.Mutex
	results []SyncTxResult
}

// NewSyncTracker creates a tracker for synchronous submission outcomes.
func NewSyncTracker() *SyncTracker {
	return &SyncTracker{results: make([]SyncTxResult, 0, 1024)}
}

// Record files the outcome of one submission.
func (t *SyncTracker) Record(txHash common.Hash, receipt *SyncReceipt, elapsed time.Duration, err error) {
	if t == nil {
		return
	}

	t.submitted.Add(1)

	result := SyncTxResult{
		TxHash:     txHash.Hex(),
		DurationMs: elapsed.Milliseconds(),
	}

	if err != nil {
		result.Error = err.Error()
		if code, ok := SyncErrorCode(err); ok {
			result.ErrorCode = code
			switch code {
			case SyncErrTimeout:
				t.timeouts.Add(1)
			case SyncErrQueued:
				t.queued.Add(1)
			case SyncErrNonceGap, SyncErrBorNonceTooHigh:
				t.nonceGaps.Add(1)
			case SyncErrBorNonceTooLow, SyncErrBorIntrinsicGas,
				SyncErrBorInsufficientFunds, SyncErrBorClientLimit:
				// Refused before the wait even started.
				t.rejected.Add(1)
			default:
				t.otherErrs.Add(1)
			}
		} else {
			t.otherErrs.Add(1)
		}

		t.append(result)
		return
	}

	if receipt == nil {
		// No error and no receipt: the node answered, but with nothing usable.
		t.otherErrs.Add(1)
		result.Error = "empty receipt"
		t.append(result)
		return
	}

	t.receipts.Add(1)
	result.Speculative = receipt.Speculative()
	if result.Speculative {
		t.speculative.Add(1)
	} else {
		t.canonical.Add(1)
		// A node can call a receipt canonical with the marker while leaving the
		// block fields out, so this stays guarded.
		if receipt.BlockNumber != nil {
			result.BlockNumber = receipt.BlockNumber.ToInt().Uint64()
		}
	}
	if receipt.GasUsed != nil {
		result.GasUsed = uint64(*receipt.GasUsed)
	}
	switch {
	case receipt.Status == nil:
		t.noStatus.Add(1)
	case receipt.Succeeded():
		result.Status = types.ReceiptStatusSuccessful
	default:
		t.reverted.Add(1)
	}

	t.append(result)
}

func (t *SyncTracker) append(result SyncTxResult) {
	t.mu.Lock()
	t.results = append(t.results, result)
	t.mu.Unlock()
}

// Results returns a copy of the per-transaction outcomes.
func (t *SyncTracker) Results() []SyncTxResult {
	if t == nil {
		return nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]SyncTxResult, len(t.results))
	copy(out, t.results)
	return out
}

// Stats logs a summary of the synchronous submissions.
func (t *SyncTracker) Stats() {
	if t == nil || t.submitted.Load() == 0 {
		return
	}

	results := t.Results()
	durations := make([]float64, 0, len(results))
	for _, r := range results {
		durations = append(durations, float64(r.DurationMs))
	}
	pct := func(p float64) float64 {
		v, err := stats.Percentile(durations, p)
		if err != nil {
			return 0
		}
		return v
	}

	log.Info().
		Uint64("submitted", t.submitted.Load()).
		Uint64("receipts", t.receipts.Load()).
		Uint64("speculative", t.speculative.Load()).
		Uint64("canonical", t.canonical.Load()).
		Uint64("reverted", t.reverted.Load()).
		Uint64("no_status", t.noStatus.Load()).
		Uint64("timeouts", t.timeouts.Load()).
		Uint64("queued", t.queued.Load()).
		Uint64("nonce_gaps", t.nonceGaps.Load()).
		Uint64("rejected", t.rejected.Load()).
		Uint64("other_errors", t.otherErrs.Load()).
		Float64("p50_ms", pct(50)).
		Float64("p90_ms", pct(90)).
		Float64("p99_ms", pct(99)).
		Msg("Synchronous transaction submission stats")
}
