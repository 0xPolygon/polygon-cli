This command will attempt to send a deposit transaction to the bridge contract.

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
    );

```

Each transaction will require manual input of parameters. Example usage:

```bash
polycli ulxly deposit-new \
        --private-key 12d7de8621a77640c9241b2595ba78ce443d05e94090365ab3bb5e19df82c625 \
        --gas-limit 300000 \
        --amount 1000000000000000000 \
        --rpc-url http://127.0.0.1:8545 \
        --bridge-address 0xD71f8F956AD979Cc2988381B8A743a2fE280537D \
        --destination-network 1 \
        --destination-address 0xE34aaF64b29273B7D567FCFc40544c014EEe9970
```
go run ulxly.go ulxly bridge-asset --private-key 846f556e5509c36b5756fa3741009310908d606db18d5916d43e6c7d22e0f5a7 --gas-limit 300000 --rpc-url https://sepolia.infura.io/v3/7c490aa772fa466293ee1ac62761330a --bridge-address 0x1348947e282138d8f377b467f7d9c2eb0f335d1f --destination-network 1 --amount 19

go run main.go ulxly bridge-message --private-key 846f556e5509c36b5756fa3741009310908d606db18d5916d43e6c7d22e0f5a7 --gas-limit 300000 --rpc-url https://sepolia.infura.io/v3/7c490aa772fa466293ee1ac62761330a --bridge-address 0x1348947e282138d8f377b467f7d9c2eb0f335d1f --destination-network 1 --amount 19 --call-data 0x18976546789087654356789654345678976543567897654356789765435468797865435468797654354678907654356879765436798765435467986543

go run main.go ulxly claim-asset --private-key 846f556e5509c36b5756fa3741009310908d606db18d5916d43e6c7d22e0f5a7 --gas-limit 30000000 --rpc-url https://sepolia.infura.io/v3/7c490aa772fa466293ee1ac62761330a --bridge-address 0x1348947e282138d8f377b467f7d9c2eb0f335d1f --claim-address 0x2536C2745Ac4A584656A830f7bdCd329c94e8F30 --claim-index 7780 --bridge-service-url https://bridge-api.internal.zkevm-rpc.com --deposit-network 1

go run main.go ulxly claim-message --private-key 846f556e5509c36b5756fa3741009310908d606db18d5916d43e6c7d
22e0f5a7 --gas-limit 300000 --rpc-url https://sepolia.infura.io/v3/7c490aa772fa466293ee1ac62761330a --bridge-address 0x1348947e282138d8f377b467f7d9c2eb0f335d1f  --claim-address 0x2536C2745Ac4A584656A830f7bdCd329c94e8F30 --claim-index 7780 --call-data 0x18976546789087654356789654345678976543567897654356789765435468797865435468797654354678907654356879765436798765435467986543 --bridge-service-url https://bridge-api.internal.zkevm-rpc.com


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
