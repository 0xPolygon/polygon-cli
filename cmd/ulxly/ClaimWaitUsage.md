This command will connect to the bridge service to check a deposit is ready to be claimed.

[Here](https://github.com/0xPolygonHermez/zkevm-contracts/blob/c8659e6282340de7bdb8fdbf7924a9bd2996bc98/contracts/v2/PolygonZkEVMBridgeV2.sol#L433-L465) is a direct link to the source code as well.

In order to check if a deposit is ready to be claimed, you need to know deposit count. Usually this is in the event data of the transaction. Alternatively, you can usually directly attempt to see the pending deposits by querying the bridge API directly. In the case of Cardona, the bridge service is running here: https://bridge-api.cardona.zkevm-rpc.com

```bash
curl -s https://bridge-api.cardona.zkevm-rpc.com/bridges/0x3878Cff9d621064d393EEF92bF1e12A944c5ba84 | jq '.'
```

In the output of the above command, I can see a deposit that looks like this:
```json
{
  "leaf_type": 0,
  "orig_net": 0,
  "orig_addr": "0x0000000000000000000000000000000000000000",
  "amount": "123456",
  "dest_net": 0,
  "dest_addr": "0x3878Cff9d621064d393EEF92bF1e12A944c5ba84",
  "block_num": "9695587",
  "deposit_cnt": 9075,
  "network_id": 1,
  "tx_hash": "0x0294dae3cfb26881e5dde9f182531aa5be0818956d029d50e9872543f020df2e",
  "claim_tx_hash": "",
  "metadata": "0x",
  "ready_for_claim": true,
  "global_index": "9075"
}
```

If we want to check if a deposit is ready to be claimed, we can use a command like this:

```bash
polycli ulxly claim wait \
    --bridge-address 0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582 \
    --bridge-service-url https://bridge-api.cardona.zkevm-rpc.com \
    --private-key 0x32430699cd4f46ab2422f1df4ad6546811be20c9725544e99253a887e971f92b \
    --deposit-network 1 \
    --deposit-count 9075 \
    --rpc-url https://sepolia.drpc.org 
```
