#!/usr/bin/env bash

function get_polycli_arg() {

    if [[ "$2" == "-1" ]]; then
        return
    fi

    if [[ "$2" == "UNSPECIFIED" ]]; then
        return
    fi


    case "$1" in
        accountFundingAmount)
            echo -n "--account-funding-amount $2"
            ;;
        adaptiveBackoffFactor)
            echo -n "--adaptive-backoff-factor $2"
            ;;
        adaptiveCycleDurationSeconds)
            echo -n "--adaptive-cycle-duration-seconds $2"
            ;;
        adaptiveRateLimit)
            if [[ "$2" == "TRUE" ]] ; then
                echo -n "--adaptive-rate-limit"
            fi
            ;;
        adaptiveRateLimitIncrement)
            echo -n "--adaptive-rate-limit-increment $2"
            ;;
        adaptiveTargetSize)
            echo -n "--adaptive-target-size $2"
            ;;
        batchSize)
            echo -n "--batch-size $2"
            ;;
        blobFeeCap)
            echo -n "--blob-fee-cap $2"
            ;;
        calldata)
            echo -n "--calldata $2"
            ;;
        chainId)
            echo -n "--chain-id $2"
            ;;
        concurrency)
            echo -n "--concurrency $2"
            ;;
        contractAddress)
            # For now we'll tightly couple contract address and calldata.. The current contract address params are random so any call data will work
            echo -n "--contract-address $2 --calldata 0xa0712d680000000000000000000000000000000000000000000000000000000000000001"
            ;;
        contractCallPayable)
            if [[ "$2" == "TRUE" ]] ; then
                echo -n "--contract-call-payable"
            fi
            ;;
        erc20Address)
            echo -n "--erc20-address $2"
            ;;
        erc721Address)
            echo -n "--erc721-address $2"
            ;;
        ethAmountInWei)
            echo -n "--eth-amount-in-wei $2"
            ;;
        ethCallOnly)
            if [[ "$2" == "TRUE" ]] ; then
                echo -n "--eth-call-only"
            fi
            ;;
        ethCallOnlyLatest)
            if [[ "$2" == "TRUE" ]] ; then
                echo -n "--eth-call-only-latest"
            fi
            ;;
        fireAndForget)
            if [[ "$2" == "TRUE" ]] ; then
                echo -n "--fire-and-forget"
            fi
            ;;
        gasLimit)
            echo -n "--gas-limit $2"
            ;;
        gasPrice)
            echo -n "--gas-price $2"
            ;;
        gasPriceMultiplier)
            echo -n "--gas-price-multiplier $2"
            ;;
        legacy)
            if [[ "$2" == "TRUE" ]] ; then
                echo -n "--legacy"
            fi
            ;;
        loadtestContractAddress)
            echo -n "--loadtest-contract-address $2"
            ;;
        maxBaseFeeWei)
            echo -n "--max-base-fee-wei $2"
            ;;
        mode)
            echo -n "--mode $2"
            ;;
        nonce)
            echo -n "--nonce $2"
            ;;
        outputMode)
            echo -n "--output-mode $2"
            ;;
        outputRawTxOnly)
            if [[ "$2" == "TRUE" ]] ; then
                echo -n "--output-raw-tx-only"
            fi
            ;;
        preFundSendingAccounts)
            if [[ "$2" == "TRUE" ]] ; then
                echo -n "--pre-fund-sending-accounts"
            fi
            ;;
        priorityGasPrice)
            echo -n "--priority-gas-price $2"
            ;;
        privateKey)
            echo -n "--private-key $2"
            ;;
        proxy)
            echo -n "--proxy $2"
            ;;
        randomRecipients)
            if [[ "$2" == "TRUE" ]] ; then
                echo -n "--random-recipients"
            fi
            ;;
        rateLimit)
            echo -n "--rate-limit $2"
            ;;
        recallBlocks)
            echo -n "--recall-blocks $2"
            ;;
        receiptRetryInitialDelayMs)
            echo -n "--receipt-retry-initial-delay-ms $2"
            ;;
        receiptRetryMax)
            echo -n "--receipt-retry-max $2"
            ;;
        refundRemainingFunds)
            if [[ "$2" == "TRUE" ]] ; then
                echo -n "--refund-remaining-funds"
            fi
            ;;
        requests)
            echo -n "--requests $2"
            ;;
        rpcUrl)
            echo -n "--rpc-url $2"
            ;;
        seed)
            echo -n "--seed $2"
            ;;
        sendOnly)
            if [[ "$2" == "TRUE" ]] ; then
                echo -n "--send-only"
            fi
            ;;
        sendingAccountsCount)
            echo -n "--sending-accounts-count $2"
            ;;
        sendingAccountsFile)
            echo -n "--sending-accounts-file $2"
            ;;
        storeDataSize)
            echo -n "--store-data-size $2"
            ;;
        summarize)
            if [[ "$2" == "TRUE" ]] ; then
                echo -n "--summarize"
            fi
            ;;
        timeLimit)
            echo -n "--time-limit $2"
            ;;
        toAddress)
            echo -n "--to-address $2"
            ;;
        waitForReceipt)
            if [[ "$2" == "TRUE" ]] ; then
                echo -n "--wait-for-receipt"
            fi
            ;;
        *)
            echo "I do not recognize $1"
            exit 1
            ;;
    esac
}

cur_cmd=""
while read line ; do
    if [[ $cur_cmd == "" ]]; then
        cur_cmd="polycli loadtest --verbosity 700"
    fi

    param_re="^[0-9]* = (.*)=(.*)$"
    if [[ $line =~ $param_re ]] ; then
        # echo $line
        cur_arg=$(get_polycli_arg ${BASH_REMATCH[1]} ${BASH_REMATCH[2]})
        cur_cmd="$cur_cmd $cur_arg"
    fi
    if [[ $line == "-------------------------------------" ]]; then
        echo $cur_cmd | sed 's/--/\\\n\t--/gi' 1>&2
        echo $cur_cmd
        cur_cmd=""
    fi
done < "${1:-/dev/stdin}"
