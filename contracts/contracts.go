package contracts

import (
	_ "embed"
	"encoding/hex"
)

// solc --version
// solc, the solidity compiler commandline interface
// Version: 0.8.15+commit.e14f2714.Darwin.appleclang
// solc loadtest.sol --bin --abi -o .

//go:embed LoadTester.bin
var LoadTesterBin string

//go:embed LoadTester.abi
var LoadTesterABI string

func GetLoadTesterBytes() ([]byte, error) {
	return hex.DecodeString(LoadTesterBin)
}
