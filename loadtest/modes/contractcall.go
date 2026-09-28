package modes

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"os"
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
type ContractCallMode struct {
	// feeder supplies per-transaction calldata when --calldata-stdin is set.
	// Nil otherwise, in which case cfg.ContractCallData is reused for every
	// transaction.
	feeder *mode.ChunkFeeder
}

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
	// The mode is a process-wide singleton, so clear any feeder left over
	// from an earlier run in the same process.
	m.feeder = nil
	if !cfg.ContractCallDataStdin {
		return nil
	}
	if cfg.ForceGasLimit == 0 {
		log.Warn().Msg("--calldata-stdin without --gas-limit estimates gas for every transaction, adding one RPC call per send")
	}
	return m.initFeeder(ctx, os.Stdin, cfg)
}

// initFeeder starts the calldata feeder on r. Split from Init so tests can
// substitute an in-memory reader for stdin.
func (m *ContractCallMode) initFeeder(ctx context.Context, r io.Reader, cfg *config.Config) error {
	buffer := int(cfg.Concurrency) * 2
	if buffer < 2 {
		buffer = 2
	}
	feeder, err := mode.NewChunkFeeder(ctx, r, int(cfg.ContractCallDataSize), buffer)
	if err != nil {
		return fmt.Errorf("failed to start calldata feeder: %w", err)
	}
	m.feeder = feeder
	return nil
}

// ReserveInput implements mode.InputReserver. With --calldata-stdin it pulls
// the next chunk before the runner takes a nonce, so an exhausted stream never
// leaves a reserved nonce unsent. Without stdin there is nothing to reserve.
func (m *ContractCallMode) ReserveInput(ctx context.Context) (any, error) {
	if m.feeder == nil {
		return nil, nil
	}
	// Returned unwrapped so the runner can match mode.ErrInputExhausted.
	return m.feeder.Next(ctx)
}

func (m *ContractCallMode) Execute(ctx context.Context, cfg *config.Config, deps *mode.Dependencies, tops *bind.TransactOpts) (start, end time.Time, txHash common.Hash, err error) {
	to := cfg.ContractETHAddress
	amount := big.NewInt(0)
	if cfg.ContractCallPayable {
		amount = cfg.SendAmount
	}

	var calldata []byte
	if m.feeder != nil {
		input, ok := mode.InputFromContext(ctx)
		if !ok {
			err = fmt.Errorf("calldata was not reserved for this request; the runner must call ReserveInput before Execute")
			return
		}
		calldata, ok = input.([]byte)
		if !ok {
			err = fmt.Errorf("reserved calldata has unexpected type %T", input)
			return
		}
	} else {
		if cfg.ContractCallData == "" {
			err = fmt.Errorf("missing calldata for function call")
			log.Error().Err(err).Msg("--calldata flag is required for contract-call mode")
			return
		}
		calldata, err = hex.DecodeString(strings.TrimPrefix(cfg.ContractCallData, "0x"))
		if err != nil {
			log.Error().Err(err).Msg("Unable to decode calldata string")
			return
		}
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

	tx := buildContractCallTx(cfg, tops, to, amount, calldata)
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
	} else {
		err = mode.SendSignedTransaction(ctx, deps, cfg, stx)
	}
	return
}

// buildContractCallTx assembles the unsigned transaction for a contract call
// using the fee model selected by cfg.LegacyTxMode.
func buildContractCallTx(cfg *config.Config, tops *bind.TransactOpts, to *common.Address, amount *big.Int, calldata []byte) *types.Transaction {
	if cfg.LegacyTxMode {
		return types.NewTx(&types.LegacyTx{
			Nonce:    tops.Nonce.Uint64(),
			To:       to,
			Value:    amount,
			Gas:      tops.GasLimit,
			GasPrice: tops.GasPrice,
			Data:     calldata,
		})
	}
	return types.NewTx(&types.DynamicFeeTx{
		ChainID:   new(big.Int).SetUint64(cfg.ChainID),
		Nonce:     tops.Nonce.Uint64(),
		To:        to,
		Gas:       tops.GasLimit,
		GasFeeCap: tops.GasFeeCap,
		GasTipCap: tops.GasTipCap,
		Data:      calldata,
		Value:     amount,
	})
}
