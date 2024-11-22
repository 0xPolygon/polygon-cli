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
	Bin: "0x608060405234801561001057600080fd5b506124d5806100206000396000f3fe608060405234801561001057600080fd5b506004361061046a5760003560e01c806380947f801161024c578063bf529ca111610146578063dd9bef60116100c3578063f279ca8111610087578063f279ca81146109f8578063f4d1fc6114610a0b578063f58fc36a14610a1e578063f6b0bbf714610a31578063fde7721c14610a5157600080fd5b8063dd9bef6014610999578063de97a363146109ac578063e9f9b3f2146109bf578063ea5141e6146109d2578063edf003cf146109e557600080fd5b8063ce3cf4ef1161010a578063ce3cf4ef1461093a578063d117320b1461094d578063d51e7b5b14610960578063d53ff3fd14610973578063d93cd5581461098657600080fd5b8063bf529ca1146108ce578063c360aba6146108e1578063c420eb61146108f4578063c4bd65d514610907578063c711e5391461091a57600080fd5b8063a18683cb116101d4578063b374012b11610198578063b374012b1461086f578063b3d847f214610882578063b7b8620714610895578063b81c1484146108a8578063bdc875fc146108bb57600080fd5b8063a18683cb1461080e578063a271b7211461082e578063a60a108714610836578063a645c9c214610849578063acaebdf61461085c57600080fd5b8063962e4dc21161021b578063962e4dc2146107af57806398456f3e146107c25780639a2b7c81146107d55780639cce7cf9146107e8578063a040aec6146107fb57600080fd5b806380947f8014610763578063880eff3914610776578063918a5fcd1461078957806391e7b2771461079c57600080fd5b80633430ec0611610368578063613d0a82116102e557806371d91d28116102a957806371d91d28146106f757806372de3cbd1461070a5780637b6e0b0e1461072a5780637c191d201461073d5780637de8c6f81461075057600080fd5b8063613d0a82146106a357806363138d4f146106b6578063659bbb4f146106c95780636e7f1fe7146106d15780636f099c8d146106e457600080fd5b806344cf3bc71161032c57806344cf3bc7146106445780634a61af1f146106575780634d2c74b31461066a5780635590c2d91461067d57806360e13cde1461069057600080fd5b80633430ec06146105f0578063371303c0146106035780633a411f121461060b5780633a425dfc1461061e57806340fe26621461063157600080fd5b806318093b46116103f6578063219cddeb116103ba578063219cddeb146105915780632294fc7f146105a45780632871ef85146105b75780632b21ef44146105ca5780632d34e798146105dd57600080fd5b806318093b461461053257806319b621d6146105455780631aba07ea146105585780631de2f3431461056b5780632007332e1461057e57600080fd5b80630ba8a73b1161043d5780630ba8a73b146104de5780631287a68c146104f1578063135d52f7146104f95780631581cf191461050c578063165821501461051f57600080fd5b8063034aef711461046f578063050082f814610498578063087b4e84146104ab5780630b3b996a146104be575b600080fd5b61048261047d366004611a7e565b610a64565b60405161048f9190611aaf565b60405180910390f35b6104826104a6366004611a7e565b610a96565b6104826104b9366004611a7e565b610abf565b6104d16104cc366004611bb4565b610af0565b60405161048f9190611c4b565b6104826104ec366004611a7e565b610b13565b600054610482565b610482610507366004611a7e565b610b38565b61048261051a366004611a7e565b610b5a565b61048261052d366004611a7e565b610b83565b610482610540366004611a7e565b610baf565b610482610553366004611a7e565b610bda565b610482610566366004611a7e565b610c46565b610482610579366004611a7e565b610c7c565b61048261058c366004611a7e565b610ca9565b61048261059f366004611a7e565b610ccb565b6104826105b2366004611a7e565b610cf4565b6104826105c5366004611a7e565b610d29565b6104826105d8366004611a7e565b610d52565b6104826105eb366004611a7e565b610d7b565b6104d16105fe366004611a7e565b610da4565b610482610e4d565b610482610619366004611a7e565b610e66565b61048261062c366004611a7e565b610e88565b61048261063f366004611a7e565b610eb3565b610482610652366004611a7e565b610ee2565b610482610665366004611a7e565b610f0b565b610482610678366004611a7e565b610f38565b61048261068b366004611a7e565b610f61565b61048261069e366004611a7e565b610f97565b6104d16106b1366004611bb4565b610fc3565b6104826106c4366004611bb4565b610fee565b610482611017565b6104826106df366004611a7e565b611051565b6104826106f2366004611a7e565b61107c565b610482610705366004611a7e565b6110a5565b61071d610718366004611e17565b6110d0565b60405161048f9190611ed9565b610482610738366004611a7e565b611152565b61048261074b366004611a7e565b61117d565b61048261075e366004611a7e565b6111a6565b610482610771366004611a7e565b6111c8565b610482610784366004611a7e565b6111f6565b610482610797366004611a7e565b611223565b6104826107aa366004611a7e565b61124c565b6104d16107bd366004611bb4565b61128a565b6104826107d0366004611a7e565b6112db565b6104826107e3366004611a7e565b61130b565b6104d16107f6366004611bb4565b61132d565b6104d1610809366004611bb4565b61134d565b61082161081c366004611bb4565b611488565b60405161048f9190611f01565b6104826114e2565b610482610844366004611a7e565b611523565b610482610857366004611a7e565b61154c565b61048261086a366004611a7e565b61156e565b61048261087d366004611f61565b611599565b610482610890366004611a7e565b6115c8565b6104826108a3366004611a7e565b6115f1565b6104826108b6366004611a7e565b61161a565b6104826108c9366004611a7e565b611643565b6104826108dc366004611a7e565b61166c565b6104826108ef366004611a7e565b6116a0565b610482610902366004611a7e565b6116c2565b610482610915366004611a7e565b6116eb565b61092d610928366004611bb4565b611711565b60405161048f9190611fb1565b610482610948366004611a7e565b611798565b61048261095b366004611a7e565b6117c3565b61048261096e366004611a7e565b6117ed565b610482610981366004611a7e565b61180f565b610482610994366004611a7e565b611831565b6104826109a7366004611a7e565b61185d565b6104826109ba366004611a7e565b61188f565b6104826109cd366004611a7e565b6118b9565b6104826109e0366004611a7e565b6118e2565b6104d16109f3366004611bb4565b611911565b610482610a06366004611a7e565b611950565b610482610a19366004611a7e565b61197a565b610482610a2c366004611a7e565b6119a5565b610a44610a3f366004611bb4565b6119d0565b60405161048f9190611fd4565b610482610a5f366004611a7e565b6119fe565b6000610a6e610e4d565b5065deadbeef00366000805b84811015610a8d57369150600101610a7a565b50909392505050565b6000610aa0610e4d565b5065deadbeef00326000805b84811015610a8d57329150600101610aac565b6000610ac9610e4d565b5065deadbeef005260005b83811015610ae9576000829052600101610ad4565b5092915050565b606060086040828451602086016000855af180610b0c57600080fd5b5050919050565b6000610b1d610e4d565b5065deadbeef000160005b83811015610ae957600101610b28565b6000610b42610e4d565b5065deadbeef001760008315610ae957600101610b28565b6000610b64610e4d565b5065deadbeef00346000805b84811015610a8d57349150600101610b70565b6000610b8d610e4d565b5065deadbeef000660005b83811015610ae95760001990910690600101610b98565b6000610bb9610e4d565b5065deadbeef00136000805b84811015610a8d576001808413925001610bc5565b6000610be4610e4d565b506001600160e01b0319600090815265deadbeef002090805b84811015610c145760046000209150600101610bfd565b507f29045a592007d0c246ef02c2223570da9522d0cf0f73282c79a1bc8f0bb2c2388114610ae9575060009392505050565b6000610c50610e4d565b5065deadbeef00a4601081905260005b83811015610ae9576004600360028360066010a4600101610c60565b6000610c86610e4d565b5065deadbeef001a6000805b84811015610a8d57600083901a9150600101610c92565b6000610cb3610e4d565b5065deadbeef001b60008315610ae957600101610b28565b6000610cd5610e4d565b5065deadbeef00426000805b84811015610a8d57429150600101610ce1565b6000610cfe610e4d565b5065deadbeef0031600030815b85811015610d1f5781319250600101610d0b565b5091949350505050565b6000610d33610e4d565b5065deadbeef00486000805b84811015610a8d57489150600101610d3f565b6000610d5c610e4d565b5065deadbeef003d6000805b84811015610a8d573d9150600101610d68565b6000610d85610e4d565b5065deadbeef00436000805b84811015610a8d57439150600101610d91565b60028181548110610db457600080fd5b906000526020600020018054909150610dcc90611ff8565b80601f0160208091040260200160405190810160405280929190818152602001828054610df890611ff8565b8015610e455780601f10610e1a57610100808354040283529160200191610e45565b820191906000526020600020905b815481529060010190602001808311610e2857829003601f168201915b505050505081565b60008054610e5c90600161203a565b6000819055919050565b6000610e70610e4d565b5065deadbeef000460008315610ae957600101610b28565b6000610e92610e4d565b5065deadbeef003760005b83811015610ae957602060008037600101610e9d565b6000610ebd610e4d565b5065deadbeef00a0601081905260005b83811015610ae95760066010a0600101610ecd565b6000610eec610e4d565b5065deadbeef00336000805b84811015610a8d57339150600101610ef8565b6000610f15610e4d565b5065deadbeef005360005b83811015610ae95763deadbeef600052600101610f20565b6000610f42610e4d565b5065deadbeef003a6000805b84811015610a8d573a9150600101610f4e565b6000610f6b610e4d565b5065deadbeef00516000818152805b84811015610f8f576000519150600101610f7a565b509392505050565b6000610fa1610e4d565b5065deadbeef001d60005b83811015610ae95760009190911d90600101610fac565b6060600560208301835160405160208183856000885af180610fe457600080fd5b5095945050505050565b600060026020830183518360208183856000885af18061100d57600080fd5b5050505050919050565b6000611021610e4d565b505b6103e85a111561104a57600180600082825461103f919061203a565b909155506110239050565b5060015490565b600061105b610e4d565b5065deadbeef00106000805b84811015610a8d576001838110925001611067565b6000611086610e4d565b5065deadbeef00446000805b84811015610a8d57449150600101611092565b60006110af610e4d565b5065deadbeef00116000805b84811015610a8d5760018084119250016110bb565b6110d8611a2c565b60096110e2611a2c565b600088885160208a0151895160208b015160408c015160608d01518c5160208e01518d60405160200161111e9a999897969594939291906120a3565b604051601f19818303018152604091825291508260d56020840186600019fa61114657600080fd5b50979650505050505050565b600061115c610e4d565b5065deadbeef003e60005b83811015610ae95760206000803e600101611167565b6000611187610e4d565b5065deadbeef00456000805b84811015610a8d57459150600101611193565b60006111b0610e4d565b5065deadbeef000260008315610ae957600101610b28565b60006111d2610e4d565b5065deadbeef000860005b83811015610ae9576000196000830891506001016111dd565b6000611200610e4d565b5065deadbeef005460008181555b83811015610ae957600054915060010161120e565b600061122d610e4d565b5065deadbeef005a6000805b84811015610a8d575a9150600101611239565b6000611256610e4d565b5065deadbeef001960005b8381101561127457901990600101611261565b5065deadbeef0019811461128457195b92915050565b606081516060146112b65760405162461bcd60e51b81526004016112ad9061217f565b60405180910390fd5b600760208301835160408482846000875af1806112d257600080fd5b50505050919050565b60006112e5610e4d565b5065deadbeef00a1601081905260005b83811015610ae9578060066010a16001016112f5565b6000611315610e4d565b5065deadbeef001660008315610ae957600101610b28565b60606004602083018351604051818183856000885af180610fe457600080fd5b6060611357611a2c565b7f48c9bdf267e6096a3ba7ca8485ae67bb2bf894fe72f36e3cf1361d5f3af54fa581527fd182e6ad7f520e511f6c3e2b8c68059b6bbd41fbabd9831f79217e1319cde05b60208201526113a8611a47565b6261626360e81b81526000602082018190526040820181905260608201526113ce611a2c565b600360f81b81526000602082018190526113ed600c85858560016110d0565b90506113f7611a2c565b7fba80a53f981c4d0d6a2797b69f12f6e94c212f14685ac4b74b12bb6fdbffa2d181527f7d87c5392aab792dc252d5de4533cc9518d38aa8dbf1925ab92386edd4009923602082015280518251146114615760405162461bcd60e51b81526004016112ad906121d5565b602081015160208301511461100d5760405162461bcd60e51b81526004016112ad90612229565b600081516080146114ab5760405162461bcd60e51b81526004016112ad9061226d565b6001602083016040840151601f1a602082015260206040516080836000865af1806114d557600080fd5b6040515195945050505050565b60006114ec610e4d565b505b6103e85a111561104a57600180600082825461150a919061203a565b909155505060015461151d904390612293565b506114ee565b600061152d610e4d565b5065deadbeef00466000805b84811015610a8d57469150600101611539565b6000611556610e4d565b5065deadbeef000560008315610ae957600101610b28565b6000611578610e4d565b5065deadbeef003960005b83811015610ae957602060008039600101611583565b600280546001810182556000918252839083906020842001916115bd919083612354565b505060025492915050565b60006115d2610e4d565b5065deadbeef00596000805b84811015610a8d575991506001016115de565b60006115fb610e4d565b5065deadbeef00386000805b84811015610a8d57389150600101611607565b6000611624610e4d565b5065deadbeef00416000805b84811015610a8d57419150600101611630565b600061164d610e4d565b5065deadbeef00306000805b84811015610a8d57309150600101611659565b6000611676610e4d565b5065deadbeef00a3601081905260005b83811015610ae957600360028260066010a3600101611686565b60006116aa610e4d565b5065deadbeef000b60008315610ae957600101610b28565b60006116cc610e4d565b5065deadbeef00476000805b84811015610a8d574791506001016116d8565b60006116f5610e4d565b5065deadbeef001c6000805b84811015610a8d57600101611701565b600061010081808285604051611727919061244d565b600060405180830381855afa9150503d8060008114611762576040513d603f01601f191681016040523d815291503d6000602084013e611767565b606091505b50915091508161177957611779612459565b60208101815181019061178c919061247e565b60011495945050505050565b60006117a2610e4d565b5065deadbeef00356000805b84811015610a8d5760003591506001016117ae565b60006117cd610e4d565b5065deadbeef005560005b83811015610ae95760008290556001016117d8565b60006117f7610e4d565b5065deadbeef001860008315610ae957600101610b28565b6000611819610e4d565b5065deadbeef000360008315610ae957600101610b28565b600061183b610e4d565b5065deadbeef000760005b83811015610ae95760001990910790600101611846565b6000611867610e4d565b5065deadbeef00a2601081905260005b83811015610ae95760028160066010a2600101611877565b6000611899610e4d565b5065deadbeef000a60005b83811015610ae95760019182900a91016118a4565b60006118c3610e4d565b5065deadbeef00146000805b84811015610a8d576001915081016118cf565b60006118ec610e4d565b5065deadbeef004060006000194301815b85811015610d1f57814092506001016118fd565b606081516080146119345760405162461bcd60e51b81526004016112ad9061217f565b600660208301835160408482846000875af1806112d257600080fd5b600061195a610e4d565b5065deadbeef00156000805b84811015610a8d5782159150600101611966565b6000611984610e4d565b5065deadbeef00126000805b84811015610a8d576001838112925001611990565b60006119af610e4d565b5065deadbeef003b600030815b85811015610d1f57813b92506001016119bc565b6000600360208301835160405160148183856000885af1806119f157600080fd5b8151979650505050505050565b6000611a08610e4d565b5065deadbeef000960005b83811015610ae957600019600183099150600101611a13565b60405160408082018152600290829080368337509192915050565b6040516080808201604052600490829080368337509192915050565b805b8114611a7057600080fd5b50565b803561128481611a63565b600060208284031215611a9357611a93600080fd5b6000611a9f8484611a73565b949350505050565b805b82525050565b602081016112848284611aa7565b634e487b7160e01b600052604160045260246000fd5b601f19601f830116810181811067ffffffffffffffff82111715611af957611af9611abd565b6040525050565b6000611b0f6000604051905090565b9050611b1b8282611ad3565b919050565b600067ffffffffffffffff821115611b3a57611b3a611abd565b601f19601f83011660200192915050565b82818337506000910152565b6000611b6a611b6584611b20565b611b00565b905082815260208101848484011115611b8557611b85600080fd5b610f8f848285611b4b565b600082601f830112611ba457611ba4600080fd5b8135611a9f848260208601611b57565b600060208284031215611bc957611bc9600080fd5b813567ffffffffffffffff811115611be357611be3600080fd5b611a9f84828501611b90565b60005b83811015611c0a578082015183820152602001611bf2565b50506000910152565b6000611c23826000815192915050565b808452602084019350611c3a818560208601611bef565b601f01601f19169290920192915050565b60208082528101611c5c8184611c13565b9392505050565b63ffffffff8116611a65565b803561128481611c63565b600067ffffffffffffffff821115611c9457611c94611abd565b5060200290565b6000611ca9611b6584611c7a565b90508060208402830185811115611cc257611cc2600080fd5b835b81811015611ce65780611cd78882611a73565b84525060209283019201611cc4565b5050509392505050565b600082601f830112611d0457611d04600080fd5b6002611a9f848285611c9b565b6000611d1f611b6584611c7a565b90508060208402830185811115611d3857611d38600080fd5b835b81811015611ce65780611d4d8882611a73565b84525060209283019201611d3a565b600082601f830112611d7057611d70600080fd5b6004611a9f848285611d11565b6001600160c01b03198116611a65565b803561128481611d7d565b6000611da6611b6584611c7a565b90508060208402830185811115611dbf57611dbf600080fd5b835b81811015611ce65780611dd48882611d8d565b84525060209283019201611dc1565b600082601f830112611df757611df7600080fd5b6002611a9f848285611d98565b801515611a65565b803561128481611e04565b60008060008060006101408688031215611e3357611e33600080fd5b6000611e3f8888611c6f565b9550506020611e5088828901611cf0565b9450506060611e6188828901611d5c565b93505060e0611e7288828901611de3565b925050610120611e8488828901611e0c565b9150509295509295909350565b6000611e9d8383611aa7565b505060200190565b6002818060005b83811015611ed1578151611ec08782611e91565b965060208301925050600101611eac565b505050505050565b604081016112848284611ea5565b60006001600160a01b038216611284565b611aa981611ee7565b602081016112848284611ef8565b60008083601f840112611f2457611f24600080fd5b50813567ffffffffffffffff811115611f3f57611f3f600080fd5b602083019150836001820283011115611f5a57611f5a600080fd5b9250929050565b60008060208385031215611f7757611f77600080fd5b823567ffffffffffffffff811115611f9157611f91600080fd5b611f9d85828601611f0f565b92509250509250929050565b801515611aa9565b602081016112848284611fa9565b6bffffffffffffffffffffffff198116611aa9565b602081016112848284611fbf565b634e487b7160e01b600052602260045260246000fd5b60028104600182168061200c57607f821691505b60208210810361201e5761201e611fe2565b50919050565b634e487b7160e01b600052601160045260246000fd5b8082018082111561128457611284612024565b60006112848260e01b90565b611aa963ffffffff821661204d565b80611aa9565b90565b6001600160c01b03198116611aa9565b60006112848260f81b90565b600061128482612081565b611aa981151561208d565b60006120af828d612059565b6004820191506120bf828c612068565b6020820191506120cf828b612068565b6020820191506120df828a612068565b6020820191506120ef8289612068565b6020820191506120ff8288612068565b60208201915061210f8287612068565b60208201915061211f8286612071565b60088201915061212f8285612071565b60088201915061213f8284612098565b506001019a9950505050505050505050565b6014815260006020820173092dcecc2d8d2c840d2dce0eae840d8cadccee8d60631b815291505b5060200190565b6020808252810161128481612151565b602681526000602082017f54657374426c616b653266202d204669727374206861736820646f65736e2774815265040dac2e8c6d60d31b602082015291505b5060400190565b602080825281016112848161218f565b602781526000602082017f54657374426c616b653266202d205365636f6e64206861736820646f65736e278152660e840dac2e8c6d60cb1b602082015291506121ce565b60208082528101611284816121e5565b601a81526000602082017f496e76616c696420696e7075742064617461206c656e6774682e00000000000081529150612178565b6020808252810161128481612239565b634e487b7160e01b600052601260045260246000fd5b6000826122a2576122a261227d565b500690565b600061128461206e8381565b6122bc836122a7565b815460001960089490940293841b1916921b91909117905550565b60006122e48184846122b3565b505050565b81811015612304576122fc6000826122d7565b6001016122e9565b5050565b601f8211156122e45761232681600081815281906020902092915050565b6020601f8501048101602085101561233b5750805b61234d6020601f8601048301826122e9565b5050505050565b8267ffffffffffffffff81111561236d5761236d611abd565b6123778254611ff8565b612382828285612308565b6000601f8311600181146123b6576000841561239e5750858201355b600019600886021c198116600286021786555061241c565b601f1984166123d086600081815281906020902092915050565b60005b828110156123f357888501358255602094850194600190920191016123d3565b8683101561240f57600019601f88166008021c19858a01351682555b6001600288020188555050505b50505050505050565b6000612435826000815192915050565b612443818560208601611bef565b9290920192915050565b6000611c5c8284612425565b634e487b7160e01b600052600160045260246000fd5b60008151905061128481611a63565b60006020828403121561249357612493600080fd5b6000611a9f848461246f56fea2646970667358221220b21578c01b69894c3c314dadf93efcfd571584a59528c1a97d77099f311bc4f464736f6c63430008170033",
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
