# `polycli ulxly`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Utilities for interacting with the uLxLy bridge

## Usage

Basic utility commands for interacting with the bridge contracts, bridge services, and generating proofs
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
- [polycli ulxly bridge](polycli_ulxly_bridge.md) - Commands for moving funds and sending messages from one chain to another

- [polycli ulxly claim](polycli_ulxly_claim.md) - Commands for claiming deposits on a particular chain

- [polycli ulxly claim-everything](polycli_ulxly_claim-everything.md) - Attempt to claim as many deposits and messages as possible

- [polycli ulxly compute-balance-nullifier-tree](polycli_ulxly_compute-balance-nullifier-tree.md) - Compute the balance tree and the nullifier tree given the deposits and claims

- [polycli ulxly compute-balance-tree](polycli_ulxly_compute-balance-tree.md) - Compute the balance tree given the deposits

- [polycli ulxly compute-nullifier-tree](polycli_ulxly_compute-nullifier-tree.md) - Compute the nullifier tree given the claims

- [polycli ulxly empty-proof](polycli_ulxly_empty-proof.md) - create an empty proof

- [polycli ulxly get-claims](polycli_ulxly_get-claims.md) - Generate ndjson for each bridge claim over a particular range of blocks

- [polycli ulxly get-deposits](polycli_ulxly_get-deposits.md) - Generate ndjson for each bridge deposit over a particular range of blocks

- [polycli ulxly get-verify-batches](polycli_ulxly_get-verify-batches.md) - Generate ndjson for each verify batch over a particular range of blocks

- [polycli ulxly proof](polycli_ulxly_proof.md) - Generate a proof for a given range of deposits

- [polycli ulxly proof](polycli_ulxly_proof.md) - Generate a proof for a given range of deposits

- [polycli ulxly rollups-proof](polycli_ulxly_rollups-proof.md) - Generate a proof for a given range of rollups

- [polycli ulxly zero-proof](polycli_ulxly_zero-proof.md) - create a proof that's filled with zeros

