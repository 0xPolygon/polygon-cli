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
	Bin: "0x608060405234801561000f575f80fd5b5060405161059438038061059483398101604081905261002e91610066565b5f81116100565760405162461bcd60e51b815260040161004d9061008c565b60405180910390fd5b5f556100de565b80515b92915050565b5f60208284031215610079576100795f80fd5b5f610084848461005d565b949350505050565b6020808252810161006081602e81527f5468652066756e64696e6720616d6f756e742073686f756c642062652067726560208201526d61746572207468616e207a65726f60901b604082015260600190565b6104a9806100eb5f395ff3fe608060405260043610610036575f3560e01c80632302440814610041578063a4626b8514610062578063aa8c217c14610081575f80fd5b3661003d57005b5f80fd5b34801561004c575f80fd5b5061006061005b366004610218565b6100ab565b005b34801561006d575f80fd5b5061006061007c36600461028c565b610178565b34801561008c575f80fd5b506100955f5481565b6040516100a291906102d1565b60405180910390f35b6001600160a01b0381166100da5760405162461bcd60e51b81526004016100d190610339565b60405180910390fd5b5f544710156100fb5760405162461bcd60e51b81526004016100d19061038e565b5f816001600160a01b03165f546040516101149061039e565b5f6040518083038185875af1925050503d805f811461014e576040519150601f19603f3d011682016040523d82523d5f602084013e610153565b606091505b50509050806101745760405162461bcd60e51b81526004016100d1906103a5565b5050565b5f546101859082906103e5565b4710156101a45760405162461bcd60e51b81526004016100d19061044f565b5f5b818110156101e0576101d88383838181106101c3576101c361045f565b905060200201602081019061005b9190610218565b6001016101a6565b505050565b5f6001600160a01b0382165b92915050565b610200816101e5565b811461020a575f80fd5b50565b80356101f1816101f7565b5f6020828403121561022b5761022b5f80fd5b5f610236848461020d565b949350505050565b5f8083601f840112610251576102515f80fd5b50813567ffffffffffffffff81111561026b5761026b5f80fd5b602083019150836020820283011115610285576102855f80fd5b9250929050565b5f80602083850312156102a0576102a05f80fd5b823567ffffffffffffffff8111156102b9576102b95f80fd5b6102c58582860161023e565b92509250509250929050565b818152602081016101f1565b603c81525f602082017f5468652066756e64656420616464726573732073686f756c642062652064696681527f666572656e74207468616e20746865207a65726f206164647265737300000000602082015291505b5060400190565b602080825281016101f1816102dd565b602981525f602082017f496e73756666696369656e7420636f6e74726163742062616c616e636520666f815268722066756e64696e6760b81b60208201529150610332565b602080825281016101f181610349565b5f816101f1565b602080825281016101f181600e81526d119d5b991a5b99c819985a5b195960921b602082015260400190565b634e487b7160e01b5f52601160045260245ffd5b8181028082158382048514176103fd576103fd6103d1565b5092915050565b602f81525f602082017f496e73756666696369656e7420636f6e74726163742062616c616e636520666f81526e722062617463682066756e64696e6760881b60208201529150610332565b602080825281016101f181610404565b634e487b7160e01b5f52603260045260245ffdfea264697066735822122056de11795a01ad6ae90ddaf642fbe18dd9bcfbd5818cef4db59577eb9109586364736f6c63430008170033",
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
