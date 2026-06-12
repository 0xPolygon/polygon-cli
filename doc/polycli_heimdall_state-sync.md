# `polycli heimdall state-sync`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Query state-sync (clerk) module endpoints.

```bash
polycli heimdall state-sync [ID] [flags]
```

## Usage

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

## Flags

```bash
  -h, --help             help for state-sync
      --watch duration   repeat every DURATION (e.g. 5s) until Ctrl-C; 0 disables
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
- [polycli heimdall state-sync count](polycli_heimdall_state-sync_count.md) - Print total state-sync (clerk) event-record count.

- [polycli heimdall state-sync get](polycli_heimdall_state-sync_get.md) - Fetch one event-record by id.

- [polycli heimdall state-sync is-old](polycli_heimdall_state-sync_is-old.md) - Check whether an L1 state-sync event was already replayed.

- [polycli heimdall state-sync latest-id](polycli_heimdall_state-sync_latest-id.md) - Latest L1 state-sync counter (needs eth_rpc_url).

- [polycli heimdall state-sync list](polycli_heimdall_state-sync_list.md) - Paginated event-record history (page-based).

- [polycli heimdall state-sync range](polycli_heimdall_state-sync_range.md) - Event-records since an id (optional time bound).

- [polycli heimdall state-sync sequence](polycli_heimdall_state-sync_sequence.md) - Dedup sequence key for an L1 state-sync event.

