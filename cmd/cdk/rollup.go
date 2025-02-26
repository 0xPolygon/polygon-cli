package cdk

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rollupCmd = &cobra.Command{
	Use:  "rollup",
	Args: cobra.NoArgs,
}

var rollupInspectCmd = &cobra.Command{
	Use:  "inspect",
	Args: cobra.NoArgs,
}

var rollupDumpCmd = &cobra.Command{
	Use:  "dump",
	Args: cobra.NoArgs,
}

var rollupMonitorCmd = &cobra.Command{
	Use:  "monitor",
	Args: cobra.NoArgs,
}

func inspectRollupInfo(rollupManager rollupManagerContractInterface, rollupID uint32) error {
	rollup_data, err := rollupManager.RollupIDToRollupData(nil, rollupID)
	if err != nil {
		return err
	}
	b, err := json.Marshal(rollup_data)
	if err != nil {
		return err
	}
	log.Info().Msg("rollup_data: " + string(b))

	return nil
}
