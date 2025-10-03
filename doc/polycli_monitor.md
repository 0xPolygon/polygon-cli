# `polycli monitor`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Monitor blocks using a JSON-RPC endpoint.

```bash
polycli monitor [flags]
```

## Usage

![GIF of `polycli monitor`](assets/monitor.gif)

If you're using the terminal UI and you'd like to be able to select text for copying, you might need to use a modifier key.

If you're experiencing missing blocks, try adjusting the `--batch-size` and `--interval` flags so that you poll for more blocks or more frequently.

## Flags

```bash
  -b, --batch-size string    number of requests per batch (default "auto")
  -c, --cache-limit int      number of cached blocks for the LRU block data structure (Min 100) (default 200)
  -h, --help                 help for monitor
  -i, --interval string      amount of time between batch block RPC calls (default "5s")
  -r, --rpc-url string       the RPC endpoint URL (default "http://localhost:8545")
  -s, --sub-batch-size int   number of requests per sub-batch (default 50)
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
