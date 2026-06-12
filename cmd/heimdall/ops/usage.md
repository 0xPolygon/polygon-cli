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
