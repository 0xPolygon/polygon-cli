Recover public key from block

```bash
polycli ecrecover -r https://polygon-mumbai-bor.publicnode.com -b 45200775
> Recovering signer from block #45200775
> 0x5082F249cDb2f2c1eE035E4f423c46EA2daB3ab1

polycli ecrecover -r https://polygon-rpc.com
> Recovering signer from block #52888893
> 0xeEDBa2484aAF940f37cd3CD21a5D7C4A7DAfbfC0

polycli ecrecover -f block-52888893.json
> Recovering signer from block #52888893
> 0xeEDBa2484aAF940f37cd3CD21a5D7C4A7DAfbfC0

cat block-52888893.json | go run main.go ecrecover
> Recovering signer from block #52888893
> 0xeEDBa2484aAF940f37cd3CD21a5D7C4A7DAfbfC0
```

JSON Data passed in follows object definition [here](https://www.quicknode.com/docs/ethereum/eth_getBlockByNumber)
