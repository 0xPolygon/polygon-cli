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
	Bin: "0x608060405234801561001057600080fd5b506040516106ff3803806106ff83398101604081905261002f9161015b565b600061003b8282610281565b5050610344565b634e487b7160e01b600052604160045260246000fd5b601f19601f83011681018181106001600160401b038211171561007d5761007d610042565b6040525050565b600061008f60405190565b905061009b8282610058565b919050565b60006001600160401b038211156100b9576100b9610042565b601f19601f83011660200192915050565b60005b838110156100e55781810151838201526020016100cd565b50506000910152565b60006101016100fc846100a0565b610084565b90508281526020810184848401111561011c5761011c600080fd5b6101278482856100ca565b509392505050565b600082601f83011261014357610143600080fd5b81516101538482602086016100ee565b949350505050565b60006020828403121561017057610170600080fd5b81516001600160401b0381111561018957610189600080fd5b6101538482850161012f565b634e487b7160e01b600052602260045260246000fd5b6002810460018216806101bf57607f821691505b6020821081036101d1576101d1610195565b50919050565b60006101e66101e38381565b90565b92915050565b6101f5836101d7565b815460001960089490940293841b1916921b91909117905550565b600061021d8184846101ec565b505050565b8181101561023d57610235600082610210565b600101610222565b5050565b601f82111561021d576000818152602090206020601f850104810160208510156102685750805b61027a6020601f860104830182610222565b5050505050565b81516001600160401b0381111561029a5761029a610042565b6102a482546101ab565b6102af828285610241565b6020601f8311600181146102e357600084156102cb5750858201515b600019600886021c198116600286021786555061033c565b600085815260208120601f198616915b8281101561031357888501518255602094850194600190920191016102f3565b8683101561032f5784890151600019601f89166008021c191682555b6001600288020188555050505b505050505050565b6103ac806103536000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806306fdde031461005c578063242e7fa11461007a57806327e235e3146100b2578063a26388bb146100df578063b6b55f25146100e9575b600080fd5b6100646100fc565b6040516100719190610257565b60405180910390f35b610064604051806040016040528060198152602001785465737420526576657274204572726f72204d65737361676560381b81525081565b6100d26100c03660046102a4565b60016020526000908152604090205481565b60405161007191906102cd565b6100e761018a565b005b6100e76100f73660046102ea565b6101da565b6000805461010990610321565b80601f016020809104026020016040519081016040528092919081815260200182805461013590610321565b80156101825780601f1061015757610100808354040283529160200191610182565b820191906000526020600020905b81548152906001019060200180831161016557829003601f168201915b505050505081565b60408051808201825260198152785465737420526576657274204572726f72204d65737361676560381b6020820152905162461bcd60e51b81526101d19190600401610257565b60405180910390fd5b33600090815260016020526040812080548392906101f9908490610363565b909155505050565b60005b8381101561021c578181015183820152602001610204565b50506000910152565b600061022f825190565b808452602084019350610246818560208601610201565b601f01601f19169290920192915050565b602080825281016102688184610225565b9392505050565b60006001600160a01b0382165b92915050565b61028b8161026f565b811461029657600080fd5b50565b803561027c81610282565b6000602082840312156102b9576102b9600080fd5b60006102c58484610299565b949350505050565b8181526020810161027c565b8061028b565b803561027c816102d9565b6000602082840312156102ff576102ff600080fd5b60006102c584846102df565b634e487b7160e01b600052602260045260246000fd5b60028104600182168061033557607f821691505b6020821081036103475761034761030b565b50919050565b634e487b7160e01b600052601160045260246000fd5b8082018082111561027c5761027c61034d56fea2646970667358221220cb58571a678e4a04b84ef57d328d854bedb20916a4a95a84ccc3bae10d9cc02e64736f6c63430008170033",
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
