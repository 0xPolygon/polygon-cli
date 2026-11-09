# `polycli p2p nodelist`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Generate a node list to seed a node.

```bash
polycli p2p nodelist [nodes.json] [flags]
```

## Flags

```bash
      --clickhouse-dsn string   ClickHouse DSN of the sensor database to read peers from, e.g. clickhouse://user:pass@host:9000/mainnet
  -h, --help                    help for nodelist
  -l, --limit int               number of unique nodes to return (default 100)
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

- [polycli p2p](polycli_p2p.md) - Set of commands related to devp2p.
