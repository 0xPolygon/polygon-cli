This command will reach the rollup manager contract and retrieve detailed information from a specific rollup.

Below is an example of how to use it

```bash
polycli cdk rollup dump
    --rpc-url https://sepolia.drpc.org
    --rollup-manager-address bali
    --rollup-id 1

polycli cdk rollup dump
    --rpc-url https://sepolia.drpc.org
    --rollup-manager-address bali
    --rollup-chain-id 2440

polycli cdk rollup dump
    --rpc-url https://sepolia.drpc.org
    --rollup-manager-address bali
    --rollup-address 0x89ba0ed947a88fe43c22ae305c0713ec8a7eb361
```
