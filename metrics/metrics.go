package metrics

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/clique"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/maticnetwork/polygon-cli/rpctypes"
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
	SortableBlocks []rpctypes.PolyBlock
)

func (a SortableBlocks) Len() int {
	return len(a)
}
func (a SortableBlocks) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SortableBlocks) Less(i, j int) bool {
	return a[i].Number().Int64() < a[j].Number().Int64()
}

func GetMeanBlockTime(blocks []rpctypes.PolyBlock) float64 {
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

func GetTxsPerBlock(blocks []rpctypes.PolyBlock) []float64 {
	bs := SortableBlocks(blocks)
	sort.Sort(bs)

	txns := make([]float64, 0)
	for _, b := range bs {
		txns = append(txns, float64(len(b.Transactions())))
	}
	return txns
}
func GetUnclesPerBlock(blocks []rpctypes.PolyBlock) []float64 {
	bs := SortableBlocks(blocks)
	sort.Sort(bs)

	uncles := make([]float64, 0)
	for _, b := range bs {
		uncles = append(uncles, float64(len(b.Uncles())))
	}
	return uncles
}

func GetSizePerBlock(blocks []rpctypes.PolyBlock) []float64 {
	bs := SortableBlocks(blocks)
	sort.Sort(bs)

	bSize := make([]float64, 0)
	for _, b := range bs {
		bSize = append(bSize, float64(b.Size()))
	}
	return bSize
}
func GetGasPerBlock(blocks []rpctypes.PolyBlock) []float64 {
	bs := SortableBlocks(blocks)
	sort.Sort(bs)

	gasUsed := make([]float64, 0)
	for _, b := range bs {
		gasUsed = append(gasUsed, float64(b.GasUsed()))
	}
	return gasUsed
}

func GetMeanGasPricePerBlock(blocks []rpctypes.PolyBlock) []float64 {
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

func GetSimpleBlockRecords(blocks []rpctypes.PolyBlock) ([]string, string) {
	bs := SortableBlocks(blocks)
	sort.Sort(bs)

	// if we ever choose to utilize terminal width for column resizing
	// width, _, err := term.GetSize(0)
	// if err != nil {
	// 	return []string{}
	// }

	headerVariables := []string{"Block #", "Timestamp", "Block Time", "Tx Count", "Gas Used", "Block Hash", "Author"}

	proportion := []int{10, 20, 5, 5, 5, 60}

	header := ""
	for i, prop := range proportion {
		header += headerVariables[i] + strings.Repeat("─", prop)
	}
	header += headerVariables[len(headerVariables)-1]

	if len(blocks) < 1 {
		return nil, header
	}

	isMined := true

	if blocks[0].Miner().String() == "0x0000000000000000000000000000000000000000" {
		isMined = false
	}

	if !isMined {
		header = strings.Replace(header, "Author", "Signer", 1)
	}

	// Set the first row to blank so that there is some space between the blocks
	// and the title.
	records := []string{""}

	for j := len(bs) - 1; j >= 0; j = j - 1 {
		author := bs[j].Miner()
		ts := bs[j].Time()
		ut := time.Unix(int64(ts), 0)
		if !isMined {
			signer, err := ecrecover(&bs[j])
			if err == nil {
				author = ethcommon.HexToAddress("0x" + hex.EncodeToString(signer))
			}
		}
		blockTime := "-"
		if j > 0 {
			blockTime = strconv.FormatUint(bs[j].Time()-bs[j-1].Time(), 10)
		}

		recordVariables := []string{
			fmt.Sprintf("%d", bs[j].Number()),
			ut.Format("02 Jan 06 15:04:05 MST"),
			blockTime,
			fmt.Sprintf("%d", len(bs[j].Transactions())),
			fmt.Sprintf("%d", bs[j].GasUsed()),
			bs[j].Hash().String(),
			author.String(),
		}

		record := " "
		for i := 0; i < len(recordVariables)-1; i++ {
			record += recordVariables[i] + strings.Repeat(" ", len(headerVariables[i])+proportion[i]-len(recordVariables[i]))
		}
		record += recordVariables[len(recordVariables)-1]

		records = append(records, record)
	}
	return records, header
}

func GetSimpleBlockFields(block rpctypes.PolyBlock) []string {
	ts := block.Time()
	ut := time.Unix(int64(ts), 0)

	author := "Mined  by"

	authorAddress := block.Miner().String()
	if authorAddress == "0x0000000000000000000000000000000000000000" {
		author = "Signed by"
		signer, err := ecrecover(&block)
		if err == nil {
			authorAddress = hex.EncodeToString(signer)
		}

	}

	return []string{
		"",
		fmt.Sprintf("Block Height: %s", block.Number()),
		fmt.Sprintf("Timestamp:    %d (%s)", ts, ut.Format(time.RFC3339)),
		fmt.Sprintf("Transactions: %d", len(block.Transactions())),
		fmt.Sprintf("%s:    %s", author, authorAddress),
		fmt.Sprintf("Difficulty:   %s", block.Difficulty()),
		fmt.Sprintf("Size:         %d", block.Size()),
		fmt.Sprintf("Uncles:       %d", len(block.Uncles())),
		fmt.Sprintf("Gas used:     %d", block.GasUsed()),
		fmt.Sprintf("Gas limit:    %d", block.GasLimit()),
		fmt.Sprintf("Base Fee:     %s", block.BaseFee()),
		fmt.Sprintf("Extra data:   %s", RawDataToASCII(block.Extra())),
		fmt.Sprintf("Hash:         %s", block.Hash()),
		fmt.Sprintf("Parent Hash:  %s", block.ParentHash()),
		fmt.Sprintf("Uncle Hash:   %s", block.UncleHash()),
		fmt.Sprintf("State Root:   %s", block.Root()),
		fmt.Sprintf("Tx Hash:      %s", block.TxHash()),
		fmt.Sprintf("Nonce:        %d", block.Nonce()),
	}
}
func GetSimpleBlockTxFields(block rpctypes.PolyBlock, chainID *big.Int) []string {
	fields := make([]string, 0)
	blank := ""
	for _, tx := range block.Transactions() {
		txFields := GetSimpleTxFields(tx, chainID, block.BaseFee())
		fields = append(fields, blank)
		fields = append(fields, txFields...)
	}
	return fields
}
func GetSimpleTxFields(tx rpctypes.PolyTransaction, chainID, baseFee *big.Int) []string {
	fields := make([]string, 0)
	fields = append(fields, fmt.Sprintf("Tx Hash: %s", tx.Hash()))

	txMethod := "Transfer"
	if tx.To().String() == "0x0000000000000000000000000000000000000000" {
		// Contract deployment
		txMethod = "Contract Deployment"
	} else if len(tx.Data()) > 4 {
		// Contract call
		txMethod = hex.EncodeToString(tx.Data()[0:4])
	}

	fields = append(fields, fmt.Sprintf("To: %s", tx.To()))
	fields = append(fields, fmt.Sprintf("From: %s", tx.From()))
	fields = append(fields, fmt.Sprintf("Method: %s", txMethod))
	fields = append(fields, fmt.Sprintf("Value: %s", tx.Value()))
	fields = append(fields, fmt.Sprintf("Gas Limit: %d", tx.Gas()))
	fields = append(fields, fmt.Sprintf("Gas Price: %s", tx.GasPrice()))
	fields = append(fields, fmt.Sprintf("Gas Tip: %d", tx.MaxPriorityFeePerGas()))
	fields = append(fields, fmt.Sprintf("Gas Fee: %d", tx.MaxFeePerGas()))
	fields = append(fields, fmt.Sprintf("Nonce: %d", tx.Nonce()))
	fields = append(fields, fmt.Sprintf("Type: %d", tx.Type()))
	fields = append(fields, fmt.Sprintf("Data Len: %d", len(tx.Data())))
	fields = append(fields, fmt.Sprintf("Data: %s", hex.EncodeToString(tx.Data())))

	return fields
}

func ecrecover(block *rpctypes.PolyBlock) ([]byte, error) {
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
	pubkey, err := ethcrypto.Ecrecover(clique.SealHash(header).Bytes(), signature)
	if err != nil {
		return nil, err
	}
	signer := ethcrypto.Keccak256(pubkey[1:])[12:]

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
