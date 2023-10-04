https://www.unixtimestamp.com/

```sh
# john's setup
# init pool
cast send \
  --legacy \
  --from 0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6 \
  --private-key 0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa \
  --rpc-url http://127.0.0.1:8545/ \
  --json \
  0x9064b7de206e6e39620326a145b53610e56aeb69 \
  'function initialize(uint160 sqrtPriceX96) external override' \
  79228162514264337593543950336

# check slot0
cast call \
  --rpc-url http://127.0.0.1:8545/ \
  0x9064b7de206e6e39620326a145b53610e56aeb69 \
  'function slot0() external view returns (uint160 sqrtPriceX96, int24 tick, uint16 observationIndex, uint16 observationCardinality, uint16 observationCardinalityNext, uint8 feeProtocol, bool unlocked)'
79228162514264337593543950336
0
0
1
1
0
true

# provide liquidity
cast send \
  --legacy \
  --from 0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6 \
  --private-key 0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa \
  --rpc-url http://127.0.0.1:8545/ \
  --json \
  0xF5A73e7cFCC83b7e8ce2e17Eb44f050E8071eE60 \
  'mint((address,address,uint24,int24,int24,uint256,uint256,uint256,uint256,address,uint256) MintParams)' \
  '(0x1c537fab97840a2fef787b75d37e6f621c870eb9,0x1f7dfc0cee2b55573bb2a3d4452693d203994274,3000,-887220,887220,5000,5000,0,0,0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6,1759474606)'

# check balances after providing a few tokens
cast call \
  --rpc-url http://127.0.0.1:8545/ \
  0x0cccc4e4fc22306cacd027a49ce78ed1972849c6 \
  'balanceOf(address account) public view virtual returns (uint256)' \
  0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6
999999985933465751988706233

cast call \
  --rpc-url http://127.0.0.1:8545 \
  0x0cccc4e4fc22306cacd027a49ce78ed1972849c6 \
  'allowance(address owner, address spender) public view virtual returns (uint256)' \
  0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6 0xF5A73e7cFCC83b7e8ce2e17Eb44f050E8071eE60
0
```

```sh
11:30AM DBG UniswapV3 deployment config config={"FactoryV3":"0xba94fa968d4f589fe649f02a923c4b4a04bb7449","Migrator":"0x311981b16534238422d301e95a42f6c27a24f346","Multicall":"0x383f70613d09544d5bbf7c66917de707a4856b3f","NFTDescriptor":"0x5a930ce789ec78f23e3bbfd8dbadda72d918523c","NFTPositionDescriptor":"0xa05a2e963b527cff651916b8a6d6d67e7b3e5b84","NonfungiblePositionManager":"0x86d9b21e49c81baefcdd05b728790e6e555eaed3","ProxyAdmin":"0xf5a73e7cfcc83b7e8ce2e17eb44f050e8071ee60","QuoterV2":"0x73a6d49037afd585a0211a7bb4e990116025b45d","Staker":"0xd5be512690ab8485c21d7030e06735dcfdd66268","SwapRouter02":"0x34cb2c25dd47f344079443cec353290441ac8ac2","TickLens":"0x84f3e2983edd66138aa8fd6dc1e482b971492992","TransparentUpgradeableProxy":"0x1c537fab97840a2fef787b75d37e6f621c870eb9","WETH9":"0xd179f002f585965435a1b62c92c2afbd3930320a"}
11:30AM DBG Token0 contract deployed address=0x7d66094751107bc366206362dc9c27b638cd8a36
11:30AM TRC Starting blocking loop blockInterval=1 numberOfBlocksToWaitFor=30 startBlockNumber=106
11:30AM TRC New block newBlock=106
11:30AM TRC Token0 contract is not available yet
11:30AM TRC Unable to execute function error="no contract code at given address" elapsedTimeSeconds=0
11:30AM TRC New block newBlock=106
11:30AM TRC New block newBlock=107
11:30AM TRC Token0 contract is not available yet
11:30AM DBG Spending approved Swapper=Token0 amount=1000000000000000000 spender=0x86d9B21e49c81BAEFcdD05B728790E6E555eaed3
11:30AM DBG Spending approved Swapper=Token0 amount=1000000000000000000 spender=0x34cB2c25Dd47F344079443Cec353290441ac8aC2
11:30AM TRC Function executed successfully elapsedTimeSeconds=2
11:30AM DBG Token1 contract deployed address=0xc0019107cb4f79d41fcb00175b9721c32f07879f
11:30AM TRC Starting blocking loop blockInterval=1 numberOfBlocksToWaitFor=30 startBlockNumber=107
11:30AM TRC New block newBlock=107
11:30AM TRC Token1 contract is not available yet
11:30AM TRC Unable to execute function error="no contract code at given address" elapsedTimeSeconds=0
11:30AM TRC New block newBlock=107
11:30AM TRC New block newBlock=108
11:30AM TRC Token1 contract is not available yet
11:30AM DBG Spending approved Swapper=Token1 amount=1000000000000000000 spender=0x86d9B21e49c81BAEFcdD05B728790E6E555eaed3
11:30AM DBG Spending approved Swapper=Token1 amount=1000000000000000000 spender=0x34cB2c25Dd47F344079443Cec353290441ac8aC2
11:30AM TRC Function executed successfully elapsedTimeSeconds=2
11:30AM DBG Pool created and initialized
11:30AM TRC Starting blocking loop blockInterval=1 numberOfBlocksToWaitFor=30 startBlockNumber=108
11:30AM TRC New block newBlock=108
11:30AM TRC Unable to execute function error="Token0-Token1 pool not deployed yet" elapsedTimeSeconds=0
11:30AM TRC New block newBlock=108
11:30AM TRC New block newBlock=109
11:30AM TRC Function executed successfully elapsedTimeSeconds=2
11:30AM DBG Token0-Token1 pool instantiated address=0x1d1c5c05a0d0735e32aaeaf8e52a6913d38bc039
11:30AM DBG Token0-Token1 pool state liquidity=0 slot0={"FeeProtocol":0,"ObservationCardinality":1,"ObservationCardinalityNext":1,"ObservationIndex":0,"SqrtPriceX96":79228162514264337593543950336,"Tick":0,"Unlocked":true}
11:30AM DBG Waiting for 5 blocks to be mined...
11:30AM DBG Waiting over
11:30AM DBG DEBUG params={"Amount0Desired":5000,"Amount0Min":0,"Amount1Desired":5000,"Amount1Min":0,"Deadline":1696412417,"Fee":3000,"Recipient":"0x85da99c8a7c2c95964c8efd687e95e632fc533d6","TickLower":-887220,"TickUpper":887220,"Token0":"0x7d66094751107bc366206362dc9c27b638cd8a36","Token1":"0xc0019107cb4f79d41fcb00175b9721c32f07879f"}
```

```sh
# my setup
# init pool
cast send \
  --legacy \
  --from 0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6 \
  --private-key 0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa \
  --rpc-url http://127.0.0.1:8545/ \
  --json \
  0x1d1c5c05a0d0735e32aaeaf8e52a6913d38bc039 \
  'function initialize(uint160 sqrtPriceX96) external override' \
  79228162514264337593543950336

# check slot0
cast call \
  --rpc-url http://127.0.0.1:8545/ \
  0x1d1c5c05a0d0735e32aaeaf8e52a6913d38bc039 \
  'function slot0() external view returns (uint160 sqrtPriceX96, int24 tick, uint16 observationIndex, uint16 observationCardinality, uint16 observationCardinalityNext, uint8 feeProtocol, bool unlocked)'
79228162514264337593543950336
0
0
1
1
0
true

# check balance of token0
cast call \
  --rpc-url http://127.0.0.1:8545/ \
  0x7d66094751107bc366206362dc9c27b638cd8a36 \
  'balanceOf(address account) public view virtual returns (uint256)' \
  0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6
1000000000000000000000000000000000000

cast call \
  --rpc-url http://127.0.0.1:8545 \
  0x7d66094751107bc366206362dc9c27b638cd8a36 \
  'allowance(address owner, address spender) public view virtual returns (uint256)' \
  0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6 0x86d9b21e49c81baefcdd05b728790e6e555eaed3
1000000000000000000

# check balance of token1
cast call \
  --rpc-url http://127.0.0.1:8545/ \
  0xc0019107cb4f79d41fcb00175b9721c32f07879f \
  'balanceOf(address account) public view virtual returns (uint256)' \
  0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6
1000000000000000000000000000000000000

cast call \
  --rpc-url http://127.0.0.1:8545 \
  0xc0019107cb4f79d41fcb00175b9721c32f07879f \
  'allowance(address owner, address spender) public view virtual returns (uint256)' \
  0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6 0x86d9b21e49c81baefcdd05b728790e6e555eaed3
1000000000000000000

# provide liquidity
cast send \
  --legacy \
  --from 0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6 \
  --private-key 0x42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa \
  --rpc-url http://127.0.0.1:8545/ \
  --json \
  0x86d9b21e49c81baefcdd05b728790e6e555eaed3 \
  'mint((address,address,uint24,int24,int24,uint256,uint256,uint256,uint256,address,uint256) MintParams)' \
  '(0x7d66094751107bc366206362dc9c27b638cd8a36,0xc0019107cb4f79d41fcb00175b9721c32f07879f,3000,-887220,887220,5000,5000,0,0,0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6,1759474606)'
```
