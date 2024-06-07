package loadtest

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maticnetwork/polygon-cli/bindings/uniswapv3"
	uniswapv3loadtest "github.com/maticnetwork/polygon-cli/cmd/loadtest/uniswapv3"
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
		inputLoadTestParams.Modes = &[]string{"v3"}

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
	switch fees := *uniswapv3LoadTestParams.PoolFees; fees {
	case float64(uniswapv3loadtest.StableTier), float64(uniswapv3loadtest.StandardTier), float64(uniswapv3loadtest.ExoticTier):
		// Fees are correct, do nothing.
	default:
		return fmt.Errorf("UniswapV3 only supports a few pool tiers which are stable: %f%%, standard: %f%%, and exotic: %f%%",
			float64(uniswapv3loadtest.StableTier), float64(uniswapv3loadtest.StandardTier), float64(uniswapv3loadtest.ExoticTier))
	}

	// Check swap amount input.
	if *uniswapv3LoadTestParams.SwapAmountInput == 0 {
		return errors.New("swap amount input has to be greater than zero")
	}
	return nil
}

type uniswap3params struct {
	UniswapFactoryV3, UniswapMulticall, UniswapProxyAdmin, UniswapTickLens, UniswapNFTLibDescriptor, UniswapNonfungibleTokenPositionDescriptor, UniswapUpgradeableProxy, UniswapNonfungiblePositionManager, UniswapMigrator, UniswapStaker, UniswapQuoterV2, UniswapSwapRouter, WETH9, UniswapPoolToken0, UniswapPoolToken1 *string
	PoolFees                                                                                                                                                                                                                                                                                                                *float64
	SwapAmountInput                                                                                                                                                                                                                                                                                                         *uint64
}

func init() {
	// Specify subcommand flags.
	params := new(uniswap3params)

	// Pre-deployed addresses.
	params.UniswapFactoryV3 = uniswapV3LoadTestCmd.Flags().String("uniswap-factory-v3-address", "", "The address of a pre-deployed UniswapFactoryV3 contract")
	params.UniswapMulticall = uniswapV3LoadTestCmd.Flags().String("uniswap-multicall-address", "", "The address of a pre-deployed Multicall contract")
	params.UniswapProxyAdmin = uniswapV3LoadTestCmd.Flags().String("uniswap-proxy-admin-address", "", "The address of a pre-deployed ProxyAdmin contract")
	params.UniswapTickLens = uniswapV3LoadTestCmd.Flags().String("uniswap-tick-lens-address", "", "The address of a pre-deployed TickLens contract")
	params.UniswapNFTLibDescriptor = uniswapV3LoadTestCmd.Flags().String("uniswap-nft-descriptor-lib-address", "", "The address of a pre-deployed NFTDescriptor library contract")
	params.UniswapNonfungibleTokenPositionDescriptor = uniswapV3LoadTestCmd.Flags().String("uniswap-nft-position-descriptor-address", "", "The address of a pre-deployed NonfungibleTokenPositionDescriptor contract")
	params.UniswapUpgradeableProxy = uniswapV3LoadTestCmd.Flags().String("uniswap-upgradeable-proxy-address", "", "The address of a pre-deployed TransparentUpgradeableProxy contract")
	params.UniswapNonfungiblePositionManager = uniswapV3LoadTestCmd.Flags().String("uniswap-non-fungible-position-manager-address", "", "The address of a pre-deployed NonfungiblePositionManager contract")
	params.UniswapMigrator = uniswapV3LoadTestCmd.Flags().String("uniswap-migrator-address", "", "The address of a pre-deployed Migrator contract")
	params.UniswapStaker = uniswapV3LoadTestCmd.Flags().String("uniswap-staker-address", "", "The address of a pre-deployed Staker contract")
	params.UniswapQuoterV2 = uniswapV3LoadTestCmd.Flags().String("uniswap-quoter-v2-address", "", "The address of a pre-deployed QuoterV2 contract")
	params.UniswapSwapRouter = uniswapV3LoadTestCmd.Flags().String("uniswap-swap-router-address", "", "The address of a pre-deployed SwapRouter contract")
	params.WETH9 = uniswapV3LoadTestCmd.Flags().String("weth9-address", "", "The address of a pre-deployed WETH9 contract")
	params.UniswapPoolToken0 = uniswapV3LoadTestCmd.Flags().String("uniswap-pool-token-0-address", "", "The address of a pre-deployed ERC20 contract used in the Uniswap pool Token0 // Token1")
	params.UniswapPoolToken1 = uniswapV3LoadTestCmd.Flags().String("uniswap-pool-token-1-address", "", "The address of a pre-deployed ERC20 contract used in the Uniswap pool Token0 // Token1")

	// Pool and swap parameters.
	params.PoolFees = uniswapV3LoadTestCmd.Flags().Float64P("pool-fees", "f", float64(uniswapv3loadtest.StandardTier), "Trading fees charged on each swap or trade made within a UniswapV3 liquidity pool (e.g. 0.3 means 0.3%)")
	params.SwapAmountInput = uniswapV3LoadTestCmd.Flags().Uint64P("swap-amount", "a", uniswapv3loadtest.SwapAmountInput.Uint64(), "The amount of inbound token given as swap input")

	uniswapv3LoadTestParams = *params
}

// Initialise UniswapV3 loadtest.
func initUniswapV3Loadtest(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts, uniswapAddresses uniswapv3loadtest.UniswapV3Addresses, recipient common.Address) (uniswapV3Config uniswapv3loadtest.UniswapV3Config, poolConfig uniswapv3loadtest.PoolConfig, err error) {
	log.Debug().Msg("🦄 Deploying UniswapV3 contracts...")
	uniswapV3Config, err = uniswapv3loadtest.DeployUniswapV3(ctx, c, tops, cops, uniswapAddresses, recipient)
	if err != nil {
		return
	}
	log.Debug().Interface("addresses", uniswapV3Config.GetAddresses()).Msg("UniswapV3 deployed")

	log.Debug().Msg("🪙 Deploying ERC20 tokens...")
	var token0 uniswapv3loadtest.ContractConfig[uniswapv3.Swapper]
	token0, err = uniswapv3loadtest.DeployERC20(
		ctx, c, tops, cops, uniswapV3Config, "SwapperA", "SA", uniswapv3loadtest.MintAmount, recipient, common.HexToAddress(*uniswapv3LoadTestParams.UniswapPoolToken0))
	if err != nil {
		return
	}

	var token1 uniswapv3loadtest.ContractConfig[uniswapv3.Swapper]
	token1, err = uniswapv3loadtest.DeployERC20(
		ctx, c, tops, cops, uniswapV3Config, "SwapperB", "SB", uniswapv3loadtest.MintAmount, recipient, common.HexToAddress(*uniswapv3LoadTestParams.UniswapPoolToken1))
	if err != nil {
		return
	}

	log.Debug().Msg("🎱 Deploying UniswapV3 liquidity pool...")
	fees := uniswapv3loadtest.PercentageToUniswapFeeTier(*uniswapv3LoadTestParams.PoolFees)
	poolConfig = *uniswapv3loadtest.NewPool(token0, token1, fees)
	if err = uniswapv3loadtest.SetupLiquidityPool(ctx, c, tops, cops, uniswapV3Config, poolConfig, recipient); err != nil {
		return
	}
	return
}

// Run UniswapV3 loadtest.
func runUniswapV3Loadtest(ctx context.Context, c *ethclient.Client, nonce uint64, uniswapV3Config uniswapv3loadtest.UniswapV3Config, poolConfig uniswapv3loadtest.PoolConfig, swapAmountIn *big.Int) (t1 time.Time, t2 time.Time, err error) {
	ltp := inputLoadTestParams
	chainID := new(big.Int).SetUint64(*ltp.ChainID)
	privateKey := ltp.ECDSAPrivateKey

	tops, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Unable create transaction signer")
		return
	}
	tops.Nonce = new(big.Int).SetUint64(nonce)
	tops = configureTransactOpts(tops)

	t1 = time.Now()
	defer func() { t2 = time.Now() }()
	err = uniswapv3loadtest.ExactInputSingleSwap(tops, uniswapV3Config.SwapRouter02.Contract, poolConfig, swapAmountIn, *ltp.FromETHAddress, nonce)
	return
}
