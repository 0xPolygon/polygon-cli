package nodelist

import (
	"encoding/json"
	"os"

	"github.com/0xPolygon/polygon-cli/p2p/database"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type (
	nodeListParams struct {
		ProjectID  string
		DatabaseID string
		OutputFile string
		Limit      int
	}
)

var (
	inputNodeListParams nodeListParams
)

var NodeListCmd = &cobra.Command{
	Use:   "nodelist [nodes.json]",
	Short: "Generate a node list to seed a node.",
	Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		inputNodeListParams.OutputFile = args[0]
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		db := database.NewDatastore(cmd.Context(), database.DatastoreOptions{
			ProjectID:  inputNodeListParams.ProjectID,
			DatabaseID: inputNodeListParams.DatabaseID,
		})

		nodes, err := db.NodeList(ctx, inputNodeListParams.Limit)
		if err != nil {
			return err
		}

		bytes, err := json.MarshalIndent(nodes, "", "    ")
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
	f := NodeListCmd.Flags()
	f.IntVarP(&inputNodeListParams.Limit, "limit", "l", 100, "number of unique nodes to return")
	f.StringVarP(&inputNodeListParams.ProjectID, "project-id", "p", "", "GCP project ID")
	f.StringVarP(&inputNodeListParams.DatabaseID, "database-id", "d", "", "datastore database ID")
	if err := NodeListCmd.MarkFlagRequired("project-id"); err != nil {
		log.Error().Err(err).Msg("Failed to mark project-id as required flag")
	}
}
