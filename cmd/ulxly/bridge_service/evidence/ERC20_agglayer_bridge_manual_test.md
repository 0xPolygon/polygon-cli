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

echo l2_a_rpc_url=\"$(kurtosis port print cdk op-el-1-op-geth-op-node-001 rpc)\"
echo l2_a_bridge_url=\"$(kurtosis port print cdk aggkit-001-bridge rest)\"

echo l2_b_rpc_url=\"$(kurtosis port print cdk op-el-1-op-geth-op-node-002 rpc)\"
echo l2_b_bridge_url=\"$(kurtosis port print cdk aggkit-002-bridge rest)\"
```

## Prepare env variables

```bash
pre_funded_private_key="0x12d7de8621a77640c9241b2595ba78ce443d05e94090365ab3bb5e19df82c625"
pre_funded_address="0xE34aaF64b29273B7D567FCFc40544c014EEe9970"

eth_address="0xced253B29527D62a1880b95C23F256CE78a73c06"
private_key="0x516d5e1c8f7e1da24379041b758b5d16fd066c8a8856791d3c5c0c79a81bad59"

bridge_address="0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6"

l1_rpc_url="http://127.0.0.1:32938"
l2_a_rpc_url="http://127.0.0.1:32948"
l2_a_bridge_url="http://127.0.0.1:32967"
l2_b_rpc_url="http://127.0.0.1:32981"
l2_b_bridge_url="http://127.0.0.1:32994"
```

## Sanity check for blocks

```bash
cast block-number --rpc-url "$l1_rpc_url"
cast block-number --rpc-url "$l2_a_rpc_url"
cast block-number --rpc-url "$l2_b_rpc_url"
```

```bash
910
1488
366
```

## Get Bridge network ID

```bash
cast call --rpc-url "$l2_a_rpc_url" "$bridge_address" "networkID()(uint32)"
cast call --rpc-url "$l2_b_rpc_url" "$bridge_address" "networkID()(uint32)"
```

```bash
1
2
```

## Fund account on L1, A and B

```bash
cast send "$eth_address" --value 1ether --private-key "$pre_funded_private_key" --rpc-url "$l1_rpc_url"
cast send "$eth_address" --value 1ether --private-key "$pre_funded_private_key" --rpc-url "$l2_a_rpc_url"
cast send "$eth_address" --value 1ether --private-key "$pre_funded_private_key" --rpc-url "$l2_b_rpc_url"
```

## Check balance on L1, A and B

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
```

```bash
1.000000000000000000
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
blockHash            0x2b7ed39f3feae5cba381482bfa847123ba823af600bf88914501d84629203480
blockNumber          916
contractAddress      0x9c85cd40541D67670aaC4D8249a55668896A6BD3
cumulativeGasUsed    1270795
effectiveGasPrice    1000000007
from                 0xE34aaF64b29273B7D567FCFc40544c014EEe9970
gasUsed              1270795
logs                 [{"address":"0x9c85cd40541d67670aac4d8249a55668896a6bd3","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x0000000000000000000000000000000000000000000000000000000000000000","0x000000000000000000000000e34aaf64b29273b7d567fcfc40544c014eee9970"],"data":"0x00000000000000000000000000000000000000000000d3c21bcecceda1000000","blockHash":"0x2b7ed39f3feae5cba381482bfa847123ba823af600bf88914501d84629203480","blockNumber":"0x394","blockTimestamp":"0x68dafb53","transactionHash":"0xb54e5e4b4fe4f3b4bc6144488dc19d8fe8133c7b46a9e3324ef05ab1198c6301","transactionIndex":"0x0","logIndex":"0x0","removed":false}]
logsBloom            0x00000800000000000000000000000800000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000020000000000000000000800000000000000000000000010004000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000001000000000000000
root
status               1 (success)
transactionHash      0xb54e5e4b4fe4f3b4bc6144488dc19d8fe8133c7b46a9e3324ef05ab1198c6301
transactionIndex     0
type                 2
blobGasPrice
blobGasUsed
```

## Set contract address env var

```bash
contract_addr="0x9c85cd40541D67670aaC4D8249a55668896A6BD3"
```

## Check current balances

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
```

```bash
0.000000000000000000
0.000000000000000000
0.000000000000000000
```

## Check ERC20 Balance

```bash
cast call "$contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l1_rpc_url"
cast call "$contract_addr" "balanceOf(address)(uint256)" "$eth_address" --rpc-url "$l1_rpc_url"
```

```bash
1000000000000000000000000 [1e24]
0
```

## Bridge from L1 to L2 - Network A

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$pre_funded_private_key" \
    --destination-network 1 \
    --token-address "$contract_addr" \
    --value 100000000000000000 \
    --destination-address "$eth_address" \
    --rpc-url "$l1_rpc_url"
```

```bash
6:40PM INF approving bridge contract to spend tokens on behalf of user amount=100000000000000000 bridgeAddress=0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6 tokenAddress=0x9c85cd40541D67670aaC4D8249a55668896A6BD3 userAddress=0xE34aaF64b29273B7D567FCFc40544c014EEe9970
6:40PM INF approveTxn: 0x5956129c482dc36c10ad11f33a438cb78d25a0dcef1a6503b8d5a6c6b56a819b
6:40PM INF transaction successful txHash=0x5956129c482dc36c10ad11f33a438cb78d25a0dcef1a6503b8d5a6c6b56a819b
6:40PM INF bridgeTxn: 0xb655523ddf8c3687c2fce01100b3ce45354731a9199b14fc0ef79be26c7caa66
6:40PM INF transaction successful txHash=0xb655523ddf8c3687c2fce01100b3ce45354731a9199b14fc0ef79be26c7caa66
6:40PM INF Bridge deposit count parsed from logs depositCount=2
```

## Check deposit is already indexed by bridge service

```bash
curl -s "$l2_a_bridge_url/bridge/v1/bridges?network_id=0&from_address=$pre_funded_address" | jq -M '.bridges[] | select(.destination_network == 1)'
```

```bash
{
  "block_num": 1108,
  "block_pos": 2,
  "from_address": "0xE34aaF64b29273B7D567FCFc40544c014EEe9970",
  "tx_hash": "0xb655523ddf8c3687c2fce01100b3ce45354731a9199b14fc0ef79be26c7caa66",
  "calldata": "0xcd5865790000000000000000000000000000000000000000000000000000000000000001000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c06000000000000000000000000000000000000000000000000016345785d8a00000000000000000000000000009c85cd40541d67670aac4d8249a55668896a6bd3000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000",
  "block_timestamp": 1759182035,
  "leaf_type": 0,
  "origin_network": 0,
  "origin_address": "0x9c85cd40541D67670aaC4D8249a55668896A6BD3",
  "destination_network": 1,
  "destination_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "amount": "100000000000000000",
  "metadata": "0x000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000001200000000000000000000000000000000000000000000000000000000000000074d79546f6b656e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034d544b0000000000000000000000000000000000000000000000000000000000",
  "deposit_count": 2,
  "is_native_token": false,
  "bridge_hash": "0x131d06b2685ddc9e6a262f719af4cd80f8e0dc1875cd9a3dd44f50015917d223"
}
```

## get wrapped token address on L2-A

```bash
cast call "$bridge_address" "getTokenWrappedAddress(uint32,address)(address)" 0 "$contract_addr" --rpc-url "$l2_a_rpc_url"
```

```bash
0x3EbC805751b5e1654f1477e0958b21F6B75b111d
```

## Set contract address env var for L2-A

```bash
l2_a_contract_addr="0x3EbC805751b5e1654f1477e0958b21F6B75b111d"
```

## Check balances to ensure auto-claimer claimed the deposit

```bash
cast call "$contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l1_rpc_url"
cast call "$contract_addr" "balanceOf(address)(uint256)" "$eth_address" --rpc-url "$l1_rpc_url"

cast call "$l2_a_contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l2_a_rpc_url"
cast call "$l2_a_contract_addr" "balanceOf(address)(uint256)" "$eth_address" --rpc-url "$l2_a_rpc_url"
```

```bash
999999900000000000000000 [9.999e23]
0
0
100000000000000000 [1e17]
```

## Bridge from L2-A to L2-B

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 2 \
    --token-address "$l2_a_contract_addr" \
    --value 100000000000 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_a_rpc_url"
```

```bash
6:54PM INF approving bridge contract to spend tokens on behalf of user amount=100000000000 bridgeAddress=0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6 tokenAddress=0x3EbC805751b5e1654f1477e0958b21F6B75b111d userAddress=0xced253B29527D62a1880b95C23F256CE78a73c06
6:54PM INF approveTxn: 0xa6cc320aadfc639fc5c4f90b0c7df3680a091ee2002be4efa0dafe46c44917b8
6:54PM INF transaction successful txHash=0xa6cc320aadfc639fc5c4f90b0c7df3680a091ee2002be4efa0dafe46c44917b8
6:54PM INF bridgeTxn: 0xa487771f94a3636f898ef821ee8710b0e61fc0a49c432a63aed726782c2e3ae5
6:54PM INF transaction successful txHash=0xa487771f94a3636f898ef821ee8710b0e61fc0a49c432a63aed726782c2e3ae5
6:54PM INF Bridge deposit count parsed from logs depositCount=0
```

## Check deposit is already indexed by bridge service on L2-A

```bash
curl -s "$l2_a_bridge_url/bridge/v1/bridges?network_id=1&from_address=$eth_address" | jq -M '.bridges[] | select(.destination_network == 2)'
```

```bash
{
  "block_num": 2704,
  "block_pos": 1,
  "from_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "tx_hash": "0xa487771f94a3636f898ef821ee8710b0e61fc0a49c432a63aed726782c2e3ae5",
  "calldata": "0xcd5865790000000000000000000000000000000000000000000000000000000000000002000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c06000000000000000000000000000000000000000000000000000000174876e8000000000000000000000000003ebc805751b5e1654f1477e0958b21f6b75b111d000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000",
  "block_timestamp": 1759182855,
  "leaf_type": 0,
  "origin_network": 0,
  "origin_address": "0x9c85cd40541D67670aaC4D8249a55668896A6BD3",
  "destination_network": 2,
  "destination_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "amount": "100000000000",
  "metadata": "0x000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000001200000000000000000000000000000000000000000000000000000000000000074d79546f6b656e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034d544b0000000000000000000000000000000000000000000000000000000000",
  "deposit_count": 0,
  "is_native_token": false,
  "bridge_hash": "0xfd2202f1ba8ac92695c686ce19306eaa72b8d1c4a4a81405f8ec04c1bbabc2c8"
}
```

```bash
polycli ulxly claim asset \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_a_bridge_url" \
    --private-key "$private_key" \
    --deposit-network 1 \
    --deposit-count 0 \
    --rpc-url "$l2_b_rpc_url"
```

```bash
7:31PM INF No destination address specified. Using private key's address destAddress=0xced253B29527D62a1880b95C23F256CE78a73c06
7:31PM INF The deposit is ready to be claimed
7:31PM INF claimTxn: 0xce81302627107ceb662f1581b3d65615a8213340218953b7597ce712100eab11
7:31PM INF transaction successful txHash=0xce81302627107ceb662f1581b3d65615a8213340218953b7597ce712100eab11
```

## get wrapped token address on L2-B

```bash
cast call "$bridge_address" "getTokenWrappedAddress(uint32,address)(address)" 0 "$contract_addr" --rpc-url "$l2_b_rpc_url"
```

```bash
0x3EbC805751b5e1654f1477e0958b21F6B75b111d
```

## Set contract address env var for L2-B

```bash
l2_b_contract_addr="0x3EbC805751b5e1654f1477e0958b21F6B75b111d"
```

## Check balances to ensure the bridge from L2-A to L2-B worked

```bash
cast call "$contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l1_rpc_url"
cast call "$contract_addr" "balanceOf(address)(uint256)" "$eth_address" --rpc-url "$l1_rpc_url"

cast call "$l2_a_contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l2_a_rpc_url"
cast call "$l2_a_contract_addr" "balanceOf(address)(uint256)" "$eth_address" --rpc-url "$l2_a_rpc_url"

cast call "$l2_b_contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l2_b_rpc_url"
cast call "$l2_b_contract_addr" "balanceOf(address)(uint256)" "$eth_address" --rpc-url "$l2_b_rpc_url"
```

```bash
999999900000000000000000 [9.999e23]
0
0
99999900000000000 [9.999e16]
0
100000000000 [1e11]
```

## Bridge from A and B to L1

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 0 \
    --token-address "$l2_a_contract_addr" \
    --value 99999900000000000 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_a_rpc_url"
```

```bash
7:41PM INF approving bridge contract to spend tokens on behalf of user amount=99999900000000000 bridgeAddress=0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6 tokenAddress=0x3EbC805751b5e1654f1477e0958b21F6B75b111d userAddress=0xced253B29527D62a1880b95C23F256CE78a73c06
7:41PM INF approveTxn: 0xcfbcd1a903d080cc525a0067dd32ec36c68762204e3c9c097cc9cfa8059ca08c
7:41PM INF transaction successful txHash=0xcfbcd1a903d080cc525a0067dd32ec36c68762204e3c9c097cc9cfa8059ca08c
7:41PM INF bridgeTxn: 0xee1ccc2a22ac17fa3545a0baefd11b92c2b667c69ae47aac5dd96b7850c0fec9
7:41PM INF transaction successful txHash=0xee1ccc2a22ac17fa3545a0baefd11b92c2b667c69ae47aac5dd96b7850c0fec9
7:41PM INF Bridge deposit count parsed from logs depositCount=1
```

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 0 \
    --token-address "$l2_b_contract_addr" \
    --value 100000000000 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_b_rpc_url"
```

```bash
7:42PM INF approving bridge contract to spend tokens on behalf of user amount=100000000000 bridgeAddress=0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6 tokenAddress=0x3EbC805751b5e1654f1477e0958b21F6B75b111d userAddress=0xced253B29527D62a1880b95C23F256CE78a73c06
7:42PM INF approveTxn: 0x9c3c5292e7692befd5e16cc316ecf87f7b0026b451395d023ef7d270684d915b
7:42PM INF transaction successful txHash=0x9c3c5292e7692befd5e16cc316ecf87f7b0026b451395d023ef7d270684d915b
7:42PM INF bridgeTxn: 0xdf8ab94c8242b7412ef88ba6217c07263a08ee74b6a731a4af6bf294db2210cc
7:42PM INF transaction successful txHash=0xdf8ab94c8242b7412ef88ba6217c07263a08ee74b6a731a4af6bf294db2210cc
7:42PM INF Bridge deposit count parsed from logs depositCount=0
```

## Check balances to ensure deposit was made from A and B to L1

```bash
cast call "$contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l1_rpc_url"
cast call "$contract_addr" "balanceOf(address)(uint256)" "$eth_address" --rpc-url "$l1_rpc_url"

cast call "$l2_a_contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l2_a_rpc_url"
cast call "$l2_a_contract_addr" "balanceOf(address)(uint256)" "$eth_address" --rpc-url "$l2_a_rpc_url"

cast call "$l2_b_contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l2_b_rpc_url"
cast call "$l2_b_contract_addr" "balanceOf(address)(uint256)" "$eth_address" --rpc-url "$l2_b_rpc_url"
```

```bash
999999900000000000000000 [9.999e23]
0
0
0
0
0
```

## Check deposit is already indexed by bridge service on A, B and C to be claimed on L1

```bash
curl -s "$l2_a_bridge_url/bridge/v1/bridges?network_id=1&from_address=$eth_address" | jq -M '.bridges[] | select(.destination_network == 0)'
curl -s "$l2_b_bridge_url/bridge/v1/bridges?network_id=2&from_address=$eth_address" | jq -M '.bridges[] | select(.destination_network == 0)'
```

```bash
{
  "block_num": 5557,
  "block_pos": 1,
  "from_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "tx_hash": "0xee1ccc2a22ac17fa3545a0baefd11b92c2b667c69ae47aac5dd96b7850c0fec9",
  "calldata": "0xcd5865790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c0600000000000000000000000000000000000000000000000001634561151318000000000000000000000000003ebc805751b5e1654f1477e0958b21f6b75b111d000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000",
  "block_timestamp": 1759185708,
  "leaf_type": 0,
  "origin_network": 0,
  "origin_address": "0x9c85cd40541D67670aaC4D8249a55668896A6BD3",
  "destination_network": 0,
  "destination_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "amount": "99999900000000000",
  "metadata": "0x000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000001200000000000000000000000000000000000000000000000000000000000000074d79546f6b656e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034d544b0000000000000000000000000000000000000000000000000000000000",
  "deposit_count": 1,
  "is_native_token": false,
  "bridge_hash": "0xf16b0de217cb801cdb50d0b59d892b5bc2cc647af294362b374425cf7b00e7b2"
}
{
  "block_num": 4465,
  "block_pos": 1,
  "from_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "tx_hash": "0xdf8ab94c8242b7412ef88ba6217c07263a08ee74b6a731a4af6bf294db2210cc",
  "calldata": "0xcd5865790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c06000000000000000000000000000000000000000000000000000000174876e8000000000000000000000000003ebc805751b5e1654f1477e0958b21f6b75b111d000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000",
  "block_timestamp": 1759185738,
  "leaf_type": 0,
  "origin_network": 0,
  "origin_address": "0x9c85cd40541D67670aaC4D8249a55668896A6BD3",
  "destination_network": 0,
  "destination_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "amount": "100000000000",
  "metadata": "0x000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000001200000000000000000000000000000000000000000000000000000000000000074d79546f6b656e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034d544b0000000000000000000000000000000000000000000000000000000000",
  "deposit_count": 0,
  "is_native_token": false,
  "bridge_hash": "0x53c93d8690de1d0a3280e6647970f8dc6d50352fee1f7f40194775e9acc66a6f"
}
```

## Claim deposit made from A and B to L1

```bash
polycli ulxly claim asset \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_a_bridge_url" \
    --private-key "$pre_funded_private_key" \
    --destination-address "$eth_address" \
    --deposit-network 1 \
    --deposit-count 1 \
    --rpc-url "$l1_rpc_url"
```

```bash
7:52PM INF The deposit is ready to be claimed
7:52PM INF claimTxn: 0xf86ad22949acf57a8f7f6f6c7793777710b072062cbb79e1141a5fc53795be26
7:52PM INF transaction successful txHash=0xf86ad22949acf57a8f7f6f6c7793777710b072062cbb79e1141a5fc53795be26
```

```bash
polycli ulxly claim asset \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_b_bridge_url" \
    --private-key "$pre_funded_private_key" \
    --destination-address "$eth_address" \
    --deposit-network 2 \
    --deposit-count 0 \
    --rpc-url "$l1_rpc_url"
```

```bash
7:52PM INF The deposit is ready to be claimed
7:52PM INF claimTxn: 0xd6a67d0a29d9f12e7d0ab0763e13d9c9fb92e75f1738752a1d6d074a333fdfc4
7:52PM INF transaction successful txHash=0xd6a67d0a29d9f12e7d0ab0763e13d9c9fb92e75f1738752a1d6d074a333fdfc4
```

## Check balances to ensure deposit was claimed from A, B and C to L1

```bash
cast call "$contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l1_rpc_url"
cast call "$contract_addr" "balanceOf(address)(uint256)" "$eth_address" --rpc-url "$l1_rpc_url"

cast call "$l2_a_contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l2_a_rpc_url"
cast call "$l2_a_contract_addr" "balanceOf(address)(uint256)" "$eth_address" --rpc-url "$l2_a_rpc_url"

cast call "$l2_b_contract_addr" "balanceOf(address)(uint256)" "$pre_funded_address" --rpc-url "$l2_b_rpc_url"
cast call "$l2_b_contract_addr" "balanceOf(address)(uint256)" "$eth_address" --rpc-url "$l2_b_rpc_url"
```

```bash
999999900000000000000000 [9.999e23]
100000000000000000 [1e17]
0
0
0
0
```
