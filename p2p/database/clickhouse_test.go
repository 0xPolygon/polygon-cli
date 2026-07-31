package database

import (
	"context"
	"math/big"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

// Test block heights live in a reserved band far above any real chain head, one
// per test, so a test row is never mistaken for chain data and two tests never
// collide on a height. Spaced by ten, so a test needing a parent or child can
// derive height-1 or height+1 without landing on another test's block -- deriving a
// parent as height-1 from consecutive constants collided, and the resulting failure
// appeared only when the whole suite ran. This matters because these tests write to a real database
// that is usually the local-stack one holding live sensor data: heights 42 and
// 4242 previously produced rows that looked exactly like a transaction included in
// four competing blocks, which is a real reorg signature.
const (
	heightWrites          = 900_000_010
	heightHeaderNoClobber = 900_000_020
	heightProductionFlags = 900_000_030
	heightParentCancel    = 900_000_040
	heightTotalDiff       = 900_000_050
	heightLowercase       = 900_000_060
	heightHeaderEvents    = 900_000_070
)

// cleanupTestHeights removes everything the given test heights produced, in
// dependency order (hash-keyed tables first, via the blocks rows). Registered by
// every test that writes blocks: block_forks has no TTL and v_forks no height
// filter, so without this each run's fresh sealing key added another competing
// hash per height -- six heights wide-N "forks" after N runs, unfilterable and
// permanent.
func cleanupTestHeights(t *testing.T, conn driver.Conn, heights ...uint64) {
	t.Helper()
	t.Cleanup(func() {
		ctx := context.Background()

		// Order matters: every statement below identifies its rows through a table
		// that a LATER statement deletes, so reversing any pair silently leaves rows
		// behind. Getting this wrong is why two earlier attempts still leaked --
		// deleting tx_events first emptied the subquery that finds the transactions.
		//
		// transactions is reachable two ways, and both are needed: WriteTransactions
		// writes tx_events for it, but writeTxs is ALSO called from WriteBlockBody,
		// which writes no tx_events at all -- so a body-only test leaves transactions
		// rows that no tx_events row points at.
		for _, h := range heights {
			for _, q := range []string{
				"ALTER TABLE transactions DELETE WHERE hash IN (SELECT tx_hash FROM block_txs WHERE block_hash IN (SELECT hash FROM blocks WHERE number = ?)) SETTINGS mutations_sync = 1",
				"ALTER TABLE tx_events_first DELETE WHERE tx_hash IN (SELECT tx_hash FROM block_txs WHERE block_hash IN (SELECT hash FROM blocks WHERE number = ?)) SETTINGS mutations_sync = 1",
				"ALTER TABLE block_txs DELETE WHERE block_hash IN (SELECT hash FROM blocks WHERE number = ?) SETTINGS mutations_sync = 1",
				"ALTER TABLE block_bodies DELETE WHERE hash IN (SELECT hash FROM blocks WHERE number = ?) SETTINGS mutations_sync = 1",
				"ALTER TABLE block_total_difficulty DELETE WHERE hash IN (SELECT hash FROM blocks WHERE number = ?) SETTINGS mutations_sync = 1",
				"ALTER TABLE blocks DELETE WHERE number = ? SETTINGS mutations_sync = 1",
				"ALTER TABLE block_events DELETE WHERE block_number = ? SETTINGS mutations_sync = 1",
				"ALTER TABLE block_events_first DELETE WHERE block_number = ? SETTINGS mutations_sync = 1",
				"ALTER TABLE block_forks DELETE WHERE number = ? SETTINGS mutations_sync = 1",
			} {
				if err := conn.Exec(ctx, q, h); err != nil {
					t.Logf("cleanup height %d: %v", h, err)
				}
			}
		}

		// Whatever the test wrote under its own sensor id, by any path. Last, because
		// the height-keyed statements above use tx_events to find nothing but this
		// catches what they could not reach.
		for _, q := range []string{
			"ALTER TABLE transactions DELETE WHERE hash IN (SELECT tx_hash FROM tx_events WHERE sensor_id LIKE 'test-sensor%') SETTINGS mutations_sync = 1",
			"ALTER TABLE tx_events_first DELETE WHERE tx_hash IN (SELECT tx_hash FROM tx_events WHERE sensor_id LIKE 'test-sensor%') SETTINGS mutations_sync = 1",
			"ALTER TABLE tx_events DELETE WHERE sensor_id LIKE 'test-sensor%' SETTINGS mutations_sync = 1",
		} {
			if err := conn.Exec(ctx, q); err != nil {
				t.Logf("tx-side cleanup: %v", err)
			}
		}
	})
}

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
	// Lowercase, because that is what the writer stores -- see addressHex.
	return header, strings.ToLower(crypto.PubkeyToAddress(priv.PublicKey).Hex())
}

// signedHeaderWithCoinbase is signedHeader for tests that need a specific coinbase.
// The seal covers Coinbase, so it must be set BEFORE signing; assigning it to an
// already-sealed header silently invalidates the signature and ecrecover then
// returns a different address entirely.
func signedHeaderWithCoinbase(t *testing.T, number int64, blockTime time.Time, coinbase common.Address) (*types.Header, string) {
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
		Coinbase:   coinbase,
		Extra:      make([]byte, crypto.SignatureLength),
	}
	sig, err := crypto.Sign(clique.SealHash(header).Bytes(), priv)
	if err != nil {
		t.Fatalf("sign header: %v", err)
	}
	copy(header.Extra[len(header.Extra)-crypto.SignatureLength:], sig)
	return header, strings.ToLower(crypto.PubkeyToAddress(priv.PublicKey).Hex())
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
	// The nonce varies per run: with it fixed, each run mapped the SAME tx hash to
	// a fresh block hash (the sealing key is fresh, so the header hash changes),
	// and block_txs accumulated one-transaction-in-N-competing-blocks -- a real
	// reorg signature, manufactured by the test suite.
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    uint64(now.UnixNano()),
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
	cleanupTestHeights(t, conn, heightWrites)
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
		types.NewTx(&types.LegacyTx{Nonce: uint64(now.UnixNano()) + 1, GasPrice: big.NewInt(1), Gas: 21_000, Value: big.NewInt(1)}),
		types.NewTx(&types.LegacyTx{Nonce: uint64(now.UnixNano()) + 2, GasPrice: big.NewInt(1), Gas: 21_000, Value: big.NewInt(2)}),
		types.NewTx(&types.LegacyTx{Nonce: uint64(now.UnixNano()) + 3, GasPrice: big.NewInt(1), Gas: 21_000, Value: big.NewInt(3)}),
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
	cleanupTestHeights(t, conn, heightHeaderNoClobber)
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
	cleanupTestHeights(t, conn, heightProductionFlags)
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
	cleanupTestHeights(t, conn, heightParentCancel)
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
	cleanupTestHeights(t, conn, heightTotalDiff)
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

// TestAddressesAreStoredLowercase guards the address casing convention.
//
// common.Address.Hex() applies the EIP-55 checksum and returns mixed case. Stored
// that way, an address column cannot be joined: ClickHouse comparison is
// case-sensitive and every validator identity in this pipeline -- the Polygon
// staking API, block-latency, data-analysis -- is lowercase. The join returns no
// rows and no error, so the wrong answer looks like a real one.
func TestAddressesAreStoredLowercase(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := newTestClickHouse(t, dsn, ctx)

	now := time.Now().UTC()
	// A coinbase whose checksummed form is mixed case, set before sealing.
	mixedCoinbase := common.HexToAddress("0x25B9fC2ED95BBAa9c030e57C860545a17694F90D")
	header, wantSigner := signedHeaderWithCoinbase(t, heightLowercase, now, mixedCoinbase)
	// Signed, so types.Sender can recover from_address -- an unsigned transaction
	// yields an empty sender and the assertion below would test nothing.
	senderKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate sender key: %v", err)
	}
	chainID := big.NewInt(137)
	tx, err := types.SignTx(
		types.NewTx(&types.LegacyTx{
			Nonce:    uint64(now.UnixNano()),
			GasPrice: big.NewInt(2_000_000_000),
			Gas:      21_000,
			To:       &mixedCoinbase,
			Value:    big.NewInt(1),
		}),
		types.LatestSignerForChainID(chainID), senderKey)
	if err != nil {
		t.Fatalf("sign tx: %v", err)
	}
	wantFrom := strings.ToLower(crypto.PubkeyToAddress(senderKey.PublicKey).Hex())
	block := types.NewBlockWithHeader(header).WithBody(types.Body{Transactions: []*types.Transaction{tx}})

	db.WriteBlock(ctx, testPeer(t), block, big.NewInt(1), now)
	db.WriteTransactions(ctx, testPeer(t), []*types.Transaction{tx}, now)
	if cerr := db.Close(); cerr != nil {
		t.Fatalf("close db: %v", cerr)
	}

	conn := verifyConn(t, dsn)
	cleanupTestHeights(t, conn, heightLowercase)

	// The signer must round-trip as the lowercase form of the recovered address,
	// proving it is lowercased rather than merely absent.
	var signer, coinbase string
	if err := conn.QueryRow(context.Background(),
		"SELECT signer, coinbase FROM blocks WHERE hash = ? LIMIT 1",
		block.Hash().Hex()).Scan(&signer, &coinbase); err != nil {
		t.Fatalf("scan blocks: %v", err)
	}
	if signer != wantSigner {
		t.Fatalf("signer: want %s got %s", wantSigner, signer)
	}
	if want := strings.ToLower(mixedCoinbase.Hex()); coinbase != want {
		t.Fatalf("coinbase: want %s got %s", want, coinbase)
	}

	var from, to string
	if err := conn.QueryRow(context.Background(),
		"SELECT from_address, to_address FROM transactions WHERE hash = ? LIMIT 1",
		tx.Hash().Hex()).Scan(&from, &to); err != nil {
		t.Fatalf("scan transactions: %v", err)
	}
	for name, got := range map[string]string{"from_address": from, "to_address": to} {
		if got == "" {
			t.Fatalf("%s was empty", name)
		}
		if got != strings.ToLower(got) {
			t.Fatalf("%s is not lowercase: %s", name, got)
		}
	}
	if from != wantFrom {
		t.Fatalf("from_address: want %s got %s", wantFrom, from)
	}
	if want := strings.ToLower(mixedCoinbase.Hex()); to != want {
		t.Fatalf("to_address: want %s got %s", want, to)
	}
}

// TestHeaderEventsDoNotRequireWriteBlocks is the third and last instance of the
// provenance-gating defect, after new_block/body and full_tx.
//
// WriteBlockHeaders returned early on !shouldWriteBlocks, before reaching
// recordsBlockEvents, which made header and header_backfill the only two provenance
// sources that also required --write-blocks. A fleet with it off still requests
// headers -- getBlockData has no such gate -- so the events were produced and thrown
// away, silently emptying v_block_provenance's header rows and losing the
// header_backfill marker that replaced Datastore's IsParent.
//
// The header row and the header event are separate concerns and are now gated
// separately.
func TestHeaderEventsDoNotRequireWriteBlocks(t *testing.T) {
	dsn := clickHouseDSN(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Block facts off, block events on.
	db := NewClickHouse(ctx, ClickHouseOptions{
		DSN:                        dsn,
		SensorID:                   "test-sensor-header-events",
		ChainID:                    137,
		MaxConcurrency:             10,
		ShouldWriteBlocks:          false,
		ShouldWriteBlockEvents:     false,
		ShouldWriteFirstBlockEvent: true,
	})

	now := time.Now().UTC()
	header, _ := signedHeader(t, heightHeaderEvents, now)
	parent, _ := signedHeader(t, heightHeaderEvents-1, now)

	db.WriteBlockHeaders(ctx, []*types.Header{header}, now, false)
	db.WriteBlockHeaders(ctx, []*types.Header{parent}, now, true)
	if cerr := db.Close(); cerr != nil {
		t.Fatalf("close db: %v", cerr)
	}

	conn := verifyConn(t, dsn)
	cleanupTestHeights(t, conn, heightHeaderEvents, heightHeaderEvents-1)

	// Both header sources must be recorded on the event flag alone.
	for hash, want := range map[string]string{
		header.Hash().Hex(): "header",
		parent.Hash().Hex(): "header_backfill",
	} {
		checkCount(t, conn,
			"SELECT count() FROM block_events WHERE block_hash = ? AND source = '"+want+"'", hash)
	}

	// And the fact row must still be suppressed -- the two gates are independent, so
	// this proves the fix separated them rather than just widening one.
	var blocks uint64
	if err := conn.QueryRow(context.Background(),
		"SELECT count() FROM blocks WHERE number IN (?, ?)",
		uint64(heightHeaderEvents), uint64(heightHeaderEvents-1)).Scan(&blocks); err != nil {
		t.Fatalf("scan blocks: %v", err)
	}
	if blocks != 0 {
		t.Fatalf("write_blocks=false should write no blocks rows, got %d", blocks)
	}
}

// TestUnavailableBackendKeepsWarning covers the degraded path, which used to be
// invisible: one error at startup and then every write silently discarded while the
// sensor peered and looked healthy. A ClickHouse auth failure did exactly that for
// an hour across two sensors.
//
// It also guards the shutdown hang found while writing it -- the warning goroutine
// listened on the caller's context while Close only cancels its own, so Close waited
// forever on a goroutine nothing could stop.
func TestUnavailableBackendKeepsWarning(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := NewClickHouse(ctx, ClickHouseOptions{
		DSN:               "clickhouse://127.0.0.1:1/nope",
		SensorID:          "warn",
		ChainID:           137,
		MaxConcurrency:    1,
		ShouldWriteBlocks: true,
	})
	ch, ok := db.(*ClickHouse)
	if !ok {
		t.Fatal("not a *ClickHouse")
	}
	if ch.conn != nil {
		t.Skip("unexpectedly connected")
	}
	// HasBlock on a nil conn must suppress backfill. It is a read, so it must NOT
	// count toward discarded -- counting it made the warning report call volume,
	// and read 0 forever under --write-blocks=false, where HasBlock is unreachable.
	if !db.HasBlock(ctx, [32]byte{1}) {
		t.Fatal("HasBlock should return true with no connection, to suppress backfill")
	}
	if got := ch.discarded.Load(); got != 0 {
		t.Fatalf("HasBlock counted toward discarded: %d", got)
	}

	// Writes are what the counter measures: a block with two transactions is three
	// rows lost.
	now := time.Now().UTC()
	header, _ := signedHeader(t, heightWrites, now)
	tx1 := types.NewTx(&types.LegacyTx{Nonce: uint64(now.UnixNano()), GasPrice: big.NewInt(1), Gas: 21_000, Value: big.NewInt(1)})
	tx2 := types.NewTx(&types.LegacyTx{Nonce: uint64(now.UnixNano()) + 1, GasPrice: big.NewInt(1), Gas: 21_000, Value: big.NewInt(1)})
	block := types.NewBlockWithHeader(header).WithBody(types.Body{Transactions: []*types.Transaction{tx1, tx2}})
	db.WriteBlock(ctx, testPeer(t), block, big.NewInt(1), now)
	// 4 block-level rows (blocks, block_bodies, block_total_difficulty, the event)
	// plus block_txs and transactions per transaction.
	if got := ch.discarded.Load(); got != 8 {
		t.Fatalf("discarded: want 8 (4 block rows + 2 per tx), got %d", got)
	}
	// The warner goroutine must be registered with the WaitGroup so Close waits.
	done := make(chan struct{})
	go func() { _ = db.Close(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Close hung: the warning goroutine is not stopping on cancel")
	}
}
