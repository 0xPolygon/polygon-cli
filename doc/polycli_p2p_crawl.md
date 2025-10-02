# `polycli p2p crawl`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Crawl a network on the devp2p layer and generate a nodes JSON file.

```bash
polycli p2p crawl [nodes file] [flags]
```

## Usage

If no nodes.json file exists, it will be created.
## Flags

```bash
  -b, --bootnodes string               comma separated nodes used for bootstrapping. At least one bootnode is
                                       required, so other nodes in the network can discover each other
  -d, --database string                node database for updating and storing client information
      --discovery-dns string           enable EIP-1459, DNS Discovery to recover node list from given ENRTree
  -h, --help                           help for crawl
  -n, --network-id uint                filter discovered nodes by this network ID
  -u, --only-urls                      only writes the enode URLs to the output (default true)
  -p, --parallel int                   how many parallel discoveries to attempt (default 16)
  -r, --revalidation-interval string   time before retrying to connect to a failed peer (default "10m")
  -t, --timeout string                 time limit for the crawl (default "30m0s")
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
