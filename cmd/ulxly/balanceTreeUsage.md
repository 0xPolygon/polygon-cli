This command will attempt to compute the root of the balnace tree based on the bridge
events that are provided.

Example usage:

```bash
polycli ulxly compute-balance-tree \
        --l2-claims-file l2-claims-0-to-11454081.ndjson \
        --l2-deposits-file l2-deposits-0-to-11454081.ndjson \
        --l2-network-id 3 | jq '.'
```

In this case we are assuming we have two files
`l2-claims-0-to-11454081.ndjson` and `l2-deposits-0-to-11454081.ndjson` that would have been generated
with a call to `polycli ulxly get-deposits` and `polycli ulxly get-claims` pointing to each network. The output will be the
root of the tree for the provided deposits and claims.

This is the response from polycli:

```json
{
  "root": "0x4516ca2a793b8e20f56ec6ba8ca6033a672330670a3772f76f2ade9bc2125150"",
}
```

Note: more info https://github.com/BrianSeong99/Agglayer_PessimisticProof_Benchmark?tab=readme-ov-file#architecture-of-pessimistic-proof