# `polycli abi decode`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Parse an ABI and print the encoded signatures.

```bash
polycli abi decode Contract.abi [flags]
```

## Flags

```bash
      --data string   Provide input data to be unpacked based on the ABI definition
      --file string   Provide a filename to read and analyze
  -h, --help          help for decode
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

- [polycli abi](polycli_abi.md) - Provides encoding and decoding functionalities with contract signatures and ABI.
