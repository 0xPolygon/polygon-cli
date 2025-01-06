# `polycli ulxly proof`

> Auto-generated documentation.

## Table of Contents

- [Description](#description)
- [Usage](#usage)
- [Flags](#flags)
- [See Also](#see-also)

## Description

generate a proof for a given range of deposits

```bash
polycli ulxly proof [flags]
```

## Usage

This command will attempt to create a merkle proof based on the bridge
events that are provided.

Example usage:

```bash
polycli ulxly proof \
        --file-name cardona-4880876-to-6028159.ndjson \
        --deposit-number 24386 | jq '.'
```

In this case we are assuming we have a file
`cardona-4880876-to-6028159.ndjson` that would have been generated
with a call to `polycli ulxly deposits`. The output will be the
sibling hashes necessary to prove inclusion of deposit `24386`.

This is a real verifiable deposit if you'd like to sanity check:

- Deposit Transaction: https://sepolia.etherscan.io/tx/0x1f950d076ad534fe588bd6a8f58904395c907df4738f92bd8aea513c19d1fa5f
- Mainnet Root: `4516CA2A793B8E20F56EC6BA8CA6033A672330670A3772F76F2ADE9BC2125150`
- Actual Claim Tx: https://cardona-zkevm.polygonscan.com/tx/0x5d4fbaca896f015801f1049b383932eaa9363d344c36b1c51e5f2e3ce20f9dc3

This is the proof response from polycli:

```json
{
  "Siblings": [
    "0x0000000000000000000000000000000000000000000000000000000000000000",
    "0x18dd49d4ac3b31a6468446597686e7164bfb88a09685d3cd31f8f4b0b91e7d86",
    "0xb4c11951957c6f8f642c4af61cd6b24640fec6dc7fc607ee8206a99e92410d30",
    "0x21ddb9a356815c3fac1026b6dec5df3124afbadb485c9ba5a3e3398a04b7ba85",
    "0xe58769b32a1beaf1ea27375a44095a0d1fb664ce2dd358e7fcbfb78c26a19344",
    "0x0eb01ebfc9ed27500cd4dfc979272d1f0913cc9f66540d7e8005811109e1cf2d",
    "0x3822fc0c89d0438d84cb8b41b6a9eecd5d9d369ce5022b6d0817d96e695fca08",
    "0xffd70157e48063fc33c97a050f7f640233bf646cc98d9524c6b92bcf3ab56f83",
    "0xd1f352a4bd8bdd17b172d7a39c2406672c4644a624fd3c5ac30111cb716c6a79",
    "0x8549075281f6770dc405ba232aa39c16992a441b2e88b01742a35bb949cb4e54",
    "0xec88c2530a54d444e5212e08a65c5b7ee7ce34ef0acd31e1eb2c1ab5aa91772b",
    "0x9b2bc070d68635596a9e3cefef6c5c7530751b415fe63cfa3c01de5d8253379c",
    "0x25e2d43e7f5545a5808a41875dc6485c4bd9e2df61396b682685d00a00ca1e88",
    "0xc1df82d9c4b87413eae2ef048f94b4d3554cea73d92b0f7af96e0271c691e2bb",
    "0x893f7fe87020f40c2edd5239c1b4be2f16862adb151deea1e379b33f35ce8bd3",
    "0xda7bce9f4e8618b6bd2f4132ce798cdc7a60e7e1460a7299e3c6342a579626d2",
    "0x2733e50f526ec2fa19a22b31e8ed50f23cd1fdf94c9154ed3a7609a2f1ff981f",
    "0xe1d3b5c807b281e4683cc6d6315cf95b9ade8641defcb32372f1c126e398ef7a",
    "0x5a2dce0a8a7f68bb74560f8f71837c2c2ebbcbf7fffb42ae1896f13f7c7479a0",
    "0xb46a28b6f55540f89444f63de0378e3d121be09e06cc9ded1c20e65876d36aa0",
    "0xc65e9645644786b620e2dd2ad648ddfcbf4a7e5b1a3a4ecfe7f64667a3f0b7e2",
    "0xf4418588ed35a2458cffeb39b93d26f18d2ab13bdce6aee58e7b99359ec2dfd9",
    "0x5a9c16dc00d6ef18b7933a6f8dc65ccb55667138776f7dea101070dc8796e377",
    "0x4df84f40ae0c8229d0d6069e5c8f39a7c299677a09d367fc7b05e3bc380ee652",
    "0xcdc72595f74c7b1043d0e1ffbab734648c838dfb0527d971b602bc216c9619ef",
    "0x0abf5ac974a1ed57f4050aa510dd9c74f508277b39d7973bb2dfccc5eeb0618d",
    "0xb8cd74046ff337f0a7bf2c8e03e10f642c1886798d71806ab1e888d9e5ee87d0",
    "0x838c5655cb21c6cb83313b5a631175dff4963772cce9108188b34ac87c81c41e",
    "0x662ee4dd2dd7b2bc707961b1e646c4047669dcb6584f0d8d770daf5d7e7deb2e",
    "0x388ab20e2573d171a88108e79d820e98f26c0b84aa8b2f4aa4968dbb818ea322",
    "0x93237c50ba75ee485f4c22adf2f741400bdf8d6a9cc7df7ecae576221665d735",
    "0x8448818bb4ae4562849e949e17ac16e0be16688e156b5cf15e098c627c0056a9"
  ],
  "Root": "0x4516ca2a793b8e20f56ec6ba8ca6033a672330670a3772f76f2ade9bc2125150",
  "DepositCount": 24386,
  "LeafHash": "0x2c42c143213fd0e36d843d9d40866ce7be02c671beec0eae3ffd3d2638acc87c"
}
```

![Sample Tree](./tree-diagram.png)

When we're creating the proof here, we're essentially storing all of
the paths to the various leafs. When we want to generate a proof, we
essentially find the appropriate sibling node in the tree to prove
that the leaf is part of the given merkle root.

## Flags

```bash
      --deposit-count uint   The deposit number to generate a proof for
      --file-name string     An ndjson file with deposit data
  -h, --help                 help for proof
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

- [polycli ulxly](polycli_ulxly.md) - Utilities for interacting with the lxly bridge
