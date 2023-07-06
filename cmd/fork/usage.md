Occasionally, we'll want to analyze the details of a side chain to understand in detail who was proposing the blocks, what was the difficultly, and just generally get better understanding of block propagation.

```bash
$ polycli fork 0x053d84d5215684c8ae810a4729f7c9b54d65a80b128a27aeddcd7dc295a0cebd https://polygon-rpc.com
```

In order to use this, you'll need to have a blockhash of a block that was part of a fork / side chain. Once you have that, you can run `fork` against a node to get the details of the fork and the canonical chain.
