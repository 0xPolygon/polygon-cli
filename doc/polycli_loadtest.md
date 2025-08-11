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
      --account-funding-amount big.Int         The amount in wei to fund the sending accounts with. Set to 0 to disable account funding (useful for eth-call-only mode or pre-funded accounts).
      --adaptive-backoff-factor float          When using adaptive rate limiting, this flag controls our multiplicative decrease value. (default 2)
      --adaptive-cycle-duration-seconds uint   When using adaptive rate limiting, this flag controls how often we check the queue size and adjust the rates (default 10)
      --adaptive-rate-limit                    Enable AIMD-style congestion control to automatically adjust request rate
      --adaptive-rate-limit-increment uint     When using adaptive rate limiting, this flag controls the size of the additive increases. (default 50)
      --adaptive-target-size uint              When using adaptive rate limiting, this value sets the target queue size. If the queue is smaller than this value, we'll speed up. If the queue is smaller than this value, we'll back off. (default 1000)
      --batch-size uint                        Number of batches to perform at a time for receipt fetching. Default is 999 requests at a time. (default 999)
      --blob-fee-cap uint                      The blob fee cap, or the maximum blob fee per chunk, in Gwei. (default 100000)
      --calldata string                        The hex encoded calldata passed in. The format is function signature + arguments encoded together. This must be paired up with --mode contract-call and --contract-address
      --chain-id uint                          The chain id for the transactions.
  -c, --concurrency int                        Number of requests to perform concurrently. Default is one request at a time. (default 1)
      --contract-address string                The address of the contract that will be used in --mode contract-call. This must be paired up with --mode contract-call and --calldata
      --contract-call-payable                  Use this flag if the function is payable, the value amount passed will be from --eth-amount-in-wei. This must be paired up with --mode contract-call and --contract-address
      --erc20-address string                   The address of a pre-deployed ERC20 contract
      --erc721-address string                  The address of a pre-deployed ERC721 contract
      --eth-amount-in-wei uint                 The amount of ether in wei to send on every transaction
      --eth-call-only                          When using this mode, rather than sending a transaction, we'll just call. This mode is incompatible with adaptive rate limiting, summarization, and a few other features.
      --eth-call-only-latest                   When using call only mode with recall, should we execute on the latest block or on the original block
      --fire-and-forget                        Send transactions and load without waiting for it to be mined.
      --function-arg strings                   The arguments that will be passed to a contract function call. This must be paired up with "--mode contract-call" and "--contract-address". Args can be passed multiple times: "--function-arg 'test' --function-arg 999" or comma separated values "--function-arg "test",9". The ordering of the arguments must match the ordering of the function parameters.
      --function-signature string              The contract's function signature that will be called. The format is '<function name>(<types...>)'. This must be paired up with '--mode contract-call' and '--contract-address'. If the function requires parameters you can pass them with '--function-arg <value>'.
      --gas-limit uint                         In environments where the gas limit can't be computed on the fly, we can specify it manually. This can also be used to avoid eth_estimateGas
      --gas-price uint                         In environments where the gas price can't be determined automatically, we can specify it manually
      --gas-price-multiplier float             A multiplier to increase or decrease the gas price (default 1)
  -h, --help                                   help for loadtest
      --inscription-content string             The inscription content that will be encoded as calldata. This must be paired up with --mode inscription (default "data:,{\"p\":\"erc-20\",\"op\":\"mint\",\"tick\":\"TEST\",\"amt\":\"1\"}")
      --keep-funds-after-test                  If set to true, the funded amount will be kept in the sending accounts. Otherwise, the funded amount will be refunded back to the account used to fund the account.
      --legacy                                 Send a legacy transaction instead of an EIP1559 transaction.
      --loadtest-contract-address string       The address of a pre-deployed load test contract
  -m, --mode strings                           The testing mode to use. It can be multiple like: "d,t"
                                               2, erc20 - Send ERC20 tokens
                                               7, erc721 - Mint ERC721 tokens
                                               b, blob - Send blob transactions
                                               cc, contract-call - Make contract calls
                                               d, deploy - Deploy contracts
                                               i, inscription - Send inscription transactions
                                               inc, increment - Increment a counter
                                               r, random - Random modes (does not include the following modes: blob, call, inscription, recall, rpc, uniswapv3)
                                               R, recall - Replay or simulate transactions
                                               rpc - Call random rpc methods
                                               s, store - Store bytes in a dynamic byte array
                                               t, transaction - Send transactions
                                               v3, uniswapv3 - Perform UniswapV3 swaps (default [t])
      --nonce uint                             Use this flag to manually set the starting nonce
      --output-mode string                     Format mode for summary output (json | text) (default "text")
      --pre-fund-sending-accounts              If set to true, the sending accounts will be funded at the start of the execution, otherwise all accounts will be funded when used for the first time.
      --priority-gas-price uint                Specify Gas Tip Price in the case of EIP-1559
      --private-key string                     The hex encoded private key that we'll use to send transactions (default "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa")
      --proxy string                           Use the proxy specified
      --random-recipients                      When doing a transfer test, should we send to random addresses rather than DEADBEEFx5
      --rate-limit float                       An overall limit to the number of requests per second. Give a number less than zero to remove this limit all together (default 4)
      --recall-blocks uint                     The number of blocks that we'll attempt to fetch for recall (default 50)
      --receipt-retry-initial-delay-ms uint    Initial delay in milliseconds for receipt polling retry. Uses exponential backoff with jitter. (default 100)
      --receipt-retry-max uint                 Maximum number of attempts to poll for transaction receipt when --wait-for-receipt is enabled. (default 30)
  -n, --requests int                           Number of requests to perform for the benchmarking session. The default is to just perform a single request which usually leads to non-representative benchmarking results. (default 1)
  -r, --rpc-url string                         The RPC endpoint url (default "http://localhost:8545")
      --seed int                               A seed for generating random values and addresses (default 123456)
      --send-only                              Alias for --fire-and-forget.
      --sending-accounts-count uint            The number of sending accounts to use. This is useful for avoiding pool account queue. (default 1)
      --sending-accounts-file string           The file containing the sending accounts private keys, one per line. This is useful for avoiding pool account queue but also to keep the same sending accounts for different execution cycles.
      --store-data-size uint                   If we're in store mode, this controls how many bytes we'll try to store in our contract (default 1024)
      --summarize                              Should we produce an execution summary after the load test has finished. If you're running a large load test, this can take a long time
  -t, --time-limit int                         Maximum number of seconds to spend for benchmarking. Use this to benchmark within a fixed total amount of time. Per default there is no time limit. (default -1)
      --to-address string                      The address that we're going to send to (default "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF")
      --wait-for-receipt                       If set to true, the load test will wait for the transaction receipt to be mined. If set to false, the load test will not wait for the transaction receipt and will just send the transaction.
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Panic
                        200 Fatal
                        300 Error
                        400 Warning
                        500 Info
                        600 Debug
                        700 Trace (default 500)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
- [polycli loadtest uniswapv3](polycli_loadtest_uniswapv3.md) - Run Uniswapv3-like load test against an Eth/EVm style JSON-RPC endpoint.

