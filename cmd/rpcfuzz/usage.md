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

## Running RPC Fuzz Tests

After setting up your RPC endpoint and funding an account, you can run the RPC fuzz tests using various output formats. The tool supports streaming output that follows Unix philosophy - results are sent to stdout and you control data persistence through shell redirection.

### Output Format Examples

All commands use the same core parameters but produce different output formats:

#### Compact Format (Default)
Real-time colored console output with pass/fail indicators:
```bash
polycli rpcfuzz --rpc-url http://localhost:8545 --private-key <YOUR_PRIVATE_KEY> --namespaces eth,web3,net --compact > results.txt
```

#### CSV Format
Structured CSV with headers for data analysis:
```bash
polycli rpcfuzz --rpc-url http://localhost:8545 --private-key <YOUR_PRIVATE_KEY> --namespaces eth,web3,net --csv > results.csv
```

#### JSON Format
Streaming JSON with detailed test execution data:
```bash
polycli rpcfuzz --rpc-url http://localhost:8545 --private-key <YOUR_PRIVATE_KEY> --namespaces eth,web3,net --json > results.json
```

#### HTML Format
Complete styled HTML document for browser viewing:
```bash
polycli rpcfuzz --rpc-url http://localhost:8545 --private-key <YOUR_PRIVATE_KEY> --namespaces eth,web3,net --html > results.html
```

#### Markdown Format
Formatted Markdown table with emoji indicators:
```bash
polycli rpcfuzz --rpc-url http://localhost:8545 --private-key <YOUR_PRIVATE_KEY> --namespaces eth,web3,net --md > results.md
```

### Command Options

- `--rpc-url`: The JSON RPC endpoint URL
- `--private-key`: Private key for account with funds for testing
- `--namespaces`: Comma-separated list of RPC method namespaces to test (e.g., `eth,web3,net`)
- Output format flags: `--compact`, `--csv`, `--json`, `--html`, `--md` (mutually exclusive)

### Example with Error Suppression

To capture only test results without debug output:
```bash
polycli rpcfuzz --rpc-url <RPC_URL> --private-key <PRIVATE_KEY> --namespaces eth,web3,net --json 2>/dev/null > clean_results.json
```

### Links

- https://ethereum.github.io/execution-apis/api-documentation/
- https://ethereum.org/en/developers/docs/apis/json-rpc/
- https://json-schema.org/
- https://www.liquid-technologies.com/online-json-to-schema-converter
