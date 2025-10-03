# `polycli ulxly empty-proof`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

create an empty proof

```bash
polycli ulxly empty-proof [flags]
```

## Usage

Use this command to print an empty proof response that's filled with zero-valued siblings like 0x0000000000000000000000000000000000000000000000000000000000000000. This can be useful when you need to submit a dummy proof.
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

- [polycli ulxly](polycli_ulxly.md) - Utilities for interacting with the uLxLy bridge
