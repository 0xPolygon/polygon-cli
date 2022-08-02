// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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
)

// LoadTesterMetaData contains all meta data concerning the LoadTester contract.
var LoadTesterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getCallCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"inc\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testADD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testADDMOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testADDRESS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testAND\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testBALANCE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testBASEFEE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testBLOCKHASH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testBYTE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCALLDATACOPY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCALLDATALOAD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCALLDATASIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCALLER\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCALLVALUE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCHAINID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCODECOPY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCODESIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCOINBASE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testDIFFICULTY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testDIV\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testEQ\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testEXP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testEXTCODESIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testGAS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testGASLIMIT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testGASPRICE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testGT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testISZERO\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testLOG0\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testLOG1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testLOG2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testLOG3\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testLOG4\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testLT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testMLOAD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testMOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testMSIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testMSTORE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testMSTORE8\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testMUL\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testMULMOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testNOT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testNUMBER\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testORIGIN\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testRETURNDATACOPY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testRETURNDATASIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSAR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSDIV\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSELFBALANCE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSGT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSHA3\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSHL\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSHR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSIGNEXTEND\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSLOAD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSLT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSMOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSSTORE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSUB\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testTIMESTAMP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testXOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061204a806100206000396000f3fe608060405234801561001057600080fd5b50600436106103c55760003560e01c80637c191d20116101ff578063c360aba61161011a578063dd9bef60116100ad578063f279ca811161007c578063f279ca8114610eb6578063f4d1fc6114610ee6578063f58fc36a14610f16578063fde7721c14610f46576103c5565b8063dd9bef6014610df6578063de97a36314610e26578063e9f9b3f214610e56578063ea5141e614610e86576103c5565b8063d117320b116100e9578063d117320b14610d36578063d51e7b5b14610d66578063d53ff3fd14610d96578063d93cd55814610dc6576103c5565b8063c360aba614610c76578063c420eb6114610ca6578063c4bd65d514610cd6578063ce3cf4ef14610d06576103c5565b8063a60a108711610192578063b7b8620711610161578063b7b8620714610bb6578063b81c148414610be6578063bdc875fc14610c16578063bf529ca114610c46576103c5565b8063a60a108714610af6578063a645c9c214610b26578063acaebdf614610b56578063b3d847f214610b86576103c5565b8063918a5fcd116101ce578063918a5fcd14610a3657806391e7b27714610a6657806398456f3e14610a965780639a2b7c8114610ac6576103c5565b80637c191d20146109765780637de8c6f8146109a657806380947f80146109d6578063880eff3914610a06576103c5565b80632871ef85116102ef5780634a61af1f116102825780636e7f1fe7116102515780636e7f1fe7146108b65780636f099c8d146108e657806371d91d28146109165780637b6e0b0e14610946576103c5565b80634a61af1f146107f65780634d2c74b3146108265780635590c2d91461085657806360e13cde14610886576103c5565b80633a411f12116102be5780633a411f12146107365780633a425dfc1461076657806340fe26621461079657806344cf3bc7146107c6576103c5565b80632871ef85146106885780632b21ef44146106b85780632d34e798146106e8578063371303c014610718576103c5565b806316582150116103675780631de2f343116103365780631de2f343146105c85780632007332e146105f8578063219cddeb146106285780632294fc7f14610658576103c5565b8063165821501461050857806318093b461461053857806319b621d6146105685780631aba07ea14610598576103c5565b80630ba8a73b116103a35780630ba8a73b1461045a5780631287a68c1461048a578063135d52f7146104a85780631581cf19146104d8576103c5565b8063034aef71146103ca578063050082f8146103fa578063087b4e841461042a575b600080fd5b6103e460048036038101906103df9190611f38565b610f76565b6040516103f19190611f74565b60405180910390f35b610414600480360381019061040f9190611f38565b610fb1565b6040516104219190611f74565b60405180910390f35b610444600480360381019061043f9190611f38565b610fec565b6040516104519190611f74565b60405180910390f35b610474600480360381019061046f9190611f38565b611026565b6040516104819190611f74565b60405180910390f35b610492611062565b60405161049f9190611f74565b60405180910390f35b6104c260048036038101906104bd9190611f38565b61106b565b6040516104cf9190611f74565b60405180910390f35b6104f260048036038101906104ed9190611f38565b6110a7565b6040516104ff9190611f74565b60405180910390f35b610522600480360381019061051d9190611f38565b6110e2565b60405161052f9190611f74565b60405180910390f35b610552600480360381019061054d9190611f38565b61113d565b60405161055f9190611f74565b60405180910390f35b610582600480360381019061057d9190611f38565b61117b565b60405161058f9190611f74565b60405180910390f35b6105b260048036038101906105ad9190611f38565b61120a565b6040516105bf9190611f74565b60405180910390f35b6105e260048036038101906105dd9190611f38565b611250565b6040516105ef9190611f74565b60405180910390f35b610612600480360381019061060d9190611f38565b61128e565b60405161061f9190611f74565b60405180910390f35b610642600480360381019061063d9190611f38565b6112ca565b60405161064f9190611f74565b60405180910390f35b610672600480360381019061066d9190611f38565b611305565b60405161067f9190611f74565b60405180910390f35b6106a2600480360381019061069d9190611f38565b611344565b6040516106af9190611f74565b60405180910390f35b6106d260048036038101906106cd9190611f38565b61137f565b6040516106df9190611f74565b60405180910390f35b61070260048036038101906106fd9190611f38565b6113ba565b60405161070f9190611f74565b60405180910390f35b6107206113f5565b60405161072d9190611f74565b60405180910390f35b610750600480360381019061074b9190611f38565b611414565b60405161075d9190611f74565b60405180910390f35b610780600480360381019061077b9190611f38565b611450565b60405161078d9190611f74565b60405180910390f35b6107b060048036038101906107ab9190611f38565b61148c565b6040516107bd9190611f74565b60405180910390f35b6107e060048036038101906107db9190611f38565b6114cb565b6040516107ed9190611f74565b60405180910390f35b610810600480360381019061080b9190611f38565b611506565b60405161081d9190611f74565b60405180910390f35b610840600480360381019061083b9190611f38565b611544565b60405161084d9190611f74565b60405180910390f35b610870600480360381019061086b9190611f38565b61157f565b60405161087d9190611f74565b60405180910390f35b6108a0600480360381019061089b9190611f38565b6115c4565b6040516108ad9190611f74565b60405180910390f35b6108d060048036038101906108cb9190611f38565b611600565b6040516108dd9190611f74565b60405180910390f35b61090060048036038101906108fb9190611f38565b61163e565b60405161090d9190611f74565b60405180910390f35b610930600480360381019061092b9190611f38565b611679565b60405161093d9190611f74565b60405180910390f35b610960600480360381019061095b9190611f38565b6116b7565b60405161096d9190611f74565b60405180910390f35b610990600480360381019061098b9190611f38565b6116f3565b60405161099d9190611f74565b60405180910390f35b6109c060048036038101906109bb9190611f38565b61172e565b6040516109cd9190611f74565b60405180910390f35b6109f060048036038101906109eb9190611f38565b61176a565b6040516109fd9190611f74565b60405180910390f35b610a206004803603810190610a1b9190611f38565b6117c7565b604051610a2d9190611f74565b60405180910390f35b610a506004803603810190610a4b9190611f38565b611806565b604051610a5d9190611f74565b60405180910390f35b610a806004803603810190610a7b9190611f38565b611841565b604051610a8d9190611f74565b60405180910390f35b610ab06004803603810190610aab9190611f38565b61188d565b604051610abd9190611f74565b60405180910390f35b610ae06004803603810190610adb9190611f38565b6118cd565b604051610aed9190611f74565b60405180910390f35b610b106004803603810190610b0b9190611f38565b611908565b604051610b1d9190611f74565b60405180910390f35b610b406004803603810190610b3b9190611f38565b611943565b604051610b4d9190611f74565b60405180910390f35b610b706004803603810190610b6b9190611f38565b61197f565b604051610b7d9190611f74565b60405180910390f35b610ba06004803603810190610b9b9190611f38565b6119bb565b604051610bad9190611f74565b60405180910390f35b610bd06004803603810190610bcb9190611f38565b6119f6565b604051610bdd9190611f74565b60405180910390f35b610c006004803603810190610bfb9190611f38565b611a31565b604051610c0d9190611f74565b60405180910390f35b610c306004803603810190610c2b9190611f38565b611a6c565b604051610c3d9190611f74565b60405180910390f35b610c606004803603810190610c5b9190611f38565b611aa7565b604051610c6d9190611f74565b60405180910390f35b610c906004803603810190610c8b9190611f38565b611aeb565b604051610c9d9190611f74565b60405180910390f35b610cc06004803603810190610cbb9190611f38565b611b27565b604051610ccd9190611f74565b60405180910390f35b610cf06004803603810190610ceb9190611f38565b611b62565b604051610cfd9190611f74565b60405180910390f35b610d206004803603810190610d1b9190611f38565b611ba0565b604051610d2d9190611f74565b60405180910390f35b610d506004803603810190610d4b9190611f38565b611bdd565b604051610d5d9190611f74565b60405180910390f35b610d806004803603810190610d7b9190611f38565b611c17565b604051610d8d9190611f74565b60405180910390f35b610db06004803603810190610dab9190611f38565b611c53565b604051610dbd9190611f74565b60405180910390f35b610de06004803603810190610ddb9190611f38565b611c8f565b604051610ded9190611f74565b60405180910390f35b610e106004803603810190610e0b9190611f38565b611cea565b604051610e1d9190611f74565b60405180910390f35b610e406004803603810190610e3b9190611f38565b611d2c565b604051610e4d9190611f74565b60405180910390f35b610e706004803603810190610e6b9190611f38565b611d68565b604051610e7d9190611f74565b60405180910390f35b610ea06004803603810190610e9b9190611f38565b611da5565b604051610ead9190611f74565b60405180910390f35b610ed06004803603810190610ecb9190611f38565b611de7565b604051610edd9190611f74565b60405180910390f35b610f006004803603810190610efb9190611f38565b611e23565b604051610f0d9190611f74565b60405180910390f35b610f306004803603810190610f2b9190611f38565b611e61565b604051610f3d9190611f74565b60405180910390f35b610f606004803603810190610f5b9190611f38565b611ea0565b604051610f6d9190611f74565b60405180910390f35b6000610f806113f5565b50600065deadbeef003690506000805b84811015610fa657369150600181019050610f90565b505080915050919050565b6000610fbb6113f5565b50600065deadbeef003290506000805b84811015610fe157329150600181019050610fcb565b505080915050919050565b6000610ff66113f5565b50600065deadbeef0052905060005b8381101561101c5781600052600181019050611005565b5080915050919050565b60006110306113f5565b50600065deadbeef0001905060005b838110156110585760008201915060018101905061103f565b5080915050919050565b60008054905090565b60006110756113f5565b50600065deadbeef0017905060005b8381101561109d57600082179150600181019050611084565b5080915050919050565b60006110b16113f5565b50600065deadbeef003490506000805b848110156110d7573491506001810190506110c1565b505080915050919050565b60006110ec6113f5565b50600065deadbeef0006905060005b83811015611133577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820691506001810190506110fb565b5080915050919050565b60006111476113f5565b50600065deadbeef001390506000805b8481101561117057600183139150600181019050611157565b505080915050919050565b60006111856113f5565b50600065deadbeef002090507fffffffff000000000000000000000000000000000000000000000000000000006000526000805b848110156111d357600460002091506001810190506111b9565b507f29045a592007d0c246ef02c2223570da9522d0cf0f73282c79a1bc8f0bb2c238811461120057600091505b5080915050919050565b60006112146113f5565b50600065deadbeef00a490508060105260005b83811015611246576004600360028360066010a4600181019050611227565b5080915050919050565b600061125a6113f5565b50600065deadbeef001a90506000805b84811015611283578260001a915060018101905061126a565b505080915050919050565b60006112986113f5565b50600065deadbeef001b905060005b838110156112c0578160001b91506001810190506112a7565b5080915050919050565b60006112d46113f5565b50600065deadbeef004290506000805b848110156112fa574291506001810190506112e4565b505080915050919050565b600061130f6113f5565b50600065deadbeef0031905060003060005b858110156113385781319250600181019050611321565b50505080915050919050565b600061134e6113f5565b50600065deadbeef004890506000805b848110156113745748915060018101905061135e565b505080915050919050565b60006113896113f5565b50600065deadbeef003d90506000805b848110156113af573d9150600181019050611399565b505080915050919050565b60006113c46113f5565b50600065deadbeef004390506000805b848110156113ea574391506001810190506113d4565b505080915050919050565b600060016000546114069190611fbe565b600081905550600054905090565b600061141e6113f5565b50600065deadbeef0004905060005b838110156114465760018204915060018101905061142d565b5080915050919050565b600061145a6113f5565b50600065deadbeef0037905060005b8381101561148257602060008037600181019050611469565b5080915050919050565b60006114966113f5565b50600065deadbeef00a090508060105260005b838110156114c15760066010a06001810190506114a9565b5080915050919050565b60006114d56113f5565b50600065deadbeef003390506000805b848110156114fb573391506001810190506114e5565b505080915050919050565b60006115106113f5565b50600065deadbeef0053905060005b8381101561153a5763deadbeef60005260018101905061151f565b5080915050919050565b600061154e6113f5565b50600065deadbeef003a90506000805b84811015611574573a915060018101905061155e565b505080915050919050565b60006115896113f5565b50600065deadbeef0051905060008160005260005b848110156115b657600051915060018101905061159e565b508091505080915050919050565b60006115ce6113f5565b50600065deadbeef001d905060005b838110156115f6578160001d91506001810190506115dd565b5080915050919050565b600061160a6113f5565b50600065deadbeef001090506000805b848110156116335782600110915060018101905061161a565b505080915050919050565b60006116486113f5565b50600065deadbeef004490506000805b8481101561166e57449150600181019050611658565b505080915050919050565b60006116836113f5565b50600065deadbeef001190506000805b848110156116ac57600183119150600181019050611693565b505080915050919050565b60006116c16113f5565b50600065deadbeef003e905060005b838110156116e95760206000803e6001810190506116d0565b5080915050919050565b60006116fd6113f5565b50600065deadbeef004590506000805b848110156117235745915060018101905061170d565b505080915050919050565b60006117386113f5565b50600065deadbeef0002905060005b8381101561176057600182029150600181019050611747565b5080915050919050565b60006117746113f5565b50600065deadbeef0008905060005b838110156117bd577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600083089150600181019050611783565b5080915050919050565b60006117d16113f5565b50600065deadbeef005490508060005560005b838110156117fc5760005491506001810190506117e4565b5080915050919050565b60006118106113f5565b50600065deadbeef005a90506000805b84811015611836575a9150600181019050611820565b505080915050919050565b600061184b6113f5565b50600065deadbeef0019905060005b83811015611871578119915060018101905061185a565b5065deadbeef0019811461188457801990505b80915050919050565b60006118976113f5565b50600065deadbeef00a190508060105260005b838110156118c3578060066010a16001810190506118aa565b5080915050919050565b60006118d76113f5565b50600065deadbeef0016905060005b838110156118fe5781821691506001810190506118e6565b5080915050919050565b60006119126113f5565b50600065deadbeef004690506000805b8481101561193857469150600181019050611922565b505080915050919050565b600061194d6113f5565b50600065deadbeef0005905060005b838110156119755760018205915060018101905061195c565b5080915050919050565b60006119896113f5565b50600065deadbeef0039905060005b838110156119b157602060008039600181019050611998565b5080915050919050565b60006119c56113f5565b50600065deadbeef005990506000805b848110156119eb575991506001810190506119d5565b505080915050919050565b6000611a006113f5565b50600065deadbeef003890506000805b84811015611a2657389150600181019050611a10565b505080915050919050565b6000611a3b6113f5565b50600065deadbeef004190506000805b84811015611a6157419150600181019050611a4b565b505080915050919050565b6000611a766113f5565b50600065deadbeef003090506000805b84811015611a9c57309150600181019050611a86565b505080915050919050565b6000611ab16113f5565b50600065deadbeef00a390508060105260005b83811015611ae157600360028260066010a3600181019050611ac4565b5080915050919050565b6000611af56113f5565b50600065deadbeef000b905060005b83811015611b1d578160200b9150600181019050611b04565b5080915050919050565b6000611b316113f5565b50600065deadbeef004790506000805b84811015611b5757479150600181019050611b41565b505080915050919050565b6000611b6c6113f5565b50600065deadbeef001c90506000805b84811015611b95578260001c9250600181019050611b7c565b505080915050919050565b6000611baa6113f5565b50600065deadbeef003590506000805b84811015611bd2576000359150600181019050611bba565b505080915050919050565b6000611be76113f5565b50600065deadbeef0055905060005b83811015611c0d5781600055600181019050611bf6565b5080915050919050565b6000611c216113f5565b50600065deadbeef0018905060005b83811015611c4957600082189150600181019050611c30565b5080915050919050565b6000611c5d6113f5565b50600065deadbeef0003905060005b83811015611c8557600082039150600181019050611c6c565b5080915050919050565b6000611c996113f5565b50600065deadbeef0007905060005b83811015611ce0577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82079150600181019050611ca8565b5080915050919050565b6000611cf46113f5565b50600065deadbeef00a290508060105260005b83811015611d225760028160066010a2600181019050611d07565b5080915050919050565b6000611d366113f5565b50600065deadbeef000a905060005b83811015611d5e576001820a9150600181019050611d45565b5080915050919050565b6000611d726113f5565b50600065deadbeef001490506000805b84811015611d9a578283149150600181019050611d82565b505080915050919050565b6000611daf6113f5565b50600065deadbeef0040905060006001430360005b85811015611ddb5781409250600181019050611dc4565b50505080915050919050565b6000611df16113f5565b50600065deadbeef001590506000805b84811015611e185782159150600181019050611e01565b505080915050919050565b6000611e2d6113f5565b50600065deadbeef001290506000805b84811015611e5657826001129150600181019050611e3d565b505080915050919050565b6000611e6b6113f5565b50600065deadbeef003b905060003060005b85811015611e9457813b9250600181019050611e7d565b50505080915050919050565b6000611eaa6113f5565b50600065deadbeef0009905060005b83811015611ef3577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600183099150600181019050611eb9565b5080915050919050565b600080fd5b6000819050919050565b611f1581611f02565b8114611f2057600080fd5b50565b600081359050611f3281611f0c565b92915050565b600060208284031215611f4e57611f4d611efd565b5b6000611f5c84828501611f23565b91505092915050565b611f6e81611f02565b82525050565b6000602082019050611f896000830184611f65565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611fc982611f02565b9150611fd483611f02565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561200957612008611f8f565b5b82820190509291505056fea2646970667358221220d8cc57ca5ddb766a9c3715937d4109a2381dca9d4ee1024c67fabd7fde18383364736f6c634300080f0033",
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
	parsed, err := abi.JSON(strings.NewReader(LoadTesterABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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
