#!/bin/bash

private_key="0xbcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31"
rpc_url="http://127.0.0.1:32925"


legacy_flag=" --legacy "
clean_up="true"
gas_limit=1000000

function normalize_address() {
        sed 's/0x//' |
         tr '[:upper:]' '[:lower:]'
}
function hex_to_dec() {
    hex_in=$(sed 's/0x//' | tr '[:lower:]' '[:upper:]')
    dec_val=$(bc <<< "ibase=16; $hex_in")
    echo $dec_val
}

function process_test_item() {
    local testfile=$1
    local test_hash=$(sha256sum $testfile | sed 's/ .*//')
    if [[ -e "/tmp/.retest-resume-$test_hash" ]]; then
        2>&1 echo "it looks like we have already tested this case. Skipping"
        return
    fi

    touch /tmp/.retest-resume-$test_hash

    local test_name=$( jq -r '.testCases[0].name' $testfile)
    2>&1 echo "processing $test_name in file at $testfile"

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

    if [[ -e last.nonce ]]; then
        # Random transaction to make sure all of the async deps are deployed before running the transactions
        cast send $legacy_flag --nonce $(cat last.nonce) --rpc-url $rpc_url --private-key $private_key --value 1 $wallet_address
        2>&1 echo "We have finished deploying the dependencies (I think)"
        rm last.nonce
    fi

    nonce=$(cast nonce --rpc-url $rpc_url $wallet_address)
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
            dec_val=$(echo $val | hex_to_dec)
            val_arg=" --value $dec_val "
        fi

        gas_arg=""
        if [[ $gas != "" ]] ; then
            dec_val=$(echo $gas | hex_to_dec)
            valid_gas=$(bc <<< "$dec_val < 30000000 && $dec_val > 0")
            if [[ $valid_gas == "1" ]] ; then
                gas_arg=" --gas-limit $dec_val "
            else
                gas_arg=" --gas-limit $gas_limit "
            fi
        fi

        local to_addr_arg=""
        if [[ $addr == "0x0000000000000000000000000000000000000000" || $addr == "" || $addr == "0000000000000000000000000000000000000000" ]] ; then
            if [[ $tx_input == "" ]]; then
                2>&1 echo "The test $name case $count seems to have a create with an empty data... skiping"
                continue
            fi
            to_addr_arg=" --create "
        else
            if [[ ! -e $addr.actual ]]; then
                2>&1 "the test file $addr.actual does not seem to exist... skipping"
                continue
            fi
            resolved_address=$(cat $addr.actual)
            to_addr_arg=" $resolved_address "
        fi

        set -x
        timeout 30 cast send $legacy_flag --async --nonce $nonce --rpc-url $rpc_url --private-key $private_key $gas_arg $val_arg $to_addr_arg $tx_input | tee tx-$count-out.json
        set +x
        if [[ $? -ne 0 ]]; then
            2>&1 "it looks like this request timed out.. it might be worth checking?!"
        fi
        nonce=$((nonce+1))
        echo "$nonce" > last.nonce
    done

    if [[ -e last.nonce ]]; then
        # Random transaction to make sure all of the async deps are deployed before running the transactions
        cast send $legacy_flag --nonce $(cat last.nonce) --rpc-url $rpc_url --private-key $private_key --value 1 $wallet_address
        rm last.nonce
    fi

    popd
    if [[ $clean_up == "true" ]] ; then
        rm -rf $tmp_dir
        rm $testfile
    fi
}

wallet_address=$(cast wallet address --private-key $private_key)
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

