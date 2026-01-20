package modes

import (
	"context"
	"io"
	"time"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/0xPolygon/polygon-cli/loadtest/mode"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func init() {
	mode.Register(&StoreMode{})
}

// StoreMode implements storing bytes in the LoadTester contract.
type StoreMode struct{}

func (m *StoreMode) Name() string {
	return "store"
}

func (m *StoreMode) Aliases() []string {
	return []string{"s"}
}

func (m *StoreMode) RequiresContract() bool {
	return true
}

func (m *StoreMode) RequiresERC20() bool {
	return false
}

func (m *StoreMode) RequiresERC721() bool {
	return false
}

func (m *StoreMode) Init(ctx context.Context, cfg *config.Config, deps *mode.Dependencies) error {
	return nil
}

func (m *StoreMode) Execute(ctx context.Context, cfg *config.Config, deps *mode.Dependencies, tops *bind.TransactOpts) (start, end time.Time, txHash common.Hash, err error) {
	inputData := make([]byte, cfg.StoreDataSize)
	_, _ = io.ReadFull(mode.NewHexwordReader(deps.RandSource), inputData)

	start = time.Now()
	defer func() { end = time.Now() }()

	if cfg.EthCallOnly {
		tops.NoSend = true
		tx, iErr := deps.LoadTesterContract.Store(tops, inputData)
		if iErr != nil {
			err = iErr
			return
		}
		msg := mode.TxToCallMsg(cfg, tx)
		_, err = deps.Client.CallContract(ctx, msg, nil)
	} else if cfg.OutputRawTxOnly {
		tops.NoSend = true
		tx, iErr := deps.LoadTesterContract.Store(tops, inputData)
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
		tx, iErr := deps.LoadTesterContract.Store(tops, inputData)
		if iErr == nil && tx != nil {
			txHash = tx.Hash()
		}
		err = iErr
	}
	return
}
