package enr

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/p2p/enode"
)

var (
	//go:embed usage.md
	usage         string
	inputFileName *string
)

var ENRCmd = &cobra.Command{
	Use:   "enr [flags]",
	Short: "Convert between ENR and Enode format",
	Long:  usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		rawData, err := getInputData(cmd, args)
		if err != nil {
			log.Error().Err(err).Msg("Unable to read input")
			return err
		}
		lines := strings.Split(string(rawData), "\n")

		for _, l := range lines {
			var node *enode.Node
			var err error
			l = strings.TrimSpace(l)
			if l == "" {
				continue
			}
			isENR := false
			genericNode := make(map[string]any, 0)

			if strings.HasPrefix(l, "enr:") {
				isENR = true
				node, err = enode.Parse(enode.V4ID{}, l)
				if err != nil {
					log.Error().Err(err).Str("line", l).Msg("Unable to parse enr record")
					continue
				}
				var eth struct {
					ForkID forkid.ID
					Tail   []rlp.RawValue `rlp:"tail"`
				}
				_ = node.Load(enr.WithEntry("eth", &eth))
				genericNode["forkid"] = ForkIDToString(eth.ForkID)
			} else {
				node, err = enode.ParseV4(l)
				if err != nil {
					log.Error().Err(err).Str("line", l).Msg("Unable to parse node record")
					continue
				}
			}

			if isENR {
				genericNode["enr"] = node.String()
			}
			genericNode["enode"] = node.URLv4()
			genericNode["id"] = node.ID().String()
			genericNode["ip"] = node.IP().String()
			genericNode["seq"] = node.Seq()
			genericNode["tcp"] = fmt.Sprintf("%d", node.TCP())
			genericNode["udp"] = fmt.Sprintf("%d", node.UDP())
			jsonOut, err := json.Marshal(genericNode)
			if err != nil {
				log.Error().Err(err).Msg("Unable to convert node to json")
				continue
			}
			fmt.Println(string(jsonOut))
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func ForkIDToString(id forkid.ID) string {
	// Convert the Hash field to a hex string
	hashHex := hex.EncodeToString(id.Hash[:])

	// Convert the Next field to a hex string
	nextHex := fmt.Sprintf("%016x", id.Next)

	// Concatenate the two hex strings
	return hashHex + nextHex
}

func init() {
	flagSet := ENRCmd.PersistentFlags()
	inputFileName = flagSet.String("file", "", "Provide a file that's holding ENRs")
}
func getInputData(cmd *cobra.Command, args []string) ([]byte, error) {
	if inputFileName != nil && *inputFileName != "" {
		return os.ReadFile(*inputFileName)
	}

	if len(args) >= 1 {
		concat := strings.Join(args, "\n")
		return []byte(concat), nil
	}

	return io.ReadAll(os.Stdin)
}
