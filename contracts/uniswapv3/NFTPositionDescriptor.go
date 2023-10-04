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

// NFTPositionDescriptorMetaData contains all meta data concerning the NFTPositionDescriptor contract.
var NFTPositionDescriptorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_WETH9\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_nativeCurrencyLabelBytes\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"WETH9\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"name\":\"flipRatio\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nativeCurrencyLabel\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nativeCurrencyLabelBytes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"name\":\"tokenRatioPriority\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractINonfungiblePositionManager\",\"name\":\"positionManager\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b506040516114f03803806114f083398101604081905261002f9161004a565b60609190911b6001600160601b03191660805260a052610082565b6000806040838503121561005c578182fd5b82516001600160a01b0381168114610072578283fd5b6020939093015192949293505050565b60805160601c60a0516114306100c06000398061027f52806102b3528061034f52508060f7528061013c52806105f2528061064652506114306000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80634aa4a4fc146100675780637e5af771146100855780639d7b0ea8146100a5578063a18246e2146100c5578063b7af3cdc146100cd578063e9dc6375146100e2575b600080fd5b61006f6100f5565b60405161007c9190611259565b60405180910390f35b610098610093366004610f5f565b610119565b60405161007c919061126d565b6100b86100b3366004610f9f565b610138565b60405161007c9190611278565b6100b861027d565b6100d56102a1565b60405161007c9190611281565b6100d56100f0366004610f9f565b6103af565b7f000000000000000000000000000000000000000000000000000000000000000081565b60006101258383610138565b61012f8584610138565b13949350505050565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316836001600160a01b0316141561017d5750606319610277565b8160011415610273576001600160a01b03831673a0b86991c6218b36c1d19d4a2e9eb0ce3606eb4814156101b4575061012c610277565b6001600160a01b03831673dac17f958d2ee523a2206206994597c13d831ec714156101e1575060c8610277565b6001600160a01b038316736b175474e89094c44da98b954eedeac495271d0f141561020e57506064610277565b6001600160a01b038316738daebade922df735c38c80c7ebd708af50815faa141561023c575060c719610277565b6001600160a01b038316732260fac5e5542a773aa44fbcfedf7c193bc2c599141561026b575061012b19610277565b506000610277565b5060005b92915050565b7f000000000000000000000000000000000000000000000000000000000000000081565b606060005b6020811080156102ee57507f000000000000000000000000000000000000000000000000000000000000000081602081106102dd57fe5b1a60f81b6001600160f81b03191615155b156102fb576001016102a6565b60008167ffffffffffffffff8111801561031457600080fd5b506040519080825280601f01601f19166020018201604052801561033f576020820181803683370190505b50905060005b828110156103a8577f0000000000000000000000000000000000000000000000000000000000000000816020811061037957fe5b1a60f81b82828151811061038957fe5b60200101906001600160f81b031916908160001a905350600101610345565b5091505090565b60606000806000806000876001600160a01b03166399fbab88886040518263ffffffff1660e01b81526004016103e59190611278565b6101806040518083038186803b1580156103fe57600080fd5b505afa158015610412573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104369190611124565b505050505096509650965096509650505060006104f4896001600160a01b031663c45a01556040518163ffffffff1660e01b815260040160206040518083038186803b15801561048557600080fd5b505afa158015610499573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104bd9190610f3c565b6040518060600160405280896001600160a01b03168152602001886001600160a01b031681526020018762ffffff168152506108bf565b9050600061050587876100936109a3565b9050600081156105155787610517565b865b9050600082156105275787610529565b885b90506000846001600160a01b0316633850c7bd6040518163ffffffff1660e01b815260040160e06040518083038186803b15801561056657600080fd5b505afa15801561057a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061059e919061107b565b505050505091505073f7012159bf761b312153e8c8d176932fe9aaa7ea63c49917d7604051806101c001604052808f8152602001866001600160a01b03168152602001856001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316876001600160a01b03161461063757610632876109a7565b61063f565b61063f6102a1565b81526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316866001600160a01b03161461068b57610686866109a7565b610693565b6106936102a1565b8152602001866001600160a01b031663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b1580156106d157600080fd5b505afa1580156106e5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610709919061110a565b60ff168152602001856001600160a01b031663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b15801561074a57600080fd5b505afa15801561075e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610782919061110a565b60ff16815260200187151581526020018a60020b81526020018960020b81526020018460020b8152602001886001600160a01b031663d0c93a7c6040518163ffffffff1660e01b815260040160206040518083038186803b1580156107e657600080fd5b505afa1580156107fa573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061081e9190610fca565b60020b81526020018b62ffffff168152602001886001600160a01b03168152506040518263ffffffff1660e01b815260040161085a9190611294565b60006040518083038186803b15801561087257600080fd5b505af4158015610886573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526108ae9190810190610fe4565b9d9c50505050505050505050505050565b600081602001516001600160a01b031682600001516001600160a01b0316106108e757600080fd5b50805160208083015160409384015184516001600160a01b0394851681850152939091168385015262ffffff166060808401919091528351808403820181526080840185528051908301206001600160f81b031960a085015294901b6bffffffffffffffffffffffff191660a183015260b58201939093527fe34f199b19b2b4f47f68442619d555527d244f78a3297ea89325f843f87b8b5460d5808301919091528251808303909101815260f5909101909152805191012090565b4690565b606060006109bc836395d89b4160e01b6109e1565b90508051600014156109d9576109d183610c09565b9150506109dc565b90505b919050565b60408051600481526024810182526020810180516001600160e01b03166001600160e01b031985161781529151815160609360009384936001600160a01b03891693919290918291908083835b60208310610a4d5780518252601f199092019160209182019101610a2e565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855afa9150503d8060008114610aad576040519150601f19603f3d011682016040523d82523d6000602084013e610ab2565b606091505b5091509150811580610ac357508051155b15610ae1576040518060200160405280600081525092505050610277565b805160201415610b19576000818060200190516020811015610b0257600080fd5b50519050610b0f81610c16565b9350505050610277565b604081511115610bf157808060200190516020811015610b3857600080fd5b8101908080516040519392919084640100000000821115610b5857600080fd5b908301906020820185811115610b6d57600080fd5b8251640100000000811182820188101715610b8757600080fd5b82525081516020918201929091019080838360005b83811015610bb4578181015183820152602001610b9c565b50505050905090810190601f168015610be15780820380516001836020036101000a031916815260200191505b5060405250505092505050610277565b50506040805160208101909152600081529392505050565b60606109d9826006610d3e565b604080516020808252818301909252606091600091906020820181803683370190505090506000805b6020811015610ca0576000858260208110610c5657fe5b1a60f81b90506001600160f81b0319811615610c975780848481518110610c7957fe5b60200101906001600160f81b031916908160001a9053506001909201915b50600101610c3f565b5060008167ffffffffffffffff81118015610cba57600080fd5b506040519080825280601f01601f191660200182016040528015610ce5576020820181803683370190505b50905060005b82811015610d3557838181518110610cff57fe5b602001015160f81c60f81b828281518110610d1657fe5b60200101906001600160f81b031916908160001a905350600101610ceb565b50949350505050565b606060028206158015610d515750600082115b8015610d5e575060288211155b610daf576040805162461bcd60e51b815260206004820152601e60248201527f41646472657373537472696e675574696c3a20494e56414c49445f4c454e0000604482015290519081900360640190fd5b60008267ffffffffffffffff81118015610dc857600080fd5b506040519080825280601f01601f191660200182016040528015610df3576020820181803683370190505b5090506001600160a01b03841660005b60028504811015610e9757600860138290030282901c600f600482901c1660f082168203610e3082610ea1565b868560020281518110610e3f57fe5b60200101906001600160f81b031916908160001a905350610e5f81610ea1565b868560020260010181518110610e7157fe5b60200101906001600160f81b031916908160001a9053505060019092019150610e039050565b5090949350505050565b6000600a8260ff161015610ebc57506030810160f81b6109dc565b506037810160f81b6109dc565b80516109dc816113e2565b8051600281900b81146109dc57600080fd5b80516fffffffffffffffffffffffffffffffff811681146109dc57600080fd5b805161ffff811681146109dc57600080fd5b805162ffffff811681146109dc57600080fd5b805160ff811681146109dc57600080fd5b600060208284031215610f4d578081fd5b8151610f58816113e2565b9392505050565b600080600060608486031215610f73578182fd5b8335610f7e816113e2565b92506020840135610f8e816113e2565b929592945050506040919091013590565b60008060408385031215610fb1578182fd5b8235610fbc816113e2565b946020939093013593505050565b600060208284031215610fdb578081fd5b610f5882610ed4565b600060208284031215610ff5578081fd5b815167ffffffffffffffff8082111561100c578283fd5b818401915084601f83011261101f578283fd5b81518181111561102b57fe5b604051601f8201601f19168101602001838111828210171561104957fe5b604052818152838201602001871015611060578485fd5b6110718260208301602087016113b2565b9695505050505050565b600080600080600080600060e0888a031215611095578283fd5b87516110a0816113e2565b96506110ae60208901610ed4565b95506110bc60408901610f06565b94506110ca60608901610f06565b93506110d860808901610f06565b92506110e660a08901610f2b565b915060c088015180151581146110fa578182fd5b8091505092959891949750929550565b60006020828403121561111b578081fd5b610f5882610f2b565b6000806000806000806000806000806000806101808d8f031215611146578485fd5b8c516bffffffffffffffffffffffff81168114611161578586fd5b9b5061116f60208e01610ec9565b9a5061117d60408e01610ec9565b995061118b60608e01610ec9565b985061119960808e01610f18565b97506111a760a08e01610ed4565b96506111b560c08e01610ed4565b95506111c360e08e01610ee6565b94506101008d015193506101208d015192506111e26101408e01610ee6565b91506111f16101608e01610ee6565b90509295989b509295989b509295989b565b6001600160a01b03169052565b15159052565b60020b9052565b600081518084526112358160208601602086016113b2565b601f01601f19169290920160200192915050565b62ffffff169052565b60ff169052565b6001600160a01b0391909116815260200190565b901515815260200190565b90815260200190565b600060208252610f58602083018461121d565b6000602082528251602083015260208301516112b36040840182611203565b5060408301516112c66060840182611203565b5060608301516101c08060808501526112e36101e085018361121d565b91506080850151601f198584030160a0860152611300838261121d565b92505060a085015161131560c0860182611252565b5060c085015161132860e0860182611252565b5060e085015161010061133d81870183611210565b860151905061012061135186820183611216565b860151905061014061136586820183611216565b860151905061016061137986820183611216565b860151905061018061138d86820183611216565b86015190506101a06113a186820183611249565b8601519050610e9785830182611203565b60005b838110156113cd5781810151838201526020016113b5565b838111156113dc576000848401525b50505050565b6001600160a01b03811681146113f757600080fd5b5056fea26469706673582212206c5b7c9e64dbe7d22211974fb6da184afe6c91cde0cddca02cb2f31a227c08db64736f6c63430007060033",
}

// NFTPositionDescriptorABI is the input ABI used to generate the binding from.
// Deprecated: Use NFTPositionDescriptorMetaData.ABI instead.
var NFTPositionDescriptorABI = NFTPositionDescriptorMetaData.ABI

// NFTPositionDescriptorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NFTPositionDescriptorMetaData.Bin instead.
var NFTPositionDescriptorBin = NFTPositionDescriptorMetaData.Bin

// DeployNFTPositionDescriptor deploys a new Ethereum contract, binding an instance of NFTPositionDescriptor to it.
func DeployNFTPositionDescriptor(auth *bind.TransactOpts, backend bind.ContractBackend, _WETH9 common.Address, _nativeCurrencyLabelBytes [32]byte, nonfungibleTokenPositionDescriptorNewBin string) (common.Address, *types.Transaction, *NFTPositionDescriptor, error) {
	parsed, err := NFTPositionDescriptorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(nonfungibleTokenPositionDescriptorNewBin), backend, _WETH9, _nativeCurrencyLabelBytes)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NFTPositionDescriptor{NFTPositionDescriptorCaller: NFTPositionDescriptorCaller{contract: contract}, NFTPositionDescriptorTransactor: NFTPositionDescriptorTransactor{contract: contract}, NFTPositionDescriptorFilterer: NFTPositionDescriptorFilterer{contract: contract}}, nil
}

// NFTPositionDescriptor is an auto generated Go binding around an Ethereum contract.
type NFTPositionDescriptor struct {
	NFTPositionDescriptorCaller     // Read-only binding to the contract
	NFTPositionDescriptorTransactor // Write-only binding to the contract
	NFTPositionDescriptorFilterer   // Log filterer for contract events
}

// NFTPositionDescriptorCaller is an auto generated read-only Go binding around an Ethereum contract.
type NFTPositionDescriptorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NFTPositionDescriptorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NFTPositionDescriptorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NFTPositionDescriptorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NFTPositionDescriptorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NFTPositionDescriptorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NFTPositionDescriptorSession struct {
	Contract     *NFTPositionDescriptor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// NFTPositionDescriptorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NFTPositionDescriptorCallerSession struct {
	Contract *NFTPositionDescriptorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// NFTPositionDescriptorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NFTPositionDescriptorTransactorSession struct {
	Contract     *NFTPositionDescriptorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// NFTPositionDescriptorRaw is an auto generated low-level Go binding around an Ethereum contract.
type NFTPositionDescriptorRaw struct {
	Contract *NFTPositionDescriptor // Generic contract binding to access the raw methods on
}

// NFTPositionDescriptorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NFTPositionDescriptorCallerRaw struct {
	Contract *NFTPositionDescriptorCaller // Generic read-only contract binding to access the raw methods on
}

// NFTPositionDescriptorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NFTPositionDescriptorTransactorRaw struct {
	Contract *NFTPositionDescriptorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNFTPositionDescriptor creates a new instance of NFTPositionDescriptor, bound to a specific deployed contract.
func NewNFTPositionDescriptor(address common.Address, backend bind.ContractBackend) (*NFTPositionDescriptor, error) {
	contract, err := bindNFTPositionDescriptor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NFTPositionDescriptor{NFTPositionDescriptorCaller: NFTPositionDescriptorCaller{contract: contract}, NFTPositionDescriptorTransactor: NFTPositionDescriptorTransactor{contract: contract}, NFTPositionDescriptorFilterer: NFTPositionDescriptorFilterer{contract: contract}}, nil
}

// NewNFTPositionDescriptorCaller creates a new read-only instance of NFTPositionDescriptor, bound to a specific deployed contract.
func NewNFTPositionDescriptorCaller(address common.Address, caller bind.ContractCaller) (*NFTPositionDescriptorCaller, error) {
	contract, err := bindNFTPositionDescriptor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NFTPositionDescriptorCaller{contract: contract}, nil
}

// NewNFTPositionDescriptorTransactor creates a new write-only instance of NFTPositionDescriptor, bound to a specific deployed contract.
func NewNFTPositionDescriptorTransactor(address common.Address, transactor bind.ContractTransactor) (*NFTPositionDescriptorTransactor, error) {
	contract, err := bindNFTPositionDescriptor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NFTPositionDescriptorTransactor{contract: contract}, nil
}

// NewNFTPositionDescriptorFilterer creates a new log filterer instance of NFTPositionDescriptor, bound to a specific deployed contract.
func NewNFTPositionDescriptorFilterer(address common.Address, filterer bind.ContractFilterer) (*NFTPositionDescriptorFilterer, error) {
	contract, err := bindNFTPositionDescriptor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NFTPositionDescriptorFilterer{contract: contract}, nil
}

// bindNFTPositionDescriptor binds a generic wrapper to an already deployed contract.
func bindNFTPositionDescriptor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NFTPositionDescriptorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NFTPositionDescriptor *NFTPositionDescriptorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NFTPositionDescriptor.Contract.NFTPositionDescriptorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NFTPositionDescriptor *NFTPositionDescriptorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NFTPositionDescriptor.Contract.NFTPositionDescriptorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NFTPositionDescriptor *NFTPositionDescriptorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NFTPositionDescriptor.Contract.NFTPositionDescriptorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NFTPositionDescriptor *NFTPositionDescriptorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NFTPositionDescriptor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NFTPositionDescriptor *NFTPositionDescriptorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NFTPositionDescriptor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NFTPositionDescriptor *NFTPositionDescriptorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NFTPositionDescriptor.Contract.contract.Transact(opts, method, params...)
}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_NFTPositionDescriptor *NFTPositionDescriptorCaller) WETH9(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NFTPositionDescriptor.contract.Call(opts, &out, "WETH9")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_NFTPositionDescriptor *NFTPositionDescriptorSession) WETH9() (common.Address, error) {
	return _NFTPositionDescriptor.Contract.WETH9(&_NFTPositionDescriptor.CallOpts)
}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_NFTPositionDescriptor *NFTPositionDescriptorCallerSession) WETH9() (common.Address, error) {
	return _NFTPositionDescriptor.Contract.WETH9(&_NFTPositionDescriptor.CallOpts)
}

// FlipRatio is a free data retrieval call binding the contract method 0x7e5af771.
//
// Solidity: function flipRatio(address token0, address token1, uint256 chainId) view returns(bool)
func (_NFTPositionDescriptor *NFTPositionDescriptorCaller) FlipRatio(opts *bind.CallOpts, token0 common.Address, token1 common.Address, chainId *big.Int) (bool, error) {
	var out []interface{}
	err := _NFTPositionDescriptor.contract.Call(opts, &out, "flipRatio", token0, token1, chainId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// FlipRatio is a free data retrieval call binding the contract method 0x7e5af771.
//
// Solidity: function flipRatio(address token0, address token1, uint256 chainId) view returns(bool)
func (_NFTPositionDescriptor *NFTPositionDescriptorSession) FlipRatio(token0 common.Address, token1 common.Address, chainId *big.Int) (bool, error) {
	return _NFTPositionDescriptor.Contract.FlipRatio(&_NFTPositionDescriptor.CallOpts, token0, token1, chainId)
}

// FlipRatio is a free data retrieval call binding the contract method 0x7e5af771.
//
// Solidity: function flipRatio(address token0, address token1, uint256 chainId) view returns(bool)
func (_NFTPositionDescriptor *NFTPositionDescriptorCallerSession) FlipRatio(token0 common.Address, token1 common.Address, chainId *big.Int) (bool, error) {
	return _NFTPositionDescriptor.Contract.FlipRatio(&_NFTPositionDescriptor.CallOpts, token0, token1, chainId)
}

// NativeCurrencyLabel is a free data retrieval call binding the contract method 0xb7af3cdc.
//
// Solidity: function nativeCurrencyLabel() view returns(string)
func (_NFTPositionDescriptor *NFTPositionDescriptorCaller) NativeCurrencyLabel(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NFTPositionDescriptor.contract.Call(opts, &out, "nativeCurrencyLabel")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NativeCurrencyLabel is a free data retrieval call binding the contract method 0xb7af3cdc.
//
// Solidity: function nativeCurrencyLabel() view returns(string)
func (_NFTPositionDescriptor *NFTPositionDescriptorSession) NativeCurrencyLabel() (string, error) {
	return _NFTPositionDescriptor.Contract.NativeCurrencyLabel(&_NFTPositionDescriptor.CallOpts)
}

// NativeCurrencyLabel is a free data retrieval call binding the contract method 0xb7af3cdc.
//
// Solidity: function nativeCurrencyLabel() view returns(string)
func (_NFTPositionDescriptor *NFTPositionDescriptorCallerSession) NativeCurrencyLabel() (string, error) {
	return _NFTPositionDescriptor.Contract.NativeCurrencyLabel(&_NFTPositionDescriptor.CallOpts)
}

// NativeCurrencyLabelBytes is a free data retrieval call binding the contract method 0xa18246e2.
//
// Solidity: function nativeCurrencyLabelBytes() view returns(bytes32)
func (_NFTPositionDescriptor *NFTPositionDescriptorCaller) NativeCurrencyLabelBytes(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NFTPositionDescriptor.contract.Call(opts, &out, "nativeCurrencyLabelBytes")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// NativeCurrencyLabelBytes is a free data retrieval call binding the contract method 0xa18246e2.
//
// Solidity: function nativeCurrencyLabelBytes() view returns(bytes32)
func (_NFTPositionDescriptor *NFTPositionDescriptorSession) NativeCurrencyLabelBytes() ([32]byte, error) {
	return _NFTPositionDescriptor.Contract.NativeCurrencyLabelBytes(&_NFTPositionDescriptor.CallOpts)
}

// NativeCurrencyLabelBytes is a free data retrieval call binding the contract method 0xa18246e2.
//
// Solidity: function nativeCurrencyLabelBytes() view returns(bytes32)
func (_NFTPositionDescriptor *NFTPositionDescriptorCallerSession) NativeCurrencyLabelBytes() ([32]byte, error) {
	return _NFTPositionDescriptor.Contract.NativeCurrencyLabelBytes(&_NFTPositionDescriptor.CallOpts)
}

// TokenRatioPriority is a free data retrieval call binding the contract method 0x9d7b0ea8.
//
// Solidity: function tokenRatioPriority(address token, uint256 chainId) view returns(int256)
func (_NFTPositionDescriptor *NFTPositionDescriptorCaller) TokenRatioPriority(opts *bind.CallOpts, token common.Address, chainId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _NFTPositionDescriptor.contract.Call(opts, &out, "tokenRatioPriority", token, chainId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenRatioPriority is a free data retrieval call binding the contract method 0x9d7b0ea8.
//
// Solidity: function tokenRatioPriority(address token, uint256 chainId) view returns(int256)
func (_NFTPositionDescriptor *NFTPositionDescriptorSession) TokenRatioPriority(token common.Address, chainId *big.Int) (*big.Int, error) {
	return _NFTPositionDescriptor.Contract.TokenRatioPriority(&_NFTPositionDescriptor.CallOpts, token, chainId)
}

// TokenRatioPriority is a free data retrieval call binding the contract method 0x9d7b0ea8.
//
// Solidity: function tokenRatioPriority(address token, uint256 chainId) view returns(int256)
func (_NFTPositionDescriptor *NFTPositionDescriptorCallerSession) TokenRatioPriority(token common.Address, chainId *big.Int) (*big.Int, error) {
	return _NFTPositionDescriptor.Contract.TokenRatioPriority(&_NFTPositionDescriptor.CallOpts, token, chainId)
}

// TokenURI is a free data retrieval call binding the contract method 0xe9dc6375.
//
// Solidity: function tokenURI(address positionManager, uint256 tokenId) view returns(string)
func (_NFTPositionDescriptor *NFTPositionDescriptorCaller) TokenURI(opts *bind.CallOpts, positionManager common.Address, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _NFTPositionDescriptor.contract.Call(opts, &out, "tokenURI", positionManager, tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xe9dc6375.
//
// Solidity: function tokenURI(address positionManager, uint256 tokenId) view returns(string)
func (_NFTPositionDescriptor *NFTPositionDescriptorSession) TokenURI(positionManager common.Address, tokenId *big.Int) (string, error) {
	return _NFTPositionDescriptor.Contract.TokenURI(&_NFTPositionDescriptor.CallOpts, positionManager, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xe9dc6375.
//
// Solidity: function tokenURI(address positionManager, uint256 tokenId) view returns(string)
func (_NFTPositionDescriptor *NFTPositionDescriptorCallerSession) TokenURI(positionManager common.Address, tokenId *big.Int) (string, error) {
	return _NFTPositionDescriptor.Contract.TokenURI(&_NFTPositionDescriptor.CallOpts, positionManager, tokenId)
}
