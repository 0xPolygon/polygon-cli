State-sync / clerk queries (`x/clerk`) against a Heimdall v2 node.

Canonical name: `state-sync`. Aliases: `clerk`, `ss`.
`state-sync <ID>` is a shorthand for `state-sync get <ID>`.

All subcommands hit the REST gateway; the `data` field is rendered as
`0x…`-hex by default. Pass `--base64` (or the global `--raw`) to
preserve the upstream base64.

```bash
# Total event-record count (cheap liveness signal)
polycli heimdall state-sync count

# Latest L1 counter + processed flag (requires eth_rpc_url on the node)
polycli heimdall state-sync latest-id

# One record by id (bare shorthand + explicit form)
polycli heimdall state-sync 36610
polycli heimdall state-sync get 36610
polycli heimdall state-sync 36610 --base64

# Paginated history — PAGE-BASED (not Cosmos-pagination). The upstream
# /clerk/event-records/list endpoint rejects page=0 with HTTP 400, so
# --page defaults to 1 and --limit is mandatory (hint is surfaced if
# --limit is omitted).
polycli heimdall state-sync list --page 1 --limit 10

# Records since an id, optionally bounded by a timestamp (uses
# pagination.limit because /clerk/time unlike /list goes through the
# cosmos-pagination middleware).
polycli heimdall state-sync range --from-id 36600 --limit 5
polycli heimdall state-sync range --from-id 36600 --to-time 2026-04-20T13:00:00Z --limit 5

# Dedup / replay keys on the bridge. Both require eth_rpc_url on the
# Heimdall node — on an L1-less node the `connection refused` / gRPC
# code-13 response is surfaced as an L1-not-configured hint.
polycli heimdall state-sync sequence 0x48bd44a3...5c6bf8 423
polycli heimdall state-sync is-old 0x48bd44a3...5c6bf8 423
```
