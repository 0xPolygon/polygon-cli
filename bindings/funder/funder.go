// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package funder

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

// FunderMetaData contains all meta data concerning the Funder contract.
var FunderMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"amount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"bulkFund\",\"inputs\":[{\"name\":\"_addresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"fund\",\"inputs\":[{\"name\":\"_address\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516105bf3803806105bf83398101604081905261002f91610069565b600081116100585760405162461bcd60e51b815260040161004f90610092565b60405180910390fd5b6000556100e4565b80515b92915050565b60006020828403121561007e5761007e600080fd5b600061008a8484610060565b949350505050565b6020808252810161006381602e81527f5468652066756e64696e6720616d6f756e742073686f756c642062652067726560208201526d61746572207468616e207a65726f60901b604082015260600190565b6104cc806100f36000396000f3fe6080604052600436106100385760003560e01c80632302440814610044578063a4626b8514610066578063aa8c217c1461008657600080fd5b3661003f57005b600080fd5b34801561005057600080fd5b5061006461005f366004610229565b6100b2565b005b34801561007257600080fd5b506100646100813660046102a4565b610185565b34801561009257600080fd5b5061009c60005481565b6040516100a991906102ec565b60405180910390f35b6001600160a01b0381166100e15760405162461bcd60e51b81526004016100d890610355565b60405180910390fd5b6000544710156101035760405162461bcd60e51b81526004016100d8906103ab565b6000816001600160a01b031660005460405161011e906103bb565b60006040518083038185875af1925050503d806000811461015b576040519150601f19603f3d011682016040523d82523d6000602084013e610160565b606091505b50509050806101815760405162461bcd60e51b81526004016100d8906103c3565b5050565b600054610193908290610405565b4710156101b25760405162461bcd60e51b81526004016100d890610470565b60005b818110156101ef576101e78383838181106101d2576101d2610480565b905060200201602081019061005f9190610229565b6001016101b5565b505050565b60006001600160a01b0382165b92915050565b610210816101f4565b811461021b57600080fd5b50565b803561020181610207565b60006020828403121561023e5761023e600080fd5b600061024a848461021e565b949350505050565b60008083601f84011261026757610267600080fd5b50813567ffffffffffffffff81111561028257610282600080fd5b60208301915083602082028301111561029d5761029d600080fd5b9250929050565b600080602083850312156102ba576102ba600080fd5b823567ffffffffffffffff8111156102d4576102d4600080fd5b6102e085828601610252565b92509250509250929050565b81815260208101610201565b603c81526000602082017f5468652066756e64656420616464726573732073686f756c642062652064696681527f666572656e74207468616e20746865207a65726f206164647265737300000000602082015291505b5060400190565b60208082528101610201816102f8565b602981526000602082017f496e73756666696369656e7420636f6e74726163742062616c616e636520666f815268722066756e64696e6760b81b6020820152915061034e565b6020808252810161020181610365565b600081610201565b6020808252810161020181600e81526d119d5b991a5b99c819985a5b195960921b602082015260400190565b634e487b7160e01b600052601160045260246000fd5b81810280821583820485141761041d5761041d6103ef565b5092915050565b602f81526000602082017f496e73756666696369656e7420636f6e74726163742062616c616e636520666f81526e722062617463682066756e64696e6760881b6020820152915061034e565b6020808252810161020181610424565b634e487b7160e01b600052603260045260246000fdfea264697066735822122014b6361a96a0ed451279b4bf8d9433e6e98dc532f9f7608842fa22b6c7813ba664736f6c63430008170033",
}

// FunderABI is the input ABI used to generate the binding from.
// Deprecated: Use FunderMetaData.ABI instead.
var FunderABI = FunderMetaData.ABI

// FunderBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FunderMetaData.Bin instead.
var FunderBin = FunderMetaData.Bin

// DeployFunder deploys a new Ethereum contract, binding an instance of Funder to it.
func DeployFunder(auth *bind.TransactOpts, backend bind.ContractBackend, _amount *big.Int) (common.Address, *types.Transaction, *Funder, error) {
	parsed, err := FunderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunderBin), backend, _amount)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Funder{FunderCaller: FunderCaller{contract: contract}, FunderTransactor: FunderTransactor{contract: contract}, FunderFilterer: FunderFilterer{contract: contract}}, nil
}

// Funder is an auto generated Go binding around an Ethereum contract.
type Funder struct {
	FunderCaller     // Read-only binding to the contract
	FunderTransactor // Write-only binding to the contract
	FunderFilterer   // Log filterer for contract events
}

// FunderCaller is an auto generated read-only Go binding around an Ethereum contract.
type FunderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FunderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FunderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FunderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FunderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FunderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FunderSession struct {
	Contract     *Funder           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FunderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FunderCallerSession struct {
	Contract *FunderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FunderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FunderTransactorSession struct {
	Contract     *FunderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FunderRaw is an auto generated low-level Go binding around an Ethereum contract.
type FunderRaw struct {
	Contract *Funder // Generic contract binding to access the raw methods on
}

// FunderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FunderCallerRaw struct {
	Contract *FunderCaller // Generic read-only contract binding to access the raw methods on
}

// FunderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FunderTransactorRaw struct {
	Contract *FunderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFunder creates a new instance of Funder, bound to a specific deployed contract.
func NewFunder(address common.Address, backend bind.ContractBackend) (*Funder, error) {
	contract, err := bindFunder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Funder{FunderCaller: FunderCaller{contract: contract}, FunderTransactor: FunderTransactor{contract: contract}, FunderFilterer: FunderFilterer{contract: contract}}, nil
}

// NewFunderCaller creates a new read-only instance of Funder, bound to a specific deployed contract.
func NewFunderCaller(address common.Address, caller bind.ContractCaller) (*FunderCaller, error) {
	contract, err := bindFunder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunderCaller{contract: contract}, nil
}

// NewFunderTransactor creates a new write-only instance of Funder, bound to a specific deployed contract.
func NewFunderTransactor(address common.Address, transactor bind.ContractTransactor) (*FunderTransactor, error) {
	contract, err := bindFunder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunderTransactor{contract: contract}, nil
}

// NewFunderFilterer creates a new log filterer instance of Funder, bound to a specific deployed contract.
func NewFunderFilterer(address common.Address, filterer bind.ContractFilterer) (*FunderFilterer, error) {
	contract, err := bindFunder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunderFilterer{contract: contract}, nil
}

// bindFunder binds a generic wrapper to an already deployed contract.
func bindFunder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Funder *FunderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Funder.Contract.FunderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Funder *FunderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Funder.Contract.FunderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Funder *FunderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Funder.Contract.FunderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Funder *FunderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Funder.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Funder *FunderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Funder.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Funder *FunderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Funder.Contract.contract.Transact(opts, method, params...)
}

// Amount is a free data retrieval call binding the contract method 0xaa8c217c.
//
// Solidity: function amount() view returns(uint256)
func (_Funder *FunderCaller) Amount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Funder.contract.Call(opts, &out, "amount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Amount is a free data retrieval call binding the contract method 0xaa8c217c.
//
// Solidity: function amount() view returns(uint256)
func (_Funder *FunderSession) Amount() (*big.Int, error) {
	return _Funder.Contract.Amount(&_Funder.CallOpts)
}

// Amount is a free data retrieval call binding the contract method 0xaa8c217c.
//
// Solidity: function amount() view returns(uint256)
func (_Funder *FunderCallerSession) Amount() (*big.Int, error) {
	return _Funder.Contract.Amount(&_Funder.CallOpts)
}

// BulkFund is a paid mutator transaction binding the contract method 0xa4626b85.
//
// Solidity: function bulkFund(address[] _addresses) returns()
func (_Funder *FunderTransactor) BulkFund(opts *bind.TransactOpts, _addresses []common.Address) (*types.Transaction, error) {
	return _Funder.contract.Transact(opts, "bulkFund", _addresses)
}

// BulkFund is a paid mutator transaction binding the contract method 0xa4626b85.
//
// Solidity: function bulkFund(address[] _addresses) returns()
func (_Funder *FunderSession) BulkFund(_addresses []common.Address) (*types.Transaction, error) {
	return _Funder.Contract.BulkFund(&_Funder.TransactOpts, _addresses)
}

// BulkFund is a paid mutator transaction binding the contract method 0xa4626b85.
//
// Solidity: function bulkFund(address[] _addresses) returns()
func (_Funder *FunderTransactorSession) BulkFund(_addresses []common.Address) (*types.Transaction, error) {
	return _Funder.Contract.BulkFund(&_Funder.TransactOpts, _addresses)
}

// Fund is a paid mutator transaction binding the contract method 0x23024408.
//
// Solidity: function fund(address _address) returns()
func (_Funder *FunderTransactor) Fund(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Funder.contract.Transact(opts, "fund", _address)
}

// Fund is a paid mutator transaction binding the contract method 0x23024408.
//
// Solidity: function fund(address _address) returns()
func (_Funder *FunderSession) Fund(_address common.Address) (*types.Transaction, error) {
	return _Funder.Contract.Fund(&_Funder.TransactOpts, _address)
}

// Fund is a paid mutator transaction binding the contract method 0x23024408.
//
// Solidity: function fund(address _address) returns()
func (_Funder *FunderTransactorSession) Fund(_address common.Address) (*types.Transaction, error) {
	return _Funder.Contract.Fund(&_Funder.TransactOpts, _address)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Funder *FunderTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Funder.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Funder *FunderSession) Receive() (*types.Transaction, error) {
	return _Funder.Contract.Receive(&_Funder.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Funder *FunderTransactorSession) Receive() (*types.Transaction, error) {
	return _Funder.Contract.Receive(&_Funder.TransactOpts)
}
