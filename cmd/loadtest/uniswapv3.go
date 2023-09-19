package loadtest

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maticnetwork/polygon-cli/contracts/uniswapv3"
	"github.com/rs/zerolog/log"
)

type UniswapV3Config struct {
	Factory struct {
		Address  ethcommon.Address
		Contract *uniswapv3.Uniswapv3
	}
}

func deployUniswapV3(c *ethclient.Client, tops *bind.TransactOpts) (UniswapV3Config, error) {
	var config UniswapV3Config
	var err error

	// 1. Deploy UniswapV3Factory.
	config.Factory.Address, config.Factory.Contract, err = deployUniswapV3Factory(c, tops)
	if err != nil {
		return UniswapV3Config{}, err
	}

	return config, nil
}

func deployUniswapV3Factory(c *ethclient.Client, tops *bind.TransactOpts) (ethcommon.Address, *uniswapv3.Uniswapv3, error) {
	address, _, _, err := uniswapv3.DeployUniswapv3(tops, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to deploy UniswapV3Factory contract")
		return ethcommon.Address{}, nil, err
	}
	log.Trace().Interface("address", address).Msg("UniswapV3Factory contract address")

	contract, err := uniswapv3.NewUniswapv3(address, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to instantiate UniswapV3Factory contract")
		return ethcommon.Address{}, nil, err
	}
	return address, contract, nil
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
