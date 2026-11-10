package database

import (
	"context"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/p2p"
)

// These tests need no ClickHouse server. They cover the paths taken with no
// connection, where every write is counted toward rows_discarded instead of
// being sent.

// chFlags selects the write settings for a disconnected ClickHouse.
type chFlags struct {
	blocks, blockEvents, firstBlockEvent, txs, txEvents, firstTxEvent, peers bool
}

func disconnectedClickHouse(f chFlags) *ClickHouse {
	return &ClickHouse{
		shouldWriteBlocks:                f.blocks,
		shouldWriteBlockEvents:           f.blockEvents,
		shouldWriteFirstBlockEvent:       f.firstBlockEvent,
		shouldWriteTransactions:          f.txs,
		shouldWriteTransactionEvents:     f.txEvents,
		shouldWriteFirstTransactionEvent: f.firstTxEvent,
		shouldWritePeers:                 f.peers,
	}
}

var allOn = chFlags{blocks: true, blockEvents: true, txs: true, txEvents: true, peers: true}

// TestDiscardedCountsOnlyEnabledRows checks, for every write method, that a
// disconnected ClickHouse counts exactly the rows it would have written with
// the given settings, and nothing for tables that are turned off.
func TestDiscardedCountsOnlyEnabledRows(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	peer := testPeer(t)
	header, _ := signedHeader(t, 1, now)
	tx1 := types.NewTx(&types.LegacyTx{Nonce: 1, GasPrice: big.NewInt(1), Gas: 21_000})
	tx2 := types.NewTx(&types.LegacyTx{Nonce: 2, GasPrice: big.NewInt(1), Gas: 21_000})
	block := types.NewBlockWithHeader(header).WithBody(types.Body{Transactions: []*types.Transaction{tx1, tx2}})
	txs := txsOf(t, block)
	body := &eth.BlockBody{}
	badUncles := &eth.BlockBody{}
	if err := badUncles.Uncles.AppendRaw([]byte{0x01}); err != nil {
		t.Fatal(err)
	}
	ann := BlockAnnouncement{Hash: block.Hash(), Number: 1}
	anns := []BlockAnnouncement{ann, ann, ann}
	hashes := []common.Hash{{1}, {2}}
	peers := []*p2p.Peer{p2p.NewPeer(peer.ID(), "a", nil), p2p.NewPeer(testPeer(t).ID(), "b", nil)}

	cases := []struct {
		name  string
		flags chFlags
		write func(*ClickHouse)
		want  uint64
	}{
		// WriteBlock: event, blocks + block_bodies + block_total_difficulty +
		// 2 block_txs, 2 transactions.
		{"block all on", allOn, func(c *ClickHouse) { c.WriteBlock(ctx, peer, block, txs, big.NewInt(1), now) }, 1 + 5 + 2},
		{"block no peer", allOn, func(c *ClickHouse) { c.WriteBlock(ctx, nil, block, txs, big.NewInt(1), now) }, 5 + 2},
		{"block zero td", allOn, func(c *ClickHouse) { c.WriteBlock(ctx, peer, block, txs, big.NewInt(0), now) }, 1 + 4 + 2},
		{"block first event only", chFlags{firstBlockEvent: true}, func(c *ClickHouse) { c.WriteBlock(ctx, peer, block, txs, big.NewInt(1), now) }, 1},
		{"block all off", chFlags{}, func(c *ClickHouse) { c.WriteBlock(ctx, peer, block, txs, big.NewInt(1), now) }, 0},

		// WriteBlockHeaders: a blocks row and an event per header.
		{"headers all on", allOn, func(c *ClickHouse) { c.WriteBlockHeaders(ctx, []*types.Header{header, header}, now, false) }, 4},
		{"headers blocks only", chFlags{blocks: true}, func(c *ClickHouse) { c.WriteBlockHeaders(ctx, []*types.Header{header, header}, now, true) }, 2},
		{"headers events only", chFlags{blockEvents: true}, func(c *ClickHouse) { c.WriteBlockHeaders(ctx, []*types.Header{header}, now, false) }, 1},
		{"headers all off", chFlags{}, func(c *ClickHouse) { c.WriteBlockHeaders(ctx, []*types.Header{header}, now, false) }, 0},

		// WriteBlockBodyArrival: one event.
		{"arrival events on", chFlags{blockEvents: true}, func(c *ClickHouse) { c.WriteBlockBodyArrival(ctx, ann, now) }, 1},
		{"arrival events off", chFlags{blocks: true, txs: true}, func(c *ClickHouse) { c.WriteBlockBodyArrival(ctx, ann, now) }, 0},

		// WriteBlockBody: block_bodies + 2 block_txs, 2 transactions.
		{"body all on", allOn, func(c *ClickHouse) { c.WriteBlockBody(ctx, body, txs, ann, now) }, 3 + 2},
		{"body blocks only", chFlags{blocks: true}, func(c *ClickHouse) { c.WriteBlockBody(ctx, body, txs, ann, now) }, 3},
		{"body txs only", chFlags{txs: true}, func(c *ClickHouse) { c.WriteBlockBody(ctx, body, txs, ann, now) }, 2},
		{"body all off", chFlags{}, func(c *ClickHouse) { c.WriteBlockBody(ctx, body, txs, ann, now) }, 0},
		// Undecodable uncles skip only the block_bodies row.
		{"body bad uncles", allOn, func(c *ClickHouse) { c.WriteBlockBody(ctx, badUncles, txs, ann, now) }, 2 + 2},

		// WriteBlockEvents and WriteTransactionEvents: one row per item, none
		// without a peer.
		{"block events", allOn, func(c *ClickHouse) { c.WriteBlockEvents(ctx, peer, anns, now) }, 3},
		{"block events no peer", allOn, func(c *ClickHouse) { c.WriteBlockEvents(ctx, nil, anns, now) }, 0},
		{"tx events", allOn, func(c *ClickHouse) { c.WriteTransactionEvents(ctx, peer, hashes, now) }, 2},
		{"tx events no peer", allOn, func(c *ClickHouse) { c.WriteTransactionEvents(ctx, nil, hashes, now) }, 0},

		// WriteTransactions: a full_tx event and a transactions row per transaction.
		{"txs all on", allOn, func(c *ClickHouse) { c.WriteTransactions(ctx, peer, block.Transactions(), now) }, 4},
		{"txs first event only", chFlags{firstTxEvent: true}, func(c *ClickHouse) { c.WriteTransactions(ctx, peer, block.Transactions(), now) }, 2},
		{"txs no peer", allOn, func(c *ClickHouse) { c.WriteTransactions(ctx, nil, block.Transactions(), now) }, 2},
		{"txs all off", chFlags{}, func(c *ClickHouse) { c.WriteTransactions(ctx, peer, block.Transactions(), now) }, 0},

		// WritePeers: one row per peer.
		{"peers on", allOn, func(c *ClickHouse) { c.WritePeers(ctx, peers, now) }, 2},
		{"peers off", chFlags{}, func(c *ClickHouse) { c.WritePeers(ctx, peers, now) }, 0},

		// WriteBlockHashFirstSeen writes nothing; the first-seen rollup is derived.
		{"hash first seen", allOn, func(c *ClickHouse) { c.WriteBlockHashFirstSeen(ctx, peer, ann.Hash, now) }, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := disconnectedClickHouse(tc.flags)
			tc.write(c)
			if got := c.discarded.Load(); got != tc.want {
				t.Fatalf("discarded = %d, want %d", got, tc.want)
			}
		})
	}
}

// TestClickHouseSettingsGetters checks that each write setting is reported as
// configured.
func TestClickHouseSettingsGetters(t *testing.T) {
	c := &ClickHouse{
		maxConcurrency:                   7,
		shouldWriteBlocks:                true,
		shouldWriteFirstBlockEvent:       true,
		shouldWriteTransactionEvents:     true,
		shouldWriteFirstTransactionEvent: true,
	}
	if c.MaxConcurrentWrites() != 7 {
		t.Errorf("MaxConcurrentWrites = %d, want 7", c.MaxConcurrentWrites())
	}
	for name, got := range map[string][2]bool{
		"ShouldWriteBlocks":                {c.ShouldWriteBlocks(), true},
		"ShouldWriteBlockEvents":           {c.ShouldWriteBlockEvents(), false},
		"ShouldWriteFirstBlockEvent":       {c.ShouldWriteFirstBlockEvent(), true},
		"ShouldWriteTransactions":          {c.ShouldWriteTransactions(), false},
		"ShouldWriteTransactionEvents":     {c.ShouldWriteTransactionEvents(), true},
		"ShouldWriteFirstTransactionEvent": {c.ShouldWriteFirstTransactionEvent(), true},
		"ShouldWritePeers":                 {c.ShouldWritePeers(), false},
	} {
		if got[0] != got[1] {
			t.Errorf("%s = %v, want %v", name, got[0], got[1])
		}
	}
}

// TestNodeListWithoutConnection checks that a disconnected ClickHouse returns an
// empty node list, while ClickHouseNodeList, which nodelist uses, fails instead
// of returning one.
func TestNodeListWithoutConnection(t *testing.T) {
	ctx := context.Background()
	nodes, err := disconnectedClickHouse(chFlags{}).NodeList(ctx, 10)
	if err != nil || len(nodes) != 0 {
		t.Fatalf("NodeList = %v, %v; want an empty list", nodes, err)
	}
	_, err = ClickHouseNodeList(ctx, "clickhouse://127.0.0.1:1/nope", 10)
	if err == nil || !strings.Contains(err.Error(), "connect to ClickHouse") {
		t.Fatalf("ClickHouseNodeList against an unreachable server = %v, want a connect error", err)
	}
	if _, err := ClickHouseNodeList(ctx, "not a dsn", 10); err == nil {
		t.Fatal("ClickHouseNodeList succeeded with an invalid DSN")
	}
}
