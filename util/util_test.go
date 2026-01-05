package util

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/0xPolygon/polygon-cli/rpctypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestGetSenderFromTx tests the sender recovery logic for all transaction types
func TestGetSenderFromTx(t *testing.T) {
	ctx := context.Background()

	// Generate a test private key and derive the expected address
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}
	expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	chainID := uint64(1337)

	tests := []struct {
		name        string
		txType      uint8
		chainID     uint64
		privateKey  *ecdsa.PrivateKey
		wantAddress common.Address
	}{
		{
			name:        "Legacy transaction with EIP-155",
			txType:      0,
			chainID:     chainID,
			privateKey:  privateKey,
			wantAddress: expectedAddress,
		},
		{
			name:        "EIP-2930 transaction",
			txType:      1,
			chainID:     chainID,
			privateKey:  privateKey,
			wantAddress: expectedAddress,
		},
		{
			name:        "EIP-1559 transaction",
			txType:      2,
			chainID:     chainID,
			privateKey:  privateKey,
			wantAddress: expectedAddress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a signed transaction based on type
			var tx *types.Transaction
			var err error

			to := common.HexToAddress("0x1234567890123456789012345678901234567890")
			value := big.NewInt(1000)
			gasLimit := uint64(21000)
			nonce := uint64(0)
			data := []byte{}

			switch tt.txType {
			case 0: // Legacy
				gasPrice := big.NewInt(1000000000)
				tx = types.NewTx(&types.LegacyTx{
					Nonce:    nonce,
					GasPrice: gasPrice,
					Gas:      gasLimit,
					To:       &to,
					Value:    value,
					Data:     data,
				})

			case 1: // EIP-2930
				gasPrice := big.NewInt(1000000000)
				tx = types.NewTx(&types.AccessListTx{
					ChainID:    big.NewInt(int64(tt.chainID)),
					Nonce:      nonce,
					GasPrice:   gasPrice,
					Gas:        gasLimit,
					To:         &to,
					Value:      value,
					Data:       data,
					AccessList: types.AccessList{},
				})

			case 2: // EIP-1559
				gasTipCap := big.NewInt(1000000000)
				gasFeeCap := big.NewInt(2000000000)
				tx = types.NewTx(&types.DynamicFeeTx{
					ChainID:    big.NewInt(int64(tt.chainID)),
					Nonce:      nonce,
					GasTipCap:  gasTipCap,
					GasFeeCap:  gasFeeCap,
					Gas:        gasLimit,
					To:         &to,
					Value:      value,
					Data:       data,
					AccessList: types.AccessList{},
				})
			}

			// Sign the transaction
			signer := types.LatestSignerForChainID(big.NewInt(int64(tt.chainID)))
			signedTx, err := types.SignTx(tx, signer, tt.privateKey)
			if err != nil {
				t.Fatalf("failed to sign transaction: %v", err)
			}

			// Convert to PolyTransaction for testing
			polyTx := createPolyTransactionFromSignedTx(t, signedTx)

			// Test GetSenderFromTx
			recoveredAddress, err := GetSenderFromTx(ctx, polyTx)
			if err != nil {
				t.Fatalf("GetSenderFromTx failed: %v", err)
			}

			if recoveredAddress != tt.wantAddress {
				t.Errorf("address mismatch: got %s, want %s", recoveredAddress.Hex(), tt.wantAddress.Hex())
			}
		})
	}
}

// TestGetSenderFromTx_PreEIP155Legacy tests legacy transactions without EIP-155 (chainID = 0)
func TestGetSenderFromTx_PreEIP155Legacy(t *testing.T) {
	t.Skip("Pre-EIP-155 transactions require different signing logic and are rarely used")
	// Note: This would require special handling as go-ethereum's SignTx doesn't easily support pre-EIP-155
}

// TestGetSenderFromTx_InvalidSignature tests error handling for invalid signatures
func TestGetSenderFromTx_InvalidSignature(t *testing.T) {
	ctx := context.Background()

	// Create a transaction with an invalid signature
	to := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// Create a mock PolyTransaction with invalid signature values
	mockTx := &mockPolyTransaction{
		txType:   0,
		chainID:  1,
		nonce:    0,
		gasPrice: big.NewInt(1000000000),
		gas:      21000,
		to:       to,
		value:    big.NewInt(1000),
		data:     []byte{},
		v:        big.NewInt(27), // Invalid v, r, s combination
		r:        big.NewInt(0),
		s:        big.NewInt(0),
	}

	_, err := GetSenderFromTx(ctx, mockTx)
	if err == nil {
		t.Error("expected error for invalid signature, got nil")
	}
}

// TestGetSenderFromTx_ContractCreation tests transactions with no 'to' address (contract creation)
func TestGetSenderFromTx_ContractCreation(t *testing.T) {
	ctx := context.Background()

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}
	expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	chainID := uint64(1337)

	// Create a contract creation transaction (to = nil)
	gasPrice := big.NewInt(1000000000)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		GasPrice: gasPrice,
		Gas:      100000,
		To:       nil, // Contract creation
		Value:    big.NewInt(0),
		Data:     []byte{0x60, 0x60, 0x60}, // Some bytecode
	})

	signer := types.LatestSignerForChainID(big.NewInt(int64(chainID)))
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}

	polyTx := createPolyTransactionFromSignedTx(t, signedTx)

	recoveredAddress, err := GetSenderFromTx(ctx, polyTx)
	if err != nil {
		t.Fatalf("GetSenderFromTx failed for contract creation: %v", err)
	}

	if recoveredAddress != expectedAddress {
		t.Errorf("address mismatch: got %s, want %s", recoveredAddress.Hex(), expectedAddress.Hex())
	}
}

// TestGetSenderFromTx_UnsupportedType tests handling of unsupported transaction types
func TestGetSenderFromTx_UnsupportedType(t *testing.T) {
	ctx := context.Background()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// Create a mock transaction with an unsupported type (e.g., type 3)
	mockTx := &mockPolyTransaction{
		txType:   3, // Unsupported type
		chainID:  1,
		from:     common.HexToAddress("0x0000000000000000000000000000000000000000"),
		nonce:    0,
		gasPrice: big.NewInt(1000000000),
		gas:      21000,
		to:       to,
		value:    big.NewInt(1000),
		data:     []byte{},
	}

	// For type > 2, GetSenderFromTx should return the 'from' field directly
	recoveredAddress, err := GetSenderFromTx(ctx, mockTx)
	if err != nil {
		t.Fatalf("GetSenderFromTx failed: %v", err)
	}

	if recoveredAddress != mockTx.from {
		t.Errorf("expected to use 'from' field for unsupported type, got %s, want %s",
			recoveredAddress.Hex(), mockTx.from.Hex())
	}
}

// TestGetSenderFromTx_DifferentChainIDs tests sender recovery with various chain IDs
func TestGetSenderFromTx_DifferentChainIDs(t *testing.T) {
	ctx := context.Background()

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}
	expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	chainIDs := []uint64{1, 137, 1337, 80001, 100000}

	for _, chainID := range chainIDs {
		t.Run(fmt.Sprintf("ChainID_%d", chainID), func(t *testing.T) {
			to := common.HexToAddress("0x1234567890123456789012345678901234567890")
			value := big.NewInt(1000)
			gasLimit := uint64(21000)
			gasPrice := big.NewInt(1000000000)

			tx := types.NewTx(&types.LegacyTx{
				Nonce:    0,
				GasPrice: gasPrice,
				Gas:      gasLimit,
				To:       &to,
				Value:    value,
				Data:     []byte{},
			})

			signer := types.LatestSignerForChainID(big.NewInt(int64(chainID)))
			signedTx, err := types.SignTx(tx, signer, privateKey)
			if err != nil {
				t.Fatalf("failed to sign transaction: %v", err)
			}

			polyTx := createPolyTransactionFromSignedTx(t, signedTx)

			recoveredAddress, err := GetSenderFromTx(ctx, polyTx)
			if err != nil {
				t.Fatalf("GetSenderFromTx failed for chainID %d: %v", chainID, err)
			}

			if recoveredAddress != expectedAddress {
				t.Errorf("address mismatch for chainID %d: got %s, want %s",
					chainID, recoveredAddress.Hex(), expectedAddress.Hex())
			}
		})
	}
}

// TestGetSenderFromTx_WithData tests transactions with various data payloads
func TestGetSenderFromTx_WithData(t *testing.T) {
	ctx := context.Background()

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}
	expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	chainID := uint64(1337)

	testCases := []struct {
		name string
		data []byte
	}{
		{
			name: "empty data",
			data: []byte{},
		},
		{
			name: "small data",
			data: []byte{0x01, 0x02, 0x03},
		},
		{
			name: "function call data",
			data: []byte{0xa9, 0x05, 0x9c, 0xbb}, // transfer(address,uint256) selector
		},
		{
			name: "large data",
			data: make([]byte, 1024),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			to := common.HexToAddress("0x1234567890123456789012345678901234567890")
			value := big.NewInt(0)
			gasLimit := uint64(100000)
			gasTipCap := big.NewInt(1000000000)
			gasFeeCap := big.NewInt(2000000000)

			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:    big.NewInt(int64(chainID)),
				Nonce:      0,
				GasTipCap:  gasTipCap,
				GasFeeCap:  gasFeeCap,
				Gas:        gasLimit,
				To:         &to,
				Value:      value,
				Data:       tc.data,
				AccessList: types.AccessList{},
			})

			signer := types.LatestSignerForChainID(big.NewInt(int64(chainID)))
			signedTx, err := types.SignTx(tx, signer, privateKey)
			if err != nil {
				t.Fatalf("failed to sign transaction: %v", err)
			}

			polyTx := createPolyTransactionFromSignedTx(t, signedTx)

			recoveredAddress, err := GetSenderFromTx(ctx, polyTx)
			if err != nil {
				t.Fatalf("GetSenderFromTx failed: %v", err)
			}

			if recoveredAddress != expectedAddress {
				t.Errorf("address mismatch: got %s, want %s", recoveredAddress.Hex(), expectedAddress.Hex())
			}
		})
	}
}

// createPolyTransactionFromSignedTx converts a signed types.Transaction to a mock PolyTransaction
func createPolyTransactionFromSignedTx(t *testing.T, signedTx *types.Transaction) rpctypes.PolyTransaction {
	t.Helper()

	v, r, s := signedTx.RawSignatureValues()

	var to common.Address
	if signedTx.To() != nil {
		to = *signedTx.To()
	}

	mockTx := &mockPolyTransaction{
		txType:  uint64(signedTx.Type()),
		chainID: signedTx.ChainId().Uint64(),
		nonce:   signedTx.Nonce(),
		gas:     signedTx.Gas(),
		to:      to,
		value:   signedTx.Value(),
		data:    signedTx.Data(),
		v:       v,
		r:       r,
		s:       s,
		hash:    signedTx.Hash(),
	}

	// Set price fields based on transaction type
	switch signedTx.Type() {
	case 0, 1: // Legacy and EIP-2930
		mockTx.gasPrice = signedTx.GasPrice()
	case 2: // EIP-1559
		mockTx.maxPriorityFeePerGas = signedTx.GasTipCap().Uint64()
		mockTx.maxFeePerGas = signedTx.GasFeeCap().Uint64()
		// For testing, use effective gas price
		mockTx.gasPrice = signedTx.GasPrice()
	}

	return mockTx
}

// mockPolyTransaction implements rpctypes.PolyTransaction for testing
type mockPolyTransaction struct {
	txType               uint64
	chainID              uint64
	from                 common.Address
	nonce                uint64
	gasPrice             *big.Int
	maxPriorityFeePerGas uint64
	maxFeePerGas         uint64
	gas                  uint64
	to                   common.Address
	value                *big.Int
	data                 []byte
	v, r, s              *big.Int
	hash                 common.Hash
	blockNumber          *big.Int
}

func (m *mockPolyTransaction) Type() uint64                 { return m.txType }
func (m *mockPolyTransaction) ChainID() uint64              { return m.chainID }
func (m *mockPolyTransaction) From() common.Address         { return m.from }
func (m *mockPolyTransaction) Nonce() uint64                { return m.nonce }
func (m *mockPolyTransaction) GasPrice() *big.Int           { return m.gasPrice }
func (m *mockPolyTransaction) MaxPriorityFeePerGas() uint64 { return m.maxPriorityFeePerGas }
func (m *mockPolyTransaction) MaxFeePerGas() uint64         { return m.maxFeePerGas }
func (m *mockPolyTransaction) Gas() uint64                  { return m.gas }
func (m *mockPolyTransaction) To() common.Address           { return m.to }
func (m *mockPolyTransaction) Value() *big.Int              { return m.value }
func (m *mockPolyTransaction) Data() []byte                 { return m.data }
func (m *mockPolyTransaction) DataStr() string {
	return "0x" + common.Bytes2Hex(m.data)
}
func (m *mockPolyTransaction) V() *big.Int       { return m.v }
func (m *mockPolyTransaction) R() *big.Int       { return m.r }
func (m *mockPolyTransaction) S() *big.Int       { return m.s }
func (m *mockPolyTransaction) Hash() common.Hash { return m.hash }
func (m *mockPolyTransaction) BlockNumber() *big.Int {
	if m.blockNumber != nil {
		return m.blockNumber
	}
	return big.NewInt(0)
}
func (m *mockPolyTransaction) String() string {
	return m.hash.Hex()
}
func (m *mockPolyTransaction) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}
