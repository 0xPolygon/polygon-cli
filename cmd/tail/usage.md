Tail full blocks from an RPC endpoint and emit each block as newline-delimited JSON.

By default this prints the latest 10 blocks and exits:

```bash
polycli tail --rpc-url http://127.0.0.1:8545
```

Tail the last 100 blocks and keep following new blocks:

```bash
polycli tail -n 100 --follow --rpc-url http://127.0.0.1:8545
```
