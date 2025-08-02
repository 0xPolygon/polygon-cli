// Package rpctypes provides pretty JSON marshaling functionality for blockchain data types.
//
// The pretty marshaling functions convert hex-encoded values to their native types
// for better human readability. For example:
//   - Timestamps: "0x55ba467c" -> 1438271100 (uint64)
//   - Gas values: "0x5208" -> 21000 (uint64)
//   - Addresses: "0x1234..." -> common.Address type
//   - Hashes: "0xabcd..." -> common.Hash type
//
// Usage example:
//
//	block := NewPolyBlock(rawBlockResponse)
//	prettyJSON, err := PolyBlockToPrettyJSON(block)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(prettyJSON))
package rpctypes

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// HexBytes represents byte data that should be JSON-marshaled as hex strings
type HexBytes []byte

// MarshalJSON implements the json.Marshaler interface for HexBytes
func (h HexBytes) MarshalJSON() ([]byte, error) {
	if h == nil {
		return []byte("null"), nil
	}
	if len(h) == 0 {
		return []byte(`"0x"`), nil
	}
	hexStr := "0x" + hex.EncodeToString(h)
	return json.Marshal(hexStr)
}

// UnmarshalJSON implements the json.Unmarshaler interface for HexBytes
func (h *HexBytes) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s == "" || s == "0x" {
		*h = HexBytes{}
		return nil
	}

	if !strings.HasPrefix(s, "0x") {
		return fmt.Errorf("hex string must start with 0x")
	}

	decoded, err := hex.DecodeString(s[2:])
	if err != nil {
		return err
	}

	*h = HexBytes(decoded)
	return nil
}

// PrettyBlock represents a human-readable block structure
type PrettyBlock struct {
	Number           *big.Int            `json:"number"`
	Hash             ethcommon.Hash      `json:"hash"`
	ParentHash       ethcommon.Hash      `json:"parentHash"`
	Nonce            uint64              `json:"nonce"`
	SHA3Uncles       ethcommon.Hash      `json:"sha3Uncles"`
	LogsBloom        HexBytes            `json:"logsBloom"`
	TransactionsRoot ethcommon.Hash      `json:"transactionsRoot"`
	StateRoot        ethcommon.Hash      `json:"stateRoot"`
	ReceiptsRoot     ethcommon.Hash      `json:"receiptsRoot"`
	Miner            ethcommon.Address   `json:"miner"`
	Difficulty       *big.Int            `json:"difficulty"`
	TotalDifficulty  *big.Int            `json:"totalDifficulty,omitempty"`
	ExtraData        HexBytes            `json:"extraData"`
	Size             uint64              `json:"size"`
	GasLimit         uint64              `json:"gasLimit"`
	GasUsed          uint64              `json:"gasUsed"`
	Timestamp        uint64              `json:"timestamp"`
	Transactions     []PrettyTransaction `json:"transactions"`
	Uncles           []ethcommon.Hash    `json:"uncles"`
	BaseFeePerGas    *big.Int            `json:"baseFeePerGas"`
	MixHash          ethcommon.Hash      `json:"mixHash"`
}

// PrettyTransaction represents a human-readable transaction structure
type PrettyTransaction struct {
	BlockHash            ethcommon.Hash    `json:"blockHash"`
	BlockNumber          *big.Int          `json:"blockNumber"`
	From                 ethcommon.Address `json:"from"`
	Gas                  uint64            `json:"gas"`
	GasPrice             *big.Int          `json:"gasPrice"`
	MaxPriorityFeePerGas uint64            `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         uint64            `json:"maxFeePerGas"`
	Hash                 ethcommon.Hash    `json:"hash"`
	Input                HexBytes          `json:"input"`
	Nonce                uint64            `json:"nonce"`
	To                   ethcommon.Address `json:"to"`
	TransactionIndex     uint64            `json:"transactionIndex"`
	Value                *big.Int          `json:"value"`
	V                    *big.Int          `json:"v"`
	R                    *big.Int          `json:"r"`
	S                    *big.Int          `json:"s"`
	Type                 uint64            `json:"type"`
	ChainID              uint64            `json:"chainId"`
	AccessList           []any             `json:"accessList"`
}

// PrettyTxLogs represents a human-readable log structure
type PrettyTxLogs struct {
	BlockHash        ethcommon.Hash    `json:"blockHash"`
	BlockNumber      uint64            `json:"blockNumber"`
	TransactionIndex uint64            `json:"transactionIndex"`
	Address          ethcommon.Address `json:"address"`
	LogIndex         uint64            `json:"logIndex"`
	Data             HexBytes          `json:"data"`
	Removed          bool              `json:"removed"`
	Topics           []ethcommon.Hash  `json:"topics"`
	TransactionHash  ethcommon.Hash    `json:"transactionHash"`
}

// PrettyReceipt represents a human-readable receipt structure
type PrettyReceipt struct {
	TransactionHash   ethcommon.Hash    `json:"transactionHash"`
	TransactionIndex  uint64            `json:"transactionIndex"`
	BlockHash         ethcommon.Hash    `json:"blockHash"`
	BlockNumber       *big.Int          `json:"blockNumber"`
	From              ethcommon.Address `json:"from"`
	To                ethcommon.Address `json:"to"`
	CumulativeGasUsed *big.Int          `json:"cumulativeGasUsed"`
	EffectiveGasPrice *big.Int          `json:"effectiveGasPrice"`
	GasUsed           *big.Int          `json:"gasUsed"`
	ContractAddress   ethcommon.Address `json:"contractAddress"`
	Logs              []PrettyTxLogs    `json:"logs"`
	LogsBloom         HexBytes          `json:"logsBloom"`
	Root              ethcommon.Hash    `json:"root"`
	Status            uint64            `json:"status"`
	BlobGasPrice      *big.Int          `json:"blobGasPrice"`
	BlobGasUsed       *big.Int          `json:"blobGasUsed"`
}

// MarshalJSONPretty returns a pretty JSON representation of a PolyBlock
func (i *implPolyBlock) MarshalJSONPretty() ([]byte, error) {
	// Convert transactions to pretty format
	prettyTxs := make([]PrettyTransaction, len(i.inner.Transactions))
	for idx, tx := range i.inner.Transactions {
		prettyTxs[idx] = PrettyTransaction{
			BlockHash:            tx.BlockHash.ToHash(),
			BlockNumber:          tx.BlockNumber.ToBigInt(),
			From:                 tx.From.ToAddress(),
			Gas:                  tx.Gas.ToUint64(),
			GasPrice:             tx.GasPrice.ToBigInt(),
			MaxPriorityFeePerGas: tx.MaxPriorityFeePerGas.ToUint64(),
			MaxFeePerGas:         tx.MaxFeePerGas.ToUint64(),
			Hash:                 tx.Hash.ToHash(),
			Input:                HexBytes(tx.Input.ToBytes()),
			Nonce:                tx.Nonce.ToUint64(),
			To:                   tx.To.ToAddress(),
			TransactionIndex:     tx.TransactionIndex.ToUint64(),
			Value:                tx.Value.ToBigInt(),
			V:                    tx.V.ToBigInt(),
			R:                    tx.R.ToBigInt(),
			S:                    tx.S.ToBigInt(),
			Type:                 tx.Type.ToUint64(),
			ChainID:              tx.ChainID.ToUint64(),
			AccessList:           tx.AccessList,
		}
	}

	// Convert uncles to hash format
	prettyUncles := make([]ethcommon.Hash, len(i.inner.Uncles))
	for idx, uncle := range i.inner.Uncles {
		prettyUncles[idx] = uncle.ToHash()
	}

	pretty := PrettyBlock{
		Number:           i.inner.Number.ToBigInt(),
		Hash:             i.inner.Hash.ToHash(),
		ParentHash:       i.inner.ParentHash.ToHash(),
		Nonce:            i.inner.Nonce.ToUint64(),
		SHA3Uncles:       i.inner.SHA3Uncles.ToHash(),
		LogsBloom:        HexBytes(i.inner.LogsBloom.ToBytes()),
		TransactionsRoot: i.inner.TransactionsRoot.ToHash(),
		StateRoot:        i.inner.StateRoot.ToHash(),
		ReceiptsRoot:     i.inner.ReceiptsRoot.ToHash(),
		Miner:            i.inner.Miner.ToAddress(),
		Difficulty:       i.inner.Difficulty.ToBigInt(),
		ExtraData:        HexBytes(i.inner.ExtraData.ToBytes()),
		Size:             i.inner.Size.ToUint64(),
		GasLimit:         i.inner.GasLimit.ToUint64(),
		GasUsed:          i.inner.GasUsed.ToUint64(),
		Timestamp:        i.inner.Timestamp.ToUint64(),
		Transactions:     prettyTxs,
		Uncles:           prettyUncles,
		BaseFeePerGas:    i.inner.BaseFeePerGas.ToBigInt(),
		MixHash:          i.inner.MixHash.ToHash(),
	}

	return json.Marshal(pretty)
}

// MarshalJSONPretty returns a pretty JSON representation of a PolyTransaction
func (i *implPolyTransaction) MarshalJSONPretty() ([]byte, error) {
	pretty := PrettyTransaction{
		BlockHash:            i.inner.BlockHash.ToHash(),
		BlockNumber:          i.inner.BlockNumber.ToBigInt(),
		From:                 i.inner.From.ToAddress(),
		Gas:                  i.inner.Gas.ToUint64(),
		GasPrice:             i.inner.GasPrice.ToBigInt(),
		MaxPriorityFeePerGas: i.inner.MaxPriorityFeePerGas.ToUint64(),
		MaxFeePerGas:         i.inner.MaxFeePerGas.ToUint64(),
		Hash:                 i.inner.Hash.ToHash(),
		Input:                HexBytes(i.inner.Input.ToBytes()),
		Nonce:                i.inner.Nonce.ToUint64(),
		To:                   i.inner.To.ToAddress(),
		TransactionIndex:     i.inner.TransactionIndex.ToUint64(),
		Value:                i.inner.Value.ToBigInt(),
		V:                    i.inner.V.ToBigInt(),
		R:                    i.inner.R.ToBigInt(),
		S:                    i.inner.S.ToBigInt(),
		Type:                 i.inner.Type.ToUint64(),
		ChainID:              i.inner.ChainID.ToUint64(),
		AccessList:           i.inner.AccessList,
	}

	return json.Marshal(pretty)
}

// MarshalJSONPretty returns a pretty JSON representation of a PolyReceipt
func (i *implPolyReceipt) MarshalJSONPretty() ([]byte, error) {
	// Convert logs to pretty format
	prettyLogs := make([]PrettyTxLogs, len(i.inner.Logs))
	for idx, log := range i.inner.Logs {
		// Convert topics to hash format
		prettyTopics := make([]ethcommon.Hash, len(log.Topics))
		for topicIdx, topic := range log.Topics {
			prettyTopics[topicIdx] = topic.ToHash()
		}

		prettyLogs[idx] = PrettyTxLogs{
			BlockHash:        log.BlockHash.ToHash(),
			BlockNumber:      log.BlockNumber.ToUint64(),
			TransactionIndex: log.TransactionIndex.ToUint64(),
			Address:          log.Address.ToAddress(),
			LogIndex:         log.LogIndex.ToUint64(),
			Data:             HexBytes(log.Data.ToBytes()),
			Removed:          log.Removed,
			Topics:           prettyTopics,
			TransactionHash:  log.TransactionHash.ToHash(),
		}
	}

	pretty := PrettyReceipt{
		TransactionHash:   i.inner.TransactionHash.ToHash(),
		TransactionIndex:  i.inner.TransactionIndex.ToUint64(),
		BlockHash:         i.inner.BlockHash.ToHash(),
		BlockNumber:       i.inner.BlockNumber.ToBigInt(),
		From:              i.inner.From.ToAddress(),
		To:                i.inner.To.ToAddress(),
		CumulativeGasUsed: i.inner.CumulativeGasUsed.ToBigInt(),
		EffectiveGasPrice: i.inner.EffectiveGasPrice.ToBigInt(),
		GasUsed:           i.inner.GasUsed.ToBigInt(),
		ContractAddress:   i.inner.ContractAddress.ToAddress(),
		Logs:              prettyLogs,
		LogsBloom:         HexBytes(i.inner.LogsBloom.ToBytes()),
		Root:              i.inner.Root.ToHash(),
		Status:            i.inner.Status.ToUint64(),
		BlobGasPrice:      i.inner.BlobGasPrice.ToBigInt(),
		BlobGasUsed:       i.inner.BlobGasUsed.ToBigInt(),
	}

	return json.Marshal(pretty)
}

// Interface extensions for pretty marshaling
type PrettyMarshaler interface {
	MarshalJSONPretty() ([]byte, error)
}

// Helper functions to convert poly types to their pretty JSON representations
func PolyBlockToPrettyJSON(block PolyBlock) ([]byte, error) {
	if impl, ok := block.(*implPolyBlock); ok {
		return impl.MarshalJSONPretty()
	}
	return json.Marshal(block)
}

func PolyTransactionToPrettyJSON(tx PolyTransaction) ([]byte, error) {
	if impl, ok := tx.(*implPolyTransaction); ok {
		return impl.MarshalJSONPretty()
	}
	return json.Marshal(tx)
}

func PolyReceiptToPrettyJSON(receipt PolyReceipt) ([]byte, error) {
	if impl, ok := receipt.(*implPolyReceipt); ok {
		return impl.MarshalJSONPretty()
	}
	return json.Marshal(receipt)
}
