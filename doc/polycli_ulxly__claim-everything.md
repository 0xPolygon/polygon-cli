# `polycli ulxly  claim-everything`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

attempt to claim any unclaimed deposits

```bash
polycli ulxly  claim-everything [flags]
```

## Flags

```bash
      --bridge-limit int             Limit the number or responses returned by the bridge service when claiming (default 25)
      --bridge-offset int            The offset to specify for pagination of the underlying bridge service deposits
      --bridge-service-map strings   Mappings between network ids and bridge service urls. E.g. '1=http://network-1-bridgeurl,7=http://network-2-bridgeurl'
  -h, --help                         help for claim-everything
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
