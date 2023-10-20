// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package uniswapv3

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

// UniswapInterfaceMulticallCall is an auto generated low-level Go binding around an user-defined struct.
type UniswapInterfaceMulticallCall struct {
	Target   common.Address
	GasLimit *big.Int
	CallData []byte
}

// UniswapInterfaceMulticallResult is an auto generated low-level Go binding around an user-defined struct.
type UniswapInterfaceMulticallResult struct {
	Success    bool
	GasUsed    *big.Int
	ReturnData []byte
}

// UniswapInterfaceMulticallMetaData contains all meta data concerning the UniswapInterfaceMulticall contract.
var UniswapInterfaceMulticallMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getCurrentBlockTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getEthBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"}],\"internalType\":\"structUniswapInterfaceMulticall.Call[]\",\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"returnData\",\"type\":\"bytes\"}],\"internalType\":\"structUniswapInterfaceMulticall.Result[]\",\"name\":\"returnData\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610567806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80630f28c97d146100465780631749e1e3146100645780634d2301cc14610085575b600080fd5b61004e610098565b60405161005b919061041f565b60405180910390f35b6100776100723660046102a7565b61009c565b60405161005b929190610428565b61004e610093366004610286565b610220565b4290565b8051439060609067ffffffffffffffff811180156100b957600080fd5b506040519080825280602002602001820160405280156100f357816020015b6100e061023a565b8152602001906001900390816100d85790505b50905060005b835181101561021a57600080600086848151811061011357fe5b60200260200101516000015187858151811061012b57fe5b60200260200101516020015188868151811061014357fe5b60200260200101516040015192509250925060005a90506000808573ffffffffffffffffffffffffffffffffffffffff1685856040516101839190610403565b60006040518083038160008787f1925050503d80600081146101c1576040519150601f19603f3d011682016040523d82523d6000602084013e6101c6565b606091505b509150915060005a8403905060405180606001604052808415158152602001828152602001838152508989815181106101fb57fe5b60200260200101819052505050505050505080806001019150506100f9565b50915091565b73ffffffffffffffffffffffffffffffffffffffff163190565b604051806060016040528060001515815260200160008152602001606081525090565b803573ffffffffffffffffffffffffffffffffffffffff8116811461028157600080fd5b919050565b600060208284031215610297578081fd5b6102a08261025d565b9392505050565b600060208083850312156102b9578182fd5b823567ffffffffffffffff808211156102d0578384fd5b818501915085601f8301126102e3578384fd5b8135818111156102ef57fe5b6102fc8485830201610506565b81815284810190848601875b848110156103f457813587017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0606081838f03011215610346578a8bfd5b60408051606081018181108b8211171561035c57fe5b8252610369848d0161025d565b8152818401358c82015260608401358a811115610384578d8efd5b8085019450508e603f850112610398578c8dfd5b8b8401358a8111156103a657fe5b6103b68d85601f84011601610506565b93508084528f838287010111156103cb578d8efd5b808386018e86013783018c018d9052908101919091528552509287019290870190600101610308565b50909998505050505050505050565b6000825161041581846020870161052a565b9190910192915050565b90815260200190565b600060408083018584526020828186015281865180845260609350838701915083838202880101838901875b838110156104f6578983037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa001855281518051151584528681015187850152880151888401889052805188850181905260806104b582828801858c0161052a565b96880196601f919091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01694909401909301925090850190600101610454565b50909a9950505050505050505050565b60405181810167ffffffffffffffff8111828210171561052257fe5b604052919050565b60005b8381101561054557818101518382015260200161052d565b83811115610554576000848401525b5050505056fea164736f6c6343000706000a",
}

// UniswapInterfaceMulticallABI is the input ABI used to generate the binding from.
// Deprecated: Use UniswapInterfaceMulticallMetaData.ABI instead.
var UniswapInterfaceMulticallABI = UniswapInterfaceMulticallMetaData.ABI

// UniswapInterfaceMulticallBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UniswapInterfaceMulticallMetaData.Bin instead.
var UniswapInterfaceMulticallBin = UniswapInterfaceMulticallMetaData.Bin

// DeployUniswapInterfaceMulticall deploys a new Ethereum contract, binding an instance of UniswapInterfaceMulticall to it.
func DeployUniswapInterfaceMulticall(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UniswapInterfaceMulticall, error) {
	parsed, err := UniswapInterfaceMulticallMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UniswapInterfaceMulticallBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UniswapInterfaceMulticall{UniswapInterfaceMulticallCaller: UniswapInterfaceMulticallCaller{contract: contract}, UniswapInterfaceMulticallTransactor: UniswapInterfaceMulticallTransactor{contract: contract}, UniswapInterfaceMulticallFilterer: UniswapInterfaceMulticallFilterer{contract: contract}}, nil
}

// UniswapInterfaceMulticall is an auto generated Go binding around an Ethereum contract.
type UniswapInterfaceMulticall struct {
	UniswapInterfaceMulticallCaller     // Read-only binding to the contract
	UniswapInterfaceMulticallTransactor // Write-only binding to the contract
	UniswapInterfaceMulticallFilterer   // Log filterer for contract events
}

// UniswapInterfaceMulticallCaller is an auto generated read-only Go binding around an Ethereum contract.
type UniswapInterfaceMulticallCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapInterfaceMulticallTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UniswapInterfaceMulticallTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapInterfaceMulticallFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UniswapInterfaceMulticallFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapInterfaceMulticallSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UniswapInterfaceMulticallSession struct {
	Contract     *UniswapInterfaceMulticall // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// UniswapInterfaceMulticallCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UniswapInterfaceMulticallCallerSession struct {
	Contract *UniswapInterfaceMulticallCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// UniswapInterfaceMulticallTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UniswapInterfaceMulticallTransactorSession struct {
	Contract     *UniswapInterfaceMulticallTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// UniswapInterfaceMulticallRaw is an auto generated low-level Go binding around an Ethereum contract.
type UniswapInterfaceMulticallRaw struct {
	Contract *UniswapInterfaceMulticall // Generic contract binding to access the raw methods on
}

// UniswapInterfaceMulticallCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UniswapInterfaceMulticallCallerRaw struct {
	Contract *UniswapInterfaceMulticallCaller // Generic read-only contract binding to access the raw methods on
}

// UniswapInterfaceMulticallTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UniswapInterfaceMulticallTransactorRaw struct {
	Contract *UniswapInterfaceMulticallTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapInterfaceMulticall creates a new instance of UniswapInterfaceMulticall, bound to a specific deployed contract.
func NewUniswapInterfaceMulticall(address common.Address, backend bind.ContractBackend) (*UniswapInterfaceMulticall, error) {
	contract, err := bindUniswapInterfaceMulticall(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniswapInterfaceMulticall{UniswapInterfaceMulticallCaller: UniswapInterfaceMulticallCaller{contract: contract}, UniswapInterfaceMulticallTransactor: UniswapInterfaceMulticallTransactor{contract: contract}, UniswapInterfaceMulticallFilterer: UniswapInterfaceMulticallFilterer{contract: contract}}, nil
}

// NewUniswapInterfaceMulticallCaller creates a new read-only instance of UniswapInterfaceMulticall, bound to a specific deployed contract.
func NewUniswapInterfaceMulticallCaller(address common.Address, caller bind.ContractCaller) (*UniswapInterfaceMulticallCaller, error) {
	contract, err := bindUniswapInterfaceMulticall(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapInterfaceMulticallCaller{contract: contract}, nil
}

// NewUniswapInterfaceMulticallTransactor creates a new write-only instance of UniswapInterfaceMulticall, bound to a specific deployed contract.
func NewUniswapInterfaceMulticallTransactor(address common.Address, transactor bind.ContractTransactor) (*UniswapInterfaceMulticallTransactor, error) {
	contract, err := bindUniswapInterfaceMulticall(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapInterfaceMulticallTransactor{contract: contract}, nil
}

// NewUniswapInterfaceMulticallFilterer creates a new log filterer instance of UniswapInterfaceMulticall, bound to a specific deployed contract.
func NewUniswapInterfaceMulticallFilterer(address common.Address, filterer bind.ContractFilterer) (*UniswapInterfaceMulticallFilterer, error) {
	contract, err := bindUniswapInterfaceMulticall(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapInterfaceMulticallFilterer{contract: contract}, nil
}

// bindUniswapInterfaceMulticall binds a generic wrapper to an already deployed contract.
func bindUniswapInterfaceMulticall(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UniswapInterfaceMulticallMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapInterfaceMulticall.Contract.UniswapInterfaceMulticallCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapInterfaceMulticall.Contract.UniswapInterfaceMulticallTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapInterfaceMulticall.Contract.UniswapInterfaceMulticallTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapInterfaceMulticall.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapInterfaceMulticall.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapInterfaceMulticall.Contract.contract.Transact(opts, method, params...)
}

// GetCurrentBlockTimestamp is a free data retrieval call binding the contract method 0x0f28c97d.
//
// Solidity: function getCurrentBlockTimestamp() view returns(uint256 timestamp)
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallCaller) GetCurrentBlockTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapInterfaceMulticall.contract.Call(opts, &out, "getCurrentBlockTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentBlockTimestamp is a free data retrieval call binding the contract method 0x0f28c97d.
//
// Solidity: function getCurrentBlockTimestamp() view returns(uint256 timestamp)
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallSession) GetCurrentBlockTimestamp() (*big.Int, error) {
	return _UniswapInterfaceMulticall.Contract.GetCurrentBlockTimestamp(&_UniswapInterfaceMulticall.CallOpts)
}

// GetCurrentBlockTimestamp is a free data retrieval call binding the contract method 0x0f28c97d.
//
// Solidity: function getCurrentBlockTimestamp() view returns(uint256 timestamp)
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallCallerSession) GetCurrentBlockTimestamp() (*big.Int, error) {
	return _UniswapInterfaceMulticall.Contract.GetCurrentBlockTimestamp(&_UniswapInterfaceMulticall.CallOpts)
}

// GetEthBalance is a free data retrieval call binding the contract method 0x4d2301cc.
//
// Solidity: function getEthBalance(address addr) view returns(uint256 balance)
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallCaller) GetEthBalance(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _UniswapInterfaceMulticall.contract.Call(opts, &out, "getEthBalance", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEthBalance is a free data retrieval call binding the contract method 0x4d2301cc.
//
// Solidity: function getEthBalance(address addr) view returns(uint256 balance)
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallSession) GetEthBalance(addr common.Address) (*big.Int, error) {
	return _UniswapInterfaceMulticall.Contract.GetEthBalance(&_UniswapInterfaceMulticall.CallOpts, addr)
}

// GetEthBalance is a free data retrieval call binding the contract method 0x4d2301cc.
//
// Solidity: function getEthBalance(address addr) view returns(uint256 balance)
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallCallerSession) GetEthBalance(addr common.Address) (*big.Int, error) {
	return _UniswapInterfaceMulticall.Contract.GetEthBalance(&_UniswapInterfaceMulticall.CallOpts, addr)
}

// Multicall is a paid mutator transaction binding the contract method 0x1749e1e3.
//
// Solidity: function multicall((address,uint256,bytes)[] calls) returns(uint256 blockNumber, (bool,uint256,bytes)[] returnData)
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallTransactor) Multicall(opts *bind.TransactOpts, calls []UniswapInterfaceMulticallCall) (*types.Transaction, error) {
	return _UniswapInterfaceMulticall.contract.Transact(opts, "multicall", calls)
}

// Multicall is a paid mutator transaction binding the contract method 0x1749e1e3.
//
// Solidity: function multicall((address,uint256,bytes)[] calls) returns(uint256 blockNumber, (bool,uint256,bytes)[] returnData)
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallSession) Multicall(calls []UniswapInterfaceMulticallCall) (*types.Transaction, error) {
	return _UniswapInterfaceMulticall.Contract.Multicall(&_UniswapInterfaceMulticall.TransactOpts, calls)
}

// Multicall is a paid mutator transaction binding the contract method 0x1749e1e3.
//
// Solidity: function multicall((address,uint256,bytes)[] calls) returns(uint256 blockNumber, (bool,uint256,bytes)[] returnData)
func (_UniswapInterfaceMulticall *UniswapInterfaceMulticallTransactorSession) Multicall(calls []UniswapInterfaceMulticallCall) (*types.Transaction, error) {
	return _UniswapInterfaceMulticall.Contract.Multicall(&_UniswapInterfaceMulticall.TransactOpts, calls)
}

