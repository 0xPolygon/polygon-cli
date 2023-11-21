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

// NonfungibleTokenPositionDescriptorMetaData contains all meta data concerning the NonfungibleTokenPositionDescriptor contract.
var NonfungibleTokenPositionDescriptorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_WETH9\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_nativeCurrencyLabelBytes\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"WETH9\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"name\":\"flipRatio\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nativeCurrencyLabel\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nativeCurrencyLabelBytes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"name\":\"tokenRatioPriority\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractINonfungiblePositionManager\",\"name\":\"positionManager\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b5060405161158038038061158083398101604081905261002f9161004a565b60609190911b6001600160601b03191660805260a052610082565b6000806040838503121561005c578182fd5b82516001600160a01b0381168114610072578283fd5b6020939093015192949293505050565b60805160601c60a0516114bf6100c16000398061028f52806102c35280610377525080610107528061014c528061061a528061066e52506114bf6000f3fe608060405234801561001057600080fd5b50600436106100725760003560e01c8063a18246e211610050578063a18246e2146100d5578063b7af3cdc146100dd578063e9dc6375146100f257610072565b80634aa4a4fc146100775780637e5af771146100955780639d7b0ea8146100b5575b600080fd5b61007f610105565b60405161008c9190611311565b60405180910390f35b6100a86100a3366004611017565b610129565b60405161008c9190611325565b6100c86100c3366004611057565b610148565b60405161008c9190611330565b6100c861028d565b6100e56102b1565b60405161008c9190611339565b6100e5610100366004611057565b6103d7565b7f000000000000000000000000000000000000000000000000000000000000000081565b60006101358383610148565b61013f8584610148565b13949350505050565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316836001600160a01b0316141561018d5750606319610287565b8160011415610283576001600160a01b03831673a0b86991c6218b36c1d19d4a2e9eb0ce3606eb4814156101c4575061012c610287565b6001600160a01b03831673dac17f958d2ee523a2206206994597c13d831ec714156101f1575060c8610287565b6001600160a01b038316736b175474e89094c44da98b954eedeac495271d0f141561021e57506064610287565b6001600160a01b038316738daebade922df735c38c80c7ebd708af50815faa141561024c575060c719610287565b6001600160a01b038316732260fac5e5542a773aa44fbcfedf7c193bc2c599141561027b575061012b19610287565b506000610287565b5060005b92915050565b7f000000000000000000000000000000000000000000000000000000000000000081565b606060005b60208110801561031657507f000000000000000000000000000000000000000000000000000000000000000081602081106102ed57fe5b1a60f81b7fff000000000000000000000000000000000000000000000000000000000000001615155b15610323576001016102b6565b60008167ffffffffffffffff8111801561033c57600080fd5b506040519080825280601f01601f191660200182016040528015610367576020820181803683370190505b50905060005b828110156103d0577f000000000000000000000000000000000000000000000000000000000000000081602081106103a157fe5b1a60f81b8282815181106103b157fe5b60200101906001600160f81b031916908160001a90535060010161036d565b5091505090565b60606000806000806000876001600160a01b03166399fbab88886040518263ffffffff1660e01b815260040161040d9190611330565b6101806040518083038186803b15801561042657600080fd5b505afa15801561043a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061045e91906111dc565b5050505050965096509650965096505050600061051c896001600160a01b031663c45a01556040518163ffffffff1660e01b815260040160206040518083038186803b1580156104ad57600080fd5b505afa1580156104c1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104e59190610ff4565b6040518060600160405280896001600160a01b03168152602001886001600160a01b031681526020018762ffffff168152506108e7565b9050600061052d87876100a36109e3565b90506000811561053d578761053f565b865b90506000821561054f5787610551565b885b90506000846001600160a01b0316633850c7bd6040518163ffffffff1660e01b815260040160e06040518083038186803b15801561058e57600080fd5b505afa1580156105a2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105c69190611133565b505050505091505073__$cea9be979eee3d87fb124d6cbb244bb0b5$__63c49917d7604051806101c001604052808f8152602001866001600160a01b03168152602001856001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316876001600160a01b03161461065f5761065a876109e7565b610667565b6106676102b1565b81526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316866001600160a01b0316146106b3576106ae866109e7565b6106bb565b6106bb6102b1565b8152602001866001600160a01b031663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b1580156106f957600080fd5b505afa15801561070d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061073191906111c2565b60ff168152602001856001600160a01b031663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b15801561077257600080fd5b505afa158015610786573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107aa91906111c2565b60ff16815260200187151581526020018a60020b81526020018960020b81526020018460020b8152602001886001600160a01b031663d0c93a7c6040518163ffffffff1660e01b815260040160206040518083038186803b15801561080e57600080fd5b505afa158015610822573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108469190611082565b60020b81526020018b62ffffff168152602001886001600160a01b03168152506040518263ffffffff1660e01b8152600401610882919061134c565b60006040518083038186803b15801561089a57600080fd5b505af41580156108ae573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526108d6919081019061109c565b9d9c50505050505050505050505050565b600081602001516001600160a01b031682600001516001600160a01b03161061090f57600080fd5b50805160208083015160409384015184516001600160a01b0394851681850152939091168385015262ffffff166060808401919091528351808403820181526080840185528051908301207fff0000000000000000000000000000000000000000000000000000000000000060a085015294901b6bffffffffffffffffffffffff191660a183015260b58201939093527fe34f199b19b2b4f47f68442619d555527d244f78a3297ea89325f843f87b8b5460d5808301919091528251808303909101815260f5909101909152805191012090565b4690565b60606000610a15837f95d89b4100000000000000000000000000000000000000000000000000000000610a3a565b9050805160001415610a3257610a2a83610c8f565b915050610a35565b90505b919050565b60408051600481526024810182526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000085161781529151815160609360009384936001600160a01b03891693919290918291908083835b60208310610ad35780518252601f199092019160209182019101610ab4565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855afa9150503d8060008114610b33576040519150601f19603f3d011682016040523d82523d6000602084013e610b38565b606091505b5091509150811580610b4957508051155b15610b67576040518060200160405280600081525092505050610287565b805160201415610b9f576000818060200190516020811015610b8857600080fd5b50519050610b9581610c9c565b9350505050610287565b604081511115610c7757808060200190516020811015610bbe57600080fd5b8101908080516040519392919084640100000000821115610bde57600080fd5b908301906020820185811115610bf357600080fd5b8251640100000000811182820188101715610c0d57600080fd5b82525081516020918201929091019080838360005b83811015610c3a578181015183820152602001610c22565b50505050905090810190601f168015610c675780820380516001836020036101000a031916815260200191505b5060405250505092505050610287565b50506040805160208101909152600081529392505050565b6060610a32826006610ddc565b604080516020808252818301909252606091600091906020820181803683370190505090506000805b6020811015610d3e576000858260208110610cdc57fe5b1a60f81b90507fff00000000000000000000000000000000000000000000000000000000000000811615610d355780848481518110610d1757fe5b60200101906001600160f81b031916908160001a9053506001909201915b50600101610cc5565b5060008167ffffffffffffffff81118015610d5857600080fd5b506040519080825280601f01601f191660200182016040528015610d83576020820181803683370190505b50905060005b82811015610dd357838181518110610d9d57fe5b602001015160f81c60f81b828281518110610db457fe5b60200101906001600160f81b031916908160001a905350600101610d89565b50949350505050565b606060028206158015610def5750600082115b8015610dfc575060288211155b610e6757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f41646472657373537472696e675574696c3a20494e56414c49445f4c454e0000604482015290519081900360640190fd5b60008267ffffffffffffffff81118015610e8057600080fd5b506040519080825280601f01601f191660200182016040528015610eab576020820181803683370190505b5090506001600160a01b03841660005b60028504811015610f4f57600860138290030282901c600f600482901c1660f082168203610ee882610f59565b868560020281518110610ef757fe5b60200101906001600160f81b031916908160001a905350610f1781610f59565b868560020260010181518110610f2957fe5b60200101906001600160f81b031916908160001a9053505060019092019150610ebb9050565b5090949350505050565b6000600a8260ff161015610f7457506030810160f81b610a35565b506037810160f81b610a35565b8051610a358161149a565b8051600281900b8114610a3557600080fd5b80516fffffffffffffffffffffffffffffffff81168114610a3557600080fd5b805161ffff81168114610a3557600080fd5b805162ffffff81168114610a3557600080fd5b805160ff81168114610a3557600080fd5b600060208284031215611005578081fd5b81516110108161149a565b9392505050565b60008060006060848603121561102b578182fd5b83356110368161149a565b925060208401356110468161149a565b929592945050506040919091013590565b60008060408385031215611069578182fd5b82356110748161149a565b946020939093013593505050565b600060208284031215611093578081fd5b61101082610f8c565b6000602082840312156110ad578081fd5b815167ffffffffffffffff808211156110c4578283fd5b818401915084601f8301126110d7578283fd5b8151818111156110e357fe5b604051601f8201601f19168101602001838111828210171561110157fe5b604052818152838201602001871015611118578485fd5b61112982602083016020870161146a565b9695505050505050565b600080600080600080600060e0888a03121561114d578283fd5b87516111588161149a565b965061116660208901610f8c565b955061117460408901610fbe565b945061118260608901610fbe565b935061119060808901610fbe565b925061119e60a08901610fe3565b915060c088015180151581146111b2578182fd5b8091505092959891949750929550565b6000602082840312156111d3578081fd5b61101082610fe3565b6000806000806000806000806000806000806101808d8f0312156111fe578485fd5b8c516bffffffffffffffffffffffff81168114611219578586fd5b9b5061122760208e01610f81565b9a5061123560408e01610f81565b995061124360608e01610f81565b985061125160808e01610fd0565b975061125f60a08e01610f8c565b965061126d60c08e01610f8c565b955061127b60e08e01610f9e565b94506101008d015193506101208d0151925061129a6101408e01610f9e565b91506112a96101608e01610f9e565b90509295989b509295989b509295989b565b6001600160a01b03169052565b15159052565b60020b9052565b600081518084526112ed81602086016020860161146a565b601f01601f19169290920160200192915050565b62ffffff169052565b60ff169052565b6001600160a01b0391909116815260200190565b901515815260200190565b90815260200190565b60006020825261101060208301846112d5565b60006020825282516020830152602083015161136b60408401826112bb565b50604083015161137e60608401826112bb565b5060608301516101c080608085015261139b6101e08501836112d5565b91506080850151601f198584030160a08601526113b883826112d5565b92505060a08501516113cd60c086018261130a565b5060c08501516113e060e086018261130a565b5060e08501516101006113f5818701836112c8565b8601519050610120611409868201836112ce565b860151905061014061141d868201836112ce565b8601519050610160611431868201836112ce565b8601519050610180611445868201836112ce565b86015190506101a061145986820183611301565b8601519050610f4f858301826112bb565b60005b8381101561148557818101518382015260200161146d565b83811115611494576000848401525b50505050565b6001600160a01b03811681146114af57600080fd5b5056fea164736f6c6343000706000a",
}

// NonfungibleTokenPositionDescriptorABI is the input ABI used to generate the binding from.
// Deprecated: Use NonfungibleTokenPositionDescriptorMetaData.ABI instead.
var NonfungibleTokenPositionDescriptorABI = NonfungibleTokenPositionDescriptorMetaData.ABI

// NonfungibleTokenPositionDescriptorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NonfungibleTokenPositionDescriptorMetaData.Bin instead.
var NonfungibleTokenPositionDescriptorBin = NonfungibleTokenPositionDescriptorMetaData.Bin

// DeployNonfungibleTokenPositionDescriptor deploys a new Ethereum contract, binding an instance of NonfungibleTokenPositionDescriptor to it.
func DeployNonfungibleTokenPositionDescriptor(auth *bind.TransactOpts, backend bind.ContractBackend, _WETH9 common.Address, _nativeCurrencyLabelBytes [32]byte) (common.Address, *types.Transaction, *NonfungibleTokenPositionDescriptor, error) {
	parsed, err := NonfungibleTokenPositionDescriptorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NonfungibleTokenPositionDescriptorBin), backend, _WETH9, _nativeCurrencyLabelBytes)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NonfungibleTokenPositionDescriptor{NonfungibleTokenPositionDescriptorCaller: NonfungibleTokenPositionDescriptorCaller{contract: contract}, NonfungibleTokenPositionDescriptorTransactor: NonfungibleTokenPositionDescriptorTransactor{contract: contract}, NonfungibleTokenPositionDescriptorFilterer: NonfungibleTokenPositionDescriptorFilterer{contract: contract}}, nil
}

// NonfungibleTokenPositionDescriptor is an auto generated Go binding around an Ethereum contract.
type NonfungibleTokenPositionDescriptor struct {
	NonfungibleTokenPositionDescriptorCaller     // Read-only binding to the contract
	NonfungibleTokenPositionDescriptorTransactor // Write-only binding to the contract
	NonfungibleTokenPositionDescriptorFilterer   // Log filterer for contract events
}

// NonfungibleTokenPositionDescriptorCaller is an auto generated read-only Go binding around an Ethereum contract.
type NonfungibleTokenPositionDescriptorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NonfungibleTokenPositionDescriptorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NonfungibleTokenPositionDescriptorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NonfungibleTokenPositionDescriptorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NonfungibleTokenPositionDescriptorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NonfungibleTokenPositionDescriptorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NonfungibleTokenPositionDescriptorSession struct {
	Contract     *NonfungibleTokenPositionDescriptor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                       // Call options to use throughout this session
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// NonfungibleTokenPositionDescriptorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NonfungibleTokenPositionDescriptorCallerSession struct {
	Contract *NonfungibleTokenPositionDescriptorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                             // Call options to use throughout this session
}

// NonfungibleTokenPositionDescriptorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NonfungibleTokenPositionDescriptorTransactorSession struct {
	Contract     *NonfungibleTokenPositionDescriptorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                             // Transaction auth options to use throughout this session
}

// NonfungibleTokenPositionDescriptorRaw is an auto generated low-level Go binding around an Ethereum contract.
type NonfungibleTokenPositionDescriptorRaw struct {
	Contract *NonfungibleTokenPositionDescriptor // Generic contract binding to access the raw methods on
}

// NonfungibleTokenPositionDescriptorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NonfungibleTokenPositionDescriptorCallerRaw struct {
	Contract *NonfungibleTokenPositionDescriptorCaller // Generic read-only contract binding to access the raw methods on
}

// NonfungibleTokenPositionDescriptorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NonfungibleTokenPositionDescriptorTransactorRaw struct {
	Contract *NonfungibleTokenPositionDescriptorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNonfungibleTokenPositionDescriptor creates a new instance of NonfungibleTokenPositionDescriptor, bound to a specific deployed contract.
func NewNonfungibleTokenPositionDescriptor(address common.Address, backend bind.ContractBackend) (*NonfungibleTokenPositionDescriptor, error) {
	contract, err := bindNonfungibleTokenPositionDescriptor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NonfungibleTokenPositionDescriptor{NonfungibleTokenPositionDescriptorCaller: NonfungibleTokenPositionDescriptorCaller{contract: contract}, NonfungibleTokenPositionDescriptorTransactor: NonfungibleTokenPositionDescriptorTransactor{contract: contract}, NonfungibleTokenPositionDescriptorFilterer: NonfungibleTokenPositionDescriptorFilterer{contract: contract}}, nil
}

// NewNonfungibleTokenPositionDescriptorCaller creates a new read-only instance of NonfungibleTokenPositionDescriptor, bound to a specific deployed contract.
func NewNonfungibleTokenPositionDescriptorCaller(address common.Address, caller bind.ContractCaller) (*NonfungibleTokenPositionDescriptorCaller, error) {
	contract, err := bindNonfungibleTokenPositionDescriptor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NonfungibleTokenPositionDescriptorCaller{contract: contract}, nil
}

// NewNonfungibleTokenPositionDescriptorTransactor creates a new write-only instance of NonfungibleTokenPositionDescriptor, bound to a specific deployed contract.
func NewNonfungibleTokenPositionDescriptorTransactor(address common.Address, transactor bind.ContractTransactor) (*NonfungibleTokenPositionDescriptorTransactor, error) {
	contract, err := bindNonfungibleTokenPositionDescriptor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NonfungibleTokenPositionDescriptorTransactor{contract: contract}, nil
}

// NewNonfungibleTokenPositionDescriptorFilterer creates a new log filterer instance of NonfungibleTokenPositionDescriptor, bound to a specific deployed contract.
func NewNonfungibleTokenPositionDescriptorFilterer(address common.Address, filterer bind.ContractFilterer) (*NonfungibleTokenPositionDescriptorFilterer, error) {
	contract, err := bindNonfungibleTokenPositionDescriptor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NonfungibleTokenPositionDescriptorFilterer{contract: contract}, nil
}

// bindNonfungibleTokenPositionDescriptor binds a generic wrapper to an already deployed contract.
func bindNonfungibleTokenPositionDescriptor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NonfungibleTokenPositionDescriptorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NonfungibleTokenPositionDescriptor.Contract.NonfungibleTokenPositionDescriptorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.NonfungibleTokenPositionDescriptorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.NonfungibleTokenPositionDescriptorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NonfungibleTokenPositionDescriptor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.contract.Transact(opts, method, params...)
}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorCaller) WETH9(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NonfungibleTokenPositionDescriptor.contract.Call(opts, &out, "WETH9")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorSession) WETH9() (common.Address, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.WETH9(&_NonfungibleTokenPositionDescriptor.CallOpts)
}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorCallerSession) WETH9() (common.Address, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.WETH9(&_NonfungibleTokenPositionDescriptor.CallOpts)
}

// FlipRatio is a free data retrieval call binding the contract method 0x7e5af771.
//
// Solidity: function flipRatio(address token0, address token1, uint256 chainId) view returns(bool)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorCaller) FlipRatio(opts *bind.CallOpts, token0 common.Address, token1 common.Address, chainId *big.Int) (bool, error) {
	var out []interface{}
	err := _NonfungibleTokenPositionDescriptor.contract.Call(opts, &out, "flipRatio", token0, token1, chainId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// FlipRatio is a free data retrieval call binding the contract method 0x7e5af771.
//
// Solidity: function flipRatio(address token0, address token1, uint256 chainId) view returns(bool)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorSession) FlipRatio(token0 common.Address, token1 common.Address, chainId *big.Int) (bool, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.FlipRatio(&_NonfungibleTokenPositionDescriptor.CallOpts, token0, token1, chainId)
}

// FlipRatio is a free data retrieval call binding the contract method 0x7e5af771.
//
// Solidity: function flipRatio(address token0, address token1, uint256 chainId) view returns(bool)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorCallerSession) FlipRatio(token0 common.Address, token1 common.Address, chainId *big.Int) (bool, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.FlipRatio(&_NonfungibleTokenPositionDescriptor.CallOpts, token0, token1, chainId)
}

// NativeCurrencyLabel is a free data retrieval call binding the contract method 0xb7af3cdc.
//
// Solidity: function nativeCurrencyLabel() view returns(string)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorCaller) NativeCurrencyLabel(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NonfungibleTokenPositionDescriptor.contract.Call(opts, &out, "nativeCurrencyLabel")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NativeCurrencyLabel is a free data retrieval call binding the contract method 0xb7af3cdc.
//
// Solidity: function nativeCurrencyLabel() view returns(string)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorSession) NativeCurrencyLabel() (string, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.NativeCurrencyLabel(&_NonfungibleTokenPositionDescriptor.CallOpts)
}

// NativeCurrencyLabel is a free data retrieval call binding the contract method 0xb7af3cdc.
//
// Solidity: function nativeCurrencyLabel() view returns(string)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorCallerSession) NativeCurrencyLabel() (string, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.NativeCurrencyLabel(&_NonfungibleTokenPositionDescriptor.CallOpts)
}

// NativeCurrencyLabelBytes is a free data retrieval call binding the contract method 0xa18246e2.
//
// Solidity: function nativeCurrencyLabelBytes() view returns(bytes32)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorCaller) NativeCurrencyLabelBytes(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NonfungibleTokenPositionDescriptor.contract.Call(opts, &out, "nativeCurrencyLabelBytes")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// NativeCurrencyLabelBytes is a free data retrieval call binding the contract method 0xa18246e2.
//
// Solidity: function nativeCurrencyLabelBytes() view returns(bytes32)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorSession) NativeCurrencyLabelBytes() ([32]byte, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.NativeCurrencyLabelBytes(&_NonfungibleTokenPositionDescriptor.CallOpts)
}

// NativeCurrencyLabelBytes is a free data retrieval call binding the contract method 0xa18246e2.
//
// Solidity: function nativeCurrencyLabelBytes() view returns(bytes32)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorCallerSession) NativeCurrencyLabelBytes() ([32]byte, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.NativeCurrencyLabelBytes(&_NonfungibleTokenPositionDescriptor.CallOpts)
}

// TokenRatioPriority is a free data retrieval call binding the contract method 0x9d7b0ea8.
//
// Solidity: function tokenRatioPriority(address token, uint256 chainId) view returns(int256)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorCaller) TokenRatioPriority(opts *bind.CallOpts, token common.Address, chainId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _NonfungibleTokenPositionDescriptor.contract.Call(opts, &out, "tokenRatioPriority", token, chainId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenRatioPriority is a free data retrieval call binding the contract method 0x9d7b0ea8.
//
// Solidity: function tokenRatioPriority(address token, uint256 chainId) view returns(int256)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorSession) TokenRatioPriority(token common.Address, chainId *big.Int) (*big.Int, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.TokenRatioPriority(&_NonfungibleTokenPositionDescriptor.CallOpts, token, chainId)
}

// TokenRatioPriority is a free data retrieval call binding the contract method 0x9d7b0ea8.
//
// Solidity: function tokenRatioPriority(address token, uint256 chainId) view returns(int256)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorCallerSession) TokenRatioPriority(token common.Address, chainId *big.Int) (*big.Int, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.TokenRatioPriority(&_NonfungibleTokenPositionDescriptor.CallOpts, token, chainId)
}

// TokenURI is a free data retrieval call binding the contract method 0xe9dc6375.
//
// Solidity: function tokenURI(address positionManager, uint256 tokenId) view returns(string)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorCaller) TokenURI(opts *bind.CallOpts, positionManager common.Address, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _NonfungibleTokenPositionDescriptor.contract.Call(opts, &out, "tokenURI", positionManager, tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xe9dc6375.
//
// Solidity: function tokenURI(address positionManager, uint256 tokenId) view returns(string)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorSession) TokenURI(positionManager common.Address, tokenId *big.Int) (string, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.TokenURI(&_NonfungibleTokenPositionDescriptor.CallOpts, positionManager, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xe9dc6375.
//
// Solidity: function tokenURI(address positionManager, uint256 tokenId) view returns(string)
func (_NonfungibleTokenPositionDescriptor *NonfungibleTokenPositionDescriptorCallerSession) TokenURI(positionManager common.Address, tokenId *big.Int) (string, error) {
	return _NonfungibleTokenPositionDescriptor.Contract.TokenURI(&_NonfungibleTokenPositionDescriptor.CallOpts, positionManager, tokenId)
}

