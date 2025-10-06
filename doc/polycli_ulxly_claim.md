# `polycli ulxly  claim`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Commands for claiming deposits on a particular chain.

## Flags

```bash
      --bridge-service-url string   URL of the bridge service
      --deposit-count uint32        deposit count of the bridge transaction
      --deposit-network uint32      rollup ID of the network where the deposit was made
      --global-index string         an override of the global index value
  -h, --help                        help for claim
      --proof-ger string            if specified and using legacy mode, the proof will be generated against this GER
      --wait duration               retry claiming until deposit is ready, up to specified duration (available for claim asset and claim message)
```

The command also inherits flags from parent commands.

```bash
      --bridge-address string              address of the lxly bridge
      --chain-id string                    chain ID to use in the transaction
      --config string                      config file (default is $HOME/.polygon-cli.yaml)
      --destination-address string         destination address for the bridge
      --dry-run                            do all of the transaction steps but do not send the transaction
      --gas-limit uint                     force specific gas limit for transaction
      --gas-price string                   gas price to use
      --insecure                           skip TLS certificate verification
      --legacy                             force usage of legacy bridge service (default true)
      --pretty-logs                        output logs in pretty format instead of JSON (default true)
      --private-key string                 hex encoded private key for sending transaction
      --rpc-url string                     RPC URL to send the transaction
      --transaction-receipt-timeout uint   timeout in seconds to wait for transaction receipt confirmation (default 60)
  -v, --verbosity string                   log level (string or int):
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

- [polycli ulxly claim asset](polycli_ulxly_claim_asset.md) - Claim a deposit.

- [polycli ulxly claim message](polycli_ulxly_claim_message.md) - Claim a message.

