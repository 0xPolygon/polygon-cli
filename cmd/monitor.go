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
	"context"
	"fmt"
	"math/big"
	"net/url"
	"sort"
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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

		Blocks            map[string]*ethtypes.Block
		MaxBlockRetrieved *big.Int
	}
	chainState struct {
		HeadBlock uint64
		ChainID   *big.Int
		PeerCount uint64
		GasPrice  *big.Int
	}
	monitorMode int
)

const (
	monitorModeHelp     monitorMode = iota
	monitorModeExplorer monitorMode = iota
	monitorModeBlock    monitorMode = iota
)

func getChainState(ctx context.Context, ec *ethclient.Client) (*chainState, error) {
	var err error
	cs := new(chainState)
	cs.HeadBlock, err = ec.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	cs.ChainID, err = ec.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	cs.PeerCount, err = ec.PeerCount(ctx)
	if err != nil {
		// log.Info().Err(err).Msg("Using fake peer count")
		cs.PeerCount = 0
	}

	cs.GasPrice, err = ec.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	return cs, nil

}

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor [rpc-url]",
	Short: "A simple terminal monitor for a blockchain",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		rpc, err := ethrpc.DialContext(ctx, args[0])
		if err != nil {
			log.Error().Err(err).Msg("Unable to dial rpc")
			return err
		}
		ec := ethclient.NewClient(rpc)

		ms := new(monitorStatus)

		ms.MaxBlockRetrieved = big.NewInt(0)
		ms.Blocks = make(map[string]*ethtypes.Block, 0)
		ms.ChainID = big.NewInt(0)

		isUiRendered := false
		errChan := make(chan error)
		go func() {
			for {
				cs, err := getChainState(ctx, ec)
				if err != nil {
					log.Error().Err(err).Msg("Encountered issue fetching network information")
					time.Sleep(5 * time.Second)
					continue
				}

				ms.HeadBlock = new(big.Int).SetUint64(cs.HeadBlock)
				ms.ChainID = cs.ChainID
				ms.PeerCount = cs.PeerCount
				ms.GasPrice = cs.GasPrice

				from := big.NewInt(0)

				// if the max block is 0, meaning we haven't fetched any blocks, we're going to start with head - 25
				if ms.MaxBlockRetrieved.Cmp(from) == 0 {
					from.Sub(ms.HeadBlock, big.NewInt(25))
				} else {
					from = ms.MaxBlockRetrieved
				}
				// log.Trace().Str("from", from.String()).Str("to", ms.HeadBlock.String()).Msg("Fetching block range")
				ms.getBlockRange(ctx, from, ms.HeadBlock, ec, args[0])
				if !isUiRendered {
					go func() {
						errChan <- renderMonitorUI(ms)
					}()
					isUiRendered = true

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
		_, err := url.Parse(args[0])
		if err != nil {
			log.Error().Err(err).Msg("Unable to parse url input error")
			return err
		}
		return nil
	},
}

func (ms *monitorStatus) getBlockRange(ctx context.Context, from, to *big.Int, c *ethclient.Client, url string) error {
	one := big.NewInt(1)
	for i := from; i.Cmp(to) != 1; i.Add(i, one) {
		block, err := c.BlockByNumber(ctx, i)
		if err != nil {
			return err
		}

		ms.Blocks[block.Number().String()] = block

		if ms.MaxBlockRetrieved.Cmp(block.Number()) == -1 {
			ms.MaxBlockRetrieved = block.Number()
		}

	}
	return nil
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}

func renderMonitorUI(ms *monitorStatus) error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	currentMode := monitorModeExplorer

	blockTable := widgets.NewTable()

	blockTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	blockTable.RowSeparator = true

	columnWidths := make([]int, 6)

	blockTable.ColumnResizer = func() {
		defaultWidth := (blockTable.Inner.Dx() - (12 + 22 + 42 + 12 + 14)) / 1
		columnWidths[0] = 12
		columnWidths[1] = 22
		columnWidths[2] = defaultWidth
		columnWidths[3] = 42
		columnWidths[4] = 12
		columnWidths[5] = 14
	}

	blockTable.ColumnWidths = columnWidths

	h0 := widgets.NewParagraph()
	h0.Title = "Current Height"

	h1 := widgets.NewParagraph()
	h1.Title = "Gas Price"

	h2 := widgets.NewParagraph()
	h2.Title = "Current Peers"

	h3 := widgets.NewParagraph()
	h3.Title = "Chain ID"

	h4 := widgets.NewParagraph()
	h4.Title = "Avg Block Time"

	sl0 := widgets.NewSparkline()
	sl0.LineColor = ui.ColorRed
	slg0 := widgets.NewSparklineGroup(sl0)
	slg0.Title = "TXs / Block"

	sl1 := widgets.NewSparkline()
	sl1.LineColor = ui.ColorGreen
	slg1 := widgets.NewSparklineGroup(sl1)
	slg1.Title = "Gas Price"

	sl2 := widgets.NewSparkline()
	sl2.LineColor = ui.ColorYellow
	slg2 := widgets.NewSparklineGroup(sl2)
	slg2.Title = "Block Size"

	sl3 := widgets.NewSparkline()
	sl3.LineColor = ui.ColorBlue
	slg3 := widgets.NewSparklineGroup(sl3)
	slg3.Title = "Uncles"

	sl4 := widgets.NewSparkline()
	sl4.LineColor = ui.ColorMagenta
	slg4 := widgets.NewSparklineGroup(sl4)
	slg4.Title = "Gas Used"

	grid := ui.NewGrid()
	blockGrid := ui.NewGrid()

	b0 := widgets.NewParagraph()
	b0.Title = "Block Headers"
	b0.Text = "Foo"

	b1 := widgets.NewList()
	b1.Title = "Block Info"
	b1.TextStyle = ui.NewStyle(ui.ColorYellow)
	b1.WrapText = false

	b2 := widgets.NewParagraph()
	b2.Title = "Transactions"
	b2.Text = "Foooo"

	blockGrid.Set(
		ui.NewRow(1.0/10, b0),

		ui.NewRow(9.0/10,
			ui.NewCol(1.0/2, b1),
			ui.NewCol(1.0/2, b2),
		),
	)

	grid.Set(
		ui.NewRow(1.0/10,
			ui.NewCol(1.0/5, h0),
			ui.NewCol(1.0/5, h1),
			ui.NewCol(1.0/5, h2),
			ui.NewCol(1.0/5, h3),
			ui.NewCol(1.0/5, h4),
		),

		ui.NewRow(4.0/10,
			ui.NewCol(1.0/5, slg0),
			ui.NewCol(1.0/5, slg1),
			ui.NewCol(1.0/5, slg2),
			ui.NewCol(1.0/5, slg3),
			ui.NewCol(1.0/5, slg4),
		),
		ui.NewRow(3.0/5, blockTable),
	)

	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	blockGrid.SetRect(0, 0, termWidth, termHeight)

	var selectedBlockIdx *int
	var selectedBlock *ethtypes.Block

	redraw := func(ms *monitorStatus) {
		if currentMode == monitorModeHelp {
			// render a help
		} else if currentMode == monitorModeBlock {
			// render a block
			b1.Rows = metrics.GetSimpleBlockFields(selectedBlock)

			ui.Clear()
			ui.Render(blockGrid)
			return
		}

		//default
		blocks := make([]*ethtypes.Block, 0)
		for _, b := range ms.Blocks {
			blocks = append(blocks, b)
		}
		recentBlocks := metrics.SortableBlocks(blocks)
		sort.Sort(recentBlocks)
		// 25 needs to be a variable / parameter
		if len(recentBlocks) > 25 {
			recentBlocks = recentBlocks[len(recentBlocks)-25 : len(recentBlocks)]
		}

		h0.Text = ms.HeadBlock.String()
		gasGwei := new(big.Int)
		gasGwei.Div(ms.GasPrice, metrics.UnitShannon)
		h1.Text = fmt.Sprintf("%s gwei", gasGwei.String())
		h2.Text = fmt.Sprintf("%d", ms.PeerCount)
		h3.Text = ms.ChainID.String()
		h4.Text = fmt.Sprintf("%0.2f", metrics.GetMeanBlockTime(recentBlocks))

		sl0.Data = metrics.GetTxsPerBlock(recentBlocks)
		sl1.Data = metrics.GetMeanGasPricePerBlock(recentBlocks)
		sl2.Data = metrics.GetSizePerBlock(recentBlocks)
		sl3.Data = metrics.GetUnclesPerBlock(recentBlocks)
		sl4.Data = metrics.GetGasPerBlock(recentBlocks)

		// assuming we haven't selected a particular row... we should get new blocks
		if selectedBlockIdx == nil {
			blockTable.Rows = metrics.GetSimpleBlockRecords(recentBlocks)
		}
		if len(columnWidths) != len(blockTable.Rows[0]) {
			// i've messed up
			panic(fmt.Sprintf("Mis matched between columns and specified widths"))
		}

		for i := 0; i < len(blockTable.Rows); i = i + 1 {
			blockTable.RowStyles[i] = ui.NewStyle(ui.ColorWhite)
		}
		if selectedBlockIdx != nil && *selectedBlockIdx > 0 && *selectedBlockIdx < len(blockTable.Rows) {
			blockTable.RowStyles[*selectedBlockIdx] = ui.NewStyle(ui.ColorWhite, ui.ColorRed, ui.ModifierBold)
			// the block table is reversed and has an extra row for the header
			selectedBlock = recentBlocks[len(recentBlocks)-*selectedBlockIdx]
		}

		ui.Render(grid)
	}

	listDraw := func() {
		l := widgets.NewList()
		l.Title = "List"
		l.Rows = []string{
			"[0] github.com/gizak/termui/v3",
			"[1] [你好，世界](fg:blue)",
			"[2] [こんにちは世界](fg:red)",
			"[3] [color](fg:white,bg:green) output",
			"[4] output.go",
			"[5] random_out.go",
			"[6] dashboard.go",
			"[7] foo",
			"[8] bar",
			"[9] baz",
		}
		l.TextStyle = ui.NewStyle(ui.ColorYellow)
		l.WrapText = false
		l.SetRect(0, 0, 25, 8)

		ui.Render(l)

	}

	currentBn := ms.HeadBlock
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C

	redraw(ms)

	currIdx := 0
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return nil
			case "<Escape>":
				selectedBlockIdx = nil
				currentMode = monitorModeExplorer
				redraw(ms)
				break
			case "<Enter>":
				// TODO
				if selectedBlockIdx != nil {
					currentMode = monitorModeBlock
				}
				_ = listDraw
				redraw(ms)
				break
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				blockGrid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				redraw(ms)
				break
			case "<Up>", "<Down>", "<Left>", "<Right>":
				if selectedBlockIdx == nil {
					currIdx = 1
					selectedBlockIdx = &currIdx
					redraw(ms)
					break
				}
				currIdx = *selectedBlockIdx

				if e.ID == "<Down>" {
					currIdx = currIdx + 1
				} else if e.ID == "<Up>" {
					currIdx = currIdx - 1
				}
				if currIdx > 0 && currIdx < 25 { // need a better way to understand how many rows are visble
					selectedBlockIdx = &currIdx
				}
				redraw(ms)
				break
			case "<MouseLeft>", "<MouseRight>", "<MouseRelease>", "<MouseWheelUp>", "<MouseWheelDown>":
				break
			default:
				log.Trace().Str("id", e.ID).Msg("Unknown ui event")
			}
		case <-ticker:
			if currentBn != ms.HeadBlock {
				currentBn = ms.HeadBlock
				redraw(ms)
				break
			}
		}
	}

	return nil
}
