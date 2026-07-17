package database

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestClickHouseWrites exercises the ClickHouse backend end-to-end against a
// real server. It is skipped unless POLYCLI_TEST_CLICKHOUSE_DSN is set, e.g.
//
//	POLYCLI_TEST_CLICKHOUSE_DSN=clickhouse://localhost:19000/sensor go test ./p2p/database/ -run TestClickHouseWrites -v
//
// The target database must already have the schema from
// sensor-network-tools/clickhouse_schema.sql applied.
func TestClickHouseWrites(t *testing.T) {
	dsn := os.Getenv("POLYCLI_TEST_CLICKHOUSE_DSN")
	if dsn == "" {
		t.Skip("POLYCLI_TEST_CLICKHOUSE_DSN not set; skipping ClickHouse integration test")
	}

	ctx, cancel := context.WithCancel(context.Background())
	db := NewClickHouse(ctx, ClickHouseOptions{
		DSN:                          dsn,
		SensorID:                     "test-sensor",
		ChainID:                      137,
		MaxConcurrency:               10,
		ShouldWriteBlocks:            true,
		ShouldWriteBlockEvents:       true,
		ShouldWriteTransactions:      true,
		ShouldWriteTransactionEvents: true,
		ShouldWritePeers:             true,
		TTL:                          14 * 24 * time.Hour,
	})

	now := time.Now().UTC()
	priv, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	wantSigner := crypto.PubkeyToAddress(priv.PublicKey)
	header := &types.Header{
		Number:     big.NewInt(42),
		Time:       uint64(now.Unix()),
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
	block := types.NewBlockWithHeader(header)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		GasPrice: big.NewInt(2_000_000_000),
		Gas:      21_000,
		Value:    big.NewInt(1),
	})

	db.WriteBlock(ctx, nil, block, big.NewInt(100), now)
	db.WriteTransactions(ctx, nil, []*types.Transaction{tx}, now)
	db.WritePeers(ctx, nil, now) // empty peer slice is fine; exercises the path

	// Trigger the drain flush and give it a moment to complete.
	cancel()
	time.Sleep(750 * time.Millisecond)

	conn, err := clickhouse.Open(mustParseDSN(t, dsn))
	if err != nil {
		t.Fatalf("open verify conn: %v", err)
	}
	defer conn.Close()

	checkCount(t, conn, "blocks", block.Hash().Hex())
	checkCount(t, conn, "transactions", tx.Hash().Hex())

	// Verify the round-tripped block fields.
	var (
		number  uint64
		gasUsed uint64
		baseFee uint64
		signer  string
	)
	row := conn.QueryRow(context.Background(),
		"SELECT number, gas_used, base_fee, signer FROM blocks WHERE hash = ? LIMIT 1", block.Hash().Hex())
	if err := row.Scan(&number, &gasUsed, &baseFee, &signer); err != nil {
		t.Fatalf("scan block: %v", err)
	}
	if number != 42 || gasUsed != 21_000 || baseFee != 1_000_000_000 {
		t.Fatalf("unexpected block fields: number=%d gas_used=%d base_fee=%d", number, gasUsed, baseFee)
	}
	if signer != wantSigner.Hex() {
		t.Fatalf("signer mismatch: want %s got %s", wantSigner.Hex(), signer)
	}
}

func mustParseDSN(t *testing.T, dsn string) *clickhouse.Options {
	t.Helper()
	opts, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		t.Fatalf("parse dsn: %v", err)
	}
	return opts
}

func checkCount(t *testing.T, conn driver.Conn, table, hash string) {
	t.Helper()
	var count uint64
	// #nosec G202 -- table is a test-only constant, not user input
	if err := conn.QueryRow(context.Background(),
		"SELECT count() FROM "+table+" WHERE hash = ?", hash).Scan(&count); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	if count == 0 {
		t.Fatalf("expected a row in %s for hash %s, got none", table, hash)
	}
}
