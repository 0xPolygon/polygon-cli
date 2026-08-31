package p2p

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"

	ds "github.com/0xPolygon/polygon-cli/p2p/datastructures"
)

// Rebroadcast filter drop reasons, used as the metric label and in logs.
const (
	dropReasonStaleNonce  = "stale_nonce"
	dropReasonLowTip      = "low_tip"
	dropReasonRateLimited = "rate_limited"
)

const (
	// nonceLookupQueueSize bounds the RPC fallback backlog. Lookups are a
	// best-effort warm-up of the nonce cache, so a full queue drops the request
	// rather than blocking the message loop.
	nonceLookupQueueSize = 4096
	// nonceLookupWorkers bounds concurrent eth_getTransactionCount calls.
	nonceLookupWorkers = 4
	// nonceLookupTimeout bounds a single fallback lookup.
	nonceLookupTimeout = 10 * time.Second
	// nonceLookupRetryAfter is how long a sender stays in the in-flight set,
	// suppressing duplicate lookups for the same address.
	nonceLookupRetryAfter = time.Minute
)

// TxFilterOptions configures the rebroadcast gates. The zero value gates
// nothing, so a filter only does work for the gates that are turned on.
type TxFilterOptions struct {
	// ChainID selects the signer used to recover transaction senders.
	ChainID uint64

	// GateStaleTxs drops transactions whose nonce is below the sender's known
	// next nonce. These can never be mined, so echoing them is pure
	// amplification -- this is the gate that catches replayed history.
	GateStaleTxs bool

	// NonceCache sizes the sender -> next-nonce map that backs the stale gate.
	NonceCache ds.LRUOptions

	// NonceRPCURL, when set, is used to look up nonces for senders that have
	// not been observed in a block yet. Lookups are asynchronous: the
	// transaction that triggered one is allowed through, and the result gates
	// later transactions from that sender.
	NonceRPCURL string

	// MinTip drops transactions offering less than this tip (in wei). On
	// Polygon, bor will not include anything below 25 gwei, so those are
	// unminable regardless of anything else. Zero disables the gate.
	MinTip uint64

	// RateLimit caps rebroadcast throughput in transactions per second, a
	// backstop that bounds amplification even if the gates above miss. Zero or
	// negative disables the cap.
	RateLimit float64
	// RateLimitBurst is the token bucket depth. Defaults to one second of
	// RateLimit when unset.
	RateLimitBurst int

	// LogOnly evaluates every gate and records what would have been dropped,
	// but rebroadcasts everything anyway. Use it to size the problem before
	// turning the gates on.
	LogOnly bool
}

// enabled reports whether any gate is configured.
func (o TxFilterOptions) enabled() bool {
	return o.GateStaleTxs || o.MinTip > 0 || o.RateLimit > 0
}

// TxFilter decides which transactions the sensor rebroadcasts. It never
// affects what is recorded to the database: the sensor observes everything and
// gates only the echo, the same split the block path uses for signer
// validation.
//
// The gates run cheapest-first, and the rate limiter runs last so that
// transactions dropped by an earlier gate do not consume budget.
type TxFilter struct {
	opts   TxFilterOptions
	signer types.Signer
	minTip *big.Int

	// nonces maps a sender to the next nonce it can validly use. It is fed for
	// free from the transactions of every block the sensor already observes,
	// with the RPC fallback filling in senders that have not appeared in one.
	nonces *ds.LRU[common.Address, uint64]

	limiter *rate.Limiter

	rpc      *ethclient.Client
	lookups  chan common.Address
	inflight *ds.LRU[common.Address, struct{}]
}

// NewTxFilter creates a transaction rebroadcast filter. It returns nil when no
// gate is enabled, and callers must treat a nil filter as "allow everything"
// (see TxFilter.Allow), so the no-gates path costs nothing.
func NewTxFilter(ctx context.Context, opts TxFilterOptions) (*TxFilter, error) {
	if !opts.enabled() {
		return nil, nil
	}

	f := &TxFilter{
		opts:   opts,
		signer: types.LatestSignerForChainID(new(big.Int).SetUint64(opts.ChainID)),
	}

	if opts.MinTip > 0 {
		f.minTip = new(big.Int).SetUint64(opts.MinTip)
	}

	if opts.RateLimit > 0 {
		burst := opts.RateLimitBurst
		if burst <= 0 {
			burst = int(opts.RateLimit)
		}
		if burst <= 0 {
			burst = 1
		}
		f.limiter = rate.NewLimiter(rate.Limit(opts.RateLimit), burst)
	}

	if opts.GateStaleTxs {
		f.nonces = ds.NewLRU[common.Address, uint64](opts.NonceCache)

		if opts.NonceRPCURL != "" {
			client, err := ethclient.DialContext(ctx, opts.NonceRPCURL)
			if err != nil {
				return nil, fmt.Errorf("dialing nonce lookup rpc %s: %w", opts.NonceRPCURL, err)
			}
			f.rpc = client
			f.lookups = make(chan common.Address, nonceLookupQueueSize)
			f.inflight = ds.NewLRU[common.Address, struct{}](ds.LRUOptions{
				MaxSize: nonceLookupQueueSize,
				TTL:     nonceLookupRetryAfter,
			})
			for range nonceLookupWorkers {
				go f.lookupWorker(ctx)
			}
		}
	}

	log.Info().
		Bool("gate_stale_txs", opts.GateStaleTxs).
		Bool("nonce_rpc_fallback", f.rpc != nil).
		Uint64("min_tip_wei", opts.MinTip).
		Float64("rate_limit", opts.RateLimit).
		Bool("log_only", opts.LogOnly).
		Msg("Transaction rebroadcast gates enabled")

	return f, nil
}

// Close releases the RPC client used for nonce fallback lookups.
func (f *TxFilter) Close() {
	if f == nil || f.rpc == nil {
		return
	}
	f.rpc.Close()
}

// ObserveMined records the sender nonces advanced by a mined block. Every block
// the sensor sees maintains the stale gate for free: a mined transaction proves
// its sender's next nonce is at least one higher.
func (f *TxFilter) ObserveMined(txs []*types.Transaction) {
	if f == nil || f.nonces == nil || len(txs) == 0 {
		return
	}

	for _, tx := range txs {
		sender, err := types.Sender(f.signer, tx)
		if err != nil {
			continue
		}
		f.setNonce(sender, tx.Nonce()+1)
	}

	txFilterKnownSenders.Set(float64(f.nonces.Len()))
}

// setNonce advances a sender's known next nonce, never lowering it. Blocks and
// RPC responses can arrive out of order, so the highest value wins.
func (f *TxFilter) setNonce(sender common.Address, next uint64) {
	f.nonces.Update(sender, func(current uint64) uint64 {
		if next > current {
			return next
		}
		return current
	})
}

// Allow returns the transactions that may be rebroadcast. A nil filter allows
// everything, so callers do not need to check for one. In LogOnly mode every
// transaction is returned, but the drop counters still move.
func (f *TxFilter) Allow(txs []*types.Transaction) []*types.Transaction {
	if f == nil || len(txs) == 0 {
		return txs
	}

	allowed := make([]*types.Transaction, 0, len(txs))
	for _, tx := range txs {
		if reason, ok := f.reject(tx); ok {
			txFilterDropped.WithLabelValues(reason).Inc()
			log.Trace().
				Stringer("hash", tx.Hash()).
				Uint64("nonce", tx.Nonce()).
				Str("reason", reason).
				Msg("Withholding transaction from rebroadcast")
			if !f.opts.LogOnly {
				continue
			}
		} else {
			txFilterAllowed.Inc()
		}
		allowed = append(allowed, tx)
	}

	return allowed
}

// reject applies each gate in increasing order of cost and reports the first
// one that fires.
func (f *TxFilter) reject(tx *types.Transaction) (reason string, rejected bool) {
	// Cheapest: a tip below the chain's floor can never be included, and needs
	// neither sender recovery nor any lookup.
	if f.minTip != nil && tx.GasTipCap().Cmp(f.minTip) < 0 {
		return dropReasonLowTip, true
	}

	// Costs an ecrecover, cached on the transaction by go-ethereum after the
	// first call, and a map hit.
	if f.nonces != nil && f.isStale(tx) {
		return dropReasonStaleNonce, true
	}

	// Last, so that transactions the gates above rejected do not spend budget.
	if f.limiter != nil && !f.limiter.Allow() {
		return dropReasonRateLimited, true
	}

	return "", false
}

// isStale reports whether a transaction's nonce is below its sender's known
// next nonce, meaning it cannot be mined. Senders that have not been observed
// are not stale as far as we know: the transaction is allowed and, when the RPC
// fallback is configured, a lookup is queued so later ones can be judged.
func (f *TxFilter) isStale(tx *types.Transaction) bool {
	sender, err := types.Sender(f.signer, tx)
	if err != nil {
		// Unsigned or malformed for this chain's signer. Not a staleness
		// judgment, and the decode path already accepted it, so leave it alone.
		return false
	}

	next, known := f.nonces.Get(sender)
	if !known {
		f.queueLookup(sender)
		return false
	}

	return tx.Nonce() < next
}

// queueLookup schedules an asynchronous nonce lookup for an unobserved sender.
// It never blocks: a full queue means the sender stays unknown for now.
func (f *TxFilter) queueLookup(sender common.Address) {
	if f.lookups == nil {
		return
	}

	// Suppress duplicate lookups for the same sender while one is outstanding.
	if existed := f.inflight.Update(sender, func(v struct{}) struct{} { return v }); existed {
		return
	}

	select {
	case f.lookups <- sender:
	default:
		txFilterNonceLookups.WithLabelValues("dropped").Inc()
		f.inflight.Remove(sender)
	}
}

// lookupWorker serves the RPC fallback queue until the context is cancelled.
func (f *TxFilter) lookupWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case sender := <-f.lookups:
			f.lookupNonce(ctx, sender)
		}
	}
}

// lookupNonce fetches a sender's next nonce and caches it. Failures are left
// uncached so the sender can be retried once its in-flight entry expires.
func (f *TxFilter) lookupNonce(ctx context.Context, sender common.Address) {
	lookupCtx, cancel := context.WithTimeout(ctx, nonceLookupTimeout)
	defer cancel()

	nonce, err := f.rpc.NonceAt(lookupCtx, sender, nil)
	if err != nil {
		txFilterNonceLookups.WithLabelValues("error").Inc()
		log.Debug().Err(err).Stringer("sender", sender).Msg("Failed to look up account nonce")
		f.inflight.Remove(sender)
		return
	}

	txFilterNonceLookups.WithLabelValues("ok").Inc()
	f.setNonce(sender, nonce)
	txFilterKnownSenders.Set(float64(f.nonces.Len()))
}
