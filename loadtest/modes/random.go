package modes

import (
	"context"
	"time"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/0xPolygon/polygon-cli/loadtest/mode"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func init() {
	mode.Register(&RandomMode{})
}

// RandomMode implements randomly selecting from other modes.
type RandomMode struct {
	modes []mode.Runner
}

func (m *RandomMode) Name() string {
	return "random"
}

func (m *RandomMode) Aliases() []string {
	return []string{"r"}
}

func (m *RandomMode) RequiresContract() bool {
	return true // random mode can select modes that require contract
}

func (m *RandomMode) RequiresERC20() bool {
	return true // random mode can select modes that require ERC20
}

func (m *RandomMode) RequiresERC721() bool {
	return true // random mode can select modes that require ERC721
}

func (m *RandomMode) Init(ctx context.Context, cfg *config.Config, deps *mode.Dependencies) error {
	// Use a deterministic, hardcoded list of modes (same as old behavior)
	// Does not include: blob, contract-call, recall, rpc, uniswapv3, random
	modeNames := []string{
		"deploy",
		"erc20",
		"erc721",
		"increment",
		"store",
		"transaction",
	}

	m.modes = make([]mode.Runner, 0, len(modeNames))
	for _, name := range modeNames {
		md, err := mode.Get(name)
		if err != nil {
			continue // skip if mode not found
		}
		m.modes = append(m.modes, md)
	}
	return nil
}

func (m *RandomMode) Execute(ctx context.Context, cfg *config.Config, deps *mode.Dependencies, tops *bind.TransactOpts) (start, end time.Time, txHash common.Hash, err error) {
	if len(m.modes) == 0 {
		start = time.Now()
		end = start
		return
	}

	// Select a random mode
	selectedMode := m.modes[deps.RandIntn(len(m.modes))]
	return selectedMode.Execute(ctx, cfg, deps, tops)
}

// GetRandomModeEnumValue returns a random Mode enum value.
// Does not include: blob, contract-call, recall, rpc, uniswapv3
func GetRandomModeEnumValue(deps *mode.Dependencies) config.Mode {
	m := []config.Mode{
		config.ModeERC20,
		config.ModeERC721,
		config.ModeDeploy,
		config.ModeIncrement,
		config.ModeStore,
		config.ModeTransaction,
	}
	return m[deps.RandIntn(len(m))]
}
