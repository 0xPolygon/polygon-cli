package uniswapv3loadtest

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maticnetwork/polygon-cli/contracts/uniswapv3"
	"github.com/rs/zerolog/log"
)

const (
	// The fee amount to enable for one basic point.
	// https://github.com/Uniswap/deploy-v3/blob/b7aac0f1c5353b36802dc0cf95c426d2ef0c3252/src/steps/add-1bp-fee-tier.ts#L5
	ONE_BP_FEE = 100

	// The spacing between ticks to be enforced for all pools with the given fee amount.
	// https://github.com/Uniswap/deploy-v3/blob/b7aac0f1c5353b36802dc0cf95c426d2ef0c3252/src/steps/add-1bp-fee-tier.ts#L6
	ONE_BP_TICK_SPACING = 1

	// The max amount of seconds into the future the incentive startTime can be set.
	// https://github.com/Uniswap/deploy-v3/blob/b7aac0f1c5353b36802dc0cf95c426d2ef0c3252/src/steps/deploy-v3-staker.ts#L11
	MAX_INCENTIVE_START_LEAD_TIME = ONE_MONTH_SECONDS

	// The max duration of an incentive in seconds.
	// https://github.com/Uniswap/deploy-v3/blob/b7aac0f1c5353b36802dc0cf95c426d2ef0c3252/src/steps/deploy-v3-staker.ts#L13
	MAX_INCENTIVE_DURATION = ONE_YEAR_SECONDS * 2

	// Time units.
	ONE_MINUTE_SECONDS = 60
	ONE_HOUR_SECONDS   = ONE_MINUTE_SECONDS * 60
	ONE_DAY_SECONDS    = ONE_HOUR_SECONDS * 24
	ONE_MONTH_SECONDS  = ONE_DAY_SECONDS * 30
	ONE_YEAR_SECONDS   = ONE_DAY_SECONDS * 365
)

var oldNFTPositionLibraryAddress = common.HexToAddress("0x73a6d49037afd585a0211a7bb4e990116025b45d")

type (
	// uniswapV3Config represents the whole UniswapV3 configuration.
	UniswapV3Config struct {
		FactoryV3                          ContractConfig[uniswapv3.UniswapV3Factory]
		Multicall                          ContractConfig[uniswapv3.UniswapInterfaceMulticall]
		ProxyAdmin                         ContractConfig[uniswapv3.ProxyAdmin]
		TickLens                           ContractConfig[uniswapv3.TickLens]
		NFTDescriptorLib                   ContractConfig[uniswapv3.NFTDescriptor]
		NonfungibleTokenPositionDescriptor ContractConfig[uniswapv3.NonfungibleTokenPositionDescriptor]
		TransparentUpgradeableProxy        ContractConfig[uniswapv3.TransparentUpgradeableProxy]
		NonfungiblePositionManager         ContractConfig[uniswapv3.NonfungiblePositionManager]
		Migrator                           ContractConfig[uniswapv3.V3Migrator]
		Staker                             ContractConfig[uniswapv3.UniswapV3Staker]
		QuoterV2                           ContractConfig[uniswapv3.QuoterV2]
		SwapRouter02                       ContractConfig[uniswapv3.SwapRouter02]

		WETH9 ContractConfig[uniswapv3.WETH9]
	}

	// Contract represents a UniswapV3 contract.
	Contract interface {
		uniswapv3.UniswapV3Factory | uniswapv3.UniswapInterfaceMulticall | uniswapv3.ProxyAdmin | uniswapv3.TickLens | uniswapv3.WETH9 | uniswapv3.NFTDescriptor | uniswapv3.NonfungibleTokenPositionDescriptor | uniswapv3.TransparentUpgradeableProxy | uniswapv3.NonfungiblePositionManager | uniswapv3.V3Migrator | uniswapv3.UniswapV3Staker | uniswapv3.QuoterV2 | uniswapv3.SwapRouter02 | uniswapv3.Swapper
	}

	// Address represents a UniswapV3 contract address (WETH9 also included).
	Address struct {
		FactoryV3, Multicall, ProxyAdmin, TickLens, NFTDescriptorLib, NonfungibleTokenPositionDescriptor, TransparentUpgradeableProxy, NonfungiblePositionManager, Migrator, Staker, QuoterV2, SwapRouter02 common.Address
		WETH9                                                                                                                                                                                               common.Address
	}

	// contractConfig represents a specific contract configuration.
	ContractConfig[T Contract] struct {
		Address  common.Address
		Contract *T
	}
)

func (c *UniswapV3Config) ToAddresses() Address {
	return Address{
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

// Deploy the full UniswapV3 contract suite in 15 different steps.
// Source: https://github.com/Uniswap/deploy-v3
func DeployUniswapV3(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts, knownAddresses Address, ownerAddress common.Address, blockblockUntilSuccessful blockUntilSuccessfulFn) (UniswapV3Config, error) {
	config := UniswapV3Config{}
	var err error

	// 0. Deploy WETH9.
	config.WETH9.Address, config.WETH9.Contract, err = DeployOrInstantiateContract(
		ctx, c, tops, cops, "Step 0: Contract WETH9 deployment", knownAddresses.WETH9,
		uniswapv3.DeployWETH9,
		uniswapv3.NewWETH9,
		func(contract *uniswapv3.WETH9) (err error) {
			_, err = contract.BalanceOf(cops, common.Address{})
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return config, err
	}

	// 1. Deploy UniswapV3Factory.
	config.FactoryV3.Address, config.FactoryV3.Contract, err = DeployOrInstantiateContract(
		ctx, c, tops, cops, "Step 1: Contract UniswapV3Factory deployment", knownAddresses.FactoryV3,
		uniswapv3.DeployUniswapV3Factory,
		uniswapv3.NewUniswapV3Factory,
		func(contract *uniswapv3.UniswapV3Factory) (err error) {
			_, err = contract.Owner(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return config, err
	}

	// 2. Enable one basic point fee tier.
	if err = enableOneBPFeeTier(config.FactoryV3.Contract, tops, cops, ONE_BP_FEE, ONE_BP_TICK_SPACING); err != nil {
		return config, err
	}

	// 3. Deploy UniswapInterfaceMulticall.
	config.Multicall.Address, config.Multicall.Contract, err = DeployOrInstantiateContract(
		ctx, c, tops, cops, "Step 3: Contract UniswapInterfaceMulticall deployment", knownAddresses.Multicall,
		uniswapv3.DeployUniswapInterfaceMulticall,
		uniswapv3.NewUniswapInterfaceMulticall,
		func(contract *uniswapv3.UniswapInterfaceMulticall) (err error) {
			_, err = contract.GetEthBalance(cops, common.Address{})
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return config, err
	}

	// 4. Deploy ProxyAdmin.
	config.ProxyAdmin.Address, config.ProxyAdmin.Contract, err = DeployOrInstantiateContract(
		ctx, c, tops, cops, "Step 4: Contract ProxyAdmin deployment", knownAddresses.ProxyAdmin,
		uniswapv3.DeployProxyAdmin,
		uniswapv3.NewProxyAdmin,
		func(contract *uniswapv3.ProxyAdmin) (err error) {
			_, err = contract.Owner(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return config, err
	}

	// 5. Deploy TickLens.
	config.TickLens.Address, config.TickLens.Contract, err = DeployOrInstantiateContract(
		ctx, c, tops, cops, "Step 5: Contract TickLens deployment", knownAddresses.TickLens,
		uniswapv3.DeployTickLens,
		uniswapv3.NewTickLens,
		func(contract *uniswapv3.TickLens) (err error) {
			// The only function we can call to check the contract is deployed is `GetPopulatedTicksInWord`.
			// Unfortunately, such call will revert because no pools are deployed yet.
			// That's why we only return a nil value here.
			return nil
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return config, err
	}

	// 6. Deploy NFTDescriptor library.
	config.NFTDescriptorLib.Address, _, err = DeployOrInstantiateContract(
		ctx, c, tops, cops, "Step 6: Library NFTDescriptor deployment", knownAddresses.NFTDescriptorLib,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *uniswapv3.NFTDescriptor, error) {
			return uniswapv3.DeployNFTDescriptor(tops, c)
		},
		uniswapv3.NewNFTDescriptor,
		func(contract *uniswapv3.NFTDescriptor) (err error) {
			// No methods to call to check if the library has been deployed.
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return config, err
	}

	// 7. Deploy NonfungibleTokenPositionDescriptor.
	config.NonfungibleTokenPositionDescriptor.Address, config.NonfungibleTokenPositionDescriptor.Contract, err = DeployOrInstantiateContract(
		ctx, c, tops, cops, "Step 7: Contract NonfungibleTokenPositionDescriptor deployment", knownAddresses.NonfungibleTokenPositionDescriptor,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *uniswapv3.NonfungibleTokenPositionDescriptor, error) {
			oldAddressFmt := "__$cea9be979eee3d87fb124d6cbb244bb0b5$__"
			newAddressFmt := strings.TrimPrefix(strings.ToLower(config.NFTDescriptorLib.Address.String()), "0x")
			newNonfungibleTokenPositionDescriptorBytecode := strings.ReplaceAll(uniswapv3.NonfungibleTokenPositionDescriptorMetaData.Bin, oldAddressFmt, newAddressFmt)
			if uniswapv3.NonfungibleTokenPositionDescriptorMetaData.Bin == newNonfungibleTokenPositionDescriptorBytecode {
				err = fmt.Errorf("the NonfungibleTokenPositionDescriptor bytecode has not been updated")
				log.Error().Err(err).Msg("NonfungibleTokenPositionDescriptor bytecode has not been updated")
				return common.Address{}, nil, nil, err
			}
			log.Debug().Interface("oldAddress", oldNFTPositionLibraryAddress).Interface("newAddress", config.NFTDescriptorLib.Address).Msg("NonfungibleTokenPositionDescriptor bytecode updated with the new NFTDescriptor library address")

			// Deploy contract.
			var nativeCurrencyLabelBytes [32]byte
			copy(nativeCurrencyLabelBytes[:], "ETH")
			uniswapv3.NonfungibleTokenPositionDescriptorMetaData.Bin = newNonfungibleTokenPositionDescriptorBytecode
			uniswapv3.NonfungibleTokenPositionDescriptorBin = newNonfungibleTokenPositionDescriptorBytecode
			return uniswapv3.DeployNonfungibleTokenPositionDescriptor(tops, c, config.WETH9.Address, nativeCurrencyLabelBytes)
		},
		uniswapv3.NewNonfungibleTokenPositionDescriptor,
		func(contract *uniswapv3.NonfungibleTokenPositionDescriptor) (err error) {
			_, err = contract.WETH9(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return config, err
	}

	// 8. Deploy TransparentUpgradeableProxy.
	config.TransparentUpgradeableProxy.Address, config.TransparentUpgradeableProxy.Contract, err = DeployOrInstantiateContract(
		ctx, c, tops, cops, "Step 8: Contract TransparentUpgradeableProxy deployment", knownAddresses.TransparentUpgradeableProxy,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *uniswapv3.TransparentUpgradeableProxy, error) {
			return uniswapv3.DeployTransparentUpgradeableProxy(tops, c, config.NonfungibleTokenPositionDescriptor.Address, config.ProxyAdmin.Address, []byte(""))
		},
		uniswapv3.NewTransparentUpgradeableProxy,
		func(contract *uniswapv3.TransparentUpgradeableProxy) (err error) {
			// The TransparentUpgradeableProxy contract methods can only be called by the admin.
			// This is not a problem when we first deploy the contract because the deployer is set to be
			// the admin by default. Thus, we can call any method of the contract to check it has been deployed.
			// But when we use pre-deployed contracts, since the TransparentUpgradeableProxy ownership
			// has been transferred, we get "execution reverted" errors when trying to call any method.
			// That's why we don't call any method in the pre-deployed contract mode.
			//if knownAddresses.TransparentUpgradeableProxy == (common.Address{}) {
			//	_, err = contract.Admin(tops)
			//}
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return config, err
	}

	// 9. Deploy NonfungiblePositionManager.
	config.NonfungiblePositionManager.Address, config.NonfungiblePositionManager.Contract, err = DeployOrInstantiateContract(
		ctx, c, tops, cops, "Step 9: Contract NonfungiblePositionManager deployment", knownAddresses.NonfungiblePositionManager,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *uniswapv3.NonfungiblePositionManager, error) {
			return uniswapv3.DeployNonfungiblePositionManager(tops, c, config.FactoryV3.Address, config.WETH9.Address, config.TransparentUpgradeableProxy.Address)
		},
		uniswapv3.NewNonfungiblePositionManager,
		func(contract *uniswapv3.NonfungiblePositionManager) (err error) {
			_, err = contract.BaseURI(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return config, err
	}

	// 10. Deploy Migrator.
	config.Migrator.Address, config.Migrator.Contract, err = DeployOrInstantiateContract(
		ctx, c, tops, cops, "Step 10: Contract V3Migrator deployment", knownAddresses.Migrator,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *uniswapv3.V3Migrator, error) {
			return uniswapv3.DeployV3Migrator(tops, c, config.FactoryV3.Address, config.WETH9.Address, config.NonfungiblePositionManager.Address)
		},
		uniswapv3.NewV3Migrator,
		func(contract *uniswapv3.V3Migrator) (err error) {
			_, err = contract.WETH9(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return config, err
	}

	// 11. Set Factory owner.
	if err = setFactoryOwner(config.FactoryV3.Contract, tops, cops, ownerAddress); err != nil {
		return config, err
	}

	// 12. Deploy Staker.
	config.Staker.Address, config.Staker.Contract, err = DeployOrInstantiateContract(
		ctx, c, tops, cops, "Step 12: Contract UniswapV3Staker deployment", knownAddresses.Staker,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *uniswapv3.UniswapV3Staker, error) {
			return uniswapv3.DeployUniswapV3Staker(tops, c, config.FactoryV3.Address, config.NonfungiblePositionManager.Address, big.NewInt(MAX_INCENTIVE_START_LEAD_TIME), big.NewInt(MAX_INCENTIVE_DURATION))
		},
		uniswapv3.NewUniswapV3Staker,
		func(contract *uniswapv3.UniswapV3Staker) (err error) {
			_, err = contract.Factory(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return config, err
	}

	// 13. Deploy QuoterV2.
	config.QuoterV2.Address, config.QuoterV2.Contract, err = DeployOrInstantiateContract(
		ctx, c, tops, cops, "Step 13: Contract QuoterV2 deployment", knownAddresses.QuoterV2,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *uniswapv3.QuoterV2, error) {
			return uniswapv3.DeployQuoterV2(tops, c, config.FactoryV3.Address, config.WETH9.Address)
		},
		uniswapv3.NewQuoterV2,
		func(contract *uniswapv3.QuoterV2) (err error) {
			_, err = contract.Factory(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return config, err
	}

	// 14. Deploy SwapRouter02.
	config.SwapRouter02.Address, config.SwapRouter02.Contract, err = DeployOrInstantiateContract(
		ctx, c, tops, cops, "Step 14: Contract SwapRouter02 deployment", knownAddresses.SwapRouter02,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *uniswapv3.SwapRouter02, error) {
			// Note: we specify an empty address for UniswapV2Factory.
			uniswapFactoryV2Address := common.Address{}
			return uniswapv3.DeploySwapRouter02(tops, c, uniswapFactoryV2Address, config.FactoryV3.Address, config.NonfungiblePositionManager.Address, config.WETH9.Address)
		},
		uniswapv3.NewSwapRouter02,
		func(contract *uniswapv3.SwapRouter02) (err error) {
			_, err = contract.Factory(cops)
			return
		},
		blockblockUntilSuccessful,
	)
	if err != nil {
		return config, err
	}

	// 15. Transfer ProxyAdmin ownership.
	if err = transferProxyAdminOwnership(config.ProxyAdmin.Contract, tops, cops, ownerAddress); err != nil {
		return config, err
	}

	return config, nil
}

// Deploy or instantiate any UniswapV3 contract.
// This method will either deploy a contract if the known address is empty (equal to `common.Address{}` or `0x0“)
// or instantiate the contract if the known address is specified.
func DeployOrInstantiateContract[T Contract](
	ctx context.Context,
	c *ethclient.Client,
	tops *bind.TransactOpts,
	cops *bind.CallOpts,
	logMessage string,
	knownAddress common.Address,
	deployFn func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *T, error),
	instantiateFn func(common.Address, bind.ContractBackend) (*T, error),
	callFn func(*T) error,
	blockUntilSuccessful blockUntilSuccessfulFn,
) (address common.Address, contract *T, err error) {
	if knownAddress == (common.Address{}) {
		// Deploy the contract if known address is empty.
		var tx *types.Transaction
		address, tx, contract, err = deployFn(tops, c)
		if err != nil {
			log.Error().Err(err).Str("logMessage", logMessage).Msg("Unable to deploy contract")
			return
		}
		reflectedContractName := reflect.TypeOf(contract).Elem().Name()
		log.Debug().Str("name", reflectedContractName).Str("logMessage", logMessage).Str("address", address.String()).Msg("Contract deployed")
		log.Trace().Str("name", reflectedContractName).Str("logMessage", logMessage).Str("hash", tx.Hash().String()).Msg("Transaction")
	} else {
		// Otherwise, instantiate the contract.
		address = knownAddress
		contract, err = instantiateFn(address, c)
		if err != nil {
			log.Error().Err(err).Str("logMessage", logMessage).Msg("Unable to instantiate contract")
			return
		}
		reflectedContractName := reflect.TypeOf(contract).Elem().Name()
		log.Debug().Str("name", reflectedContractName).Str("logMessage", logMessage).Msg("Contract instantiated")
	}

	// Check that the contract is deployed and ready.
	if err = blockUntilSuccessful(ctx, c, func() error {
		log.Trace().Str("logMessage", logMessage).Msg("Contract is not available yet")
		return callFn(contract)
	}); err != nil {
		return
	}
	return
}

func enableOneBPFeeTier(contract *uniswapv3.UniswapV3Factory, tops *bind.TransactOpts, cops *bind.CallOpts, fee, tickSpacing int64) error {
	// Check the current tick spacing for this fee amount.
	currentTickSpacing, err := contract.FeeAmountTickSpacing(cops, big.NewInt(fee))
	if err != nil {
		return err
	}

	newTickSpacing := big.NewInt(tickSpacing)
	if currentTickSpacing.Cmp(newTickSpacing) == 0 {
		// If those are the same, it means it has already been enabled.
		log.Debug().Msg("One basic point fee tier already enabled")
	} else {
		// If those are not the same, it means it should be enabled.
		tx, err := contract.EnableFeeAmount(tops, big.NewInt(fee), big.NewInt(tickSpacing))
		if err != nil {
			log.Error().Err(err).Msg("Unable to enable one basic point fee tier")
			return err
		}
		log.Debug().Msg("Enable one basic point fee tier")
		log.Trace().Interface("hash", tx.Hash()).Msg("Transaction")
	}
	return nil
}

func setFactoryOwner(contract *uniswapv3.UniswapV3Factory, tops *bind.TransactOpts, cops *bind.CallOpts, newOwner common.Address) error {
	currentOwner, err := contract.Owner(cops)
	if err != nil {
		return err
	}
	if currentOwner == newOwner {
		log.Debug().Msg("Factory contract already owned by this address")
	} else {
		tx, err := contract.SetOwner(tops, newOwner)
		if err != nil {
			log.Error().Err(err).Msg("Unable to set a new owner for the Factory contract")
			return err
		}
		log.Debug().Msg("Set new owner for Factory contract")
		log.Trace().Interface("hash", tx.Hash()).Msg("Transaction")
	}
	return nil
}

func transferProxyAdminOwnership(contract *uniswapv3.ProxyAdmin, tops *bind.TransactOpts, cops *bind.CallOpts, newOwner common.Address) error {
	currentOwner, err := contract.Owner(cops)
	if err != nil {
		return err
	}
	if currentOwner == newOwner {
		log.Debug().Msg("ProxyAdmin contract already owned by this address")
	} else {
		tx, err := contract.TransferOwnership(tops, newOwner)
		if err != nil {
			log.Error().Err(err).Msg("Unable to transfer ProxyAdmin ownership")
			return err
		}
		log.Debug().Msg("Transfer ProxyAdmin ownership")
		log.Trace().Interface("hash", tx.Hash()).Msg("Transaction")
	}
	return nil
}
