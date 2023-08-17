# `polycli sensor nodelist`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Generate a node list to seed a node

```bash
polycli sensor nodelist [nodes.json] [flags]
```

## Flags

```bash
  -h, --help        help for nodelist
  -l, --limit int   Number of unique nodes to return (default 100)
```

The command also inherits flags from parent commands.

```bash
      --config string       config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs         Should logs be in pretty format or JSON (default true)
  -P, --project-id string   GCP project ID
  -v, --verbosity int       0 - Silent
                            100 Fatal
                            200 Error
                            300 Warning
                            400 Info
                            500 Debug
                            600 Trace (default 400)
```

## See also

- [polycli sensor](polycli_sensor.md) - Start a devp2p sensor that discovers other peers and will receive blocks and transactions.
