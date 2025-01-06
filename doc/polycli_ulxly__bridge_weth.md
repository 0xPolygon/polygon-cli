# `polycli ulxly  bridge weth`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

send some WETH into the bridge

```bash
polycli ulxly  bridge weth [flags]
```

## Usage

This command will attempt to send a deposit transaction to the bridge contract.

```solidity
    /**
     * @notice Bridge message and send ETH value
     * note User/UI must be aware of the existing/available networks when choosing the destination network
     * @param destinationNetwork Network destination
     * @param destinationAddress Address destination
     * @param amountWETH Amount of WETH tokens
     * @param forceUpdateGlobalExitRoot Indicates if the new global exit root is updated or not
     * @param metadata Message metadata
     */
    function bridgeMessageWETH(
        uint32 destinationNetwork,
        address destinationAddress,
        uint256 amountWETH,
        bool forceUpdateGlobalExitRoot,
        bytes calldata metadata
    );
```

Each transaction will require manual input of parameters. Example usage:

```bash
polycli ulxly bridge-message-weth \
        --private-key 12d7de8621a77640c9241b2595ba78ce443d05e94090365ab3bb5e19df82c625 \
        --gas-limit 300000 \
        --amount 1000000000000000000 \
        --rpc-url http://127.0.0.1:8545 \
        --bridge-address 0xD71f8F956AD979Cc2988381B8A743a2fE280537D \
        --destination-network 1 \
        --destination-address 0xE34aaF64b29273B7D567FCFc40544c014EEe9970
        --call-data 0x001010109200090028979743971976836486868648629808961824738090896826764980866fac97863898ca08928fc7279643
```

This command would use the supplied private key and attempt to send a deposit transaction to the bridge contract address with the input flags.
Successful deposit transaction will output logs like below:

```bash
Deposit Transaction Successful: 0x8c9b82e8abdfb4aad5fccd91879397acfa73e4261282c8dc634734d05ad889d3
```

Upon successful deposit, the transaction can be queried using `polycli ulxly deposit-get` command


Failed deposit transactions will output logs like below: 

```bash
Deposit Transaction Failed: 0x60385209b0e9db359c24c88c2fb8a5c9e4628fffe8d5fb2b5e64dfac3a2b7639
Try increasing the gas limit:
Current gas limit: 100000
Cumulative gas used for transaction: 98641
```

The reason for failing may likely be due to the `out of gas` error. Increasing the `--gas-limit` flag value will likely resolve this. 

## Flags

```bash
  -h, --help   help for weth
```

The command also inherits flags from parent commands.

```bash
      --bridge-address string              the address of the lxly bridge
      --call-data string                   call data to be passed directly with bridge-message or as an ERC20 Permit (default "0x")
      --chain-id string                    set the chain id to be used in the transaction
      --config string                      config file (default is $HOME/.polygon-cli.yaml)
      --destination-address string         the address where the bridge will be sent to
      --destination-network uint32         the rollup id of the destination network
      --dry-run                            do all of the transaction steps but do not send the transaction
      --force-update-root                  indicates if the new global exit root is updated or not (default true)
      --gas-limit uint                     force a gas limit when sending a transaction
      --gas-price string                   the gas price to be used
      --pretty-logs                        Should logs be in pretty format or JSON (default true)
      --private-key string                 the hex encoded private key to be used when sending the tx
      --rpc-url string                     the URL of the RPC to send the transaction
      --token-address string               the address of an ERC20 token to be used (default "0x0000000000000000000000000000000000000000")
      --transaction-receipt-timeout uint   the amount of time to wait while trying to confirm a transaction receipt (default 60)
      --value string                       the amount in wei to be sent along with the transaction
  -v, --verbosity int                      0 - Silent
                                           100 Panic
                                           200 Fatal
                                           300 Error
                                           400 Warning
                                           500 Info
                                           600 Debug
                                           700 Trace (default 500)
```

## See also

- [polycli ulxly  bridge](polycli_ulxly__bridge.md) - commands for making deposits to the uLxLy bridge
