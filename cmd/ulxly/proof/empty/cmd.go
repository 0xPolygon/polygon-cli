// Package empty provides the proof empty command.
package empty

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
	Use:          "empty-proof",
	Short:        "Create an empty proof.",
	Long:         usage,
	RunE:         runEmptyProof,
	SilenceUsage: true,
}

func runEmptyProof(_ *cobra.Command, _ []string) error {
	return emptyProof()
}

func emptyProof() error {
	p := new(ulxlycommon.Proof)

	e := ulxlycommon.GenerateEmptyHashes(ulxlycommon.TreeDepth)
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
