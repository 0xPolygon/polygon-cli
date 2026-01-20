package modes

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/0xPolygon/polygon-cli/loadtest/mode"
	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/0xPolygon/polygon-cli/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
)

func init() {
	mode.Register(&RecallMode{})
}

// RecallMode implements replaying historical transactions.
type RecallMode struct{}

func (m *RecallMode) Name() string {
	return "recall"
}

func (m *RecallMode) Aliases() []string {
	return []string{"R"}
}

func (m *RecallMode) RequiresContract() bool {
	return false
}

func (m *RecallMode) RequiresERC20() bool {
	return false
}

func (m *RecallMode) RequiresERC721() bool {
	return false
}

func (m *RecallMode) Init(ctx context.Context, cfg *config.Config, deps *mode.Dependencies) error {
	return nil
}

func (m *RecallMode) Execute(ctx context.Context, cfg *config.Config, deps *mode.Dependencies, tops *bind.TransactOpts) (start, end time.Time, txHash common.Hash, err error) {
	if len(deps.RecallTransactions) == 0 {
		err = fmt.Errorf("no recall transactions available")
		log.Error().Err(err).Msg("No recall transactions")
		return
	}

	originalTx := deps.RecallTransactions[int(tops.Nonce.Uint64())%len(deps.RecallTransactions)]

	// For EIP-1559 transactions, use GasFeeCap instead of GasPrice (which is nil for dynamic fee txs)
	gasPrice := tops.GasPrice
	if gasPrice == nil && tops.GasFeeCap != nil {
		gasPrice = tops.GasFeeCap
	}
	tx := RawTransactionToNewTx(originalTx, tops.Nonce.Uint64(), gasPrice, tops.GasTipCap)

	stx, err := tops.Signer(tops.From, tx)
	if err != nil {
		log.Error().Err(err).Msg("Unable to sign transaction")
		return
	}
	log.Trace().Str("txID", originalTx.Hash().String()).Bool("callOnly", cfg.EthCallOnly).Msg("Attempting to replay transaction")
	txHash = stx.Hash()

	start = time.Now()
	defer func() { end = time.Now() }()

	if cfg.EthCallOnly {
		callMsg := mode.TxToCallMsg(cfg, stx)
		callMsg.From = originalTx.From()
		callMsg.Gas = originalTx.Gas()
		if cfg.EthCallOnlyLatestBlock {
			_, err = deps.Client.CallContract(ctx, callMsg, nil)
		} else {
			callMsg.GasFeeCap = new(big.Int).SetUint64(originalTx.MaxFeePerGas())
			callMsg.GasTipCap = new(big.Int).SetUint64(originalTx.MaxPriorityFeePerGas())
			if originalTx.MaxFeePerGas() == 0 && originalTx.MaxPriorityFeePerGas() == 0 {
				callMsg.GasPrice = originalTx.GasPrice()
				callMsg.GasFeeCap = nil
				callMsg.GasTipCap = nil
			} else {
				callMsg.GasPrice = nil
			}

			_, err = deps.Client.CallContract(ctx, callMsg, originalTx.BlockNumber())
		}
		if err != nil {
			log.Warn().Err(err).Msg("Recall failure")
		}
		// we're not going to return the error in the case because there is no point retrying
		err = nil
	} else if cfg.OutputRawTxOnly {
		err = mode.OutputRawTransaction(stx)
	} else {
		err = deps.Client.SendTransaction(ctx, stx)
	}
	return
}

// RawTransactionToNewTx converts a PolyTransaction to a new transaction with updated nonce and gas prices.
func RawTransactionToNewTx(pt rpctypes.PolyTransaction, nonce uint64, price, tipCap *big.Int) *types.Transaction {
	if pt.MaxFeePerGas() != 0 || pt.ChainID() != 0 {
		return rawTransactionToDynamicFeeTx(pt, nonce, price, tipCap)
	}
	return rawTransactionToLegacyTx(pt, nonce, price)
}

func rawTransactionToDynamicFeeTx(pt rpctypes.PolyTransaction, nonce uint64, price, tipCap *big.Int) *types.Transaction {
	toAddr := pt.To()
	chainID := new(big.Int).SetUint64(pt.ChainID())
	dynamicFeeTx := &types.DynamicFeeTx{
		ChainID:   chainID,
		To:        &toAddr,
		Data:      pt.Data(),
		Value:     pt.Value(),
		Gas:       pt.Gas(),
		GasFeeCap: price,
		GasTipCap: tipCap,
		Nonce:     nonce,
	}
	tx := types.NewTx(dynamicFeeTx)
	return tx
}

func rawTransactionToLegacyTx(pt rpctypes.PolyTransaction, nonce uint64, price *big.Int) *types.Transaction {
	toAddr := pt.To()
	tx := types.NewTx(&types.LegacyTx{
		To:       &toAddr,
		Value:    pt.Value(),
		Data:     pt.Data(),
		Gas:      pt.Gas(),
		Nonce:    nonce,
		GasPrice: price,
	})
	return tx
}

// GetRecentBlocks fetches recent blocks from the chain.
func GetRecentBlocks(ctx context.Context, ec *ethclient.Client, c *ethrpc.Client, recallLength, blockBatchSize uint64, onlyTxHashes bool) ([]*json.RawMessage, error) {
	bn, err := ec.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	rawBlocks, err := util.GetBlockRangeInPages(ctx, bn-recallLength, bn, blockBatchSize, c, onlyTxHashes)
	return rawBlocks, err
}

// GetRecallTransactions fetches transactions from recent blocks for replay.
func GetRecallTransactions(ctx context.Context, c *ethclient.Client, rpc *ethrpc.Client, recallLength, blockBatchSize uint64) ([]rpctypes.PolyTransaction, error) {
	rb, err := GetRecentBlocks(ctx, c, rpc, recallLength, blockBatchSize, false)
	if err != nil {
		return nil, err
	}
	txs := make([]rpctypes.PolyTransaction, 0)
	for _, v := range rb {
		pb := new(rpctypes.RawBlockResponse)
		err := json.Unmarshal(*v, pb)
		if err != nil {
			return nil, err
		}
		for k := range pb.Transactions {
			pt := rpctypes.NewPolyTransaction(&pb.Transactions[k])
			txs = append(txs, pt)
		}
	}
	return txs, nil
}

// GetIndexedRecentActivity builds indexed activity data from recent blocks for RPC mode.
func GetIndexedRecentActivity(ctx context.Context, ec *ethclient.Client, c *ethrpc.Client, recallLength, blockBatchSize uint64) (*mode.IndexedActivity, error) {
	blockData, err := GetRecentBlocks(ctx, ec, c, recallLength, blockBatchSize, false)
	if err != nil {
		return nil, err
	}

	ia := new(mode.IndexedActivity)
	ia.BlockNumbers = make([]string, 0)
	ia.TransactionIDs = make([]string, 0)
	ia.Transactions = make([]rpctypes.PolyTransaction, 0)
	ia.BlockIDs = make([]string, 0)
	ia.Addresses = make([]string, 0)
	ia.ERC20Addresses = make([]string, 0)
	ia.ERC721Addresses = make([]string, 0)
	ia.Contracts = make([]string, 0)
	for _, block := range blockData {
		pb := new(rpctypes.RawBlockResponse)
		err = json.Unmarshal(*block, pb)
		if err != nil {
			return nil, err
		}
		ia.BlockIDs = append(ia.BlockIDs, string(pb.Hash))
		ia.BlockNumbers = append(ia.BlockNumbers, string(pb.Number))
		for k := range pb.Transactions {
			pt := rpctypes.NewPolyTransaction(&pb.Transactions[k])
			ia.TransactionIDs = append(ia.TransactionIDs, pt.Hash().String())
			ia.Transactions = append(ia.Transactions, pt)
			ia.Addresses = append(ia.Addresses, pt.From().String(), pt.To().String())

			// balanceOf(address)
			if strings.HasPrefix(string(pt.Data()), "0x70a08231") {
				ia.ERC20Addresses = append(ia.ERC20Addresses, pt.To().String())
			}
			if strings.HasPrefix(string(pt.Data()), "0xc87b56dd") {
				ia.ERC721Addresses = append(ia.ERC721Addresses, pt.To().String())
			}
			if len(string(pt.Data())) > 10 {
				ia.Contracts = append(ia.Contracts, pt.To().String())
			}
		}
	}
	ia.BlockNumbers = deduplicate(ia.BlockNumbers)
	ia.TransactionIDs = deduplicate(ia.TransactionIDs)
	ia.BlockIDs = deduplicate(ia.BlockIDs)
	ia.Addresses = deduplicate(ia.Addresses)
	ia.ERC20Addresses = deduplicate(ia.ERC20Addresses)
	ia.ERC721Addresses = deduplicate(ia.ERC721Addresses)
	ia.Contracts = deduplicate(ia.Contracts)

	ia.BlockNumber, err = ec.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	return ia, nil
}

func deduplicate(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	result := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}
