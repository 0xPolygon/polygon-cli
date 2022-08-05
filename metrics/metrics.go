package metrics

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"time"

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

type (
	SortableBlocks []*ethtypes.Block
)

func (a SortableBlocks) Len() int {
	return len(a)
}
func (a SortableBlocks) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SortableBlocks) Less(i, j int) bool {
	return a[i].Time() < a[j].Time()
}

func GetMeanBlockTime(blocks []*ethtypes.Block) float64 {
	if len(blocks) < 2 {
		return 0
	}
	times := make([]int, 0)
	for _, block := range blocks {
		blockTime := block.Time()
		times = append(times, int(blockTime))
	}

	sortTimes := sort.IntSlice(times)
	sortTimes.Sort()

	minTime := sortTimes[0]
	maxTime := sortTimes[len(sortTimes)-1]

	return float64(maxTime-minTime) / float64(len(sortTimes)-1)
}

func GetTxsPerBlock(blocks []*ethtypes.Block) []float64 {
	bs := SortableBlocks(blocks)
	sort.Sort(bs)

	txns := make([]float64, 0)
	for _, b := range bs {
		txns = append(txns, float64(len(b.Transactions())))
	}
	return txns
}
func GetUnclesPerBlock(blocks []*ethtypes.Block) []float64 {
	bs := SortableBlocks(blocks)
	sort.Sort(bs)

	uncles := make([]float64, 0)
	for _, b := range bs {
		uncles = append(uncles, float64(len(b.Uncles())))
	}
	return uncles
}

func GetSizePerBlock(blocks []*ethtypes.Block) []float64 {
	bs := SortableBlocks(blocks)
	sort.Sort(bs)

	bSize := make([]float64, 0)
	for _, b := range bs {
		bSize = append(bSize, float64(b.Size()))
	}
	return bSize
}
func GetGasPerBlock(blocks []*ethtypes.Block) []float64 {
	bs := SortableBlocks(blocks)
	sort.Sort(bs)

	gasUsed := make([]float64, 0)
	for _, b := range bs {
		gasUsed = append(gasUsed, float64(b.GasUsed()))
	}
	return gasUsed
}

func GetMeanGasPricePerBlock(blocks []*ethtypes.Block) []float64 {
	bs := SortableBlocks(blocks)
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

func GetSimpleBlockRecords(blocks []*ethtypes.Block) [][]string {
	bs := SortableBlocks(blocks)
	sort.Sort(bs)

	header := []string{
		"Block #",
		"Timestamp",
		"Block Hash",
		"Author",
		"Tx Count",
		"Gas Used",
	}

	if len(blocks) < 1 {
		return [][]string{header}
	}

	isMined := true

	blockHeader := blocks[0].Header()

	if blockHeader.Coinbase.String() == "0x0000000000000000000000000000000000000000" {
		isMined = false
	}

	if !isMined {
		header[3] = "Signer"
	}

	records := make([][]string, 0)
	records = append(records, header)
	for j := len(bs) - 1; j >= 0; j = j - 1 {
		author := bs[j].Header().Coinbase.String()
		ts := bs[j].Time()
		ut := time.Unix(int64(ts), 0)
		if !isMined {
			signer, err := ecrecover(bs[j].Header())
			if err == nil {
				author = "0x" + hex.EncodeToString(signer)
			}
		}
		record := []string{
			fmt.Sprintf("%d", bs[j].Number()),
			ut.Format("02 Jan 06 15:04:05 MST"),
			bs[j].Hash().String(),
			author,
			fmt.Sprintf("%d", len(bs[j].Transactions())),
			fmt.Sprintf("%d", bs[j].GasUsed()),
		}
		records = append(records, record)
	}
	return records
}

func GetSimpleBlockFields(block *ethtypes.Block) []string {
	ts := block.Time()
	ut := time.Unix(int64(ts), 0)

	blockHeader := block.Header()
	author := "Mined  by"

	authorAddress := blockHeader.Coinbase.String()

	if authorAddress == "0x0000000000000000000000000000000000000000" {
		author = "Signed by"
		signer, _ := ecrecover(blockHeader)
		authorAddress = hex.EncodeToString(signer)
	}

	return []string{
		"",
		fmt.Sprintf("Block Height: %s", block.Number()),
		fmt.Sprintf("Timestamp:    %d (%s)", ts, ut.Format(time.RFC3339)),
		fmt.Sprintf("Transactions: %d", len(block.Transactions())),
		fmt.Sprintf("%s:    %s", author, authorAddress),
		fmt.Sprintf("Difficulty:   %s", block.Difficulty()),
		fmt.Sprintf("Size:         %s", block.Size()),
		fmt.Sprintf("Uncles:       %d", len(block.Uncles())),
		fmt.Sprintf("Gas used:     %d", block.GasUsed()),
		fmt.Sprintf("Gas limit:    %d", block.GasLimit()),
		fmt.Sprintf("Base Fee:     %s", block.BaseFee()),
		fmt.Sprintf("Extra data:   %s", string(block.Extra())),
		fmt.Sprintf("Hash:         %s", block.Hash()),
		fmt.Sprintf("Parent Hash:  %s", block.ParentHash()),
		fmt.Sprintf("Uncle Hash:   %s", block.UncleHash()),
		fmt.Sprintf("State Root:   %s", block.Root()),
		fmt.Sprintf("Tx Hash:      %s", block.TxHash()),
		fmt.Sprintf("Nonce:        %d", block.Nonce()),
	}
}

func ecrecover(header *ethtypes.Header) ([]byte, error) {
	signature := header.Extra[len(header.Extra)-ethcrypto.SignatureLength:]
	pubkey, err := ethcrypto.Ecrecover(clique.SealHash(header).Bytes(), signature)
	if err != nil {
		return nil, err
	}
	signer := ethcrypto.Keccak256(pubkey[1:])[12:]

	return signer, nil
}
