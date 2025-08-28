# `polycli contract`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Interact with smart contracts and fetch contract information from the blockchain

```bash
polycli contract [flags]
```

## Usage

The `contract` is meant to help gathering smart contract information that is not directly available through the RPC endpoints

```bash
$ polycli contract --rpc-url "http://localhost:8545" --address "0x0000000000000000000000000000000000000001"
```

## Flags

```bash
      --address string   The contract address
  -h, --help             help for contract
      --rpc-url string   The RPC URL of the network containing the contract (default "http://localhost:8545")
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
