package agglayer

import (
	"math/big"
	"strings"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	"github.com/ethereum/go-ethereum/common"
)

type getBridgesResponse struct {
	Bridges []bridgeResponse `json:"bridges"`
	Count   int              `json:"count"`
}

type bridgeResponse struct {
	LeafType uint8 `json:"leaf_type"`

	OrigNet     uint32 `json:"origin_network"`
	OrigAddr    string `json:"origin_address"`
	Amount      string `json:"amount"`
	DestNet     uint32 `json:"destination_network"`
	DestAddr    string `json:"destination_address"`
	BlockNum    uint64 `json:"block_num"`
	DepositCnt  uint32 `json:"deposit_count"`
	TxHash      string `json:"tx_hash"`
	Metadata    string `json:"metadata"`
	GlobalIndex string `json:"global_index"`

	// Additional fields from the AggLayer API
	BridgeHash     string `json:"bridge_hash"`
	BlockPos       uint64 `json:"block_pos"`
	BlockTimestamp uint64 `json:"block_timestamp"`
	Calldata       string `json:"calldata"`
	FromAddress    string `json:"from_address"`
}

func (r *bridgeResponse) ToDeposit(networkID uint32, isReadyForClaim bool, claimTx string) *bridge_service.Deposit {
	d := &bridge_service.Deposit{}
	d.BlockNum = r.BlockNum

	d.GlobalIndex = new(big.Int)
	d.GlobalIndex.SetString(r.GlobalIndex, 10)

	d.Amount = new(big.Int)
	d.Amount.SetString(r.Amount, 10)

	d.TxHash = common.HexToHash(r.TxHash)
	if len(claimTx) > 0 {
		claimTxHash := common.HexToHash(claimTx)
		d.ClaimTxHash = &claimTxHash
	}

	d.OrigAddr = common.HexToAddress(r.OrigAddr)
	d.DestAddr = common.HexToAddress(r.DestAddr)

	d.Metadata = common.Hex2Bytes(strings.TrimPrefix(r.Metadata, "0x"))

	d.LeafType = r.LeafType
	d.OrigNet = r.OrigNet
	d.DestNet = r.DestNet
	d.NetworkID = networkID
	d.DepositCnt = r.DepositCnt
	d.ReadyForClaim = isReadyForClaim

	return d
}
