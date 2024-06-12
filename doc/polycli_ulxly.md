# `polycli ulxly`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Utilities for interacting with the lxly bridge

## Usage

These are low level tools for directly scanning bridge events and constructing proofs.
## Flags

```bash
  -h, --help   help for ulxly
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

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
- [polycli ulxly empty-proof](polycli_ulxly_empty-proof.md) - print an empty proof structure

- [polycli ulxly get-deposits](polycli_ulxly_get-deposits.md) - Get a range of deposits

- [polycli ulxly proof](polycli_ulxly_proof.md) - generate a merkle proof

- [polycli ulxly zero-proof](polycli_ulxly_zero-proof.md) - print a proof structure with the zero hashes

