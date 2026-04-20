# `polycli heimdall`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Query and interact with a Heimdall v2 node.

## Usage

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

## Flags

```bash
      --amoy                     shortcut for --network amoy (default)
      --chain-id string          chain id used for signing
      --color string             color mode (auto|always|never) (default "auto")
      --curl                     print the equivalent curl command instead of executing
      --denom string             fee denom
      --heimdall-config string   path to heimdall config TOML (default ~/.polycli/heimdall.toml)
  -h, --help                     help for heimdall
  -k, --insecure                 accept invalid TLS certs
      --json                     emit JSON instead of key/value
      --mainnet                  shortcut for --network mainnet
  -N, --network string           named network preset (amoy|mainnet)
      --no-color                 disable color output
      --raw                      preserve raw bytes (no 0x-hex normalization)
  -r, --rest-url string          heimdall REST gateway URL
      --rpc-headers string       extra request headers, comma-separated key=value pairs
  -R, --rpc-url string           cometBFT RPC URL
      --timeout int              HTTP timeout in seconds
```

The command also inherits flags from parent commands.

```bash
      --config string      config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs        output logs in pretty format instead of JSON (default true)
  -v, --verbosity string   log level (string or int):
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

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
- [polycli heimdall age](polycli_heimdall_age.md) - Show the timestamp of a CometBFT block.

- [polycli heimdall balance](polycli_heimdall_balance.md) - Show an account's balance for a denom.

- [polycli heimdall block](polycli_heimdall_block.md) - Show a CometBFT block by height (or latest).

- [polycli heimdall block-number](polycli_heimdall_block-number.md) - Print the latest CometBFT block height.

- [polycli heimdall chain](polycli_heimdall_chain.md) - Print the human-readable chain name.

- [polycli heimdall chain-id](polycli_heimdall_chain-id.md) - Print the CometBFT chain id.

- [polycli heimdall checkpoint](polycli_heimdall_checkpoint.md) - Query checkpoint module endpoints.

- [polycli heimdall client](polycli_heimdall_client.md) - Show Heimdall app + CometBFT versions.

- [polycli heimdall find-block](polycli_heimdall_find-block.md) - Find the block height closest to a timestamp.

- [polycli heimdall logs](polycli_heimdall_logs.md) - Query the CometBFT tx index.

- [polycli heimdall nonce](polycli_heimdall_nonce.md) - Print an account's sequence number.

- [polycli heimdall publish](polycli_heimdall_publish.md) - Broadcast a signed TxRaw (base64 or hex).

- [polycli heimdall receipt](polycli_heimdall_receipt.md) - Show a transaction receipt (events + logs).

- [polycli heimdall rpc](polycli_heimdall_rpc.md) - Invoke an arbitrary CometBFT JSON-RPC method.

- [polycli heimdall sequence](polycli_heimdall_sequence.md) - Alias of nonce; print an account's sequence.

- [polycli heimdall span](polycli_heimdall_span.md) - Query bor/span module endpoints.

- [polycli heimdall tx](polycli_heimdall_tx.md) - Show a transaction by hash.

