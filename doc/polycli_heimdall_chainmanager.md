# `polycli heimdall chainmanager`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Query chainmanager module endpoints.

## Usage

Chainmanager module queries (`x/chainmanager`) against a Heimdall v2 node.

The chainmanager module holds the L1 / L2 chain ids, the tx confirmation
depths, and the Ethereum contract addresses Heimdall uses to interact
with the root chain. Only one REST route exists upstream
(`/chainmanager/params`); `addresses` is a derived view for quick
copy-paste into a block explorer.

```bash
# Full params (chain ids + confirmations + contract addresses)
polycli heimdall chainmanager params

# Pluck a single field
polycli heimdall chainmanager params --field params.chain_params.root_chain_address

# Just the address map, one address per line, for etherscan copy-paste
polycli heimdall chainmanager addresses

# Alias
polycli heimdall cm params
```

Endpoints covered (confirmed from heimdall-v2 `proto/heimdallv2/chainmanager/query.proto`):

- `GET /chainmanager/params`

## Flags

```bash
  -h, --help   help for chainmanager
```

The command also inherits flags from parent commands.

```bash
      --amoy                     shortcut for --network amoy (default)
      --chain-id string          chain id used for signing
      --color string             color mode (auto|always|never) (default "auto")
      --config string            config file (default is $HOME/.polygon-cli.yaml)
      --curl                     print the equivalent curl command instead of executing
      --denom string             fee denom
      --heimdall-config string   path to heimdall config TOML (default ~/.polycli/heimdall.toml)
  -k, --insecure                 accept invalid TLS certs
      --json                     emit JSON instead of key/value
      --mainnet                  shortcut for --network mainnet
  -N, --network string           named network preset (amoy|mainnet)
      --no-color                 disable color output
      --pretty-logs              output logs in pretty format instead of JSON (default true)
      --raw                      preserve raw bytes (no 0x-hex normalization)
  -r, --rest-url string          heimdall REST gateway URL
      --rpc-headers string       extra request headers, comma-separated key=value pairs
  -R, --rpc-url string           cometBFT RPC URL
      --timeout int              HTTP timeout in seconds
  -v, --verbosity string         log level (string or int):
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

- [polycli heimdall](polycli_heimdall.md) - Query and interact with a Heimdall v2 node.
- [polycli heimdall chainmanager addresses](polycli_heimdall_chainmanager_addresses.md) - Print L1 contract addresses and chain ids.

- [polycli heimdall chainmanager params](polycli_heimdall_chainmanager_params.md) - Fetch the chainmanager module parameters.

