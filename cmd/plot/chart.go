package plot

import (
	"fmt"
	"image/color"
	"math"
	"slices"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/font/liberation"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func init() {
	// Use Liberation Sans (sans-serif) instead of the default Liberation Serif
	plot.DefaultFont = font.Font{Typeface: "Liberation", Variant: "Sans"}
	plotter.DefaultFont = font.Font{Typeface: "Liberation", Variant: "Sans"}

	// Register the Liberation font collection
	font.DefaultCache.Add(liberation.Collection())
}

var (
	lineThickness = vg.Points(2)

	gasBlockLimitLineColor = color.NRGBA{130, 38, 89, 220}
	gasTxsLimitLineColor   = color.NRGBA{255, 0, 189, 220}
	gasUsedLineColor       = color.NRGBA{0, 255, 133, 220}
	avgGasUsedLineColor    = color.NRGBA{255, 193, 7, 220}
	avgGasPriceLineColor   = color.NRGBA{30, 144, 255, 220}

	txDotsColor = color.NRGBA{0, 0, 0, 25}
	txDotsSizes = []vg.Length{
		vg.Points(3), // gasLimit <= 1M
		vg.Points(4), // gasLimit <= 2M
		vg.Points(5), // gasLimit <= 3M
		vg.Points(6), // gasLimit <= 4M
		vg.Points(7), // gasLimit <= 5M
		vg.Points(8), // gasLimit > 5M
	}

	targetTxDotsThickness = vg.Points(2)
	targetTxDotsSize      = vg.Points(8)
	targetTxDotsColor     = color.NRGBA{255, 0, 0, 255}
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

	title := fmt.Sprintf("ChainID: %d | Blocks %d - %d (%d) | Txs: %d | Scale: %s",
		metadata.chainID, metadata.startBlock, metadata.endBlock,
		metadata.endBlock-metadata.startBlock, metadata.blocksMetadata.txCount, scale)
	if len(metadata.targetAddr) > 0 {
		title += fmt.Sprintf("\nTarget: %s (%d txs)", metadata.targetAddr, metadata.blocksMetadata.targetTxCount)
	}

	p.Title.Text = title

	// Configure legend
	p.Legend.Top = true
	p.Legend.Left = true
	p.Legend.TextStyle.Font.Size = vg.Points(10)
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
				ticks[i].Label = formatSI(ticks[i].Value)
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
				ticks[i].Label = formatSI(ticks[i].Value)
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

	// other transactions (groups 1-6, index 0 is for target txs)
	for group := 1; group <= 6; group++ {
		points := txGroups[uint64(group)]
		if len(points) == 0 {
			continue
		}
		sc, err := plotter.NewScatter(points)
		if err != nil {
			log.Error().Err(err).Int("group", group).Msg("Failed to create scatter plot")
			continue
		}
		sc.GlyphStyle.Color = txDotsColor
		sc.GlyphStyle.Shape = draw.CircleGlyph{}
		sc.GlyphStyle.Radius = txDotsSizes[group-1]
		p.Add(sc)
	}

	// target transactions
	if len(txGroups[0]) > 0 {
		sc, err := plotter.NewScatter(txGroups[0])
		if err != nil {
			log.Error().Err(err).Msg("Failed to create target tx scatter plot")
			return
		}
		sc.GlyphStyle.Color = targetTxDotsColor
		sc.GlyphStyle.Shape = ThickCrossGlyph{Width: targetTxDotsThickness}
		sc.GlyphStyle.Radius = targetTxDotsSize
		p.Add(sc)
		p.Legend.Add("Target transactions", sc)
	}
}

// createLines creates line plots for various gas metrics.
func createLines(p *plot.Plot, metadata txGasChartMetadata) {
	numBlocks := len(metadata.blocksMetadata.blocks)
	blocks := make([]uint64, numBlocks)
	perBlockAvgGasPrice := make(map[uint64]float64, numBlocks)
	pointsBlockGasLimit := make(plotter.XYs, numBlocks)
	pointsTxsGasLimit := make(plotter.XYs, numBlocks)
	pointsAvgGasUsed := make(plotter.XYs, numBlocks)
	pointsGasUsed := make(plotter.XYs, numBlocks)

	for i, b := range metadata.blocksMetadata.blocks {
		blocks[i] = b.number

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

	addLine := func(points plotter.XYs, c color.Color, label string) {
		line, err := plotter.NewLine(points)
		if err != nil {
			log.Error().Err(err).Str("label", label).Msg("Failed to create line plot")
			return
		}
		line.Color = c
		line.Width = lineThickness
		p.Add(line)
		p.Legend.Add(label, line)
	}

	addLine(rollingMean(blocks, perBlockAvgGasPrice, 30), avgGasPriceLineColor, "30-block avg gas price")
	addLine(pointsGasUsed, gasUsedLineColor, "Block gas used")
	addLine(pointsTxsGasLimit, gasTxsLimitLineColor, "Txs gas limit")
	addLine(pointsBlockGasLimit, gasBlockLimitLineColor, "Block gas limit")
	addLine(pointsAvgGasUsed, avgGasUsedLineColor, "Avg block gas used")
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

// formatSI formats a number with intuitive suffixes (k, M, B, T).
func formatSI(v float64) string {
	if v < 1 {
		return fmt.Sprintf("%g", v)
	}

	var suffix string
	var scaled float64

	switch {
	case v >= 1e12:
		scaled = v / 1e12
		suffix = "T"
	case v >= 1e9:
		scaled = v / 1e9
		suffix = "B"
	case v >= 1e6:
		scaled = v / 1e6
		suffix = "M"
	case v >= 1e3:
		scaled = v / 1e3
		suffix = "k"
	default:
		return fmt.Sprintf("%.0f", v)
	}

	// Remove trailing zeros and decimal point if not needed
	formatted := fmt.Sprintf("%.2f", scaled)
	formatted = strings.TrimRight(formatted, "0")
	formatted = strings.TrimRight(formatted, ".")
	return formatted + suffix
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
