package p2p

import (
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const testChainID = 137

// testHead returns a head header resembling a Polygon block: post-London with a
// base fee, a 30M gas limit, and a non-zero difficulty.
func testHead() *types.Header {
	return &types.Header{
		Number:     big.NewInt(100),
		GasLimit:   30_000_000,
		Time:       1_700_000_000,
		BaseFee:    big.NewInt(25_000_000_000),
		Difficulty: big.NewInt(1),
	}
}

// signTx signs tx for chainID with a fresh key.
func signTx(t *testing.T, chainID uint64, inner types.TxData) *types.Transaction {
	t.Helper()

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	signer := types.LatestSigner(broadcastChainConfig(chainID))
	tx, err := types.SignNewTx(key, signer, inner)
	if err != nil {
		t.Fatalf("sign tx: %v", err)
	}
	return tx
}

// dynamicFeeTx builds a well-formed transfer that passes every check.
func dynamicFeeTx(t *testing.T) types.TxData {
	t.Helper()

	return dynamicFeeTxFor(t, testChainID)
}

// dynamicFeeTxFor builds a well-formed transfer for an arbitrary chain.
func dynamicFeeTxFor(t *testing.T, chainID uint64) types.TxData {
	t.Helper()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	return &types.DynamicFeeTx{
		ChainID:   new(big.Int).SetUint64(chainID),
		Nonce:     1,
		GasTipCap: big.NewInt(30_000_000_000),
		GasFeeCap: big.NewInt(60_000_000_000),
		Gas:       21_000,
		To:        &to,
		Value:     big.NewInt(1),
	}
}

func newTestValidator(t *testing.T, opts TxValidatorOptions) *TxValidator {
	t.Helper()

	if opts.ChainID == 0 {
		opts.ChainID = testChainID
	}
	v, err := NewTxValidator(opts)
	if err != nil {
		t.Fatalf("new validator: %v", err)
	}
	return v
}

func TestNewTxValidatorRejectsBadOptions(t *testing.T) {
	tests := []struct {
		name string
		opts TxValidatorOptions
	}{
		{"zero chain ID", TxValidatorOptions{}},
		{"negative gas price", TxValidatorOptions{ChainID: testChainID, MinGasPrice: big.NewInt(-1)}},
		{"negative tip", TxValidatorOptions{ChainID: testChainID, MinTip: big.NewInt(-1)}},
		{"negative ratio", TxValidatorOptions{ChainID: testChainID, MinBaseFeeRatio: -1}},
		{"ratio too large", TxValidatorOptions{ChainID: testChainID, MinBaseFeeRatio: maxBaseFeeRatio + 1}},
		{"NaN ratio", TxValidatorOptions{ChainID: testChainID, MinBaseFeeRatio: math.NaN()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewTxValidator(tt.opts); err == nil {
				t.Fatal("want error, got nil")
			}
		})
	}
}

func TestValidateAcceptsWellFormedTxs(t *testing.T) {
	v := newTestValidator(t, TxValidatorOptions{MinBaseFeeRatio: 1})
	head := testHead()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	tests := []struct {
		name  string
		inner types.TxData
	}{
		{"dynamic fee", dynamicFeeTx(t)},
		{"legacy", &types.LegacyTx{
			Nonce:    1,
			GasPrice: big.NewInt(60_000_000_000),
			Gas:      21_000,
			To:       &to,
			Value:    big.NewInt(1),
		}},
		{"access list", &types.AccessListTx{
			ChainID:  big.NewInt(testChainID),
			Nonce:    1,
			GasPrice: big.NewInt(60_000_000_000),
			Gas:      21_000,
			To:       &to,
			Value:    big.NewInt(1),
		}},
		{"contract creation", &types.DynamicFeeTx{
			ChainID:   big.NewInt(testChainID),
			Nonce:     1,
			GasTipCap: big.NewInt(30_000_000_000),
			GasFeeCap: big.NewInt(60_000_000_000),
			Gas:       1_000_000,
			Data:      make([]byte, 1024),
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Validate(signTx(t, testChainID, tt.inner), head); err != nil {
				t.Fatalf("want accepted, got %v", err)
			}
		})
	}
}

func TestValidateRejects(t *testing.T) {
	to := common.HexToAddress("0x1234567890123456789012345678901234567890")

	tests := []struct {
		name    string
		opts    TxValidatorOptions
		chainID uint64
		inner   types.TxData
		reason  string
	}{
		{
			name:    "wrong chain ID",
			chainID: testChainID + 1,
			inner: &types.DynamicFeeTx{
				ChainID:   big.NewInt(testChainID + 1),
				Nonce:     1,
				GasTipCap: big.NewInt(30_000_000_000),
				GasFeeCap: big.NewInt(60_000_000_000),
				Gas:       21_000,
				To:        &to,
				Value:     big.NewInt(1),
			},
			reason: "invalid_signature",
		},
		{
			name: "gas below intrinsic",
			inner: &types.DynamicFeeTx{
				ChainID:   big.NewInt(testChainID),
				Nonce:     1,
				GasTipCap: big.NewInt(30_000_000_000),
				GasFeeCap: big.NewInt(60_000_000_000),
				Gas:       20_999,
				To:        &to,
				Value:     big.NewInt(1),
			},
			reason: "intrinsic_gas",
		},
		{
			name: "gas above block limit",
			inner: &types.DynamicFeeTx{
				ChainID:   big.NewInt(testChainID),
				Nonce:     1,
				GasTipCap: big.NewInt(30_000_000_000),
				GasFeeCap: big.NewInt(60_000_000_000),
				Gas:       30_000_001,
				To:        &to,
				Value:     big.NewInt(1),
			},
			reason: "gas_limit",
		},
		{
			name: "tip above fee cap",
			inner: &types.DynamicFeeTx{
				ChainID:   big.NewInt(testChainID),
				Nonce:     1,
				GasTipCap: big.NewInt(90_000_000_000),
				GasFeeCap: big.NewInt(60_000_000_000),
				Gas:       21_000,
				To:        &to,
				Value:     big.NewInt(1),
			},
			reason: "tip_above_fee_cap",
		},
		{
			name: "oversized",
			inner: &types.DynamicFeeTx{
				ChainID:   big.NewInt(testChainID),
				Nonce:     1,
				GasTipCap: big.NewInt(30_000_000_000),
				GasFeeCap: big.NewInt(60_000_000_000),
				Gas:       25_000_000,
				To:        &to,
				Data:      make([]byte, maxBroadcastTxSize+1),
			},
			reason: "oversized",
		},
		{
			name: "tip below floor",
			opts: TxValidatorOptions{MinTip: big.NewInt(30_000_000_001)},
			inner: &types.DynamicFeeTx{
				ChainID:   big.NewInt(testChainID),
				Nonce:     1,
				GasTipCap: big.NewInt(30_000_000_000),
				GasFeeCap: big.NewInt(60_000_000_000),
				Gas:       21_000,
				To:        &to,
				Value:     big.NewInt(1),
			},
			reason: "tip_too_low",
		},
		{
			name:   "fee cap below minimum",
			opts:   TxValidatorOptions{MinGasPrice: big.NewInt(60_000_000_001)},
			inner:  dynamicFeeTx(t),
			reason: "below_min_gas_price",
		},
		{
			name: "fee cap below base fee",
			opts: TxValidatorOptions{MinBaseFeeRatio: 1},
			inner: &types.DynamicFeeTx{
				ChainID:   big.NewInt(testChainID),
				Nonce:     1,
				GasTipCap: big.NewInt(1),
				GasFeeCap: big.NewInt(24_999_999_999),
				Gas:       21_000,
				To:        &to,
				Value:     big.NewInt(1),
			},
			reason: "fee_cap_below_basefee",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := newTestValidator(t, tt.opts)
			chainID := uint64(testChainID)
			if tt.chainID != 0 {
				chainID = tt.chainID
			}

			err := v.Validate(signTx(t, chainID, tt.inner), testHead())
			if err == nil {
				t.Fatal("want rejection, got nil")
			}
			if got := RejectReason(err); got != tt.reason {
				t.Fatalf("want reason %q, got %q (%v)", tt.reason, got, err)
			}
		})
	}
}

// TestValidateBaseFeeRatio checks the fractional floor, which lets an operator
// keep forwarding transactions that are slightly under the current base fee.
func TestValidateBaseFeeRatio(t *testing.T) {
	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	// Head base fee is 25 gwei, so a 0.5 ratio puts the floor at 12.5 gwei.
	tx := signTx(t, testChainID, &types.DynamicFeeTx{
		ChainID:   big.NewInt(testChainID),
		Nonce:     1,
		GasTipCap: big.NewInt(1),
		GasFeeCap: big.NewInt(13_000_000_000),
		Gas:       21_000,
		To:        &to,
		Value:     big.NewInt(1),
	})

	if err := newTestValidator(t, TxValidatorOptions{MinBaseFeeRatio: 0.5}).Validate(tx, testHead()); err != nil {
		t.Fatalf("ratio 0.5: want accepted, got %v", err)
	}
	if err := newTestValidator(t, TxValidatorOptions{MinBaseFeeRatio: 1}).Validate(tx, testHead()); err == nil {
		t.Fatal("ratio 1: want rejection, got nil")
	}
	if err := newTestValidator(t, TxValidatorOptions{MinBaseFeeRatio: 0}).Validate(tx, testHead()); err != nil {
		t.Fatalf("ratio 0 (disabled): want accepted, got %v", err)
	}
}

// TestValidateMissingHeadFields covers headers that reach the validator without
// the fields go-ethereum dereferences unconditionally. A nil difficulty used to
// be the kind of input that panics on a peer's protocol goroutine.
func TestValidateMissingHeadFields(t *testing.T) {
	v := newTestValidator(t, TxValidatorOptions{MinBaseFeeRatio: 1})
	tx := signTx(t, testChainID, dynamicFeeTx(t))

	head := testHead()
	head.Difficulty = nil
	if err := v.Validate(tx, head); err != nil {
		t.Fatalf("nil difficulty: want accepted, got %v", err)
	}
	if head.Difficulty != nil {
		t.Fatal("validator mutated the caller's header")
	}

	head = testHead()
	head.BaseFee = nil
	if err := v.Validate(tx, head); err != nil {
		t.Fatalf("nil base fee: want accepted, got %v", err)
	}
}

func TestValidateNilInputs(t *testing.T) {
	v := newTestValidator(t, TxValidatorOptions{})

	if err := v.Validate(nil, testHead()); err == nil {
		t.Fatal("nil tx: want error, got nil")
	}
	if err := v.Validate(signTx(t, testChainID, dynamicFeeTx(t)), nil); err == nil {
		t.Fatal("nil head: want error, got nil")
	}
}

func TestRejectReasonUnknownErrors(t *testing.T) {
	if got := RejectReason(nil); got != "" {
		t.Fatalf("nil error: want empty reason, got %q", got)
	}
	if got := RejectReason(errBelowMinGasPrice); got != "below_min_gas_price" {
		t.Fatalf("want below_min_gas_price, got %q", got)
	}
	if got := RejectReason(types.ErrInvalidSig); got != "other" {
		t.Fatalf("unmapped error: want other, got %q", got)
	}
}
