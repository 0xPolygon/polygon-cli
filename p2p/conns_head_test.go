package p2p

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	ds "github.com/0xPolygon/polygon-cli/p2p/datastructures"
)

// signedHeader returns a header sealed by key, the way a Bor block carries its
// signer in the last 65 bytes of Extra.
func signedHeader(t *testing.T, key *ecdsa.PrivateKey, number int64, baseFee *big.Int) *types.Header {
	t.Helper()

	header := &types.Header{
		Number:     big.NewInt(number),
		Difficulty: big.NewInt(1),
		GasLimit:   30_000_000,
		BaseFee:    baseFee,
		Extra:      make([]byte, crypto.SignatureLength),
	}
	sig, err := crypto.Sign(clique.SealHash(header).Bytes(), key)
	if err != nil {
		t.Fatalf("signing header: %v", err)
	}
	copy(header.Extra[len(header.Extra)-crypto.SignatureLength:], sig)
	return header
}

func headPacket(header *types.Header, td int64) NewBlockPacket {
	return NewBlockPacket{Block: types.NewBlockWithHeader(header), TD: big.NewInt(td)}
}

// headConns returns a Conns with an initialised head and no validator set,
// built directly rather than through NewConns, which registers Prometheus
// collectors and can only run once per test binary.
func headConns() *Conns {
	return &Conns{head: &ds.Locked[NewBlockPacket]{}}
}

// validatorConns returns a Conns whose validator set contains exactly signers.
func validatorConns(signers ...common.Address) *Conns {
	c := headConns()
	c.validators = &ValidatorSet{}
	set := make(map[common.Address]struct{}, len(signers))
	for _, s := range signers {
		set[s] = struct{}{}
	}
	c.validators.signers.Set(set)
	return c
}

func TestUpdateHeadBlockWithoutValidatorSet(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	// No validator set: RecoverSigner reports everything as known, so the head
	// tracks whatever arrives, exactly as it did before the signer gate.
	c := headConns()
	if !c.UpdateHeadBlock(headPacket(signedHeader(t, key, 100, big.NewInt(1)), 100)) {
		t.Fatal("first block: want head updated")
	}
	if !c.UpdateHeadBlock(headPacket(&types.Header{Number: big.NewInt(101), Difficulty: big.NewInt(1)}, 101)) {
		t.Fatal("unsigned block with no validator set: want head updated")
	}
	if got := c.HeadBlock().Block.NumberU64(); got != 101 {
		t.Fatalf("want head 101, got %d", got)
	}
}

// TestUpdateHeadBlockRejectsUnknownSigner is the regression test for a head
// poisoned by a peer. The head is read by the transaction validator, so a block
// a peer simply declared -- with an outlandish base fee and an unbeatable total
// difficulty -- would have rejected every honest transaction from then on, and
// no later block could have taken the head back.
func TestUpdateHeadBlockRejectsUnknownSigner(t *testing.T) {
	validator, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	attacker, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	c := validatorConns(crypto.PubkeyToAddress(validator.PublicKey))

	honest := signedHeader(t, validator, 100, big.NewInt(25_000_000_000))
	if !c.UpdateHeadBlock(headPacket(honest, 100)) {
		t.Fatal("validator-signed block: want head updated")
	}

	poison := signedHeader(t, attacker, 101, new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil))
	if c.UpdateHeadBlock(headPacket(poison, 1<<62)) {
		t.Fatal("block from a signer outside the validator set: want head unchanged")
	}

	head := c.HeadBlock().Block
	if head.NumberU64() != 100 {
		t.Fatalf("want head still 100, got %d", head.NumberU64())
	}
	if head.BaseFee().Cmp(big.NewInt(25_000_000_000)) != 0 {
		t.Fatalf("want head base fee untouched, got %v", head.BaseFee())
	}

	// A later validator-signed block still advances it, so the gate rejects the
	// block rather than wedging the head.
	next := signedHeader(t, validator, 102, big.NewInt(26_000_000_000))
	if !c.UpdateHeadBlock(headPacket(next, 102)) {
		t.Fatal("later validator-signed block: want head updated")
	}
	if got := c.HeadBlock().Block.NumberU64(); got != 102 {
		t.Fatalf("want head 102, got %d", got)
	}
}

// TestUpdateHeadBlockUnrecoverableSigner covers a header whose Extra is too
// short to hold a signature: Ecrecover fails, and a failure to recover must not
// read as a pass.
func TestUpdateHeadBlockUnrecoverableSigner(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	c := validatorConns(crypto.PubkeyToAddress(key.PublicKey))
	if c.UpdateHeadBlock(headPacket(&types.Header{Number: big.NewInt(1), Difficulty: big.NewInt(1)}, 1)) {
		t.Fatal("header with no recoverable signature: want head unchanged")
	}
	if c.HeadBlock().Block != nil {
		t.Fatal("want head still unset")
	}
}

func TestUpdateHeadBlockIgnoresOlderAndNil(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	c := validatorConns(crypto.PubkeyToAddress(key.PublicKey))
	if !c.UpdateHeadBlock(headPacket(signedHeader(t, key, 100, big.NewInt(1)), 100)) {
		t.Fatal("first block: want head updated")
	}

	tests := []struct {
		name   string
		packet NewBlockPacket
	}{
		{"lower number", headPacket(signedHeader(t, key, 99, big.NewInt(1)), 200)},
		{"same number", headPacket(signedHeader(t, key, 100, big.NewInt(1)), 200)},
		{"higher number, lower td", headPacket(signedHeader(t, key, 101, big.NewInt(1)), 99)},
		{"nil block", NewBlockPacket{TD: big.NewInt(1 << 62)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if c.UpdateHeadBlock(tt.packet) {
				t.Fatal("want head unchanged")
			}
			if got := c.HeadBlock().Block.NumberU64(); got != 100 {
				t.Fatalf("want head still 100, got %d", got)
			}
		})
	}
}
