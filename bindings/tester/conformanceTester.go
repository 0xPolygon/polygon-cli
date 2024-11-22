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
	Bin: "0x608060405234801561000f575f80fd5b506040516106cc3803806106cc83398101604081905261002e9161014f565b5f610039828261026b565b505061032a565b634e487b7160e01b5f52604160045260245ffd5b601f19601f83011681018181106001600160401b038211171561007957610079610040565b6040525050565b5f61008a60405190565b90506100968282610054565b919050565b5f6001600160401b038211156100b3576100b3610040565b601f19601f83011660200192915050565b5f5b838110156100de5781810151838201526020016100c6565b50505f910152565b5f6100f86100f38461009b565b610080565b905082815260208101848484011115610112576101125f80fd5b61011d8482856100c4565b509392505050565b5f82601f830112610137576101375f80fd5b81516101478482602086016100e6565b949350505050565b5f60208284031215610162576101625f80fd5b81516001600160401b0381111561017a5761017a5f80fd5b61014784828501610125565b634e487b7160e01b5f52602260045260245ffd5b6002810460018216806101ae57607f821691505b6020821081036101c0576101c0610186565b50919050565b5f6101d46101d18381565b90565b92915050565b6101e3836101c6565b81545f1960089490940293841b1916921b91909117905550565b5f6102098184846101da565b505050565b81811015610228576102205f826101fd565b60010161020e565b5050565b601f821115610209575f818152602090206020601f850104810160208510156102525750805b6102646020601f86010483018261020e565b5050505050565b81516001600160401b0381111561028457610284610040565b61028e825461019a565b61029982828561022c565b6020601f8311600181146102cb575f84156102b45750858201515b5f19600886021c1981166002860217865550610322565b5f85815260208120601f198616915b828110156102fa57888501518255602094850194600190920191016102da565b8683101561031557848901515f19601f89166008021c191682555b6001600288020188555050505b505050505050565b610395806103375f395ff3fe608060405234801561000f575f80fd5b5060043610610055575f3560e01c806306fdde0314610059578063242e7fa11461007757806327e235e3146100af578063a26388bb146100db578063b6b55f25146100e5575b5f80fd5b6100616100f8565b60405161006e919061024c565b60405180910390f35b610061604051806040016040528060198152602001785465737420526576657274204572726f72204d65737361676560381b81525081565b6100ce6100bd366004610297565b60016020525f908152604090205481565b60405161006e91906102bd565b6100e3610183565b005b6100e36100f33660046102da565b6101d3565b5f80546101049061030c565b80601f01602080910402602001604051908101604052809291908181526020018280546101309061030c565b801561017b5780601f106101525761010080835404028352916020019161017b565b820191905f5260205f20905b81548152906001019060200180831161015e57829003601f168201915b505050505081565b60408051808201825260198152785465737420526576657274204572726f72204d65737361676560381b6020820152905162461bcd60e51b81526101ca919060040161024c565b60405180910390fd5b335f90815260016020526040812080548392906101f190849061034c565b909155505050565b5f5b838110156102135781810151838201526020016101fb565b50505f910152565b5f610224825190565b80845260208401935061023b8185602086016101f9565b601f01601f19169290920192915050565b6020808252810161025d818461021b565b9392505050565b5f6001600160a01b0382165b92915050565b61027f81610264565b8114610289575f80fd5b50565b803561027081610276565b5f602082840312156102aa576102aa5f80fd5b5f6102b5848461028c565b949350505050565b81815260208101610270565b8061027f565b8035610270816102c9565b5f602082840312156102ed576102ed5f80fd5b5f6102b584846102cf565b634e487b7160e01b5f52602260045260245ffd5b60028104600182168061032057607f821691505b602082108103610332576103326102f8565b50919050565b634e487b7160e01b5f52601160045260245ffd5b808201808211156102705761027061033856fea2646970667358221220c0196d207e13598e27fa830bcc84f4d2535e6700306401a16ac6fb24c62a04fc64736f6c63430008170033",
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
