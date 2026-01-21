This command will directly attempt to make a deposit on the uLxLy bridge. This call responds to the method defined below:

```solidity
/**
 * @notice Deposit add a new leaf to the merkle tree
 * note If this function is called with a reentrant token, it would be possible to `claimTokens` in the same call
 * Reducing the supply of tokens on this contract, and actually locking tokens in the contract.
 * Therefore we recommend to third parties bridges that if they do implement reentrant call of `beforeTransfer` of some reentrant tokens
 * do not call any external address in that case
 * note User/UI must be aware of the existing/available networks when choosing the destination network
 * @param destinationNetwork Network destination
 * @param destinationAddress Address destination
 * @param amount Amount of tokens
 * @param token Token address, 0 address is reserved for ether
 * @param forceUpdateGlobalExitRoot Indicates if the new global exit root is updated or not
 * @param permitData Raw data of the call `permit` of the token
 */
function bridgeAsset(
    uint32 destinationNetwork,
    address destinationAddress,
    uint256 amount,
    address token,
    bool forceUpdateGlobalExitRoot,
    bytes calldata permitData
) public payable virtual ifNotEmergencyState nonReentrant {
```

The source of this method is [here](https://github.com/0xPolygonHermez/zkevm-contracts/blob/c8659e6282340de7bdb8fdbf7924a9bd2996bc98/contracts/v2/PolygonZkEVMBridgeV2.sol#L198-L219).
Below is an example of how we would make simple bridge of native ETH from Sepolia (L1) into Cardona (L2).

```bash
polycli ulxly bridge asset \
    --bridge-address 0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582 \
    --private-key 0x32430699cd4f46ab2422f1df4ad6546811be20c9725544e99253a887e971f92b \
    --destination-network 1 \
    --value 10000000000000000 \
    --rpc-url https://sepolia.drpc.org
```

[This](https://sepolia.etherscan.io/tx/0xf57b8171b2f62dce3eedbe3e50d5ee8413d61438af64286b5017ed9d5d154816) is the transaction that was created and mined from running this command.

Here is another example that will bridge a [test ERC20 token](https://sepolia.etherscan.io/address/0xC92AeF5873d058a76685140F3328B0DED79733Af) from Sepolia (L1) into Cardona (L2). In order for this to work, the token would need to have an [approval](https://sepolia.etherscan.io/tx/0x028513b13a2a7899de4db56e60d1dad66c7b7e29f91c54f385fdfdfc8f14b8b4#eventlog) for the bridge to spend tokens for that particular user.

```bash
polycli ulxly bridge asset \
    --bridge-address 0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582 \
    --private-key 0x32430699cd4f46ab2422f1df4ad6546811be20c9725544e99253a887e971f92b \
    --destination-network 1 \
    --value 10000000000000000 \
    --token-address 0xC92AeF5873d058a76685140F3328B0DED79733Af \
    --destination-address 0x3878Cff9d621064d393EEF92bF1e12A944c5ba84 \
    --rpc-url https://sepolia.drpc.org
```

[This](https://sepolia.etherscan.io/tx/0x8ed1c2c0f2e994c86867f401c86fea3c709a28a18629d473cf683049f176fa93) is the transaction that was created and mined from running this command.

Assuming you have funds on L2, a bridge from L2 to L1 looks pretty much the same.
The command below will bridge `123456` of the native ETH on Cardona (L2) back to network 0 which corresponds to Sepolia (L1).

```bash
polycli ulxly bridge asset \
    --bridge-address 0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582 \
    --private-key 0x32430699cd4f46ab2422f1df4ad6546811be20c9725544e99253a887e971f92b \
    --destination-network 0 \
    --value 123456 \
    --destination-address 0x3878Cff9d621064d393EEF92bF1e12A944c5ba84 \
    --rpc-url https://rpc.cardona.zkevm-rpc.com
```

[This](https://cardona-zkevm.polygonscan.com/tx/0x0294dae3cfb26881e5dde9f182531aa5be0818956d029d50e9872543f020df2e) is the transaction that was created and mined from running this command.
