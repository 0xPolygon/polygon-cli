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
  -a, --addr ip           Address to bind discovery listener (default 127.0.0.1)
  -h, --help              help for ping
      --key string        Hex-encoded private key (cannot be set with --key-file)
  -k, --key-file string   Private key file (cannot be set with --key)
  -l, --listen            Keep the connection open and listen to the peer(s) (default true)
  -o, --output string     Write ping results to output file (default stdout)
  -p, --parallel int      How many parallel pings to attempt (default 16)
  -P, --port int          Port for discovery protocol (default 30303)
  -w, --wit               Whether to enable the wit/1 capability
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

- [polycli p2p](polycli_p2p.md) - Set of commands related to devp2p.
