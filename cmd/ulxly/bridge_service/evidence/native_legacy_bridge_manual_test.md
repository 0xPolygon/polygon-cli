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
echo l2_a_bridge_url=\"$(kurtosis port print cdk zkevm-bridge-service-001 rpc)\"

echo l2_b_rpc_url=\"$(kurtosis port print cdk op-el-1-op-geth-op-node-002 rpc)\"
echo l2_b_bridge_url=\"$(kurtosis port print cdk zkevm-bridge-service-002 rpc)\"
```

## Prepare env variables

```bash
pre_funded_private_key="0x12d7de8621a77640c9241b2595ba78ce443d05e94090365ab3bb5e19df82c625"

eth_address="0xced253B29527D62a1880b95C23F256CE78a73c06"
private_key="0x516d5e1c8f7e1da24379041b758b5d16fd066c8a8856791d3c5c0c79a81bad59"

bridge_address="0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6"

l1_rpc_url="http://127.0.0.1:32867"
l2_a_rpc_url="http://127.0.0.1:32877"
l2_a_bridge_url="http://127.0.0.1:32906"
l2_b_rpc_url="http://127.0.0.1:32910"
l2_b_bridge_url="http://127.0.0.1:32933"
```

## Sanity check for blocks

```bash
cast block-number --rpc-url "$l1_rpc_url"
cast block-number --rpc-url "$l2_a_rpc_url"
cast block-number --rpc-url "$l2_b_rpc_url"
```

```bash
523
710
136
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

## Fund account on L1

```bash
cast send "$eth_address" --value 1ether --private-key "$pre_funded_private_key" --rpc-url "$l1_rpc_url"
```

```bash

```

## Check balance again to identify balance on L1

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
```

```bash
blockHash            0x32e29f612185e0e388d2143b26aa5957f8a72a0e89953ead5cdbacbbb1fdbd39
blockNumber          543
contractAddress
cumulativeGasUsed    21000
effectiveGasPrice    1000000007
from                 0xE34aaF64b29273B7D567FCFc40544c014EEe9970
gasUsed              21000
logs                 []
logsBloom            0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
root
status               1 (success)
transactionHash      0x12832eab0c81da82544da8bb59c4dd1469353aa7cc6ef17f851cf50bac983137
transactionIndex     0
type                 2
blobGasPrice
blobGasUsed
to                   0xced253B29527D62a1880b95C23F256CE78a73c06
```

-----

## Bridge from L1 to L2 - Network A

```bash
polycli ulxly bridge asset \
    --legacy \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 1 \
    --value 100000000000000000 \
    --destination-address "$eth_address" \
    --rpc-url "$l1_rpc_url"
```

```bash
4:40PM INF bridgeTxn: 0x87970b535abf1880b56722bd9b4780adf7f584fb8ad07f67053aa4130f0c3a7c
4:40PM INF transaction successful txHash=0x87970b535abf1880b56722bd9b4780adf7f584fb8ad07f67053aa4130f0c3a7c
4:40PM INF Bridge deposit count parsed from logs depositCount=2
```

## Check deposit is already indexed by bridge service

```bash
curl -s "$l2_a_bridge_url/bridges/$eth_address" | jq -M '.deposits[] | select(.dest_net == 1)'
```

```bash
{
  "leaf_type": 0,
  "orig_net": 0,
  "orig_addr": "0x0000000000000000000000000000000000000000",
  "amount": "100000000000000000",
  "dest_net": 1,
  "dest_addr": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "block_num": "551",
  "deposit_cnt": 2,
  "network_id": 0,
  "tx_hash": "0x87970b535abf1880b56722bd9b4780adf7f584fb8ad07f67053aa4130f0c3a7c",
  "claim_tx_hash": "",
  "metadata": "0x",
  "ready_for_claim": false,
  "global_index": "18446744073709551618"
}
```

## Check balances to ensure auto-claimer claimed the deposit

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
```

```bash
0.899999999998761552
0.100000000000000000
0.000000000000000000
```

## Bridge from L1 to L2 - Network B

```bash
polycli ulxly bridge asset \
    --legacy \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 2 \
    --value 100000000000000000 \
    --destination-address "$eth_address" \
    --rpc-url "$l1_rpc_url"
```

```bash
4:41PM INF bridgeTxn: 0x1fd70f42929f28e2836900195c144fcf5e8d0f8da283600acfc77e78dd246f2e
4:41PM INF transaction successful txHash=0x1fd70f42929f28e2836900195c144fcf5e8d0f8da283600acfc77e78dd246f2e
4:41PM INF Bridge deposit count parsed from logs depositCount=3
```

## Check deposit is already indexed by bridge service

```bash
curl -s "$l2_b_bridge_url/bridges/$eth_address" | jq -M '.deposits[] | select(.dest_net == 2)'
```

```bash
{
  "leaf_type": 0,
  "orig_net": 0,
  "orig_addr": "0x0000000000000000000000000000000000000000",
  "amount": "100000000000000000",
  "dest_net": 2,
  "dest_addr": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "block_num": "585",
  "deposit_cnt": 3,
  "network_id": 0,
  "tx_hash": "0x1fd70f42929f28e2836900195c144fcf5e8d0f8da283600acfc77e78dd246f2e",
  "claim_tx_hash": "",
  "metadata": "0x",
  "ready_for_claim": false,
  "global_index": "18446744073709551619"
}
```

## Check balances to ensure auto-claimer claimed the deposit

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
```

```bash
0.799999999997210336
0.100000000000000000
0.100000000000000000
```

-------

## Bridge from L2 to L1

```bash
polycli ulxly bridge asset \
    --legacy \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 0 \
    --value 1000001 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_a_rpc_url"
```

```bash
4:42PM INF bridgeTxn: 0x844311ee7be01727ed10af42162c03597d2de2e8830ae9fe1f6733dea4c3691d
4:42PM INF transaction successful txHash=0x844311ee7be01727ed10af42162c03597d2de2e8830ae9fe1f6733dea4c3691d
4:42PM INF Bridge deposit count parsed from logs depositCount=0
```

```bash
polycli ulxly bridge asset \
    --legacy \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 0 \
    --value 1000001 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_b_rpc_url"
```

```bash
4:42PM INF bridgeTxn: 0x8cb33394da318557bb3e758de9eff2b9b6596251d96ab23dec0df7d5e16a9483
4:42PM INF transaction successful txHash=0x8cb33394da318557bb3e758de9eff2b9b6596251d96ab23dec0df7d5e16a9483
4:42PM INF Bridge deposit count parsed from logs depositCount=0
```

## Claim deposit on L1

```bash
polycli ulxly claim asset \
    --legacy \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_a_bridge_url" \
    --private-key "$private_key" \
    --deposit-network 1 \
    --deposit-count 0 \
    --rpc-url "$l1_rpc_url"
```

```bash
4:47PM INF No destination address specified. Using private key's address destAddress=0xced253B29527D62a1880b95C23F256CE78a73c06
4:47PM INF The deposit is ready to be claimed
4:47PM INF claimTxn: 0xdd2cb27bc6e532ec8732243880d47b200f2c6b03ca4cd41ab53e1923045a25c2
4:47PM INF transaction successful txHash=0xdd2cb27bc6e532ec8732243880d47b200f2c6b03ca4cd41ab53e1923045a25c2
```

```bash
polycli ulxly claim asset \
    --legacy \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_b_bridge_url" \
    --private-key "$private_key" \
    --deposit-network 2 \
    --deposit-count 0 \
    --rpc-url "$l1_rpc_url"
```

```bash
4:50PM INF No destination address specified. Using private key's address destAddress=0xced253B29527D62a1880b95C23F256CE78a73c06
4:50PM INF The deposit is ready to be claimed
4:50PM INF claimTxn: 0xd7462481689232b5d4d501bac39795a7ea883c93f290e0cb60bf670e5d217b80
4:50PM INF transaction successful txHash=0xd7462481689232b5d4d501bac39795a7ea883c93f290e0cb60bf670e5d217b80
```

## Check balances after claims on L1

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
```

```bash
0.799999999997108098
0.099996219810141652
0.099966001821780277
```

----- 

## Bridge from L2-A to L2-B

```bash
polycli ulxly bridge asset \
    --legacy \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 2 \
    --value 1000002 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_a_rpc_url"
```

```bash
4:51PM INF bridgeTxn: 0x9f6b135d044b8466b08723c69ba65bbe5876e14d4a99185d85979df69f676ca3
4:51PM INF transaction successful txHash=0x9f6b135d044b8466b08723c69ba65bbe5876e14d4a99185d85979df69f676ca3
4:51PM INF Bridge deposit count parsed from logs depositCount=1
```

```bash
polycli ulxly claim asset \
    --legacy \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_a_bridge_url" \
    --private-key "$private_key" \
    --deposit-network 1 \
    --deposit-count 1 \
    --rpc-url "$l2_b_rpc_url"
```

```bash
4:54PM INF No destination address specified. Using private key's address destAddress=0xced253B29527D62a1880b95C23F256CE78a73c06
4:54PM INF The deposit is ready to be claimed
4:54PM INF claimTxn: 0xa09d87072b514defa383c31aefed6b5e5ffcc627af7d70579c850e92e8c37729
4:54PM INF transaction successful txHash=0xa09d87072b514defa383c31aefed6b5e5ffcc627af7d70579c850e92e8c37729
```

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
```

```bash
0.799999999997108098
0.099995755435525731
0.099963645641196397
```

## Bridge from L2-B to L2-A

```bash
polycli ulxly bridge asset \
    --legacy \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 1 \
    --value 1000002 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_b_rpc_url"
```

```bash
4:55PM INF bridgeTxn: 0x66344a689a16f858fbc3d43ce54747b03539ce5c7b4ec2f9b84223e1e2fc2e74
4:55PM INF transaction successful txHash=0x66344a689a16f858fbc3d43ce54747b03539ce5c7b4ec2f9b84223e1e2fc2e74
4:55PM INF Bridge deposit count parsed from logs depositCount=1
```

```bash
polycli ulxly claim asset \
    --legacy \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_b_bridge_url" \
    --private-key "$private_key" \
    --deposit-network 2 \
    --deposit-count 1 \
    --rpc-url "$l2_a_rpc_url"
```

```bash
5:03PM INF No destination address specified. Using private key's address destAddress=0xced253B29527D62a1880b95C23F256CE78a73c06
5:03PM INF The deposit is ready to be claimed
5:03PM INF claimTxn: 0x0cb920c4af8be6f6dd362a96254781876847db5db67e01c00bfb1c608c5fb9db
5:03PM INF transaction successful txHash=0x0cb920c4af8be6f6dd362a96254781876847db5db67e01c00bfb1c608c5fb9db
```

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
```

```bash
0.799999999997108098
0.099995562970002721
0.099962345360592177
```
