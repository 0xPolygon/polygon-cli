package nodelist

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/0xPolygon/polygon-cli/flag"
	"github.com/0xPolygon/polygon-cli/p2p/database"
	"github.com/spf13/cobra"
)

type (
	nodeListParams struct {
		ProjectID     string
		DatabaseID    string
		ClickHouseDSN string
		OutputFile    string
		Limit         int
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
		if inputNodeListParams.Limit <= 0 {
			return fmt.Errorf("--limit must be greater than zero, got %d", inputNodeListParams.Limit)
		}
		inputNodeListParams.OutputFile = args[0]
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		nodes, err := database.ClickHouseNodeList(ctx, inputNodeListParams.ClickHouseDSN, inputNodeListParams.Limit)
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
	f.StringVar(&inputNodeListParams.ClickHouseDSN, "clickhouse-dsn", "",
		"ClickHouse DSN of the sensor database to read peers from, e.g. clickhouse://user:pass@host:9000/mainnet")
	flag.MarkFlagsRequired(NodeListCmd, "clickhouse-dsn")
	// The Datastore backend is retired. Its flags are still accepted so existing
	// invocations keep parsing, but they do nothing.
	f.StringVarP(&inputNodeListParams.ProjectID, "project-id", "p", "", "GCP project ID")
	f.StringVarP(&inputNodeListParams.DatabaseID, "database-id", "d", "", "datastore database ID")
	flag.MarkDeprecated(NodeListCmd, "project-id", "the Datastore backend is retired; use --clickhouse-dsn")
	flag.MarkDeprecated(NodeListCmd, "database-id", "the Datastore backend is retired; use --clickhouse-dsn")
}
