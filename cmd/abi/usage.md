# ABI Encode

When calling contract functions, the interaction primarily goes through the Contract Application Binary Interface (ABI). This interface follows a specific encoding. To quickly create that encoding, we can utilize `abi encode`:

```bash
$ polycli abi encode "someFunctionSignature(string,uint256)" "string input" 10
```

The command follows the format of:
```
polycli abi encode [function-signature] [args...]
```
Where the signature is a fragment in the form `<function name>(<types...>)`.

The encoding adheres to [Solidity Contract ABI Specifications](https://docs.soliditylang.org/en/latest/abi-spec.html).


# ABI Decode

When looking at raw contract calls, sometimes we have an ABI and we just want to quickly figure out which method is being called. This is a quick way to get all of the function selectors for an ABI.

We can the command like this to get the function signatures and selectors.

```bash
$ polycli abi decode --file contract.abi
```

This would output some information that would let us know the various function selectors for this contract.

```bash
$ polycli abi decode --file ./bindings/tester/LoadTester.abi
Selector:3430ec06       Signature:dumpster(uint256)(bytes)
Selector:a271b721       Signature:loopBlockHashUntilLimit()(uint256)
Selector:0ba8a73b       Signature:testADD(uint256)(uint256)
Selector:fde7721c       Signature:testMULMOD(uint256)(uint256)
Selector:63138d4f       Signature:testSHA256(bytes)(bytes32)
...
```

If we want to break down input data we can run something like this.

```bash
$ polycli abi decode --data 0xd53ff3fd0000000000000000000000000000000000000000000000000000000000000063 < ./bindings/tester/LoadTester.abi
Selector:a60a1087       Signature:testCHAINID(uint256)(uint256)
Selector:b7b86207       Signature:testCODESIZE(uint256)(uint256)
...
Selector:d53ff3fd       Signature:testSUB(uint256)(uint256)
Selector:962e4dc2       Signature:testECMul(bytes)(bytes)
Selector:de97a363       Signature:testEXP(uint256)(uint256)
Selector:9cce7cf9       Signature:testIdentity(bytes)(bytes)
Selector:19b621d6       Signature:testSHA3(uint256)(uint256)
...
id: d53ff3fd, 0000000000000000000000000000000000000000000000000000000000000063
Input data:
{
  "x": 99
}
Signature and Input
testSUB(uint256)(uint256) 99
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
