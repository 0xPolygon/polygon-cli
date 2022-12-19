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
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/cenkalti/backoff"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/maticnetwork/polygon-cli/metrics"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var inputBatchSize *uint64
var verbosity *int64

type (
	monitorStatus struct {
		ChainID   *big.Int
		HeadBlock *big.Int
		PeerCount uint64
		GasPrice  *big.Int

		Blocks            map[string]rpctypes.PolyBlock
		BlocksLock        sync.RWMutex
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

	defaultBatchSize uint64 = 25
)

func getChainState(ctx context.Context, ec *ethclient.Client) (*chainState, error) {
	var err error
	cs := new(chainState)
	cs.HeadBlock, err = ec.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch block number: %s", err.Error())
	}

	cs.ChainID, err = ec.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch chain id: %s", err.Error())
	}

	cs.PeerCount, err = ec.PeerCount(ctx)
	if err != nil {
		// log.Info().Err(err).Msg("Using fake peer count")
		cs.PeerCount = 0
	}

	cs.GasPrice, err = ec.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't estimate gas: %s", err.Error())
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
		ms.BlocksLock.Lock()
		ms.Blocks = make(map[string]rpctypes.PolyBlock, 0)
		ms.BlocksLock.Unlock()
		ms.ChainID = big.NewInt(0)
		zero := big.NewInt(0)

		isUiRendered := false
		errChan := make(chan error)
		go func() {
			for {
				var cs *chainState
				cs, err = getChainState(ctx, ec)
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

				// batchSize := *inputBatchSize
				// if *inputBatchSize > 0 {
				// 	batchSize = *inputBatchSize
				// }

				// if the max block is 0, meaning we haven't fetched any blocks, we're going to start with head - batchSize
				if ms.MaxBlockRetrieved.Cmp(from) == 0 {
					headBlockMinusBatchSize := new(big.Int).SetUint64(*inputBatchSize + 100 - 1)
					from.Sub(ms.HeadBlock, headBlockMinusBatchSize)
				} else {
					from = ms.MaxBlockRetrieved
				}

				if from.Cmp(zero) < 0 {
					from.SetInt64(0)
				}
				err = ms.getBlockRange(ctx, from, ms.HeadBlock, rpc, args[0])
				if err != nil {
					log.Error().Err(err).Msg("there was an issue fetching the block range")
				}
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
			return fmt.Errorf("too many arguments")
		}

		// validate url argument
		_, err := url.Parse(args[0])
		if err != nil {
			log.Error().Err(err).Msg("Unable to parse url input error")
			return err
		}

		// validate batch-size flag
		if *inputBatchSize == 0 {
			return fmt.Errorf("batch-size can't be equal to zero")
		}

		setMonitorLogLevel(*verbosity)

		return nil
	},
}

func (ms *monitorStatus) getBlockRange(ctx context.Context, from, to *big.Int, c *ethrpc.Client, url string) error {
	one := big.NewInt(1)
	blms := make([]ethrpc.BatchElem, 0)
	for i := from; i.Cmp(to) != 1; i.Add(i, one) {
		r := new(rpctypes.RawBlockResponse)
		var err error
		blms = append(blms, ethrpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{"0x" + i.Text(16), true},
			Result: r,
			Error:  err,
		})
	}
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 3 * time.Minute
	retryable := func() error {
		err := c.BatchCallContext(ctx, blms)
		return err
	}
	err := backoff.Retry(retryable, b)
	if err != nil {
		return err
	}
	for _, b := range blms {
		if b.Error != nil {
			return b.Error
		}
		pb := rpctypes.NewPolyBlock(b.Result.(*rpctypes.RawBlockResponse))

		ms.BlocksLock.Lock()
		ms.Blocks[pb.Number().String()] = pb
		ms.BlocksLock.Unlock()

		if ms.MaxBlockRetrieved.Cmp(pb.Number()) == -1 {
			ms.MaxBlockRetrieved = pb.Number()
		}

	}
	return nil
}

func init() {
	rootCmd.AddCommand(monitorCmd)

	inputBatchSize = monitorCmd.PersistentFlags().Uint64P("batch-size", "b", defaultBatchSize, "Number of requests per batch")
	verbosity = monitorCmd.PersistentFlags().Int64P("verbosity", "v", 200, "0 - Silent\n100 Fatals\n200 Errors\n300 Warnings\n400 INFO\n500 Debug\n600 Trace")
}

func renderMonitorUI(ms *monitorStatus) error {
	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	currentMode := monitorModeExplorer

	blockTable := widgets.NewList()

	blockTable.TextStyle = ui.NewStyle(ui.ColorWhite)

	h0 := widgets.NewParagraph()
	h0.Title = "Current"

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
	b0.Text = "Use the arrow keys to scroll through the transactions. Press <Esc> to go back to the explorer view"

	b1 := widgets.NewList()
	b1.Title = "Block Info"
	b1.TextStyle = ui.NewStyle(ui.ColorYellow)
	b1.WrapText = false

	b2 := widgets.NewList()
	b2.Title = "Transactions"
	b2.TextStyle = ui.NewStyle(ui.ColorGreen)
	b2.WrapText = true

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
		ui.NewRow(2.5/5, blockTable),
	)

	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	blockGrid.SetRect(0, 0, termWidth, termHeight)

	var selectedBlock rpctypes.PolyBlock
	var setBlock = false
	var allBlocks metrics.SortableBlocks
	var recentBlocks metrics.SortableBlocks

	redraw := func(ms *monitorStatus) {
		if currentMode == monitorModeHelp {
			// TODO add some help context?
		} else if currentMode == monitorModeBlock {
			// render a block
			b1.Rows = metrics.GetSimpleBlockFields(selectedBlock)
			b2.Rows = metrics.GetSimpleBlockTxFields(selectedBlock, ms.ChainID)

			ui.Clear()
			ui.Render(blockGrid)
			return
		}

		if blockTable.SelectedRow == 0 {
			//default
			blocks := make([]rpctypes.PolyBlock, 0)

			ms.BlocksLock.RLock()
			for _, b := range ms.Blocks {
				blocks = append(blocks, b)
			}
			ms.BlocksLock.RUnlock()

			allBlocks = metrics.SortableBlocks(blocks)
			sort.Sort(allBlocks)
		}

		if uint64(len(allBlocks)) > *inputBatchSize {
			recentBlocks = allBlocks[uint64(len(allBlocks))-*inputBatchSize:]
		}

		h0.Text = fmt.Sprintf("Height: %s\nTime: %s", ms.HeadBlock.String(), time.Now().Format("02 Jan 06 15:04:05 MST"))
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
		rows, title := metrics.GetSimpleBlockRecords(recentBlocks)
		blockTable.Rows = rows
		blockTable.Title = title

		blockTable.TextStyle = ui.NewStyle(ui.ColorWhite)
		blockTable.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorRed, ui.ModifierBold)
		if blockTable.SelectedRow > 0 && blockTable.SelectedRow <= len(blockTable.Rows) {
			// only changed the selected block when the user presses the up down keys. Otherwise this will adjust when the table is updated automatically
			if setBlock {
				selectedBlock = recentBlocks[len(recentBlocks)-blockTable.SelectedRow]
				setBlock = false
			}

		}

		ui.Render(grid)
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
				blockTable.SelectedRow = 0
				currentMode = monitorModeExplorer
				redraw(ms)
			case "<Enter>":
				// TODO
				if blockTable.SelectedRow > 0 {
					currentMode = monitorModeBlock
				}
				redraw(ms)
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				blockGrid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				redraw(ms)
			case "<PageDown>", "<PageUp>":
				if currentMode == monitorModeBlock {
					if e.ID == "<PageDown>" {
						b2.ScrollPageDown()
					} else if e.ID == "<PageUp>" {
						b2.ScrollPageUp()
					}
					redraw(ms)
					break
				}
				if blockTable.SelectedRow == 0 {
					currIdx = 1
					blockTable.SelectedRow = currIdx
					setBlock = true
					redraw(ms)
					break
				}
				currIdx = blockTable.SelectedRow

				if e.ID == "<PageDown>" {
					currIdx = currIdx + 1
					setBlock = true
				} else if e.ID == "<PageUp>" {
					currIdx = currIdx - 1
					setBlock = true
				}
				if currIdx >= 0 && uint64(currIdx) <= *inputBatchSize { // need a better way to understand how many rows are visible
					blockTable.SelectedRow = currIdx
				}

				redraw(ms)
			case "<Up>", "<Down>", "<Left>", "<Right>":
				if currentMode == monitorModeBlock {
					if e.ID == "<Down>" {
						b2.ScrollDown()
					} else if e.ID == "<Up>" {
						b2.ScrollUp()
					}
					redraw(ms)
					break
				}
				if blockTable.SelectedRow == 0 {
					currIdx = 1
					blockTable.SelectedRow = currIdx
					setBlock = true
					redraw(ms)
					break
				}
				currIdx = blockTable.SelectedRow

				if e.ID == "<Down>" {
					if currIdx > int(*inputBatchSize)-1 {
						if int(*inputBatchSize)+10 < len(allBlocks) {
							*inputBatchSize = *inputBatchSize + 10
						} else {
							*inputBatchSize = uint64(len(allBlocks)) - 1
							break
						}
					}
					currIdx = currIdx + 1
					setBlock = true
				} else if e.ID == "<Up>" {
					currIdx = currIdx - 1
					setBlock = true
				}
				if currIdx >= 0 && uint64(currIdx) <= *inputBatchSize { // need a better way to understand how many rows are visble
					blockTable.SelectedRow = currIdx
				}

				redraw(ms)
			case "<MouseLeft>", "<MouseRight>", "<MouseRelease>", "<MouseWheelUp>", "<MouseWheelDown>":
				break
			default:
				log.Trace().Str("id", e.ID).Msg("Unknown ui event")
			}
		case <-ticker:
			if currentBn != ms.HeadBlock {
				currentBn = ms.HeadBlock
				redraw(ms)
			}
		}
	}
}

func setMonitorLogLevel(verbosity int64) {
	if verbosity < 100 {
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	} else if verbosity < 200 {
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	} else if verbosity < 300 {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if verbosity < 400 {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if verbosity < 500 {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if verbosity < 600 {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
}
