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

// IApproveAndCallIncreaseLiquidityParams is an auto generated low-level Go binding around an user-defined struct.
type IApproveAndCallIncreaseLiquidityParams struct {
	Token0     common.Address
	Token1     common.Address
	TokenId    *big.Int
	Amount0Min *big.Int
	Amount1Min *big.Int
}

// IApproveAndCallMintParams is an auto generated low-level Go binding around an user-defined struct.
type IApproveAndCallMintParams struct {
	Token0     common.Address
	Token1     common.Address
	Fee        *big.Int
	TickLower  *big.Int
	TickUpper  *big.Int
	Amount0Min *big.Int
	Amount1Min *big.Int
	Recipient  common.Address
}

// IV3SwapRouterExactInputParams is an auto generated low-level Go binding around an user-defined struct.
type IV3SwapRouterExactInputParams struct {
	Path             []byte
	Recipient        common.Address
	AmountIn         *big.Int
	AmountOutMinimum *big.Int
}

// IV3SwapRouterExactInputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type IV3SwapRouterExactInputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	Fee               *big.Int
	Recipient         common.Address
	AmountIn          *big.Int
	AmountOutMinimum  *big.Int
	SqrtPriceLimitX96 *big.Int
}

// IV3SwapRouterExactOutputParams is an auto generated low-level Go binding around an user-defined struct.
type IV3SwapRouterExactOutputParams struct {
	Path            []byte
	Recipient       common.Address
	AmountOut       *big.Int
	AmountInMaximum *big.Int
}

// IV3SwapRouterExactOutputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type IV3SwapRouterExactOutputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	Fee               *big.Int
	Recipient         common.Address
	AmountOut         *big.Int
	AmountInMaximum   *big.Int
	SqrtPriceLimitX96 *big.Int
}

// SwapRouter02MetaData contains all meta data concerning the SwapRouter02 contract.
var SwapRouter02MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factoryV2\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"factoryV3\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_positionManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_WETH9\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"WETH9\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"approveMax\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"approveMaxMinusOne\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"approveZeroThenMax\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"approveZeroThenMaxMinusOne\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"callPositionManager\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"result\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"paths\",\"type\":\"bytes[]\"},{\"internalType\":\"uint128[]\",\"name\":\"amounts\",\"type\":\"uint128[]\"},{\"internalType\":\"uint24\",\"name\":\"maximumTickDivergence\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"secondsAgo\",\"type\":\"uint32\"}],\"name\":\"checkOracleSlippage\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"uint24\",\"name\":\"maximumTickDivergence\",\"type\":\"uint24\"},{\"internalType\":\"uint32\",\"name\":\"secondsAgo\",\"type\":\"uint32\"}],\"name\":\"checkOracleSlippage\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"}],\"internalType\":\"structIV3SwapRouter.ExactInputParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structIV3SwapRouter.ExactInputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactInputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"path\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMaximum\",\"type\":\"uint256\"}],\"internalType\":\"structIV3SwapRouter.ExactOutputParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMaximum\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structIV3SwapRouter.ExactOutputSingleParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"exactOutputSingle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factoryV2\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"getApprovalType\",\"outputs\":[{\"internalType\":\"enumIApproveAndCall.ApprovalType\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount0Min\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1Min\",\"type\":\"uint256\"}],\"internalType\":\"structIApproveAndCall.IncreaseLiquidityParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"increaseLiquidity\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"result\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"internalType\":\"uint256\",\"name\":\"amount0Min\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1Min\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"internalType\":\"structIApproveAndCall.MintParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"result\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"previousBlockhash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"positionManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"pull\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"refundETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitAllowed\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitAllowedIfNecessary\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"selfPermitIfNecessary\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"swapExactTokensForTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMax\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"swapTokensForExactTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"sweepToken\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"}],\"name\":\"sweepToken\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeBips\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"}],\"name\":\"sweepTokenWithFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeBips\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"}],\"name\":\"sweepTokenWithFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"amount0Delta\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"amount1Delta\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"uniswapV3SwapCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"unwrapWETH9\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"}],\"name\":\"unwrapWETH9\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeBips\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"}],\"name\":\"unwrapWETH9WithFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountMinimum\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeBips\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"}],\"name\":\"unwrapWETH9WithFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"wrapETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6101006040526000196000553480156200001857600080fd5b5060405162006135380380620061358339810160408190526200003b9162000087565b6001600160601b0319606094851b811660805291841b821660a05291831b811660c052911b1660e052620000e3565b80516001600160a01b03811681146200008257600080fd5b919050565b600080600080608085870312156200009d578384fd5b620000a8856200006a565b9350620000b8602086016200006a565b9250620000c8604086016200006a565b9150620000d8606086016200006a565b905092959194509250565b60805160601c60a05160601c60c05160601c60e05160601c615fb162000184600039806102c15280610b3c52806112ad52806113d7528061147e52806116af52806117d95280612d8f5280612def5280612e70525080611e4c52806124df5280613cdb52508061166f5280611b1a5280611e9c52806132a6525080610c625280610d365280610fe2528061164b5280612fc252806131855250615fb16000f3fe6080604052600436106102a45760003560e01c80639b2c0a371161016e578063dee00f35116100cb578063f100b2051161007f578063f2d5d56b11610064578063f2d5d56b1461066e578063f3995c6714610681578063fa461e33146106945761034f565b8063f100b2051461063b578063f25801a71461064e5761034f565b8063e0e189a0116100b0578063e0e189a0146105f5578063e90a182f14610608578063efdeed8e1461061b5761034f565b8063dee00f35146105b5578063df2ab5bb146105e25761034f565b8063b858183f11610122578063c45a015511610107578063c45a01551461057a578063cab372ce1461058f578063d4ef38de146105a25761034f565b8063b858183f14610554578063c2e3140a146105675761034f565b8063ab3fdd5011610153578063ab3fdd501461051b578063ac9650d81461052e578063b3a2af13146105415761034f565b80639b2c0a37146104f5578063a4a78f0c146105085761034f565b8063472b43f31161021c578063571ac8b0116101d0578063639d71a9116101b5578063639d71a9146104b857806368e0d4e1146104cb578063791b98bc146104e05761034f565b8063571ac8b0146104925780635ae401dc146104a55761034f565b80634961699711610201578063496169971461044a5780634aa4a4fc1461045d5780635023b4df1461047f5761034f565b8063472b43f31461042457806349404b7c146104375761034f565b80631c58db4f116102735780633068c554116102585780633068c554146103eb57806342712a67146103fe5780634659a494146104115761034f565b80631c58db4f146103b85780631f0464d1146103cb5761034f565b806304e45aaf1461035457806309b813461461037d57806311ed56c91461039057806312210e8a146103b05761034f565b3661034f573373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461034d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f4e6f742057455448390000000000000000000000000000000000000000000000604482015290519081900360640190fd5b005b600080fd5b610367610362366004615543565b6106b4565b6040516103749190615dfd565b60405180910390f35b61036761038b3660046155de565b61083c565b6103a361039e366004615638565b61091c565b6040516103749190615b7a565b61034d610b28565b61034d6103c63660046157bb565b610b3a565b6103de6103d93660046152a7565b610bbe565b6040516103749190615afc565b61034d6103f93660046150d8565b610c48565b61036761040c366004615885565b610c5b565b61034d61041f366004615121565b610e35565b610367610432366004615885565b610ef5565b61034d6104453660046157eb565b6112a9565b61034d6104583660046157bb565b61146f565b34801561046957600080fd5b5061047261147c565b6040516103749190615a3c565b61036761048d366004615616565b6114a0565b61034d6104a0366004614feb565b611589565b6103de6104b33660046152a7565b6115bc565b61034d6104c6366004614feb565b611635565b3480156104d757600080fd5b50610472611649565b3480156104ec57600080fd5b5061047261166d565b61034d61050336600461581a565b611691565b61034d610516366004615121565b6118a7565b61034d610529366004614feb565b61197c565b6103de61053c36600461517c565b6119ba565b6103a361054f3660046152f1565b611b14565b61036761056236600461549d565b611bd2565b61034d610575366004615121565b611d95565b34801561058657600080fd5b50610472611e4a565b61034d61059d366004614feb565b611990565b61034d6105b0366004615858565b611e6e565b3480156105c157600080fd5b506105d56105d036600461500e565b611e7a565b6040516103749190615b8d565b61034d6105f0366004615039565b612027565b61034d61060336600461507a565b61213e565b61034d61061636600461500e565b6122a4565b34801561062757600080fd5b5061034d6106363660046151bc565b6122b3565b6103a3610649366004615627565b612305565b34801561065a57600080fd5b5061034d610669366004615324565b6123a5565b61034d61067c36600461500e565b6123f6565b61034d61068f366004615121565b612402565b3480156106a057600080fd5b5061034d6106af3660046153b8565b61249a565b600080600083608001511415610771575081516040517f70a0823100000000000000000000000000000000000000000000000000000000815260019173ffffffffffffffffffffffffffffffffffffffff16906370a082319061071b903090600401615a3c565b60206040518083038186803b15801561073357600080fd5b505afa158015610747573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061076b91906157d3565b60808401525b6107ed836080015184606001518560c001516040518060400160405280886000015189604001518a602001516040516020016107af939291906159aa565b6040516020818303038152906040528152602001866107ce57336107d0565b305b73ffffffffffffffffffffffffffffffffffffffff1690526125de565b91508260a00151821015610836576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161082d90615c7d565b60405180910390fd5b50919050565b60006108b0604083018035906108559060208601614feb565b604080518082019091526000908061086d8880615e41565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505050908252503360209091015261278f565b505060005460608201358111156108f3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161082d90615c0f565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600055919050565b604080516101608101909152606090610b20907f8831645600000000000000000000000000000000000000000000000000000000908061095f6020870187614feb565b73ffffffffffffffffffffffffffffffffffffffff16815260200185602001602081019061098d9190614feb565b73ffffffffffffffffffffffffffffffffffffffff1681526020016109b860608701604088016157a1565b62ffffff1681526020016109d26080870160608801615379565b60020b81526020016109ea60a0870160808801615379565b60020b8152602090810190610a0a90610a0590880188614feb565b612976565b8152602001610a25866020016020810190610a059190614feb565b815260a0860135602082015260c08601356040820152606001610a4f610100870160e08801614feb565b73ffffffffffffffffffffffffffffffffffffffff1681526020017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff815250604051602401610a9e9190615cf8565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152611b14565b90505b919050565b4715610b3857610b383347612a1b565b565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663d0e30db0826040518263ffffffff1660e01b81526004016000604051808303818588803b158015610ba257600080fd5b505af1158015610bb6573d6000803e3d6000fd5b505050505050565b60608380600143034014610c3357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f426c6f636b686173680000000000000000000000000000000000000000000000604482015290519081900360640190fd5b610c3d84846119ba565b91505b509392505050565b610c55848433858561213e565b50505050565b6000610cbb7f000000000000000000000000000000000000000000000000000000000000000087868680806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250612b6992505050565b600081518110610cc757fe5b6020026020010151905084811115610d0b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161082d90615c0f565b610da484846000818110610d1b57fe5b9050602002016020810190610d309190614feb565b33610d9e7f000000000000000000000000000000000000000000000000000000000000000088886000818110610d6257fe5b9050602002016020810190610d779190614feb565b89896001818110610d8457fe5b9050602002016020810190610d999190614feb565b612ca2565b84612d8d565b73ffffffffffffffffffffffffffffffffffffffff821660011415610dcb57339150610dee565b73ffffffffffffffffffffffffffffffffffffffff821660021415610dee573091505b610e2c848480806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250869250612f6b915050565b95945050505050565b604080517f8fcbaf0c00000000000000000000000000000000000000000000000000000000815233600482015230602482015260448101879052606481018690526001608482015260ff851660a482015260c4810184905260e48101839052905173ffffffffffffffffffffffffffffffffffffffff881691638fcbaf0c9161010480830192600092919082900301818387803b158015610ed557600080fd5b505af1158015610ee9573d6000803e3d6000fd5b50505050505050505050565b60008086610fab575060018484600081610f0b57fe5b9050602002016020810190610f209190614feb565b73ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401610f589190615a3c565b60206040518083038186803b158015610f7057600080fd5b505afa158015610f84573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fa891906157d3565b96505b61103685856000818110610fbb57fe5b9050602002016020810190610fd09190614feb565b82610fdb5733610fdd565b305b6110307f00000000000000000000000000000000000000000000000000000000000000008989600081811061100e57fe5b90506020020160208101906110239190614feb565b8a8a6001818110610d8457fe5b8a612d8d565b73ffffffffffffffffffffffffffffffffffffffff83166001141561105d57339250611080565b73ffffffffffffffffffffffffffffffffffffffff831660021415611080573092505b600085857fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81018181106110b057fe5b90506020020160208101906110c59190614feb565b73ffffffffffffffffffffffffffffffffffffffff166370a08231856040518263ffffffff1660e01b81526004016110fd9190615a3c565b60206040518083038186803b15801561111557600080fd5b505afa158015611129573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061114d91906157d3565b905061118d868680806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250889250612f6b915050565b6112628187877fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81018181106111bf57fe5b90506020020160208101906111d49190614feb565b73ffffffffffffffffffffffffffffffffffffffff166370a08231876040518263ffffffff1660e01b815260040161120c9190615a3c565b60206040518083038186803b15801561122457600080fd5b505afa158015611238573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061125c91906157d3565b90613270565b92508683101561129e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161082d90615c7d565b505095945050505050565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561133257600080fd5b505afa158015611346573d6000803e3d6000fd5b505050506040513d602081101561135c57600080fd5b50519050828110156113cf57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f496e73756666696369656e742057455448390000000000000000000000000000604482015290519081900360640190fd5b801561146a577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16632e1a7d4d826040518263ffffffff1660e01b815260040180828152602001915050600060405180830381600087803b15801561144857600080fd5b505af115801561145c573d6000803e3d6000fd5b5050505061146a8282612a1b565b505050565b61147981336112a9565b50565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000611549608083018035906114b99060608601614feb565b6114c960e0860160c08701614feb565b60405180604001604052808760200160208101906114e79190614feb565b6114f760608a0160408b016157a1565b61150460208b018b614feb565b604051602001611516939291906159aa565b60405160208183030381529060405281526020013373ffffffffffffffffffffffffffffffffffffffff1681525061278f565b90508160a001358111156108f3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161082d90615c0f565b6115b3817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff613280565b61147957600080fd5b606083806115c86133cc565b1115610c3357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f5472616e73616374696f6e20746f6f206f6c6400000000000000000000000000604482015290519081900360640190fd5b611640816000613280565b61158957600080fd5b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000821180156116a2575060648211155b6116ab57600080fd5b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561173457600080fd5b505afa158015611748573d6000803e3d6000fd5b505050506040513d602081101561175e57600080fd5b50519050848110156117d157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f496e73756666696369656e742057455448390000000000000000000000000000604482015290519081900360640190fd5b80156118a0577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16632e1a7d4d826040518263ffffffff1660e01b815260040180828152602001915050600060405180830381600087803b15801561184a57600080fd5b505af115801561185e573d6000803e3d6000fd5b50505050600061271061187a85846133d090919063ffffffff16565b8161188157fe5b0490508015611894576118948382612a1b565b610bb685828403612a1b565b5050505050565b604080517fdd62ed3e00000000000000000000000000000000000000000000000000000000815233600482015230602482015290517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9173ffffffffffffffffffffffffffffffffffffffff89169163dd62ed3e91604480820192602092909190829003018186803b15801561193c57600080fd5b505afa158015611950573d6000803e3d6000fd5b505050506040513d602081101561196657600080fd5b50511015610bb657610bb6868686868686610e35565b611987816000613280565b61199057600080fd5b6115b3817ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe613280565b60608167ffffffffffffffff811180156119d357600080fd5b50604051908082528060200260200182016040528015611a0757816020015b60608152602001906001900390816119f25790505b50905060005b82811015611b0d5760008030868685818110611a2557fe5b9050602002810190611a379190615e41565b604051611a45929190615a10565b600060405180830381855af49150503d8060008114611a80576040519150601f19603f3d011682016040523d82523d6000602084013e611a85565b606091505b509150915081611aeb57604481511015611a9e57600080fd5b60048101905080806020019051810190611ab89190615433565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161082d9190615b7a565b80848481518110611af857fe5b60209081029190910101525050600101611a0d565b5092915050565b606060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1683604051611b5d9190615a20565b6000604051808303816000865af19150503d8060008114611b9a576040519150601f19603f3d011682016040523d82523d6000602084013e611b9f565b606091505b50925090508061083657604482511015611bb857600080fd5b60048201915081806020019051810190611ab89190615433565b600080600083604001511415611ca357600190506000611bf584600001516133f4565b50506040517f70a0823100000000000000000000000000000000000000000000000000000000815290915073ffffffffffffffffffffffffffffffffffffffff8216906370a0823190611c4c903090600401615a3c565b60206040518083038186803b158015611c6457600080fd5b505afa158015611c78573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c9c91906157d3565b6040850152505b600081611cb05733611cb2565b305b90505b6000611cc48560000151613425565b9050611d1d856040015182611cdd578660200151611cdf565b305b60006040518060400160405280611cf98b6000015161342d565b81526020018773ffffffffffffffffffffffffffffffffffffffff168152506125de565b60408601528015611d3d578451309250611d369061343c565b8552611d4a565b8460400151935050611d50565b50611cb5565b8360600151831015611d8e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161082d90615c7d565b5050919050565b604080517fdd62ed3e0000000000000000000000000000000000000000000000000000000081523360048201523060248201529051869173ffffffffffffffffffffffffffffffffffffffff89169163dd62ed3e91604480820192602092909190829003018186803b158015611e0a57600080fd5b505afa158015611e1e573d6000803e3d6000fd5b505050506040513d6020811015611e3457600080fd5b50511015610bb657610bb6868686868686612402565b7f000000000000000000000000000000000000000000000000000000000000000081565b61146a83338484611691565b6000818373ffffffffffffffffffffffffffffffffffffffff1663dd62ed3e307f00000000000000000000000000000000000000000000000000000000000000006040518363ffffffff1660e01b8152600401611ed8929190615a5d565b60206040518083038186803b158015611ef057600080fd5b505afa158015611f04573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f2891906157d3565b10611f3557506000612021565b611f5f837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff613280565b15611f6c57506001612021565b611f96837ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe613280565b15611fa357506002612021565b611fae836000613280565b611fb757600080fd5b611fe1837fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff613280565b15611fee57506003612021565b612018837ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe613280565b1561034f575060045b92915050565b60008373ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561209057600080fd5b505afa1580156120a4573d6000803e3d6000fd5b505050506040513d60208110156120ba57600080fd5b505190508281101561212d57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f496e73756666696369656e7420746f6b656e0000000000000000000000000000604482015290519081900360640190fd5b8015610c5557610c55848383613471565b60008211801561214f575060648211155b61215857600080fd5b60008573ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156121c157600080fd5b505afa1580156121d5573d6000803e3d6000fd5b505050506040513d60208110156121eb57600080fd5b505190508481101561225e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f496e73756666696369656e7420746f6b656e0000000000000000000000000000604482015290519081900360640190fd5b8015610bb657600061271061227383866133d0565b8161227a57fe5b049050801561228e5761228e878483613471565b61229b8786838503613471565b50505050505050565b6122af828233612027565b5050565b6000806122c1868685613646565b915091508362ffffff1681830312610bb6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161082d90615c46565b6060610b2063219f5d1760e01b6040518060c001604052808560400135815260200161233d866000016020810190610a059190614feb565b8152602001612358866020016020810190610a059190614feb565b815260200185606001358152602001856080013581526020017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff815250604051602401610a9e9190615cb4565b6000806123b28584613859565b915091508362ffffff16818303126118a0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161082d90615c46565b6122af82333084613ae1565b604080517fd505accf000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018790526064810186905260ff8516608482015260a4810184905260c48101839052905173ffffffffffffffffffffffffffffffffffffffff88169163d505accf9160e480830192600092919082900301818387803b158015610ed557600080fd5b60008413806124a95750600083135b6124b257600080fd5b60006124c08284018461564a565b905060008060006124d484600001516133f4565b9250925092506125067f0000000000000000000000000000000000000000000000000000000000000000848484613cbe565b5060008060008a13612547578473ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161089612578565b8373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff16108a5b915091508115612597576125928587602001513384612d8d565b610ee9565b85516125a290613425565b156125c75785516125b29061343c565b86526125c1813360008961278f565b50610ee9565b80600081905550610ee98487602001513384612d8d565b600073ffffffffffffffffffffffffffffffffffffffff8416600114156126075733935061262a565b73ffffffffffffffffffffffffffffffffffffffff84166002141561262a573093505b600080600061263c85600001516133f4565b9194509250905073ffffffffffffffffffffffffffffffffffffffff8083169084161060008061266d868686613cd4565b73ffffffffffffffffffffffffffffffffffffffff1663128acb088b856126938f613d12565b73ffffffffffffffffffffffffffffffffffffffff8e16156126b5578d6126db565b876126d45773fffd8963efd1fc6a506488495d951d5263988d256126db565b6401000276a45b8d6040516020016126ec9190615da6565b6040516020818303038152906040526040518663ffffffff1660e01b815260040161271b959493929190615a84565b6040805180830381600087803b15801561273457600080fd5b505af1158015612748573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061276c9190615395565b915091508261277b578161277d565b805b6000039b9a5050505050505050505050565b600073ffffffffffffffffffffffffffffffffffffffff8416600114156127b8573393506127db565b73ffffffffffffffffffffffffffffffffffffffff8416600214156127db573093505b60008060006127ed85600001516133f4565b9194509250905073ffffffffffffffffffffffffffffffffffffffff8084169083161060008061281e858786613cd4565b73ffffffffffffffffffffffffffffffffffffffff1663128acb088b856128448f613d12565b60000373ffffffffffffffffffffffffffffffffffffffff8e1615612869578d61288f565b876128885773fffd8963efd1fc6a506488495d951d5263988d2561288f565b6401000276a45b8d6040516020016128a09190615da6565b6040516020818303038152906040526040518663ffffffff1660e01b81526004016128cf959493929190615a84565b6040805180830381600087803b1580156128e857600080fd5b505af11580156128fc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129209190615395565b9150915060008361293557818360000361293b565b82826000035b909850905073ffffffffffffffffffffffffffffffffffffffff8a16612967578b811461296757600080fd5b50505050505050949350505050565b6040517f70a0823100000000000000000000000000000000000000000000000000000000815260009073ffffffffffffffffffffffffffffffffffffffff8316906370a08231906129cb903090600401615a3c565b60206040518083038186803b1580156129e357600080fd5b505afa1580156129f7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b2091906157d3565b6040805160008082526020820190925273ffffffffffffffffffffffffffffffffffffffff84169083906040518082805190602001908083835b60208310612a9257805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101612a55565b6001836020036101000a03801982511681845116808217855250505050505090500191505060006040518083038185875af1925050503d8060008114612af4576040519150601f19603f3d011682016040523d82523d6000602084013e612af9565b606091505b505090508061146a57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600360248201527f5354450000000000000000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6060600282511015612b7a57600080fd5b815167ffffffffffffffff81118015612b9257600080fd5b50604051908082528060200260200182016040528015612bbc578160200160208202803683370190505b5090508281600183510381518110612bd057fe5b602090810291909101015281517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff015b8015610c4057600080612c3d87866001860381518110612c1c57fe5b6020026020010151878681518110612c3057fe5b6020026020010151613d44565b91509150612c5f848481518110612c5057fe5b60200260200101518383613e2c565b846001850381518110612c6e57fe5b602090810291909101015250507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01612c00565b6000806000612cb18585613f02565b604080517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606094851b811660208084019190915293851b81166034830152825160288184030181526048830184528051908501207fff0000000000000000000000000000000000000000000000000000000000000060688401529a90941b9093166069840152607d8301989098527f96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f609d808401919091528851808403909101815260bd909201909752805196019590952095945050505050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16148015612de85750804710155b15612f31577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663d0e30db0826040518263ffffffff1660e01b81526004016000604051808303818588803b158015612e5557600080fd5b505af1158015612e69573d6000803e3d6000fd5b50505050507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb83836040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015612eff57600080fd5b505af1158015612f13573d6000803e3d6000fd5b505050506040513d6020811015612f2957600080fd5b50610c559050565b73ffffffffffffffffffffffffffffffffffffffff8316301415612f5f57612f5a848383613471565b610c55565b610c5584848484613ae1565b60005b600183510381101561146a57600080848381518110612f8957fe5b6020026020010151858460010181518110612fa057fe5b6020026020010151915091506000612fb88383613f02565b5090506000612fe87f00000000000000000000000000000000000000000000000000000000000000008585612ca2565b90506000806000808473ffffffffffffffffffffffffffffffffffffffff16630902f1ac6040518163ffffffff1660e01b815260040160606040518083038186803b15801561303657600080fd5b505afa15801561304a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061306e91906156da565b506dffffffffffffffffffffffffffff1691506dffffffffffffffffffffffffffff1691506000808773ffffffffffffffffffffffffffffffffffffffff168a73ffffffffffffffffffffffffffffffffffffffff16146130d05782846130d3565b83835b91509150613114828b73ffffffffffffffffffffffffffffffffffffffff166370a082318a6040518263ffffffff1660e01b815260040161120c9190615a3c565b9550613121868383613fa7565b9450505050506000808573ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff161461316557826000613169565b6000835b91509150600060028c51038a10613180578a6131c1565b6131c17f0000000000000000000000000000000000000000000000000000000000000000898e8d600201815181106131b457fe5b6020026020010151612ca2565b604080516000815260208101918290527f022c0d9f0000000000000000000000000000000000000000000000000000000090915290915073ffffffffffffffffffffffffffffffffffffffff87169063022c0d9f906132299086908690869060248101615e06565b600060405180830381600087803b15801561324357600080fd5b505af1158015613257573d6000803e3d6000fd5b50506001909b019a50612f6e9950505050505050505050565b8082038281111561202157600080fd5b60008060008473ffffffffffffffffffffffffffffffffffffffff1663095ea7b360e01b7f0000000000000000000000000000000000000000000000000000000000000000866040516024016132d7929190615ad6565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009094169390931790925290516133609190615a20565b6000604051808303816000865af19150503d806000811461339d576040519150601f19603f3d011682016040523d82523d6000602084013e6133a2565b606091505b5091509150818015610e2c575080511580610e2c575080806020019051810190610e2c919061528d565b4290565b60008215806133eb575050818102818382816133e857fe5b04145b61202157600080fd5b60008080613402848261407d565b925061340f84601461417d565b905061341c84601761407d565b91509193909250565b516042111590565b6060610b20826000602b61426d565b8051606090610b209083906017907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe90161426d565b6040805173ffffffffffffffffffffffffffffffffffffffff8481166024830152604480830185905283518084039091018152606490920183526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb000000000000000000000000000000000000000000000000000000001781529251825160009485949389169392918291908083835b6020831061354657805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101613509565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d80600081146135a8576040519150601f19603f3d011682016040523d82523d6000602084013e6135ad565b606091505b50915091508180156135db5750805115806135db57508080602001905160208110156135d857600080fd5b50515b6118a057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600260248201527f5354000000000000000000000000000000000000000000000000000000000000604482015290519081900360640190fd5b600080835185511461365757600080fd5b6000855167ffffffffffffffff8111801561367157600080fd5b506040519080825280602002602001820160405280156136ab57816020015b613698614e34565b8152602001906001900390816136905790505b5090506000865167ffffffffffffffff811180156136c857600080fd5b5060405190808252806020026020018201604052801561370257816020015b6136ef614e34565b8152602001906001900390816136e75790505b50905060005b8751811015613832576000806137318a848151811061372357fe5b602002602001015189613859565b9150915061373e82614454565b85848151811061374a57fe5b60200260200101516000019060020b908160020b8152505061376b81614454565b84848151811061377757fe5b60200260200101516000019060020b908160020b8152505088838151811061379b57fe5b60200260200101518584815181106137af57fe5b6020026020010151602001906fffffffffffffffffffffffffffffffff1690816fffffffffffffffffffffffffffffffff16815250508883815181106137f157fe5b602002602001015184848151811061380557fe5b6020908102919091018101516fffffffffffffffffffffffffffffffff9092169101525050600101613708565b5061383c82614465565b60020b935061384a81614465565b60020b92505050935093915050565b6000806000806138688661454d565b90506000805b82811015613a865760008060006138848b6133f4565b9250925092506000613897848484613cd4565b905060008063ffffffff8d166138c0576138b083614578565b600291820b9350900b9050613962565b6138ca838e614810565b8160020b915050809250508273ffffffffffffffffffffffffffffffffffffffff16633850c7bd6040518163ffffffff1660e01b815260040160e06040518083038186803b15801561391b57600080fd5b505afa15801561392f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906139539190615715565b50505060029290920b93505050505b600189038714156139a3578473ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff161099506139b2565b6139ac8e61343c565b9d508597505b6000871580613a5357508673ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff1610613a23578673ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff1610613a53565b8573ffffffffffffffffffffffffffffffffffffffff168773ffffffffffffffffffffffffffffffffffffffff16105b90508015613a68579b82019b9a81019a613a73565b828d039c50818c039b505b50506001909501945061386e9350505050565b5082613ad7577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff850294507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff840293505b5050509250929050565b6040805173ffffffffffffffffffffffffffffffffffffffff85811660248301528481166044830152606480830185905283518084039091018152608490920183526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f23b872dd00000000000000000000000000000000000000000000000000000000178152925182516000948594938a169392918291908083835b60208310613bbe57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101613b81565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114613c20576040519150601f19603f3d011682016040523d82523d6000602084013e613c25565b606091505b5091509150818015613c53575080511580613c535750808060200190516020811015613c5057600080fd5b50515b610bb657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600360248201527f5354460000000000000000000000000000000000000000000000000000000000604482015290519081900360640190fd5b6000610e2c85613ccf868686614c41565b614cbe565b6000613d0a7f0000000000000000000000000000000000000000000000000000000000000000613d05868686614c41565b614cee565b949350505050565b60007f80000000000000000000000000000000000000000000000000000000000000008210613d4057600080fd5b5090565b6000806000613d538585613f02565b509050600080613d64888888612ca2565b73ffffffffffffffffffffffffffffffffffffffff16630902f1ac6040518163ffffffff1660e01b815260040160606040518083038186803b158015613da957600080fd5b505afa158015613dbd573d6000803e3d6000fd5b505050506040513d6060811015613dd357600080fd5b5080516020909101516dffffffffffffffffffffffffffff918216935016905073ffffffffffffffffffffffffffffffffffffffff87811690841614613e1a578082613e1d565b81815b90999098509650505050505050565b6000808411613e9c57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f494e53554646494349454e545f4f55545055545f414d4f554e54000000000000604482015290519081900360640190fd5b600083118015613eac5750600082115b613eb557600080fd5b6000613ecd6103e8613ec786886133d0565b906133d0565b90506000613ee16103e5613ec78689613270565b9050613ef86001828481613ef157fe5b0490614e24565b9695505050505050565b6000808273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161415613f3e57600080fd5b8273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1610613f78578284613f7b565b83835b909250905073ffffffffffffffffffffffffffffffffffffffff8216613fa057600080fd5b9250929050565b600080841161401757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f494e53554646494349454e545f494e5055545f414d4f554e5400000000000000604482015290519081900360640190fd5b6000831180156140275750600082115b61403057600080fd5b600061403e856103e56133d0565b9050600061404c82856133d0565b9050600061406683614060886103e86133d0565b90614e24565b905080828161407157fe5b04979650505050505050565b6000818260140110156140f157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f746f416464726573735f6f766572666c6f770000000000000000000000000000604482015290519081900360640190fd5b816014018351101561416457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f746f416464726573735f6f75744f66426f756e64730000000000000000000000604482015290519081900360640190fd5b5001602001516c01000000000000000000000000900490565b6000818260030110156141f157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f746f55696e7432345f6f766572666c6f77000000000000000000000000000000604482015290519081900360640190fd5b816003018351101561426457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f746f55696e7432345f6f75744f66426f756e6473000000000000000000000000604482015290519081900360640190fd5b50016003015190565b60608182601f0110156142e157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f736c6963655f6f766572666c6f77000000000000000000000000000000000000604482015290519081900360640190fd5b82828401101561435257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f736c6963655f6f766572666c6f77000000000000000000000000000000000000604482015290519081900360640190fd5b818301845110156143c457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f736c6963655f6f75744f66426f756e6473000000000000000000000000000000604482015290519081900360640190fd5b6060821580156143e3576040519150600082526020820160405261444b565b6040519150601f8416801560200281840101858101878315602002848b0101015b8183101561441c578051835260209283019201614404565b5050858452601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016604052505b50949350505050565b80600281900b8114610b2357600080fd5b6000806000805b84518110156144fa5784818151811061448157fe5b6020026020010151602001516fffffffffffffffffffffffffffffffff168582815181106144ab57fe5b60200260200101516000015160020b02830192508481815181106144cb57fe5b6020026020010151602001516fffffffffffffffffffffffffffffffff1682019150808060010191505061446c565b5080828161450457fe5b05925060008212801561451f575080828161451b57fe5b0715155b15611d8e5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01919050565b5160177fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec9091010490565b6000806000808473ffffffffffffffffffffffffffffffffffffffff16633850c7bd6040518163ffffffff1660e01b815260040160e06040518083038186803b1580156145c457600080fd5b505afa1580156145d8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906145fc9190615715565b50939750919550935050600161ffff84161191506146489050576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161082d90615bd8565b6000808673ffffffffffffffffffffffffffffffffffffffff1663252c09d7856040518263ffffffff1660e01b81526004016146849190615dee565b60806040518083038186803b15801561469c57600080fd5b505afa1580156146b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906146d491906158e0565b5050915091506146e26133cc565b63ffffffff168263ffffffff16146146fc57849550614807565b60008361ffff1660018561ffff168761ffff1601038161471857fe5b06905060008060008a73ffffffffffffffffffffffffffffffffffffffff1663252c09d7856040518263ffffffff1660e01b81526004016147599190615dfd565b60806040518083038186803b15801561477157600080fd5b505afa158015614785573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906147a991906158e0565b93505092509250806147e7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161082d90615ba1565b82860363ffffffff811683870360060b816147fe57fe5b059a5050505050505b50505050915091565b60008063ffffffff831661488557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600260248201527f4250000000000000000000000000000000000000000000000000000000000000604482015290519081900360640190fd5b60408051600280825260608201835260009260208301908036833701905050905083816000815181106148b457fe5b602002602001019063ffffffff16908163ffffffff16815250506000816001815181106148dd57fe5b63ffffffff9092166020928302919091018201526040517f883bdbfd00000000000000000000000000000000000000000000000000000000815260048101828152835160248301528351600093849373ffffffffffffffffffffffffffffffffffffffff8b169363883bdbfd9388939192839260449091019185820191028083838b5b83811015614978578181015183820152602001614960565b505050509050019250505060006040518083038186803b15801561499b57600080fd5b505afa1580156149af573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160409081528110156149f657600080fd5b8101908080516040519392919084640100000000821115614a1657600080fd5b908301906020820185811115614a2b57600080fd5b8251866020820283011164010000000082111715614a4857600080fd5b82525081516020918201928201910280838360005b83811015614a75578181015183820152602001614a5d565b5050505090500160405260200180516040519392919084640100000000821115614a9e57600080fd5b908301906020820185811115614ab357600080fd5b8251866020820283011164010000000082111715614ad057600080fd5b82525081516020918201928201910280838360005b83811015614afd578181015183820152602001614ae5565b5050505090500160405250505091509150600082600081518110614b1d57fe5b602002602001015183600181518110614b3257fe5b6020026020010151039050600082600081518110614b4c57fe5b602002602001015183600181518110614b6157fe5b60200260200101510390508763ffffffff168260060b81614b7e57fe5b05965060008260060b128015614ba857508763ffffffff168260060b81614ba157fe5b0760060b15155b15614bd3577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff909601955b63ffffffff881673ffffffffffffffffffffffffffffffffffffffff0277ffffffffffffffffffffffffffffffffffffffff00000000602083901b1677ffffffffffffffffffffffffffffffffffffffffffffffff821681614c3157fe5b0496505050505050509250929050565b614c49614e4b565b8273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161115614c81579192915b506040805160608101825273ffffffffffffffffffffffffffffffffffffffff948516815292909316602083015262ffffff169181019190915290565b6000614cca8383614cee565b90503373ffffffffffffffffffffffffffffffffffffffff82161461202157600080fd5b6000816020015173ffffffffffffffffffffffffffffffffffffffff16826000015173ffffffffffffffffffffffffffffffffffffffff1610614d3057600080fd5b508051602080830151604093840151845173ffffffffffffffffffffffffffffffffffffffff94851681850152939091168385015262ffffff166060808401919091528351808403820181526080840185528051908301207fff0000000000000000000000000000000000000000000000000000000000000060a085015294901b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001660a183015260b58201939093527fe34f199b19b2b4f47f68442619d555527d244f78a3297ea89325f843f87b8b5460d5808301919091528251808303909101815260f5909101909152805191012090565b8082018281101561202157600080fd5b604080518082019091526000808252602082015290565b604080516060810182526000808252602082018190529181019190915290565b8035610b2381615f52565b60008083601f840112614e87578182fd5b50813567ffffffffffffffff811115614e9e578182fd5b6020830191508360208083028501011115613fa057600080fd5b600082601f830112614ec8578081fd5b81356020614edd614ed883615ec8565b615ea4565b8281528181019085830183850287018401881015614ef9578586fd5b855b85811015614f345781356fffffffffffffffffffffffffffffffff81168114614f22578788fd5b84529284019290840190600101614efb565b5090979650505050505050565b80518015158114610b2357600080fd5b600082601f830112614f61578081fd5b8135614f6f614ed882615ee6565b818152846020838601011115614f83578283fd5b816020850160208301379081016020019190915292915050565b80516dffffffffffffffffffffffffffff81168114610b2357600080fd5b805161ffff81168114610b2357600080fd5b803562ffffff81168114610b2357600080fd5b8035610b2381615f83565b600060208284031215614ffc578081fd5b813561500781615f52565b9392505050565b60008060408385031215615020578081fd5b823561502b81615f52565b946020939093013593505050565b60008060006060848603121561504d578081fd5b833561505881615f52565b925060208401359150604084013561506f81615f52565b809150509250925092565b600080600080600060a08688031215615091578283fd5b853561509c81615f52565b94506020860135935060408601356150b381615f52565b92506060860135915060808601356150ca81615f52565b809150509295509295909350565b600080600080608085870312156150ed578182fd5b84356150f881615f52565b93506020850135925060408501359150606085013561511681615f52565b939692955090935050565b60008060008060008060c08789031215615139578384fd5b863561514481615f52565b95506020870135945060408701359350606087013561516281615f95565b9598949750929560808101359460a0909101359350915050565b6000806020838503121561518e578182fd5b823567ffffffffffffffff8111156151a4578283fd5b6151b085828601614e76565b90969095509350505050565b600080600080608085870312156151d1578182fd5b843567ffffffffffffffff808211156151e8578384fd5b818701915087601f8301126151fb578384fd5b8135602061520b614ed883615ec8565b82815281810190858301885b858110156152405761522e8e8684358b0101614f51565b84529284019290840190600101615217565b50909950505088013592505080821115615258578384fd5b5061526587828801614eb8565b93505061527460408601614fcd565b915061528260608601614fe0565b905092959194509250565b60006020828403121561529e578081fd5b61500782614f41565b6000806000604084860312156152bb578081fd5b83359250602084013567ffffffffffffffff8111156152d8578182fd5b6152e486828701614e76565b9497909650939450505050565b600060208284031215615302578081fd5b813567ffffffffffffffff811115615318578182fd5b613d0a84828501614f51565b600080600060608486031215615338578081fd5b833567ffffffffffffffff81111561534e578182fd5b61535a86828701614f51565b93505061536960208501614fcd565b9150604084013561506f81615f83565b60006020828403121561538a578081fd5b813561500781615f74565b600080604083850312156153a7578182fd5b505080516020909101519092909150565b600080600080606085870312156153cd578182fd5b8435935060208501359250604085013567ffffffffffffffff808211156153f2578384fd5b818701915087601f830112615405578384fd5b813581811115615413578485fd5b886020828501011115615424578485fd5b95989497505060200194505050565b600060208284031215615444578081fd5b815167ffffffffffffffff81111561545a578182fd5b8201601f8101841361546a578182fd5b8051615478614ed882615ee6565b81815285602083850101111561548c578384fd5b610e2c826020830160208601615f26565b6000602082840312156154ae578081fd5b813567ffffffffffffffff808211156154c5578283fd5b90830190608082860312156154d8578283fd5b6040516080810181811083821117156154ed57fe5b6040528235828111156154fe578485fd5b61550a87828601614f51565b8252506020830135915061551d82615f52565b816020820152604083013560408201526060830135606082015280935050505092915050565b600060e08284031215615554578081fd5b60405160e0810181811067ffffffffffffffff8211171561557157fe5b60405261557d83614e6b565b815261558b60208401614e6b565b602082015261559c60408401614fcd565b60408201526155ad60608401614e6b565b60608201526080830135608082015260a083013560a08201526155d260c08401614e6b565b60c08201529392505050565b6000602082840312156155ef578081fd5b813567ffffffffffffffff811115615605578182fd5b820160808185031215615007578182fd5b600060e08284031215610836578081fd5b600060a08284031215610836578081fd5b60006101008284031215610836578081fd5b60006020828403121561565b578081fd5b813567ffffffffffffffff80821115615672578283fd5b9083019060408286031215615685578283fd5b60405160408101818110838211171561569a57fe5b6040528235828111156156ab578485fd5b6156b787828601614f51565b825250602083013592506156ca83615f52565b6020810192909252509392505050565b6000806000606084860312156156ee578081fd5b6156f784614f9d565b925061570560208501614f9d565b9150604084015161506f81615f83565b600080600080600080600060e0888a03121561572f578485fd5b875161573a81615f52565b602089015190975061574b81615f74565b955061575960408901614fbb565b945061576760608901614fbb565b935061577560808901614fbb565b925060a088015161578581615f95565b915061579360c08901614f41565b905092959891949750929550565b6000602082840312156157b2578081fd5b61500782614fcd565b6000602082840312156157cc578081fd5b5035919050565b6000602082840312156157e4578081fd5b5051919050565b600080604083850312156157fd578182fd5b82359150602083013561580f81615f52565b809150509250929050565b6000806000806080858703121561582f578182fd5b84359350602085013561584181615f52565b925060408501359150606085013561511681615f52565b60008060006060848603121561586c578081fd5b8335925060208401359150604084013561506f81615f52565b60008060008060006080868803121561589c578283fd5b8535945060208601359350604086013567ffffffffffffffff8111156158c0578384fd5b6158cc88828901614e76565b90945092505060608601356150ca81615f52565b600080600080608085870312156158f5578182fd5b845161590081615f83565b8094505060208501518060060b8114615917578283fd5b604086015190935061592881615f52565b915061528260608601614f41565b73ffffffffffffffffffffffffffffffffffffffff169052565b60008151808452615968816020860160208601615f26565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60020b9052565b62ffffff169052565b606093841b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000908116825260e89390931b7fffffff0000000000000000000000000000000000000000000000000000000000166014820152921b166017820152602b0190565b6000828483379101908152919050565b60008251615a32818460208701615f26565b9190910192915050565b73ffffffffffffffffffffffffffffffffffffffff91909116815260200190565b73ffffffffffffffffffffffffffffffffffffffff92831681529116602082015260400190565b600073ffffffffffffffffffffffffffffffffffffffff8088168352861515602084015285604084015280851660608401525060a06080830152615acb60a0830184615950565b979650505050505050565b73ffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b6000602080830181845280855180835260408601915060408482028701019250838701855b82811015615b6d577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452615b5b858351615950565b94509285019290850190600101615b21565b5092979650505050505050565b6000602082526150076020830184615950565b6020810160058310615b9b57fe5b91905290565b60208082526003908201527f4f4e490000000000000000000000000000000000000000000000000000000000604082015260600190565b60208082526003908201527f4e454f0000000000000000000000000000000000000000000000000000000000604082015260600190565b60208082526012908201527f546f6f206d756368207265717565737465640000000000000000000000000000604082015260600190565b60208082526002908201527f5444000000000000000000000000000000000000000000000000000000000000604082015260600190565b60208082526013908201527f546f6f206c6974746c6520726563656976656400000000000000000000000000604082015260600190565b600060c082019050825182526020830151602083015260408301516040830152606083015160608301526080830151608083015260a083015160a083015292915050565b600061016082019050615d0c828451615936565b6020830151615d1e6020840182615936565b506040830151615d3160408401826159a1565b506060830151615d44606084018261599a565b506080830151615d57608084018261599a565b5060a083015160a083015260c083015160c083015260e083015160e083015261010080840151818401525061012080840151615d9582850182615936565b505061014092830151919092015290565b600060208252825160406020840152615dc26060840182615950565b905073ffffffffffffffffffffffffffffffffffffffff60208501511660408401528091505092915050565b61ffff91909116815260200190565b90815260200190565b600085825284602083015273ffffffffffffffffffffffffffffffffffffffff8416604083015260806060830152613ef86080830184615950565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112615e75578283fd5b83018035915067ffffffffffffffff821115615e8f578283fd5b602001915036819003821315613fa057600080fd5b60405181810167ffffffffffffffff81118282101715615ec057fe5b604052919050565b600067ffffffffffffffff821115615edc57fe5b5060209081020190565b600067ffffffffffffffff821115615efa57fe5b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b60005b83811015615f41578181015183820152602001615f29565b83811115610c555750506000910152565b73ffffffffffffffffffffffffffffffffffffffff8116811461147957600080fd5b8060020b811461147957600080fd5b63ffffffff8116811461147957600080fd5b60ff8116811461147957600080fdfea164736f6c6343000706000a",
}

// SwapRouter02ABI is the input ABI used to generate the binding from.
// Deprecated: Use SwapRouter02MetaData.ABI instead.
var SwapRouter02ABI = SwapRouter02MetaData.ABI

// SwapRouter02Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SwapRouter02MetaData.Bin instead.
var SwapRouter02Bin = SwapRouter02MetaData.Bin

// DeploySwapRouter02 deploys a new Ethereum contract, binding an instance of SwapRouter02 to it.
func DeploySwapRouter02(auth *bind.TransactOpts, backend bind.ContractBackend, _factoryV2 common.Address, factoryV3 common.Address, _positionManager common.Address, _WETH9 common.Address) (common.Address, *types.Transaction, *SwapRouter02, error) {
	parsed, err := SwapRouter02MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SwapRouter02Bin), backend, _factoryV2, factoryV3, _positionManager, _WETH9)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SwapRouter02{SwapRouter02Caller: SwapRouter02Caller{contract: contract}, SwapRouter02Transactor: SwapRouter02Transactor{contract: contract}, SwapRouter02Filterer: SwapRouter02Filterer{contract: contract}}, nil
}

// SwapRouter02 is an auto generated Go binding around an Ethereum contract.
type SwapRouter02 struct {
	SwapRouter02Caller     // Read-only binding to the contract
	SwapRouter02Transactor // Write-only binding to the contract
	SwapRouter02Filterer   // Log filterer for contract events
}

// SwapRouter02Caller is an auto generated read-only Go binding around an Ethereum contract.
type SwapRouter02Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapRouter02Transactor is an auto generated write-only Go binding around an Ethereum contract.
type SwapRouter02Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapRouter02Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SwapRouter02Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapRouter02Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SwapRouter02Session struct {
	Contract     *SwapRouter02     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SwapRouter02CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SwapRouter02CallerSession struct {
	Contract *SwapRouter02Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// SwapRouter02TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SwapRouter02TransactorSession struct {
	Contract     *SwapRouter02Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// SwapRouter02Raw is an auto generated low-level Go binding around an Ethereum contract.
type SwapRouter02Raw struct {
	Contract *SwapRouter02 // Generic contract binding to access the raw methods on
}

// SwapRouter02CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SwapRouter02CallerRaw struct {
	Contract *SwapRouter02Caller // Generic read-only contract binding to access the raw methods on
}

// SwapRouter02TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SwapRouter02TransactorRaw struct {
	Contract *SwapRouter02Transactor // Generic write-only contract binding to access the raw methods on
}

// NewSwapRouter02 creates a new instance of SwapRouter02, bound to a specific deployed contract.
func NewSwapRouter02(address common.Address, backend bind.ContractBackend) (*SwapRouter02, error) {
	contract, err := bindSwapRouter02(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SwapRouter02{SwapRouter02Caller: SwapRouter02Caller{contract: contract}, SwapRouter02Transactor: SwapRouter02Transactor{contract: contract}, SwapRouter02Filterer: SwapRouter02Filterer{contract: contract}}, nil
}

// NewSwapRouter02Caller creates a new read-only instance of SwapRouter02, bound to a specific deployed contract.
func NewSwapRouter02Caller(address common.Address, caller bind.ContractCaller) (*SwapRouter02Caller, error) {
	contract, err := bindSwapRouter02(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SwapRouter02Caller{contract: contract}, nil
}

// NewSwapRouter02Transactor creates a new write-only instance of SwapRouter02, bound to a specific deployed contract.
func NewSwapRouter02Transactor(address common.Address, transactor bind.ContractTransactor) (*SwapRouter02Transactor, error) {
	contract, err := bindSwapRouter02(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SwapRouter02Transactor{contract: contract}, nil
}

// NewSwapRouter02Filterer creates a new log filterer instance of SwapRouter02, bound to a specific deployed contract.
func NewSwapRouter02Filterer(address common.Address, filterer bind.ContractFilterer) (*SwapRouter02Filterer, error) {
	contract, err := bindSwapRouter02(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SwapRouter02Filterer{contract: contract}, nil
}

// bindSwapRouter02 binds a generic wrapper to an already deployed contract.
func bindSwapRouter02(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SwapRouter02MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwapRouter02 *SwapRouter02Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SwapRouter02.Contract.SwapRouter02Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwapRouter02 *SwapRouter02Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SwapRouter02Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwapRouter02 *SwapRouter02Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SwapRouter02Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwapRouter02 *SwapRouter02CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SwapRouter02.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwapRouter02 *SwapRouter02TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapRouter02.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwapRouter02 *SwapRouter02TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwapRouter02.Contract.contract.Transact(opts, method, params...)
}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_SwapRouter02 *SwapRouter02Caller) WETH9(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SwapRouter02.contract.Call(opts, &out, "WETH9")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_SwapRouter02 *SwapRouter02Session) WETH9() (common.Address, error) {
	return _SwapRouter02.Contract.WETH9(&_SwapRouter02.CallOpts)
}

// WETH9 is a free data retrieval call binding the contract method 0x4aa4a4fc.
//
// Solidity: function WETH9() view returns(address)
func (_SwapRouter02 *SwapRouter02CallerSession) WETH9() (common.Address, error) {
	return _SwapRouter02.Contract.WETH9(&_SwapRouter02.CallOpts)
}

// CheckOracleSlippage is a free data retrieval call binding the contract method 0xefdeed8e.
//
// Solidity: function checkOracleSlippage(bytes[] paths, uint128[] amounts, uint24 maximumTickDivergence, uint32 secondsAgo) view returns()
func (_SwapRouter02 *SwapRouter02Caller) CheckOracleSlippage(opts *bind.CallOpts, paths [][]byte, amounts []*big.Int, maximumTickDivergence *big.Int, secondsAgo uint32) error {
	var out []interface{}
	err := _SwapRouter02.contract.Call(opts, &out, "checkOracleSlippage", paths, amounts, maximumTickDivergence, secondsAgo)

	if err != nil {
		return err
	}

	return err

}

// CheckOracleSlippage is a free data retrieval call binding the contract method 0xefdeed8e.
//
// Solidity: function checkOracleSlippage(bytes[] paths, uint128[] amounts, uint24 maximumTickDivergence, uint32 secondsAgo) view returns()
func (_SwapRouter02 *SwapRouter02Session) CheckOracleSlippage(paths [][]byte, amounts []*big.Int, maximumTickDivergence *big.Int, secondsAgo uint32) error {
	return _SwapRouter02.Contract.CheckOracleSlippage(&_SwapRouter02.CallOpts, paths, amounts, maximumTickDivergence, secondsAgo)
}

// CheckOracleSlippage is a free data retrieval call binding the contract method 0xefdeed8e.
//
// Solidity: function checkOracleSlippage(bytes[] paths, uint128[] amounts, uint24 maximumTickDivergence, uint32 secondsAgo) view returns()
func (_SwapRouter02 *SwapRouter02CallerSession) CheckOracleSlippage(paths [][]byte, amounts []*big.Int, maximumTickDivergence *big.Int, secondsAgo uint32) error {
	return _SwapRouter02.Contract.CheckOracleSlippage(&_SwapRouter02.CallOpts, paths, amounts, maximumTickDivergence, secondsAgo)
}

// CheckOracleSlippage0 is a free data retrieval call binding the contract method 0xf25801a7.
//
// Solidity: function checkOracleSlippage(bytes path, uint24 maximumTickDivergence, uint32 secondsAgo) view returns()
func (_SwapRouter02 *SwapRouter02Caller) CheckOracleSlippage0(opts *bind.CallOpts, path []byte, maximumTickDivergence *big.Int, secondsAgo uint32) error {
	var out []interface{}
	err := _SwapRouter02.contract.Call(opts, &out, "checkOracleSlippage0", path, maximumTickDivergence, secondsAgo)

	if err != nil {
		return err
	}

	return err

}

// CheckOracleSlippage0 is a free data retrieval call binding the contract method 0xf25801a7.
//
// Solidity: function checkOracleSlippage(bytes path, uint24 maximumTickDivergence, uint32 secondsAgo) view returns()
func (_SwapRouter02 *SwapRouter02Session) CheckOracleSlippage0(path []byte, maximumTickDivergence *big.Int, secondsAgo uint32) error {
	return _SwapRouter02.Contract.CheckOracleSlippage0(&_SwapRouter02.CallOpts, path, maximumTickDivergence, secondsAgo)
}

// CheckOracleSlippage0 is a free data retrieval call binding the contract method 0xf25801a7.
//
// Solidity: function checkOracleSlippage(bytes path, uint24 maximumTickDivergence, uint32 secondsAgo) view returns()
func (_SwapRouter02 *SwapRouter02CallerSession) CheckOracleSlippage0(path []byte, maximumTickDivergence *big.Int, secondsAgo uint32) error {
	return _SwapRouter02.Contract.CheckOracleSlippage0(&_SwapRouter02.CallOpts, path, maximumTickDivergence, secondsAgo)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_SwapRouter02 *SwapRouter02Caller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SwapRouter02.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_SwapRouter02 *SwapRouter02Session) Factory() (common.Address, error) {
	return _SwapRouter02.Contract.Factory(&_SwapRouter02.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_SwapRouter02 *SwapRouter02CallerSession) Factory() (common.Address, error) {
	return _SwapRouter02.Contract.Factory(&_SwapRouter02.CallOpts)
}

// FactoryV2 is a free data retrieval call binding the contract method 0x68e0d4e1.
//
// Solidity: function factoryV2() view returns(address)
func (_SwapRouter02 *SwapRouter02Caller) FactoryV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SwapRouter02.contract.Call(opts, &out, "factoryV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FactoryV2 is a free data retrieval call binding the contract method 0x68e0d4e1.
//
// Solidity: function factoryV2() view returns(address)
func (_SwapRouter02 *SwapRouter02Session) FactoryV2() (common.Address, error) {
	return _SwapRouter02.Contract.FactoryV2(&_SwapRouter02.CallOpts)
}

// FactoryV2 is a free data retrieval call binding the contract method 0x68e0d4e1.
//
// Solidity: function factoryV2() view returns(address)
func (_SwapRouter02 *SwapRouter02CallerSession) FactoryV2() (common.Address, error) {
	return _SwapRouter02.Contract.FactoryV2(&_SwapRouter02.CallOpts)
}

// PositionManager is a free data retrieval call binding the contract method 0x791b98bc.
//
// Solidity: function positionManager() view returns(address)
func (_SwapRouter02 *SwapRouter02Caller) PositionManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SwapRouter02.contract.Call(opts, &out, "positionManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PositionManager is a free data retrieval call binding the contract method 0x791b98bc.
//
// Solidity: function positionManager() view returns(address)
func (_SwapRouter02 *SwapRouter02Session) PositionManager() (common.Address, error) {
	return _SwapRouter02.Contract.PositionManager(&_SwapRouter02.CallOpts)
}

// PositionManager is a free data retrieval call binding the contract method 0x791b98bc.
//
// Solidity: function positionManager() view returns(address)
func (_SwapRouter02 *SwapRouter02CallerSession) PositionManager() (common.Address, error) {
	return _SwapRouter02.Contract.PositionManager(&_SwapRouter02.CallOpts)
}

// ApproveMax is a paid mutator transaction binding the contract method 0x571ac8b0.
//
// Solidity: function approveMax(address token) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) ApproveMax(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "approveMax", token)
}

// ApproveMax is a paid mutator transaction binding the contract method 0x571ac8b0.
//
// Solidity: function approveMax(address token) payable returns()
func (_SwapRouter02 *SwapRouter02Session) ApproveMax(token common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ApproveMax(&_SwapRouter02.TransactOpts, token)
}

// ApproveMax is a paid mutator transaction binding the contract method 0x571ac8b0.
//
// Solidity: function approveMax(address token) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) ApproveMax(token common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ApproveMax(&_SwapRouter02.TransactOpts, token)
}

// ApproveMaxMinusOne is a paid mutator transaction binding the contract method 0xcab372ce.
//
// Solidity: function approveMaxMinusOne(address token) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) ApproveMaxMinusOne(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "approveMaxMinusOne", token)
}

// ApproveMaxMinusOne is a paid mutator transaction binding the contract method 0xcab372ce.
//
// Solidity: function approveMaxMinusOne(address token) payable returns()
func (_SwapRouter02 *SwapRouter02Session) ApproveMaxMinusOne(token common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ApproveMaxMinusOne(&_SwapRouter02.TransactOpts, token)
}

// ApproveMaxMinusOne is a paid mutator transaction binding the contract method 0xcab372ce.
//
// Solidity: function approveMaxMinusOne(address token) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) ApproveMaxMinusOne(token common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ApproveMaxMinusOne(&_SwapRouter02.TransactOpts, token)
}

// ApproveZeroThenMax is a paid mutator transaction binding the contract method 0x639d71a9.
//
// Solidity: function approveZeroThenMax(address token) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) ApproveZeroThenMax(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "approveZeroThenMax", token)
}

// ApproveZeroThenMax is a paid mutator transaction binding the contract method 0x639d71a9.
//
// Solidity: function approveZeroThenMax(address token) payable returns()
func (_SwapRouter02 *SwapRouter02Session) ApproveZeroThenMax(token common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ApproveZeroThenMax(&_SwapRouter02.TransactOpts, token)
}

// ApproveZeroThenMax is a paid mutator transaction binding the contract method 0x639d71a9.
//
// Solidity: function approveZeroThenMax(address token) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) ApproveZeroThenMax(token common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ApproveZeroThenMax(&_SwapRouter02.TransactOpts, token)
}

// ApproveZeroThenMaxMinusOne is a paid mutator transaction binding the contract method 0xab3fdd50.
//
// Solidity: function approveZeroThenMaxMinusOne(address token) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) ApproveZeroThenMaxMinusOne(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "approveZeroThenMaxMinusOne", token)
}

// ApproveZeroThenMaxMinusOne is a paid mutator transaction binding the contract method 0xab3fdd50.
//
// Solidity: function approveZeroThenMaxMinusOne(address token) payable returns()
func (_SwapRouter02 *SwapRouter02Session) ApproveZeroThenMaxMinusOne(token common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ApproveZeroThenMaxMinusOne(&_SwapRouter02.TransactOpts, token)
}

// ApproveZeroThenMaxMinusOne is a paid mutator transaction binding the contract method 0xab3fdd50.
//
// Solidity: function approveZeroThenMaxMinusOne(address token) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) ApproveZeroThenMaxMinusOne(token common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ApproveZeroThenMaxMinusOne(&_SwapRouter02.TransactOpts, token)
}

// CallPositionManager is a paid mutator transaction binding the contract method 0xb3a2af13.
//
// Solidity: function callPositionManager(bytes data) payable returns(bytes result)
func (_SwapRouter02 *SwapRouter02Transactor) CallPositionManager(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "callPositionManager", data)
}

// CallPositionManager is a paid mutator transaction binding the contract method 0xb3a2af13.
//
// Solidity: function callPositionManager(bytes data) payable returns(bytes result)
func (_SwapRouter02 *SwapRouter02Session) CallPositionManager(data []byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.CallPositionManager(&_SwapRouter02.TransactOpts, data)
}

// CallPositionManager is a paid mutator transaction binding the contract method 0xb3a2af13.
//
// Solidity: function callPositionManager(bytes data) payable returns(bytes result)
func (_SwapRouter02 *SwapRouter02TransactorSession) CallPositionManager(data []byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.CallPositionManager(&_SwapRouter02.TransactOpts, data)
}

// ExactInput is a paid mutator transaction binding the contract method 0xb858183f.
//
// Solidity: function exactInput((bytes,address,uint256,uint256) params) payable returns(uint256 amountOut)
func (_SwapRouter02 *SwapRouter02Transactor) ExactInput(opts *bind.TransactOpts, params IV3SwapRouterExactInputParams) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "exactInput", params)
}

// ExactInput is a paid mutator transaction binding the contract method 0xb858183f.
//
// Solidity: function exactInput((bytes,address,uint256,uint256) params) payable returns(uint256 amountOut)
func (_SwapRouter02 *SwapRouter02Session) ExactInput(params IV3SwapRouterExactInputParams) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ExactInput(&_SwapRouter02.TransactOpts, params)
}

// ExactInput is a paid mutator transaction binding the contract method 0xb858183f.
//
// Solidity: function exactInput((bytes,address,uint256,uint256) params) payable returns(uint256 amountOut)
func (_SwapRouter02 *SwapRouter02TransactorSession) ExactInput(params IV3SwapRouterExactInputParams) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ExactInput(&_SwapRouter02.TransactOpts, params)
}

// ExactInputSingle is a paid mutator transaction binding the contract method 0x04e45aaf.
//
// Solidity: function exactInputSingle((address,address,uint24,address,uint256,uint256,uint160) params) payable returns(uint256 amountOut)
func (_SwapRouter02 *SwapRouter02Transactor) ExactInputSingle(opts *bind.TransactOpts, params IV3SwapRouterExactInputSingleParams) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "exactInputSingle", params)
}

// ExactInputSingle is a paid mutator transaction binding the contract method 0x04e45aaf.
//
// Solidity: function exactInputSingle((address,address,uint24,address,uint256,uint256,uint160) params) payable returns(uint256 amountOut)
func (_SwapRouter02 *SwapRouter02Session) ExactInputSingle(params IV3SwapRouterExactInputSingleParams) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ExactInputSingle(&_SwapRouter02.TransactOpts, params)
}

// ExactInputSingle is a paid mutator transaction binding the contract method 0x04e45aaf.
//
// Solidity: function exactInputSingle((address,address,uint24,address,uint256,uint256,uint160) params) payable returns(uint256 amountOut)
func (_SwapRouter02 *SwapRouter02TransactorSession) ExactInputSingle(params IV3SwapRouterExactInputSingleParams) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ExactInputSingle(&_SwapRouter02.TransactOpts, params)
}

// ExactOutput is a paid mutator transaction binding the contract method 0x09b81346.
//
// Solidity: function exactOutput((bytes,address,uint256,uint256) params) payable returns(uint256 amountIn)
func (_SwapRouter02 *SwapRouter02Transactor) ExactOutput(opts *bind.TransactOpts, params IV3SwapRouterExactOutputParams) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "exactOutput", params)
}

// ExactOutput is a paid mutator transaction binding the contract method 0x09b81346.
//
// Solidity: function exactOutput((bytes,address,uint256,uint256) params) payable returns(uint256 amountIn)
func (_SwapRouter02 *SwapRouter02Session) ExactOutput(params IV3SwapRouterExactOutputParams) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ExactOutput(&_SwapRouter02.TransactOpts, params)
}

// ExactOutput is a paid mutator transaction binding the contract method 0x09b81346.
//
// Solidity: function exactOutput((bytes,address,uint256,uint256) params) payable returns(uint256 amountIn)
func (_SwapRouter02 *SwapRouter02TransactorSession) ExactOutput(params IV3SwapRouterExactOutputParams) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ExactOutput(&_SwapRouter02.TransactOpts, params)
}

// ExactOutputSingle is a paid mutator transaction binding the contract method 0x5023b4df.
//
// Solidity: function exactOutputSingle((address,address,uint24,address,uint256,uint256,uint160) params) payable returns(uint256 amountIn)
func (_SwapRouter02 *SwapRouter02Transactor) ExactOutputSingle(opts *bind.TransactOpts, params IV3SwapRouterExactOutputSingleParams) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "exactOutputSingle", params)
}

// ExactOutputSingle is a paid mutator transaction binding the contract method 0x5023b4df.
//
// Solidity: function exactOutputSingle((address,address,uint24,address,uint256,uint256,uint160) params) payable returns(uint256 amountIn)
func (_SwapRouter02 *SwapRouter02Session) ExactOutputSingle(params IV3SwapRouterExactOutputSingleParams) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ExactOutputSingle(&_SwapRouter02.TransactOpts, params)
}

// ExactOutputSingle is a paid mutator transaction binding the contract method 0x5023b4df.
//
// Solidity: function exactOutputSingle((address,address,uint24,address,uint256,uint256,uint160) params) payable returns(uint256 amountIn)
func (_SwapRouter02 *SwapRouter02TransactorSession) ExactOutputSingle(params IV3SwapRouterExactOutputSingleParams) (*types.Transaction, error) {
	return _SwapRouter02.Contract.ExactOutputSingle(&_SwapRouter02.TransactOpts, params)
}

// GetApprovalType is a paid mutator transaction binding the contract method 0xdee00f35.
//
// Solidity: function getApprovalType(address token, uint256 amount) returns(uint8)
func (_SwapRouter02 *SwapRouter02Transactor) GetApprovalType(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "getApprovalType", token, amount)
}

// GetApprovalType is a paid mutator transaction binding the contract method 0xdee00f35.
//
// Solidity: function getApprovalType(address token, uint256 amount) returns(uint8)
func (_SwapRouter02 *SwapRouter02Session) GetApprovalType(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.Contract.GetApprovalType(&_SwapRouter02.TransactOpts, token, amount)
}

// GetApprovalType is a paid mutator transaction binding the contract method 0xdee00f35.
//
// Solidity: function getApprovalType(address token, uint256 amount) returns(uint8)
func (_SwapRouter02 *SwapRouter02TransactorSession) GetApprovalType(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.Contract.GetApprovalType(&_SwapRouter02.TransactOpts, token, amount)
}

// IncreaseLiquidity is a paid mutator transaction binding the contract method 0xf100b205.
//
// Solidity: function increaseLiquidity((address,address,uint256,uint256,uint256) params) payable returns(bytes result)
func (_SwapRouter02 *SwapRouter02Transactor) IncreaseLiquidity(opts *bind.TransactOpts, params IApproveAndCallIncreaseLiquidityParams) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "increaseLiquidity", params)
}

// IncreaseLiquidity is a paid mutator transaction binding the contract method 0xf100b205.
//
// Solidity: function increaseLiquidity((address,address,uint256,uint256,uint256) params) payable returns(bytes result)
func (_SwapRouter02 *SwapRouter02Session) IncreaseLiquidity(params IApproveAndCallIncreaseLiquidityParams) (*types.Transaction, error) {
	return _SwapRouter02.Contract.IncreaseLiquidity(&_SwapRouter02.TransactOpts, params)
}

// IncreaseLiquidity is a paid mutator transaction binding the contract method 0xf100b205.
//
// Solidity: function increaseLiquidity((address,address,uint256,uint256,uint256) params) payable returns(bytes result)
func (_SwapRouter02 *SwapRouter02TransactorSession) IncreaseLiquidity(params IApproveAndCallIncreaseLiquidityParams) (*types.Transaction, error) {
	return _SwapRouter02.Contract.IncreaseLiquidity(&_SwapRouter02.TransactOpts, params)
}

// Mint is a paid mutator transaction binding the contract method 0x11ed56c9.
//
// Solidity: function mint((address,address,uint24,int24,int24,uint256,uint256,address) params) payable returns(bytes result)
func (_SwapRouter02 *SwapRouter02Transactor) Mint(opts *bind.TransactOpts, params IApproveAndCallMintParams) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "mint", params)
}

// Mint is a paid mutator transaction binding the contract method 0x11ed56c9.
//
// Solidity: function mint((address,address,uint24,int24,int24,uint256,uint256,address) params) payable returns(bytes result)
func (_SwapRouter02 *SwapRouter02Session) Mint(params IApproveAndCallMintParams) (*types.Transaction, error) {
	return _SwapRouter02.Contract.Mint(&_SwapRouter02.TransactOpts, params)
}

// Mint is a paid mutator transaction binding the contract method 0x11ed56c9.
//
// Solidity: function mint((address,address,uint24,int24,int24,uint256,uint256,address) params) payable returns(bytes result)
func (_SwapRouter02 *SwapRouter02TransactorSession) Mint(params IApproveAndCallMintParams) (*types.Transaction, error) {
	return _SwapRouter02.Contract.Mint(&_SwapRouter02.TransactOpts, params)
}

// Multicall is a paid mutator transaction binding the contract method 0x1f0464d1.
//
// Solidity: function multicall(bytes32 previousBlockhash, bytes[] data) payable returns(bytes[])
func (_SwapRouter02 *SwapRouter02Transactor) Multicall(opts *bind.TransactOpts, previousBlockhash [32]byte, data [][]byte) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "multicall", previousBlockhash, data)
}

// Multicall is a paid mutator transaction binding the contract method 0x1f0464d1.
//
// Solidity: function multicall(bytes32 previousBlockhash, bytes[] data) payable returns(bytes[])
func (_SwapRouter02 *SwapRouter02Session) Multicall(previousBlockhash [32]byte, data [][]byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.Multicall(&_SwapRouter02.TransactOpts, previousBlockhash, data)
}

// Multicall is a paid mutator transaction binding the contract method 0x1f0464d1.
//
// Solidity: function multicall(bytes32 previousBlockhash, bytes[] data) payable returns(bytes[])
func (_SwapRouter02 *SwapRouter02TransactorSession) Multicall(previousBlockhash [32]byte, data [][]byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.Multicall(&_SwapRouter02.TransactOpts, previousBlockhash, data)
}

// Multicall0 is a paid mutator transaction binding the contract method 0x5ae401dc.
//
// Solidity: function multicall(uint256 deadline, bytes[] data) payable returns(bytes[])
func (_SwapRouter02 *SwapRouter02Transactor) Multicall0(opts *bind.TransactOpts, deadline *big.Int, data [][]byte) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "multicall0", deadline, data)
}

// Multicall0 is a paid mutator transaction binding the contract method 0x5ae401dc.
//
// Solidity: function multicall(uint256 deadline, bytes[] data) payable returns(bytes[])
func (_SwapRouter02 *SwapRouter02Session) Multicall0(deadline *big.Int, data [][]byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.Multicall0(&_SwapRouter02.TransactOpts, deadline, data)
}

// Multicall0 is a paid mutator transaction binding the contract method 0x5ae401dc.
//
// Solidity: function multicall(uint256 deadline, bytes[] data) payable returns(bytes[])
func (_SwapRouter02 *SwapRouter02TransactorSession) Multicall0(deadline *big.Int, data [][]byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.Multicall0(&_SwapRouter02.TransactOpts, deadline, data)
}

// Multicall1 is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_SwapRouter02 *SwapRouter02Transactor) Multicall1(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "multicall1", data)
}

// Multicall1 is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_SwapRouter02 *SwapRouter02Session) Multicall1(data [][]byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.Multicall1(&_SwapRouter02.TransactOpts, data)
}

// Multicall1 is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_SwapRouter02 *SwapRouter02TransactorSession) Multicall1(data [][]byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.Multicall1(&_SwapRouter02.TransactOpts, data)
}

// Pull is a paid mutator transaction binding the contract method 0xf2d5d56b.
//
// Solidity: function pull(address token, uint256 value) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) Pull(opts *bind.TransactOpts, token common.Address, value *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "pull", token, value)
}

// Pull is a paid mutator transaction binding the contract method 0xf2d5d56b.
//
// Solidity: function pull(address token, uint256 value) payable returns()
func (_SwapRouter02 *SwapRouter02Session) Pull(token common.Address, value *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.Contract.Pull(&_SwapRouter02.TransactOpts, token, value)
}

// Pull is a paid mutator transaction binding the contract method 0xf2d5d56b.
//
// Solidity: function pull(address token, uint256 value) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) Pull(token common.Address, value *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.Contract.Pull(&_SwapRouter02.TransactOpts, token, value)
}

// RefundETH is a paid mutator transaction binding the contract method 0x12210e8a.
//
// Solidity: function refundETH() payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) RefundETH(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "refundETH")
}

// RefundETH is a paid mutator transaction binding the contract method 0x12210e8a.
//
// Solidity: function refundETH() payable returns()
func (_SwapRouter02 *SwapRouter02Session) RefundETH() (*types.Transaction, error) {
	return _SwapRouter02.Contract.RefundETH(&_SwapRouter02.TransactOpts)
}

// RefundETH is a paid mutator transaction binding the contract method 0x12210e8a.
//
// Solidity: function refundETH() payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) RefundETH() (*types.Transaction, error) {
	return _SwapRouter02.Contract.RefundETH(&_SwapRouter02.TransactOpts)
}

// SelfPermit is a paid mutator transaction binding the contract method 0xf3995c67.
//
// Solidity: function selfPermit(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) SelfPermit(opts *bind.TransactOpts, token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "selfPermit", token, value, deadline, v, r, s)
}

// SelfPermit is a paid mutator transaction binding the contract method 0xf3995c67.
//
// Solidity: function selfPermit(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_SwapRouter02 *SwapRouter02Session) SelfPermit(token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SelfPermit(&_SwapRouter02.TransactOpts, token, value, deadline, v, r, s)
}

// SelfPermit is a paid mutator transaction binding the contract method 0xf3995c67.
//
// Solidity: function selfPermit(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) SelfPermit(token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SelfPermit(&_SwapRouter02.TransactOpts, token, value, deadline, v, r, s)
}

// SelfPermitAllowed is a paid mutator transaction binding the contract method 0x4659a494.
//
// Solidity: function selfPermitAllowed(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) SelfPermitAllowed(opts *bind.TransactOpts, token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "selfPermitAllowed", token, nonce, expiry, v, r, s)
}

// SelfPermitAllowed is a paid mutator transaction binding the contract method 0x4659a494.
//
// Solidity: function selfPermitAllowed(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_SwapRouter02 *SwapRouter02Session) SelfPermitAllowed(token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SelfPermitAllowed(&_SwapRouter02.TransactOpts, token, nonce, expiry, v, r, s)
}

// SelfPermitAllowed is a paid mutator transaction binding the contract method 0x4659a494.
//
// Solidity: function selfPermitAllowed(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) SelfPermitAllowed(token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SelfPermitAllowed(&_SwapRouter02.TransactOpts, token, nonce, expiry, v, r, s)
}

// SelfPermitAllowedIfNecessary is a paid mutator transaction binding the contract method 0xa4a78f0c.
//
// Solidity: function selfPermitAllowedIfNecessary(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) SelfPermitAllowedIfNecessary(opts *bind.TransactOpts, token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "selfPermitAllowedIfNecessary", token, nonce, expiry, v, r, s)
}

// SelfPermitAllowedIfNecessary is a paid mutator transaction binding the contract method 0xa4a78f0c.
//
// Solidity: function selfPermitAllowedIfNecessary(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_SwapRouter02 *SwapRouter02Session) SelfPermitAllowedIfNecessary(token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SelfPermitAllowedIfNecessary(&_SwapRouter02.TransactOpts, token, nonce, expiry, v, r, s)
}

// SelfPermitAllowedIfNecessary is a paid mutator transaction binding the contract method 0xa4a78f0c.
//
// Solidity: function selfPermitAllowedIfNecessary(address token, uint256 nonce, uint256 expiry, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) SelfPermitAllowedIfNecessary(token common.Address, nonce *big.Int, expiry *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SelfPermitAllowedIfNecessary(&_SwapRouter02.TransactOpts, token, nonce, expiry, v, r, s)
}

// SelfPermitIfNecessary is a paid mutator transaction binding the contract method 0xc2e3140a.
//
// Solidity: function selfPermitIfNecessary(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) SelfPermitIfNecessary(opts *bind.TransactOpts, token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "selfPermitIfNecessary", token, value, deadline, v, r, s)
}

// SelfPermitIfNecessary is a paid mutator transaction binding the contract method 0xc2e3140a.
//
// Solidity: function selfPermitIfNecessary(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_SwapRouter02 *SwapRouter02Session) SelfPermitIfNecessary(token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SelfPermitIfNecessary(&_SwapRouter02.TransactOpts, token, value, deadline, v, r, s)
}

// SelfPermitIfNecessary is a paid mutator transaction binding the contract method 0xc2e3140a.
//
// Solidity: function selfPermitIfNecessary(address token, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) SelfPermitIfNecessary(token common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SelfPermitIfNecessary(&_SwapRouter02.TransactOpts, token, value, deadline, v, r, s)
}

// SwapExactTokensForTokens is a paid mutator transaction binding the contract method 0x472b43f3.
//
// Solidity: function swapExactTokensForTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to) payable returns(uint256 amountOut)
func (_SwapRouter02 *SwapRouter02Transactor) SwapExactTokensForTokens(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "swapExactTokensForTokens", amountIn, amountOutMin, path, to)
}

// SwapExactTokensForTokens is a paid mutator transaction binding the contract method 0x472b43f3.
//
// Solidity: function swapExactTokensForTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to) payable returns(uint256 amountOut)
func (_SwapRouter02 *SwapRouter02Session) SwapExactTokensForTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SwapExactTokensForTokens(&_SwapRouter02.TransactOpts, amountIn, amountOutMin, path, to)
}

// SwapExactTokensForTokens is a paid mutator transaction binding the contract method 0x472b43f3.
//
// Solidity: function swapExactTokensForTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to) payable returns(uint256 amountOut)
func (_SwapRouter02 *SwapRouter02TransactorSession) SwapExactTokensForTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SwapExactTokensForTokens(&_SwapRouter02.TransactOpts, amountIn, amountOutMin, path, to)
}

// SwapTokensForExactTokens is a paid mutator transaction binding the contract method 0x42712a67.
//
// Solidity: function swapTokensForExactTokens(uint256 amountOut, uint256 amountInMax, address[] path, address to) payable returns(uint256 amountIn)
func (_SwapRouter02 *SwapRouter02Transactor) SwapTokensForExactTokens(opts *bind.TransactOpts, amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "swapTokensForExactTokens", amountOut, amountInMax, path, to)
}

// SwapTokensForExactTokens is a paid mutator transaction binding the contract method 0x42712a67.
//
// Solidity: function swapTokensForExactTokens(uint256 amountOut, uint256 amountInMax, address[] path, address to) payable returns(uint256 amountIn)
func (_SwapRouter02 *SwapRouter02Session) SwapTokensForExactTokens(amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SwapTokensForExactTokens(&_SwapRouter02.TransactOpts, amountOut, amountInMax, path, to)
}

// SwapTokensForExactTokens is a paid mutator transaction binding the contract method 0x42712a67.
//
// Solidity: function swapTokensForExactTokens(uint256 amountOut, uint256 amountInMax, address[] path, address to) payable returns(uint256 amountIn)
func (_SwapRouter02 *SwapRouter02TransactorSession) SwapTokensForExactTokens(amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SwapTokensForExactTokens(&_SwapRouter02.TransactOpts, amountOut, amountInMax, path, to)
}

// SweepToken is a paid mutator transaction binding the contract method 0xdf2ab5bb.
//
// Solidity: function sweepToken(address token, uint256 amountMinimum, address recipient) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) SweepToken(opts *bind.TransactOpts, token common.Address, amountMinimum *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "sweepToken", token, amountMinimum, recipient)
}

// SweepToken is a paid mutator transaction binding the contract method 0xdf2ab5bb.
//
// Solidity: function sweepToken(address token, uint256 amountMinimum, address recipient) payable returns()
func (_SwapRouter02 *SwapRouter02Session) SweepToken(token common.Address, amountMinimum *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SweepToken(&_SwapRouter02.TransactOpts, token, amountMinimum, recipient)
}

// SweepToken is a paid mutator transaction binding the contract method 0xdf2ab5bb.
//
// Solidity: function sweepToken(address token, uint256 amountMinimum, address recipient) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) SweepToken(token common.Address, amountMinimum *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SweepToken(&_SwapRouter02.TransactOpts, token, amountMinimum, recipient)
}

// SweepToken0 is a paid mutator transaction binding the contract method 0xe90a182f.
//
// Solidity: function sweepToken(address token, uint256 amountMinimum) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) SweepToken0(opts *bind.TransactOpts, token common.Address, amountMinimum *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "sweepToken0", token, amountMinimum)
}

// SweepToken0 is a paid mutator transaction binding the contract method 0xe90a182f.
//
// Solidity: function sweepToken(address token, uint256 amountMinimum) payable returns()
func (_SwapRouter02 *SwapRouter02Session) SweepToken0(token common.Address, amountMinimum *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SweepToken0(&_SwapRouter02.TransactOpts, token, amountMinimum)
}

// SweepToken0 is a paid mutator transaction binding the contract method 0xe90a182f.
//
// Solidity: function sweepToken(address token, uint256 amountMinimum) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) SweepToken0(token common.Address, amountMinimum *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SweepToken0(&_SwapRouter02.TransactOpts, token, amountMinimum)
}

// SweepTokenWithFee is a paid mutator transaction binding the contract method 0x3068c554.
//
// Solidity: function sweepTokenWithFee(address token, uint256 amountMinimum, uint256 feeBips, address feeRecipient) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) SweepTokenWithFee(opts *bind.TransactOpts, token common.Address, amountMinimum *big.Int, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "sweepTokenWithFee", token, amountMinimum, feeBips, feeRecipient)
}

// SweepTokenWithFee is a paid mutator transaction binding the contract method 0x3068c554.
//
// Solidity: function sweepTokenWithFee(address token, uint256 amountMinimum, uint256 feeBips, address feeRecipient) payable returns()
func (_SwapRouter02 *SwapRouter02Session) SweepTokenWithFee(token common.Address, amountMinimum *big.Int, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SweepTokenWithFee(&_SwapRouter02.TransactOpts, token, amountMinimum, feeBips, feeRecipient)
}

// SweepTokenWithFee is a paid mutator transaction binding the contract method 0x3068c554.
//
// Solidity: function sweepTokenWithFee(address token, uint256 amountMinimum, uint256 feeBips, address feeRecipient) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) SweepTokenWithFee(token common.Address, amountMinimum *big.Int, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SweepTokenWithFee(&_SwapRouter02.TransactOpts, token, amountMinimum, feeBips, feeRecipient)
}

// SweepTokenWithFee0 is a paid mutator transaction binding the contract method 0xe0e189a0.
//
// Solidity: function sweepTokenWithFee(address token, uint256 amountMinimum, address recipient, uint256 feeBips, address feeRecipient) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) SweepTokenWithFee0(opts *bind.TransactOpts, token common.Address, amountMinimum *big.Int, recipient common.Address, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "sweepTokenWithFee0", token, amountMinimum, recipient, feeBips, feeRecipient)
}

// SweepTokenWithFee0 is a paid mutator transaction binding the contract method 0xe0e189a0.
//
// Solidity: function sweepTokenWithFee(address token, uint256 amountMinimum, address recipient, uint256 feeBips, address feeRecipient) payable returns()
func (_SwapRouter02 *SwapRouter02Session) SweepTokenWithFee0(token common.Address, amountMinimum *big.Int, recipient common.Address, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SweepTokenWithFee0(&_SwapRouter02.TransactOpts, token, amountMinimum, recipient, feeBips, feeRecipient)
}

// SweepTokenWithFee0 is a paid mutator transaction binding the contract method 0xe0e189a0.
//
// Solidity: function sweepTokenWithFee(address token, uint256 amountMinimum, address recipient, uint256 feeBips, address feeRecipient) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) SweepTokenWithFee0(token common.Address, amountMinimum *big.Int, recipient common.Address, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.SweepTokenWithFee0(&_SwapRouter02.TransactOpts, token, amountMinimum, recipient, feeBips, feeRecipient)
}

// UniswapV3SwapCallback is a paid mutator transaction binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes _data) returns()
func (_SwapRouter02 *SwapRouter02Transactor) UniswapV3SwapCallback(opts *bind.TransactOpts, amount0Delta *big.Int, amount1Delta *big.Int, _data []byte) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "uniswapV3SwapCallback", amount0Delta, amount1Delta, _data)
}

// UniswapV3SwapCallback is a paid mutator transaction binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes _data) returns()
func (_SwapRouter02 *SwapRouter02Session) UniswapV3SwapCallback(amount0Delta *big.Int, amount1Delta *big.Int, _data []byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.UniswapV3SwapCallback(&_SwapRouter02.TransactOpts, amount0Delta, amount1Delta, _data)
}

// UniswapV3SwapCallback is a paid mutator transaction binding the contract method 0xfa461e33.
//
// Solidity: function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes _data) returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) UniswapV3SwapCallback(amount0Delta *big.Int, amount1Delta *big.Int, _data []byte) (*types.Transaction, error) {
	return _SwapRouter02.Contract.UniswapV3SwapCallback(&_SwapRouter02.TransactOpts, amount0Delta, amount1Delta, _data)
}

// UnwrapWETH9 is a paid mutator transaction binding the contract method 0x49404b7c.
//
// Solidity: function unwrapWETH9(uint256 amountMinimum, address recipient) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) UnwrapWETH9(opts *bind.TransactOpts, amountMinimum *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "unwrapWETH9", amountMinimum, recipient)
}

// UnwrapWETH9 is a paid mutator transaction binding the contract method 0x49404b7c.
//
// Solidity: function unwrapWETH9(uint256 amountMinimum, address recipient) payable returns()
func (_SwapRouter02 *SwapRouter02Session) UnwrapWETH9(amountMinimum *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.UnwrapWETH9(&_SwapRouter02.TransactOpts, amountMinimum, recipient)
}

// UnwrapWETH9 is a paid mutator transaction binding the contract method 0x49404b7c.
//
// Solidity: function unwrapWETH9(uint256 amountMinimum, address recipient) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) UnwrapWETH9(amountMinimum *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.UnwrapWETH9(&_SwapRouter02.TransactOpts, amountMinimum, recipient)
}

// UnwrapWETH90 is a paid mutator transaction binding the contract method 0x49616997.
//
// Solidity: function unwrapWETH9(uint256 amountMinimum) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) UnwrapWETH90(opts *bind.TransactOpts, amountMinimum *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "unwrapWETH90", amountMinimum)
}

// UnwrapWETH90 is a paid mutator transaction binding the contract method 0x49616997.
//
// Solidity: function unwrapWETH9(uint256 amountMinimum) payable returns()
func (_SwapRouter02 *SwapRouter02Session) UnwrapWETH90(amountMinimum *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.Contract.UnwrapWETH90(&_SwapRouter02.TransactOpts, amountMinimum)
}

// UnwrapWETH90 is a paid mutator transaction binding the contract method 0x49616997.
//
// Solidity: function unwrapWETH9(uint256 amountMinimum) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) UnwrapWETH90(amountMinimum *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.Contract.UnwrapWETH90(&_SwapRouter02.TransactOpts, amountMinimum)
}

// UnwrapWETH9WithFee is a paid mutator transaction binding the contract method 0x9b2c0a37.
//
// Solidity: function unwrapWETH9WithFee(uint256 amountMinimum, address recipient, uint256 feeBips, address feeRecipient) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) UnwrapWETH9WithFee(opts *bind.TransactOpts, amountMinimum *big.Int, recipient common.Address, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "unwrapWETH9WithFee", amountMinimum, recipient, feeBips, feeRecipient)
}

// UnwrapWETH9WithFee is a paid mutator transaction binding the contract method 0x9b2c0a37.
//
// Solidity: function unwrapWETH9WithFee(uint256 amountMinimum, address recipient, uint256 feeBips, address feeRecipient) payable returns()
func (_SwapRouter02 *SwapRouter02Session) UnwrapWETH9WithFee(amountMinimum *big.Int, recipient common.Address, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.UnwrapWETH9WithFee(&_SwapRouter02.TransactOpts, amountMinimum, recipient, feeBips, feeRecipient)
}

// UnwrapWETH9WithFee is a paid mutator transaction binding the contract method 0x9b2c0a37.
//
// Solidity: function unwrapWETH9WithFee(uint256 amountMinimum, address recipient, uint256 feeBips, address feeRecipient) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) UnwrapWETH9WithFee(amountMinimum *big.Int, recipient common.Address, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.UnwrapWETH9WithFee(&_SwapRouter02.TransactOpts, amountMinimum, recipient, feeBips, feeRecipient)
}

// UnwrapWETH9WithFee0 is a paid mutator transaction binding the contract method 0xd4ef38de.
//
// Solidity: function unwrapWETH9WithFee(uint256 amountMinimum, uint256 feeBips, address feeRecipient) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) UnwrapWETH9WithFee0(opts *bind.TransactOpts, amountMinimum *big.Int, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "unwrapWETH9WithFee0", amountMinimum, feeBips, feeRecipient)
}

// UnwrapWETH9WithFee0 is a paid mutator transaction binding the contract method 0xd4ef38de.
//
// Solidity: function unwrapWETH9WithFee(uint256 amountMinimum, uint256 feeBips, address feeRecipient) payable returns()
func (_SwapRouter02 *SwapRouter02Session) UnwrapWETH9WithFee0(amountMinimum *big.Int, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.UnwrapWETH9WithFee0(&_SwapRouter02.TransactOpts, amountMinimum, feeBips, feeRecipient)
}

// UnwrapWETH9WithFee0 is a paid mutator transaction binding the contract method 0xd4ef38de.
//
// Solidity: function unwrapWETH9WithFee(uint256 amountMinimum, uint256 feeBips, address feeRecipient) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) UnwrapWETH9WithFee0(amountMinimum *big.Int, feeBips *big.Int, feeRecipient common.Address) (*types.Transaction, error) {
	return _SwapRouter02.Contract.UnwrapWETH9WithFee0(&_SwapRouter02.TransactOpts, amountMinimum, feeBips, feeRecipient)
}

// WrapETH is a paid mutator transaction binding the contract method 0x1c58db4f.
//
// Solidity: function wrapETH(uint256 value) payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) WrapETH(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.contract.Transact(opts, "wrapETH", value)
}

// WrapETH is a paid mutator transaction binding the contract method 0x1c58db4f.
//
// Solidity: function wrapETH(uint256 value) payable returns()
func (_SwapRouter02 *SwapRouter02Session) WrapETH(value *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.Contract.WrapETH(&_SwapRouter02.TransactOpts, value)
}

// WrapETH is a paid mutator transaction binding the contract method 0x1c58db4f.
//
// Solidity: function wrapETH(uint256 value) payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) WrapETH(value *big.Int) (*types.Transaction, error) {
	return _SwapRouter02.Contract.WrapETH(&_SwapRouter02.TransactOpts, value)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SwapRouter02 *SwapRouter02Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapRouter02.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SwapRouter02 *SwapRouter02Session) Receive() (*types.Transaction, error) {
	return _SwapRouter02.Contract.Receive(&_SwapRouter02.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_SwapRouter02 *SwapRouter02TransactorSession) Receive() (*types.Transaction, error) {
	return _SwapRouter02.Contract.Receive(&_SwapRouter02.TransactOpts)
}

