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

## Prepare env variables

```bash
pre_funded_private_key="0x12d7de8621a77640c9241b2595ba78ce443d05e94090365ab3bb5e19df82c625"

eth_address="0xced253B29527D62a1880b95C23F256CE78a73c06"
private_key="0x516d5e1c8f7e1da24379041b758b5d16fd066c8a8856791d3c5c0c79a81bad59"

l2_a_rpc_url=http://127.0.0.1:32904
l2_a_bridge_url=http://127.0.0.1:32923

l2_b_rpc_url=http://127.0.0.1:32937
l2_b_bridge_url=http://127.0.0.1:32950

l1_rpc_url=http://127.0.0.1:32894
bridge_address="0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6"
```

## Sanity check for blocks

```bash
cast block-number --rpc-url "$l2_a_rpc_url"
cast block-number --rpc-url "$l2_b_rpc_url"
cast block-number --rpc-url "$l1_rpc_url"
```

```bash
1367
971
848
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
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
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
blockHash            0xda81134302191a61c05525de94e4e7463e163ed0fb1a68908cc6555232af8daa
blockNumber          868
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

## Check balance again to identify balance on L1

```bash
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
```

```bash
0.000000000000000000
0.000000000000000000
1.000000000000000000
```

-----

## Bridge from L1 to L2 - Network A

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
4:00PM INF bridgeTxn: 0x87970b535abf1880b56722bd9b4780adf7f584fb8ad07f67053aa4130f0c3a7c
4:00PM INF transaction successful txHash=0x87970b535abf1880b56722bd9b4780adf7f584fb8ad07f67053aa4130f0c3a7c
4:00PM INF Bridge deposit count parsed from logs depositCount=2
```

## Get L1 to L2 bridge TX details

```bash
cast tx 0x87970b535abf1880b56722bd9b4780adf7f584fb8ad07f67053aa4130f0c3a7c --rpc-url "$l1_rpc_url"
```

```bash
blockHash            0x43394162da5b9d22599f0f43d33df03a3dc70986871e386492d357819ac57be0
blockNumber          884
from                 0xced253B29527D62a1880b95C23F256CE78a73c06
transactionIndex     0
effectiveGasPrice    8

accessList           []
chainId              271828
gasLimit             162444
hash                 0x87970b535abf1880b56722bd9b4780adf7f584fb8ad07f67053aa4130f0c3a7c
input                0xcd5865790000000000000000000000000000000000000000000000000000000000000001000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c06000000000000000000000000000000000000000000000000016345785d8a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000
maxFeePerGas         15
maxPriorityFeePerGas 1
nonce                0
r                    0x35ef839cd88965bf29f2960bd38830b27dc44556f843f710cb2e57aa5d666a40
s                    0x194f0f8cb66f20946304d824a99c3c53fd32de1dd03f5df231236fd01f39920e
to                   0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6
type                 2
value                100000000000000000
yParity              0
```

## Check deposit is already indexed by bridge service

```bash
curl -s "$l2_a_bridge_url/bridge/v1/bridges?network_id=0&from_address=$eth_address" | jq -M '.bridges[] | select(.destination_network == 1)'
```

```bash
{
  "block_num": 884,
  "block_pos": 0,
  "from_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "tx_hash": "0x87970b535abf1880b56722bd9b4780adf7f584fb8ad07f67053aa4130f0c3a7c",
  "calldata": "0xcd5865790000000000000000000000000000000000000000000000000000000000000001000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c06000000000000000000000000000000000000000000000000016345785d8a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000",
  "block_timestamp": 1758740453,
  "leaf_type": 0,
  "origin_network": 0,
  "origin_address": "0x0000000000000000000000000000000000000000",
  "destination_network": 1,
  "destination_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "amount": "100000000000000000",
  "metadata": "0x",
  "deposit_count": 2,
  "is_native_token": true,
  "bridge_hash": "0x069a7b3e0a20d6cba4e9696f8eed44c611d945263ab4a0a3800d2cf5a2c9cd7e"
}
```

## Check balances to ensure auto-claimer claimed the deposit

```bash
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
```

```bash
0.100000000000000000
0.000000000000000000
0.899999999998761552
```

## Bridge from L1 to L2 - Network B

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 2 \
    --value 100000000000000000 \
    --destination-address "$eth_address" \
    --rpc-url "$l1_rpc_url"
```

```bash
4:04PM INF bridgeTxn: 0x1fd70f42929f28e2836900195c144fcf5e8d0f8da283600acfc77e78dd246f2e
4:04PM INF transaction successful txHash=0x1fd70f42929f28e2836900195c144fcf5e8d0f8da283600acfc77e78dd246f2e
4:04PM INF Bridge deposit count parsed from logs depositCount=3
```

## Get L1 to L2 bridge TX details

```bash
cast tx 0x1fd70f42929f28e2836900195c144fcf5e8d0f8da283600acfc77e78dd246f2e --rpc-url "$l1_rpc_url"
```

```bash
blockHash            0x0665b50b58022b9d354ccb6ede874036946e51ebbff3150e9a2bab2141589e32
blockNumber          998
from                 0xced253B29527D62a1880b95C23F256CE78a73c06
transactionIndex     0
effectiveGasPrice    8

accessList           []
chainId              271828
gasLimit             203739
hash                 0x1fd70f42929f28e2836900195c144fcf5e8d0f8da283600acfc77e78dd246f2e
input                0xcd5865790000000000000000000000000000000000000000000000000000000000000002000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c06000000000000000000000000000000000000000000000000016345785d8a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000
maxFeePerGas         15
maxPriorityFeePerGas 1
nonce                1
r                    0x97cd13ea90658477079544e4a870f6109ccd247946b2eb47db9a63578f569ff2
s                    0x6ccd9ccae36b6ffba36d1656fb3eee2277ff4694a52a85a1390e329603a92b21
to                   0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6
type                 2
value                100000000000000000
yParity              1
```

## Check deposit is already indexed by bridge service

```bash
curl -s "$l2_b_bridge_url/bridge/v1/bridges?network_id=0&from_address=$eth_address" | jq -M '.bridges[] | select(.destination_network == 2)'
```

```bash
{
  "block_num": 998,
  "block_pos": 0,
  "from_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "tx_hash": "0x1fd70f42929f28e2836900195c144fcf5e8d0f8da283600acfc77e78dd246f2e",
  "calldata": "0xcd5865790000000000000000000000000000000000000000000000000000000000000002000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c06000000000000000000000000000000000000000000000000016345785d8a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000",
  "block_timestamp": 1758740681,
  "leaf_type": 0,
  "origin_network": 0,
  "origin_address": "0x0000000000000000000000000000000000000000",
  "destination_network": 2,
  "destination_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "amount": "100000000000000000",
  "metadata": "0x",
  "deposit_count": 3,
  "is_native_token": true,
  "bridge_hash": "0x356c12c3e53e767acce842aff5513bb67f564b9192d5daf2bd8bd850bfcd0f1f"
}
```

## Check balances to ensure auto-claimer claimed the deposit

```bash
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
```

```bash
0.100000000000000000
0.100000000000000000
0.799999999997210336
```

-------

## Bridge from L2 to L1

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 0 \
    --value 1000001 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_a_rpc_url"
```

```bash
4:06PM INF bridgeTxn: 0x864a562c2ab63c7ebdd7cfff3174b2b521bef1cda69b9cb655ca7e5a6e9b5162
4:06PM INF transaction successful txHash=0x864a562c2ab63c7ebdd7cfff3174b2b521bef1cda69b9cb655ca7e5a6e9b5162
4:06PM INF Bridge deposit count parsed from logs depositCount=0
```

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 0 \
    --value 1000001 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_b_rpc_url"
```

```bash
4:06PM INF bridgeTxn: 0xa8146e7a218db2eddf86f35e84dbc438c5e5fa8810c5f1ea79258c11eff64e03
4:06PM INF transaction successful txHash=0xa8146e7a218db2eddf86f35e84dbc438c5e5fa8810c5f1ea79258c11eff64e03
4:06PM INF Bridge deposit count parsed from logs depositCount=0
```


## Get L2 to L1 bridge TXs details

```bash
cast tx 0x864a562c2ab63c7ebdd7cfff3174b2b521bef1cda69b9cb655ca7e5a6e9b5162 --rpc-url "$l2_a_rpc_url"
```

```bash
blockHash            0xee4a0d51712e111b7c26886887366d145aeaccf7cd9bb305c7612a63fda0b1fe
blockNumber          1770
from                 0xced253B29527D62a1880b95C23F256CE78a73c06
transactionIndex     1
effectiveGasPrice    1858969

accessList           []
chainId              20201
gasLimit             145752
hash                 0x864a562c2ab63c7ebdd7cfff3174b2b521bef1cda69b9cb655ca7e5a6e9b5162
input                0xcd5865790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c0600000000000000000000000000000000000000000000000000000000000f42410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000
maxFeePerGas         2731696
maxPriorityFeePerGas 1000000
nonce                0
r                    0x631ae4ae09a6e232fbb14f1d0637e7f15d7c247804493196b61f08eea9e07e93
s                    0x6415dd13d488001fbd3902e27ae245515ede1446bd2c87f6059496f98a922a26
to                   0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6
type                 2
value                1000001
yParity              1
```

```bash
cast tx 0xa8146e7a218db2eddf86f35e84dbc438c5e5fa8810c5f1ea79258c11eff64e03 --rpc-url "$l2_b_rpc_url"
```

```bash
blockHash            0xcde72efdc6b1cd7d2511ff004a1ded419c9c4a81d1ef2517c9637774273a7bbf
blockNumber          1389
from                 0xced253B29527D62a1880b95C23F256CE78a73c06
transactionIndex     1
effectiveGasPrice    4925873

accessList           []
chainId              20202
gasLimit             145752
hash                 0xa8146e7a218db2eddf86f35e84dbc438c5e5fa8810c5f1ea79258c11eff64e03
input                0xcd5865790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c0600000000000000000000000000000000000000000000000000000000000f42410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000
maxFeePerGas         8914644
maxPriorityFeePerGas 1000000
nonce                0
r                    0xab3b0f57eca586a00f7c677db8efd2230825dd3668293550a0b1e6e1891ce213
s                    0x44fb8b7c0830fd733575a0e0fdc583341eae0245550fa20fe14c091c5dc99836
to                   0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6
type                 2
value                1000001
yParity              0
```

## Get the root from the block where the bridge from L2 to L1 was mined

```bash
cast call --block 1770 --rpc-url "$l2_a_rpc_url" "$bridge_address" "getRoot()(bytes32)"
cast call --block 1389 --rpc-url "$l2_b_rpc_url" "$bridge_address" "getRoot()(bytes32)"
```

```bash
0x698d3b07df614b290697dc156b1b0257afc0fe784b64b7be6702eddbec6b4984
0x698d3b07df614b290697dc156b1b0257afc0fe784b64b7be6702eddbec6b4984
```

## Claim deposit on L1

```bash
polycli ulxly claim asset \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_a_bridge_url" \
    --private-key "$private_key" \
    --deposit-network 1 \
    --deposit-count 0 \
    --rpc-url "$l1_rpc_url"
```

```bash
4:09PM INF No destination address specified. Using private key's address destAddress=0xced253B29527D62a1880b95C23F256CE78a73c06
4:09PM INF The deposit is ready to be claimed
4:09PM INF claimTxn: 0xdd2cb27bc6e532ec8732243880d47b200f2c6b03ca4cd41ab53e1923045a25c2
4:09PM INF transaction successful txHash=0xdd2cb27bc6e532ec8732243880d47b200f2c6b03ca4cd41ab53e1923045a25c2
```

```bash
polycli ulxly claim asset \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_b_bridge_url" \
    --private-key "$private_key" \
    --deposit-network 2 \
    --deposit-count 0 \
    --rpc-url "$l1_rpc_url"
```

```bash
4:09PM INF No destination address specified. Using private key's address destAddress=0xced253B29527D62a1880b95C23F256CE78a73c06
4:09PM INF The deposit is ready to be claimed
4:09PM INF claimTxn: 0xd7462481689232b5d4d501bac39795a7ea883c93f290e0cb60bf670e5d217b80
4:09PM INF transaction successful txHash=0xd7462481689232b5d4d501bac39795a7ea883c93f290e0cb60bf670e5d217b80
```

## Check balances after claims on L1

```bash
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
```

```bash
0.099999742764164528
0.099999318381323528
0.799999999997108098
```

----- 

## Bridge from L2-A to L2-B

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 2 \
    --value 1000002 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_a_rpc_url"
```

```bash
4:12PM INF bridgeTxn: 0x1e7d0fd70ae6597d4be9a220d106cba2907b741da7ca0c6fef733bcdb555d425
4:12PM INF transaction successful txHash=0x1e7d0fd70ae6597d4be9a220d106cba2907b741da7ca0c6fef733bcdb555d425
4:12PM INF Bridge deposit count parsed from logs depositCount=1
```

```bash
polycli ulxly claim asset \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_a_bridge_url" \
    --private-key "$private_key" \
    --deposit-network 1 \
    --deposit-count 1 \
    --rpc-url "$l2_b_rpc_url"
```

```bash
4:48PM INF No destination address specified. Using private key's address destAddress=0xced253B29527D62a1880b95C23F256CE78a73c06
4:48PM INF The deposit is ready to be claimed
4:48PM INF claimTxn: 0xea18ed66cd417e50654958d19b7c01975c8342c30d979eddeded87858ef6d5c6
4:48PM INF transaction successful txHash=0xea18ed66cd417e50654958d19b7c01975c8342c30d979eddeded87858ef6d5c6
```

```bash
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
```

```bash
0.099999612884329807
0.099999153034088944
0.799999999997108098
```

## Bridge from L2-B to L2-A

```bash
polycli ulxly bridge asset \
    --bridge-address "$bridge_address" \
    --private-key "$private_key" \
    --destination-network 1 \
    --value 1000002 \
    --destination-address "$eth_address" \
    --rpc-url "$l2_b_rpc_url"
```

```bash
4:52PM INF bridgeTxn: 0x206747461a57e72466f30df4039ead81f6a42733b74c1b500f30a759e3232170
4:52PM INF transaction successful txHash=0x206747461a57e72466f30df4039ead81f6a42733b74c1b500f30a759e3232170
4:52PM INF Bridge deposit count parsed from logs depositCount=1
```

```bash
polycli ulxly claim asset \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_b_bridge_url" \
    --private-key "$private_key" \
    --deposit-network 2 \
    --deposit-count 1 \
    --rpc-url "$l2_a_rpc_url"
```

```bash
4:55PM INF No destination address specified. Using private key's address destAddress=0xced253B29527D62a1880b95C23F256CE78a73c06
4:55PM INF The deposit is ready to be claimed
4:55PM INF claimTxn: 0x3aad216a98754efafc357e57af2e409050a09f625438d067fadf7f3ef1b11796
4:55PM INF transaction successful txHash=0x3aad216a98754efafc357e57af2e409050a09f625438d067fadf7f3ef1b11796
```
