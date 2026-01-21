# `polycli ulxly claim-everything`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Attempt to claim as many deposits and messages as possible.

```bash
polycli ulxly claim-everything [flags]
```

## Flags

```bash
      --bridge-address string              address of the lxly bridge
      --bridge-limit int                   limit the number or responses returned by the bridge service when claiming (default 25)
      --bridge-offset int                  offset to specify for pagination of underlying bridge service deposits
      --bridge-service-map strings         network ID to bridge service URL mappings (e.g. '1=http://network-1-bridgeurl,7=http://network-2-bridgeurl')
      --chain-id string                    chain ID to use in the transaction
      --concurrency uint                   worker pool size for claims (default 1)
      --destination-address string         destination address for the bridge
      --dry-run                            do all of the transaction steps but do not send the transaction
      --gas-limit uint                     force specific gas limit for transaction
      --gas-price string                   gas price to use
  -h, --help                               help for claim-everything
      --insecure                           skip TLS certificate verification
      --legacy                             force usage of legacy bridge service (default true)
      --private-key string                 hex encoded private key for sending transaction
      --rpc-url string                     RPC URL to send the transaction
      --transaction-receipt-timeout uint   timeout in seconds to wait for transaction receipt confirmation (default 60)
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

- [polycli ulxly](polycli_ulxly.md) - Utilities for interacting with the uLxLy bridge.
