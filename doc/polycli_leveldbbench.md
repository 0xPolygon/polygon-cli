# `polycli leveldbbench`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Perform a level db benchmark

```bash
polycli leveldbbench [flags]
```

## Usage

This command is meant to give us a sense of the system level performance for leveldb.

```bash
go run main.go leveldbbench --degree-of-parallelism 2 | jq '.' > result.json
```



## Flags

```bash
      --degree-of-parallelism uint8   The number of concurrent iops we'll perform (default 1)
  -h, --help                          help for leveldbbench
      --key-size uint                 The byte length of the keys that we'll use (default 8)
      --large-fill-limit uint         The number of large entries to write in the db (default 2000)
      --large-value-size uint         the number of random bytes to store for large tests (default 102400)
      --no-merge-write                allows disabling write merge
      --read-limit uint               the number of reads will attempt to complete in a given test (default 10000000)
      --small-fill-limit uint         The number of small entries to write in the db (default 1000000)
      --small-value-size uint         the number of random bytes to store (default 32)
      --sync-writes                   sync each write
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Fatal
                        200 Error
                        300 Warning
                        400 Info
                        500 Debug
                        600 Trace (default 400)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
