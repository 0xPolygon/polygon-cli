# `polycli heimdall milestone`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Query milestone module endpoints.

```bash
polycli heimdall milestone [NUMBER] [flags]
```

## Usage

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

## Flags

```bash
  -h, --help   help for milestone
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
- [polycli heimdall milestone count](polycli_heimdall_milestone_count.md) - Print total milestone count.

- [polycli heimdall milestone get](polycli_heimdall_milestone_get.md) - Fetch one milestone by sequence number.

- [polycli heimdall milestone latest](polycli_heimdall_milestone_latest.md) - Show the latest milestone.

- [polycli heimdall milestone params](polycli_heimdall_milestone_params.md) - Show milestone module parameters.

