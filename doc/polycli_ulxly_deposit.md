# `polycli ulxly deposit`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Make a uLxLy deposit transaction

```bash
polycli ulxly deposit [flags]
```

## Usage

Make a uLxLy deposit transaction
## Flags

```bash
      --bridge-address string        The address of the bridge contract.
      --destination-address string   The address of receiver in destination network.
      --destination-network uint32   The destination network number. (default 1)
      --forced                       The deposit transaction is forced. (default true)
      --gas-limit uint               The gas limit for the transaction. (default 300000)
  -h, --help                         help for deposit
      --metabytes string             Metabytes to append. (default "0x")
      --private-key string           The private key of the sender account.
      --rpc-url string               The RPC endpoint of the network (default "http://127.0.0.1:8545")
      --token-address string         The address of the token to send. (default "0x0000000000000000000000000000000000000000")
      --value int                    The amount to send.
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

- [polycli ulxly](polycli_ulxly.md) - Utilities for interacting with the lxly bridge
