package uniswapv3loadtest

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0xPolygon/polygon-cli/bindings/uniswapv3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

// The amount of inbound token given as swap input.
var SwapAmountInput = big.NewInt(1_000)

// ExactInputSingleSwap performs a UniswapV3 swap using the `ExactInputSingle` method which swaps a fixed amount of
// one token for a maximum possible amount of another token. The direction of the swap is determined
// by the nonce value.
func ExactInputSingleSwap(tops *bind.TransactOpts, swapRouter *uniswapv3.SwapRouter02, poolConfig PoolConfig, amountIn *big.Int, recipient common.Address, nonce uint64) (tx *types.Transaction, err error) {
	// Determine the direction of the swap.
	swapDirection := getSwapDirection(nonce, poolConfig)

	// Perform swap.
	slippageFactor := new(big.Float).SetFloat64(0.75)
	amountInFloat := new(big.Float).SetInt(amountIn)
	amountInFloat.Mul(amountInFloat, slippageFactor)
	amountOut := new(big.Int)
	amountInFloat.Int(amountOut)

	tx, err = swapRouter.ExactInputSingle(tops, uniswapv3.IV3SwapRouterExactInputSingleParams{
		// The contract address of the inbound token.
		TokenIn: swapDirection.tokenIn,
		// The contract address of the outbound token.
		TokenOut: swapDirection.tokenOut,
		// The fee tier of the pool, used to determine the correct pool contract in which to execute the swap.
		Fee: poolConfig.Fees,
		// The destination address of the outbound token.
		Recipient: recipient,
		// The amount of inbound token given as swap input.
		AmountIn: amountIn,
		// The minimum amount of outbound token received as swap output.
		AmountOutMinimum: amountOut,
		// The limit for the price swap.
		// Note: we set this to zero which makes the parameter inactive. In production, this value can
		// be used to protect against price impact.
		SqrtPriceLimitX96: big.NewInt(0),
	})
	if err != nil {
		log.Error().Err(err).Str("tokenIn", swapDirection.tokenInName).Str("tokenOut", swapDirection.tokenOutName).Interface("amountIn", amountIn).Msg("Unable to swap")
		return
	}
	log.Trace().Str("tokenIn", swapDirection.tokenInName).Str("tokenOut", swapDirection.tokenOutName).Interface("amountIn", amountIn).Msg("Successful swap")
	return
}

// swapDirection represents a swap direction with the inbound and outbound tokens.
type uniswapDirection struct {
	tokenIn, tokenOut         common.Address
	tokenInName, tokenOutName string
}

// Return the direction of the swap given the nonce value.
func getSwapDirection(nonce uint64, poolConfig PoolConfig) uniswapDirection {
	if nonce%2 == 0 {
		return uniswapDirection{
			tokenIn:     poolConfig.Token0.Address,
			tokenInName: "token0",

			tokenOut:     poolConfig.Token1.Address,
			tokenOutName: "token1",
		}
	}
	return uniswapDirection{
		tokenIn:     poolConfig.Token1.Address,
		tokenInName: "token1",

		tokenOut:     poolConfig.Token0.Address,
		tokenOutName: "token0",
	}
}
