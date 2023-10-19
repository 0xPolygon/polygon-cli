package uniswapv3loadtest

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

// blockUntilSuccessfulFn is designed to wait until a specified number of Ethereum blocks have been
// mined, periodically checking for the completion of a given function within each block interval.
type blockUntilSuccessfulFn func(ctx context.Context, c *ethclient.Client, f func() error) error
