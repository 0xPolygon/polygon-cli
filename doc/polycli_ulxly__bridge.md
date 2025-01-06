# `polycli ulxly  bridge`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

commands for making deposits to the uLxLy bridge

## Flags

```bash
      --call-data string             call data to be passed directly with bridge-message or as an ERC20 Permit (default "0x")
      --destination-network uint32   the rollup id of the destination network
      --force-update-root            indicates if the new global exit root is updated or not (default true)
  -h, --help                         help for bridge
      --token-address string         the address of an ERC20 token to be used (default "0x0000000000000000000000000000000000000000")
      --value string                 the amount in wei to be sent along with the transaction
```

The command also inherits flags from parent commands.

```bash
      --bridge-address string              the address of the lxly bridge
      --chain-id string                    set the chain id to be used in the transaction
      --config string                      config file (default is $HOME/.polygon-cli.yaml)
      --destination-address string         the address where the bridge will be sent to
      --dry-run                            do all of the transaction steps but do not send the transaction
      --gas-limit uint                     force a gas limit when sending a transaction
      --gas-price string                   the gas price to be used
      --pretty-logs                        Should logs be in pretty format or JSON (default true)
      --private-key string                 the hex encoded private key to be used when sending the tx
      --rpc-url string                     the URL of the RPC to send the transaction
      --transaction-receipt-timeout uint   the amount of time to wait while trying to confirm a transaction receipt (default 60)
  -v, --verbosity int                      0 - Silent
                                           100 Panic
                                           200 Fatal
                                           300 Error
                                           400 Warning
                                           500 Info
                                           600 Debug
                                           700 Trace (default 500)
```

## See also

- [polycli ulxly ](polycli_ulxly_.md) - 
- [polycli ulxly  bridge asset](polycli_ulxly__bridge_asset.md) - send a single deposit of value or an ERC20 into the bridge

- [polycli ulxly  bridge message](polycli_ulxly__bridge_message.md) - send some value along with call data into the bridge

- [polycli ulxly  bridge weth](polycli_ulxly__bridge_weth.md) - send some WETH into the bridge

