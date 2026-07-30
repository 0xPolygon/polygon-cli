package p2p

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// TestDecodeTxSurvivesHostileInput is the regression test for a remote process kill.
//
// decodeTx logged bytes[0] after rlp.DecodeBytes succeeded. 0x80 is a valid RLP
// empty string, so DecodeBytes returns an empty slice with no error and the index
// panicked -- on the peer's own Protocol.Run goroutine, which geth does not recover
// around, from four unauthenticated message paths (TransactionsMsg,
// PooledTransactionsMsg, NewBlockMsg, BlockBodiesMsg). A three-byte message killed
// the sensor.
//
// This is the same failure class as the clique panic util.Ecrecover now converts to
// an error; that fix's comment names the mechanism exactly.
func TestDecodeTxSurvivesHostileInput(t *testing.T) {
	c := &conn{}

	for _, tc := range []struct {
		name string
		raw  []byte
	}{
		{"empty_rlp_string", []byte{0x80}},          // decodes to []byte{} with no error
		{"empty_rlp_list", []byte{0xc0}},            // decodes to a list, not bytes
		{"single_zero_byte", []byte{0x00}},          // valid RLP, one zero byte
		{"truncated_typed_tx", []byte{0x02}},        // EIP-1559 prefix, nothing after
		{"long_string_header_only", []byte{0xb8}},   // claims a length, supplies none
		{"nested_empty", []byte{0xc1, 0x80}},        // list containing an empty string
		{"garbage", []byte{0xff, 0xff, 0xff, 0xff}}, // not valid RLP at all
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Must return nil rather than panic; a peer controls this input entirely.
			if tx := c.decodeTx(tc.raw); tx != nil {
				t.Fatalf("unexpectedly decoded a transaction from %x", tc.raw)
			}
		})
	}
}

// The list paths must survive the same input, since that is how a peer actually
// delivers it: a Transactions packet is a list of these blobs.
func TestDecodeTxsSurvivesHostileInput(t *testing.T) {
	c := &conn{}
	raws := []rlp.RawValue{{0x80}, {0xc0}, {0x02}, {0xff}}
	if got := c.decodeTxs(raws); len(got) != 0 {
		t.Fatalf("expected nothing to decode, got %d transactions", len(got))
	}
}

// TestBuildBlockBodyRejectsUndecodableTx is the regression test for silent block
// corruption.
//
// buildBlockBody re-encoded whatever survived the lenient decodeTxs, so one garbage
// blob in a peer's BlockBodies response produced a body that was not the body for
// that hash: block_bodies got a tx_count short by the drops, and writeBlockTxs
// indexes by loop position so every later tx_index shifted down. Keyed by a real
// block hash, in ReplacingMergeTree tables with no version column, that row could
// win the merge against an honest sensor's and persist.
func TestBuildBlockBodyRejectsUndecodableTx(t *testing.T) {
	c := &conn{}

	good, err := types.SignTx(
		types.NewTx(&types.LegacyTx{Nonce: 1, GasPrice: big.NewInt(1), Gas: 21_000, Value: big.NewInt(1)}),
		types.LatestSignerForChainID(big.NewInt(137)), mustKey(t))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	goodRLP, err := good.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	encGood, err := rlp.EncodeToBytes(goodRLP)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	// A body whose first transaction is fine and whose second is garbage.
	body := struct {
		Transactions []rlp.RawValue
		Uncles       []*types.Header
	}{
		Transactions: []rlp.RawValue{encGood, {0x80}},
	}
	raw, err := rlp.EncodeToBytes(&body)
	if err != nil {
		t.Fatalf("encode body: %v", err)
	}

	if _, err := c.buildBlockBody(raw); err == nil {
		t.Fatal("accepted a body with an undecodable transaction; the tx_index mapping would be silently shifted")
	}

	// The all-good case must still work, or the guard is useless.
	body.Transactions = []rlp.RawValue{encGood}
	raw, err = rlp.EncodeToBytes(&body)
	if err != nil {
		t.Fatalf("encode body: %v", err)
	}
	built, err := c.buildBlockBody(raw)
	if err != nil {
		t.Fatalf("rejected a valid body: %v", err)
	}
	items, err := built.Transactions.Items()
	if err != nil {
		t.Fatalf("items: %v", err)
	}
	if len(items) != 1 || items[0].Hash() != good.Hash() {
		t.Fatalf("round trip lost the transaction: %d items", len(items))
	}
}

func mustKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	return k
}
