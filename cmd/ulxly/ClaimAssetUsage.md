This command will connect to the bridge service, generate a proof, and then attempt to claim the deposit on which never network is referred to in the `--rpc-url` argument.
This is the corresponding interface in the bridge contract:

```solidity
/**
 * @notice Verify merkle proof and withdraw tokens/ether
 * @param smtProofLocalExitRoot Smt proof to proof the leaf against the network exit root
 * @param smtProofRollupExitRoot Smt proof to proof the rollupLocalExitRoot against the rollups exit root
 * @param globalIndex Global index is defined as:
 * | 191 bits |    1 bit     |   32 bits   |     32 bits    |
 * |    0     |  mainnetFlag | rollupIndex | localRootIndex |
 * note that only the rollup index will be used only in case the mainnet flag is 0
 * note that global index do not assert the unused bits to 0.
 * This means that when synching the events, the globalIndex must be decoded the same way that in the Smart contract
 * to avoid possible synch attacks
 * @param mainnetExitRoot Mainnet exit root
 * @param rollupExitRoot Rollup exit root
 * @param originNetwork Origin network
 * @param originTokenAddress  Origin token address, 0 address is reserved for gas token address. If WETH address is zero, means this gas token is ether, else means is a custom erc20 gas token
 * @param destinationNetwork Network destination
 * @param destinationAddress Address destination
 * @param amount Amount of tokens
 * @param metadata Abi encoded metadata if any, empty otherwise
 */
function claimAsset(
    bytes32[_DEPOSIT_CONTRACT_TREE_DEPTH] calldata smtProofLocalExitRoot,
    bytes32[_DEPOSIT_CONTRACT_TREE_DEPTH] calldata smtProofRollupExitRoot,
    uint256 globalIndex,
    bytes32 mainnetExitRoot,
    bytes32 rollupExitRoot,
    uint32 originNetwork,
    address originTokenAddress,
    uint32 destinationNetwork,
    address destinationAddress,
    uint256 amount,
    bytes calldata metadata
) external ifNotEmergencyState {
```

[Here](https://github.com/0xPolygonHermez/zkevm-contracts/blob/c8659e6282340de7bdb8fdbf7924a9bd2996bc98/contracts/v2/PolygonZkEVMBridgeV2.sol#L433-L465) is a direct link to the source code as well.

In order to claim an asset or a message, you need to know deposit count. Usually this is in the event data of the transaction. Alternatively, you can usually directly attempt to see the pending deposits by querying the bridge API directly. In the case of Cardona, the bridge service is running here: https://bridge-api.cardona.zkevm-rpc.com

```bash
curl -s https://bridge-api.cardona.zkevm-rpc.com/bridges/0x3878Cff9d621064d393EEF92bF1e12A944c5ba84 | jq '.'
```

In the output of the above command, I can see a deposit that looks like this:
```json
{
  "leaf_type": 0,
  "orig_net": 0,
  "orig_addr": "0x0000000000000000000000000000000000000000",
  "amount": "123456",
  "dest_net": 0,
  "dest_addr": "0x3878Cff9d621064d393EEF92bF1e12A944c5ba84",
  "block_num": "9695587",
  "deposit_cnt": 9075,
  "network_id": 1,
  "tx_hash": "0x0294dae3cfb26881e5dde9f182531aa5be0818956d029d50e9872543f020df2e",
  "claim_tx_hash": "",
  "metadata": "0x",
  "ready_for_claim": true,
  "global_index": "9075"
}
```

If we want to claim this deposit, we can use a command like this:

```bash
polycli ulxly claim asset \
    --bridge-address 0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582 \
    --bridge-service-url https://bridge-api.cardona.zkevm-rpc.com \
    --private-key 0x32430699cd4f46ab2422f1df4ad6546811be20c9725544e99253a887e971f92b \
    --deposit-network 1 \
    --deposit-count 9075 \
    --rpc-url https://sepolia.drpc.org
```

[Here](https://sepolia.etherscan.io/tx/0x21fee6e47a3b6733034fb963b20fe7accb0fb168257450f8f0053d6af8e4bc76) is the transaction that was created and mined based on this command.