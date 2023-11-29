package tester

import (
	"context"
	_ "embed"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maticnetwork/polygon-cli/util"
	"github.com/rs/zerolog/log"
)

// solc --version
// solc, the solidity compiler commandline interface
// Version: 0.8.15+commit.e14f2714.Darwin.appleclang
// solc LoadTester.sol --bin --abi -o . --overwrite
// From within `polygon-cli/contracts/loadtester` directory:
// ~/code/go-ethereum/build/bin/abigen --abi LoadTester.abi --pkg contracts --type LoadTester --bin LoadTester.bin --out ../loadtester.go

//go:embed LoadTester.bin
var RawLoadTesterBin string

//go:embed LoadTester.abi
var RawLoadTesterABI string

var randSrc *rand.Rand

func GetLoadTesterBytes() ([]byte, error) {
	return hex.DecodeString(RawLoadTesterBin)
}

func DeployConformanceContract(ctx context.Context, c *ethclient.Client, tops *bind.TransactOpts, cops *bind.CallOpts) (conformanceContractAddr ethcommon.Address, conformanceContract *ConformanceTester, err error) {
	conformanceContractAddr, _, _, err = DeployConformanceTester(tops, c, "ConformanceTesterContractName")
	if err != nil {
		log.Error().Err(err).Msg("Unable to deploy ConformanceTester contract")
		return
	}
	log.Info().Interface("conformanceContractAddr", conformanceContractAddr).Msg("Conformance contract deployed")

	conformanceContract, err = NewConformanceTester(conformanceContractAddr, c)
	if err != nil {
		log.Error().Err(err).Msg("Unable to instantiate new conformance contract")
		return
	}
	log.Trace().Msg("Conformance contract instantiated")

	err = util.BlockUntilSuccessful(ctx, c, func() error {
		_, err := conformanceContract.Name(cops)
		return err
	})

	return
}

func CallLoadTestFunctionByOpCode(shortCode uint64, lt *LoadTester, opts *bind.TransactOpts, iterations uint64) (*ethtypes.Transaction, error) {
	x := new(big.Int).SetUint64(iterations)
	var longCode = 0xDEADBEEF0000 | shortCode

	switch longCode {
	case 0xDEADBEEF0001:
		log.Trace().Str("method", "TestADD").Msg("Executing contract method")
		return lt.TestADD(opts, x)

	case 0xDEADBEEF0002:
		log.Trace().Str("method", "TestMUL").Msg("Executing contract method")
		return lt.TestMUL(opts, x)

	case 0xDEADBEEF0003:
		log.Trace().Str("method", "TestSUB").Msg("Executing contract method")
		return lt.TestSUB(opts, x)

	case 0xDEADBEEF0004:
		log.Trace().Str("method", "TestDIV").Msg("Executing contract method")
		return lt.TestDIV(opts, x)

	case 0xDEADBEEF0005:
		log.Trace().Str("method", "TestSDIV").Msg("Executing contract method")
		return lt.TestSDIV(opts, x)

	case 0xDEADBEEF0006:
		log.Trace().Str("method", "TestMOD").Msg("Executing contract method")
		return lt.TestMOD(opts, x)

	case 0xDEADBEEF0007:
		log.Trace().Str("method", "TestSMOD").Msg("Executing contract method")
		return lt.TestSMOD(opts, x)

	case 0xDEADBEEF0008:
		log.Trace().Str("method", "TestADDMOD").Msg("Executing contract method")
		return lt.TestADDMOD(opts, x)

	case 0xDEADBEEF0009:
		log.Trace().Str("method", "TestMULMOD").Msg("Executing contract method")
		return lt.TestMULMOD(opts, x)

	case 0xDEADBEEF000A:
		log.Trace().Str("method", "TestEXP").Msg("Executing contract method")
		return lt.TestEXP(opts, x)

	case 0xDEADBEEF000B:
		log.Trace().Str("method", "TestSIGNEXTEND").Msg("Executing contract method")
		return lt.TestSIGNEXTEND(opts, x)

	case 0xDEADBEEF0010:
		log.Trace().Str("method", "TestLT").Msg("Executing contract method")
		return lt.TestLT(opts, x)

	case 0xDEADBEEF0011:
		log.Trace().Str("method", "TestGT").Msg("Executing contract method")
		return lt.TestGT(opts, x)

	case 0xDEADBEEF0012:
		log.Trace().Str("method", "TestSLT").Msg("Executing contract method")
		return lt.TestSLT(opts, x)

	case 0xDEADBEEF0013:
		log.Trace().Str("method", "TestSGT").Msg("Executing contract method")
		return lt.TestSGT(opts, x)

	case 0xDEADBEEF0014:
		log.Trace().Str("method", "TestEQ").Msg("Executing contract method")
		return lt.TestEQ(opts, x)

	case 0xDEADBEEF0015:
		log.Trace().Str("method", "TestISZERO").Msg("Executing contract method")
		return lt.TestISZERO(opts, x)

	case 0xDEADBEEF0016:
		log.Trace().Str("method", "TestAND").Msg("Executing contract method")
		return lt.TestAND(opts, x)

	case 0xDEADBEEF0017:
		log.Trace().Str("method", "TestOR").Msg("Executing contract method")
		return lt.TestOR(opts, x)

	case 0xDEADBEEF0018:
		log.Trace().Str("method", "TestXOR").Msg("Executing contract method")
		return lt.TestXOR(opts, x)

	case 0xDEADBEEF0019:
		log.Trace().Str("method", "TestNOT").Msg("Executing contract method")
		return lt.TestNOT(opts, x)

	case 0xDEADBEEF001A:
		log.Trace().Str("method", "TestBYTE").Msg("Executing contract method")
		return lt.TestBYTE(opts, x)

	case 0xDEADBEEF001B:
		log.Trace().Str("method", "TestSHL").Msg("Executing contract method")
		return lt.TestSHL(opts, x)

	case 0xDEADBEEF001C:
		log.Trace().Str("method", "TestSHR").Msg("Executing contract method")
		return lt.TestSHR(opts, x)

	case 0xDEADBEEF001D:
		log.Trace().Str("method", "TestSAR").Msg("Executing contract method")
		return lt.TestSAR(opts, x)

	case 0xDEADBEEF0020:
		log.Trace().Str("method", "TestSHA3").Msg("Executing contract method")
		return lt.TestSHA3(opts, x)

	case 0xDEADBEEF0030:
		log.Trace().Str("method", "TestADDRESS").Msg("Executing contract method")
		return lt.TestADDRESS(opts, x)

	case 0xDEADBEEF0031:
		log.Trace().Str("method", "TestBALANCE").Msg("Executing contract method")
		return lt.TestBALANCE(opts, x)

	case 0xDEADBEEF0032:
		log.Trace().Str("method", "TestORIGIN").Msg("Executing contract method")
		return lt.TestORIGIN(opts, x)

	case 0xDEADBEEF0033:
		log.Trace().Str("method", "TestCALLER").Msg("Executing contract method")
		return lt.TestCALLER(opts, x)

	case 0xDEADBEEF0034:
		log.Trace().Str("method", "TestCALLVALUE").Msg("Executing contract method")
		return lt.TestCALLVALUE(opts, x)

	case 0xDEADBEEF0035:
		log.Trace().Str("method", "TestCALLDATALOAD").Msg("Executing contract method")
		return lt.TestCALLDATALOAD(opts, x)

	case 0xDEADBEEF0036:
		log.Trace().Str("method", "TestCALLDATASIZE").Msg("Executing contract method")
		return lt.TestCALLDATASIZE(opts, x)

	case 0xDEADBEEF0037:
		log.Trace().Str("method", "TestCALLDATACOPY").Msg("Executing contract method")
		return lt.TestCALLDATACOPY(opts, x)

	case 0xDEADBEEF0038:
		log.Trace().Str("method", "TestCODESIZE").Msg("Executing contract method")
		return lt.TestCODESIZE(opts, x)

	case 0xDEADBEEF0039:
		log.Trace().Str("method", "TestCODECOPY").Msg("Executing contract method")
		return lt.TestCODECOPY(opts, x)

	case 0xDEADBEEF003A:
		log.Trace().Str("method", "TestGASPRICE").Msg("Executing contract method")
		return lt.TestGASPRICE(opts, x)

	case 0xDEADBEEF003B:
		log.Trace().Str("method", "TestEXTCODESIZE").Msg("Executing contract method")
		return lt.TestEXTCODESIZE(opts, x)

	case 0xDEADBEEF003D:
		log.Trace().Str("method", "TestRETURNDATASIZE").Msg("Executing contract method")
		return lt.TestRETURNDATASIZE(opts, x)

	case 0xDEADBEEF003E:
		log.Trace().Str("method", "TestRETURNDATACOPY").Msg("Executing contract method")
		return lt.TestRETURNDATACOPY(opts, x)

	case 0xDEADBEEF0040:
		log.Trace().Str("method", "TestBLOCKHASH").Msg("Executing contract method")
		return lt.TestBLOCKHASH(opts, x)

	case 0xDEADBEEF0041:
		log.Trace().Str("method", "TestCOINBASE").Msg("Executing contract method")
		return lt.TestCOINBASE(opts, x)

	case 0xDEADBEEF0042:
		log.Trace().Str("method", "TestTIMESTAMP").Msg("Executing contract method")
		return lt.TestTIMESTAMP(opts, x)

	case 0xDEADBEEF0043:
		log.Trace().Str("method", "TestNUMBER").Msg("Executing contract method")
		return lt.TestNUMBER(opts, x)

	case 0xDEADBEEF0044:
		log.Trace().Str("method", "TestDIFFICULTY").Msg("Executing contract method")
		return lt.TestDIFFICULTY(opts, x)

	case 0xDEADBEEF0045:
		log.Trace().Str("method", "TestGASLIMIT").Msg("Executing contract method")
		return lt.TestGASLIMIT(opts, x)

	case 0xDEADBEEF0046:
		log.Trace().Str("method", "TestCHAINID").Msg("Executing contract method")
		return lt.TestCHAINID(opts, x)

	case 0xDEADBEEF0047:
		log.Trace().Str("method", "TestSELFBALANCE").Msg("Executing contract method")
		return lt.TestSELFBALANCE(opts, x)

	case 0xDEADBEEF0048:
		log.Trace().Str("method", "TestBASEFEE").Msg("Executing contract method")
		return lt.TestBASEFEE(opts, x)

	case 0xDEADBEEF0051:
		log.Trace().Str("method", "TestMLOAD").Msg("Executing contract method")
		return lt.TestMLOAD(opts, x)

	case 0xDEADBEEF0052:
		log.Trace().Str("method", "TestMSTORE").Msg("Executing contract method")
		return lt.TestMSTORE(opts, x)

	case 0xDEADBEEF0053:
		log.Trace().Str("method", "TestMSTORE8").Msg("Executing contract method")
		return lt.TestMSTORE8(opts, x)

	case 0xDEADBEEF0054:
		log.Trace().Str("method", "TestSLOAD").Msg("Executing contract method")
		return lt.TestSLOAD(opts, x)

	case 0xDEADBEEF0055:
		log.Trace().Str("method", "TestSSTORE").Msg("Executing contract method")
		return lt.TestSSTORE(opts, x)

	case 0xDEADBEEF0059:
		log.Trace().Str("method", "TestMSIZE").Msg("Executing contract method")
		return lt.TestMSIZE(opts, x)

	case 0xDEADBEEF005A:
		log.Trace().Str("method", "TestGAS").Msg("Executing contract method")
		return lt.TestGAS(opts, x)

	case 0xDEADBEEF00A0:
		log.Trace().Str("method", "TestLOG0").Msg("Executing contract method")
		return lt.TestLOG0(opts, x)

	case 0xDEADBEEF00A1:
		log.Trace().Str("method", "TestLOG1").Msg("Executing contract method")
		return lt.TestLOG1(opts, x)

	case 0xDEADBEEF00A2:
		log.Trace().Str("method", "TestLOG2").Msg("Executing contract method")
		return lt.TestLOG2(opts, x)

	case 0xDEADBEEF00A3:
		log.Trace().Str("method", "TestLOG3").Msg("Executing contract method")
		return lt.TestLOG3(opts, x)

	case 0xDEADBEEF00A4:
		log.Trace().Str("method", "TestLOG4").Msg("Executing contract method")
		return lt.TestLOG4(opts, x)

	}
	return nil, fmt.Errorf("the tx code %d was unrecognized", shortCode)
}

func GetRandomOPCode() uint64 {
	codes := []uint64{
		0x01,
		0x02,
		0x03,
		0x04,
		0x05,
		0x06,
		0x07,
		0x08,
		0x09,
		0x0A,
		0x0B,
		0x10,
		0x11,
		0x12,
		0x13,
		0x14,
		0x15,
		0x16,
		0x17,
		0x18,
		0x19,
		0x1A,
		0x1B,
		0x1C,
		0x1D,
		0x20,
		0x30,
		0x31,
		0x32,
		0x33,
		0x34,
		0x35,
		0x36,
		0x37,
		0x38,
		0x39,
		0x3A,
		0x3B,
		0x3D,
		// 0x3E, // return data copy is buggy in the test contract.
		0x40,
		0x41,
		0x42,
		0x43,
		0x44,
		0x45,
		0x46,
		0x47,
		0x48,
		0x51,
		0x52,
		0x53,
		0x54,
		0x55,
		0x59,
		0x5A,
		0xA0,
		0xA1,
		0xA2,
		0xA3,
		0xA4,
	}

	return codes[randSrc.Intn(len(codes))]
}

func init() {
	randSrc = rand.New(rand.NewSource(time.Now().Unix()))
}
