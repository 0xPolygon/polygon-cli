package txgaschart

import (
	"context"
	"fmt"
	"image/color"
	"math"
	"math/big"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/dustin/go-humanize"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"

	vgd "gonum.org/v1/plot/vg/draw"
)

func buildChart(cmd *cobra.Command) error {
	ctx := context.Background()
	log.Info().
		Str("rpc_url", inputArgs.rpcURL).
		Float64("rate_limit", inputArgs.rateLimit).
		Msg("RPC connection parameters")

	log.Info().
		Uint64("start_block", inputArgs.startBlock).
		Uint64("end_block", inputArgs.endBlock).
		Str("target_address", inputArgs.targetAddr).
		Msg("Chart generation parameters")

	target := common.HexToAddress(inputArgs.targetAddr)

	client, err := ethclient.DialContext(ctx, inputArgs.rpcURL)
	if err != nil {
		return err
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return err
	}

	rl := rate.NewLimiter(rate.Limit(inputArgs.rateLimit), 1)
	if inputArgs.rateLimit <= 0.0 {
		rl = nil
	}

	startBlock := inputArgs.startBlock
	endBlock := inputArgs.endBlock
	if endBlock == math.MaxUint64 {
		h, err := client.HeaderByNumber(ctx, nil)
		if err != nil {
			return err
		}
		endBlock = h.Number.Uint64()
		log.Warn().Uint64("end_block", endBlock).Msg("end block was set to max uint64 or not set, defaulting to latest block")
	}
	if startBlock > endBlock {
		return fmt.Errorf("start block %d cannot be greater than end block %d", startBlock, endBlock)
	}
	if startBlock == 0 {
		if endBlock < 1000 {
			startBlock = 0
		} else {
			startBlock = endBlock - 1000
		}

		log.Warn().Uint64("start_block", startBlock).
			Msg("start block was 0, defaulting to last 1000 blocks")
	}

	perBlockAvgMu := &sync.Mutex{}
	perBlockAvg := make(map[uint64]float64)

	wrk := make(chan struct{}, 15) // max concurrent requests
	for i := 0; i < cap(wrk); i++ {
		wrk <- struct{}{}
	}
	txs := make(chan TxPoint, 1000)
	wg := sync.WaitGroup{}
	wg.Add(1)
	var minGL, maxGL atomic.Uint64
	minGL.Store(math.MaxUint64)
	maxGL.Store(0.0)
	log.Info().Msg("reading blocks")
	go func() {
		wgb := sync.WaitGroup{}
		for b := startBlock; b <= endBlock; b++ {
			wgb.Add(1)
			go func(b uint64) {
				defer wgb.Done()
				<-wrk
				log.Trace().Uint64("block_number", b).Msg("processing block")
				defer func() { wrk <- struct{}{} }()
				_ = rl.Wait(ctx)
				blk, err := client.BlockByNumber(ctx, new(big.Int).SetUint64(b))
				if err != nil {
					log.Error().Err(err).Uint64("block_number", b).Msg("failed to fetch block")
					return
				}
				list := blk.Transactions()
				if len(list) == 0 {
					return
				}

				sum := 0.0
				for _, tx := range list {
					wg.Add(1)
					signer := types.LatestSignerForChainID(chainID)
					from, _ := types.Sender(signer, tx)
					isTarget := strings.EqualFold(from.String(), target.String())
					if !isTarget {
						isTarget = tx.To() != nil && strings.EqualFold(tx.To().String(), target.String())
					}
					gp := f64FromBigIntWei(tx.GasPrice())
					gl := tx.Gas()
					txs <- TxPoint{b, tx.Hash(), gp, gl, isTarget}

					minGL.Store(min(minGL.Load(), gl))
					maxGL.Store(max(maxGL.Load(), gl))
					sum += gp
				}
				perBlockAvgMu.Lock()
				perBlockAvg[b] = sum / float64(len(list))
				perBlockAvgMu.Unlock()
			}(b)
		}
		wgb.Wait()
		close(txs)
	}()

	// build plot
	p := plot.New()
	p.Title.Text = fmt.Sprintf("RPC: %s | Transactions from block %d to %d\nTransactions with red * are from/to: %s | Blue line is 30-block rolling avg gas price",
		inputArgs.rpcURL, startBlock, endBlock, target.String())
	p.Title.TextStyle.Font.Size = vg.Points(14)
	p.X.Label.Text = "Block Number"
	p.Y.Label.Text = "Gas Price (wei, log10)"
	p.Y.Scale = plot.LogScale{}
	p.Y.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
		ticks := plot.LogTicks{}.Ticks(min, max)
		for i := range ticks {
			if ticks[i].Label == "" {
				continue
			}
			v := ticks[i].Value // this is the real value (e.g., 10000), not an exponent
			if v >= 1 {
				ticks[i].Label = humanize.Comma(int64(math.Round(v)))
			} else {
				ticks[i].Label = fmt.Sprintf("%g", v) // 0.1, 0.01, etc.
			}
		}
		return ticks
	})

	// scatter points
	grayColor := color.RGBA{0, 0, 0, 25}
	redColor := color.RGBA{255, 0, 0, 255}

	log.Info().Msg("processing transactions for plotting")
	go func() {
		for t := range txs {
			sc, err := plotter.NewScatter(plotter.XYs{{X: float64(t.Block), Y: t.GasPrice}})
			if err != nil {
				log.Error().Err(err).
					Uint64("block", t.Block).
					Stringer("txHash", t.TxHash).
					Msg("failed to create scatter point")
				wg.Done()
				continue
			}
			if t.IsTarget {
				sc.GlyphStyle.Color = redColor
				sc.GlyphStyle.Radius = vg.Points(15)
				sc.GlyphStyle.Shape = ThickCrossGlyph{Width: vg.Points(4)} // thickness
				log.Info().
					Uint64("block", t.Block).
					Stringer("txHash", t.TxHash).
					Float64("gas_price_wei", t.GasPrice).
					Uint64("gas_limit", t.GasLimit).
					Msg("target tx found")
				p.Add(sc)
			} else {
				sc.GlyphStyle.Color = grayColor
				sc.GlyphStyle.Radius = scale(minGL.Load(), maxGL.Load(), t.GasLimit)
				sc.GlyphStyle.Shape = draw.CircleGlyph{}
			}
			p.Add(sc)
			wg.Done()
		}
		wg.Done()
	}()
	wg.Wait()

	// rolling avg line
	var blocks []uint64
	for b := range perBlockAvg {
		blocks = append(blocks, b)
	}
	lineXY := rollingMean(blocks, perBlockAvg, 30)
	line, _ := plotter.NewLine(lineXY)
	blueColor := color.RGBA{0, 100, 255, 255}
	line.Color = blueColor
	line.Width = vg.Points(3)
	p.Add(line)

	if err := p.Save(1600, 900, inputArgs.output); err != nil {
		return err
	}
	log.Info().
		Str("file", inputArgs.output).
		Msg("Chart saved successfully")

	return nil
}

// ThickCrossGlyph draws an 'X' with configurable stroke width.
type ThickCrossGlyph struct {
	Width vg.Length
}

func (g ThickCrossGlyph) DrawGlyph(c *vgd.Canvas, sty vgd.GlyphStyle, p vg.Point) {
	if !c.Contains(p) {
		return
	}
	r := sty.Radius
	ls := vgd.LineStyle{Color: sty.Color, Width: g.Width}

	// Horizontal
	h := []vg.Point{{X: p.X - r, Y: p.Y}, {X: p.X + r, Y: p.Y}}
	// Vertical
	v := []vg.Point{{X: p.X, Y: p.Y - r}, {X: p.X, Y: p.Y + r}}
	// Diagonal 1 (top-left -> bottom-right)
	d1 := []vg.Point{{X: p.X - r, Y: p.Y + r}, {X: p.X + r, Y: p.Y - r}}
	// Diagonal 2 (bottom-left -> top-right)
	d2 := []vg.Point{{X: p.X - r, Y: p.Y - r}, {X: p.X + r, Y: p.Y + r}}

	c.StrokeLines(ls, h)
	c.StrokeLines(ls, v)
	c.StrokeLines(ls, d1)
	c.StrokeLines(ls, d2)
}

type TxPoint struct {
	Block    uint64
	TxHash   common.Hash
	GasPrice float64
	GasLimit uint64
	IsTarget bool
}

func f64FromBigIntWei(x *big.Int) float64 {
	bf := new(big.Float).SetInt(x)
	f, _ := bf.Float64()
	return f
}

func rollingMean(blocks []uint64, perBlockAvg map[uint64]float64, window int) plotter.XYs {
	slices.Sort(blocks)
	points := make(plotter.XYs, len(blocks))
	sum := 0.0
	buffer := make([]float64, 0, window)
	for i, b := range blocks {
		val := perBlockAvg[b]
		buffer = append(buffer, val)
		sum += val
		if len(buffer) > window {
			sum -= buffer[0]
			buffer = buffer[1:]
		}
		points[i].X = float64(b)
		points[i].Y = sum / float64(len(buffer))
	}
	return points
}

// normalize sizes
func scale(minGL, maxGL, gl uint64) vg.Length {
	if maxGL == minGL {
		return vg.Points(4)
	}
	norm := (gl - minGL) / (maxGL - minGL)
	return vg.Points(float64(3 + 8*norm))
}
