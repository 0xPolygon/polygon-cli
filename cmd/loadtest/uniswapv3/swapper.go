package uniswapv3loadtest

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maticnetwork/polygon-cli/contracts/uniswapv3"
	"github.com/rs/zerolog/log"
)

// The amount of token to approve a spender to use on behalf of the token owner.
// We use a very high amount to avoid frequent approval transactions. Don't use in production.
var approvalAmount = big.NewInt(999_999_999_999_999_999)

// Deploy an ERC20 token.
func DeployERC20(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts, uniswapV3Config UniswapV3Config, tokenName, tokenSymbol string, recipient common.Address, tokenKnownAddress common.Address, blockUntilSuccessful blockUntilSuccessfulFn) (tokenConfig ContractConfig[uniswapv3.Swapper], err error) {
	tokenConfig.Address, tokenConfig.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		tokenKnownAddress,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *uniswapv3.Swapper, error) {
			return uniswapv3.DeploySwapper(tops, c)
		},
		uniswapv3.NewSwapper,
		func(contract *uniswapv3.Swapper) error {
			// After the contract has been deployed, we autorize a few UniswapV3 addresses to spend those ERC20 tokens.
			// This is required to be able to perform swaps later.
			addressesToApprove := []common.Address{uniswapV3Config.NonfungiblePositionManager.Address, uniswapV3Config.SwapRouter02.Address}
			return approveSwapperSpendingsByUniswap(ctx, c, contract, tops, cops, addressesToApprove, recipient, blockUntilSuccessful)
		},
		blockUntilSuccessful,
	)
	if err != nil {
		return
	}
	return
}

// Approve a slice of addresses to spend tokens on behalf of the token owner.
func approveSwapperSpendingsByUniswap(ctx context.Context, c *ethclient.Client, contract *uniswapv3.Swapper, tops *bind.TransactOpts, cops *bind.CallOpts, addresses []common.Address, owner common.Address, blockUntilSuccessful blockUntilSuccessfulFn) error {
	// Get the ERC20 contract name.
	name, err := contract.Name(cops)
	if err != nil {
		return err

	}

	// Set allowances.
	for _, address := range addresses {
		// Approve the spender to spend the tokens on behalf of the owner.
		if _, err = contract.Approve(tops, address, approvalAmount); err != nil {
			log.Error().Err(err).Interface("address", address).Msg("Unable to set the allowance")
			return err
		}

		// Check that the allowance is set.
		if err = blockUntilSuccessful(ctx, c, func() (err error) {
			allowance, err := contract.Allowance(cops, owner, address)
			if err != nil {
				return err
			}
			if allowance.Cmp(approvalAmount) != 0 {
				return errors.New("allowance has not been set properly")
			}
			return nil
		}); err != nil {
			log.Error().Err(err).Msg("Unable to verify that the allowance has been set")
			return err
		}
		log.Debug().Str("name", name).Str("spender", address.String()).Interface("amount", approvalAmount).Msg("Allowance set")
	}
	return nil
}
