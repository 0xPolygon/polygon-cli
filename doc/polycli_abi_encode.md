# `polycli abi encode`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

ABI encodes a function signature and the inputs

```bash
polycli abi encode [function signature] [args...] [flags]
```

## Usage

[function-signature] is required and is a fragment in the form <function name>(<types...>). If the function signature has parameters, then those values would have to be passed as arguments after the function signature.
## Flags

```bash
  -h, --help   help for encode
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
