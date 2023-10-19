package uniswapv3loadtest

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/cenkalti/backoff"
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

func DeploySwapperContract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts, config UniswapV3Config, tokenName, tokenSymbol string, recipient common.Address, tokenKnownAddress common.Address, blockUntilSuccessful blockUntilSuccessfulFn) (ContractConfig[uniswapv3.Swapper], error) {
	var token ContractConfig[uniswapv3.Swapper]
	var err error
	addressesToApprove := []common.Address{config.NonfungiblePositionManager.Address, config.SwapRouter02.Address}
	token.Address, token.Contract, err = deployOrInstantiateContract(
		ctx, c, tops, cops,
		tokenKnownAddress,
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
		tx, err := contract.Approve(tops, address, approvalAmount)
		if err != nil {
			log.Error().Err(err).Interface("address", address).Msg("Unable to approve spendings")
			return err
		}

		if err := backoff.Retry(func() error {
			allowance, err := contract.Allowance(cops, owner, address)
			if err != nil {
				return err
			}
			zero := big.NewInt(0)
			if allowance.Cmp(zero) == 0 {
				return fmt.Errorf("allowance is zero")
			}
			return nil
		}, backoff.NewConstantBackOff(time.Second*2)); err != nil {
			return err
		}

		log.Debug().Str("Swapper", name).Str("spender", address.String()).Interface("amount", approvalAmount).Msg("Spending approved")
		log.Trace().Interface("hash", tx.Hash()).Msg("Transaction")
	}
	return nil
}
