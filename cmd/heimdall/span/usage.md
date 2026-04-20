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
