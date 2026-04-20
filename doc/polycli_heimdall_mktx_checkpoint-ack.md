# `polycli heimdall mktx checkpoint-ack`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Acknowledge a checkpoint on L2 (MsgCpAck, L1-mirroring).

```bash
polycli heimdall mktx checkpoint-ack [flags]
```

## Usage

Build, sign, and optionally broadcast a heimdallv2.checkpoint.MsgCpAck.

MsgCpAck is produced by the bridge after observing an L1 event. Manual
use is a replay that competes with the real bridge path; the command
refuses to run without --force. --l1-tx identifies the L1 tx hash the
operator intends to mirror (advisory — not part of the proto).
## Flags

```bash
      --account string           address or index into keystore (overrides --from for key lookup)
      --account-number uint      override fetched account number
      --derivation-path string   BIP-32 derivation path (default m/44'/60'/0'/0/<index>)
      --end-block uint           bor end block number
      --fee string               explicit fee coin amount, e.g. 10000pol (overrides --gas-price)
      --force                    bypass safety guards for L1-mirroring message types
      --from string              signer address (20-byte hex)
      --from-msg string          MsgCpAck.from address (default: signer)
      --gas uint                 gas limit (0 means estimate via simulation)
      --gas-adjustment float     multiplier applied to simulated gas to pick final gas limit (default 1.3)
      --gas-price float          fee price per gas unit in the default denom
  -h, --help                     help for checkpoint-ack
      --json                     emit JSON instead of key/value output
      --keystore-dir string      keystore directory (overrides ETH_KEYSTORE)
      --keystore-file string     explicit keystore JSON file path
      --l1-tx string             L1 transaction hash being mirrored (32 bytes hex)
      --memo string              optional tx memo
      --mnemonic string          BIP-39 mnemonic used to derive the signing key
      --mnemonic-index uint32    address index when deriving from --mnemonic
      --number uint              checkpoint number on Heimdall
      --password string          keystore password (mutually exclusive with --password-file)
      --password-file string     path to file containing keystore password
      --private-key string       hex-encoded secp256k1 private key (unsafe outside local dev)
      --proposer string          original proposer address of the checkpoint
      --root-hash string         32-byte root hash (hex)
      --sequence uint            override fetched sequence
      --sign-mode string         signing mode (direct|amino-json) (default "direct")
      --start-block uint         bor start block number
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
