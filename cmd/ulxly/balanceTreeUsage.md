This command will attempt to compute the root of the balnace tree based on the bridge
events that are provided.

Example usage:

```bash
polycli ulxly compute-balance-tree \
        --l1-deposits-file l1-cardona-4880876-to-6028159.ndjson \
        --l2-deposits-file l2-network-4880876-to-6028159.ndjson \
        --destination-network 19 | jq '.'
```

In this case we are assuming we have two files
`l1-cardona-4880876-to-6028159.ndjson` and `l2-network-4880876-to-6028159.ndjson` that would have been generated
with a call to `polycli ulxly get-deposits` pointing to each network. The output will be the
root of the tree for the provided deposits.

This is the response from polycli:

```json
{
  "root": "0x4516ca2a793b8e20f56ec6ba8ca6033a672330670a3772f76f2ade9bc2125150"",
}
```

Note: more info https://github.com/BrianSeong99/Agglayer_PessimisticProof_Benchmark?tab=readme-ov-file#architecture-of-pessimistic-proof