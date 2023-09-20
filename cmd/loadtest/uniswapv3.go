package loadtest

import (
	"context"
	"fmt"
	"math/big"

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
	ONE_BP_FEE int64 = 100

	// The spacing between ticks to be enforced for all pools with the given fee amount.
	// https://github.com/Uniswap/deploy-v3/blob/b7aac0f1c5353b36802dc0cf95c426d2ef0c3252/src/steps/add-1bp-fee-tier.ts#L6
	ONE_BP_TICK_SPACING int64 = 1
)

type UniswapV3Addresses struct {
	Factory, Multicall, ProxyAdmin, TickLens, NFTDescriptionLib common.Address
}

type UniswapV3Config struct {
	Factory           uniswapV3ContractDeployment[uniswapv3.UniswapV3Factory]
	Multicall         uniswapV3ContractDeployment[uniswapv3.UniswapInterfaceMulticall]
	ProxyAdmin        uniswapV3ContractDeployment[uniswapv3.ProxyAdmin]
	TickLens          uniswapV3ContractDeployment[uniswapv3.TickLens]
	NFTDescriptionLib uniswapV3ContractDeployment[uniswapv3.NFTDescriptor]
}

type uniswapV3ContractDeployment[T uniswapV3Contract] struct {
	Address  common.Address
	Contract *T
}

type uniswapV3Contract interface {
	uniswapv3.UniswapV3Factory | uniswapv3.UniswapInterfaceMulticall | uniswapv3.ProxyAdmin | uniswapv3.TickLens | uniswapv3.NFTDescriptor
}

func deployUniswapV3(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts, knownAddresses UniswapV3Addresses) (UniswapV3Config, error) {
	var config UniswapV3Config
	var err error

	// 1. Deploy UniswapV3Factory.
	config.Factory.Address, config.Factory.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops, "Factory", knownAddresses.Factory,
		uniswapv3.DeployUniswapV3Factory,
		uniswapv3.NewUniswapV3Factory,
		func(contract *uniswapv3.UniswapV3Factory) (err error) {
			_, err = contract.Owner(cops)
			return
		},
	)
	if err != nil {
		return UniswapV3Config{}, err
	}

	// 2. Enable one basic point fee tier.
	err = enableOneBPFeeTier(config.Factory.Contract, tops, ONE_BP_FEE, ONE_BP_TICK_SPACING)
	if err != nil {
		return UniswapV3Config{}, err
	}

	// 3. Deploy UniswapInterfaceMulticall.
	config.Multicall.Address, config.Multicall.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops, "Multicall", knownAddresses.Multicall,
		uniswapv3.DeployUniswapInterfaceMulticall,
		uniswapv3.NewUniswapInterfaceMulticall,
		func(contract *uniswapv3.UniswapInterfaceMulticall) (err error) {
			_, err = contract.GetEthBalance(cops, common.Address{})
			return
		},
	)
	if err != nil {
		return UniswapV3Config{}, err
	}

	// 4. Deploy ProxyAdmin.
	config.ProxyAdmin.Address, config.ProxyAdmin.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops, "ProxyAdmin", knownAddresses.ProxyAdmin,
		uniswapv3.DeployProxyAdmin,
		uniswapv3.NewProxyAdmin,
		func(contract *uniswapv3.ProxyAdmin) (err error) {
			_, err = contract.Owner(cops)
			return
		},
	)
	if err != nil {
		return UniswapV3Config{}, err
	}

	// 5. Deploy TickLens.
	config.TickLens.Address, config.TickLens.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops, "TickLens", knownAddresses.TickLens,
		uniswapv3.DeployTickLens,
		uniswapv3.NewTickLens,
		func(contract *uniswapv3.TickLens) (err error) {
			// This call will revert because no ticks are populated yet.
			_, err = contract.GetPopulatedTicksInWord(cops, common.Address{}, int16(1))
			// TODO: Compare with error instead of string.
			if err.Error() == "execution reverted" {
				err = nil
			}
			return
		},
	)
	if err != nil {
		return UniswapV3Config{}, err
	}

	// 6. Deploy NFTDescriptionLib.
	config.NFTDescriptionLib.Address, config.NFTDescriptionLib.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops, "NFTDescriptionLib", knownAddresses.NFTDescriptionLib,
		uniswapv3.DeployNFTDescriptor,
		uniswapv3.NewNFTDescriptor,
		func(contract *uniswapv3.NFTDescriptor) (err error) {
			// FIXME: This call will cause a panic.
			//_, err = contract.ConstructTokenURI(cops, uniswapv3.NFTDescriptorConstructTokenURIParams{})
			return
		},
	)
	if err != nil {
		return UniswapV3Config{}, err
	}

	return config, nil
}

// Deploy or instantiate any UniswapV3 contract.
// This method will either deploy a contract if the known address is empty (equal to `common.Address{}` or `0x0“)
// or instantiate the contract if the known address is specified.
func deployOrInstantiateContract[T uniswapV3Contract](
	ctx context.Context,
	c *ethclient.Client,
	tops *bind.TransactOpts,
	cops *bind.CallOpts,
	name string,
	knownAddress common.Address,
	deployFunc func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *T, error),
	instantiateFunc func(common.Address, bind.ContractBackend) (*T, error),
	callFunc func(*T) error,
) (address common.Address, contract *T, err error) {
	if knownAddress == (common.Address{}) {
		// Deploy the contract if known address is empty.
		address, _, contract, err = deployFunc(tops, c)
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Unable to deploy %s contract", name))
			return
		}
		log.Trace().Interface("address", address).Msg(fmt.Sprintf("%s contract deployed", name))
	} else {
		// Otherwise, instantiate the contract.
		address = knownAddress
		contract, err = instantiateFunc(address, c)
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Unable to instantiate %s contract", name))
			return
		}
		log.Trace().Msg(fmt.Sprintf("%s contract instantiated", name))
	}

	// Check that the contract is deployed and ready.
	if err = blockUntilSuccessful(ctx, c, func() error { return callFunc(contract) }); err != nil {
		return
	}
	return
}

func enableOneBPFeeTier(contract *uniswapv3.UniswapV3Factory, tops *bind.TransactOpts, fee, tickSpacing int64) error {
	if _, err := contract.EnableFeeAmount(tops, big.NewInt(fee), big.NewInt(tickSpacing)); err != nil {
		return err
	}
	log.Trace().Msg("Enable a one basic point fee tier")
	return nil
}

// Create and initialise an ERC20 pool between two ERC20 contracts.
// Note that this will also deploy both ERC20 contracts.
func createPool() {
	// TODO
}

func swapTokenAForTokenB() {
	// TODO
}

func swapTokenBForTokenA() {
	// TODO
}
