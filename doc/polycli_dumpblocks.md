# `polycli dumpblocks`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Export a range of blocks from a JSON-RPC endpoint.

```bash
polycli dumpblocks start end [flags]
```

## Usage

For various reasons, we might want to dump a large range of blocks for analytics or replay purposes. This is a simple util to export over RPC a range of blocks.

The following command would download the first 500K blocks and zip them and then look for blocks with transactions that create an account.

```bash
$ polycli dumpblocks 0 500000 --rpc-url http://172.26.26.12:8545/ | gzip > foo.gz
$ zcat < foo.gz | jq '. | select(.transactions | length > 0) | select(.transactions[].to == null)'
```

Dumpblocks can also output to protobuf format.

If you wish to make changes to the protobuf.

1. Install the protobuf compiler

On GNU/Linux:

```bash
$ sudo apt install protoc-gen-go
```

On a MAC:

```bash
$ brew install protobuf
```

2. Install the protobuf plugin

```bash
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.9
```

3. Compile the proto file

```bash
$ make generate
```

4. Depending on what endpoint and chain you're querying, you may be missing some fields in the proto definition.

```json
{
  "level": "error",
  "error": "proto: (line 1:813): unknown field \"nonce\"",
  "time": "2023-01-17T13:35:53-05:00",
  "message": "failed to unmarshal proto message"
}
```

To solve this, add the unknown fields to the `.proto` files and recompile them (step 3).

## Flags

```bash
  -b, --batch-size uint    batch size for requests (most providers cap at 1000) (default 150)
  -c, --concurrency uint   how many go routines to leverage (default 1)
  -B, --dump-blocks        dump blocks to output (default true)
      --dump-receipts      dump receipts to output (default true)
  -f, --filename string    where to write the output to (default stdout)
  -F, --filter string      filter output based on tx to and from, not setting a filter means all are allowed (default "{}")
  -h, --help               help for dumpblocks
  -m, --mode string        the output format [json, proto] (default "json")
  -r, --rpc-url string     the RPC endpoint URL (default "http://localhost:8545")
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
