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
	Bin: "0x60e06040523480156200001157600080fd5b5060405162001a3638038062001a36833981016040819052620000349162000079565b6001600160601b0319606093841b811660805291831b821660a05290911b1660c052620000c2565b80516001600160a01b03811681146200007457600080fd5b919050565b6000806000606084860312156200008e578283fd5b62000099846200005c565b9250620000a9602085016200005c565b9150620000b9604085016200005c565b90509250925092565b60805160601c60a05160601c60c05160601c6119086200012e600039806107735280610a3e5280610a785280610aa25280610c4b52508060a552806105765280610c975280610cee5280610dc95280610e2052508061020852806102cf528061082652506119086000f3fe6080604052600436106100955760003560e01c8063b44a272211610059578063b44a272214610176578063c2e3140a1461018b578063c45a01551461019e578063d44f2bf2146101b3578063f3995c67146101d3576100ed565b806313ead562146100f25780634659a4941461011b5780634aa4a4fc1461012e578063a4a78f0c14610143578063ac9650d814610156576100ed565b366100ed57336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146100eb5760405162461bcd60e51b81526004016100e290611728565b60405180910390fd5b005b600080fd5b610105610100366004611325565b6101e6565b604051610112919061164f565b60405180910390f35b6100eb61012936600461137e565b6104da565b34801561013a57600080fd5b50610105610574565b6100eb61015136600461137e565b610598565b6101696101643660046113d7565b610631565b6040516101129190611687565b34801561018257600080fd5b50610105610771565b6100eb61019936600461137e565b610795565b3480156101aa57600080fd5b50610105610824565b3480156101bf57600080fd5b506100eb6101ce366004611536565b610848565b6100eb6101e136600461137e565b610eb4565b6000836001600160a01b0316856001600160a01b03161061020657600080fd5b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316631698ee828686866040518463ffffffff1660e01b815260040180846001600160a01b03168152602001836001600160a01b031681526020018262ffffff168152602001935050505060206040518083038186803b15801561029157600080fd5b505afa1580156102a5573d6000803e3d6000fd5b505050506040513d60208110156102bb57600080fd5b505190506001600160a01b0381166103f1577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663a16712958686866040518463ffffffff1660e01b815260040180846001600160a01b03168152602001836001600160a01b031681526020018262ffffff1681526020019350505050602060405180830381600087803b15801561035a57600080fd5b505af115801561036e573d6000803e3d6000fd5b505050506040513d602081101561038457600080fd5b50516040805163f637731d60e01b81526001600160a01b03858116600483015291519293509083169163f637731d9160248082019260009290919082900301818387803b1580156103d457600080fd5b505af11580156103e8573d6000803e3d6000fd5b505050506104d2565b6000816001600160a01b0316633850c7bd6040518163ffffffff1660e01b815260040160e06040518083038186803b15801561042c57600080fd5b505afa158015610440573d6000803e3d6000fd5b505050506040513d60e081101561045657600080fd5b505190506001600160a01b0381166104d057816001600160a01b031663f637731d846040518263ffffffff1660e01b815260040180826001600160a01b03168152602001915050600060405180830381600087803b1580156104b757600080fd5b505af11580156104cb573d6000803e3d6000fd5b505050505b505b949350505050565b604080516323f2ebc360e21b815233600482015230602482015260448101879052606481018690526001608482015260ff851660a482015260c4810184905260e4810183905290516001600160a01b03881691638fcbaf0c9161010480830192600092919082900301818387803b15801561055457600080fd5b505af1158015610568573d6000803e3d6000fd5b50505050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b60408051636eb1769f60e11b81523360048201523060248201529051600019916001600160a01b0389169163dd62ed3e91604480820192602092909190829003018186803b1580156105e957600080fd5b505afa1580156105fd573d6000803e3d6000fd5b505050506040513d602081101561061357600080fd5b50511015610629576106298686868686866104da565b505050505050565b60608167ffffffffffffffff8111801561064a57600080fd5b5060405190808252806020026020018201604052801561067e57816020015b60608152602001906001900390816106695790505b50905060005b8281101561076a576000803086868581811061069c57fe5b90506020028101906106ae9190611830565b6040516106bc92919061163f565b600060405180830381855af49150503d80600081146106f7576040519150601f19603f3d011682016040523d82523d6000602084013e6106fc565b606091505b5091509150816107485760448151101561071557600080fd5b6004810190508080602001905181019061072f919061149f565b60405162461bcd60e51b81526004016100e291906116e7565b8084848151811061075557fe5b60209081029190910101525050600101610684565b5092915050565b7f000000000000000000000000000000000000000000000000000000000000000081565b60408051636eb1769f60e11b8152336004820152306024820152905186916001600160a01b0389169163dd62ed3e91604480820192602092909190829003018186803b1580156107e457600080fd5b505afa1580156107f8573d6000803e3d6000fd5b505050506040513d602081101561080e57600080fd5b5051101561062957610629868686868686610eb4565b7f000000000000000000000000000000000000000000000000000000000000000081565b600061085a60608301604084016115dc565b60ff161161087a5760405162461bcd60e51b81526004016100e2906116fa565b606461088c60608301604084016115dc565b60ff1611156108ad5760405162461bcd60e51b81526004016100e29061174b565b6108ba6020820182611302565b6001600160a01b03166323b872dd336108d66020850185611302565b84602001356040518463ffffffff1660e01b81526004016108f993929190611663565b602060405180830381600087803b15801561091357600080fd5b505af1158015610927573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061094b9190611462565b5060008061095c6020840184611302565b6001600160a01b03166389afcb44306040518263ffffffff1660e01b8152600401610987919061164f565b6040805180830381600087803b1580156109a057600080fd5b505af11580156109b4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109d891906115b9565b9092509050600060646109fe6109f460608701604088016115dc565b859060ff16610f26565b81610a0557fe5b04905060006064610a1f6109f460608801604089016115dc565b81610a2657fe5b049050610a63610a3c6080870160608801611302565b7f000000000000000000000000000000000000000000000000000000000000000084610f50565b610a9d610a7660a0870160808801611302565b7f000000000000000000000000000000000000000000000000000000000000000083610f50565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663883164566040518061016001604052808a6060016020810190610aee9190611302565b6001600160a01b03168152602001610b0c60a08c0160808d01611302565b6001600160a01b03168152602001610b2a60c08c0160a08d0161154e565b62ffffff168152602001610b4460e08c0160c08d0161147e565b60020b8152602001610b5d6101008c0160e08d0161147e565b60020b815260208101889052604081018790526101008b013560608201526101208b0135608082015260a001610b9b6101608c016101408d01611302565b6001600160a01b031681526020018a61016001358152506040518263ffffffff1660e01b8152600401610bce9190611779565b608060405180830381600087803b158015610be857600080fd5b505af1158015610bfc573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c209190611568565b93509350505085821015610d805783821015610c7157610c71610c496080890160608a01611302565b7f00000000000000000000000000000000000000000000000000000000000000006000610f50565b818603610c866101a089016101808a01611446565b8015610cd257506001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016610cc760808a0160608b01611302565b6001600160a01b0316145b15610d6457604051632e1a7d4d60e01b81526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690632e1a7d4d90610d23908490600401611827565b600060405180830381600087803b158015610d3d57600080fd5b505af1158015610d51573d6000803e3d6000fd5b50505050610d5f338261109e565b610d7e565b610d7e610d7760808a0160608b01611302565b3383611192565b505b84811015610eab5782811015610da357610da3610c4960a0890160808a01611302565b808503610db86101a089016101808a01611446565b8015610e0457506001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016610df960a08a0160808b01611302565b6001600160a01b0316145b15610e9657604051632e1a7d4d60e01b81526001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001690632e1a7d4d90610e55908490600401611827565b600060405180830381600087803b158015610e6f57600080fd5b505af1158015610e83573d6000803e3d6000fd5b50505050610e91338261109e565b610ea9565b610ea9610d7760a08a0160808b01611302565b505b50505050505050565b6040805163d505accf60e01b8152336004820152306024820152604481018790526064810186905260ff8516608482015260a4810184905260c4810183905290516001600160a01b0388169163d505accf9160e480830192600092919082900301818387803b15801561055457600080fd5b6000821580610f4157505081810281838281610f3e57fe5b04145b610f4a57600080fd5b92915050565b604080516001600160a01b038481166024830152604480830185905283518084039091018152606490920183526020820180516001600160e01b031663095ea7b360e01b1781529251825160009485949389169392918291908083835b60208310610fcc5780518252601f199092019160209182019101610fad565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461102e576040519150601f19603f3d011682016040523d82523d6000602084013e611033565b606091505b5091509150818015611061575080511580611061575080806020019051602081101561105e57600080fd5b50515b611097576040805162461bcd60e51b8152602060048201526002602482015261534160f01b604482015290519081900360640190fd5b5050505050565b604080516000808252602082019092526001600160a01b0384169083906040518082805190602001908083835b602083106110ea5780518252601f1990920191602091820191016110cb565b6001836020036101000a03801982511681845116808217855250505050505090500191505060006040518083038185875af1925050503d806000811461114c576040519150601f19603f3d011682016040523d82523d6000602084013e611151565b606091505b505090508061118d576040805162461bcd60e51b815260206004820152600360248201526253544560e81b604482015290519081900360640190fd5b505050565b604080516001600160a01b038481166024830152604480830185905283518084039091018152606490920183526020820180516001600160e01b031663a9059cbb60e01b1781529251825160009485949389169392918291908083835b6020831061120e5780518252601f1990920191602091820191016111ef565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114611270576040519150601f19603f3d011682016040523d82523d6000602084013e611275565b606091505b50915091508180156112a35750805115806112a357508080602001905160208110156112a057600080fd5b50515b611097576040805162461bcd60e51b815260206004820152600260248201526114d560f21b604482015290519081900360640190fd5b803562ffffff811681146112ec57600080fd5b919050565b803560ff811681146112ec57600080fd5b600060208284031215611313578081fd5b813561131e816118ac565b9392505050565b6000806000806080858703121561133a578283fd5b8435611345816118ac565b93506020850135611355816118ac565b9250611363604086016112d9565b91506060850135611373816118ac565b939692955090935050565b60008060008060008060c08789031215611396578182fd5b86356113a1816118ac565b955060208701359450604087013593506113bd606088016112f1565b92506080870135915060a087013590509295509295509295565b600080602083850312156113e9578182fd5b823567ffffffffffffffff80821115611400578384fd5b818501915085601f830112611413578384fd5b813581811115611421578485fd5b8660208083028501011115611434578485fd5b60209290920196919550909350505050565b600060208284031215611457578081fd5b813561131e816118c4565b600060208284031215611473578081fd5b815161131e816118c4565b60006020828403121561148f578081fd5b81358060020b811461131e578182fd5b6000602082840312156114b0578081fd5b815167ffffffffffffffff808211156114c7578283fd5b818401915084601f8301126114da578283fd5b8151818111156114e657fe5b604051601f8201601f19168101602001838111828210171561150457fe5b60405281815283820160200187101561151b578485fd5b61152c82602083016020870161187c565b9695505050505050565b60006101a08284031215611548578081fd5b50919050565b60006020828403121561155f578081fd5b61131e826112d9565b6000806000806080858703121561157d578384fd5b8451935060208501516fffffffffffffffffffffffffffffffff811681146115a3578384fd5b6040860151606090960151949790965092505050565b600080604083850312156115cb578182fd5b505080516020909101519092909150565b6000602082840312156115ed578081fd5b61131e826112f1565b6001600160a01b03169052565b6000815180845261161b81602086016020860161187c565b601f01601f19169290920160200192915050565b60020b9052565b62ffffff169052565b6000828483379101908152919050565b6001600160a01b0391909116815260200190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b6000602080830181845280855180835260408601915060408482028701019250838701855b828110156116da57603f198886030184526116c8858351611603565b945092850192908501906001016116ac565b5092979650505050505050565b60006020825261131e6020830184611603565b60208082526014908201527314195c98d95b9d1859d9481d1bdbc81cdb585b1b60621b604082015260600190565b6020808252600990820152684e6f7420574554483960b81b604082015260600190565b60208082526014908201527350657263656e7461676520746f6f206c6172676560601b604082015260600190565b60006101608201905061178d8284516115f6565b602083015161179f60208401826115f6565b5060408301516117b26040840182611636565b5060608301516117c5606084018261162f565b5060808301516117d8608084018261162f565b5060a083015160a083015260c083015160c083015260e083015160e083015261010080840151818401525061012080840151611816828501826115f6565b505061014092830151919092015290565b90815260200190565b6000808335601e19843603018112611846578283fd5b83018035915067ffffffffffffffff821115611860578283fd5b60200191503681900382131561187557600080fd5b9250929050565b60005b8381101561189757818101518382015260200161187f565b838111156118a6576000848401525b50505050565b6001600160a01b03811681146118c157600080fd5b50565b80151581146118c157600080fdfea2646970667358221220292d833aa33e7af11cbdaba59dea48c0434d735a917ba6a5b6a516c028bfa34864736f6c63430007060033",
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
