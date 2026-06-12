`polycli sendtx` reads pre-signed raw transactions from a file and sends them to a JSON-RPC endpoint using batch `eth_sendRawTransaction` requests.

The command is designed for high-throughput transaction injection. It reads the input file, groups transactions into batches, and sends them concurrently via HTTP POST.

## Usage

```bash
polycli sendtx --file txs.txt --rpc-url http://localhost:8545
```

## Input File Format

One hex-encoded raw transaction per line (with or without `0x` prefix):

```
0x02f86f8301388280...
0x02f86f8301388280...
```

Empty lines are skipped.
