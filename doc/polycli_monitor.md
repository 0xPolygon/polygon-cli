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
polycli monitor url [flags]
```

## Usage

![GIF of `polycli monitor`](assets/monitor.gif)

If you're using the terminal UI and you'd like to be able to select text for copying, you might need to use a modifier key.

If you're experiencing missing blocks, try adjusting the `--batch-size` and `--interval` flags so that you poll for more blocks or more frequently.

## Flags

```bash
  -b, --batch-size uint   Number of requests per batch (default 25)
  -h, --help              help for monitor
  -i, --interval string   Amount of time between batch block rpc calls (default "5s")
  -w, --window-size int   Number of blocks visible in the window (default 25)
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
