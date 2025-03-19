# `polycli ulxly  bridge message`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Send some ETH along with data from one chain to another chain

```bash
polycli ulxly  bridge message [flags]
```

## Usage

This command is very similar to `polycli ulxly bridge asset`, but instead this is a more generic interface that can be used to transfer ETH and make a contract call. This is the underlying solidity interface that we're  referencing.

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
) external payable ifNotEmergencyState {
```

The source code for this particular method is [here](https://github.com/0xPolygonHermez/zkevm-contracts/blob/c8659e6282340de7bdb8fdbf7924a9bd2996bc98/contracts/v2/PolygonZkEVMBridgeV2.sol#L324-L337).

Below is a simple example of using this command to bridge a small amount of ETH from Sepolia (L1) to Cardona (L2). In this case, we're not including any call data, so it's essentially equivalent to a `bridge asset` call, but the deposit will not be automatically claimed on L2.

```bash
polycli ulxly bridge message \
    --bridge-address 0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582 \
    --private-key 0x32430699cd4f46ab2422f1df4ad6546811be20c9725544e99253a887e971f92b \
    --destination-network 1 \
    --value 10000000000000000 \
    --rpc-url https://sepolia.drpc.org
```

[This](https://sepolia.etherscan.io/tx/0x1a6e2be69fa65e866889d95403b2fe820f08b6a07b96c6afbde646b8092addb2) is the transaction that was generated and mined from this command.

In most cases, you'll want to specify some `call-data` and a `destination-address` in order for a contract to be called on the destination chain. For example:
```bash
polycli ulxly bridge message \
    --bridge-address 0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582 \
    --private-key 0x32430699cd4f46ab2422f1df4ad6546811be20c9725544e99253a887e971f92b \
    --destination-network 1 \
    --destination-address 0xC92AeF5873d058a76685140F3328B0DED79733Af \
    --call-data 0x40c10f190000000000000000000000003878cff9d621064d393eef92bf1e12a944c5ba84000000000000000000000000000000000000000000000000002386f26fc10000 \
    --value 0 \
    --rpc-url https://sepolia.drpc.org
```
[This](https://sepolia.etherscan.io/tx/0x517b9d827a3a81770d608a6b997e230d992e1e0cabc0fd2797285693b1cc6a9f) is the transaction that was created and mined from running the above command.

In this case, I've configured the destination address to be a test contract I've deployed on L2.
```soldity
// SPDX-License-Identifier: AGPL-3.0
pragma solidity 0.8.20;

contract MessageEmitter {
    event MessageReceived (address originAddress, uint32 originNetwork, bytes data);

    function onMessageReceived(address originAddress, uint32 originNetwork, bytes memory data) external payable {
        emit MessageReceived(originAddress, originNetwork, data);
    }
}
```

The idea is to have minimal contract that will meet the expected interface of the bridge contract: https://github.com/0xPolygonHermez/zkevm-contracts/blob/v9.0.0-rc.3-pp/contracts/interfaces/IBridgeMessageReceiver.sol

In this case, I didn't bother implementing the proxy to an ERC20 or extending some ERC20 contract. I'm just emitting an event to know that the transaction actually fired as expected.
The calldata comes from running this command `cast calldata 'mint(address account, uint256 amount)' 0x3878Cff9d621064d393EEF92bF1e12A944c5ba84  10000000000000000`. Again, in this case no ERC20 will be minted because I didn't set it up.


## Flags

```bash
  -h, --help   help for message
```

The command also inherits flags from parent commands.

```bash
      --bridge-address string              the address of the lxly bridge
      --call-data string                   call data to be passed directly with bridge-message or as an ERC20 Permit (default "0x")
      --call-data-file string              a file containing hex encoded call data
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
      --value string                       the amount in wei to be sent along with the transaction (default "0")
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
