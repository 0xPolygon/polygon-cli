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

If no nodes.json file exists, run `echo "[]" >> nodes.json` to get started.
## Flags

```bash
  -b, --bootnodes string               Comma separated nodes used for bootstrapping. At least one bootnode is
                                       required, so other nodes in the network can discover each other.
  -d, --database string                Node database for updating and storing client information.
  -h, --help                           help for crawl
  -n, --network-id uint                Filter discovered nodes by this network id.
  -p, --parallel int                   How many parallel discoveries to attempt. (default 16)
  -r, --revalidation-interval string   The amount of time it takes to retry connecting to a failed peer. (default "10m")
  -t, --timeout string                 Time limit for the crawl. (default "30m0s")
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
