# `polycli heimdall wallet sign`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Sign a message with a keystore key.

```bash
polycli heimdall wallet sign <message-or-hash> [flags]
```

## Flags

```bash
      --address string                address of the keystore key to sign with
  -h, --help                          help for sign
      --keystore-dir string           keystore directory (overrides ETH_KEYSTORE, ~/.foundry/keystores, ~/.polycli/keystores)
      --keystore-file string          explicit keystore JSON file path
      --ledger cast wallet --ledger   not supported; use cast wallet --ledger
      --password string               keystore password (mutually exclusive with --password-file)
      --password-file string          path to a file containing the keystore password
      --private-key string            hex-encoded private key (skips the keystore)
      --raw                           sign the 32-byte hash directly (no EIP-191 framing)
      --trezor cast wallet --trezor   not supported; use cast wallet --trezor
      --yes                           skip confirmation prompts
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

- [polycli heimdall wallet](polycli_heimdall_wallet.md) - Manage keystores, keys, and message signatures.
