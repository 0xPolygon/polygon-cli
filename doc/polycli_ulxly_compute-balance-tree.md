# `polycli ulxly compute-balance-tree`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Compute the balance tree given the deposits.

```bash
polycli ulxly compute-balance-tree [flags]
```

## Usage

This command will attempt to compute the root of the balnace tree based on the bridge
events that are provided.

Example usage:

```bash
polycli ulxly compute-balance-tree \
        --l2-claims-file l2-claims-0-to-11454081.ndjson \
        --l2-deposits-file l2-deposits-0-to-11454081.ndjson \
        --l2-network-id 3
        --bridge-address 0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582 \
        --rpc-url http://localhost:8213 | jq '.'
```

In this case we are assuming we have two files
`l2-claims-0-to-11454081.ndjson` and `l2-deposits-0-to-11454081.ndjson` that would have been generated
with a call to `polycli ulxly get-deposits` and `polycli ulxly get-claims` pointing to each network. The output will be the
root of the tree for the provided deposits and claims.

This is the response from polycli:

```json
{
  "root": "0x4516ca2a793b8e20f56ec6ba8ca6033a672330670a3772f76f2ade9bc2125150"",
}
```

Note: more info https://github.com/BrianSeong99/Agglayer_PessimisticProof_Benchmark?tab=readme-ov-file#architecture-of-pessimistic-proof
## Flags

```bash
      --bridge-address string     bridge address
  -h, --help                      help for compute-balance-tree
      --insecure                  skip TLS certificate verification
      --l2-claims-file string     ndjson file with l2 claim events data
      --l2-deposits-file string   ndjson file with l2 deposit events data
      --l2-network-id uint32      L2 network ID
  -r, --rpc-url string            RPC URL
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

- [polycli ulxly](polycli_ulxly.md) - Utilities for interacting with the uLxLy bridge.
