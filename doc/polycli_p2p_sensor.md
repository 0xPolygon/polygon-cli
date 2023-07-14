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

If no nodes.json file exists, run `echo "{}" >> nodes.json` to get started.
## Flags

```bash
  -b, --bootnodes string               Comma separated nodes used for bootstrapping. At least one bootnode is
                                       required, so other nodes in the network can discover each other.
  -d, --database string                Node database for updating and storing client information.
  -h, --help                           help for sensor
  -D, --max-db-writes int              The maximum number of concurrent database writes to perform. Increasing
                                       this will result in less chance of missing data (i.e. broken pipes) but
                                       can significantly increase memory usage. (default 100)
  -m, --max-peers int                  Maximum number of peers to connect to. (default 200)
  -n, --network-id uint                Filter discovered nodes by this network ID.
  -p, --parallel int                   How many parallel discoveries to attempt. (default 16)
      --pprof                          Whether to run pprof.
      --pprof-port uint                The port to run pprof on. (default 6060)
  -P, --project-id string              GCP project ID.
  -r, --revalidation-interval string   The amount of time it takes to retry connecting to a failed peer. (default "10m")
  -s, --sensor-id string               Sensor ID.
      --write-block-events             Whether to write block events to the database. (default true)
  -B, --write-blocks                   Whether to write blocks to the database. (default true)
      --write-tx-events                Whether to write transaction events to the database. This option could significantly
                                       increase CPU and memory usage. (default true)
  -t, --write-txs                      Whether to write transactions to the database. This option could significantly
                                       increase CPU and memory usage. (default true)
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Fatal
                        200 Error
                        300 Warning
                        400 Info
                        500 Debug
                        600 Trace (default 400)
```

## See also

- [polycli p2p](polycli_p2p.md) - Set of commands related to devp2p.
