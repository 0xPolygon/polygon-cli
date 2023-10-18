package loadtest

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"time"

	"github.com/spf13/cobra"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	uniswapv3loadtest "github.com/maticnetwork/polygon-cli/cmd/loadtest/uniswapv3"
	"github.com/maticnetwork/polygon-cli/contracts/uniswapv3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	//go:embed uniswapv3Usage.md
	uniswapv3Usage          string
	uniswapv3LoadTestParams params

	tokenPoolSize, _ = big.NewInt(0).SetString("100000000000000000000000000", 10)
)

type params struct {
	UniswapFactoryV3, UniswapMulticall, UniswapProxyAdmin, UniswapTickLens, UniswapNFTLibDescriptor, UniswapNonfungibleTokenPositionDescriptor, UniswapUpgradeableProxy, UniswapNonfungiblePositionManager, UniswapMigrator, UniswapStaker, UniswapQuoterV2, UniswapSwapRouter *string
	WETH9, UniswapPoolToken0, UniswapPoolToken1                                                                                                                                                                                                                                *string
}

var uniswapV3LoadTestCmd = &cobra.Command{
	Use:   "uniswapv3 url",
	Short: "Run Uniswapv3-like load test against an Eth/EVm style JSON-RPC endpoint.",
	Long:  uniswapv3Usage,
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
	Args: func(cmd *cobra.Command, args []string) error {
		zerolog.DurationFieldUnit = time.Second
		zerolog.DurationFieldInteger = true

		if len(args) != 1 {
			return fmt.Errorf("expected exactly one argument")
		}

		url, err := validateUrl(args[0])
		if err != nil {
			return err
		}
		inputLoadTestParams.URL = url

		return nil
	},
}

func validateUrl(input string) (*url.URL, error) {
	url, err := url.Parse(input)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse url input error")
		return nil, err
	}

	if url.Scheme == "" {
		return nil, errors.New("the scheme has not been specified")
	}
	switch url.Scheme {
	case "http", "https", "ws", "wss":
		return url, nil
	default:
		return nil, fmt.Errorf("the scheme %s is not supported", url.Scheme)
	}
}

func init() {
	// Specify subcommand flags.
	params := new(params)
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
	uniswapv3LoadTestParams = *params
}

func deploySwapperContract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts, config uniswapv3loadtest.UniswapV3Config, tokenName, tokenSymbol string, recipient common.Address, tokenKnownAddress common.Address) (uniswapv3loadtest.ContractConfig[uniswapv3.Swapper], error) {
	var token uniswapv3loadtest.ContractConfig[uniswapv3.Swapper]
	var err error
	addressesToApprove := []common.Address{config.NonfungiblePositionManager.Address, config.SwapRouter02.Address}
	token.Address, token.Contract, err = uniswapv3loadtest.DeployOrInstantiateContract(
		ctx, c, tops, cops, tokenName, tokenKnownAddress,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *uniswapv3.Swapper, error) {
			return uniswapv3.DeploySwapper(tops, c)
		},
		uniswapv3.NewSwapper,
		func(contract *uniswapv3.Swapper) error {
			return approveSwapperSpendingsByUniswap(ctx, contract, tops, cops, addressesToApprove, recipient)
		},
		blockUntilSuccessful,
	)
	if err != nil {
		return token, err
	}
	return token, nil
}

func approveSwapperSpendingsByUniswap(ctx context.Context, contract *uniswapv3.Swapper, tops *bind.TransactOpts, cops *bind.CallOpts, addresses []common.Address, owner common.Address) error {
	name, err := contract.Name(cops)
	if err != nil {
		return err

	}
	for _, address := range addresses {
		tx, err := contract.Approve(tops, address, tokenPoolSize)
		if err != nil {
			log.Error().Err(err).Interface("address", address).Msg("Unable to approve spendings")
			return err
		}

		backoff.Retry(func() error {
			allowance, err := contract.Allowance(cops, owner, address)
			if err != nil {
				return err
			}
			zero := big.NewInt(0)
			if allowance.Cmp(zero) == 0 {
				return fmt.Errorf("allowance is zero")
			}
			return nil
		}, backoff.NewConstantBackOff(time.Second*2))

		log.Debug().Str("Swapper", name).Str("spender", address.String()).Str("amount", tokenPoolSize.String()).Msg("Spending approved")
		log.Trace().Interface("hash", tx.Hash()).Msg("Transaction")
	}
	return nil
}

// Run UniswapV3 loadtest by performing swaps.
func loadTestUniswapV3(ctx context.Context, c *ethclient.Client, nonce uint64, uniswapV3Config uniswapv3loadtest.UniswapV3Config, poolConfig uniswapv3loadtest.PoolConfig) (t1 time.Time, t2 time.Time, err error) {
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
	err = uniswapv3loadtest.ExactInputSingleSwap(tops, uniswapV3Config.SwapRouter02.Contract, poolConfig, *ltp.FromETHAddress, nonce)
	return
}
