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

// IV3MigratorMigrateParams is an auto generated low-level Go binding around an user-defined struct.
type IV3MigratorMigrateParams struct {
	Pair                common.Address
	LiquidityToMigrate  *big.Int
	PercentageToMigrate uint8
	Token0              common.Address
	Token1              common.Address
	Fee                 *big.Int
	TickLower           *big.Int
	TickUpper           *big.Int
	Amount0Min          *big.Int
	Amount1Min          *big.Int
	Recipient           common.Address
	Deadline            *big.Int
	RefundAsETH         bool
}

// V3MigratorMetaData contains all meta data concerning the V3Migrator contract.
var V3MigratorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_WETH9\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_nonfungiblePositionManager\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"WETH9\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96\",\"type\":\"uint160\"}],\"name\":\"createAndInitializePoolIfNecessary\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidityToMigrate\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"percentageToMigrate\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"internalType\":\"uint256\",\"name\":\"amount0Min\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1Min\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"refundAsETH\",\"type\":\"bool\"}],\"internalType\":\"structIV3Migrator.MigrateParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nonfungiblePositionManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitAllowed\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitAllowedIfNecessary\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitIfNecessary\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162001f4238038062001f42833981016040819052620000349162000079565b6001600160601b0319606093841b811660805291831b821660a05290911b1660c052620000c2565b80516001600160a01b03811681146200007457600080fd5b919050565b6000806000606084860312156200008e578283fd5b62000099846200005c565b9250620000a9602085016200005c565b9150620000b9604085016200005c565b90509250925092565b60805160601c60a05160601c60c05160601c611e146200012e6000398061090d5280610c4c5280610c865280610cb05280610e8d52508060d252806106b25280610ee65280610f70528061105852806110e2525080610269528061036452806109e65250611e146000f3fe6080604052600436106100b55760003560e01c8063b44a272211610069578063c45a01551161004e578063c45a0155146101e5578063d44f2bf2146101fa578063f3995c671461021a57610134565b8063b44a2722146101bd578063c2e3140a146101d257610134565b80634aa4a4fc1161009a5780634aa4a4fc14610175578063a4a78f0c1461018a578063ac9650d81461019d57610134565b806313ead562146101395780634659a4941461016257610134565b36610134573373ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001614610132576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161012990611c15565b60405180910390fd5b005b600080fd5b61014c610147366004611788565b61022d565b6040516101599190611afb565b60405180910390f35b6101326101703660046117e1565b6105f0565b34801561018157600080fd5b5061014c6106b0565b6101326101983660046117e1565b6106d4565b6101b06101ab36600461183a565b6107b1565b6040516101599190611b4d565b3480156101c957600080fd5b5061014c61090b565b6101326101e03660046117e1565b61092f565b3480156101f157600080fd5b5061014c6109e4565b34801561020657600080fd5b506101326102153660046119b7565b610a08565b6101326102283660046117e1565b611176565b60008373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff161061026757600080fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16631698ee828686866040518463ffffffff1660e01b8152600401808473ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1681526020018262ffffff168152602001935050505060206040518083038186803b15801561031957600080fd5b505afa15801561032d573d6000803e3d6000fd5b505050506040513d602081101561034357600080fd5b5051905073ffffffffffffffffffffffffffffffffffffffff81166104d3577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a16712958686866040518463ffffffff1660e01b8152600401808473ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1681526020018262ffffff1681526020019350505050602060405180830381600087803b15801561041657600080fd5b505af115801561042a573d6000803e3d6000fd5b505050506040513d602081101561044057600080fd5b5051604080517ff637731d00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff858116600483015291519293509083169163f637731d9160248082019260009290919082900301818387803b1580156104b657600080fd5b505af11580156104ca573d6000803e3d6000fd5b505050506105e8565b60008173ffffffffffffffffffffffffffffffffffffffff16633850c7bd6040518163ffffffff1660e01b815260040160e06040518083038186803b15801561051b57600080fd5b505afa15801561052f573d6000803e3d6000fd5b505050506040513d60e081101561054557600080fd5b5051905073ffffffffffffffffffffffffffffffffffffffff81166105e6578173ffffffffffffffffffffffffffffffffffffffff1663f637731d846040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff168152602001915050600060405180830381600087803b1580156105cd57600080fd5b505af11580156105e1573d6000803e3d6000fd5b505050505b505b949350505050565b604080517f8fcbaf0c00000000000000000000000000000000000000000000000000000000815233600482015230602482015260448101879052606481018690526001608482015260ff851660a482015260c4810184905260e48101839052905173ffffffffffffffffffffffffffffffffffffffff881691638fcbaf0c9161010480830192600092919082900301818387803b15801561069057600080fd5b505af11580156106a4573d6000803e3d6000fd5b50505050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b604080517fdd62ed3e00000000000000000000000000000000000000000000000000000000815233600482015230602482015290517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9173ffffffffffffffffffffffffffffffffffffffff89169163dd62ed3e91604480820192602092909190829003018186803b15801561076957600080fd5b505afa15801561077d573d6000803e3d6000fd5b505050506040513d602081101561079357600080fd5b505110156107a9576107a98686868686866105f0565b505050505050565b60608167ffffffffffffffff811180156107ca57600080fd5b506040519080825280602002602001820160405280156107fe57816020015b60608152602001906001900390816107e95790505b50905060005b82811015610904576000803086868581811061081c57fe5b905060200281019061082e9190611d3a565b60405161083c929190611aeb565b600060405180830381855af49150503d8060008114610877576040519150601f19603f3d011682016040523d82523d6000602084013e61087c565b606091505b5091509150816108e25760448151101561089557600080fd5b600481019050808060200190518101906108af9190611902565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101299190611bcb565b808484815181106108ef57fe5b60209081029190910101525050600101610804565b5092915050565b7f000000000000000000000000000000000000000000000000000000000000000081565b604080517fdd62ed3e0000000000000000000000000000000000000000000000000000000081523360048201523060248201529051869173ffffffffffffffffffffffffffffffffffffffff89169163dd62ed3e91604480820192602092909190829003018186803b1580156109a457600080fd5b505afa1580156109b8573d6000803e3d6000fd5b505050506040513d60208110156109ce57600080fd5b505110156107a9576107a9868686868686611176565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000610a1a6060830160408401611a5d565b60ff1611610a54576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161012990611bde565b6064610a666060830160408401611a5d565b60ff161115610aa1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161012990611c4c565b610aae6020820182611765565b73ffffffffffffffffffffffffffffffffffffffff166323b872dd33610ad76020850185611765565b84602001356040518463ffffffff1660e01b8152600401610afa93929190611b1c565b602060405180830381600087803b158015610b1457600080fd5b505af1158015610b28573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b4c91906118c5565b50600080610b5d6020840184611765565b73ffffffffffffffffffffffffffffffffffffffff166389afcb44306040518263ffffffff1660e01b8152600401610b959190611afb565b6040805180830381600087803b158015610bae57600080fd5b505af1158015610bc2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610be69190611a3a565b909250905060006064610c0c610c026060870160408801611a5d565b859060ff1661120e565b81610c1357fe5b04905060006064610c2d610c026060880160408901611a5d565b81610c3457fe5b049050610c71610c4a6080870160608801611765565b7f000000000000000000000000000000000000000000000000000000000000000084611238565b610cab610c8460a0870160808801611765565b7f000000000000000000000000000000000000000000000000000000000000000083611238565b6000807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663883164566040518061016001604052808a6060016020810190610d099190611765565b73ffffffffffffffffffffffffffffffffffffffff168152602001610d3460a08c0160808d01611765565b73ffffffffffffffffffffffffffffffffffffffff168152602001610d5f60c08c0160a08d016119cf565b62ffffff168152602001610d7960e08c0160c08d016118e1565b60020b8152602001610d926101008c0160e08d016118e1565b60020b815260208101889052604081018790526101008b013560608201526101208b0135608082015260a001610dd06101608c016101408d01611765565b73ffffffffffffffffffffffffffffffffffffffff1681526020018a61016001358152506040518263ffffffff1660e01b8152600401610e109190611c83565b608060405180830381600087803b158015610e2a57600080fd5b505af1158015610e3e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e6291906119e9565b935093505050858210156110025783821015610eb357610eb3610e8b6080890160608a01611765565b7f00000000000000000000000000000000000000000000000000000000000000006000611238565b818603610ec86101a089016101808a016118a9565b8015610f2e575073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016610f1660808a0160608b01611765565b73ffffffffffffffffffffffffffffffffffffffff16145b15610fe6576040517f2e1a7d4d00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690632e1a7d4d90610fa5908490600401611d31565b600060405180830381600087803b158015610fbf57600080fd5b505af1158015610fd3573d6000803e3d6000fd5b50505050610fe13382611414565b611000565b611000610ff960808a0160608b01611765565b3383611567565b505b8481101561116d578281101561102557611025610e8b60a0890160808a01611765565b80850361103a6101a089016101808a016118a9565b80156110a0575073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001661108860a08a0160808b01611765565b73ffffffffffffffffffffffffffffffffffffffff16145b15611158576040517f2e1a7d4d00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690632e1a7d4d90611117908490600401611d31565b600060405180830381600087803b15801561113157600080fd5b505af1158015611145573d6000803e3d6000fd5b505050506111533382611414565b61116b565b61116b610ff960a08a0160808b01611765565b505b50505050505050565b604080517fd505accf000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018790526064810186905260ff8516608482015260a4810184905260c48101839052905173ffffffffffffffffffffffffffffffffffffffff88169163d505accf9160e480830192600092919082900301818387803b15801561069057600080fd5b60008215806112295750508181028183828161122657fe5b04145b61123257600080fd5b92915050565b6040805173ffffffffffffffffffffffffffffffffffffffff8481166024830152604480830185905283518084039091018152606490920183526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f095ea7b3000000000000000000000000000000000000000000000000000000001781529251825160009485949389169392918291908083835b6020831061130d57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016112d0565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461136f576040519150601f19603f3d011682016040523d82523d6000602084013e611374565b606091505b50915091508180156113a25750805115806113a2575080806020019051602081101561139f57600080fd5b50515b61140d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600260248201527f5341000000000000000000000000000000000000000000000000000000000000604482015290519081900360640190fd5b5050505050565b6040805160008082526020820190925273ffffffffffffffffffffffffffffffffffffffff84169083906040518082805190602001908083835b6020831061148b57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0909201916020918201910161144e565b6001836020036101000a03801982511681845116808217855250505050505090500191505060006040518083038185875af1925050503d80600081146114ed576040519150601f19603f3d011682016040523d82523d6000602084013e6114f2565b606091505b505090508061156257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600360248201527f5354450000000000000000000000000000000000000000000000000000000000604482015290519081900360640190fd5b505050565b6040805173ffffffffffffffffffffffffffffffffffffffff8481166024830152604480830185905283518084039091018152606490920183526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb000000000000000000000000000000000000000000000000000000001781529251825160009485949389169392918291908083835b6020831061163c57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016115ff565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461169e576040519150601f19603f3d011682016040523d82523d6000602084013e6116a3565b606091505b50915091508180156116d15750805115806116d157508080602001905160208110156116ce57600080fd5b50515b61140d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600260248201527f5354000000000000000000000000000000000000000000000000000000000000604482015290519081900360640190fd5b803562ffffff8116811461174f57600080fd5b919050565b803560ff8116811461174f57600080fd5b600060208284031215611776578081fd5b813561178181611dd4565b9392505050565b6000806000806080858703121561179d578283fd5b84356117a881611dd4565b935060208501356117b881611dd4565b92506117c66040860161173c565b915060608501356117d681611dd4565b939692955090935050565b60008060008060008060c087890312156117f9578182fd5b863561180481611dd4565b9550602087013594506040870135935061182060608801611754565b92506080870135915060a087013590509295509295509295565b6000806020838503121561184c578182fd5b823567ffffffffffffffff80821115611863578384fd5b818501915085601f830112611876578384fd5b813581811115611884578485fd5b8660208083028501011115611897578485fd5b60209290920196919550909350505050565b6000602082840312156118ba578081fd5b813561178181611df9565b6000602082840312156118d6578081fd5b815161178181611df9565b6000602082840312156118f2578081fd5b81358060020b8114611781578182fd5b600060208284031215611913578081fd5b815167ffffffffffffffff8082111561192a578283fd5b818401915084601f83011261193d578283fd5b81518181111561194957fe5b60405160207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116820101818110848211171561198557fe5b60405281815283820160200187101561199c578485fd5b6119ad826020830160208701611da4565b9695505050505050565b60006101a082840312156119c9578081fd5b50919050565b6000602082840312156119e0578081fd5b6117818261173c565b600080600080608085870312156119fe578384fd5b8451935060208501516fffffffffffffffffffffffffffffffff81168114611a24578384fd5b6040860151606090960151949790965092505050565b60008060408385031215611a4c578182fd5b505080516020909101519092909150565b600060208284031215611a6e578081fd5b61178182611754565b73ffffffffffffffffffffffffffffffffffffffff169052565b60008151808452611aa9816020860160208601611da4565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60020b9052565b62ffffff169052565b6000828483379101908152919050565b73ffffffffffffffffffffffffffffffffffffffff91909116815260200190565b73ffffffffffffffffffffffffffffffffffffffff9384168152919092166020820152604081019190915260600190565b6000602080830181845280855180835260408601915060408482028701019250838701855b82811015611bbe577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452611bac858351611a91565b94509285019290850190600101611b72565b5092979650505050505050565b6000602082526117816020830184611a91565b60208082526014908201527f50657263656e7461676520746f6f20736d616c6c000000000000000000000000604082015260600190565b60208082526009908201527f4e6f742057455448390000000000000000000000000000000000000000000000604082015260600190565b60208082526014908201527f50657263656e7461676520746f6f206c61726765000000000000000000000000604082015260600190565b600061016082019050611c97828451611a77565b6020830151611ca96020840182611a77565b506040830151611cbc6040840182611ae2565b506060830151611ccf6060840182611adb565b506080830151611ce26080840182611adb565b5060a083015160a083015260c083015160c083015260e083015160e083015261010080840151818401525061012080840151611d2082850182611a77565b505061014092830151919092015290565b90815260200190565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112611d6e578283fd5b83018035915067ffffffffffffffff821115611d88578283fd5b602001915036819003821315611d9d57600080fd5b9250929050565b60005b83811015611dbf578181015183820152602001611da7565b83811115611dce576000848401525b50505050565b73ffffffffffffffffffffffffffffffffffffffff81168114611df657600080fd5b50565b8015158114611df657600080fdfea164736f6c6343000706000a",
}

// V3MigratorABI is the input ABI used to generate the binding from.
// Deprecated: Use V3MigratorMetaData.ABI instead.
var V3MigratorABI = V3MigratorMetaData.ABI

// V3MigratorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use V3MigratorMetaData.Bin instead.
var V3MigratorBin = V3MigratorMetaData.Bin

// DeployV3Migrator deploys a new Ethereum contract, binding an instance of V3Migrator to it.
func DeployV3Migrator(auth *bind.TransactOpts, backend bind.ContractBackend, _factory common.Address, _WETH9 common.Address, _nonfungiblePositionManager common.Address) (common.Address, *types.Transaction, *V3Migrator, error) {
	parsed, err := V3MigratorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(V3MigratorBin), backend, _factory, _WETH9, _nonfungiblePositionManager)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &V3Migrator{V3MigratorCaller: V3MigratorCaller{contract: contract}, V3MigratorTransactor: V3MigratorTransactor{contract: contract}, V3MigratorFilterer: V3MigratorFilterer{contract: contract}}, nil
}

// V3Migrator is an auto generated Go binding around an Ethereum contract.
type V3Migrator struct {
	V3MigratorCaller     // Read-only binding to the contract
	V3MigratorTransactor // Write-only binding to the contract
	V3MigratorFilterer   // Log filterer for contract events
}

// V3MigratorCaller is an auto generated read-only Go binding around an Ethereum contract.
type V3MigratorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V3MigratorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type V3MigratorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V3MigratorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type V3MigratorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V3MigratorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type V3MigratorSession struct {
	Contract     *V3Migrator       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// V3MigratorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type V3MigratorCallerSession struct {
	Contract *V3MigratorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// V3MigratorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type V3MigratorTransactorSession struct {
	Contract     *V3MigratorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// V3MigratorRaw is an auto generated low-level Go binding around an Ethereum contract.
type V3MigratorRaw struct {
	Contract *V3Migrator // Generic contract binding to access the raw methods on
}

// V3MigratorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type V3MigratorCallerRaw struct {
	Contract *V3MigratorCaller // Generic read-only contract binding to access the raw methods on
}

// V3MigratorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type V3MigratorTransactorRaw struct {
	Contract *V3MigratorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewV3Migrator creates a new instance of V3Migrator, bound to a specific deployed contract.
func NewV3Migrator(address common.Address, backend bind.ContractBackend) (*V3Migrator, error) {
	contract, err := bindV3Migrator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &V3Migrator{V3MigratorCaller: V3MigratorCaller{contract: contract}, V3MigratorTransactor: V3MigratorTransactor{contract: contract}, V3MigratorFilterer: V3MigratorFilterer{contract: contract}}, nil
}

// NewV3MigratorCaller creates a new read-only instance of V3Migrator, bound to a specific deployed contract.
func NewV3MigratorCaller(address common.Address, caller bind.ContractCaller) (*V3MigratorCaller, error) {
	contract, err := bindV3Migrator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &V3MigratorCaller{contract: contract}, nil
}

// NewV3MigratorTransactor creates a new write-only instance of V3Migrator, bound to a specific deployed contract.
func NewV3MigratorTransactor(address common.Address, transactor bind.ContractTransactor) (*V3MigratorTransactor, error) {
	contract, err := bindV3Migrator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &V3MigratorTransactor{contract: contract}, nil
}

// NewV3MigratorFilterer creates a new log filterer instance of V3Migrator, bound to a specific deployed contract.
func NewV3MigratorFilterer(address common.Address, filterer bind.ContractFilterer) (*V3MigratorFilterer, error) {
	contract, err := bindV3Migrator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &V3MigratorFilterer{contract: contract}, nil
}

// bindV3Migrator binds a generic wrapper to an already deployed contract.
func bindV3Migrator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := V3MigratorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_V3Migrator *V3MigratorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _V3Migrator.Contract.V3MigratorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_V3Migrator *V3MigratorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V3Migrator.Contract.V3MigratorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_V3Migrator *V3MigratorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _V3Migrator.Contract.V3MigratorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_V3Migrator *V3MigratorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _V3Migrator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_V3Migrator *V3MigratorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V3Migrator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_V3Migrator *V3MigratorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _V3Migrator.Contract.contract.Transact(opts, method, params...)
}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_V3Migrator *V3MigratorCaller) WETH9(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V3Migrator.contract.Call(opts, &out, "WETH9")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_V3Migrator *V3MigratorSession) WETH9() (common.Address, error) {
	return _V3Migrator.Contract.WETH9(&_V3Migrator.CallOpts)
}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_V3Migrator *V3MigratorCallerSession) WETH9() (common.Address, error) {
	return _V3Migrator.Contract.WETH9(&_V3Migrator.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_V3Migrator *V3MigratorCaller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V3Migrator.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_V3Migrator *V3MigratorSession) Factory() (common.Address, error) {
	return _V3Migrator.Contract.Factory(&_V3Migrator.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_V3Migrator *V3MigratorCallerSession) Factory() (common.Address, error) {
	return _V3Migrator.Contract.Factory(&_V3Migrator.CallOpts)
}

// NonfungiblePositionManager is a free data retrieval call binding the contract method 0xb44a2722.
//
// Solidity: function nonfungiblePositionManager() view returns(address)
func (_V3Migrator *V3MigratorCaller) NonfungiblePositionManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V3Migrator.contract.Call(opts, &out, "nonfungiblePositionManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NonfungiblePositionManager is a free data retrieval call binding the contract method 0xb44a2722.
//
// Solidity: function nonfungiblePositionManager() view returns(address)
func (_V3Migrator *V3MigratorSession) NonfungiblePositionManager() (common.Address, error) {
	return _V3Migrator.Contract.NonfungiblePositionManager(&_V3Migrator.CallOpts)
}

// NonfungiblePositionManager is a free data retrieval call binding the contract method 0xb44a2722.
//
// Solidity: function nonfungiblePositionManager() view returns(address)
func (_V3Migrator *V3MigratorCallerSession) NonfungiblePositionManager() (common.Address, error) {
	return _V3Migrator.Contract.NonfungiblePositionManager(&_V3Migrator.CallOpts)
}

// CreateAndInitializePoolIfNecessary is a paid mutator transaction binding the contract method 0x13ead562.
//
// Solidity: function createAndInitializePoolIfNecessary(address token0, address token1, uint24 fee, uint160 sqrtPriceX96) payable returns(address pool)
func (_V3Migrator *V3MigratorTransactor) CreateAndInitializePoolIfNecessary(opts *bind.TransactOpts, token0 common.Address, token1 common.Address, fee *big.Int, sqrtPriceX96 *big.Int) (*types.Transaction, error) {
	return _V3Migrator.contract.Transact(opts, "createAndInitializePoolIfNecessary", token0, token1, fee, sqrtPriceX96)
}

// CreateAndInitializePoolIfNecessary is a paid mutator transaction binding the contract method 0x13ead562.
//
// Solidity: function createAndInitializePoolIfNecessary(address token0, address token1, uint24 fee, uint160 sqrtPriceX96) payable returns(address pool)
func (_V3Migrator *V3MigratorSession) CreateAndInitializePoolIfNecessary(token0 common.Address, token1 common.Address, fee *big.Int, sqrtPriceX96 *big.Int) (*types.Transaction, error) {
	return _V3Migrator.Contract.CreateAndInitializePoolIfNecessary(&_V3Migrator.TransactOpts, token0, token1, fee, sqrtPriceX96)
}

// CreateAndInitializePoolIfNecessary is a paid mutator transaction binding the contract method 0x13ead562.
//
// Solidity: function createAndInitializePoolIfNecessary(address token0, address token1, uint24 fee, uint160 sqrtPriceX96) payable returns(address pool)
func (_V3Migrator *V3MigratorTransactorSession) CreateAndInitializePoolIfNecessary(token0 common.Address, token1 common.Address, fee *big.Int, sqrtPriceX96 *big.Int) (*types.Transaction, error) {
	return _V3Migrator.Contract.CreateAndInitializePoolIfNecessary(&_V3Migrator.TransactOpts, token0, token1, fee, sqrtPriceX96)
}

// Migrate is a paid mutator transaction binding the contract method 0xd44f2bf2.
//
// Solidity: function migrate((address,uint256,uint8,address,address,uint24,int24,int24,uint256,uint256,address,uint256,bool) params) returns()
func (_V3Migrator *V3MigratorTransactor) Migrate(opts *bind.TransactOpts, params IV3MigratorMigrateParams) (*types.Transaction, error) {
	return _V3Migrator.contract.Transact(opts, "migrate", params)
}

// Migrate is a paid mutator transaction binding the contract method 0xd44f2bf2.
//
// Solidity: function migrate((address,uint256,uint8,address,address,uint24,int24,int24,uint256,uint256,address,uint256,bool) params) returns()
func (_V3Migrator *V3MigratorSession) Migrate(params IV3MigratorMigrateParams) (*types.Transaction, error) {
	return _V3Migrator.Contract.Migrate(&_V3Migrator.TransactOpts, params)
}

// Migrate is a paid mutator transaction binding the contract method 0xd44f2bf2.
//
// Solidity: function migrate((address,uint256,uint8,address,address,uint24,int24,int24,uint256,uint256,address,uint256,bool) params) returns()
func (_V3Migrator *V3MigratorTransactorSession) Migrate(params IV3MigratorMigrateParams) (*types.Transaction, error) {
	return _V3Migrator.Contract.Migrate(&_V3Migrator.TransactOpts, params)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_V3Migrator *V3MigratorTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _V3Migrator.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_V3Migrator *V3MigratorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _V3Migrator.Contract.Multicall(&_V3Migrator.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_V3Migrator *V3MigratorTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _V3Migrator.Contract.Multicall(&_V3Migrator.TransactOpts, data)
}

// SelfPermit is a paid mutator transaction binding the contract method 0xf3995c67.
//
// Solidity: function selfPermit(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3Migrator *V3MigratorTransactor) SelfPermit(opts *bind.TransactOpts, token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3Migrator.contract.Transact(opts, "selfPermit", token, value, deadline, v, r, s)
}

// SelfPermit is a paid mutator transaction binding the contract method 0xf3995c67.
//
// Solidity: function selfPermit(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3Migrator *V3MigratorSession) SelfPermit(token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3Migrator.Contract.SelfPermit(&_V3Migrator.TransactOpts, token, value, deadline, v, r, s)
}

// SelfPermit is a paid mutator transaction binding the contract method 0xf3995c67.
//
// Solidity: function selfPermit(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3Migrator *V3MigratorTransactorSession) SelfPermit(token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3Migrator.Contract.SelfPermit(&_V3Migrator.TransactOpts, token, value, deadline, v, r, s)
}

// SelfPermitAllowed is a paid mutator transaction binding the contract method 0x4659a494.
//
// Solidity: function selfPermitAllowed(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3Migrator *V3MigratorTransactor) SelfPermitAllowed(opts *bind.TransactOpts, token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3Migrator.contract.Transact(opts, "selfPermitAllowed", token, nonce, expiry, v, r, s)
}

// SelfPermitAllowed is a paid mutator transaction binding the contract method 0x4659a494.
//
// Solidity: function selfPermitAllowed(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3Migrator *V3MigratorSession) SelfPermitAllowed(token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3Migrator.Contract.SelfPermitAllowed(&_V3Migrator.TransactOpts, token, nonce, expiry, v, r, s)
}

// SelfPermitAllowed is a paid mutator transaction binding the contract method 0x4659a494.
//
// Solidity: function selfPermitAllowed(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3Migrator *V3MigratorTransactorSession) SelfPermitAllowed(token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3Migrator.Contract.SelfPermitAllowed(&_V3Migrator.TransactOpts, token, nonce, expiry, v, r, s)
}

// SelfPermitAllowedIfNecessary is a paid mutator transaction binding the contract method 0xa4a78f0c.
//
// Solidity: function selfPermitAllowedIfNecessary(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3Migrator *V3MigratorTransactor) SelfPermitAllowedIfNecessary(opts *bind.TransactOpts, token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3Migrator.contract.Transact(opts, "selfPermitAllowedIfNecessary", token, nonce, expiry, v, r, s)
}

// SelfPermitAllowedIfNecessary is a paid mutator transaction binding the contract method 0xa4a78f0c.
//
// Solidity: function selfPermitAllowedIfNecessary(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3Migrator *V3MigratorSession) SelfPermitAllowedIfNecessary(token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3Migrator.Contract.SelfPermitAllowedIfNecessary(&_V3Migrator.TransactOpts, token, nonce, expiry, v, r, s)
}

// SelfPermitAllowedIfNecessary is a paid mutator transaction binding the contract method 0xa4a78f0c.
//
// Solidity: function selfPermitAllowedIfNecessary(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3Migrator *V3MigratorTransactorSession) SelfPermitAllowedIfNecessary(token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3Migrator.Contract.SelfPermitAllowedIfNecessary(&_V3Migrator.TransactOpts, token, nonce, expiry, v, r, s)
}

// SelfPermitIfNecessary is a paid mutator transaction binding the contract method 0xc2e3140a.
//
// Solidity: function selfPermitIfNecessary(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3Migrator *V3MigratorTransactor) SelfPermitIfNecessary(opts *bind.TransactOpts, token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3Migrator.contract.Transact(opts, "selfPermitIfNecessary", token, value, deadline, v, r, s)
}

// SelfPermitIfNecessary is a paid mutator transaction binding the contract method 0xc2e3140a.
//
// Solidity: function selfPermitIfNecessary(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3Migrator *V3MigratorSession) SelfPermitIfNecessary(token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3Migrator.Contract.SelfPermitIfNecessary(&_V3Migrator.TransactOpts, token, value, deadline, v, r, s)
}

// SelfPermitIfNecessary is a paid mutator transaction binding the contract method 0xc2e3140a.
//
// Solidity: function selfPermitIfNecessary(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_V3Migrator *V3MigratorTransactorSession) SelfPermitIfNecessary(token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _V3Migrator.Contract.SelfPermitIfNecessary(&_V3Migrator.TransactOpts, token, value, deadline, v, r, s)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_V3Migrator *V3MigratorTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V3Migrator.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_V3Migrator *V3MigratorSession) Receive() (*types.Transaction, error) {
	return _V3Migrator.Contract.Receive(&_V3Migrator.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_V3Migrator *V3MigratorTransactorSession) Receive() (*types.Transaction, error) {
	return _V3Migrator.Contract.Receive(&_V3Migrator.TransactOpts)
}

