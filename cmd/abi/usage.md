When looking at raw contract calls, sometimes we have an ABI and we just want to quickly figure out which method is being called. This is a quick way to get all of the function selectors for an ABI.

This is a simple test.

We can the command like this to get the function signatures and selectors.

```bash
$ polycli abi --file contract.abi
```

This would output some information that would let us know the various function selectors for this contract.

```bash
Selector:19d8ac61	Signature:function lastTimestamp() view returns(uint64)
Selector:a066215c	Signature:function setVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget) returns()
Selector:715018a6	Signature:function renounceOwnership() returns()
Selector:cfa8ed47	Signature:function trustedSequencer() view returns(address)
```

If we want to break down input data we can run something like this.

```bash
$ polycli abi --data 0x3c158267000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000063ed0f8f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006eec03843b9aca0082520894d2f852ec7b4e457f6e7ff8b243c49ff5692926ea87038d7ea4c68000808204c58080642dfe2cca094f2419aad1322ec68e3b37974bd9c918e0686b9bbf02b8bd1145622a3dd64202da71549c010494fd1475d3bf232aa9028204a872fd2e531abfd31c000000000000000000000000000000000000 < contract.abi
```

In addition to the function selector data, we'll also get a breakdown of input data:

```json
{
  "batches": [
    {
      "transactions": "7AOEO5rKAIJSCJTS+FLse05Ff25/+LJDxJ/1aSkm6ocDjX6kxoAAgIIExYCAZC3+LMoJTyQZqtEyLsaOOzeXS9nJGOBoa5u/Ari9EUViKj3WQgLacVScAQSU/RR1078jKqkCggSocv0uUxq/0xw=",
      "globalExitRoot": [
        0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
        0, 0, 0, 0, 0, 0, 0, 0
      ],
      "timestamp": 1676480399,
      "minForcedTimestamp": 0
    }
  ]
}
```
