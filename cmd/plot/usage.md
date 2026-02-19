# plot

`plot` generates interactive HTML charts analyzing transaction gas prices and limits across a range of blocks. It fetches block data from an Ethereum-compatible RPC endpoint and produces an HTML chart showing gas price distribution, transaction gas limits, block gas limits, and block gas usage over time.

## Basic Usage

Generate a chart for the last 500 blocks:

```bash
polycli plot --rpc-url http://localhost:8545
```

This will create a file named `plot.html` in the current directory.

## Analyzing Specific Block Ranges

Analyze blocks 9356826 to 9358826:

```bash
polycli plot --rpc-url https://sepolia.infura.io/v3/YOUR_API_KEY \
  --start-block 9356826 \
  --end-block 9358826 \
  --output "sepolia_analysis.html"
```

## Highlighting Target Address Transactions

Track transactions involving a specific address (either sent from or to):

```bash
polycli plot --rpc-url http://localhost:8545 \
  --target-address "0xeE76bECaF80fFe451c8B8AFEec0c21518Def02f9" \
  --start-block 1000 \
  --end-block 2000
```

Target transactions will be highlighted in the chart and logged during execution.

## Performance Options

When fetching large block ranges, adjust rate limiting and concurrency.

Process 10 blocks concurrently with 10 requests/second rate limit:

```bash
polycli plot --rpc-url http://localhost:8545 \
  --concurrency 10 \
  --rate-limit 10 \
  --start-block 1000 \
  --end-block 5000
```

Remove rate limiting entirely (use with caution):

```bash
polycli plot --rpc-url http://localhost:8545 \
  --rate-limit -1 \
  --concurrency 20
```

## Renderer Options

Choose between SVG (default) and Canvas rendering. Both output HTML files but use different rendering technologies.

Use SVG renderer (default, sharper at any zoom level):

```bash
polycli plot --rpc-url http://localhost:8545 \
  --renderer "svg" \
  --output "svg_chart.html"
```

Use Canvas renderer (better performance for large datasets):

```bash
polycli plot --rpc-url http://localhost:8545 \
  --renderer "canvas" \
  --output "canvas_chart.html"
```

## Understanding the Chart

The generated chart displays key metrics and is interactive:

1. **Transaction Gas Prices**: Individual transaction gas prices plotted as points, with target address transactions highlighted
2. **Transaction Gas Limits**: Gas limits for individual transactions (grouped by size)
3. **Block Gas Limits**: Maximum gas limit per block
4. **Block Gas Used**: Actual gas consumed per block
5. **30-block Avg Gas Price**: Rolling average of gas prices

Interactive features:
- **Tooltips**: Hover over data points for details
- **Zoom**: Use the slider or scroll to zoom into specific block ranges
- **Legend Toggle**: Click legend items to show/hide series

## Example Use Cases

Analyzing gas price patterns during network congestion:

```bash
polycli plot --rpc-url https://eth-mainnet.g.alchemy.com/v2/YOUR_KEY \
  --start-block 18000000 \
  --end-block 18001000 \
  --output mainnet_congestion.html
```

Tracking your contract deployment gas costs:

```bash
polycli plot --rpc-url http://localhost:8545 \
  --target-address 0xYourContractAddress \
  --output my_contract_gas.html
```

Analyzing test network behavior:

```bash
polycli plot --rpc-url http://localhost:8545 \
  --concurrency 1 \
  --rate-limit 4 \
  --output local_test.html
```

## Caching Block Data

To avoid re-fetching block data from the RPC endpoint, use the `--cache` flag to store data in an NDJSON file:

```bash
# First run: fetches from RPC and writes to cache
polycli plot --rpc-url http://localhost:8545 \
  --start-block 1000 \
  --end-block 2000 \
  --cache blocks.ndjson

# Subsequent runs: reads from cache (much faster)
polycli plot --rpc-url http://localhost:8545 \
  --start-block 1000 \
  --end-block 2000 \
  --cache blocks.ndjson

# Target address filtering works with cache
polycli plot --rpc-url http://localhost:8545 \
  --cache blocks.ndjson \
  --target-address 0xYourAddress
```

This is useful when iterating on chart options (renderer, target-address) without waiting for RPC queries.

## Notes

- If `--start-block` is not set, the command analyzes the last 500 blocks
- If `--end-block` exceeds the latest block or is not set, it defaults to the latest block
- The command logs the 20 most frequently used gas prices at debug level
- Charts are saved in HTML format with interactive features
- Cache files use NDJSON format (one JSON object per line) for efficient streaming
