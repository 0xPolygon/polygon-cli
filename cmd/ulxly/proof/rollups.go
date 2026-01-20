package proof

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

	"github.com/0xPolygon/polygon-cli/bindings/ulxly/polygonrollupmanager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	ArgRollupID   = "rollup-id"
	ArgCompleteMT = "complete-merkle-tree"
)

var (
	//go:embed rollupsProofUsage.md
	rollupsProofUsage string

	rollupsProofOptions *rollupsProofArgs
	rollupsFileOptions  *fileArgs
)

type rollupsProofArgs struct {
	RollupID           uint32
	CompleteMerkleTree bool
}

var RollupsProofCmd = &cobra.Command{
	Use:          "rollups-proof",
	Short:        "Generate a proof for a given range of rollups.",
	Long:         rollupsProofUsage,
	RunE:         runRollupsProof,
	SilenceUsage: true,
}

func init() {
	rollupsProofOptions = &rollupsProofArgs{}
	rollupsFileOptions = &fileArgs{}

	f := RollupsProofCmd.Flags()
	f.StringVar(&rollupsFileOptions.FileName, ArgFileName, "", "ndjson file with events data")
	f.Uint32Var(&rollupsProofOptions.RollupID, ArgRollupID, 0, "rollup ID number to generate a proof for")
	f.BoolVar(&rollupsProofOptions.CompleteMerkleTree, ArgCompleteMT, false, "get proof for a leave higher than the highest rollup ID")
}

func runRollupsProof(_ *cobra.Command, args []string) error {
	return rollupsExitRootProof(args)
}

func rollupsExitRootProof(args []string) error {
	rollupID := rollupsProofOptions.RollupID
	completeMT := rollupsProofOptions.CompleteMerkleTree
	rawLeavesData, err := getRollupsInputData(args)
	if err != nil {
		return err
	}
	return readRollupsExitRootLeaves(rawLeavesData, rollupID, completeMT)
}

func getRollupsInputData(args []string) ([]byte, error) {
	fileName := rollupsFileOptions.FileName
	if fileName != "" {
		return os.ReadFile(fileName)
	}

	if len(args) > 1 {
		concat := strings.Join(args[1:], " ")
		return []byte(concat), nil
	}

	return io.ReadAll(os.Stdin)
}

func readRollupsExitRootLeaves(rawLeaves []byte, rollupID uint32, completeMT bool) error {
	buf := bytes.NewBuffer(rawLeaves)
	scanner := bufio.NewScanner(buf)
	scannerBuf := make([]byte, 0)
	scanner.Buffer(scannerBuf, 1024*1024)
	leaves := make(map[uint32]*polygonrollupmanager.PolygonrollupmanagerVerifyBatchesTrustedAggregator, 0)
	highestRollupID := uint32(0)
	for scanner.Scan() {
		evt := new(polygonrollupmanager.PolygonrollupmanagerVerifyBatchesTrustedAggregator)
		err := json.Unmarshal(scanner.Bytes(), evt)
		if err != nil {
			return err
		}
		if highestRollupID < evt.RollupID {
			highestRollupID = evt.RollupID
		}
		leaves[evt.RollupID] = evt
	}
	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("there was an error reading the deposit file")
		return err
	}
	if rollupID > highestRollupID && !completeMT {
		return fmt.Errorf("rollupID %d required is higher than the highest rollupID %d provided in the file. Please use --complete-merkle-tree option if you know what you are doing", rollupID, highestRollupID)
	} else if completeMT {
		highestRollupID = rollupID
	}
	var ls []common.Hash
	var i uint32 = 0
	for ; i <= highestRollupID; i++ {
		var exitRoot common.Hash
		if leaf, exists := leaves[i]; exists {
			exitRoot = leaf.ExitRoot
			log.Info().
				Uint64("block-number", leaf.Raw.BlockNumber).
				Uint32("rollupID", leaf.RollupID).
				Str("exitRoot", exitRoot.String()).
				Str("tx-hash", leaf.Raw.TxHash.String()).
				Msg("latest event received for the tree")
		} else {
			log.Warn().Uint32("rollupID", i).Msg("No event found for this rollup")
		}
		ls = append(ls, exitRoot)
	}
	p, err := ulxlycommon.ComputeSiblings(rollupID, ls, ulxlycommon.TreeDepth)
	if err != nil {
		return err
	}
	log.Info().Str("root", p.Root.String()).Msg("finished")
	fmt.Println(proofString(p))
	return nil
}
