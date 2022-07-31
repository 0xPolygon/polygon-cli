/*
Copyright © 2022 Polygon <engineering@polygon.technology>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"math/big"
	"net/url"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/maticnetwork/polygon-cli/jsonrpc"
	"github.com/maticnetwork/polygon-cli/metrics"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type (
	monitorStatus struct {
		ChainID   *big.Int
		HeadBlock *big.Int
		PeerCount uint64
		GasPrice  *big.Int

		Blocks            map[string]*jsonrpc.RawBlockResponse
		MaxBlockRetrieved *big.Int
	}
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor [rpc-url]",
	Short: "A simple terminal monitor for a blockchain",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := jsonrpc.NewClient()
		ms := new(monitorStatus)
		_, err := c.MakeRequestBatch(args[0], []string{"eth_blockNumber", "net_version", "net_peerCount", "eth_gasPrice"}, [][]any{nil, nil, nil, nil})
		if err != nil {
			return err
		}
		ms.MaxBlockRetrieved = big.NewInt(0)
		ms.Blocks = make(map[string]*jsonrpc.RawBlockResponse, 0)
		ms.ChainID = big.NewInt(0)

		isUiRendered := false
		errChan := make(chan error)
		go func() {
			for {
				resps, err := c.MakeRequestBatch(args[0], []string{"eth_blockNumber", "net_version", "net_peerCount", "eth_gasPrice"}, [][]any{nil, nil, nil, nil})
				if err != nil {
					log.Error().Err(err).Msg("Encountered issue fetching network information")
					continue
				}

				ms.HeadBlock = jsonrpc.MustConvHexToBigInt(resps[0].Result)
				ms.ChainID.SetString(resps[1].Result.(string), 10)
				ms.PeerCount = jsonrpc.MustConvHexToUint64(resps[2].Result)
				ms.GasPrice = jsonrpc.MustConvHexToBigInt(resps[3].Result)

				from := big.NewInt(0)

				// if the max block is 0, meaning we haven't fetched any blocks, we're going to start with head - 25
				if ms.MaxBlockRetrieved.Cmp(from) == 0 {
					from.Sub(ms.HeadBlock, big.NewInt(25))
				} else {
					from = ms.MaxBlockRetrieved
				}
				ms.getBlockRange(from, ms.HeadBlock, c, args[0])
				if !isUiRendered {
					go func() {
						errChan <- renderMonitorUI(ms)
					}()

				}
				time.Sleep(5 * time.Second)
			}

		}()

		err = <-errChan
		return err
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Expected exactly one argument")
		}
		url, err := url.Parse(args[0])
		if err != nil {
			log.Error().Err(err).Msg("Unable to parse url input error")
			return err
		}
		if url.Scheme != "http" && url.Scheme != "https" {
			return fmt.Errorf("The scheme %s is not supported", url.Scheme)
		}
		return nil
	},
}

func (ms *monitorStatus) getBlockRange(from, to *big.Int, c *jsonrpc.Client, url string) (any, error) {
	one := big.NewInt(1)
	methods := make([]string, 0)
	params := make([][]any, 0)
	for i := from; i.Cmp(to) != 1; i.Add(i, one) {
		methods = append(methods, "eth_getBlockByNumber")
		params = append(params, []any{"0x" + i.Text(16), true})
	}
	var resps []jsonrpc.RPCBlockResp
	err := c.MakeRequestBatchGenric(url, methods, params, &resps)
	if err != nil {
		return nil, err
	}
	for _, r := range resps {
		block := r.Result
		ms.Blocks[string(block.Number)] = &block
		bi, err := jsonrpc.ConvHexToBigInt(block.Number)
		if err != nil {
			// unclear why this would happen
			log.Error().Err(err).Msg("Could not convert block number")
		}
		if ms.MaxBlockRetrieved.Cmp(bi) == -1 {
			ms.MaxBlockRetrieved = bi

		}
	}
	return nil, nil
}

func init() {
	rootCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func renderMonitorUI(ms *monitorStatus) error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	paintMonitor(ms)
	currentBn := ms.HeadBlock
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(2 * time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID { // event string/identifier
			case "q", "<C-c>": // press 'q' or 'C-c' to quit
				return nil
			}
			switch e.Type {
			case ui.ResizeEvent:
				paintMonitor(ms)
				break

			}
		// use Go's built-in tickers for updating and drawing data
		case <-ticker:
			if currentBn != ms.HeadBlock {
				paintMonitor(ms)
				currentBn = ms.HeadBlock
			}

		}
	}
	return nil
}

func paintMonitor(ms *monitorStatus) error {
	blocks := make([]jsonrpc.RawBlockResponse, 0)
	for _, b := range ms.Blocks {
		blocks = append(blocks, *b)
	}

	table3 := widgets.NewTable()
	table3.Rows = metrics.GetSimpleBlockRecords(blocks)
	table3.TextStyle = ui.NewStyle(ui.ColorWhite)
	table3.RowSeparator = true

	columnWidths := make([]int, 5)
	if len(columnWidths) != len(table3.Rows[0]) {
		// i've messed up
		return fmt.Errorf("Something went wrong with table setup")
	}

	table3.ColumnResizer = func() {
		defaultWidth := (table3.Inner.Dx() - (12*3 + 42)) / 1
		columnWidths[0] = 12
		columnWidths[1] = defaultWidth
		columnWidths[2] = 42
		columnWidths[3] = 12
		columnWidths[4] = 12
	}

	table3.ColumnWidths = columnWidths

	h0 := widgets.NewParagraph()
	h0.Text = ms.HeadBlock.String()
	h0.Title = "Current Height"

	h1 := widgets.NewParagraph()
	h1.Text = fmt.Sprintf("%s wei", ms.GasPrice.String())
	h1.Title = "Gas Price"

	h2 := widgets.NewParagraph()
	h2.Text = fmt.Sprintf("%d", ms.PeerCount)
	h2.Title = "Current Peers"

	h3 := widgets.NewParagraph()
	h3.Text = ms.ChainID.String()
	h3.Title = "Chain ID"

	h4 := widgets.NewParagraph()
	h4.Text = fmt.Sprintf("%0.2f", metrics.GetMeanBlockTime(blocks))
	h4.Title = "Avg Block Time"

	sl0 := widgets.NewSparkline()
	sl0.Data = metrics.GetTxsPerBlock(blocks)
	sl0.LineColor = ui.ColorRed
	slg0 := widgets.NewSparklineGroup(sl0)
	slg0.Title = "TXs / Block"

	sl1 := widgets.NewSparkline()
	sl1.Data = metrics.GetMeanGasPricePerBlock(blocks)
	sl1.LineColor = ui.ColorGreen
	slg1 := widgets.NewSparklineGroup(sl1)
	slg1.Title = "Gas Price"

	sl2 := widgets.NewSparkline()
	sl2.Data = metrics.GetSizePerBlock(blocks)
	sl2.LineColor = ui.ColorYellow
	slg2 := widgets.NewSparklineGroup(sl2)
	slg2.Title = "Block Size"

	sl3 := widgets.NewSparkline()
	sl3.Data = metrics.GetUnclesPerBlock(blocks)
	sl3.LineColor = ui.ColorBlue
	slg3 := widgets.NewSparklineGroup(sl3)
	slg3.Title = "Uncles"

	sl4 := widgets.NewSparkline()
	sl4.Data = metrics.GetGasPerBlock(blocks)
	sl4.LineColor = ui.ColorMagenta
	slg4 := widgets.NewSparklineGroup(sl4)
	slg4.Title = "Gas Used"

	p0 := widgets.NewParagraph()
	p0.Text = ""
	p0.Title = "no impl"

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/8,
			ui.NewCol(1.0/5, h0),
			ui.NewCol(1.0/5, h1),
			ui.NewCol(1.0/5, h2),
			ui.NewCol(1.0/5, h3),
			ui.NewCol(1.0/5, h4),
		),

		ui.NewRow(3.0/8,
			ui.NewCol(1.0/5, slg0),
			ui.NewCol(1.0/5, slg1),
			ui.NewCol(1.0/5, slg2),
			ui.NewCol(1.0/5, slg3),
			ui.NewCol(1.0/5, slg4),
		),
		ui.NewRow(1.0/2, table3),
	)

	ui.Render(grid)
	return nil
}
