package modes

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/0xPolygon/polygon-cli/bindings/tokens"
	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/0xPolygon/polygon-cli/loadtest/mode"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/rs/zerolog/log"
)

func init() {
	mode.Register(&RPCMode{})
}

// RPCMode implements random RPC method calls.
type RPCMode struct{}

func (m *RPCMode) Name() string {
	return "rpc"
}

func (m *RPCMode) Aliases() []string {
	return []string{}
}

func (m *RPCMode) RequiresContract() bool {
	return false
}

func (m *RPCMode) RequiresERC20() bool {
	return true
}

func (m *RPCMode) RequiresERC721() bool {
	return true
}

func (m *RPCMode) Init(ctx context.Context, cfg *config.Config, deps *mode.Dependencies) error {
	return nil
}

func (m *RPCMode) Execute(ctx context.Context, cfg *config.Config, deps *mode.Dependencies, tops *bind.TransactOpts) (start, end time.Time, txHash common.Hash, err error) {
	ia := deps.IndexedActivity
	if ia == nil {
		start = time.Now()
		end = start
		return
	}

	funcNum := deps.RandSource.Intn(300)
	start = time.Now()
	defer func() { end = time.Now() }()

	if funcNum < 10 {
		log.Trace().Msg("eth_gasPrice")
		_, err = deps.Client.SuggestGasPrice(ctx)
	} else if funcNum < 21 {
		log.Trace().Msg("eth_estimateGas")
		var rawTxData []byte
		pt := ia.Transactions[deps.RandSource.Intn(len(ia.Transactions))]
		rawTxData, err = pt.MarshalJSON()
		if err != nil {
			log.Error().Err(err).Str("txHash", pt.Hash().String()).Msg("issue converting poly transaction to json")
			return
		}
		var txArgs apitypes.SendTxArgs
		if err = json.Unmarshal(rawTxData, &txArgs); err != nil {
			log.Error().Err(err).Str("txHash", pt.Hash().String()).Msg("unable to unmarshal poly transaction to json")
			return
		}
		var tx *ethtypes.Transaction
		tx, err = txArgs.ToTransaction()
		if err != nil {
			log.Error().Err(err).Str("txArgs", txArgs.String()).Msg("unable to convert the arguments to a transaction")
			return
		}
		cm := mode.TxToCallMsg(cfg, tx)
		cm.From = pt.From()
		_, err = deps.Client.EstimateGas(ctx, cm)
	} else if funcNum < 33 {
		log.Trace().Msg("eth_getTransactionCount")
		_, err = deps.Client.NonceAt(ctx, common.HexToAddress(ia.Addresses[deps.RandSource.Intn(len(ia.Addresses))]), nil)
	} else if funcNum < 47 {
		log.Trace().Msg("eth_getCode")
		_, err = deps.Client.CodeAt(ctx, common.HexToAddress(ia.Contracts[deps.RandSource.Intn(len(ia.Contracts))]), nil)
	} else if funcNum < 64 {
		log.Trace().Msg("eth_getBlockByNumber")
		_, err = deps.Client.BlockByNumber(ctx, big.NewInt(int64(deps.RandSource.Intn(int(ia.BlockNumber)))))
	} else if funcNum < 84 {
		log.Trace().Msg("eth_getTransactionByHash")
		_, _, err = deps.Client.TransactionByHash(ctx, common.HexToHash(ia.TransactionIDs[deps.RandSource.Intn(len(ia.TransactionIDs))]))
	} else if funcNum < 109 {
		log.Trace().Msg("eth_getBalance")
		_, err = deps.Client.BalanceAt(ctx, common.HexToAddress(ia.Addresses[deps.RandSource.Intn(len(ia.Addresses))]), nil)
	} else if funcNum < 142 {
		log.Trace().Msg("eth_getTransactionReceipt")
		_, err = deps.Client.TransactionReceipt(ctx, common.HexToHash(ia.TransactionIDs[deps.RandSource.Intn(len(ia.TransactionIDs))]))
	} else if funcNum < 192 {
		log.Trace().Msg("eth_getLogs")
		h := common.HexToHash(ia.BlockIDs[deps.RandSource.Intn(len(ia.BlockIDs))])
		_, err = deps.Client.FilterLogs(ctx, ethereum.FilterQuery{BlockHash: &h})
	} else {
		log.Trace().Msg("eth_call")

		if len(ia.ERC20Addresses) != 0 {
			erc20Str := string(ia.ERC20Addresses[deps.RandSource.Intn(len(ia.ERC20Addresses))])
			erc20Addr := common.HexToAddress(erc20Str)

			log.Trace().
				Str("erc20str", erc20Str).
				Stringer("erc20addr", erc20Addr).
				Msg("Retrieve contract addresses")
			cops := new(bind.CallOpts)
			cops.Context = ctx
			var erc20Contract *tokens.ERC20

			erc20Contract, err = tokens.NewERC20(erc20Addr, deps.Client)
			if err != nil {
				log.Error().Err(err).Msg("Unable to instantiate new erc20 contract")
				return
			}
			start = time.Now()

			_, err = erc20Contract.BalanceOf(cops, *cfg.FromETHAddress)
			if err != nil && err == bind.ErrNoCode {
				err = nil
			}
		} else {
			log.Warn().Msg("Unable to find deployed erc20 contract, skipping making calls...")
		}

		if len(ia.ERC721Addresses) != 0 {
			erc721Str := string(ia.ERC721Addresses[deps.RandSource.Intn(len(ia.ERC721Addresses))])
			erc721Addr := common.HexToAddress(erc721Str)

			log.Trace().
				Str("erc721str", erc721Str).
				Stringer("erc721addr", erc721Addr).
				Msg("Retrieve contract addresses")
			cops := new(bind.CallOpts)
			cops.Context = ctx
			var erc721Contract *tokens.ERC721

			erc721Contract, err = tokens.NewERC721(erc721Addr, deps.Client)
			if err != nil {
				log.Error().Err(err).Msg("Unable to instantiate new erc721 contract")
				return
			}
			start = time.Now()

			_, err = erc721Contract.BalanceOf(cops, *cfg.FromETHAddress)
			if err != nil && err == bind.ErrNoCode {
				err = nil
			}
		} else {
			log.Warn().Msg("Unable to find deployed erc721 contract, skipping making calls...")
		}
	}

	return
}
