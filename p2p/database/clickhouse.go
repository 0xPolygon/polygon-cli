package database

import (
	"context"
	"fmt"
	"math/big"
	"sync"
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
	chFlushInterval = 1 * time.Second
	chFlushTimeout  = 30 * time.Second
	// chMaxFlushAttempts bounds retries of a failed batch insert before the
	// batch is dropped. Retries are immediate: a fresh connection is acquired
	// each attempt, which recovers from stale-connection and transient errors
	// without delaying shutdown.
	chMaxFlushAttempts = 3
	chBlockBatch       = 5000
	chBlockEventBatch  = 50000
	chTxBatch          = 20000
	chTxEventBatch     = 50000
	chPeerBatch        = 2000
)

// ClickHouse implements the Database interface backed by a ClickHouse cluster.
// The table definitions this writer targets (and the block_first_seen
// materialized view) live in the sensor-network-tools repo
// (clickhouse_schema.sql), not this repo.
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

	// cancel stops the batcher goroutines; wg tracks them so Close can wait for
	// their final drain flush before the connection is closed.
	cancel context.CancelFunc
	wg     sync.WaitGroup
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
}

// NewClickHouse connects to ClickHouse, verifies connectivity, and starts the
// background batch flushers. The flusher goroutines run until either the
// provided context is cancelled or Close is called, at which point they flush
// any buffered rows and exit. Callers should defer Close to guarantee buffered
// rows are drained before shutdown. If the connection cannot be established the
// returned Database no-ops all writes (mirroring the Datastore backend) so the
// sensor keeps running.
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

	conn, err := connectClickHouse(ctx, opts.DSN)
	if err != nil {
		log.Error().Err(err).Msg("Could not initialize ClickHouse connection")
		return c
	}
	c.conn = conn

	// Derive a cancellable context so Close can stop the batchers independently
	// of the parent context.
	bctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	c.startBatchers(bctx)

	return c
}

// Close stops the batcher goroutines, waits for their final drain flush to
// complete, and closes the connection. It is safe to call on a no-op instance
// (connection never established).
func (c *ClickHouse) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// connectClickHouse parses the DSN, opens a connection, and verifies
// connectivity with a ping.
func connectClickHouse(ctx context.Context, dsn string) (driver.Conn, error) {
	chOpts, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("could not parse ClickHouse DSN: %w", err)
	}

	// Apply writer-friendly defaults when the DSN doesn't set them. LZ4
	// compression is a large network win on the wide blocks table, and the pool
	// is sized for the concurrent per-table flushers plus the occasional query.
	if chOpts.Compression == nil {
		chOpts.Compression = &clickhouse.Compression{Method: clickhouse.CompressionLZ4}
	}
	if chOpts.MaxIdleConns == 0 {
		chOpts.MaxIdleConns = 10
	}
	if chOpts.MaxOpenConns == 0 {
		chOpts.MaxOpenConns = 20
	}

	conn, err := clickhouse.Open(chOpts)
	if err != nil {
		return nil, fmt.Errorf("could not connect to ClickHouse: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("could not ping ClickHouse: %w", err)
	}
	return conn, nil
}

// startBatchers creates the background row batchers, one per target table. Each
// batcher's append closure only needs to map a row to its column values; the
// surrounding batch/flush/error handling lives in newInsertBatcher.
func (c *ClickHouse) startBatchers(ctx context.Context) {
	c.blocks = newInsertBatcher(ctx, c, "blocks", chBlockBatch,
		"INSERT INTO blocks (hash, number, parent_hash, block_time, coinbase, signer, difficulty, total_difficulty, gas_used, gas_limit, base_fee, tx_count, uncle_count, uncle_hash, state_root, tx_root, receipt_root, logs_bloom, extra_data, mix_digest, nonce, sensor_id, ingested_at, is_parent)",
		func(b driver.Batch, r chBlock) error {
			return b.Append(r.hash, r.number, r.parentHash, r.blockTime, r.coinbase, r.signer, r.difficulty, r.totalDifficulty, r.gasUsed, r.gasLimit, r.baseFee, r.txCount, r.uncleCount, r.uncleHash, r.stateRoot, r.txRoot, r.receiptRoot, r.logsBloom, r.extraData, r.mixDigest, r.nonce, c.sensorID, r.ingestedAt, r.isParent)
		})
	c.blockEvt = newInsertBatcher(ctx, c, "block_events", chBlockEventBatch,
		"INSERT INTO block_events (block_hash, sensor_id, peer_id, seen_at)",
		func(b driver.Batch, r chEvent) error {
			return b.Append(r.hash, c.sensorID, r.peerID, r.seenAt)
		})
	c.txs = newInsertBatcher(ctx, c, "transactions", chTxBatch,
		"INSERT INTO transactions (hash, from_address, to_address, value, gas, gas_price, gas_fee_cap, gas_tip_cap, nonce, tx_type, sensor_id, ingested_at)",
		func(b driver.Batch, r chTx) error {
			return b.Append(r.hash, r.from, r.to, r.value, r.gas, r.gasPrice, r.gasFeeCap, r.gasTipCap, r.nonce, r.txType, c.sensorID, r.ingestedAt)
		})
	c.txEvt = newInsertBatcher(ctx, c, "transaction_events", chTxEventBatch,
		"INSERT INTO transaction_events (tx_hash, sensor_id, peer_id, seen_at)",
		func(b driver.Batch, r chEvent) error {
			return b.Append(r.hash, c.sensorID, r.peerID, r.seenAt)
		})
	c.peers = newInsertBatcher(ctx, c, "peers", chPeerBatch,
		"INSERT INTO peers (peer_id, name, url, caps, last_seen_by, time_last_seen)",
		func(b driver.Batch, r chPeer) error {
			return b.Append(r.peerID, r.name, r.url, r.caps, c.sensorID, r.timeLastSeen)
		})
}

// newInsertBatcher wraps newRowBatcher with the common flush behaviour: prepare
// the INSERT, append each row via appendRow, and send, retrying transient
// failures. Only appendRow varies per table.
func newInsertBatcher[T any](ctx context.Context, c *ClickHouse, name string, maxRows int, query string, appendRow func(driver.Batch, T) error) *rowBatcher[T] {
	return newRowBatcher(ctx, &c.wg, name, maxRows, func(rows []T) error {
		var err error
		for attempt := 1; attempt <= chMaxFlushAttempts; attempt++ {
			if err = flushBatch(c.conn, query, rows, appendRow); err == nil {
				return nil
			}
			if attempt < chMaxFlushAttempts {
				log.Warn().Err(err).Str("table", name).Int("attempt", attempt).Int("rows", len(rows)).
					Msg("ClickHouse batch insert failed; retrying")
			}
		}
		return err
	})
}

// flushBatch prepares, fills, and sends a single INSERT. It runs on a detached,
// time-bounded context so a flush triggered during shutdown (parent context
// already cancelled) still completes.
func flushBatch[T any](conn driver.Conn, query string, rows []T, appendRow func(driver.Batch, T) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), chFlushTimeout)
	defer cancel()

	b, err := conn.PrepareBatch(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare batch: %w", err)
	}
	for _, r := range rows {
		if err := appendRow(b, r); err != nil {
			return fmt.Errorf("append row: %w", err)
		}
	}
	if err := b.Send(); err != nil {
		return fmt.Errorf("send batch: %w", err)
	}
	return nil
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
	uncleHash       string
	stateRoot       string
	txRoot          string
	receiptRoot     string
	logsBloom       []byte
	extraData       []byte
	mixDigest       string
	nonce           uint64
	ingestedAt      time.Time
	isParent        bool
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
		c.blocks.add(c.newBlock(block.Header(), td, tfs, len(block.Transactions()), len(block.Uncles()), false))
	}
	if c.shouldWriteTransactions {
		c.writeTxs(block.Transactions(), tfs)
	}
}

func (c *ClickHouse) WriteBlockHeaders(ctx context.Context, headers []*types.Header, tfs time.Time, isParent bool) {
	if c.conn == nil || !c.shouldWriteBlocks {
		return
	}
	// A header alone carries no transaction/uncle counts, so they are written as
	// 0; the full-block (NewBlock) path writes a separate row with the real
	// counts. isParent marks headers fetched as ancestors during backfill so
	// parent-vs-live analysis can distinguish them.
	for _, h := range headers {
		c.blocks.add(c.newBlock(h, big.NewInt(0), tfs, 0, 0, isParent))
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
	// Record a transaction event when either the full event stream or the
	// first-seen-only mode is enabled. processTransactions dedups upstream, so
	// this fires once per first-seen tx in both cases (mirroring the
	// first-block-event behavior for blocks).
	if peer != nil && (c.shouldWriteTransactionEvents || c.shouldWriteFirstTransactionEvent) {
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
		return nil, fmt.Errorf("query node list: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close ClickHouse rows")
		}
	}()

	nodelist := []string{}
	for rows.Next() {
		var peerID string
		if err := rows.Scan(&peerID); err != nil {
			log.Error().Err(err).Msg("Failed to scan peer_id")
			continue
		}
		nodelist = append(nodelist, peerID)
	}
	if err := rows.Err(); err != nil {
		return nodelist, fmt.Errorf("iterate node list: %w", err)
	}
	return nodelist, nil
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

func (c *ClickHouse) newBlock(h *types.Header, td *big.Int, tfs time.Time, txCount, uncleCount int, isParent bool) chBlock {
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
		uncleHash:       h.UncleHash.Hex(),
		stateRoot:       h.Root.Hex(),
		txRoot:          h.TxHash.Hex(),
		receiptRoot:     h.ReceiptHash.Hex(),
		logsBloom:       h.Bloom.Bytes(),
		extraData:       h.Extra,
		mixDigest:       h.MixDigest.Hex(),
		nonce:           h.Nonce.Uint64(),
		ingestedAt:      tfs,
		isParent:        isParent,
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
	flush   func([]T) error
	dropped atomic.Uint64
}

func newRowBatcher[T any](ctx context.Context, wg *sync.WaitGroup, name string, maxRows int, flush func([]T) error) *rowBatcher[T] {
	b := &rowBatcher[T]{
		name:    name,
		in:      make(chan T, maxRows*2),
		maxRows: maxRows,
		flush:   flush,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		b.loop(ctx)
	}()
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
		if err := b.flush(buf); err != nil {
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
