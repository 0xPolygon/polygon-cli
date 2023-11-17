package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockUntilSuccessfulFn is designed to wait until a specified number of Ethereum blocks have been
// mined, periodically checking for the completion of a given function within each block interval.
type BlockUntilSuccessfulFn func(ctx context.Context, c *ethclient.Client, f func() error) error
