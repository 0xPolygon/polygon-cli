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
