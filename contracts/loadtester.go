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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testADD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testADDMOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testADDRESS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testAND\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testBALANCE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testBASEFEE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testBLOCKHASH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testBYTE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCALLDATACOPY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCALLDATALOAD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCALLDATASIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCALLER\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCALLVALUE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCHAINID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCODECOPY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCODESIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testCOINBASE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testDIFFICULTY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testDIV\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testEQ\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testEXP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testEXTCODESIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testGAS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testGASLIMIT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testGASPRICE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testGT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testISZERO\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testLOG0\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testLOG1\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testLOG2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testLOG3\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testLOG4\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testLT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testMLOAD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testMOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testMSIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testMSTORE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testMSTORE8\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testMUL\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testMULMOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testNOT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testNUMBER\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testORIGIN\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testRETURNDATACOPY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testRETURNDATASIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSAR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSDIV\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSELFBALANCE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSGT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSHA3\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSHL\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSHR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSIGNEXTEND\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSLOAD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSLT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSMOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSSTORE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testSUB\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testTIMESTAMP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"testXOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// LoadTesterABI is the input ABI used to generate the binding from.
// Deprecated: Use LoadTesterMetaData.ABI instead.
var LoadTesterABI = LoadTesterMetaData.ABI

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

// TestADD is a free data retrieval call binding the contract method 0x0ba8a73b.
//
// Solidity: function testADD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestADD(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testADD", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestADD is a free data retrieval call binding the contract method 0x0ba8a73b.
//
// Solidity: function testADD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestADD(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestADD(&_LoadTester.CallOpts, x)
}

// TestADD is a free data retrieval call binding the contract method 0x0ba8a73b.
//
// Solidity: function testADD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestADD(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestADD(&_LoadTester.CallOpts, x)
}

// TestADDMOD is a free data retrieval call binding the contract method 0x80947f80.
//
// Solidity: function testADDMOD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestADDMOD(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testADDMOD", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestADDMOD is a free data retrieval call binding the contract method 0x80947f80.
//
// Solidity: function testADDMOD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestADDMOD(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestADDMOD(&_LoadTester.CallOpts, x)
}

// TestADDMOD is a free data retrieval call binding the contract method 0x80947f80.
//
// Solidity: function testADDMOD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestADDMOD(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestADDMOD(&_LoadTester.CallOpts, x)
}

// TestADDRESS is a free data retrieval call binding the contract method 0xbdc875fc.
//
// Solidity: function testADDRESS(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestADDRESS(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testADDRESS", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestADDRESS is a free data retrieval call binding the contract method 0xbdc875fc.
//
// Solidity: function testADDRESS(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestADDRESS(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestADDRESS(&_LoadTester.CallOpts, x)
}

// TestADDRESS is a free data retrieval call binding the contract method 0xbdc875fc.
//
// Solidity: function testADDRESS(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestADDRESS(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestADDRESS(&_LoadTester.CallOpts, x)
}

// TestAND is a free data retrieval call binding the contract method 0x9a2b7c81.
//
// Solidity: function testAND(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestAND(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testAND", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestAND is a free data retrieval call binding the contract method 0x9a2b7c81.
//
// Solidity: function testAND(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestAND(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestAND(&_LoadTester.CallOpts, x)
}

// TestAND is a free data retrieval call binding the contract method 0x9a2b7c81.
//
// Solidity: function testAND(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestAND(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestAND(&_LoadTester.CallOpts, x)
}

// TestBALANCE is a free data retrieval call binding the contract method 0x2294fc7f.
//
// Solidity: function testBALANCE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestBALANCE(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testBALANCE", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestBALANCE is a free data retrieval call binding the contract method 0x2294fc7f.
//
// Solidity: function testBALANCE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestBALANCE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestBALANCE(&_LoadTester.CallOpts, x)
}

// TestBALANCE is a free data retrieval call binding the contract method 0x2294fc7f.
//
// Solidity: function testBALANCE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestBALANCE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestBALANCE(&_LoadTester.CallOpts, x)
}

// TestBASEFEE is a free data retrieval call binding the contract method 0x2871ef85.
//
// Solidity: function testBASEFEE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestBASEFEE(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testBASEFEE", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestBASEFEE is a free data retrieval call binding the contract method 0x2871ef85.
//
// Solidity: function testBASEFEE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestBASEFEE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestBASEFEE(&_LoadTester.CallOpts, x)
}

// TestBASEFEE is a free data retrieval call binding the contract method 0x2871ef85.
//
// Solidity: function testBASEFEE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestBASEFEE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestBASEFEE(&_LoadTester.CallOpts, x)
}

// TestBLOCKHASH is a free data retrieval call binding the contract method 0xea5141e6.
//
// Solidity: function testBLOCKHASH(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestBLOCKHASH(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testBLOCKHASH", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestBLOCKHASH is a free data retrieval call binding the contract method 0xea5141e6.
//
// Solidity: function testBLOCKHASH(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestBLOCKHASH(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestBLOCKHASH(&_LoadTester.CallOpts, x)
}

// TestBLOCKHASH is a free data retrieval call binding the contract method 0xea5141e6.
//
// Solidity: function testBLOCKHASH(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestBLOCKHASH(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestBLOCKHASH(&_LoadTester.CallOpts, x)
}

// TestBYTE is a free data retrieval call binding the contract method 0x1de2f343.
//
// Solidity: function testBYTE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestBYTE(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testBYTE", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestBYTE is a free data retrieval call binding the contract method 0x1de2f343.
//
// Solidity: function testBYTE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestBYTE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestBYTE(&_LoadTester.CallOpts, x)
}

// TestBYTE is a free data retrieval call binding the contract method 0x1de2f343.
//
// Solidity: function testBYTE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestBYTE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestBYTE(&_LoadTester.CallOpts, x)
}

// TestCALLDATACOPY is a free data retrieval call binding the contract method 0x3a425dfc.
//
// Solidity: function testCALLDATACOPY(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestCALLDATACOPY(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testCALLDATACOPY", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestCALLDATACOPY is a free data retrieval call binding the contract method 0x3a425dfc.
//
// Solidity: function testCALLDATACOPY(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestCALLDATACOPY(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCALLDATACOPY(&_LoadTester.CallOpts, x)
}

// TestCALLDATACOPY is a free data retrieval call binding the contract method 0x3a425dfc.
//
// Solidity: function testCALLDATACOPY(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestCALLDATACOPY(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCALLDATACOPY(&_LoadTester.CallOpts, x)
}

// TestCALLDATALOAD is a free data retrieval call binding the contract method 0xce3cf4ef.
//
// Solidity: function testCALLDATALOAD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestCALLDATALOAD(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testCALLDATALOAD", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestCALLDATALOAD is a free data retrieval call binding the contract method 0xce3cf4ef.
//
// Solidity: function testCALLDATALOAD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestCALLDATALOAD(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCALLDATALOAD(&_LoadTester.CallOpts, x)
}

// TestCALLDATALOAD is a free data retrieval call binding the contract method 0xce3cf4ef.
//
// Solidity: function testCALLDATALOAD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestCALLDATALOAD(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCALLDATALOAD(&_LoadTester.CallOpts, x)
}

// TestCALLDATASIZE is a free data retrieval call binding the contract method 0x034aef71.
//
// Solidity: function testCALLDATASIZE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestCALLDATASIZE(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testCALLDATASIZE", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestCALLDATASIZE is a free data retrieval call binding the contract method 0x034aef71.
//
// Solidity: function testCALLDATASIZE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestCALLDATASIZE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCALLDATASIZE(&_LoadTester.CallOpts, x)
}

// TestCALLDATASIZE is a free data retrieval call binding the contract method 0x034aef71.
//
// Solidity: function testCALLDATASIZE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestCALLDATASIZE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCALLDATASIZE(&_LoadTester.CallOpts, x)
}

// TestCALLER is a free data retrieval call binding the contract method 0x44cf3bc7.
//
// Solidity: function testCALLER(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestCALLER(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testCALLER", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestCALLER is a free data retrieval call binding the contract method 0x44cf3bc7.
//
// Solidity: function testCALLER(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestCALLER(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCALLER(&_LoadTester.CallOpts, x)
}

// TestCALLER is a free data retrieval call binding the contract method 0x44cf3bc7.
//
// Solidity: function testCALLER(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestCALLER(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCALLER(&_LoadTester.CallOpts, x)
}

// TestCALLVALUE is a free data retrieval call binding the contract method 0x1581cf19.
//
// Solidity: function testCALLVALUE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestCALLVALUE(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testCALLVALUE", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestCALLVALUE is a free data retrieval call binding the contract method 0x1581cf19.
//
// Solidity: function testCALLVALUE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestCALLVALUE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCALLVALUE(&_LoadTester.CallOpts, x)
}

// TestCALLVALUE is a free data retrieval call binding the contract method 0x1581cf19.
//
// Solidity: function testCALLVALUE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestCALLVALUE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCALLVALUE(&_LoadTester.CallOpts, x)
}

// TestCHAINID is a free data retrieval call binding the contract method 0xa60a1087.
//
// Solidity: function testCHAINID(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestCHAINID(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testCHAINID", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestCHAINID is a free data retrieval call binding the contract method 0xa60a1087.
//
// Solidity: function testCHAINID(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestCHAINID(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCHAINID(&_LoadTester.CallOpts, x)
}

// TestCHAINID is a free data retrieval call binding the contract method 0xa60a1087.
//
// Solidity: function testCHAINID(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestCHAINID(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCHAINID(&_LoadTester.CallOpts, x)
}

// TestCODECOPY is a free data retrieval call binding the contract method 0xacaebdf6.
//
// Solidity: function testCODECOPY(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestCODECOPY(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testCODECOPY", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestCODECOPY is a free data retrieval call binding the contract method 0xacaebdf6.
//
// Solidity: function testCODECOPY(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestCODECOPY(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCODECOPY(&_LoadTester.CallOpts, x)
}

// TestCODECOPY is a free data retrieval call binding the contract method 0xacaebdf6.
//
// Solidity: function testCODECOPY(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestCODECOPY(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCODECOPY(&_LoadTester.CallOpts, x)
}

// TestCODESIZE is a free data retrieval call binding the contract method 0xb7b86207.
//
// Solidity: function testCODESIZE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestCODESIZE(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testCODESIZE", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestCODESIZE is a free data retrieval call binding the contract method 0xb7b86207.
//
// Solidity: function testCODESIZE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestCODESIZE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCODESIZE(&_LoadTester.CallOpts, x)
}

// TestCODESIZE is a free data retrieval call binding the contract method 0xb7b86207.
//
// Solidity: function testCODESIZE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestCODESIZE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCODESIZE(&_LoadTester.CallOpts, x)
}

// TestCOINBASE is a free data retrieval call binding the contract method 0xb81c1484.
//
// Solidity: function testCOINBASE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestCOINBASE(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testCOINBASE", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestCOINBASE is a free data retrieval call binding the contract method 0xb81c1484.
//
// Solidity: function testCOINBASE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestCOINBASE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCOINBASE(&_LoadTester.CallOpts, x)
}

// TestCOINBASE is a free data retrieval call binding the contract method 0xb81c1484.
//
// Solidity: function testCOINBASE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestCOINBASE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestCOINBASE(&_LoadTester.CallOpts, x)
}

// TestDIFFICULTY is a free data retrieval call binding the contract method 0x6f099c8d.
//
// Solidity: function testDIFFICULTY(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestDIFFICULTY(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testDIFFICULTY", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestDIFFICULTY is a free data retrieval call binding the contract method 0x6f099c8d.
//
// Solidity: function testDIFFICULTY(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestDIFFICULTY(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestDIFFICULTY(&_LoadTester.CallOpts, x)
}

// TestDIFFICULTY is a free data retrieval call binding the contract method 0x6f099c8d.
//
// Solidity: function testDIFFICULTY(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestDIFFICULTY(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestDIFFICULTY(&_LoadTester.CallOpts, x)
}

// TestDIV is a free data retrieval call binding the contract method 0x3a411f12.
//
// Solidity: function testDIV(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestDIV(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testDIV", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestDIV is a free data retrieval call binding the contract method 0x3a411f12.
//
// Solidity: function testDIV(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestDIV(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestDIV(&_LoadTester.CallOpts, x)
}

// TestDIV is a free data retrieval call binding the contract method 0x3a411f12.
//
// Solidity: function testDIV(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestDIV(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestDIV(&_LoadTester.CallOpts, x)
}

// TestEQ is a free data retrieval call binding the contract method 0xe9f9b3f2.
//
// Solidity: function testEQ(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestEQ(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testEQ", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestEQ is a free data retrieval call binding the contract method 0xe9f9b3f2.
//
// Solidity: function testEQ(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestEQ(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestEQ(&_LoadTester.CallOpts, x)
}

// TestEQ is a free data retrieval call binding the contract method 0xe9f9b3f2.
//
// Solidity: function testEQ(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestEQ(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestEQ(&_LoadTester.CallOpts, x)
}

// TestEXP is a free data retrieval call binding the contract method 0xde97a363.
//
// Solidity: function testEXP(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestEXP(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testEXP", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestEXP is a free data retrieval call binding the contract method 0xde97a363.
//
// Solidity: function testEXP(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestEXP(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestEXP(&_LoadTester.CallOpts, x)
}

// TestEXP is a free data retrieval call binding the contract method 0xde97a363.
//
// Solidity: function testEXP(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestEXP(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestEXP(&_LoadTester.CallOpts, x)
}

// TestEXTCODESIZE is a free data retrieval call binding the contract method 0xf58fc36a.
//
// Solidity: function testEXTCODESIZE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestEXTCODESIZE(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testEXTCODESIZE", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestEXTCODESIZE is a free data retrieval call binding the contract method 0xf58fc36a.
//
// Solidity: function testEXTCODESIZE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestEXTCODESIZE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestEXTCODESIZE(&_LoadTester.CallOpts, x)
}

// TestEXTCODESIZE is a free data retrieval call binding the contract method 0xf58fc36a.
//
// Solidity: function testEXTCODESIZE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestEXTCODESIZE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestEXTCODESIZE(&_LoadTester.CallOpts, x)
}

// TestGAS is a free data retrieval call binding the contract method 0x918a5fcd.
//
// Solidity: function testGAS(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestGAS(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testGAS", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestGAS is a free data retrieval call binding the contract method 0x918a5fcd.
//
// Solidity: function testGAS(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestGAS(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestGAS(&_LoadTester.CallOpts, x)
}

// TestGAS is a free data retrieval call binding the contract method 0x918a5fcd.
//
// Solidity: function testGAS(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestGAS(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestGAS(&_LoadTester.CallOpts, x)
}

// TestGASLIMIT is a free data retrieval call binding the contract method 0x7c191d20.
//
// Solidity: function testGASLIMIT(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestGASLIMIT(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testGASLIMIT", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestGASLIMIT is a free data retrieval call binding the contract method 0x7c191d20.
//
// Solidity: function testGASLIMIT(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestGASLIMIT(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestGASLIMIT(&_LoadTester.CallOpts, x)
}

// TestGASLIMIT is a free data retrieval call binding the contract method 0x7c191d20.
//
// Solidity: function testGASLIMIT(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestGASLIMIT(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestGASLIMIT(&_LoadTester.CallOpts, x)
}

// TestGASPRICE is a free data retrieval call binding the contract method 0x4d2c74b3.
//
// Solidity: function testGASPRICE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestGASPRICE(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testGASPRICE", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestGASPRICE is a free data retrieval call binding the contract method 0x4d2c74b3.
//
// Solidity: function testGASPRICE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestGASPRICE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestGASPRICE(&_LoadTester.CallOpts, x)
}

// TestGASPRICE is a free data retrieval call binding the contract method 0x4d2c74b3.
//
// Solidity: function testGASPRICE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestGASPRICE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestGASPRICE(&_LoadTester.CallOpts, x)
}

// TestGT is a free data retrieval call binding the contract method 0x71d91d28.
//
// Solidity: function testGT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestGT(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testGT", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestGT is a free data retrieval call binding the contract method 0x71d91d28.
//
// Solidity: function testGT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestGT(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestGT(&_LoadTester.CallOpts, x)
}

// TestGT is a free data retrieval call binding the contract method 0x71d91d28.
//
// Solidity: function testGT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestGT(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestGT(&_LoadTester.CallOpts, x)
}

// TestISZERO is a free data retrieval call binding the contract method 0xf279ca81.
//
// Solidity: function testISZERO(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestISZERO(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testISZERO", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestISZERO is a free data retrieval call binding the contract method 0xf279ca81.
//
// Solidity: function testISZERO(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestISZERO(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestISZERO(&_LoadTester.CallOpts, x)
}

// TestISZERO is a free data retrieval call binding the contract method 0xf279ca81.
//
// Solidity: function testISZERO(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestISZERO(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestISZERO(&_LoadTester.CallOpts, x)
}

// TestLT is a free data retrieval call binding the contract method 0x6e7f1fe7.
//
// Solidity: function testLT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestLT(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testLT", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestLT is a free data retrieval call binding the contract method 0x6e7f1fe7.
//
// Solidity: function testLT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestLT(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestLT(&_LoadTester.CallOpts, x)
}

// TestLT is a free data retrieval call binding the contract method 0x6e7f1fe7.
//
// Solidity: function testLT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestLT(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestLT(&_LoadTester.CallOpts, x)
}

// TestMLOAD is a free data retrieval call binding the contract method 0x5590c2d9.
//
// Solidity: function testMLOAD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestMLOAD(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testMLOAD", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestMLOAD is a free data retrieval call binding the contract method 0x5590c2d9.
//
// Solidity: function testMLOAD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestMLOAD(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestMLOAD(&_LoadTester.CallOpts, x)
}

// TestMLOAD is a free data retrieval call binding the contract method 0x5590c2d9.
//
// Solidity: function testMLOAD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestMLOAD(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestMLOAD(&_LoadTester.CallOpts, x)
}

// TestMOD is a free data retrieval call binding the contract method 0x16582150.
//
// Solidity: function testMOD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestMOD(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testMOD", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestMOD is a free data retrieval call binding the contract method 0x16582150.
//
// Solidity: function testMOD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestMOD(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestMOD(&_LoadTester.CallOpts, x)
}

// TestMOD is a free data retrieval call binding the contract method 0x16582150.
//
// Solidity: function testMOD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestMOD(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestMOD(&_LoadTester.CallOpts, x)
}

// TestMSIZE is a free data retrieval call binding the contract method 0xb3d847f2.
//
// Solidity: function testMSIZE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestMSIZE(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testMSIZE", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestMSIZE is a free data retrieval call binding the contract method 0xb3d847f2.
//
// Solidity: function testMSIZE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestMSIZE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestMSIZE(&_LoadTester.CallOpts, x)
}

// TestMSIZE is a free data retrieval call binding the contract method 0xb3d847f2.
//
// Solidity: function testMSIZE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestMSIZE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestMSIZE(&_LoadTester.CallOpts, x)
}

// TestMSTORE is a free data retrieval call binding the contract method 0x087b4e84.
//
// Solidity: function testMSTORE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestMSTORE(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testMSTORE", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestMSTORE is a free data retrieval call binding the contract method 0x087b4e84.
//
// Solidity: function testMSTORE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestMSTORE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestMSTORE(&_LoadTester.CallOpts, x)
}

// TestMSTORE is a free data retrieval call binding the contract method 0x087b4e84.
//
// Solidity: function testMSTORE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestMSTORE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestMSTORE(&_LoadTester.CallOpts, x)
}

// TestMSTORE8 is a free data retrieval call binding the contract method 0x4a61af1f.
//
// Solidity: function testMSTORE8(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestMSTORE8(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testMSTORE8", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestMSTORE8 is a free data retrieval call binding the contract method 0x4a61af1f.
//
// Solidity: function testMSTORE8(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestMSTORE8(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestMSTORE8(&_LoadTester.CallOpts, x)
}

// TestMSTORE8 is a free data retrieval call binding the contract method 0x4a61af1f.
//
// Solidity: function testMSTORE8(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestMSTORE8(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestMSTORE8(&_LoadTester.CallOpts, x)
}

// TestMUL is a free data retrieval call binding the contract method 0x7de8c6f8.
//
// Solidity: function testMUL(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestMUL(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testMUL", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestMUL is a free data retrieval call binding the contract method 0x7de8c6f8.
//
// Solidity: function testMUL(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestMUL(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestMUL(&_LoadTester.CallOpts, x)
}

// TestMUL is a free data retrieval call binding the contract method 0x7de8c6f8.
//
// Solidity: function testMUL(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestMUL(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestMUL(&_LoadTester.CallOpts, x)
}

// TestMULMOD is a free data retrieval call binding the contract method 0xfde7721c.
//
// Solidity: function testMULMOD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestMULMOD(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testMULMOD", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestMULMOD is a free data retrieval call binding the contract method 0xfde7721c.
//
// Solidity: function testMULMOD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestMULMOD(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestMULMOD(&_LoadTester.CallOpts, x)
}

// TestMULMOD is a free data retrieval call binding the contract method 0xfde7721c.
//
// Solidity: function testMULMOD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestMULMOD(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestMULMOD(&_LoadTester.CallOpts, x)
}

// TestNOT is a free data retrieval call binding the contract method 0x91e7b277.
//
// Solidity: function testNOT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestNOT(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testNOT", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestNOT is a free data retrieval call binding the contract method 0x91e7b277.
//
// Solidity: function testNOT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestNOT(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestNOT(&_LoadTester.CallOpts, x)
}

// TestNOT is a free data retrieval call binding the contract method 0x91e7b277.
//
// Solidity: function testNOT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestNOT(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestNOT(&_LoadTester.CallOpts, x)
}

// TestNUMBER is a free data retrieval call binding the contract method 0x2d34e798.
//
// Solidity: function testNUMBER(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestNUMBER(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testNUMBER", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestNUMBER is a free data retrieval call binding the contract method 0x2d34e798.
//
// Solidity: function testNUMBER(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestNUMBER(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestNUMBER(&_LoadTester.CallOpts, x)
}

// TestNUMBER is a free data retrieval call binding the contract method 0x2d34e798.
//
// Solidity: function testNUMBER(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestNUMBER(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestNUMBER(&_LoadTester.CallOpts, x)
}

// TestOR is a free data retrieval call binding the contract method 0x135d52f7.
//
// Solidity: function testOR(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestOR(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testOR", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestOR is a free data retrieval call binding the contract method 0x135d52f7.
//
// Solidity: function testOR(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestOR(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestOR(&_LoadTester.CallOpts, x)
}

// TestOR is a free data retrieval call binding the contract method 0x135d52f7.
//
// Solidity: function testOR(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestOR(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestOR(&_LoadTester.CallOpts, x)
}

// TestORIGIN is a free data retrieval call binding the contract method 0x050082f8.
//
// Solidity: function testORIGIN(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestORIGIN(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testORIGIN", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestORIGIN is a free data retrieval call binding the contract method 0x050082f8.
//
// Solidity: function testORIGIN(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestORIGIN(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestORIGIN(&_LoadTester.CallOpts, x)
}

// TestORIGIN is a free data retrieval call binding the contract method 0x050082f8.
//
// Solidity: function testORIGIN(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestORIGIN(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestORIGIN(&_LoadTester.CallOpts, x)
}

// TestRETURNDATACOPY is a free data retrieval call binding the contract method 0x7b6e0b0e.
//
// Solidity: function testRETURNDATACOPY(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestRETURNDATACOPY(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testRETURNDATACOPY", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestRETURNDATACOPY is a free data retrieval call binding the contract method 0x7b6e0b0e.
//
// Solidity: function testRETURNDATACOPY(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestRETURNDATACOPY(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestRETURNDATACOPY(&_LoadTester.CallOpts, x)
}

// TestRETURNDATACOPY is a free data retrieval call binding the contract method 0x7b6e0b0e.
//
// Solidity: function testRETURNDATACOPY(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestRETURNDATACOPY(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestRETURNDATACOPY(&_LoadTester.CallOpts, x)
}

// TestRETURNDATASIZE is a free data retrieval call binding the contract method 0x2b21ef44.
//
// Solidity: function testRETURNDATASIZE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestRETURNDATASIZE(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testRETURNDATASIZE", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestRETURNDATASIZE is a free data retrieval call binding the contract method 0x2b21ef44.
//
// Solidity: function testRETURNDATASIZE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestRETURNDATASIZE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestRETURNDATASIZE(&_LoadTester.CallOpts, x)
}

// TestRETURNDATASIZE is a free data retrieval call binding the contract method 0x2b21ef44.
//
// Solidity: function testRETURNDATASIZE(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestRETURNDATASIZE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestRETURNDATASIZE(&_LoadTester.CallOpts, x)
}

// TestSAR is a free data retrieval call binding the contract method 0x60e13cde.
//
// Solidity: function testSAR(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestSAR(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testSAR", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestSAR is a free data retrieval call binding the contract method 0x60e13cde.
//
// Solidity: function testSAR(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestSAR(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSAR(&_LoadTester.CallOpts, x)
}

// TestSAR is a free data retrieval call binding the contract method 0x60e13cde.
//
// Solidity: function testSAR(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestSAR(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSAR(&_LoadTester.CallOpts, x)
}

// TestSDIV is a free data retrieval call binding the contract method 0xa645c9c2.
//
// Solidity: function testSDIV(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestSDIV(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testSDIV", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestSDIV is a free data retrieval call binding the contract method 0xa645c9c2.
//
// Solidity: function testSDIV(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestSDIV(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSDIV(&_LoadTester.CallOpts, x)
}

// TestSDIV is a free data retrieval call binding the contract method 0xa645c9c2.
//
// Solidity: function testSDIV(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestSDIV(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSDIV(&_LoadTester.CallOpts, x)
}

// TestSELFBALANCE is a free data retrieval call binding the contract method 0xc420eb61.
//
// Solidity: function testSELFBALANCE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestSELFBALANCE(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testSELFBALANCE", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestSELFBALANCE is a free data retrieval call binding the contract method 0xc420eb61.
//
// Solidity: function testSELFBALANCE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestSELFBALANCE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSELFBALANCE(&_LoadTester.CallOpts, x)
}

// TestSELFBALANCE is a free data retrieval call binding the contract method 0xc420eb61.
//
// Solidity: function testSELFBALANCE(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestSELFBALANCE(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSELFBALANCE(&_LoadTester.CallOpts, x)
}

// TestSGT is a free data retrieval call binding the contract method 0x18093b46.
//
// Solidity: function testSGT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestSGT(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testSGT", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestSGT is a free data retrieval call binding the contract method 0x18093b46.
//
// Solidity: function testSGT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestSGT(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSGT(&_LoadTester.CallOpts, x)
}

// TestSGT is a free data retrieval call binding the contract method 0x18093b46.
//
// Solidity: function testSGT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestSGT(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSGT(&_LoadTester.CallOpts, x)
}

// TestSHA3 is a free data retrieval call binding the contract method 0x19b621d6.
//
// Solidity: function testSHA3(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestSHA3(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testSHA3", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestSHA3 is a free data retrieval call binding the contract method 0x19b621d6.
//
// Solidity: function testSHA3(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestSHA3(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSHA3(&_LoadTester.CallOpts, x)
}

// TestSHA3 is a free data retrieval call binding the contract method 0x19b621d6.
//
// Solidity: function testSHA3(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestSHA3(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSHA3(&_LoadTester.CallOpts, x)
}

// TestSHL is a free data retrieval call binding the contract method 0x2007332e.
//
// Solidity: function testSHL(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestSHL(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testSHL", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestSHL is a free data retrieval call binding the contract method 0x2007332e.
//
// Solidity: function testSHL(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestSHL(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSHL(&_LoadTester.CallOpts, x)
}

// TestSHL is a free data retrieval call binding the contract method 0x2007332e.
//
// Solidity: function testSHL(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestSHL(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSHL(&_LoadTester.CallOpts, x)
}

// TestSHR is a free data retrieval call binding the contract method 0xc4bd65d5.
//
// Solidity: function testSHR(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestSHR(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testSHR", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestSHR is a free data retrieval call binding the contract method 0xc4bd65d5.
//
// Solidity: function testSHR(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestSHR(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSHR(&_LoadTester.CallOpts, x)
}

// TestSHR is a free data retrieval call binding the contract method 0xc4bd65d5.
//
// Solidity: function testSHR(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestSHR(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSHR(&_LoadTester.CallOpts, x)
}

// TestSIGNEXTEND is a free data retrieval call binding the contract method 0xc360aba6.
//
// Solidity: function testSIGNEXTEND(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestSIGNEXTEND(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testSIGNEXTEND", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestSIGNEXTEND is a free data retrieval call binding the contract method 0xc360aba6.
//
// Solidity: function testSIGNEXTEND(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestSIGNEXTEND(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSIGNEXTEND(&_LoadTester.CallOpts, x)
}

// TestSIGNEXTEND is a free data retrieval call binding the contract method 0xc360aba6.
//
// Solidity: function testSIGNEXTEND(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestSIGNEXTEND(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSIGNEXTEND(&_LoadTester.CallOpts, x)
}

// TestSLT is a free data retrieval call binding the contract method 0xf4d1fc61.
//
// Solidity: function testSLT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestSLT(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testSLT", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestSLT is a free data retrieval call binding the contract method 0xf4d1fc61.
//
// Solidity: function testSLT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestSLT(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSLT(&_LoadTester.CallOpts, x)
}

// TestSLT is a free data retrieval call binding the contract method 0xf4d1fc61.
//
// Solidity: function testSLT(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestSLT(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSLT(&_LoadTester.CallOpts, x)
}

// TestSMOD is a free data retrieval call binding the contract method 0xd93cd558.
//
// Solidity: function testSMOD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestSMOD(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testSMOD", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestSMOD is a free data retrieval call binding the contract method 0xd93cd558.
//
// Solidity: function testSMOD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestSMOD(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSMOD(&_LoadTester.CallOpts, x)
}

// TestSMOD is a free data retrieval call binding the contract method 0xd93cd558.
//
// Solidity: function testSMOD(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestSMOD(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSMOD(&_LoadTester.CallOpts, x)
}

// TestSUB is a free data retrieval call binding the contract method 0xd53ff3fd.
//
// Solidity: function testSUB(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestSUB(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testSUB", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestSUB is a free data retrieval call binding the contract method 0xd53ff3fd.
//
// Solidity: function testSUB(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestSUB(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSUB(&_LoadTester.CallOpts, x)
}

// TestSUB is a free data retrieval call binding the contract method 0xd53ff3fd.
//
// Solidity: function testSUB(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestSUB(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestSUB(&_LoadTester.CallOpts, x)
}

// TestTIMESTAMP is a free data retrieval call binding the contract method 0x219cddeb.
//
// Solidity: function testTIMESTAMP(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCaller) TestTIMESTAMP(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testTIMESTAMP", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestTIMESTAMP is a free data retrieval call binding the contract method 0x219cddeb.
//
// Solidity: function testTIMESTAMP(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterSession) TestTIMESTAMP(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestTIMESTAMP(&_LoadTester.CallOpts, x)
}

// TestTIMESTAMP is a free data retrieval call binding the contract method 0x219cddeb.
//
// Solidity: function testTIMESTAMP(uint256 x) view returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestTIMESTAMP(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestTIMESTAMP(&_LoadTester.CallOpts, x)
}

// TestXOR is a free data retrieval call binding the contract method 0xd51e7b5b.
//
// Solidity: function testXOR(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCaller) TestXOR(opts *bind.CallOpts, x *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _LoadTester.contract.Call(opts, &out, "testXOR", x)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestXOR is a free data retrieval call binding the contract method 0xd51e7b5b.
//
// Solidity: function testXOR(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterSession) TestXOR(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestXOR(&_LoadTester.CallOpts, x)
}

// TestXOR is a free data retrieval call binding the contract method 0xd51e7b5b.
//
// Solidity: function testXOR(uint256 x) pure returns(uint256)
func (_LoadTester *LoadTesterCallerSession) TestXOR(x *big.Int) (*big.Int, error) {
	return _LoadTester.Contract.TestXOR(&_LoadTester.CallOpts, x)
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
