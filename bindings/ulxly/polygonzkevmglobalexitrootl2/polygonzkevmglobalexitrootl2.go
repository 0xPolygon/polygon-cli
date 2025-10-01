// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package polygonzkevmglobalexitrootl2

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

// Polygonzkevmglobalexitrootl2MetaData contains all meta data concerning the Polygonzkevmglobalexitrootl2 contract.
var Polygonzkevmglobalexitrootl2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"GlobalExitRootAlreadySet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GlobalExitRootNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyAllowedContracts\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyGlobalExitRootRemover\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyGlobalExitRootUpdater\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"globalExitRootMap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRollupExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newRoot\",\"type\":\"bytes32\"}],\"name\":\"updateExitRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561000f575f5ffd5b506040516103fe3803806103fe833981810160405281019061003191906100c9565b8073ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff1681525050506100f4565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6100988261006f565b9050919050565b6100a88161008e565b81146100b2575f5ffd5b50565b5f815190506100c38161009f565b92915050565b5f602082840312156100de576100dd61006b565b5b5f6100eb848285016100b5565b91505092915050565b6080516102ec6101125f395f818160f2015261018101526102ec5ff3fe608060405234801561000f575f5ffd5b506004361061004a575f3560e01c806301fd90441461004e578063257b36321461006c57806333d6247d1461009c578063a3c573eb146100b8575b5f5ffd5b6100566100d6565b60405161006391906101bb565b60405180910390f35b61008660048036038101906100819190610202565b6100dc565b6040516100939190610245565b60405180910390f35b6100b660048036038101906100b19190610202565b6100f0565b005b6100c061017f565b6040516100cd919061029d565b60405180910390f35b60015481565b5f602052805f5260405f205f915090505481565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610175576040517fb49365dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060018190555050565b7f000000000000000000000000000000000000000000000000000000000000000081565b5f819050919050565b6101b5816101a3565b82525050565b5f6020820190506101ce5f8301846101ac565b92915050565b5f5ffd5b6101e1816101a3565b81146101eb575f5ffd5b50565b5f813590506101fc816101d8565b92915050565b5f60208284031215610217576102166101d4565b5b5f610224848285016101ee565b91505092915050565b5f819050919050565b61023f8161022d565b82525050565b5f6020820190506102585f830184610236565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6102878261025e565b9050919050565b6102978161027d565b82525050565b5f6020820190506102b05f83018461028e565b9291505056fea2646970667358221220c393973a9ae757dd7bd5220cac1b2d60f5f06c1013be10c2e6464de440e474f864736f6c634300081c0033",
}

// Polygonzkevmglobalexitrootl2ABI is the input ABI used to generate the binding from.
// Deprecated: Use Polygonzkevmglobalexitrootl2MetaData.ABI instead.
var Polygonzkevmglobalexitrootl2ABI = Polygonzkevmglobalexitrootl2MetaData.ABI

// Polygonzkevmglobalexitrootl2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Polygonzkevmglobalexitrootl2MetaData.Bin instead.
var Polygonzkevmglobalexitrootl2Bin = Polygonzkevmglobalexitrootl2MetaData.Bin

// DeployPolygonzkevmglobalexitrootl2 deploys a new Ethereum contract, binding an instance of Polygonzkevmglobalexitrootl2 to it.
func DeployPolygonzkevmglobalexitrootl2(auth *bind.TransactOpts, backend bind.ContractBackend, _bridgeAddress common.Address) (common.Address, *types.Transaction, *Polygonzkevmglobalexitrootl2, error) {
	parsed, err := Polygonzkevmglobalexitrootl2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Polygonzkevmglobalexitrootl2Bin), backend, _bridgeAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Polygonzkevmglobalexitrootl2{Polygonzkevmglobalexitrootl2Caller: Polygonzkevmglobalexitrootl2Caller{contract: contract}, Polygonzkevmglobalexitrootl2Transactor: Polygonzkevmglobalexitrootl2Transactor{contract: contract}, Polygonzkevmglobalexitrootl2Filterer: Polygonzkevmglobalexitrootl2Filterer{contract: contract}}, nil
}

// Polygonzkevmglobalexitrootl2 is an auto generated Go binding around an Ethereum contract.
type Polygonzkevmglobalexitrootl2 struct {
	Polygonzkevmglobalexitrootl2Caller     // Read-only binding to the contract
	Polygonzkevmglobalexitrootl2Transactor // Write-only binding to the contract
	Polygonzkevmglobalexitrootl2Filterer   // Log filterer for contract events
}

// Polygonzkevmglobalexitrootl2Caller is an auto generated read-only Go binding around an Ethereum contract.
type Polygonzkevmglobalexitrootl2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Polygonzkevmglobalexitrootl2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Polygonzkevmglobalexitrootl2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Polygonzkevmglobalexitrootl2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Polygonzkevmglobalexitrootl2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Polygonzkevmglobalexitrootl2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Polygonzkevmglobalexitrootl2Session struct {
	Contract     *Polygonzkevmglobalexitrootl2 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// Polygonzkevmglobalexitrootl2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Polygonzkevmglobalexitrootl2CallerSession struct {
	Contract *Polygonzkevmglobalexitrootl2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// Polygonzkevmglobalexitrootl2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Polygonzkevmglobalexitrootl2TransactorSession struct {
	Contract     *Polygonzkevmglobalexitrootl2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// Polygonzkevmglobalexitrootl2Raw is an auto generated low-level Go binding around an Ethereum contract.
type Polygonzkevmglobalexitrootl2Raw struct {
	Contract *Polygonzkevmglobalexitrootl2 // Generic contract binding to access the raw methods on
}

// Polygonzkevmglobalexitrootl2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Polygonzkevmglobalexitrootl2CallerRaw struct {
	Contract *Polygonzkevmglobalexitrootl2Caller // Generic read-only contract binding to access the raw methods on
}

// Polygonzkevmglobalexitrootl2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Polygonzkevmglobalexitrootl2TransactorRaw struct {
	Contract *Polygonzkevmglobalexitrootl2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewPolygonzkevmglobalexitrootl2 creates a new instance of Polygonzkevmglobalexitrootl2, bound to a specific deployed contract.
func NewPolygonzkevmglobalexitrootl2(address common.Address, backend bind.ContractBackend) (*Polygonzkevmglobalexitrootl2, error) {
	contract, err := bindPolygonzkevmglobalexitrootl2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Polygonzkevmglobalexitrootl2{Polygonzkevmglobalexitrootl2Caller: Polygonzkevmglobalexitrootl2Caller{contract: contract}, Polygonzkevmglobalexitrootl2Transactor: Polygonzkevmglobalexitrootl2Transactor{contract: contract}, Polygonzkevmglobalexitrootl2Filterer: Polygonzkevmglobalexitrootl2Filterer{contract: contract}}, nil
}

// NewPolygonzkevmglobalexitrootl2Caller creates a new read-only instance of Polygonzkevmglobalexitrootl2, bound to a specific deployed contract.
func NewPolygonzkevmglobalexitrootl2Caller(address common.Address, caller bind.ContractCaller) (*Polygonzkevmglobalexitrootl2Caller, error) {
	contract, err := bindPolygonzkevmglobalexitrootl2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Polygonzkevmglobalexitrootl2Caller{contract: contract}, nil
}

// NewPolygonzkevmglobalexitrootl2Transactor creates a new write-only instance of Polygonzkevmglobalexitrootl2, bound to a specific deployed contract.
func NewPolygonzkevmglobalexitrootl2Transactor(address common.Address, transactor bind.ContractTransactor) (*Polygonzkevmglobalexitrootl2Transactor, error) {
	contract, err := bindPolygonzkevmglobalexitrootl2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Polygonzkevmglobalexitrootl2Transactor{contract: contract}, nil
}

// NewPolygonzkevmglobalexitrootl2Filterer creates a new log filterer instance of Polygonzkevmglobalexitrootl2, bound to a specific deployed contract.
func NewPolygonzkevmglobalexitrootl2Filterer(address common.Address, filterer bind.ContractFilterer) (*Polygonzkevmglobalexitrootl2Filterer, error) {
	contract, err := bindPolygonzkevmglobalexitrootl2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Polygonzkevmglobalexitrootl2Filterer{contract: contract}, nil
}

// bindPolygonzkevmglobalexitrootl2 binds a generic wrapper to an already deployed contract.
func bindPolygonzkevmglobalexitrootl2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Polygonzkevmglobalexitrootl2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Polygonzkevmglobalexitrootl2.Contract.Polygonzkevmglobalexitrootl2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitrootl2.Contract.Polygonzkevmglobalexitrootl2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitrootl2.Contract.Polygonzkevmglobalexitrootl2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Polygonzkevmglobalexitrootl2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitrootl2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitrootl2.Contract.contract.Transact(opts, method, params...)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2Caller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonzkevmglobalexitrootl2.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2Session) BridgeAddress() (common.Address, error) {
	return _Polygonzkevmglobalexitrootl2.Contract.BridgeAddress(&_Polygonzkevmglobalexitrootl2.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2CallerSession) BridgeAddress() (common.Address, error) {
	return _Polygonzkevmglobalexitrootl2.Contract.BridgeAddress(&_Polygonzkevmglobalexitrootl2.CallOpts)
}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2Caller) GlobalExitRootMap(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Polygonzkevmglobalexitrootl2.contract.Call(opts, &out, "globalExitRootMap", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2Session) GlobalExitRootMap(arg0 [32]byte) (*big.Int, error) {
	return _Polygonzkevmglobalexitrootl2.Contract.GlobalExitRootMap(&_Polygonzkevmglobalexitrootl2.CallOpts, arg0)
}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2CallerSession) GlobalExitRootMap(arg0 [32]byte) (*big.Int, error) {
	return _Polygonzkevmglobalexitrootl2.Contract.GlobalExitRootMap(&_Polygonzkevmglobalexitrootl2.CallOpts, arg0)
}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2Caller) LastRollupExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Polygonzkevmglobalexitrootl2.contract.Call(opts, &out, "lastRollupExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2Session) LastRollupExitRoot() ([32]byte, error) {
	return _Polygonzkevmglobalexitrootl2.Contract.LastRollupExitRoot(&_Polygonzkevmglobalexitrootl2.CallOpts)
}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2CallerSession) LastRollupExitRoot() ([32]byte, error) {
	return _Polygonzkevmglobalexitrootl2.Contract.LastRollupExitRoot(&_Polygonzkevmglobalexitrootl2.CallOpts)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2Transactor) UpdateExitRoot(opts *bind.TransactOpts, newRoot [32]byte) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitrootl2.contract.Transact(opts, "updateExitRoot", newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2Session) UpdateExitRoot(newRoot [32]byte) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitrootl2.Contract.UpdateExitRoot(&_Polygonzkevmglobalexitrootl2.TransactOpts, newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Polygonzkevmglobalexitrootl2 *Polygonzkevmglobalexitrootl2TransactorSession) UpdateExitRoot(newRoot [32]byte) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitrootl2.Contract.UpdateExitRoot(&_Polygonzkevmglobalexitrootl2.TransactOpts, newRoot)
}
