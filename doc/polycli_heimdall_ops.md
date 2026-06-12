# `polycli heimdall ops`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Node-operator commands backed by CometBFT JSON-RPC.

## Usage

# heimdall ops

Operator-facing commands backed by the CometBFT JSON-RPC endpoint
(`:26657`). Covers liveness, sync, peers, mempool, consensus, and
validator set inspection — the CometBFT-layer view that sits under
Heimdall.

## Examples

```bash
# single-shot liveness check (exits 0 only if /health returns OK)
polycli heimdall ops health

# one-line status snapshot
polycli heimdall ops status

# list peers
polycli heimdall ops peers

# full peer detail
polycli heimdall ops peers --verbose

# CometBFT-layer validator set (NOT Heimdall x/stake; see `heimdall validator`)
polycli heimdall ops validators-cometbft

# signed header for a height
polycli heimdall ops commit 32634175

# pending tx count and list
polycli heimdall ops tx-pool
polycli heimdall ops tx-pool --list

# app identity / last block hash
polycli heimdall ops abci-info

# consensus round/step summary (expensive on a busy node)
polycli heimdall ops consensus
```

All subcommands honour the heimdall root flags (`--rpc-url`,
`--network`, `--json`, `--field`, `--curl`, `--timeout`, etc.).

## Flags

```bash
  -h, --help   help for ops
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
- [polycli heimdall ops abci-info](polycli_heimdall_ops_abci-info.md) - Show CometBFT /abci_info app identity.

- [polycli heimdall ops commit](polycli_heimdall_ops_commit.md) - Fetch signed CometBFT commit header.

- [polycli heimdall ops consensus](polycli_heimdall_ops_consensus.md) - Summarise CometBFT /dump_consensus_state.

- [polycli heimdall ops health](polycli_heimdall_ops_health.md) - Probe CometBFT /health; exit 0 on success.

- [polycli heimdall ops peers](polycli_heimdall_ops_peers.md) - List peers from CometBFT /net_info.

- [polycli heimdall ops status](polycli_heimdall_ops_status.md) - Show CometBFT /status: height, sync, moniker, own validator.

- [polycli heimdall ops tx-pool](polycli_heimdall_ops_tx-pool.md) - Show CometBFT mempool size (--list for hashes).

- [polycli heimdall ops validators-cometbft](polycli_heimdall_ops_validators-cometbft.md) - List CometBFT consensus validators (NOT Heimdall x/stake).

