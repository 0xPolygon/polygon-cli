The `report` command analyzes a range of blocks from an Ethereum-compatible blockchain and generates a comprehensive report with statistics and visualizations.

## Features

- **Stateless Operation**: All data is queried from the blockchain via RPC, no local storage required
- **Smart Defaults**: Automatically analyzes the latest 500 blocks if no range is specified
- **JSON Output**: Always generates a structured JSON report for programmatic analysis
- **HTML Visualization**: Optionally generates a visual HTML report with charts and tables
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

## Notes

- To analyze a single block, set both start and end to the same block number (e.g., `--start-block 100 --end-block 100`)
- The command queries blocks concurrently with rate limiting to avoid overwhelming the RPC endpoint:
  - `--concurrency` controls the number of concurrent RPC requests (default: 10)
  - `--rate-limit` controls the maximum requests per second (default: 4)
  - Adjust these values based on your RPC endpoint's capacity
- Progress is logged every 100 blocks
- Blocks that cannot be fetched are skipped with a warning
- HTML reports include interactive hover tooltips on charts
- For large block ranges, consider running the command with a dedicated RPC endpoint
