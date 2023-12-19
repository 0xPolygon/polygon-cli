// Code generated by "stringer -type=loadTestMode"; DO NOT EDIT.

package loadtest

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[loadTestModeTransaction-0]
	_ = x[loadTestModeDeploy-1]
	_ = x[loadTestModeCall-2]
	_ = x[loadTestModeFunction-3]
	_ = x[loadTestModeInc-4]
	_ = x[loadTestModeStore-5]
	_ = x[loadTestModeERC20-6]
	_ = x[loadTestModeERC721-7]
	_ = x[loadTestModePrecompiledContracts-8]
	_ = x[loadTestModePrecompiledContract-9]
	_ = x[loadTestModeRandom-10]
	_ = x[loadTestModeRecall-11]
	_ = x[loadTestModeRPC-12]
	_ = x[loadTestModeContractCall-13]
	_ = x[loadTestModeInscription-14]
	_ = x[loadTestModeUniswapV3-15]
}

const _loadTestMode_name = "loadTestModeTransactionloadTestModeDeployloadTestModeCallloadTestModeFunctionloadTestModeIncloadTestModeStoreloadTestModeERC20loadTestModeERC721loadTestModePrecompiledContractsloadTestModePrecompiledContractloadTestModeRandomloadTestModeRecallloadTestModeRPCloadTestModeContractCallloadTestModeInscriptionloadTestModeUniswapV3"

var _loadTestMode_index = [...]uint16{0, 23, 41, 57, 77, 92, 109, 126, 144, 176, 207, 225, 243, 258, 282, 305, 326}

func (i loadTestMode) String() string {
	if i < 0 || i >= loadTestMode(len(_loadTestMode_index)-1) {
		return "loadTestMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _loadTestMode_name[_loadTestMode_index[i]:_loadTestMode_index[i+1]]
}
