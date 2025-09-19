package legacy

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	"github.com/ethereum/go-ethereum/common"
)

type GetDepositResponse struct {
	Deposit DepositResponse `json:"deposit"`
	Code    *int            `json:"code"`
	Message *string         `json:"message"`
}

type DepositResponse struct {
	LeafType      uint8  `json:"leaf_type"`
	OrigNet       uint32 `json:"orig_net"`
	OrigAddr      string `json:"orig_addr"`
	Amount        string `json:"amount"`
	DestNet       uint32 `json:"dest_net"`
	DestAddr      string `json:"dest_addr"`
	BlockNum      string `json:"block_num"`
	DepositCnt    uint32 `json:"deposit_cnt"`
	NetworkID     uint32 `json:"network_id"`
	TxHash        string `json:"tx_hash"`
	ClaimTxHash   string `json:"claim_tx_hash"`
	Metadata      string `json:"metadata"`
	ReadyForClaim bool   `json:"ready_for_claim"`
	GlobalIndex   string `json:"global_index"`
}

func (r *DepositResponse) ToDeposit() (*bridge_service.Deposit, error) {
	d := &bridge_service.Deposit{}
	var err error
	d.BlockNum, err = strconv.ParseUint(r.BlockNum, 10, 64)
	if err != nil {
		return nil, err
	}

	d.GlobalIndex = new(big.Int)
	d.GlobalIndex.SetString(r.GlobalIndex, 10)

	d.Amount = new(big.Int)
	d.Amount.SetString(r.Amount, 10)

	d.TxHash = common.HexToHash(r.TxHash)
	if len(r.ClaimTxHash) > 0 {
		claimTxHash := common.HexToHash(r.ClaimTxHash)
		d.ClaimTxHash = &claimTxHash
	}

	d.OrigAddr = common.HexToAddress(r.OrigAddr)
	d.DestAddr = common.HexToAddress(r.DestAddr)

	d.Metadata = common.Hex2Bytes(strings.TrimPrefix(r.Metadata, "0x"))

	d.LeafType = r.LeafType
	d.OrigNet = r.OrigNet
	d.DestNet = r.DestNet
	d.NetworkID = r.NetworkID
	d.DepositCnt = r.DepositCnt
	d.ReadyForClaim = r.ReadyForClaim

	return d, nil
}

type GetDepositsResponse struct {
	Deposits []DepositResponse `json:"deposits"`
	Total    int               `json:"total_cnt,string"`
}
