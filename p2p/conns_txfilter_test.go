package p2p

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// withBroadcastHead points the shared Conns at a usable head and enables
// transaction broadcasting for the duration of a test, restoring both after.
func withBroadcastHead(t *testing.T, c *Conns, validator *TxValidator) {
	t.Helper()

	head, broadcast, prev := c.head.Get(), c.shouldBroadcastTx, c.txValidator
	t.Cleanup(func() {
		c.head.Set(head)
		c.shouldBroadcastTx = broadcast
		c.txValidator = prev
	})

	c.head.Set(NewBlockPacket{Block: types.NewBlockWithHeader(testHead()), TD: big.NewInt(1)})
	c.shouldBroadcastTx = true
	c.txValidator = validator
}

// foreignTx returns a transaction signed for a different chain: well formed,
// but something no peer on this network should be asked to relay.
func foreignTx(t *testing.T) *types.Transaction {
	t.Helper()

	return signTx(t, testChainID+1, dynamicFeeTxFor(t, testChainID+1))
}

func TestFilterBroadcastableTxsDisabled(t *testing.T) {
	conns := sharedTestConns(t, false)
	withBroadcastHead(t, conns, newTestValidator(t, TxValidatorOptions{}))
	conns.shouldBroadcastTx = false

	txs, hashes := conns.FilterBroadcastableTxs([]*types.Transaction{signTx(t, testChainID, dynamicFeeTx(t))})
	if txs != nil || hashes != nil {
		t.Fatalf("broadcast disabled: want no txs, got %d txs and %d hashes", len(txs), len(hashes))
	}
}

func TestFilterBroadcastableTxsNoValidator(t *testing.T) {
	conns := sharedTestConns(t, false)
	withBroadcastHead(t, conns, nil)

	// A transaction signed for another chain: without a validator it is still
	// forwarded, which is the behavior the validator exists to change.
	tx := foreignTx(t)

	txs, hashes := conns.FilterBroadcastableTxs([]*types.Transaction{tx})
	if len(txs) != 1 || len(hashes) != 1 {
		t.Fatalf("want 1 tx and 1 hash, got %d and %d", len(txs), len(hashes))
	}
	if hashes[0] != tx.Hash() {
		t.Fatalf("want hash %v, got %v", tx.Hash(), hashes[0])
	}
}

func TestFilterBroadcastableTxsDropsInvalid(t *testing.T) {
	conns := sharedTestConns(t, false)
	withBroadcastHead(t, conns, newTestValidator(t, TxValidatorOptions{MinBaseFeeRatio: 1}))

	valid := signTx(t, testChainID, dynamicFeeTx(t))
	forged := foreignTx(t)

	before := testutil.ToFloat64(conns.metrics.txsRejected.WithLabelValues("invalid_signature"))

	txs, hashes := conns.FilterBroadcastableTxs([]*types.Transaction{forged, valid})
	if len(txs) != 1 || len(hashes) != 1 {
		t.Fatalf("want 1 tx and 1 hash, got %d and %d", len(txs), len(hashes))
	}
	if txs[0].Hash() != valid.Hash() || hashes[0] != valid.Hash() {
		t.Fatalf("want the valid tx %v, got tx %v and hash %v", valid.Hash(), txs[0].Hash(), hashes[0])
	}

	if got := testutil.ToFloat64(conns.metrics.txsRejected.WithLabelValues("invalid_signature")) - before; got != 1 {
		t.Fatalf("want 1 rejection recorded, got %v", got)
	}
}

// TestFilterBroadcastableTxsNoHead covers the startup window before a head is
// known: the sensor forwards rather than stalling broadcast entirely.
func TestFilterBroadcastableTxsNoHead(t *testing.T) {
	conns := sharedTestConns(t, false)
	withBroadcastHead(t, conns, newTestValidator(t, TxValidatorOptions{}))
	conns.head.Set(NewBlockPacket{})

	txs, hashes := conns.FilterBroadcastableTxs([]*types.Transaction{foreignTx(t)})
	if len(txs) != 1 || len(hashes) != 1 {
		t.Fatalf("want 1 tx and 1 hash, got %d and %d", len(txs), len(hashes))
	}
}

func TestFilterBroadcastableTxsEmpty(t *testing.T) {
	conns := sharedTestConns(t, false)
	withBroadcastHead(t, conns, newTestValidator(t, TxValidatorOptions{}))

	txs, hashes := conns.FilterBroadcastableTxs(nil)
	if txs != nil || hashes != nil {
		t.Fatalf("want no txs, got %d txs and %d hashes", len(txs), len(hashes))
	}
}
