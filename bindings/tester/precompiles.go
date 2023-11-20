package tester

import (
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
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

func GenerateECMulInput() []byte {
	x1 := big.NewInt(1)
	y1 := big.NewInt(2)
	s := big.NewInt(2)

	// Convert the x and y coordinates to 32-byte arrays in big-endian format
	x1Bytes := math.PaddedBigBytes(x1, 32)
	y1Bytes := math.PaddedBigBytes(y1, 32)
	sBytes := math.PaddedBigBytes(s, 32)

	inputData := append(append(x1Bytes, y1Bytes...), sBytes...)
	return inputData
}

func GenerateECPairingInput() []byte {
	x1 := "2cf44499d5d27bb186308b7af7af02ac5bc9eeb6a3d147c186b21fb1b76e18da"
	y1 := "2c0f001f52110ccfe69108924926e45f0b0c868df0e7bde1fe16d3242dc715f6"
	x2 := "1fb19bb476f6b9e44e2a32234da8212f61cd63919354bc06aef31e3cfaff3ebc"
	y2 := "22606845ff186793914e03e21df544c34ffe2f2f3504de8a79d9159eca2d98d9"
	x3 := "2bd368e28381e8eccb5fa81fc26cf3f048eea9abfdd85d7ed3ab3698d63e4f90"
	y3 := "2fe02e47887507adf0ff1743cbac6ba291e66f59be6bd763950bb16041a0a85e"

	x4 := "0000000000000000000000000000000000000000000000000000000000000001"
	y4 := "30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd45"
	x5 := "1971ff0471b09fa93caaf13cbf443c1aede09cc4328f5a62aad45f40ec133eb4"
	y5 := "091058a3141822985733cbdddfed0fd8d6c104e9e9eff40bf5abfef9ab163bc7"
	x6 := "2a23af9a5ce2ba2796c1f4e453a370eb0af8c212d9dc9acd8fc02c2e907baea2"
	y6 := "23a8eb0b0996252cb548a4487da97b02422ebc0e834613f954de6c7e0afdc1fc"

	inputHex := x1 + y1 + x2 + y2 + x3 + y3 + x4 + y4 + x5 + y5 + x6 + y6
	inputData, err := hex.DecodeString(inputHex)
	if err != nil {
		panic(err)
	}

	return inputData
}

func GenerateBlake2FInput() []byte {
	rounds := uint32(12)
	h := [8]uint64{1, 2, 3, 4, 5, 6, 7, 8}
	m := [16]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	t := [2]uint64{1, 2}
	f := byte(1)

	inputData := make([]byte, 213)
	binary.BigEndian.PutUint32(inputData[0:4], rounds)

	for i := 0; i < 8; i++ {
		binary.LittleEndian.PutUint64(inputData[4+(i*8):4+((i+1)*8)], h[i])
	}

	for i := 0; i < 16; i++ {
		binary.LittleEndian.PutUint64(inputData[68+(i*8):68+((i+1)*8)], m[i])
	}

	for i := 0; i < 2; i++ {
		binary.LittleEndian.PutUint64(inputData[196+(i*8):196+((i+1)*8)], t[i])
	}

	inputData[212] = f

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
	case 7:
		log.Trace().Str("method", "TestECMul").Msg("Executing contract method")
		inputData = GenerateECMulInput()
		return lt.TestECMul(opts, inputData)
	case 8:
		log.Trace().Str("method", "TestECPairing").Msg("Executing contract method")
		inputData = GenerateECPairingInput()
		return lt.TestECPairing(opts, inputData)
	case 9:
		log.Trace().Str("method", "TestBlake2f").Msg("Executing contract method")
		inputData = GenerateECPairingInput()
		return lt.TestBlake2f(opts, inputData)
	}

	return nil, fmt.Errorf("unrecognized precompiled address %d", address)
}

func GetRandomPrecompiledContractAddress() int {
	codes := []int{
		1,
		2,
		3,
		4,
		5,
		// 6, // NOTE: ecAdd requires a lot of gas and buggy
		// 7, // NOTE: ecMul requires a lot of gas and buggy
		8,
		9,
	}

	return codes[rand.Intn(len(codes))]
}
