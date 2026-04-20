# `polycli heimdall checkpoint`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Query checkpoint module endpoints.

```bash
polycli heimdall checkpoint [ID] [flags]
```

## Usage

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

## Flags

```bash
  -h, --help   help for checkpoint
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
- [polycli heimdall checkpoint buffer](polycli_heimdall_checkpoint_buffer.md) - Show the in-flight (buffered) checkpoint.

- [polycli heimdall checkpoint count](polycli_heimdall_checkpoint_count.md) - Print total acked checkpoint count.

- [polycli heimdall checkpoint get](polycli_heimdall_checkpoint_get.md) - Fetch one checkpoint by id.

- [polycli heimdall checkpoint last-no-ack](polycli_heimdall_checkpoint_last-no-ack.md) - Print the timestamp of the last no-ack.

- [polycli heimdall checkpoint latest](polycli_heimdall_checkpoint_latest.md) - Show the latest acked checkpoint.

- [polycli heimdall checkpoint list](polycli_heimdall_checkpoint_list.md) - Paginated checkpoint history.

- [polycli heimdall checkpoint next](polycli_heimdall_checkpoint_next.md) - Compute the next checkpoint to propose.

- [polycli heimdall checkpoint overview](polycli_heimdall_checkpoint_overview.md) - Checkpoint module dashboard bundle.

- [polycli heimdall checkpoint params](polycli_heimdall_checkpoint_params.md) - Show checkpoint module parameters.

- [polycli heimdall checkpoint signatures](polycli_heimdall_checkpoint_signatures.md) - Aggregated validator signatures for a checkpoint tx.

