package modes

import (
	"context"
	"time"

	"github.com/0xPolygon/polygon-cli/bindings/tester"
	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/0xPolygon/polygon-cli/loadtest/mode"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func init() {
	mode.Register(&DeployMode{})
}

// DeployMode implements contract deployment.
type DeployMode struct{}

func (m *DeployMode) Name() string {
	return "deploy"
}

func (m *DeployMode) Aliases() []string {
	return []string{"d"}
}

func (m *DeployMode) RequiresContract() bool {
	return false
}

func (m *DeployMode) RequiresERC20() bool {
	return false
}

func (m *DeployMode) RequiresERC721() bool {
	return false
}

func (m *DeployMode) Init(ctx context.Context, cfg *config.Config, deps *mode.Dependencies) error {
	return nil
}

func (m *DeployMode) Execute(ctx context.Context, cfg *config.Config, deps *mode.Dependencies, tops *bind.TransactOpts) (start, end time.Time, txHash common.Hash, err error) {
	start = time.Now()
	defer func() { end = time.Now() }()

	if cfg.EthCallOnly {
		msg := mode.TransactOptsToCallMsg(cfg, tops.GasLimit)
		msg.Data = common.FromHex(tester.LoadTesterMetaData.Bin)
		_, err = deps.Client.CallContract(ctx, msg, nil)
	} else if cfg.OutputRawTxOnly {
		tops.NoSend = true
		var tx any
		_, tx, _, err = tester.DeployLoadTester(tops, deps.Client)
		if err != nil {
			return
		}
		if tx != nil {
			// Type assert to get hash and output
			if typedTx, ok := tx.(interface{ Hash() common.Hash }); ok {
				txHash = typedTx.Hash()
			}
			if typedTx, ok := tx.(interface{ MarshalBinary() ([]byte, error) }); ok {
				rawTx, marshalErr := typedTx.MarshalBinary()
				if marshalErr != nil {
					err = marshalErr
					return
				}
				err = mode.OutputRawBytes(rawTx)
			}
		}
	} else {
		var tx any
		_, tx, _, err = tester.DeployLoadTester(tops, deps.Client)
		if err == nil && tx != nil {
			if typedTx, ok := tx.(interface{ Hash() common.Hash }); ok {
				txHash = typedTx.Hash()
			}
		}
	}
	return
}
