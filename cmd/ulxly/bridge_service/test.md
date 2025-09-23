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

l2_a_rpc_url=http://127.0.0.1:32817
l2_a_bridge_url=http://127.0.0.1:32837

# l2_b_rpc_url=https://rpc-forknet-testnet.t.conduit.xyz
# l2_b_bridge_url=https://rpc-bridge-katana-bokuto.t.conduit.xyz

l1_rpc_url=http://127.0.0.1:32807
bridge_address="0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6"
```

## Sanity check for blocks

```bash
cast block-number --rpc-url "$l2_a_rpc_url"
# cast block-number --rpc-url "$l2_b_rpc_url"
cast block-number --rpc-url "$l1_rpc_url"
```

```bash
544
448
```

## Get Bridge network ID

```bash
cast call --rpc-url "$l2_a_rpc_url" "$bridge_address" "networkID()(uint32)"
# cast call --rpc-url "$l2_b_rpc_url" "$bridge_address" "networkID()(uint32)"
```

```bash
1
```

## Check current balances

```bash
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
# cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
```

```bash
0.000000000000000000
0.000000000000000000
```

## Fund account on L1

```bash
cast send "$eth_address" --value 1ether --private-key "$pre_funded_private_key" --rpc-url "$l1_rpc_url"
```

```bash
blockHash            0x06eca6e22e1a7f0f27bc8cfa9abab909d399431a1d0c1df31e54606f2bad5621
blockNumber          854
contractAddress
cumulativeGasUsed    21000
effectiveGasPrice    8
from                 0xE34aaF64b29273B7D567FCFc40544c014EEe9970
gasUsed              21000
logs                 []
logsBloom            0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
root
status               1 (success)
transactionHash      0xa490addb1c4d52c5b7f7eefcbd6b31ecc03a14294aac7350911adadb465bc50b
transactionIndex     0
type                 2
blobGasPrice
blobGasUsed
to                   0xced253B29527D62a1880b95C23F256CE78a73c06
```

## Check balance again to identify balance on L1 

```bash
cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
# cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
```

```bash
0.000000000000000000
1.000000000000000000
```

## Bridge from L1 to L2

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
1:29PM INF bridgeTxn: 0xd5a98966831391c61525b44bbc822b2ba93edcb1f170731a4680d5558e809959
1:29PM INF transaction successful txHash=0xd5a98966831391c61525b44bbc822b2ba93edcb1f170731a4680d5558e809959
1:29PM INF Bridge deposit count parsed from logs depositCount=222
```

## Get L1 to L2 bridge TX details

```bash
cast tx 0xd5a98966831391c61525b44bbc822b2ba93edcb1f170731a4680d5558e809959 --rpc-url "$l1_rpc_url"
```

```bash
blockHash            0x13394a79803c671188624eb06239b9e890f0b3bfd4be25179f78e633e2f5e6d4
blockNumber          979
from                 0xced253B29527D62a1880b95C23F256CE78a73c06
transactionIndex     0
effectiveGasPrice    8

accessList           []
chainId              271828
gasLimit             179949
hash                 0xd5a98966831391c61525b44bbc822b2ba93edcb1f170731a4680d5558e809959
input                0xcd5865790000000000000000000000000000000000000000000000000000000000000001000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c06000000000000000000000000000000000000000000000000016345785d8a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000
maxFeePerGas         15
maxPriorityFeePerGas 1
nonce                0
r                    0xe4d7eeaadc4b5c021904dce37a60148fc5502a282c8f5265add47bb2aa62819f
s                    0x2115d5d21dd30d5a8c8d3fb8662ea2510771082e7d899a18053c8a62fd86aae6
to                   0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6
type                 2
value                100000000000000000
yParity              1
```

## Check deposit is already indexed by bridge service

```bash
curl -s "$l2_a_bridge_url/bridge/v1/bridges?network_id=0&from_address=$eth_address" | jq -M '.bridges[] | select(.destination_network == 1'
```

```bash
{
  "block_num": 979,
  "block_pos": 0,
  "from_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "tx_hash": "0xd5a98966831391c61525b44bbc822b2ba93edcb1f170731a4680d5558e809959",
  "calldata": "0xcd5865790000000000000000000000000000000000000000000000000000000000000001000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c06000000000000000000000000000000000000000000000000016345785d8a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000",
  "block_timestamp": 1758644958,
  "leaf_type": 0,
  "origin_network": 0,
  "origin_address": "0x0000000000000000000000000000000000000000",
  "destination_network": 1,
  "destination_address": "0xced253B29527D62a1880b95C23F256CE78a73c06",
  "amount": "100000000000000000",
  "metadata": "0x",
  "deposit_count": 222,
  "is_native_token": true,
  "bridge_hash": "0x069a7b3e0a20d6cba4e9696f8eed44c611d945263ab4a0a3800d2cf5a2c9cd7e"
}
```

## Check balances to ensure auto-claimer claimed the deposit

```bash
 cast balance --rpc-url "$l2_a_rpc_url" --ether "$eth_address"
# cast balance --rpc-url "$l2_b_rpc_url" --ether "$eth_address"
cast balance --rpc-url "$l1_rpc_url" --ether "$eth_address"
```

```bash
0.100000000000000000
0.899999999998623696
```

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
1:58PM INF bridgeTxn: 0x6b93091338bedb097c0cd9efa1b5224f74847009f05dfa1b6b7d35ebd790dce6
1:58PM INF transaction successful txHash=0x6b93091338bedb097c0cd9efa1b5224f74847009f05dfa1b6b7d35ebd790dce6
1:58PM INF Bridge deposit count parsed from logs depositCount=482
```

## Get L2 to L1 bridge TX details

```bash
cast tx 0x6b93091338bedb097c0cd9efa1b5224f74847009f05dfa1b6b7d35ebd790dce6 --rpc-url "$l2_a_rpc_url"
```

```bash
blockHash            0x8c4a50a2e11846f991cc86687d5589d8f273c4c4f7370a1821ce423f4225e7fc
blockNumber          3344
from                 0xced253B29527D62a1880b95C23F256CE78a73c06
transactionIndex     1
effectiveGasPrice    1001815

accessList           []
chainId              2151908
gasLimit             104579
hash                 0x6b93091338bedb097c0cd9efa1b5224f74847009f05dfa1b6b7d35ebd790dce6
input                0xcd5865790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ced253b29527d62a1880b95c23f256ce78a73c0600000000000000000000000000000000000000000000000000000000000f42410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000
maxFeePerGas         1003658
maxPriorityFeePerGas 1000000
nonce                0
r                    0x5e93f67ffaa34c1ce3fd9d27a63df34483098077222b9aa4eaea23f64ac35cc1
s                    0x4c928e86d91b6384ac7dc8d6afe700c6571fc8507e1d5cee7441b56226b73161
to                   0x78908F7A87d589fdB46bdd5EfE7892C5aD6001b6
type                 2
value                1000001
yParity              0
```

## Get the root from the block where the bridge from L2 to L1 was mined

```bash
cast call --block 3344 --rpc-url "$l2_a_rpc_url" "$bridge_address" "getRoot()(bytes32)"
```

```bash
0x9f9c7aae2107f703dda1a7cf6c2847d01eb0a54bb92d4c37c6a20155f090d968
```

## Claim deposit on L1

```bash
polycli ulxly claim asset \
    --bridge-address "$bridge_address" \
    --bridge-service-url "$l2_a_bridge_url" \
    --private-key "$private_key" \
    --deposit-network 1 \
    --deposit-count 482 \
    --rpc-url "$l1_rpc_url"
```
