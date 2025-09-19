package agglayer

import (
	"math/big"
	"strings"

	"github.com/0xPolygon/polygon-cli/cmd/ulxly/bridge_service"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type GetBridgeResponse struct {
	Bridges []BridgeResponse `json:"bridges"`
	Count   int              `json:"count"`
}

type BridgeResponse struct {
	LeafType uint8 `json:"leaf_type"`

	OrigNet    uint32 `json:"origin_network"`
	OrigAddr   string `json:"origin_address"`
	Amount     string `json:"amount"`
	DestNet    uint32 `json:"destination_network"`
	DestAddr   string `json:"destination_address"`
	BlockNum   uint64 `json:"block_num"`
	DepositCnt uint32 `json:"deposit_count"`
	NetworkID  uint32 `json:"network_id"`
	TxHash     string `json:"tx_hash"`
	Metadata   string `json:"metadata"`

	// Additional fields from the AggLayer API
	BridgeHash     string `json:"bridge_hash"`
	BlockPos       uint64 `json:"block_pos"`
	BlockTimestamp uint64 `json:"block_timestamp"`
	Calldata       string `json:"calldata"`
	FromAddress    string `json:"from_address"`

	// TODO: Thiago REVIEW THESE INFORMATION USED IN THE LEGACY BRIDGE SERVICE
	// ClaimTxHash   string `json:"claim_tx_hash"`
	GlobalIndex string `json:"global_index"`
}

func (r *BridgeResponse) ToDeposit(isReadyForClaim bool) (*bridge_service.Deposit, error) {
	d := &bridge_service.Deposit{}
	d.BlockNum = r.BlockNum

	var err error
	d.GlobalIndex, err = r.generateGlobalIndex(false, false)
	if err != nil {
		return nil, err
	}

	d.Amount = new(big.Int)
	d.Amount.SetString(r.Amount, 10)

	d.TxHash = common.HexToHash(r.TxHash)
	// TODO: Thiago review
	// if len(r.ClaimTxHash) > 0 {
	// 	claimTxHash := common.HexToHash(r.ClaimTxHash)
	// 	d.ClaimTxHash = &claimTxHash
	// }

	d.OrigAddr = common.HexToAddress(r.OrigAddr)
	d.DestAddr = common.HexToAddress(r.DestAddr)

	d.Metadata = common.Hex2Bytes(strings.TrimPrefix(r.Metadata, "0x"))

	d.LeafType = r.LeafType
	d.OrigNet = r.OrigNet
	d.DestNet = r.DestNet
	d.NetworkID = r.NetworkID
	d.DepositCnt = r.DepositCnt
	d.ReadyForClaim = isReadyForClaim

	return d, nil
}

// GenerateGlobalIndex converts the bash `generate_global_index` into Go.
// - bridgeInfoJSON: JSON string that contains a "deposit_count"
// - sourceNetworkID: 32-bit source network id
// - manipulatedUnusedBits / manipulatedRollupID: same flags as in the bash script
//
// Returns the computed 256-bit value as *big.Int.
func (r *BridgeResponse) generateGlobalIndex(manipulatedUnusedBits, manipulatedRollupID bool) (*big.Int, error) {
	// ---- Parse deposit_count from JSON (supports number or string) ----
	depositCount := r.DepositCnt

	// ---- Mask both to 32 bits (match bash behavior) ----
	srcMasked := uint64(r.OrigNet) & 0xFFFFFFFF
	depMasked := uint64(depositCount) & 0xFFFFFFFF

	final := new(big.Int) // starts at 0

	// precompute some powers of two we need
	twoPow32 := new(big.Int).Lsh(big.NewInt(1), 32)   // 2^32
	twoPow64 := new(big.Int).Lsh(big.NewInt(1), 64)   // 2^64
	twoPow128 := new(big.Int).Lsh(big.NewInt(1), 128) // 2^128

	// Offsets used when flags are set (match bash constants)
	unusedBitsOffset := new(big.Int).Mul(big.NewInt(10), twoPow128) // 10 * 2^128
	rollupIDOffset := new(big.Int).Mul(big.NewInt(10), twoPow32)    // 10 * 2^32

	// ---- 192nd bit logic (bash adds 2^64 when source_network_id == 0) ----
	if srcMasked == 0 {
		// final_value += 2^64
		final.Add(final, twoPow64)

		if manipulatedUnusedBits {
			log.Trace().Msg("-------------------------- Manipulated unused bits: true")
			final.Add(final, unusedBitsOffset)
		}
		if manipulatedRollupID {
			log.Trace().Msg("-------------------------- Manipulated rollup id: true")
			final.Add(final, rollupIDOffset)
		}
	}

	// ---- 193-224 bits: if mainnet is 0 -> 0; otherwise (source_network_id - 1) * 2^32 ----
	if srcMasked != 0 {
		// (source_network_id - 1) * 2^32
		mult := new(big.Int).SetUint64(srcMasked - 1)
		destShifted := new(big.Int).Mul(mult, twoPow32)
		final.Add(final, destShifted)

		if manipulatedUnusedBits {
			log.Trace().Msg("🔍 -------------------------- Manipulated unused bits: true")
			final.Add(final, unusedBitsOffset)
		}
	}

	// ---- 225-256 bits: add deposit_count (lower 32 bits) ----
	final.Add(final, new(big.Int).SetUint64(depMasked))

	return final, nil
}
