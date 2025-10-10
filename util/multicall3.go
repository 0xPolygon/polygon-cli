package util

import (
	"context"
	"math/big"

	"github.com/0xPolygon/polygon-cli/bindings/multicall3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

const Multicall3Addr = "0xcA11bde05977b3631167028862bE2a173976CA11"

// Contract CALL (Multicall3) must pay all CALL-path costs, including:
// New account surcharge: 25,000 gas when value creates a previously empty account.
// EIP Library
// Cold account access: 2,600 gas on first touch (EIP-2929).
// Ethereum Improvement Proposals
// +1
// CALL base ~700 and value transfer cost 9,000. (Standard CALL metering.)
// GitHub
// Net: for brand-new recipients, a Multicall3 transfer is ~37k gas each (≈25k + 9k + 2.6k + 0.7k)
const gasToFundAccountAnAccount = 40000 // estimated gas per account to fund with multicall3
const maxAccsToFundPerTx = 700          // arbitrary limit to avoid too large transactions

func Multicall3Deploy(c *ethclient.Client, sender *bind.TransactOpts) (common.Address, *types.Transaction, *multicall3.Multicall3, error) {
	address, tx, instance, err := multicall3.DeployMulticall3(sender, c)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, instance, nil
}

func Multicall3Exists(ctx context.Context, c *ethclient.Client, customAddr *common.Address) (bool, error) {
	scAddr := customAddr
	if scAddr == nil {
		addr := common.HexToAddress(Multicall3Addr)
		scAddr = &addr
	}

	code, err := c.CodeAt(ctx, *scAddr, nil)
	if err != nil {
		return false, err
	}

	if len(code) == 0 {
		return false, nil
	}

	sc, err := multicall3.NewMulticall3(*scAddr, c)
	if err != nil {
		return false, err
	}

	_, err = sc.GetBlockNumber(&bind.CallOpts{Context: ctx})
	if err != nil {
		return false, err
	}

	return true, nil
}

func Multicall3New(c *ethclient.Client, customAddr *common.Address) (*multicall3.Multicall3, error) {
	scAddr := customAddr
	if scAddr == nil {
		addr := common.HexToAddress(Multicall3Addr)
		scAddr = &addr
	}

	sc, err := multicall3.NewMulticall3(*scAddr, c)
	if err != nil {
		return nil, err
	}

	return sc, nil
}

func IsMulticall3Supported(ctx context.Context, c *ethclient.Client, tryDeployIfNotExist bool, tops *bind.TransactOpts, customAddr *common.Address) (*common.Address, bool) {
	exists, err := Multicall3Exists(ctx, c, customAddr)
	if err != nil {
		return nil, false
	}

	if exists {
		if customAddr != nil {
			return customAddr, true
		}
		addr := common.HexToAddress(Multicall3Addr)
		return &addr, true
	}

	if !tryDeployIfNotExist {
		return nil, false
	}

	if tops == nil {
		return nil, false
	}

	addr, tx, _, err := Multicall3Deploy(c, tops)
	if err != nil {
		return nil, false
	}

	receipt, err := bind.WaitMined(ctx, c, tx.Hash())
	if err != nil || receipt == nil || receipt.Status != 1 {
		return nil, false
	}

	return &addr, true
}

func Multicall3MaxAccountsToFundPerTx(ctx context.Context, c *ethclient.Client) (uint64, error) {
	latestBlock, err := c.BlockByNumber(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to get block gas limit")
		return 0, err
	}
	return min(latestBlock.GasLimit()/gasToFundAccountAnAccount, maxAccsToFundPerTx), nil
}

func Multicall3FundAccountsWithNativeToken(c *ethclient.Client, sender *bind.TransactOpts, accounts []common.Address, amount *big.Int, customAddr *common.Address) (*types.Transaction, error) {
	sc, err := Multicall3New(c, customAddr)
	if err != nil {
		return nil, err
	}

	calls := make([]multicall3.Multicall3Call3Value, 0, len(accounts))
	for _, account := range accounts {
		calls = append(calls, multicall3.Multicall3Call3Value{
			Target:       account,
			AllowFailure: false,
			Value:        amount,
		})
	}

	sender.Value = big.NewInt(0).Mul(amount, big.NewInt(int64(len(accounts))))

	return sc.Aggregate3Value(sender, calls)
}
