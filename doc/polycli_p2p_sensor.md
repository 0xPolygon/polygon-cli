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

If no nodes.json file exists, it will be created.
## Flags

```bash
      --api-port uint            port the API server will listen on (default 8080)
  -b, --bootnodes string         comma separated nodes used for bootstrapping
      --database string          which database to persist data to, options are:
                                   - datastore (GCP Datastore)
                                   - json (output to stdout)
                                   - none (no persistence) (default "none")
  -d, --database-id string       datastore database ID
      --dial-ratio int           ratio of inbound to dialed connections. A dial ratio of 2 allows 1/2 of
                                 connections to be dialed. Setting this to 0 defaults it to 3
      --discovery-dns string     DNS discovery ENR tree URL
      --discovery-port int       UDP P2P discovery port (default 30303)
      --fork-id bytesHex         the hex encoded fork ID (omit the 0x) (default F097BC13)
      --genesis-hash string      the genesis block hash (default "0xa9c28ce2141b56c474f1dc504bee9b01eb1bd7d1a507580d5519d4437a97de1b")
  -h, --help                     help for sensor
      --key string               hex-encoded private key (cannot be set with --key-file)
  -k, --key-file string          private key file (cannot be set with --key)
  -D, --max-db-concurrency int   maximum number of concurrent database operations to perform. Increasing this
                                 will result in less chance of missing data (i.e. broken pipes) but can
                                 significantly increase memory usage (default 10000)
  -m, --max-peers int            maximum number of peers to connect to (default 2000)
      --nat string               NAT port mapping mechanism (any|none|upnp|pmp|pmp:<IP>|extip:<IP>) (default "any")
  -n, --network-id uint          filter discovered nodes by this network ID
      --port int                 TCP network listening port (default 30303)
      --pprof                    run pprof server
      --pprof-port uint          port pprof runs on (default 6060)
  -p, --project-id string        GCP project ID
      --prom                     run Prometheus server (default true)
      --prom-port uint           port Prometheus runs on (default 2112)
      --quick-start              whether to load the nodes.json as static nodes to quickly start the network.
                                 This produces faster development cycles but can prevent the sensor from being to
                                 connect to new peers if the nodes.json file is large
      --rpc string               RPC endpoint used to fetch the latest block (default "https://polygon-rpc.com")
  -s, --sensor-id string         sensor ID when writing block/tx events
      --trusted-nodes string     trusted nodes file
      --ttl duration             time to live (default 336h0m0s)
      --write-block-events       write block events to the database (default true)
  -B, --write-blocks             write blocks to the database (default true)
      --write-peers              write peers to the database (default true)
      --write-tx-events          write transaction events to the database (can significantly increase CPU and memory usage) (default true)
  -t, --write-txs                write transactions to the database (can significantly increase CPU and memory usage) (default true)
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     output logs in pretty format instead of JSON (default true)
  -v, --verbosity int   0 - silent
                        100 panic
                        200 fatal
                        300 error
                        400 warning
                        500 info
                        600 debug
                        700 trace (default 500)
```

## See also

- [polycli p2p](polycli_p2p.md) - Set of commands related to devp2p.
