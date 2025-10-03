# `polycli loadtest`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Run a generic load test against an Eth/EVM style JSON-RPC endpoint.

```bash
polycli loadtest [flags]
```

## Usage

The `loadtest` tool is meant to generate various types of load against RPC end points. It leverages the [`ethclient`](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient) library Go Ethereum to interact with the blockchain.x

```bash
$ polycli wallet inspect  --mnemonic "code code code code code code code code code code code quality" --addresses 1
```

The `--mode` flag is important for this command.

- `t`/`transaction` will perform ETH transfers. This is the simplest
  and cheapest transaction that can be performed.
- `d`/`deploy` will deploy the load testing contract over and over
  again.
- `2`/`erc20` will run an ERC20 transfer test. The process initializes
  by minting a large amount of tokens then transferring it in small
  amounts. Each transaction is a single transfer.
- `7`/`erc721` will run an ERC721 mint test which will mint an NFT
  over and over again.
- `i`/`inc`/`increment` will call the increment function repeatedly on
  the load test contract. It's a minimal example of a contract call
  that will require an update to a contract's storage.
- `s`/`store` is used to store random data in the smart contract
  storage. The amount of data stored per transaction is controlled
  with the `store-data-size` flag.
- `R`/`recall` will attempt to replay all of the transactions from the
  previous blocks. You can use `--recall-blocks` to specify how many
  previous blocks should be used to seed transaction history. It's
  expected that many of the transactions in this mode would fail.
- `r`/`random` will call any of the other modes randomly. This mode
  shouldn't be used in combination with other modes. Ideally this is a
  good way to generate a lot of random activity on a test network.
- `rpc` is a unique mode that won't just simulate transactions, it
  will simulate RPC traffic (e.g. calls to get transaction receipt or
  filter logs). This is meant to stress test RPC servers rather than
  full blockchain networks. The approach is similar to `recall` mode
  where we'll fetch some recent blocks and then use that data to
  generate a variety of calls to the RPC server.

The default private key is: `42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa`. We can use `wallet inspect` to get more information about this address, in particular its `ETHAddress` if you want to check balance or pre-mine value for this particular account.

Here is a simple example that runs 1000 requests at a max rate of 1 request per second against the http rpc endpoint on localhost. It's running in transaction mode so it will perform simple transactions send to the default address.

```bash
$ polycli loadtest --verbosity 700 --chain-id 1256 --concurrency 1 --requests 1000 --rate-limit 1 --mode t --rpc-url http://localhost:8888
```

### Load Test Contract

The codebase has a contract that used for load testing. It's written in Solidity. The workflow for modifying this contract is.

1. Make changes to <file:contracts/LoadTester.sol>
2. Compile the contracts:
   - `$ solc LoadTester.sol --bin --abi -o . --overwrite`
3. Run `abigen`
   - `$ abigen --abi LoadTester.abi --pkg contracts --type LoadTester --bin LoadTester.bin --out loadtester.go`
4. Run the loadtester to ensure it deploys and runs successfully
   - `$ polycli loadtest --verbosity 700 --rpc-url http://127.0.0.1:8541`

## Flags

```bash
      --account-funding-amount big.Int         amount in wei to fund sending accounts (set to 0 to disable)
      --adaptive-backoff-factor float          multiplicative decrease factor for adaptive rate limiting (default 2)
      --adaptive-cycle-duration-seconds uint   interval in seconds to check queue size and adjust rates for adaptive rate limiting (default 10)
      --adaptive-rate-limit                    enable AIMD-style congestion control to automatically adjust request rate
      --adaptive-rate-limit-increment uint     size of additive increases for adaptive rate limiting (default 50)
      --adaptive-target-size uint              target queue size for adaptive rate limiting (speed up if smaller, back off if larger) (default 1000)
      --batch-size uint                        batch size for receipt fetching (default: 999) (default 999)
      --blob-fee-cap uint                      blob fee cap, or maximum blob fee per chunk, in Gwei (default 100000)
      --calldata string                        hex encoded calldata: function signature + encoded arguments (requires --mode contract-call and --contract-address)
      --chain-id uint                          chain ID for the transactions
  -c, --concurrency int                        number of requests to perform concurrently (default: one at a time) (default 1)
      --contract-address string                contract address for --mode contract-call (requires --calldata)
      --contract-call-payable                  mark function as payable using value from --eth-amount-in-wei (requires --mode contract-call and --contract-address)
      --erc20-address string                   address of pre-deployed ERC20 contract
      --erc721-address string                  address of pre-deployed ERC721 contract
      --eth-amount-in-wei uint                 amount of ether in wei to send per transaction
      --eth-call-only                          call contracts without sending transactions (incompatible with adaptive rate limiting and summarization)
      --eth-call-only-latest                   execute on latest block instead of original block in call-only mode with recall
      --fire-and-forget                        send transactions and load without waiting for it to be mined
      --gas-limit uint                         manually specify gas limit (useful to avoid eth_estimateGas or when auto-computation fails)
      --gas-price uint                         manually specify gas price (useful when auto-detection fails)
      --gas-price-multiplier float             a multiplier to increase or decrease the gas price (default 1)
  -h, --help                                   help for loadtest
      --legacy                                 send a legacy transaction instead of an EIP1559 transaction
      --loadtest-contract-address string       address of pre-deployed load test contract
      --max-base-fee-wei uint                  maximum base fee in wei (pause sending new transactions when exceeded, useful during network congestion)
  -m, --mode strings                           testing mode (can specify multiple like "d,t"):
                                               2, erc20 - send ERC20 tokens
                                               7, erc721 - mint ERC721 tokens
                                               b, blob - send blob transactions
                                               cc, contract-call - make contract calls
                                               d, deploy - deploy contracts
                                               inc, increment - increment a counter
                                               r, random - random modes (excludes: blob, call, inscription, recall, rpc, uniswapv3)
                                               R, recall - replay or simulate transactions
                                               rpc - call random rpc methods
                                               s, store - store bytes in a dynamic byte array
                                               t, transaction - send transactions
                                               v3, uniswapv3 - perform UniswapV3 swaps (default [t])
      --nonce uint                             use this flag to manually set the starting nonce
      --output-mode string                     format mode for summary output (json | text) (default "text")
      --output-raw-tx-only                     output raw signed transaction hex without sending (works with most modes except RPC and UniswapV3)
      --pre-fund-sending-accounts              fund all sending accounts at start instead of on first use
      --priority-gas-price uint                gas tip price for EIP-1559 transactions
      --private-key string                     hex encoded private key to use for sending transactions (default "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa")
      --proxy string                           use the proxy specified
      --random-recipients                      send to random addresses instead of fixed address in transfer tests
      --rate-limit float                       requests per second limit (use negative value to remove limit) (default 4)
      --recall-blocks uint                     number of blocks that we'll attempt to fetch for recall (default 50)
      --receipt-retry-initial-delay-ms uint    initial delay in milliseconds for receipt polling (uses exponential backoff with jitter) (default 100)
      --receipt-retry-max uint                 maximum polling attempts for transaction receipt with --wait-for-receipt (default 30)
      --refund-remaining-funds                 refund remaining balance to funding account after completion
  -n, --requests int                           number of requests to perform for the benchmarking session (default of 1 leads to non-representative results) (default 1)
  -r, --rpc-url string                         the RPC endpoint URL (default "http://localhost:8545")
      --seed int                               a seed for generating random values and addresses (default 123456)
      --send-only                              alias for --fire-and-forget
      --sending-accounts-count uint            number of sending accounts to use (avoids pool account queue)
      --sending-accounts-file string           file with sending account private keys, one per line (avoids pool queue and preserves accounts across runs)
      --store-data-size uint                   number of bytes to store in contract for store mode (default 1024)
      --summarize                              produce execution summary after load test (can take a long time for large tests)
  -t, --time-limit int                         maximum seconds to spend benchmarking (default: no limit) (default -1)
      --to-address string                      recipient address for transactions (default "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF")
      --wait-for-receipt                       wait for transaction receipt to be mined instead of just sending
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
- [polycli loadtest uniswapv3](polycli_loadtest_uniswapv3.md) - Run Uniswapv3-like load test against an Eth/EVm style JSON-RPC endpoint.

