# `polycli sendtx`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Send raw transactions to a JSON-RPC endpoint in batches.

```bash
polycli sendtx [flags]
```

## Usage

`polycli sendtx` reads pre-signed raw transactions from a file and sends them to a JSON-RPC endpoint using batch `eth_sendRawTransaction` requests.

The command is designed for high-throughput transaction injection. It reads the input file, groups transactions into batches, and sends them concurrently via HTTP POST.

## Usage

```bash
polycli sendtx --file txs.txt --rpc-url http://localhost:8545
```

## Input File Format

One hex-encoded raw transaction per line (with or without `0x` prefix):

```
0x02f86f8301388280...
0x02f86f8301388280...
```

Empty lines are skipped.

## Flags

```bash
  -b, --batch-size int    transactions per batch request (default 1000)
  -c, --concurrency int   concurrent batch requests (default: number of CPUs)
  -f, --file string       file containing raw transactions, one per line
  -h, --help              help for sendtx
  -r, --rpc-url string    RPC endpoint URL (default "http://localhost:8545")
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
