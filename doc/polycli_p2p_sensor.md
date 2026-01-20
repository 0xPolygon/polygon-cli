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

The sensor can persist data to various backends including Google Cloud Datastore
or JSON output. If no nodes.json file exists at the specified path, it will be
created automatically.

The bootnodes may change, so refer to the [Polygon Knowledge Layer][bootnodes]
if the sensor is not discovering peers.

## Metrics

The sensor exposes Prometheus metrics at `http://localhost:2112/metrics`
(configurable via `--prom-port`). For a complete list of available metrics, see
[polycli_p2p_sensor_metrics.md](polycli_p2p_sensor_metrics.md).

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
      --api-port uint                    port API server will listen on (default 8080)
      --blocks-cache-ttl duration        time to live for block cache entries (0 for no expiration) (default 10m0s)
  -b, --bootnodes string                 comma separated nodes used for bootstrapping
      --broadcast-block-hashes           broadcast block hashes to peers
      --broadcast-blocks                 broadcast full blocks to peers
      --broadcast-tx                     broadcast full transactions to peers
      --broadcast-tx-hashes              broadcast transaction hashes to peers
      --database string                  which database to persist data to, options are:
                                           - datastore (GCP Datastore)
                                           - json (output to stdout)
                                           - none (no persistence) (default "none")
  -d, --database-id string               datastore database ID
      --dial-ratio int                   ratio of inbound to dialed connections (dial ratio of 2 allows 1/2 of connections to be dialed, setting to 0 defaults to 3)
      --discovery-dns string             DNS discovery ENR tree URL
      --discovery-port int               UDP P2P discovery port (default 30303)
      --fork-id bytesHex                 hex encoded fork ID (omit 0x) (default F097BC13)
      --genesis-hash string              genesis block hash (default "0xa9c28ce2141b56c474f1dc504bee9b01eb1bd7d1a507580d5519d4437a97de1b")
  -h, --help                             help for sensor
      --key string                       hex-encoded private key (cannot be set with --key-file)
  -k, --key-file string                  private key file (cannot be set with --key)
      --known-blocks-cache-ttl duration  time to live for known block cache entries (0 for no expiration) (default 5m0s)
      --known-txs-cache-ttl duration     time to live for known transaction cache entries (0 for no expiration) (default 5m0s)
      --max-blocks int                   maximum blocks to track across all peers (0 for no limit) (default 1024)
  -D, --max-db-concurrency int           maximum number of concurrent database operations to perform (increasing this
                                         will result in less chance of missing data but can significantly increase memory usage) (default 10000)
      --max-known-blocks int             maximum block hashes to track per peer (0 for no limit) (default 1024)
      --max-known-txs int                maximum transaction hashes to track per peer (0 for no limit) (default 8192)
      --max-parents int                  maximum parent block hashes to track per peer (0 for no limit) (default 1024)
  -m, --max-peers int                    maximum number of peers to connect to (default 2000)
      --max-requests int                 maximum request IDs to track per peer (0 for no limit) (default 2048)
      --max-txs int                      maximum transactions to cache for serving to peers (0 for no limit) (default 8192)
      --nat string                       NAT port mapping mechanism (any|none|upnp|pmp|pmp:<IP>|extip:<IP>) (default "any")
  -n, --network-id uint                  filter discovered nodes by this network ID
      --no-discovery                     disable P2P peer discovery
      --parents-cache-ttl duration       time to live for parent hash cache entries (0 for no expiration) (default 5m0s)
      --port int                         TCP network listening port (default 30303)
      --pprof                            run pprof server
      --pprof-port uint                  port pprof runs on (default 6060)
  -p, --project-id string                GCP project ID
      --prom                             run Prometheus server (default true)
      --prom-port uint                   port Prometheus runs on (default 2112)
      --requests-cache-ttl duration      time to live for requests cache entries (0 for no expiration) (default 5m0s)
      --rpc string                       RPC endpoint used to fetch latest block (default "https://polygon-rpc.com")
      --rpc-port uint                    port for JSON-RPC server to receive transactions (default 8545)
  -s, --sensor-id string                 sensor ID when writing block/tx events
      --static-nodes string              static nodes file
      --trusted-nodes string             trusted nodes file
      --ttl duration                     time to live (default 336h0m0s)
      --txs-cache-ttl duration           time to live for transaction cache entries (0 for no expiration) (default 10m0s)
      --write-block-events               write block events to database (default true)
  -B, --write-blocks                     write blocks to database (default true)
      --write-peers                      write peers to database (default true)
      --write-tx-events                  write transaction events to database (this option can significantly increase CPU and memory usage) (default true)
  -t, --write-txs                        write transactions to database (this option can significantly increase CPU and memory usage) (default true)
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
