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
	swapAmountIn = big.NewInt(1_000)

	// The minimum amount of outbound token received as swap output.
	// It is set to a percentage of `swapAmountIn` (regardless of pool fees).
	swapAmountOutMinimum = new(big.Int).Mul(swapAmountIn, new(big.Int).Div(big.NewInt(98), big.NewInt(100)))
)

// ExactInputSingleSwap performs a UniswapV3 swap using the `ExactInputSingle` method which swaps a fixed amount of
// one token for a maximum possible amount of another token. The direction of the swap is determined
// by the nonce value.
func ExactInputSingleSwap(tops *bind.TransactOpts, swapRouter *uniswapv3.SwapRouter02, poolConfig PoolConfig, recipient common.Address, nonce uint64) error {
	// Determine the direction of the swap.
	swapDirection := getSwapDirection(nonce, poolConfig)

	// Perform swap.
	tx, err := swapRouter.ExactInputSingle(tops, uniswapv3.IV3SwapRouterExactInputSingleParams{
		// The contract address of the inbound token.
		TokenIn: swapDirection.tokenIn,
		// The contract address of the outbound token.
		TokenOut: swapDirection.tokenOut,
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
		log.Error().Err(err).Str("tokenIn", swapDirection.tokenInName).Str("tokenOut", swapDirection.tokenOutName).Msg("Unable to swap")
		return err
	}
	log.Debug().Str("tokenIn", swapDirection.tokenInName).Str("tokenOut", swapDirection.tokenOutName).Msg("Successful swap")
	log.Trace().Interface("hash", tx.Hash()).Msg("Transaction")
	return nil
}

// swapDirection represents a swap direction with the inbound and outbound tokens.
type swapDirection struct {
	tokenIn, tokenOut         common.Address
	tokenInName, tokenOutName string
}

// Return the direction of the swap given the nonce value.
func getSwapDirection(nonce uint64, poolConfig PoolConfig) swapDirection {
	if nonce%2 == 0 {
		return swapDirection{
			tokenIn:     poolConfig.Token0.Address,
			tokenInName: "token0",

			tokenOut:     poolConfig.Token1.Address,
			tokenOutName: "token1",
		}
	}
	return swapDirection{
		tokenIn:     poolConfig.Token1.Address,
		tokenInName: "token1",

		tokenOut:     poolConfig.Token0.Address,
		tokenOutName: "token0",
	}
}
