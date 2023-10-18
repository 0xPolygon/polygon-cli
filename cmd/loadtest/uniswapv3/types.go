package uniswapv3loadtest

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/polygon-cli/contracts/uniswapv3"
)

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
