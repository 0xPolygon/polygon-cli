Cast-familiar block and chain queries against Heimdall's CometBFT RPC.

All subcommands talk to `/block`, `/status`, and `/abci_info` on the
CometBFT endpoint; REST is never used here. Height arguments accept a
bare integer, `latest`, or `earliest` — `finalized`, `safe`, and
`pending` are rejected with a hint (Heimdall has instant finality and
no pending queue at the consensus layer; those tags belong on
Ethereum).

```bash
# Latest block summary
polycli heimdall block

# A specific height, including the tx list
polycli heimdall block 32620627 --full

# Just the tip height, as a bare integer (for scripts)
polycli heimdall block-number

# Human time of a block
polycli heimdall age 32620627

# Find the height closest to a wall-clock time
polycli heimdall find-block 2026-04-20T15:10:00Z
polycli heimdall find-block 1776640500

# Chain identity
polycli heimdall chain-id
polycli heimdall chain

# Node software versions
polycli heimdall client
```
