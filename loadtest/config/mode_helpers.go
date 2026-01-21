package config

import (
	"fmt"
	"slices"
	"strings"
)

// ParseMode converts a mode string to a Mode enum.
// Note: "R" (capital) is the alias for recall mode, while "r" is for random mode.
func ParseMode(modeStr string) (Mode, error) {
	// Handle case-sensitive aliases first
	switch modeStr {
	case "R":
		return ModeRecall, nil
	case "r":
		return ModeRandom, nil
	}

	// Handle case-insensitive mode names
	switch strings.ToLower(modeStr) {
	case "2", "erc20":
		return ModeERC20, nil
	case "7", "erc721":
		return ModeERC721, nil
	case "b", "blob":
		return ModeBlob, nil
	case "cc", "contract-call":
		return ModeContractCall, nil
	case "d", "deploy":
		return ModeDeploy, nil
	case "inc", "increment":
		return ModeIncrement, nil
	case "random":
		return ModeRandom, nil
	case "recall":
		return ModeRecall, nil
	case "rpc":
		return ModeRPC, nil
	case "s", "store":
		return ModeStore, nil
	case "t", "transaction":
		return ModeTransaction, nil
	case "v3", "uniswapv3":
		return ModeUniswapV3, nil
	default:
		return 0, fmt.Errorf("unrecognized load test mode: %s", modeStr)
	}
}

// RequiresLoadTestContract returns true if the mode requires the LoadTester contract.
func RequiresLoadTestContract(m Mode) bool {
	return m == ModeIncrement || m == ModeRandom || m == ModeStore
}

// RequiresERC20 returns true if the mode requires an ERC20 contract.
func RequiresERC20(m Mode) bool {
	return m == ModeERC20 || m == ModeRandom || m == ModeRPC
}

// RequiresERC721 returns true if the mode requires an ERC721 contract.
func RequiresERC721(m Mode) bool {
	return m == ModeERC721 || m == ModeRandom || m == ModeRPC
}

// HasMode returns true if the given mode is in the slice.
func HasMode(mode Mode, m []Mode) bool {
	return slices.Contains(m, mode)
}

// HasUniqueModes returns true if all modes in the slice are unique.
func HasUniqueModes(m []Mode) bool {
	seen := make(map[Mode]bool, len(m))
	for _, mode := range m {
		if seen[mode] {
			return false
		}
		seen[mode] = true
	}
	return true
}

// AnyRequiresLoadTestContract returns true if any mode requires the LoadTester contract.
func AnyRequiresLoadTestContract(m []Mode) bool {
	return slices.ContainsFunc(m, RequiresLoadTestContract)
}

// AnyRequiresERC20 returns true if any mode requires an ERC20 contract.
func AnyRequiresERC20(m []Mode) bool {
	return slices.ContainsFunc(m, RequiresERC20)
}

// AnyRequiresERC721 returns true if any mode requires an ERC721 contract.
func AnyRequiresERC721(m []Mode) bool {
	return slices.ContainsFunc(m, RequiresERC721)
}
