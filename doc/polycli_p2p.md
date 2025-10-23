# `polycli p2p`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Set of commands related to devp2p.

## Flags

```bash
  -h, --help   help for p2p
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
- [polycli p2p crawl](polycli_p2p_crawl.md) - Crawl a network on the devp2p layer and generate a nodes JSON file.

- [polycli p2p nodelist](polycli_p2p_nodelist.md) - Generate a node list to seed a node.

- [polycli p2p ping](polycli_p2p_ping.md) - Ping node(s) and return the output.

- [polycli p2p query](polycli_p2p_query.md) - Query block header(s) from node and prints the output.

- [polycli p2p sensor](polycli_p2p_sensor.md) - Start a devp2p sensor that discovers other peers and will receive blocks and transactions.

