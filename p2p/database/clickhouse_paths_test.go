package database

import (
	"context"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"
)

// Integration tests for the write paths not covered elsewhere. Like the rest of
// this package's ClickHouse tests they run only with POLYCLI_TEST_CLICKHOUSE_DSN
// set, and use their own reserved heights.
const (
	heightBody        = 900_000_090
	heightBadUncles   = 900_000_100
	heightBodyFlags   = 900_000_110
	heightAnnounce    = 900_000_120
	heightHasBlock    = 900_000_130
	heightTxFlagsOnly = 900_000_140
)

// countRows returns the number of rows a query matches.
func countRows(t *testing.T, conn driver.Conn, query string, args ...any) uint64 {
	t.Helper()
	var n uint64
	if err := conn.QueryRow(context.Background(), query, args...).Scan(&n); err != nil {
		t.Fatalf("count via %q: %v", query, err)
	}
	return n
}

// cleanupTestPeers removes peer rows written under the test sensor ids.
func cleanupTestPeers(t *testing.T, conn driver.Conn) {
	t.Helper()
	t.Cleanup(func() {
		for _, q := range []string{
			"ALTER TABLE peers DELETE WHERE sensor_id LIKE 'test-sensor%' SETTINGS mutations_sync = 1",
			"ALTER TABLE peers_current DELETE WHERE sensor_id LIKE 'test-sensor%' SETTINGS mutations_sync = 1",
		} {
			if err := conn.Exec(context.Background(), q); err != nil {
				t.Logf("peer cleanup: %v", err)
			}
		}
	})
}

// testBody returns a block at height with two transactions and one uncle, its
// raw body as a BlockBodies response carries it, and the decoded transactions.
func testBody(t *testing.T, height int64) (*types.Block, *eth.BlockBody, []*types.Transaction) {
	t.Helper()
	now := time.Now().UTC()
	header, _ := signedHeader(t, height, now)
	uncle, _ := signedHeader(t, height-1, now)
	txs := []*types.Transaction{
		types.NewTx(&types.LegacyTx{Nonce: uint64(now.UnixNano()), GasPrice: big.NewInt(1), Gas: 21_000, Value: big.NewInt(1)}),
		types.NewTx(&types.LegacyTx{Nonce: uint64(now.UnixNano()) + 1, GasPrice: big.NewInt(1), Gas: 21_000, Value: big.NewInt(1)}),
	}
	block := types.NewBlockWithHeader(header).WithBody(types.Body{Transactions: txs, Uncles: []*types.Header{uncle}})
	txList, err := rlp.EncodeToRawList(txs)
	if err != nil {
		t.Fatal(err)
	}
	uncleList, err := rlp.EncodeToRawList([]*types.Header{uncle})
	if err != nil {
		t.Fatal(err)
	}
	return block, &eth.BlockBody{Transactions: txList, Uncles: uncleList}, txs
}

// TestClickHouseWritesBlockBody covers the BlockBodies path: the body event on
// arrival, then the body facts once the body has been checked.
func TestClickHouseWritesBlockBody(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestClickHouse(t, dsn, ctx)

	block, body, txs := testBody(t, heightBody)
	ann := BlockAnnouncement{Hash: block.Hash(), Number: heightBody}
	now := time.Now().UTC()
	// The header gives the cleanup a blocks row to find the body rows by.
	db.WriteBlockHeaders(ctx, []*types.Header{block.Header()}, now, false)
	db.WriteBlockBodyArrival(ctx, ann, now)
	db.WriteBlockBody(ctx, body, txsOf(t, block), ann, now)
	if err := db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	conn := verifyConn(t, dsn)
	cleanupTestHeights(t, conn, heightBody)
	hash := block.Hash().Hex()

	var txCount uint32
	var uncleCount uint16
	if err := conn.QueryRow(context.Background(),
		"SELECT tx_count, uncle_count FROM block_bodies FINAL WHERE hash = ?", hash).Scan(&txCount, &uncleCount); err != nil {
		t.Fatalf("scan block_bodies: %v", err)
	}
	if txCount != 2 || uncleCount != 1 {
		t.Fatalf("block_bodies tx_count=%d uncle_count=%d, want 2 and 1", txCount, uncleCount)
	}
	if n := countRows(t, conn, "SELECT count() FROM block_txs FINAL WHERE block_hash = ?", hash); n != 2 {
		t.Fatalf("block_txs rows = %d, want 2", n)
	}
	for _, tx := range txs {
		checkCount(t, conn, "SELECT count() FROM transactions WHERE hash = ?", tx.Hash().Hex())
	}
	if n := countRows(t, conn, "SELECT count() FROM block_events WHERE block_hash = ? AND source = 'body'", hash); n != 1 {
		t.Fatalf("body events = %d, want 1", n)
	}
}

// TestClickHouseBlockBodyWithBadUncles covers a body whose uncles do not decode:
// the block_bodies row, which carries uncle facts, is skipped, but block_txs and
// transactions, which depend only on the transactions, are still written.
func TestClickHouseBlockBodyWithBadUncles(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestClickHouse(t, dsn, ctx)

	block, body, txs := testBody(t, heightBadUncles)
	if err := body.Uncles.AppendRaw([]byte{0x01}); err != nil {
		t.Fatal(err)
	}
	ann := BlockAnnouncement{Hash: block.Hash(), Number: heightBadUncles}
	now := time.Now().UTC()
	db.WriteBlockHeaders(ctx, []*types.Header{block.Header()}, now, false)
	db.WriteBlockBody(ctx, body, txsOf(t, block), ann, now)
	if err := db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	conn := verifyConn(t, dsn)
	cleanupTestHeights(t, conn, heightBadUncles)
	hash := block.Hash().Hex()
	if n := countRows(t, conn, "SELECT count() FROM block_bodies WHERE hash = ?", hash); n != 0 {
		t.Fatalf("block_bodies rows = %d, want 0 for undecodable uncles", n)
	}
	if n := countRows(t, conn, "SELECT count() FROM block_txs FINAL WHERE block_hash = ?", hash); n != 2 {
		t.Fatalf("block_txs rows = %d, want 2", n)
	}
	checkCount(t, conn, "SELECT count() FROM transactions WHERE hash = ?", txs[0].Hash().Hex())
}

// TestClickHouseBlockBodyRespectsWriteSettings checks that the body path writes
// nothing to a table whose setting is off: with only transactions on, there is
// no body event and no block_bodies or block_txs row.
func TestClickHouseBlockBodyRespectsWriteSettings(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := NewClickHouse(ctx, ClickHouseOptions{
		DSN:                     dsn,
		SensorID:                "test-sensor-flags",
		ChainID:                 137,
		MaxConcurrency:          1,
		ShouldWriteTransactions: true,
	})

	block, body, txs := testBody(t, heightBodyFlags)
	ann := BlockAnnouncement{Hash: block.Hash(), Number: heightBodyFlags}
	now := time.Now().UTC()
	db.WriteBlockBodyArrival(ctx, ann, now)
	db.WriteBlockBody(ctx, body, txsOf(t, block), ann, now)
	if err := db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	conn := verifyConn(t, dsn)
	cleanupTestHeights(t, conn, heightBodyFlags)
	t.Cleanup(func() {
		for _, tx := range txs {
			if err := conn.Exec(context.Background(),
				"ALTER TABLE transactions DELETE WHERE hash = ? SETTINGS mutations_sync = 1", tx.Hash().Hex()); err != nil {
				t.Logf("transactions cleanup: %v", err)
			}
		}
	})
	hash := block.Hash().Hex()
	for table, query := range map[string]string{
		"block_events": "SELECT count() FROM block_events WHERE block_hash = ?",
		"block_bodies": "SELECT count() FROM block_bodies WHERE hash = ?",
		"block_txs":    "SELECT count() FROM block_txs WHERE block_hash = ?",
	} {
		if n := countRows(t, conn, query, hash); n != 0 {
			t.Errorf("%s rows = %d, want 0 with its setting off", table, n)
		}
	}
	for _, tx := range txs {
		checkCount(t, conn, "SELECT count() FROM transactions WHERE hash = ?", tx.Hash().Hex())
	}
}

// TestClickHouseWritesAnnouncements covers the per-peer announcement streams:
// block hash announcements and transaction hash announcements are recorded with
// the announcing peer's node id, and an announcement with no peer is not.
func TestClickHouseWritesAnnouncements(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestClickHouse(t, dsn, ctx)

	now := time.Now().UTC()
	peer := testPeer(t)
	blockHash := common.BigToHash(big.NewInt(now.UnixNano()))
	txHash := common.BigToHash(big.NewInt(now.UnixNano() + 1))
	ann := BlockAnnouncement{Hash: blockHash, Number: heightAnnounce}

	db.WriteBlockEvents(ctx, peer, []BlockAnnouncement{ann}, now)
	db.WriteBlockEvents(ctx, nil, []BlockAnnouncement{ann}, now)
	db.WriteTransactionEvents(ctx, peer, []common.Hash{txHash}, now)
	db.WriteTransactionEvents(ctx, nil, []common.Hash{txHash}, now)
	db.WriteBlockHashFirstSeen(ctx, peer, blockHash, now)
	if err := db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	conn := verifyConn(t, dsn)
	cleanupTestHeights(t, conn, heightAnnounce)
	t.Cleanup(func() {
		if err := conn.Exec(context.Background(),
			"ALTER TABLE tx_events DELETE WHERE tx_hash = ? SETTINGS mutations_sync = 1", txHash.Hex()); err != nil {
			t.Logf("tx_events cleanup: %v", err)
		}
	})
	nodeID := peer.ID().String()
	if n := countRows(t, conn,
		"SELECT count() FROM block_events WHERE block_hash = ? AND source = 'hash_announce' AND node_id = ?",
		blockHash.Hex(), nodeID); n != 1 {
		t.Fatalf("block hash_announce events = %d, want 1 from the peer", n)
	}
	if n := countRows(t, conn,
		"SELECT count() FROM tx_events WHERE tx_hash = ? AND source = 'hash_announce' AND node_id = ?",
		txHash.Hex(), nodeID); n != 1 {
		t.Fatalf("tx hash_announce events = %d, want 1 from the peer", n)
	}
}

// TestClickHouseTransactionsRespectWriteSettings checks WriteTransactions with
// events off: the transactions row is written, the full_tx event is not.
func TestClickHouseTransactionsRespectWriteSettings(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := NewClickHouse(ctx, ClickHouseOptions{
		DSN:                     dsn,
		SensorID:                "test-sensor-txonly",
		ChainID:                 137,
		MaxConcurrency:          1,
		ShouldWriteTransactions: true,
	})

	_, _, txs := testBody(t, heightTxFlagsOnly)
	db.WriteTransactions(ctx, testPeer(t), txs, time.Now().UTC())
	if err := db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	conn := verifyConn(t, dsn)
	t.Cleanup(func() {
		for _, tx := range txs {
			if err := conn.Exec(context.Background(),
				"ALTER TABLE transactions DELETE WHERE hash = ? SETTINGS mutations_sync = 1", tx.Hash().Hex()); err != nil {
				t.Logf("transactions cleanup: %v", err)
			}
		}
	})
	for _, tx := range txs {
		checkCount(t, conn, "SELECT count() FROM transactions WHERE hash = ?", tx.Hash().Hex())
		if n := countRows(t, conn, "SELECT count() FROM tx_events WHERE tx_hash = ?", tx.Hash().Hex()); n != 0 {
			t.Fatalf("tx_events rows = %d, want 0 with transaction events off", n)
		}
	}
}

// TestClickHouseWritesPeersAndNodeList covers the peer snapshot and reading it
// back: peers are written with their node id, name and enode URL, and both
// NodeList and ClickHouseNodeList return their URLs.
func TestClickHouseWritesPeersAndNodeList(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestClickHouse(t, dsn, ctx)

	nodes := []*p2p.Peer{p2p.NewPeer(testPeer(t).ID(), "peer-a", []p2p.Cap{{Name: "eth", Version: 68}})}
	db.WritePeers(ctx, nodes, time.Now().UTC())
	if err := db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	conn := verifyConn(t, dsn)
	cleanupTestPeers(t, conn)
	nodeID := nodes[0].ID().String()
	var name, url string
	var caps []string
	if err := conn.QueryRow(context.Background(),
		"SELECT name, url, caps FROM peers_current FINAL WHERE sensor_id = 'test-sensor' AND node_id = ?", nodeID).Scan(&name, &url, &caps); err != nil {
		t.Fatalf("scan peers_current: %v", err)
	}
	if name != "peer-a" || url != nodes[0].Node().URLv4() || len(caps) != 1 || caps[0] != "eth/68" {
		t.Fatalf("peers_current row = %q %q %v", name, url, caps)
	}

	list, err := ClickHouseNodeList(context.Background(), dsn, 1000)
	if err != nil {
		t.Fatalf("ClickHouseNodeList: %v", err)
	}
	found := false
	for _, u := range list {
		found = found || u == url
	}
	if !found {
		t.Fatalf("ClickHouseNodeList did not return the written peer %s", url)
	}
}

// TestClickHouseHasBlock covers HasBlock's query: true for a stored block,
// false for an unknown one, and true when the query fails, so a failed check
// does not trigger a backfill.
func TestClickHouseHasBlock(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	header, _ := signedHeader(t, heightHasBlock, time.Now().UTC())
	writer := newTestClickHouse(t, dsn, ctx)
	writer.WriteBlockHeaders(ctx, []*types.Header{header}, time.Now().UTC(), false)
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	conn := verifyConn(t, dsn)
	cleanupTestHeights(t, conn, heightHasBlock)

	// A fresh instance, so the answer comes from the database and not from the
	// writer's record of pending blocks.
	reader := newTestClickHouse(t, dsn, ctx)
	defer func() { _ = reader.Close() }()
	if !reader.HasBlock(ctx, header.Hash()) {
		t.Fatal("HasBlock = false for a stored block")
	}
	if reader.HasBlock(ctx, common.Hash{0xde, 0xad}) {
		t.Fatal("HasBlock = true for an unknown block")
	}
	cancelled, stop := context.WithCancel(ctx)
	stop()
	if !reader.HasBlock(cancelled, common.Hash{0xde, 0xad}) {
		t.Fatal("HasBlock = false when the query failed; it should assume the block is stored")
	}
}

// TestClickHouseRecordsSender checks that a signed transaction's sender and
// recipient are stored, lowercase.
func TestClickHouseRecordsSender(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestClickHouse(t, dsn, ctx)

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	to := common.HexToAddress("0x00000000000000000000000000000000000000Ab")
	tx, err := types.SignNewTx(key, types.LatestSignerForChainID(big.NewInt(137)), &types.DynamicFeeTx{
		ChainID:   big.NewInt(137),
		Nonce:     uint64(time.Now().UnixNano()),
		GasTipCap: big.NewInt(1),
		GasFeeCap: big.NewInt(2),
		Gas:       21_000,
		To:        &to,
		Value:     big.NewInt(1),
	})
	if err != nil {
		t.Fatal(err)
	}
	// A pre-EIP-155 transaction carries no chain id, so the sensor's is used.
	legacy, err := types.SignNewTx(key, types.HomesteadSigner{}, &types.LegacyTx{
		Nonce:    uint64(time.Now().UnixNano()),
		GasPrice: big.NewInt(1),
		Gas:      21_000,
		To:       &to,
		Value:    big.NewInt(1),
	})
	if err != nil {
		t.Fatal(err)
	}
	db.WriteTransactions(ctx, testPeer(t), []*types.Transaction{tx, legacy}, time.Now().UTC())
	if err := db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	conn := verifyConn(t, dsn)
	cleanupTestHeights(t, conn)
	wantFrom := strings.ToLower(crypto.PubkeyToAddress(key.PublicKey).Hex())
	wantTo := strings.ToLower(to.Hex())
	for _, tx := range []*types.Transaction{tx, legacy} {
		var from, gotTo string
		if err := conn.QueryRow(context.Background(),
			"SELECT from_address, to_address FROM transactions WHERE hash = ? LIMIT 1", tx.Hash().Hex()).Scan(&from, &gotTo); err != nil {
			t.Fatalf("scan transaction type %d: %v", tx.Type(), err)
		}
		if from != wantFrom || gotTo != wantTo {
			t.Fatalf("type %d transaction from=%s to=%s, want %s and %s", tx.Type(), from, gotTo, wantFrom, wantTo)
		}
	}
}

// TestClickHouseTransactionEventsOnly checks WriteTransactions with only
// transaction events on: the full_tx event is written, the transactions row
// is not.
func TestClickHouseTransactionEventsOnly(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := NewClickHouse(ctx, ClickHouseOptions{
		DSN:                          dsn,
		SensorID:                     "test-sensor-txevents",
		ChainID:                      137,
		MaxConcurrency:               1,
		ShouldWriteTransactionEvents: true,
	})

	_, _, txs := testBody(t, heightTxFlagsOnly)
	peer := testPeer(t)
	db.WriteTransactions(ctx, peer, txs, time.Now().UTC())
	if err := db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	conn := verifyConn(t, dsn)
	cleanupTestHeights(t, conn)
	for _, tx := range txs {
		if n := countRows(t, conn,
			"SELECT count() FROM tx_events WHERE tx_hash = ? AND source = 'full_tx' AND node_id = ?",
			tx.Hash().Hex(), peer.ID().String()); n != 1 {
			t.Fatalf("full_tx events = %d, want 1", n)
		}
		if n := countRows(t, conn, "SELECT count() FROM transactions WHERE hash = ?", tx.Hash().Hex()); n != 0 {
			t.Fatalf("transactions rows = %d, want 0 with transactions off", n)
		}
	}
}

// TestClickHouseFlushErrors covers a failed flush: an INSERT that cannot be
// prepared is retried and then reported, and a row that does not fit the
// INSERT fails on append.
func TestClickHouseFlushErrors(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestClickHouse(t, dsn, ctx)
	defer func() { _ = db.Close() }()
	ch := db.(*ClickHouse)

	// The batcher joins ch's WaitGroup, so it is stopped before ch is closed.
	bctx, stop := context.WithCancel(ctx)
	defer stop()
	attempts := 0
	b := newInsertBatcher(bctx, ch, "missing", 1,
		"INSERT INTO no_such_table (x)", func(batch driver.Batch, r int) error {
			attempts++
			return batch.Append(r)
		})
	if err := b.flush([]int{1}); err == nil || !strings.Contains(err.Error(), "prepare batch") {
		t.Fatalf("flush into a missing table = %v, want a prepare error", err)
	}

	err := flushBatch(ch.conn, "INSERT INTO peers (sensor_id, node_id, name, url, caps, seen_at)", []int{1},
		func(batch driver.Batch, r int) error { return batch.Append(r) })
	if err == nil || !strings.Contains(err.Error(), "append row") {
		t.Fatalf("flush with a mismatched row = %v, want an append error", err)
	}
	if attempts != 0 {
		t.Fatalf("appended %d rows to an INSERT that could not be prepared", attempts)
	}
}

// TestClickHouseNodeListQueryError checks that a failed node list query is
// returned as an error, not as an empty list.
func TestClickHouseNodeListQueryError(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := newTestClickHouse(t, dsn, ctx)
	defer func() { _ = db.Close() }()

	cancelled, stop := context.WithCancel(ctx)
	stop()
	if _, err := db.(*ClickHouse).NodeList(cancelled, 10); err == nil {
		t.Fatal("NodeList succeeded with a cancelled context")
	}
}
