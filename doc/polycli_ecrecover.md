# `polycli ecrecover`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Recovers and returns the public key of the signature

```bash
polycli ecrecover [flags]
```

## Usage

![GIF of `polycli monitor`](assets/monitor.gif)

If you're using the terminal UI and you'd like to be able to select text for copying, you might need to use a modifier key.

If you're experiencing missing blocks, try adjusting the `--batch-size` and `--interval` flags so that you poll for more blocks or more frequently.

## Flags

```bash
  -b, --block-number uint   Block number to check the extra data for (default: latest)
  -h, --help                help for ecrecover
  -r, --rpc-url string      The RPC endpoint url (default "http://localhost:8545")
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
