package contracts

import (
	_ "embed"
	"encoding/hex"
)

// solc --version
// solc, the solidity compiler commandline interface
// Version: 0.8.15+commit.e14f2714.Darwin.appleclang
// solc loadtest.sol --bin --abi -o .
// ~/code/go-ethereum/build/bin/abigen --abi LoadTester.abi --pkg contracts --type LoadTester --out loadtester.go

//go:embed LoadTester.bin
var RawLoadTesterBin string

//go:embed LoadTester.abi
var RawLoadTesterABI string

func GetLoadTesterBytes() ([]byte, error) {
	return hex.DecodeString(RawLoadTesterBin)
}
