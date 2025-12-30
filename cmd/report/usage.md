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

Adjust concurrency for rate-limited endpoints:
```bash
polycli report --rpc-url https://public-rpc.example.com \
  --start-block 1000000 \
  --end-block 1001000 \
  --concurrency 5 \
  --rate-limit 2 \
  --format html
```

## Notes

- The `--end-block` flag is required; you must explicitly specify the block range to analyze
- The `--start-block` flag defaults to 0 (genesis block), which is a valid starting point
- To analyze a single block, set both start and end to the same block number (e.g., `--start-block 100 --end-block 100`)
- The command queries blocks concurrently with rate limiting to avoid overwhelming the RPC endpoint:
  - `--concurrency` controls the number of concurrent RPC requests (default: 10)
  - `--rate-limit` controls the maximum requests per second (default: 4)
  - Adjust these values based on your RPC endpoint's capacity
- Progress is logged every 100 blocks
- Blocks that cannot be fetched are skipped with a warning
- HTML reports include interactive hover tooltips on charts
- For large block ranges, consider running the command with a dedicated RPC endpoint
