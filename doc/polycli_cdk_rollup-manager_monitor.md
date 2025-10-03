# `polycli cdk rollup-manager monitor`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Watch for rollup manager events and display them on the fly

```bash
polycli cdk rollup-manager monitor [flags]
```

## Usage

This command will keep watching for rollup manager events on chain and print them on the fly.

Below is an example of how to use it

```bash
polycli cdk rollup-manager monitor
    --rpc-url https://sepolia.drpc.org
    --rollup-manager-address bali
```

## Flags

```bash
  -h, --help   help for monitor
```

The command also inherits flags from parent commands.

```bash
      --config string                   config file (default is $HOME/.polygon-cli.yaml)
      --fork-id string                  fork ID of CDK networks (default "12")
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

- [polycli cdk rollup-manager](polycli_cdk_rollup-manager.md) - Utilities for interacting with CDK rollup manager contract
