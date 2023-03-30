package contracts

import (
	"crypto/ecdsa"
	"fmt"
	"math/rand"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/rs/zerolog/log"
)

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
		// case 3:
		// case 4:
		// case 5:
		// case 6:
		// case 7:
		// case 8:
		// case 9:
		// case 10:
		// case 11:
		// case 12:
		// case 13:
		// case 14:
		// case 15:
		// case 16:
		// case 17:
	}

	return nil, fmt.Errorf("Unrecognized precompiled address %d", address)
}

func GetRandomPrecompiledContractAddress() int {
	n := 2
	return rand.Intn(n) + 1 // [1, n + 1)
	// return rand.Intn(17) + 1
}
