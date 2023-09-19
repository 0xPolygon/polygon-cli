package loadtest

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

type UniswapV3Config struct {
	Factory    uniswapV3ContractDeployment[uniswapv3.UniswapV3Factory]
	Multicall  uniswapV3ContractDeployment[uniswapv3.UniswapInterfaceMulticall]
	ProxyAdmin uniswapV3ContractDeployment[uniswapv3.ProxyAdmin]
	TickLens   uniswapV3ContractDeployment[uniswapv3.TickLens]
}

type uniswapV3ContractDeployment[T uniswapV3Contract] struct {
	Address  common.Address
	Contract *T
}

type uniswapV3Contract interface {
	uniswapv3.UniswapV3Factory | uniswapv3.UniswapInterfaceMulticall | uniswapv3.ProxyAdmin | uniswapv3.TickLens
}

func deployUniswapV3(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts) (UniswapV3Config, error) {
	var config UniswapV3Config
	var err error

	config.Factory.Address, _, config.Factory.Contract, err = uniswapv3.DeployUniswapV3Factory(tops, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to deploy UniswapV3Factory contract")
		return UniswapV3Config{}, err
	}
	log.Trace().Interface("address", config.Factory.Address).Msg("UniswapV3Factory contract deployed")

	err = enableOneBPFeeTier(config.Factory.Contract, tops, ONE_BP_FEE, ONE_BP_TICK_SPACING)
	if err != nil {
		return UniswapV3Config{}, err
	}

	config.Multicall.Address, _, config.Multicall.Contract, err = uniswapv3.DeployUniswapInterfaceMulticall(tops, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to deploy UniswapInterfaceMulticall contract")
		return UniswapV3Config{}, err
	}
	log.Trace().Interface("address", config.Multicall.Address).Msg("UniswapInterfaceMulticall contract deployed")

	config.ProxyAdmin.Address, _, config.ProxyAdmin.Contract, err = uniswapv3.DeployProxyAdmin(tops, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to deploy ProxyAdmin contract")
		return UniswapV3Config{}, err
	}
	log.Trace().Interface("address", config.ProxyAdmin.Address).Msg("ProxyAdmin contract deployed")

	config.TickLens.Address, _, config.TickLens.Contract, err = uniswapv3.DeployTickLens(tops, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to deploy TickLens contract")
		return UniswapV3Config{}, err
	}
	log.Trace().Interface("address", config.TickLens.Address).Msg("TickLens contract deployed")

	return config, nil
}

// https://github.com/Uniswap/deploy-v3/blob/b7aac0f1c5353b36802dc0cf95c426d2ef0c3252/src/steps/add-1bp-fee-tier.ts
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
