package modes

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/0xPolygon/polygon-cli/loadtest/mode"
	"github.com/0xPolygon/polygon-cli/loadtest/uniswapv3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func init() {
	mode.Register(&UniswapV3Mode{})
}

// UniswapV3Mode implements UniswapV3 swap operations.
type UniswapV3Mode struct{}

func (m *UniswapV3Mode) Name() string {
	return "uniswapv3"
}

func (m *UniswapV3Mode) Aliases() []string {
	return []string{"v3"}
}

func (m *UniswapV3Mode) RequiresContract() bool {
	return false
}

func (m *UniswapV3Mode) RequiresERC20() bool {
	return false
}

func (m *UniswapV3Mode) RequiresERC721() bool {
	return false
}

func (m *UniswapV3Mode) Init(ctx context.Context, cfg *config.Config, deps *mode.Dependencies) error {
	return nil
}

func (m *UniswapV3Mode) Execute(ctx context.Context, cfg *config.Config, deps *mode.Dependencies, tops *bind.TransactOpts) (start, end time.Time, txHash common.Hash, err error) {
	if deps.UniswapV3Config == nil || deps.UniswapV3PoolConfig == nil {
		start = time.Now()
		end = start
		err = fmt.Errorf("uniswapv3 mode requires UniswapV3Config and UniswapV3PoolConfig to be set")
		return
	}

	swapAmountIn := big.NewInt(int64(cfg.UniswapV3.SwapAmountInput))
	recipient := *cfg.FromETHAddress

	return uniswapv3.Run(ctx, deps.Client, tops, *deps.UniswapV3Config, *deps.UniswapV3PoolConfig, swapAmountIn, recipient)
}
