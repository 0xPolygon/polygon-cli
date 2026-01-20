# Aggkit Bridge Service Manual Test

## Create account

```bash
cast wallet new-mnemonic
```

```bash
Generating mnemonic from provided entropy...
Successfully generated a new mnemonic.
Phrase:
plastic cram delay outdoor metal kit carry radar vital retreat embark happy

Accounts:
- Account 0:
Address:     0xced253B29527D62a1880b95C23F256CE78a73c06
Private key: 0x516d5e1c8f7e1da24379041b758b5d16fd066c8a8856791d3c5c0c79a81bad59
```

## Prepare env variables

```bash
pre_funded_private_key="0x12d7de8621a77640c9241b2595ba78ce443d05e94090365ab3bb5e19df82c625"
pre_funded_address="0xE34aaF64b29273B7D567FCFc40544c014EEe9970"

funder_private_key="0x516d5e1c8f7e1da24379041b758b5d16fd066c8a8856791d3c5c0c79a81bad59"
funder_address="0xced253B29527D62a1880b95C23F256CE78a73c06"

l1_rpc_url="http://$(kurtosis port print cdk el-1-geth-lighthouse rpc)"
l2_rpc_url="$(kurtosis port print cdk op-el-1-op-geth-op-node-001 rpc)"
```

## Sanity check for blocks

```bash
cast block-number --rpc-url "$l1_rpc_url"
cast block-number --rpc-url "$l2_rpc_url"
```

```bash
259
140
```

## Check balance on L1 and L2

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$pre_funded_address"
cast balance --rpc-url "$l2_rpc_url" --ether "$pre_funded_address"
cast balance --rpc-url "$l1_rpc_url" --ether "$funder_address"
cast balance --rpc-url "$l2_rpc_url" --ether "$funder_address"
```

```bash
1999999.899962054234851121
100300.000000000000000000
0.000000000000000000
0.000000000000000000
```

## Fund account on L1

```bash
polycli fund --addresses "$funder_address" --eth-amount 1000000000000000000 --private-key "$pre_funded_private_key" --rpc-url "$l1_rpc_url"
```

```bash
10:57AM INF Starting bulk funding wallets
10:57AM INF Using addresses provided by the user
10:57AM INF multicall3 is supported and will be used to fund wallets address=0x62bf798EdaE1B7FDe524276864757cc424A5c3dD
10:57AM INF multicall3 transaction to fund accounts sent done=1 of=1 txHash=0xbcc4356af43a1c418d676a2aa00b60852599913cecbdeadd172b2d383d292a36
10:57AM INF all funding transactions sent, waiting for confirmation...
10:57AM INF transaction confirmed txHash=0xbcc4356af43a1c418d676a2aa00b60852599913cecbdeadd172b2d383d292a36
10:57AM INF Wallet(s) funded! 💸
10:57AM INF Total execution time: 4.921481542s
```

## Fund account on L2

```bash
polycli fund --addresses "$funder_address" --eth-amount 1000000000000000000 --private-key "$pre_funded_private_key" --rpc-url "$l2_rpc_url"
```

```bash
10:58AM INF Starting bulk funding wallets
10:58AM INF Using addresses provided by the user
10:58AM INF multicall3 is supported and will be used to fund wallets address=0xcA11bde05977b3631167028862bE2a173976CA11
10:58AM INF multicall3 transaction to fund accounts sent done=1 of=1 txHash=0x08a708dd4a9c2e6a795bda456deae29a5f638bc30cdc2d1c75db7fdf679a9d1f
10:58AM INF all funding transactions sent, waiting for confirmation...
10:58AM INF transaction confirmed txHash=0x08a708dd4a9c2e6a795bda456deae29a5f638bc30cdc2d1c75db7fdf679a9d1f
10:58AM INF Wallet(s) funded! 💸
10:58AM INF Total execution time: 2.282900083s
```

## Check balance on L1 and L2 after fund

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$pre_funded_address"
cast balance --rpc-url "$l2_rpc_url" --ether "$pre_funded_address"
cast balance --rpc-url "$l1_rpc_url" --ether "$funder_address"
cast balance --rpc-url "$l2_rpc_url" --ether "$funder_address"
```

```bash
1999998.899962054224359841
100298.999974038322447806
1.000000000000000000
1.000000000000000000
```

## Deploy ERC20 smart contract to L1

```bash
cd ./contracts/src/tokens
cast send --private-key "$pre_funded_private_key" --rpc-url "$l1_rpc_url" --create $(forge build ERC20.sol --json | jq -r ".contracts.\"/Users/thiago/github.com/0xPolygon/polygon-cli/contracts/ERC20.sol\".ERC20[0].contract.evm.bytecode.object")
cd -
```

```bash
blockHash            0xb98cd5153a007c32c3fe3ffc2d1d956e68e11be2a7a7cdae5f8914761e941ef5
blockNumber          304
contractAddress      0x9ceA3ee97f9eB1c39F3196060f24B7ED52bb7Ca3
cumulativeGasUsed    1270795
effectiveGasPrice    1000000007
from                 0xE34aaF64b29273B7D567FCFc40544c014EEe9970
gasUsed              1270795
logs                 [{"address":"0x9cea3ee97f9eb1c39f3196060f24b7ed52bb7ca3","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x0000000000000000000000000000000000000000000000000000000000000000","0x000000000000000000000000e34aaf64b29273b7d567fcfc40544c014eee9970"],"data":"0x00000000000000000000000000000000000000000000d3c21bcecceda1000000","blockHash":"0xb98cd5153a007c32c3fe3ffc2d1d956e68e11be2a7a7cdae5f8914761e941ef5","blockNumber":"0x130","blockTimestamp":"0x68f0fa11","transactionHash":"0xedef82b2c4ca2e00b31610c84571e760c3a81733f0e40e23f95d0ca102229b31","transactionIndex":"0x0","logIndex":"0x0","removed":false}]
logsBloom            0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000001000000000000000000000000000000000000008000000000000000000000000000000000000000000100000020000000000000000000800000000000000000000000010000000000000000000000000000000000000000000000000000000200000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000001000000000000000
root
status               1 (success)
transactionHash      0xedef82b2c4ca2e00b31610c84571e760c3a81733f0e40e23f95d0ca102229b31
transactionIndex     0
type                 2
blobGasPrice
blobGasUsed
```

## Deploy ERC20 smart contract to L2

```bash
cd ./contracts/src/tokens
cast send --private-key "$pre_funded_private_key" --rpc-url "$l2_rpc_url" --create $(forge build ERC20.sol --json | jq -r ".contracts.\"/Users/thiago/github.com/0xPolygon/polygon-cli/contracts/ERC20.sol\".ERC20[0].contract.evm.bytecode.object")
cd -
```

```bash
blockHash            0x561b9a5eb1175cadc62dacb8f2d8f774a089272a5cbfef795e7c5c5da1bfb02d
blockNumber          234
contractAddress      0x1f7ad7caA53e35b4f0D138dC5CBF91aC108a2674
cumulativeGasUsed    1316879
effectiveGasPrice    393291005
from                 0xE34aaF64b29273B7D567FCFc40544c014EEe9970
gasUsed              1270795
logs                 [{"address":"0x1f7ad7caa53e35b4f0d138dc5cbf91ac108a2674","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x0000000000000000000000000000000000000000000000000000000000000000","0x000000000000000000000000e34aaf64b29273b7d567fcfc40544c014eee9970"],"data":"0x00000000000000000000000000000000000000000000d3c21bcecceda1000000","blockHash":"0x561b9a5eb1175cadc62dacb8f2d8f774a089272a5cbfef795e7c5c5da1bfb02d","blockNumber":"0xea","transactionHash":"0x3c9883145b121d80c65d74a4b097d0f0d8b095eab981eed586c1e94b3240d963","transactionIndex":"0x1","logIndex":"0x0","removed":false}]
logsBloom            0x00000000000000000000000008000000000000000000000000000000000000000008000000000000000008000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000020000000000000000000800000000000000000000000010000000000000000000000000000000000100000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000001000000000000000
root
status               1 (success)
transactionHash      0x3c9883145b121d80c65d74a4b097d0f0d8b095eab981eed586c1e94b3240d963
transactionIndex     1
type                 2
blobGasPrice
blobGasUsed
l1BaseFeeScalar      1368
l1BlobBaseFee        1
l1BlobBaseFeeScalar  810949
l1Fee                3088
l1GasPrice           7
l1GasUsed            51248
```

## Set contract address env var

```bash
l1_contract_addr="0x9ceA3ee97f9eB1c39F3196060f24B7ED52bb7Ca3"
l2_contract_addr="0x1f7ad7caA53e35b4f0D138dC5CBF91aC108a2674"
```

## Check ERC20 Balance

```bash
cast call "$l1_contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l1_rpc_url"
cast call "$l2_contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l2_rpc_url"
cast call "$l1_contract_addr" "balanceOf(address)(uint256)" "$funder_address" --rpc-url "$l1_rpc_url"
cast call "$l2_contract_addr" "balanceOf(address)(uint256)" "$funder_address" --rpc-url "$l2_rpc_url"
```

```bash
1000000000000000000000000 [1e24]
1000000000000000000000000 [1e24]
0
0
```

## Mint ERC20 tokens to funder

```bash
cast send "$l1_contract_addr" "mint(uint256)" "1000000000000" --private-key "$funder_private_key" --rpc-url "$l1_rpc_url"
cast send "$l2_contract_addr" "mint(uint256)" "1000000000000" --private-key "$funder_private_key" --rpc-url "$l2_rpc_url"
```

## Check ERC20 Balance after mint

```bash
cast call "$l1_contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l1_rpc_url"
cast call "$l2_contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l2_rpc_url"
cast call "$l1_contract_addr" "balanceOf(address)(uint256)" "$funder_address" --rpc-url "$l1_rpc_url"
cast call "$l2_contract_addr" "balanceOf(address)(uint256)" "$funder_address" --rpc-url "$l2_rpc_url"
```

```bash
1000000000000000000000000 [1e24]
1000000000000000000000000 [1e24]
1000000000000 [1e12]
1000000000000 [1e12]
```

## Fund ERC20 to accounts on L1

```bash
polycli fund --rpc-url "$l1_rpc_url" --private-key "$funder_private_key" --seed "ephemeral_test_l1" --token-address "$l1_contract_addr" --token-amount 1000 --number 5 --file "wallets-funded-l1.json"
```

```bash
11:00AM INF Starting bulk funding wallets
11:00AM INF Generating wallets from seed numWallets=5 seed=ephemeral_test_l1
11:00AM INF Wallet(s) generated from seed count=5
11:00AM INF Wallets' address(es) and private key(s) saved to file fileName=wallets-funded-l1.json
11:00AM INF multicall3 is supported and will be used to fund wallets address=0xaEd7FE0a652395C4d8F9AbD038375b13e632BF85
11:00AM INF transaction to approve ERC20 token spending by multicall3 sent done=5 of=5 txHash=0xb25298b5a87cdced5d9e68b2670b4dbd5359527b2a664d637ea49e367ca80470
11:00AM INF multicall3 transaction to fund accounts sent done=5 of=5 txHash=0x21a81ac529acaab8ada53e72a1e523cd1291fe5e91875c5df73ede01992ff5b4
11:00AM INF all funding transactions sent, waiting for confirmation...
11:00AM INF transaction confirmed txHash=0xb25298b5a87cdced5d9e68b2670b4dbd5359527b2a664d637ea49e367ca80470
11:00AM INF transaction confirmed txHash=0x21a81ac529acaab8ada53e72a1e523cd1291fe5e91875c5df73ede01992ff5b4
11:00AM INF Wallet(s) funded! 💸
11:00AM INF Total execution time: 6.397935959s
```

## Fund ERC20 to accounts on L2

```bash
polycli fund --rpc-url "$l2_rpc_url" --private-key "$funder_private_key" --seed "ephemeral_test_l2" --token-address "$l2_contract_addr" --token-amount 1000 --number 5 --file "wallets-funded-l2.json"
```

```bash
11:00AM INF Starting bulk funding wallets
11:00AM INF Generating wallets from seed numWallets=5 seed=ephemeral_test_l2
11:00AM INF Wallet(s) generated from seed count=5
11:00AM INF Wallets' address(es) and private key(s) saved to file fileName=wallets-funded-l2.json
11:00AM INF multicall3 is supported and will be used to fund wallets address=0xcA11bde05977b3631167028862bE2a173976CA11
11:00AM INF transaction to approve ERC20 token spending by multicall3 sent done=5 of=5 txHash=0x5ee4b0ad173ef28bdeb34047cfbe01e645fc266e97001fcdc07b73e1b44b3e71
11:00AM INF multicall3 transaction to fund accounts sent done=5 of=5 txHash=0xa81c5305c4d7d08b4a0b340d6ce6efa904d0310d2c7f919e12d4b3f27df2357a
11:00AM INF all funding transactions sent, waiting for confirmation...
11:00AM INF transaction confirmed txHash=0x5ee4b0ad173ef28bdeb34047cfbe01e645fc266e97001fcdc07b73e1b44b3e71
11:00AM INF transaction confirmed txHash=0xa81c5305c4d7d08b4a0b340d6ce6efa904d0310d2c7f919e12d4b3f27df2357a
11:00AM INF Wallet(s) funded! 💸
11:00AM INF Total execution time: 3.990132417s
```

## Check wallets balances

```bash
jq -r '.[].Address' wallets-funded-l1.json \
| while read -r addr; do
    bal=$(cast call "$l1_contract_addr" "balanceOf(address)(uint256)" "$addr" --rpc-url "$l1_rpc_url")
    printf '%s: %s\n' "$addr" "$bal"
  done
echo ""
jq -r '.[].Address' wallets-funded-l2.json \
| while read -r addr; do
    bal=$(cast call "$l2_contract_addr" "balanceOf(address)(uint256)" "$addr" --rpc-url "$l2_rpc_url")
    printf '%s: %s\n' "$addr" "$bal"
  done
```

```bash
0xA663Fc82FF0e336014f2e51265845DeB90FDC67E: 1000
0xA7FEBaBde379b056E82fC3780EAFB2564346e110: 1000
0xb3bF95FF2598FD6B7713216e631Fa66A11cc59AE: 1000
0xD29dbA4Eb80A5514Ca5e60B7017461C63Ab66671: 1000
0x432456de21797DFb3A21DFaa555143cEAc9c106B: 1000

0x0EDD8143d2519326eC89FCF3d20C1b67d6793AB8: 1000
0x7c387d821aD71B5B34C456fa4907Ba3065372290: 1000
0xAfD4Da29Ac08C205d0c514704F80bDF762BaC6A3: 1000
0x05725B2f384e7FA1c92C15e95d575070a6202B02: 1000
0x232D5A236D1810928FEB0D9F7b1AFa178598e34a: 1000
```
