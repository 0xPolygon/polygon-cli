package contracts

import (
	"crypto/ecdsa"
	"fmt"
	"math/rand"
	// "encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/rs/zerolog/log"
)

func byteSize(value *big.Int) int {
	// Calculate the byte size of the input value
	bytes := value.Bytes()
	if len(bytes) == 0 {
		return 0
	}
	return len(bytes)
}

func GenerateECRecoverInput(privateKey *ecdsa.PrivateKey) []byte {
	message := []byte("Test ecRecover")
	messageHash := crypto.Keccak256Hash(message)
	signature, err := crypto.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		panic(err)
	}

	// Prepare input data for ecRecover precompiled contract
	inputData := make([]byte, 128)
	copy(inputData[0:32], messageHash.Bytes())
	copy(inputData[32:64], common.LeftPadBytes(signature[64:65], 32))
	copy(inputData[64:96], common.LeftPadBytes(signature[0:32], 32))
	copy(inputData[96:128], common.LeftPadBytes(signature[32:64], 32))

	return inputData
}

func GenerateSHA256Input() []byte {
	inputData := []byte("Test")
	paddedInput := common.RightPadBytes(inputData, 32) // pad input to 32 bytes

	return paddedInput
}

func GenerateRIPEMD160Input() []byte {
	inputData := []byte("Test")

	return inputData
}

func GenerateIdentityInput() []byte {
	inputData := []byte("Test")

	return inputData
}

func GenerateModExpInput() []byte {
	base := big.NewInt(8)
	exponent := big.NewInt(9)
	modulus := big.NewInt(10)
	bSize := len(base.Bytes())
	eSize := len(exponent.Bytes())
	mSize := len(modulus.Bytes())
	inputData := append(
		common.LeftPadBytes(big.NewInt(int64(bSize)).Bytes(), 32),
		append(
			common.LeftPadBytes(big.NewInt(int64(eSize)).Bytes(), 32),
			append(
				common.LeftPadBytes(big.NewInt(int64(mSize)).Bytes(), 32),
				append(
					common.LeftPadBytes(base.Bytes(), bSize),
					append(
						common.LeftPadBytes(exponent.Bytes(), eSize),
						common.LeftPadBytes(modulus.Bytes(), mSize)...,
					)...,
				)...,
			)...,
		)...,
	)

	return inputData
}

func GenerateECAddInput() []byte {
	x1 := big.NewInt(1)
	y1 := big.NewInt(2)
	x2 := big.NewInt(1)
	y2 := big.NewInt(2)

	// Convert the x and y coordinates to 32-byte arrays in big-endian format
	x1Bytes := math.PaddedBigBytes(x1, 32)
	y1Bytes := math.PaddedBigBytes(y1, 32)
	x2Bytes := math.PaddedBigBytes(x2, 32)
	y2Bytes := math.PaddedBigBytes(y2, 32)

	inputData := append(append(append(x1Bytes, y1Bytes...), x2Bytes...), y2Bytes...)
	return inputData
}

func CallPrecompiledContracts(address int, lt *LoadTester, opts *bind.TransactOpts, iterations uint64, privateKey *ecdsa.PrivateKey) (*ethtypes.Transaction, error) {
	var inputData []byte

	switch address {
	case 1:
		log.Trace().Str("method", "TestECRecover").Msg("Executing contract method")
		inputData = GenerateECRecoverInput(privateKey)
		return lt.TestECRecover(opts, inputData)
	case 2:
		log.Trace().Str("method", "TestSHA256").Msg("Executing contract method")
		inputData = GenerateSHA256Input()
		return lt.TestSHA256(opts, inputData)
	case 3:
		log.Trace().Str("method", "TestRipemd160").Msg("Executing contract method")
		inputData = GenerateRIPEMD160Input()
		return lt.TestRipemd160(opts, inputData)
	case 4:
		log.Trace().Str("method", "TestIdentity").Msg("Executing contract method")
		inputData = GenerateIdentityInput()
		return lt.TestIdentity(opts, inputData)
	case 5:
		log.Trace().Str("method", "TestModExp").Msg("Executing contract method")
		inputData = GenerateModExpInput()
		return lt.TestModExp(opts, inputData)
	case 6:
		log.Trace().Str("method", "TestECAdd").Msg("Executing contract method")
		inputData = GenerateECAddInput()
		return lt.TestECAdd(opts, inputData)
		// case 7:
		// case 8:
		// case 9:
	}

	return nil, fmt.Errorf("Unrecognized precompiled address %d", address)
}

func GetRandomPrecompiledContractAddress() int {
	n := 6
	return rand.Intn(n) + 1 // [1, n + 1)
}
