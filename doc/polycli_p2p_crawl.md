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

To crawl the network for nodes and write the output json to a file. This will
not engage in block or transaction propagation, but it can give a good indicator
of network size, and the output json can be used to quick start other nodes.

## Example

```bash
polycli p2p crawl nodes.json \
  --bootnodes "enode://0cb82b395094ee4a2915e9714894627de9ed8498fb881cec6db7c65e8b9a5bd7f2f25cc84e71e89d0947e51c76e85d0847de848c7782b13c0255247a6758178c@44.232.55.71:30303,enode://88116f4295f5a31538ae409e4d44ad40d22e44ee9342869e7d68bdec55b0f83c1530355ce8b41fbec0928a7d75a5745d528450d30aec92066ab6ba1ee351d710@159.203.9.164:30303,enode://4be7248c3a12c5f95d4ef5fff37f7c44ad1072fdb59701b2e5987c5f3846ef448ce7eabc941c5575b13db0fb016552c1fa5cca0dda1a8008cf6d63874c0f3eb7@3.93.224.197:30303,enode://32dd20eaf75513cf84ffc9940972ab17a62e88ea753b0780ea5eca9f40f9254064dacb99508337043d944c2a41b561a17deaad45c53ea0be02663e55e6a302b2@3.212.183.151:30303" \
  --network-id 137
```

## Flags

```bash
  -b, --bootnodes string               comma separated nodes used for bootstrapping. At least one bootnode is
                                       required, so other nodes in the network can discover each other
  -d, --database string                node database for updating and storing client information
      --discovery-dns string           enable EIP-1459, DNS Discovery to recover node list from given ENRTree
  -h, --help                           help for crawl
  -n, --network-id uint                filter discovered nodes by this network ID
  -u, --only-urls                      only writes enode URLs to output (default true)
  -p, --parallel int                   how many parallel discoveries to attempt (default 16)
  -r, --revalidation-interval string   time before retrying to connect to a failed peer (default "10m")
  -t, --timeout string                 time limit for the crawl (default "30m0s")
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
