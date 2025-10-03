# `polycli cdk ger dump`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

List detailed information about the global exit root manager

```bash
polycli cdk ger dump [flags]
```

## Usage

This command will reach the global exit root contract and retrieve detailed information.

Below is an example of how to use it

```bash
polycli cdk ger dump
    --rpc-url https://sepolia.drpc.org
    --rollup-manager-address bali
```

## Flags

```bash
  -h, --help   help for dump
```

The command also inherits flags from parent commands.

```bash
      --config string                   config file (default is $HOME/.polygon-cli.yaml)
      --fork-id string                  fork ID of CDK networks (default "12")
      --ger-address string              address of GER contract
      --pretty-logs                     output logs in pretty format instead of JSON (default true)
      --rollup-manager-address string   address of rollup contract
      --rpc-url string                  RPC URL of network containing CDK contracts (default "http://localhost:8545")
  -v, --verbosity int                   0 - silent
                                        100 panic
                                        200 fatal
                                        300 error
                                        400 warning
                                        500 info
                                        600 debug
                                        700 trace (default 500)
```

## See also

- [polycli cdk ger](polycli_cdk_ger.md) - Utilities for interacting with CDK global exit root manager contract
