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
- [polycli ulxly bridge](polycli_ulxly_bridge.md) - commands for making deposits to the uLxLy bridge

- [polycli ulxly claim](polycli_ulxly_claim.md) - commands for making claims of deposits from the uLxLy bridge

- [polycli ulxly claim-everything](polycli_ulxly_claim-everything.md) - attempt to claim any unclaimed deposits

- [polycli ulxly empty-proof](polycli_ulxly_empty-proof.md) - create an empty proof

- [polycli ulxly get-deposits](polycli_ulxly_get-deposits.md) - generate ndjson for each bridge deposit over a particular range of blocks

- [polycli ulxly proof](polycli_ulxly_proof.md) - generate a proof for a given range of deposits

- [polycli ulxly zero-proof](polycli_ulxly_zero-proof.md) - create a proof that's filled with zeros

