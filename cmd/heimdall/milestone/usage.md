Milestone module queries (`x/milestone`) against a Heimdall v2 node.

Alias: `ms`. `milestone <NUMBER>` is a shorthand for `milestone get
<NUMBER>`.

All subcommands hit the REST gateway; the `hash` field is rendered as
`0x…`-hex by default and `--raw` preserves the upstream base64.

```bash
# Thresholds + interval configured on this chain
polycli heimdall milestone params

# Total milestone count (cheap liveness signal)
polycli heimdall milestone count

# Latest milestone (hash decoded to hex)
polycli heimdall milestone latest

# One milestone by sequence number
polycli heimdall milestone 11602043
polycli heimdall milestone get 11602043
```

**Footgun.** The URL path (`/milestones/{number}`) uses a sequence
number that counts from 1 up to `milestone count`. The `milestone_id`
field inside the response body is **not** the same value — it is an
on-chain identifier minted by the proposer at milestone-creation time
(either a hex digest or a `uuid - 0x…` string, depending on vintage).
Both labels are printed to head off confusion.

An out-of-range `get` (number 0, or > count) surfaces a hint that the
valid range is `1..count`.
