# `polycli heimdall mktx stake-join`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Register a validator (MsgValidatorJoin, L1-mirroring).

```bash
polycli heimdall mktx stake-join [flags]
```

## Usage

Build, sign, and optionally broadcast a heimdallv2.stake.MsgValidatorJoin.

Produced by the bridge after a StakingInfo event; manual use requires
--force.
## Flags

```bash
      --account string           address or index into keystore (overrides --from for key lookup)
      --account-number uint      override fetched account number
      --activation-epoch uint    activation epoch
      --amount string            stake amount (decimal string)
      --block-number uint        L1 block number
      --derivation-path string   BIP-32 derivation path (default m/44'/60'/0'/0/<index>)
      --fee string               explicit fee coin amount, e.g. 10000pol (overrides --gas-price)
      --force                    bypass safety guards for L1-mirroring message types
      --from string              signer address (20-byte hex)
      --from-msg string          MsgValidatorJoin.from address (default: signer)
      --gas uint                 gas limit (0 means estimate via simulation)
      --gas-adjustment float     multiplier applied to simulated gas to pick final gas limit (default 1.3)
      --gas-price float          fee price per gas unit in the default denom
  -h, --help                     help for stake-join
      --json                     emit JSON instead of key/value output
      --keystore-dir string      keystore directory (overrides ETH_KEYSTORE)
      --keystore-file string     explicit keystore JSON file path
      --log-index uint           L1 log index
      --memo string              optional tx memo
      --mnemonic string          BIP-39 mnemonic used to derive the signing key
      --mnemonic-index uint32    address index when deriving from --mnemonic
      --nonce-l1 uint            L1 stake nonce
      --password string          keystore password (mutually exclusive with --password-file)
      --password-file string     path to file containing keystore password
      --private-key string       hex-encoded secp256k1 private key (unsafe outside local dev)
      --sequence uint            override fetched sequence
      --sign-mode string         signing mode (direct|amino-json) (default "direct")
      --signer-pub-key string    validator signer pubkey (hex)
      --tx-hash string           L1 tx hash (32 bytes hex)
      --val-id uint              validator id
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

- [polycli heimdall mktx](polycli_heimdall_mktx.md) - Build a signed TxRaw without broadcasting.
