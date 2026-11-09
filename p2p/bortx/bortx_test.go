package bortx

import (
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// Mainnet block 97,029,696, fetched with debug_getRawBlock. Its last
// transaction (index 48) is a Bor state-sync transaction.
const (
	fixtureBlock       = "testdata/block_97029696.rlp.hex"
	fixtureBlockHash   = "0x13e4fb76b91bf5e371177f4a7d475e341ec8f117f0125bb233fa6b9da26c84f7"
	fixtureTxRoot      = "0xa103efc5bcd5db9db456df27b3150c77ee5cf723e0019524902c5c29c009e009"
	fixtureStateSyncTx = "0xb6f56cd70a8d7a3090a46ecee86525410fe2eff3ff4db379cca7332f56a2e218"
	fixtureTxCount     = 49
)

type rawBlock struct {
	Header *types.Header
	Txs    []rlp.RawValue
	Uncles []*types.Header
}

func loadFixture(t *testing.T) rawBlock {
	t.Helper()
	data, err := os.ReadFile(fixtureBlock)
	if err != nil {
		t.Fatal(err)
	}
	enc, err := hex.DecodeString(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatal(err)
	}
	var block rawBlock
	if err := rlp.DecodeBytes(enc, &block); err != nil {
		t.Fatal(err)
	}
	if got := block.Header.Hash().Hex(); got != fixtureBlockHash {
		t.Fatalf("fixture header hash = %s, want %s", got, fixtureBlockHash)
	}
	return block
}

func TestDecodeMainnetBlockWithStateSyncTx(t *testing.T) {
	block := loadFixture(t)

	body, err := Decode(block.Txs)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if len(body.Hashes) != fixtureTxCount {
		t.Fatalf("hashes = %d, want %d", len(body.Hashes), fixtureTxCount)
	}
	if len(body.Txs) != fixtureTxCount-1 || body.StateSyncCount() != 1 {
		t.Fatalf("txs = %d, state-sync = %d, want %d and 1", len(body.Txs), body.StateSyncCount(), fixtureTxCount-1)
	}
	if got := body.Hashes[fixtureTxCount-1].Hex(); got != fixtureStateSyncTx {
		t.Errorf("state-sync hash = %s, want %s", got, fixtureStateSyncTx)
	}
	for i, tx := range body.Txs {
		if tx.Hash() != body.Hashes[i] {
			t.Errorf("tx %d hash %s does not match body hash %s", i, tx.Hash().Hex(), body.Hashes[i].Hex())
		}
	}
	if got := body.Root().Hex(); got != fixtureTxRoot {
		t.Errorf("root = %s, want %s", got, fixtureTxRoot)
	}
	if body.Root() != block.Header.TxHash {
		t.Errorf("root %s does not match header TxHash %s", body.Root().Hex(), block.Header.TxHash.Hex())
	}
}

func TestDecodeListMatchesDecode(t *testing.T) {
	block := loadFixture(t)

	var list rlp.RawList[*types.Transaction]
	for _, item := range block.Txs {
		if err := list.AppendRaw(item); err != nil {
			t.Fatal(err)
		}
	}
	body, err := DecodeList(&list)
	if err != nil {
		t.Fatalf("DecodeList: %v", err)
	}
	if len(body.Hashes) != fixtureTxCount || body.StateSyncCount() != 1 {
		t.Fatalf("hashes = %d, state-sync = %d, want %d and 1", len(body.Hashes), body.StateSyncCount(), fixtureTxCount)
	}
	if body.Root() != block.Header.TxHash {
		t.Errorf("root %s does not match header TxHash %s", body.Root().Hex(), block.Header.TxHash.Hex())
	}
}

func TestDecodeRejectsUnknownTxType(t *testing.T) {
	block := loadFixture(t)

	items := append([]rlp.RawValue{}, block.Txs...)
	bogus, err := rlp.EncodeToBytes([]byte{0x7e, 0xc0})
	if err != nil {
		t.Fatal(err)
	}
	items[3] = bogus
	if _, err := Decode(items); err == nil {
		t.Fatal("Decode accepted a body with an unknown transaction type")
	}
}

func TestIsStateSync(t *testing.T) {
	valid, err := rlp.EncodeToBytes([]stateSyncData{{ID: 1, Contract: common.Address{1}, Data: []byte{2}, TxHash: common.Hash{3}}})
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name   string
		binary []byte
		want   bool
	}{
		{"valid", append([]byte{StateSyncTxType}, valid...), true},
		{"empty", nil, false},
		{"wrong type", append([]byte{0x02}, valid...), false},
		{"malformed payload", []byte{StateSyncTxType, 0xc5, 0x01}, false},
	}
	for _, tc := range cases {
		if got := IsStateSync(tc.binary); got != tc.want {
			t.Errorf("%s: IsStateSync = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestTxByBodyIndex(t *testing.T) {
	block := loadFixture(t)
	body, err := Decode(block.Txs)
	if err != nil {
		t.Fatal(err)
	}
	if tx := body.Tx(fixtureTxCount - 1); tx != nil {
		t.Fatalf("Tx(%d) = %s, want nil for the state-sync transaction", fixtureTxCount-1, tx.Hash().Hex())
	}
	for i := 0; i < fixtureTxCount-1; i++ {
		if tx := body.Tx(i); tx == nil || tx.Hash() != body.Hashes[i] {
			t.Fatalf("Tx(%d) does not match body hash %s", i, body.Hashes[i].Hex())
		}
	}
}

func TestDecodeErrorDescribesTheFailingTx(t *testing.T) {
	block := loadFixture(t)
	items := append([]rlp.RawValue{}, block.Txs...)
	binary := []byte{0x7e, 0xc0}
	bogus, err := rlp.EncodeToBytes(binary)
	if err != nil {
		t.Fatal(err)
	}
	items[3] = bogus

	_, err = Decode(items)
	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("Decode error %v is not a *DecodeError", err)
	}
	if decodeErr.Index != 3 || decodeErr.Count != len(items) || decodeErr.Type != 0x7e || decodeErr.Size != len(binary) {
		t.Fatalf("DecodeError = %+v, want index 3, type 0x7e, size %d", decodeErr, len(binary))
	}
	if decodeErr.Hash != crypto.Keccak256Hash(binary) {
		t.Fatalf("DecodeError hash = %s, want %s", decodeErr.Hash.Hex(), crypto.Keccak256Hash(binary).Hex())
	}
}
