# `polycli loadtest uniswapv3`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Run Uniswapv3-like load test against an Eth/EVm style JSON-RPC endpoint.

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
  -f, --pool-fees float                                        Trading fees charged on each swap or trade made within a UniswapV3 liquidity pool (e.g. 0.3 means 0.3%) (default 0.3)
  -a, --swap-amount uint                                       The amount of inbound token given as swap input (default 1000)
      --uniswap-factory-v3-address string                      The address of a pre-deployed UniswapFactoryV3 contract
      --uniswap-migrator-address string                        The address of a pre-deployed Migrator contract
      --uniswap-multicall-address string                       The address of a pre-deployed Multicall contract
      --uniswap-nft-descriptor-lib-address string              The address of a pre-deployed NFTDescriptor library contract
      --uniswap-nft-position-descriptor-address string         The address of a pre-deployed NonfungibleTokenPositionDescriptor contract
      --uniswap-non-fungible-position-manager-address string   The address of a pre-deployed NonfungiblePositionManager contract
      --uniswap-pool-token-0-address string                    The address of a pre-deployed ERC20 contract used in the Uniswap pool Token0 // Token1
      --uniswap-pool-token-1-address string                    The address of a pre-deployed ERC20 contract used in the Uniswap pool Token0 // Token1
      --uniswap-proxy-admin-address string                     The address of a pre-deployed ProxyAdmin contract
      --uniswap-quoter-v2-address string                       The address of a pre-deployed QuoterV2 contract
      --uniswap-staker-address string                          The address of a pre-deployed Staker contract
      --uniswap-swap-router-address string                     The address of a pre-deployed SwapRouter contract
      --uniswap-tick-lens-address string                       The address of a pre-deployed TickLens contract
      --uniswap-upgradeable-proxy-address string               The address of a pre-deployed TransparentUpgradeableProxy contract
      --weth9-address string                                   The address of a pre-deployed WETH9 contract
```

The command also inherits flags from parent commands.

```bash
      --adaptive-backoff-factor float          When using adaptive rate limiting, this flag controls our multiplicative decrease value. (default 2)
      --adaptive-cycle-duration-seconds uint   When using adaptive rate limiting, this flag controls how often we check the queue size and adjust the rates (default 10)
      --adaptive-rate-limit                    Enable AIMD-style congestion control to automatically adjust request rate
      --adaptive-rate-limit-increment uint     When using adaptive rate limiting, this flag controls the size of the additive increases. (default 50)
      --batch-size uint                        Number of batches to perform at a time for receipt fetching. Default is 999 requests at a time. (default 999)
      --call-only                              When using this mode, rather than sending a transaction, we'll just call. This mode is incompatible with adaptive rate limiting, summarization, and a few other features.
      --call-only-latest                       When using call only mode with recall, should we execute on the latest block or on the original block
      --chain-id uint                          The chain id for the transactions.
  -c, --concurrency int                        Number of requests to perform concurrently. Default is one request at a time. (default 1)
      --config string                          config file (default is $HOME/.polygon-cli.yaml)
      --eth-amount float                       The amount of ether to send on every transaction (default 0.001)
      --gas-limit uint                         In environments where the gas limit can't be computed on the fly, we can specify it manually. This can also be used to avoid eth_estimateGas
      --gas-price uint                         In environments where the gas price can't be determined automatically, we can specify it manually
  -i, --iterations uint                        If we're making contract calls, this controls how many times the contract will execute the instruction in a loop. If we are making ERC721 Mints, this indicates the minting batch size (default 1)
      --legacy                                 Send a legacy transaction instead of an EIP1559 transaction.
      --output-mode string                     Format mode for summary output (json | text) (default "text")
      --pretty-logs                            Should logs be in pretty format or JSON (default true)
      --priority-gas-price uint                Specify Gas Tip Price in the case of EIP-1559
      --private-key string                     The hex encoded private key that we'll use to send transactions (default "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa")
      --rate-limit float                       An overall limit to the number of requests per second. Give a number less than zero to remove this limit all together (default 4)
  -n, --requests int                           Number of requests to perform for the benchmarking session. The default is to just perform a single request which usually leads to non-representative benchmarking results. (default 1)
  -r, --rpc-url string                         The RPC endpoint url (default "http://localhost:8545")
      --seed int                               A seed for generating random values and addresses (default 123456)
      --send-only                              Send transactions and load without waiting for it to be mined.
      --steady-state-tx-pool-size uint         When using adaptive rate limiting, this value sets the target queue size. If the queue is smaller than this value, we'll speed up. If the queue is smaller than this value, we'll back off. (default 1000)
      --summarize                              Should we produce an execution summary after the load test has finished. If you're running a large load test, this can take a long time
  -t, --time-limit int                         Maximum number of seconds to spend for benchmarking. Use this to benchmark within a fixed total amount of time. Per default there is no time limit. (default -1)
      --to-address string                      The address that we're going to send to (default "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF")
      --to-random                              When doing a transfer test, should we send to random addresses rather than DEADBEEFx5
  -v, --verbosity int                          0 - Silent
                                               100 Panic
                                               200 Fatal
                                               300 Error
                                               400 Warning
                                               500 Info
                                               600 Debug
                                               700 Trace (default 500)
```

## See also

- [polycli loadtest](polycli_loadtest.md) - Run a generic load test against an Eth/EVM style JSON-RPC endpoint.
