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

// ConformanceTesterMetaData contains all meta data concerning the ConformanceTester contract.
var ConformanceTesterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_name\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"RevertErrorMessage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balances\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"testRevert\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"pure\"}]",
	Bin: "0x608060405234801562000010575f80fd5b5060405162000a8338038062000a838339818101604052810190620000369190620001d3565b805f908162000046919062000459565b50506200053d565b5f604051905090565b5f80fd5b5f80fd5b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b620000af8262000067565b810181811067ffffffffffffffff82111715620000d157620000d062000077565b5b80604052505050565b5f620000e56200004e565b9050620000f38282620000a4565b919050565b5f67ffffffffffffffff82111562000115576200011462000077565b5b620001208262000067565b9050602081019050919050565b5f5b838110156200014c5780820151818401526020810190506200012f565b5f8484015250505050565b5f6200016d6200016784620000f8565b620000da565b9050828152602081018484840111156200018c576200018b62000063565b5b620001998482856200012d565b509392505050565b5f82601f830112620001b857620001b76200005f565b5b8151620001ca84826020860162000157565b91505092915050565b5f60208284031215620001eb57620001ea62000057565b5b5f82015167ffffffffffffffff8111156200020b576200020a6200005b565b5b6200021984828501620001a1565b91505092915050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806200027157607f821691505b6020821081036200028757620002866200022c565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f60088302620002eb7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82620002ae565b620002f78683620002ae565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f620003416200033b62000335846200030f565b62000318565b6200030f565b9050919050565b5f819050919050565b6200035c8362000321565b620003746200036b8262000348565b848454620002ba565b825550505050565b5f90565b6200038a6200037c565b6200039781848462000351565b505050565b5b81811015620003be57620003b25f8262000380565b6001810190506200039d565b5050565b601f8211156200040d57620003d7816200028d565b620003e2846200029f565b81016020851015620003f2578190505b6200040a62000401856200029f565b8301826200039c565b50505b505050565b5f82821c905092915050565b5f6200042f5f198460080262000412565b1980831691505092915050565b5f6200044983836200041e565b9150826002028217905092915050565b620004648262000222565b67ffffffffffffffff81111562000480576200047f62000077565b5b6200048c825462000259565b62000499828285620003c2565b5f60209050601f831160018114620004cf575f8415620004ba578287015190505b620004c685826200043c565b86555062000535565b601f198416620004df866200028d565b5f5b828110156200050857848901518255600182019150602085019450602081019050620004e1565b8683101562000528578489015162000524601f8916826200041e565b8355505b6001600288020188555050505b505050505050565b610538806200054b5f395ff3fe608060405234801561000f575f80fd5b5060043610610055575f3560e01c806306fdde0314610059578063242e7fa11461007757806327e235e314610095578063a26388bb146100c5578063b6b55f25146100cf575b5f80fd5b6100616100eb565b60405161006e9190610316565b60405180910390f35b61007f610176565b60405161008c9190610316565b60405180910390f35b6100af60048036038101906100aa9190610394565b6101af565b6040516100bc91906103d7565b60405180910390f35b6100cd6101c4565b005b6100e960048036038101906100e4919061041a565b610236565b005b5f80546100f790610472565b80601f016020809104026020016040519081016040528092919081815260200182805461012390610472565b801561016e5780601f106101455761010080835404028352916020019161016e565b820191905f5260205f20905b81548152906001019060200180831161015157829003601f168201915b505050505081565b6040518060400160405280601981526020017f5465737420526576657274204572726f72204d6573736167650000000000000081525081565b6001602052805f5260405f205f915090505481565b6040518060400160405280601981526020017f5465737420526576657274204572726f72204d657373616765000000000000008152506040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161022d9190610316565b60405180910390fd5b8060015f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461028291906104cf565b9250508190555050565b5f81519050919050565b5f82825260208201905092915050565b5f5b838110156102c35780820151818401526020810190506102a8565b5f8484015250505050565b5f601f19601f8301169050919050565b5f6102e88261028c565b6102f28185610296565b93506103028185602086016102a6565b61030b816102ce565b840191505092915050565b5f6020820190508181035f83015261032e81846102de565b905092915050565b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6103638261033a565b9050919050565b61037381610359565b811461037d575f80fd5b50565b5f8135905061038e8161036a565b92915050565b5f602082840312156103a9576103a8610336565b5b5f6103b684828501610380565b91505092915050565b5f819050919050565b6103d1816103bf565b82525050565b5f6020820190506103ea5f8301846103c8565b92915050565b6103f9816103bf565b8114610403575f80fd5b50565b5f81359050610414816103f0565b92915050565b5f6020828403121561042f5761042e610336565b5b5f61043c84828501610406565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061048957607f821691505b60208210810361049c5761049b610445565b5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6104d9826103bf565b91506104e4836103bf565b92508282019050808211156104fc576104fb6104a2565b5b9291505056fea26469706673582212204f6eddedd8603edfb81c671534c31cecace7232c7bdf0fdf38aed9903576a6c064736f6c63430008170033",
}

// ConformanceTesterABI is the input ABI used to generate the binding from.
// Deprecated: Use ConformanceTesterMetaData.ABI instead.
var ConformanceTesterABI = ConformanceTesterMetaData.ABI

// ConformanceTesterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConformanceTesterMetaData.Bin instead.
var ConformanceTesterBin = ConformanceTesterMetaData.Bin

// DeployConformanceTester deploys a new Ethereum contract, binding an instance of ConformanceTester to it.
func DeployConformanceTester(auth *bind.TransactOpts, backend bind.ContractBackend, _name string) (common.Address, *types.Transaction, *ConformanceTester, error) {
	parsed, err := ConformanceTesterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConformanceTesterBin), backend, _name)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConformanceTester{ConformanceTesterCaller: ConformanceTesterCaller{contract: contract}, ConformanceTesterTransactor: ConformanceTesterTransactor{contract: contract}, ConformanceTesterFilterer: ConformanceTesterFilterer{contract: contract}}, nil
}

// ConformanceTester is an auto generated Go binding around an Ethereum contract.
type ConformanceTester struct {
	ConformanceTesterCaller     // Read-only binding to the contract
	ConformanceTesterTransactor // Write-only binding to the contract
	ConformanceTesterFilterer   // Log filterer for contract events
}

// ConformanceTesterCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConformanceTesterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConformanceTesterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConformanceTesterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConformanceTesterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConformanceTesterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConformanceTesterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConformanceTesterSession struct {
	Contract     *ConformanceTester // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ConformanceTesterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConformanceTesterCallerSession struct {
	Contract *ConformanceTesterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ConformanceTesterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConformanceTesterTransactorSession struct {
	Contract     *ConformanceTesterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ConformanceTesterRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConformanceTesterRaw struct {
	Contract *ConformanceTester // Generic contract binding to access the raw methods on
}

// ConformanceTesterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConformanceTesterCallerRaw struct {
	Contract *ConformanceTesterCaller // Generic read-only contract binding to access the raw methods on
}

// ConformanceTesterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConformanceTesterTransactorRaw struct {
	Contract *ConformanceTesterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConformanceTester creates a new instance of ConformanceTester, bound to a specific deployed contract.
func NewConformanceTester(address common.Address, backend bind.ContractBackend) (*ConformanceTester, error) {
	contract, err := bindConformanceTester(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConformanceTester{ConformanceTesterCaller: ConformanceTesterCaller{contract: contract}, ConformanceTesterTransactor: ConformanceTesterTransactor{contract: contract}, ConformanceTesterFilterer: ConformanceTesterFilterer{contract: contract}}, nil
}

// NewConformanceTesterCaller creates a new read-only instance of ConformanceTester, bound to a specific deployed contract.
func NewConformanceTesterCaller(address common.Address, caller bind.ContractCaller) (*ConformanceTesterCaller, error) {
	contract, err := bindConformanceTester(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConformanceTesterCaller{contract: contract}, nil
}

// NewConformanceTesterTransactor creates a new write-only instance of ConformanceTester, bound to a specific deployed contract.
func NewConformanceTesterTransactor(address common.Address, transactor bind.ContractTransactor) (*ConformanceTesterTransactor, error) {
	contract, err := bindConformanceTester(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConformanceTesterTransactor{contract: contract}, nil
}

// NewConformanceTesterFilterer creates a new log filterer instance of ConformanceTester, bound to a specific deployed contract.
func NewConformanceTesterFilterer(address common.Address, filterer bind.ContractFilterer) (*ConformanceTesterFilterer, error) {
	contract, err := bindConformanceTester(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConformanceTesterFilterer{contract: contract}, nil
}

// bindConformanceTester binds a generic wrapper to an already deployed contract.
func bindConformanceTester(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConformanceTesterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConformanceTester *ConformanceTesterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConformanceTester.Contract.ConformanceTesterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConformanceTester *ConformanceTesterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConformanceTester.Contract.ConformanceTesterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConformanceTester *ConformanceTesterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConformanceTester.Contract.ConformanceTesterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConformanceTester *ConformanceTesterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConformanceTester.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConformanceTester *ConformanceTesterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConformanceTester.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConformanceTester *ConformanceTesterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConformanceTester.Contract.contract.Transact(opts, method, params...)
}

// RevertErrorMessage is a free data retrieval call binding the contract method 0x242e7fa1.
//
// Solidity: function RevertErrorMessage() view returns(string)
func (_ConformanceTester *ConformanceTesterCaller) RevertErrorMessage(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConformanceTester.contract.Call(opts, &out, "RevertErrorMessage")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// RevertErrorMessage is a free data retrieval call binding the contract method 0x242e7fa1.
//
// Solidity: function RevertErrorMessage() view returns(string)
func (_ConformanceTester *ConformanceTesterSession) RevertErrorMessage() (string, error) {
	return _ConformanceTester.Contract.RevertErrorMessage(&_ConformanceTester.CallOpts)
}

// RevertErrorMessage is a free data retrieval call binding the contract method 0x242e7fa1.
//
// Solidity: function RevertErrorMessage() view returns(string)
func (_ConformanceTester *ConformanceTesterCallerSession) RevertErrorMessage() (string, error) {
	return _ConformanceTester.Contract.RevertErrorMessage(&_ConformanceTester.CallOpts)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_ConformanceTester *ConformanceTesterCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ConformanceTester.contract.Call(opts, &out, "balances", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_ConformanceTester *ConformanceTesterSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _ConformanceTester.Contract.Balances(&_ConformanceTester.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_ConformanceTester *ConformanceTesterCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _ConformanceTester.Contract.Balances(&_ConformanceTester.CallOpts, arg0)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConformanceTester *ConformanceTesterCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConformanceTester.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConformanceTester *ConformanceTesterSession) Name() (string, error) {
	return _ConformanceTester.Contract.Name(&_ConformanceTester.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConformanceTester *ConformanceTesterCallerSession) Name() (string, error) {
	return _ConformanceTester.Contract.Name(&_ConformanceTester.CallOpts)
}

// TestRevert is a free data retrieval call binding the contract method 0xa26388bb.
//
// Solidity: function testRevert() pure returns()
func (_ConformanceTester *ConformanceTesterCaller) TestRevert(opts *bind.CallOpts) error {
	var out []interface{}
	err := _ConformanceTester.contract.Call(opts, &out, "testRevert")

	if err != nil {
		return err
	}

	return err

}

// TestRevert is a free data retrieval call binding the contract method 0xa26388bb.
//
// Solidity: function testRevert() pure returns()
func (_ConformanceTester *ConformanceTesterSession) TestRevert() error {
	return _ConformanceTester.Contract.TestRevert(&_ConformanceTester.CallOpts)
}

// TestRevert is a free data retrieval call binding the contract method 0xa26388bb.
//
// Solidity: function testRevert() pure returns()
func (_ConformanceTester *ConformanceTesterCallerSession) TestRevert() error {
	return _ConformanceTester.Contract.TestRevert(&_ConformanceTester.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_ConformanceTester *ConformanceTesterTransactor) Deposit(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _ConformanceTester.contract.Transact(opts, "deposit", amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_ConformanceTester *ConformanceTesterSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _ConformanceTester.Contract.Deposit(&_ConformanceTester.TransactOpts, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_ConformanceTester *ConformanceTesterTransactorSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _ConformanceTester.Contract.Deposit(&_ConformanceTester.TransactOpts, amount)
}
