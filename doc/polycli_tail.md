# `polycli tail`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Tail full blocks from a JSON-RPC endpoint as NDJSON.

```bash
polycli tail [flags]
```

## Usage

Tail full blocks from an RPC endpoint and emit each block as newline-delimited JSON.

By default this prints the latest 10 blocks and exits:

```bash
polycli tail --rpc-url http://127.0.0.1:8545
```

Tail the last 100 blocks and keep following new blocks:

```bash
polycli tail -n 100 --follow --rpc-url http://127.0.0.1:8545
```

## Flags

```bash
  -b, --batch-size uint          batch size for block requests (default 150)
  -n, --blocks-back uint         number of latest blocks to output before following (default 10)
      --follow                   poll for and stream newly produced blocks
  -h, --help                     help for tail
      --poll-interval duration   poll interval when --follow is enabled (default 2s)
  -r, --rpc-url string           the RPC endpoint URL (default "http://localhost:8545")
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
