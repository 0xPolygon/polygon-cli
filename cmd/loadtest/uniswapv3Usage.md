The `uniswapv3` command is a subcommand of the `loadtest` tool. It is meant to generate UniswapV3-like load against JSON-RPC endpoints.

You can either chose to deploy the full UniswapV3 contract suite.

```sh
polycli loadtest uniswapv3
```

Or to use pre-deployed contracts to speed up the process.

```bash
polycli loadtest uniswapv3 \
  --uniswap-factory-v3-address 0xc5f46e00822c828e1edcc12cf98b5a7b50c9e81b \
  --uniswap-migrator-address 0x24951726c5d22a3569d5474a1e74734a09046cd9 \
  --uniswap-multicall-address 0x0e695f36ade2a12abea51622e80f105e125d1d6e \
  --uniswap-nft-descriptor-lib-address 0x23050ec03bb24308c788300428a8f9c247f28b25 \
  --uniswap-nft-position-descriptor-address 0xea43847a98b671211b0e412849b69bbd7d53fd00 \
  --uniswap-non-fungible-position-manager-address 0x58eabc23408fb7896b7ce943828cc00044786449 \
  --uniswap-proxy-admin-address 0xdba55eb96288eac85974376b25b3c3f3d67399b7 \
  --uniswap-quoter-v2-address 0x91464a00c4aae9dca6d503a2c24b1dfb8c279e50 \
  --uniswap-staker-address 0xc87383ece9ee3ad3f5158998c4fc04833ba1336e \
  --uniswap-swap-router-address 0x46096eb627d30125f9eaaeefeecaa4e237a04a97 \
  --uniswap-tick-lens-address 0xc73dfb5055874cc7b1cf06ae83f7fe8f6facdb19 \
  --uniswap-upgradeable-proxy-address 0x28656635b0ecd600801600475d61e3ec1534de6e \
  --weth9-address 0x5570d4fd7cce73f0135536d83b8d49e6b77bb76c \
  --uniswap-pool-token-0-address 0x1ce270d0380fbbead12371286aff578a1227d1d7 \
  --uniswap-pool-token-1-address 0x060f7db3146f3d6748822fb4c69489a04b5f3278
```

Contracts are cloned from the different Uniswap repositories, compiled with a specific version of `solc` and go bindings are generated using `abigen`. To learn more about this process, make sure to check out `contracts/uniswapv3/README.org`.
