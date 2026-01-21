# `polycli ulxly empty-proof`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Create an empty proof.

```bash
polycli ulxly empty-proof [flags]
```

## Usage

This command prints a proof response where all siblings are set to
the zero value (`0x0000000000000000000000000000000000000000000000000000000000000000`).

This is useful when you need to submit a dummy proof, for example in
testing scenarios or when a proof placeholder is required but the
actual proof data is not needed.

Note: This is different from `zero-proof`, which generates the actual
zero hashes of a Merkle tree (intermediate hashes computed from empty
subtrees). Empty proof simply fills all siblings with the zero value.

Example usage:

```bash
polycli ulxly empty-proof
```

Example output:

```json
{
  "siblings": [
    "0x0000000000000000000000000000000000000000000000000000000000000000",
    "0x0000000000000000000000000000000000000000000000000000000000000000",
    "0x0000000000000000000000000000000000000000000000000000000000000000",
    ...
  ]
}
```

## Flags

```bash
  -h, --help   help for empty-proof
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
