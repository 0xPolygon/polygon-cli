package report

import (
	_ "embed"
	"fmt"
	"strings"
	"time"
)

//go:embed template.html
var htmlTemplate string

// generateHTML creates an HTML report from the BlockReport data
func generateHTML(report *BlockReport) string {
	html := htmlTemplate

	// Replace metadata placeholders
	html = strings.ReplaceAll(html, "{{CHAIN_ID}}", fmt.Sprintf("%d", report.ChainID))
	html = strings.ReplaceAll(html, "{{RPC_URL}}", report.RpcUrl)
	html = strings.ReplaceAll(html, "{{BLOCK_RANGE}}", fmt.Sprintf("%d - %d", report.StartBlock, report.EndBlock))
	html = strings.ReplaceAll(html, "{{GENERATED_AT}}", report.GeneratedAt.Format(time.RFC3339))
	html = strings.ReplaceAll(html, "{{TOTAL_BLOCKS}}", formatNumber(report.Summary.TotalBlocks))

	// Generate and replace stat cards
	html = strings.ReplaceAll(html, "{{STAT_CARDS}}", generateStatCards(report))

	// Generate and replace charts
	html = strings.ReplaceAll(html, "{{TX_COUNT_CHART}}", generateTxCountChart(report))
	html = strings.ReplaceAll(html, "{{GAS_USAGE_CHART}}", generateGasUsageChart(report))

	// Generate and replace top 10 sections
	html = strings.ReplaceAll(html, "{{TOP_10_SECTIONS}}", generateTop10Sections(report))

	return html
}

// generateStatCards creates the statistics cards HTML
func generateStatCards(report *BlockReport) string {
	var sb strings.Builder

	cards := []struct {
		title string
		value string
	}{
		{"Total Blocks", formatNumber(report.Summary.TotalBlocks)},
		{"Total Transactions", formatNumber(report.Summary.TotalTransactions)},
		{"Unique Senders", formatNumber(report.Summary.UniqueSenders)},
		{"Unique Recipients", formatNumber(report.Summary.UniqueRecipients)},
		{"Average Tx/Block", fmt.Sprintf("%.2f", report.Summary.AvgTxPerBlock)},
		{"Total Gas Used", formatNumber(report.Summary.TotalGasUsed)},
		{"Average Gas/Block", formatNumber(uint64(report.Summary.AvgGasPerBlock))},
	}

	// Add base fee card if available
	if report.Summary.AvgBaseFeePerGas > 0 {
		cards = append(cards, struct {
			title string
			value string
		}{"Avg Base Fee (Gwei)", fmt.Sprintf("%.2f", float64(report.Summary.AvgBaseFeePerGas)/1e9)})
	}

	for _, card := range cards {
		sb.WriteString(fmt.Sprintf(`
            <div class="stat-card">
                <h3>%s</h3>
                <div class="value">%s</div>
            </div>`, card.title, card.value))
	}

	return sb.String()
}

// generateTxCountChart creates a line chart for transaction counts
func generateTxCountChart(report *BlockReport) string {
	if len(report.Blocks) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`
        <div class="section">
            <h2>Transaction Count by Block</h2>
            <div class="chart-container">`)

	// Find max tx count for scaling
	maxTx := uint64(1)
	for _, block := range report.Blocks {
		if block.TxCount > maxTx {
			maxTx = block.TxCount
		}
	}

	// Limit the number of points to avoid overcrowding
	step := 1
	if len(report.Blocks) > 500 {
		step = len(report.Blocks) / 500
	}

	// Generate SVG line chart
	width := 1200.0
	height := 300.0
	padding := 40.0
	chartWidth := width - 2*padding
	chartHeight := height - 2*padding

	sb.WriteString(fmt.Sprintf(`
            <svg width="100%%" height="%.0f" viewBox="0 0 %.0f %.0f" class="line-chart">
                <!-- Grid lines -->
                <line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ddd" stroke-width="1"/>
                <line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ddd" stroke-width="1"/>`,
		height, width, height,
		padding, padding, padding, height-padding,
		padding, height-padding, width-padding, height-padding))

	// Build points for the line
	var points []string
	var circles strings.Builder
	numPoints := 0
	for i := 0; i < len(report.Blocks); i += step {
		block := report.Blocks[i]
		x := padding + (float64(numPoints) / float64((len(report.Blocks)-1)/step)) * chartWidth
		y := height - padding - (float64(block.TxCount) / float64(maxTx)) * chartHeight

		points = append(points, fmt.Sprintf("%.2f,%.2f", x, y))
		circles.WriteString(fmt.Sprintf(`
                <circle cx="%.2f" cy="%.2f" r="3" fill="#3498db" class="chart-point">
                    <title>Block %d: %d transactions</title>
                </circle>`, x, y, block.Number, block.TxCount))
		numPoints++
	}

	// Draw the line
	sb.WriteString(fmt.Sprintf(`
                <polyline points="%s" fill="none" stroke="#3498db" stroke-width="2"/>`,
		strings.Join(points, " ")))

	// Draw circles (points)
	sb.WriteString(circles.String())

	// Add axis labels
	sb.WriteString(fmt.Sprintf(`
                <text x="%.0f" y="%.0f" text-anchor="middle" font-size="12" fill="#666">Block Number</text>
                <text x="20" y="%.0f" text-anchor="middle" font-size="12" fill="#666" transform="rotate(-90 20 %.0f)">Transactions</text>
                <text x="%.0f" y="%.0f" text-anchor="start" font-size="10" fill="#666">%d</text>
                <text x="%.0f" y="%.0f" text-anchor="end" font-size="10" fill="#666">%d</text>
                <text x="%.0f" y="%.0f" text-anchor="start" font-size="10" fill="#666">0</text>
                <text x="%.0f" y="%.0f" text-anchor="start" font-size="10" fill="#666">%d</text>
            </svg>
        </div>`,
		width/2, height-10,
		height/2, height/2,
		padding, height-padding+15, report.Blocks[0].Number,
		width-padding, height-padding+15, report.Blocks[len(report.Blocks)-1].Number,
		padding-35, height-padding+5,
		padding-35, padding+5, maxTx))

	sb.WriteString(`
        </div>`)

	return sb.String()
}

// generateGasUsageChart creates a line chart for gas usage
func generateGasUsageChart(report *BlockReport) string {
	if len(report.Blocks) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`
        <div class="section">
            <h2>Gas Usage by Block</h2>
            <div class="chart-container">`)

	// Find max gas for scaling
	maxGas := uint64(1)
	for _, block := range report.Blocks {
		if block.GasUsed > maxGas {
			maxGas = block.GasUsed
		}
	}

	// Limit the number of points to avoid overcrowding
	step := 1
	if len(report.Blocks) > 500 {
		step = len(report.Blocks) / 500
	}

	// Generate SVG line chart
	width := 1200.0
	height := 300.0
	padding := 40.0
	chartWidth := width - 2*padding
	chartHeight := height - 2*padding

	sb.WriteString(fmt.Sprintf(`
            <svg width="100%%" height="%.0f" viewBox="0 0 %.0f %.0f" class="line-chart">
                <!-- Grid lines -->
                <line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ddd" stroke-width="1"/>
                <line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ddd" stroke-width="1"/>`,
		height, width, height,
		padding, padding, padding, height-padding,
		padding, height-padding, width-padding, height-padding))

	// Build points for the line
	var points []string
	var circles strings.Builder
	numPoints := 0
	for i := 0; i < len(report.Blocks); i += step {
		block := report.Blocks[i]
		x := padding + (float64(numPoints) / float64((len(report.Blocks)-1)/step)) * chartWidth
		y := height - padding - (float64(block.GasUsed) / float64(maxGas)) * chartHeight

		points = append(points, fmt.Sprintf("%.2f,%.2f", x, y))
		circles.WriteString(fmt.Sprintf(`
                <circle cx="%.2f" cy="%.2f" r="3" fill="#9b59b6" class="chart-point">
                    <title>Block %d: %s gas</title>
                </circle>`, x, y, block.Number, formatNumber(block.GasUsed)))
		numPoints++
	}

	// Draw the line
	sb.WriteString(fmt.Sprintf(`
                <polyline points="%s" fill="none" stroke="#9b59b6" stroke-width="2"/>`,
		strings.Join(points, " ")))

	// Draw circles (points)
	sb.WriteString(circles.String())

	// Add axis labels
	sb.WriteString(fmt.Sprintf(`
                <text x="%.0f" y="%.0f" text-anchor="middle" font-size="12" fill="#666">Block Number</text>
                <text x="20" y="%.0f" text-anchor="middle" font-size="12" fill="#666" transform="rotate(-90 20 %.0f)">Gas Used</text>
                <text x="%.0f" y="%.0f" text-anchor="start" font-size="10" fill="#666">%d</text>
                <text x="%.0f" y="%.0f" text-anchor="end" font-size="10" fill="#666">%d</text>
                <text x="%.0f" y="%.0f" text-anchor="start" font-size="10" fill="#666">0</text>
                <text x="%.0f" y="%.0f" text-anchor="start" font-size="10" fill="#666">%s</text>
            </svg>
        </div>`,
		width/2, height-10,
		height/2, height/2,
		padding, height-padding+15, report.Blocks[0].Number,
		width-padding, height-padding+15, report.Blocks[len(report.Blocks)-1].Number,
		padding-35, height-padding+5,
		padding-35, padding+5, formatNumber(maxGas)))

	sb.WriteString(`
        </div>`)

	return sb.String()
}

// generateBlocksTable creates a table with detailed block information
func generateBlocksTable(report *BlockReport) string {
	if len(report.Blocks) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`
        <h2>Block Details</h2>
        <table>
            <thead>
                <tr>
                    <th>Block Number</th>
                    <th>Timestamp</th>
                    <th>Transactions</th>
                    <th>Gas Used</th>
                    <th>Gas Limit</th>
                    <th>Gas Used %</th>`)

	// Check if any block has base fee
	hasBaseFee := false
	for _, block := range report.Blocks {
		if block.BaseFeePerGas != nil {
			hasBaseFee = true
			break
		}
	}

	if hasBaseFee {
		sb.WriteString(`
                    <th>Base Fee (Gwei)</th>`)
	}

	sb.WriteString(`
                </tr>
            </thead>
            <tbody>`)

	// Limit table rows if there are too many blocks
	blocks := report.Blocks
	showEllipsis := false
	if len(blocks) > 1000 {
		// Show first 500 and last 500
		blocks = append(report.Blocks[:500], report.Blocks[len(report.Blocks)-500:]...)
		showEllipsis = true
	}

	for i, block := range blocks {
		// Insert ellipsis row after first 500
		if showEllipsis && i == 500 {
			colSpan := 6
			if hasBaseFee {
				colSpan = 7
			}
			sb.WriteString(fmt.Sprintf(`
                <tr>
                    <td colspan="%d" style="text-align: center; font-style: italic;">
                        ... (showing first 500 and last 500 blocks of %d total)
                    </td>
                </tr>`, colSpan, len(report.Blocks)))
		}

		timestamp := time.Unix(int64(block.Timestamp), 0).Format("2006-01-02 15:04:05")
		gasUsedPercent := 0.0
		if block.GasLimit > 0 {
			gasUsedPercent = (float64(block.GasUsed) / float64(block.GasLimit)) * 100
		}

		sb.WriteString(fmt.Sprintf(`
                <tr>
                    <td>%d</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%.2f%%</td>`,
			block.Number,
			timestamp,
			formatNumber(block.TxCount),
			formatNumber(block.GasUsed),
			formatNumber(block.GasLimit),
			gasUsedPercent))

		if hasBaseFee {
			baseFeeGwei := "-"
			if block.BaseFeePerGas != nil {
				baseFeeGwei = fmt.Sprintf("%.2f", float64(block.BaseFeePerGas.Uint64())/1e9)
			}
			sb.WriteString(fmt.Sprintf(`
                    <td>%s</td>`, baseFeeGwei))
		}

		sb.WriteString(`
                </tr>`)
	}

	sb.WriteString(`
            </tbody>
        </table>`)

	return sb.String()
}

// formatNumber adds thousand separators to numbers
func formatNumber(n uint64) string {
	str := fmt.Sprintf("%d", n)
	if len(str) <= 3 {
		return str
	}

	var result strings.Builder
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(digit)
	}
	return result.String()
}

// formatNumberWithUnits formats large numbers with units (K, M, B, T, Q)
func formatNumberWithUnits(n uint64) string {
	if n == 0 {
		return "0"
	}

	units := []struct {
		suffix    string
		threshold uint64
	}{
		{"Q", 1e15},  // Quadrillion
		{"T", 1e12},  // Trillion
		{"B", 1e9},   // Billion
		{"M", 1e6},   // Million
		{"K", 1e3},   // Thousand
	}

	for _, unit := range units {
		if n >= unit.threshold {
			value := float64(n) / float64(unit.threshold)
			// Format with appropriate precision
			if value >= 100 {
				return fmt.Sprintf("%.0f%s", value, unit.suffix)
			} else if value >= 10 {
				return fmt.Sprintf("%.1f%s", value, unit.suffix)
			} else {
				return fmt.Sprintf("%.2f%s", value, unit.suffix)
			}
		}
	}

	return formatNumber(n)
}

// generateTop10Sections creates all top 10 sections HTML
func generateTop10Sections(report *BlockReport) string {
	var sb strings.Builder

	sb.WriteString(`
        <div class="section">
            <h2>Top 10 Analysis</h2>`)

	// Top 10 blocks by transaction count
	if len(report.Top10.BlocksByTxCount) > 0 {
		sb.WriteString(`
        <div class="subsection">
            <h3>Top 10 Blocks by Transaction Count</h3>
            <table>
            <thead>
                <tr>
                    <th>Rank</th>
                    <th>Block Number</th>
                    <th>Transaction Count</th>
                </tr>
            </thead>
            <tbody>`)

		for i, block := range report.Top10.BlocksByTxCount {
			sb.WriteString(fmt.Sprintf(`
                <tr>
                    <td>%d</td>
                    <td>%s</td>
                    <td>%s</td>
                </tr>`, i+1, formatNumber(block.Number), formatNumber(block.TxCount)))
		}

		sb.WriteString(`
            </tbody>
        </table>
        </div>`)
	}

	// Top 10 blocks by gas used
	if len(report.Top10.BlocksByGasUsed) > 0 {
		sb.WriteString(`
        <div class="subsection">
            <h3>Top 10 Blocks by Gas Used</h3>
            <table>
            <thead>
                <tr>
                    <th>Rank</th>
                    <th>Block Number</th>
                    <th>Gas Used (Wei)</th>
                    <th>Gas Limit</th>
                    <th>Gas Used %</th>
                </tr>
            </thead>
            <tbody>`)

		for i, block := range report.Top10.BlocksByGasUsed {
			sb.WriteString(fmt.Sprintf(`
                <tr>
                    <td>%d</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%.2f%%</td>
                </tr>`, i+1, formatNumber(block.Number), formatNumber(block.GasUsed), formatNumberWithUnits(block.GasLimit), block.GasUsedPercent))
		}

		sb.WriteString(`
            </tbody>
        </table>
        </div>`)
	}

	// Top 10 transactions by gas used
	if len(report.Top10.TransactionsByGas) > 0 {
		sb.WriteString(`
        <div class="subsection">
            <h3>Top 10 Transactions by Gas Used</h3>
            <table>
            <thead>
                <tr>
                    <th>Rank</th>
                    <th>Transaction Hash</th>
                    <th>Block Number</th>
                    <th>Gas Limit</th>
                    <th>Gas Used (Wei)</th>
                </tr>
            </thead>
            <tbody>`)

		for i, tx := range report.Top10.TransactionsByGas {
			sb.WriteString(fmt.Sprintf(`
                <tr>
                    <td>%d</td>
                    <td><code>%s</code></td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                </tr>`, i+1, tx.Hash, formatNumber(tx.BlockNumber), formatNumberWithUnits(tx.GasLimit), formatNumber(tx.GasUsed)))
		}

		sb.WriteString(`
            </tbody>
        </table>
        </div>`)
	}

	// Top 10 transactions by gas limit
	if len(report.Top10.TransactionsByGasLimit) > 0 {
		sb.WriteString(`
        <div class="subsection">
            <h3>Top 10 Transactions by Gas Limit</h3>
            <table>
            <thead>
                <tr>
                    <th>Rank</th>
                    <th>Transaction Hash</th>
                    <th>Block Number</th>
                    <th>Gas Limit</th>
                    <th>Gas Used (Wei)</th>
                </tr>
            </thead>
            <tbody>`)

		for i, tx := range report.Top10.TransactionsByGasLimit {
			sb.WriteString(fmt.Sprintf(`
                <tr>
                    <td>%d</td>
                    <td><code>%s</code></td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                </tr>`, i+1, tx.Hash, formatNumber(tx.BlockNumber), formatNumberWithUnits(tx.GasLimit), formatNumber(tx.GasUsed)))
		}

		sb.WriteString(`
            </tbody>
        </table>
        </div>`)
	}

	// Top 10 most used gas prices
	if len(report.Top10.MostUsedGasPrices) > 0 {
		sb.WriteString(`
        <div class="subsection">
            <h3>Top 10 Most Used Gas Prices</h3>
            <table>
            <thead>
                <tr>
                    <th>Rank</th>
                    <th>Gas Price (Wei)</th>
                    <th>Transaction Count</th>
                </tr>
            </thead>
            <tbody>`)

		for i, gp := range report.Top10.MostUsedGasPrices {
			sb.WriteString(fmt.Sprintf(`
                <tr>
                    <td>%d</td>
                    <td>%s</td>
                    <td>%s</td>
                </tr>`, i+1, formatNumber(gp.GasPrice), formatNumber(gp.Count)))
		}

		sb.WriteString(`
            </tbody>
        </table>
        </div>`)
	}

	// Top 10 most used gas limits
	if len(report.Top10.MostUsedGasLimits) > 0 {
		sb.WriteString(`
        <div class="subsection">
            <h3>Top 10 Most Used Gas Limits</h3>
            <table>
            <thead>
                <tr>
                    <th>Rank</th>
                    <th>Gas Limit</th>
                    <th>Transaction Count</th>
                </tr>
            </thead>
            <tbody>`)

		for i, gl := range report.Top10.MostUsedGasLimits {
			sb.WriteString(fmt.Sprintf(`
                <tr>
                    <td>%d</td>
                    <td>%s</td>
                    <td>%s</td>
                </tr>`, i+1, formatNumberWithUnits(gl.GasLimit), formatNumber(gl.Count)))
		}

		sb.WriteString(`
            </tbody>
        </table>
        </div>`)
	}

	sb.WriteString(`
        </div>`)

	return sb.String()
}
