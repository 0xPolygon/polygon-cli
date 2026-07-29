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
	header, wantSigner := signedHeader(t, 42, now)
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
	if number != 42 || gasUsed != 21_000 || baseFee.Uint64() != 1_000_000_000 {
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
	header, _ := signedHeader(t, 4242, now)
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
