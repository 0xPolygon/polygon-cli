package metrics

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/consensus/clique"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/maticnetwork/polygon-cli/jsonrpc"
)

type (
	BlockSlice []jsonrpc.RawBlockResponse
)

func (a BlockSlice) Len() int {
	return len(a)
}
func (a BlockSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a BlockSlice) Less(i, j int) bool {
	return a[i].Timestamp.ToUint64() < a[j].Timestamp.ToUint64()
}

func GetMeanBlockTime(blocks []jsonrpc.RawBlockResponse) float64 {
	if len(blocks) < 2 {
		return 0
	}

	times := make([]int, 0)
	for _, block := range blocks {
		blockTime := jsonrpc.MustConvHexToUint64(block.Timestamp)
		times = append(times, int(blockTime))
	}

	sortTimes := sort.IntSlice(times)
	sortTimes.Sort()

	minTime := sortTimes[0]
	maxTime := sortTimes[len(sortTimes)-1]

	return float64(maxTime-minTime) / float64(len(sortTimes)-1)
}

func GetTxsPerBlock(blocks []jsonrpc.RawBlockResponse) []float64 {
	bs := BlockSlice(blocks)
	sort.Sort(bs)

	txns := make([]float64, 0)
	for _, b := range bs {
		txns = append(txns, float64(len(b.Transactions)))
	}
	return txns
}
func GetUnclesPerBlock(blocks []jsonrpc.RawBlockResponse) []float64 {
	bs := BlockSlice(blocks)
	sort.Sort(bs)

	uncles := make([]float64, 0)
	for _, b := range bs {
		uncles = append(uncles, float64(len(b.Uncles)))
	}
	return uncles
}

func GetSizePerBlock(blocks []jsonrpc.RawBlockResponse) []float64 {
	bs := BlockSlice(blocks)
	sort.Sort(bs)

	bSize := make([]float64, 0)
	for _, b := range bs {
		bSize = append(bSize, float64(b.Size.ToUint64()))
	}
	return bSize
}
func GetGasPerBlock(blocks []jsonrpc.RawBlockResponse) []float64 {
	bs := BlockSlice(blocks)
	sort.Sort(bs)

	gasUsed := make([]float64, 0)
	for _, b := range bs {
		gasUsed = append(gasUsed, float64(b.GasUsed.ToUint64()))
	}
	return gasUsed
}

func GetMeanGasPricePerBlock(blocks []jsonrpc.RawBlockResponse) []float64 {
	bs := BlockSlice(blocks)
	sort.Sort(bs)

	gasPrices := make([]float64, 0)
	for _, b := range bs {
		var totGas uint64 = 0
		for _, tx := range b.Transactions {
			totGas += tx.GasPrice.ToUint64()
		}
		meanGas := float64(totGas) / float64(len(b.Transactions))
		gasPrices = append(gasPrices, meanGas)
	}
	return gasPrices
}

func GetSimpleBlockRecords(blocks []jsonrpc.RawBlockResponse) [][]string {
	bs := BlockSlice(blocks)
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

	if string(blocks[0].Miner) == "0x0000000000000000000000000000000000000000" {
		isMined = false
	}

	if !isMined {
		header[3] = "Signer"
	}

	records := make([][]string, 0)
	records = append(records, header)
	for j := len(bs) - 1; j >= 0; j = j - 1 {
		author := string(bs[j].Miner)
		ts := bs[j].Timestamp.ToInt64()
		ut := time.Unix(ts, 0)
		if !isMined {
			signer, err := getBlockSigner(bs[j])
			if err == nil {
				author = "0x" + hex.EncodeToString(signer)
			}
		}
		record := []string{
			fmt.Sprintf("%d", bs[j].Number.ToUint64()),
			ut.Format(time.RFC822),
			string(bs[j].Hash),
			author,
			fmt.Sprintf("%d", len(bs[j].Transactions)),
			fmt.Sprintf("%d", bs[j].GasUsed.ToUint64()),
		}
		records = append(records, record)
	}
	return records
}

func getBlockSigner(b jsonrpc.RawBlockResponse) ([]byte, error) {
	var h *ethtypes.Header = new(ethtypes.Header)
	// the common interface here is json, so I'm going to convert from my raw type to to a smart geth type for clique validtor information
	jsonData, _ := json.Marshal(b)
	err := h.UnmarshalJSON(jsonData)
	if err != nil {
		return nil, err
	}
	return ecrecover(h)
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
