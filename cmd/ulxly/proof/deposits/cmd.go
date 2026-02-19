// Package deposits provides the proof deposits command.
package deposits

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"

	"github.com/0xPolygon/polygon-cli/bindings/ulxly"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	ArgDepositCount = "deposit-count"
	ArgFileName     = "file-name"
)

var (
	//go:embed usage.md
	usage string

	proofOptions *proofArgs
	fileOptions  *fileArgs
)

type proofArgs struct {
	DepositCount uint32
}

type fileArgs struct {
	FileName string
}

var Cmd = &cobra.Command{
	Use:          "proof",
	Short:        "Generate a proof for a given range of deposits.",
	Long:         usage,
	RunE:         runProof,
	SilenceUsage: true,
}

func init() {
	proofOptions = &proofArgs{}
	fileOptions = &fileArgs{}

	f := Cmd.Flags()
	f.StringVar(&fileOptions.FileName, ArgFileName, "", "ndjson file with events data")
	f.Uint32Var(&proofOptions.DepositCount, ArgDepositCount, 0, "deposit number to generate a proof for")
}

func runProof(_ *cobra.Command, args []string) error {
	return proof(args)
}

func proof(args []string) error {
	depositNumber := proofOptions.DepositCount
	rawDepositData, err := getInputData(args)
	if err != nil {
		return err
	}
	return readDeposits(rawDepositData, uint32(depositNumber))
}

func getInputData(args []string) ([]byte, error) {
	fileName := fileOptions.FileName
	if fileName != "" {
		return os.ReadFile(fileName)
	}

	if len(args) > 1 {
		concat := strings.Join(args[1:], " ")
		return []byte(concat), nil
	}

	return io.ReadAll(os.Stdin)
}

func readDeposits(rawDeposits []byte, depositNumber uint32) error {
	buf := bytes.NewBuffer(rawDeposits)
	scanner := bufio.NewScanner(buf)
	scannerBuf := make([]byte, 0)
	scanner.Buffer(scannerBuf, 1024*1024)
	imt := new(ulxlycommon.IMT)
	imt.Init()
	seenDeposit := make(map[uint32]common.Hash, 0)
	lastDeposit := uint32(0)
	for scanner.Scan() {
		evt := new(ulxly.UlxlyBridgeEvent)
		err := json.Unmarshal(scanner.Bytes(), evt)
		if err != nil {
			return err
		}
		if _, hasBeenSeen := seenDeposit[evt.DepositCount]; hasBeenSeen {
			log.Warn().Uint32("deposit", evt.DepositCount).Str("tx-hash", evt.Raw.TxHash.String()).Msg("Skipping duplicate deposit")
			continue
		}
		seenDeposit[evt.DepositCount] = evt.Raw.TxHash
		if lastDeposit+1 != evt.DepositCount && lastDeposit != 0 {
			log.Error().Uint32("missing-deposit", lastDeposit+1).Uint32("current-deposit", evt.DepositCount).Msg("Missing deposit")
			return fmt.Errorf("missing deposit: %d", lastDeposit+1)
		}
		lastDeposit = evt.DepositCount
		leaf := ulxlycommon.HashDeposit(evt)
		log.Debug().Str("leaf-hash", common.Bytes2Hex(leaf[:])).Msg("Leaf hash calculated")
		imt.AddLeaf(leaf, evt.DepositCount)
		log.Info().
			Uint64("block-number", evt.Raw.BlockNumber).
			Uint32("deposit-count", evt.DepositCount).
			Str("tx-hash", evt.Raw.TxHash.String()).
			Str("root", common.Hash(imt.Roots[len(imt.Roots)-1]).String()).
			Msg("adding event to tree")
		// There's no point adding more leaves if we can prove the deposit already?
		if evt.DepositCount >= depositNumber {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("there was an error reading the deposit file")
		return err
	}

	log.Info().Msg("finished")
	p := imt.GetProof(depositNumber)
	fmt.Println(proofString(p))
	return nil
}

// proofString will create the json representation of the proof
func proofString[T any](p T) string {
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		log.Error().Err(err).Msg("error marshalling proof to json")
		return ""
	}
	return string(jsonBytes)
}
