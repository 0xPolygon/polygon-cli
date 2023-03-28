package contracts

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	// "math/rand"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

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

func CallPrecompiledContracts(ctx context.Context, c *ethclient.Client, address int, tops *bind.TransactOpts, iterations uint64, fromAddress common.Address, privateKey *ecdsa.PrivateKey) error {
	h := fmt.Sprintf("0x%x", address)
	contractPrecompiledAddress := common.HexToAddress(h)
	var inputData []byte

	switch address {
	case 1:
		inputData = GenerateECRecoverInput(privateKey)
	case 2:
		inputData = GenerateSHA256Input()
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
	// case 18:
	default:
		return fmt.Errorf("Unrecognized precompiled address %d", address)
	}

	// gasPrice
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("Failed to estimate gasPrice")
	}
	fmt.Printf("gasPrice: %d\n", gasPrice)

	fmt.Println("fromAddress: ", fromAddress)
	// Prepare the call message
	callMsg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &contractPrecompiledAddress,
		Gas:      10000000,
		GasPrice: gasPrice,
		Data:     inputData,
	}

	/**
		// Estimate gas
		gas, err := c.EstimateGas(ctx, callMsg)
		if err != nil {
	        return fmt.Errorf("Failed to estimate gas")
		}

		// Increase the gas limit
		tops.GasLimit = gas * 2
	*/

	// Call the precompiled contract
	callResult, err := c.CallContract(ctx, callMsg, nil)
	if err != nil {
		return fmt.Errorf("Failed to call precompiled contract 0x%x", contractPrecompiledAddress)
	}

	// Print the result
	log.Trace().Str("method", h).Msg("Executing contract method")
	fmt.Printf("callResult: 0x%x\n", callResult)

	return nil
}

func GetRandomPrecompiledContractAddress() int {
	return 2
	// return rand.Intn(17) + 1
}
