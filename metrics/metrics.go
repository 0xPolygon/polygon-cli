package metrics

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/ethereum/go-ethereum/consensus/clique"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

var (
	UnitWei       = big.NewInt(1)                                    // | 1 | 1 | wei | Wei
	UnitBabbage   = new(big.Int).Mul(UnitWei, big.NewInt(1000))      // | 1,000 | 10^3^ | Babbage | Kilowei or femtoether
	UnitLovelace  = new(big.Int).Mul(UnitBabbage, big.NewInt(1000))  // | 1,000,000 | 10^6^ | Lovelace | Megawei or picoether
	UnitShannon   = new(big.Int).Mul(UnitLovelace, big.NewInt(1000)) // | 1,000,000,000 | 10^9^ | Shannon | Gigawei or nanoether
	UnitSzabo     = new(big.Int).Mul(UnitShannon, big.NewInt(1000))  // | 1,000,000,000,000 | 10^12^ | Szabo | Microether or micro
	UnitFinney    = new(big.Int).Mul(UnitSzabo, big.NewInt(1000))    // | 1,000,000,000,000,000 | 10^15^ | Finney | Milliether or milli
	UnitEther     = new(big.Int).Mul(UnitFinney, big.NewInt(1000))   // | 1,000,000,000,000,000,000 | 10^18^ | Ether | Ether
	UnitGrand     = new(big.Int).Mul(UnitEther, big.NewInt(1000))    // | 1,000,000,000,000,000,000,000 | 10^21^ | Grand | Kiloether
	UnitMegaether = new(big.Int).Mul(UnitGrand, big.NewInt(1000))    // | 1,000,000,000,000,000,000,000,000 | 10^24^ | | Megaether
)

func GetMeanBlockTime(blocks []rpctypes.PolyBlock) float64 {
	if len(blocks) < 2 {
		return 0
	}

	t0 := blocks[0].Time()
	minTs, maxTs := t0, t0

	for i := 1; i < len(blocks); i++ {
		ts := blocks[i].Time()

		if ts < minTs {
			minTs = ts
		}

		if ts > maxTs {
			maxTs = ts
		}
	}

	return float64(maxTs-minTs) / float64(len(blocks)-1)
}

func GetTxsPerBlock(blocks []rpctypes.PolyBlock) []float64 {
	bs := rpctypes.SortableBlocks(blocks)
	sort.Sort(bs)

	txns := make([]float64, 0)
	for _, b := range bs {
		txns = append(txns, float64(len(b.Transactions())))
	}
	return txns
}

func GetUnclesPerBlock(blocks []rpctypes.PolyBlock) []float64 {
	bs := rpctypes.SortableBlocks(blocks)
	sort.Sort(bs)

	uncles := make([]float64, 0)
	for _, b := range bs {
		uncles = append(uncles, float64(len(b.Uncles())))
	}
	return uncles
}

func GetSizePerBlock(blocks []rpctypes.PolyBlock) []float64 {
	bs := rpctypes.SortableBlocks(blocks)
	sort.Sort(bs)

	bSize := make([]float64, 0)
	for _, b := range bs {
		bSize = append(bSize, float64(b.Size()))
	}
	return bSize
}

func GetGasPerBlock(blocks []rpctypes.PolyBlock) []float64 {
	bs := rpctypes.SortableBlocks(blocks)
	sort.Sort(bs)

	gasUsed := make([]float64, 0)
	for _, b := range bs {
		gasUsed = append(gasUsed, float64(b.GasUsed()))
	}
	return gasUsed
}

func GetMeanGasPricePerBlock(blocks []rpctypes.PolyBlock) []float64 {
	bs := rpctypes.SortableBlocks(blocks)
	sort.Sort(bs)

	gasPrices := make([]float64, 0)
	for _, b := range bs {
		totGas := big.NewInt(0)
		txs := b.Transactions()
		if len(txs) < 1 {
			gasPrices = append(gasPrices, 0.0)
			continue
		}
		for _, tx := range txs {
			totGas.Add(totGas, tx.GasPrice())
		}
		meanGas := totGas.Div(totGas, big.NewInt(int64(len(txs))))
		gasPrices = append(gasPrices, float64(meanGas.Int64()))
	}
	return gasPrices
}

// GetMeanBaseFeePerBlock calculates the mean base fee for each block in the provided slice of PolyBlock.
// It returns a slice of float64 values representing the base fee for each block.
func GetMeanBaseFeePerBlock(blocks []rpctypes.PolyBlock) []float64 {
	bs := rpctypes.SortableBlocks(blocks)
	sort.Sort(bs)

	fee := make([]float64, 0, len(bs))
	for _, b := range bs {
		if b.BaseFee() != nil {
			fee = append(fee, float64(b.BaseFee().Uint64()))
		} else {
			fee = append(fee, 0.0)
		}
	}
	return fee
}

func TruncateHexString(hexStr string, totalLength int) string {
	hexStr = strings.TrimPrefix(hexStr, "0x")

	visibleLength := max(totalLength-5, 0)

	if len(hexStr) <= visibleLength {
		return "0x" + hexStr
	}

	beginning := hexStr[:visibleLength/2]
	end := hexStr[len(hexStr)-visibleLength/2:]

	return "0x" + beginning + "..." + end
}

func Ecrecover(block *rpctypes.PolyBlock) ([]byte, error) {
	input, err := json.Marshal(*block)
	if err != nil {
		return nil, err
	}
	header := new(ethtypes.Header)
	err = header.UnmarshalJSON(input)
	if err != nil {
		return nil, err
	}
	sigStart := len(header.Extra) - ethcrypto.SignatureLength
	if sigStart < 0 || sigStart > len(header.Extra) {
		return nil, fmt.Errorf("unable to recover signature")
	}
	signature := header.Extra[sigStart:]
	pubKey, err := ethcrypto.Ecrecover(clique.SealHash(header).Bytes(), signature)
	if err != nil {
		return nil, err
	}
	signer := ethcrypto.Keccak256(pubKey[1:])[12:]

	return signer, nil
}

func RawDataToASCII(data []byte) string {
	retString := ""
	for _, b := range data {
		if b >= 32 && b < 127 {
			retString = retString + string(b)
		} else {
			retString = retString + fmt.Sprintf("\\x%X", b)
		}
	}
	return retString
}
