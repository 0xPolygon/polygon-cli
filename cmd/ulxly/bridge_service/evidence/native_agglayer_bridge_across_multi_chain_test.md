# Agglayer Bridge Service Manual Test

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

echo l2_c_rpc_url=\"$(kurtosis port print cdk op-el-1-op-geth-op-node-003 rpc)\"
echo l2_c_bridge_url=\"$(kurtosis port print cdk aggkit-003-bridge rest)\"
```

## Prepare env variables

```bash
pre_funded_private_key="0x12d7de8621a77640c9241b2595ba78ce443d05e94090365ab3bb5e19df82c625"

eth_address="0xced253B29527D62a1880b95C23F256CE78a73c06"
private_key="0x516d5e1c8f7e1da24379041b758b5d16fd066c8a8856791d3c5c0c79a81bad59"

bridge_address="0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6"

l1_rpc_url="http://127.0.0.1:32769"
l2_a_rpc_url="http://127.0.0.1:32779"
l2_a_bridge_url="http://127.0.0.1:32798"
l2_b_rpc_url="http://127.0.0.1:32812"
l2_b_bridge_url="http://127.0.0.1:32825"
l2_c_rpc_url="http://127.0.0.1:32839"
l2_c_bridge_url="http://127.0.0.1:32852"
```

## Sanity check for blocks

```bash
cast block-number --rpc-url "$l1_rpc_url"
cast block-number --rpc-url "$l2_a_rpc_url"
cast block-number --rpc-url "$l2_b_rpc_url"
cast block-number --rpc-url "$l2_c_rpc_url"
```

```bash
726
1117
739
389
```

## Get Bridge network ID

```bash
cast call --rpc-url "$l2_a_rpc_url" "$bridge_address" "networkID()(uint32)"
cast call --rpc-url "$l2_b_rpc_url" "$bridge_address" "networkID()(uint32)"
cast call --rpc-url "$l2_c_rpc_url" "$bridge_address" "networkID()(uint32)"
```

```bash
1
2
3
```

## Check current balances

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_c_rpc_url" --ether "$eth_address"
```

```bash
0.000000000000000000
0.000000000000000000
0.000000000000000000
0.000000000000000000
```

## Fund account on L1

```bash
cast send "$eth_address" --value 1ether --private-key "$pre_funded_private_key" --rpc-url "$l1_rpc_url"
```

```bash
blockHash            0xcd2faac64b5cd41d217902df76401104457c96e36d518bd251cd7a9cb8fa0056
blockNumber          786
contractAddress
cumulativeGasUsed    21000
effectiveGasPrice    1000000007
from                 0xE34aaF64b29273B7D567FCFc40544c014EEe9970
gasUsed              21000
logs                 []
logsBloom            0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
root
status               1 (success)
transactionHash      0xc4b94e54233e299cf4a62dc3b1a452afd77a668ea9a32c0a87275cb299308eba
transactionIndex     0
type                 2
blobGasPrice
blobGasUsed
to                   0xced253B29527D62a1880b95C23F256CE78a73c06
```

## Check balance again to identify balance on L1

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_c_rpc_url" --ether "$eth_address"
```

```bash
1.000000000000000000
0.000000000000000000
0.000000000000000000
0.000000000000000000
```

## Bridge from L1 to A

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 1 \
    --value 100000000000000000 \
    --destination-address "$eth_address" \
    --rpc-url "$l1_rpc_url"
```

```bash
6:26PM INF bridgeTxn: 0x65f9396a8f6811220c1f5aeaecda3c0600359f7ad946dd50adb5367adfa6990c
6:26PM INF transaction successful txHash=0x65f9396a8f6811220c1f5aeaecda3c0600359f7ad946dd50adb5367adfa6990c
6:26PM INF Bridge deposit count parsed from logs depositCount=3
```

## Check deposit is already indexed by bridge service on A

```bash
curl -s "$l2_a_bridge_url/bridge/v1/bridges?network_id=0&from_address=$eth_address" | jq -M '.bridges[] | select(.destination_network == 1)'
```

```bash
{
  "block_num": 835,
  "block_pos": 0,
  "from_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "tx_hash": "0x65f9396a8f6811220c1f5aeaecda3c0600359f7ad946dd50adb5367adfa6990c",
  "calldata": "0xcd5865790000000000000000000000000000000000000000000000000000000000000001000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c06000000000000000000000000000000000000000000000000016345785d8a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000",
  "block_timestamp": 1759094798,
  "leaf_type": 0,
  "origin_network": 0,
  "origin_address": "0x0000000000000000000000000000000000000000",
  "destination_network": 1,
  "destination_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "amount": "100000000000000000",
  "metadata": "0x",
  "deposit_count": 3,
  "is_native_token": true,
  "bridge_hash": "0x069a7b3e0a20d6cba4e9696f8eed44c611d945263ab4a0a3800d2cf5a2c9cd7e"
}
```

## Check balances to ensure deposit was auto claimed on A

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_c_rpc_url" --ether "$eth_address"
```

```bash
0.899806097998642686
0.100000000000000000
0.000000000000000000
0.000000000000000000
```

## Bridge from A to B

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 2 \
    --value 10000000000000000 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_a_rpc_url"
```

```bash
6:30PM INF bridgeTxn: 0x2fc2c9cc3c1e63a2067ad71915f8eb669a9ebc2c1b104415ce54e7938d1e9aa3
6:30PM INF transaction successful txHash=0x2fc2c9cc3c1e63a2067ad71915f8eb669a9ebc2c1b104415ce54e7938d1e9aa3
6:30PM INF Bridge deposit count parsed from logs depositCount=0
```

## Check balances to ensure deposit was made from A to B

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_c_rpc_url" --ether "$eth_address"
```

```bash
0.899806097998642686
0.089999566411625131
0.000000000000000000
0.000000000000000000
```

## Check deposit is already indexed by bridge service on A to be claimed on B

```bash
curl -s "$l2_a_bridge_url/bridge/v1/bridges?network_id=1&from_address=$eth_address" | jq -M '.bridges[] | select(.destination_network == 2)'
```

```bash
{
  "block_num": 1542,
  "block_pos": 0,
  "from_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "tx_hash": "0x2fc2c9cc3c1e63a2067ad71915f8eb669a9ebc2c1b104415ce54e7938d1e9aa3",
  "calldata": "0xcd5865790000000000000000000000000000000000000000000000000000000000000002000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c06000000000000000000000000000000000000000000000000002386f26fc100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000",
  "block_timestamp": 1759095006,
  "leaf_type": 0,
  "origin_network": 0,
  "origin_address": "0x0000000000000000000000000000000000000000",
  "destination_network": 2,
  "destination_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "amount": "10000000000000000",
  "metadata": "0x",
  "deposit_count": 0,
  "is_native_token": true,
  "bridge_hash": "0x0f3a128d185351c55765282e062f086b69d863254eff332bd3ec02c7e8c10d2f"
}
```

## Claim deposit made from A to B

```bash
polycli ulxly claim asset \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_a_bridge_url" \
    --private-key "$pre_funded_private_key" \
    --destination-address "$eth_address" \
    --deposit-network 1 \
    --deposit-count 0 \
    --rpc-url "$l2_b_rpc_url"
```

```bash
6:46PM INF The deposit is ready to be claimed
6:46PM INF claimTxn: 0x44a69405c1a6c4d4b02e645536d7c3454f2505a70fcb57760c0c5db98d249262
6:46PM INF transaction successful txHash=0x44a69405c1a6c4d4b02e645536d7c3454f2505a70fcb57760c0c5db98d249262
```

## Check balances to ensure deposit was claimed from A to B

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_c_rpc_url" --ether "$eth_address"
```

```bash
0.899806097998642686
0.089999566411625131
0.010000000000000000
0.000000000000000000
```

## Bridge from B to C

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 3 \
    --value 1000000000000000 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_b_rpc_url"
```

```bash
6:48PM INF bridgeTxn: 0x59969834a85a1f0b465016d8937bdfd7e768cd69ed8b57915fe0b31673f33ffd
6:48PM INF transaction successful txHash=0x59969834a85a1f0b465016d8937bdfd7e768cd69ed8b57915fe0b31673f33ffd
6:48PM INF Bridge deposit count parsed from logs depositCount=0
```

## Check balances to ensure deposit was made from B to C

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_c_rpc_url" --ether "$eth_address"
```

```bash
0.899806097998642686
0.089999566411625131
0.008999845825708303
0.000000000000000000
```

## Check deposit is already indexed by bridge service on B to be claimed on C

```bash
curl -s "$l2_b_bridge_url/bridge/v1/bridges?network_id=2&from_address=$eth_address" | jq -M '.bridges[] | select(.destination_network == 3)'
```

```bash
{
  "block_num": 2277,
  "block_pos": 0,
  "from_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "tx_hash": "0x59969834a85a1f0b465016d8937bdfd7e768cd69ed8b57915fe0b31673f33ffd",
  "calldata": "0xcd5865790000000000000000000000000000000000000000000000000000000000000003000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c0600000000000000000000000000000000000000000000000000038d7ea4c680000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000",
  "block_timestamp": 1759096119,
  "leaf_type": 0,
  "origin_network": 0,
  "origin_address": "0x0000000000000000000000000000000000000000",
  "destination_network": 3,
  "destination_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "amount": "1000000000000000",
  "metadata": "0x",
  "deposit_count": 0,
  "is_native_token": true,
  "bridge_hash": "0xadfa76521a290e60ab31b81909ab6c1ec7be4d7f39371285e3e464d783f046a6"
}
```

## Claim deposit made from B to C

```bash
polycli ulxly claim asset \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_b_bridge_url" \
    --private-key "$pre_funded_private_key" \
    --destination-address "$eth_address" \
    --deposit-network 2 \
    --deposit-count 0 \
    --rpc-url "$l2_c_rpc_url"
```

```bash
6:57PM INF The deposit is ready to be claimed
6:57PM INF claimTxn: 0x61027c8e10f3da9142fb3ed6076571f27b16c3d9ece83e710f428b5ed6222e07
6:57PM INF transaction successful txHash=0x61027c8e10f3da9142fb3ed6076571f27b16c3d9ece83e710f428b5ed6222e07
```

## Check balances to ensure deposit was claimed from B to C

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_c_rpc_url" --ether "$eth_address"
```

```bash
0.899806097998642686
0.089999566411625131
0.008999845825708303
0.001000000000000000
```

## Bridge from C to A

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 1 \
    --value 100000000000000 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_c_rpc_url"
```

```bash
6:59PM INF bridgeTxn: 0x5b0dc9efebd1c2c516eee0d357e76b709bfe7411b3fa1b3eb6bb8387ab019251
6:59PM INF transaction successful txHash=0x5b0dc9efebd1c2c516eee0d357e76b709bfe7411b3fa1b3eb6bb8387ab019251
6:59PM INF Bridge deposit count parsed from logs depositCount=0
```

## Check balances to ensure deposit was made from C to a

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_c_rpc_url" --ether "$eth_address"
```

```bash
0.899806097998642686
0.089999566411625131
0.008999845825708303
0.000899856360939603
```

## Check deposit is already indexed by bridge service on C to be claimed on A

```bash
curl -s "$l2_c_bridge_url/bridge/v1/bridges?network_id=3&from_address=$eth_address" | jq -M '.bridges[] | select(.destination_network == 1)'
```

```bash
{
  "block_num": 2554,
  "block_pos": 0,
  "from_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "tx_hash": "0x5b0dc9efebd1c2c516eee0d357e76b709bfe7411b3fa1b3eb6bb8387ab019251",
  "calldata": "0xcd5865790000000000000000000000000000000000000000000000000000000000000001000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c0600000000000000000000000000000000000000000000000000005af3107a40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000",
  "block_timestamp": 1759096746,
  "leaf_type": 0,
  "origin_network": 0,
  "origin_address": "0x0000000000000000000000000000000000000000",
  "destination_network": 1,
  "destination_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "amount": "100000000000000",
  "metadata": "0x",
  "deposit_count": 0,
  "is_native_token": true,
  "bridge_hash": "0xfea460151f4cb4379a0970869d13350dcd6b9769bc8540bb1a1deaf7338dd254"
}
```

## Claim deposit made from C to A

```bash
polycli ulxly claim asset \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_c_bridge_url" \
    --private-key "$pre_funded_private_key" \
    --destination-address "$eth_address" \
    --deposit-network 3 \
    --deposit-count 0 \
    --rpc-url "$l2_a_rpc_url"
```

```bash
7:03PM INF The deposit is ready to be claimed
7:03PM INF claimTxn: 0x9438bcef7fe845b36cf6f80f599d08cd5098a2351a1e3dc95059d1452e129c5f
7:03PM INF transaction successful txHash=0x9438bcef7fe845b36cf6f80f599d08cd5098a2351a1e3dc95059d1452e129c5f
```

## Check balances to ensure deposit was claimed from C to A

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_c_rpc_url" --ether "$eth_address"
```

```bash
0.899806097998642686
0.090099566411625131
0.008999845825708303
0.000899856360939603
```

## Bridge from A, B and C to L1

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 0 \
    --value 1000 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_a_rpc_url"
```

```bash
7:15PM INF bridgeTxn: 0xee45d59c25c3cef8e2f9e394ce8a4248a5ea857d6e8a95cec05c04c50d719720
7:15PM INF transaction successful txHash=0xee45d59c25c3cef8e2f9e394ce8a4248a5ea857d6e8a95cec05c04c50d719720
7:15PM INF Bridge deposit count parsed from logs depositCount=2
```

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 0 \
    --value 1000 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_b_rpc_url"
```

```bash
7:16PM INF bridgeTxn: 0x80c7394a589f76644658dc80639f59f59772b07278de946d0c223768b636cd1b
7:16PM INF transaction successful txHash=0x80c7394a589f76644658dc80639f59f59772b07278de946d0c223768b636cd1b
7:16PM INF Bridge deposit count parsed from logs depositCount=1
```

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 0 \
    --value 1000 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_c_rpc_url"
```

```bash
7:17PM INF bridgeTxn: 0x07e5431f948287e8b0660bf4b87d17e4385ddc250c02cca11ce82f4a48ea7deb
7:17PM INF transaction successful txHash=0x07e5431f948287e8b0660bf4b87d17e4385ddc250c02cca11ce82f4a48ea7deb
7:17PM INF Bridge deposit count parsed from logs depositCount=1
```

## Check balances to ensure deposit was made from A, B and C to L1

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_c_rpc_url" --ether "$eth_address"
```

```bash
0.899806097998642686
0.090099477172230568
0.008999739334984743
0.000899749830292043
```

## Check deposit is already indexed by bridge service on A, B and C to be claimed on L1

```bash
curl -s "$l2_a_bridge_url/bridge/v1/bridges?network_id=1&from_address=$eth_address" | jq -M '.bridges[] | select(.destination_network == 0)'

curl -s "$l2_b_bridge_url/bridge/v1/bridges?network_id=2&from_address=$eth_address" | jq -M '.bridges[] | select(.destination_network == 0)'

curl -s "$l2_c_bridge_url/bridge/v1/bridges?network_id=3&from_address=$eth_address" | jq -M '.bridges[] | select(.destination_network == 0)'
```

```bash
{
  "block_num": 4259,
  "block_pos": 0,
  "from_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "tx_hash": "0xee45d59c25c3cef8e2f9e394ce8a4248a5ea857d6e8a95cec05c04c50d719720",
  "calldata": "0xcd5865790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c0600000000000000000000000000000000000000000000000000000000000003e80000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000",
  "block_timestamp": 1759097723,
  "leaf_type": 0,
  "origin_network": 0,
  "origin_address": "0x0000000000000000000000000000000000000000",
  "destination_network": 0,
  "destination_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "amount": "1000",
  "metadata": "0x",
  "deposit_count": 2,
  "is_native_token": true,
  "bridge_hash": "0xc44a7a21a9c15aa1702e7e5e5871f6d96b01a92dc71c6ac6b47507b409cd6538"
}
{
  "block_num": 3963,
  "block_pos": 0,
  "from_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "tx_hash": "0x80c7394a589f76644658dc80639f59f59772b07278de946d0c223768b636cd1b",
  "calldata": "0xcd5865790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c0600000000000000000000000000000000000000000000000000000000000003e80000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000",
  "block_timestamp": 1759097805,
  "leaf_type": 0,
  "origin_network": 0,
  "origin_address": "0x0000000000000000000000000000000000000000",
  "destination_network": 0,
  "destination_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "amount": "1000",
  "metadata": "0x",
  "deposit_count": 1,
  "is_native_token": true,
  "bridge_hash": "0xc44a7a21a9c15aa1702e7e5e5871f6d96b01a92dc71c6ac6b47507b409cd6538"
}
{
  "block_num": 3642,
  "block_pos": 0,
  "from_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "tx_hash": "0x07e5431f948287e8b0660bf4b87d17e4385ddc250c02cca11ce82f4a48ea7deb",
  "calldata": "0xcd5865790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c0600000000000000000000000000000000000000000000000000000000000003e80000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000",
  "block_timestamp": 1759097834,
  "leaf_type": 0,
  "origin_network": 0,
  "origin_address": "0x0000000000000000000000000000000000000000",
  "destination_network": 0,
  "destination_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "amount": "1000",
  "metadata": "0x",
  "deposit_count": 1,
  "is_native_token": true,
  "bridge_hash": "0xc44a7a21a9c15aa1702e7e5e5871f6d96b01a92dc71c6ac6b47507b409cd6538"
}
```

## Claim deposit made from A, B and C to L1

```bash
polycli ulxly claim asset \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_a_bridge_url" \
    --private-key "$pre_funded_private_key" \
    --destination-address "$eth_address" \
    --deposit-network 1 \
    --deposit-count 2 \
    --rpc-url "$l1_rpc_url"
```

```bash
7:23PM INF The deposit is ready to be claimed
7:23PM INF claimTxn: 0x413151e17c071cd508ba79316ad299686c51849ed07e02f009666ff9a1c72541
7:23PM INF transaction successful txHash=0x413151e17c071cd508ba79316ad299686c51849ed07e02f009666ff9a1c72541
```

```bash
polycli ulxly claim asset \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_b_bridge_url" \
    --private-key "$pre_funded_private_key" \
    --destination-address "$eth_address" \
    --deposit-network 2 \
    --deposit-count 1 \
    --rpc-url "$l1_rpc_url"
```

```bash
7:24PM INF The deposit is ready to be claimed
7:24PM INF claimTxn: 0x9e1b132c45ee14135fa6275d00adc1aea9b22bfc1fc9a9346a806dd3f47dd3d8
7:24PM INF transaction successful txHash=0x9e1b132c45ee14135fa6275d00adc1aea9b22bfc1fc9a9346a806dd3f47dd3d8
```

```bash
polycli ulxly claim asset \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_c_bridge_url" \
    --private-key "$pre_funded_private_key" \
    --destination-address "$eth_address" \
    --deposit-network 3 \
    --deposit-count 1 \
    --rpc-url "$l1_rpc_url"
```

```bash
7:24PM INF The deposit is ready to be claimed
7:24PM INF claimTxn: 0xff0480c84b6b703cb7bf0b9e5b572b5b020009f8be97f773b501c404f33a7bb8
7:24PM INF transaction successful txHash=0xff0480c84b6b703cb7bf0b9e5b572b5b020009f8be97f773b501c404f33a7bb8
```

## Check balances to ensure deposit was claimed from A, B and C to L1

```bash
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_c_rpc_url" --ether "$eth_address"
```

```bash
0.899806097998645686
0.090099477172230568
0.008999739334984743
0.000899749830292043
```
