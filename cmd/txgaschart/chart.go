package txgaschart

import (
	"fmt"
	"image/color"
	"math"
	"slices"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

var (
	lineThickness = 2

	gasBlockLimitLineColor     = color.NRGBA{130, 38, 89, 220}
	gasBlockLimitLineColorName = "Purple"
	gasBlockLimitLineThickness = lineThickness

	gasTxsLimitLineColor     = color.NRGBA{255, 0, 189, 220}
	gasTxsLimitLineColorName = "Pink"
	gasTxsLimitLineThickness = lineThickness

	gasUsedLineColor     = color.NRGBA{0, 255, 133, 220}
	gasUsedLineColorName = "Green"
	gasUsedLineThickness = lineThickness

	avgGasUsedLineColor     = color.NRGBA{255, 193, 7, 220}
	avgGasUsedLineColorName = "Orange"
	avgGasUsedLineThickness = lineThickness

	avgGasPriceAvgLineColor     = color.NRGBA{30, 144, 255, 220}
	avgGasPriceAvgLineColorName = "Blue"
	avgGasPriceAvgLineThickness = lineThickness

	txDotsColor = color.NRGBA{0, 0, 0, 25}
	txDotsSize1 = 3
	txDotsSize2 = 4
	txDotsSize3 = 5
	txDotsSize4 = 6
	txDotsSize5 = 7
	txDotsSize6 = 8

	targetTxDotsThickness = 2
	targetTxDotsSize      = 8
	targetTxDotsColor     = color.NRGBA{255, 0, 0, 255}
	targetTxDotsColorName = "Red"
)

// txGasChartMetadata holds metadata for generating the transaction gas chart.
type txGasChartMetadata struct {
	rpcURL  string
	chainID uint64

	targetAddr string
	startBlock uint64
	endBlock   uint64

	blocksMetadata blocksMetadata

	scale string

	outputPath string
}

// plotChart generates and saves the transaction gas chart based on the provided metadata.
func plotChart(metadata txGasChartMetadata) error {
	p := plot.New()
	createHeader(p, metadata)
	createTxsDots(p, metadata)
	createLines(p, metadata)

	p.X.Min = float64(metadata.startBlock)
	p.X.Max = float64(metadata.endBlock) + (float64(metadata.endBlock-metadata.startBlock) * 0.02)

	// Protect min and max for logarithmic scale
	if metadata.blocksMetadata.minTxGasPrice == 0 {
		metadata.blocksMetadata.minTxGasPrice = 1
	}
	if metadata.blocksMetadata.maxTxGasPrice == 0 {
		metadata.blocksMetadata.maxTxGasPrice = 1
	}

	p.Y.Min = float64(metadata.blocksMetadata.minTxGasPrice)
	p.Y.Max = float64(metadata.blocksMetadata.maxTxGasPrice) * 1.02

	return save(p, metadata)
}

// createHeader sets the title and header information for the plot.
func createHeader(p *plot.Plot, metadata txGasChartMetadata) {
	p.Title.TextStyle.Font.Size = vg.Points(14)
	scale := "logarithmic"
	if !strings.EqualFold(metadata.scale, "log") {
		scale = "linear"
	}

	title := fmt.Sprintf("ChainID: %d | Block range %d - %d | Scale: %s\n", metadata.chainID, metadata.startBlock, metadata.endBlock, scale)
	title += fmt.Sprintf("Blocks: %d | Txs: %d", metadata.endBlock-metadata.startBlock, metadata.blocksMetadata.txCount)
	if len(metadata.targetAddr) > 0 {
		title += fmt.Sprintf(" | Target Txs: %d | Target Addr: %s\n", metadata.blocksMetadata.targetTxCount, metadata.targetAddr)
	} else {
		title += "\n"
	}

	title += fmt.Sprintf("%s stars are target transactions | %s line is 30-block rolling avg gas price\n", targetTxDotsColorName, avgGasPriceAvgLineColorName)
	title += fmt.Sprintf("Gas in %%, Y height = 100%% | %s line is block gas used | %s line is avg block gas used | %s line is block gas limit | %s line is txs gas limit",
		gasUsedLineColorName, avgGasUsedLineColorName, gasBlockLimitLineColorName, gasTxsLimitLineColorName)

	p.Title.Text = title
}

// createTxsDots creates scatter plots for transaction gas prices.
func createTxsDots(p *plot.Plot, metadata txGasChartMetadata) {
	p.X.Label.Text = "Block Number"

	// Custom ticker for X axis (block numbers must be integers)
	p.X.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
		ticks := plot.DefaultTicks{}.Ticks(min, max)
		for i := range ticks {
			if ticks[i].Label == "" {
				continue
			}
			// Format as integer without decimal places
			ticks[i].Label = humanize.Comma(int64(math.Round(ticks[i].Value)))
		}
		return ticks
	})

	if strings.EqualFold(metadata.scale, "log") {
		p.Y.Scale = plot.LogScale{}
		p.Y.Label.Text = "Gas Price (wei, log)"

		// Custom ticker for logarithmic scale
		p.Y.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
			// Protect against values <= 0
			if min <= 0 {
				min = 1
			}
			if max <= 0 {
				max = 1
			}

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
	} else {
		p.Y.Scale = plot.LinearScale{}
		p.Y.Label.Text = "Gas Price (wei, linear)"

		// Custom ticker for linear scale
		p.Y.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
			ticks := plot.DefaultTicks{}.Ticks(min, max)
			for i := range ticks {
				if ticks[i].Label == "" {
					continue
				}
				v := ticks[i].Value
				if v >= 1 {
					ticks[i].Label = humanize.Comma(int64(math.Round(v)))
				} else {
					ticks[i].Label = fmt.Sprintf("%g", v) // 0.1, 0.01, etc.
				}
			}
			return ticks
		})
	}

	txGroups := make(map[uint64]plotter.XYs)
	txGroups[0] = make(plotter.XYs, 0) // target txs
	txGroups[1] = make(plotter.XYs, 0) // gasLimit <= 1,000,000
	txGroups[2] = make(plotter.XYs, 0) // gasLimit <= 2,000,000
	txGroups[3] = make(plotter.XYs, 0) // gasLimit <= 3,000,000
	txGroups[4] = make(plotter.XYs, 0) // gasLimit <= 4,000,000
	txGroups[5] = make(plotter.XYs, 0) // gasLimit <= 5,000,000
	txGroups[6] = make(plotter.XYs, 0) // gasLimit > 5,000,000

	for _, b := range metadata.blocksMetadata.blocks {
		for _, t := range b.txs {
			// For plotting on a logarithmic Y scale we must avoid zero/negative values.
			// Clamp gasPrice to at least 1 for visualization purposes; this means the
			// plotted gas price may differ from the original t.gasPrice when it is <= 0.
			gasPrice := t.gasPrice
			if gasPrice <= 0 {
				gasPrice = 1
			}

			// Use the local gasPrice variable (protected) in all appends
			if t.target {
				txGroups[0] = append(txGroups[0], plotter.XY{X: float64(b.number), Y: float64(gasPrice)})
			} else if t.gasLimit <= 1000000 {
				txGroups[1] = append(txGroups[1], plotter.XY{X: float64(b.number), Y: float64(gasPrice)})
			} else if t.gasLimit <= 2000000 {
				txGroups[2] = append(txGroups[2], plotter.XY{X: float64(b.number), Y: float64(gasPrice)})
			} else if t.gasLimit <= 3000000 {
				txGroups[3] = append(txGroups[3], plotter.XY{X: float64(b.number), Y: float64(gasPrice)})
			} else if t.gasLimit <= 4000000 {
				txGroups[4] = append(txGroups[4], plotter.XY{X: float64(b.number), Y: float64(gasPrice)})
			} else if t.gasLimit <= 5000000 {
				txGroups[5] = append(txGroups[5], plotter.XY{X: float64(b.number), Y: float64(gasPrice)})
			} else {
				txGroups[6] = append(txGroups[6], plotter.XY{X: float64(b.number), Y: float64(gasPrice)})
			}
		}
	}

	// other transactions
	for group := range len(txGroups) {
		if group == 0 {
			continue
		}
		if len(txGroups[uint64(group)]) == 0 {
			continue
		}
		sc, err := plotter.NewScatter(txGroups[uint64(group)])
		if err != nil {
			log.Error().Err(err).Msg("Failed to create regular tx scatter plot")
		} else {
			sc.GlyphStyle.Color = txDotsColor
			sc.GlyphStyle.Shape = draw.CircleGlyph{}
			switch group {
			case 1:
				sc.GlyphStyle.Radius = vg.Points(float64(txDotsSize1))
			case 2:
				sc.GlyphStyle.Radius = vg.Points(float64(txDotsSize2))
			case 3:
				sc.GlyphStyle.Radius = vg.Points(float64(txDotsSize3))
			case 4:
				sc.GlyphStyle.Radius = vg.Points(float64(txDotsSize4))
			case 5:
				sc.GlyphStyle.Radius = vg.Points(float64(txDotsSize5))
			case 6:
				sc.GlyphStyle.Radius = vg.Points(float64(txDotsSize6))
			}
			p.Add(sc)
		}
	}

	// target transactions
	if len(txGroups[0]) > 0 {
		sc, err := plotter.NewScatter(txGroups[0])
		if err != nil {
			log.Error().Err(err).Msg("Failed to create target tx scatter plot")
		} else {
			sc.GlyphStyle.Color = targetTxDotsColor
			sc.GlyphStyle.Shape = ThickCrossGlyph{Width: vg.Points(float64(targetTxDotsThickness))}
			sc.GlyphStyle.Radius = vg.Points(float64(targetTxDotsSize))
			p.Add(sc)
		}
	}
}

// createLines creates line plots for various gas metrics.
func createLines(p *plot.Plot, metadata txGasChartMetadata) {
	var blocks []uint64
	var perBlockAvgGasPrice = make(map[uint64]float64)
	pointsBlockGasLimit := make(plotter.XYs, len(metadata.blocksMetadata.blocks))
	pointsTxsGasLimit := make(plotter.XYs, len(metadata.blocksMetadata.blocks))
	pointsAvgGasUsed := make(plotter.XYs, len(metadata.blocksMetadata.blocks))
	pointsGasUsed := make(plotter.XYs, len(metadata.blocksMetadata.blocks))
	for i, b := range metadata.blocksMetadata.blocks {
		blocks = append(blocks, b.number)

		// Protect avgGasPrice for logarithmic scale
		avgGasPrice := float64(b.avgGasPrice)
		if avgGasPrice <= 0 {
			avgGasPrice = 1
		}
		perBlockAvgGasPrice[b.number] = avgGasPrice

		pointsBlockGasLimit[i].X = float64(b.number)
		pointsBlockGasLimit[i].Y = scaleGasToGasPrice(b.gasLimit, metadata)

		pointsTxsGasLimit[i].X = float64(b.number)
		pointsTxsGasLimit[i].Y = scaleGasToGasPrice(b.txsGasLimit, metadata)

		pointsAvgGasUsed[i].X = float64(b.number)
		pointsAvgGasUsed[i].Y = scaleGasToGasPrice(metadata.blocksMetadata.avgBlockGasUsed, metadata)

		pointsGasUsed[i].X = float64(b.number)
		pointsGasUsed[i].Y = scaleGasToGasPrice(b.gasUsed, metadata)
	}

	lineXY := rollingMean(blocks, perBlockAvgGasPrice, 30)
	line, _ := plotter.NewLine(lineXY)
	line.Color = avgGasPriceAvgLineColor
	line.Width = vg.Points(float64(avgGasPriceAvgLineThickness))
	p.Add(line)

	line, _ = plotter.NewLine(pointsGasUsed)
	line.Color = gasUsedLineColor
	line.Width = vg.Points(float64(gasUsedLineThickness))
	p.Add(line)

	line, _ = plotter.NewLine(pointsTxsGasLimit)
	line.Color = gasTxsLimitLineColor
	line.Width = vg.Points(float64(gasTxsLimitLineThickness))
	p.Add(line)

	line, _ = plotter.NewLine(pointsBlockGasLimit)
	line.Color = gasBlockLimitLineColor
	line.Width = vg.Points(float64(gasBlockLimitLineThickness))
	p.Add(line)

	line, _ = plotter.NewLine(pointsAvgGasUsed)
	line.Color = avgGasUsedLineColor
	line.Width = vg.Points(float64(avgGasUsedLineThickness))
	p.Add(line)
}

// save saves the plot to the specified output path.
func save(p *plot.Plot, metadata txGasChartMetadata) error {
	if err := p.Save(1600, 900, metadata.outputPath); err != nil {
		return err
	}
	log.Info().
		Str("file", metadata.outputPath).
		Msg("Chart saved successfully")
	return nil
}

// rollingMean calculates the rolling mean of per-block average gas prices over a specified window.
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

// scaleGasToGasPrice scales the gas limit to a corresponding gas price based on the provided metadata.
func scaleGasToGasPrice(gasLimit uint64, metadata txGasChartMetadata) float64 {
	minTxGasPrice := metadata.blocksMetadata.minTxGasPrice
	maxTxGasPrice := metadata.blocksMetadata.maxTxGasPrice

	maxBlockGasLimit := metadata.blocksMetadata.maxBlockGasLimit

	if maxBlockGasLimit == 0 {
		return 1
	}

	yRange := float64(maxTxGasPrice) - float64(minTxGasPrice)
	proportion := float64(gasLimit) / float64(maxBlockGasLimit)
	y := proportion*yRange + float64(minTxGasPrice)

	if y < 1 {
		y = 1
	}
	return y
}

// ThickCrossGlyph draws an 'X' with configurable stroke width.
type ThickCrossGlyph struct {
	Width vg.Length
}

// DrawGlyph implements the GlyphDrawer interface.
func (g ThickCrossGlyph) DrawGlyph(c *draw.Canvas, sty draw.GlyphStyle, p vg.Point) {
	if !c.Contains(p) {
		return
	}
	r := sty.Radius
	ls := draw.LineStyle{Color: sty.Color, Width: g.Width}

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
