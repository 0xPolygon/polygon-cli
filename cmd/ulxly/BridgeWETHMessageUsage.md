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

