Milestone module queries (`x/milestone`) against a Heimdall v2 node.

Alias: `ms`. `milestone <NUMBER>` is a shorthand for `milestone get
<NUMBER>`.

The query subcommands hit the REST gateway; the `hash` field is
rendered as `0x…`-hex by default and `--raw` preserves the upstream
base64. `votes` additionally scans blocks over the CometBFT RPC.

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

# Per-validator milestone votes for the last 1000 heights (default)
polycli heimdall milestone votes

# A specific range of vote heights, machine-readable
polycli heimdall milestone votes --from 30001 --to 30200 --json

# A time window, aggregated to one row per validator
polycli heimdall milestone votes --from-time 2026-06-11T06:00:00Z \
  --to-time 2026-06-11T07:00:00Z --summary

# Only the problem rows: validators that did not sign, or signed
# without a milestone proposition
polycli heimdall milestone votes --missing-only
```

**How `votes` reads the chain.** Validators vote on milestones through
CometBFT vote extensions: the extension attached to each validator's
pre-commit on height `H` carries a `MilestoneProposition` (its view of
the bor chain). The proposer of `H+1` embeds the previous height's
extended commit — every validator's vote, including the absent ones —
as the first transaction of block `H+1`. `--from`/`--to` are vote
heights `H`; the command fetches block (and block_results) `H+1`.

Per row: `flag` is the commit flag (`ABSENT` = the validator never
delivered a pre-commit — the strongest "node down" signal), `prop_*`
describe the proposed bor block range (`-` for a commit without a
proposition), `lag` is the distance between the 2/3-majority end block
and the validator's proposed end block at that height, and `milestone`
is the milestone number finalized from that height's votes (`miss` =
the validator's proposition did not cover the finalized end block).

**Caveats.** `val_id` is resolved from the *current* validator set
(`/stake/validators-set`); validators that rotated out render as `-`.
Under `--curl` only one representative request pair is printed. The
scan is all-or-nothing: a height that keeps failing aborts the run
after retries so the output is never silently incomplete.

**Footgun.** The URL path (`/milestones/{number}`) uses a sequence
number that counts from 1 up to `milestone count`. The `milestone_id`
field inside the response body is **not** the same value — it is an
on-chain identifier minted by the proposer at milestone-creation time
(either a hex digest or a `uuid - 0x…` string, depending on vintage).
Both labels are printed to head off confusion.

An out-of-range `get` (number 0, or > count) surfaces a hint that the
valid range is `1..count`.
