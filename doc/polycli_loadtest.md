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

The `loadtest` tool is meant to generate various types of load against RPC end points. It leverages the [`ethclient`](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient) library Go Ethereum to interact with the blockchain.

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
- `inc`/`increment` will call the increment function repeatedly on
  the load test contract. It's a minimal example of a contract call
  that will require an update to a contract's storage.
- `s`/`store` is used to store random data in the smart contract
  storage. The amount of data stored per transaction is controlled
  with the `store-data-size` flag.
- `b`/`blob` will send EIP-4844 blob transactions. Use `--blob-fee-cap`
  to set the maximum blob fee per chunk.
- `cc`/`contract-call` will call a specific contract function. Requires
  `--contract-address` and either `--calldata` (hex string) or
  `--calldata-file` (path to a file containing the hex calldata). Use
  `--contract-call-payable` if the function is payable.
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
- `v3`/`uniswapv3` will deploy UniswapV3 contracts and perform token
  swaps. This mode can also be run as a subcommand (`polycli loadtest
  uniswapv3`) which provides additional flags for specifying
  pre-deployed contract addresses, pool fees, and swap amounts.

The default private key is: `42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa`. We can use `wallet inspect` to get more information about this address, in particular its `ETHAddress` if you want to check balance or pre-mine value for this particular account.

Here is a simple example that runs 1000 requests at a max rate of 1 request per second against the http rpc endpoint on localhost. It's running in transaction mode so it will perform simple transactions send to the default address.

```bash
$ polycli loadtest --verbosity 700 --chain-id 1256 --concurrency 1 --requests 1000 --rate-limit 1 --mode t --rpc-url http://localhost:8888
```

### Separate Broadcast Endpoint

By default, all RPC calls (gas estimation, chain ID, nonces, receipts, and transaction broadcast) go to `--rpc-url`. The `--send-rpc-url` flag routes only the transaction broadcast (`eth_sendRawTransaction`, or `eth_sendRawTransactionPrivate` when combined with `--private-txs`) to a secondary endpoint while everything else, including account funding, stays on `--rpc-url`. This is useful for:

- **Private mempools**: an endpoint that only accepts `eth_sendRawTransactionPrivate`.
- **Gossip-only broadcasters**: a light client connected to the p2p network that can broadcast transactions but has no chain state to answer other queries.

```bash
$ polycli loadtest --rpc-url http://fullnode:8545 --send-rpc-url http://broadcaster:8545 --mode t
```

Like `--private-txs`, this flag is only supported by the modes that broadcast transactions explicitly: `transaction`, `blob`, `contract-call`, and `recall`.

### Synchronous Sending (EIP-7966)

`--sync-txs` broadcasts with [`eth_sendRawTransactionSync`][eip-7966] instead of
`eth_sendRawTransaction`. That call does not return when the transaction is accepted; it
returns when the node has a receipt for it, or fails when its timeout elapses. On a chain
that preconfirms, the receipt can arrive well before canonical inclusion, which is what
makes this useful for measuring preconfirmation latency.

```bash
$ polycli loadtest --rpc-url http://localhost:8545 --mode t --sync-txs --concurrency 20
```

`--sync-tx-timeout` sets how long the node should wait, in whole milliseconds. Leave it at
`0` to omit the parameter so the node applies its own default. EIP-7966 recommends 2
seconds; bor defaults to 20 seconds (`--rpc.txsync.defaulttimeout`) and silently clamps
anything above `--rpc.txsync.maxtimeout`, 1 minute by default, rather than rejecting it.

The value goes on the wire as a hex quantity, because bor takes the parameter as
`*hexutil.Uint64`, which unmarshals only from a quoted hex string and rejects a bare JSON
number. EIP-7966 describes the parameter as an integer instead, and servers implementing it
literally reject the hex form, so `--sync-tx-timeout-int` switches encodings. Against bor,
leave it off.

**This changes what the latency numbers mean.** Every other send path measures time to
*accept* a transaction; with `--sync-txs` the recorded request duration is time to
*receipt*. Because each request now blocks for the life of the transaction, `--concurrency`
becomes the number of transactions in flight rather than the number of submissions per
moment, and you will need a much higher value to reach the same send rate.

At the end of a run the sync tracker logs a summary:

- `receipts`, split into `speculative` and `canonical`
- `reverted` and `no_status` for receipts that came back unsuccessful or without a status
- `timeouts`, `queued`, `nonce_gaps` — EIP-7966 error codes 4, 5 and 6. A timeout means the
  transaction reached the mempool but produced no receipt in time; it may still be mined
  afterwards. A nonce gap carries the node's expected nonce in the error data.
- `rejected` — the transaction was refused before the wait began
- `p50_ms`, `p90_ms`, `p99_ms` of the synchronous call

Bor implements only code 4 of that set. A transaction it refuses to accept fails before the
wait starts and comes back with bor's own codes, which are counted too: `-38011`
(nonce too high) as a nonce gap, and `-38010` (nonce too low), `-38013` (intrinsic gas),
`-38014` (insufficient funds) and `-38026` (client limit exceeded) as `rejected`.

The speculative/canonical split comes from the receipt itself where possible. Bor's
preconfirmation pipeline marks a preconfirmed receipt with `"preconfirmation": true` and
leaves its `blockHash` null, and that marker is taken as authoritative. EIP-7966 defines no
such field, so for any other node the split falls back to reading the block fields: a
receipt with no block cannot be canonical. A node that both omits the marker and fills in a
speculative block number is indistinguishable from one answering canonically. To confirm
canonical inclusion independently, combine with `--wait-for-receipt`, which polls for the
receipt after the synchronous call returns — that pairing measures the speculative receipt
and the canonical one separately.

Note that stock bor without preconfirmations enabled only ever returns canonical receipts:
it waits on chain events and requires both a block number and a block hash before
answering, so `speculative` stays at zero. The speculative path needs bor's
preconfirmation pipeline (`SubmitTxForPreconf` and the preconfirmation receipt index).

At verbosity 700 (trace) each synchronous submission is also logged individually: the
node's raw receipt verbatim as JSON on success, or the RPC error with its code and data on
failure. Pair with `--pretty-logs=false` for machine-parseable lines that can be checked for
correctness later, e.g. by diffing against `eth_getTransactionReceipt`:

```bash
$ polycli loadtest --rpc-url http://localhost:8545 --mode t --sync-txs -v 700 --pretty-logs=false
```

Like `--private-txs`, `--sync-txs` is only supported by `transaction`, `blob`,
`contract-call`, and `recall`, and the two flags are mutually exclusive.

[eip-7966]: https://eips.ethereum.org/EIPS/eip-7966

### Waiting for Receipts

`--wait-for-receipt` polls `eth_getTransactionReceipt` after each send, in the same
goroutine that sent the transaction, so each worker blocks until its transaction is mined
before sending the next one. By default polling uses exponential backoff with jitter,
starting at `--receipt-retry-initial-delay-ms` and giving up after `--receipt-retry-max`
attempts or one minute, whichever comes first.

`--receipt-poll-interval` switches to polling at a fixed interval instead. In this mode
`--receipt-retry-max` is ignored and polling is bounded only by the one-minute timeout.
Backoff can overshoot the moment the receipt appeared by the length of the current backoff
step, so a small fixed interval also makes the wait duration a usable receipt-latency
measurement, at the cost of steadier RPC load:

```bash
$ polycli loadtest --rpc-url http://localhost:8545 --mode t --wait-for-receipt --receipt-poll-interval 50ms
```

At verbosity 700 (trace) each raw receipt is logged verbatim as JSON with the tx hash and
the wait duration, the same shape as the `--sync-txs` receipt logs. Combined with
`--sync-txs`, this logs the speculative receipt and the canonical one separately per
transaction.

### Gas Manager

The loadtest command includes an optional gas manager for controlling transaction gas limits and pricing. Enable it with `--gas-manager-enabled`, then use the `--gas-manager-*` flags to:

- **Oscillate gas limits** with wave patterns (flat, sine, square, triangle, sawtooth)
- **Control gas pricing** with strategies (estimated, fixed, dynamic)

Example with sine wave oscillation:
```bash
$ polycli loadtest --rpc-url http://localhost:8545 \
  --gas-manager-enabled \
  --gas-manager-oscillation-wave sine \
  --gas-manager-target 20000000 \
  --gas-manager-amplitude 10000000 \
  --gas-manager-period 100
```

See [Gas Manager README](../../loadtest/gasmanager/README.md) for detailed documentation.

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
      --account-funding-amount big.Int                   amount in wei to fund sending accounts (set to 0 to disable)
      --accounts-per-funding-tx uint                     number of accounts to fund per multicall3 transaction (default 400)
      --adaptive-backoff-factor float                    multiplicative decrease factor for adaptive rate limiting (default 2)
      --adaptive-cycle-duration-seconds uint             interval in seconds to check queue size and adjust rates for adaptive rate limiting (default 10)
      --adaptive-rate-limit                              enable AIMD-style congestion control to automatically adjust request rate
      --adaptive-rate-limit-increment uint               size of additive increases for adaptive rate limiting (default 50)
      --adaptive-target-size uint                        target queue size for adaptive rate limiting (speed up if smaller, back off if larger) (default 1000)
      --batch-size uint                                  batch size for receipt fetching (default: 999) (default 999)
      --blob-fee-cap uint                                blob fee cap, or maximum blob fee per chunk, in Gwei (default 100000)
      --block-batch-size uint                            number of blocks to fetch per RPC batch request for recall and rpc modes (default 25)
      --calldata string                                  hex encoded calldata: function signature + encoded arguments (requires --mode contract-call and --contract-address)
      --calldata-file string                             path to a file containing hex encoded calldata (alternative to --calldata; mutually exclusive with it)
      --chain-id uint                                    chain ID for the transactions
      --check-balance-before-funding                     check account balance before funding sending accounts (saves gas when accounts are already funded)
      --check-preconf                                    check for preconf status after sending tx
  -c, --concurrency int                                  number of requests to perform concurrently (default: one at a time) (default 1)
      --contract-address string                          contract address for --mode contract-call (requires --calldata)
      --contract-call-payable                            mark function as payable using value from --eth-amount-in-wei (requires --mode contract-call and --contract-address)
      --dump-sending-accounts-file string                file path to dump generated private keys when using --sending-accounts-count
      --duplicate-nonce-rate float                       ratio of duplicate-nonce txs to fresh txs (0 disables; 1 = 50% duplicates, 4 = 80%); requires --fire-and-forget
      --erc20-address string                             address of pre-deployed ERC20 contract
      --erc721-address string                            address of pre-deployed ERC721 contract
      --eth-amount-in-wei uint                           amount of ether in wei to send per transaction
      --eth-call-only                                    call contracts without sending transactions (incompatible with adaptive rate limiting and summarization)
      --eth-call-only-latest                             execute on latest block instead of original block in call-only mode with recall
      --fire-and-forget                                  send transactions and load without waiting for it to be mined
      --gas-limit uint                                   manually specify gas limit (useful to avoid eth_estimateGas or when auto-computation fails)
      --gas-manager-amplitude uint                       amplitude for oscillation wave
      --gas-manager-dynamic-gas-prices-variation float   variation percentage for dynamic strategy (default 0.3)
      --gas-manager-dynamic-gas-prices-wei string        comma-separated gas prices in wei for dynamic strategy (default "0,1000000,0,10000000,0,100000000")
      --gas-manager-enabled                              enable block-based gas manager (oscillation wave + gas budget vault)
      --gas-manager-fixed-gas-price-wei uint             fixed gas price in wei (default 300000000)
      --gas-manager-oscillation-wave string              type of oscillation wave (flat | sine | square | triangle | sawtooth) (default "flat")
      --gas-manager-period uint                          period in blocks for oscillation wave (default 1)
      --gas-manager-price-strategy string                gas price strategy (estimated | fixed | dynamic) (default "estimated")
      --gas-manager-target uint                          target gas limit for oscillation wave (default 30000000)
      --gas-price gas                                    gas price with unit support (e.g., "100gwei", "1000000000")
      --gas-price-multiplier float                       a multiplier to increase or decrease the gas price (default 1)
  -h, --help                                             help for loadtest
      --legacy                                           send a legacy transaction instead of an EIP1559 transaction
      --loadtest-contract-address string                 address of pre-deployed load test contract
      --max-base-fee-wei uint                            maximum base fee in wei (pause sending new transactions when exceeded, useful during network congestion)
  -m, --mode strings                                     testing mode (can specify multiple like "d,t"):
                                                         2, erc20 - send ERC20 tokens
                                                         7, erc721 - mint ERC721 tokens
                                                         b, blob - send blob transactions
                                                         cc, contract-call - make contract calls
                                                         d, deploy - deploy contracts
                                                         inc, increment - increment a counter
                                                         r, random - random modes (excludes: blob, call, recall, rpc, uniswapv3)
                                                         R, recall - replay or simulate transactions
                                                         rpc - call random rpc methods
                                                         s, store - store bytes in a dynamic byte array
                                                         t, transaction - send transactions
                                                         v3, uniswapv3 - perform UniswapV3 swaps (default [t])
      --nonce uint                                       use this flag to manually set the starting nonce
      --output-mode string                               format mode for summary output (json | text) (default "text")
      --output-raw-tx-only                               output raw signed transaction hex without sending (works with most modes except RPC and UniswapV3)
      --pre-fund-sending-accounts                        fund all sending accounts at start instead of on first use
      --preconf-stats-file string                        path for preconf stats JSON output, updated every 2 seconds
      --priority-gas-price gas                           gas tip for EIP-1559 with unit support (e.g., "2gwei")
      --private-key string                               hex encoded private key to use for sending transactions (default "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa")
      --private-txs                                      send transactions via eth_sendRawTransactionPrivate
      --proxy string                                     use the proxy specified
      --random-recipients                                send to random addresses instead of fixed address in transfer tests
      --rate-limit float                                 requests per second limit (use negative value to remove limit) (default 4)
      --rate-limit-ramp-duration duration                linearly ramp rate limit from max(1% of --rate-limit, 1 TPS) to full --rate-limit over this duration (e.g. 3m; 0 disables ramp)
      --recall-blocks uint                               number of blocks that we'll attempt to fetch for recall (default 50)
      --receipt-poll-interval duration                   fixed interval between receipt polls with --wait-for-receipt; when set, polling is
                                                         bounded only by the receipt timeout and --receipt-retry-max is ignored (0 uses
                                                         exponential backoff with jitter)
      --receipt-retry-initial-delay-ms uint              initial delay in milliseconds for receipt polling (uses exponential backoff with jitter) (default 100)
      --receipt-retry-max uint                           maximum polling attempts for transaction receipt with --wait-for-receipt (default 30)
      --refund-remaining-funds                           refund remaining balance to funding account after completion
  -n, --requests int                                     number of requests to perform for the benchmarking session (default of 1 leads to non-representative results) (default 1)
      --reverse-nonce-order                              send each account's txs in descending nonce order, from highest planned nonce down to the current one, to stress queued vs pending txpool dynamics; total requests must divide evenly across accounts; requires --fire-and-forget
      --rpc-headers string                               custom HTTP headers for RPC requests (format: "key1:value1,key2:value2")
  -r, --rpc-url string                                   the RPC endpoint URL (default "http://localhost:8545")
      --seed int                                         a seed for generating random values and addresses (default 123456)
      --send-only                                        alias for --fire-and-forget
      --send-rpc-url string                              secondary RPC endpoint used only to broadcast transactions (eth_sendRawTransaction / eth_sendRawTransactionPrivate); all other calls use --rpc-url
      --sending-accounts-count uint                      number of sending accounts to use (avoids pool account queue)
      --sending-accounts-file string                     file with sending account private keys, one per line (avoids pool queue and preserves accounts across runs)
      --sequential-nonce-fetch                           fetch nonces one at a time through the rate limiter instead of in parallel bounded by --concurrency
      --stop-on-insufficient-funds                       stop sending from account when it encounters insufficient funds error
      --store-data-size uint                             number of bytes to store in contract for store mode (default 1024)
      --summarize                                        produce execution summary after load test (can take a long time for large tests)
      --sync-tx-timeout duration                         maximum time the node should wait for a receipt with --sync-txs, sent in whole
                                                         milliseconds (0 omits the parameter so the node applies its own default)
      --sync-tx-timeout-int                              send the --sync-tx-timeout value as a bare JSON integer instead of a hex quantity;
                                                         bor wants hex (the default), while servers implementing EIP-7966 literally want an integer
      --sync-txs                                         send transactions via eth_sendRawTransactionSync (EIP-7966), which blocks until
                                                         the node has a receipt; useful for measuring preconfirmation latency
  -t, --time-limit int                                   maximum seconds to spend benchmarking (default: no limit) (default -1)
      --to-address string                                recipient address for transactions (default "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF")
      --wait-for-receipt                                 wait for transaction receipt to be mined instead of just sending
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
- [polycli loadtest uniswapv3](polycli_loadtest_uniswapv3.md) - Run UniswapV3-like load test against an Eth/EVM style JSON-RPC endpoint.

