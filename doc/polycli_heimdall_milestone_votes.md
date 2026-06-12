# `polycli heimdall milestone votes`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Dump per-validator milestone votes over a height range.

```bash
polycli heimdall milestone votes [flags]
```

## Usage

Scan a range of heimdall heights and report, for every validator at every
height, whether it signed the commit and which milestone proposition (bor
block range) its vote extension carried. Heights where a milestone was
finalized are correlated so late or missing validators stand out.

--from/--to are vote heights H; the data is read from block H+1, which
carries the previous height's vote extensions as its first transaction.
## Flags

```bash
      --concurrency int     number of heights fetched in parallel (default 8)
  -f, --field stringArray   pluck one or more fields (repeatable)
      --from string         first vote height to scan
      --from-time string    start of scan as unix seconds or RFC3339 (resolved to a height)
  -h, --help                help for votes
      --missing-only        only show votes that did not commit or carried no proposition
      --summary             aggregate to one row per validator
      --to string           last vote height to scan
      --to-time string      end of scan as unix seconds or RFC3339 (resolved to a height)
      --watch duration      repeat every DURATION (e.g. 5s) until Ctrl-C; 0 disables
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

- [polycli heimdall milestone](polycli_heimdall_milestone.md) - Query milestone module endpoints.
