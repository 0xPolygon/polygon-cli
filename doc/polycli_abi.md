# `polycli abi`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

Provides encoding and decoding functionalities with contract signatures and ABI.

## Usage

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

## Encode Usage Guide:

### Basic type inputs:

- `string`: `"f(string)" "a string input"`
> Note: All strings must be enclosed with double quotation `""`.

- `uint<M>`: `"f(uint256)" 999`
- `int<M>`: `"f(int256,int8)" -- 222 -8`

> Note: the `--` is intended here to indicate end of command line flags. If you
> don't have it, `Cobra` will assume the negative number is a flag.

- `bool`: `"f(bool)" true`
- `address`: `"f(address)" 0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6`
- `bytes<M>`: `"f(bytes3)" 0x123456`

> Note: The `0x` prefix is optional

More Examples:

```bash
$ polycli abi encode "someFunctionSignature(string,uint256,bytes3,int8,bool,address,bytes)" "string input" 10 0x123456 222 true 0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6 99999999
```

### Array Types

- `<type>[]`: `"f(int256[])" '[1,2,3]'`

> Note: the argument of structs and arrays, must be enclosed within quotes so
> that Cobra can parse it as a whole argument.

### Struct Types

For structs, each element is represented as an element of a tuple. For example:
```
struct School {
    string   name;
    string[] addresses;
    uint256  numOfStudents;
    bool     isPrivate;
}
```
Would be represented as: 
```
(string,string[],uint256,bool)
```

An example of a call with that type of input would be:
```
$ polycli abi encode "newSchool((string,string[],uint256,bool))" '("matic",["123 street ave.","321 ave st."], 9999, false)'
0x5866fb060000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000270f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000056d61746963000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000f31323320737472656574206176652e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b333231206176652073742e000000000000000000000000000000000000000000
```

For arrays inputs are 

For more nested structures like:
```
struct School {
    string     name;
    string[]   addresses;
    uint256    numOfStudents;
    bool       isPrivate;
    Students[] students
}
struct Students {
    string     name;
    uint256[]  grades;
    bool       isActive;
}
```
Just nest the tuples as needed:
```
(string,string[],uint256,bool,(string,uint256[],bool)[])
```

An example run:
```
$ polycli abi encode "newSchool((string,string[],uint256,bool,(string,uint256[],bool)[]))" '("matic",["123 street ave.","321 ave st."], 9999, false,[("Alice",[90,89],true),("Bob",[99,99,70],true)])'
0x197289f8000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000270f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001c000000000000000000000000000000000000000000000000000000000000000056d61746963000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000f31323320737472656574206176652e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b333231206176652073742e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000005416c6963650000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000005a0000000000000000000000000000000000000000000000000000000000000059000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000003426f6200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006300000000000000000000000000000000000000000000000000000000000000630000000000000000000000000000000000000000000000000000000000000046
```

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

## Flags

```bash
  -h, --help   help for abi
```

The command also inherits flags from parent commands.

```bash
      --config string   config file (default is $HOME/.polygon-cli.yaml)
      --pretty-logs     Should logs be in pretty format or JSON (default true)
  -v, --verbosity int   0 - Silent
                        100 Panic
                        200 Fatal
                        300 Error
                        400 Warning
                        500 Info
                        600 Debug
                        700 Trace (default 500)
```

## See also

- [polycli](polycli.md) - A Swiss Army knife of blockchain tools.
- [polycli abi decode](polycli_abi_decode.md) - Parse an ABI and print the encoded signatures.

- [polycli abi encode](polycli_abi_encode.md) - ABI encodes a function signature and the inputs

