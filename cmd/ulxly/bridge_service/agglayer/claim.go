package agglayer

type getClaimsResponse struct {
	Claims []claimResponse `json:"claims"`
	Count  int             `json:"count"`
}

type claimResponse struct {
	Amount              string   `json:"amount"`
	BlockNum            uint64   `json:"block_num"`
	BlockTimestamp      uint64   `json:"block_timestamp"`
	DestinationAddress  string   `json:"destination_address"`
	DestinationNetwork  uint32   `json:"destination_network"`
	FromAddress         string   `json:"from_address"`
	GlobalExitRoot      string   `json:"global_exit_root"`
	GlobalIndex         string   `json:"global_index"`
	MainnetExitRoot     string   `json:"mainnet_exit_root"`
	Metadata            string   `json:"metadata"`
	OriginAddress       string   `json:"origin_address"`
	OriginNetwork       uint32   `json:"origin_network"`
	ProofLocalExitRoot  []string `json:"proof_local_exit_root"`
	ProofRollupExitRoot []string `json:"proof_rollup_exit_root"`
	RollupExitRoot      string   `json:"rollup_exit_root"`
	TxHash              string   `json:"tx_hash"`
}
