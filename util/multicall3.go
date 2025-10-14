package util

import (
	"context"
	"math/big"
	"strings"

	"github.com/0xPolygon/polygon-cli/bindings/multicall3"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

const DefaultMulticall3Addr = "0xcA11bde05977b3631167028862bE2a173976CA11"

// Contract CALL (Multicall3) must pay all CALL-path costs, including:
// New account surcharge: 25,000 gas when value creates a previously empty account.
// EIP Library
// Cold account access: 2,600 gas on first touch (EIP-2929).
// Ethereum Improvement Proposals
// +1
// CALL base ~700 and value transfer cost 9,000. (Standard CALL metering.)
// GitHub
// Net: for brand-new recipients, a Multicall3 transfer is ~37k gas each (≈25k + 9k + 2.6k + 0.7k)
const estimatedGasNeededToFundASingleAccount = 40000 // estimated gas per account to fund with multicall3 + margin(3k)
const maxAccsToFundPerTx = 700                       // arbitrary limit to avoid too large transactions

func Multicall3Deploy(c *ethclient.Client, sender *bind.TransactOpts) (common.Address, *types.Transaction, *multicall3.Multicall3, error) {
	address, tx, instance, err := multicall3.DeployMulticall3(sender, c)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, instance, nil
}

func Multicall3Exists(ctx context.Context, c *ethclient.Client, customAddr *common.Address) (common.Address, bool, error) {
	scAddr := customAddr
	if scAddr == nil {
		addr := common.HexToAddress(DefaultMulticall3Addr)
		scAddr = &addr
	}

	code, err := c.CodeAt(ctx, *scAddr, nil)
	if err != nil {
		return *scAddr, false, err
	}

	if len(code) == 0 {
		return *scAddr, false, nil
	}

	sc, err := multicall3.NewMulticall3(*scAddr, c)
	if err != nil {
		return *scAddr, false, err
	}

	_, err = sc.GetBlockNumber(&bind.CallOpts{Context: ctx})
	if err != nil {
		return *scAddr, false, err
	}

	return *scAddr, true, nil
}

func Multicall3New(c *ethclient.Client, customAddr *common.Address) (common.Address, *multicall3.Multicall3, error) {
	scAddr := customAddr
	if scAddr == nil {
		addr := common.HexToAddress(DefaultMulticall3Addr)
		scAddr = &addr
	}

	sc, err := multicall3.NewMulticall3(*scAddr, c)
	if err != nil {
		return common.Address{}, nil, err
	}

	return *scAddr, sc, nil
}

func IsMulticall3Supported(ctx context.Context, c *ethclient.Client, tryDeployIfNotExist bool, tops *bind.TransactOpts, customAddr *common.Address) (*common.Address, bool) {
	scAddr, exists, err := Multicall3Exists(ctx, c, customAddr)
	if err != nil {
		return nil, false
	}

	if exists {
		return &scAddr, true
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
	latestBlock, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to get block gas limit")
		return 0, err
	}
	return min(latestBlock.GasLimit/estimatedGasNeededToFundASingleAccount, maxAccsToFundPerTx), nil
}

func Multicall3FundAccountsWithNativeToken(c *ethclient.Client, tops *bind.TransactOpts, accounts []common.Address, amount *big.Int, customAddr *common.Address) (*types.Transaction, error) {
	_, sc, err := Multicall3New(c, customAddr)
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

	tops.Value = big.NewInt(0).Mul(amount, big.NewInt(int64(len(accounts))))

	return sc.Aggregate3Value(tops, calls)
}

func Multicall3MintERC20ToAccounts(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, accounts []common.Address, tokenAddress common.Address, amount *big.Int, customAddr *common.Address) (*types.Transaction, error) {
	_, sc, err := Multicall3New(c, customAddr)
	if err != nil {
		return nil, err
	}

	// Create ABI for mint(address, uint256) function
	mintABI, err := abi.JSON(strings.NewReader(`[{"type":"function","name":"mint","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[],"stateMutability":"nonpayable"}]`))
	if err != nil {
		return nil, err
	}

	calls := make([]multicall3.Multicall3Call3, 0, len(accounts))
	for _, account := range accounts {
		callData, iErr := mintABI.Pack("mint", account, amount)
		if iErr != nil {
			return nil, iErr
		}

		calls = append(calls, multicall3.Multicall3Call3{
			Target:       tokenAddress,
			AllowFailure: false,
			CallData:     callData,
		})
	}

	return sc.Aggregate3(tops, calls)
}
