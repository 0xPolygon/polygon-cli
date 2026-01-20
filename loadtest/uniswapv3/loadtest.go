package uniswapv3

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygon/polygon-cli/bindings/tokens"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

// InitParams holds parameters for initializing the UniswapV3 loadtest.
type InitParams struct {
	// Pre-deployed token addresses (optional - if empty, new tokens are deployed)
	PoolToken0Address common.Address
	PoolToken1Address common.Address

	// Pool configuration (fee tier as *big.Int from PercentageToUniswapFeeTier)
	PoolFees *big.Int
}

// Init initializes the UniswapV3 loadtest by deploying contracts and setting up the pool.
// Returns the UniswapV3 config and pool config needed for running swaps.
func Init(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts, uniswapAddresses UniswapV3Addresses, recipient common.Address, params InitParams) (uniswapV3Config UniswapV3Config, poolConfig PoolConfig, err error) {
	log.Info().Msg("Deploying UniswapV3 contracts...")
	uniswapV3Config, err = DeployUniswapV3(ctx, c, tops, cops, uniswapAddresses, recipient)
	if err != nil {
		return
	}
	log.Info().Interface("addresses", uniswapV3Config.GetAddresses()).Msg("UniswapV3 deployed")

	log.Info().Msg("Deploying ERC20 tokens...")
	var token0 ContractConfig[tokens.ERC20]
	token0, err = DeployERC20(
		ctx, c, tops, cops, uniswapV3Config, "SwapperA", "SA", MintAmount, recipient, params.PoolToken0Address)
	if err != nil {
		return
	}

	var token1 ContractConfig[tokens.ERC20]
	token1, err = DeployERC20(
		ctx, c, tops, cops, uniswapV3Config, "SwapperB", "SB", MintAmount, recipient, params.PoolToken1Address)
	if err != nil {
		return
	}

	poolConfig = *NewPool(token0, token1, params.PoolFees)

	// Skip liquidity pool setup if using pre-deployed tokens
	if params.PoolToken0Address != (common.Address{}) {
		return
	}

	if err = SetupLiquidityPool(ctx, c, tops, cops, uniswapV3Config, poolConfig, recipient); err != nil {
		return
	}

	log.Info().
		Stringer("--uniswap-factory-v3-address", uniswapV3Config.FactoryV3.Address).
		Stringer("--uniswap-migrator-address", uniswapV3Config.Migrator.Address).
		Stringer("--uniswap-multicall-address", uniswapV3Config.Multicall.Address).
		Stringer("--uniswap-nft-descriptor-lib-address", uniswapV3Config.NFTDescriptorLib.Address).
		Stringer("--uniswap-nft-position-descriptor-address", uniswapV3Config.NonfungibleTokenPositionDescriptor.Address).
		Stringer("--uniswap-non-fungible-position-manager-address", uniswapV3Config.NonfungiblePositionManager.Address).
		Stringer("--uniswap-pool-token-0-address", poolConfig.Token0.Address).
		Stringer("--uniswap-pool-token-1-address", poolConfig.Token1.Address).
		Stringer("--uniswap-proxy-admin-address", uniswapV3Config.ProxyAdmin.Address).
		Stringer("--uniswap-quoter-v2-address", uniswapV3Config.QuoterV2.Address).
		Stringer("--uniswap-staker-address", uniswapV3Config.Staker.Address).
		Stringer("--uniswap-swap-router-address", uniswapV3Config.SwapRouter02.Address).
		Stringer("--uniswap-tick-lens-address", uniswapV3Config.TickLens.Address).
		Stringer("--uniswap-upgradeable-proxy-address", uniswapV3Config.TransparentUpgradeableProxy.Address).
		Stringer("--weth9-address", uniswapV3Config.WETH9.Address).
		Msg("Parameters to re-run")

	return
}

// Run performs a single UniswapV3 swap operation.
// Returns the start time, end time, transaction hash, and any error.
func Run(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, uniswapV3Config UniswapV3Config, poolConfig PoolConfig, swapAmountIn *big.Int, recipient common.Address) (start, end time.Time, txHash common.Hash, err error) {
	var tx *ethtypes.Transaction

	start = time.Now()
	defer func() { end = time.Now() }()

	tx, err = ExactInputSingleSwap(tops, uniswapV3Config.SwapRouter02.Contract, poolConfig, swapAmountIn, recipient, tops.Nonce.Uint64())
	if err == nil && tx != nil {
		txHash = tx.Hash()
	}
	return
}
