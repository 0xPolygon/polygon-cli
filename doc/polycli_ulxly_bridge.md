# `polycli ulxly  bridge`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Commands for moving funds and sending messages from one chain to another

## Flags

```bash
      --call-data string             call data to be passed directly with bridge-message or as an ERC20 Permit (default "0x")
      --call-data-file string        a file containing hex encoded call data
      --destination-network uint32   rollup ID of the destination network
      --force-update-root            update the new global exit root (default true)
  -h, --help                         help for bridge
      --token-address string         address of ERC20 token to use (default "0x0000000000000000000000000000000000000000")
      --value string                 amount in wei to send with the transaction (default "0")
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

- [polycli ulxly bridge asset](polycli_ulxly_bridge_asset.md) - Move ETH or an ERC20 between to chains

- [polycli ulxly bridge message](polycli_ulxly_bridge_message.md) - Send some ETH along with data from one chain to another chain

- [polycli ulxly bridge weth](polycli_ulxly_bridge_weth.md) - For L2's that use a gas token, use this to transfer WETH to another chain

