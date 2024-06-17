# `polycli ulxly deposit-claim`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Make a uLxLy claim transaction

```bash
polycli ulxly deposit-claim [flags]
```

## Usage

Make a uLxLy claim transaction
## Flags

```bash
      --bridge-address string                The address of the bridge contract.
      --bridge-service-url string            The RPC endpoint of the bridge service component.
      --chain-id string                      The chainID.
      --claim-address string                 The address that is receiving the bridged asset.
      --claim-index string                   The deposit count, or index to initiate a claim transaction for. (default "0")
      --gas-limit uint                       The gas limit for the transaction. (default 300000)
  -h, --help                                 help for deposit-claim
      --origin-network string                The network ID of the origin network. (default "0")
      --private-key string                   The private key of the sender account.
      --rpc-url string                       The RPC endpoint of the destination network (default "http://127.0.0.1:8545")
      --transaction-receipt-timeout uint32   The timeout limit to check for the transaction receipt of the claim. (default 60)
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
