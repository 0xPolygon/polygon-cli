package loadtest

import (
	"context"
	"encoding/json"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/maticnetwork/polygon-cli/rpctypes"
	"github.com/maticnetwork/polygon-cli/util"
	"math/big"
	"strings"
)

// TODO allow this to be pre-specified with an input file
func getRecentBlocks(ctx context.Context, ec *ethclient.Client, c *ethrpc.Client) ([]*json.RawMessage, error) {
	bn, err := ec.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	// FIXME the batch size of 25 is hard coded and probably should at least be a constant or a parameter. This limit is
	// different than the actual json RPC batch size of 999. Because we're fetching blocks, its' more likely that we hit
	// a response size limit rather than a batch length limit
	rawBlocks, err := util.GetBlockRangeInPages(ctx, bn-*inputLoadTestParams.RecallLength, bn, 25, c)
	return rawBlocks, err
}

func getRecallTransactions(ctx context.Context, c *ethclient.Client, rpc *ethrpc.Client) ([]rpctypes.PolyTransaction, error) {
	rb, err := getRecentBlocks(ctx, c, rpc)
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

// IndexedActivity is used to hold a bunch of values for testing an RPC
type IndexedActivity struct {
	BlockNumbers    []string
	TransactionIDs  []string
	BlockIDs        []string
	Addresses       []string
	ERC20Addresses  []string
	ERC721Addresses []string
	Contracts       []string
	BlockNumber     uint64
	Transactions    []rpctypes.PolyTransaction
}

func getIndexedRecentActivity(ctx context.Context, ec *ethclient.Client, c *ethrpc.Client) (*IndexedActivity, error) {
	blockData, err := getRecentBlocks(ctx, ec, c)
	if err != nil {
		return nil, err
	}

	ia := new(IndexedActivity)
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
			if strings.HasPrefix("0x70a08231", string(pt.Data())) {
				ia.ERC20Addresses = append(ia.ERC20Addresses, pt.To().String())
			}
			if strings.HasPrefix("0xc87b56dd", string(pt.Data())) {
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

func deduplicate(slice []string) []string {
	seen := make(map[string]struct{}) // struct{} takes no memory
	var result []string

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func rawTransactionToNewTx(pt rpctypes.PolyTransaction, nonce uint64, price, tipCap *big.Int) *ethtypes.Transaction {
	if pt.MaxFeePerGas() != 0 || pt.ChainID() != 0 {
		return rawTransactionToDynamicFeeTx(pt, nonce, price, tipCap)
	}
	return rawTransactionToLegacyTx(pt, nonce, price)
}

func rawTransactionToDynamicFeeTx(pt rpctypes.PolyTransaction, nonce uint64, price, tipCap *big.Int) *ethtypes.Transaction {
	toAddr := pt.To()
	chainId := new(big.Int).SetUint64(pt.ChainID())
	dynamicFeeTx := &ethtypes.DynamicFeeTx{
		ChainID:   chainId,
		To:        &toAddr,
		Data:      pt.Data(),
		Value:     pt.Value(),
		Gas:       pt.Gas(),
		GasFeeCap: price,
		GasTipCap: tipCap,
		Nonce:     nonce,
	}
	tx := ethtypes.NewTx(dynamicFeeTx)
	return tx
}

func rawTransactionToLegacyTx(pt rpctypes.PolyTransaction, nonce uint64, price *big.Int) *ethtypes.Transaction {
	toAddr := pt.To()
	tx := ethtypes.NewTx(&ethtypes.LegacyTx{
		To:       &toAddr,
		Value:    pt.Value(),
		Data:     pt.Data(),
		Gas:      pt.Gas(),
		Nonce:    nonce,
		GasPrice: price,
	})
	return tx
}
