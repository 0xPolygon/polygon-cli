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

- [polycli abi](polycli_abi.md) - Provides encoding and decoding functionalities with contract signatures and ABI.
