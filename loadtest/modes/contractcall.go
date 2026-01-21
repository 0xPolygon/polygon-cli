package modes

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/0xPolygon/polygon-cli/loadtest/mode"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

func init() {
	mode.Register(&ContractCallMode{})
}

// ContractCallMode implements generic contract calls.
type ContractCallMode struct{}

func (m *ContractCallMode) Name() string {
	return "contract-call"
}

func (m *ContractCallMode) Aliases() []string {
	return []string{"cc"}
}

func (m *ContractCallMode) RequiresContract() bool {
	return false
}

func (m *ContractCallMode) RequiresERC20() bool {
	return false
}

func (m *ContractCallMode) RequiresERC721() bool {
	return false
}

func (m *ContractCallMode) Init(ctx context.Context, cfg *config.Config, deps *mode.Dependencies) error {
	return nil
}

func (m *ContractCallMode) Execute(ctx context.Context, cfg *config.Config, deps *mode.Dependencies, tops *bind.TransactOpts) (start, end time.Time, txHash common.Hash, err error) {
	to := cfg.ContractETHAddress
	chainID := new(big.Int).SetUint64(cfg.ChainID)
	amount := big.NewInt(0)
	if cfg.ContractCallPayable {
		amount = cfg.SendAmount
	}

	if cfg.ContractCallData == "" {
		err = fmt.Errorf("missing calldata for function call")
		log.Error().Err(err).Msg("--calldata flag is required for contract-call mode")
		return
	}

	calldata, err := hex.DecodeString(strings.TrimPrefix(cfg.ContractCallData, "0x"))
	if err != nil {
		log.Error().Err(err).Msg("Unable to decode calldata string")
		return
	}

	if tops.GasLimit == 0 {
		estimateInput := ethereum.CallMsg{
			From:      tops.From,
			To:        to,
			Value:     amount,
			GasPrice:  tops.GasPrice,
			GasTipCap: tops.GasTipCap,
			GasFeeCap: tops.GasFeeCap,
			Data:      calldata,
		}
		tops.GasLimit, err = deps.Client.EstimateGas(ctx, estimateInput)
		if err != nil {
			log.Error().Err(err).Msg("Unable to estimate gas for transaction. Manually setting gas-limit might be required")
			return
		}
	}

	var tx *types.Transaction
	if cfg.LegacyTxMode {
		tx = types.NewTx(&types.LegacyTx{
			Nonce:    tops.Nonce.Uint64(),
			To:       to,
			Value:    amount,
			Gas:      tops.GasLimit,
			GasPrice: tops.GasPrice,
			Data:     calldata,
		})
	} else {
		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     tops.Nonce.Uint64(),
			To:        to,
			Gas:       tops.GasLimit,
			GasFeeCap: tops.GasFeeCap,
			GasTipCap: tops.GasTipCap,
			Data:      calldata,
			Value:     amount,
		})
	}
	log.Trace().Interface("tx", tx).Msg("Contract call data")

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
