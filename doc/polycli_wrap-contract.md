# `polycli wrap-contract`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Wrap deployed bytecode into create bytecode.

```bash
polycli wrap-contract bytecode [flags]
```

## Usage

This command takes the runtime bytecode, the bytecode deployed on-chain, as input and converts it into creation bytecode, the bytecode used to create the contract

```bash
$ polycli wrap-contract 69602a60005260206000f3600052600a6016f3

```

The resulting bytecode will be formatted this way:

		0x??   // storage initialization code if any
		63??   // PUSH4 to indicate the size of the data that should be copied into memory
		63??   // PUSH4 to indicate the offset in the call data to start the copy
		6000   // PUSH1 00 to indicate the destination offset in memory
		39     // CODECOPY
		63??   // PUSH4 to indicate the size of the data to be returned from memory
		6000   // PUSH1 00 to indicate that it starts from offset 0
		F3     // RETURN
		??,    // Deployed Bytecode

## Flags

```bash
  -h, --help   help for wrap-contract
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
