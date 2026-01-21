package modes

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/0xPolygon/polygon-cli/loadtest/mode"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

func init() {
	mode.Register(&TransactionMode{})
}

// TransactionMode implements basic ETH transfer transactions.
type TransactionMode struct{}

func (m *TransactionMode) Name() string {
	return "transaction"
}

func (m *TransactionMode) Aliases() []string {
	return []string{"t"}
}

func (m *TransactionMode) RequiresContract() bool {
	return false
}

func (m *TransactionMode) RequiresERC20() bool {
	return false
}

func (m *TransactionMode) RequiresERC721() bool {
	return false
}

func (m *TransactionMode) Init(ctx context.Context, cfg *config.Config, deps *mode.Dependencies) error {
	return nil
}

func (m *TransactionMode) Execute(ctx context.Context, cfg *config.Config, deps *mode.Dependencies, tops *bind.TransactOpts) (start, end time.Time, txHash common.Hash, err error) {
	to := cfg.ToETHAddress
	if cfg.RandomRecipients {
		to = mode.GetRandomAddress(deps)
	}

	tops.GasLimit = uint64(21000)

	amount := cfg.SendAmount
	chainID := new(big.Int).SetUint64(cfg.ChainID)

	var tx *types.Transaction
	if cfg.LegacyTxMode {
		tx = types.NewTx(&types.LegacyTx{
			Nonce:    tops.Nonce.Uint64(),
			To:       to,
			Value:    amount,
			Gas:      tops.GasLimit,
			GasPrice: tops.GasPrice,
			Data:     nil,
		})
	} else {
		dynamicFeeTx := &types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     tops.Nonce.Uint64(),
			To:        to,
			Gas:       tops.GasLimit,
			GasFeeCap: tops.GasFeeCap,
			GasTipCap: tops.GasTipCap,
			Data:      nil,
			Value:     amount,
		}
		tx = types.NewTx(dynamicFeeTx)
	}

	stx, err := tops.Signer(tops.From, tx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to sign transaction")
		return
	}

	txHash = stx.Hash()

	start = time.Now()
	defer func() { end = time.Now() }()
	if cfg.EthCallOnly {
		_, err = deps.Client.CallContract(ctx, mode.TxToCallMsg(cfg, stx), nil)
	} else if cfg.OutputRawTxOnly {
		err = mode.OutputRawTransaction(stx)
	} else {
		err = deps.Client.SendTransaction(ctx, stx)
	}

	return
}
