package uniswapv3loadtest

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	openzeppelin "github.com/maticnetwork/polygon-cli/contracts/src/openzeppelin"
	tokens "github.com/maticnetwork/polygon-cli/contracts/src/tokens"
	v3core "github.com/maticnetwork/polygon-cli/contracts/src/uniswap/v3core"
	v3periphery "github.com/maticnetwork/polygon-cli/contracts/src/uniswap/v3periphery"
	v3router "github.com/maticnetwork/polygon-cli/contracts/src/uniswap/v3router"
	v3staker "github.com/maticnetwork/polygon-cli/contracts/src/uniswap/v3staker"
	"github.com/rs/zerolog/log"
)

const (
	// The NFTPositionLib link (or address) in NFTPositionDescriptor bytecode.
	// When recompiling the contracts and updating the go bindings, make sure to update this value.
	oldNFTPositionLibraryAddress = "__$cea9be979eee3d87fb124d6cbb244bb0b5$__"

	// The fee amount to enable for one basic point.
	oneBPFee = 100
	// The spacing between ticks to be enforced for all pools with the given fee amount.
	oneBPTickSpacing = 1

	// The max amount of seconds into the future the incentive startTime can be set.
	maxIncentiveStartLeadTime = 30 * 24 * 60 * 60 // 1 month
	// The max duration of an incentive in seconds.
	maxIncentiveDuration = 2 * 365 * 23 * 60 * 60 * 2 // 2 years
)

type (
	// UniswapV3Config represents the whole UniswapV3 configuration (contracts and addresses), including WETH9.
	UniswapV3Config struct {
		FactoryV3                          ContractConfig[v3core.UniswapV3Factory]
		Multicall                          ContractConfig[v3periphery.UniswapInterfaceMulticall]
		ProxyAdmin                         ContractConfig[openzeppelin.ProxyAdmin]
		TickLens                           ContractConfig[v3periphery.TickLens]
		NFTDescriptorLib                   ContractConfig[v3periphery.NFTDescriptor]
		NonfungibleTokenPositionDescriptor ContractConfig[v3periphery.NonfungibleTokenPositionDescriptor]
		TransparentUpgradeableProxy        ContractConfig[openzeppelin.TransparentUpgradeableProxy]
		NonfungiblePositionManager         ContractConfig[v3periphery.NonfungiblePositionManager]
		Migrator                           ContractConfig[v3periphery.V3Migrator]
		Staker                             ContractConfig[v3staker.UniswapV3Staker]
		QuoterV2                           ContractConfig[v3router.QuoterV2]
		SwapRouter02                       ContractConfig[v3router.SwapRouter02]

		WETH9 ContractConfig[tokens.WETH9]
	}

	// UniswapV3Addresses is a subset of UniswapV3Config. It represents the addresses of the whole
	// UniswapV3 configuration, including WETH9.
	UniswapV3Addresses struct {
		FactoryV3, Multicall, ProxyAdmin, TickLens, NFTDescriptorLib, NonfungibleTokenPositionDescriptor, TransparentUpgradeableProxy, NonfungiblePositionManager, Migrator, Staker, QuoterV2, SwapRouter02, WETH9 common.Address
	}

	// ContractConfig represents a contract and its address.
	ContractConfig[T Contract] struct {
		Address  common.Address
		Contract *T
	}

	// Contract represents a UniswapV3 contract (including WETH9 and Swapper).
	Contract interface {
		tokens.WETH9 | v3core.UniswapV3Factory | v3periphery.UniswapInterfaceMulticall | openzeppelin.ProxyAdmin | v3periphery.TickLens | v3periphery.NFTDescriptor | v3periphery.NonfungibleTokenPositionDescriptor | openzeppelin.TransparentUpgradeableProxy | v3periphery.NonfungiblePositionManager | v3periphery.V3Migrator | v3staker.UniswapV3Staker | v3router.QuoterV2 | v3router.SwapRouter02 | tokens.Swapper
	}
)

// Deploy the full UniswapV3 contract suite in 15 different steps.
// Source: https://github.com/Uniswap/deploy-v3
func DeployUniswapV3(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts, knownAddresses UniswapV3Addresses, ownerAddress common.Address, blockblockUntilSuccessful blockUntilSuccessfulFn) (config UniswapV3Config, err error) {
	log.Debug().Msg("Step 0: WETH9 deployment")
	config.WETH9.Address, config.WETH9.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		knownAddresses.WETH9,
		tokens.DeployWETH9,
		tokens.NewWETH9,
		func(contract *tokens.WETH9) (err error) {
			_, err = contract.BalanceOf(cops, common.Address{})
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return
	}

	log.Debug().Msg("Step 1: UniswapV3Factory deployment")
	config.FactoryV3.Address, config.FactoryV3.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		knownAddresses.FactoryV3,
		v3core.DeployUniswapV3Factory,
		v3core.NewUniswapV3Factory,
		func(contract *v3core.UniswapV3Factory) (err error) {
			_, err = contract.Owner(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return
	}

	log.Debug().Msg("Step 2: Enable fee amount")
	if err = enableFeeAmount(config.FactoryV3.Contract, tops, cops, oneBPFee, oneBPTickSpacing); err != nil {
		return
	}

	log.Debug().Msg("Step 3: UniswapInterfaceMulticall deployment")
	config.Multicall.Address, config.Multicall.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		knownAddresses.Multicall,
		v3periphery.DeployUniswapInterfaceMulticall,
		v3periphery.NewUniswapInterfaceMulticall,
		func(contract *v3periphery.UniswapInterfaceMulticall) (err error) {
			_, err = contract.GetEthBalance(cops, common.Address{})
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return
	}

	log.Debug().Msg("Step 4: ProxyAdmin deployment")
	config.ProxyAdmin.Address, config.ProxyAdmin.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		knownAddresses.ProxyAdmin,
		openzeppelin.DeployProxyAdmin,
		openzeppelin.NewProxyAdmin,
		func(contract *openzeppelin.ProxyAdmin) (err error) {
			_, err = contract.Owner(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return
	}

	log.Debug().Msg("Step 5: TickLens deployment")
	config.TickLens.Address, config.TickLens.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		knownAddresses.TickLens,
		v3periphery.DeployTickLens,
		v3periphery.NewTickLens,
		func(contract *v3periphery.TickLens) (err error) {
			// The only function we can call to check the contract is deployed is `GetPopulatedTicksInWord`.
			// Unfortunately, such call will revert because no pools are deployed yet.
			// That's why we only return a nil value here.
			return nil
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return
	}

	log.Debug().Msg("Step 6: NFTDescriptorLib deployment")
	config.NFTDescriptorLib.Address, _, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		knownAddresses.NFTDescriptorLib,
		v3periphery.DeployNFTDescriptor,
		v3periphery.NewNFTDescriptor,
		func(contract *v3periphery.NFTDescriptor) (err error) {
			// The only method we could call requires a pool to be deployed.
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return
	}

	log.Debug().Msg("Step 7: NFTPositionDescriptor deployment")
	config.NonfungibleTokenPositionDescriptor.Address, config.NonfungibleTokenPositionDescriptor.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		knownAddresses.NonfungibleTokenPositionDescriptor,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *v3periphery.NonfungibleTokenPositionDescriptor, error) {
			// Update NFTPosition library address in NFTPositionDescriptor bytecode.
			newNFTPositionLibraryAddress := strings.TrimPrefix(strings.ToLower(config.NFTDescriptorLib.Address.String()), "0x")
			newNonfungibleTokenPositionDescriptorBytecode := strings.ReplaceAll(v3periphery.NonfungibleTokenPositionDescriptorMetaData.Bin, oldNFTPositionLibraryAddress, newNFTPositionLibraryAddress)
			if v3periphery.NonfungibleTokenPositionDescriptorMetaData.Bin == newNonfungibleTokenPositionDescriptorBytecode {
				return common.Address{}, nil, nil, errors.New("NFTPositionDescriptor bytecode has not been updated")
			}

			var nativeCurrencyLabelBytes [32]byte
			copy(nativeCurrencyLabelBytes[:], "ETH")
			v3periphery.NonfungibleTokenPositionDescriptorMetaData.Bin = newNonfungibleTokenPositionDescriptorBytecode
			v3periphery.NonfungibleTokenPositionDescriptorBin = newNonfungibleTokenPositionDescriptorBytecode
			log.Trace().Interface("oldAddress", oldNFTPositionLibraryAddress).Interface("newAddress", config.NFTDescriptorLib.Address).Msg("NFTPositionDescriptor bytecode updated with the new NFTDescriptor library address")

			// Deploy NFTPositionDescriptor contract.
			return v3periphery.DeployNonfungibleTokenPositionDescriptor(tops, c, config.WETH9.Address, nativeCurrencyLabelBytes)
		},
		v3periphery.NewNonfungibleTokenPositionDescriptor,
		func(contract *v3periphery.NonfungibleTokenPositionDescriptor) (err error) {
			_, err = contract.WETH9(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return
	}

	log.Debug().Msg("Step 8: TransparentUpgradeableProxy deployment")
	config.TransparentUpgradeableProxy.Address, config.TransparentUpgradeableProxy.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		knownAddresses.TransparentUpgradeableProxy,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *openzeppelin.TransparentUpgradeableProxy, error) {
			return openzeppelin.DeployTransparentUpgradeableProxy(tops, c, config.NonfungibleTokenPositionDescriptor.Address, config.ProxyAdmin.Address, []byte(""))
		},
		openzeppelin.NewTransparentUpgradeableProxy,
		func(contract *openzeppelin.TransparentUpgradeableProxy) (err error) {
			// The TransparentUpgradeableProxy contract methods can only be called by the admin.
			// This is not a problem when we first deploy the contract because the deployer is set to be
			// the admin by default. Thus, we can call any method of the contract to check it has been deployed.
			// But when we use pre-deployed contracts, since the TransparentUpgradeableProxy ownership
			// has been transferred, we get "execution reverted" errors when trying to call any method.
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return
	}

	log.Debug().Msg("Step 9: NonfungiblePositionManager deployment")
	config.NonfungiblePositionManager.Address, config.NonfungiblePositionManager.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		knownAddresses.NonfungiblePositionManager,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *v3periphery.NonfungiblePositionManager, error) {
			return v3periphery.DeployNonfungiblePositionManager(tops, c, config.FactoryV3.Address, config.WETH9.Address, config.TransparentUpgradeableProxy.Address)
		},
		v3periphery.NewNonfungiblePositionManager,
		func(contract *v3periphery.NonfungiblePositionManager) (err error) {
			_, err = contract.BaseURI(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return
	}

	log.Debug().Msg("Step 10: V3Migrator deployment")
	config.Migrator.Address, config.Migrator.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		knownAddresses.Migrator,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *v3periphery.V3Migrator, error) {
			return v3periphery.DeployV3Migrator(tops, c, config.FactoryV3.Address, config.WETH9.Address, config.NonfungiblePositionManager.Address)
		},
		v3periphery.NewV3Migrator,
		func(contract *v3periphery.V3Migrator) (err error) {
			_, err = contract.WETH9(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return
	}

	log.Debug().Msg("Step 11: Transfer UniswapV3Factory ownership")
	if err = transferUniswapV3FactoryOwnership(config.FactoryV3.Contract, tops, cops, ownerAddress); err != nil {
		return
	}

	log.Debug().Msg("Step 12: UniswapV3Staker deployment")
	config.Staker.Address, config.Staker.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		knownAddresses.Staker,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *v3staker.UniswapV3Staker, error) {
			return v3staker.DeployUniswapV3Staker(tops, c, config.FactoryV3.Address, config.NonfungiblePositionManager.Address, big.NewInt(maxIncentiveStartLeadTime), big.NewInt(maxIncentiveDuration))
		},
		v3staker.NewUniswapV3Staker,
		func(contract *v3staker.UniswapV3Staker) (err error) {
			_, err = contract.Factory(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return
	}

	log.Debug().Msg("Step 13: QuoterV2 deployment")
	config.QuoterV2.Address, config.QuoterV2.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		knownAddresses.QuoterV2,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *v3router.QuoterV2, error) {
			return v3router.DeployQuoterV2(tops, c, config.FactoryV3.Address, config.WETH9.Address)
		},
		v3router.NewQuoterV2,
		func(contract *v3router.QuoterV2) (err error) {
			_, err = contract.Factory(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return
	}

	log.Debug().Msg("Step 14: SwapRouter02 deployment")
	config.SwapRouter02.Address, config.SwapRouter02.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		knownAddresses.SwapRouter02,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *v3router.SwapRouter02, error) {
			uniswapFactoryV2Address := common.Address{} // Note: we specify an empty address for UniswapV2Factory since we don't deploy it.
			return v3router.DeploySwapRouter02(tops, c, uniswapFactoryV2Address, config.FactoryV3.Address, config.NonfungiblePositionManager.Address, config.WETH9.Address)
		},
		v3router.NewSwapRouter02,
		func(contract *v3router.SwapRouter02) (err error) {
			_, err = contract.Factory(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return
	}

	log.Debug().Msg("Step 15: Transfer ProxyAdmin ownership")
	if err = transferProxyAdminOwnership(config.ProxyAdmin.Contract, tops, cops, ownerAddress); err != nil {
		return
	}

	return
}

// deployOrInstantiateContract deploys or instantiates a UniswapV3 contract.
// If knownAddress is empty, it deploys the contract; otherwise, it instantiates it.
func deployOrInstantiateContract[T Contract](
	ctx context.Context,
	c *ethclient.Client,
	tops *bind.TransactOpts,
	cops *bind.CallOpts,
	knownAddress common.Address,
	deploy func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *T, error),
	instantiate func(common.Address, bind.ContractBackend) (*T, error),
	call func(*T) error,
	blockUntilSuccessful blockUntilSuccessfulFn,
) (address common.Address, contract *T, err error) {
	if knownAddress == (common.Address{}) {
		// Deploy the contract if known address is empty.
		address, _, contract, err = deploy(tops, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to deploy contract")
			return
		}
		reflectedContractName := reflect.TypeOf(contract).Elem().Name()
		log.Debug().Str("name", reflectedContractName).Interface("address", address).Msg("Contract deployed")
	} else {
		// Otherwise, instantiate the contract.
		address = knownAddress
		contract, err = instantiate(address, c)
		if err != nil {
			log.Error().Err(err).Msg("Unable to instantiate contract")
			return
		}
		reflectedContractName := reflect.TypeOf(contract).Elem().Name()
		log.Debug().Str("name", reflectedContractName).Msg("Contract instantiated")
	}

	// Check that the contract can be called.
	err = blockUntilSuccessful(ctx, c, func() error {
		log.Trace().Msg("Contract is not available yet")
		return call(contract)
	})
	return
}

// Ensure the UniswapV3Factory fee tier is enabled, activating it if it hasn't been enabled already.
func enableFeeAmount(contract *v3core.UniswapV3Factory, tops *bind.TransactOpts, cops *bind.CallOpts, fee, tickSpacing int64) error {
	// Check the current tick spacing for this fee amount.
	currentTickSpacing, err := contract.FeeAmountTickSpacing(cops, big.NewInt(fee))
	if err != nil {
		return err
	}

	// Enable the fee amount if needed.
	newTickSpacing := big.NewInt(tickSpacing)
	if currentTickSpacing.Cmp(newTickSpacing) == 0 {
		log.Debug().Msg("Fee amount already enabled")
	} else {
		_, err = contract.EnableFeeAmount(tops, big.NewInt(fee), big.NewInt(tickSpacing))
		if err != nil {
			log.Error().Err(err).Msg("Unable to enable fee amount")
			return err
		}
		log.Debug().Msg("Fee amount enabled")
	}
	return nil
}

// Transfer UniswapV3Factory ownership to a new address.
func transferUniswapV3FactoryOwnership(contract *v3core.UniswapV3Factory, tops *bind.TransactOpts, cops *bind.CallOpts, newOwner common.Address) error {
	currentOwner, err := contract.Owner(cops)
	if err != nil {
		return err
	}
	if currentOwner == newOwner {
		log.Debug().Msg("Factory contract already owned by this address")
	} else {
		_, err = contract.SetOwner(tops, newOwner)
		if err != nil {
			log.Error().Err(err).Msg("Unable to set a new owner for the Factory contract")
			return err
		}
		log.Debug().Msg("Set new owner for Factory contract")
	}
	return nil
}

// Transfer ProxyAdmin ownership to a new address.
func transferProxyAdminOwnership(contract *openzeppelin.ProxyAdmin, tops *bind.TransactOpts, cops *bind.CallOpts, newOwner common.Address) error {
	currentOwner, err := contract.Owner(cops)
	if err != nil {
		return err
	}
	if currentOwner == newOwner {
		log.Debug().Msg("ProxyAdmin contract already owned by this address")
	} else {
		_, err = contract.TransferOwnership(tops, newOwner)
		if err != nil {
			log.Error().Err(err).Msg("Unable to transfer ProxyAdmin ownership")
			return err
		}
		log.Debug().Msg("Transfer ProxyAdmin ownership")
	}
	return nil
}

// Return contracts addresses from the UniswapV3 configuration.
func (c *UniswapV3Config) GetAddresses() UniswapV3Addresses {
	return UniswapV3Addresses{
		FactoryV3:                          c.FactoryV3.Address,
		Multicall:                          c.Multicall.Address,
		ProxyAdmin:                         c.ProxyAdmin.Address,
		TickLens:                           c.TickLens.Address,
		NFTDescriptorLib:                   c.NFTDescriptorLib.Address,
		NonfungibleTokenPositionDescriptor: c.NonfungibleTokenPositionDescriptor.Address,
		TransparentUpgradeableProxy:        c.TransparentUpgradeableProxy.Address,
		NonfungiblePositionManager:         c.NonfungiblePositionManager.Address,
		Migrator:                           c.Migrator.Address,
		Staker:                             c.Staker.Address,
		QuoterV2:                           c.QuoterV2.Address,
		SwapRouter02:                       c.SwapRouter02.Address,
		WETH9:                              c.WETH9.Address,
	}
}
