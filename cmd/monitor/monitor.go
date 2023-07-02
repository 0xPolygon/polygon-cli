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
package monitor

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"sort"
	"sync"
	"time"

	_ "embed"

	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/cenkalti/backoff"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/maticnetwork/polygon-cli/metrics"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage       string
	batchSize   uint64
	windowSize  int
	intervalStr string
	interval    time.Duration

	one  = big.NewInt(1)
	zero = big.NewInt(0)
)

type (
	monitorStatus struct {
		ChainID   *big.Int
		HeadBlock *big.Int
		PeerCount uint64
		GasPrice  *big.Int

		Blocks            map[string]rpctypes.PolyBlock `json:"-"`
		BlocksLock        sync.RWMutex                  `json:"-"`
		MaxBlockRetrieved *big.Int
		MinBlockRetrieved *big.Int
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
	monitorModeHelp monitorMode = iota
	monitorModeExplorer
	monitorModeBlock
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
		log.Debug().Err(err).Msg("Using fake peer count")
		cs.PeerCount = 0
	}

	cs.GasPrice, err = ec.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't estimate gas: %s", err.Error())
	}

	return cs, nil

}

// monitorCmd represents the monitor command
var MonitorCmd = &cobra.Command{
	Use:   "monitor url",
	Short: "Monitor blocks using a JSON-RPC endpoint.",
	Long:  usage,
	Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// validate url argument
		_, err := url.Parse(args[0])
		if err != nil {
			return err
		}

		// validate batch-size flag
		if batchSize == 0 {
			return fmt.Errorf("batch-size can't be equal to zero")
		}

		// validate interval duration
		if interval, err = time.ParseDuration(intervalStr); err != nil {
			return err
		}

		return nil
	},
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

		isUiRendered := false
		errChan := make(chan error)
		go func() {
			for {
				var cs *chainState
				cs, err = getChainState(ctx, ec)
				if err != nil {
					log.Error().Err(err).Msg("Encountered issue fetching network information")
					time.Sleep(interval)
					continue
				}

				ms.HeadBlock = new(big.Int).SetUint64(cs.HeadBlock)
				ms.ChainID = cs.ChainID
				ms.PeerCount = cs.PeerCount
				ms.GasPrice = cs.GasPrice
				batchSize := new(big.Int).SetUint64(batchSize - 1)
				from := new(big.Int).Sub(ms.HeadBlock, batchSize)
				// Prevent getBlockRange from fetching duplicate blocks.
				if ms.MaxBlockRetrieved.Cmp(from) == 1 {
					from.Add(ms.MaxBlockRetrieved, big.NewInt(1))
				}

				if from.Cmp(zero) < 0 {
					from.SetInt64(0)
				}

				log.Debug().
					Int64("from", from.Int64()).
					Int64("to", ms.HeadBlock.Int64()).
					Int64("max", ms.MaxBlockRetrieved.Int64()).
					Msg("Fetching latest blocks")

				err = ms.getBlockRange(ctx, from, ms.HeadBlock, rpc)
				if err != nil {
					log.Error().Err(err).Msg("There was an issue fetching the block range")
				}

				to := new(big.Int).Sub(ms.MinBlockRetrieved, one)
				if from = new(big.Int).Sub(to, batchSize); from.Cmp(zero) < 0 {
					from.SetInt64(0)
				}

				log.Debug().
					Int64("from", from.Int64()).
					Int64("to", to.Int64()).
					Int64("min", ms.MinBlockRetrieved.Int64()).
					Msg("Fetching older blocks")

				err = ms.getBlockRange(ctx, from, to, rpc)
				if err != nil {
					log.Error().Err(err).Msg("There was an issue fetching the block range")
				}

				if !isUiRendered {
					go func() {
						errChan <- renderMonitorUI(ms)
					}()
					isUiRendered = true
				}

				time.Sleep(interval)
			}
		}()

		err = <-errChan
		return err
	},
}

func (ms *monitorStatus) getBlockRange(ctx context.Context, from, to *big.Int, c *ethrpc.Client) error {
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
		if ms.MinBlockRetrieved == nil || (ms.MinBlockRetrieved.Cmp(pb.Number()) == 1 && pb.Number().Cmp(zero) == 1) {
			ms.MinBlockRetrieved = pb.Number()
		}
	}

	return nil
}

func init() {
	MonitorCmd.PersistentFlags().Uint64VarP(&batchSize, "batch-size", "b", 25, "Number of requests per batch")
	MonitorCmd.PersistentFlags().StringVarP(&intervalStr, "interval", "i", "5s", "Amount of time between batch block rpc calls")
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
		ui.NewRow(5.0/10, blockTable),
	)

	termWidth, termHeight := ui.TerminalDimensions()
	windowSize = termHeight/2 - 4
	log.Info().Msg(fmt.Sprintln(windowSize))
	grid.SetRect(0, 0, termWidth, termHeight)
	blockGrid.SetRect(0, 0, termWidth, termHeight)

	var selectedBlock rpctypes.PolyBlock
	var setBlock = false
	var allBlocks metrics.SortableBlocks
	var renderedBlocks metrics.SortableBlocks
	windowOffset := 0

	redraw := func(ms *monitorStatus) {
		log.Debug().Interface("ms", ms).Msg("Redrawing")

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
			// default
			blocks := make([]rpctypes.PolyBlock, 0)

			ms.BlocksLock.RLock()
			for _, b := range ms.Blocks {
				blocks = append(blocks, b)
			}
			ms.BlocksLock.RUnlock()

			allBlocks = metrics.SortableBlocks(blocks)
			sort.Sort(allBlocks)
		}

		_, termHeight := ui.TerminalDimensions()
		// log.Info().Msg(fmt.Sprintln(termHeight/2 - 4))
		windowSize = termHeight/2 - 4
		start := len(allBlocks) - windowSize - windowOffset
		if start < 0 {
			start = 0
		}
		end := len(allBlocks) - windowOffset
		renderedBlocks = allBlocks[start:end]

		h0.Text = fmt.Sprintf("Height: %s\nTime: %s", ms.HeadBlock.String(), time.Now().Format("02 Jan 06 15:04:05 MST"))
		gasGwei := new(big.Int).Div(ms.GasPrice, metrics.UnitShannon)
		h1.Text = fmt.Sprintf("%s gwei", gasGwei.String())
		h2.Text = fmt.Sprintf("%d", ms.PeerCount)
		h3.Text = ms.ChainID.String()
		h4.Text = fmt.Sprintf("%0.2f", metrics.GetMeanBlockTime(renderedBlocks))

		sl0.Data = metrics.GetTxsPerBlock(renderedBlocks)
		sl1.Data = metrics.GetMeanGasPricePerBlock(renderedBlocks)
		sl2.Data = metrics.GetSizePerBlock(renderedBlocks)
		sl3.Data = metrics.GetUnclesPerBlock(renderedBlocks)
		sl4.Data = metrics.GetGasPerBlock(renderedBlocks)

		// If a row has not been selected, continue to update the list with new blocks.
		rows, title := metrics.GetSimpleBlockRecords(renderedBlocks)
		blockTable.Rows = rows
		blockTable.Title = title

		blockTable.TextStyle = ui.NewStyle(ui.ColorWhite)
		blockTable.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorRed, ui.ModifierBold)
		if blockTable.SelectedRow > 0 && blockTable.SelectedRow <= len(blockTable.Rows) {
			// Only changed the selected block when the user presses the up down keys.
			// Otherwise this will adjust when the table is updated automatically.
			if setBlock {
				selectedBlock = renderedBlocks[len(renderedBlocks)-blockTable.SelectedRow]
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
	previousKey := ""
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return nil
			case "<Escape>":
				blockTable.SelectedRow = 0
				currentMode = monitorModeExplorer
				windowOffset = 0
			case "<Enter>":
				if blockTable.SelectedRow > 0 {
					currentMode = monitorModeBlock
				}
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				blockGrid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
			case "<Up>", "<Down>":
				if currentMode == monitorModeBlock {
					if len(b2.Rows) != 0 && e.ID == "<Down>" {
						b2.ScrollDown()
					} else if len(b2.Rows) != 0 && e.ID == "<Up>" {
						b2.ScrollUp()
					}
					break
				}

				if blockTable.SelectedRow == 0 {
					currIdx = 1
					blockTable.SelectedRow = currIdx
					setBlock = true
					break
				}
				currIdx = blockTable.SelectedRow

				if e.ID == "<Down>" {
					log.Debug().
						Int("currIdx", currIdx).
						Int("windowSize", windowSize).
						Int("renderedBlocks", len(renderedBlocks)).
						Int("dy", blockTable.Dy()).
						Int("windowOffset", windowOffset).
						Msg("Down")

					if currIdx > windowSize-1 && windowOffset < len(allBlocks)-windowSize {
						windowOffset += 1
						break
					}
					currIdx += 1
					setBlock = true
				} else if e.ID == "<Up>" {
					log.Debug().Int("currIdx", currIdx).Int("windowSize", windowSize).Msg("Up")
					if currIdx <= 1 && windowOffset > 0 {
						windowOffset -= 1
						break
					}
					currIdx -= 1
					setBlock = true
				}
				// need a better way to understand how many rows are visible
				if currIdx > 0 && currIdx <= windowSize && currIdx <= len(renderedBlocks) {
					blockTable.SelectedRow = currIdx
				}
			case "<Home>":
				windowOffset = 0
				blockTable.SelectedRow = 1
				setBlock = true
			case "g":
				if previousKey == "g" {
					windowOffset = 0
					blockTable.SelectedRow = 1
					setBlock = true
				}
			case "G", "<End>":
				if len(renderedBlocks) < windowSize {
					windowOffset = 0
					blockTable.SelectedRow = len(renderedBlocks)
				} else {
					windowOffset = len(allBlocks) - windowSize
					blockTable.SelectedRow = max(windowSize, len(renderedBlocks))
				}
				setBlock = true
			case "<C-f>", "<PageDown>":
				if len(renderedBlocks) < windowSize {
					windowOffset = 0
					blockTable.SelectedRow = len(renderedBlocks)
					break
				}

				windowOffset += windowSize
				if windowOffset > len(allBlocks)-windowSize {
					windowOffset = len(allBlocks) - windowSize
					blockTable.SelectedRow = windowSize
				}
			case "<C-b>", "<PageUp>":
				windowOffset -= windowSize
				if windowOffset < 0 {
					windowOffset = 0
					blockTable.SelectedRow = 1
				}
			default:
				log.Trace().Str("id", e.ID).Msg("Unknown ui event")
			}

			if previousKey == "g" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}

			redraw(ms)
		case <-ticker:
			if currentBn != ms.HeadBlock {
				currentBn = ms.HeadBlock
				redraw(ms)
			}
		}
	}
}

func max(nums ...int) int {
	m := nums[0]
	for _, n := range nums {
		if m < n {
			m = n
		}
	}
	return m
}
