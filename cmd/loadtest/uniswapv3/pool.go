package uniswapv3loadtest

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maticnetwork/polygon-cli/bindings/uniswapv3"
	"github.com/maticnetwork/polygon-cli/util"
	"github.com/rs/zerolog/log"
)

var (
	// Reserve of a token in a UniswapV3 pool.
	poolReserveForOneToken = big.NewInt(1_000_000_000_000)

	// The timeout of the mint operation (liquidity providing).
	mintOperationTimeout = 1 * time.Hour
)

// The maximum tick that may be passed to `getSqrtRatioAtTick` computed from log base 1.0001 of 2**128.
const maxTick = 887272

type feeTier float64

var (
	// Only a few fee tiers are possible in UniswapV3. They are represented in percentage.
	// https://uniswapv3book.com/docs/milestone_5/swap-fees/#accruing-swap-fees
	StableTier   feeTier = 0.05 // 500
	StandardTier feeTier = 0.3  // 3_000
	ExoticTier   feeTier = 1    // 10_000
)

// PercentageToUniswapFeeTier takes a percentage and returns the corresponding UniswapV3 fee tier.
func PercentageToUniswapFeeTier(p float64) *big.Int {
	var fees int64
	switch p {
	case float64(StableTier):
		fees = 500
	case float64(StandardTier):
		fees = 3_000
	case float64(ExoticTier):
		fees = 10_000
	}
	return big.NewInt(fees)
}

// PoolConfig represents the configuration of a UniswapV3 pool.
type PoolConfig struct {
	Token0, Token1     ContractConfig[uniswapv3.Swapper]
	ReserveA, ReserveB *big.Int
	Fees               *big.Int
}

// Create a new `PoolConfig` object.
func NewPool(token0, token1 ContractConfig[uniswapv3.Swapper], fees *big.Int) *PoolConfig {
	p := PoolConfig{
		ReserveA: poolReserveForOneToken,
		ReserveB: poolReserveForOneToken,
		Fees:     fees,
	}

	// Make sure the token pair is sorted.
	if token0.Address.Hex() < token1.Address.Hex() {
		p.Token0 = token0
		p.Token1 = token1
	} else {
		p.Token0 = token1
		p.Token1 = token0
	}

	return &p
}

// slot represents the state of a UniswapV3 pool.
type slot struct {
	SqrtPriceX96               *big.Int
	Tick                       *big.Int
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint8
	Unlocked                   bool
}

// SetupLiquidityPool sets up a UniswapV3 liquidity pool, creating and initializing it if needed,
// and providing liquidity in case none exists.
func SetupLiquidityPool(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts, uniswapV3Config UniswapV3Config, poolConfig PoolConfig, recipient common.Address) error {
	// Create and initialise pool.
	poolContract, err := createPool(ctx, c, tops, cops, uniswapV3Config, poolConfig)
	if err != nil {
		return err
	}

	// Get pool state.
	var slot0 slot
	var liquidity *big.Int
	slot0, liquidity, err = getPoolState(cops, poolContract)
	if err != nil {
		return err
	}
	log.Trace().Interface("slot0", slot0).Interface("liquidity", liquidity).Msg("Pool state")

	// Provide liquidity if there's none.
	if liquidity.Cmp(big.NewInt(0)) == 0 {
		if provideLiquidity(ctx, c, tops, cops, poolContract, poolConfig, recipient, uniswapV3Config.NonfungiblePositionManager.Contract) != nil {
			return err
		}
	} else {
		log.Debug().Msg("Liquidity already provided to the pool")
	}

	return nil
}

// createPool creates and initialises the UniswapV3 liquidity pool if needed.
func createPool(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts, uniswapV3Config UniswapV3Config, poolConfig PoolConfig) (*uniswapv3.IUniswapV3Pool, error) {
	// Create and initialize the pool.
	sqrtPriceX96 := computeSqrtPriceX96(poolConfig.ReserveA, poolConfig.ReserveB)
	if _, err := uniswapV3Config.NonfungiblePositionManager.Contract.CreateAndInitializePoolIfNecessary(tops, poolConfig.Token0.Address, poolConfig.Token1.Address, poolConfig.Fees, sqrtPriceX96); err != nil {
		log.Error().Err(err).Msg("Unable to create and initialize the pool")
		return nil, err
	}
	log.Debug().Interface("fees", poolConfig.Fees).Msg("Pool created and initialized")

	// Retrieve the pool address.
	var poolAddress common.Address
	err := util.BlockUntilSuccessful(ctx, c, func() (err error) {
		poolAddress, err = uniswapV3Config.FactoryV3.Contract.GetPool(cops, poolConfig.Token0.Address, poolConfig.Token1.Address, poolConfig.Fees)
		if poolAddress == (common.Address{}) {
			return errors.New("pool not deployed yet")
		}
		return
	})
	if err != nil {
		log.Error().Err(err).Msg("Unable to retrieve the address of the pool")
		return nil, err
	}

	// Instantiate the pool contract.
	contract, err := uniswapv3.NewIUniswapV3Pool(poolAddress, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to instantiate the pool")
		return nil, err
	}
	log.Debug().Interface("address", poolAddress).Msg("Pool instantiated")
	return contract, nil
}

// computeSqrtPriceX96 calcules the square root of the price ratio of two reserves in a UniswapV3 pool.
// https://uniswapv3book.com/docs/milestone_1/calculating-liquidity/#price-range-calculation
func computeSqrtPriceX96(reserveA, reserveB *big.Int) *big.Int {
	sqrtReserveA := new(big.Int).Sqrt(reserveA)
	sqrtReserveB := new(big.Int).Sqrt(reserveB)
	q96 := new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
	sqrtPriceX96 := new(big.Int).Mul(sqrtReserveB, q96)
	sqrtPriceX96.Div(sqrtPriceX96, sqrtReserveA)
	return sqrtPriceX96
}

// getPoolState returns UniswapV3 pool's slot0 and liquidity.
func getPoolState(cops *bind.CallOpts, contract *uniswapv3.IUniswapV3Pool) (slot, *big.Int, error) {
	// Get pool state.
	var slot0 slot
	var err error
	slot0, err = contract.Slot0(cops)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get pool's slot0")
		return slot{}, nil, err
	}

	// Get pool's liquidity.
	var liquidity *big.Int
	liquidity, err = contract.Liquidity(cops)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get pool's liquidity")
		return slot{}, nil, err
	}

	return slot0, liquidity, nil
}

// provideLiquidity provides liquidity to the UniswapV3 pool.
func provideLiquidity(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts, poolContract *uniswapv3.IUniswapV3Pool, poolConfig PoolConfig, recipient common.Address, nftPositionManagerContract *uniswapv3.NonfungiblePositionManager) error {
	// Compute the tick lower and upper for providing liquidity.
	// The default tick spacing is set to 60 for the 0.3% fee tier and unfortunately, `MIN_TICK` and
	// `MAX_TICK` are not divisible by this amount. The solution is to use a multiple of 60 instead.
	tickSpacing, err := poolContract.TickSpacing(cops)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get tick spacing")
		return err
	}
	// tickUpper = (MAX_TICK / tickSpacing) * tickSpacing
	tickUpper := new(big.Int).Div(big.NewInt(maxTick), tickSpacing)
	tickUpper.Mul(tickUpper, tickSpacing)
	// tickLower = - tickUpper
	tickLower := new(big.Int).Neg(tickUpper)

	// Compute deadline.
	var latestBlockTimestamp *big.Int
	latestBlockTimestamp, err = getLatestBlockTimestamp(ctx, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get timestamp of latest block")
		return err
	}
	deadline := new(big.Int).Add(latestBlockTimestamp, big.NewInt(int64(mintOperationTimeout.Seconds()))) // only for one minute.

	// Provide liquidity.
	poolSize := new(big.Int).Add(poolConfig.ReserveA, poolConfig.ReserveB)
	mintParams := uniswapv3.INonfungiblePositionManagerMintParams{
		// The address of the token0, first token of the pool.
		Token0: poolConfig.Token0.Address,
		// The address of the token1, second token of the pool.
		Token1: poolConfig.Token1.Address,
		// The fee associated with the pool.
		Fee: poolConfig.Fees,

		// The lower end of the tick range for the position.
		// Here, we provide liquidity across the whole possible range (divisible by tick spacing).
		TickLower: tickLower,
		// The higher end of the tick range for the position.
		TickUpper: tickUpper,

		// The desired amount of token0 to be sent to the pool during the minting operation.
		Amount0Desired: poolSize,
		// The desired amount of token1 to be sent to the pool during the minting operation.
		Amount1Desired: poolSize,
		// The minimum acceptable amount of token0 to be sent to the pool. This represents the slippage
		// protection for token0 during the minting. Here we don't want to lose any token.
		Amount0Min: poolSize,
		// The minimum acceptable amount of token1 to be sent to the pool. This represents the slippage
		// protection for token1 during the minting. Here we don't want to lose any token.
		Amount1Min: poolSize,

		// The destination address of the pool fees.
		Recipient: recipient,

		// The unix time after which the mint will fail, to protect against long-pending transactions
		// and wild swings in prices.
		Deadline: deadline,
	}
	var liquidity *big.Int
	err = util.BlockUntilSuccessful(ctx, c, func() (err error) {
		// Mint tokens.
		_, err = nftPositionManagerContract.Mint(tops, mintParams)
		if err != nil {
			return err
		}

		// Check that liquidity has been added to the pool.
		liquidity, err = poolContract.Liquidity(cops)
		if err != nil {
			return err
		}
		if liquidity.Cmp(big.NewInt(0)) == 0 {
			return errors.New("pool has no liquidity")
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("Unable to provide liquidity to the pool")
		return err
	}
	log.Debug().Interface("liquidity", liquidity).Msg("Liquidity provided to the pool")
	return nil
}

// Get the timestamp of the latest block.
func getLatestBlockTimestamp(ctx context.Context, c *ethclient.Client) (*big.Int, error) {
	// Get latest block number.
	blockNumber, err := c.BlockNumber(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to get the latest block number")
		return big.NewInt(0), err
	}

	// Get latest block.
	var block *types.Block
	block, err = c.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		log.Error().Err(err).Msg("Unable to get the latest block")
		return big.NewInt(0), err
	}
	return big.NewInt(int64(block.Time())), nil
}
