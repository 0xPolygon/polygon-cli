Checkpoint module queries (`x/checkpoint`) against a Heimdall v2 node.

Alias: `cp`. `checkpoint <ID>` is a shorthand for `checkpoint get <ID>`.

All subcommands hit the REST gateway; `root_hash` is rendered as
`0x…`-hex by default and `--raw` preserves the upstream base64.

```bash
# Current and historical checkpoints
polycli heimdall checkpoint count
polycli heimdall checkpoint latest
polycli heimdall checkpoint 38871
polycli heimdall checkpoint get 38871

# In-flight / system state
polycli heimdall checkpoint buffer          # prints `empty` for zero-address proposer
polycli heimdall checkpoint last-no-ack     # unix seconds + human age
polycli heimdall checkpoint next            # prepare-next (requires node to have L1 RPC configured)

# Paginated history
polycli heimdall checkpoint list --limit 20 --reverse
polycli heimdall checkpoint list --page AAAA...

# Signatures for a specific checkpoint-ack tx hash (0x prefix optional)
polycli heimdall checkpoint signatures 0x94297F18F736A0C018E4871A5257384450673AC8441F8F7956523231D74D2A29

# Dashboard bundle
polycli heimdall checkpoint overview
```
