This command will attempt to send a claim transaction to the bridge contract.

```solidity
    /**
     * @notice Bridge message and send ETH value
     * note User/UI must be aware of the existing/available networks when choosing the destination network
     * @param destinationNetwork Network destination
     * @param destinationAddress Address destination
     * @param forceUpdateGlobalExitRoot Indicates if the new global exit root is updated or not
     * @param metadata Message metadata
     */
    function bridgeMessage(
        uint32 destinationNetwork,
        address destinationAddress,
        bool forceUpdateGlobalExitRoot,
        bytes calldata metadata
    );

```

Each transaction will require manual input of parameters. Example usage:

```bash
polycli ulxly deposit-claim \
        --bridge-address 0xD71f8F956AD979Cc2988381B8A743a2fE280537D \
        --private-key 12d7de8621a77640c9241b2595ba78ce443d05e94090365ab3bb5e19df82c625 \
        --claim-index 0 \
        --claim-address 0xE34aaF64b29273B7D567FCFc40544c014EEe9970 \
        --claim-network 0 \
        --rpc-url http://127.0.0.1:32790 \
        --bridge-service-url http://127.0.0.1:32804
```

This command would use the supplied private key and attempt to send a claim transaction to the bridge contract address with the input flags.
Successful deposit transaction will output logs like below:

```bash
Claim Transaction Successful: 0x7180201b19e1aa596503d8541137d6f341e682835bf7a54aab6422c89158866b
```

Upon successful claim, the transferred funds can be queried in the destination network using tools like `cast balance <claim-address> --rpc-url <destination-network-url>`


Failed deposit transactions will output logs like below: 

```bash
Claim Transaction Failed: 0x32ac34797159c79e57ae801c350bccfe5f8105d4dd3b717e31d811397e98036a
```

The reason for failing may be very difficult to debug. I have personally spun up a bridge-ui and compared the byte data of a successful transaction to the byte data of a failing claim transaction queried using:

```!
curl http://127.0.0.1:32790 \
-X POST \
-H "Content-Type: application/json" \
--data '{"method":"debug_traceTransaction","params":["0x32ac34797159c79e57ae801c350bccfe5f8105d4dd3b717e31d811397e98036a", {"tracer": "callTracer"}], "id":1,"jsonrpc":"2.0"}' | jq '.'
```
