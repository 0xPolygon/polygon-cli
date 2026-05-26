package tx

import (
	"fmt"

	hproto "github.com/0xPolygon/polygon-cli/internal/heimdall/proto"
)

// CheckpointMsg wraps heimdallv2.checkpoint.MsgCheckpoint for the
// builder. All fields must be supplied by the caller; the builder
// performs no semantic validation beyond "proposer non-empty".
type CheckpointMsg struct {
	Proposer        string
	StartBlock      uint64
	EndBlock        uint64
	RootHash        []byte
	AccountRootHash []byte
	BorChainID      string
}

// TypeURL implements Msg.
func (m *CheckpointMsg) TypeURL() string { return hproto.MsgCheckpointTypeURL }

// Marshal implements Msg.
func (m *CheckpointMsg) Marshal() ([]byte, error) {
	if m.Proposer == "" {
		return nil, fmt.Errorf("CheckpointMsg: proposer is required")
	}
	p := &hproto.MsgCheckpoint{
		Proposer:        m.Proposer,
		StartBlock:      m.StartBlock,
		EndBlock:        m.EndBlock,
		RootHash:        m.RootHash,
		AccountRootHash: m.AccountRootHash,
		BorChainID:      m.BorChainID,
	}
	return p.Marshal(), nil
}

// AminoName implements Msg.
func (m *CheckpointMsg) AminoName() string { return "heimdallv2/checkpoint/MsgCheckpoint" }

// AminoJSON implements Msg.
func (m *CheckpointMsg) AminoJSON() (any, error) {
	return map[string]any{
		"proposer":          m.Proposer,
		"start_block":       fmt.Sprintf("%d", m.StartBlock),
		"end_block":         fmt.Sprintf("%d", m.EndBlock),
		"root_hash":         m.RootHash,
		"account_root_hash": m.AccountRootHash,
		"bor_chain_id":      m.BorChainID,
	}, nil
}

// CpAckMsg wraps MsgCpAck.
type CpAckMsg struct {
	From       string
	Number     uint64
	Proposer   string
	StartBlock uint64
	EndBlock   uint64
	RootHash   []byte
}

func (m *CpAckMsg) TypeURL() string { return hproto.MsgCpAckTypeURL }
func (m *CpAckMsg) Marshal() ([]byte, error) {
	if m.From == "" {
		return nil, fmt.Errorf("CpAckMsg: from is required")
	}
	p := &hproto.MsgCpAck{
		From: m.From, Number: m.Number, Proposer: m.Proposer,
		StartBlock: m.StartBlock, EndBlock: m.EndBlock, RootHash: m.RootHash,
	}
	return p.Marshal(), nil
}
func (m *CpAckMsg) AminoName() string { return "heimdallv2/checkpoint/MsgCpAck" }
func (m *CpAckMsg) AminoJSON() (any, error) {
	return map[string]any{
		"from":        m.From,
		"number":      fmt.Sprintf("%d", m.Number),
		"proposer":    m.Proposer,
		"start_block": fmt.Sprintf("%d", m.StartBlock),
		"end_block":   fmt.Sprintf("%d", m.EndBlock),
		"root_hash":   m.RootHash,
	}, nil
}

// CpNoAckMsg wraps MsgCpNoAck.
type CpNoAckMsg struct {
	From string
}

func (m *CpNoAckMsg) TypeURL() string { return hproto.MsgCpNoAckTypeURL }
func (m *CpNoAckMsg) Marshal() ([]byte, error) {
	if m.From == "" {
		return nil, fmt.Errorf("CpNoAckMsg: from is required")
	}
	p := &hproto.MsgCpNoAck{From: m.From}
	return p.Marshal(), nil
}
func (m *CpNoAckMsg) AminoName() string { return "heimdallv2/checkpoint/MsgCpNoAck" }
func (m *CpNoAckMsg) AminoJSON() (any, error) {
	return map[string]any{"from": m.From}, nil
}

// ProposeSpanMsg wraps MsgProposeSpan.
type ProposeSpanMsg struct {
	SpanID     uint64
	Proposer   string
	StartBlock uint64
	EndBlock   uint64
	ChainID    string
	Seed       []byte
	SeedAuthor string
}

func (m *ProposeSpanMsg) TypeURL() string { return hproto.MsgProposeSpanTypeURL }
func (m *ProposeSpanMsg) Marshal() ([]byte, error) {
	if m.Proposer == "" {
		return nil, fmt.Errorf("ProposeSpanMsg: proposer is required")
	}
	p := &hproto.MsgProposeSpan{
		SpanID: m.SpanID, Proposer: m.Proposer,
		StartBlock: m.StartBlock, EndBlock: m.EndBlock,
		ChainID: m.ChainID, Seed: m.Seed, SeedAuthor: m.SeedAuthor,
	}
	return p.Marshal(), nil
}
func (m *ProposeSpanMsg) AminoName() string { return "heimdallv2/bor/MsgProposeSpan" }
func (m *ProposeSpanMsg) AminoJSON() (any, error) {
	return map[string]any{
		"span_id":     fmt.Sprintf("%d", m.SpanID),
		"proposer":    m.Proposer,
		"start_block": fmt.Sprintf("%d", m.StartBlock),
		"end_block":   fmt.Sprintf("%d", m.EndBlock),
		"chain_id":    m.ChainID,
		"seed":        m.Seed,
		"seed_author": m.SeedAuthor,
	}, nil
}

// BackfillSpansMsg wraps MsgBackfillSpans.
type BackfillSpansMsg struct {
	Proposer        string
	ChainID         string
	LatestSpanID    uint64
	LatestBorSpanID uint64
}

func (m *BackfillSpansMsg) TypeURL() string { return hproto.MsgBackfillSpansTypeURL }
func (m *BackfillSpansMsg) Marshal() ([]byte, error) {
	if m.Proposer == "" {
		return nil, fmt.Errorf("BackfillSpansMsg: proposer is required")
	}
	p := &hproto.MsgBackfillSpans{
		Proposer: m.Proposer, ChainID: m.ChainID,
		LatestSpanID: m.LatestSpanID, LatestBorSpanID: m.LatestBorSpanID,
	}
	return p.Marshal(), nil
}
func (m *BackfillSpansMsg) AminoName() string { return "heimdallv2/bor/MsgBackfillSpans" }
func (m *BackfillSpansMsg) AminoJSON() (any, error) {
	return map[string]any{
		"proposer":           m.Proposer,
		"chain_id":           m.ChainID,
		"latest_span_id":     fmt.Sprintf("%d", m.LatestSpanID),
		"latest_bor_span_id": fmt.Sprintf("%d", m.LatestBorSpanID),
	}, nil
}

// VoteProducersMsg wraps MsgVoteProducers.
type VoteProducersMsg struct {
	Voter   string
	VoterID uint64
	Votes   []uint64
}

func (m *VoteProducersMsg) TypeURL() string { return hproto.MsgVoteProducersTypeURL }
func (m *VoteProducersMsg) Marshal() ([]byte, error) {
	if m.Voter == "" {
		return nil, fmt.Errorf("VoteProducersMsg: voter is required")
	}
	p := &hproto.MsgVoteProducers{
		Voter: m.Voter, VoterID: m.VoterID,
		Votes: hproto.ProducerVotes{Votes: m.Votes},
	}
	return p.Marshal(), nil
}
func (m *VoteProducersMsg) AminoName() string { return "heimdallv2/bor/MsgVoteProducers" }
func (m *VoteProducersMsg) AminoJSON() (any, error) {
	votes := make([]string, 0, len(m.Votes))
	for _, v := range m.Votes {
		votes = append(votes, fmt.Sprintf("%d", v))
	}
	return map[string]any{
		"voter":    m.Voter,
		"voter_id": fmt.Sprintf("%d", m.VoterID),
		"votes":    map[string]any{"votes": votes},
	}, nil
}

// SetProducerDowntimeMsg wraps MsgSetProducerDowntime.
type SetProducerDowntimeMsg struct {
	Producer   string
	StartBlock uint64
	EndBlock   uint64
}

func (m *SetProducerDowntimeMsg) TypeURL() string { return hproto.MsgSetProducerDowntimeTypeURL }
func (m *SetProducerDowntimeMsg) Marshal() ([]byte, error) {
	if m.Producer == "" {
		return nil, fmt.Errorf("SetProducerDowntimeMsg: producer is required")
	}
	p := &hproto.MsgSetProducerDowntime{
		Producer: m.Producer,
		DowntimeRange: hproto.BlockRange{
			StartBlock: m.StartBlock,
			EndBlock:   m.EndBlock,
		},
	}
	return p.Marshal(), nil
}
func (m *SetProducerDowntimeMsg) AminoName() string { return "heimdallv2/bor/MsgSetProducerDowntime" }
func (m *SetProducerDowntimeMsg) AminoJSON() (any, error) {
	return map[string]any{
		"producer": m.Producer,
		"downtime_range": map[string]any{
			"start_block": fmt.Sprintf("%d", m.StartBlock),
			"end_block":   fmt.Sprintf("%d", m.EndBlock),
		},
	}, nil
}

// TopupMsg wraps MsgTopupTx (L1-mirroring; gated by RequireForce).
type TopupMsg struct {
	Proposer    string
	User        string
	Fee         string
	TxHash      []byte
	LogIndex    uint64
	BlockNumber uint64
}

func (m *TopupMsg) TypeURL() string { return hproto.MsgTopupTxTypeURL }
func (m *TopupMsg) Marshal() ([]byte, error) {
	if m.Proposer == "" {
		return nil, fmt.Errorf("TopupMsg: proposer is required")
	}
	if m.Fee == "" {
		return nil, fmt.Errorf("TopupMsg: fee is required")
	}
	p := &hproto.MsgTopupTx{
		Proposer: m.Proposer, User: m.User, Fee: m.Fee,
		TxHash: m.TxHash, LogIndex: m.LogIndex, BlockNumber: m.BlockNumber,
	}
	return p.Marshal(), nil
}
func (m *TopupMsg) AminoName() string { return "heimdallv2/topup/MsgTopupTx" }
func (m *TopupMsg) AminoJSON() (any, error) {
	return map[string]any{
		"proposer":     m.Proposer,
		"user":         m.User,
		"fee":          m.Fee,
		"tx_hash":      m.TxHash,
		"log_index":    fmt.Sprintf("%d", m.LogIndex),
		"block_number": fmt.Sprintf("%d", m.BlockNumber),
	}, nil
}

// ValidatorJoinMsg wraps MsgValidatorJoin.
type ValidatorJoinMsg struct {
	From            string
	ValID           uint64
	ActivationEpoch uint64
	Amount          string
	SignerPubKey    []byte
	TxHash          []byte
	LogIndex        uint64
	BlockNumber     uint64
	Nonce           uint64
}

func (m *ValidatorJoinMsg) TypeURL() string { return hproto.MsgValidatorJoinTypeURL }
func (m *ValidatorJoinMsg) Marshal() ([]byte, error) {
	if m.From == "" {
		return nil, fmt.Errorf("ValidatorJoinMsg: from is required")
	}
	p := &hproto.MsgValidatorJoin{
		From: m.From, ValID: m.ValID, ActivationEpoch: m.ActivationEpoch,
		Amount: m.Amount, SignerPubKey: m.SignerPubKey,
		TxHash: m.TxHash, LogIndex: m.LogIndex,
		BlockNumber: m.BlockNumber, Nonce: m.Nonce,
	}
	return p.Marshal(), nil
}
func (m *ValidatorJoinMsg) AminoName() string { return "heimdallv2/stake/MsgValidatorJoin" }
func (m *ValidatorJoinMsg) AminoJSON() (any, error) {
	return map[string]any{
		"from":             m.From,
		"val_id":           fmt.Sprintf("%d", m.ValID),
		"activation_epoch": fmt.Sprintf("%d", m.ActivationEpoch),
		"amount":           m.Amount,
		"signer_pub_key":   m.SignerPubKey,
		"tx_hash":          m.TxHash,
		"log_index":        fmt.Sprintf("%d", m.LogIndex),
		"block_number":     fmt.Sprintf("%d", m.BlockNumber),
		"nonce":            fmt.Sprintf("%d", m.Nonce),
	}, nil
}

// StakeUpdateMsg wraps MsgStakeUpdate.
type StakeUpdateMsg struct {
	From        string
	ValID       uint64
	NewAmount   string
	TxHash      []byte
	LogIndex    uint64
	BlockNumber uint64
	Nonce       uint64
}

func (m *StakeUpdateMsg) TypeURL() string { return hproto.MsgStakeUpdateTypeURL }
func (m *StakeUpdateMsg) Marshal() ([]byte, error) {
	if m.From == "" {
		return nil, fmt.Errorf("StakeUpdateMsg: from is required")
	}
	p := &hproto.MsgStakeUpdate{
		From: m.From, ValID: m.ValID, NewAmount: m.NewAmount,
		TxHash: m.TxHash, LogIndex: m.LogIndex,
		BlockNumber: m.BlockNumber, Nonce: m.Nonce,
	}
	return p.Marshal(), nil
}
func (m *StakeUpdateMsg) AminoName() string { return "heimdallv2/stake/MsgStakeUpdate" }
func (m *StakeUpdateMsg) AminoJSON() (any, error) {
	return map[string]any{
		"from":         m.From,
		"val_id":       fmt.Sprintf("%d", m.ValID),
		"new_amount":   m.NewAmount,
		"tx_hash":      m.TxHash,
		"log_index":    fmt.Sprintf("%d", m.LogIndex),
		"block_number": fmt.Sprintf("%d", m.BlockNumber),
		"nonce":        fmt.Sprintf("%d", m.Nonce),
	}, nil
}

// SignerUpdateMsg wraps MsgSignerUpdate.
type SignerUpdateMsg struct {
	From            string
	ValID           uint64
	NewSignerPubKey []byte
	TxHash          []byte
	LogIndex        uint64
	BlockNumber     uint64
	Nonce           uint64
}

func (m *SignerUpdateMsg) TypeURL() string { return hproto.MsgSignerUpdateTypeURL }
func (m *SignerUpdateMsg) Marshal() ([]byte, error) {
	if m.From == "" {
		return nil, fmt.Errorf("SignerUpdateMsg: from is required")
	}
	p := &hproto.MsgSignerUpdate{
		From: m.From, ValID: m.ValID, NewSignerPubKey: m.NewSignerPubKey,
		TxHash: m.TxHash, LogIndex: m.LogIndex,
		BlockNumber: m.BlockNumber, Nonce: m.Nonce,
	}
	return p.Marshal(), nil
}
func (m *SignerUpdateMsg) AminoName() string { return "heimdallv2/stake/MsgSignerUpdate" }
func (m *SignerUpdateMsg) AminoJSON() (any, error) {
	return map[string]any{
		"from":               m.From,
		"val_id":             fmt.Sprintf("%d", m.ValID),
		"new_signer_pub_key": m.NewSignerPubKey,
		"tx_hash":            m.TxHash,
		"log_index":          fmt.Sprintf("%d", m.LogIndex),
		"block_number":       fmt.Sprintf("%d", m.BlockNumber),
		"nonce":              fmt.Sprintf("%d", m.Nonce),
	}, nil
}

// ValidatorExitMsg wraps MsgValidatorExit.
type ValidatorExitMsg struct {
	From              string
	ValID             uint64
	DeactivationEpoch uint64
	TxHash            []byte
	LogIndex          uint64
	BlockNumber       uint64
	Nonce             uint64
}

func (m *ValidatorExitMsg) TypeURL() string { return hproto.MsgValidatorExitTypeURL }
func (m *ValidatorExitMsg) Marshal() ([]byte, error) {
	if m.From == "" {
		return nil, fmt.Errorf("ValidatorExitMsg: from is required")
	}
	p := &hproto.MsgValidatorExit{
		From: m.From, ValID: m.ValID, DeactivationEpoch: m.DeactivationEpoch,
		TxHash: m.TxHash, LogIndex: m.LogIndex,
		BlockNumber: m.BlockNumber, Nonce: m.Nonce,
	}
	return p.Marshal(), nil
}
func (m *ValidatorExitMsg) AminoName() string { return "heimdallv2/stake/MsgValidatorExit" }
func (m *ValidatorExitMsg) AminoJSON() (any, error) {
	return map[string]any{
		"from":               m.From,
		"val_id":             fmt.Sprintf("%d", m.ValID),
		"deactivation_epoch": fmt.Sprintf("%d", m.DeactivationEpoch),
		"tx_hash":            m.TxHash,
		"log_index":          fmt.Sprintf("%d", m.LogIndex),
		"block_number":       fmt.Sprintf("%d", m.BlockNumber),
		"nonce":              fmt.Sprintf("%d", m.Nonce),
	}, nil
}

// ClerkEventRecordMsg wraps MsgEventRecord.
type ClerkEventRecordMsg struct {
	From            string
	TxHash          string
	LogIndex        uint64
	BlockNumber     uint64
	ContractAddress string
	Data            []byte
	ID              uint64
	ChainID         string
}

func (m *ClerkEventRecordMsg) TypeURL() string { return hproto.MsgEventRecordTypeURL }
func (m *ClerkEventRecordMsg) Marshal() ([]byte, error) {
	if m.From == "" {
		return nil, fmt.Errorf("ClerkEventRecordMsg: from is required")
	}
	p := &hproto.MsgEventRecord{
		From: m.From, TxHash: m.TxHash, LogIndex: m.LogIndex,
		BlockNumber: m.BlockNumber, ContractAddress: m.ContractAddress,
		Data: m.Data, ID: m.ID, ChainID: m.ChainID,
	}
	return p.Marshal(), nil
}
func (m *ClerkEventRecordMsg) AminoName() string { return "heimdallv2/clerk/MsgEventRecord" }
func (m *ClerkEventRecordMsg) AminoJSON() (any, error) {
	return map[string]any{
		"from":             m.From,
		"tx_hash":          m.TxHash,
		"log_index":        fmt.Sprintf("%d", m.LogIndex),
		"block_number":     fmt.Sprintf("%d", m.BlockNumber),
		"contract_address": m.ContractAddress,
		"data":             m.Data,
		"id":               fmt.Sprintf("%d", m.ID),
		"chain_id":         m.ChainID,
	}, nil
}
