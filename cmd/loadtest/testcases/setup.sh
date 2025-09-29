#!/usr/bin/env bash

anvil --block-base-fee-per-gas 100 --balance 10000000 &> anvil.out &
echo $! > anvil.pid

mitmdump -p 8484 --mode reverse:http://127.0.0.1:8545 &> mitm.out &
echo $! > mitm.pid



rpc_url="http://127.0.0.1:8545"
private_key="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
eth_address="$(cast wallet address --private-key "$private_key")"


deterministic_deployer_addr=0x4e59b44847b379578588920ca78fbf26c0b4956c

loadtester_code=$(cat ../../../contracts/out/LoadTester.sol/LoadTester.json  | jq -r '.bytecode.object')
cast send --private-key "$private_key" \
     --rpc-url "$rpc_url" \
     "$deterministic_deployer_addr" \
     $(cast concat-hex $(cast hz) $(echo "$loadtester_code"))

loadtester_addr=$(cast create2 --init-code "$loadtester_code" --salt $(cast hz))
cast code --rpc-url "$rpc_url" "$loadtester_addr"


erc_20_code=$(cat ../../../contracts/out/ERC20.sol/ERC20.json  | jq -r '.bytecode.object')
cast send --private-key "$private_key" \
     --rpc-url "$rpc_url" \
     "$deterministic_deployer_addr" \
     $(cast concat-hex $(cast hz) $(echo "$erc_20_code"))

erc_20_addr=$(cast create2 --init-code "$erc_20_code" --salt $(cast hz))

erc_721_code=$(cat ../../../contracts/out/ERC721.sol/ERC721.json  | jq -r '.bytecode.object')
cast send --private-key "$private_key" \
     --rpc-url "$rpc_url" \
     "$deterministic_deployer_addr" \
     $(cast concat-hex $(cast hz) $(echo "$erc_721_code"))

erc_721_addr=$(cast create2 --init-code "$erc_721_code" --salt $(cast hz))


cast send --private-key "$private_key" \
     --rpc-url "$rpc_url" \
     --value 1000000ether \
     0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6

keys_to_warm=(
    0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
    0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a
    0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa
)

for pk in ${keys_to_warm[@]}; do
    addr=$(cast wallet address --private-key "$pk")
    echo $addr
    for i in {1..5}; do
        cast send --private-key "$pk" --rpc-url "$rpc_url" --value 1 $addr
        cast send --private-key "$pk" \
             --rpc-url "$rpc_url" \
             "$erc_721_addr" \
             'mintBatch(address,uint256)' "$addr" 2
        cast send --private-key "$pk" \
             --rpc-url "$rpc_url" \
             "$erc_20_addr" \
             'mint(uint256)' 1000000000000000000
    done
done

shuf test-commands.sh | while read polyclicmd ; do
    echo $polyclicmd > cur.out
    if ! timeout 120s bash -c "$polyclicmd" &>> cur.out ; then
        rc=$?
        if [[ $rc -eq 124 ]]; then
            mv cur.out timeout.$(date +%s).out
        else
            mv cur.out failure.$(date +%s).out
        fi
    fi
done


kill "$(cat anvil.pid)"
kill "$(cat mitm.pid)"
