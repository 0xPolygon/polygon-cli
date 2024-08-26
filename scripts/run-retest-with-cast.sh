#!/bin/bash

private_key="0xbcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31"
rpc_url="http://127.0.0.1:32887"
legacy_flag=" --legacy "
clean_up=0 # 0 will cause a cleanup 1 would leave the files
gas_limit=30000000



function normalize_address() {
        sed 's/0x//' |
         tr '[:upper:]' '[:lower:]'
}

function process_test_item() {
    local testfile=$1
    2>&1 echo "processing file at $testfile"
    local tmp_dir=$(mktemp -p /tmp -d retest-work-XXXXXXXXXXXX)
    pushd $tmp_dir

    local nonce=$(cast nonce --rpc-url $rpc_url $wallet_address)
    local count=0
    jq -c '.dependencies[]' $testfile | while read pre ; do
        count=$((count+1))
        echo $pre | jq '.' > dep-$count.json
        local code_to_deploy=$(jq -r '.code' dep-$count.json)
        local reference_address=$(jq -r '.addr' dep-$count.json | normalize_address)
        2>&1 echo "deploying dependency $count for $reference_address"
        echo $nonce
        cast send $legacy_flag --async --nonce $nonce --rpc-url $rpc_url --private-key $private_key --create "$code_to_deploy" | tee $reference_address.txhash
        cast compute-address --nonce $nonce $wallet_address | sed 's/^.*0x/0x/' > $reference_address.actual

        nonce=$((nonce+1))
        echo "$nonce" > last.nonce
    done
    # Random transaction to make sure all of the async deps are deployed before running the transactions

    cast send $legacy_flag --nonce $(cat last.nonce) --rpc-url $rpc_url --private-key $private_key --value 1 $wallet_address
    2>&1 echo "We have finished deploying the dependencies (I think)"

    count=0
    jq -c '.testCases[]' $testfile | while read test_case ; do
        count=$((count+1))
        echo $test_case | jq '.' > test_case_$count.json
        tx_input=$(jq -r '.input' test_case_$count.json)
        local name=$(jq -r '.name' test_case_$count.json)
        local addr=$(jq -r '.to' test_case_$count.json | normalize_address)
        local gas=$(jq -r '.gas' test_case_$count.json) # this value can be obscenely high in the test cases
        local val=$(jq -r '.value' test_case_$count.json)
        val_arg=""
        if [[ $val != "0x0" ]] ; then
            hex_in=$(echo $val | sed 's/0x//' | tr '[:lower:]' '[:upper:]')
            dec_val=$(bc <<< "ibase=16; $hex_in")
            val_arg=" --value $dec_val "
        fi

        local to_addr_arg=""
        if [[ $addr == "0x0000000000000000000000000000000000000000" || $addr == "" || $addr == "0000000000000000000000000000000000000000" ]] ; then
            to_addr_arg=" --create "
        else
            resolved_address=$(cat $addr.actual)
            to_addr_arg=" $resolved_address "
        fi

        set -x
        cast send $legacy_flag --rpc-url $rpc_url --private-key $private_key --gas-limit $gas_limit $val_arg $to_addr_arg $tx_input | tee tx-$count-out.json
        set +x
    done

    popd
    if $clean_up ; then
        rm -rf $tmp_dir
        rm $testfile
    fi
}

wallet_address=$(cast wallet address --private-key 0xbcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31)
wallet_balance=$(cast balance --rpc-url $rpc_url $wallet_address)
wallet_nonce=$(cast nonce --rpc-url $rpc_url $wallet_address)
2>&1 echo "Address $wallet_address has a balance of $wallet_balance and nonce $wallet_nonce"

# Create a temp file to store the output of the test cases
tmpfile=$(mktemp retest-jq-XXXXXXXXXXXX)

# Break down each test into different files
jq -c '.[]' | while read test_item ; do
    testfile=$(mktemp -p /tmp retest-item-jq-XXXXXXXXXXXX)
    echo $test_item > $testfile
    process_test_item $testfile
done

rm $tmpfile

