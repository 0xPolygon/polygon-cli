package plot

import (
	"bytes"
	_ "embed"
	"fmt"
	"html"
	"html/template"
	"os"
	"regexp"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/rs/zerolog/log"
)

//go:embed template.html
var chartTemplate string

const (
	// Gas limit thresholds for grouping transactions (in gas units)
	gasLimit1M = 1_000_000
	gasLimit2M = 2_000_000
	gasLimit3M = 3_000_000
	gasLimit4M = 4_000_000
	gasLimit5M = 5_000_000

	// Number of gas limit groups (excluding target txs group)
	numGasLimitGroups = 6
)

var (
	// Colors as hex strings for go-echarts
	gasBlockLimitLineColor = "#822659"
	gasTxsLimitLineColor   = "#FF00BD"
	gasUsedLineColor       = "#00FF85"
	avgGasUsedLineColor    = "#FFC107"

	txDotsColor   = "#4682B4" // Steel blue
	targetTxColor = "#FF0000" // Red

	// Symbol sizes for gas limit groups (in pixels)
	txDotsSizes = []int{6, 8, 10, 12, 14, 16}

	// Gas limit group names for legend
	gasLimitGroupNames = []string{
		"Gas ≤1M",
		"Gas ≤2M",
		"Gas ≤3M",
		"Gas ≤4M",
		"Gas ≤5M",
		"Gas >5M",
	}
)

// txGasChartMetadata holds metadata for generating the transaction gas chart.
type txGasChartMetadata struct {
	chainID uint64

	targetAddr string
	startBlock uint64
	endBlock   uint64

	blocksMetadata blocksMetadata

	renderer string

	outputPath string
}

// formatSI formats a number with SI suffixes (k, M, B, T).
func formatSI(v float64) string {
	if v < 1 {
		return fmt.Sprintf("%.2f", v)
	}
	switch {
	case v >= 1e12:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", v/1e12), "0"), ".") + "T"
	case v >= 1e9:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", v/1e9), "0"), ".") + "B"
	case v >= 1e6:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", v/1e6), "0"), ".") + "M"
	case v >= 1e3:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", v/1e3), "0"), ".") + "k"
	default:
		return fmt.Sprintf("%.0f", v)
	}
}

// formatGasPrice formats a gas price in wei to the most convenient unit (wei, gwei, or ether).
func formatGasPrice(wei uint64) string {
	const (
		gwei  = 1e9
		ether = 1e18
	)
	v := float64(wei)
	switch {
	case v >= ether:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.6f", v/ether), "0"), ".") + " ether"
	case v >= gwei:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", v/gwei), "0"), ".") + " gwei"
	default:
		return fmt.Sprintf("%d wei", wei)
	}
}

// formatTxTooltip formats a tooltip for a transaction scatter point.
func formatTxTooltip(seriesName string, blockNum uint64, blockHash string, txHash string, gasLimit, gasPrice uint64) string {
	return fmt.Sprintf("<b>%s</b><br/>Block: %d<br/>Block Hash: %s<br/>Gas Limit: %s gas<br/>Gas Price: %s<br/>Transaction: %s",
		html.EscapeString(seriesName),
		blockNum,
		html.EscapeString(blockHash),
		formatSI(float64(gasLimit)),
		formatGasPrice(gasPrice),
		html.EscapeString(txHash))
}

// formatLineTooltip formats a tooltip for a line chart point.
func formatLineTooltip(seriesName string, blockNum uint64, blockHash string, txCount int, value float64) string {
	return fmt.Sprintf("<b>%s</b><br/>Block: %d<br/>Block Hash: %s<br/>Transaction Count: %d<br/>Value: %s gas",
		html.EscapeString(seriesName),
		blockNum,
		html.EscapeString(blockHash),
		txCount,
		formatSI(value))
}

// plotChart generates and saves the transaction gas chart based on the provided metadata.
func plotChart(metadata txGasChartMetadata) error {
	scatter := createScatterChart(metadata)
	line := createLineChart(metadata)

	// Overlay line chart on scatter chart
	scatter.Overlap(line)

	return save(scatter, metadata)
}

// createScatterChart creates scatter series for transaction gas prices.
func createScatterChart(metadata txGasChartMetadata) *charts.Scatter {
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:    "1600px",
			Height:   "900px",
			Renderer: metadata.renderer,
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    createTitle(metadata),
			Left:     "center",
			Top:      "2%",
			Subtitle: createSubtitle(metadata),
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:   opts.Bool(true),
			Top:    "8%",
			Left:   "center",
			Orient: "horizontal",
		}),
		charts.WithGridOpts(opts.Grid{
			Top:    "15%", // Leave room for title and legend
			Bottom: "10%", // Leave room for x-axis label and slider
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show: opts.Bool(true),
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  opts.Bool(true),
					Title: "Save",
				},
				Restore: &opts.ToolBoxFeatureRestore{
					Show:  opts.Bool(true),
					Title: "Restore",
				},
				DataZoom: &opts.ToolBoxFeatureDataZoom{
					Show:  opts.Bool(true),
					Title: map[string]string{"zoom": "Zoom", "back": "Reset"},
				},
				DataView: &opts.ToolBoxFeatureDataView{
					Show:  opts.Bool(true),
					Title: "Data",
				},
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      opts.Bool(true),
			Trigger:   "item",
			Enterable: opts.Bool(true),
			Formatter: opts.FuncOpts(`function(params) { return params.name || params.seriesName; }`),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name:         "Block Number",
			NameLocation: "center",
			NameGap:      30,
			Type:         "value",
			Min:          metadata.startBlock,
			Max:          metadata.endBlock,
			AxisLabel: &opts.AxisLabel{
				Formatter: "{value}",
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:     "Gas",
			Type:     "log",
			Min:      1, // For log scale, min must be at least 1
			Max:      metadata.blocksMetadata.maxBlockGasLimit,
			Position: "left",
		}),
		// X-axis zoom (inside + slider)
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
		// Y-axis zoom (slider only, no mousewheel)
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      0,
			End:        100,
			YAxisIndex: []int{0},
			Orient:     "vertical",
		}),
	)

	// Group transactions by gas limit
	txGroups, targetTxs := groupTransactionsByGasLimit(metadata)

	// Add series for each gas limit group
	for i, name := range gasLimitGroupNames {
		if len(txGroups[i]) == 0 {
			continue
		}
		scatter.AddSeries(name, txGroups[i],
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color:   txDotsColor,
				Opacity: opts.Float(0.6),
			}),
			charts.WithScatterChartOpts(opts.ScatterChart{
				SymbolSize: txDotsSizes[i],
			}),
		)
	}

	// Add target transactions series
	if len(targetTxs) > 0 {
		scatter.AddSeries("Target Txs", targetTxs,
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color:   targetTxColor,
				Opacity: opts.Float(1),
			}),
			charts.WithScatterChartOpts(opts.ScatterChart{
				Symbol:     "diamond",
				SymbolSize: 16,
			}),
		)
	}

	return scatter
}

// createLineChart creates line series for various gas metrics.
func createLineChart(metadata txGasChartMetadata) *charts.Line {
	line := charts.NewLine()

	numBlocks := len(metadata.blocksMetadata.blocks)

	pointsBlockGasLimit := make([]opts.LineData, numBlocks)
	pointsTxsGasLimit := make([]opts.LineData, numBlocks)
	pointsAvgGasUsed := make([]opts.LineData, numBlocks)
	pointsGasUsed := make([]opts.LineData, numBlocks)

	for i, b := range metadata.blocksMetadata.blocks {
		blockHash := b.Hash.Hex()
		txCount := len(b.Txs)

		pointsBlockGasLimit[i] = opts.LineData{
			Value: []any{b.Number, b.GasLimit},
			Name:  formatLineTooltip("Block Gas Limit", b.Number, blockHash, txCount, float64(b.GasLimit)),
		}
		pointsTxsGasLimit[i] = opts.LineData{
			Value: []any{b.Number, b.TxsGasLimit},
			Name:  formatLineTooltip("Transaction Gas Limit", b.Number, blockHash, txCount, float64(b.TxsGasLimit)),
		}
		pointsAvgGasUsed[i] = opts.LineData{
			Value: []any{b.Number, metadata.blocksMetadata.avgBlockGasUsed},
			Name:  formatLineTooltip("Avg Block Gas Used", b.Number, blockHash, txCount, float64(metadata.blocksMetadata.avgBlockGasUsed)),
		}
		pointsGasUsed[i] = opts.LineData{
			Value: []any{b.Number, b.GasUsed},
			Name:  formatLineTooltip("Block Gas Used", b.Number, blockHash, txCount, float64(b.GasUsed)),
		}
	}

	// Line chart options (all use the same Y-axis)
	lineChartOpts := charts.WithLineChartOpts(opts.LineChart{
		Symbol:     "circle",
		SymbolSize: 4,
		ShowSymbol: opts.Bool(true),
		YAxisIndex: 0,
	})

	line.AddSeries("Block Gas Used", pointsGasUsed,
		charts.WithLineStyleOpts(opts.LineStyle{Color: gasUsedLineColor, Width: 2}),
		charts.WithItemStyleOpts(opts.ItemStyle{Color: gasUsedLineColor}),
		lineChartOpts,
	).
		AddSeries("Transaction Gas Limit", pointsTxsGasLimit,
			charts.WithLineStyleOpts(opts.LineStyle{Color: gasTxsLimitLineColor, Width: 2}),
			charts.WithItemStyleOpts(opts.ItemStyle{Color: gasTxsLimitLineColor}),
			lineChartOpts,
		).
		AddSeries("Block Gas Limit", pointsBlockGasLimit,
			charts.WithLineStyleOpts(opts.LineStyle{Color: gasBlockLimitLineColor, Width: 2}),
			charts.WithItemStyleOpts(opts.ItemStyle{Color: gasBlockLimitLineColor}),
			lineChartOpts,
		).
		AddSeries("Avg Block Gas Used", pointsAvgGasUsed,
			charts.WithLineStyleOpts(opts.LineStyle{Color: avgGasUsedLineColor, Width: 2}),
			charts.WithItemStyleOpts(opts.ItemStyle{Color: avgGasUsedLineColor}),
			lineChartOpts,
		)

	return line
}

// createTitle creates the chart title.
func createTitle(metadata txGasChartMetadata) string {
	return fmt.Sprintf("ChainID: %d | Blocks %d - %d (%d) | Transactions: %d",
		metadata.chainID, metadata.startBlock, metadata.endBlock,
		metadata.endBlock-metadata.startBlock, metadata.blocksMetadata.txCount)
}

// createSubtitle creates the chart subtitle (target address info if specified).
func createSubtitle(metadata txGasChartMetadata) string {
	if len(metadata.targetAddr) > 0 {
		return fmt.Sprintf("Target: %s (%d transactions)", metadata.targetAddr, metadata.blocksMetadata.targetTxCount)
	}
	return ""
}

// groupTransactionsByGasLimit groups transactions by their gas limit into scatter data.
// Returns a slice of groups (indexed 0-5 for ≤1M through >5M) and a separate slice for target txs.
func groupTransactionsByGasLimit(metadata txGasChartMetadata) ([][]opts.ScatterData, []opts.ScatterData) {
	txGroups := make([][]opts.ScatterData, numGasLimitGroups)
	for i := range txGroups {
		txGroups[i] = make([]opts.ScatterData, 0)
	}
	targetTxs := make([]opts.ScatterData, 0)

	for _, b := range metadata.blocksMetadata.blocks {
		blockHash := b.Hash.Hex()
		for _, t := range b.Txs {
			// Clamp gasLimit to at least 1 for logarithmic Y scale
			gasLimit := t.GasLimit
			if gasLimit <= 0 {
				gasLimit = 1
			}

			groupIdx := gasLimitToGroupIndex(t.GasLimit)
			seriesName := gasLimitGroupNames[groupIdx]
			if t.Target {
				seriesName = "Target Txs"
			}

			// Value: [blockNumber, gasLimit] with pre-formatted tooltip in Name
			point := opts.ScatterData{
				Value: []any{b.Number, gasLimit},
				Name:  formatTxTooltip(seriesName, b.Number, blockHash, t.Hash.Hex(), t.GasLimit, t.GasPrice),
			}

			if t.Target {
				targetTxs = append(targetTxs, point)
				continue
			}

			txGroups[groupIdx] = append(txGroups[groupIdx], point)
		}
	}

	return txGroups, targetTxs
}

// gasLimitToGroupIndex returns the group index (0-5) for a given gas limit.
func gasLimitToGroupIndex(gasLimit uint64) int {
	switch {
	case gasLimit <= gasLimit1M:
		return 0
	case gasLimit <= gasLimit2M:
		return 1
	case gasLimit <= gasLimit3M:
		return 2
	case gasLimit <= gasLimit4M:
		return 3
	case gasLimit <= gasLimit5M:
		return 4
	default:
		return 5
	}
}

// templateData holds the data for rendering the chart template.
type templateData struct {
	Width        string
	Height       string
	Renderer     string
	ChartOptions template.JS
}

// save saves the chart to the specified output path using the custom template.
func save(chart *charts.Scatter, metadata txGasChartMetadata) error {
	// Render chart to buffer to extract options JSON
	var buf bytes.Buffer
	if err := chart.Render(&buf); err != nil {
		return fmt.Errorf("failed to render chart: %w", err)
	}

	// Extract the chart options JSON from the rendered HTML
	chartOptions, err := extractChartOptions(buf.String())
	if err != nil {
		return fmt.Errorf("failed to extract chart options: %w", err)
	}

	// Prepare template data
	data := templateData{
		Width:        "1600px",
		Height:       "900px",
		Renderer:     metadata.renderer,
		ChartOptions: template.JS(chartOptions),
	}

	// Parse and execute template
	tmpl, err := template.New("chart").Parse(chartTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	f, err := os.Create(metadata.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	log.Info().
		Str("file", metadata.outputPath).
		Msg("Chart saved successfully")
	return nil
}

// extractChartOptions extracts the echarts options JSON from rendered go-echarts HTML.
func extractChartOptions(htmlContent string) (string, error) {
	// go-echarts renders: let option_XXXXX = {...};
	// Find the start of the options assignment
	re := regexp.MustCompile(`let option_\w+ = `)
	loc := re.FindStringIndex(htmlContent)
	if loc == nil {
		return "", fmt.Errorf("could not find chart options in rendered HTML")
	}

	// Start after "let option_XXX = "
	start := loc[1]

	// Find the matching closing brace by counting braces
	braceCount := 0
	end := start
	for i := start; i < len(htmlContent); i++ {
		switch htmlContent[i] {
		case '{':
			braceCount++
		case '}':
			braceCount--
			if braceCount == 0 {
				end = i + 1
				return htmlContent[start:end], nil
			}
		}
	}

	return "", fmt.Errorf("could not find matching closing brace for chart options")
}
