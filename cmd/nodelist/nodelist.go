package nodelist

import (
	"encoding/json"
	"os"

	"github.com/maticnetwork/polygon-cli/p2p/database"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const jsonIndent = "    "

type (
	nodeListParams struct {
		ProjectID  string
		OutputFile string
		Limit      int
	}
)

var (
	inputNodeListParams nodeListParams
)

var NodeListCmd = &cobra.Command{
	Use:   "nodelist [nodes.json]",
	Short: "Query the database for ",
	Long: `Ping nodes by either giving a single enode/enr or an entire nodes file.

This command will establish a handshake and status exchange to get the Hello and
Status messages and output JSON. If providing a enode/enr rather than a nodes
file, then the connection will remain open by default (--listen=true), and you
can see other messages the peer sends (e.g. blocks, transactions, etc.).`,
	Args: cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		inputNodeListParams.OutputFile = args[0]
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		db := database.NewDatastore(cmd.Context(), database.DatastoreOptions{
			ProjectID: inputNodeListParams.ProjectID,
		})

		nodes, err := db.NodeList(ctx, inputNodeListParams.Limit)
		if err != nil {
			return err
		}

		bytes, err := json.MarshalIndent(nodes, "", jsonIndent)
		if err != nil {
			return err
		}

		if err = os.WriteFile(inputNodeListParams.OutputFile, bytes, 0644); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	NodeListCmd.PersistentFlags().IntVarP(&inputNodeListParams.Limit, "limit", "l", 100, "Number of unique nodes to return")
	NodeListCmd.PersistentFlags().StringVarP(&inputNodeListParams.ProjectID, "project-id", "p", "", "GCP project ID")
	if err := NodeListCmd.MarkPersistentFlagRequired("project-id"); err != nil {
		log.Error().Err(err).Msg("Failed to mark project-id as required persistent flag")
	}
}
