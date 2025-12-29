# `polycli report`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Generate a report analyzing a range of blocks from an Ethereum-compatible blockchain.

```bash
polycli report [flags]
```

## Usage

The `report` command analyzes a range of blocks from an Ethereum-compatible blockchain and generates a comprehensive report with statistics and visualizations.

## Features

- **Stateless Operation**: All data is queried from the blockchain via RPC, no local storage required
- **JSON Output**: Always generates a structured JSON report for programmatic analysis
- **HTML Visualization**: Optionally generates a visual HTML report with charts and tables
- **Block Range Analysis**: Analyze any range of blocks from start to end
- **Transaction Metrics**: Track transaction counts, gas usage, and other key metrics

## Basic Usage

Generate a JSON report for blocks 1000 to 2000:

```bash
polycli report --rpc-url http://localhost:8545 --start-block 1000 --end-block 2000
```

Generate an HTML report:

```bash
polycli report --rpc-url http://localhost:8545 --start-block 1000 --end-block 2000 --format html
```

Save JSON output to a file:

```bash
polycli report --rpc-url http://localhost:8545 --start-block 1000 --end-block 2000 --output report.json
```

## Report Contents

The report includes:

### Summary Statistics
- Total number of blocks analyzed
- Total transaction count
- Average transactions per block
- Total gas used across all blocks
- Average gas used per block
- Average base fee per gas (for EIP-1559 compatible chains)

### Block Details
For each block in the range:
- Block number and timestamp
- Transaction count
- Gas used and gas limit
- Base fee per gas (if available)

### Visualizations (HTML format only)
- Transaction count chart showing distribution across blocks
- Gas usage chart showing gas consumption patterns
- Detailed table with all block information

## Examples

Analyze recent blocks:
```bash
polycli report --rpc-url http://localhost:8545 --start-block 19000000 --end-block 19000100 --format html -o analysis.html
```

Generate JSON for automated processing:
```bash
polycli report --rpc-url https://eth-mainnet.g.alchemy.com/v2/YOUR_API_KEY \
  --start-block 18000000 \
  --end-block 18001000 \
  --output mainnet-analysis.json
```

Quick analysis to stdout:
```bash
polycli report --rpc-url http://localhost:8545 --start-block 1000 --end-block 1100 | jq '.summary'
```

## Notes

- The command queries blocks sequentially to avoid overwhelming the RPC endpoint
- Progress is logged every 100 blocks
- Blocks that cannot be fetched are skipped with a warning
- HTML reports include interactive hover tooltips on charts
- For large block ranges, consider running the command with a dedicated RPC endpoint

## Flags

```bash
      --concurrency int    number of concurrent RPC requests (default 10)
      --end-block uint     ending block number for analysis
  -f, --format string      output format [json, html, pdf] (default "json")
  -h, --help               help for report
  -o, --output string      output file path (default: stdout for JSON, report.html for HTML, report.pdf for PDF)
      --rate-limit float   requests per second limit (default 4)
      --rpc-url string     RPC endpoint URL (default "http://localhost:8545")
      --start-block uint   starting block number for analysis
```

The command also inherits flags from parent commands.

```bash
      --config string      config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs        output logs in pretty format instead of JSON (default true)
  -v, --verbosity string   log level (string or int):
                             0   - silent
                             100 - panic
                             200 - fatal
                             300 - error
                             400 - warn
                             500 - info (default)
                             600 - debug
                             700 - trace (default "info")
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
