# `polycli fork`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Take a forked block and walk up the chain to do analysis.

```bash
polycli fork blockhash url [flags]
```

## Flags

```bash
  -h, --help   help for fork
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Panic
                        200 Fatal
                        300 Error
                        400 Warning
                        500 Info
                        600 Debug
                        700 Trace (default 500)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
