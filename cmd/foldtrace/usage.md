This command is meant to take a transaction op code trace and convert it into a folded output that can be easily visualized with Flamegraph tools.

```bash
# First grab a trace from an RPC that supports the debug namespace
cast rpc --rpc-url http://127.0.0.1:18545 debug_traceTransaction 0x12f63f489213f5bd5b88fbfb12960b8248f61e2062a369ba41d8a3c96bb74d57 > trace.json

# Read the trace and use the `fold-trace` command and write the output
polycli fold-trace --metric actualgas < trace.json > folded-trace.out

# Convert the folded trace into a flame graph
flamegraph.pl --title "Gas Profile for 0x7405fc5e254352350bebcadc1392bd06f158aa88e9fb01733389621a29db5f08" --width 1920 --countname folded-trace.out > flame.svg
```