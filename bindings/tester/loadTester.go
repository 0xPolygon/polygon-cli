// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tester

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// LoadTesterMetaData contains all meta data concerning the LoadTester contract.
var LoadTesterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"F\",\"inputs\":[{\"name\":\"rounds\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"h\",\"type\":\"bytes32[2]\",\"internalType\":\"bytes32[2]\"},{\"name\":\"m\",\"type\":\"bytes32[4]\",\"internalType\":\"bytes32[4]\"},{\"name\":\"t\",\"type\":\"bytes8[2]\",\"internalType\":\"bytes8[2]\"},{\"name\":\"f\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[2]\",\"internalType\":\"bytes32[2]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dumpster\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCallCounter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"inc\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"loopBlockHashUntilLimit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"loopUntilLimit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"store\",\"inputs\":[{\"name\":\"trash\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testADD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testADDMOD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testADDRESS\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testAND\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testBALANCE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testBASEFEE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testBLOCKHASH\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testBYTE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testBlake2f\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCALLDATACOPY\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCALLDATALOAD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCALLDATASIZE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCALLER\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCALLVALUE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCHAINID\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCODECOPY\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCODESIZE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCOINBASE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testDIFFICULTY\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testDIV\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testECAdd\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testECMul\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testECPairing\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testECRecover\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testEQ\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testEXP\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testEXTCODESIZE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testGAS\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testGASLIMIT\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testGASPRICE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testGT\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testISZERO\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testIdentity\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testLOG0\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testLOG1\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testLOG2\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testLOG3\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testLOG4\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testLT\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testMLOAD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testMOD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testMSIZE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testMSTORE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testMSTORE8\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testMUL\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testMULMOD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testModExp\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testNOT\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testNUMBER\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testOR\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testORIGIN\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testP256Verify\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testRETURNDATACOPY\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testRETURNDATASIZE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testRipemd160\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSAR\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSDIV\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSELFBALANCE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSGT\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSHA256\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSHA3\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSHL\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSHR\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSIGNEXTEND\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSLOAD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSLT\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSMOD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSSTORE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSUB\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testTIMESTAMP\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testXOR\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x608060405234801561000f575f80fd5b506123b98061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610459575f3560e01c806380947f8011610242578063bf529ca111610140578063dd9bef60116100bf578063f279ca8111610084578063f279ca81146109e5578063f4d1fc61146109f8578063f58fc36a14610a0b578063f6b0bbf714610a1e578063fde7721c14610a3e575f80fd5b8063dd9bef6014610986578063de97a36314610999578063e9f9b3f2146109ac578063ea5141e6146109bf578063edf003cf146109d2575f80fd5b8063ce3cf4ef11610105578063ce3cf4ef14610927578063d117320b1461093a578063d51e7b5b1461094d578063d53ff3fd14610960578063d93cd55814610973575f80fd5b8063bf529ca1146108bb578063c360aba6146108ce578063c420eb61146108e1578063c4bd65d5146108f4578063c711e53914610907575f80fd5b8063a18683cb116101cc578063b374012b11610191578063b374012b1461085c578063b3d847f21461086f578063b7b8620714610882578063b81c148414610895578063bdc875fc146108a8575f80fd5b8063a18683cb146107fb578063a271b7211461081b578063a60a108714610823578063a645c9c214610836578063acaebdf614610849575f80fd5b8063962e4dc211610212578063962e4dc21461079c57806398456f3e146107af5780639a2b7c81146107c25780639cce7cf9146107d5578063a040aec6146107e8575f80fd5b806380947f8014610750578063880eff3914610763578063918a5fcd1461077657806391e7b27714610789575f80fd5b80633430ec061161035a578063613d0a82116102d957806371d91d281161029e57806371d91d28146106e457806372de3cbd146106f75780637b6e0b0e146107175780637c191d201461072a5780637de8c6f81461073d575f80fd5b8063613d0a821461069057806363138d4f146106a3578063659bbb4f146106b65780636e7f1fe7146106be5780636f099c8d146106d1575f80fd5b806344cf3bc71161031f57806344cf3bc7146106315780634a61af1f146106445780634d2c74b3146106575780635590c2d91461066a57806360e13cde1461067d575f80fd5b80633430ec06146105dd578063371303c0146105f05780633a411f12146105f85780633a425dfc1461060b57806340fe26621461061e575f80fd5b806318093b46116103e6578063219cddeb116103ab578063219cddeb1461057e5780632294fc7f146105915780632871ef85146105a45780632b21ef44146105b75780632d34e798146105ca575f80fd5b806318093b461461051f57806319b621d6146105325780631aba07ea146105455780631de2f343146105585780632007332e1461056b575f80fd5b80630ba8a73b1161042c5780630ba8a73b146104cc5780631287a68c146104df578063135d52f7146104e65780631581cf19146104f9578063165821501461050c575f80fd5b8063034aef711461045d578063050082f814610486578063087b4e84146104995780630b3b996a146104ac575b5f80fd5b61047061046b3660046119b5565b610a51565b60405161047d91906119e3565b60405180910390f35b6104706104943660046119b5565b610a81565b6104706104a73660046119b5565b610aa8565b6104bf6104ba366004611ade565b610ad6565b60405161047d9190611b6e565b6104706104da3660046119b5565b610af7565b5f54610470565b6104706104f43660046119b5565b610b1a565b6104706105073660046119b5565b610b3a565b61047061051a3660046119b5565b610b61565b61047061052d3660046119b5565b610b8a565b6104706105403660046119b5565b610bb3565b6104706105533660046119b5565b610c1b565b6104706105663660046119b5565b610c4f565b6104706105793660046119b5565b610c79565b61047061058c3660046119b5565b610c99565b61047061059f3660046119b5565b610cc0565b6104706105b23660046119b5565b610cf3565b6104706105c53660046119b5565b610d1a565b6104706105d83660046119b5565b610d41565b6104bf6105eb3660046119b5565b610d68565b610470610e0c565b6104706106063660046119b5565b610e23565b6104706106193660046119b5565b610e43565b61047061062c3660046119b5565b610e6b565b61047061063f3660046119b5565b610e98565b6104706106523660046119b5565b610ebf565b6104706106653660046119b5565b610ee9565b6104706106783660046119b5565b610f10565b61047061068b3660046119b5565b610f43565b6104bf61069e366004611ade565b610f6c565b6104706106b1366004611ade565b610f95565b610470610fbb565b6104706106cc3660046119b5565b610ff3565b6104706106df3660046119b5565b61101c565b6104706106f23660046119b5565b611043565b61070a610705366004611d2d565b61106c565b60405161047d9190611de8565b6104706107253660046119b5565b6110eb565b6104706107383660046119b5565b611113565b61047061074b3660046119b5565b61113a565b61047061075e3660046119b5565b61115a565b6104706107713660046119b5565b611184565b6104706107843660046119b5565b6111ae565b6104706107973660046119b5565b6111d5565b6104bf6107aa366004611ade565b611211565b6104706107bd3660046119b5565b611260565b6104706107d03660046119b5565b61128e565b6104bf6107e3366004611ade565b6112ae565b6104bf6107f6366004611ade565b6112cc565b61080e610809366004611ade565b611405565b60405161047d9190611e0f565b61047061145c565b6104706108313660046119b5565b61149b565b6104706108443660046119b5565b6114c2565b6104706108573660046119b5565b6114e2565b61047061086a366004611e6b565b61150a565b61047061087d3660046119b5565b611538565b6104706108903660046119b5565b61155f565b6104706108a33660046119b5565b611586565b6104706108b63660046119b5565b6115ad565b6104706108c93660046119b5565b6115d4565b6104706108dc3660046119b5565b611606565b6104706108ef3660046119b5565b611626565b6104706109023660046119b5565b61164d565b61091a610915366004611ade565b611671565b60405161047d9190611eb8565b6104706109353660046119b5565b6116f4565b6104706109483660046119b5565b61171c565b61047061095b3660046119b5565b611743565b61047061096e3660046119b5565b611763565b6104706109813660046119b5565b611783565b6104706109943660046119b5565b6117ac565b6104706109a73660046119b5565b6117dc565b6104706109ba3660046119b5565b611804565b6104706109cd3660046119b5565b61182b565b6104bf6109e0366004611ade565b611857565b6104706109f33660046119b5565b611894565b610470610a063660046119b5565b6118bc565b610470610a193660046119b5565b6118e5565b610a31610a2c366004611ade565b61190e565b60405161047d9190611edb565b610470610a4c3660046119b5565b611939565b5f610a5a610e0c565b5065deadbeef00365f805b84811015610a7857369150600101610a65565b50909392505050565b5f610a8a610e0c565b5065deadbeef00325f805b84811015610a7857329150600101610a95565b5f610ab1610e0c565b5065deadbeef00525f5b83811015610acf575f829052600101610abb565b5092915050565b606060086040828451602086015f855af180610af0575f80fd5b5050919050565b5f610b00610e0c565b5065deadbeef00015f5b83811015610acf57600101610b0a565b5f610b23610e0c565b5065deadbeef00175f8315610acf57600101610b0a565b5f610b43610e0c565b5065deadbeef00345f805b84811015610a7857349150600101610b4e565b5f610b6a610e0c565b5065deadbeef00065f5b83811015610acf575f1990910690600101610b74565b5f610b93610e0c565b5065deadbeef00135f805b84811015610a78576001808413925001610b9e565b5f610bbc610e0c565b506001600160e01b03195f90815265deadbeef002090805b84811015610bea5760045f209150600101610bd4565b507f29045a592007d0c246ef02c2223570da9522d0cf0f73282c79a1bc8f0bb2c2388114610acf57505f9392505050565b5f610c24610e0c565b5065deadbeef00a460108190525f5b83811015610acf576004600360028360066010a4600101610c33565b5f610c58610e0c565b5065deadbeef001a5f805b84811015610a78575f83901a9150600101610c63565b5f610c82610e0c565b5065deadbeef001b5f8315610acf57600101610b0a565b5f610ca2610e0c565b5065deadbeef00425f805b84811015610a7857429150600101610cad565b5f610cc9610e0c565b5065deadbeef00315f30815b85811015610ce95781319250600101610cd5565b5091949350505050565b5f610cfc610e0c565b5065deadbeef00485f805b84811015610a7857489150600101610d07565b5f610d23610e0c565b5065deadbeef003d5f805b84811015610a78573d9150600101610d2e565b5f610d4a610e0c565b5065deadbeef00435f805b84811015610a7857439150600101610d55565b60028181548110610d77575f80fd5b905f5260205f20018054909150610d8d90611efd565b80601f0160208091040260200160405190810160405280929190818152602001828054610db990611efd565b8015610e045780601f10610ddb57610100808354040283529160200191610e04565b820191905f5260205f20905b815481529060010190602001808311610de757829003601f168201915b505050505081565b5f8054610e1a906001611f3d565b5f819055919050565b5f610e2c610e0c565b5065deadbeef00045f8315610acf57600101610b0a565b5f610e4c610e0c565b5065deadbeef00375f5b83811015610acf5760205f8037600101610e56565b5f610e74610e0c565b5065deadbeef00a060108190525f5b83811015610acf5760066010a0600101610e83565b5f610ea1610e0c565b5065deadbeef00335f805b84811015610a7857339150600101610eac565b5f610ec8610e0c565b5065deadbeef00535f5b83811015610acf5763deadbeef5f52600101610ed2565b5f610ef2610e0c565b5065deadbeef003a5f805b84811015610a78573a9150600101610efd565b5f610f19610e0c565b5065deadbeef00515f818152805b84811015610f3b575f519150600101610f27565b509392505050565b5f610f4c610e0c565b5065deadbeef001d5f5b83811015610acf575f9190911d90600101610f56565b6060600560208301835160405160208183855f885af180610f8b575f80fd5b5095945050505050565b5f60026020830183518360208183855f885af180610fb1575f80fd5b5050505050919050565b5f610fc4610e0c565b505b6103e85a1115610fec576001805f828254610fe19190611f3d565b90915550610fc69050565b5060015490565b5f610ffc610e0c565b5065deadbeef00105f805b84811015610a78576001838110925001611007565b5f611025610e0c565b5065deadbeef00445f805b84811015610a7857449150600101611030565b5f61104c610e0c565b5065deadbeef00115f805b84811015610a78576001808411925001611057565b611074611964565b600961107e611964565b5f88885160208a0151895160208b015160408c015160608d01518c5160208e01518d6040516020016110b99a99989796959493929190611fa3565b604051601f19818303018152604091825291508260d560208401865f19fa6110df575f80fd5b50979650505050505050565b5f6110f4610e0c565b5065deadbeef003e5f5b83811015610acf5760205f803e6001016110fe565b5f61111c610e0c565b5065deadbeef00455f805b84811015610a7857459150600101611127565b5f611143610e0c565b5065deadbeef00025f8315610acf57600101610b0a565b5f611163610e0c565b5065deadbeef00085f5b83811015610acf575f195f8308915060010161116d565b5f61118d610e0c565b5065deadbeef00545f8181555b83811015610acf575f54915060010161119a565b5f6111b7610e0c565b5065deadbeef005a5f805b84811015610a78575a91506001016111c2565b5f6111de610e0c565b5065deadbeef00195f5b838110156111fb579019906001016111e8565b5065deadbeef0019811461120b57195b92915050565b6060815160601461123d5760405162461bcd60e51b81526004016112349061207d565b60405180910390fd5b600760208301835160408482845f875af180611257575f80fd5b50505050919050565b5f611269610e0c565b5065deadbeef00a160108190525f5b83811015610acf578060066010a1600101611278565b5f611297610e0c565b5065deadbeef00165f8315610acf57600101610b0a565b60606004602083018351604051818183855f885af180610f8b575f80fd5b60606112d6611964565b7f48c9bdf267e6096a3ba7ca8485ae67bb2bf894fe72f36e3cf1361d5f3af54fa581527fd182e6ad7f520e511f6c3e2b8c68059b6bbd41fbabd9831f79217e1319cde05b602082015261132761197f565b6261626360e81b81525f6020820181905260408201819052606082015261134c611964565b600360f81b81525f6020820181905261136a600c858585600161106c565b9050611374611964565b7fba80a53f981c4d0d6a2797b69f12f6e94c212f14685ac4b74b12bb6fdbffa2d181527f7d87c5392aab792dc252d5de4533cc9518d38aa8dbf1925ab92386edd4009923602082015280518251146113de5760405162461bcd60e51b8152600401611234906120d2565b6020810151602083015114610fb15760405162461bcd60e51b815260040161123490612125565b5f81516080146114275760405162461bcd60e51b815260040161123490612168565b6001602083016040840151601f1a602082015260206040516080835f865af18061144f575f80fd5b6040515195945050505050565b5f611465610e0c565b505b6103e85a1115610fec576001805f8282546114829190611f3d565b909155505060015461149590439061218c565b50611467565b5f6114a4610e0c565b5065deadbeef00465f805b84811015610a78574691506001016114af565b5f6114cb610e0c565b5065deadbeef00055f8315610acf57600101610b0a565b5f6114eb610e0c565b5065deadbeef00395f5b83811015610acf5760205f80396001016114f5565b600280546001810182555f9182528390839060208420019161152d919083612247565b505060025492915050565b5f611541610e0c565b5065deadbeef00595f805b84811015610a785759915060010161154c565b5f611568610e0c565b5065deadbeef00385f805b84811015610a7857389150600101611573565b5f61158f610e0c565b5065deadbeef00415f805b84811015610a785741915060010161159a565b5f6115b6610e0c565b5065deadbeef00305f805b84811015610a78573091506001016115c1565b5f6115dd610e0c565b5065deadbeef00a360108190525f5b83811015610acf57600360028260066010a36001016115ec565b5f61160f610e0c565b5065deadbeef000b5f8315610acf57600101610b0a565b5f61162f610e0c565b5065deadbeef00475f805b84811015610a785747915060010161163a565b5f611656610e0c565b5065deadbeef001c5f805b84811015610a7857600101611661565b5f610100818082856040516116869190612338565b5f60405180830381855afa9150503d805f81146116be576040513d603f01601f191681016040523d815291503d5f602084013e6116c3565b606091505b5091509150816116d5576116d5612343565b6020810181518101906116e89190612365565b60011495945050505050565b5f6116fd610e0c565b5065deadbeef00355f805b84811015610a78575f359150600101611708565b5f611725610e0c565b5065deadbeef00555f5b83811015610acf575f82905560010161172f565b5f61174c610e0c565b5065deadbeef00185f8315610acf57600101610b0a565b5f61176c610e0c565b5065deadbeef00035f8315610acf57600101610b0a565b5f61178c610e0c565b5065deadbeef00075f5b83811015610acf575f1990910790600101611796565b5f6117b5610e0c565b5065deadbeef00a260108190525f5b83811015610acf5760028160066010a26001016117c4565b5f6117e5610e0c565b5065deadbeef000a5f5b83811015610acf5760019182900a91016117ef565b5f61180d610e0c565b5065deadbeef00145f805b84811015610a7857600191508101611818565b5f611834610e0c565b5065deadbeef00405f5f194301815b85811015610ce95781409250600101611843565b6060815160801461187a5760405162461bcd60e51b81526004016112349061207d565b600660208301835160408482845f875af180611257575f80fd5b5f61189d610e0c565b5065deadbeef00155f805b84811015610a7857821591506001016118a8565b5f6118c5610e0c565b5065deadbeef00125f805b84811015610a785760018381129250016118d0565b5f6118ee610e0c565b5065deadbeef003b5f30815b85811015610ce957813b92506001016118fa565b5f600360208301835160405160148183855f885af18061192c575f80fd5b8151979650505050505050565b5f611942610e0c565b5065deadbeef00095f5b83811015610acf575f1960018309915060010161194c565b60405160408082018152600290829080368337509192915050565b6040516080808201604052600490829080368337509192915050565b805b81146119a7575f80fd5b50565b803561120b8161199b565b5f602082840312156119c8576119c85f80fd5b5f6119d384846119aa565b949350505050565b805b82525050565b6020810161120b82846119db565b634e487b7160e01b5f52604160045260245ffd5b601f19601f830116810181811067ffffffffffffffff82111715611a2b57611a2b6119f1565b6040525050565b5f611a3f5f604051905090565b9050611a4b8282611a05565b919050565b5f67ffffffffffffffff821115611a6957611a696119f1565b601f19601f83011660200192915050565b82818337505f910152565b5f611a97611a9284611a50565b611a32565b905082815260208101848484011115611ab157611ab15f80fd5b610f3b848285611a7a565b5f82601f830112611ace57611ace5f80fd5b81356119d3848260208601611a85565b5f60208284031215611af157611af15f80fd5b813567ffffffffffffffff811115611b0a57611b0a5f80fd5b6119d384828501611abc565b5f5b83811015611b30578082015183820152602001611b18565b50505f910152565b5f611b46825f815192915050565b808452602084019350611b5d818560208601611b16565b601f01601f19169290920192915050565b60208082528101611b7f8184611b38565b9392505050565b63ffffffff811661199d565b803561120b81611b86565b5f67ffffffffffffffff821115611bb657611bb66119f1565b5060200290565b5f611bca611a9284611b9d565b90508060208402830185811115611be257611be25f80fd5b835b81811015611c065780611bf788826119aa565b84525060209283019201611be4565b5050509392505050565b5f82601f830112611c2257611c225f80fd5b60026119d3848285611bbd565b5f611c3c611a9284611b9d565b90508060208402830185811115611c5457611c545f80fd5b835b81811015611c065780611c6988826119aa565b84525060209283019201611c56565b5f82601f830112611c8a57611c8a5f80fd5b60046119d3848285611c2f565b6001600160c01b0319811661199d565b803561120b81611c97565b5f611cbf611a9284611b9d565b90508060208402830185811115611cd757611cd75f80fd5b835b81811015611c065780611cec8882611ca7565b84525060209283019201611cd9565b5f82601f830112611d0d57611d0d5f80fd5b60026119d3848285611cb2565b80151561199d565b803561120b81611d1a565b5f805f805f6101408688031215611d4557611d455f80fd5b5f611d508888611b92565b9550506020611d6188828901611c10565b9450506060611d7288828901611c78565b93505060e0611d8388828901611cfb565b925050610120611d9588828901611d22565b9150509295509295909350565b5f611dad83836119db565b505060200190565b600281805f5b83811015611de0578151611dcf8782611da2565b965060208301925050600101611dbb565b505050505050565b6040810161120b8284611db5565b5f6001600160a01b03821661120b565b6119dd81611df6565b6020810161120b8284611e06565b5f8083601f840112611e3057611e305f80fd5b50813567ffffffffffffffff811115611e4a57611e4a5f80fd5b602083019150836001820283011115611e6457611e645f80fd5b9250929050565b5f8060208385031215611e7f57611e7f5f80fd5b823567ffffffffffffffff811115611e9857611e985f80fd5b611ea485828601611e1d565b92509250509250929050565b8015156119dd565b6020810161120b8284611eb0565b6bffffffffffffffffffffffff1981166119dd565b6020810161120b8284611ec6565b634e487b7160e01b5f52602260045260245ffd5b600281046001821680611f1157607f821691505b602082108103611f2357611f23611ee9565b50919050565b634e487b7160e01b5f52601160045260245ffd5b8082018082111561120b5761120b611f29565b5f61120b8260e01b90565b6119dd63ffffffff8216611f50565b806119dd565b90565b6001600160c01b031981166119dd565b5f61120b8260f81b90565b5f61120b82611f83565b6119dd811515611f8e565b5f611fae828d611f5b565b600482019150611fbe828c611f6a565b602082019150611fce828b611f6a565b602082019150611fde828a611f6a565b602082019150611fee8289611f6a565b602082019150611ffe8288611f6a565b60208201915061200e8287611f6a565b60208201915061201e8286611f73565b60088201915061202e8285611f73565b60088201915061203e8284611f98565b506001019a9950505050505050505050565b601481525f6020820173092dcecc2d8d2c840d2dce0eae840d8cadccee8d60631b815291505b5060200190565b6020808252810161120b81612050565b602681525f602082017f54657374426c616b653266202d204669727374206861736820646f65736e2774815265040dac2e8c6d60d31b602082015291505b5060400190565b6020808252810161120b8161208d565b602781525f602082017f54657374426c616b653266202d205365636f6e64206861736820646f65736e278152660e840dac2e8c6d60cb1b602082015291506120cb565b6020808252810161120b816120e2565b601a81525f602082017f496e76616c696420696e7075742064617461206c656e6774682e00000000000081529150612076565b6020808252810161120b81612135565b634e487b7160e01b5f52601260045260245ffd5b5f8261219a5761219a612178565b500690565b5f61120b611f708381565b6121b38361219f565b81545f1960089490940293841b1916921b91909117905550565b5f6121d98184846121aa565b505050565b818110156121f8576121f05f826121cd565b6001016121de565b5050565b601f8211156121d957612219815f81815281906020902092915050565b6020601f8501048101602085101561222e5750805b6122406020601f8601048301826121de565b5050505050565b8267ffffffffffffffff811115612260576122606119f1565b61226a8254611efd565b6122758282856121fc565b5f601f8311600181146122a6575f841561228f5750858201355b5f19600886021c1981166002860217865550612309565b601f1984166122bf865f81815281906020902092915050565b5f5b828110156122e157888501358255602094850194600190920191016122c1565b868310156122fc575f19601f88166008021c19858a01351682555b6001600288020188555050505b50505050505050565b5f612320825f815192915050565b61232e818560208601611b16565b9290920192915050565b5f611b7f8284612312565b634e487b7160e01b5f52600160045260245ffd5b5f8151905061120b8161199b565b5f60208284031215612378576123785f80fd5b5f6119d3848461235756fea26469706673582212208edc99ef8989da55c4d6a5a266b0da5e84a8a9a1924c7d37709880241768efc764736f6c63430008170033",
}

// LoadTesterABI is the input ABI used to generate the binding from.
// Deprecated: Use LoadTesterMetaData.ABI instead.
var LoadTesterABI = LoadTesterMetaData.ABI

// LoadTesterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LoadTesterMetaData.Bin instead.
var LoadTesterBin = LoadTesterMetaData.Bin

// DeployLoadTester deploys a new Ethereum contract, binding an instance of LoadTester to it.
func DeployLoadTester(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LoadTester, error) {
	parsed, err := LoadTesterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LoadTesterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LoadTester{LoadTesterCaller: LoadTesterCaller{contract: contract}, LoadTesterTransactor: LoadTesterTransactor{contract: contract}, LoadTesterFilterer: LoadTesterFilterer{contract: contract}}, nil
}

// LoadTester is an auto generated Go binding around an Ethereum contract.
type LoadTester struct {
	LoadTesterCaller     // Read-only binding to the contract
	LoadTesterTransactor // Write-only binding to the contract
	LoadTesterFilterer   // Log filterer for contract events
}

// LoadTesterCaller is an auto generated read-only Go binding around an Ethereum contract.
type LoadTesterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LoadTesterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LoadTesterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LoadTesterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LoadTesterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LoadTesterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LoadTesterSession struct {
	Contract     *LoadTester       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LoadTesterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LoadTesterCallerSession struct {
	Contract *LoadTesterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// LoadTesterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LoadTesterTransactorSession struct {
	Contract     *LoadTesterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// LoadTesterRaw is an auto generated low-level Go binding around an Ethereum contract.
type LoadTesterRaw struct {
	Contract *LoadTester // Generic contract binding to access the raw methods on
}

// LoadTesterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LoadTesterCallerRaw struct {
	Contract *LoadTesterCaller // Generic read-only contract binding to access the raw methods on
}

// LoadTesterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LoadTesterTransactorRaw struct {
	Contract *LoadTesterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLoadTester creates a new instance of LoadTester, bound to a specific deployed contract.
func NewLoadTester(address common.Address, backend bind.ContractBackend) (*LoadTester, error) {
	contract, err := bindLoadTester(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LoadTester{LoadTesterCaller: LoadTesterCaller{contract: contract}, LoadTesterTransactor: LoadTesterTransactor{contract: contract}, LoadTesterFilterer: LoadTesterFilterer{contract: contract}}, nil
}

// NewLoadTesterCaller creates a new read-only instance of LoadTester, bound to a specific deployed contract.
func NewLoadTesterCaller(address common.Address, caller bind.ContractCaller) (*LoadTesterCaller, error) {
	contract, err := bindLoadTester(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LoadTesterCaller{contract: contract}, nil
}

// NewLoadTesterTransactor creates a new write-only instance of LoadTester, bound to a specific deployed contract.
func NewLoadTesterTransactor(address common.Address, transactor bind.ContractTransactor) (*LoadTesterTransactor, error) {
	contract, err := bindLoadTester(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LoadTesterTransactor{contract: contract}, nil
}

// NewLoadTesterFilterer creates a new log filterer instance of LoadTester, bound to a specific deployed contract.
func NewLoadTesterFilterer(address common.Address, filterer bind.ContractFilterer) (*LoadTesterFilterer, error) {
	contract, err := bindLoadTester(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LoadTesterFilterer{contract: contract}, nil
}

// bindLoadTester binds a generic wrapper to an already deployed contract.
func bindLoadTester(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LoadTesterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LoadTester *LoadTesterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LoadTester.Contract.LoadTesterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LoadTester *LoadTesterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoadTester.Contract.LoadTesterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LoadTester *LoadTesterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LoadTester.Contract.LoadTesterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LoadTester *LoadTesterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LoadTester.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LoadTester *LoadTesterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoadTester.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LoadTester *LoadTesterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LoadTester.Contract.contract.Transact(opts, method, params...)
}

// F is a free data retrieval call binding the contract method 0x72de3cbd.
//
// Solidity: function F(uint32 rounds, bytes32[2] h, bytes32[4] m, bytes8[2] t, bool f) view returns(bytes32[2])
func (_LoadTester *LoadTesterCaller) F(opts *bind.CallOpts, rounds uint32, h [2][32]byte, m [4][32]byte, t [2][8]byte, f bool) ([2][32]byte, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "F", rounds, h, m, t, f)

	if err != nil {
		return *new([2][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([2][32]byte)).(*[2][32]byte)

	return out0, err

}

// F is a free data retrieval call binding the contract method 0x72de3cbd.
//
// Solidity: function F(uint32 rounds, bytes32[2] h, bytes32[4] m, bytes8[2] t, bool f) view returns(bytes32[2])
func (_LoadTester *LoadTesterSession) F(rounds uint32, h [2][32]byte, m [4][32]byte, t [2][8]byte, f bool) ([2][32]byte, error) {
	return _LoadTester.Contract.F(&_LoadTester.CallOpts, rounds, h, m, t, f)
}

// F is a free data retrieval call binding the contract method 0x72de3cbd.
//
// Solidity: function F(uint32 rounds, bytes32[2] h, bytes32[4] m, bytes8[2] t, bool f) view returns(bytes32[2])
func (_LoadTester *LoadTesterCallerSession) F(rounds uint32, h [2][32]byte, m [4][32]byte, t [2][8]byte, f bool) ([2][32]byte, error) {
	return _LoadTester.Contract.F(&_LoadTester.CallOpts, rounds, h, m, t, f)
}

// Dumpster is a free data retrieval call binding the contract method 0x3430ec06.
//
// Solidity: function dumpster(uint256 ) view returns(bytes)
func (_LoadTester *LoadTesterCaller) Dumpster(opts *bind.CallOpts, arg0 *big.Int) ([]byte, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "dumpster", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Dumpster is a free data retrieval call binding the contract method 0x3430ec06.
//
// Solidity: function dumpster(uint256 ) view returns(bytes)
func (_LoadTester *LoadTesterSession) Dumpster(arg0 *big.Int) ([]byte, error) {
	return _LoadTester.Contract.Dumpster(&_LoadTester.CallOpts, arg0)
}

// Dumpster is a free data retrieval call binding the contract method 0x3430ec06.
//
// Solidity: function dumpster(uint256 ) view returns(bytes)
func (_LoadTester *LoadTesterCallerSession) Dumpster(arg0 *big.Int) ([]byte, error) {
	return _LoadTester.Contract.Dumpster(&_LoadTester.CallOpts, arg0)
}

// GetCallCounter is a free data retrieval call binding the contract method 0x1287a68c.
//
// Solidity: function getCallCounter() view returns(uint256)
func (_LoadTester *LoadTesterCaller) GetCallCounter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "getCallCounter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCallCounter is a free data retrieval call binding the contract method 0x1287a68c.
//
// Solidity: function getCallCounter() view returns(uint256)
func (_LoadTester *LoadTesterSession) GetCallCounter() (*big.Int, error) {
	return _LoadTester.Contract.GetCallCounter(&_LoadTester.CallOpts)
}

// GetCallCounter is a free data retrieval call binding the contract method 0x1287a68c.
//
// Solidity: function getCallCounter() view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) GetCallCounter() (*big.Int, error) {
	return _LoadTester.Contract.GetCallCounter(&_LoadTester.CallOpts)
}

// Inc is a paid mutator transaction binding the contract method 0x371303c0.
//
// Solidity: function inc() returns(uint256)
func (_LoadTester *LoadTesterTransactor) Inc(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "inc")
}

// Inc is a paid mutator transaction binding the contract method 0x371303c0.
//
// Solidity: function inc() returns(uint256)
func (_LoadTester *LoadTesterSession) Inc() (*types.Transaction, error) {
	return _LoadTester.Contract.Inc(&_LoadTester.TransactOpts)
}

// Inc is a paid mutator transaction binding the contract method 0x371303c0.
//
// Solidity: function inc() returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) Inc() (*types.Transaction, error) {
	return _LoadTester.Contract.Inc(&_LoadTester.TransactOpts)
}

// LoopBlockHashUntilLimit is a paid mutator transaction binding the contract method 0xa271b721.
//
// Solidity: function loopBlockHashUntilLimit() returns(uint256)
func (_LoadTester *LoadTesterTransactor) LoopBlockHashUntilLimit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "loopBlockHashUntilLimit")
}

// LoopBlockHashUntilLimit is a paid mutator transaction binding the contract method 0xa271b721.
//
// Solidity: function loopBlockHashUntilLimit() returns(uint256)
func (_LoadTester *LoadTesterSession) LoopBlockHashUntilLimit() (*types.Transaction, error) {
	return _LoadTester.Contract.LoopBlockHashUntilLimit(&_LoadTester.TransactOpts)
}

// LoopBlockHashUntilLimit is a paid mutator transaction binding the contract method 0xa271b721.
//
// Solidity: function loopBlockHashUntilLimit() returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) LoopBlockHashUntilLimit() (*types.Transaction, error) {
	return _LoadTester.Contract.LoopBlockHashUntilLimit(&_LoadTester.TransactOpts)
}

// LoopUntilLimit is a paid mutator transaction binding the contract method 0x659bbb4f.
//
// Solidity: function loopUntilLimit() returns(uint256)
func (_LoadTester *LoadTesterTransactor) LoopUntilLimit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "loopUntilLimit")
}

// LoopUntilLimit is a paid mutator transaction binding the contract method 0x659bbb4f.
//
// Solidity: function loopUntilLimit() returns(uint256)
func (_LoadTester *LoadTesterSession) LoopUntilLimit() (*types.Transaction, error) {
	return _LoadTester.Contract.LoopUntilLimit(&_LoadTester.TransactOpts)
}

// LoopUntilLimit is a paid mutator transaction binding the contract method 0x659bbb4f.
//
// Solidity: function loopUntilLimit() returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) LoopUntilLimit() (*types.Transaction, error) {
	return _LoadTester.Contract.LoopUntilLimit(&_LoadTester.TransactOpts)
}

// Store is a paid mutator transaction binding the contract method 0xb374012b.
//
// Solidity: function store(bytes trash) returns(uint256)
func (_LoadTester *LoadTesterTransactor) Store(opts *bind.TransactOpts, trash []byte) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "store", trash)
}

// Store is a paid mutator transaction binding the contract method 0xb374012b.
//
// Solidity: function store(bytes trash) returns(uint256)
func (_LoadTester *LoadTesterSession) Store(trash []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.Store(&_LoadTester.TransactOpts, trash)
}

// Store is a paid mutator transaction binding the contract method 0xb374012b.
//
// Solidity: function store(bytes trash) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) Store(trash []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.Store(&_LoadTester.TransactOpts, trash)
}

// TestADD is a paid mutator transaction binding the contract method 0x0ba8a73b.
//
// Solidity: function testADD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestADD(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testADD", x)
}

// TestADD is a paid mutator transaction binding the contract method 0x0ba8a73b.
//
// Solidity: function testADD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestADD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestADD(&_LoadTester.TransactOpts, x)
}

// TestADD is a paid mutator transaction binding the contract method 0x0ba8a73b.
//
// Solidity: function testADD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestADD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestADD(&_LoadTester.TransactOpts, x)
}

// TestADDMOD is a paid mutator transaction binding the contract method 0x80947f80.
//
// Solidity: function testADDMOD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestADDMOD(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testADDMOD", x)
}

// TestADDMOD is a paid mutator transaction binding the contract method 0x80947f80.
//
// Solidity: function testADDMOD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestADDMOD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestADDMOD(&_LoadTester.TransactOpts, x)
}

// TestADDMOD is a paid mutator transaction binding the contract method 0x80947f80.
//
// Solidity: function testADDMOD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestADDMOD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestADDMOD(&_LoadTester.TransactOpts, x)
}

// TestADDRESS is a paid mutator transaction binding the contract method 0xbdc875fc.
//
// Solidity: function testADDRESS(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestADDRESS(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testADDRESS", x)
}

// TestADDRESS is a paid mutator transaction binding the contract method 0xbdc875fc.
//
// Solidity: function testADDRESS(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestADDRESS(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestADDRESS(&_LoadTester.TransactOpts, x)
}

// TestADDRESS is a paid mutator transaction binding the contract method 0xbdc875fc.
//
// Solidity: function testADDRESS(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestADDRESS(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestADDRESS(&_LoadTester.TransactOpts, x)
}

// TestAND is a paid mutator transaction binding the contract method 0x9a2b7c81.
//
// Solidity: function testAND(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestAND(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testAND", x)
}

// TestAND is a paid mutator transaction binding the contract method 0x9a2b7c81.
//
// Solidity: function testAND(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestAND(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestAND(&_LoadTester.TransactOpts, x)
}

// TestAND is a paid mutator transaction binding the contract method 0x9a2b7c81.
//
// Solidity: function testAND(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestAND(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestAND(&_LoadTester.TransactOpts, x)
}

// TestBALANCE is a paid mutator transaction binding the contract method 0x2294fc7f.
//
// Solidity: function testBALANCE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestBALANCE(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testBALANCE", x)
}

// TestBALANCE is a paid mutator transaction binding the contract method 0x2294fc7f.
//
// Solidity: function testBALANCE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestBALANCE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestBALANCE(&_LoadTester.TransactOpts, x)
}

// TestBALANCE is a paid mutator transaction binding the contract method 0x2294fc7f.
//
// Solidity: function testBALANCE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestBALANCE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestBALANCE(&_LoadTester.TransactOpts, x)
}

// TestBASEFEE is a paid mutator transaction binding the contract method 0x2871ef85.
//
// Solidity: function testBASEFEE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestBASEFEE(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testBASEFEE", x)
}

// TestBASEFEE is a paid mutator transaction binding the contract method 0x2871ef85.
//
// Solidity: function testBASEFEE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestBASEFEE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestBASEFEE(&_LoadTester.TransactOpts, x)
}

// TestBASEFEE is a paid mutator transaction binding the contract method 0x2871ef85.
//
// Solidity: function testBASEFEE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestBASEFEE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestBASEFEE(&_LoadTester.TransactOpts, x)
}

// TestBLOCKHASH is a paid mutator transaction binding the contract method 0xea5141e6.
//
// Solidity: function testBLOCKHASH(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestBLOCKHASH(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testBLOCKHASH", x)
}

// TestBLOCKHASH is a paid mutator transaction binding the contract method 0xea5141e6.
//
// Solidity: function testBLOCKHASH(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestBLOCKHASH(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestBLOCKHASH(&_LoadTester.TransactOpts, x)
}

// TestBLOCKHASH is a paid mutator transaction binding the contract method 0xea5141e6.
//
// Solidity: function testBLOCKHASH(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestBLOCKHASH(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestBLOCKHASH(&_LoadTester.TransactOpts, x)
}

// TestBYTE is a paid mutator transaction binding the contract method 0x1de2f343.
//
// Solidity: function testBYTE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestBYTE(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testBYTE", x)
}

// TestBYTE is a paid mutator transaction binding the contract method 0x1de2f343.
//
// Solidity: function testBYTE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestBYTE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestBYTE(&_LoadTester.TransactOpts, x)
}

// TestBYTE is a paid mutator transaction binding the contract method 0x1de2f343.
//
// Solidity: function testBYTE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestBYTE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestBYTE(&_LoadTester.TransactOpts, x)
}

// TestBlake2f is a paid mutator transaction binding the contract method 0xa040aec6.
//
// Solidity: function testBlake2f(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterTransactor) TestBlake2f(opts *bind.TransactOpts, inputData []byte) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testBlake2f", inputData)
}

// TestBlake2f is a paid mutator transaction binding the contract method 0xa040aec6.
//
// Solidity: function testBlake2f(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterSession) TestBlake2f(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestBlake2f(&_LoadTester.TransactOpts, inputData)
}

// TestBlake2f is a paid mutator transaction binding the contract method 0xa040aec6.
//
// Solidity: function testBlake2f(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterTransactorSession) TestBlake2f(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestBlake2f(&_LoadTester.TransactOpts, inputData)
}

// TestCALLDATACOPY is a paid mutator transaction binding the contract method 0x3a425dfc.
//
// Solidity: function testCALLDATACOPY(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestCALLDATACOPY(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testCALLDATACOPY", x)
}

// TestCALLDATACOPY is a paid mutator transaction binding the contract method 0x3a425dfc.
//
// Solidity: function testCALLDATACOPY(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestCALLDATACOPY(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCALLDATACOPY(&_LoadTester.TransactOpts, x)
}

// TestCALLDATACOPY is a paid mutator transaction binding the contract method 0x3a425dfc.
//
// Solidity: function testCALLDATACOPY(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestCALLDATACOPY(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCALLDATACOPY(&_LoadTester.TransactOpts, x)
}

// TestCALLDATALOAD is a paid mutator transaction binding the contract method 0xce3cf4ef.
//
// Solidity: function testCALLDATALOAD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestCALLDATALOAD(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testCALLDATALOAD", x)
}

// TestCALLDATALOAD is a paid mutator transaction binding the contract method 0xce3cf4ef.
//
// Solidity: function testCALLDATALOAD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestCALLDATALOAD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCALLDATALOAD(&_LoadTester.TransactOpts, x)
}

// TestCALLDATALOAD is a paid mutator transaction binding the contract method 0xce3cf4ef.
//
// Solidity: function testCALLDATALOAD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestCALLDATALOAD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCALLDATALOAD(&_LoadTester.TransactOpts, x)
}

// TestCALLDATASIZE is a paid mutator transaction binding the contract method 0x034aef71.
//
// Solidity: function testCALLDATASIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestCALLDATASIZE(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testCALLDATASIZE", x)
}

// TestCALLDATASIZE is a paid mutator transaction binding the contract method 0x034aef71.
//
// Solidity: function testCALLDATASIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestCALLDATASIZE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCALLDATASIZE(&_LoadTester.TransactOpts, x)
}

// TestCALLDATASIZE is a paid mutator transaction binding the contract method 0x034aef71.
//
// Solidity: function testCALLDATASIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestCALLDATASIZE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCALLDATASIZE(&_LoadTester.TransactOpts, x)
}

// TestCALLER is a paid mutator transaction binding the contract method 0x44cf3bc7.
//
// Solidity: function testCALLER(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestCALLER(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testCALLER", x)
}

// TestCALLER is a paid mutator transaction binding the contract method 0x44cf3bc7.
//
// Solidity: function testCALLER(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestCALLER(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCALLER(&_LoadTester.TransactOpts, x)
}

// TestCALLER is a paid mutator transaction binding the contract method 0x44cf3bc7.
//
// Solidity: function testCALLER(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestCALLER(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCALLER(&_LoadTester.TransactOpts, x)
}

// TestCALLVALUE is a paid mutator transaction binding the contract method 0x1581cf19.
//
// Solidity: function testCALLVALUE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestCALLVALUE(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testCALLVALUE", x)
}

// TestCALLVALUE is a paid mutator transaction binding the contract method 0x1581cf19.
//
// Solidity: function testCALLVALUE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestCALLVALUE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCALLVALUE(&_LoadTester.TransactOpts, x)
}

// TestCALLVALUE is a paid mutator transaction binding the contract method 0x1581cf19.
//
// Solidity: function testCALLVALUE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestCALLVALUE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCALLVALUE(&_LoadTester.TransactOpts, x)
}

// TestCHAINID is a paid mutator transaction binding the contract method 0xa60a1087.
//
// Solidity: function testCHAINID(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestCHAINID(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testCHAINID", x)
}

// TestCHAINID is a paid mutator transaction binding the contract method 0xa60a1087.
//
// Solidity: function testCHAINID(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestCHAINID(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCHAINID(&_LoadTester.TransactOpts, x)
}

// TestCHAINID is a paid mutator transaction binding the contract method 0xa60a1087.
//
// Solidity: function testCHAINID(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestCHAINID(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCHAINID(&_LoadTester.TransactOpts, x)
}

// TestCODECOPY is a paid mutator transaction binding the contract method 0xacaebdf6.
//
// Solidity: function testCODECOPY(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestCODECOPY(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testCODECOPY", x)
}

// TestCODECOPY is a paid mutator transaction binding the contract method 0xacaebdf6.
//
// Solidity: function testCODECOPY(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestCODECOPY(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCODECOPY(&_LoadTester.TransactOpts, x)
}

// TestCODECOPY is a paid mutator transaction binding the contract method 0xacaebdf6.
//
// Solidity: function testCODECOPY(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestCODECOPY(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCODECOPY(&_LoadTester.TransactOpts, x)
}

// TestCODESIZE is a paid mutator transaction binding the contract method 0xb7b86207.
//
// Solidity: function testCODESIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestCODESIZE(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testCODESIZE", x)
}

// TestCODESIZE is a paid mutator transaction binding the contract method 0xb7b86207.
//
// Solidity: function testCODESIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestCODESIZE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCODESIZE(&_LoadTester.TransactOpts, x)
}

// TestCODESIZE is a paid mutator transaction binding the contract method 0xb7b86207.
//
// Solidity: function testCODESIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestCODESIZE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCODESIZE(&_LoadTester.TransactOpts, x)
}

// TestCOINBASE is a paid mutator transaction binding the contract method 0xb81c1484.
//
// Solidity: function testCOINBASE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestCOINBASE(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testCOINBASE", x)
}

// TestCOINBASE is a paid mutator transaction binding the contract method 0xb81c1484.
//
// Solidity: function testCOINBASE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestCOINBASE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCOINBASE(&_LoadTester.TransactOpts, x)
}

// TestCOINBASE is a paid mutator transaction binding the contract method 0xb81c1484.
//
// Solidity: function testCOINBASE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestCOINBASE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestCOINBASE(&_LoadTester.TransactOpts, x)
}

// TestDIFFICULTY is a paid mutator transaction binding the contract method 0x6f099c8d.
//
// Solidity: function testDIFFICULTY(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestDIFFICULTY(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testDIFFICULTY", x)
}

// TestDIFFICULTY is a paid mutator transaction binding the contract method 0x6f099c8d.
//
// Solidity: function testDIFFICULTY(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestDIFFICULTY(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestDIFFICULTY(&_LoadTester.TransactOpts, x)
}

// TestDIFFICULTY is a paid mutator transaction binding the contract method 0x6f099c8d.
//
// Solidity: function testDIFFICULTY(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestDIFFICULTY(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestDIFFICULTY(&_LoadTester.TransactOpts, x)
}

// TestDIV is a paid mutator transaction binding the contract method 0x3a411f12.
//
// Solidity: function testDIV(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestDIV(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testDIV", x)
}

// TestDIV is a paid mutator transaction binding the contract method 0x3a411f12.
//
// Solidity: function testDIV(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestDIV(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestDIV(&_LoadTester.TransactOpts, x)
}

// TestDIV is a paid mutator transaction binding the contract method 0x3a411f12.
//
// Solidity: function testDIV(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestDIV(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestDIV(&_LoadTester.TransactOpts, x)
}

// TestECAdd is a paid mutator transaction binding the contract method 0xedf003cf.
//
// Solidity: function testECAdd(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterTransactor) TestECAdd(opts *bind.TransactOpts, inputData []byte) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testECAdd", inputData)
}

// TestECAdd is a paid mutator transaction binding the contract method 0xedf003cf.
//
// Solidity: function testECAdd(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterSession) TestECAdd(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestECAdd(&_LoadTester.TransactOpts, inputData)
}

// TestECAdd is a paid mutator transaction binding the contract method 0xedf003cf.
//
// Solidity: function testECAdd(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterTransactorSession) TestECAdd(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestECAdd(&_LoadTester.TransactOpts, inputData)
}

// TestECMul is a paid mutator transaction binding the contract method 0x962e4dc2.
//
// Solidity: function testECMul(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterTransactor) TestECMul(opts *bind.TransactOpts, inputData []byte) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testECMul", inputData)
}

// TestECMul is a paid mutator transaction binding the contract method 0x962e4dc2.
//
// Solidity: function testECMul(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterSession) TestECMul(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestECMul(&_LoadTester.TransactOpts, inputData)
}

// TestECMul is a paid mutator transaction binding the contract method 0x962e4dc2.
//
// Solidity: function testECMul(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterTransactorSession) TestECMul(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestECMul(&_LoadTester.TransactOpts, inputData)
}

// TestECPairing is a paid mutator transaction binding the contract method 0x0b3b996a.
//
// Solidity: function testECPairing(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterTransactor) TestECPairing(opts *bind.TransactOpts, inputData []byte) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testECPairing", inputData)
}

// TestECPairing is a paid mutator transaction binding the contract method 0x0b3b996a.
//
// Solidity: function testECPairing(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterSession) TestECPairing(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestECPairing(&_LoadTester.TransactOpts, inputData)
}

// TestECPairing is a paid mutator transaction binding the contract method 0x0b3b996a.
//
// Solidity: function testECPairing(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterTransactorSession) TestECPairing(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestECPairing(&_LoadTester.TransactOpts, inputData)
}

// TestECRecover is a paid mutator transaction binding the contract method 0xa18683cb.
//
// Solidity: function testECRecover(bytes inputData) returns(address result)
func (_LoadTester *LoadTesterTransactor) TestECRecover(opts *bind.TransactOpts, inputData []byte) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testECRecover", inputData)
}

// TestECRecover is a paid mutator transaction binding the contract method 0xa18683cb.
//
// Solidity: function testECRecover(bytes inputData) returns(address result)
func (_LoadTester *LoadTesterSession) TestECRecover(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestECRecover(&_LoadTester.TransactOpts, inputData)
}

// TestECRecover is a paid mutator transaction binding the contract method 0xa18683cb.
//
// Solidity: function testECRecover(bytes inputData) returns(address result)
func (_LoadTester *LoadTesterTransactorSession) TestECRecover(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestECRecover(&_LoadTester.TransactOpts, inputData)
}

// TestEQ is a paid mutator transaction binding the contract method 0xe9f9b3f2.
//
// Solidity: function testEQ(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestEQ(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testEQ", x)
}

// TestEQ is a paid mutator transaction binding the contract method 0xe9f9b3f2.
//
// Solidity: function testEQ(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestEQ(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestEQ(&_LoadTester.TransactOpts, x)
}

// TestEQ is a paid mutator transaction binding the contract method 0xe9f9b3f2.
//
// Solidity: function testEQ(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestEQ(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestEQ(&_LoadTester.TransactOpts, x)
}

// TestEXP is a paid mutator transaction binding the contract method 0xde97a363.
//
// Solidity: function testEXP(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestEXP(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testEXP", x)
}

// TestEXP is a paid mutator transaction binding the contract method 0xde97a363.
//
// Solidity: function testEXP(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestEXP(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestEXP(&_LoadTester.TransactOpts, x)
}

// TestEXP is a paid mutator transaction binding the contract method 0xde97a363.
//
// Solidity: function testEXP(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestEXP(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestEXP(&_LoadTester.TransactOpts, x)
}

// TestEXTCODESIZE is a paid mutator transaction binding the contract method 0xf58fc36a.
//
// Solidity: function testEXTCODESIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestEXTCODESIZE(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testEXTCODESIZE", x)
}

// TestEXTCODESIZE is a paid mutator transaction binding the contract method 0xf58fc36a.
//
// Solidity: function testEXTCODESIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestEXTCODESIZE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestEXTCODESIZE(&_LoadTester.TransactOpts, x)
}

// TestEXTCODESIZE is a paid mutator transaction binding the contract method 0xf58fc36a.
//
// Solidity: function testEXTCODESIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestEXTCODESIZE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestEXTCODESIZE(&_LoadTester.TransactOpts, x)
}

// TestGAS is a paid mutator transaction binding the contract method 0x918a5fcd.
//
// Solidity: function testGAS(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestGAS(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testGAS", x)
}

// TestGAS is a paid mutator transaction binding the contract method 0x918a5fcd.
//
// Solidity: function testGAS(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestGAS(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestGAS(&_LoadTester.TransactOpts, x)
}

// TestGAS is a paid mutator transaction binding the contract method 0x918a5fcd.
//
// Solidity: function testGAS(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestGAS(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestGAS(&_LoadTester.TransactOpts, x)
}

// TestGASLIMIT is a paid mutator transaction binding the contract method 0x7c191d20.
//
// Solidity: function testGASLIMIT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestGASLIMIT(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testGASLIMIT", x)
}

// TestGASLIMIT is a paid mutator transaction binding the contract method 0x7c191d20.
//
// Solidity: function testGASLIMIT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestGASLIMIT(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestGASLIMIT(&_LoadTester.TransactOpts, x)
}

// TestGASLIMIT is a paid mutator transaction binding the contract method 0x7c191d20.
//
// Solidity: function testGASLIMIT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestGASLIMIT(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestGASLIMIT(&_LoadTester.TransactOpts, x)
}

// TestGASPRICE is a paid mutator transaction binding the contract method 0x4d2c74b3.
//
// Solidity: function testGASPRICE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestGASPRICE(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testGASPRICE", x)
}

// TestGASPRICE is a paid mutator transaction binding the contract method 0x4d2c74b3.
//
// Solidity: function testGASPRICE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestGASPRICE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestGASPRICE(&_LoadTester.TransactOpts, x)
}

// TestGASPRICE is a paid mutator transaction binding the contract method 0x4d2c74b3.
//
// Solidity: function testGASPRICE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestGASPRICE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestGASPRICE(&_LoadTester.TransactOpts, x)
}

// TestGT is a paid mutator transaction binding the contract method 0x71d91d28.
//
// Solidity: function testGT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestGT(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testGT", x)
}

// TestGT is a paid mutator transaction binding the contract method 0x71d91d28.
//
// Solidity: function testGT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestGT(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestGT(&_LoadTester.TransactOpts, x)
}

// TestGT is a paid mutator transaction binding the contract method 0x71d91d28.
//
// Solidity: function testGT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestGT(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestGT(&_LoadTester.TransactOpts, x)
}

// TestISZERO is a paid mutator transaction binding the contract method 0xf279ca81.
//
// Solidity: function testISZERO(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestISZERO(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testISZERO", x)
}

// TestISZERO is a paid mutator transaction binding the contract method 0xf279ca81.
//
// Solidity: function testISZERO(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestISZERO(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestISZERO(&_LoadTester.TransactOpts, x)
}

// TestISZERO is a paid mutator transaction binding the contract method 0xf279ca81.
//
// Solidity: function testISZERO(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestISZERO(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestISZERO(&_LoadTester.TransactOpts, x)
}

// TestIdentity is a paid mutator transaction binding the contract method 0x9cce7cf9.
//
// Solidity: function testIdentity(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterTransactor) TestIdentity(opts *bind.TransactOpts, inputData []byte) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testIdentity", inputData)
}

// TestIdentity is a paid mutator transaction binding the contract method 0x9cce7cf9.
//
// Solidity: function testIdentity(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterSession) TestIdentity(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestIdentity(&_LoadTester.TransactOpts, inputData)
}

// TestIdentity is a paid mutator transaction binding the contract method 0x9cce7cf9.
//
// Solidity: function testIdentity(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterTransactorSession) TestIdentity(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestIdentity(&_LoadTester.TransactOpts, inputData)
}

// TestLOG0 is a paid mutator transaction binding the contract method 0x40fe2662.
//
// Solidity: function testLOG0(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestLOG0(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testLOG0", x)
}

// TestLOG0 is a paid mutator transaction binding the contract method 0x40fe2662.
//
// Solidity: function testLOG0(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestLOG0(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestLOG0(&_LoadTester.TransactOpts, x)
}

// TestLOG0 is a paid mutator transaction binding the contract method 0x40fe2662.
//
// Solidity: function testLOG0(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestLOG0(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestLOG0(&_LoadTester.TransactOpts, x)
}

// TestLOG1 is a paid mutator transaction binding the contract method 0x98456f3e.
//
// Solidity: function testLOG1(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestLOG1(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testLOG1", x)
}

// TestLOG1 is a paid mutator transaction binding the contract method 0x98456f3e.
//
// Solidity: function testLOG1(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestLOG1(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestLOG1(&_LoadTester.TransactOpts, x)
}

// TestLOG1 is a paid mutator transaction binding the contract method 0x98456f3e.
//
// Solidity: function testLOG1(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestLOG1(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestLOG1(&_LoadTester.TransactOpts, x)
}

// TestLOG2 is a paid mutator transaction binding the contract method 0xdd9bef60.
//
// Solidity: function testLOG2(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestLOG2(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testLOG2", x)
}

// TestLOG2 is a paid mutator transaction binding the contract method 0xdd9bef60.
//
// Solidity: function testLOG2(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestLOG2(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestLOG2(&_LoadTester.TransactOpts, x)
}

// TestLOG2 is a paid mutator transaction binding the contract method 0xdd9bef60.
//
// Solidity: function testLOG2(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestLOG2(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestLOG2(&_LoadTester.TransactOpts, x)
}

// TestLOG3 is a paid mutator transaction binding the contract method 0xbf529ca1.
//
// Solidity: function testLOG3(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestLOG3(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testLOG3", x)
}

// TestLOG3 is a paid mutator transaction binding the contract method 0xbf529ca1.
//
// Solidity: function testLOG3(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestLOG3(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestLOG3(&_LoadTester.TransactOpts, x)
}

// TestLOG3 is a paid mutator transaction binding the contract method 0xbf529ca1.
//
// Solidity: function testLOG3(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestLOG3(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestLOG3(&_LoadTester.TransactOpts, x)
}

// TestLOG4 is a paid mutator transaction binding the contract method 0x1aba07ea.
//
// Solidity: function testLOG4(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestLOG4(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testLOG4", x)
}

// TestLOG4 is a paid mutator transaction binding the contract method 0x1aba07ea.
//
// Solidity: function testLOG4(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestLOG4(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestLOG4(&_LoadTester.TransactOpts, x)
}

// TestLOG4 is a paid mutator transaction binding the contract method 0x1aba07ea.
//
// Solidity: function testLOG4(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestLOG4(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestLOG4(&_LoadTester.TransactOpts, x)
}

// TestLT is a paid mutator transaction binding the contract method 0x6e7f1fe7.
//
// Solidity: function testLT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestLT(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testLT", x)
}

// TestLT is a paid mutator transaction binding the contract method 0x6e7f1fe7.
//
// Solidity: function testLT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestLT(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestLT(&_LoadTester.TransactOpts, x)
}

// TestLT is a paid mutator transaction binding the contract method 0x6e7f1fe7.
//
// Solidity: function testLT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestLT(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestLT(&_LoadTester.TransactOpts, x)
}

// TestMLOAD is a paid mutator transaction binding the contract method 0x5590c2d9.
//
// Solidity: function testMLOAD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestMLOAD(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testMLOAD", x)
}

// TestMLOAD is a paid mutator transaction binding the contract method 0x5590c2d9.
//
// Solidity: function testMLOAD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestMLOAD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestMLOAD(&_LoadTester.TransactOpts, x)
}

// TestMLOAD is a paid mutator transaction binding the contract method 0x5590c2d9.
//
// Solidity: function testMLOAD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestMLOAD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestMLOAD(&_LoadTester.TransactOpts, x)
}

// TestMOD is a paid mutator transaction binding the contract method 0x16582150.
//
// Solidity: function testMOD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestMOD(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testMOD", x)
}

// TestMOD is a paid mutator transaction binding the contract method 0x16582150.
//
// Solidity: function testMOD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestMOD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestMOD(&_LoadTester.TransactOpts, x)
}

// TestMOD is a paid mutator transaction binding the contract method 0x16582150.
//
// Solidity: function testMOD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestMOD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestMOD(&_LoadTester.TransactOpts, x)
}

// TestMSIZE is a paid mutator transaction binding the contract method 0xb3d847f2.
//
// Solidity: function testMSIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestMSIZE(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testMSIZE", x)
}

// TestMSIZE is a paid mutator transaction binding the contract method 0xb3d847f2.
//
// Solidity: function testMSIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestMSIZE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestMSIZE(&_LoadTester.TransactOpts, x)
}

// TestMSIZE is a paid mutator transaction binding the contract method 0xb3d847f2.
//
// Solidity: function testMSIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestMSIZE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestMSIZE(&_LoadTester.TransactOpts, x)
}

// TestMSTORE is a paid mutator transaction binding the contract method 0x087b4e84.
//
// Solidity: function testMSTORE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestMSTORE(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testMSTORE", x)
}

// TestMSTORE is a paid mutator transaction binding the contract method 0x087b4e84.
//
// Solidity: function testMSTORE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestMSTORE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestMSTORE(&_LoadTester.TransactOpts, x)
}

// TestMSTORE is a paid mutator transaction binding the contract method 0x087b4e84.
//
// Solidity: function testMSTORE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestMSTORE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestMSTORE(&_LoadTester.TransactOpts, x)
}

// TestMSTORE8 is a paid mutator transaction binding the contract method 0x4a61af1f.
//
// Solidity: function testMSTORE8(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestMSTORE8(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testMSTORE8", x)
}

// TestMSTORE8 is a paid mutator transaction binding the contract method 0x4a61af1f.
//
// Solidity: function testMSTORE8(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestMSTORE8(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestMSTORE8(&_LoadTester.TransactOpts, x)
}

// TestMSTORE8 is a paid mutator transaction binding the contract method 0x4a61af1f.
//
// Solidity: function testMSTORE8(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestMSTORE8(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestMSTORE8(&_LoadTester.TransactOpts, x)
}

// TestMUL is a paid mutator transaction binding the contract method 0x7de8c6f8.
//
// Solidity: function testMUL(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestMUL(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testMUL", x)
}

// TestMUL is a paid mutator transaction binding the contract method 0x7de8c6f8.
//
// Solidity: function testMUL(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestMUL(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestMUL(&_LoadTester.TransactOpts, x)
}

// TestMUL is a paid mutator transaction binding the contract method 0x7de8c6f8.
//
// Solidity: function testMUL(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestMUL(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestMUL(&_LoadTester.TransactOpts, x)
}

// TestMULMOD is a paid mutator transaction binding the contract method 0xfde7721c.
//
// Solidity: function testMULMOD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestMULMOD(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testMULMOD", x)
}

// TestMULMOD is a paid mutator transaction binding the contract method 0xfde7721c.
//
// Solidity: function testMULMOD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestMULMOD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestMULMOD(&_LoadTester.TransactOpts, x)
}

// TestMULMOD is a paid mutator transaction binding the contract method 0xfde7721c.
//
// Solidity: function testMULMOD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestMULMOD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestMULMOD(&_LoadTester.TransactOpts, x)
}

// TestModExp is a paid mutator transaction binding the contract method 0x613d0a82.
//
// Solidity: function testModExp(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterTransactor) TestModExp(opts *bind.TransactOpts, inputData []byte) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testModExp", inputData)
}

// TestModExp is a paid mutator transaction binding the contract method 0x613d0a82.
//
// Solidity: function testModExp(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterSession) TestModExp(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestModExp(&_LoadTester.TransactOpts, inputData)
}

// TestModExp is a paid mutator transaction binding the contract method 0x613d0a82.
//
// Solidity: function testModExp(bytes inputData) returns(bytes result)
func (_LoadTester *LoadTesterTransactorSession) TestModExp(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestModExp(&_LoadTester.TransactOpts, inputData)
}

// TestNOT is a paid mutator transaction binding the contract method 0x91e7b277.
//
// Solidity: function testNOT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestNOT(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testNOT", x)
}

// TestNOT is a paid mutator transaction binding the contract method 0x91e7b277.
//
// Solidity: function testNOT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestNOT(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestNOT(&_LoadTester.TransactOpts, x)
}

// TestNOT is a paid mutator transaction binding the contract method 0x91e7b277.
//
// Solidity: function testNOT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestNOT(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestNOT(&_LoadTester.TransactOpts, x)
}

// TestNUMBER is a paid mutator transaction binding the contract method 0x2d34e798.
//
// Solidity: function testNUMBER(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestNUMBER(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testNUMBER", x)
}

// TestNUMBER is a paid mutator transaction binding the contract method 0x2d34e798.
//
// Solidity: function testNUMBER(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestNUMBER(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestNUMBER(&_LoadTester.TransactOpts, x)
}

// TestNUMBER is a paid mutator transaction binding the contract method 0x2d34e798.
//
// Solidity: function testNUMBER(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestNUMBER(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestNUMBER(&_LoadTester.TransactOpts, x)
}

// TestOR is a paid mutator transaction binding the contract method 0x135d52f7.
//
// Solidity: function testOR(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestOR(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testOR", x)
}

// TestOR is a paid mutator transaction binding the contract method 0x135d52f7.
//
// Solidity: function testOR(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestOR(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestOR(&_LoadTester.TransactOpts, x)
}

// TestOR is a paid mutator transaction binding the contract method 0x135d52f7.
//
// Solidity: function testOR(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestOR(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestOR(&_LoadTester.TransactOpts, x)
}

// TestORIGIN is a paid mutator transaction binding the contract method 0x050082f8.
//
// Solidity: function testORIGIN(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestORIGIN(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testORIGIN", x)
}

// TestORIGIN is a paid mutator transaction binding the contract method 0x050082f8.
//
// Solidity: function testORIGIN(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestORIGIN(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestORIGIN(&_LoadTester.TransactOpts, x)
}

// TestORIGIN is a paid mutator transaction binding the contract method 0x050082f8.
//
// Solidity: function testORIGIN(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestORIGIN(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestORIGIN(&_LoadTester.TransactOpts, x)
}

// TestP256Verify is a paid mutator transaction binding the contract method 0xc711e539.
//
// Solidity: function testP256Verify(bytes inputData) returns(bool)
func (_LoadTester *LoadTesterTransactor) TestP256Verify(opts *bind.TransactOpts, inputData []byte) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testP256Verify", inputData)
}

// TestP256Verify is a paid mutator transaction binding the contract method 0xc711e539.
//
// Solidity: function testP256Verify(bytes inputData) returns(bool)
func (_LoadTester *LoadTesterSession) TestP256Verify(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestP256Verify(&_LoadTester.TransactOpts, inputData)
}

// TestP256Verify is a paid mutator transaction binding the contract method 0xc711e539.
//
// Solidity: function testP256Verify(bytes inputData) returns(bool)
func (_LoadTester *LoadTesterTransactorSession) TestP256Verify(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestP256Verify(&_LoadTester.TransactOpts, inputData)
}

// TestRETURNDATACOPY is a paid mutator transaction binding the contract method 0x7b6e0b0e.
//
// Solidity: function testRETURNDATACOPY(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestRETURNDATACOPY(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testRETURNDATACOPY", x)
}

// TestRETURNDATACOPY is a paid mutator transaction binding the contract method 0x7b6e0b0e.
//
// Solidity: function testRETURNDATACOPY(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestRETURNDATACOPY(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestRETURNDATACOPY(&_LoadTester.TransactOpts, x)
}

// TestRETURNDATACOPY is a paid mutator transaction binding the contract method 0x7b6e0b0e.
//
// Solidity: function testRETURNDATACOPY(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestRETURNDATACOPY(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestRETURNDATACOPY(&_LoadTester.TransactOpts, x)
}

// TestRETURNDATASIZE is a paid mutator transaction binding the contract method 0x2b21ef44.
//
// Solidity: function testRETURNDATASIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestRETURNDATASIZE(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testRETURNDATASIZE", x)
}

// TestRETURNDATASIZE is a paid mutator transaction binding the contract method 0x2b21ef44.
//
// Solidity: function testRETURNDATASIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestRETURNDATASIZE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestRETURNDATASIZE(&_LoadTester.TransactOpts, x)
}

// TestRETURNDATASIZE is a paid mutator transaction binding the contract method 0x2b21ef44.
//
// Solidity: function testRETURNDATASIZE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestRETURNDATASIZE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestRETURNDATASIZE(&_LoadTester.TransactOpts, x)
}

// TestRipemd160 is a paid mutator transaction binding the contract method 0xf6b0bbf7.
//
// Solidity: function testRipemd160(bytes inputData) returns(bytes20 result)
func (_LoadTester *LoadTesterTransactor) TestRipemd160(opts *bind.TransactOpts, inputData []byte) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testRipemd160", inputData)
}

// TestRipemd160 is a paid mutator transaction binding the contract method 0xf6b0bbf7.
//
// Solidity: function testRipemd160(bytes inputData) returns(bytes20 result)
func (_LoadTester *LoadTesterSession) TestRipemd160(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestRipemd160(&_LoadTester.TransactOpts, inputData)
}

// TestRipemd160 is a paid mutator transaction binding the contract method 0xf6b0bbf7.
//
// Solidity: function testRipemd160(bytes inputData) returns(bytes20 result)
func (_LoadTester *LoadTesterTransactorSession) TestRipemd160(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestRipemd160(&_LoadTester.TransactOpts, inputData)
}

// TestSAR is a paid mutator transaction binding the contract method 0x60e13cde.
//
// Solidity: function testSAR(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestSAR(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testSAR", x)
}

// TestSAR is a paid mutator transaction binding the contract method 0x60e13cde.
//
// Solidity: function testSAR(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestSAR(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSAR(&_LoadTester.TransactOpts, x)
}

// TestSAR is a paid mutator transaction binding the contract method 0x60e13cde.
//
// Solidity: function testSAR(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestSAR(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSAR(&_LoadTester.TransactOpts, x)
}

// TestSDIV is a paid mutator transaction binding the contract method 0xa645c9c2.
//
// Solidity: function testSDIV(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestSDIV(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testSDIV", x)
}

// TestSDIV is a paid mutator transaction binding the contract method 0xa645c9c2.
//
// Solidity: function testSDIV(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestSDIV(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSDIV(&_LoadTester.TransactOpts, x)
}

// TestSDIV is a paid mutator transaction binding the contract method 0xa645c9c2.
//
// Solidity: function testSDIV(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestSDIV(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSDIV(&_LoadTester.TransactOpts, x)
}

// TestSELFBALANCE is a paid mutator transaction binding the contract method 0xc420eb61.
//
// Solidity: function testSELFBALANCE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestSELFBALANCE(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testSELFBALANCE", x)
}

// TestSELFBALANCE is a paid mutator transaction binding the contract method 0xc420eb61.
//
// Solidity: function testSELFBALANCE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestSELFBALANCE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSELFBALANCE(&_LoadTester.TransactOpts, x)
}

// TestSELFBALANCE is a paid mutator transaction binding the contract method 0xc420eb61.
//
// Solidity: function testSELFBALANCE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestSELFBALANCE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSELFBALANCE(&_LoadTester.TransactOpts, x)
}

// TestSGT is a paid mutator transaction binding the contract method 0x18093b46.
//
// Solidity: function testSGT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestSGT(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testSGT", x)
}

// TestSGT is a paid mutator transaction binding the contract method 0x18093b46.
//
// Solidity: function testSGT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestSGT(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSGT(&_LoadTester.TransactOpts, x)
}

// TestSGT is a paid mutator transaction binding the contract method 0x18093b46.
//
// Solidity: function testSGT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestSGT(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSGT(&_LoadTester.TransactOpts, x)
}

// TestSHA256 is a paid mutator transaction binding the contract method 0x63138d4f.
//
// Solidity: function testSHA256(bytes inputData) returns(bytes32 result)
func (_LoadTester *LoadTesterTransactor) TestSHA256(opts *bind.TransactOpts, inputData []byte) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testSHA256", inputData)
}

// TestSHA256 is a paid mutator transaction binding the contract method 0x63138d4f.
//
// Solidity: function testSHA256(bytes inputData) returns(bytes32 result)
func (_LoadTester *LoadTesterSession) TestSHA256(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSHA256(&_LoadTester.TransactOpts, inputData)
}

// TestSHA256 is a paid mutator transaction binding the contract method 0x63138d4f.
//
// Solidity: function testSHA256(bytes inputData) returns(bytes32 result)
func (_LoadTester *LoadTesterTransactorSession) TestSHA256(inputData []byte) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSHA256(&_LoadTester.TransactOpts, inputData)
}

// TestSHA3 is a paid mutator transaction binding the contract method 0x19b621d6.
//
// Solidity: function testSHA3(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestSHA3(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testSHA3", x)
}

// TestSHA3 is a paid mutator transaction binding the contract method 0x19b621d6.
//
// Solidity: function testSHA3(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestSHA3(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSHA3(&_LoadTester.TransactOpts, x)
}

// TestSHA3 is a paid mutator transaction binding the contract method 0x19b621d6.
//
// Solidity: function testSHA3(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestSHA3(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSHA3(&_LoadTester.TransactOpts, x)
}

// TestSHL is a paid mutator transaction binding the contract method 0x2007332e.
//
// Solidity: function testSHL(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestSHL(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testSHL", x)
}

// TestSHL is a paid mutator transaction binding the contract method 0x2007332e.
//
// Solidity: function testSHL(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestSHL(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSHL(&_LoadTester.TransactOpts, x)
}

// TestSHL is a paid mutator transaction binding the contract method 0x2007332e.
//
// Solidity: function testSHL(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestSHL(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSHL(&_LoadTester.TransactOpts, x)
}

// TestSHR is a paid mutator transaction binding the contract method 0xc4bd65d5.
//
// Solidity: function testSHR(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestSHR(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testSHR", x)
}

// TestSHR is a paid mutator transaction binding the contract method 0xc4bd65d5.
//
// Solidity: function testSHR(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestSHR(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSHR(&_LoadTester.TransactOpts, x)
}

// TestSHR is a paid mutator transaction binding the contract method 0xc4bd65d5.
//
// Solidity: function testSHR(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestSHR(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSHR(&_LoadTester.TransactOpts, x)
}

// TestSIGNEXTEND is a paid mutator transaction binding the contract method 0xc360aba6.
//
// Solidity: function testSIGNEXTEND(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestSIGNEXTEND(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testSIGNEXTEND", x)
}

// TestSIGNEXTEND is a paid mutator transaction binding the contract method 0xc360aba6.
//
// Solidity: function testSIGNEXTEND(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestSIGNEXTEND(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSIGNEXTEND(&_LoadTester.TransactOpts, x)
}

// TestSIGNEXTEND is a paid mutator transaction binding the contract method 0xc360aba6.
//
// Solidity: function testSIGNEXTEND(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestSIGNEXTEND(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSIGNEXTEND(&_LoadTester.TransactOpts, x)
}

// TestSLOAD is a paid mutator transaction binding the contract method 0x880eff39.
//
// Solidity: function testSLOAD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestSLOAD(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testSLOAD", x)
}

// TestSLOAD is a paid mutator transaction binding the contract method 0x880eff39.
//
// Solidity: function testSLOAD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestSLOAD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSLOAD(&_LoadTester.TransactOpts, x)
}

// TestSLOAD is a paid mutator transaction binding the contract method 0x880eff39.
//
// Solidity: function testSLOAD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestSLOAD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSLOAD(&_LoadTester.TransactOpts, x)
}

// TestSLT is a paid mutator transaction binding the contract method 0xf4d1fc61.
//
// Solidity: function testSLT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestSLT(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testSLT", x)
}

// TestSLT is a paid mutator transaction binding the contract method 0xf4d1fc61.
//
// Solidity: function testSLT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestSLT(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSLT(&_LoadTester.TransactOpts, x)
}

// TestSLT is a paid mutator transaction binding the contract method 0xf4d1fc61.
//
// Solidity: function testSLT(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestSLT(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSLT(&_LoadTester.TransactOpts, x)
}

// TestSMOD is a paid mutator transaction binding the contract method 0xd93cd558.
//
// Solidity: function testSMOD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestSMOD(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testSMOD", x)
}

// TestSMOD is a paid mutator transaction binding the contract method 0xd93cd558.
//
// Solidity: function testSMOD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestSMOD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSMOD(&_LoadTester.TransactOpts, x)
}

// TestSMOD is a paid mutator transaction binding the contract method 0xd93cd558.
//
// Solidity: function testSMOD(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestSMOD(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSMOD(&_LoadTester.TransactOpts, x)
}

// TestSSTORE is a paid mutator transaction binding the contract method 0xd117320b.
//
// Solidity: function testSSTORE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestSSTORE(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testSSTORE", x)
}

// TestSSTORE is a paid mutator transaction binding the contract method 0xd117320b.
//
// Solidity: function testSSTORE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestSSTORE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSSTORE(&_LoadTester.TransactOpts, x)
}

// TestSSTORE is a paid mutator transaction binding the contract method 0xd117320b.
//
// Solidity: function testSSTORE(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestSSTORE(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSSTORE(&_LoadTester.TransactOpts, x)
}

// TestSUB is a paid mutator transaction binding the contract method 0xd53ff3fd.
//
// Solidity: function testSUB(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestSUB(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testSUB", x)
}

// TestSUB is a paid mutator transaction binding the contract method 0xd53ff3fd.
//
// Solidity: function testSUB(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestSUB(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSUB(&_LoadTester.TransactOpts, x)
}

// TestSUB is a paid mutator transaction binding the contract method 0xd53ff3fd.
//
// Solidity: function testSUB(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestSUB(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestSUB(&_LoadTester.TransactOpts, x)
}

// TestTIMESTAMP is a paid mutator transaction binding the contract method 0x219cddeb.
//
// Solidity: function testTIMESTAMP(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestTIMESTAMP(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testTIMESTAMP", x)
}

// TestTIMESTAMP is a paid mutator transaction binding the contract method 0x219cddeb.
//
// Solidity: function testTIMESTAMP(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestTIMESTAMP(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestTIMESTAMP(&_LoadTester.TransactOpts, x)
}

// TestTIMESTAMP is a paid mutator transaction binding the contract method 0x219cddeb.
//
// Solidity: function testTIMESTAMP(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestTIMESTAMP(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestTIMESTAMP(&_LoadTester.TransactOpts, x)
}

// TestXOR is a paid mutator transaction binding the contract method 0xd51e7b5b.
//
// Solidity: function testXOR(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactor) TestXOR(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _LoadTester.contract.Transact(opts, "testXOR", x)
}

// TestXOR is a paid mutator transaction binding the contract method 0xd51e7b5b.
//
// Solidity: function testXOR(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterSession) TestXOR(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestXOR(&_LoadTester.TransactOpts, x)
}

// TestXOR is a paid mutator transaction binding the contract method 0xd51e7b5b.
//
// Solidity: function testXOR(uint256 x) returns(uint256)
func (_LoadTester *LoadTesterTransactorSession) TestXOR(x *big.Int) (*types.Transaction, error) {
	return _LoadTester.Contract.TestXOR(&_LoadTester.TransactOpts, x)
}
