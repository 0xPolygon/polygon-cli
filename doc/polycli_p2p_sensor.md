# `polycli p2p sensor`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Start a devp2p sensor that discovers other peers and will receive blocks and transactions.

```bash
polycli p2p sensor [nodes file] [flags]
```

## Usage

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

### Transaction rebroadcast gates

Rebroadcasting a transaction is only useful if that transaction could still be
mined. Echoing one that cannot — a replayed historical transaction, or one
priced below what the chain will include — amplifies junk across the network:
the sensor fetches the body, then announces it to every peer. Because a sensor
carries far more peers than an ordinary node, that amplification is large.

Three gates apply when transaction rebroadcasting is on (`--broadcast-txs` or
`--broadcast-tx-hashes`). They run cheapest-first and only ever withhold an
*echo*: everything the sensor sees is still recorded to the database and still
served to peers that ask for it, the same split `--validate-block-signer` uses
for blocks.

- **Tip floor** (`--rebroadcast-min-tip`, disabled by default): withholds
  transactions offering less than the given tip, which accepts units — e.g.
  `--rebroadcast-min-tip=25gwei`. Bor will not include anything below 25 gwei on
  Polygon, so those transactions are unminable no matter what else is true. The
  check needs no sender recovery and no lookup.
- **Stale nonce** (`--gate-stale-txs`, enabled by default): withholds
  transactions whose nonce is below their sender's next usable nonce. The sensor
  tracks sender nonces from the transactions of every block it already observes,
  so the common case costs nothing extra. Senders that have not appeared in an
  observed block are looked up against `--rpc` asynchronously
  (`--gate-stale-txs-rpc`, enabled by default); the transaction that triggers a
  lookup is still rebroadcast, and the result gates later ones from that sender.
  `--max-account-nonces` and `--account-nonces-ttl` size that map. Only blocks
  that pass the signer check feed it, so a peer that is not a known validator
  cannot poison a sender's nonce with a fabricated block and get that sender's
  real transactions withheld.
- **Rate cap** (`--rebroadcast-rate-limit`, disabled by default): a token bucket
  over rebroadcast transactions per second, with `--rebroadcast-burst` as the
  depth. It bounds worst-case amplification even when the gates above miss
  something. It runs last, so transactions the other gates rejected do not spend
  the budget.

To size the problem before enforcing anything, run with
`--rebroadcast-gate-log-only`: every gate is evaluated and its
`sensor_rebroadcast_filtered` counters move, but nothing is actually withheld.

Withheld transactions are counted by
`sensor_rebroadcast_filtered{reason="stale_nonce"|"low_tip"|"rate_limited"}` and
passing ones by `sensor_rebroadcast_allowed`. The nonce map is reported by
`sensor_rebroadcast_known_senders`, and fallback lookups by
`sensor_rebroadcast_nonce_lookups{result="ok"|"error"|"dropped"}`.

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

## Flags

```bash
      --account-nonces-ttl duration       time to live for tracked sender nonces (0 for no expiration) (default 1h0m0s)
      --api-port uint                     port API server will listen on (default 8080)
      --blocks-cache-ttl duration         time to live for block cache entries (0 for no expiration) (default 10m0s)
  -b, --bootnodes string                  comma separated nodes used for bootstrapping
      --broadcast-block-hashes            broadcast block hashes to peers
      --broadcast-blocks                  broadcast full blocks to peers
      --broadcast-tx-hashes               broadcast transaction hashes to peers
      --broadcast-txs                     broadcast full transactions to peers
      --broadcast-workers int             number of concurrent broadcast workers (default 4)
      --cache-only-validated-blocks       only cache and serve blocks signed by a known validator (unknown-signer blocks are still recorded to the database); has no effect without --validate-block-signer (default true)
      --clickhouse-dsn string             ClickHouse DSN, e.g. clickhouse://user:pass@host:9000/sensor (used with --database=clickhouse)
      --database string                   which database to persist data to, options are:
                                            - datastore (GCP Datastore)
                                            - clickhouse (ClickHouse, see --clickhouse-dsn)
                                            - json (output to stdout)
                                            - none (no persistence) (default "none")
  -d, --database-id string                datastore database ID
      --dial-ratio int                    ratio of inbound to dialed connections (dial ratio of 2 allows 1/2 of connections to be dialed, setting to 0 defaults to 3)
      --discovery-dns string              DNS discovery ENR tree URL
      --discovery-port int                UDP P2P discovery port (default 30303)
      --fork-id bytesHex                  hex encoded fork ID (omit 0x) (default 22D523B2)
      --gate-stale-txs                    only rebroadcast transactions whose nonce is at or above the sender's known
                                          next nonce; stale ones can never be mined, so echoing them only amplifies replay
                                          traffic (they are still recorded to the database) (default true)
      --gate-stale-txs-rpc                look up nonces from --rpc for senders not yet seen in a block; lookups are
                                          asynchronous and the transaction that triggers one is still rebroadcast (default true)
      --genesis-hash string               genesis block hash (default "0xa9c28ce2141b56c474f1dc504bee9b01eb1bd7d1a507580d5519d4437a97de1b")
      --heimdall-url string               heimdall REST URL for the validator set (used to validate blocks before rebroadcast) (default "https://heimdall-api.polygon.technology")
  -h, --help                              help for sensor
      --key string                        hex-encoded private key (cannot be set with --key-file)
  -k, --key-file string                   private key file (cannot be set with --key)
      --known-txs-bloom-hashes uint       number of hash functions for known txs bloom filter (default 7)
      --known-txs-bloom-size uint         bloom filter size in bits for tracking known transactions per peer (default ~40KB per filter,
                                          optimized for ~32K elements with ~1% false positive rate) (default 327680)
      --max-account-nonces int            maximum sender nonces to track for --gate-stale-txs (0 for no limit) (default 65536)
      --max-blocks int                    maximum blocks to track across all peers (0 for no limit) (default 1024)
  -D, --max-db-concurrency int            maximum number of concurrent database operations to perform (increasing this
                                          will result in less chance of missing data but can significantly increase memory usage) (default 10000)
      --max-known-blocks int              maximum block hashes to track per peer (0 for no limit) (default 1024)
      --max-parents int                   maximum parent block hashes to track per peer (0 for no limit) (default 1024)
  -m, --max-peers int                     maximum number of peers to connect to (default 2000)
      --max-queued-txs int                maximum transaction announcements to queue per peer (default 4096)
      --max-requests int                  maximum request IDs to track per peer (0 for no limit) (default 2048)
      --max-tx-packet-size int            target size in bytes for transaction broadcast packets (default 102400)
      --max-txs int                       maximum transactions to cache for serving to peers (0 for no limit) (default 32768)
      --nat string                        NAT port mapping mechanism (any|none|upnp|pmp|pmp:<IP>|extip:<IP>) (default "any")
  -n, --network-id uint                   filter discovered nodes by this network ID
      --no-discovery                      disable P2P peer discovery
      --parents-cache-ttl duration        time to live for parent hash cache entries (0 for no expiration) (default 5m0s)
      --peer-snapshot-interval duration   how often to persist the connected-peer set (requires --write-peers); lower
                                          values multiply write volume by up to --max-peers rows per tick (default 30s)
      --port int                          TCP network listening port (default 30303)
      --pprof                             run pprof server
      --pprof-port uint                   port pprof runs on (default 6060)
  -p, --project-id string                 GCP project ID
      --prom                              run Prometheus server (default true)
      --prom-port uint                    port Prometheus runs on (default 2112)
      --proxy-rpc                         proxy unsupported RPC methods to the --rpc endpoint
      --proxy-rpc-timeout duration        timeout for proxied RPC requests (default 30s)
      --rebroadcast-burst int             token bucket depth for --rebroadcast-rate-limit (defaults to one second of the rate)
      --rebroadcast-gate-log-only         evaluate the rebroadcast gates and count what they would drop, but rebroadcast
                                          everything anyway; use it to size the problem before enforcing
      --rebroadcast-min-tip gas           withhold transactions offering less than this tip from rebroadcast, with unit
                                          support (e.g. "25gwei"); bor will not include anything below 25gwei on Polygon, so
                                          those are unminable regardless (0 disables)
      --rebroadcast-rate-limit float      cap rebroadcast throughput in transactions per second, a backstop that bounds
                                          amplification even when the gates above miss (0 for no limit)
      --requests-cache-ttl duration       time to live for requests cache entries (0 for no expiration) (default 5m0s)
      --rpc string                        RPC endpoint used to fetch latest block (default "https://polygon-rpc.com")
      --rpc-port uint                     port for JSON-RPC server to receive transactions (default 8545)
  -s, --sensor-id string                  sensor ID when writing block/tx events
      --static-nodes string               static nodes file
      --trusted-nodes string              trusted nodes file
      --ttl duration                      time to live (default 336h0m0s)
      --tx-batch-timeout duration         timeout for batching transactions before broadcast (default 500ms)
      --tx-broadcast-queue-size int       capacity of transaction broadcast queue (default 100000)
      --txs-cache-ttl duration            time to live for transaction cache entries (0 for no expiration) (default 10m0s)
      --validate-block-signer             only rebroadcast blocks signed by a validator in the heimdall validator set (default true)
      --validator-set-refresh duration    interval to refresh the validator set from heimdall (default 5m0s)
      --write-block-events                write block events to database (default true)
  -B, --write-blocks                      write blocks to database (default true)
      --write-first-block-event           write one block event on first-seen only; ignored when --write-block-events is set
      --write-first-tx-event              write one transaction event on first-seen only; ignored when --write-tx-events is set
      --write-peers                       write peers to database (default true)
      --write-tx-events                   write transaction events to database (this option can significantly increase CPU and memory usage) (default true)
  -t, --write-txs                         write transactions to database (this option can significantly increase CPU and memory usage) (default true)
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

- [polycli p2p](polycli_p2p.md) - Set of commands related to devp2p.
