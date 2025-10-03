# `polycli ulxly compute-nullifier-tree`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Compute the nullifier tree given the claims

```bash
polycli ulxly compute-nullifier-tree [flags]
```

## Usage

This command will attempt to computethe nullifierTree based on the claims that are provided.

Example usage:

```bash
polycli ulxly compute-nullifier-tree \
        --file-name claims-cardona-4880876-to-6028159.ndjson | jq '.'
```

In this case we are assuming we have a file
`claims-cardona-4880876-to-6028159.ndjson` that would have been generated
with a call to `polycli ulxly get-claims`. The output will be the
claims necessary to compute the nullifier tree.

This is the response from polycli:

```json
{
  "root": "0x4516ca2a793b8e20f56ec6ba8ca6033a672330670a3772f76f2ade9bc2125150"",
}
```

Note: more info https://github.com/BrianSeong99/Agglayer_PessimisticProof_Benchmark?tab=readme-ov-file#architecture-of-pessimistic-proof
## Flags

```bash
      --file-name string   ndjson file with events data
  -h, --help               help for compute-nullifier-tree
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

- [polycli ulxly](polycli_ulxly.md) - Utilities for interacting with the uLxLy bridge
