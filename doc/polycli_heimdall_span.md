# `polycli heimdall span`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Query bor/span module endpoints.

```bash
polycli heimdall span [ID] [flags]
```

## Usage

Bor/span module queries (`x/bor`) against a Heimdall v2 node.

Alias: `sp`. `span <ID>` is a shorthand for `span get <ID>`.

All subcommands hit the REST gateway; byte-valued fields (pub keys,
signers, seeds) are rendered as `0x…`-hex by default and `--raw`
preserves the upstream base64.

```bash
# Module parameters and current/historical spans
polycli heimdall span params
polycli heimdall span latest
polycli heimdall span 5982
polycli heimdall span get 5982
polycli heimdall span list --limit 20
polycli heimdall span list --reverse=false          # oldest-first

# Derived / per-span helpers
polycli heimdall span producers 5982                 # selected_producers[] only
polycli heimdall span seed 5982                      # seed + seed_author

# Producer-set votes
polycli heimdall span votes                          # all votes
polycli heimdall span votes 4                        # votes for a single voter id

# Planned downtime and performance scores
polycli heimdall span downtime 4                     # prints `none` on 404
polycli heimdall span scores                         # performance-score map, desc

# Operator query: who produced this Bor block?
polycli heimdall span find 36985000                  # designated sprint producer
```

**Veblop caveat (post-Rio).** `span find` prints the *designated*
sprint producer from Heimdall span state. After Rio, Veblop rotates
producers based on performance scores and planned downtime, so the
actual block author may differ. To see the on-chain author, query the
Bor block itself (e.g. `cast block <N> --rpc-url <bor>` against a Bor
node). polycli heimdall does not talk to Bor.

## Flags

```bash
  -h, --help   help for span
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
- [polycli heimdall span downtime](polycli_heimdall_span_downtime.md) - Show planned downtime for a producer (or `none`).

- [polycli heimdall span find](polycli_heimdall_span_find.md) - Find the span covering a Bor block and its designated producer.

- [polycli heimdall span get](polycli_heimdall_span_get.md) - Fetch one span by id.

- [polycli heimdall span latest](polycli_heimdall_span_latest.md) - Show the current (latest) span.

- [polycli heimdall span list](polycli_heimdall_span_list.md) - Paginated span history.

- [polycli heimdall span params](polycli_heimdall_span_params.md) - Show bor module parameters.

- [polycli heimdall span producers](polycli_heimdall_span_producers.md) - List selected producers for a span.

- [polycli heimdall span scores](polycli_heimdall_span_scores.md) - Show validator performance scores (desc).

- [polycli heimdall span seed](polycli_heimdall_span_seed.md) - Show seed and seed_author for a span.

- [polycli heimdall span votes](polycli_heimdall_span_votes.md) - Show producer-set votes (all or by voter id).

