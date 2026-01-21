# `polycli ulxly zero-proof`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Create a proof that's filled with zeros.

```bash
polycli ulxly zero-proof [flags]
```

## Usage

This command prints a proof response filled with the zero hashes for
a Merkle tree. Zero hashes are the intermediate hashes computed when
all leaves below a node are empty (zero-valued).

These values are helpful for:

- **Debugging**: Understanding how populated a tree is and identifying
  which leaves and siblings are empty
- **Sanity checking**: Verifying if a hash in a proof response is a
  zero hash (indicating an empty subtree) or an actual intermediate
  hash from real data

Example usage:

```bash
polycli ulxly zero-proof
```

Example output:

```json
{
  "siblings": [
    "0x0000000000000000000000000000000000000000000000000000000000000000",
    "0xad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5",
    "0xb4c11951957c6f8f642c4af61cd6b24640fec6dc7fc607ee8206a99e92410d30",
    ...
  ]
}
```

Each sibling at index `i` represents the hash of an empty subtree at
depth `i`. These are computed by recursively hashing pairs of zero
hashes: `hash(zero[i-1], zero[i-1])`.

## Flags

```bash
  -h, --help   help for zero-proof
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
