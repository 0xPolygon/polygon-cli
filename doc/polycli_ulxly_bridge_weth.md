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
      --bridge-address string              address of the lxly bridge
      --call-data string                   call data to be passed directly with bridge-message or as an ERC20 Permit (default "0x")
      --call-data-file string              a file containing hex encoded call data
      --chain-id string                    chain ID to use in the transaction
      --config string                      config file (default is $HOME/.polygon-cli.yaml)
      --destination-address string         destination address for the bridge
      --destination-network uint32         rollup ID of the destination network
      --dry-run                            do all of the transaction steps but do not send the transaction
      --force-update-root                  update the new global exit root (default true)
      --gas-limit uint                     force specific gas limit for transaction
      --gas-price string                   gas price to use
      --insecure                           skip TLS certificate verification
      --pretty-logs                        output logs in pretty format instead of JSON (default true)
      --private-key string                 hex encoded private key for sending transaction
      --rpc-url string                     RPC URL to send the transaction
      --token-address string               address of ERC20 token to use (default "0x0000000000000000000000000000000000000000")
      --transaction-receipt-timeout uint   timeout in seconds to wait for transaction receipt confirmation (default 60)
      --value string                       amount in wei to send with the transaction (default "0")
  -v, --verbosity int                      0 - silent
                                           100 panic
                                           200 fatal
                                           300 error
                                           400 warning
                                           500 info
                                           600 debug
                                           700 trace (default 500)
```

## See also

- [polycli ulxly bridge](polycli_ulxly_bridge.md) - Commands for moving funds and sending messages from one chain to another
