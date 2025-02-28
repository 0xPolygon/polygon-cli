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
      --fork-id string                  The ForkID of the cdk networks (default "12")
  -h, --help                            help for cdk
      --rollup-manager-address string   The address of the rollup contract
      --rpc-url string                  The RPC URL of the network containing the CDK contracts (default "http://localhost:8545")
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
- [polycli cdk bridge](polycli_cdk_bridge.md) - Utilities for interacting with CDK bridge contract

- [polycli cdk ger](polycli_cdk_ger.md) - Utilities for interacting with CDK global exit root manager contract

- [polycli cdk rollup](polycli_cdk_rollup.md) - Utilities for interacting with CDK rollup manager to get rollup specific information

- [polycli cdk rollup-manager](polycli_cdk_rollup-manager.md) - Utilities for interacting with CDK rollup manager contract

