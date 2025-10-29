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
	gasLimitLineColor     = color.NRGBA{255, 0, 0, 100}
	gasLimitLineThickness = 2

	gasUsedLineColor     = color.NRGBA{0, 255, 0, 150}
	gasUsedLineThickness = 1

	avgGasUsedLineColor     = color.NRGBA{255, 165, 0, 150}
	avgGasUsedLineThickness = 2

	avgGasPriceAvgLineColor     = color.NRGBA{0, 0, 255, 150}
	avgGasPriceAvgLineThickness = 3

	txDotsColor       = color.NRGBA{0, 0, 0, 25}
	targetTxDotsColor = color.NRGBA{255, 0, 0, 255}
)

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

func plotChart(metadata txGasChartMetadata) error {
	p := plot.New()
	createHeader(p, metadata)
	createLines(p, metadata)
	createTxsDots(p, metadata)
	return save(p, metadata)
}

func createHeader(p *plot.Plot, metadata txGasChartMetadata) {
	p.Title.TextStyle.Font.Size = vg.Points(14)
	title := fmt.Sprintf("ChainID: %d | Block range %d - %d\n", metadata.chainID, metadata.startBlock, metadata.endBlock)
	title += fmt.Sprintf("Blocks: %d | Txs: %d | Target Txs: %d", metadata.endBlock-metadata.startBlock, metadata.blocksMetadata.txCount, metadata.blocksMetadata.targetTxCount)
	if len(metadata.targetAddr) > 0 {
		title += fmt.Sprintf(" | Target Addr: %s\n", metadata.targetAddr)
	} else {
		title += "\n"
	}
	title += "Red stars are target transactions | Blue line is 30-block rolling avg gas price\nGas are % Y height means 100% | Green line is block gas used | Orange line is avg block gas used | Red line is block gas limit"
	p.Title.Text = title
}

func createTxsDots(p *plot.Plot, metadata txGasChartMetadata) {
	p.X.Label.Text = "Block Number"

	if strings.EqualFold(metadata.scale, "log") {
		p.Y.Scale = plot.LogScale{}
		p.Y.Label.Text = "Gas Price (wei, log)"
	} else {
		p.Y.Scale = plot.LinearScale{}
		p.Y.Label.Text = "Gas Price (wei, linear)"
	}

	p.X.Min = float64(metadata.startBlock)
	p.X.Max = float64(metadata.endBlock) + (float64(metadata.endBlock-metadata.startBlock) * 0.02)
	p.Y.Min = float64(metadata.blocksMetadata.minTxGasPrice)
	p.Y.Max = float64(metadata.blocksMetadata.maxTxGasPrice) * 1.02
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
			if t.gasPrice <= 0 {
				t.gasPrice = 1
			}
			if t.target {
				txGroups[0] = append(txGroups[0], plotter.XY{X: float64(b.number), Y: float64(t.gasPrice)})
			} else if t.gasLimit <= 1000000 {
				txGroups[1] = append(txGroups[1], plotter.XY{X: float64(b.number), Y: float64(t.gasPrice)})
			} else if t.gasLimit <= 2000000 {
				txGroups[2] = append(txGroups[2], plotter.XY{X: float64(b.number), Y: float64(t.gasPrice)})
			} else if t.gasLimit <= 3000000 {
				txGroups[3] = append(txGroups[3], plotter.XY{X: float64(b.number), Y: float64(t.gasPrice)})
			} else if t.gasLimit <= 4000000 {
				txGroups[4] = append(txGroups[4], plotter.XY{X: float64(b.number), Y: float64(t.gasPrice)})
			} else if t.gasLimit <= 5000000 {
				txGroups[5] = append(txGroups[5], plotter.XY{X: float64(b.number), Y: float64(t.gasPrice)})
			} else {
				txGroups[6] = append(txGroups[6], plotter.XY{X: float64(b.number), Y: float64(t.gasPrice)})
			}
		}
	}

	if len(txGroups[0]) > 0 {
		sc, err := plotter.NewScatter(txGroups[0])
		if err != nil {
			log.Error().Err(err).Msg("failed to create target tx scatter plot")
		} else {
			sc.GlyphStyle.Color = targetTxDotsColor
			sc.GlyphStyle.Shape = ThickCrossGlyph{Width: vg.Points(4)}
			sc.GlyphStyle.Radius = vg.Points(15)
			p.Add(sc)
		}
	}

	for group := range len(txGroups) {
		if group == 0 {
			continue
		}
		if len(txGroups[uint64(group)]) == 0 {
			continue
		}
		sc, err := plotter.NewScatter(txGroups[uint64(group)])
		if err != nil {
			log.Error().Err(err).Msg("failed to create regular tx scatter plot")
		} else {
			sc.GlyphStyle.Color = txDotsColor
			sc.GlyphStyle.Shape = draw.CircleGlyph{}
			switch group {
			case 1:
				sc.GlyphStyle.Radius = vg.Points(3)
			case 2:
				sc.GlyphStyle.Radius = vg.Points(4)
			case 3:
				sc.GlyphStyle.Radius = vg.Points(5)
			case 4:
				sc.GlyphStyle.Radius = vg.Points(6)
			case 5:
				sc.GlyphStyle.Radius = vg.Points(7)
			case 6:
				sc.GlyphStyle.Radius = vg.Points(8)
			}
			p.Add(sc)
		}
	}
}

func createLines(p *plot.Plot, metadata txGasChartMetadata) {
	var blocks []uint64
	var perBlockAvgGasPrice = make(map[uint64]float64)
	pointsGasLimit := make(plotter.XYs, len(metadata.blocksMetadata.blocks))
	pointsAvgGasUsed := make(plotter.XYs, len(metadata.blocksMetadata.blocks))
	pointsGasUsed := make(plotter.XYs, len(metadata.blocksMetadata.blocks))
	for i, b := range metadata.blocksMetadata.blocks {
		blocks = append(blocks, b.number)

		perBlockAvgGasPrice[b.number] = float64(b.avgGasPrice)

		pointsGasLimit[i].X = float64(b.number)
		pointsGasLimit[i].Y = scaleGasToGasPrice(b.gasLimit, metadata)

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

	line, _ = plotter.NewLine(pointsGasLimit)
	line.Color = gasLimitLineColor
	line.Width = vg.Points(float64(gasLimitLineThickness))
	p.Add(line)

	line, _ = plotter.NewLine(pointsAvgGasUsed)
	line.Color = avgGasUsedLineColor
	line.Width = vg.Points(float64(avgGasUsedLineThickness))
	p.Add(line)
}

func save(p *plot.Plot, metadata txGasChartMetadata) error {
	if err := p.Save(1600, 900, metadata.outputPath); err != nil {
		return err
	}
	log.Info().
		Str("file", metadata.outputPath).
		Msg("Chart saved successfully")
	return nil
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
