package report

import (
	"math/big"
	"time"
)

// BlockReport represents the complete report for a range of blocks
type BlockReport struct {
	ChainID     uint64       `json:"chain_id"`
	RpcUrl      string       `json:"rpc_url"`
	StartBlock  uint64       `json:"start_block"`
	EndBlock    uint64       `json:"end_block"`
	GeneratedAt time.Time    `json:"generated_at"`
	Summary     SummaryStats `json:"summary"`
	Top10       Top10Stats   `json:"top_10"`
	Blocks      []BlockInfo  `json:"-"` // Internal use only, not exported
}

// SummaryStats contains aggregate statistics for the block range
type SummaryStats struct {
	TotalBlocks       uint64  `json:"total_blocks"`
	TotalTransactions uint64  `json:"total_transactions"`
	TotalGasUsed      uint64  `json:"total_gas_used"`
	AvgTxPerBlock     float64 `json:"avg_tx_per_block"`
	AvgGasPerBlock    float64 `json:"avg_gas_per_block"`
	AvgBaseFeePerGas  uint64  `json:"avg_base_fee_per_gas,omitempty"`
	UniqueSenders     uint64  `json:"unique_senders"`
	UniqueRecipients  uint64  `json:"unique_recipients"`
}

// BlockInfo contains information about a single block
type BlockInfo struct {
	Number        uint64            `json:"number"`
	Timestamp     uint64            `json:"timestamp"`
	TxCount       uint64            `json:"tx_count"`
	GasUsed       uint64            `json:"gas_used"`
	GasLimit      uint64            `json:"gas_limit"`
	BaseFeePerGas *big.Int          `json:"base_fee_per_gas,omitempty"`
	Transactions  []TransactionInfo `json:"-"` // Internal use only
}

// TransactionInfo contains information about a single transaction
type TransactionInfo struct {
	Hash           string  `json:"hash"`
	From           string  `json:"from"`
	To             string  `json:"to"`
	BlockNumber    uint64  `json:"block_number"`
	GasUsed        uint64  `json:"gas_used"`
	GasLimit       uint64  `json:"gas_limit"`
	GasPrice       uint64  `json:"gas_price"`
	BlockGasLimit  uint64  `json:"block_gas_limit"`
	GasUsedPercent float64 `json:"gas_used_percent"`
}

// Top10Stats contains top 10 lists for various metrics
type Top10Stats struct {
	BlocksByTxCount        []TopBlock       `json:"blocks_by_tx_count"`
	BlocksByGasUsed        []TopBlock       `json:"blocks_by_gas_used"`
	TransactionsByGas      []TopTransaction `json:"transactions_by_gas"`
	TransactionsByGasLimit []TopTransaction `json:"transactions_by_gas_limit"`
	MostUsedGasPrices      []GasPriceFreq   `json:"most_used_gas_prices"`
	MostUsedGasLimits      []GasLimitFreq   `json:"most_used_gas_limits"`
}

// TopBlock represents a block in a top 10 list
type TopBlock struct {
	Number         uint64  `json:"number"`
	TxCount        uint64  `json:"tx_count,omitempty"`
	GasUsed        uint64  `json:"gas_used,omitempty"`
	GasLimit       uint64  `json:"gas_limit,omitempty"`
	GasUsedPercent float64 `json:"gas_used_percent,omitempty"`
}

// TopTransaction represents a transaction in a top 10 list
type TopTransaction struct {
	Hash           string  `json:"hash"`
	BlockNumber    uint64  `json:"block_number"`
	GasUsed        uint64  `json:"gas_used,omitempty"`
	GasLimit       uint64  `json:"gas_limit,omitempty"`
	BlockGasLimit  uint64  `json:"block_gas_limit,omitempty"`
	GasUsedPercent float64 `json:"gas_used_percent,omitempty"`
}

// GasPriceFreq represents the frequency of a specific gas price
type GasPriceFreq struct {
	GasPrice uint64 `json:"gas_price"`
	Count    uint64 `json:"count"`
}

// GasLimitFreq represents the frequency of a specific gas limit
type GasLimitFreq struct {
	GasLimit uint64 `json:"gas_limit"`
	Count    uint64 `json:"count"`
}
