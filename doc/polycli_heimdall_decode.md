# `polycli heimdall decode`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Offline proto decoders for Heimdall bytes.

## Usage

# decode

Local, offline proto decoders for Heimdall v2 transactions, messages,
and vote extensions. All commands accept base64 (default) or 0x-prefixed
hex and never touch the network.

## Subcommands

- `decode tx <B64_OR_HEX>`
  Parses `cosmos.tx.v1beta1.TxRaw`, resolves each `Any.type_url` via the
  internal registry, and pretty-prints body + auth info + signature
  metadata.

- `decode msg <TYPE_URL> <B64>`
  Decodes a single `Any.value` for the provided type URL (the registry
  includes every Msg implemented by `polycli heimdall send`).

- `decode hash-tx <B64_OR_HEX>`
  Returns the upper-case `SHA256(txraw)` hash CometBFT uses to address
  transactions (what `polycli heimdall tx` looks up).

- `decode ve <HEX>`
  Parses CometBFT vote-extension bytes as
  `heimdallv2.sidetxs.VoteExtension`. Input is hex because vote
  extensions surface as hex in CometBFT logs.

Type URLs registered locally are listed at
`polycli heimdall decode msg --list`.

## Flags

```bash
  -h, --help   help for decode
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
- [polycli heimdall decode hash-tx](polycli_heimdall_decode_hash-tx.md) - CometBFT SHA-256 hash of a TxRaw.

- [polycli heimdall decode msg](polycli_heimdall_decode_msg.md) - Decode a single Any.value for type-url (base64 value).

- [polycli heimdall decode tx](polycli_heimdall_decode_tx.md) - Decode a TxRaw (base64 or 0x-hex).

- [polycli heimdall decode ve](polycli_heimdall_decode_ve.md) - Decode CometBFT vote-extension bytes.

