This command prints a proof response where all siblings are set to
the zero value (`0x0000000000000000000000000000000000000000000000000000000000000000`).

This is useful when you need to submit a dummy proof, for example in
testing scenarios or when a proof placeholder is required but the
actual proof data is not needed.

Note: This is different from `zero-proof`, which generates the actual
zero hashes of a Merkle tree (intermediate hashes computed from empty
subtrees). Empty proof simply fills all siblings with the zero value.

Example usage:

```bash
polycli ulxly empty-proof
```

Example output:

```json
{
  "siblings": [
    "0x0000000000000000000000000000000000000000000000000000000000000000",
    "0x0000000000000000000000000000000000000000000000000000000000000000",
    "0x0000000000000000000000000000000000000000000000000000000000000000",
    ...
  ]
}
```
