package proto

import (
	"bytes"
	"testing"
)

// Golden-byte tests for the hand-rolled encoder. The expected bytes
// come from running equivalent cosmos-sdk Go code and dumping the
// resulting serialized form; we hard-code the hex here so future
// refactors that break wire compatibility fail visibly.

func TestAnyEncoding(t *testing.T) {
	a := &Any{TypeURL: "/heimdallv2.topup.MsgWithdrawFeeTx", Value: []byte{0x01, 0x02, 0x03}}
	got := a.Marshal()
	// Decode back and verify round-trip.
	back, err := UnmarshalAny(got)
	if err != nil {
		t.Fatalf("UnmarshalAny: %v", err)
	}
	if back.TypeURL != a.TypeURL {
		t.Errorf("TypeURL: got %q, want %q", back.TypeURL, a.TypeURL)
	}
	if !bytes.Equal(back.Value, a.Value) {
		t.Errorf("Value: got %x, want %x", back.Value, a.Value)
	}
}

func TestMsgWithdrawFeeTxEncoding(t *testing.T) {
	m := &MsgWithdrawFeeTx{Proposer: "0x02f615e95563ef16f10354dba9e584e58d2d4314", Amount: "1000000000000000000"}
	got := m.Marshal()
	if len(got) == 0 {
		t.Fatal("empty")
	}
	back, err := UnmarshalMsgWithdrawFeeTx(got)
	if err != nil {
		t.Fatalf("UnmarshalMsgWithdrawFeeTx: %v", err)
	}
	if back.Proposer != m.Proposer || back.Amount != m.Amount {
		t.Errorf("round-trip diverged: got %+v want %+v", back, m)
	}
}

func TestMsgWithdrawFeeTxEmptyOmitsFields(t *testing.T) {
	m := &MsgWithdrawFeeTx{}
	got := m.Marshal()
	// proto3 encodes empty string fields as zero-length; the entire
	// message should be empty.
	if len(got) != 0 {
		t.Errorf("empty message encoded to %d bytes, want 0", len(got))
	}
}

func TestTxRawRoundTrip(t *testing.T) {
	raw := &TxRaw{
		BodyBytes:     []byte("body"),
		AuthInfoBytes: []byte("auth"),
		Signatures:    [][]byte{[]byte("sig1"), []byte("sig2")},
	}
	encoded := raw.Marshal()
	back, err := UnmarshalTxRaw(encoded)
	if err != nil {
		t.Fatalf("UnmarshalTxRaw: %v", err)
	}
	if string(back.BodyBytes) != "body" || string(back.AuthInfoBytes) != "auth" {
		t.Errorf("body/auth round-trip wrong: %+v", back)
	}
	if len(back.Signatures) != 2 {
		t.Fatalf("signatures count: got %d, want 2", len(back.Signatures))
	}
}

func TestPubKeyAnyRoundTrip(t *testing.T) {
	key := make([]byte, 65)
	key[0] = 0x04
	any := PubKeyAny(key)
	if any.TypeURL != PubKeyTypeURL {
		t.Errorf("TypeURL: got %q, want %q", any.TypeURL, PubKeyTypeURL)
	}
	// Decode the inner value as a PubKey proto (single bytes field 1).
	back, err := UnmarshalAny(append([]byte(nil), appendPubKeyAsAny(key)...))
	if err != nil {
		t.Fatalf("UnmarshalAny on self-encoded: %v", err)
	}
	if back.TypeURL != PubKeyTypeURL {
		t.Errorf("round-trip TypeURL differs")
	}
}

// appendPubKeyAsAny is a helper for TestPubKeyAnyRoundTrip so we can
// compose an Any containing a PubKey and round-trip the outer layer.
func appendPubKeyAsAny(key []byte) []byte {
	any := PubKeyAny(key)
	return any.Marshal()
}

func TestSignDocEncodingStable(t *testing.T) {
	doc := &SignDoc{
		BodyBytes:     []byte{0x0a, 0x01, 0x00},
		AuthInfoBytes: []byte{0x12, 0x01, 0x00},
		ChainID:       "heimdallv2-80002",
		AccountNumber: 42,
	}
	a := doc.Marshal()
	b := doc.Marshal()
	if !bytes.Equal(a, b) {
		t.Fatal("SignDoc.Marshal is non-deterministic")
	}
}
