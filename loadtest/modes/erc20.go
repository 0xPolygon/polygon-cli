package modes

import (
	"context"
	"time"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/0xPolygon/polygon-cli/loadtest/mode"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

func init() {
	mode.Register(&ERC20Mode{})
}

// ERC20Mode implements ERC20 token transfers.
type ERC20Mode struct{}

func (m *ERC20Mode) Name() string {
	return "erc20"
}

func (m *ERC20Mode) Aliases() []string {
	return []string{"2"}
}

func (m *ERC20Mode) RequiresContract() bool {
	return false
}

func (m *ERC20Mode) RequiresERC20() bool {
	return true
}

func (m *ERC20Mode) RequiresERC721() bool {
	return false
}

func (m *ERC20Mode) Init(ctx context.Context, cfg *config.Config, deps *mode.Dependencies) error {
	return nil
}

func (m *ERC20Mode) Execute(ctx context.Context, cfg *config.Config, deps *mode.Dependencies, tops *bind.TransactOpts) (start, end time.Time, txHash common.Hash, err error) {
	to := cfg.ToETHAddress
	if cfg.RandomRecipients {
		to = mode.GetRandomAddress(deps)
	}
	amount := cfg.SendAmount

	start = time.Now()
	defer func() { end = time.Now() }()

	if cfg.EthCallOnly {
		tops.NoSend = true
		tx, iErr := deps.ERC20Contract.Transfer(tops, *to, amount)
		if iErr != nil {
			err = iErr
			return
		}
		msg := mode.TxToCallMsg(cfg, tx)
		_, err = deps.Client.CallContract(ctx, msg, nil)
	} else if cfg.OutputRawTxOnly {
		tops.NoSend = true
		tx, iErr := deps.ERC20Contract.Transfer(tops, *to, amount)
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
		tx, iErr := deps.ERC20Contract.Transfer(tops, *to, amount)
		if iErr == nil && tx != nil {
			txHash = tx.Hash()
		}
		err = iErr
	}

	if err != nil {
		log.Error().Err(err).Msg("ERC20 transfer failed")
	}
	return
}
