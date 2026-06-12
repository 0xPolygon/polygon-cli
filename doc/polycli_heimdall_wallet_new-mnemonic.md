# `polycli heimdall wallet new-mnemonic`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Generate a new BIP-39 mnemonic and derive a key.

```bash
polycli heimdall wallet new-mnemonic [flags]
```

## Flags

```bash
      --bip39-passphrase string       optional BIP-39 passphrase (not the keystore password)
  -h, --help                          help for new-mnemonic
      --index uint32                  address index used when --path is not set
      --keystore-dir string           keystore directory (overrides ETH_KEYSTORE, ~/.foundry/keystores, ~/.polycli/keystores)
      --keystore-file string          explicit keystore JSON file path
      --ledger cast wallet --ledger   not supported; use cast wallet --ledger
      --password string               keystore password (mutually exclusive with --password-file)
      --password-file string          path to a file containing the keystore password
      --path string                   derivation path (default m/44'/60'/0'/0/<index>)
      --print-only                    print mnemonic and derived address without writing to keystore
      --trezor cast wallet --trezor   not supported; use cast wallet --trezor
      --words int                     mnemonic word count (12, 15, 18, 21, 24) (default 12)
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

- [polycli heimdall wallet](polycli_heimdall_wallet.md) - Manage keystores, keys, and message signatures.
