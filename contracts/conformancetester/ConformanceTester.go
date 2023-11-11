// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package conformancetester

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
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"RevertErrorMessage\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testRevert\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162000883380380620008838339818101604052810190620000379190620001e3565b80600090816200004891906200047f565b505062000566565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b620000b9826200006e565b810181811067ffffffffffffffff82111715620000db57620000da6200007f565b5b80604052505050565b6000620000f062000050565b9050620000fe8282620000ae565b919050565b600067ffffffffffffffff8211156200012157620001206200007f565b5b6200012c826200006e565b9050602081019050919050565b60005b83811015620001595780820151818401526020810190506200013c565b60008484015250505050565b60006200017c620001768462000103565b620000e4565b9050828152602081018484840111156200019b576200019a62000069565b5b620001a884828562000139565b509392505050565b600082601f830112620001c857620001c762000064565b5b8151620001da84826020860162000165565b91505092915050565b600060208284031215620001fc57620001fb6200005a565b5b600082015167ffffffffffffffff8111156200021d576200021c6200005f565b5b6200022b84828501620001b0565b91505092915050565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806200028757607f821691505b6020821081036200029d576200029c6200023f565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302620003077fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82620002c8565b620003138683620002c8565b95508019841693508086168417925050509392505050565b6000819050919050565b6000819050919050565b6000620003606200035a62000354846200032b565b62000335565b6200032b565b9050919050565b6000819050919050565b6200037c836200033f565b620003946200038b8262000367565b848454620002d5565b825550505050565b600090565b620003ab6200039c565b620003b881848462000371565b505050565b5b81811015620003e057620003d4600082620003a1565b600181019050620003be565b5050565b601f8211156200042f57620003f981620002a3565b6200040484620002b8565b8101602085101562000414578190505b6200042c6200042385620002b8565b830182620003bd565b50505b505050565b600082821c905092915050565b6000620004546000198460080262000434565b1980831691505092915050565b60006200046f838362000441565b9150826002028217905092915050565b6200048a8262000234565b67ffffffffffffffff811115620004a657620004a56200007f565b5b620004b282546200026e565b620004bf828285620003e4565b600060209050601f831160018114620004f75760008415620004e2578287015190505b620004ee858262000461565b8655506200055e565b601f1984166200050786620002a3565b60005b8281101562000531578489015182556001820191506020850194506020810190506200050a565b868310156200055157848901516200054d601f89168262000441565b8355505b6001600288020188555050505b505050505050565b61030d80620005766000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806306fdde0314610046578063242e7fa114610064578063a26388bb14610082575b600080fd5b61004e61008c565b60405161005b9190610255565b60405180910390f35b61006c61011a565b6040516100799190610255565b60405180910390f35b61008a610153565b005b60008054610099906102a6565b80601f01602080910402602001604051908101604052809291908181526020018280546100c5906102a6565b80156101125780601f106100e757610100808354040283529160200191610112565b820191906000526020600020905b8154815290600101906020018083116100f557829003601f168201915b505050505081565b6040518060400160405280601981526020017f5465737420526576657274204572726f72204d6573736167650000000000000081525081565b6040518060400160405280601981526020017f5465737420526576657274204572726f72204d657373616765000000000000008152506040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101bc9190610255565b60405180910390fd5b600081519050919050565b600082825260208201905092915050565b60005b838110156101ff5780820151818401526020810190506101e4565b60008484015250505050565b6000601f19601f8301169050919050565b6000610227826101c5565b61023181856101d0565b93506102418185602086016101e1565b61024a8161020b565b840191505092915050565b6000602082019050818103600083015261026f818461021c565b905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806102be57607f821691505b6020821081036102d1576102d0610277565b5b5091905056fea26469706673582212202704472a01eddebdaeb571ea2f6c0e48fbd140c2982399ccaf6b350827fdfc4064736f6c63430008150033",
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
