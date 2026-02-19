package modes

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/0xPolygon/polygon-cli/loadtest/mode"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

func init() {
	mode.Register(&ERC721Mode{})
}

// ERC721Mode implements ERC721 token minting.
type ERC721Mode struct{}

func (m *ERC721Mode) Name() string {
	return "erc721"
}

func (m *ERC721Mode) Aliases() []string {
	return []string{"7"}
}

func (m *ERC721Mode) RequiresContract() bool {
	return false
}

func (m *ERC721Mode) RequiresERC20() bool {
	return false
}

func (m *ERC721Mode) RequiresERC721() bool {
	return true
}

func (m *ERC721Mode) Init(ctx context.Context, cfg *config.Config, deps *mode.Dependencies) error {
	return nil
}

func (m *ERC721Mode) Execute(ctx context.Context, cfg *config.Config, deps *mode.Dependencies, tops *bind.TransactOpts) (start, end time.Time, txHash common.Hash, err error) {
	to := cfg.ToETHAddress
	if cfg.RandomRecipients {
		to = mode.GetRandomAddress(deps)
	}

	start = time.Now()
	defer func() { end = time.Now() }()

	if cfg.EthCallOnly {
		tops.NoSend = true
		tx, iErr := deps.ERC721Contract.MintBatch(tops, *to, big.NewInt(1))
		if iErr != nil {
			err = iErr
			return
		}
		msg := mode.TxToCallMsg(cfg, tx)
		_, err = deps.Client.CallContract(ctx, msg, nil)
	} else if cfg.OutputRawTxOnly {
		tops.NoSend = true
		tx, iErr := deps.ERC721Contract.MintBatch(tops, *to, big.NewInt(1))
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
		tx, iErr := deps.ERC721Contract.MintBatch(tops, *to, big.NewInt(1))
		if iErr == nil && tx != nil {
			txHash = tx.Hash()
		}
		err = iErr
	}

	if err != nil {
		log.Error().Err(err).Msg("ERC721 mint failed")
	}
	return
}
