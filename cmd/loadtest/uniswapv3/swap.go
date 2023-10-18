package uniswapv3loadtest

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/polygon-cli/contracts/uniswapv3"
	"github.com/rs/zerolog/log"
)

var (
	// The amount of inbound token given as swap input.
	swapAmountIn = big.NewInt(1000)
	// The minimum amount of outbound token received as swap output.
	swapAmountOutMinimum = big.NewInt(996)
)

// ExactInputSingleSwap performs a UniswapV3 swap using the `ExactInputSingle` method which swaps a fixed amount of
// one token for a maximum possible amount of another token. The direction of the swap is determined
// by the nonce value.
func ExactInputSingleSwap(tops *bind.TransactOpts, swapRouter *uniswapv3.SwapRouter02, poolConfig PoolConfig, recipient common.Address, nonce uint64) error {
	// Determine the direction of the swap.
	tIn := poolConfig.Token0
	tInName := "token0"
	tOut := poolConfig.Token1
	tOutName := "token1"
	if nonce%2 == 0 {
		tIn = poolConfig.Token1
		tInName = "token1"
		tOut = poolConfig.Token0
		tOutName = "token0"
	}

	// Perform swap.
	tx, err := swapRouter.ExactInputSingle(tops, uniswapv3.IV3SwapRouterExactInputSingleParams{
		// The contract address of the inbound token.
		TokenIn: tIn.Address,
		// The contract address of the outbound token.
		TokenOut: tOut.Address,
		// The fee tier of the pool, used to determine the correct pool contract in which to execute the swap.
		Fee: poolConfig.Fees,
		// The destination address of the outbound token.
		Recipient: recipient,
		// The amount of inbound token given as swap input.
		AmountIn: swapAmountIn,
		// The minimum amount of outbound token received as swap output.
		AmountOutMinimum: swapAmountOutMinimum,
		// The limit for the price swap.
		// Note: we set this to zero which makes the parameter inactive. In production, this value can
		// be used to protect against price impact.
		SqrtPriceLimitX96: big.NewInt(0),
	})
	if err != nil {
		log.Error().Err(err).Str("tokenIn", tInName).Str("tokenOut", tOutName).Msg("Unable to swap")
		return err
	}
	log.Debug().Str("tokenIn", tInName).Str("tokenOut", tOutName).Msg("Successful swap")
	log.Trace().Interface("hash", tx.Hash()).Msg("Transaction")
	return nil
}
