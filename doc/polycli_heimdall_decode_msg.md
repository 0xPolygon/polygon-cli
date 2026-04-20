# `polycli heimdall decode msg`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Decode a single Any.value for type-url (base64 value).

```bash
polycli heimdall decode msg <type-url> <value-b64> [flags]
```

## Usage

Decode a single Any.value for a registered type URL.

Example:
  polycli heimdall decode msg /heimdallv2.topup.MsgWithdrawFeeTx \
    CioweDAxNzE3MDAyN2YwYzVjZDE5MDRmOGI0MDU1OGRhZjUwN2FiNGViNjJhEgEw
## Flags

```bash
  -h, --help   help for msg
      --json   emit single-line JSON
      --list   print every registered type URL and exit
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

- [polycli heimdall decode](polycli_heimdall_decode.md) - Offline proto decoders for Heimdall bytes.
