// Package zero provides the proof zero command.
package zero

import (
	_ "embed"
	"encoding/json"
	"fmt"

	ulxlycommon "github.com/0xPolygon/polygon-cli/cmd/ulxly/common"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//go:embed usage.md
var usage string

var Cmd = &cobra.Command{
	Use:          "zero-proof",
	Short:        "Create a proof that's filled with zeros.",
	Long:         usage,
	RunE:         runZeroProof,
	SilenceUsage: true,
}

func runZeroProof(_ *cobra.Command, _ []string) error {
	return zeroProof()
}

func zeroProof() error {
	p := new(ulxlycommon.Proof)

	e := ulxlycommon.GenerateZeroHashes(ulxlycommon.TreeDepth)
	copy(p.Siblings[:], e)
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
