# `polycli ulxly deposit-claim`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Make a uLxLy claim transaction

```bash
polycli ulxly deposit-claim [flags]
```

## Usage

This command will attempt to send a claim transaction to the bridge contract.

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
     * @param originTokenAddress  Origin token address, 0 address is reserved for ether
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
    )

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
      --bridge-address string                The address of the bridge contract.
      --bridge-service-url string            The RPC endpoint of the bridge service component.
      --chain-id string                      The chainID.
      --claim-address string                 The address that is receiving the bridged asset.
      --claim-index string                   The deposit count, or index to initiate a claim transaction for. (default "0")
      --claim-message                        Claim a message instead of an asset.
      --claim-weth                           Claim a weth instead of an asset.
      --destination-network string           The network ID of the destination network. (default "1")
      --gas-limit uint                       The gas limit for the transaction. Setting this value to 0 will estimate the gas limit.
  -h, --help                                 help for deposit-claim
      --origin-network string                The network ID of the origin network. (default "0")
      --private-key string                   The private key of the sender account.
      --rpc-url string                       The RPC endpoint of the destination network (default "http://127.0.0.1:8545")
      --transaction-receipt-timeout uint32   The timeout limit to check for the transaction receipt of the claim. (default 60)
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Panic
                        200 Fatal
                        300 Error
                        400 Warning
                        500 Info
                        600 Debug
                        700 Trace (default 500)
```

## See also

- [polycli ulxly](polycli_ulxly.md) - Utilities for interacting with the lxly bridge
