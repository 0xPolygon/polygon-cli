This is a simple tool to avoid typing JSON on the command line while making RPC calls. The implementation is generic and this is meant to be a complete generic RPC tool.

```bash
$ polycli rpc https://polygon-rpc.com eth_blockNumber
$ polycli rpc https://polygon-rpc.com eth_getBlockByNumber 0x1e99576 true
$ polycli rpc http://127.0.0.1:8541 eth_getBlockByNumber 0x10E true
$ polycli rpc http://127.0.0.1:8541 eth_getBlockByHash 0x15206ab0a5b408214127f5c445a86b7cfe6ae48fdcd9172b14e013dae7a7f470 true
$ polycli rpc http://127.0.0.1:8541 eth_getTransactionByHash 0x97c070cb07bfac783ca73f08fb5999ae1ab509bf644197ef4a2c4e4f4a3c1516
```
