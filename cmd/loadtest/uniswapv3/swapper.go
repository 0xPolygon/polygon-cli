package uniswapv3loadtest

import (
	"context"
	"errors"
	"fmt"
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
var allowanceAmount = big.NewInt(999_999_999_999_999_999)

// Deploy an ERC20 token.
func DeployERC20(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts, uniswapV3Config UniswapV3Config, tokenName string, recipient common.Address, tokenKnownAddress common.Address, blockUntilSuccessful blockUntilSuccessfulFn) (tokenConfig ContractConfig[uniswapv3.Swapper], err error) {
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
			uniswapV3Addresses := map[string]common.Address{
				"NFTPositionManager": uniswapV3Config.NonfungiblePositionManager.Address,
				"SwapRouter02":       uniswapV3Config.SwapRouter02.Address,
			}
			return setUniswapV3Allowances(ctx, c, contract, tops, cops, tokenName, uniswapV3Addresses, recipient, blockUntilSuccessful)
		},
		blockUntilSuccessful,
	)
	if err != nil {
		return
	}
	return
}

// Approve some UniswapV3 addresses to spend tokens on behalf of the token owner.
func setUniswapV3Allowances(ctx context.Context, c *ethclient.Client, contract *uniswapv3.Swapper, tops *bind.TransactOpts, cops *bind.CallOpts, tokenName string, addresses map[string]common.Address, owner common.Address, blockUntilSuccessful blockUntilSuccessfulFn) error {
	// Get the ERC20 contract name.
	erc20Name, err := contract.Name(cops)
	if err != nil {
		return err

	}

	for spenderName, spenderAddress := range addresses {
		// Approve the spender to spend the tokens on behalf of the owner.
		if _, err = contract.Approve(tops, spenderAddress, allowanceAmount); err != nil {
			log.Error().Err(err).
				Str("tokenName", fmt.Sprintf("%s_%s", erc20Name, tokenName)).
				Interface("spenderAddress", spenderAddress).Str("spenderName", spenderName).
				Interface("amount", allowanceAmount).
				Msg("Unable to set the allowance")
			return err
		}

		// Check that the allowance is set.
		err = blockUntilSuccessful(ctx, c, func() (err error) {
			allowance, err := contract.Allowance(cops, owner, spenderAddress)
			if err != nil {
				return err
			}
			if allowance.Cmp(big.NewInt(0)) == 0 { // allowance == 0
				return errors.New("allowance is set to zero")
			}
			if allowance.Cmp(allowanceAmount) == -1 { // allowance < allowanceAmount
				return errors.New("allowance has not been set properly")
			}
			return nil
		})
		if err != nil {
			log.Error().Err(err).
				Str("tokenName", fmt.Sprintf("%s_%s", erc20Name, tokenName)).
				Interface("spenderAddress", spenderAddress).Str("spenderName", spenderName).
				Interface("amount", allowanceAmount).
				Msg("Unable to verify that the allowance has been set")
			return err
		} else {
			log.Debug().
				Str("tokenName", fmt.Sprintf("%s_%s", erc20Name, tokenName)).
				Interface("spenderAddress", spenderAddress).Str("spenderName", spenderName).
				Interface("amount", allowanceAmount).
				Msg("Allowance set")
		}
	}
	return nil
}
