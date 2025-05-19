This command will attempt to computethe nullifierTree based on the claims that are provided.

Example usage:

```bash
polycli ulxly compute-nullifier-tree \
        --file-name claims-cardona-4880876-to-6028159.ndjson | jq '.'
```

In this case we are assuming we have a file
`claims-cardona-4880876-to-6028159.ndjson` that would have been generated
with a call to `polycli ulxly get-claims`. The output will be the
claims necessary to compute the nullifier tree.

This is the response from polycli:

```json
{
  "root": "0x4516ca2a793b8e20f56ec6ba8ca6033a672330670a3772f76f2ade9bc2125150"",
}
```

Note: more info https://github.com/BrianSeong99/Agglayer_PessimisticProof_Benchmark?tab=readme-ov-file#architecture-of-pessimistic-proof