# `polycli heimdall util`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Local helpers: addr, b64, version, completions.

## Usage

Utility commands for Heimdall v2 developers.

Small local helpers that convert between address and encoding formats
commonly used when bouncing between Heimdall REST, CometBFT RPC, and
Polygon tooling. These subcommands do not touch the network unless
explicitly flagged.

```bash
# Convert a 0x-hex address to bech32 (cosmos1…) and back.
polycli heimdall util addr 0x02f615e95563ef16f10354dba9e584e58d2d4314
polycli heimdall util addr cosmos1qtmpt624v0h3dugr2nd6nevyukxj6sc54tvenp

# Print both forms.
polycli heimdall util addr 0x02f615e95563ef16f10354dba9e584e58d2d4314 --all

# Convert base64 blobs to 0x-hex and back (auto-detected, --to overrides).
polycli heimdall util b64 AQIDBA==
polycli heimdall util b64 0x01020304
polycli heimdall util b64 AQIDBA== --to hex

# Show polycli version; add --node to also fetch CometBFT /status.
polycli heimdall util version
polycli heimdall util version --node

# Emit shell completions.
polycli heimdall util completions bash > /etc/bash_completion.d/polycli
polycli heimdall util completions zsh > "${fpath[1]}/_polycli"
```

The bech32 human-readable part defaults to `cosmos` (Heimdall v2 inherits
the default cosmos-sdk account prefix per the API reference). Override
with `--hrp` if the target node uses a custom prefix.

## Flags

```bash
  -h, --help   help for util
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
- [polycli heimdall util addr](polycli_heimdall_util_addr.md) - Convert an address between 0x-hex and bech32.

- [polycli heimdall util b64](polycli_heimdall_util_b64.md) - Convert between base64 and 0x-hex.

- [polycli heimdall util completions](polycli_heimdall_util_completions.md) - Generate shell completion script.

- [polycli heimdall util version](polycli_heimdall_util_version.md) - Print polycli and (optionally) node version.

