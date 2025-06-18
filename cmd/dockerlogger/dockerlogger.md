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