package sensor

import (
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/0xPolygon/polygon-cli/p2p/bortx"
)

// TestFormatBodyTxResponseCoversStateSyncTx checks that every body index of
// mainnet block 97,029,696 is answerable, including the Bor state-sync
// transaction at index 48, so the hash list, the transaction count and the
// per-transaction lookups agree.
func TestFormatBodyTxResponseCoversStateSyncTx(t *testing.T) {
	data, err := os.ReadFile("../../../p2p/bortx/testdata/block_97029696.rlp.hex")
	if err != nil {
		t.Fatal(err)
	}
	enc, err := hex.DecodeString(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatal(err)
	}
	var block struct {
		Header *types.Header
		Txs    []rlp.RawValue
		Uncles []*types.Header
	}
	if err := rlp.DecodeBytes(enc, &block); err != nil {
		t.Fatal(err)
	}
	body, err := bortx.Decode(block.Txs)
	if err != nil {
		t.Fatal(err)
	}
	blockHash := block.Header.Hash()

	for i, txHash := range body.Hashes {
		resp := formatBodyTxResponse(body, i, blockHash, block.Header)
		if resp["hash"] != txHash.Hex() {
			t.Fatalf("index %d hash = %v, want %s", i, resp["hash"], txHash.Hex())
		}
		if resp["transactionIndex"] != hexutil.EncodeUint64(uint64(i)) {
			t.Fatalf("index %d transactionIndex = %v", i, resp["transactionIndex"])
		}
	}

	stateSync := formatBodyTxResponse(body, len(body.Hashes)-1, blockHash, block.Header)
	if stateSync["type"] != "0x7f" {
		t.Fatalf("state-sync type = %v, want 0x7f", stateSync["type"])
	}
	if stateSync["hash"] != "0xb6f56cd70a8d7a3090a46ecee86525410fe2eff3ff4db379cca7332f56a2e218" {
		t.Fatalf("state-sync hash = %v", stateSync["hash"])
	}
}
