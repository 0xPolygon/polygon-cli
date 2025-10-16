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

## Load env variables from kurtosis env

```bash
echo l1_rpc_url=\""http://"$(kurtosis port print cdk el-1-geth-lighthouse rpc)\"

echo l2_rpc_url=\"$(kurtosis port print cdk op-el-1-op-geth-op-node-001 rpc)\"
```

## Prepare env variables

```bash
pre_funded_private_key="0x12d7de8621a77640c9241b2595ba78ce443d05e94090365ab3bb5e19df82c625"
pre_funded_address="0xE34aaF64b29273B7D567FCFc40544c014EEe9970"

l1_rpc_url="http://127.0.0.1:32804"
l2_rpc_url="http://127.0.0.1:32814"
```

## Sanity check for blocks

```bash
cast block-number --rpc-url "$l1_rpc_url"
cast block-number --rpc-url "$l2_rpc_url"
```

```bash
1029
1694
```

## Check balance on L1, A and B

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$pre_funded_address"
cast balance --rpc-url "$l2_rpc_url" --ether "$pre_funded_address"
```

```bash
1999999.899935270835960337
100300.000000000000000000
```

## Deploy ERC20 smart contract to L1

```bash
cd ./contracts/src/tokens
cast send --private-key "$pre_funded_private_key" --rpc-url "$l1_rpc_url" --create $(forge build ERC20.sol --json | jq -r ".contracts.\"/Users/thiago/github.com/0xPolygon/polygon-cli/contracts/ERC20.sol\".ERC20[0].contract.evm.bytecode.object")
cd -
```

```bash
blockHash            0xfb70aeb59707f5806e393f0ac12693d37b67092e25e6fc387b3cbbbbcb54499e
blockNumber          1058
contractAddress      0x62bf798EdaE1B7FDe524276864757cc424A5c3dD
cumulativeGasUsed    1270795
effectiveGasPrice    1000000007
from                 0xE34aaF64b29273B7D567FCFc40544c014EEe9970
gasUsed              1270795
logs                 [{"address":"0x62bf798edae1b7fde524276864757cc424a5c3dd","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x0000000000000000000000000000000000000000000000000000000000000000","0x000000000000000000000000e34aaf64b29273b7d567fcfc40544c014eee9970"],"data":"0x00000000000000000000000000000000000000000000d3c21bcecceda1000000","blockHash":"0xfb70aeb59707f5806e393f0ac12693d37b67092e25e6fc387b3cbbbbcb54499e","blockNumber":"0x422","blockTimestamp":"0x68edf964","transactionHash":"0xae7059755680615d60bec2f1529c7ddef0a81be33e80ba1976ee47540aacca50","transactionIndex":"0x0","logIndex":"0x0","removed":false}]
logsBloom            0x00000000000000000000000000000200000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000020000000000000000000800000000000000000000000010100000000000000000000000000000040000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000001000000000000000
root
status               1 (success)
transactionHash      0xae7059755680615d60bec2f1529c7ddef0a81be33e80ba1976ee47540aacca50
transactionIndex     0
type                 2
blobGasPrice
blobGasUsed
```

## Set contract address env var

```bash
l1_contract_addr="0x62bf798EdaE1B7FDe524276864757cc424A5c3dD"
```

## Check ERC20 Balance

```bash
cast call "$l1_contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l1_rpc_url"
```

```bash
1000000000000000000000000 [1e24]
```

## Fund ERC20 to accounts

```bash
polycli fund  --verbosity 700 --rate-limit 2000 --rpc-url "$l1_rpc_url" --private-key "$pre_funded_private_key" --seed "ephemeral_test" --token-address "$l1_contract_addr" --token-amount 1000 --number 5 --file "wallets-funded.json"
```

```bash
4:27AM TRC Starting logger in console mode
4:27AM INF Starting bulk funding wallets
4:27AM TRC Input parameters params={"ApproveAmount":1000000000000000000000,"ApproveSpender":"","FunderAddress":"","FundingAmountInWei":50000000000000000,"KeyFile":"","Multicall3Address":"","OutputFile":"wallets-funded.json","PrivateKey":"0x12d7de8621a77640c9241b2595ba78ce443d05e94090365ab3bb5e19df82c625","RateLimit":2000,"RpcUrl":"http://127.0.0.1:32804","Seed":"ephemeral_test","TokenAddress":"0x62bf798EdaE1B7FDe524276864757cc424A5c3dD","TokenAmount":1000,"UseHDDerivation":true,"WalletAddresses":null,"WalletsNumber":5}
4:27AM TRC Detected chain ID chainID=271828
4:27AM INF Generating wallets from seed numWallets=5 seed=ephemeral_test
4:27AM TRC New wallet generated from seed address=0x6705a7352c76b0c9d204b5ad9e9dc92e57c4d44c privateKey=ad7a39a4c7e10a4d55dd5307b10d42cb7afe8b43ba54e4fe70ea4e5efc67b958 seedWithIndex=ephemeral_test_0_20251014
4:27AM TRC New wallet generated from seed address=0x26ee261781092a1833a4f09210bacf5826c975d1 privateKey=8c7d33f8401b89973d43d3da1457c02fc6502f46578ddaf4263c8d2edaffc5cb seedWithIndex=ephemeral_test_1_20251014
4:27AM TRC New wallet generated from seed address=0x002b6fada0efeeed24b1b7197c3785426c6bdd75 privateKey=4e589278f20c280ebf570e215b6a416a8b169f0949a39ddf6186c88340667d01 seedWithIndex=ephemeral_test_2_20251014
4:27AM TRC New wallet generated from seed address=0x4330ba531b0a8415552f1fafa02ccf390e66400b privateKey=f69c2a4d2374d3cc818aaca25e054f314db5b6b3ff98aef7b8e1966c87470382 seedWithIndex=ephemeral_test_3_20251014
4:27AM TRC New wallet generated from seed address=0x80c70ba9ee3abd5f17282cb989730c5ee5d282e4 privateKey=99d057ffd5d15acf91766cc7df1ba568c06c8b14d2ea7c52e9683533387dae24 seedWithIndex=ephemeral_test_4_20251014
4:27AM INF Wallet(s) generated from seed count=5
4:27AM DBG checking if multicall3 is supported
4:27AM INF Wallets' address(es) and private key(s) saved to file fileName=wallets-funded.json
4:27AM INF multicall3 is supported and will be used to fund wallets address=0xe293A6b8F558422813499bb5C89B60adD8c54636
4:27AM DBG funding wallets with multicall3
4:27AM DBG multicall3 max accounts to fund per tx accsToFundPerTx=700
4:27AM INF transaction to approve ERC20 token spending by multicall3 sent done=5 of=5 txHash=0xc0ef0c47461cac76609cb3880facd27e0fbebc02c7762fada8b57247035c4937
4:27AM INF multicall3 transaction to fund accounts sent done=5 of=5 txHash=0x6afe1bf620b22b99f2483ed43e88846ae3deba945da57848d94838f10eba38c6
4:27AM INF all funding transactions sent, waiting for confirmation...
4:27AM INF transaction confirmed txHash=0xc0ef0c47461cac76609cb3880facd27e0fbebc02c7762fada8b57247035c4937
4:27AM INF transaction confirmed txHash=0x6afe1bf620b22b99f2483ed43e88846ae3deba945da57848d94838f10eba38c6
4:27AM INF Wallet(s) funded! 💸
4:27AM INF Total execution time: 6.613848833s
```

## Check wallets balances

```bash
jq -r '.[].Address' wallets-funded.json \
| while read -r addr; do
    bal=$(cast call "$l1_contract_addr" "balanceOf(address)(uint256)" "$addr" --rpc-url "$l1_rpc_url")
    printf '%s: %s\n' "$addr" "$bal"
  done
```

```bash
0x6705a7352c76b0c9D204b5AD9E9Dc92E57c4D44C: 1000
0x26ee261781092a1833A4f09210bACf5826C975D1: 1000
0x002B6fADA0EFEEED24B1b7197C3785426C6BDd75: 1000
0x4330ba531b0a8415552F1fAFa02cCF390e66400B: 1000
0x80C70Ba9Ee3Abd5F17282Cb989730c5ee5D282e4: 1000
```
