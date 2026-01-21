This command will attempt to scan a range of blocks and look for uLxLy
Claim Events. This is the specific signature that we're interested
in:

```solidity
    /**
     * @dev Emitted when a claim is done from another network
     */
    event ClaimEvent(
        uint256 globalIndex,
        uint32 originNetwork,
        address originAddress,
        address destinationAddress,
        uint256 amount
    );
```

If you're looking at the raw topics from on chain or in an explorer, this is the associated value:

`0x1df3f2a973a00d6635911755c260704e95e8a5876997546798770f76396fda4d`

Each event that we counter will be parsed and written as JSON to
stdout. Example usage:

```bash
polycli ulxly get-claims \
        --bridge-address 0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582 \
        --rpc-url https://eth-sepolia.g.alchemy.com/v2/demo \
        --from-block 4880876 \
        --to-block 6028159 \
        --filter-size 9999 > cardona-4880876-to-6028159.ndjson
```

This command will look for claim events from block `4880876` to
block `6028159` in increments of `9999` blocks at a time for the
contract address `0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582`. The
output will be written as newline delimited JSON.

This command is very specific for the ulxly bridge, and it's meant to
serve as the input to the proof command.

