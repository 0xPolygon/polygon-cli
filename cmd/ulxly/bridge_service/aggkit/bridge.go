package aggkit

import (
	"fmt"
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

	OrigNet    uint32 `json:"origin_network"`
	OrigAddr   string `json:"origin_address"`
	Amount     string `json:"amount"`
	DestNet    uint32 `json:"destination_network"`
	DestAddr   string `json:"destination_address"`
	BlockNum   uint64 `json:"block_num"`
	DepositCnt uint32 `json:"deposit_count"`
	TxHash     string `json:"tx_hash"`
	Metadata   string `json:"metadata"`

	// Additional fields from the Aggkit API
	BridgeHash     string `json:"bridge_hash"`
	BlockPos       uint64 `json:"block_pos"`
	BlockTimestamp uint64 `json:"block_timestamp"`
	Calldata       string `json:"calldata"`
	FromAddress    string `json:"from_address"`
}

func (r *bridgeResponse) ToDeposit(networkID uint32) (*bridge_service.Deposit, error) {
	d := &bridge_service.Deposit{}
	d.BlockNum = r.BlockNum

	d.GlobalIndex = r.generateGlobalIndex(networkID)

	d.Amount = new(big.Int)
	_, ok := d.Amount.SetString(r.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount: %s", r.Amount)
	}

	d.TxHash = common.HexToHash(r.TxHash)

	d.OrigAddr = common.HexToAddress(r.OrigAddr)
	d.DestAddr = common.HexToAddress(r.DestAddr)

	d.Metadata = common.Hex2Bytes(strings.TrimPrefix(r.Metadata, "0x"))

	d.LeafType = r.LeafType
	d.OrigNet = r.OrigNet
	d.DestNet = r.DestNet
	d.NetworkID = networkID
	d.DepositCnt = r.DepositCnt

	return d, nil
}

func (r *bridgeResponse) generateGlobalIndex(networkID uint32) *big.Int {
	result := new(big.Int)
	if networkID == 0 {
		// Mainnet: set bit 64 (add 2^64)
		result.SetBit(result, 64, 1)
	} else {
		// Non-mainnet: (networkID - 1) shifted left by 32 bits
		result.SetUint64(uint64(networkID-1) << 32)
	}
	// Add deposit count to the lower 32 bits
	result.Add(result, big.NewInt(int64(r.DepositCnt)))
	return result
}
