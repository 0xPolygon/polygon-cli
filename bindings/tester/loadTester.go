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
	Bin: "0x608060405234801561000f575f80fd5b50613b608061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610468575f3560e01c806380947f801161024a578063bf529ca111610144578063dd9bef60116100c1578063f279ca8111610085578063f279ca81146111d4578063f4d1fc6114611204578063f58fc36a14611234578063f6b0bbf714611264578063fde7721c1461129457610468565b8063dd9bef60146110e4578063de97a36314611114578063e9f9b3f214611144578063ea5141e614611174578063edf003cf146111a457610468565b8063ce3cf4ef11610108578063ce3cf4ef14610ff4578063d117320b14611024578063d51e7b5b14611054578063d53ff3fd14611084578063d93cd558146110b457610468565b8063bf529ca114610f04578063c360aba614610f34578063c420eb6114610f64578063c4bd65d514610f94578063c711e53914610fc457610468565b8063a18683cb116101d2578063b374012b11610196578063b374012b14610e14578063b3d847f214610e44578063b7b8620714610e74578063b81c148414610ea4578063bdc875fc14610ed457610468565b8063a18683cb14610d36578063a271b72114610d66578063a60a108714610d84578063a645c9c214610db4578063acaebdf614610de457610468565b8063962e4dc211610219578063962e4dc214610c4657806398456f3e14610c765780639a2b7c8114610ca65780639cce7cf914610cd6578063a040aec614610d0657610468565b806380947f8014610b86578063880eff3914610bb6578063918a5fcd14610be657806391e7b27714610c1657610468565b80633430ec0611610366578063613d0a82116102e357806371d91d28116102a757806371d91d2814610a9657806372de3cbd14610ac65780637b6e0b0e14610af65780637c191d2014610b265780637de8c6f814610b5657610468565b8063613d0a82146109b857806363138d4f146109e8578063659bbb4f14610a185780636e7f1fe714610a365780636f099c8d14610a6657610468565b806344cf3bc71161032a57806344cf3bc7146108c85780634a61af1f146108f85780634d2c74b3146109285780635590c2d91461095857806360e13cde1461098857610468565b80633430ec06146107ea578063371303c01461081a5780633a411f12146108385780633a425dfc1461086857806340fe26621461089857610468565b806318093b46116103f4578063219cddeb116103b8578063219cddeb146106fa5780632294fc7f1461072a5780632871ef851461075a5780632b21ef441461078a5780632d34e798146107ba57610468565b806318093b461461060a57806319b621d61461063a5780631aba07ea1461066a5780631de2f3431461069a5780632007332e146106ca57610468565b80630ba8a73b1161043b5780630ba8a73b1461052c5780631287a68c1461055c578063135d52f71461057a5780631581cf19146105aa57806316582150146105da57610468565b8063034aef711461046c578063050082f81461049c578063087b4e84146104cc5780630b3b996a146104fc575b5f80fd5b61048660048036038101906104819190612b01565b6112c4565b6040516104939190612b3b565b60405180910390f35b6104b660048036038101906104b19190612b01565b6112fc565b6040516104c39190612b3b565b60405180910390f35b6104e660048036038101906104e19190612b01565b611334565b6040516104f39190612b3b565b60405180910390f35b61051660048036038101906105119190612c90565b61136a565b6040516105239190612d51565b60405180910390f35b61054660048036038101906105419190612b01565b61138e565b6040516105539190612b3b565b60405180910390f35b6105646113c6565b6040516105719190612b3b565b60405180910390f35b610594600480360381019061058f9190612b01565b6113ce565b6040516105a19190612b3b565b60405180910390f35b6105c460048036038101906105bf9190612b01565b611406565b6040516105d19190612b3b565b60405180910390f35b6105f460048036038101906105ef9190612b01565b61143e565b6040516106019190612b3b565b60405180910390f35b610624600480360381019061061f9190612b01565b611496565b6040516106319190612b3b565b60405180910390f35b610654600480360381019061064f9190612b01565b6114d1565b6040516106619190612b3b565b60405180910390f35b610684600480360381019061067f9190612b01565b61155a565b6040516106919190612b3b565b60405180910390f35b6106b460048036038101906106af9190612b01565b61159d565b6040516106c19190612b3b565b60405180910390f35b6106e460048036038101906106df9190612b01565b6115d7565b6040516106f19190612b3b565b60405180910390f35b610714600480360381019061070f9190612b01565b61160f565b6040516107219190612b3b565b60405180910390f35b610744600480360381019061073f9190612b01565b611647565b6040516107519190612b3b565b60405180910390f35b610774600480360381019061076f9190612b01565b611682565b6040516107819190612b3b565b60405180910390f35b6107a4600480360381019061079f9190612b01565b6116ba565b6040516107b19190612b3b565b60405180910390f35b6107d460048036038101906107cf9190612b01565b6116f2565b6040516107e19190612b3b565b60405180910390f35b61080460048036038101906107ff9190612b01565b61172a565b6040516108119190612d51565b60405180910390f35b6108226117d0565b60405161082f9190612b3b565b60405180910390f35b610852600480360381019061084d9190612b01565b6117eb565b60405161085f9190612b3b565b60405180910390f35b610882600480360381019061087d9190612b01565b611824565b60405161088f9190612b3b565b60405180910390f35b6108b260048036038101906108ad9190612b01565b61185c565b6040516108bf9190612b3b565b60405180910390f35b6108e260048036038101906108dd9190612b01565b611898565b6040516108ef9190612b3b565b60405180910390f35b610912600480360381019061090d9190612b01565b6118d0565b60405161091f9190612b3b565b60405180910390f35b610942600480360381019061093d9190612b01565b61190a565b60405161094f9190612b3b565b60405180910390f35b610972600480360381019061096d9190612b01565b611942565b60405161097f9190612b3b565b60405180910390f35b6109a2600480360381019061099d9190612b01565b611981565b6040516109af9190612b3b565b60405180910390f35b6109d260048036038101906109cd9190612c90565b6119b9565b6040516109df9190612d51565b60405180910390f35b610a0260048036038101906109fd9190612c90565b6119e8565b604051610a0f9190612d89565b60405180910390f35b610a20611a11565b604051610a2d9190612b3b565b60405180910390f35b610a506004803603810190610a4b9190612b01565b611a4b565b604051610a5d9190612b3b565b60405180910390f35b610a806004803603810190610a7b9190612b01565b611a86565b604051610a8d9190612b3b565b60405180910390f35b610ab06004803603810190610aab9190612b01565b611abe565b604051610abd9190612b3b565b60405180910390f35b610ae06004803603810190610adb919061309d565b611af9565b604051610aed91906131bb565b60405180910390f35b610b106004803603810190610b0b9190612b01565b611c24565b604051610b1d9190612b3b565b60405180910390f35b610b406004803603810190610b3b9190612b01565b611c5c565b604051610b4d9190612b3b565b60405180910390f35b610b706004803603810190610b6b9190612b01565b611c94565b604051610b7d9190612b3b565b60405180910390f35b610ba06004803603810190610b9b9190612b01565b611ccd565b604051610bad9190612b3b565b60405180910390f35b610bd06004803603810190610bcb9190612b01565b611d26565b604051610bdd9190612b3b565b60405180910390f35b610c006004803603810190610bfb9190612b01565b611d60565b604051610c0d9190612b3b565b60405180910390f35b610c306004803603810190610c2b9190612b01565b611d98565b604051610c3d9190612b3b565b60405180910390f35b610c606004803603810190610c5b9190612c90565b611de1565b604051610c6d9190612d51565b60405180910390f35b610c906004803603810190610c8b9190612b01565b611e4c565b604051610c9d9190612b3b565b60405180910390f35b610cc06004803603810190610cbb9190612b01565b611e89565b604051610ccd9190612b3b565b60405180910390f35b610cf06004803603810190610ceb9190612c90565b611ec1565b604051610cfd9190612d51565b60405180910390f35b610d206004803603810190610d1b9190612c90565b611eef565b604051610d2d9190612d51565b60405180910390f35b610d506004803603810190610d4b9190612c90565b612267565b604051610d5d9190613213565b60405180910390f35b610d6e6122e5565b604051610d7b9190612b3b565b60405180910390f35b610d9e6004803603810190610d999190612b01565b61232e565b604051610dab9190612b3b565b60405180910390f35b610dce6004803603810190610dc99190612b01565b612366565b604051610ddb9190612b3b565b60405180910390f35b610dfe6004803603810190610df99190612b01565b61239f565b604051610e0b9190612b3b565b60405180910390f35b610e2e6004803603810190610e299190613285565b6123d7565b604051610e3b9190612b3b565b60405180910390f35b610e5e6004803603810190610e599190612b01565b612421565b604051610e6b9190612b3b565b60405180910390f35b610e8e6004803603810190610e899190612b01565b612459565b604051610e9b9190612b3b565b60405180910390f35b610ebe6004803603810190610eb99190612b01565b612491565b604051610ecb9190612b3b565b60405180910390f35b610eee6004803603810190610ee99190612b01565b6124c9565b604051610efb9190612b3b565b60405180910390f35b610f1e6004803603810190610f199190612b01565b612501565b604051610f2b9190612b3b565b60405180910390f35b610f4e6004803603810190610f499190612b01565b612542565b604051610f5b9190612b3b565b60405180910390f35b610f7e6004803603810190610f799190612b01565b61257b565b604051610f8b9190612b3b565b60405180910390f35b610fae6004803603810190610fa99190612b01565b6125b3565b604051610fbb9190612b3b565b60405180910390f35b610fde6004803603810190610fd99190612c90565b6125ed565b604051610feb91906132df565b60405180910390f35b61100e60048036038101906110099190612b01565b61268c565b60405161101b9190612b3b565b60405180910390f35b61103e60048036038101906110399190612b01565b6126c5565b60405161104b9190612b3b565b60405180910390f35b61106e60048036038101906110699190612b01565b6126fb565b60405161107b9190612b3b565b60405180910390f35b61109e60048036038101906110999190612b01565b612733565b6040516110ab9190612b3b565b60405180910390f35b6110ce60048036038101906110c99190612b01565b61276b565b6040516110db9190612b3b565b60405180910390f35b6110fe60048036038101906110f99190612b01565b6127c3565b60405161110b9190612b3b565b60405180910390f35b61112e60048036038101906111299190612b01565b612802565b60405161113b9190612b3b565b60405180910390f35b61115e60048036038101906111599190612b01565b61283b565b60405161116b9190612b3b565b60405180910390f35b61118e60048036038101906111899190612b01565b612875565b60405161119b9190612b3b565b60405180910390f35b6111be60048036038101906111b99190612c90565b6128b3565b6040516111cb9190612d51565b60405180910390f35b6111ee60048036038101906111e99190612b01565b61291f565b6040516111fb9190612b3b565b60405180910390f35b61121e60048036038101906112199190612b01565b612958565b60405161122b9190612b3b565b60405180910390f35b61124e60048036038101906112499190612b01565b612993565b60405161125b9190612b3b565b60405180910390f35b61127e60048036038101906112799190612c90565b6129ce565b60405161128b9190613332565b60405180910390f35b6112ae60048036038101906112a99190612b01565b6129fd565b6040516112bb9190612b3b565b60405180910390f35b5f6112cd6117d0565b505f65deadbeef003690505f805b848110156112f1573691506001810190506112db565b505080915050919050565b5f6113056117d0565b505f65deadbeef003290505f805b8481101561132957329150600181019050611313565b505080915050919050565b5f61133d6117d0565b505f65deadbeef005290505f5b8381101561136057815f5260018101905061134a565b5080915050919050565b60605f600890506040828451602086015f855af180611387575f80fd5b5050919050565b5f6113976117d0565b505f65deadbeef000190505f5b838110156113bc575f820191506001810190506113a4565b5080915050919050565b5f8054905090565b5f6113d76117d0565b505f65deadbeef001790505f5b838110156113fc575f821791506001810190506113e4565b5080915050919050565b5f61140f6117d0565b505f65deadbeef003490505f805b848110156114335734915060018101905061141d565b505080915050919050565b5f6114476117d0565b505f65deadbeef000690505f5b8381101561148c577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82069150600181019050611454565b5080915050919050565b5f61149f6117d0565b505f65deadbeef001390505f805b848110156114c6576001831391506001810190506114ad565b505080915050919050565b5f6114da6117d0565b505f65deadbeef002090507fffffffff000000000000000000000000000000000000000000000000000000005f525f805b848110156115245760045f20915060018101905061150b565b507f29045a592007d0c246ef02c2223570da9522d0cf0f73282c79a1bc8f0bb2c2388114611550575f91505b5080915050919050565b5f6115636117d0565b505f65deadbeef00a49050806010525f5b83811015611593576004600360028360066010a4600181019050611574565b5080915050919050565b5f6115a66117d0565b505f65deadbeef001a90505f805b848110156115cc57825f1a91506001810190506115b4565b505080915050919050565b5f6115e06117d0565b505f65deadbeef001b90505f5b8381101561160557815f1b91506001810190506115ed565b5080915050919050565b5f6116186117d0565b505f65deadbeef004290505f805b8481101561163c57429150600181019050611626565b505080915050919050565b5f6116506117d0565b505f65deadbeef003190505f305f5b85811015611676578131925060018101905061165f565b50505080915050919050565b5f61168b6117d0565b505f65deadbeef004890505f805b848110156116af57489150600181019050611699565b505080915050919050565b5f6116c36117d0565b505f65deadbeef003d90505f805b848110156116e7573d91506001810190506116d1565b505080915050919050565b5f6116fb6117d0565b505f65deadbeef004390505f805b8481101561171f57439150600181019050611709565b505080915050919050565b60028181548110611739575f80fd5b905f5260205f20015f91509050805461175190613378565b80601f016020809104026020016040519081016040528092919081815260200182805461177d90613378565b80156117c85780601f1061179f576101008083540402835291602001916117c8565b820191905f5260205f20905b8154815290600101906020018083116117ab57829003601f168201915b505050505081565b5f60015f546117df91906133d5565b5f819055505f54905090565b5f6117f46117d0565b505f65deadbeef000490505f5b8381101561181a57600182049150600181019050611801565b5080915050919050565b5f61182d6117d0565b505f65deadbeef003790505f5b838110156118525760205f803760018101905061183a565b5080915050919050565b5f6118656117d0565b505f65deadbeef00a09050806010525f5b8381101561188e5760066010a0600181019050611876565b5080915050919050565b5f6118a16117d0565b505f65deadbeef003390505f805b848110156118c5573391506001810190506118af565b505080915050919050565b5f6118d96117d0565b505f65deadbeef005390505f5b838110156119005763deadbeef5f526001810190506118e6565b5080915050919050565b5f6119136117d0565b505f65deadbeef003a90505f805b84811015611937573a9150600181019050611921565b505080915050919050565b5f61194b6117d0565b505f65deadbeef005190505f815f525f5b84811015611973575f51915060018101905061195c565b508091505080915050919050565b5f61198a6117d0565b505f65deadbeef001d90505f5b838110156119af57815f1d9150600181019050611997565b5080915050919050565b60605f6005905060208301835160405160208183855f885af1806119db575f80fd5b8195505050505050919050565b5f80600290506020830183518360208183855f885af180611a07575f80fd5b5050505050919050565b5f611a1a6117d0565b505b6103e85a1115611a43576001805f828254611a3791906133d5565b92505081905550611a1c565b600154905090565b5f611a546117d0565b505f65deadbeef001090505f805b84811015611a7b57826001109150600181019050611a62565b505080915050919050565b5f611a8f6117d0565b505f65deadbeef004490505f805b84811015611ab357449150600181019050611a9d565b505080915050919050565b5f611ac76117d0565b505f65deadbeef001190505f805b84811015611aee57600183119150600181019050611ad5565b505080915050919050565b611b01612a57565b5f60099050611b0e612a57565b5f88885f60028110611b2357611b22613408565b5b602002015189600160028110611b3c57611b3b613408565b5b6020020151895f60048110611b5457611b53613408565b5b60200201518a600160048110611b6d57611b6c613408565b5b60200201518b600260048110611b8657611b85613408565b5b60200201518c600360048110611b9f57611b9e613408565b5b60200201518c5f60028110611bb757611bb6613408565b5b60200201518d600160028110611bd057611bcf613408565b5b60200201518d604051602001611bef9a999897969594939291906134ee565b604051602081830303815290604052905060408260d560208401865f19fa611c15575f80fd5b81935050505095945050505050565b5f611c2d6117d0565b505f65deadbeef003e90505f5b83811015611c525760205f803e600181019050611c3a565b5080915050919050565b5f611c656117d0565b505f65deadbeef004590505f805b84811015611c8957459150600181019050611c73565b505080915050919050565b5f611c9d6117d0565b505f65deadbeef000290505f5b83811015611cc357600182029150600181019050611caa565b5080915050919050565b5f611cd66117d0565b505f65deadbeef000890505f5b83811015611d1c577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff5f83089150600181019050611ce3565b5080915050919050565b5f611d2f6117d0565b505f65deadbeef00549050805f555f5b83811015611d56575f549150600181019050611d3f565b5080915050919050565b5f611d696117d0565b505f65deadbeef005a90505f805b84811015611d8d575a9150600181019050611d77565b505080915050919050565b5f611da16117d0565b505f65deadbeef001990505f5b83811015611dc55781199150600181019050611dae565b5065deadbeef00198114611dd857801990505b80915050919050565b606080825114611e26576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611e1d906135fb565b60405180910390fd5b5f6007905060208301835160408482845f875af180611e43575f80fd5b50505050919050565b5f611e556117d0565b505f65deadbeef00a19050806010525f5b83811015611e7f578060066010a1600181019050611e66565b5080915050919050565b5f611e926117d0565b505f65deadbeef001690505f5b83811015611eb7578182169150600181019050611e9f565b5080915050919050565b60605f60049050602083018351604051818183855f885af180611ee2575f80fd5b8195505050505050919050565b6060611ef9612a57565b7f48c9bdf267e6096a3ba7ca8485ae67bb2bf894fe72f36e3cf1361d5f3af54fa5815f60028110611f2d57611f2c613408565b5b6020020181815250507fd182e6ad7f520e511f6c3e2b8c68059b6bbd41fbabd9831f79217e1319cde05b81600160028110611f6b57611f6a613408565b5b602002018181525050611f7c612a79565b7f6162630000000000000000000000000000000000000000000000000000000000815f60048110611fb057611faf613408565b5b6020020181815250505f81600160048110611fce57611fcd613408565b5b6020020181815250505f81600260048110611fec57611feb613408565b5b6020020181815250505f8160036004811061200a57612009613408565b5b60200201818152505061201b612a9b565b7f0300000000000000000000000000000000000000000000000000000000000000815f6002811061204f5761204e613408565b5b602002019077ffffffffffffffffffffffffffffffffffffffffffffffff1916908177ffffffffffffffffffffffffffffffffffffffffffffffff1916815250505f816001600281106120a5576120a4613408565b5b602002019077ffffffffffffffffffffffffffffffffffffffffffffffff1916908177ffffffffffffffffffffffffffffffffffffffffffffffff1916815250505f6120f6600c8585856001611af9565b9050612100612a57565b7fba80a53f981c4d0d6a2797b69f12f6e94c212f14685ac4b74b12bb6fdbffa2d1815f6002811061213457612133613408565b5b6020020181815250507f7d87c5392aab792dc252d5de4533cc9518d38aa8dbf1925ab92386edd40099238160016002811061217257612171613408565b5b602002018181525050805f6002811061218e5761218d613408565b5b6020020151825f600281106121a6576121a5613408565b5b6020020151146121eb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016121e290613689565b60405180910390fd5b806001600281106121ff576121fe613408565b5b60200201518260016002811061221857612217613408565b5b60200201511461225d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161225490613717565b60405180910390fd5b5050505050919050565b5f60808251146122ac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016122a39061377f565b60405180910390fd5b5f60019050602083016020810151601f1a602082015260206040516080835f865af1806122d7575f80fd5b604051519350505050919050565b5f6122ee6117d0565b505b6103e85a1115612326576001805f82825461230b91906133d5565b925050819055504360015461232091906137ca565b506122f0565b600154905090565b5f6123376117d0565b505f65deadbeef004690505f805b8481101561235b57469150600181019050612345565b505080915050919050565b5f61236f6117d0565b505f65deadbeef000590505f5b838110156123955760018205915060018101905061237c565b5080915050919050565b5f6123a86117d0565b505f65deadbeef003990505f5b838110156123cd5760205f80396001810190506123b5565b5080915050919050565b5f6002838390918060018154018082558091505060019003905f5260205f20015f9091929091929091929091925091826124129291906139a1565b50600280549050905092915050565b5f61242a6117d0565b505f65deadbeef005990505f805b8481101561244e57599150600181019050612438565b505080915050919050565b5f6124626117d0565b505f65deadbeef003890505f805b8481101561248657389150600181019050612470565b505080915050919050565b5f61249a6117d0565b505f65deadbeef004190505f805b848110156124be574191506001810190506124a8565b505080915050919050565b5f6124d26117d0565b505f65deadbeef003090505f805b848110156124f6573091506001810190506124e0565b505080915050919050565b5f61250a6117d0565b505f65deadbeef00a39050806010525f5b8381101561253857600360028260066010a360018101905061251b565b5080915050919050565b5f61254b6117d0565b505f65deadbeef000b90505f5b83811015612571578160200b9150600181019050612558565b5080915050919050565b5f6125846117d0565b505f65deadbeef004790505f805b848110156125a857479150600181019050612592565b505080915050919050565b5f6125bc6117d0565b505f65deadbeef001c90505f805b848110156125e257825f1c92506001810190506125ca565b505080915050919050565b5f8061010090505f808273ffffffffffffffffffffffffffffffffffffffff168560405161261b9190613aa8565b5f60405180830381855afa9150503d805f8114612653576040519150601f19603f3d011682016040523d82523d5f602084013e612658565b606091505b50915091508161266b5761266a613abe565b5b6001818060200190518101906126819190613aff565b149350505050919050565b5f6126956117d0565b505f65deadbeef003590505f805b848110156126ba575f3591506001810190506126a3565b505080915050919050565b5f6126ce6117d0565b505f65deadbeef005590505f5b838110156126f157815f556001810190506126db565b5080915050919050565b5f6127046117d0565b505f65deadbeef001890505f5b83811015612729575f82189150600181019050612711565b5080915050919050565b5f61273c6117d0565b505f65deadbeef000390505f5b83811015612761575f82039150600181019050612749565b5080915050919050565b5f6127746117d0565b505f65deadbeef000790505f5b838110156127b9577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82079150600181019050612781565b5080915050919050565b5f6127cc6117d0565b505f65deadbeef00a29050806010525f5b838110156127f85760028160066010a26001810190506127dd565b5080915050919050565b5f61280b6117d0565b505f65deadbeef000a90505f5b83811015612831576001820a9150600181019050612818565b5080915050919050565b5f6128446117d0565b505f65deadbeef001490505f805b8481101561286a578283149150600181019050612852565b505080915050919050565b5f61287e6117d0565b505f65deadbeef004090505f600143035f5b858110156128a75781409250600181019050612890565b50505080915050919050565b606060808251146128f9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016128f0906135fb565b60405180910390fd5b5f6006905060208301835160408482845f875af180612916575f80fd5b50505050919050565b5f6129286117d0565b505f65deadbeef001590505f805b8481101561294d5782159150600181019050612936565b505080915050919050565b5f6129616117d0565b505f65deadbeef001290505f805b848110156129885782600112915060018101905061296f565b505080915050919050565b5f61299c6117d0565b505f65deadbeef003b90505f305f5b858110156129c257813b92506001810190506129ab565b50505080915050919050565b5f806003905060208301835160405160148183855f885af1806129ef575f80fd5b815195505050505050919050565b5f612a066117d0565b505f65deadbeef000990505f5b83811015612a4d577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600183099150600181019050612a13565b5080915050919050565b6040518060400160405280600290602082028036833780820191505090505090565b6040518060800160405280600490602082028036833780820191505090505090565b6040518060400160405280600290602082028036833780820191505090505090565b5f604051905090565b5f80fd5b5f80fd5b5f819050919050565b612ae081612ace565b8114612aea575f80fd5b50565b5f81359050612afb81612ad7565b92915050565b5f60208284031215612b1657612b15612ac6565b5b5f612b2384828501612aed565b91505092915050565b612b3581612ace565b82525050565b5f602082019050612b4e5f830184612b2c565b92915050565b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b612ba282612b5c565b810181811067ffffffffffffffff82111715612bc157612bc0612b6c565b5b80604052505050565b5f612bd3612abd565b9050612bdf8282612b99565b919050565b5f67ffffffffffffffff821115612bfe57612bfd612b6c565b5b612c0782612b5c565b9050602081019050919050565b828183375f83830152505050565b5f612c34612c2f84612be4565b612bca565b905082815260208101848484011115612c5057612c4f612b58565b5b612c5b848285612c14565b509392505050565b5f82601f830112612c7757612c76612b54565b5b8135612c87848260208601612c22565b91505092915050565b5f60208284031215612ca557612ca4612ac6565b5b5f82013567ffffffffffffffff811115612cc257612cc1612aca565b5b612cce84828501612c63565b91505092915050565b5f81519050919050565b5f82825260208201905092915050565b5f5b83811015612d0e578082015181840152602081019050612cf3565b5f8484015250505050565b5f612d2382612cd7565b612d2d8185612ce1565b9350612d3d818560208601612cf1565b612d4681612b5c565b840191505092915050565b5f6020820190508181035f830152612d698184612d19565b905092915050565b5f819050919050565b612d8381612d71565b82525050565b5f602082019050612d9c5f830184612d7a565b92915050565b5f63ffffffff82169050919050565b612dba81612da2565b8114612dc4575f80fd5b50565b5f81359050612dd581612db1565b92915050565b5f67ffffffffffffffff821115612df557612df4612b6c565b5b602082029050919050565b5f80fd5b612e0d81612d71565b8114612e17575f80fd5b50565b5f81359050612e2881612e04565b92915050565b5f612e40612e3b84612ddb565b612bca565b90508060208402830185811115612e5a57612e59612e00565b5b835b81811015612e835780612e6f8882612e1a565b845260208401935050602081019050612e5c565b5050509392505050565b5f82601f830112612ea157612ea0612b54565b5b6002612eae848285612e2e565b91505092915050565b5f67ffffffffffffffff821115612ed157612ed0612b6c565b5b602082029050919050565b5f612eee612ee984612eb7565b612bca565b90508060208402830185811115612f0857612f07612e00565b5b835b81811015612f315780612f1d8882612e1a565b845260208401935050602081019050612f0a565b5050509392505050565b5f82601f830112612f4f57612f4e612b54565b5b6004612f5c848285612edc565b91505092915050565b5f67ffffffffffffffff821115612f7f57612f7e612b6c565b5b602082029050919050565b5f7fffffffffffffffff00000000000000000000000000000000000000000000000082169050919050565b612fbe81612f8a565b8114612fc8575f80fd5b50565b5f81359050612fd981612fb5565b92915050565b5f612ff1612fec84612f65565b612bca565b9050806020840283018581111561300b5761300a612e00565b5b835b8181101561303457806130208882612fcb565b84526020840193505060208101905061300d565b5050509392505050565b5f82601f83011261305257613051612b54565b5b600261305f848285612fdf565b91505092915050565b5f8115159050919050565b61307c81613068565b8114613086575f80fd5b50565b5f8135905061309781613073565b92915050565b5f805f805f61014086880312156130b7576130b6612ac6565b5b5f6130c488828901612dc7565b95505060206130d588828901612e8d565b94505060606130e688828901612f3b565b93505060e06130f78882890161303e565b92505061012061310988828901613089565b9150509295509295909350565b5f60029050919050565b5f81905092915050565b5f819050919050565b61313c81612d71565b82525050565b5f61314d8383613133565b60208301905092915050565b5f602082019050919050565b61316e81613116565b6131788184613120565b92506131838261312a565b805f5b838110156131b357815161319a8782613142565b96506131a583613159565b925050600181019050613186565b505050505050565b5f6040820190506131ce5f830184613165565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6131fd826131d4565b9050919050565b61320d816131f3565b82525050565b5f6020820190506132265f830184613204565b92915050565b5f80fd5b5f8083601f84011261324557613244612b54565b5b8235905067ffffffffffffffff8111156132625761326161322c565b5b60208301915083600182028301111561327e5761327d612e00565b5b9250929050565b5f806020838503121561329b5761329a612ac6565b5b5f83013567ffffffffffffffff8111156132b8576132b7612aca565b5b6132c485828601613230565b92509250509250929050565b6132d981613068565b82525050565b5f6020820190506132f25f8301846132d0565b92915050565b5f7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000082169050919050565b61332c816132f8565b82525050565b5f6020820190506133455f830184613323565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061338f57607f821691505b6020821081036133a2576133a161334b565b5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6133df82612ace565b91506133ea83612ace565b9250828201905080821115613402576134016133a8565b5b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f8160e01b9050919050565b5f61344b82613435565b9050919050565b61346361345e82612da2565b613441565b82525050565b5f819050919050565b61348361347e82612d71565b613469565b82525050565b5f819050919050565b6134a361349e82612f8a565b613489565b82525050565b5f8160f81b9050919050565b5f6134bf826134a9565b9050919050565b5f6134d0826134b5565b9050919050565b6134e86134e382613068565b6134c6565b82525050565b5f6134f9828d613452565b600482019150613509828c613472565b602082019150613519828b613472565b602082019150613529828a613472565b6020820191506135398289613472565b6020820191506135498288613472565b6020820191506135598287613472565b6020820191506135698286613492565b6008820191506135798285613492565b60088201915061358982846134d7565b6001820191508190509b9a5050505050505050505050565b5f82825260208201905092915050565b7f496e76616c696420696e707574206c656e6774680000000000000000000000005f82015250565b5f6135e56014836135a1565b91506135f0826135b1565b602082019050919050565b5f6020820190508181035f830152613612816135d9565b9050919050565b7f54657374426c616b653266202d204669727374206861736820646f65736e27745f8201527f206d617463680000000000000000000000000000000000000000000000000000602082015250565b5f6136736026836135a1565b915061367e82613619565b604082019050919050565b5f6020820190508181035f8301526136a081613667565b9050919050565b7f54657374426c616b653266202d205365636f6e64206861736820646f65736e275f8201527f74206d6174636800000000000000000000000000000000000000000000000000602082015250565b5f6137016027836135a1565b915061370c826136a7565b604082019050919050565b5f6020820190508181035f83015261372e816136f5565b9050919050565b7f496e76616c696420696e7075742064617461206c656e6774682e0000000000005f82015250565b5f613769601a836135a1565b915061377482613735565b602082019050919050565b5f6020820190508181035f8301526137968161375d565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f6137d482612ace565b91506137df83612ace565b9250826137ef576137ee61379d565b5b828206905092915050565b5f82905092915050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026138607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82613825565b61386a8683613825565b95508019841693508086168417925050509392505050565b5f819050919050565b5f6138a56138a061389b84612ace565b613882565b612ace565b9050919050565b5f819050919050565b6138be8361388b565b6138d26138ca826138ac565b848454613831565b825550505050565b5f90565b6138e66138da565b6138f18184846138b5565b505050565b5b81811015613914576139095f826138de565b6001810190506138f7565b5050565b601f8211156139595761392a81613804565b61393384613816565b81016020851015613942578190505b61395661394e85613816565b8301826138f6565b50505b505050565b5f82821c905092915050565b5f6139795f198460080261395e565b1980831691505092915050565b5f613991838361396a565b9150826002028217905092915050565b6139ab83836137fa565b67ffffffffffffffff8111156139c4576139c3612b6c565b5b6139ce8254613378565b6139d9828285613918565b5f601f831160018114613a06575f84156139f4578287013590505b6139fe8582613986565b865550613a65565b601f198416613a1486613804565b5f5b82811015613a3b57848901358255600182019150602085019450602081019050613a16565b86831015613a585784890135613a54601f89168261396a565b8355505b6001600288020188555050505b50505050505050565b5f81905092915050565b5f613a8282612cd7565b613a8c8185613a6e565b9350613a9c818560208601612cf1565b80840191505092915050565b5f613ab38284613a78565b915081905092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52600160045260245ffd5b5f81519050613af981612ad7565b92915050565b5f60208284031215613b1457613b13612ac6565b5b5f613b2184828501613aeb565b9150509291505056fea2646970667358221220d2046a05353f8a265314b3a3e3edba20ebcfaf40b1fc0bfd127770a9167946b464736f6c63430008170033",
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
