package database

import (
	"context"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/0xPolygon/polygon-cli/util"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"
)

// Default batching parameters. ClickHouse strongly prefers large, infrequent
// inserts over many small ones, so every Write* call only enqueues a row into
// an in-memory buffer; a background goroutine flushes it in batches. This keeps
// the sensor hot path non-blocking and turns the write pattern into the bulk
// append ClickHouse is built for. Buffers are fixed-size and drop-on-full so a
// slow/unavailable database can never exhaust memory or stall the sensor.
const (
	chFlushInterval   = 1 * time.Second
	chFlushTimeout    = 30 * time.Second
	chBlockBatch      = 5000
	chBlockEventBatch = 50000
	chTxBatch         = 20000
	chTxEventBatch    = 50000
	chPeerBatch       = 2000
)

// ClickHouse implements the Database interface backed by a ClickHouse cluster.
// See clickhouse_schema.sql for the table definitions this writer targets.
type ClickHouse struct {
	conn                             driver.Conn
	sensorID                         string
	chainID                          *big.Int
	maxConcurrency                   int
	shouldWriteBlocks                bool
	shouldWriteBlockEvents           bool
	shouldWriteFirstBlockEvent       bool
	shouldWriteTransactions          bool
	shouldWriteTransactionEvents     bool
	shouldWriteFirstTransactionEvent bool
	shouldWritePeers                 bool

	blocks   *rowBatcher[chBlock]
	blockEvt *rowBatcher[chEvent]
	txs      *rowBatcher[chTx]
	txEvt    *rowBatcher[chEvent]
	peers    *rowBatcher[chPeer]
}

// ClickHouseOptions is used when creating a NewClickHouse.
type ClickHouseOptions struct {
	DSN                              string
	SensorID                         string
	ChainID                          uint64
	MaxConcurrency                   int
	ShouldWriteBlocks                bool
	ShouldWriteBlockEvents           bool
	ShouldWriteFirstBlockEvent       bool
	ShouldWriteTransactions          bool
	ShouldWriteTransactionEvents     bool
	ShouldWriteFirstTransactionEvent bool
	ShouldWritePeers                 bool
	TTL                              time.Duration
}

// NewClickHouse connects to ClickHouse, verifies connectivity, and starts the
// background batch flushers. The provided context governs the lifetime of the
// flusher goroutines; when it is cancelled they flush any buffered rows and
// exit. If the connection cannot be established the returned Database no-ops all
// writes (mirroring the Datastore backend) so the sensor keeps running.
func NewClickHouse(ctx context.Context, opts ClickHouseOptions) Database {
	c := &ClickHouse{
		sensorID:                         opts.SensorID,
		chainID:                          new(big.Int).SetUint64(opts.ChainID),
		maxConcurrency:                   opts.MaxConcurrency,
		shouldWriteBlocks:                opts.ShouldWriteBlocks,
		shouldWriteBlockEvents:           opts.ShouldWriteBlockEvents,
		shouldWriteFirstBlockEvent:       opts.ShouldWriteFirstBlockEvent,
		shouldWriteTransactions:          opts.ShouldWriteTransactions,
		shouldWriteTransactionEvents:     opts.ShouldWriteTransactionEvents,
		shouldWriteFirstTransactionEvent: opts.ShouldWriteFirstTransactionEvent,
		shouldWritePeers:                 opts.ShouldWritePeers,
	}

	chOpts, err := clickhouse.ParseDSN(opts.DSN)
	if err != nil {
		log.Error().Err(err).Msg("Could not parse ClickHouse DSN")
		return c
	}

	conn, err := clickhouse.Open(chOpts)
	if err != nil {
		log.Error().Err(err).Msg("Could not connect to ClickHouse")
		return c
	}

	if err := conn.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("Could not ping ClickHouse")
		return c
	}
	c.conn = conn

	c.blocks = newRowBatcher(ctx, "blocks", chBlockBatch, func(fctx context.Context, rows []chBlock) error {
		return c.flush(fctx, "INSERT INTO blocks (hash, number, parent_hash, block_time, coinbase, signer, difficulty, total_difficulty, gas_used, gas_limit, base_fee, tx_count, uncle_count, sensor_id, ingested_at)", func(b driver.Batch) error {
			for _, r := range rows {
				if err := b.Append(r.hash, r.number, r.parentHash, r.blockTime, r.coinbase, r.signer, r.difficulty, r.totalDifficulty, r.gasUsed, r.gasLimit, r.baseFee, r.txCount, r.uncleCount, c.sensorID, r.ingestedAt); err != nil {
					return err
				}
			}
			return nil
		})
	})
	c.blockEvt = newRowBatcher(ctx, "block_events", chBlockEventBatch, func(fctx context.Context, rows []chEvent) error {
		return c.flush(fctx, "INSERT INTO block_events (block_hash, sensor_id, peer_id, seen_at)", func(b driver.Batch) error {
			for _, r := range rows {
				if err := b.Append(r.hash, c.sensorID, r.peerID, r.seenAt); err != nil {
					return err
				}
			}
			return nil
		})
	})
	c.txs = newRowBatcher(ctx, "transactions", chTxBatch, func(fctx context.Context, rows []chTx) error {
		return c.flush(fctx, "INSERT INTO transactions (hash, from_address, to_address, value, gas, gas_price, gas_fee_cap, gas_tip_cap, nonce, tx_type, first_seen, sensor_id, ingested_at)", func(b driver.Batch) error {
			for _, r := range rows {
				if err := b.Append(r.hash, r.from, r.to, r.value, r.gas, r.gasPrice, r.gasFeeCap, r.gasTipCap, r.nonce, r.txType, r.firstSeen, c.sensorID, r.ingestedAt); err != nil {
					return err
				}
			}
			return nil
		})
	})
	c.txEvt = newRowBatcher(ctx, "transaction_events", chTxEventBatch, func(fctx context.Context, rows []chEvent) error {
		return c.flush(fctx, "INSERT INTO transaction_events (tx_hash, sensor_id, peer_id, seen_at)", func(b driver.Batch) error {
			for _, r := range rows {
				if err := b.Append(r.hash, c.sensorID, r.peerID, r.seenAt); err != nil {
					return err
				}
			}
			return nil
		})
	})
	c.peers = newRowBatcher(ctx, "peers", chPeerBatch, func(fctx context.Context, rows []chPeer) error {
		return c.flush(fctx, "INSERT INTO peers (peer_id, name, url, caps, last_seen_by, time_last_seen)", func(b driver.Batch) error {
			for _, r := range rows {
				if err := b.Append(r.peerID, r.name, r.url, r.caps, c.sensorID, r.timeLastSeen); err != nil {
					return err
				}
			}
			return nil
		})
	})

	return c
}

// flush prepares a batch for the given INSERT, appends every row via the
// provided callback, and sends it. It uses a detached, time-bounded context so
// a flush triggered during shutdown (parent context already cancelled) still
// completes.
func (c *ClickHouse) flush(_ context.Context, query string, append func(driver.Batch) error) error {
	fctx, cancel := context.WithTimeout(context.Background(), chFlushTimeout)
	defer cancel()

	b, err := c.conn.PrepareBatch(fctx, query)
	if err != nil {
		return err
	}
	if err := append(b); err != nil {
		return err
	}
	return b.Send()
}

// --- row types -------------------------------------------------------------

type chBlock struct {
	hash            string
	number          uint64
	parentHash      string
	blockTime       time.Time
	coinbase        string
	signer          string
	difficulty      uint64
	totalDifficulty *big.Int
	gasUsed         uint64
	gasLimit        uint64
	baseFee         uint64
	txCount         uint32
	uncleCount      uint16
	ingestedAt      time.Time
}

type chEvent struct {
	hash   string
	peerID string
	seenAt time.Time
}

type chTx struct {
	hash       string
	from       string
	to         string
	value      *big.Int
	gas        uint64
	gasPrice   *big.Int
	gasFeeCap  *big.Int
	gasTipCap  *big.Int
	nonce      uint64
	txType     uint8
	firstSeen  time.Time
	ingestedAt time.Time
}

type chPeer struct {
	peerID       string
	name         string
	url          string
	caps         []string
	timeLastSeen time.Time
}

// --- Database interface ----------------------------------------------------

func (c *ClickHouse) WriteBlock(ctx context.Context, peer *enode.Node, block *types.Block, td *big.Int, tfs time.Time) {
	if c.conn == nil {
		return
	}
	if c.shouldWriteBlockEvents && peer != nil {
		c.blockEvt.add(chEvent{hash: block.Hash().Hex(), peerID: peer.URLv4(), seenAt: tfs})
	}
	if c.shouldWriteBlocks {
		c.blocks.add(c.newBlock(block.Header(), td, tfs, len(block.Transactions()), len(block.Uncles())))
	}
	if c.shouldWriteTransactions {
		c.writeTxs(block.Transactions(), tfs)
	}
}

func (c *ClickHouse) WriteBlockHeaders(ctx context.Context, headers []*types.Header, tfs time.Time, isParent bool) {
	if c.conn == nil || !c.shouldWriteBlocks {
		return
	}
	// A header alone carries no transaction/uncle counts; they default to 0 and
	// are populated by the full-block (NewBlock) path when available.
	for _, h := range headers {
		c.blocks.add(c.newBlock(h, big.NewInt(0), tfs, 0, 0))
	}
}

func (c *ClickHouse) WriteBlockBody(ctx context.Context, body *eth.BlockBody, hash common.Hash, tfs time.Time) {
	if c.conn == nil || !c.shouldWriteTransactions {
		return
	}
	// The block row is written from the header path; here we only persist the
	// transactions carried in the body (no read-modify-write on blocks).
	txs, err := body.Transactions.Items()
	if err != nil {
		log.Error().Err(err).Str("hash", hash.Hex()).Msg("Failed to decode transactions from block body")
		return
	}
	c.writeTxs(txs, tfs)
}

func (c *ClickHouse) WriteBlockHashes(ctx context.Context, peer *enode.Node, hashes []common.Hash, tfs time.Time) {
	if c.conn == nil || !c.shouldWriteBlockEvents || peer == nil {
		return
	}
	for _, hash := range hashes {
		c.blockEvt.add(chEvent{hash: hash.Hex(), peerID: peer.URLv4(), seenAt: tfs})
	}
}

func (c *ClickHouse) WriteBlockHashFirstSeen(ctx context.Context, peer *enode.Node, hash common.Hash, tfsh time.Time) {
	if c.conn == nil {
		return
	}
	// Earliest first-seen is derived at query time from block_events (see the
	// block_first_seen materialized view), so we only need to record the event.
	// This mirrors the Datastore backend's first-block-event behavior.
	if c.shouldWriteFirstBlockEvent && !c.shouldWriteBlockEvents && peer != nil {
		c.blockEvt.add(chEvent{hash: hash.Hex(), peerID: peer.URLv4(), seenAt: tfsh})
	}
}

func (c *ClickHouse) WriteTransactions(ctx context.Context, peer *enode.Node, txs []*types.Transaction, tfs time.Time) {
	if c.conn == nil {
		return
	}
	if c.shouldWriteTransactions {
		c.writeTxs(txs, tfs)
	}
	if c.shouldWriteTransactionEvents && peer != nil {
		for _, tx := range txs {
			c.txEvt.add(chEvent{hash: tx.Hash().Hex(), peerID: peer.URLv4(), seenAt: tfs})
		}
	}
}

func (c *ClickHouse) WritePeers(ctx context.Context, peers []*p2p.Peer, tls time.Time) {
	if c.conn == nil || !c.shouldWritePeers {
		return
	}
	for _, peer := range peers {
		c.peers.add(chPeer{
			peerID:       peer.ID().String(),
			name:         peer.Fullname(),
			url:          peer.Node().URLv4(),
			caps:         peer.Info().Caps,
			timeLastSeen: tls,
		})
	}
}

// HasBlock reports whether the block already exists. It is called once per new
// block (not per event), so a lightweight indexed point lookup is cheap and
// preserves the parent-backfill behavior of the Datastore backend.
func (c *ClickHouse) HasBlock(ctx context.Context, hash common.Hash) bool {
	if c.conn == nil {
		return true
	}
	var exists uint8
	err := c.conn.QueryRow(ctx, "SELECT 1 FROM blocks WHERE hash = ? LIMIT 1", hash.Hex()).Scan(&exists)
	return err == nil && exists == 1
}

func (c *ClickHouse) NodeList(ctx context.Context, limit int) ([]string, error) {
	if c.conn == nil {
		return []string{}, nil
	}
	rows, err := c.conn.Query(ctx,
		"SELECT peer_id FROM block_events GROUP BY peer_id ORDER BY max(seen_at) DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodelist := []string{}
	for rows.Next() {
		var peerID string
		if err := rows.Scan(&peerID); err != nil {
			log.Error().Err(err).Msg("Failed to scan peer_id")
			continue
		}
		nodelist = append(nodelist, peerID)
	}
	return nodelist, rows.Err()
}

func (c *ClickHouse) MaxConcurrentWrites() int           { return c.maxConcurrency }
func (c *ClickHouse) ShouldWriteBlocks() bool            { return c.shouldWriteBlocks }
func (c *ClickHouse) ShouldWriteBlockEvents() bool       { return c.shouldWriteBlockEvents }
func (c *ClickHouse) ShouldWriteFirstBlockEvent() bool   { return c.shouldWriteFirstBlockEvent }
func (c *ClickHouse) ShouldWriteTransactions() bool      { return c.shouldWriteTransactions }
func (c *ClickHouse) ShouldWriteTransactionEvents() bool { return c.shouldWriteTransactionEvents }
func (c *ClickHouse) ShouldWriteFirstTransactionEvent() bool {
	return c.shouldWriteFirstTransactionEvent
}
func (c *ClickHouse) ShouldWritePeers() bool { return c.shouldWritePeers }

// --- helpers ---------------------------------------------------------------

func (c *ClickHouse) newBlock(h *types.Header, td *big.Int, tfs time.Time, txCount, uncleCount int) chBlock {
	baseFee := uint64(0)
	if h.BaseFee != nil {
		baseFee = h.BaseFee.Uint64()
	}
	if td == nil {
		td = big.NewInt(0)
	}
	// Recover the block signer from the header seal so signer-based analytics
	// (double-signers, stolen/sealing blocks, reorgs) don't have to ecrecover on
	// every query. Uses the same recovery as the sensor's block-signer validation
	// and the data-analysis tools. Left empty when it can't be recovered.
	var signer string
	if sig, err := util.Ecrecover(h); err == nil {
		signer = common.BytesToAddress(sig).Hex()
	}
	return chBlock{
		hash:            h.Hash().Hex(),
		number:          h.Number.Uint64(),
		parentHash:      h.ParentHash.Hex(),
		blockTime:       time.Unix(int64(h.Time), 0).UTC(),
		coinbase:        h.Coinbase.Hex(),
		signer:          signer,
		difficulty:      h.Difficulty.Uint64(),
		totalDifficulty: new(big.Int).Set(td),
		gasUsed:         h.GasUsed,
		gasLimit:        h.GasLimit,
		baseFee:         baseFee,
		txCount:         uint32(txCount),
		uncleCount:      uint16(uncleCount),
		ingestedAt:      tfs,
	}
}

func (c *ClickHouse) writeTxs(txs []*types.Transaction, tfs time.Time) {
	for _, tx := range txs {
		var from, to string
		chainID := tx.ChainId()
		if chainID == nil || chainID.Sign() <= 0 {
			chainID = c.chainID
		}
		if addr, err := types.Sender(types.LatestSignerForChainID(chainID), tx); err == nil {
			from = addr.Hex()
		}
		if tx.To() != nil {
			to = tx.To().Hex()
		}
		c.txs.add(chTx{
			hash:       tx.Hash().Hex(),
			from:       from,
			to:         to,
			value:      new(big.Int).Set(tx.Value()),
			gas:        tx.Gas(),
			gasPrice:   new(big.Int).Set(tx.GasPrice()),
			gasFeeCap:  new(big.Int).Set(tx.GasFeeCap()),
			gasTipCap:  new(big.Int).Set(tx.GasTipCap()),
			nonce:      tx.Nonce(),
			txType:     tx.Type(),
			firstSeen:  tfs,
			ingestedAt: tfs,
		})
	}
}

// --- batching --------------------------------------------------------------

// rowBatcher buffers rows and flushes them in bulk, either when the buffer
// reaches maxRows or on a fixed interval. add is non-blocking: when the buffer
// is full rows are dropped and counted (logged periodically) so the sensor hot
// path is never stalled by a slow database. The buffer is a fixed size.
type rowBatcher[T any] struct {
	name    string
	in      chan T
	maxRows int
	flush   func(context.Context, []T) error
	dropped atomic.Uint64
}

func newRowBatcher[T any](ctx context.Context, name string, maxRows int, flush func(context.Context, []T) error) *rowBatcher[T] {
	b := &rowBatcher[T]{
		name:    name,
		in:      make(chan T, maxRows*2),
		maxRows: maxRows,
		flush:   flush,
	}
	go b.loop(ctx)
	return b
}

func (b *rowBatcher[T]) add(v T) {
	select {
	case b.in <- v:
	default:
		b.dropped.Add(1)
	}
}

func (b *rowBatcher[T]) loop(ctx context.Context) {
	ticker := time.NewTicker(chFlushInterval)
	defer ticker.Stop()

	buf := make([]T, 0, b.maxRows)
	doFlush := func() {
		if len(buf) == 0 {
			return
		}
		if err := b.flush(ctx, buf); err != nil {
			log.Error().Err(err).Str("table", b.name).Int("rows", len(buf)).Msg("ClickHouse batch insert failed")
		}
		buf = buf[:0]
	}

	for {
		select {
		case <-ctx.Done():
			// Drain any buffered rows before exiting.
			for {
				select {
				case v := <-b.in:
					buf = append(buf, v)
					if len(buf) >= b.maxRows {
						doFlush()
					}
				default:
					doFlush()
					return
				}
			}
		case v := <-b.in:
			buf = append(buf, v)
			if len(buf) >= b.maxRows {
				doFlush()
			}
		case <-ticker.C:
			doFlush()
			if d := b.dropped.Swap(0); d > 0 {
				log.Warn().Uint64("dropped", d).Str("table", b.name).Msg("ClickHouse batcher dropped rows (buffer full)")
			}
		}
	}
}
