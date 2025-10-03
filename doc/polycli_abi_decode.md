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
      --data string   input data to be unpacked based on ABI definition
      --file string   filename to read and analyze
  -h, --help          help for decode
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

- [polycli abi](polycli_abi.md) - Provides encoding and decoding functionalities with contract signatures and ABI.
