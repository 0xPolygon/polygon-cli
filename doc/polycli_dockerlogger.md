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
      --all              show all logs
      --debug            show debug logs
      --errors           show error logs
      --filter string    additional keywords to filter, comma-separated
  -h, --help             help for dockerlogger
      --info             show info logs
      --levels string    comma-separated log levels to show (error,warn,info,debug)
      --network string   docker network name to monitor
      --service string   filter logs by service names (comma-separated, partial match)
      --warnings         show warning logs
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
