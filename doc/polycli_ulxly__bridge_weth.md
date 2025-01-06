# `polycli ulxly  bridge weth`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

For L2's that use a gas token, use this to transfer WETH to another chain

```bash
polycli ulxly  bridge weth [flags]
```

## Usage

This command is not used very often but can be used on L2 networks that have a gas token.

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
) external ifNotEmergencyState {
```
[Here](https://github.com/0xPolygonHermez/zkevm-contracts/blob/c8659e6282340de7bdb8fdbf7924a9bd2996bc98/contracts/v2/PolygonZkEVMBridgeV2.sol#L352-L367) is the source code that corresponds to this interface.

Assuming the network is configured with a gas token, you could call this method like this:

```bash
polycli ulxly bridge weth \
        --bridge-address 0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582 \
        --destination-address 0x3878Cff9d621064d393EEF92bF1e12A944c5ba84 \
        --private-key 0x32430699cd4f46ab2422f1df4ad6546811be20c9725544e99253a887e971f92b \
        --value 123456 \
        --destination-network 1 \
        --rpc-url http://l2-rpc-url.invalid \
        --token-address $WETH_ADDRESS
```


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

- [polycli ulxly  bridge](polycli_ulxly__bridge.md) - Commands for moving funds and sending messages from one chain to another
