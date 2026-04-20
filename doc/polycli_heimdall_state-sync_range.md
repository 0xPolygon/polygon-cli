# `polycli heimdall state-sync range`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Event-records since an id (optional time bound).

```bash
polycli heimdall state-sync range [flags]
```

## Flags

```bash
      --base64 data         preserve raw base64 for data (default 0x-hex)
  -f, --field stringArray   pluck one or more fields (repeatable, --json only)
      --from-id uint        lowest event-record id to return (required)
  -h, --help                help for range
      --limit int           maximum entries to return
      --to-time string      RFC3339 upper bound on record_time
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

- [polycli heimdall state-sync](polycli_heimdall_state-sync.md) - Query state-sync (clerk) module endpoints.
