#!/bin/bash
# ./scripts/run-retest-with-cast.sh < simple-test-out-new.json 2>&1 | tee -a local-test-aug-28-2.logs
# find /tmp -type f -newer /tmp/.retest-resume-a01b837809ca1555757ab2edebead6321eaf569c78a54fc81985b84be71eacb0 -name '.retest-resume-*' | xargs rm

private_key="0x12d7de8621a77640c9241b2595ba78ce443d05e94090365ab3bb5e19df82c625"
rpc_url="http://127.0.0.1:33237"


legacy_flag=" --legacy "
clean_up="true"
gas_limit=1000000

function print_warning() {
    2>&1 echo -e "\e[32m$1\e[0m"
}

function normalize_address() {
        sed 's/0x//' |
         tr '[:upper:]' '[:lower:]'
}

function hex_to_dec() {
    hex_in=$(sed 's/0x//' | tr '[:lower:]' '[:upper:]')
    dec_val=$(bc <<< "ibase=16; $hex_in")
    echo "$dec_val"
}

function mark_progress() {
    local bn
    local cur_time
    local test_name
    local test_file
    local test_hash
    local test_counter

    test_name=$1
    test_file=$2
    test_hash=$3
    test_counter=$4

    bn=$(cast block-number --rpc-url "$rpc_url")
    cur_time=$(date -R)
    2>&1 printf "\n\n"
    2>&1 echo "################################################################################"
    2>&1 echo "Starting test #$test_counter $test_name at block $bn at $cur_time"
    2>&1 echo "Test source $test_file with test lock /tmp/.retest-resume-$test_hash"
    2>&1 echo "################################################################################"
}

# We'll attempt to send a synchronous transaction. If that works
# (doesn't time out), it means that the test is no longer pending. If
# it failed, it means that we might need to clear out some pending
# transactions. If we don't do this, it's very easy for one test to
# intere with the excution of the next test.
function clear_pending_txs() {
    local last_nonce
    local current_nonce

    last_nonce=$1
    # shellcheck disable=SC2086
    timeout 30 cast send $legacy_flag --nonce "$(cat last.nonce)" --rpc-url "$rpc_url" --private-key "$private_key" --value 1 "$wallet_address"
    ret_code=$?

    if [[ $ret_code -eq 0 ]]; then
        return
    fi

    print_warning "The transaction to clear pending txs is stuck.. attemping to clear all stuck transaction. This means the previous test did not execute properly"
    current_nonce=$(cast nonce --rpc-url "$rpc_url" "$wallet_address")
    print_warning "Attemping relacements from $current_nonce to $last_nonce"

    for ((i = current_nonce ; i <= last_nonce ; i++)); do
        # shellcheck disable=SC2086
        cast send $legacy_flag --nonce "$i" --gas-price 100gwei --rpc-url "$rpc_url" --private-key "$private_key" --value 1 "$wallet_address"
    done
}

function process_test_item() {
    local testfile
    local test_hash
    local test_name
    local tmp_dir
    local nonce
    local count
    local test_counter

    testfile=$1
    test_counter=$2

    test_hash=$(sha256sum "$testfile" | sed 's/ .*//')
    if [[ -e "/tmp/.retest-resume-$test_hash" ]]; then
        2>&1 echo "it looks like we have already tested this case. Skipping"
        return
    fi

    touch "/tmp/.retest-resume-$test_hash"

    test_name=$(jq -r '.testCases[0].name' "$testfile")
    mark_progress "$test_name" "$testfile" "$test_hash" "$test_counter"

    tmp_dir=$(mktemp -p /tmp -d retest-work-XXXXXXXXXXXX)
    pushd "$tmp_dir" || exit 1

    nonce=$(cast nonce --rpc-url "$rpc_url" "$wallet_address")
    count=0
    jq -c '.dependencies[]' "$testfile" | while read -r pre ; do
        local reference_address
        local code_to_deploy

        count=$((count+1))
        echo "$pre" | jq '.' > "dep-$count.json"
        code_to_deploy=$(jq -r '.code' "dep-$count.json")
        reference_address=$(jq -r '.addr' "dep-$count.json" | normalize_address)
        2>&1 echo "deploying dependency $count for $reference_address"
        2>&1 echo "current nonce: $nonce"
        # shellcheck disable=SC2086
        cast send $legacy_flag: --async --nonce "$nonce" --rpc-url "$rpc_url" --private-key "$private_key" --create "$code_to_deploy" | tee "$reference_address.txhash"
        cast compute-address --nonce "$nonce" "$wallet_address" | sed 's/^.*0x/0x/' > "$reference_address.actual"

        nonce=$((nonce+1))
        echo "$nonce" > last.nonce
    done

    if [[ -e last.nonce ]]; then
        clear_pending_txs "$(cat last.nonce)"
        rm last.nonce
    fi

    nonce=$(cast nonce --rpc-url "$rpc_url" "$wallet_address")
    count=0
    jq -c '.testCases[]' "$testfile" | while read -r "test_case" ; do
        local name
        local addr
        local gas
        local val

        count=$((count+1))
        echo "$test_case" | jq '.' > "test_case_$count.json"
        tx_input=$(jq -r '.input' test_case_$count.json)
        name=$(jq -r '.name' test_case_$count.json)
        addr=$(jq -r '.to' test_case_$count.json | normalize_address)
        gas=$(jq -r '.gas' test_case_$count.json) # this value can be obscenely high in the test cases
        val=$(jq -r '.value' test_case_$count.json)
        val_arg=""
        if [[ $val != "0x0" ]] ; then
            dec_val=$(echo "$val" | hex_to_dec)
            val_arg=" --value $dec_val "
        fi

        gas_arg=""
        if [[ $gas != "" ]] ; then
            dec_val=$(echo "$gas" | hex_to_dec)
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
                print_warning "The test $name case $count seems to have a create with an empty data... skiping"
                continue
            fi
            to_addr_arg=" --create "
        else
            if [[ ! -e $addr.actual ]]; then
                print_warning "the test file $addr.actual does not seem to exist... skipping"
                continue
            fi
            resolved_address=$(cat "$addr.actual")
            to_addr_arg=" $resolved_address "
        fi

        2>&1 echo "executing tx $count for $name to alias of $addr"
        2>&1 echo "current nonce: $nonce"

        set -x
        # shellcheck disable=SC2086
        timeout 30 cast send $legacy_flag --async --nonce "$nonce" --rpc-url "$rpc_url" --private-key "$private_key" $gas_arg $val_arg $to_addr_arg $tx_input | tee "tx-$count-out.json"
        ret_code=$?
        set +x
        if [[ $ret_code -ne 0 ]]; then
            print_warning "it looks like this request timed out.. it might be worth checking?!"
        fi
        nonce=$((nonce+1))
        echo "$nonce" > last.nonce
    done

    if [[ -e last.nonce ]]; then
        clear_pending_txs "$(cat last.nonce)"
        rm last.nonce
    fi

    popd || exit 1

    if [[ $clean_up == "true" ]] ; then
        rm -rf "$tmp_dir"
        rm "$testfile"
    fi
}

wallet_address=$(cast wallet address --private-key "$private_key")
wallet_balance=$(cast balance --rpc-url $rpc_url "$wallet_address")
wallet_nonce=$(cast nonce --rpc-url $rpc_url "$wallet_address")
2>&1 echo "Address $wallet_address has a balance of $wallet_balance and nonce $wallet_nonce"

test_counter=0
# Break down each test into different files
jq -c '.[]' | while read -r test_item ; do
    testfile=$(mktemp -p /tmp retest-item-jq-XXXXXXXXXXXX)
    echo "$test_item" > "$testfile"
    test_counter=$((test_counter+1))
    process_test_item "$testfile" "$test_counter"
done

