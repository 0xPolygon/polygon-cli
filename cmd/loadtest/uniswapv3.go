package loadtest

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maticnetwork/polygon-cli/contracts/uniswapv3"
	"github.com/rs/zerolog/log"
)

func deployUniswapV3Factory(c *ethclient.Client, tops *bind.TransactOpts) (ethcommon.Address, *uniswapv3.Uniswapv3, error) {
	// Deploy the UniswapV3Factory contract.
	address, _, _, err := uniswapv3.DeployUniswapv3(tops, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to deploy UniswapV3Factory contract")
		return ethcommon.Address{}, nil, err
	}
	log.Trace().Interface("address", address).Msg("UniswapV3Factory contract address")

	// Create a new instance of the contract.
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
