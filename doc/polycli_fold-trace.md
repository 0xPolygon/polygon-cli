# `polycli fold-trace`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Trace an execution trace and fold it for visualization

```bash
polycli fold-trace [flags]
```

## Usage

This command is meant to take a transaction op code trace and convert it into a folded output that can be easily visualized with Flamegraph tools.

```bash
# First grab a trace from an RPC that supports the debug namespace
cast rpc --rpc-url http://127.0.0.1:18545 debug_traceTransaction 0x12f63f489213f5bd5b88fbfb12960b8248f61e2062a369ba41d8a3c96bb74d57 > trace.json

# Read the trace and use the `fold-trace` command and write the output
polycli fold-trace --metric actualgas < trace.json > folded-trace.out

# Convert the folded trace into a flame graph
flamegraph.pl --title "Gas Profile for 0x7405fc5e254352350bebcadc1392bd06f158aa88e9fb01733389621a29db5f08" --width 1920 --countname folded-trace.out > flame.svg
```
## Flags

```bash
      --file string           filename to read and hash
  -h, --help                  help for fold-trace
      --metric string         metric name for analysis: gas, count, actualgas (default "gas")
      --root-context string   name for top most initial context (default "root context")
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

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
