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
      --pretty-logs     output logs in pretty format instead of JSON (default true)
  -v, --verbosity int   0 - silent
                        100 panic
                        200 fatal
                        300 error
                        400 warning
                        500 info
                        600 debug
                        700 trace (default 500)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
