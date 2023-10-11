This command will run a series of RPC calls against a given JSON RPC endpoint. The idea is to be able to check for various features and function to see if the RPC generally conforms to typical geth standards for the RPC.

Some setup might be needed depending on how you're testing. We'll demonstrate with geth.

In order to quickly test this, you can run geth in dev mode.

```bash
$ geth \
    --dev \
    --dev.period 5 \
    --http \
    --http.addr localhost \
    --http.port 8545 \
    --http.api 'admin,debug,web3,eth,txpool,personal,miner,net' \
    --verbosity 5 \
    --rpc.gascap 50000000 \
    --rpc.txfeecap 0 \
    --miner.gaslimit 10 \
    --miner.gasprice 1 \
    --gpo.blocks 1 \
    --gpo.percentile 1 \
    --gpo.maxprice 10 \
    --gpo.ignoreprice 2 \
    --dev.gaslimit 50000000
```

If we wanted to use erigon for testing, we could do something like this as well.

```bash
$ erigon \
    --chain dev \
    --dev.period 5 \
    --http \
    --http.addr localhost \
    --http.port 8545 \
    --http.api 'admin,debug,web3,eth,txpool,clique,net' \
    --verbosity 5 \
    --rpc.gascap 50000000 \
    --miner.gaslimit 10 \
    --gpo.blocks 1 \
    --gpo.percentile 1 \
    --mine
```

Once your Eth client is running and the RPC is functional, you'll need to transfer some amount of ether to a known account that ca be used for testing.

```
$ cast send \
    --from "$(cast rpc --rpc-url localhost:8545 eth_coinbase | jq -r '.')" \
    --rpc-url localhost:8545 \
    --unlocked \
    --value 100ether \
    --json \
    0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6 | jq
```

Then we might want to deploy some test smart contracts. For the purposes of testing we'll our ERC20 contract.

```bash
$ cast send \
    --from 0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6 \
    --private-key 0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa \
    --rpc-url localhost:8545 \
    --json \
    --create \
    "$(cat ./contracts/tokens/ERC20/ERC20.bin)" | jq
```

Once this has been completed this will be the address of the contract: `0x6fda56c57b0acadb96ed5624ac500c0429d59429`.

```bash
$  docker run -v $PWD/contracts:/contracts ethereum/solc:stable --storage-layout /contracts/tokens/ERC20/ERC20.sol
```

### Links

- https://ethereum.github.io/execution-apis/api-documentation/
- https://ethereum.org/en/developers/docs/apis/json-rpc/
- https://json-schema.org/
- https://www.liquid-technologies.com/online-json-to-schema-converter
