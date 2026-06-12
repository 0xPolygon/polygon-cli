Utility commands for Heimdall v2 developers.

Small local helpers that convert between address and encoding formats
commonly used when bouncing between Heimdall REST, CometBFT RPC, and
Polygon tooling. These subcommands do not touch the network unless
explicitly flagged.

```bash
# Convert a 0x-hex address to bech32 (cosmos1…) and back.
polycli heimdall util addr 0x02f615e95563ef16f10354dba9e584e58d2d4314
polycli heimdall util addr cosmos1qtmpt624v0h3dugr2nd6nevyukxj6sc54tvenp

# Print both forms.
polycli heimdall util addr 0x02f615e95563ef16f10354dba9e584e58d2d4314 --all

# Convert base64 blobs to 0x-hex and back (auto-detected, --to overrides).
polycli heimdall util b64 AQIDBA==
polycli heimdall util b64 0x01020304
polycli heimdall util b64 AQIDBA== --to hex

# Show polycli version; add --node to also fetch CometBFT /status.
polycli heimdall util version
polycli heimdall util version --node

# Emit shell completions.
polycli heimdall util completions bash > /etc/bash_completion.d/polycli
polycli heimdall util completions zsh > "${fpath[1]}/_polycli"
```

The bech32 human-readable part defaults to `cosmos` (Heimdall v2 inherits
the default cosmos-sdk account prefix per the API reference). Override
with `--hrp` if the target node uses a custom prefix.
