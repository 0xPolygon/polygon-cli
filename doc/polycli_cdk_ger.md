# `polycli cdk ger`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Utilities for interacting with CDK global exit root manager contract

## Flags

```bash
      --ger-address string   address of GER contract
  -h, --help                 help for ger
```

The command also inherits flags from parent commands.

```bash
      --config string                   config file (default is $HOME/.polygon-cli.yaml)
      --fork-id string                  fork ID of CDK networks (default "12")
      --pretty-logs                     output logs in pretty format instead of JSON (default true)
      --rollup-manager-address string   address of rollup contract
      --rpc-url string                  RPC URL of network containing CDK contracts (default "http://localhost:8545")
  -v, --verbosity string                log level (string or int):
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

- [polycli cdk](polycli_cdk.md) - Utilities for interacting with CDK networks
- [polycli cdk ger dump](polycli_cdk_ger_dump.md) - List detailed information about the global exit root manager

- [polycli cdk ger inspect](polycli_cdk_ger_inspect.md) - List some basic information about the global exit root manager

- [polycli cdk ger monitor](polycli_cdk_ger_monitor.md) - Watch for global exit root manager events and display them on the fly

