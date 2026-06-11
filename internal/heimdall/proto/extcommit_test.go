package proto

import (
	"bytes"
	"testing"

	"google.golang.org/protobuf/encoding/protowire"
)

func sampleVoteExtensionBytes(t *testing.T) []byte {
	t.Helper()
	ve := &VoteExtension{
		BlockHash: bytes.Repeat([]byte{0xaa}, 32),
		Height:    30001,
		MilestoneProposition: &MilestoneProposition{
			BlockHashes:      [][]byte{bytes.Repeat([]byte{0x01}, 32), bytes.Repeat([]byte{0x02}, 32)},
			StartBlockNumber: 72100001,
			ParentHash:       bytes.Repeat([]byte{0x03}, 32),
			BlockTDs:         []uint64{10, 11},
		},
	}
	return ve.Marshal()
}

func TestExtendedCommitInfoRoundTrip(t *testing.T) {
	in := &ExtendedCommitInfo{
		Round: 1,
		Votes: []ExtendedVoteInfo{
			{
				Validator:               ExtValidator{Address: bytes.Repeat([]byte{0x11}, 20), Power: 71302411},
				VoteExtension:           sampleVoteExtensionBytes(t),
				ExtensionSignature:      []byte("sig-a"),
				BlockIDFlag:             BlockIDFlagCommit,
				NonRpVoteExtension:      []byte("nonrp"),
				NonRpExtensionSignature: []byte("nonrp-sig"),
			},
			{
				// Absent validator: no extension, no signature.
				Validator:   ExtValidator{Address: bytes.Repeat([]byte{0x22}, 20), Power: 50000000},
				BlockIDFlag: BlockIDFlagAbsent,
			},
		},
	}
	back, err := UnmarshalExtendedCommitInfo(in.Marshal())
	if err != nil {
		t.Fatalf("UnmarshalExtendedCommitInfo: %v", err)
	}
	if back.Round != in.Round {
		t.Errorf("Round: got %d, want %d", back.Round, in.Round)
	}
	if len(back.Votes) != len(in.Votes) {
		t.Fatalf("Votes: got %d, want %d", len(back.Votes), len(in.Votes))
	}
	for i := range in.Votes {
		got, want := back.Votes[i], in.Votes[i]
		if !bytes.Equal(got.Validator.Address, want.Validator.Address) {
			t.Errorf("vote %d address: got %x, want %x", i, got.Validator.Address, want.Validator.Address)
		}
		if got.Validator.Power != want.Validator.Power {
			t.Errorf("vote %d power: got %d, want %d", i, got.Validator.Power, want.Validator.Power)
		}
		if !bytes.Equal(got.VoteExtension, want.VoteExtension) {
			t.Errorf("vote %d extension diverged", i)
		}
		if !bytes.Equal(got.ExtensionSignature, want.ExtensionSignature) {
			t.Errorf("vote %d signature diverged", i)
		}
		if got.BlockIDFlag != want.BlockIDFlag {
			t.Errorf("vote %d flag: got %v, want %v", i, got.BlockIDFlag, want.BlockIDFlag)
		}
		if !bytes.Equal(got.NonRpVoteExtension, want.NonRpVoteExtension) {
			t.Errorf("vote %d non-rp extension diverged", i)
		}
		if !bytes.Equal(got.NonRpExtensionSignature, want.NonRpExtensionSignature) {
			t.Errorf("vote %d non-rp signature diverged", i)
		}
	}

	// The committed vote's extension must decode as a VoteExtension with
	// the proposition intact.
	ve, err := UnmarshalVoteExtension(back.Votes[0].VoteExtension)
	if err != nil {
		t.Fatalf("UnmarshalVoteExtension: %v", err)
	}
	if ve.MilestoneProposition == nil || ve.MilestoneProposition.StartBlockNumber != 72100001 {
		t.Errorf("nested proposition diverged: %+v", ve.MilestoneProposition)
	}
}

func TestExtendedCommitInfoSkipsUnknownFields(t *testing.T) {
	vote := ExtendedVoteInfo{
		Validator:   ExtValidator{Address: bytes.Repeat([]byte{0x33}, 20), Power: 7},
		BlockIDFlag: BlockIDFlagCommit,
	}
	inner := vote.Marshal()
	// Splice in an unknown field (number 9, bytes) to simulate a future
	// fork extension.
	inner = protowire.AppendTag(inner, 9, protowire.BytesType)
	inner = protowire.AppendBytes(inner, []byte("future"))

	var raw []byte
	raw = protowire.AppendTag(raw, 2, protowire.BytesType)
	raw = protowire.AppendBytes(raw, inner)

	back, err := UnmarshalExtendedCommitInfo(raw)
	if err != nil {
		t.Fatalf("UnmarshalExtendedCommitInfo: %v", err)
	}
	if len(back.Votes) != 1 || back.Votes[0].Validator.Power != 7 || back.Votes[0].BlockIDFlag != BlockIDFlagCommit {
		t.Errorf("vote diverged after unknown field: %+v", back.Votes)
	}
}

func TestExtendedCommitInfoMalformed(t *testing.T) {
	// A bytes tag with a length that overruns the buffer.
	bad := []byte{0x12, 0xff, 0x01, 0x02}
	if _, err := UnmarshalExtendedCommitInfo(bad); err == nil {
		t.Fatal("expected error on malformed input")
	}
}

func TestBlockIDFlagString(t *testing.T) {
	cases := map[BlockIDFlag]string{
		BlockIDFlagUnknown: "UNKNOWN",
		BlockIDFlagAbsent:  "ABSENT",
		BlockIDFlagCommit:  "COMMIT",
		BlockIDFlagNil:     "NIL",
		BlockIDFlag(9):     "BlockIDFlag(9)",
	}
	for flag, want := range cases {
		if got := flag.String(); got != want {
			t.Errorf("BlockIDFlag(%d).String(): got %q, want %q", int32(flag), got, want)
		}
	}
}
