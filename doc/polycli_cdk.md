# `polycli cdk`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Utilities for interacting with CDK networks

## Usage

Basic utility commands for interacting with the cdk contracts
## Flags

```bash
      --fork-id string                  fork ID of CDK networks (default "12")
  -h, --help                            help for cdk
      --rollup-manager-address string   address of rollup contract
      --rpc-url string                  RPC URL of network containing CDK contracts (default "http://localhost:8545")
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
- [polycli cdk bridge](polycli_cdk_bridge.md) - Utilities for interacting with CDK bridge contract

- [polycli cdk ger](polycli_cdk_ger.md) - Utilities for interacting with CDK global exit root manager contract

- [polycli cdk rollup](polycli_cdk_rollup.md) - Utilities for interacting with CDK rollup manager to get rollup specific information

- [polycli cdk rollup-manager](polycli_cdk_rollup-manager.md) - Utilities for interacting with CDK rollup manager contract

