package modes

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"
	"testing"

	"github.com/0xPolygon/polygon-cli/loadtest/config"
	"github.com/0xPolygon/polygon-cli/loadtest/mode"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func testTransactor(t *testing.T, chainID uint64) (*bind.TransactOpts, *ecdsa.PrivateKey) {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	tops, err := bind.NewKeyedTransactorWithChainID(key, new(big.Int).SetUint64(chainID))
	if err != nil {
		t.Fatal(err)
	}
	tops.Nonce = big.NewInt(7)
	tops.GasLimit = 100_000
	tops.GasFeeCap = big.NewInt(30)
	tops.GasTipCap = big.NewInt(2)
	tops.GasPrice = big.NewInt(25)
	return tops, key
}

func TestBuildContractCallTx(t *testing.T) {
	to := common.HexToAddress("0x1111111111111111111111111111111111111111")
	calldata := []byte{0xde, 0xad, 0xbe, 0xef}
	amount := big.NewInt(42)

	t.Run("dynamic fee", func(t *testing.T) {
		cfg := &config.Config{ChainID: 1337}
		tops, _ := testTransactor(t, cfg.ChainID)
		tx := buildContractCallTx(cfg, tops, &to, amount, calldata)
		if tx.Type() != types.DynamicFeeTxType {
			t.Fatalf("tx type = %d, want dynamic fee", tx.Type())
		}
		if tx.ChainId().Uint64() != 1337 || tx.Nonce() != 7 || tx.Gas() != 100_000 {
			t.Fatalf("unexpected chainID/nonce/gas: %d/%d/%d", tx.ChainId(), tx.Nonce(), tx.Gas())
		}
		if tx.GasFeeCap().Int64() != 30 || tx.GasTipCap().Int64() != 2 {
			t.Fatalf("unexpected fee caps: %s/%s", tx.GasFeeCap(), tx.GasTipCap())
		}
		if *tx.To() != to || tx.Value().Cmp(amount) != 0 || !bytes.Equal(tx.Data(), calldata) {
			t.Fatalf("unexpected to/value/data: %s/%s/%x", tx.To(), tx.Value(), tx.Data())
		}
	})

	t.Run("legacy", func(t *testing.T) {
		cfg := &config.Config{ChainID: 1337, LegacyTxMode: true}
		tops, _ := testTransactor(t, cfg.ChainID)
		tx := buildContractCallTx(cfg, tops, &to, amount, calldata)
		if tx.Type() != types.LegacyTxType {
			t.Fatalf("tx type = %d, want legacy", tx.Type())
		}
		if tx.GasPrice().Int64() != 25 || tx.Nonce() != 7 || !bytes.Equal(tx.Data(), calldata) {
			t.Fatalf("unexpected gasPrice/nonce/data: %s/%d/%x", tx.GasPrice(), tx.Nonce(), tx.Data())
		}
	})
}

// TestExecuteWithStdinCalldata drives Execute with --output-raw-tx-only so no
// RPC client is needed, and checks that each transaction carries the next
// chunk from the feeder and that the run ends with ErrInputExhausted.
func TestExecuteWithStdinCalldata(t *testing.T) {
	const size = 4
	input := []byte("aaaabbbbcccc")
	to := common.HexToAddress("0x1111111111111111111111111111111111111111")
	cfg := &config.Config{
		ChainID:               1337,
		Concurrency:           1,
		OutputRawTxOnly:       true,
		ContractETHAddress:    &to,
		ContractCallDataStdin: true,
		ContractCallDataSize:  size,
	}

	// Silence the raw tx hex that Execute prints to stdout.
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	origStdout := os.Stdout
	os.Stdout = devNull
	t.Cleanup(func() {
		os.Stdout = origStdout
		_ = devNull.Close()
	})

	m := &ContractCallMode{}
	if err = m.initFeeder(t.Context(), bytes.NewReader(input), cfg); err != nil {
		t.Fatal(err)
	}

	tops, key := testTransactor(t, cfg.ChainID)
	signer := types.LatestSignerForChainID(new(big.Int).SetUint64(cfg.ChainID))
	deps := &mode.Dependencies{}

	var hashes []common.Hash
	for i := 0; i < len(input)/size; i++ {
		want := input[i*size : (i+1)*size]
		reserved, reserveErr := m.ReserveInput(t.Context())
		if reserveErr != nil {
			t.Fatalf("ReserveInput() %d error: %v", i, reserveErr)
		}
		if !bytes.Equal(reserved.([]byte), want) {
			t.Fatalf("ReserveInput() %d = %q, want %q", i, reserved, want)
		}
		_, _, txHash, execErr := m.Execute(mode.WithInput(t.Context(), reserved), cfg, deps, tops)
		if execErr != nil {
			t.Fatalf("Execute() %d error: %v", i, execErr)
		}
		// Rebuild the expected signed tx from the same inputs. A matching hash
		// proves the chunk ended up as the calldata.
		expected, signErr := types.SignTx(buildContractCallTx(cfg, tops, &to, big.NewInt(0), want), signer, key)
		if signErr != nil {
			t.Fatal(signErr)
		}
		if txHash != expected.Hash() {
			t.Fatalf("Execute() %d hash = %s, want %s (calldata %q)", i, txHash, expected.Hash(), want)
		}
		hashes = append(hashes, txHash)
	}
	for i := 1; i < len(hashes); i++ {
		if hashes[i] == hashes[i-1] {
			t.Fatalf("consecutive transactions %d and %d have identical hashes", i-1, i)
		}
	}

	_, reserveErr := m.ReserveInput(t.Context())
	if !errors.Is(reserveErr, mode.ErrInputExhausted) {
		t.Fatalf("ReserveInput() after EOF = %v, want ErrInputExhausted", reserveErr)
	}

	// Execute without a reservation must fail rather than block on stdin or
	// silently reuse a payload.
	if _, _, _, execErr := m.Execute(t.Context(), cfg, deps, tops); execErr == nil {
		t.Fatal("Execute() without a reserved chunk succeeded, want error")
	}
}

func TestReserveInputWithoutStdinIsNoop(t *testing.T) {
	m := &ContractCallMode{}
	input, err := m.ReserveInput(t.Context())
	if err != nil || input != nil {
		t.Fatalf("ReserveInput() = %v, %v; want nil, nil", input, err)
	}
}

func TestInitClearsStaleFeeder(t *testing.T) {
	m := &ContractCallMode{}
	stdinCfg := &config.Config{Concurrency: 1, ContractCallDataStdin: true, ContractCallDataSize: 4}
	if err := m.initFeeder(t.Context(), bytes.NewReader(nil), stdinCfg); err != nil {
		t.Fatal(err)
	}
	if m.feeder == nil {
		t.Fatal("feeder not set")
	}
	if err := m.Init(t.Context(), &config.Config{}, &mode.Dependencies{}); err != nil {
		t.Fatal(err)
	}
	if m.feeder != nil {
		t.Fatal("Init without --calldata-stdin left a stale feeder in place")
	}
}

func TestExecuteWithoutCalldataFails(t *testing.T) {
	to := common.HexToAddress("0x1111111111111111111111111111111111111111")
	cfg := &config.Config{ChainID: 1337, OutputRawTxOnly: true, ContractETHAddress: &to}
	tops, _ := testTransactor(t, cfg.ChainID)
	m := &ContractCallMode{}
	if _, _, _, err := m.Execute(t.Context(), cfg, &mode.Dependencies{}, tops); err == nil {
		t.Fatal("expected error when no calldata source is configured")
	}
}
