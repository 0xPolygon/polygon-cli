package database

import (
	"context"
	"encoding/hex"
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

// Default batching parameters. Write* calls only enqueue rows; a background
// goroutine flushes them in batches, since ClickHouse prefers large, infrequent
// inserts. Buffers are fixed-size and drop-on-full so a slow database can never
// stall the sensor or exhaust memory.
const (
	chFlushInterval = 1 * time.Second
	chFlushTimeout  = 30 * time.Second
	// chMaxFlushAttempts bounds retries of a failed batch insert. Each retry is
	// immediate with a fresh connection, recovering from stale-connection and
	// transient errors without delaying shutdown.
	chMaxFlushAttempts = 3
	chBlockBatch       = 5000
	chBlockBodyBatch   = 5000
	chBlockTxBatch     = 20000
	chBlockEventBatch  = 50000
	chTxBatch          = 20000
	chTxEventBatch     = 50000
	chPeerBatch        = 2000
)

// Event sources, so consumers can tell a hash announcement from a delivered
// header or body.
const (
	srcHashAnnounce  = "hash_announce"
	srcNewBlock      = "new_block"
	srcHeader        = "header"
	srcHeaderBackfil = "header_backfill"
	srcFullTx        = "full_tx"
)

// ClickHouse implements the Database interface backed by a ClickHouse cluster.
// The table definitions live in clickhouse_schema.sql, in the sensor-network-tools
// and polygon-infrastructure repos rather than here.
//
// The schema separates content-addressed facts from observations and this writer
// matches it: every row written to a fact table (blocks, block_bodies, block_txs,
// transactions) is complete and a pure function of its hash, so two sensors emit
// byte-identical rows and no write can partially overwrite another. Everything
// observational goes to the event streams.
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

	blocks      *rowBatcher[chBlock]
	blockBodies *rowBatcher[chBlockBody]
	blockTxs    *rowBatcher[chBlockTx]
	blockEvt    *rowBatcher[chBlockEvent]
	txs         *rowBatcher[chTx]
	txEvt       *rowBatcher[chTxEvent]
	peers       *rowBatcher[chPeerSnapshot]

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
// background batch flushers. Callers should defer Close to drain buffered rows
// on shutdown. If the connection cannot be established the returned Database
// no-ops all writes (mirroring the Datastore backend) so the sensor keeps running.
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

	// Writer-friendly defaults when the DSN omits them: LZ4 compression is a
	// network win on the wide blocks table, and the pool is sized for the
	// concurrent per-table flushers plus the occasional query.
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
// batcher's append closure maps a row to its column values; the batch/flush/error
// handling lives in newInsertBatcher.
func (c *ClickHouse) startBatchers(ctx context.Context) {
	c.blocks = newInsertBatcher(ctx, c, "blocks", chBlockBatch,
		"INSERT INTO blocks (number, hash, parent_hash, block_time, signer, coinbase, difficulty, gas_used, gas_limit, base_fee, uncle_hash, state_root, tx_root, receipt_root, logs_bloom, extra_data, mix_digest, nonce, withdrawals_root, blob_gas_used, excess_blob_gas, parent_beacon_root)",
		func(b driver.Batch, r chBlock) error {
			return b.Append(r.number, r.hash, r.parentHash, r.blockTime, r.signer, r.coinbase, r.difficulty, r.gasUsed, r.gasLimit, r.baseFee, r.uncleHash, r.stateRoot, r.txRoot, r.receiptRoot, r.logsBloom, r.extraData, r.mixDigest, r.nonce, r.withdrawalsRoot, r.blobGasUsed, r.excessBlobGas, r.parentBeaconRoot)
		})
	c.blockBodies = newInsertBatcher(ctx, c, "block_bodies", chBlockBodyBatch,
		"INSERT INTO block_bodies (hash, tx_count, uncle_count, uncles, size_bytes)",
		func(b driver.Batch, r chBlockBody) error {
			return b.Append(r.hash, r.txCount, r.uncleCount, r.uncles, r.sizeBytes)
		})
	c.blockTxs = newInsertBatcher(ctx, c, "block_txs", chBlockTxBatch,
		"INSERT INTO block_txs (block_hash, tx_index, tx_hash, seen_date)",
		func(b driver.Batch, r chBlockTx) error {
			return b.Append(r.blockHash, r.txIndex, r.txHash, r.seenDate)
		})
	c.blockEvt = newInsertBatcher(ctx, c, "block_events", chBlockEventBatch,
		"INSERT INTO block_events (block_number, block_hash, sensor_id, node_id, source, seen_at, total_difficulty)",
		func(b driver.Batch, r chBlockEvent) error {
			return b.Append(r.blockNumber, r.blockHash, c.sensorID, r.nodeID, r.source, r.seenAt, r.totalDifficulty)
		})
	c.txs = newInsertBatcher(ctx, c, "transactions", chTxBatch,
		"INSERT INTO transactions (hash, from_address, to_address, value, gas, gas_price, gas_fee_cap, gas_tip_cap, nonce, tx_type, chain_id, input_selector, input_size, access_list_size, blob_count, auth_list_size, seen_date)",
		func(b driver.Batch, r chTx) error {
			return b.Append(r.hash, r.from, r.to, r.value, r.gas, r.gasPrice, r.gasFeeCap, r.gasTipCap, r.nonce, r.txType, r.chainID, r.inputSelector, r.inputSize, r.accessListSize, r.blobCount, r.authListSize, r.seenDate)
		})
	c.txEvt = newInsertBatcher(ctx, c, "tx_events", chTxEventBatch,
		"INSERT INTO tx_events (tx_hash, sensor_id, node_id, source, seen_at)",
		func(b driver.Batch, r chTxEvent) error {
			return b.Append(r.txHash, c.sensorID, r.nodeID, r.source, r.seenAt)
		})
	c.peers = newInsertBatcher(ctx, c, "peer_snapshots", chPeerBatch,
		"INSERT INTO peer_snapshots (sensor_id, node_id, name, url, caps, seen_at)",
		func(b driver.Batch, r chPeerSnapshot) error {
			return b.Append(c.sensorID, r.nodeID, r.name, r.url, r.caps, r.seenAt)
		})
}

// newInsertBatcher wraps newRowBatcher with the common flush behaviour: prepare
// the INSERT, append rows, send, retrying transient failures. Only appendRow
// varies per table.
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

// chBlock is a header row. Every field is derived from the header itself, so any
// two sensors that see this hash produce an identical row. Nothing observational
// (which sensor, when, how it was learned) belongs here -- see chBlockEvent.
type chBlock struct {
	number           uint64
	hash             string
	parentHash       string
	blockTime        time.Time
	signer           string
	coinbase         string
	difficulty       uint64
	gasUsed          uint64
	gasLimit         uint64
	baseFee          *big.Int
	uncleHash        string
	stateRoot        string
	txRoot           string
	receiptRoot      string
	logsBloom        []byte
	extraData        []byte
	mixDigest        string
	nonce            uint64
	withdrawalsRoot  string
	blobGasUsed      uint64
	excessBlobGas    uint64
	parentBeaconRoot string
}

// chBlockBody holds the facts a header cannot carry. Written only when a body or
// a full block is actually delivered, which is what keeps the header path from
// having to invent a tx_count.
type chBlockBody struct {
	hash       string
	txCount    uint32
	uncleCount uint16
	uncles     []string
	sizeBytes  uint32
}

type chBlockTx struct {
	blockHash string
	txIndex   uint32
	txHash    string
	seenDate  time.Time
}

type chBlockEvent struct {
	blockNumber     uint64
	blockHash       string
	nodeID          string
	source          string
	seenAt          time.Time
	totalDifficulty *big.Int
}

type chTx struct {
	hash           string
	from           string
	to             string
	value          *big.Int
	gas            uint64
	gasPrice       *big.Int
	gasFeeCap      *big.Int
	gasTipCap      *big.Int
	nonce          uint64
	txType         uint8
	chainID        uint64
	inputSelector  string
	inputSize      uint32
	accessListSize uint16
	blobCount      uint8
	authListSize   uint8
	seenDate       time.Time
}

type chTxEvent struct {
	txHash string
	nodeID string
	source string
	seenAt time.Time
}

type chPeerSnapshot struct {
	nodeID string
	name   string
	url    string
	caps   []string
	seenAt time.Time
}

// --- Database interface ----------------------------------------------------

func (c *ClickHouse) WriteBlock(ctx context.Context, peer *enode.Node, block *types.Block, td *big.Int, tfs time.Time) {
	if c.conn == nil {
		return
	}
	// Announced total difficulty describes the announcement, not the block, so it
	// rides on the event.
	if c.shouldWriteBlockEvents && peer != nil {
		c.blockEvt.add(chBlockEvent{
			blockNumber:     block.NumberU64(),
			blockHash:       block.Hash().Hex(),
			nodeID:          peer.ID().String(),
			source:          srcNewBlock,
			seenAt:          tfs,
			totalDifficulty: td,
		})
	}
	if c.shouldWriteBlocks {
		c.blocks.add(newChBlock(block.Header()))
		c.writeBlockBody(block.Hash(), block.Transactions(), block.Uncles(), block.Size(), tfs)
	}
	if c.shouldWriteTransactions {
		c.writeTxs(block.Transactions(), tfs)
	}
}

func (c *ClickHouse) WriteBlockHeaders(ctx context.Context, headers []*types.Header, tfs time.Time, isParent bool) {
	if c.conn == nil || !c.shouldWriteBlocks {
		return
	}
	// A header row is complete on its own: the fields it cannot carry (tx/uncle
	// counts, size) live in block_bodies, so this path can never overwrite them.
	source := srcHeader
	if isParent {
		source = srcHeaderBackfil
	}
	for _, h := range headers {
		c.blocks.add(newChBlock(h))
		if c.shouldWriteBlockEvents {
			c.blockEvt.add(chBlockEvent{
				blockNumber:     h.Number.Uint64(),
				blockHash:       h.Hash().Hex(),
				source:          source,
				seenAt:          tfs,
				totalDifficulty: big.NewInt(0),
			})
		}
	}
}

func (c *ClickHouse) WriteBlockBody(ctx context.Context, body *eth.BlockBody, hash common.Hash, tfs time.Time) {
	if c.conn == nil {
		return
	}
	txs, err := body.Transactions.Items()
	if err != nil {
		log.Error().Err(err).Str("hash", hash.Hex()).Msg("Failed to decode transactions from block body")
		return
	}
	// The header row comes from the header path; record only the body facts (keyed
	// by hash, since the height is unknown here) and the transactions.
	if c.shouldWriteBlocks {
		uncles, err := body.Uncles.Items()
		if err != nil {
			log.Error().Err(err).Str("hash", hash.Hex()).Msg("Failed to decode uncles from block body")
			uncles = nil
		}
		c.writeBlockBody(hash, txs, uncles, 0, tfs)
	}
	if c.shouldWriteTransactions {
		c.writeTxs(txs, tfs)
	}
}

func (c *ClickHouse) WriteBlockEvents(ctx context.Context, peer *enode.Node, anns []BlockAnnouncement, tfs time.Time) {
	if c.conn == nil || peer == nil {
		return
	}
	nodeID := peer.ID().String()
	for _, ann := range anns {
		c.blockEvt.add(chBlockEvent{
			blockNumber:     ann.Number,
			blockHash:       ann.Hash.Hex(),
			nodeID:          nodeID,
			source:          srcHashAnnounce,
			seenAt:          tfs,
			totalDifficulty: big.NewInt(0),
		})
	}
}

// WriteBlockHashFirstSeen is a no-op: earliest first-seen is derived from
// block_events by the block_events_first materialized view.
func (c *ClickHouse) WriteBlockHashFirstSeen(ctx context.Context, peer *enode.Node, hash common.Hash, tfsh time.Time) {
}

func (c *ClickHouse) WriteTransactionEvents(ctx context.Context, peer *enode.Node, hashes []common.Hash, tfs time.Time) {
	if c.conn == nil || peer == nil {
		return
	}
	nodeID := peer.ID().String()
	for _, hash := range hashes {
		c.txEvt.add(chTxEvent{
			txHash: hash.Hex(),
			nodeID: nodeID,
			source: srcHashAnnounce,
			seenAt: tfs,
		})
	}
}

func (c *ClickHouse) WriteTransactions(ctx context.Context, peer *enode.Node, txs []*types.Transaction, tfs time.Time) {
	if c.conn == nil {
		return
	}
	// A delivered body is a distinct event from a hash announcement, but recorded
	// only under the full per-peer stream flag: in first-event-only mode it would
	// multiply rows on the largest table without adding a new first-event.
	if c.shouldWriteTransactionEvents && peer != nil {
		nodeID := peer.ID().String()
		for _, tx := range txs {
			c.txEvt.add(chTxEvent{
				txHash: tx.Hash().Hex(),
				nodeID: nodeID,
				source: srcFullTx,
				seenAt: tfs,
			})
		}
	}
	if !c.shouldWriteTransactions {
		return
	}
	c.writeTxs(txs, tfs)
}

func (c *ClickHouse) WritePeers(ctx context.Context, peers []*p2p.Peer, tls time.Time) {
	if c.conn == nil || !c.shouldWritePeers {
		return
	}
	// node_id matches what the event streams record, so peers and events are
	// joinable. The enode URL is a column, not the key.
	for _, peer := range peers {
		c.peers.add(chPeerSnapshot{
			nodeID: peer.ID().String(),
			name:   peer.Fullname(),
			url:    peer.Node().URLv4(),
			caps:   peer.Info().Caps,
			seenAt: tls,
		})
	}
}

// HasBlock reports whether the block already exists. Called once per new block
// (not per event), so an indexed point lookup is cheap. Without a connection it
// reports true so the sensor never attempts a backfill it could not persist.
func (c *ClickHouse) HasBlock(ctx context.Context, hash common.Hash) bool {
	if c.conn == nil {
		return true
	}
	var exists uint8
	err := c.conn.QueryRow(ctx, "SELECT 1 FROM blocks WHERE hash = ? LIMIT 1", hash.Hex()).Scan(&exists)
	return err == nil && exists == 1
}

// NodeList returns the most recently seen peers' enode URLs, from the narrow
// peers_current rollup rather than grouping over the event firehose.
func (c *ClickHouse) NodeList(ctx context.Context, limit int) ([]string, error) {
	if c.conn == nil {
		return []string{}, nil
	}
	rows, err := c.conn.Query(ctx,
		"SELECT url FROM peers_current FINAL WHERE url != '' ORDER BY last_seen DESC LIMIT ?", limit)
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

// newChBlock maps a header to a blocks-table row. Takes nothing but the header, so
// everything it writes is a pure function of it -- which is what makes duplicate
// rows for a hash byte-identical.
func newChBlock(h *types.Header) chBlock {
	baseFee := new(big.Int)
	if h.BaseFee != nil {
		baseFee.Set(h.BaseFee)
	}
	// Recover the block signer from the header seal so signer-based analytics
	// don't have to ecrecover on every query. Left empty when it can't be recovered.
	var signer string
	if sig, err := util.Ecrecover(h); err == nil {
		signer = common.BytesToAddress(sig).Hex()
	}
	b := chBlock{
		number:      h.Number.Uint64(),
		hash:        h.Hash().Hex(),
		parentHash:  h.ParentHash.Hex(),
		blockTime:   time.Unix(int64(h.Time), 0).UTC(),
		signer:      signer,
		coinbase:    h.Coinbase.Hex(),
		difficulty:  h.Difficulty.Uint64(),
		gasUsed:     h.GasUsed,
		gasLimit:    h.GasLimit,
		baseFee:     baseFee,
		uncleHash:   h.UncleHash.Hex(),
		stateRoot:   h.Root.Hex(),
		txRoot:      h.TxHash.Hex(),
		receiptRoot: h.ReceiptHash.Hex(),
		logsBloom:   h.Bloom.Bytes(),
		extraData:   h.Extra,
		mixDigest:   h.MixDigest.Hex(),
		nonce:       h.Nonce.Uint64(),
	}
	// Post-Shanghai/Cancun fields are pointers, absent on older chains; the columns
	// default to '' / 0.
	if h.WithdrawalsHash != nil {
		b.withdrawalsRoot = h.WithdrawalsHash.Hex()
	}
	if h.BlobGasUsed != nil {
		b.blobGasUsed = *h.BlobGasUsed
	}
	if h.ExcessBlobGas != nil {
		b.excessBlobGas = *h.ExcessBlobGas
	}
	if h.ParentBeaconRoot != nil {
		b.parentBeaconRoot = h.ParentBeaconRoot.Hex()
	}
	return b
}

// writeBlockBody records the body facts and the ordered block -> tx mapping. size
// is 0 when the body arrived alone, since the encoded size needs the whole block.
func (c *ClickHouse) writeBlockBody(hash common.Hash, txs []*types.Transaction, uncles []*types.Header, size uint64, tfs time.Time) {
	uncleHashes := make([]string, 0, len(uncles))
	for _, u := range uncles {
		uncleHashes = append(uncleHashes, u.Hash().Hex())
	}
	c.blockBodies.add(chBlockBody{
		hash:       hash.Hex(),
		txCount:    uint32(len(txs)),
		uncleCount: uint16(len(uncles)),
		uncles:     uncleHashes,
		sizeBytes:  uint32(size),
	})

	blockHash := hash.Hex()
	seenDate := tfs.UTC().Truncate(24 * time.Hour)
	for i, tx := range txs {
		c.blockTxs.add(chBlockTx{
			blockHash: blockHash,
			txIndex:   uint32(i),
			txHash:    tx.Hash().Hex(),
			seenDate:  seenDate,
		})
	}
}

func (c *ClickHouse) writeTxs(txs []*types.Transaction, tfs time.Time) {
	seenDate := tfs.UTC().Truncate(24 * time.Hour)
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
		// Selector and size but not the calldata: enough for contract-interaction
		// analysis at negligible cost.
		var selector string
		data := tx.Data()
		if len(data) >= 4 {
			selector = "0x" + hex.EncodeToString(data[:4])
		}
		txChainID := uint64(0)
		if id := tx.ChainId(); id != nil && id.IsUint64() {
			txChainID = id.Uint64()
		}
		c.txs.add(chTx{
			hash:           tx.Hash().Hex(),
			from:           from,
			to:             to,
			value:          new(big.Int).Set(tx.Value()),
			gas:            tx.Gas(),
			gasPrice:       new(big.Int).Set(tx.GasPrice()),
			gasFeeCap:      new(big.Int).Set(tx.GasFeeCap()),
			gasTipCap:      new(big.Int).Set(tx.GasTipCap()),
			nonce:          tx.Nonce(),
			txType:         tx.Type(),
			chainID:        txChainID,
			inputSelector:  selector,
			inputSize:      uint32(len(data)),
			accessListSize: uint16(len(tx.AccessList())),
			blobCount:      uint8(len(tx.BlobHashes())),
			authListSize:   uint8(len(tx.SetCodeAuthorizations())),
			seenDate:       seenDate,
		})
	}
}

// --- batching --------------------------------------------------------------

// rowBatcher buffers rows and flushes them in bulk when the buffer reaches
// maxRows or on a fixed interval. add is non-blocking: rows are dropped and
// counted when the fixed-size buffer is full, so a slow database never stalls
// the sensor hot path.
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
