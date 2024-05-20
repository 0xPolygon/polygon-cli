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
	ABI: "[{\"type\":\"function\",\"name\":\"dumpster\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCallCounter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"inc\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"loopBlockHashUntilLimit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"loopUntilLimit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"store\",\"inputs\":[{\"name\":\"trash\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testADD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testADDMOD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testADDRESS\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testAND\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testBALANCE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testBASEFEE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testBLOCKHASH\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testBYTE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testBlake2f\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCALLDATACOPY\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCALLDATALOAD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCALLDATASIZE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCALLER\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCALLVALUE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCHAINID\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCODECOPY\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCODESIZE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testCOINBASE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testDIFFICULTY\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testDIV\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testECAdd\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testECMul\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testECPairing\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testECRecover\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testEQ\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testEXP\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testEXTCODESIZE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testGAS\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testGASLIMIT\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testGASPRICE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testGT\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testISZERO\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testIdentity\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testLOG0\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testLOG1\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testLOG2\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testLOG3\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testLOG4\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testLT\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testMLOAD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testMOD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testMSIZE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testMSTORE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testMSTORE8\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testMUL\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testMULMOD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testModExp\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testNOT\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testNUMBER\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testOR\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testORIGIN\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testRETURNDATACOPY\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testRETURNDATASIZE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testRipemd160\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSAR\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSDIV\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSELFBALANCE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSGT\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSHA256\",\"inputs\":[{\"name\":\"inputData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"result\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSHA3\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSHL\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSHR\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSIGNEXTEND\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSLOAD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSLT\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSMOD\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSSTORE\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testSUB\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testTIMESTAMP\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"testXOR\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611d1e806100206000396000f3fe608060405234801561001057600080fd5b50600436106104545760003560e01c806380947f8011610241578063bf529ca11161013b578063dd9bef60116100c3578063f279ca8111610087578063f279ca811461098f578063f4d1fc61146109a2578063f58fc36a146109b5578063f6b0bbf7146109c8578063fde7721c146109e857600080fd5b8063dd9bef6014610930578063de97a36314610943578063e9f9b3f214610956578063ea5141e614610969578063edf003cf1461097c57600080fd5b8063ce3cf4ef1161010a578063ce3cf4ef146108d1578063d117320b146108e4578063d51e7b5b146108f7578063d53ff3fd1461090a578063d93cd5581461091d57600080fd5b8063bf529ca114610885578063c360aba614610898578063c420eb61146108ab578063c4bd65d5146108be57600080fd5b8063a18683cb116101c9578063b374012b1161018d578063b374012b14610826578063b3d847f214610839578063b7b862071461084c578063b81c14841461085f578063bdc875fc1461087257600080fd5b8063a18683cb146107c5578063a271b721146107e5578063a60a1087146107ed578063a645c9c214610800578063acaebdf61461081357600080fd5b8063962e4dc211610210578063962e4dc21461077957806398456f3e1461078c5780639a2b7c811461079f5780639cce7cf9146107b2578063a040aec6146104a857600080fd5b806380947f801461072d578063880eff3914610740578063918a5fcd1461075357806391e7b2771461076657600080fd5b80633430ec061161035257806360e13cde116102da5780636f099c8d1161029e5780636f099c8d146106ce57806371d91d28146106e15780637b6e0b0e146106f45780637c191d20146107075780637de8c6f81461071a57600080fd5b806360e13cde1461067a578063613d0a821461068d57806363138d4f146106a0578063659bbb4f146106b35780636e7f1fe7146106bb57600080fd5b806340fe26621161032157806340fe26621461061b57806344cf3bc71461062e5780634a61af1f146106415780634d2c74b3146106545780635590c2d91461066757600080fd5b80633430ec06146105da578063371303c0146105ed5780633a411f12146105f55780633a425dfc1461060857600080fd5b806318093b46116103e0578063219cddeb116103a4578063219cddeb1461057b5780632294fc7f1461058e5780632871ef85146105a15780632b21ef44146105b45780632d34e798146105c757600080fd5b806318093b461461051c57806319b621d61461052f5780631aba07ea146105425780631de2f343146105555780632007332e1461056857600080fd5b80630ba8a73b116104275780630ba8a73b146104c85780631287a68c146104db578063135d52f7146104e35780631581cf19146104f6578063165821501461050957600080fd5b8063034aef7114610459578063050082f814610482578063087b4e84146104955780630b3b996a146104a8575b600080fd5b61046c610467366004611786565b6109fb565b60405161047991906117b7565b60405180910390f35b61046c610490366004611786565b610a2d565b61046c6104a3366004611786565b610a56565b6104bb6104b63660046118bc565b610a87565b6040516104799190611953565b61046c6104d6366004611786565b610aaa565b60005461046c565b61046c6104f1366004611786565b610acf565b61046c610504366004611786565b610af1565b61046c610517366004611786565b610b1a565b61046c61052a366004611786565b610b46565b61046c61053d366004611786565b610b71565b61046c610550366004611786565b610bdd565b61046c610563366004611786565b610c13565b61046c610576366004611786565b610c40565b61046c610589366004611786565b610c62565b61046c61059c366004611786565b610c8b565b61046c6105af366004611786565b610cc0565b61046c6105c2366004611786565b610ce9565b61046c6105d5366004611786565b610d12565b6104bb6105e8366004611786565b610d3b565b61046c610de4565b61046c610603366004611786565b610dfd565b61046c610616366004611786565b610e1f565b61046c610629366004611786565b610e4a565b61046c61063c366004611786565b610e79565b61046c61064f366004611786565b610ea2565b61046c610662366004611786565b610ecf565b61046c610675366004611786565b610ef8565b61046c610688366004611786565b610f2e565b6104bb61069b3660046118bc565b610f5a565b61046c6106ae3660046118bc565b610f85565b61046c610fae565b61046c6106c9366004611786565b610fe8565b61046c6106dc366004611786565b611013565b61046c6106ef366004611786565b61103c565b61046c610702366004611786565b611067565b61046c610715366004611786565b611092565b61046c610728366004611786565b6110bb565b61046c61073b366004611786565b6110dd565b61046c61074e366004611786565b61110b565b61046c610761366004611786565b611138565b61046c610774366004611786565b611161565b6104bb6107873660046118bc565b61119f565b61046c61079a366004611786565b6111f0565b61046c6107ad366004611786565b611220565b6104bb6107c03660046118bc565b611242565b6107d86107d33660046118bc565b611262565b6040516104799190611985565b61046c6112bc565b61046c6107fb366004611786565b6112fd565b61046c61080e366004611786565b611326565b61046c610821366004611786565b611348565b61046c6108343660046119e5565b611373565b61046c610847366004611786565b6113a2565b61046c61085a366004611786565b6113cb565b61046c61086d366004611786565b6113f4565b61046c610880366004611786565b61141d565b61046c610893366004611786565b611446565b61046c6108a6366004611786565b61147a565b61046c6108b9366004611786565b61149c565b61046c6108cc366004611786565b6114c5565b61046c6108df366004611786565b6114eb565b61046c6108f2366004611786565b611516565b61046c610905366004611786565b611540565b61046c610918366004611786565b611562565b61046c61092b366004611786565b611584565b61046c61093e366004611786565b6115b0565b61046c610951366004611786565b6115e2565b61046c610964366004611786565b61160c565b61046c610977366004611786565b611635565b6104bb61098a3660046118bc565b611664565b61046c61099d366004611786565b6116a3565b61046c6109b0366004611786565b6116cd565b61046c6109c3366004611786565b6116f8565b6109db6109d63660046118bc565b611723565b6040516104799190611a42565b61046c6109f6366004611786565b611751565b6000610a05610de4565b5065deadbeef00366000805b84811015610a2457369150600101610a11565b50909392505050565b6000610a37610de4565b5065deadbeef00326000805b84811015610a2457329150600101610a43565b6000610a60610de4565b5065deadbeef005260005b83811015610a80576000829052600101610a6b565b5092915050565b606060086040828451602086016000855af180610aa357600080fd5b5050919050565b6000610ab4610de4565b5065deadbeef000160005b83811015610a8057600101610abf565b6000610ad9610de4565b5065deadbeef001760008315610a8057600101610abf565b6000610afb610de4565b5065deadbeef00346000805b84811015610a2457349150600101610b07565b6000610b24610de4565b5065deadbeef000660005b83811015610a805760001990910690600101610b2f565b6000610b50610de4565b5065deadbeef00136000805b84811015610a24576001808413925001610b5c565b6000610b7b610de4565b506001600160e01b0319600090815265deadbeef002090805b84811015610bab5760046000209150600101610b94565b507f29045a592007d0c246ef02c2223570da9522d0cf0f73282c79a1bc8f0bb2c2388114610a80575060009392505050565b6000610be7610de4565b5065deadbeef00a4601081905260005b83811015610a80576004600360028360066010a4600101610bf7565b6000610c1d610de4565b5065deadbeef001a6000805b84811015610a2457600083901a9150600101610c29565b6000610c4a610de4565b5065deadbeef001b60008315610a8057600101610abf565b6000610c6c610de4565b5065deadbeef00426000805b84811015610a2457429150600101610c78565b6000610c95610de4565b5065deadbeef0031600030815b85811015610cb65781319250600101610ca2565b5091949350505050565b6000610cca610de4565b5065deadbeef00486000805b84811015610a2457489150600101610cd6565b6000610cf3610de4565b5065deadbeef003d6000805b84811015610a24573d9150600101610cff565b6000610d1c610de4565b5065deadbeef00436000805b84811015610a2457439150600101610d28565b60028181548110610d4b57600080fd5b906000526020600020018054909150610d6390611a66565b80601f0160208091040260200160405190810160405280929190818152602001828054610d8f90611a66565b8015610ddc5780601f10610db157610100808354040283529160200191610ddc565b820191906000526020600020905b815481529060010190602001808311610dbf57829003601f168201915b505050505081565b60008054610df3906001611aa8565b6000819055919050565b6000610e07610de4565b5065deadbeef000460008315610a8057600101610abf565b6000610e29610de4565b5065deadbeef003760005b83811015610a8057602060008037600101610e34565b6000610e54610de4565b5065deadbeef00a0601081905260005b83811015610a805760066010a0600101610e64565b6000610e83610de4565b5065deadbeef00336000805b84811015610a2457339150600101610e8f565b6000610eac610de4565b5065deadbeef005360005b83811015610a805763deadbeef600052600101610eb7565b6000610ed9610de4565b5065deadbeef003a6000805b84811015610a24573a9150600101610ee5565b6000610f02610de4565b5065deadbeef00516000818152805b84811015610f26576000519150600101610f11565b509392505050565b6000610f38610de4565b5065deadbeef001d60005b83811015610a805760009190911d90600101610f43565b6060600560208301835160405160208183856000885af180610f7b57600080fd5b5095945050505050565b600060026020830183518360208183856000885af180610fa457600080fd5b5050505050919050565b6000610fb8610de4565b505b6103e85a1115610fe1576001806000828254610fd69190611aa8565b90915550610fba9050565b5060015490565b6000610ff2610de4565b5065deadbeef00106000805b84811015610a24576001838110925001610ffe565b600061101d610de4565b5065deadbeef00446000805b84811015610a2457449150600101611029565b6000611046610de4565b5065deadbeef00116000805b84811015610a24576001808411925001611052565b6000611071610de4565b5065deadbeef003e60005b83811015610a805760206000803e60010161107c565b600061109c610de4565b5065deadbeef00456000805b84811015610a24574591506001016110a8565b60006110c5610de4565b5065deadbeef000260008315610a8057600101610abf565b60006110e7610de4565b5065deadbeef000860005b83811015610a80576000196000830891506001016110f2565b6000611115610de4565b5065deadbeef005460008181555b83811015610a80576000549150600101611123565b6000611142610de4565b5065deadbeef005a6000805b84811015610a24575a915060010161114e565b600061116b610de4565b5065deadbeef001960005b8381101561118957901990600101611176565b5065deadbeef0019811461119957195b92915050565b606081516060146111cb5760405162461bcd60e51b81526004016111c290611ae9565b60405180910390fd5b600760208301835160408482846000875af1806111e757600080fd5b50505050919050565b60006111fa610de4565b5065deadbeef00a1601081905260005b83811015610a80578060066010a160010161120a565b600061122a610de4565b5065deadbeef001660008315610a8057600101610abf565b60606004602083018351604051818183856000885af180610f7b57600080fd5b600081516080146112855760405162461bcd60e51b81526004016111c290611b2d565b6001602083016040840151601f1a602082015260206040516080836000865af1806112af57600080fd5b6040515195945050505050565b60006112c6610de4565b505b6103e85a1115610fe15760018060008282546112e49190611aa8565b90915550506001546112f7904390611b53565b506112c8565b6000611307610de4565b5065deadbeef00466000805b84811015610a2457469150600101611313565b6000611330610de4565b5065deadbeef000560008315610a8057600101610abf565b6000611352610de4565b5065deadbeef003960005b83811015610a805760206000803960010161135d565b60028054600181018255600091825283908390602084200191611397919083611c17565b505060025492915050565b60006113ac610de4565b5065deadbeef00596000805b84811015610a24575991506001016113b8565b60006113d5610de4565b5065deadbeef00386000805b84811015610a24573891506001016113e1565b60006113fe610de4565b5065deadbeef00416000805b84811015610a245741915060010161140a565b6000611427610de4565b5065deadbeef00306000805b84811015610a2457309150600101611433565b6000611450610de4565b5065deadbeef00a3601081905260005b83811015610a8057600360028260066010a3600101611460565b6000611484610de4565b5065deadbeef000b60008315610a8057600101610abf565b60006114a6610de4565b5065deadbeef00476000805b84811015610a24574791506001016114b2565b60006114cf610de4565b5065deadbeef001c6000805b84811015610a24576001016114db565b60006114f5610de4565b5065deadbeef00356000805b84811015610a24576000359150600101611501565b6000611520610de4565b5065deadbeef005560005b83811015610a8057600082905560010161152b565b600061154a610de4565b5065deadbeef001860008315610a8057600101610abf565b600061156c610de4565b5065deadbeef000360008315610a8057600101610abf565b600061158e610de4565b5065deadbeef000760005b83811015610a805760001990910790600101611599565b60006115ba610de4565b5065deadbeef00a2601081905260005b83811015610a805760028160066010a26001016115ca565b60006115ec610de4565b5065deadbeef000a60005b83811015610a805760019182900a91016115f7565b6000611616610de4565b5065deadbeef00146000805b84811015610a2457600191508101611622565b600061163f610de4565b5065deadbeef004060006000194301815b85811015610cb65781409250600101611650565b606081516080146116875760405162461bcd60e51b81526004016111c290611ae9565b600660208301835160408482846000875af1806111e757600080fd5b60006116ad610de4565b5065deadbeef00156000805b84811015610a2457821591506001016116b9565b60006116d7610de4565b5065deadbeef00126000805b84811015610a245760018381129250016116e3565b6000611702610de4565b5065deadbeef003b600030815b85811015610cb657813b925060010161170f565b6000600360208301835160405160148183856000885af18061174457600080fd5b8151979650505050505050565b600061175b610de4565b5065deadbeef000960005b83811015610a8057600019600183099150600101611766565b8035611199565b60006020828403121561179b5761179b600080fd5b60006117a7848461177f565b949350505050565b805b82525050565b6020810161119982846117af565b634e487b7160e01b600052604160045260246000fd5b601f19601f830116810181811067ffffffffffffffff82111715611801576118016117c5565b6040525050565b60006118176000604051905090565b905061182382826117db565b919050565b600067ffffffffffffffff821115611842576118426117c5565b601f19601f83011660200192915050565b82818337506000910152565b600061187261186d84611828565b611808565b90508281526020810184848401111561188d5761188d600080fd5b610f26848285611853565b600082601f8301126118ac576118ac600080fd5b81356117a784826020860161185f565b6000602082840312156118d1576118d1600080fd5b813567ffffffffffffffff8111156118eb576118eb600080fd5b6117a784828501611898565b60005b838110156119125780820151838201526020016118fa565b50506000910152565b600061192b826000815192915050565b8084526020840193506119428185602086016118f7565b601f01601f19169290920192915050565b60208082528101611964818461191b565b9392505050565b60006001600160a01b038216611199565b6117b18161196b565b60208101611199828461197c565b60008083601f8401126119a8576119a8600080fd5b50813567ffffffffffffffff8111156119c3576119c3600080fd5b6020830191508360018202830111156119de576119de600080fd5b9250929050565b600080602083850312156119fb576119fb600080fd5b823567ffffffffffffffff811115611a1557611a15600080fd5b611a2185828601611993565b92509250509250929050565b6bffffffffffffffffffffffff1981166117b1565b602081016111998284611a2d565b634e487b7160e01b600052602260045260246000fd5b600281046001821680611a7a57607f821691505b602082108103611a8c57611a8c611a50565b50919050565b634e487b7160e01b600052601160045260246000fd5b8082018082111561119957611199611a92565b6014815260006020820173092dcecc2d8d2c840d2dce0eae840d8cadccee8d60631b815291505b5060200190565b6020808252810161119981611abb565b601a81526000602082017f496e76616c696420696e7075742064617461206c656e6774682e00000000000081529150611ae2565b6020808252810161119981611af9565b634e487b7160e01b600052601260045260246000fd5b600082611b6257611b62611b3d565b500690565b6000611199611b738381565b90565b611b7f83611b67565b815460001960089490940293841b1916921b91909117905550565b6000611ba7818484611b76565b505050565b81811015611bc757611bbf600082611b9a565b600101611bac565b5050565b601f821115611ba757611be981600081815281906020902092915050565b6020601f85010481016020851015611bfe5750805b611c106020601f860104830182611bac565b5050505050565b8267ffffffffffffffff811115611c3057611c306117c5565b611c3a8254611a66565b611c45828285611bcb565b6000601f831160018114611c795760008415611c615750858201355b600019600886021c1981166002860217865550611cdf565b601f198416611c9386600081815281906020902092915050565b60005b82811015611cb65788850135825560209485019460019092019101611c96565b86831015611cd257600019601f88166008021c19858a01351682555b6001600288020188555050505b5050505050505056fea2646970667358221220b3a835504b6ee4829d77d46cf13a8e6b8f6628dead0c6e15b8148e4555ee87e864736f6c63430008170033",
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
