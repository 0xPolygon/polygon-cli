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
polycli wrap-contract bytecode|file [flags]
```

## Usage

This command takes the runtime bytecode, the bytecode deployed on-chain, as input and converts it into creation bytecode, the bytecode used to create the contract

```bash
$ polycli wrap-contract 69602a60005260206000f3600052600a6016f3
$ echo 69602a60005260206000f3600052600a6016f3 | polycli wrap-contract 

```

You can also provide a path to a file, and the bytecode while be read from there.

```bash
$ polycli wrap-contract bytecode.txt
$ polycli wrap-contract ../bytecode.txt
$ polycli wrap-contract /tmp/bytecode.txt
$ echo /tmp/bytecode.txt | polycli wrap-contract
```

Additionally, you can provide storage for the contract in JSON
```bash
$ polycli wrap-contract 0x4455 --storage '{"0x01":"0x0034"}'
$ polycli wrap-contract 0x4455 --storage '{"0x01":"0x0034", "0x02": "0xFF"}'
$ echo 69602a60005260206000f3600052600a6016f3 | polycli wrap-contract --storage '{"0x01":"0x0034", "0x02": "0xFF"}'
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
  -h, --help             help for wrap-contract
      --storage string   storage slots in JSON format k:v
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

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
