# `polycli ulxly  claim message`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

perform a claim of a given message in the bridge

```bash
polycli ulxly  claim message [flags]
```

## Usage

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

## Flags

```bash
  -h, --help   help for message
```

The command also inherits flags from parent commands.

```bash
      --bridge-address string              the address of the lxly bridge
      --bridge-service-url string          the URL of the bridge service
      --chain-id string                    set the chain id to be used in the transaction
      --config string                      config file (default is $HOME/.polygon-cli.yaml)
      --deposit-count uint                 the deposit count of the bridge transaction
      --deposit-network uint               the rollup id of the network where the bridge is being claimed
      --destination-address string         the address where the bridge will be sent to
      --dry-run                            do all of the transaction steps but do not send the transaction
      --gas-limit uint                     force a gas limit when sending a transaction
      --gas-price string                   the gas price to be used
      --global-index string                an override of the global index value
      --pretty-logs                        Should logs be in pretty format or JSON (default true)
      --private-key string                 the hex encoded private key to be used when sending the tx
      --rpc-url string                     the URL of the RPC to send the transaction
      --transaction-receipt-timeout uint   the amount of time to wait while trying to confirm a transaction receipt (default 60)
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

- [polycli ulxly  claim](polycli_ulxly__claim.md) - commands for making claims of deposits from the uLxLy bridge
