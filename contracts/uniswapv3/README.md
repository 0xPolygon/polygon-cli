# Requirements

1. Make sure you have `solc@0.7.6` installed. This is required to build UniswapV3 contracts. A handy way to manage `solc` versions is to use [crytic/solc-select](https://github.com/crytic/solc-select).

2. Build the core, periphery and swap-router contracts using the `build.sh`script.

```sh
$ sh build.sh
solc, the solidity compiler commandline interface
Version: 0.7.6+commit.7338295f.Darwin.appleclang

🏗️  Building v3-core contracts...
Cloning into 'v3-core'...
remote: Enumerating objects: 8244, done.
remote: Counting objects: 100% (4/4), done.
remote: Compressing objects: 100% (4/4), done.
remote: Total 8244 (delta 0), reused 2 (delta 0), pack-reused 8240
Receiving objects: 100% (8244/8244), 6.37 MiB | 13.20 MiB/s, done.
Resolving deltas: 100% (6276/6276), done.
Compiler run successful. Artifact(s) can be found in directory tmp/v3-core.
✅ Successfully built v3-core contracts...
...
```

3. Generate Go bindings for Uniswap contracts using the `bindings.sh` script.

```sh
$ sh bindings.sh
abigen version 1.13.1-stable
✅ UniswapV3Factory bindings generated.
...
```
