package proto

import (
	"bytes"
	"reflect"
	"testing"
)

// TestCheckpointRoundTrip exercises a full marshal/unmarshal cycle for
// every Msg added in W4. The goal is not to assert an exact byte
// pattern (that would require pinning on the Go protobuf encoder's
// output) but to ensure each field survives the round trip intact.
func TestCheckpointRoundTrip(t *testing.T) {
	in := &MsgCheckpoint{
		Proposer:        "0xabc",
		StartBlock:      100,
		EndBlock:        200,
		RootHash:        []byte{0x01, 0x02, 0x03},
		AccountRootHash: []byte{0x0a, 0x0b},
		BorChainID:      "137",
	}
	raw := in.Marshal()
	got, err := UnmarshalMsgCheckpoint(raw)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(in, got) {
		t.Errorf("round trip diverged:\nin=%+v\ngot=%+v", in, got)
	}
}

func TestCpAckRoundTrip(t *testing.T) {
	in := &MsgCpAck{From: "0x1", Number: 42, Proposer: "0x2", StartBlock: 10, EndBlock: 20, RootHash: []byte{0xff}}
	got, err := UnmarshalMsgCpAck(in.Marshal())
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(in, got) {
		t.Errorf("mismatch: %+v vs %+v", in, got)
	}
}

func TestCpNoAckRoundTrip(t *testing.T) {
	in := &MsgCpNoAck{From: "0xabc"}
	got, err := UnmarshalMsgCpNoAck(in.Marshal())
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if in.From != got.From {
		t.Errorf("mismatch: %q vs %q", in.From, got.From)
	}
}

func TestProposeSpanRoundTrip(t *testing.T) {
	in := &MsgProposeSpan{
		SpanID: 3, Proposer: "0xaaa", StartBlock: 1, EndBlock: 2,
		ChainID: "137", Seed: []byte{1, 2, 3, 4}, SeedAuthor: "0xbbb",
	}
	got, err := UnmarshalMsgProposeSpan(in.Marshal())
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(in, got) {
		t.Errorf("mismatch")
	}
}

func TestBackfillSpansRoundTrip(t *testing.T) {
	in := &MsgBackfillSpans{Proposer: "x", ChainID: "1", LatestSpanID: 5, LatestBorSpanID: 6}
	got, err := UnmarshalMsgBackfillSpans(in.Marshal())
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(in, got) {
		t.Errorf("mismatch")
	}
}

func TestVoteProducersRoundTripPacked(t *testing.T) {
	in := &MsgVoteProducers{
		Voter: "0x1", VoterID: 42,
		Votes: ProducerVotes{Votes: []uint64{1, 2, 3, 4, 5}},
	}
	got, err := UnmarshalMsgVoteProducers(in.Marshal())
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if in.Voter != got.Voter || in.VoterID != got.VoterID {
		t.Fatalf("scalar mismatch: %+v vs %+v", in, got)
	}
	if !reflect.DeepEqual(in.Votes.Votes, got.Votes.Votes) {
		t.Errorf("votes mismatch: %v vs %v", in.Votes.Votes, got.Votes.Votes)
	}
}

func TestSetProducerDowntimeRoundTrip(t *testing.T) {
	in := &MsgSetProducerDowntime{
		Producer: "0xaaa",
		DowntimeRange: BlockRange{StartBlock: 100, EndBlock: 200},
	}
	got, err := UnmarshalMsgSetProducerDowntime(in.Marshal())
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(in, got) {
		t.Errorf("mismatch: %+v vs %+v", in, got)
	}
}

func TestTopupRoundTrip(t *testing.T) {
	in := &MsgTopupTx{
		Proposer: "0x1", User: "0x2", Fee: "1000",
		TxHash: []byte{0xde, 0xad, 0xbe, 0xef},
		LogIndex: 1, BlockNumber: 100,
	}
	got, err := UnmarshalMsgTopupTx(in.Marshal())
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(in, got) {
		t.Errorf("mismatch")
	}
}

func TestValidatorJoinRoundTrip(t *testing.T) {
	in := &MsgValidatorJoin{
		From: "0x1", ValID: 42, ActivationEpoch: 1, Amount: "1000",
		SignerPubKey: []byte{0x04, 0x01, 0x02},
		TxHash: []byte{0xbe, 0xef},
		LogIndex: 1, BlockNumber: 100, Nonce: 5,
	}
	got, err := UnmarshalMsgValidatorJoin(in.Marshal())
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(in, got) {
		t.Errorf("mismatch")
	}
}

func TestStakeUpdateRoundTrip(t *testing.T) {
	in := &MsgStakeUpdate{From: "x", ValID: 1, NewAmount: "10", TxHash: []byte{1}, LogIndex: 1, BlockNumber: 1, Nonce: 1}
	got, err := UnmarshalMsgStakeUpdate(in.Marshal())
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(in, got) {
		t.Errorf("mismatch")
	}
}

func TestSignerUpdateRoundTrip(t *testing.T) {
	in := &MsgSignerUpdate{From: "x", ValID: 1, NewSignerPubKey: []byte{0x04, 0xff}, TxHash: []byte{1}, LogIndex: 1, BlockNumber: 1, Nonce: 1}
	got, err := UnmarshalMsgSignerUpdate(in.Marshal())
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(in, got) {
		t.Errorf("mismatch")
	}
}

func TestValidatorExitRoundTrip(t *testing.T) {
	in := &MsgValidatorExit{From: "x", ValID: 1, DeactivationEpoch: 99, TxHash: []byte{1}, LogIndex: 1, BlockNumber: 1, Nonce: 1}
	got, err := UnmarshalMsgValidatorExit(in.Marshal())
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(in, got) {
		t.Errorf("mismatch")
	}
}

func TestEventRecordRoundTrip(t *testing.T) {
	in := &MsgEventRecord{
		From: "0x1", TxHash: "0xabc", LogIndex: 1, BlockNumber: 100,
		ContractAddress: "0x2", Data: []byte{0xde, 0xad}, ID: 99, ChainID: "1",
	}
	got, err := UnmarshalMsgEventRecord(in.Marshal())
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(in, got) {
		t.Errorf("mismatch")
	}
}

func TestVoteExtensionRoundTrip(t *testing.T) {
	in := &VoteExtension{
		BlockHash: []byte{0x01, 0x02, 0x03},
		Height:    42,
		SideTxResponses: []SideTxResponse{
			{TxHash: []byte{0xde, 0xad}, Result: VoteYes},
			{TxHash: []byte{0xbe, 0xef}, Result: VoteNo},
		},
		MilestoneProposition: &MilestoneProposition{
			BlockHashes:      [][]byte{{0x0a}, {0x0b}},
			StartBlockNumber: 100,
			ParentHash:       []byte{0xff},
			BlockTDs:         []uint64{10, 20, 30},
		},
	}
	raw := in.Marshal()
	got, err := UnmarshalVoteExtension(raw)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !bytes.Equal(in.BlockHash, got.BlockHash) || in.Height != got.Height {
		t.Fatalf("scalars diverge: %+v vs %+v", in, got)
	}
	if len(in.SideTxResponses) != len(got.SideTxResponses) {
		t.Fatalf("side tx count: %d vs %d", len(in.SideTxResponses), len(got.SideTxResponses))
	}
	if in.MilestoneProposition != nil {
		if got.MilestoneProposition == nil {
			t.Fatalf("milestone proposition lost")
		}
		if !reflect.DeepEqual(in.MilestoneProposition.BlockTDs, got.MilestoneProposition.BlockTDs) {
			t.Errorf("block TDs mismatch: %v vs %v",
				in.MilestoneProposition.BlockTDs, got.MilestoneProposition.BlockTDs)
		}
	}
}

// TestRegistryCoversAllMsgs asserts that every Msg exposes its type URL
// via the registry. This catches the common mistake of adding a new
// Msg's proto file without also updating registry.go.
func TestRegistryCoversAllMsgs(t *testing.T) {
	want := []string{
		MsgWithdrawFeeTxTypeURL,
		MsgTopupTxTypeURL,
		MsgCheckpointTypeURL,
		MsgCpAckTypeURL,
		MsgCpNoAckTypeURL,
		MsgProposeSpanTypeURL,
		MsgBackfillSpansTypeURL,
		MsgVoteProducersTypeURL,
		MsgSetProducerDowntimeTypeURL,
		MsgValidatorJoinTypeURL,
		MsgStakeUpdateTypeURL,
		MsgSignerUpdateTypeURL,
		MsgValidatorExitTypeURL,
		MsgEventRecordTypeURL,
	}
	got := KnownTypeURLs()
	have := map[string]bool{}
	for _, u := range got {
		have[u] = true
	}
	for _, u := range want {
		if !have[u] {
			t.Errorf("registry missing type URL %s", u)
		}
	}
}

// TestRegistryDecodesViaLookup hits the full Decode path (lookup +
// unmarshal) for a representative Msg in each module. A mistake in
// registry wiring would surface here as a type assertion / nil pointer
// failure on the returned value.
func TestRegistryDecodesViaLookup(t *testing.T) {
	cases := []struct {
		typeURL string
		raw     []byte
	}{
		{MsgWithdrawFeeTxTypeURL, (&MsgWithdrawFeeTx{Proposer: "x", Amount: "1"}).Marshal()},
		{MsgCheckpointTypeURL, (&MsgCheckpoint{Proposer: "x", BorChainID: "1"}).Marshal()},
		{MsgProposeSpanTypeURL, (&MsgProposeSpan{Proposer: "x", ChainID: "1"}).Marshal()},
		{MsgValidatorJoinTypeURL, (&MsgValidatorJoin{From: "x"}).Marshal()},
		{MsgEventRecordTypeURL, (&MsgEventRecord{From: "x"}).Marshal()},
	}
	for _, c := range cases {
		t.Run(c.typeURL, func(t *testing.T) {
			v, err := Decode(c.typeURL, c.raw)
			if err != nil {
				t.Fatalf("Decode: %v", err)
			}
			if v == nil {
				t.Fatalf("Decode returned nil value")
			}
		})
	}
}
