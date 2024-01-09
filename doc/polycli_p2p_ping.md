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
  -h, --help            help for ping
  -l, --listen          Keep the connection open and listen to the peer. This only works if the first
                        argument is an enode/enr, not a nodes file. (default true)
  -o, --output string   Write ping results to output file (default stdout)
  -p, --parallel int    How many parallel pings to attempt (default 16)
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
