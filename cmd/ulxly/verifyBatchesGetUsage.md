This command will attempt to scan a range of blocks and look for verify batch Events. This is the specific signature that we're interested
in:

```solidity
    /**
     * @dev Emitted when the trusted aggregator verifies batches
     */
    event VerifyBatchesTrustedAggregator(
        uint32 indexed rollupID,
        uint64 numBatch,
        bytes32 stateRoot,
        bytes32 exitRoot,
        address indexed aggregator
    );

```

If you're looking at the raw topics from on chain or in an explorer, this is the associated value:

`0xd1ec3a1216f08b6eff72e169ceb548b782db18a6614852618d86bb19f3f9b0d3`

Each event that we counter will be parsed and written as JSON to
stdout. Example usage:

```bash
polycli ulxly get-verify-batches \
        --rollup-manager-address 0x32d33d5137a7cffb54c5bf8371172bcec5f310ff \
        --rpc-url https://eth-sepolia.g.alchemy.com/v2/demo \
        --from-block 4880876 \
        --to-block 6028159 \
        --filter-size 9999 > verify-batches-cardona-4880876-to-6028159.ndjson
```

This command will look for verify batch events from block `4880876` to
block `6028159` in increments of `9999` blocks at a time for the
contract address `0x32d33d5137a7cffb54c5bf8371172bcec5f310ff`. The
output will be written as newline delimited JSON.

This command is very specific for the ulxly bridge, and it's meant to
serve as the input to the rollup-proof command.


