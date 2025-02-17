This command will check the account current nonce against the max nonce found in the pool. In case of a nonce gap is found, txs will be sent to fill those gaps.

To fix a nonce gap, we can use a command like this:

```bash
polycli account nonce fix-gap \
    --rpc-url https://sepolia.drpc.org
    --private-key 0x32430699cd4f46ab2422f1df4ad6546811be20c9725544e99253a887e971f92b
```
