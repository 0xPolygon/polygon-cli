# `polycli heimdall topup`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Query topup (dividend account) module endpoints.

## Usage

Topup module queries (`x/topup`) against a Heimdall v2 node.

All subcommands hit the REST gateway. Byte fields (`account_root_hash`,
`account_proof`, `tx_hash`) are rendered as `0x…`-hex by default; pass
the global `--raw` to preserve the upstream base64.

```bash
# Merkle root of all dividend accounts
polycli heimdall topup root

# Dividend account for an address (balance / fee_amount)
polycli heimdall topup account 0x02f615e95563ef16f10354dba9e584e58d2d4314

# Merkle proof for a dividend account (requires eth_rpc_url on the
# Heimdall node — an L1-less node surfaces an L1-not-configured hint)
polycli heimdall topup proof 0x02f615e95563ef16f10354dba9e584e58d2d4314

# Verify a submitted proof (proof is hex with or without 0x prefix)
polycli heimdall topup verify 0x02f615e95563ef16f10354dba9e584e58d2d4314 0x0000…0000

# Sequence / is-old replay keys for an L1 topup tx. Both require
# eth_rpc_url on the Heimdall node.
polycli heimdall topup sequence 0x48bd44a3…5c6bf8 423
polycli heimdall topup is-old   0x48bd44a3…5c6bf8 423
```

Endpoints covered (confirmed from heimdall-v2 `proto/heimdallv2/topup/query.proto`):

- `GET /topup/dividend-account-root`
- `GET /topup/dividend-account/{address}`
- `GET /topup/account-proof/{address}`
- `GET /topup/account-proof/{address}/verify?proof=…`
- `GET /topup/sequence?tx_hash=…&log_index=…`
- `GET /topup/is-old-tx?tx_hash=…&log_index=…`

## Flags

```bash
  -h, --help   help for topup
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
- [polycli heimdall topup account](polycli_heimdall_topup_account.md) - Fetch the dividend account for an address.

- [polycli heimdall topup is-old](polycli_heimdall_topup_is-old.md) - Check whether an L1 topup tx was already processed.

- [polycli heimdall topup proof](polycli_heimdall_topup_proof.md) - Fetch the Merkle proof for a dividend account.

- [polycli heimdall topup root](polycli_heimdall_topup_root.md) - Print the Merkle root of all dividend accounts.

- [polycli heimdall topup sequence](polycli_heimdall_topup_sequence.md) - Dedup sequence key for an L1 topup tx.

- [polycli heimdall topup verify](polycli_heimdall_topup_verify.md) - Verify a submitted Merkle proof for a dividend account.

