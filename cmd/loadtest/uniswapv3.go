package loadtest

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygon/polygon-cli/bindings/tokens"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/spf13/cobra"

	uniswapv3loadtest "github.com/0xPolygon/polygon-cli/cmd/loadtest/uniswapv3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

var (
	//go:embed uniswapv3Usage.md
	uniswapv3Usage          string
	uniswapv3LoadTestParams uniswap3params
)

var uniswapV3LoadTestCmd = &cobra.Command{
	Use:   "uniswapv3",
	Short: "Run Uniswapv3-like load test against an Eth/EVm style JSON-RPC endpoint.",
	Long:  uniswapv3Usage,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkUniswapV3LoadtestFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Override root command `mode` flag.
		inputLoadTestParams.Modes = []string{"v3"}

		// Run load test.
		err := runLoadTest(cmd.Context())
		if err != nil {
			return err
		}
		return nil
	},
}

func checkUniswapV3LoadtestFlags() error {
	// Check pool fees.
	switch fees := uniswapv3LoadTestParams.PoolFees; fees {
	case float64(uniswapv3loadtest.StableTier), float64(uniswapv3loadtest.StandardTier), float64(uniswapv3loadtest.ExoticTier):
		// Fees are correct, do nothing.
	default:
		return fmt.Errorf("UniswapV3 only supports a few pool tiers which are stable: %f%%, standard: %f%%, and exotic: %f%%",
			float64(uniswapv3loadtest.StableTier), float64(uniswapv3loadtest.StandardTier), float64(uniswapv3loadtest.ExoticTier))
	}

	// Check swap amount input.
	if uniswapv3LoadTestParams.SwapAmountInput == 0 {
		return errors.New("swap amount input has to be greater than zero")
	}

	if (uniswapv3LoadTestParams.UniswapPoolToken0 != "") != (uniswapv3LoadTestParams.UniswapPoolToken1 != "") {
		return errors.New("both pool tokens must be empty or specified. Specifying only one token is not allowed")
	}
	return nil
}

type uniswap3params struct {
	UniswapFactoryV3, UniswapMulticall, UniswapProxyAdmin, UniswapTickLens, UniswapNFTLibDescriptor, UniswapNonfungibleTokenPositionDescriptor, UniswapUpgradeableProxy, UniswapNonfungiblePositionManager, UniswapMigrator, UniswapStaker, UniswapQuoterV2, UniswapSwapRouter, WETH9, UniswapPoolToken0, UniswapPoolToken1 string
	PoolFees                                                                                                                                                                                                                                                                                                                float64
	SwapAmountInput                                                                                                                                                                                                                                                                                                         uint64
}

func init() {
	// Specify subcommand flags.
	params := &uniswapv3LoadTestParams

	// Pre-deployed addresses.
	f := uniswapV3LoadTestCmd.Flags()
	f.StringVar(&params.UniswapFactoryV3, "uniswap-factory-v3-address", "", "address of pre-deployed UniswapFactoryV3 contract")
	f.StringVar(&params.UniswapMulticall, "uniswap-multicall-address", "", "address of pre-deployed Multicall contract")
	f.StringVar(&params.UniswapProxyAdmin, "uniswap-proxy-admin-address", "", "address of pre-deployed ProxyAdmin contract")
	f.StringVar(&params.UniswapTickLens, "uniswap-tick-lens-address", "", "address of pre-deployed TickLens contract")
	f.StringVar(&params.UniswapNFTLibDescriptor, "uniswap-nft-descriptor-lib-address", "", "address of pre-deployed NFTDescriptor library contract")
	f.StringVar(&params.UniswapNonfungibleTokenPositionDescriptor, "uniswap-nft-position-descriptor-address", "", "address of pre-deployed NonfungibleTokenPositionDescriptor contract")
	f.StringVar(&params.UniswapUpgradeableProxy, "uniswap-upgradeable-proxy-address", "", "address of pre-deployed TransparentUpgradeableProxy contract")
	f.StringVar(&params.UniswapNonfungiblePositionManager, "uniswap-non-fungible-position-manager-address", "", "address of pre-deployed NonfungiblePositionManager contract")
	f.StringVar(&params.UniswapMigrator, "uniswap-migrator-address", "", "address of pre-deployed Migrator contract")
	f.StringVar(&params.UniswapStaker, "uniswap-staker-address", "", "address of pre-deployed Staker contract")
	f.StringVar(&params.UniswapQuoterV2, "uniswap-quoter-v2-address", "", "address of pre-deployed QuoterV2 contract")
	f.StringVar(&params.UniswapSwapRouter, "uniswap-swap-router-address", "", "address of pre-deployed SwapRouter contract")
	f.StringVar(&params.WETH9, "weth9-address", "", "address of pre-deployed WETH9 contract")
	f.StringVar(&params.UniswapPoolToken0, "uniswap-pool-token-0-address", "", "address of pre-deployed ERC20 contract used in Uniswap pool Token0 // Token1")
	f.StringVar(&params.UniswapPoolToken1, "uniswap-pool-token-1-address", "", "address of pre-deployed ERC20 contract used in Uniswap pool Token0 // Token1")

	// Pool and swap parameters.
	f.Float64VarP(&params.PoolFees, "pool-fees", "f", float64(uniswapv3loadtest.StandardTier), "trading fees for UniswapV3 liquidity pool swaps (e.g. 0.3 means 0.3%)")
	f.Uint64VarP(&params.SwapAmountInput, "swap-amount", "a", uniswapv3loadtest.SwapAmountInput.Uint64(), "amount of inbound token given as swap input")
}

// Initialise UniswapV3 loadtest.
func initUniswapV3Loadtest(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts, uniswapAddresses uniswapv3loadtest.UniswapV3Addresses, recipient common.Address) (uniswapV3Config uniswapv3loadtest.UniswapV3Config, poolConfig uniswapv3loadtest.PoolConfig, err error) {
	log.Info().Msg("Deploying UniswapV3 contracts...")
	uniswapV3Config, err = uniswapv3loadtest.DeployUniswapV3(ctx, c, tops, cops, uniswapAddresses, recipient)
	if err != nil {
		return
	}
	log.Info().Interface("addresses", uniswapV3Config.GetAddresses()).Msg("UniswapV3 deployed")

	log.Info().Msg("Deploying ERC20 tokens...")
	var token0 uniswapv3loadtest.ContractConfig[tokens.ERC20]
	token0, err = uniswapv3loadtest.DeployERC20(
		ctx, c, tops, cops, uniswapV3Config, "SwapperA", "SA", uniswapv3loadtest.MintAmount, recipient, common.HexToAddress(uniswapv3LoadTestParams.UniswapPoolToken0))
	if err != nil {
		return
	}

	var token1 uniswapv3loadtest.ContractConfig[tokens.ERC20]
	token1, err = uniswapv3loadtest.DeployERC20(
		ctx, c, tops, cops, uniswapV3Config, "SwapperB", "SB", uniswapv3loadtest.MintAmount, recipient, common.HexToAddress(uniswapv3LoadTestParams.UniswapPoolToken1))
	if err != nil {
		return
	}

	fees := uniswapv3loadtest.PercentageToUniswapFeeTier(uniswapv3LoadTestParams.PoolFees)
	poolConfig = *uniswapv3loadtest.NewPool(token0, token1, fees)
	if uniswapv3LoadTestParams.UniswapPoolToken0 != "" {
		return
	}
	if err = uniswapv3loadtest.SetupLiquidityPool(ctx, c, tops, cops, uniswapV3Config, poolConfig, recipient); err != nil {
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
		Stringer("--weth9-address", uniswapV3Config.WETH9.Address).Msg("Parameters to re-run")

	return
}

// Run UniswapV3 loadtest.
func runUniswapV3Loadtest(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, uniswapV3Config uniswapv3loadtest.UniswapV3Config, poolConfig uniswapv3loadtest.PoolConfig, swapAmountIn *big.Int) (t1 time.Time, t2 time.Time, txHash common.Hash, err error) {
	var tx *ethtypes.Transaction

	ltp := inputLoadTestParams

	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	tx, err = uniswapv3loadtest.ExactInputSingleSwap(tops, uniswapV3Config.SwapRouter02.Contract, poolConfig, swapAmountIn, *ltp.FromETHAddress, tops.Nonce.Uint64())
	if err == nil && tx != nil {
		txHash = tx.Hash()
	}
	return
}
