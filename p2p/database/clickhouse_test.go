package database

import (
	"context"
	"math/big"
	"net"
	"os"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

// Test block heights live in a reserved band far above any real chain head, one
// per test, so a test row is never mistaken for chain data and two tests never
// collide on a height. This matters because these tests write to a real database
// that is usually the local-stack one holding live sensor data: heights 42 and
// 4242 previously produced rows that looked exactly like a transaction included in
// four competing blocks, which is a real reorg signature.
const (
	heightWrites          = 900_000_001
	heightHeaderNoClobber = 900_000_002
	heightProductionFlags = 900_000_003
	heightParentCancel    = 900_000_004
	heightTotalDiff       = 900_000_005
)

// The integration tests in this file are skipped unless
// POLYCLI_TEST_CLICKHOUSE_DSN is set, e.g.
//
//	POLYCLI_TEST_CLICKHOUSE_DSN=clickhouse://localhost:19000/sensor go test ./p2p/database/ -run TestClickHouse -v
//
// The target database must already have clickhouse_schema.sql applied.

func clickHouseDSN(t *testing.T) string {
	t.Helper()
	dsn := os.Getenv("POLYCLI_TEST_CLICKHOUSE_DSN")
	if dsn == "" {
		t.Skip("POLYCLI_TEST_CLICKHOUSE_DSN not set; skipping ClickHouse integration test")
	}
	return dsn
}

// signedHeader returns a clique-sealed header plus the address that sealed it, so
// tests can assert the ecrecover-at-ingest path.
func signedHeader(t *testing.T, number int64, blockTime time.Time) (*types.Header, string) {
	t.Helper()
	priv, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	header := &types.Header{
		Number:     big.NewInt(number),
		Time:       uint64(blockTime.Unix()),
		Difficulty: big.NewInt(7),
		GasLimit:   30_000_000,
		GasUsed:    21_000,
		BaseFee:    big.NewInt(1_000_000_000),
		Extra:      make([]byte, crypto.SignatureLength),
	}
	sig, err := crypto.Sign(clique.SealHash(header).Bytes(), priv)
	if err != nil {
		t.Fatalf("sign header: %v", err)
	}
	copy(header.Extra[len(header.Extra)-crypto.SignatureLength:], sig)
	return header, crypto.PubkeyToAddress(priv.PublicKey).Hex()
}

func testPeer(t *testing.T) *enode.Node {
	t.Helper()
	priv, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate peer key: %v", err)
	}
	return enode.NewV4(&priv.PublicKey, net.IPv4(127, 0, 0, 1), 30303, 30303)
}

func newTestClickHouse(t *testing.T, dsn string, ctx context.Context) Database {
	t.Helper()
	return NewClickHouse(ctx, ClickHouseOptions{
		DSN:                          dsn,
		SensorID:                     "test-sensor",
		ChainID:                      137,
		MaxConcurrency:               10,
		ShouldWriteBlocks:            true,
		ShouldWriteBlockEvents:       true,
		ShouldWriteTransactions:      true,
		ShouldWriteTransactionEvents: true,
		ShouldWritePeers:             true,
	})
}

// TestClickHouseWrites exercises the whole write path against a real server and
// asserts every table the sensor targets receives its row.
func TestClickHouseWrites(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := newTestClickHouse(t, dsn, ctx)

	now := time.Now().UTC()
	header, wantSigner := signedHeader(t, heightWrites, now)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		GasPrice: big.NewInt(2_000_000_000),
		Gas:      21_000,
		Value:    big.NewInt(1),
		Data:     []byte{0xa9, 0x05, 0x9c, 0xbb, 0x01, 0x02},
	})
	block := types.NewBlockWithHeader(header).WithBody(types.Body{Transactions: []*types.Transaction{tx}})
	peer := testPeer(t)

	db.WriteBlock(ctx, peer, block, big.NewInt(100), now)
	db.WriteTransactions(ctx, peer, []*types.Transaction{tx}, now)
	db.WritePeers(ctx, nil, now) // empty peer slice is fine; exercises the path

	// Close drains the buffered rows synchronously before we verify them.
	if cerr := db.Close(); cerr != nil {
		t.Fatalf("close db: %v", cerr)
	}

	conn := verifyConn(t, dsn)
	blockHash := block.Hash().Hex()

	// Fact tables.
	checkCount(t, conn, "SELECT count() FROM blocks WHERE hash = ?", blockHash)
	checkCount(t, conn, "SELECT count() FROM block_bodies WHERE hash = ?", blockHash)
	checkCount(t, conn, "SELECT count() FROM block_txs WHERE block_hash = ?", blockHash)
	checkCount(t, conn, "SELECT count() FROM transactions WHERE hash = ?", tx.Hash().Hex())

	// Observation stream, and the rollup the materialized view maintains off it.
	checkCount(t, conn, "SELECT count() FROM block_events WHERE block_hash = ?", blockHash)
	checkCount(t, conn, "SELECT count() FROM block_events_first WHERE block_hash = ?", blockHash)

	// Round-tripped header fields. base_fee is UInt256, so it must be scanned as
	// a big.Int rather than the uint64 the old schema used.
	var (
		number  uint64
		gasUsed uint64
		baseFee big.Int
		signer  string
	)
	row := conn.QueryRow(context.Background(),
		"SELECT number, gas_used, base_fee, signer FROM blocks WHERE hash = ? LIMIT 1", blockHash)
	if err := row.Scan(&number, &gasUsed, &baseFee, &signer); err != nil {
		t.Fatalf("scan block: %v", err)
	}
	if number != heightWrites || gasUsed != 21_000 || baseFee.Uint64() != 1_000_000_000 {
		t.Fatalf("unexpected block fields: number=%d gas_used=%d base_fee=%s", number, gasUsed, baseFee.String())
	}
	if signer != wantSigner {
		t.Fatalf("signer mismatch: want %s got %s", wantSigner, signer)
	}

	// The event must carry the announced total difficulty and the peer's node
	// id (not its enode URL), since node id is what joins to the peer tables.
	var (
		td     big.Int
		nodeID string
		source string
	)
	row = conn.QueryRow(context.Background(),
		"SELECT total_difficulty, node_id, source FROM block_events WHERE block_hash = ? AND source = 'new_block' LIMIT 1", blockHash)
	if err := row.Scan(&td, &nodeID, &source); err != nil {
		t.Fatalf("scan event: %v", err)
	}
	if td.Uint64() != 100 {
		t.Fatalf("total_difficulty: want 100 got %s", td.String())
	}
	if nodeID != peer.ID().String() {
		t.Fatalf("node_id: want %s got %s", peer.ID().String(), nodeID)
	}

	// The block -> tx mapping, which the previous schema could not express.
	var mapped string
	if err := conn.QueryRow(context.Background(),
		"SELECT tx_hash FROM block_txs WHERE block_hash = ? AND tx_index = 0 LIMIT 1", blockHash).Scan(&mapped); err != nil {
		t.Fatalf("scan block_txs: %v", err)
	}
	if mapped != tx.Hash().Hex() {
		t.Fatalf("block_txs tx_hash: want %s got %s", tx.Hash().Hex(), mapped)
	}

	// Calldata shape is recorded without the calldata.
	var (
		selector  string
		inputSize uint32
	)
	if err := conn.QueryRow(context.Background(),
		"SELECT input_selector, input_size FROM transactions WHERE hash = ? LIMIT 1", tx.Hash().Hex()).Scan(&selector, &inputSize); err != nil {
		t.Fatalf("scan transaction: %v", err)
	}
	if selector != "0xa9059cbb" || inputSize != 6 {
		t.Fatalf("calldata shape: got selector=%s size=%d", selector, inputSize)
	}
}

// TestClickHouseHeaderDoesNotClobberBody is the regression test for the defect
// that motivated splitting body facts out of the blocks table.
//
// Under the previous schema a header carried no transaction count, so the header
// path wrote tx_count = 0 onto the block row; because the engine kept the row
// with the newest version column, a header arriving after the full block --
// routine during parent backfill -- permanently replaced the real counts with
// zeros. Body facts now live in their own hash-keyed table that the header path
// never writes, so the ordering cannot matter.
func TestClickHouseHeaderDoesNotClobberBody(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := newTestClickHouse(t, dsn, ctx)

	now := time.Now().UTC()
	header, _ := signedHeader(t, heightHeaderNoClobber, now)
	txs := []*types.Transaction{
		types.NewTx(&types.LegacyTx{Nonce: 1, GasPrice: big.NewInt(1), Gas: 21_000, Value: big.NewInt(1)}),
		types.NewTx(&types.LegacyTx{Nonce: 2, GasPrice: big.NewInt(1), Gas: 21_000, Value: big.NewInt(2)}),
		types.NewTx(&types.LegacyTx{Nonce: 3, GasPrice: big.NewInt(1), Gas: 21_000, Value: big.NewInt(3)}),
	}
	block := types.NewBlockWithHeader(header).WithBody(types.Body{Transactions: txs})
	peer := testPeer(t)

	// Full block first, then the same header again a second later -- the ordering
	// that used to lose the counts.
	db.WriteBlock(ctx, peer, block, big.NewInt(100), now)
	db.WriteBlockHeaders(ctx, []*types.Header{block.Header()}, now.Add(time.Second), true)

	if cerr := db.Close(); cerr != nil {
		t.Fatalf("close db: %v", cerr)
	}

	conn := verifyConn(t, dsn)
	blockHash := block.Hash().Hex()

	var txCount uint32
	if err := conn.QueryRow(context.Background(),
		"SELECT max(tx_count) FROM block_bodies WHERE hash = ?", blockHash).Scan(&txCount); err != nil {
		t.Fatalf("scan block_bodies: %v", err)
	}
	if txCount != uint32(len(txs)) {
		t.Fatalf("tx_count was clobbered by the header write: want %d got %d", len(txs), txCount)
	}

	// Both events should be recorded, distinguishable by source, and the
	// backfilled header must be tagged as such.
	var sources []string
	rows, err := conn.Query(context.Background(),
		"SELECT DISTINCT source FROM block_events WHERE block_hash = ? ORDER BY source", blockHash)
	if err != nil {
		t.Fatalf("query sources: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			t.Errorf("close rows: %v", err)
		}
	}()
	for rows.Next() {
		var source string
		if err := rows.Scan(&source); err != nil {
			t.Fatalf("scan source: %v", err)
		}
		sources = append(sources, source)
	}
	if len(sources) != 2 || sources[0] != "header_backfill" || sources[1] != "new_block" {
		t.Fatalf("want sources [header_backfill new_block], got %v", sources)
	}
}

func verifyConn(t *testing.T, dsn string) driver.Conn {
	t.Helper()
	opts, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		t.Fatalf("parse dsn: %v", err)
	}
	conn, err := clickhouse.Open(opts)
	if err != nil {
		t.Fatalf("open verify conn: %v", err)
	}
	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Errorf("close verify conn: %v", err)
		}
	})
	return conn
}

func checkCount(t *testing.T, conn driver.Conn, query, arg string) {
	t.Helper()
	var count uint64
	if err := conn.QueryRow(context.Background(), query, arg).Scan(&count); err != nil {
		t.Fatalf("count via %q: %v", query, err)
	}
	if count == 0 {
		t.Fatalf("expected at least one row from %q for %s, got none", query, arg)
	}
}

// TestClickHouseProductionFlagsRecordProvenance is the regression test for two
// instances of the same defect: provenance events gated on the flag that controls
// the full per-peer announcement stream, rather than on "either event flag is set".
//
// Production runs write_block_events=false / write_first_block_event=true and the
// same pair for transactions. Under that configuration new_block, header,
// header_backfill, body and full_tx were all silently dropped -- the sensors looked
// healthy and the volume-bounded first-event rows kept arriving, so nothing pointed
// at the gap. The block sources were fixed first; full_tx survived one more round
// and is why this test asserts both sides together.
//
// These sources are not volume. Per (block, sensor) the per-peer streams are
// hash_announce at ~52 rows and tx hash_announce at ~8, both scaling with
// --max-peers; the provenance sources are 1-8 rows and full_tx ~2 per transaction.
func TestClickHouseProductionFlagsRecordProvenance(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// The exact terraform default flag set, which is the point of the test.
	db := NewClickHouse(ctx, ClickHouseOptions{
		DSN:                              dsn,
		SensorID:                         "test-sensor-prod-flags",
		ChainID:                          137,
		MaxConcurrency:                   10,
		ShouldWriteBlocks:                true,
		ShouldWriteBlockEvents:           false,
		ShouldWriteFirstBlockEvent:       true,
		ShouldWriteTransactions:          true,
		ShouldWriteTransactionEvents:     false,
		ShouldWriteFirstTransactionEvent: true,
		ShouldWritePeers:                 true,
	})

	now := time.Now().UTC()
	header, _ := signedHeader(t, heightProductionFlags, now)

	// The nonce must vary per run. These tables are append-only and the test asserts
	// on row presence, so a fixed nonce yields a fixed transaction hash and rows left
	// behind by an earlier run satisfy the assertion -- the test then passes even with
	// the defect reintroduced, which is how the first draft of it failed to catch the
	// very bug it exists for. Block hashes are already unique per run because
	// signedHeader seals with a fresh key.
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    uint64(now.UnixNano()),
		GasPrice: big.NewInt(2_000_000_000),
		Gas:      21_000,
		Value:    big.NewInt(1),
	})
	block := types.NewBlockWithHeader(header).WithBody(types.Body{Transactions: []*types.Transaction{tx}})
	peer := testPeer(t)

	db.WriteBlock(ctx, peer, block, big.NewInt(555), now)
	db.WriteBlockHeaders(ctx, []*types.Header{header}, now, false)
	db.WriteTransactions(ctx, peer, []*types.Transaction{tx}, now)

	if cerr := db.Close(); cerr != nil {
		t.Fatalf("close db: %v", cerr)
	}

	conn := verifyConn(t, dsn)
	blockHash := block.Hash().Hex()

	for _, source := range []string{"new_block", "header"} {
		checkCount(t, conn,
			"SELECT count() FROM block_events WHERE block_hash = ? AND source = '"+source+"'", blockHash)
	}

	// The one that regressed after the block-side fix.
	checkCount(t, conn,
		"SELECT count() FROM tx_events WHERE tx_hash = ? AND source = 'full_tx'", tx.Hash().Hex())

	// total_difficulty rides on new_block alone, so losing that source loses the
	// column entirely -- assert the value survived, not just the row.
	var td big.Int
	if err := conn.QueryRow(context.Background(),
		"SELECT total_difficulty FROM block_events WHERE block_hash = ? AND source = 'new_block' LIMIT 1",
		blockHash).Scan(&td); err != nil {
		t.Fatalf("scan total_difficulty: %v", err)
	}
	if td.Uint64() != 555 {
		t.Fatalf("total_difficulty: want 555 got %s", td.String())
	}
}

// TestClickHouseWritesSurviveParentContextCancel is the regression test for silent
// shutdown data loss.
//
// The sensor shuts down by cancelling its signal context, and only stops serving
// peers afterwards (stopServer/conns.Close are deferred later than db.Close, so
// they run first). While the p2p server winds down, peers keep delivering blocks.
// When the batcher context inherited the caller's, cancellation drained and
// stopped the batchers at the instant of SIGINT, so every row written during that
// window went into the buffered channel with no reader -- never flushed, and not
// counted as dropped because the channel had capacity, so nothing logged it.
//
// Close, and only Close, may stop the batchers.
func TestClickHouseWritesSurviveParentContextCancel(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())

	db := newTestClickHouse(t, dsn, ctx)

	// SIGINT arrives.
	cancel()
	time.Sleep(300 * time.Millisecond)

	// A peer delivers a block while the p2p server is still winding down.
	now := time.Now().UTC()
	header, _ := signedHeader(t, heightParentCancel, now)
	db.WriteBlockHeaders(ctx, []*types.Header{header}, now, false)

	// Only now does the sensor close the database, which must drain.
	if cerr := db.Close(); cerr != nil {
		t.Fatalf("close db: %v", cerr)
	}

	conn := verifyConn(t, dsn)
	blockHash := header.Hash().Hex()
	checkCount(t, conn, "SELECT count() FROM blocks WHERE hash = ?", blockHash)
}

// TestClickHouseTotalDifficultySurvivesEventExpiry is the regression test for
// total_difficulty decaying to zero.
//
// total_difficulty reaches the sensor only on a NewBlock announcement, so it used
// to be carried on block_events_first purely to outlive the raw event stream. When
// retention was normalised that rollup became 14 days itself, while blocks stayed
// forever -- so v_blocks returned 0 for every block older than the TTL, which is
// also the documented value for "no peer ever announced it to us". The two cases
// were indistinguishable.
//
// It now lives in its own forever-kept, hash-keyed table, and absence rather than 0
// means never announced. Deleting the events stands in for the TTL expiring them.
func TestClickHouseTotalDifficultySurvivesEventExpiry(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := newTestClickHouse(t, dsn, ctx)

	now := time.Now().UTC()
	header, _ := signedHeader(t, heightTotalDiff, now)
	block := types.NewBlockWithHeader(header)
	wantTD := big.NewInt(987654321)

	db.WriteBlock(ctx, testPeer(t), block, wantTD, now)
	if cerr := db.Close(); cerr != nil {
		t.Fatalf("close db: %v", cerr)
	}

	conn := verifyConn(t, dsn)
	blockHash := block.Hash().Hex()

	// Expire every observation of this block, which is what the 14-day TTL does.
	for _, tbl := range []string{"block_events", "block_events_first"} {
		if err := conn.Exec(context.Background(),
			"ALTER TABLE "+tbl+" DELETE WHERE block_hash = ? SETTINGS mutations_sync = 1",
			blockHash); err != nil {
			t.Fatalf("expire %s: %v", tbl, err)
		}
	}

	var (
		have bool
		td   *big.Int
	)
	if err := conn.QueryRow(context.Background(),
		"SELECT have_total_difficulty, total_difficulty FROM v_blocks WHERE hash = ? LIMIT 1",
		blockHash).Scan(&have, &td); err != nil {
		t.Fatalf("scan v_blocks: %v", err)
	}
	if !have {
		t.Fatal("total difficulty was lost with the events it was announced in")
	}
	if td == nil || td.Cmp(wantTD) != 0 {
		t.Fatalf("total_difficulty: want %s got %v", wantTD, td)
	}
}
