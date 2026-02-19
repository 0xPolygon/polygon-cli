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
	mode.Register(&IncrementMode{})
}

// IncrementMode implements counter incrementing via LoadTester contract.
type IncrementMode struct{}

func (m *IncrementMode) Name() string {
	return "increment"
}

func (m *IncrementMode) Aliases() []string {
	return []string{"inc"}
}

func (m *IncrementMode) RequiresContract() bool {
	return true
}

func (m *IncrementMode) RequiresERC20() bool {
	return false
}

func (m *IncrementMode) RequiresERC721() bool {
	return false
}

func (m *IncrementMode) Init(ctx context.Context, cfg *config.Config, deps *mode.Dependencies) error {
	return nil
}

func (m *IncrementMode) Execute(ctx context.Context, cfg *config.Config, deps *mode.Dependencies, tops *bind.TransactOpts) (start, end time.Time, txHash common.Hash, err error) {
	start = time.Now()
	defer func() { end = time.Now() }()

	if cfg.EthCallOnly {
		tops.NoSend = true
		tx, iErr := deps.LoadTesterContract.Inc(tops)
		if iErr != nil {
			err = iErr
			return
		}
		msg := mode.TxToCallMsg(cfg, tx)
		_, err = deps.Client.CallContract(ctx, msg, nil)
	} else if cfg.OutputRawTxOnly {
		tops.NoSend = true
		tx, iErr := deps.LoadTesterContract.Inc(tops)
		if iErr != nil {
			err = iErr
			return
		}
		signedTx, signErr := tops.Signer(tops.From, tx)
		if signErr != nil {
			err = signErr
			return
		}
		txHash = signedTx.Hash()
		err = mode.OutputRawTransaction(signedTx)
	} else {
		tx, iErr := deps.LoadTesterContract.Inc(tops)
		if iErr == nil && tx != nil {
			txHash = tx.Hash()
		}
		err = iErr
	}
	return
}
