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

**⚠️ Important Requirements**:
- RPC endpoint with `eth_getBlockReceipts` support (see [RPC Requirements](#rpc-requirements))
- Chrome or Chromium for PDF generation (see [System Requirements](#system-requirements))

## Features

- **Stateless Operation**: All data is queried from the blockchain via RPC, no local storage required
- **Smart Defaults**: Automatically analyzes the latest 500 blocks if no range is specified
- **JSON Output**: Always generates a structured JSON report for programmatic analysis
- **HTML Visualization**: Optionally generates a visual HTML report with charts and tables
- **PDF Export**: Generate PDF reports (requires Chrome/Chromium installed)
- **Flexible Block Range**: Analyze any range of blocks with automatic range completion
- **Transaction Metrics**: Track transaction counts, gas usage, and other key metrics

## Basic Usage

Analyze the latest 500 blocks (no range specified):

```bash
polycli report --rpc-url http://localhost:8545
```

Generate a JSON report for blocks 1000 to 2000:

```bash
polycli report --rpc-url http://localhost:8545 --start-block 1000 --end-block 2000
```

Analyze 500 blocks starting from block 1000:

```bash
polycli report --rpc-url http://localhost:8545 --start-block 1000
```

Analyze the previous 500 blocks ending at block 2000:

```bash
polycli report --rpc-url http://localhost:8545 --end-block 2000
```

Analyze only the genesis block (block 0):

```bash
polycli report --rpc-url http://localhost:8545 --start-block 0 --end-block 0
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

Adjust concurrency for rate-limited endpoints:
```bash
polycli report --rpc-url https://public-rpc.example.com \
  --start-block 1000000 \
  --end-block 1001000 \
  --concurrency 5 \
  --rate-limit 2 \
  --format html
```

## Block Range Behavior

The command uses smart defaults for block ranges:

- **No flags specified**: Analyzes the latest 500 blocks on the chain
- **Only `--start-block` specified**: Analyzes 500 blocks starting from the specified block, or up to the latest block if fewer than 500 blocks remain
- **Only `--end-block` specified**: Analyzes 500 blocks ending at the specified block (or from block 0 if the chain has fewer than 500 blocks)
- **Both flags specified**: Analyzes the exact range specified (e.g., `--start-block 0 --end-block 0` analyzes only the genesis block)

The default range of 500 blocks can be modified by changing the `DefaultBlockRange` constant in the code.

**Note**: Block 0 (genesis) can be explicitly specified. To analyze only the genesis block, use `--start-block 0 --end-block 0`.

## RPC Requirements

**IMPORTANT**: This command requires an RPC endpoint that supports the `eth_getBlockReceipts` method. This is a non-standard but widely implemented extension to the JSON-RPC API.

### Supported RPC Providers
- ✅ Geth (full nodes)
- ✅ Erigon
- ✅ Polygon nodes
- ✅ Most self-hosted nodes
- ✅ Alchemy (premium endpoints)
- ✅ QuickNode
- ❌ Many public/free RPC endpoints (may not support `eth_getBlockReceipts`)
- ❌ Infura (does not support `eth_getBlockReceipts`)

If your RPC endpoint does not support `eth_getBlockReceipts`, the command will fail with an error like:
```
failed to fetch block receipts: method eth_getBlockReceipts does not exist/is not available
```

**Recommendation**: Use a self-hosted node or a premium RPC provider that supports this method.

## System Requirements

### PDF Generation

**IMPORTANT**: PDF report generation requires Google Chrome or Chromium to be installed on your system. The command uses Chrome's headless mode to render the HTML report as a PDF.

**Installing Chrome/Chromium:**

- **macOS**:
  ```bash
  brew install --cask google-chrome
  # or
  brew install chromium
  ```

- **Ubuntu/Debian**:
  ```bash
  sudo apt-get update
  sudo apt-get install chromium-browser
  # or
  sudo apt-get install google-chrome-stable
  ```

- **RHEL/CentOS/Fedora**:
  ```bash
  sudo dnf install chromium
  # or install Chrome from official RPM
  ```

- **Windows**: Download and install from [google.com/chrome](https://www.google.com/chrome/)

If Chrome/Chromium is not installed, PDF generation will fail with an error message indicating that Chrome could not be found.

**Alternative**: If you need PDF reports but cannot install Chrome, you can generate an HTML report and convert it to PDF using another tool.

## Notes

- To analyze a single block, set both start and end to the same block number (e.g., `--start-block 100 --end-block 100`)
- The command queries blocks concurrently with rate limiting to avoid overwhelming the RPC endpoint:
  - `--concurrency` controls the number of concurrent RPC requests (default: 10)
  - `--rate-limit` controls the maximum requests per second (default: 4)
  - Adjust these values based on your RPC endpoint's capacity
- Progress is logged every 100 blocks
- **Data Integrity**: The command automatically retries failed block fetches up to 3 times. If any blocks cannot be fetched after all retry attempts, the command fails with an error listing the failed blocks. This ensures reports are deterministic and complete - the same parameters always produce the same report.
- HTML reports include interactive hover tooltips on charts
- For large block ranges, consider running the command with a dedicated RPC endpoint

## Flags

```bash
      --concurrency int    number of concurrent RPC requests (default 10)
      --end-block uint     ending block number (default: auto-detect based on start-block or latest) (default 18446744073709551615)
  -f, --format string      output format [json, html, pdf] (default "json")
  -h, --help               help for report
  -o, --output string      output file path (default: stdout for JSON, report.html for HTML, report.pdf for PDF)
      --rate-limit float   requests per second limit (default 4)
      --rpc-url string     RPC endpoint URL (default "http://localhost:8545")
      --start-block uint   starting block number (default: auto-detect based on end-block or latest) (default 18446744073709551615)
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
