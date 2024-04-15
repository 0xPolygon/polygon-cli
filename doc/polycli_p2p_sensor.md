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
  -b, --bootnodes string         Comma separated nodes used for bootstrapping
  -d, --database-id string       Datastore database ID
      --dial-ratio int           Ratio of inbound to dialed connections. A dial ratio of 2 allows 1/2 of
                                 connections to be dialed. Setting this to 0 defaults it to 3.
      --discovery-port int       UDP P2P discovery port (default 30303)
      --fork-id bytesHex         The hex encoded fork id (omit the 0x) (default F097BC13)
      --genesis-hash string      The genesis block hash (default "0xa9c28ce2141b56c474f1dc504bee9b01eb1bd7d1a507580d5519d4437a97de1b")
  -h, --help                     help for sensor
  -k, --key-file string          Private key file
  -D, --max-db-concurrency int   Maximum number of concurrent database operations to perform. Increasing this
                                 will result in less chance of missing data (i.e. broken pipes) but can
                                 significantly increase memory usage. (default 10000)
  -m, --max-peers int            Maximum number of peers to connect to (default 200)
      --nat string               NAT port mapping mechanism (any|none|upnp|pmp|pmp:<IP>|extip:<IP>) (default "any")
  -n, --network-id uint          Filter discovered nodes by this network ID
      --port int                 TCP network listening port (default 30303)
      --pprof                    Whether to run pprof
      --pprof-port uint          Port pprof runs on (default 6060)
  -p, --project-id string        GCP project ID
      --prom                     Whether to run Prometheus (default true)
      --prom-port uint           Port Prometheus runs on (default 2112)
      --quick-start              Whether to load the nodes.json as static nodes to quickly start the network.
                                 This produces faster development cycles but can prevent the sensor from being to
                                 connect to new peers if the nodes.json file is large.
      --rpc string               RPC endpoint used to fetch the latest block (default "https://polygon-rpc.com")
  -s, --sensor-id string         Sensor ID when writing block/tx events
      --trusted-nodes string     Trusted nodes file
      --ttl duration             Time to live (default 336h0m0s)
      --write-block-events       Whether to write block events to the database (default true)
  -B, --write-blocks             Whether to write blocks to the database (default true)
      --write-tx-events          Whether to write transaction events to the database. This option could
                                 significantly increase CPU and memory usage. (default true)
  -t, --write-txs                Whether to write transactions to the database. This option could significantly
                                 increase CPU and memory usage. (default true)
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

- [polycli p2p](polycli_p2p.md) - Set of commands related to devp2p.
