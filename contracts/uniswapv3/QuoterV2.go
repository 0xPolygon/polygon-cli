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

// IQuoterV2QuoteExactInputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type IQuoterV2QuoteExactInputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	AmountIn          *big.Int
	Fee               *big.Int
	SqrtPriceLimitX96 *big.Int
}

// IQuoterV2QuoteExactOutputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type IQuoterV2QuoteExactOutputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	Amount            *big.Int
	Fee               *big.Int
	SqrtPriceLimitX96 *big.Int
}

// QuoterV2MetaData contains all meta data concerning the QuoterV2 contract.
var QuoterV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_WETH9\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"WETH9\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"quoteExactInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint160[]\",\"name\":\"sqrtPriceX96AfterList\",\"type\":\"uint160[]\"},{\"internalType\":\"uint32[]\",\"name\":\"initializedTicksCrossedList\",\"type\":\"uint32[]\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structIQuoterV2.QuoteExactInputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"quoteExactInputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96After\",\"type\":\"uint160\"},{\"internalType\":\"uint32\",\"name\":\"initializedTicksCrossed\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"quoteExactOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint160[]\",\"name\":\"sqrtPriceX96AfterList\",\"type\":\"uint160[]\"},{\"internalType\":\"uint32[]\",\"name\":\"initializedTicksCrossedList\",\"type\":\"uint32[]\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structIQuoterV2.QuoteExactOutputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"quoteExactOutputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96After\",\"type\":\"uint160\"},{\"internalType\":\"uint32\",\"name\":\"initializedTicksCrossed\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"gasEstimate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"amount0Delta\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"amount1Delta\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"}],\"name\":\"uniswapV3SwapCallback\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162002ad238038062002ad28339818101604052810190620000379190620000c8565b81818173ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff1660601b815250508073ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff1660601b815250505050505062000157565b600081519050620000c2816200013d565b92915050565b60008060408385031215620000dc57600080fd5b6000620000ec85828601620000b1565b9250506020620000ff85828601620000b1565b9150509250929050565b600062000116826200011d565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b620001488162000109565b81146200015457600080fd5b50565b60805160601c60a05160601c6129476200018b600039806103b75250806106185280610a5d5280610ca152506129476000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063c45a01551161005b578063c45a015514610106578063c6a5026a14610124578063cdca175314610157578063fa461e331461018a5761007d565b80632f80bb1d146100825780634aa4a4fc146100b5578063bd21704a146100d3575b600080fd5b61009c60048036038101906100979190611f82565b6101a6565b6040516100ac9493929190612531565b60405180910390f35b6100bd6103b5565b6040516100ca919061247a565b60405180910390f35b6100ed60048036038101906100e891906120e3565b6103d9565b6040516100fd9493929190612584565b60405180910390f35b61010e610616565b60405161011b919061247a565b60405180910390f35b61013e600480360381019061013991906120ba565b61063a565b60405161014e9493929190612584565b60405180910390f35b610171600480360381019061016c9190611f82565b61081d565b6040516101819493929190612531565b60405180910390f35b6101a4600480360381019061019f9190612012565b610a2c565b005b600060608060006101b686610bea565b67ffffffffffffffff811180156101cc57600080fd5b506040519080825280602002602001820160405280156101fb5781602001602082028036833780820191505090505b50925061020786610bea565b67ffffffffffffffff8111801561021d57600080fd5b5060405190808252806020026020018201604052801561024c5781602001602082028036833780820191505090505b50915060005b6001156103aa5760008060006102678a610c05565b9250925092506000806000806102ea6040518060a001604052808873ffffffffffffffffffffffffffffffffffffffff1681526020018973ffffffffffffffffffffffffffffffffffffffff1681526020018f81526020018762ffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff168152506103d9565b9350935093509350828b89815181106102ff57fe5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050818a898151811061034657fe5b602002602001019063ffffffff16908163ffffffff1681525050839c50808901985087806001019850506103798e610c56565b1561038e576103878e610c71565b9d5061039e565b8c9b5050505050505050506103ac565b50505050505050610252565b505b92959194509250565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000806000806000856020015173ffffffffffffffffffffffffffffffffffffffff16866000015173ffffffffffffffffffffffffffffffffffffffff161090506000610433876000015188602001518960600151610c9a565b90506000876080015173ffffffffffffffffffffffffffffffffffffffff1614156104645786604001516000819055505b60005a90508173ffffffffffffffffffffffffffffffffffffffff1663128acb0830856104948c60400151610cd9565b60000360008d6080015173ffffffffffffffffffffffffffffffffffffffff16146104c3578c608001516104f0565b876104e557600173fffd8963efd1fc6a506488495d951d5263988d26036104ef565b60016401000276a3015b5b8d602001518e606001518f600001516040516020016105119392919061243d565b6040516020818303038152906040526040518663ffffffff1660e01b8152600401610540959493929190612495565b6040805180830381600087803b15801561055957600080fd5b505af192505050801561058a57506040513d601f19601f820116820180604052508101906105879190611fd6565b60015b610609573d80600081146105ba576040519150601f19603f3d011682016040523d82523d6000602084013e6105bf565b606091505b505a820394506000896080015173ffffffffffffffffffffffffffffffffffffffff1614156105ed57600080555b6105f8818487610d0f565b97509750975097505050505061060f565b50505050505b9193509193565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000806000806000856020015173ffffffffffffffffffffffffffffffffffffffff16866000015173ffffffffffffffffffffffffffffffffffffffff161090506000610694876000015188602001518960600151610c9a565b905060005a90508173ffffffffffffffffffffffffffffffffffffffff1663128acb0830856106c68c60400151610cd9565b60008d6080015173ffffffffffffffffffffffffffffffffffffffff16146106f2578c6080015161071f565b8761071457600173fffd8963efd1fc6a506488495d951d5263988d260361071e565b60016401000276a3015b5b8d600001518e606001518f602001516040516020016107409392919061243d565b6040516020818303038152906040526040518663ffffffff1660e01b815260040161076f959493929190612495565b6040805180830381600087803b15801561078857600080fd5b505af19250505080156107b957506040513d601f19601f820116820180604052508101906107b69190611fd6565b60015b610810573d80600081146107e9576040519150601f19603f3d011682016040523d82523d6000602084013e6107ee565b606091505b505a820394506107ff818487610d0f565b975097509750975050505050610816565b50505050505b9193509193565b6000606080600061082d86610bea565b67ffffffffffffffff8111801561084357600080fd5b506040519080825280602002602001820160405280156108725781602001602082028036833780820191505090505b50925061087e86610bea565b67ffffffffffffffff8111801561089457600080fd5b506040519080825280602002602001820160405280156108c35781602001602082028036833780820191505090505b50915060005b600115610a215760008060006108de8a610c05565b9250925092506000806000806109616040518060a001604052808973ffffffffffffffffffffffffffffffffffffffff1681526020018873ffffffffffffffffffffffffffffffffffffffff1681526020018f81526020018762ffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff1681525061063a565b9350935093509350828b898151811061097657fe5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050818a89815181106109bd57fe5b602002602001019063ffffffff16908163ffffffff1681525050839c50808901985087806001019850506109f08e610c56565b15610a05576109fe8e610c71565b9d50610a15565b8c9b505050505050505050610a23565b505050505050506108c9565b505b92959194509250565b6000831380610a3b5750600082135b610a4457600080fd5b6000806000610a5284610c05565b925092509250610a847f0000000000000000000000000000000000000000000000000000000000000000848484610e09565b506000806000808913610aca578573ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1610888a600003610aff565b8473ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff161089896000035b9250925092506000610b12878787610c9a565b90506000808273ffffffffffffffffffffffffffffffffffffffff16633850c7bd6040518163ffffffff1660e01b815260040160e06040518083038186803b158015610b5d57600080fd5b505afa158015610b71573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b95919061210c565b5050505050915091508515610bbb57604051848152826020820152816040820152606081fd5b6000805414610bd3576000548414610bd257600080fd5b5b604051858152826020820152816040820152606081fd5b60006003601401601483510381610bfd57fe5b049050919050565b6000806000610c1e600085610e2990919063ffffffff16565b9250610c34601485610f4290919063ffffffff16565b9050610c4d600360140185610e2990919063ffffffff16565b91509193909250565b60006003601401601460036014010101825110159050919050565b6060610c93600360140160036014018451038461104c9092919063ffffffff16565b9050919050565b6000610cd07f0000000000000000000000000000000000000000000000000000000000000000610ccb868686611236565b6112d2565b90509392505050565b60007f80000000000000000000000000000000000000000000000000000000000000008210610d0757600080fd5b819050919050565b6000806000806000808773ffffffffffffffffffffffffffffffffffffffff16633850c7bd6040518163ffffffff1660e01b815260040160e06040518083038186803b158015610d5e57600080fd5b505afa158015610d72573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d96919061210c565b9091929394955090919293509091925090915090505080925050610db98961142d565b809350819750829850505050610df082828a73ffffffffffffffffffffffffffffffffffffffff166114f79092919063ffffffff16565b9350858585899550955095509550505093509350935093565b6000610e1f85610e1a868686611236565b611bda565b9050949350505050565b600081601483011015610ea4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f746f416464726573735f6f766572666c6f77000000000000000000000000000081525060200191505060405180910390fd5b6014820183511015610f1e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f746f416464726573735f6f75744f66426f756e6473000000000000000000000081525060200191505060405180910390fd5b60006c01000000000000000000000000836020860101510490508091505092915050565b600081600383011015610fbd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f746f55696e7432345f6f766572666c6f7700000000000000000000000000000081525060200191505060405180910390fd5b6003820183511015611037576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f746f55696e7432345f6f75744f66426f756e647300000000000000000000000081525060200191505060405180910390fd5b60008260038501015190508091505092915050565b606081601f830110156110c7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600e8152602001807f736c6963655f6f766572666c6f7700000000000000000000000000000000000081525060200191505060405180910390fd5b82828401101561113f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600e8152602001807f736c6963655f6f766572666c6f7700000000000000000000000000000000000081525060200191505060405180910390fd5b818301845110156111b8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f736c6963655f6f75744f66426f756e647300000000000000000000000000000081525060200191505060405180910390fd5b60608215600081146111d9576040519150600082526020820160405261122a565b6040519150601f8416801560200281840101858101878315602002848b0101015b8183101561121757805183526020830192506020810190506111fa565b50868552601f19601f8301166040525050505b50809150509392505050565b61123e611c54565b8273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16111561127d57828480945081955050505b60405180606001604052808573ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1681526020018362ffffff1681525090509392505050565b6000816020015173ffffffffffffffffffffffffffffffffffffffff16826000015173ffffffffffffffffffffffffffffffffffffffff161061131457600080fd5b82826000015183602001518460400151604051602001808473ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1681526020018262ffffff1681526020019350505050604051602081830303815290604052805190602001207fe34f199b19b2b4f47f68442619d555527d244f78a3297ea89325f843f87b8b5460001b60405160200180807fff000000000000000000000000000000000000000000000000000000000000008152506001018473ffffffffffffffffffffffffffffffffffffffff1660601b815260140183815260200182815260200193505050506040516020818303038152906040528051906020012060001c905092915050565b600080600060608451146114d657604484511015611480576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161147790612511565b60405180910390fd5b6004840193508380602001905181019061149a9190612079565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016114cd91906124ef565b60405180910390fd5b838060200190518101906114ea91906121aa565b9250925092509193909250565b60008060008060008060008060088b73ffffffffffffffffffffffffffffffffffffffff1663d0c93a7c6040518163ffffffff1660e01b815260040160206040518083038186803b15801561154b57600080fd5b505afa15801561155f573d6000803e3d6000fd5b505050506040513d602081101561157557600080fd5b810190808051906020019092919050505060020b8b60020b8161159457fe5b0560020b901d905060006101008c73ffffffffffffffffffffffffffffffffffffffff1663d0c93a7c6040518163ffffffff1660e01b815260040160206040518083038186803b1580156115e757600080fd5b505afa1580156115fb573d6000803e3d6000fd5b505050506040513d602081101561161157600080fd5b810190808051906020019092919050505060020b8c60020b8161163057fe5b0560020b8161163b57fe5b079050600060088d73ffffffffffffffffffffffffffffffffffffffff1663d0c93a7c6040518163ffffffff1660e01b815260040160206040518083038186803b15801561168857600080fd5b505afa15801561169c573d6000803e3d6000fd5b505050506040513d60208110156116b257600080fd5b810190808051906020019092919050505060020b8c60020b816116d157fe5b0560020b901d905060006101008e73ffffffffffffffffffffffffffffffffffffffff1663d0c93a7c6040518163ffffffff1660e01b815260040160206040518083038186803b15801561172457600080fd5b505afa158015611738573d6000803e3d6000fd5b505050506040513d602081101561174e57600080fd5b810190808051906020019092919050505060020b8d60020b8161176d57fe5b0560020b8161177857fe5b07905060008160ff166001901b8f73ffffffffffffffffffffffffffffffffffffffff16635339c296856040518263ffffffff1660e01b8152600401808260010b815260200191505060206040518083038186803b1580156117d957600080fd5b505afa1580156117ed573d6000803e3d6000fd5b505050506040513d602081101561180357600080fd5b8101908080519060200190929190505050161180156118b4575060008e73ffffffffffffffffffffffffffffffffffffffff1663d0c93a7c6040518163ffffffff1660e01b815260040160206040518083038186803b15801561186557600080fd5b505afa158015611879573d6000803e3d6000fd5b505050506040513d602081101561188f57600080fd5b810190808051906020019092919050505060020b8d60020b816118ae57fe5b0760020b145b80156118c557508b60020b8d60020b135b945060008360ff166001901b8f73ffffffffffffffffffffffffffffffffffffffff16635339c296876040518263ffffffff1660e01b8152600401808260010b815260200191505060206040518083038186803b15801561192557600080fd5b505afa158015611939573d6000803e3d6000fd5b505050506040513d602081101561194f57600080fd5b810190808051906020019092919050505016118015611a00575060008e73ffffffffffffffffffffffffffffffffffffffff1663d0c93a7c6040518163ffffffff1660e01b815260040160206040518083038186803b1580156119b157600080fd5b505afa1580156119c5573d6000803e3d6000fd5b505050506040513d60208110156119db57600080fd5b810190808051906020019092919050505060020b8e60020b816119fa57fe5b0760020b145b8015611a1157508b60020b8d60020b125b95508160010b8460010b1280611a3e57508160010b8460010b148015611a3d57508060ff168360ff1611155b5b15611a5457839950829750819850809650611a61565b8199508097508398508296505b5050505060008460ff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff901b90505b8560010b8760010b13611bb2578560010b8760010b1415611adb578360ff0360ff167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff901c811690505b6000818c73ffffffffffffffffffffffffffffffffffffffff16635339c2968a6040518263ffffffff1660e01b8152600401808260010b815260200191505060206040518083038186803b158015611b3257600080fd5b505afa158015611b46573d6000803e3d6000fd5b505050506040513d6020811015611b5c57600080fd5b8101908080519060200190929190505050169050611b7981611c26565b61ffff168901985087806001019850507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff915050611a91565b8115611bbf576001880397505b8215611bcc576001880397505b505050505050509392505050565b6000611be683836112d2565b90508073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611c2057600080fd5b92915050565b600080600090505b60008314611c4b5780806001019150506001830383169250611c2e565b80915050919050565b6040518060600160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600062ffffff1681525090565b6000611cb9611cb4846125fa565b6125c9565b905082815260208101848484011115611cd157600080fd5b611cdc84828561279d565b509392505050565b6000611cf7611cf28461262a565b6125c9565b905082815260208101848484011115611d0f57600080fd5b611d1a8482856127ac565b509392505050565b600081359050611d3181612842565b92915050565b600081519050611d4681612859565b92915050565b600082601f830112611d5d57600080fd5b8135611d6d848260208601611ca6565b91505092915050565b600081519050611d8581612870565b92915050565b600081359050611d9a81612887565b92915050565b600081519050611daf81612887565b92915050565b600082601f830112611dc657600080fd5b8151611dd6848260208601611ce4565b91505092915050565b600060a08284031215611df157600080fd5b611dfb60a06125c9565b90506000611e0b84828501611d22565b6000830152506020611e1f84828501611d22565b6020830152506040611e3384828501611f43565b6040830152506060611e4784828501611f2e565b6060830152506080611e5b84828501611eef565b60808301525092915050565b600060a08284031215611e7957600080fd5b611e8360a06125c9565b90506000611e9384828501611d22565b6000830152506020611ea784828501611d22565b6020830152506040611ebb84828501611f43565b6040830152506060611ecf84828501611f2e565b6060830152506080611ee384828501611eef565b60808301525092915050565b600081359050611efe816128b5565b92915050565b600081519050611f13816128b5565b92915050565b600081519050611f288161289e565b92915050565b600081359050611f3d816128cc565b92915050565b600081359050611f52816128e3565b92915050565b600081519050611f67816128e3565b92915050565b600081519050611f7c816128fa565b92915050565b60008060408385031215611f9557600080fd5b600083013567ffffffffffffffff811115611faf57600080fd5b611fbb85828601611d4c565b9250506020611fcc85828601611f43565b9150509250929050565b60008060408385031215611fe957600080fd5b6000611ff785828601611da0565b925050602061200885828601611da0565b9150509250929050565b60008060006060848603121561202757600080fd5b600061203586828701611d8b565b935050602061204686828701611d8b565b925050604084013567ffffffffffffffff81111561206357600080fd5b61206f86828701611d4c565b9150509250925092565b60006020828403121561208b57600080fd5b600082015167ffffffffffffffff8111156120a557600080fd5b6120b184828501611db5565b91505092915050565b600060a082840312156120cc57600080fd5b60006120da84828501611ddf565b91505092915050565b600060a082840312156120f557600080fd5b600061210384828501611e67565b91505092915050565b600080600080600080600060e0888a03121561212757600080fd5b60006121358a828b01611f04565b97505060206121468a828b01611d76565b96505060406121578a828b01611f19565b95505060606121688a828b01611f19565b94505060806121798a828b01611f19565b93505060a061218a8a828b01611f6d565b92505060c061219b8a828b01611d37565b91505092959891949750929550565b6000806000606084860312156121bf57600080fd5b60006121cd86828701611f58565b93505060206121de86828701611f04565b92505060406121ef86828701611d76565b9150509250925092565b600061220583836123db565b60208301905092915050565b600061221d838361241f565b60208301905092915050565b61223281612704565b82525050565b61224961224482612704565b6127df565b82525050565b600061225a8261267a565b61226481856126c0565b935061226f8361265a565b8060005b838110156122a057815161228788826121f9565b9750612292836126a6565b925050600181019050612273565b5085935050505092915050565b60006122b882612685565b6122c281856126d1565b93506122cd8361266a565b8060005b838110156122fe5781516122e58882612211565b97506122f0836126b3565b9250506001810190506122d1565b5085935050505092915050565b61231481612716565b82525050565b600061232582612690565b61232f81856126e2565b935061233f8185602086016127ac565b61234881612817565b840191505092915050565b61235c8161272f565b82525050565b600061236d8261269b565b61237781856126f3565b93506123878185602086016127ac565b61239081612817565b840191505092915050565b60006123a86010836126f3565b91507f556e6578706563746564206572726f72000000000000000000000000000000006000830152602082019050919050565b6123e481612747565b82525050565b6123f381612747565b82525050565b61240a61240582612767565b612803565b82525050565b61241981612776565b82525050565b61242881612780565b82525050565b61243781612780565b82525050565b60006124498286612238565b60148201915061245982856123f9565b6003820191506124698284612238565b601482019150819050949350505050565b600060208201905061248f6000830184612229565b92915050565b600060a0820190506124aa6000830188612229565b6124b7602083018761230b565b6124c46040830186612353565b6124d160608301856123ea565b81810360808301526124e3818461231a565b90509695505050505050565b600060208201905081810360008301526125098184612362565b905092915050565b6000602082019050818103600083015261252a8161239b565b9050919050565b60006080820190506125466000830187612410565b8181036020830152612558818661224f565b9050818103604083015261256c81856122ad565b905061257b6060830184612410565b95945050505050565b60006080820190506125996000830187612410565b6125a660208301866123ea565b6125b3604083018561242e565b6125c06060830184612410565b95945050505050565b6000604051905081810181811067ffffffffffffffff821117156125f0576125ef612815565b5b8060405250919050565b600067ffffffffffffffff82111561261557612614612815565b5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff82111561264557612644612815565b5b601f19601f8301169050602081019050919050565b6000819050602082019050919050565b6000819050602082019050919050565b600081519050919050565b600081519050919050565b600081519050919050565b600081519050919050565b6000602082019050919050565b6000602082019050919050565b600082825260208201905092915050565b600082825260208201905092915050565b600082825260208201905092915050565b600082825260208201905092915050565b600061270f82612747565b9050919050565b60008115159050919050565b60008160020b9050919050565b6000819050919050565b600061ffff82169050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600062ffffff82169050919050565b6000819050919050565b600063ffffffff82169050919050565b600060ff82169050919050565b82818337600083830152505050565b60005b838110156127ca5780820151818401526020810190506127af565b838111156127d9576000848401525b50505050565b60006127ea826127f1565b9050919050565b60006127fc82612835565b9050919050565b600061280e82612828565b9050919050565bfe5b6000601f19601f8301169050919050565b60008160e81b9050919050565b60008160601b9050919050565b61284b81612704565b811461285657600080fd5b50565b61286281612716565b811461286d57600080fd5b50565b61287981612722565b811461288457600080fd5b50565b6128908161272f565b811461289b57600080fd5b50565b6128a781612739565b81146128b257600080fd5b50565b6128be81612747565b81146128c957600080fd5b50565b6128d581612767565b81146128e057600080fd5b50565b6128ec81612776565b81146128f757600080fd5b50565b61290381612790565b811461290e57600080fd5b5056fea26469706673582212206ebbc2af0790c7997ec6949499a8244896db475d0086966828e1c3097214598664736f6c63430007060033",
}

// QuoterV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use QuoterV2MetaData.ABI instead.
var QuoterV2ABI = QuoterV2MetaData.ABI

// QuoterV2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use QuoterV2MetaData.Bin instead.
var QuoterV2Bin = QuoterV2MetaData.Bin

// DeployQuoterV2 deploys a new Ethereum contract, binding an instance of QuoterV2 to it.
func DeployQuoterV2(auth *bind.TransactOpts, backend bind.ContractBackend, _factory common.Address, _WETH9 common.Address) (common.Address, *types.Transaction, *QuoterV2, error) {
	parsed, err := QuoterV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(QuoterV2Bin), backend, _factory, _WETH9)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &QuoterV2{QuoterV2Caller: QuoterV2Caller{contract: contract}, QuoterV2Transactor: QuoterV2Transactor{contract: contract}, QuoterV2Filterer: QuoterV2Filterer{contract: contract}}, nil
}

// QuoterV2 is an auto generated Go binding around an Ethereum contract.
type QuoterV2 struct {
	QuoterV2Caller     // Read-only binding to the contract
	QuoterV2Transactor // Write-only binding to the contract
	QuoterV2Filterer   // Log filterer for contract events
}

// QuoterV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type QuoterV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuoterV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type QuoterV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuoterV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type QuoterV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuoterV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type QuoterV2Session struct {
	Contract     *QuoterV2         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// QuoterV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type QuoterV2CallerSession struct {
	Contract *QuoterV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// QuoterV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type QuoterV2TransactorSession struct {
	Contract     *QuoterV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// QuoterV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type QuoterV2Raw struct {
	Contract *QuoterV2 // Generic contract binding to access the raw methods on
}

// QuoterV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type QuoterV2CallerRaw struct {
	Contract *QuoterV2Caller // Generic read-only contract binding to access the raw methods on
}

// QuoterV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type QuoterV2TransactorRaw struct {
	Contract *QuoterV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewQuoterV2 creates a new instance of QuoterV2, bound to a specific deployed contract.
func NewQuoterV2(address common.Address, backend bind.ContractBackend) (*QuoterV2, error) {
	contract, err := bindQuoterV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &QuoterV2{QuoterV2Caller: QuoterV2Caller{contract: contract}, QuoterV2Transactor: QuoterV2Transactor{contract: contract}, QuoterV2Filterer: QuoterV2Filterer{contract: contract}}, nil
}

// NewQuoterV2Caller creates a new read-only instance of QuoterV2, bound to a specific deployed contract.
func NewQuoterV2Caller(address common.Address, caller bind.ContractCaller) (*QuoterV2Caller, error) {
	contract, err := bindQuoterV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &QuoterV2Caller{contract: contract}, nil
}

// NewQuoterV2Transactor creates a new write-only instance of QuoterV2, bound to a specific deployed contract.
func NewQuoterV2Transactor(address common.Address, transactor bind.ContractTransactor) (*QuoterV2Transactor, error) {
	contract, err := bindQuoterV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &QuoterV2Transactor{contract: contract}, nil
}

// NewQuoterV2Filterer creates a new log filterer instance of QuoterV2, bound to a specific deployed contract.
func NewQuoterV2Filterer(address common.Address, filterer bind.ContractFilterer) (*QuoterV2Filterer, error) {
	contract, err := bindQuoterV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &QuoterV2Filterer{contract: contract}, nil
}

// bindQuoterV2 binds a generic wrapper to an already deployed contract.
func bindQuoterV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := QuoterV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_QuoterV2 *QuoterV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _QuoterV2.Contract.QuoterV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_QuoterV2 *QuoterV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QuoterV2.Contract.QuoterV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_QuoterV2 *QuoterV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _QuoterV2.Contract.QuoterV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_QuoterV2 *QuoterV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _QuoterV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_QuoterV2 *QuoterV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _QuoterV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_QuoterV2 *QuoterV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _QuoterV2.Contract.contract.Transact(opts, method, params...)
}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_QuoterV2 *QuoterV2Caller) WETH9(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _QuoterV2.contract.Call(opts, &out, "WETH9")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_QuoterV2 *QuoterV2Session) WETH9() (common.Address, error) {
	return _QuoterV2.Contract.WETH9(&_QuoterV2.CallOpts)
}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_QuoterV2 *QuoterV2CallerSession) WETH9() (common.Address, error) {
	return _QuoterV2.Contract.WETH9(&_QuoterV2.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_QuoterV2 *QuoterV2Caller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _QuoterV2.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_QuoterV2 *QuoterV2Session) Factory() (common.Address, error) {
	return _QuoterV2.Contract.Factory(&_QuoterV2.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_QuoterV2 *QuoterV2CallerSession) Factory() (common.Address, error) {
	return _QuoterV2.Contract.Factory(&_QuoterV2.CallOpts)
}

// UniswapV3SwapCallback is a free data retrieval call binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes path) view returns()
func (_QuoterV2 *QuoterV2Caller) UniswapV3SwapCallback(opts *bind.CallOpts, amount0Delta *big.Int, amount1Delta *big.Int, path []byte) error {
	var out []interface{}
	err := _QuoterV2.contract.Call(opts, &out, "uniswapV3SwapCallback", amount0Delta, amount1Delta, path)

	if err != nil {
		return err
	}

	return err

}

// UniswapV3SwapCallback is a free data retrieval call binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes path) view returns()
func (_QuoterV2 *QuoterV2Session) UniswapV3SwapCallback(amount0Delta *big.Int, amount1Delta *big.Int, path []byte) error {
	return _QuoterV2.Contract.UniswapV3SwapCallback(&_QuoterV2.CallOpts, amount0Delta, amount1Delta, path)
}

// UniswapV3SwapCallback is a free data retrieval call binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes path) view returns()
func (_QuoterV2 *QuoterV2CallerSession) UniswapV3SwapCallback(amount0Delta *big.Int, amount1Delta *big.Int, path []byte) error {
	return _QuoterV2.Contract.UniswapV3SwapCallback(&_QuoterV2.CallOpts, amount0Delta, amount1Delta, path)
}

// QuoteExactInput is a paid mutator transaction binding the contract method 0xcdca1753.
//
// Solidity: function quoteExactInput(bytes path, uint256 amountIn) returns(uint256 amountOut, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_QuoterV2 *QuoterV2Transactor) QuoteExactInput(opts *bind.TransactOpts, path []byte, amountIn *big.Int) (*types.Transaction, error) {
	return _QuoterV2.contract.Transact(opts, "quoteExactInput", path, amountIn)
}

// QuoteExactInput is a paid mutator transaction binding the contract method 0xcdca1753.
//
// Solidity: function quoteExactInput(bytes path, uint256 amountIn) returns(uint256 amountOut, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_QuoterV2 *QuoterV2Session) QuoteExactInput(path []byte, amountIn *big.Int) (*types.Transaction, error) {
	return _QuoterV2.Contract.QuoteExactInput(&_QuoterV2.TransactOpts, path, amountIn)
}

// QuoteExactInput is a paid mutator transaction binding the contract method 0xcdca1753.
//
// Solidity: function quoteExactInput(bytes path, uint256 amountIn) returns(uint256 amountOut, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_QuoterV2 *QuoterV2TransactorSession) QuoteExactInput(path []byte, amountIn *big.Int) (*types.Transaction, error) {
	return _QuoterV2.Contract.QuoteExactInput(&_QuoterV2.TransactOpts, path, amountIn)
}

// QuoteExactInputSingle is a paid mutator transaction binding the contract method 0xc6a5026a.
//
// Solidity: function quoteExactInputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountOut, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV2 *QuoterV2Transactor) QuoteExactInputSingle(opts *bind.TransactOpts, params IQuoterV2QuoteExactInputSingleParams) (*types.Transaction, error) {
	return _QuoterV2.contract.Transact(opts, "quoteExactInputSingle", params)
}

// QuoteExactInputSingle is a paid mutator transaction binding the contract method 0xc6a5026a.
//
// Solidity: function quoteExactInputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountOut, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV2 *QuoterV2Session) QuoteExactInputSingle(params IQuoterV2QuoteExactInputSingleParams) (*types.Transaction, error) {
	return _QuoterV2.Contract.QuoteExactInputSingle(&_QuoterV2.TransactOpts, params)
}

// QuoteExactInputSingle is a paid mutator transaction binding the contract method 0xc6a5026a.
//
// Solidity: function quoteExactInputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountOut, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV2 *QuoterV2TransactorSession) QuoteExactInputSingle(params IQuoterV2QuoteExactInputSingleParams) (*types.Transaction, error) {
	return _QuoterV2.Contract.QuoteExactInputSingle(&_QuoterV2.TransactOpts, params)
}

// QuoteExactOutput is a paid mutator transaction binding the contract method 0x2f80bb1d.
//
// Solidity: function quoteExactOutput(bytes path, uint256 amountOut) returns(uint256 amountIn, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_QuoterV2 *QuoterV2Transactor) QuoteExactOutput(opts *bind.TransactOpts, path []byte, amountOut *big.Int) (*types.Transaction, error) {
	return _QuoterV2.contract.Transact(opts, "quoteExactOutput", path, amountOut)
}

// QuoteExactOutput is a paid mutator transaction binding the contract method 0x2f80bb1d.
//
// Solidity: function quoteExactOutput(bytes path, uint256 amountOut) returns(uint256 amountIn, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_QuoterV2 *QuoterV2Session) QuoteExactOutput(path []byte, amountOut *big.Int) (*types.Transaction, error) {
	return _QuoterV2.Contract.QuoteExactOutput(&_QuoterV2.TransactOpts, path, amountOut)
}

// QuoteExactOutput is a paid mutator transaction binding the contract method 0x2f80bb1d.
//
// Solidity: function quoteExactOutput(bytes path, uint256 amountOut) returns(uint256 amountIn, uint160[] sqrtPriceX96AfterList, uint32[] initializedTicksCrossedList, uint256 gasEstimate)
func (_QuoterV2 *QuoterV2TransactorSession) QuoteExactOutput(path []byte, amountOut *big.Int) (*types.Transaction, error) {
	return _QuoterV2.Contract.QuoteExactOutput(&_QuoterV2.TransactOpts, path, amountOut)
}

// QuoteExactOutputSingle is a paid mutator transaction binding the contract method 0xbd21704a.
//
// Solidity: function quoteExactOutputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV2 *QuoterV2Transactor) QuoteExactOutputSingle(opts *bind.TransactOpts, params IQuoterV2QuoteExactOutputSingleParams) (*types.Transaction, error) {
	return _QuoterV2.contract.Transact(opts, "quoteExactOutputSingle", params)
}

// QuoteExactOutputSingle is a paid mutator transaction binding the contract method 0xbd21704a.
//
// Solidity: function quoteExactOutputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV2 *QuoterV2Session) QuoteExactOutputSingle(params IQuoterV2QuoteExactOutputSingleParams) (*types.Transaction, error) {
	return _QuoterV2.Contract.QuoteExactOutputSingle(&_QuoterV2.TransactOpts, params)
}

// QuoteExactOutputSingle is a paid mutator transaction binding the contract method 0xbd21704a.
//
// Solidity: function quoteExactOutputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (_QuoterV2 *QuoterV2TransactorSession) QuoteExactOutputSingle(params IQuoterV2QuoteExactOutputSingleParams) (*types.Transaction, error) {
	return _QuoterV2.Contract.QuoteExactOutputSingle(&_QuoterV2.TransactOpts, params)
}
