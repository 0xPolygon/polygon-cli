package ui

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/maticnetwork/polygon-cli/metrics"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/rs/zerolog/log"
)

type UiSkeleton struct {
	Current         *widgets.Paragraph
	TxPerBlockChart *widgets.Sparkline
	GasPriceChart   *widgets.Sparkline
	BlockSizeChart  *widgets.Sparkline
	PendingTxChart  *widgets.Sparkline
	GasChart        *widgets.Sparkline
	BlockInfo       *widgets.List
	TransactionList *widgets.List
}

func GetCurrentBlockInfo(headBlock *big.Int, gasPrice *big.Int, peerCount uint64, pendingCount uint64, chainID *big.Int, blocks []rpctypes.PolyBlock) string {
	height := fmt.Sprintf("Height: %s", headBlock.String())
	time := fmt.Sprintf("Time: %s", time.Now().Format("02 Jan 06 15:04:05 MST"))
	gasPriceString := fmt.Sprintf("Gas Price: %s gwei", new(big.Int).Div(gasPrice, metrics.UnitShannon).String())
	peers := fmt.Sprintf("Peers: %d", peerCount)
	pendingTx := fmt.Sprintf("Pending Tx: %d", pendingCount)
	chainIdString := fmt.Sprintf("Chain ID: %s", chainID.String())
	blockTime := fmt.Sprintf("Avg Block Time: %0.2f", metrics.GetMeanBlockTime(blocks))
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s", height, time, gasPriceString, peers, pendingTx, chainIdString, blockTime)

}

func GetBlocksList(blocks []rpctypes.PolyBlock) ([]string, string) {
	bs := rpctypes.SortableBlocks(blocks)
	sort.Sort(bs)

	// if we ever choose to utilize terminal width for column resizing
	// width, _, err := term.GetSize(0)
	// if err != nil {
	// 	return []string{}
	// }

	zone, _ := time.Now().Zone()
	headerVariables := []string{"#", fmt.Sprintf("TIME (%s)", zone), "BLK TIME", "TXN #", "GAS USED", "HASH", "AUTHOR"}

	proportion := []int{10, 10, 5, 5, 10, 20}

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
		header = strings.Replace(header, "AUTHOR", "SIGNER", 1)
	}

	// Set the first row to blank so that there is some space between the blocks
	// and the title.
	records := []string{""}

	for j := len(bs) - 1; j >= 0; j = j - 1 {
		author := bs[j].Miner()
		ts := bs[j].Time()
		ut := time.Unix(int64(ts), 0)
		if !isMined {
			signer, err := metrics.Ecrecover(&bs[j])
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
			ut.Format("02 Jan 06 15:04:05"),
			fmt.Sprintf("%ss", blockTime),
			fmt.Sprintf("%d", len(bs[j].Transactions())),
			fmt.Sprintf("%d", bs[j].GasUsed()),
			metrics.TruncateHexString(bs[j].Hash().String(), 14),
			metrics.TruncateHexString(author.String(), 14),
		}

		record := " "
		for i := 0; i < len(recordVariables)-1; i++ {
			spaceOffset := len(headerVariables[i]) + proportion[i] - len(recordVariables[i])
			if spaceOffset < 0 {
				spaceOffset = 0
				log.Error().Str("record", recordVariables[i]).Str("column", headerVariables[i]).Msg("Column width exceed header width")
			}
			record += recordVariables[i] + strings.Repeat(" ", spaceOffset)
		}
		record += recordVariables[len(recordVariables)-1]

		records = append(records, record)
	}
	return records, header
}

func GetSimpleBlockFields(block rpctypes.PolyBlock) []string {
	ts := block.Time()
	ut := time.Unix(int64(ts), 0)

	author := "Mined by"

	authorAddress := block.Miner().String()
	if authorAddress == "0x0000000000000000000000000000000000000000" {
		author = "Signed by"
		signer, err := metrics.Ecrecover(&block)
		if err == nil {
			authorAddress = hex.EncodeToString(signer)
		}

	}

	return []string{
		"",
		fmt.Sprintf("Block Height:      %s", block.Number()),
		fmt.Sprintf("Timestamp:         %d (%s)", ts, ut.Format(time.RFC3339)),
		fmt.Sprintf("Transactions:      %d", len(block.Transactions())),
		fmt.Sprintf("%s:         %s", author, authorAddress),
		fmt.Sprintf("Difficulty:        %s", block.Difficulty()),
		fmt.Sprintf("Size:              %d", block.Size()),
		fmt.Sprintf("Uncles:            %d", len(block.Uncles())),
		fmt.Sprintf("Gas used:          %d", block.GasUsed()),
		fmt.Sprintf("Gas limit:         %d", block.GasLimit()),
		fmt.Sprintf("Base Fee per gas:  %s", block.BaseFee()),
		fmt.Sprintf("Extra data:        %s", metrics.RawDataToASCII(block.Extra())),
		fmt.Sprintf("Hash:              %s", block.Hash()),
		fmt.Sprintf("Parent Hash:       %s", block.ParentHash()),
		fmt.Sprintf("Uncle Hash:        %s", block.UncleHash()),
		fmt.Sprintf("State Root:        %s", block.Root()),
		fmt.Sprintf("Tx Hash:           %s", block.TxHash()),
		fmt.Sprintf("Nonce:             %d", block.Nonce()),
	}
}

func GetBlockTxTable(block rpctypes.PolyBlock, chainID *big.Int) [][]string {
	fields := make([][]string, 0)
	header := []string{"Txn Hash", "Method", "Timestamp", "From", "To", "Value", "Gas Price"}
	fields = append(fields, header)
	for _, tx := range block.Transactions() {
		txFields := getTxTable(tx, chainID, block.BaseFee())
		fields = append(fields, txFields)
	}
	return fields
}

func GetTxMethod(tx rpctypes.PolyTransaction) string {
	txMethod := "Transfer"
	if tx.To().String() == "0x0000000000000000000000000000000000000000" {
		// Contract deployment
		txMethod = "Contract Deployment"
	} else if len(tx.Data()) > 4 {
		// Contract call
		txMethod = hex.EncodeToString(tx.Data()[0:4])
	}

	return txMethod
}

func getTxTable(tx rpctypes.PolyTransaction, chainID, baseFee *big.Int) []string {
	fields := make([]string, 0)
	fields = append(fields, fmt.Sprintf("%s", tx.Hash()))

	txMethod := GetTxMethod(tx)

	fields = append(fields, txMethod)
	fields = append(fields, "TIME")
	fields = append(fields, fmt.Sprintf("%s", tx.From()))
	fields = append(fields, fmt.Sprintf("%s", tx.To()))
	fields = append(fields, fmt.Sprintf("%s", tx.Value()))
	fields = append(fields, fmt.Sprintf("%s", tx.GasPrice()))

	return fields
}

func GetTransactionsList(block rpctypes.PolyBlock, chainID *big.Int) ([]string, string) {
	txs := block.Transactions()

	headerVariables := []string{"Txn Hash", "Method", "From", "To", "Value", "Gas Price"}
	proportion := []int{60, 5, 15, 15, 30}

	header := ""
	for i, prop := range proportion {
		header += headerVariables[i] + strings.Repeat("─", prop)
	}
	header += headerVariables[len(headerVariables)-1]

	records := []string{""}

	for _, tx := range txs {
		txMethod := GetTxMethod(tx)
		recordVariables := []string{
			fmt.Sprintf("%s", tx.Hash()),
			txMethod,
			metrics.TruncateHexString(fmt.Sprintf("%s", tx.From()), 14),
			metrics.TruncateHexString(fmt.Sprintf("%s", tx.To()), 14),
			fmt.Sprintf("%s", tx.Value()),
			fmt.Sprintf("%s", tx.GasPrice()),
		}

		record := " "
		for i := 0; i < len(recordVariables)-1; i++ {
			spaceOffset := len(headerVariables[i]) + proportion[i] - len(recordVariables[i])
			if spaceOffset < 0 {
				spaceOffset = 0
				log.Error().Str("record", recordVariables[i]).Str("column", headerVariables[i]).Msg("Column width exceed header width")
			}
			record += recordVariables[i] + strings.Repeat(" ", spaceOffset)
		}
		record += recordVariables[len(recordVariables)-1]

		records = append(records, record)
	}
	return records, header
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

func SetUISkeleton() (blockList *widgets.List, blockInfo *widgets.List, transactionInfo *widgets.Table, grid *ui.Grid, blockGrid *ui.Grid, termUi UiSkeleton) {
	// help := widgets.NewParagraph()
	// help.Title = "Block Headers"
	// help.Text = "Use the arrow keys to scroll through the transactions. Press <Esc> to go back to the explorer view"

	blockList = widgets.NewList()
	blockList.TextStyle = ui.NewStyle(ui.ColorWhite)

	blockInfo = widgets.NewList()
	blockInfo.TextStyle = ui.NewStyle(ui.ColorWhite)
	blockInfo.Title = "Block Information"

	transactionInfo = widgets.NewTable()
	transactionInfo.TextStyle = ui.NewStyle(ui.ColorWhite)
	transactionInfo.FillRow = true
	transactionInfo.Title = "Latest Transactions"
	transactionInfo.Rows = [][]string{{""}, {""}}

	termUi = UiSkeleton{}

	termUi.Current = widgets.NewParagraph()
	termUi.Current.Title = "Current"

	termUi.TxPerBlockChart = widgets.NewSparkline()
	termUi.TxPerBlockChart.LineColor = ui.ColorRed
	termUi.TxPerBlockChart.MaxHeight = 1000
	slg0 := widgets.NewSparklineGroup(termUi.TxPerBlockChart)
	slg0.Title = "TXs / Block"

	termUi.GasPriceChart = widgets.NewSparkline()
	termUi.GasPriceChart.LineColor = ui.ColorGreen
	termUi.GasPriceChart.MaxHeight = 1000
	slg1 := widgets.NewSparklineGroup(termUi.GasPriceChart)
	slg1.Title = "Gas Price"

	termUi.BlockSizeChart = widgets.NewSparkline()
	termUi.BlockSizeChart.LineColor = ui.ColorYellow
	termUi.BlockSizeChart.MaxHeight = 1000
	slg2 := widgets.NewSparklineGroup(termUi.BlockSizeChart)
	slg2.Title = "Block Size"

	termUi.PendingTxChart = widgets.NewSparkline()
	termUi.PendingTxChart.LineColor = ui.ColorBlue
	termUi.PendingTxChart.MaxHeight = 1000
	slg3 := widgets.NewSparklineGroup(termUi.PendingTxChart)
	slg3.Title = "Pending Tx"

	termUi.GasChart = widgets.NewSparkline()
	termUi.GasChart.LineColor = ui.ColorMagenta
	termUi.GasChart.MaxHeight = 1000
	slg4 := widgets.NewSparklineGroup(termUi.GasChart)
	slg4.Title = "Gas Used"

	grid = ui.NewGrid()
	blockGrid = ui.NewGrid()

	// b0 := widgets.NewParagraph()
	// b0.Title = "Block Headers"
	// b0.Text = "Use the arrow keys to scroll through the transactions. Press <Esc> to go back to the explorer view"

	termUi.BlockInfo = widgets.NewList()
	termUi.BlockInfo.Title = "Block Info"
	termUi.BlockInfo.TextStyle = ui.NewStyle(ui.ColorYellow)
	termUi.BlockInfo.WrapText = false

	termUi.TransactionList = widgets.NewList()
	termUi.TransactionList.Title = "Transactions"
	termUi.TransactionList.TextStyle = ui.NewStyle(ui.ColorGreen)
	termUi.TransactionList.WrapText = true

	blockGrid.Set(
		// ui.NewRow(1.0/10, b0),
		ui.NewRow(2.0/10, termUi.BlockInfo),
		ui.NewRow(8.0/10, termUi.TransactionList),
	)

	grid.Set(
		ui.NewRow(1.0/10, termUi.Current),

		ui.NewRow(2.0/10,
			ui.NewCol(1.0/5, slg0),
			ui.NewCol(1.0/5, slg1),
			ui.NewCol(1.0/5, slg2),
			ui.NewCol(1.0/5, slg3),
			ui.NewCol(1.0/5, slg4),
		),

		ui.NewRow(5.0/10,
			ui.NewCol(3.0/5, blockList),
			ui.NewCol(2.0/5, blockInfo),
		),

		ui.NewRow(2.0/10,
			ui.NewCol(5.0/5, transactionInfo),
		),
	)

	return
}
