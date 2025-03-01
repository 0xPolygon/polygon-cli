# `polycli cdk rollup-manager list-rollup-types`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

List some basic information about each rollup type

```bash
polycli cdk rollup-manager list-rollup-types [flags]
```

## Usage

This command will reach the rollup manager contract and retrieve all information about the rollups types registered.

Below is an example of how to use it

```bash
polycli cdk rollup-manager list-rollup-types
    --rpc-url https://sepolia.drpc.org
    --rollup-manager-address bali
```

## Flags

```bash
  -h, --help   help for list-rollup-types
```

The command also inherits flags from parent commands.

```bash
      --config string                   config file (default is $HOME/.polygon-cli.yaml)
      --fork-id string                  The ForkID of the cdk networks (default "12")
      --pretty-logs                     Should logs be in pretty format or JSON (default true)
      --rollup-manager-address string   The address of the rollup contract
      --rpc-url string                  The RPC URL of the network containing the CDK contracts (default "http://localhost:8545")
  -v, --verbosity int                   0 - Silent
                                        100 Panic
                                        200 Fatal
                                        300 Error
                                        400 Warning
                                        500 Info
                                        600 Debug
                                        700 Trace (default 500)
```

## See also

- [polycli cdk rollup-manager](polycli_cdk_rollup-manager.md) - Utilities for interacting with CDK rollup manager contract
