# `polycli heimdall wallet import`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Import an existing key into the keystore.

```bash
polycli heimdall wallet import [flags]
```

## Flags

```bash
      --bip39-passphrase string       optional BIP-39 passphrase
  -h, --help                          help for import
      --index uint32                  address index when --path is not set
      --keystore-dir string           keystore directory (overrides ETH_KEYSTORE, ~/.foundry/keystores, ~/.polycli/keystores)
      --keystore-file string          explicit keystore JSON file path
      --mnemonic string               BIP-39 mnemonic
      --mnemonic-file string          file containing a BIP-39 mnemonic
      --password string               keystore password (mutually exclusive with --password-file)
      --password-file string          path to a file containing the keystore password
      --path string                   derivation path (default m/44'/60'/0'/0/<index>)
      --private-key string            hex-encoded secp256k1 private key
      --source-keystore-file string   path to an existing v3 JSON keystore to import
      --source-password-file string   file with the existing keystore's password
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
