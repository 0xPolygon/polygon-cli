# `polycli p2p ping`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Ping node(s) and return the output.

```bash
polycli p2p ping [enode/enr or nodes file] [flags]
```

## Usage

Ping nodes by either giving a single enode/enr or an entire nodes file.

This command will establish a handshake and status exchange to get the Hello and
Status messages and output JSON. If providing a enode/enr rather than a nodes
file, then the connection will remain open by default (--listen=true), and you
can see other messages the peer sends (e.g. blocks, transactions, etc.).
## Flags

```bash
  -a, --addr ip           address to bind discovery listener (default 127.0.0.1)
  -h, --help              help for ping
      --key string        hex-encoded private key (cannot be set with --key-file)
  -k, --key-file string   private key file (cannot be set with --key)
  -l, --listen            keep connection open and listen to peer. This only works if first
                          argument is an enode/enr, not a nodes file (default true)
  -o, --output string     write ping results to output file (default stdout)
  -p, --parallel int      how many parallel pings to attempt (default 16)
  -P, --port int          port for discovery protocol (default 30303)
  -w, --wit               enable wit/1 capability
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
