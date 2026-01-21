# `polycli loadtest uniswapv3`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Run UniswapV3-like load test against an Eth/EVM style JSON-RPC endpoint.

```bash
polycli loadtest uniswapv3 [flags]
```

## Usage

The `uniswapv3` command is a subcommand of the `loadtest` tool. It is meant to generate UniswapV3-like load against JSON-RPC endpoints.

You can either chose to deploy the full UniswapV3 contract suite.

```sh
polycli loadtest uniswapv3
```

Or to use pre-deployed contracts to speed up the process.

```bash
polycli loadtest uniswapv3 \
  --uniswap-factory-v3-address 0xc5f46e00822c828e1edcc12cf98b5a7b50c9e81b \
  --uniswap-migrator-address 0x24951726c5d22a3569d5474a1e74734a09046cd9 \
  --uniswap-multicall-address 0x0e695f36ade2a12abea51622e80f105e125d1d6e \
  --uniswap-nft-descriptor-lib-address 0x23050ec03bb24308c788300428a8f9c247f28b25 \
  --uniswap-nft-position-descriptor-address 0xea43847a98b671211b0e412849b69bbd7d53fd00 \
  --uniswap-non-fungible-position-manager-address 0x58eabc23408fb7896b7ce943828cc00044786449 \
  --uniswap-proxy-admin-address 0xdba55eb96288eac85974376b25b3c3f3d67399b7 \
  --uniswap-quoter-v2-address 0x91464a00c4aae9dca6d503a2c24b1dfb8c279e50 \
  --uniswap-staker-address 0xc87383ece9ee3ad3f5158998c4fc04833ba1336e \
  --uniswap-swap-router-address 0x46096eb627d30125f9eaaeefeecaa4e237a04a97 \
  --uniswap-tick-lens-address 0xc73dfb5055874cc7b1cf06ae83f7fe8f6facdb19 \
  --uniswap-upgradeable-proxy-address 0x28656635b0ecd600801600475d61e3ec1534de6e \
  --weth9-address 0x5570d4fd7cce73f0135536d83b8d49e6b77bb76c \
  --uniswap-pool-token-0-address 0x1ce270d0380fbbead12371286aff578a1227d1d7 \
  --uniswap-pool-token-1-address 0x060f7db3146f3d6748822fb4c69489a04b5f3278
```

Contracts are cloned from the different Uniswap repositories, compiled with a specific version of `solc` and go bindings are generated using `abigen`. To learn more about this process, make sure to check out `contracts/uniswapv3/README.org`.

## Flags

```bash
  -h, --help                                                   help for uniswapv3
  -f, --pool-fees float                                        trading fees for UniswapV3 liquidity pool swaps (e.g. 0.3 means 0.3%) (default 0.3)
  -a, --swap-amount uint                                       amount of inbound token given as swap input (default 1000)
      --uniswap-factory-v3-address string                      address of pre-deployed UniswapFactoryV3 contract
      --uniswap-migrator-address string                        address of pre-deployed Migrator contract
      --uniswap-multicall-address string                       address of pre-deployed Multicall contract
      --uniswap-nft-descriptor-lib-address string              address of pre-deployed NFTDescriptor library contract
      --uniswap-nft-position-descriptor-address string         address of pre-deployed NonfungibleTokenPositionDescriptor contract
      --uniswap-non-fungible-position-manager-address string   address of pre-deployed NonfungiblePositionManager contract
      --uniswap-pool-token-0-address string                    address of pre-deployed ERC20 contract used in Uniswap pool Token0 // Token1
      --uniswap-pool-token-1-address string                    address of pre-deployed ERC20 contract used in Uniswap pool Token0 // Token1
      --uniswap-proxy-admin-address string                     address of pre-deployed ProxyAdmin contract
      --uniswap-quoter-v2-address string                       address of pre-deployed QuoterV2 contract
      --uniswap-staker-address string                          address of pre-deployed Staker contract
      --uniswap-swap-router-address string                     address of pre-deployed SwapRouter contract
      --uniswap-tick-lens-address string                       address of pre-deployed TickLens contract
      --uniswap-upgradeable-proxy-address string               address of pre-deployed TransparentUpgradeableProxy contract
      --weth9-address string                                   address of pre-deployed WETH9 contract
```

The command also inherits flags from parent commands.

```bash
      --adaptive-backoff-factor float          multiplicative decrease factor for adaptive rate limiting (default 2)
      --adaptive-cycle-duration-seconds uint   interval in seconds to check queue size and adjust rates for adaptive rate limiting (default 10)
      --adaptive-rate-limit                    enable AIMD-style congestion control to automatically adjust request rate
      --adaptive-rate-limit-increment uint     size of additive increases for adaptive rate limiting (default 50)
      --adaptive-target-size uint              target queue size for adaptive rate limiting (speed up if smaller, back off if larger) (default 1000)
      --batch-size uint                        batch size for receipt fetching (default: 999) (default 999)
      --chain-id uint                          chain ID for the transactions
  -c, --concurrency int                        number of requests to perform concurrently (default: one at a time) (default 1)
      --config string                          config file (default is $HOME/.polygon-cli.yaml)
      --eth-amount-in-wei uint                 amount of ether in wei to send per transaction
      --eth-call-only                          call contracts without sending transactions (incompatible with adaptive rate limiting and summarization)
      --eth-call-only-latest                   execute on latest block instead of original block in call-only mode with recall
      --fire-and-forget                        send transactions and load without waiting for it to be mined
      --gas-limit uint                         manually specify gas limit (useful to avoid eth_estimateGas or when auto-computation fails)
      --gas-price uint                         manually specify gas price (useful when auto-detection fails)
      --gas-price-multiplier float             a multiplier to increase or decrease the gas price (default 1)
      --legacy                                 send a legacy transaction instead of an EIP1559 transaction
      --nonce uint                             use this flag to manually set the starting nonce
      --output-mode string                     format mode for summary output (json | text) (default "text")
      --output-raw-tx-only                     output raw signed transaction hex without sending (works with most modes except RPC and UniswapV3)
      --pretty-logs                            output logs in pretty format instead of JSON (default true)
      --priority-gas-price uint                gas tip price for EIP-1559 transactions
      --private-key string                     hex encoded private key to use for sending transactions (default "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa")
      --random-recipients                      send to random addresses instead of fixed address in transfer tests
      --rate-limit float                       requests per second limit (use negative value to remove limit) (default 4)
  -n, --requests int                           number of requests to perform for the benchmarking session (default of 1 leads to non-representative results) (default 1)
  -r, --rpc-url string                         the RPC endpoint URL (default "http://localhost:8545")
      --seed int                               a seed for generating random values and addresses (default 123456)
      --send-only                              alias for --fire-and-forget
      --summarize                              produce execution summary after load test (can take a long time for large tests)
  -t, --time-limit int                         maximum seconds to spend benchmarking (default: no limit) (default -1)
      --to-address string                      recipient address for transactions (default "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF")
  -v, --verbosity string                       log level (string or int):
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

- [polycli loadtest](polycli_loadtest.md) - Run a generic load test against an Eth/EVM style JSON-RPC endpoint.
