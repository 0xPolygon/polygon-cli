Running the sensor will do peer discovery and continue to watch for blocks and
transactions from those peers. This is useful for observing the network for
forks and reorgs without the need to run the entire full node infrastructure.

The sensor can persist data to various backends including ClickHouse, Google Cloud
Datastore, or JSON output. If no nodes.json file exists at the specified path, it
will be created automatically.

The bootnodes may change, so refer to the [Polygon Knowledge Layer][bootnodes]
if the sensor is not discovering peers.

## Data Model

The two persistent backends store different shapes, and are documented separately:

- [ClickHouse data model](/cmd/p2p/sensor/clickhouse.md) — tables and write paths
- [Datastore data model](/cmd/p2p/sensor/datastore.md) — kinds and write paths

The ClickHouse backend is append-only and batched; the Datastore backend does a
read-modify-write per entity. Select one with `--database`, and for ClickHouse pass
`--clickhouse-dsn`. The ClickHouse DDL lives in `clickhouse/schema.sql` in the
sensor-network-tools repo, not here.

## JSON-RPC Server

The sensor runs a JSON-RPC server on port 8545 (configurable via `--rpc-port`)
that supports a subset of Ethereum JSON-RPC methods using cached data.

### Supported Methods

| Method                                  | Description                                        |
| --------------------------------------- | -------------------------------------------------- |
| `eth_chainId`                           | Returns the chain ID                               |
| `eth_blockNumber`                       | Returns the current head block number              |
| `eth_gasPrice`                          | Returns suggested gas price based on recent blocks |
| `eth_getBlockByHash`                    | Returns block by hash                              |
| `eth_getBlockByNumber`                  | Returns block by number (if cached)                |
| `eth_getTransactionByHash`              | Returns transaction by hash                        |
| `eth_getTransactionByBlockHashAndIndex` | Returns transaction at index in block              |
| `eth_getBlockTransactionCountByHash`    | Returns transaction count in block                 |
| `eth_getUncleCountByBlockHash`          | Returns uncle count in block                       |
| `eth_sendRawTransaction`                | Broadcasts signed transaction to peers             |

### Limitations

Methods requiring state or receipts are not supported:

- `eth_getBalance`, `eth_getCode`, `eth_call`, `eth_estimateGas`
- `eth_getTransactionReceipt`, `eth_getLogs`

Data is served from an LRU cache, so older blocks/transactions may not be available.

## Metrics

The sensor exposes Prometheus metrics at `http://localhost:2112/metrics`
(configurable via `--prom-port`). For a complete list of available metrics, see
[polycli_p2p_sensor_metrics.md](polycli_p2p_sensor_metrics.md).

## Rebroadcasting

The sensor can rebroadcast the transactions and blocks it receives back to its
peers via `--broadcast-txs`, `--broadcast-tx-hashes`, `--broadcast-blocks`, and
`--broadcast-block-hashes`.

When block rebroadcasting is enabled (`--broadcast-blocks` or
`--broadcast-block-hashes`), the sensor validates the block signer before
rebroadcasting by default (`--validate-block-signer`, enabled by default): it
recovers the signer from the block header and only rebroadcasts blocks signed by
an address in the current Heimdall validator set. Set
`--validate-block-signer=false` to rebroadcast every block regardless of signer.

By default (`--cache-only-validated-blocks`), blocks from unknown signers are
still recorded to the database and their headers/bodies are still requested, but
they are not kept in the in-memory serving cache — so the sensor neither
rebroadcasts them nor serves them to peers on request, and they cannot evict
legitimate blocks from the cache. Set `--cache-only-validated-blocks=false` to
cache every block while still gating rebroadcast by signer.

The validator set is fetched from `--heimdall-url` at startup (the sensor aborts
if this initial fetch fails) and refreshed on the `--validator-set-refresh`
interval.

### Transaction Validation

When transaction rebroadcasting is enabled (`--broadcast-txs` or
`--broadcast-tx-hashes`), the sensor validates each transaction before
forwarding it (`--validate-broadcast-txs`, enabled by default). Without this,
any peer can push malformed or unmineable transactions through the sensor to
every other peer it is connected to.

Validation is the stateless half of a node's transaction admission rules, the
same checks `go-ethereum` applies before a transaction enters its pool:

- signature recovery against the `--network-id` chain ID, which rejects forged
  signatures and transactions signed for another chain
- transaction type (blob transactions are never forwarded, since the sensor has
  no sidecar to forward with them)
- encoded size, capped at 128KB
- gas below the intrinsic cost, or above the head block's gas limit
- init code size for contract creations
- oversized fee fields, and a tip cap above the fee cap

Nonce and balance are deliberately not checked. The sensor holds no chain state,
and those checks would cost an RPC round trip per sender.

The head block gates two of those checks (gas limit, base fee) and the head is
peer-supplied, so `UpdateHeadBlock` only accepts a block whose signer is in the
validator set when one is configured (`--validate-block-signer`). Without that,
a peer could declare a head with an absurd base fee and have every honest
transaction rejected. On a chain with no validator set the head is unguarded, as
it was before.

Fork rules are assumed current: every fork through Prague is treated as active,
since neither `--network-id` nor `--fork-id` yields a fork schedule. Later forks
mostly add transaction types, but two of them tighten — Shanghai caps init code
size and Prague adds the calldata floor gas cost — so on a chain that has not
adopted those, large deployments and calldata-heavy transactions are dropped as
`init_code_too_large` or `intrinsic_gas`. Turn validation off on such a chain.

The signer is bound to `--network-id`. On the networks this targets that is also
the chain ID; on a chain where the two differ, every transaction fails sender
recovery and nothing is forwarded. The sensor logs the chain ID it validates
against at startup.

Three fee floors are configurable on top of those rules:

| Flag                            | Default | Effect                                                            |
| ------------------------------- | ------- | ----------------------------------------------------------------- |
| `--broadcast-min-basefee-ratio` | `1.0`   | Drops transactions whose fee cap is under this fraction of the head block's base fee; `1.0` means "cannot be included right now", `0` disables |
| `--broadcast-min-gas-price`     | `0`     | Absolute floor in wei on the fee cap                               |
| `--broadcast-min-tip`           | `0`     | Absolute floor in wei on the tip cap, equivalent to a node's `--txpool.pricelimit` |

Validation gates rebroadcasting only. Rejected transactions are still cached,
served on request, and written to the database, so the sensor keeps a complete
record of the spam it declines to amplify. Transactions submitted to the
sensor's own `eth_sendRawTransaction` endpoint bypass these checks.

Two metrics track what is being dropped: `sensor_broadcast_txs_validated`
(labeled `result="accepted"|"rejected"`) gives the drop ratio, and
`sensor_broadcast_txs_rejected` breaks the rejections down by `reason`.
Transactions that fail to decode at all never reach validation and are counted
separately by `sensor_tx_decode_errors`.

## Examples

### Mainnet

To run a Polygon Mainnet sensor, copy the `genesis.json` from [here][mainnet-genesis].

```bash
polycli p2p sensor nodes.json \
  --bootnodes "enode://b8f1cc9c5d4403703fbf377116469667d2b1823c0daf16b7250aa576bacf399e42c3930ccfcb02c5df6879565a2b8931335565f0e8d3f8e72385ecf4a4bf160a@3.36.224.80:30303,enode://8729e0c825f3d9cad382555f3e46dcff21af323e89025a0e6312df541f4a9e73abfa562d64906f5e59c51fe6f0501b3e61b07979606c56329c020ed739910759@54.194.245.5:30303" \
  --network-id 137 \
  --sensor-id "sensor" \
  --write-blocks=true \
  --write-block-events=true \
  --write-txs=true \
  --write-tx-events=true \
  --genesis-hash "0xa9c28ce2141b56c474f1dc504bee9b01eb1bd7d1a507580d5519d4437a97de1b" \
  --fork-id "22d523b2" \
  --rpc "https://polygon-rpc.com" \
  --discovery-dns "enrtree://AKUEZKN7PSKVNR65FZDHECMKOJQSGPARGTPPBI7WS2VUL4EGR6XPC@pos.polygon-peers.io" \
  --pprof \
  --verbosity 700 \
  --pretty-logs=true \
  --database "json"
```

### Amoy

To run a Polygon Amoy sensor, copy the `genesis.json` from [here][amoy-genesis].

```bash
polycli p2p sensor amoy-nodes.json \
  --bootnodes "enode://b8f1cc9c5d4403703fbf377116469667d2b1823c0daf16b7250aa576bacf399e42c3930ccfcb02c5df6879565a2b8931335565f0e8d3f8e72385ecf4a4bf160a@3.36.224.80:30303,enode://8729e0c825f3d9cad382555f3e46dcff21af323e89025a0e6312df541f4a9e73abfa562d64906f5e59c51fe6f0501b3e61b07979606c56329c020ed739910759@54.194.245.5:30303" \
  --network-id 80002 \
  --sensor-id "sensor-amoy" \
  --write-blocks=true \
  --write-block-events=true \
  --write-txs=true \
  --write-tx-events=true \
  --genesis-hash "0x7202b2b53c5a0836e773e319d18922cc756dd67432f9a1f65352b61f4406c697" \
  --fork-id "8b7e4175" \
  --rpc "https://rpc-amoy.polygon.technology" \
  --discovery-dns "enrtree://AKUEZKN7PSKVNR65FZDHECMKOJQSGPARGTPPBI7WS2VUL4EGR6XPC@amoy.polygon-peers.io" \
  --pprof \
  --verbosity 700 \
  --pretty-logs=true \
  --database "json"
```

[mainnet-genesis]: https://github.com/0xPolygon/bor/blob/master/builder/files/genesis-mainnet-v1.json
[amoy-genesis]: https://github.com/0xPolygon/bor/blob/master/builder/files/genesis-amoy.json
[bootnodes]: https://docs.polygon.technology/pos/reference/seed-and-bootnodes/
