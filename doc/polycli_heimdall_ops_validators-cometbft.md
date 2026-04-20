# `polycli heimdall ops validators-cometbft`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

List CometBFT consensus validators (NOT Heimdall x/stake).

```bash
polycli heimdall ops validators-cometbft [HEIGHT] [flags]
```

## Usage

List validators from CometBFT's /validators endpoint at a given
height (default latest). Output is the consensus layer's view: the
20-byte consensus address, the validator's Secp256k1-eth pubkey,
voting power, and proposer priority.

This is distinct from the Heimdall x/stake validator set. Use
'polycli heimdall validator' for staking info (operator address, moniker,
jailed status, etc).
## Flags

```bash
  -f, --field stringArray   pluck one or more fields (repeatable)
  -h, --help                help for validators-cometbft
      --page int            page number (1-indexed) (default 1)
      --per-page int        validators per page (CometBFT default is 30, max 100) (default 100)
      --watch duration      repeat every DURATION (e.g. 5s) until Ctrl-C; 0 disables
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

- [polycli heimdall ops](polycli_heimdall_ops.md) - Node-operator commands backed by CometBFT JSON-RPC.
