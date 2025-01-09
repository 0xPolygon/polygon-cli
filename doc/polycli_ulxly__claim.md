# `polycli ulxly  claim`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Commands for claiming deposits on a particular chain

## Flags

```bash
      --bridge-service-url string   the URL of the bridge service
      --deposit-count uint          the deposit count of the bridge transaction
      --deposit-network uint        the rollup id of the network where the deposit was initially made
      --global-index string         an override of the global index value
  -h, --help                        help for claim
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
- [polycli ulxly  claim asset](polycli_ulxly__claim_asset.md) - Claim a deposit

- [polycli ulxly  claim message](polycli_ulxly__claim_message.md) - Claim a message

