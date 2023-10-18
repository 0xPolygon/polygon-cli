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

func DeploySwapperContract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts, config UniswapV3Config, tokenName, tokenSymbol string, amount *big.Int, recipient common.Address, tokenKnownAddress common.Address, blockUntilSuccessful blockUntilSuccessfulFn) (ContractConfig[uniswapv3.Swapper], error) {
	var token ContractConfig[uniswapv3.Swapper]
	var err error
	addressesToApprove := []common.Address{config.NonfungiblePositionManager.Address, config.SwapRouter02.Address}
	token.Address, token.Contract, err = DeployOrInstantiateContract(
		ctx, c, tops, cops, tokenName, tokenKnownAddress,
		func(*bind.TransactOpts, bind.ContractBackend) (common.Address, *types.Transaction, *uniswapv3.Swapper, error) {
			return uniswapv3.DeploySwapper(tops, c)
		},
		uniswapv3.NewSwapper,
		func(contract *uniswapv3.Swapper) error {
			return approveSwapperSpendingsByUniswap(ctx, contract, tops, cops, addressesToApprove, amount, recipient)
		},
		blockUntilSuccessful,
	)
	if err != nil {
		return token, err
	}
	return token, nil
}

func approveSwapperSpendingsByUniswap(ctx context.Context, contract *uniswapv3.Swapper, tops *bind.TransactOpts, cops *bind.CallOpts, addresses []common.Address, amount *big.Int, owner common.Address) error {
	name, err := contract.Name(cops)
	if err != nil {
		return err

	}
	for _, address := range addresses {
		tx, err := contract.Approve(tops, address, amount)
		if err != nil {
			log.Error().Err(err).Interface("address", address).Msg("Unable to approve spendings")
			return err
		}

		backoff.Retry(func() error {
			allowance, err := contract.Allowance(cops, owner, address)
			if err != nil {
				return err
			}
			zero := big.NewInt(0)
			if allowance.Cmp(zero) == 0 {
				return fmt.Errorf("allowance is zero")
			}
			return nil
		}, backoff.NewConstantBackOff(time.Second*2))

		log.Debug().Str("Swapper", name).Str("spender", address.String()).Interface("amount", amount).Msg("Spending approved")
		log.Trace().Interface("hash", tx.Hash()).Msg("Transaction")
	}
	return nil
}
