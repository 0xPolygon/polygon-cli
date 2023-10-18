package uniswapv3loadtest

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

type blockUntilSuccessfulFn func(ctx context.Context, c *ethclient.Client, f func() error) error
