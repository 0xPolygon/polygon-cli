package database

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"
)

// JSONDatabase outputs data as JSON to stdout.
// Each record is output as a single line of JSON (newline-delimited JSON).
type JSONDatabase struct {
	sensorID                     string
	chainID                      *big.Int
	shouldWriteBlocks            bool
	shouldWriteBlockEvents       bool
	shouldWriteTransactions      bool
	shouldWriteTransactionEvents bool
	shouldWritePeers             bool
}

// JSONDatabaseOptions is used when creating a NewJSONDatabase.
type JSONDatabaseOptions struct {
	SensorID                     string
	ChainID                      uint64
	MaxConcurrency               int
	ShouldWriteBlocks            bool
	ShouldWriteBlockEvents       bool
	ShouldWriteTransactions      bool
	ShouldWriteTransactionEvents bool
	ShouldWritePeers             bool
}

// NewJSONDatabase creates a new JSONDatabase instance.
func NewJSONDatabase(opts JSONDatabaseOptions) Database {
	return &JSONDatabase{
		sensorID:                     opts.SensorID,
		chainID:                      new(big.Int).SetUint64(opts.ChainID),
		shouldWriteBlocks:            opts.ShouldWriteBlocks,
		shouldWriteBlockEvents:       opts.ShouldWriteBlockEvents,
		shouldWriteTransactions:      opts.ShouldWriteTransactions,
		shouldWriteTransactionEvents: opts.ShouldWriteTransactionEvents,
		shouldWritePeers:             opts.ShouldWritePeers,
	}
}

// JSONBlock represents a block in JSON format.
type JSONBlock struct {
	Type                string    `json:"type"`
	SensorID            string    `json:"sensor_id"`
	Hash                string    `json:"hash"`
	ParentHash          string    `json:"parent_hash"`
	Number              uint64    `json:"number"`
	Timestamp           uint64    `json:"timestamp"`
	GasLimit            uint64    `json:"gas_limit"`
	GasUsed             uint64    `json:"gas_used"`
	Difficulty          string    `json:"difficulty,omitempty"`
	TotalDifficulty     string    `json:"total_difficulty,omitempty"`
	BaseFee             string    `json:"base_fee,omitempty"`
	TxCount             int       `json:"tx_count"`
	UncleCount          int       `json:"uncle_count"`
	TimeFirstSeen       time.Time `json:"time_first_seen"`
	TimeFirstSeenHash   time.Time `json:"time_first_seen_hash"`
	IsParent            bool      `json:"is_parent"`
	SensorFirstSeen     string    `json:"sensor_first_seen"`
	SensorFirstSeenHash string    `json:"sensor_first_seen_hash,omitempty"`
}

// JSONBlockEvent represents a block event in JSON format.
type JSONBlockEvent struct {
	Type      string    `json:"type"`
	SensorID  string    `json:"sensor_id"`
	PeerID    string    `json:"peer_id"`
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"`
}

// JSONTransaction represents a transaction in JSON format.
type JSONTransaction struct {
	Type          string    `json:"type"`
	SensorID      string    `json:"sensor_id"`
	Hash          string    `json:"hash"`
	From          string    `json:"from,omitempty"`
	To            string    `json:"to,omitempty"`
	Value         string    `json:"value"`
	Gas           uint64    `json:"gas"`
	GasPrice      string    `json:"gas_price"`
	GasFeeCap     string    `json:"gas_fee_cap,omitempty"`
	GasTipCap     string    `json:"gas_tip_cap,omitempty"`
	Nonce         uint64    `json:"nonce"`
	TxType        uint8     `json:"tx_type"`
	TimeFirstSeen time.Time `json:"time_first_seen"`
}

// JSONTransactionEvent represents a transaction event in JSON format.
type JSONTransactionEvent struct {
	Type      string    `json:"type"`
	SensorID  string    `json:"sensor_id"`
	PeerID    string    `json:"peer_id"`
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"`
}

// JSONPeer represents a peer in JSON format.
type JSONPeer struct {
	Type         string    `json:"type"`
	SensorID     string    `json:"sensor_id"`
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	Caps         []string  `json:"caps"`
	TimeLastSeen time.Time `json:"time_last_seen"`
}

// Write safely outputs JSON to stdout.
func (j *JSONDatabase) Write(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal JSON")
		return
	}

	fmt.Fprintln(os.Stdout, string(data))
}

// WriteBlock writes the block and the block event as JSON.
func (j *JSONDatabase) WriteBlock(_ context.Context, peer *enode.Node, block *types.Block, td *big.Int, tfs time.Time) {
	j.writeBlockEvent(peer, block, tfs)
	j.writeBlock(block, td, tfs)
}

func (j *JSONDatabase) writeBlockEvent(peer *enode.Node, block *types.Block, tfs time.Time) {
	if !j.ShouldWriteBlockEvents() || peer == nil {
		return
	}

	event := JSONBlockEvent{
		Type:      "block_event",
		SensorID:  j.sensorID,
		PeerID:    peer.URLv4(),
		Hash:      block.Hash().Hex(),
		Timestamp: tfs,
	}

	j.Write(event)
}

func (j *JSONDatabase) writeBlock(block *types.Block, td *big.Int, tfs time.Time) {
	if !j.ShouldWriteBlocks() {
		return
	}

	jsonBlock := JSONBlock{
		Type:            "block",
		SensorID:        j.sensorID,
		Hash:            block.Hash().Hex(),
		ParentHash:      block.ParentHash().Hex(),
		Number:          block.NumberU64(),
		Timestamp:       block.Time(),
		GasLimit:        block.GasLimit(),
		GasUsed:         block.GasUsed(),
		Difficulty:      block.Difficulty().String(),
		TotalDifficulty: td.String(),
		TxCount:         len(block.Transactions()),
		UncleCount:      len(block.Uncles()),
		TimeFirstSeen:   tfs,
		SensorFirstSeen: j.sensorID,
	}

	if block.BaseFee() != nil {
		jsonBlock.BaseFee = block.BaseFee().String()
	}

	j.Write(jsonBlock)
}

// WriteBlockHeaders writes the block headers as JSON.
// The isParent parameter indicates if these headers were fetched as parent blocks.
func (j *JSONDatabase) WriteBlockHeaders(ctx context.Context, headers []*types.Header, tfs time.Time, isParent bool) {
	if !j.ShouldWriteBlocks() {
		return
	}

	for _, header := range headers {
		block := JSONBlock{
			Type:            "block_header",
			SensorID:        j.sensorID,
			Hash:            header.Hash().Hex(),
			ParentHash:      header.ParentHash.Hex(),
			Number:          header.Number.Uint64(),
			Timestamp:       header.Time,
			GasLimit:        header.GasLimit,
			GasUsed:         header.GasUsed,
			Difficulty:      header.Difficulty.String(),
			TimeFirstSeen:   tfs,
			IsParent:        isParent,
			SensorFirstSeen: j.sensorID,
		}

		if header.BaseFee != nil {
			block.BaseFee = header.BaseFee.String()
		}

		j.Write(block)
	}
}

// WriteBlockHashes writes the block events as JSON.
func (j *JSONDatabase) WriteBlockHashes(ctx context.Context, peer *enode.Node, hashes []common.Hash, tfs time.Time) {
	if !j.ShouldWriteBlockEvents() || len(hashes) == 0 || peer == nil {
		return
	}

	for _, hash := range hashes {
		event := JSONBlockEvent{
			Type:      "block_hash",
			SensorID:  j.sensorID,
			PeerID:    peer.URLv4(),
			Hash:      hash.Hex(),
			Timestamp: tfs,
		}

		j.Write(event)
	}
}

// WriteBlockHashFirstSeen writes a partial block entry with just the hash
// first seen time. For JSON output, this writes a separate record type.
func (j *JSONDatabase) WriteBlockHashFirstSeen(ctx context.Context, hash common.Hash, tfsh time.Time) {
	if !j.ShouldWriteBlocks() {
		return
	}

	partial := map[string]any{
		"type":                   "block_hash_first_seen",
		"sensor_id":              j.sensorID,
		"hash":                   hash.Hex(),
		"time_first_seen_hash":   tfsh,
		"sensor_first_seen_hash": j.sensorID,
	}

	j.Write(partial)
}

// WriteBlockBody writes the block body as JSON.
func (j *JSONDatabase) WriteBlockBody(ctx context.Context, body *eth.BlockBody, hash common.Hash, tfs time.Time) {
	if !j.ShouldWriteBlocks() {
		return
	}

	jsonBody := map[string]any{
		"type":            "block_body",
		"sensor_id":       j.sensorID,
		"hash":            hash.Hex(),
		"tx_count":        len(body.Transactions),
		"uncle_count":     len(body.Uncles),
		"time_first_seen": tfs,
	}

	j.Write(jsonBody)
}

// WriteTransactions writes the transactions and transaction events as JSON.
func (j *JSONDatabase) WriteTransactions(_ context.Context, peer *enode.Node, txs []*types.Transaction, tfs time.Time) {
	j.writeTxs(txs, tfs)
	j.writeTxEvents(peer, txs, tfs)
}

func (j *JSONDatabase) writeTxs(txs []*types.Transaction, tfs time.Time) {
	if !j.ShouldWriteTransactions() {
		return
	}

	for _, tx := range txs {
		chainID := tx.ChainId()
		if tx.ChainId() == nil || tx.ChainId().Sign() <= 0 {
			chainID = j.chainID
		}

		var from string
		addr, err := types.Sender(types.LatestSignerForChainID(chainID), tx)
		if err == nil {
			from = addr.Hex()
		}

		jsonTx := JSONTransaction{
			Type:          "transaction",
			SensorID:      j.sensorID,
			Hash:          tx.Hash().Hex(),
			From:          from,
			Value:         tx.Value().String(),
			Gas:           tx.Gas(),
			GasPrice:      tx.GasPrice().String(),
			GasFeeCap:     tx.GasFeeCap().String(),
			GasTipCap:     tx.GasTipCap().String(),
			Nonce:         tx.Nonce(),
			TxType:        tx.Type(),
			TimeFirstSeen: tfs,
		}

		if tx.To() != nil {
			jsonTx.To = tx.To().Hex()
		}

		j.Write(jsonTx)
	}
}

func (j *JSONDatabase) writeTxEvents(peer *enode.Node, txs []*types.Transaction, tfs time.Time) {
	if !j.ShouldWriteTransactionEvents() || peer == nil {
		return
	}

	for _, tx := range txs {
		event := JSONTransactionEvent{
			Type:      "transaction_event",
			SensorID:  j.sensorID,
			PeerID:    peer.URLv4(),
			Hash:      tx.Hash().Hex(),
			Timestamp: tfs,
		}

		j.Write(event)
	}
}

// WritePeers writes the connected peers as JSON.
func (j *JSONDatabase) WritePeers(ctx context.Context, peers []*p2p.Peer, tls time.Time) {
	if !j.ShouldWritePeers() {
		return
	}

	for _, peer := range peers {
		jsonPeer := JSONPeer{
			Type:         "peer",
			SensorID:     j.sensorID,
			ID:           peer.ID().String(),
			Name:         peer.Fullname(),
			URL:          peer.Node().URLv4(),
			Caps:         peer.Info().Caps,
			TimeLastSeen: tls,
		}

		j.Write(jsonPeer)
	}
}

// HasBlock always returns true to avoid unnecessary parent block fetching for JSON output.
func (j *JSONDatabase) HasBlock(ctx context.Context, hash common.Hash) bool {
	return false
}

// MaxConcurrentWrites returns the max concurrency.
func (j *JSONDatabase) MaxConcurrentWrites() int {
	return 1
}

// ShouldWriteBlocks returns the configured value.
func (j *JSONDatabase) ShouldWriteBlocks() bool {
	return j.shouldWriteBlocks
}

// ShouldWriteBlockEvents returns the configured value.
func (j *JSONDatabase) ShouldWriteBlockEvents() bool {
	return j.shouldWriteBlockEvents
}

// ShouldWriteTransactions returns the configured value.
func (j *JSONDatabase) ShouldWriteTransactions() bool {
	return j.shouldWriteTransactions
}

// ShouldWriteTransactionEvents returns the configured value.
func (j *JSONDatabase) ShouldWriteTransactionEvents() bool {
	return j.shouldWriteTransactionEvents
}

// ShouldWritePeers returns the configured value.
func (j *JSONDatabase) ShouldWritePeers() bool {
	return j.shouldWritePeers
}

// NodeList returns an empty list as JSON database doesn't store nodes.
func (j *JSONDatabase) NodeList(ctx context.Context, limit int) ([]string, error) {
	return []string{}, nil
}
