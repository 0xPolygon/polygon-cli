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
  -b, --bootnodes string               Comma separated nodes used for bootstrapping. At least one bootnode is
                                       required, so other nodes in the network can discover each other.
  -d, --database string                Node database for updating and storing client information
  -h, --help                           help for crawl
  -n, --network-id uint                Filter discovered nodes by this network id
  -u, --only-urls                      Only writes the enode URLs to the output (default true)
  -p, --parallel int                   How many parallel discoveries to attempt (default 16)
  -r, --revalidation-interval string   Time before retrying to connect to a failed peer (default "10m")
  -t, --timeout string                 Time limit for the crawl (default "30m0s")
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
