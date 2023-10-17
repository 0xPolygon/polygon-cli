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

// ITickLensPopulatedTick is an auto generated low-level Go binding around an user-defined struct.
type ITickLensPopulatedTick struct {
	Tick           *big.Int
	LiquidityNet   *big.Int
	LiquidityGross *big.Int
}

// TickLensMetaData contains all meta data concerning the TickLens contract.
var TickLensMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"int16\",\"name\":\"tickBitmapIndex\",\"type\":\"int16\"}],\"name\":\"getPopulatedTicksInWord\",\"outputs\":[{\"components\":[{\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"},{\"internalType\":\"int128\",\"name\":\"liquidityNet\",\"type\":\"int128\"},{\"internalType\":\"uint128\",\"name\":\"liquidityGross\",\"type\":\"uint128\"}],\"internalType\":\"structITickLens.PopulatedTick[]\",\"name\":\"populatedTicks\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610569806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063351fb47814610030575b600080fd5b61004361003e36600461037c565b610059565b60405161005091906104aa565b60405180910390f35b606060008373ffffffffffffffffffffffffffffffffffffffff16635339c296846040518263ffffffff1660e01b8152600401610096919061051b565b60206040518083038186803b1580156100ae57600080fd5b505afa1580156100c2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100e69190610492565b90506000805b610100811015610110576001811b831615610108576001909101905b6001016100ec565b5060008573ffffffffffffffffffffffffffffffffffffffff1663d0c93a7c6040518163ffffffff1660e01b815260040160206040518083038186803b15801561015957600080fd5b505afa15801561016d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061019191906103ba565b90508167ffffffffffffffff811180156101aa57600080fd5b506040519080825280602002602001820160405280156101e457816020015b6101d1610328565b8152602001906001900390816101c95790505b50935060005b61010081101561031e576001811b841615610316576040517ff30dba93000000000000000000000000000000000000000000000000000000008152600187900b60020b60081b8201830290600090819073ffffffffffffffffffffffffffffffffffffffff8b169063f30dba9390610266908690600401610529565b6101006040518083038186803b15801561027f57600080fd5b505afa158015610293573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102b791906103e2565b5050505050509150915060405180606001604052808460020b815260200182600f0b8152602001836fffffffffffffffffffffffffffffffff168152508887600190039750878151811061030757fe5b60200260200101819052505050505b6001016101ea565b5050505092915050565b604080516060810182526000808252602082018190529181019190915290565b8051801515811461035857600080fd5b919050565b805161035881610537565b805163ffffffff8116811461035857600080fd5b6000806040838503121561038e578182fd5b823561039981610537565b91506020830135600181900b81146103af578182fd5b809150509250929050565b6000602082840312156103cb578081fd5b81518060020b81146103db578182fd5b9392505050565b600080600080600080600080610100898b0312156103fe578384fd5b88516fffffffffffffffffffffffffffffffff8116811461041d578485fd5b80985050602089015180600f0b8114610434578485fd5b80975050604089015195506060890151945060808901518060060b8114610459578485fd5b935061046760a08a0161035d565b925061047560c08a01610368565b915061048360e08a01610348565b90509295985092959890939650565b6000602082840312156104a3578081fd5b5051919050565b602080825282518282018190526000919060409081850190868401855b8281101561050e578151805160020b855286810151600f0b878601528501516fffffffffffffffffffffffffffffffff1685850152606090930192908501906001016104c7565b5091979650505050505050565b60019190910b815260200190565b60029190910b815260200190565b73ffffffffffffffffffffffffffffffffffffffff8116811461055957600080fd5b5056fea164736f6c6343000706000a",
}

// TickLensABI is the input ABI used to generate the binding from.
// Deprecated: Use TickLensMetaData.ABI instead.
var TickLensABI = TickLensMetaData.ABI

// TickLensBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TickLensMetaData.Bin instead.
var TickLensBin = TickLensMetaData.Bin

// DeployTickLens deploys a new Ethereum contract, binding an instance of TickLens to it.
func DeployTickLens(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TickLens, error) {
	parsed, err := TickLensMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TickLensBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TickLens{TickLensCaller: TickLensCaller{contract: contract}, TickLensTransactor: TickLensTransactor{contract: contract}, TickLensFilterer: TickLensFilterer{contract: contract}}, nil
}

// TickLens is an auto generated Go binding around an Ethereum contract.
type TickLens struct {
	TickLensCaller     // Read-only binding to the contract
	TickLensTransactor // Write-only binding to the contract
	TickLensFilterer   // Log filterer for contract events
}

// TickLensCaller is an auto generated read-only Go binding around an Ethereum contract.
type TickLensCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TickLensTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TickLensTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TickLensFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TickLensFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TickLensSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TickLensSession struct {
	Contract     *TickLens         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TickLensCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TickLensCallerSession struct {
	Contract *TickLensCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TickLensTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TickLensTransactorSession struct {
	Contract     *TickLensTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TickLensRaw is an auto generated low-level Go binding around an Ethereum contract.
type TickLensRaw struct {
	Contract *TickLens // Generic contract binding to access the raw methods on
}

// TickLensCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TickLensCallerRaw struct {
	Contract *TickLensCaller // Generic read-only contract binding to access the raw methods on
}

// TickLensTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TickLensTransactorRaw struct {
	Contract *TickLensTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTickLens creates a new instance of TickLens, bound to a specific deployed contract.
func NewTickLens(address common.Address, backend bind.ContractBackend) (*TickLens, error) {
	contract, err := bindTickLens(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TickLens{TickLensCaller: TickLensCaller{contract: contract}, TickLensTransactor: TickLensTransactor{contract: contract}, TickLensFilterer: TickLensFilterer{contract: contract}}, nil
}

// NewTickLensCaller creates a new read-only instance of TickLens, bound to a specific deployed contract.
func NewTickLensCaller(address common.Address, caller bind.ContractCaller) (*TickLensCaller, error) {
	contract, err := bindTickLens(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TickLensCaller{contract: contract}, nil
}

// NewTickLensTransactor creates a new write-only instance of TickLens, bound to a specific deployed contract.
func NewTickLensTransactor(address common.Address, transactor bind.ContractTransactor) (*TickLensTransactor, error) {
	contract, err := bindTickLens(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TickLensTransactor{contract: contract}, nil
}

// NewTickLensFilterer creates a new log filterer instance of TickLens, bound to a specific deployed contract.
func NewTickLensFilterer(address common.Address, filterer bind.ContractFilterer) (*TickLensFilterer, error) {
	contract, err := bindTickLens(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TickLensFilterer{contract: contract}, nil
}

// bindTickLens binds a generic wrapper to an already deployed contract.
func bindTickLens(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TickLensMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TickLens *TickLensRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TickLens.Contract.TickLensCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TickLens *TickLensRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TickLens.Contract.TickLensTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TickLens *TickLensRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TickLens.Contract.TickLensTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TickLens *TickLensCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TickLens.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TickLens *TickLensTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TickLens.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TickLens *TickLensTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TickLens.Contract.contract.Transact(opts, method, params...)
}

// GetPopulatedTicksInWord is a free data retrieval call binding the contract method 0x351fb478.
//
// Solidity: function getPopulatedTicksInWord(address pool, int16 tickBitmapIndex) view returns((int24,int128,uint128)[] populatedTicks)
func (_TickLens *TickLensCaller) GetPopulatedTicksInWord(opts *bind.CallOpts, pool common.Address, tickBitmapIndex int16) ([]ITickLensPopulatedTick, error) {
	var out []interface{}
	err := _TickLens.contract.Call(opts, &out, "getPopulatedTicksInWord", pool, tickBitmapIndex)

	if err != nil {
		return *new([]ITickLensPopulatedTick), err
	}

	out0 := *abi.ConvertType(out[0], new([]ITickLensPopulatedTick)).(*[]ITickLensPopulatedTick)

	return out0, err

}

// GetPopulatedTicksInWord is a free data retrieval call binding the contract method 0x351fb478.
//
// Solidity: function getPopulatedTicksInWord(address pool, int16 tickBitmapIndex) view returns((int24,int128,uint128)[] populatedTicks)
func (_TickLens *TickLensSession) GetPopulatedTicksInWord(pool common.Address, tickBitmapIndex int16) ([]ITickLensPopulatedTick, error) {
	return _TickLens.Contract.GetPopulatedTicksInWord(&_TickLens.CallOpts, pool, tickBitmapIndex)
}

// GetPopulatedTicksInWord is a free data retrieval call binding the contract method 0x351fb478.
//
// Solidity: function getPopulatedTicksInWord(address pool, int16 tickBitmapIndex) view returns((int24,int128,uint128)[] populatedTicks)
func (_TickLens *TickLensCallerSession) GetPopulatedTicksInWord(pool common.Address, tickBitmapIndex int16) ([]ITickLensPopulatedTick, error) {
	return _TickLens.Contract.GetPopulatedTicksInWord(&_TickLens.CallOpts, pool, tickBitmapIndex)
}

