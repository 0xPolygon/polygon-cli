Cast-like subcommands for interacting with a Heimdall v2 node. Targets
Polygon PoS node operators and validators who already have a REST
gateway (`:1317`) and a CometBFT RPC endpoint (`:26657`) and want to
inspect consensus state, query checkpoints/spans/milestones, or
broadcast the occasional signed transaction without reaching for
`curl + jq` or the `heimdalld` CLI.

The default network is `amoy` (Polygon testnet). Override with
`--mainnet`, `--network <name>`, or with explicit `--rest-url` /
`--rpc-url` flags.

```bash
# Liveness
polycli heimdall status
polycli heimdall block-number

# Checkpoints
polycli heimdall checkpoint latest
polycli heimdall checkpoint count

# Spans and validators
polycli heimdall span latest
polycli heimdall validator proposer
```

See `HEIMDALLCAST_REQUIREMENTS.md` for the full command catalogue.
