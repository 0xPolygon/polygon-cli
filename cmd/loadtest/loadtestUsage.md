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

Like `--private-txs`, `--sync-txs` is only supported by `transaction`, `blob`,
`contract-call`, and `recall`, and the two flags are mutually exclusive.

[eip-7966]: https://eips.ethereum.org/EIPS/eip-7966

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
