# `polycli dockerlogger`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Monitor and filter Docker container logs

```bash
polycli dockerlogger [flags]
```

## Usage

# Docker Logger

A tool to monitor and filter Docker container logs.

## Usage

```bash
dockerlogger --network <network-name> [flags]

Flags:
  --network string     Docker network name to monitor
  --all               Show all logs
  --errors            Show error logs
  --warnings          Show warning logs
  --info              Show info logs
  --debug             Show debug logs
  --filter string     Additional keywords to filter (comma-separated)
  --levels string     Comma-separated log levels to show (error,warn,info,debug)
  --service string    Filter logs by service names (comma-separated, partial match)
```
## Flags

```bash
      --all              Show all logs
      --debug            Show debug logs
      --errors           Show error logs
      --filter string    Additional keywords to filter, comma-separated
  -h, --help             help for dockerlogger
      --info             Show info logs
      --levels string    Comma-separated log levels to show (error,warn,info,debug)
      --network string   Docker network name to monitor
      --service string   Filter logs by service names (comma-separated, partial match)
      --warnings         Show warning logs
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
