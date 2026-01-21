# plot

`plot` generates visual charts analyzing transaction gas prices and limits across a range of blocks. It fetches block data from an Ethereum-compatible RPC endpoint and produces a PNG chart showing gas price distribution, transaction gas limits, block gas limits, and block gas usage over time.

## Basic Usage

Generate a chart for the last 500 blocks:

```bash
polycli plot --rpc-url http://localhost:8545
```

This will create a file named `tx_gasprice_chart.png` in the current directory.

## Analyzing Specific Block Ranges

Analyze blocks 9356826 to 9358826:

```bash
polycli plot --rpc-url https://sepolia.infura.io/v3/YOUR_API_KEY \
  --start-block 9356826 \
  --end-block 9358826 \
  --output "sepolia_analysis.png"
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

## Chart Scale Options

Choose between logarithmic (default) and linear scale for the gas price axis.

Use linear scale for gas prices:

```bash
polycli plot --rpc-url http://localhost:8545 \
  --scale "linear" \
  --output "linear_chart.png"
```

Use logarithmic scale (default):

```bash
polycli plot --rpc-url http://localhost:8545 \
  --scale "log" \
  --output "log_chart.png"
```

## Understanding the Chart

The generated chart displays four key metrics:

1. **Transaction Gas Prices**: Individual transaction gas prices plotted as points, with target address transactions highlighted
2. **Transaction Gas Limits**: Gas limits for individual transactions
3. **Block Gas Limits**: Maximum gas limit per block
4. **Block Gas Used**: Actual gas consumed per block

## Example Use Cases

Analyzing gas price patterns during network congestion:

```bash
polycli plot --rpc-url https://eth-mainnet.g.alchemy.com/v2/YOUR_KEY \
  --start-block 18000000 \
  --end-block 18001000 \
  --scale log \
  --output mainnet_congestion.png
```

Tracking your contract deployment gas costs:

```bash
polycli plot --rpc-url http://localhost:8545 \
  --target-address 0xYourContractAddress \
  --output my_contract_gas.png
```

Analyzing test network behavior:

```bash
polycli plot --rpc-url http://localhost:8545 \
  --concurrency 1 \
  --rate-limit 4 \
  --output local_test.png
```

## Notes

- If `--start-block` is not set, the command analyzes the last 500 blocks
- If `--end-block` exceeds the latest block or is not set, it defaults to the latest block
- The command logs the 20 most frequently used gas prices at debug level
- Charts are saved in PNG format
