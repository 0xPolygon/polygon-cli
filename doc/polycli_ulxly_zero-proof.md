# `polycli ulxly zero-proof`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

create a proof that's filled with zeros

```bash
polycli ulxly zero-proof [flags]
```

## Usage

Use this command to print a proof response that's filled with the zero
hashes. This values are very helpful for debugging because it would
tell you how populated the tree is and roughly which leaves and
siblings are empty. It's also helpful for sanity checking a proof
response to understand if the hashed value is part of the zero hashes
or if it's actually an intermediate hash.
## Flags

```bash
  -h, --help   help for zero-proof
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

- [polycli ulxly](polycli_ulxly.md) - Utilities for interacting with the uLxLy bridge
