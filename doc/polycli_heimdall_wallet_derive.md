# `polycli heimdall wallet derive`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Derive addresses from a BIP-39 mnemonic.

```bash
polycli heimdall wallet derive [flags]
```

## Flags

```bash
      --bip39-passphrase string   optional BIP-39 passphrase
      --count uint32              number of sequential addresses to derive (default 1)
  -h, --help                      help for derive
      --index uint32              starting address index when --path is not set
      --mnemonic string           BIP-39 mnemonic
      --mnemonic-file string      file containing a BIP-39 mnemonic
      --path string               derivation path (default m/44'/60'/0'/0/<index>)
      --show-private-key          also emit the derived private key on each line
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
