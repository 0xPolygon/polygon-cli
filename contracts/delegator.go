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

// DelegatorMetaData contains all meta data concerning the Delegator contract.
var DelegatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"packedCall\",\"type\":\"bytes\"}],\"name\":\"call\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"packedCall\",\"type\":\"bytes\"}],\"name\":\"delegateCall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"packedCall\",\"type\":\"bytes\"}],\"name\":\"loopCall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"packedCall\",\"type\":\"bytes\"}],\"name\":\"loopDelegateCall\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061052a806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80631b8b921d1461005157806356e7b7aa1461008157806362a1ffbd146100b1578063757bd1e5146100e1575b600080fd5b61006b60048036038101906100669190610406565b610111565b6040516100789190610481565b60405180910390f35b61009b60048036038101906100969190610406565b610194565b6040516100a89190610481565b60405180910390f35b6100cb60048036038101906100c69190610406565b610215565b6040516100d89190610481565b60405180910390f35b6100fb60048036038101906100f69190610406565b6102a8565b6040516101089190610481565b60405180910390f35b60008060608573ffffffffffffffffffffffffffffffffffffffff16858560405161013d9291906104db565b6000604051808303816000865af19150503d806000811461017a576040519150601f19603f3d011682016040523d82523d6000602084013e61017f565b606091505b50809250819350505081925050509392505050565b60008060608573ffffffffffffffffffffffffffffffffffffffff1685856040516101c09291906104db565b600060405180830381855af49150503d80600081146101fb576040519150601f19603f3d011682016040523d82523d6000602084013e610200565b606091505b50809250819350505081925050509392505050565b60008060605b6103e85a111561029c578573ffffffffffffffffffffffffffffffffffffffff16858560405161024c9291906104db565b6000604051808303816000865af19150503d8060008114610289576040519150601f19603f3d011682016040523d82523d6000602084013e61028e565b606091505b50809250819350505061021b565b81925050509392505050565b60008060605b6103e85a111561032d578573ffffffffffffffffffffffffffffffffffffffff1685856040516102df9291906104db565b600060405180830381855af49150503d806000811461031a576040519150601f19603f3d011682016040523d82523d6000602084013e61031f565b606091505b5080925081935050506102ae565b81925050509392505050565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061036e82610343565b9050919050565b61037e81610363565b811461038957600080fd5b50565b60008135905061039b81610375565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f8401126103c6576103c56103a1565b5b8235905067ffffffffffffffff8111156103e3576103e26103a6565b5b6020830191508360018202830111156103ff576103fe6103ab565b5b9250929050565b60008060006040848603121561041f5761041e610339565b5b600061042d8682870161038c565b935050602084013567ffffffffffffffff81111561044e5761044d61033e565b5b61045a868287016103b0565b92509250509250925092565b60008115159050919050565b61047b81610466565b82525050565b60006020820190506104966000830184610472565b92915050565b600081905092915050565b82818337600083830152505050565b60006104c2838561049c565b93506104cf8385846104a7565b82840190509392505050565b60006104e88284866104b6565b9150819050939250505056fea264697066735822122026359e91fa0fb5826a461a3e171ba836040a8f6089b79691f9cdfe45f6bc99e264736f6c634300080f0033",
}

// DelegatorABI is the input ABI used to generate the binding from.
// Deprecated: Use DelegatorMetaData.ABI instead.
var DelegatorABI = DelegatorMetaData.ABI

// DelegatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DelegatorMetaData.Bin instead.
var DelegatorBin = DelegatorMetaData.Bin

// DeployDelegator deploys a new Ethereum contract, binding an instance of Delegator to it.
func DeployDelegator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Delegator, error) {
	parsed, err := DelegatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DelegatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Delegator{DelegatorCaller: DelegatorCaller{contract: contract}, DelegatorTransactor: DelegatorTransactor{contract: contract}, DelegatorFilterer: DelegatorFilterer{contract: contract}}, nil
}

// Delegator is an auto generated Go binding around an Ethereum contract.
type Delegator struct {
	DelegatorCaller     // Read-only binding to the contract
	DelegatorTransactor // Write-only binding to the contract
	DelegatorFilterer   // Log filterer for contract events
}

// DelegatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelegatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelegatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelegatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelegatorSession struct {
	Contract     *Delegator        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DelegatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelegatorCallerSession struct {
	Contract *DelegatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// DelegatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelegatorTransactorSession struct {
	Contract     *DelegatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// DelegatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelegatorRaw struct {
	Contract *Delegator // Generic contract binding to access the raw methods on
}

// DelegatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelegatorCallerRaw struct {
	Contract *DelegatorCaller // Generic read-only contract binding to access the raw methods on
}

// DelegatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelegatorTransactorRaw struct {
	Contract *DelegatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelegator creates a new instance of Delegator, bound to a specific deployed contract.
func NewDelegator(address common.Address, backend bind.ContractBackend) (*Delegator, error) {
	contract, err := bindDelegator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Delegator{DelegatorCaller: DelegatorCaller{contract: contract}, DelegatorTransactor: DelegatorTransactor{contract: contract}, DelegatorFilterer: DelegatorFilterer{contract: contract}}, nil
}

// NewDelegatorCaller creates a new read-only instance of Delegator, bound to a specific deployed contract.
func NewDelegatorCaller(address common.Address, caller bind.ContractCaller) (*DelegatorCaller, error) {
	contract, err := bindDelegator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelegatorCaller{contract: contract}, nil
}

// NewDelegatorTransactor creates a new write-only instance of Delegator, bound to a specific deployed contract.
func NewDelegatorTransactor(address common.Address, transactor bind.ContractTransactor) (*DelegatorTransactor, error) {
	contract, err := bindDelegator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelegatorTransactor{contract: contract}, nil
}

// NewDelegatorFilterer creates a new log filterer instance of Delegator, bound to a specific deployed contract.
func NewDelegatorFilterer(address common.Address, filterer bind.ContractFilterer) (*DelegatorFilterer, error) {
	contract, err := bindDelegator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelegatorFilterer{contract: contract}, nil
}

// bindDelegator binds a generic wrapper to an already deployed contract.
func bindDelegator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DelegatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegator *DelegatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Delegator.Contract.DelegatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegator *DelegatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegator.Contract.DelegatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegator *DelegatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegator.Contract.DelegatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegator *DelegatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Delegator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegator *DelegatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegator *DelegatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegator.Contract.contract.Transact(opts, method, params...)
}

// Call is a paid mutator transaction binding the contract method 0x1b8b921d.
//
// Solidity: function call(address contractAddress, bytes packedCall) returns(bool)
func (_Delegator *DelegatorTransactor) Call(opts *bind.TransactOpts, contractAddress common.Address, packedCall []byte) (*types.Transaction, error) {
	return _Delegator.contract.Transact(opts, "call", contractAddress, packedCall)
}

// Call is a paid mutator transaction binding the contract method 0x1b8b921d.
//
// Solidity: function call(address contractAddress, bytes packedCall) returns(bool)
func (_Delegator *DelegatorSession) Call(contractAddress common.Address, packedCall []byte) (*types.Transaction, error) {
	return _Delegator.Contract.Call(&_Delegator.TransactOpts, contractAddress, packedCall)
}

// Call is a paid mutator transaction binding the contract method 0x1b8b921d.
//
// Solidity: function call(address contractAddress, bytes packedCall) returns(bool)
func (_Delegator *DelegatorTransactorSession) Call(contractAddress common.Address, packedCall []byte) (*types.Transaction, error) {
	return _Delegator.Contract.Call(&_Delegator.TransactOpts, contractAddress, packedCall)
}

// DelegateCall is a paid mutator transaction binding the contract method 0x56e7b7aa.
//
// Solidity: function delegateCall(address contractAddress, bytes packedCall) returns(bool)
func (_Delegator *DelegatorTransactor) DelegateCall(opts *bind.TransactOpts, contractAddress common.Address, packedCall []byte) (*types.Transaction, error) {
	return _Delegator.contract.Transact(opts, "delegateCall", contractAddress, packedCall)
}

// DelegateCall is a paid mutator transaction binding the contract method 0x56e7b7aa.
//
// Solidity: function delegateCall(address contractAddress, bytes packedCall) returns(bool)
func (_Delegator *DelegatorSession) DelegateCall(contractAddress common.Address, packedCall []byte) (*types.Transaction, error) {
	return _Delegator.Contract.DelegateCall(&_Delegator.TransactOpts, contractAddress, packedCall)
}

// DelegateCall is a paid mutator transaction binding the contract method 0x56e7b7aa.
//
// Solidity: function delegateCall(address contractAddress, bytes packedCall) returns(bool)
func (_Delegator *DelegatorTransactorSession) DelegateCall(contractAddress common.Address, packedCall []byte) (*types.Transaction, error) {
	return _Delegator.Contract.DelegateCall(&_Delegator.TransactOpts, contractAddress, packedCall)
}

// LoopCall is a paid mutator transaction binding the contract method 0x62a1ffbd.
//
// Solidity: function loopCall(address contractAddress, bytes packedCall) returns(bool)
func (_Delegator *DelegatorTransactor) LoopCall(opts *bind.TransactOpts, contractAddress common.Address, packedCall []byte) (*types.Transaction, error) {
	return _Delegator.contract.Transact(opts, "loopCall", contractAddress, packedCall)
}

// LoopCall is a paid mutator transaction binding the contract method 0x62a1ffbd.
//
// Solidity: function loopCall(address contractAddress, bytes packedCall) returns(bool)
func (_Delegator *DelegatorSession) LoopCall(contractAddress common.Address, packedCall []byte) (*types.Transaction, error) {
	return _Delegator.Contract.LoopCall(&_Delegator.TransactOpts, contractAddress, packedCall)
}

// LoopCall is a paid mutator transaction binding the contract method 0x62a1ffbd.
//
// Solidity: function loopCall(address contractAddress, bytes packedCall) returns(bool)
func (_Delegator *DelegatorTransactorSession) LoopCall(contractAddress common.Address, packedCall []byte) (*types.Transaction, error) {
	return _Delegator.Contract.LoopCall(&_Delegator.TransactOpts, contractAddress, packedCall)
}

// LoopDelegateCall is a paid mutator transaction binding the contract method 0x757bd1e5.
//
// Solidity: function loopDelegateCall(address contractAddress, bytes packedCall) returns(bool)
func (_Delegator *DelegatorTransactor) LoopDelegateCall(opts *bind.TransactOpts, contractAddress common.Address, packedCall []byte) (*types.Transaction, error) {
	return _Delegator.contract.Transact(opts, "loopDelegateCall", contractAddress, packedCall)
}

// LoopDelegateCall is a paid mutator transaction binding the contract method 0x757bd1e5.
//
// Solidity: function loopDelegateCall(address contractAddress, bytes packedCall) returns(bool)
func (_Delegator *DelegatorSession) LoopDelegateCall(contractAddress common.Address, packedCall []byte) (*types.Transaction, error) {
	return _Delegator.Contract.LoopDelegateCall(&_Delegator.TransactOpts, contractAddress, packedCall)
}

// LoopDelegateCall is a paid mutator transaction binding the contract method 0x757bd1e5.
//
// Solidity: function loopDelegateCall(address contractAddress, bytes packedCall) returns(bool)
func (_Delegator *DelegatorTransactorSession) LoopDelegateCall(contractAddress common.Address, packedCall []byte) (*types.Transaction, error) {
	return _Delegator.Contract.LoopDelegateCall(&_Delegator.TransactOpts, contractAddress, packedCall)
}
