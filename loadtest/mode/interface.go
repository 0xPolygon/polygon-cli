package mode

import (
	"context"
	"time"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// Runner defines the interface that all load test modes must implement.
type Runner interface {
	// Name returns the canonical name of the mode.
	Name() string

	// Aliases returns alternative names/shortcuts for this mode.
	Aliases() []string

	// RequiresContract returns true if this mode needs the LoadTester contract.
	RequiresContract() bool

	// RequiresERC20 returns true if this mode needs an ERC20 contract.
	RequiresERC20() bool

	// RequiresERC721 returns true if this mode needs an ERC721 contract.
	RequiresERC721() bool

	// Init sets up any mode-specific state before execution.
	Init(ctx context.Context, cfg *config.Config, deps *Dependencies) error

	// Execute performs a single load test operation.
	// Returns start time, end time, transaction hash, and any error.
	Execute(ctx context.Context, cfg *config.Config, deps *Dependencies, opts *bind.TransactOpts) (start, end time.Time, txHash common.Hash, err error)
}

// InputReserver is implemented by modes that draw each transaction from an
// external input stream. The runner calls ReserveInput before consuming any
// per-request resource (nonce, gas budget) so that an exhausted stream is
// detected while nothing has been taken. The reserved input is handed to
// Execute through the context via WithInput. ErrInputExhausted stops the
// test cleanly; any other error fails it.
type InputReserver interface {
	ReserveInput(ctx context.Context) (any, error)
}

type inputKey struct{}

// WithInput attaches a reserved per-request input to ctx for Execute.
func WithInput(ctx context.Context, input any) context.Context {
	return context.WithValue(ctx, inputKey{}, input)
}

// InputFromContext returns the input attached by WithInput, if any.
func InputFromContext(ctx context.Context) (any, bool) {
	input := ctx.Value(inputKey{})
	return input, input != nil
}
